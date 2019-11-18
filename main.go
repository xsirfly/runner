package main

import (
	"flag"
	"runner/context"
	"bytes"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"runner/model"
	"encoding/json"
	"os"
	"errors"
)

func main() {
	input := flag.String("input", "", "test case input")
	expected := flag.String("expected", "", "test case expected")
	key := flag.String("key", "", "program language")
	timeout := flag.Int("timeout", 2000, "timeout in milliseconds")
	flag.Parse()

	var (
		ctx *context.Context
		ok bool
	)

	r := new(model.Result)
	if ctx, ok = context.GetRunContextByProjectKey(*key); !ok {
		logResult(r.GetSystemErrorTaskResult(errors.New("not support project")))
		return
	}

	// compile
	if ctx.CompileCommand != "" {
		var stdout, stderr bytes.Buffer
		cmd := exec.Command(ctx.CompileCommand, ctx.CompileArgs...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			if stderr.Len() > 0 {
				logResult(r.GetCompileErrorTaskResult(stderr.String()))
			} else {
				logResult(r.GetSystemErrorTaskResult(err))
			}
			return
		}
	}

	// run and judge
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(ctx.RunCommand, ctx.RunArgs...)
	cmd.Stdin = strings.NewReader(*input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	time.AfterFunc(time.Duration(*timeout)*time.Millisecond, func() {
		logResult(r.GetTimeLimitExceededErrorTaskResult())
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})

	startTime := time.Now().UnixNano() / 1e6
	if err := cmd.Run(); err != nil {
		if r.Status > 0 {
			// time limit exceeded, receive kill signal, do nothing
		} else if stderr.Len() > 0 {
			// be be judged program panic
			logResult(r.GetRuntimeErrorTaskResult(stderr.String()))
		} else {
			// system error
			logResult(r.GetSystemErrorTaskResult(err))
		}
		return
	}
	endTime := time.Now().UnixNano() / 1e6

	if stderr.Len() > 0 {
		logResult(r.GetRuntimeErrorTaskResult(stderr.String()))
		return
	}

	output := strings.TrimSpace(stdout.String())
	if output == *expected {
		// ms, MB
		timeCost, memoryCost := endTime-startTime, cmd.ProcessState.SysUsage().(*syscall.Rusage).Maxrss/1024
		// timeCost value 0 will be omitted
		if timeCost == 0 {
			timeCost = 1
		}

		logResult(r.GetAcceptedTaskResult(timeCost, memoryCost, *input, output, *expected))
		return
	} else {
		logResult(r.GetWrongAnswerTaskResult(*input, output, *expected))
		return
	}
}

func logResult(r *model.Result) {
	result, _ := json.Marshal(r)
	_, _ = os.Stdout.Write(result)
}
