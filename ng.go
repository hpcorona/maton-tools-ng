package main

import (
  "fmt"
  "flag"
  "github.com/hpcorona/go-v8/v8"
  "os"
  "time"
  "io/ioutil"
)

var context *v8.V8Context
var startTime time.Time

func localTime() time.Time {
  return time.Now()
}

func finish() {
  endTime := localTime()

  diffs := endTime.Sub(startTime)

  fmt.Printf("Time spent: %s\n", diffs.String())
  fmt.Printf("Script finished\n")
}

func main() {
  startTime = localTime()

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

  fmt.Printf("Match Test: %v\n",
    match("output/**/source.java", "output/file/one/two/three/source.java", true))
  fmt.Printf("Match Test: %v\n",
    match("**/*.java", "my/pkg/File.java", true))
  fmt.Printf("Match Test: %v\n",
    match("*.java", "my/pkg/File.java", true))

  fmt.Printf("Nailgun\n")
  fmt.Printf("ng file:    %s\n", ngfile)
  fmt.Printf("work dir:   %s\n", work)

  os.Chdir(work)

  context = v8.NewContext()
  loadFunctions(context)

  context.Eval(string(ngdata))

  defer finish()
}

