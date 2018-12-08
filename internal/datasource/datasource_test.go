package datasource

import (
  "testing"
)

var path string = "nbs:testdata::test"

func TestDatasource(t *testing.T) {
  d, err := New(path)
  if err != nil {
    t.Error(err)
  }
  err = d.Init()
  if err != nil {
    t.Error(err)
  }
}
