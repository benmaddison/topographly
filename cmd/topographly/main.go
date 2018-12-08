package main

import (
  "fmt"
  "os"
  "net/http"
  "github.com/benmaddison/topographly/internal/datasource"
  "github.com/benmaddison/topographly/internal/gql"
)

func main() {
  var err error

  path := "nbs:data::topology"

  // connect to datastore
  d, err := datasource.New(path)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not get datasource: %s\n", err)
    return
  }
  // initialise topology datastore
  err = d.Init()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not initialise topology: %s\n", err)
    return
  }
  // get request handler
  h, err := gql.GetHandler(d)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not contruct the graphql handler: %s\n", err)
    return
  }
  // start the server
  http.Handle("/graphql", h)
  endpoint := "[::1]:8080"
  fmt.Fprintf(os.Stdout, "Listening on http://%s\n", endpoint)
  err = http.ListenAndServe(endpoint, nil)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Could not start http server at http://%s\n", endpoint)
    return
  }
}
