package main

import (
  "encoding/json"
  "text/template"
  "bytes"
  "fmt"
)

type Task struct {
  Name string
  Description string
  Depends []string
  Function string
  Executed bool
}

var tasks = make(map[string]*Task)
var queuedTasks = make([]string, 0, 10)
var defTask string = ""

var runTaskTpl = template.Must(template.New("runTask").Parse(`
{{.Name}}({{.Params}});
`))

func newTask(name, description string, depends []string, function string) bool {
  task := &Task {
    Name: name,
    Description: description,
    Depends: depends,
    Function: function,
    Executed: false,
    }
  tasks[name] = task

  return true
}

func defaultTask(name string) bool {
  _, found := tasks[name]
  if !found {
    return false
  }

  defTask = name
  return true
}

func queueTask(name string) bool {
  task := tasks[name]
  if task == nil {
    fmt.Printf("Task %s not found\n", name)
    return false
  }
  if task.Executed {
    return true
  }
  task.Executed = true

  for i := 0; i < len(task.Depends); i++ {
    if !queueTask(task.Depends[i]) {
      return false
    }
  }
  queuedTasks = append(queuedTasks, name)

  return true
}

func runTask(name string, params []string) bool {
  if !queueTask(name) {
    return false
  }

  v, _ := json.Marshal(params)

  for i := 0; i < len(queuedTasks); i++ {
    task := tasks[queuedTasks[i]]
    fmt.Printf("TASK '%s'\n", task.Name)

    b := bytes.NewBufferString("")
    runTaskTpl.Execute(b, map[string]interface{} {
        "Name": task.Function,
        "Params": string(v),
      })
    _, err := context.Eval(b.String())
    if err != nil {
      fmt.Printf("%v\n", err)
      return false
    }
  }

  return true
}

func showHelp() {
  if len(tasks) == 0 {
    fmt.Printf("No tasks defined.\n")
    return
  }

  fmt.Printf("Tasks available:\n")
  for _, task := range tasks {
    fmt.Printf("\t%s\n\t\t%s\n", task.Name, task.Description)
  }

  fmt.Printf("The default task is: %s\n", defTask);
}

