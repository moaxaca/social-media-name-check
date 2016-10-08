package main

import (
  "encoding/json"
	"fmt"
	"net/http"
  "sync"
)

type callback func(bool)

func fetch(url string, fn callback) {
  isAvailable := true
  // Send request
  res, err := http.Get(url)
  switch res.Status {
    case "200 OK":
      isAvailable = false
  }
  // Handle Exception
  if (err != nil) {
    fmt.Printf("Fetch against: %v failed to complete.\n", url)
  }
  // Execute Callback with status
  fn(isAvailable)
}

func checkAvailability(name string) map[string]bool {
  // Available social networks to check against
  SOCIAL_NETWORK_URLS := [4] string {
   "https://instagram.com/",
   "https://twitter.com/",
   "https://github.com/",
   "https://dribbble.com/",
  }
  results := make(map[string]bool)
  // WaitGroup
  var wg sync.WaitGroup
  wg.Add(len(SOCIAL_NETWORK_URLS))
  for key := range SOCIAL_NETWORK_URLS {
    url := SOCIAL_NETWORK_URLS[key]
    // fetch async - light weight thread
    go fetch(url+name, (func (isAvailable bool) {
      results[url] = isAvailable
      wg.Done()
    }))
  }
  // Wait till complete then return results
  wg.Wait()
  return results
}

func checkJSONDecorator(name string) []byte {
  results := checkAvailability(name)
  jsonb, err := json.Marshal(results)
  // Handle Exception
  if (err != nil) {
    fmt.Printf("Failed to JSON encode. Error: %v.\n", err)
  }
  return jsonb
}

func rootRouteHandler(writer http.ResponseWriter, request *http.Request) {
  name := request.URL.Path[1:]
  jsonb := checkJSONDecorator(name)
  fmt.Fprintf(writer, string(jsonb))
}

func main() {
  fmt.Printf("Initializing server on port 8080.\n")
  http.HandleFunc("/", rootRouteHandler)
  http.ListenAndServe(":8080", nil)
}
