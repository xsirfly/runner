package model

type Result struct {
	Runtime      int64  `json:"runtime,omitempty"`
	Memory       int64  `json:"memory,omitempty"`
	Status       int32  `json:"status"`
	Input        string `json:"input,omitempty"`
	Output       string `json:"output,omitempty"`
	Expected     string `json:"expected,omitempty"`
	SystemError  string `json:"system_error,omitempty"`
	CompileError string `json:"compile_error,omitempty"`
	RunError     string `json:"run_error,omitempty"`
}

const (
	StatusAc = iota
	_
	StatusRe
	StatusTle
	_
	StatusWa

	_
	StatusCe
)

func (r *Result) GetAcceptedTaskResult(runtime, memory int64, input, output, expected string) *Result {
	r.Status = StatusAc
	r.Runtime = runtime
	r.Memory = memory
	r.Input = input
	r.Output = output
	r.Expected = expected
	return r
}

func (r *Result) GetRuntimeErrorTaskResult(err string) *Result {
	r.Status = StatusRe
	r.RunError = err
	return r
}

func (r *Result) GetTimeLimitExceededErrorTaskResult() *Result {
	r.Status = StatusTle
	return r
}

func (r *Result) GetWrongAnswerTaskResult(input, output, expected string) *Result {
	r.Status = StatusWa
	r.Input = input
	r.Output = output
	r.Expected = expected
	return r
}

func (r *Result) GetSystemErrorTaskResult(err error) *Result {
	r.SystemError = err.Error()
	return r
}

func (r *Result) GetCompileErrorTaskResult(err string) *Result {
	r.Status = StatusCe
	r.CompileError = err
	return r
}
