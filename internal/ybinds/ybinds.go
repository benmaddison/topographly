package ybinds

import (
  "github.com/openconfig/ygot/ytypes"
  "github.com/benmaddison/topographly/internal/ybinds/v1"
)

// name of the current schema version
var currentSchema string = "topology-v1"

// aliases for current binding types
type Root = v1.Root
type Topology = v1.Topology
type Node = v1.Topology_Node
type Link = v1.Topology_Link

var schemaRegister = map[string]func()(*ytypes.Schema, error){
  "topology-v1": v1.Schema,
}
