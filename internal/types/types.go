package types

import (
  "fmt"
  "os"
  "encoding/json"
  "github.com/openconfig/ygot/ygot"
  nomsmarshal "github.com/attic-labs/noms/go/marshal"
  nomstypes "github.com/attic-labs/noms/go/types"
)

// to generate Go-bindings from yang schema, using ygot generator:
//   $ go get github.com/openconfig/ygot/generator
//   $ go generate github.com/benmaddison/topographly/internal/types
//
//go:generate generator -path=../../schema -package_name=types -output_file=generated.go -generate_delete -generate_fakeroot -fakeroot_name=root -compress_paths ../../schema/topology.yang

func NewRoot() (root *Root) {
  return &Root{
    Topology: &Topology{},
  }
}

func (root Root) MarshalNoms(vrw nomstypes.ValueReadWriter) (val nomstypes.Value, err error) {
  err = root.Validate()
  if err != nil {
    return
  }
  j, err := ygot.EmitJSON(&root, &ygot.EmitJSONConfig{
    Format: ygot.RFC7951,
  })
  if err != nil {
    return
  }
  fmt.Fprintf(os.Stdout, "json value:\n%s\n", j)
  i := make(map[string]interface{})
  err = json.Unmarshal([]byte(j), &i)
  if err != nil {
    return
  }
  val, err = nomsmarshal.Marshal(vrw, i)
  return
}

func (root *Root) UnmarshalNoms(v nomstypes.Value) (err error) {
  i := make(map[string]interface{})
  err = nomsmarshal.Unmarshal(v, &i)
  if err != nil {
    return
  }
  j, err := json.Marshal(i)
  if err != nil {
    return
  }
  err = Unmarshal(j, root)
  return
}
