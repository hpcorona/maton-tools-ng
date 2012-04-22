package main

import (
  "strings"
	"github.com/hpcorona/go-v8"
	"io/ioutil"
	"path/filepath"
)

func load_ng_functions(ctx *v8.V8Context) {
	// ng functions
  ctx.AddFunc("_ng_include", ng_include)
  ctx.AddFunc("_ng_task", ng_task)
  ctx.AddFunc("_ng_default", ng_default)
	ctx.AddFunc("_ng_wd", ng_wd)

  ctx.Eval(`
    this.ng = {
      "include": function() { return _ng_include.apply(this, arguments); },
      "task": function() { return _ng_task.apply(this, arguments); },
      "default": function(tn) { return _ng_default(tn); },
			"wd": function() { return _ng_wd(); }
    };
  `)
}

func push_ngdir(file string) {
  abs := absfile(file)
  dir, _ := filepath.Split(abs)

  ngdirs = append(ngdirs, dir)
}

func pop_ngdir() {
  ngdirs = ngdirs[: len(ngdirs) - 1]
}

func get_ngdir() string {
  return ngdirs[len(ngdirs) - 1]
}

func get_absjs(file string) string {
  file = filepath.Clean(file)
  if strings.HasPrefix(file, "/") {
    return file
  }
  return filepath.Join(get_ngdir(), file)
}

func load_js(file string) {
  file = get_absjs(file)
  push_ngdir(file)
  ngdata, err := ioutil.ReadFile(file)
  if err != nil {
    panic(err.Error())
  }
  context.Eval(string(ngdata))
  pop_ngdir()
}

func ng_include(value ...interface{}) interface{} {
  paramMin(value, 1, "ng.include")

  for i := 0; i < len(value); i++ {
    load_js(value[i].(string))
  }

  return true
}

func ng_task(value ...interface{}) interface{} {
  paramMin(value, 3, "ng.task")
  paramMax(value, 4, "ng.task")

  funcName := value[0].(string)
  taskName := value[1].(string)
  desc := value[2].(string)

  var deps []string
  if len(value) < 4 {
    deps = make([]string, 0)
  } else {
    intd, ok := value[3].([]interface{})
    if ok {
      deps = make([]string, len(intd))
      for i := 0; i < len(intd); i++ {
        deps[i] = intd[i].(string)
      }
    } else {
      deps = make([]string, 1)
      deps[0] = value[3].(string)
    }
  }
  return newTask(taskName, desc, deps, funcName)
}

func ng_default(value ...interface{}) interface{} {
  paramCount(value, 1, "ng.default")

  return defaultTask(value[0].(string))
}

func ng_wd(value ...interface{}) interface{} {
	paramCount(value, 0, "ng.wd")
	
	return get_ngdir()
}
