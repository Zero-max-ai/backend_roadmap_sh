package main

import (
  "os"
  "fmt"
  "log"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

type Event struct {
  Type      string    `json: type`
  Repo      struct  {
    Name    string    `json: name`
  } `json: repo`
  Payload   struct {
    Size    int       `json: size,omitempty`
  } `json: payload`
}

// the OG main function
func main() {
  username := os.Args[1]
  apiUrl := fmt.Sprintf("https://api.github.com/users/%s/events", username)

  // response from the github api
  resp, err := http.Get(apiUrl)
  if err != nil {
    log.Fatalf("Error making request: %v", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode != http.StatusOK {
    log.Fatalf("Unexpected status code: %d)", resp.StatusCode)
  }

  // read all the response body
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println("Error while reading response: ", err)
    return
  }

  var events []Event
  err = json.Unmarshal(body, &events)
  if err != nil {
    fmt.Println("Error while parsing error: ", err)
    return
  }

  if len(events) == 0 {
    fmt.Printf("No recent activity found for user: %s\n", username)
    return
  }

  for _, event := range events {
    switch event.Type {
      case "PushEvent": fmt.Printf("Pushed %d commits to %s\n", event.Payload.Size, event.Repo.Name)
      case "IssuesEvent": fmt.Printf("Opened a new issue in %s\n", event.Repo.Name)
      case "WatchEvent": fmt.Printf("Starred as %s\n", event.Repo.Name)
      default:  fmt.Printf("Performed %s on %s\n", event.Type, event.Repo.Name)
    }
  }

}
