package main

import (
  "fmt"
  "os"
  "time"
  "strconv"
  "encoding/json"
)

type Task struct {
  ID            int64   `json: id`
  Description   string  `json: description`
  Status        string  `json: status`
  Created_At    string  `json: created_At`
  Updated_At    string  `json: updated_At`
}

func loadTasks() ([]Task, error) {
  var tasks []Task
  file, err := os.Open("tasks.json")
  if err != nil {
    if os.IsNotExist(err) {
      return tasks, nil
    }
    return nil, err
  }
  defer file.Close()

  decoder := json.NewDecoder(file)
  err = decoder.Decode(&tasks)
  if err != nil {
    return nil, err
  }

  for i := range tasks {
    tasks[i].ID = int64(i + 1)
  }
  return tasks, nil
}

func saveTasks(tasks []Task) error {
  file, err := os.Create("tasks.json")
  if err != nil {
    return err
  }
  defer file.Close()

  encoder := json.NewEncoder(file)
  encoder.SetIndent("", " ")
  return encoder.Encode(tasks)
}


func addTask(tasks []Task) []Task {
  taskName := os.Args[2]
  currentTime := time.Now().Format("2006-01-02 15:04:05")
  newTask := Task{
    ID: int64(len(tasks) + 1), Description: taskName, Status: "todo", Created_At: currentTime, Updated_At: currentTime,
  }
  tasks = append(tasks, newTask)
  return tasks
}

func listTask(tasks []Task) {
  if len(tasks) == 0 {
    fmt.Printf("--No task added in QUEUE--\n")
    return
  }

  if len(os.Args) > 2 {
    filteredListTask(tasks)
    return
  }

  for _, task := range tasks {
    fmt.Printf("%d) %s: [%s] {%s - %s}\n", task.ID, task.Description, task.Status, task.Created_At, task.Updated_At)
  }
}

func filteredListTask(tasks []Task) {
  status := os.Args[2]
  flag := false
  for _, task := range tasks {
    if task.Status == status {
      fmt.Printf("%d) %s: [%s] {%s - %s}\n", task.ID, task.Description, task.Status, task.Created_At, task.Updated_At)
      flag = true
    }
  }
  if !flag {
    fmt.Printf("No task with %s founded.\n", status)
  }
}

func updateTask(tasks []Task) {
  id, err := strconv.ParseInt(os.Args[2], 10, 8)
  if err != nil {
    fmt.Println("Error while parsing ID: ", err)
    return
  }

  taskName := os.Args[3]
  flag := false
  for i := range tasks {
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    if tasks[i].ID == id {
      tasks[i].Description = taskName
      tasks[i].Updated_At = currentTime
      flag = true
      break
    }
  }
  
  if !flag {
    fmt.Printf("No task found by Id{%d}\n", id)
    }
}

func deleteTask(tasks []Task) []Task {
  
  if len(tasks) == 0 {
    fmt.Println("No tasks exists in the bucket")
    return tasks
  }

  id, err := strconv.ParseInt(os.Args[2], 10, 8)
  if err != nil {
    fmt.Println("Error while parsing ID: ", err)
    return tasks
  }

  if id < 1 || id > int64(len(tasks)) {
    fmt.Println("No tasks founded by that id")
    return tasks
  }
  tasks = append(tasks[:id-1], tasks[id:]...)
  return tasks
}

func markTodo(tasks []Task) {

  if len(tasks) == 0 {
    fmt.Println("No tasks exists in the bucket")
    return
  }

  id, err := strconv.ParseInt(os.Args[2], 10, 8)
  if err != nil {
    fmt.Println("Error while parsing ID: ", err)
  }
  tasks[id-1].Status = "todo"
}

func markInProgress(tasks []Task) {
  
  if len(tasks) == 0 {
    fmt.Println("No tasks exists in the bucket")
    return
  }

  id, err := strconv.ParseInt(os.Args[2], 10, 8)
  if err != nil {
    fmt.Println("Error while parsing ID: ", err)
    return
  }
  tasks[id-1].Status = "in-progress"
}

func markDone(tasks []Task) {
  if len(tasks) == 0 {
    fmt.Println("No tasks exists in the bucket")
    return
  }
  
  id, err := strconv.ParseInt(os.Args[2], 10, 8)
  if err != nil {
    fmt.Println("Error while parsing ID: ", err)
    return
  }
  tasks[id-1].Status = "done"
}

func main() {
  // load all the tasks
  tasks, err := loadTasks()
  if err != nil {
    fmt.Println("Error while loading the task: ", err)
    return
  }

  // as per positional args run the corresponding function
  typeOp := os.Args[1]
  switch typeOp {
    case "add" : tasks = addTask(tasks)
    case "list": listTask(tasks)
    case "update": updateTask(tasks)
    case "delete": tasks = deleteTask(tasks)
    case "mark-todo":  markTodo(tasks)
    case "mark-in-progress": markInProgress(tasks)
    case "mark-done": markDone(tasks)
  }

  // save the tasks in file
  err = saveTasks(tasks)
  if err != nil {
    fmt.Println("Error while saving tasks: ", err)
  }

  // bye bye
}
