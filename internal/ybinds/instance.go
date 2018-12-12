package ybinds

import (
  "fmt"
  "os"
  "encoding/json"
  "github.com/openconfig/ygot/ygot"
  "github.com/openconfig/ygot/ytypes"
  nomsmarshal "github.com/attic-labs/noms/go/marshal"
  nomstypes "github.com/attic-labs/noms/go/types"
)

type Instance struct {
  SchemaName string
  Schema *ytypes.Schema
  Root *Root
}

func NewInstance() (ins *Instance, err error) {
  ins = &Instance{}
  err = ins.Init()
  return
}

func (ins *Instance) Init() (err error) {
  schemaName := currentSchema
  makeSchema, ok := schemaRegister[schemaName]
  if !ok {
    err = fmt.Errorf("No schema with name %s found", schemaName)
    return
  }
  ins.SchemaName = schemaName
  ins.Root = &Root{}
  schema, err := makeSchema()
  if err != nil {
    return
  }
  ins.Schema = schema
  return
}

func (ins Instance) MarshalNoms(vrw nomstypes.ValueReadWriter) (val nomstypes.Value, err error) {
  j, err := ygot.EmitJSON(ins.Root, &ygot.EmitJSONConfig{
    Format: ygot.RFC7951,
  })
  if err != nil {
    return
  }
  i := make(map[string]interface{})
  err = json.Unmarshal([]byte(j), &i)
  if err != nil {
    return
  }
  val, err = nomsmarshal.Marshal(vrw, i)
  return
}

func (ins *Instance) UnmarshalNoms(v nomstypes.Value) (err error) {
  commit, ok := v.(nomstypes.Struct)
  if !ok {
    err = fmt.Errorf("Can only unmarshal Commit Structs, got %v\n", v)
    return
  }
  meta, ok := commit.MaybeGet("meta")
  if !ok {
    err = fmt.Errorf("Could not get 'meta' from noms Struct")
    return
  }
  value, ok := commit.MaybeGet("value")
  if !ok {
    err = fmt.Errorf("Could not get 'value' from noms Struct")
    return
  }
  schemaName := currentSchema
  metaSchemaName, ok := meta.(nomstypes.Struct).MaybeGet("schema")
  if ok {
    schemaName = string(metaSchemaName.(nomstypes.String))
  }
  inmap := make(map[string]interface{})
  err = nomsmarshal.Unmarshal(value, &inmap)
  if err != nil {
    return
  }
  outmap, err := normaliseToCurrent(schemaName, inmap)
  if err != nil {
    return
  }
  jsonInstance, err := json.Marshal(outmap)
  if err != nil {
    return
  }
  fmt.Fprintf(os.Stdout, "%s\n", jsonInstance)
  err = ins.Init()
  if err != nil {
    return
  }
  err = ins.Schema.Unmarshal(jsonInstance, ins.Root)
  return
}
