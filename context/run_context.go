package context

type Context struct {
	CompileCommand string
	CompileArgs    []string
	RunCommand     string
	RunArgs        []string
}

var contextMap map[string]*Context

func init() {
	contextMap = make(map[string]*Context)
	contextMap["java"] = &Context{
		CompileCommand: "javac",
		CompileArgs: []string{
			"Main.java",
		},
		RunCommand: "java",
		RunArgs: []string{
			"Main",
		},
	}

	contextMap["c"] = &Context{
		CompileCommand: "gcc",
		CompileArgs: []string{
			"Main.c",
			"-save-temps",
			"-std=gnu11",
			"-o",
			"Main",
		},
		RunCommand: "./Main",
	}

	contextMap["c++"] = &Context{
		CompileCommand: "gcc",
		CompileArgs: []string{
			"Main.cpp",
			"-save-temps",
			"-std=gnull",
			"-fmax-errors=10",
			"-static",
			"-o",
			"Main",
		},
		RunCommand: "./Main",
	}
}

func GetRunContextByProjectKey(projectKey string) (*Context, bool) {
	c, ok := contextMap[projectKey]
	return c, ok
}


