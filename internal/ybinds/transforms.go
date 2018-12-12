package ybinds

import (
  "fmt"
)

func normaliseToCurrent(schemaName string, inval map[string]interface{}) (outval map[string]interface{}, err error) {
  switch schemaName {
    case currentSchema:
      outval = inval
      return
    default:
      err = fmt.Errorf("Could not normalise from schema %s\n", schemaName)
      return
  }
}
