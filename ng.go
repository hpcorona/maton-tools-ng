package main

import (
  "fmt"
  "flag"
  "github.com/hpcorona/go-v8/v8"
  "os"
  "time"
  "io/ioutil"
  "path/filepath"
  "strings"
)

var context *v8.V8Context
var startTime time.Time
var ngdirs []string

func localTime() time.Time {
  return time.Now()
}

func finish() {
  endTime := localTime()

  diffs := endTime.Sub(startTime)

  fmt.Printf("Time spent: %s\n", diffs.String())
  fmt.Printf("Script finished\n")
}

func absfile(file string) string {
  cleaned := filepath.Clean(file)

  if strings.HasPrefix(cleaned, "/") {
    return cleaned
  }
  wd, _ := os.Getwd()
  return filepath.Join(wd, cleaned)
}

func main() {
  startTime = localTime()
  ngdirs = make([]string, 1, 20)

  var ngfile = "ng.js"
  var work, _ = os.Getwd()
  flag.StringVar(&ngfile, "file", "ng.js", "override the default 'ng.js' file in the current directory")
  flag.StringVar(&work, "work", work, "change the initial working directory")
  flag.Parse()

  ngdata, err := ioutil.ReadFile(ngfile)
  if err != nil {
    fmt.Printf("File %s not found\n", ngfile)
    return
  }

  ngdirs[0], ngfile = filepath.Split(absfile(ngfile))

  fmt.Printf("Nailgun\n")
  fmt.Printf("ng file:    %s\n", ngfile)
  fmt.Printf("ng dir:     %s\n", ngdirs[0])
  fmt.Printf("work dir:   %s\n", work)

  os.Chdir(work)

  context = v8.NewContext()
  loadFunctions(context)

  _, err = context.Eval(string(ngdata))
  if err != nil {
    fmt.Printf("=====\nERROR\n=====\n%s:%s", ngfile, err.Error())
    return
  }

  task := defTask
  if flag.NArg() > 0 {
    task = flag.Arg(0)
  }
  if task == "" || task == "?" {
    if task == "" {
      fmt.Printf("No default task defined.\n")
    }
    showHelp()
  } else {
    tst := runTask(task)
    if !tst {
      fmt.Printf("Error while trying to run the task %s\n", task)
    }
  }

  defer finish()
}

