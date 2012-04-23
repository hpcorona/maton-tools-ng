package main

import (
	"fmt"
	"github.com/hpcorona/go-v8"
	"os"
	"os/exec"
	"io"
)

func load_os_functions(ctx *v8.V8Context) {
	// operating system functions
	ctx.AddFunc("_os_log", os_log)
	ctx.AddFunc("_os_env", os_env)
	ctx.AddFunc("_os_exit", os_exit)
	ctx.AddFunc("_os_run", os_run)
	ctx.AddFunc("_os_findCmd", os_findCmd)
	ctx.AddFunc("_os_error", os_error)

	ctx.Eval(`
  this.os = 
    {
      "log": function() { _os_log.apply(this, arguments); },
      "env": function(args) { return _os_env(args); },
      "exit": function(args) { _os_exit(args); },
      "run": function() { return _os_run.apply(this, arguments); },
      "findCmd": function(args) { return _os_findCmd(args); },
			"error": function(args) { _os_error(args); }
    };
  `)
}

func os_log(value ...interface{}) interface{} {
	for i := 0; i < len(value); i++ {
		fmt.Printf("%v", value[i])
	}

	fmt.Printf("\n")

	return true
}

func os_findCmd(value ...interface{}) interface{} {
	paramCount(value, 1, "os.findCmd")

	path, err := exec.LookPath(value[0].(string))
	if err != nil {
		panic(err.Error())
	}

	return path
}

func os_run(value ...interface{}) interface{} {
	paramMin(value, 1, "os.run")

  fmt.Printf("Run '%s'", value[0].(string))

  strs := make([]string, len(value) - 1)
  for i := 0; i < len(value) - 1; i++ {
    strs[i] = value[i+1].(string)
    fmt.Printf(" %s", strs[i])
  }
  fmt.Println()

	cmd := exec.Command(value[0].(string), strs...)
	stdout, err := cmd.StdoutPipe()
  if err != nil {
		panic(err.Error())
  }
  stderr, err := cmd.StderrPipe()
  if err != nil {
		panic(err.Error())
  }
  err = cmd.Start()
  if err != nil {
		panic(err.Error())
  }
	go io.Copy(os.Stdout, stdout) 
	go io.Copy(os.Stderr, stderr) 
  cmd.Wait()

  return true
}

func os_env(value ...interface{}) interface{} {
	paramCount(value, 1, "os.env")

	return os.Getenv(value[0].(string))
}

func os_exit(value ...interface{}) interface{} {
	paramCount(value, 1, "os.exit")

  fmt.Printf("System exit '%d'.\n", value[0].(int))

	os.Exit(value[0].(int))

	return true
}

func os_error(value ...interface{}) interface{} {
	paramCount(value, 1, "os.error")
	
	panic(value[0].(string))
	
	return true
}
