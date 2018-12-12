package gql

import (
  "fmt"
  "github.com/graphql-go/graphql"
  "github.com/openconfig/ygot/ygot"
  "github.com/benmaddison/topographly/internal/datasource"
  "github.com/benmaddison/topographly/internal/ybinds"
)

func getInstance(p graphql.ResolveParams) (ins *ybinds.Instance, err error) {
  ins, ok := p.Info.RootValue.(map[string]interface{})["instance"].(*ybinds.Instance)
  if !ok {
    err = fmt.Errorf("Could not get data instance from root object: %v\n", p.Info.RootValue)
  }
  return
}

func getDatasource(p graphql.ResolveParams) (d *datasource.Datasource, err error) {
  d, ok := p.Info.RootValue.(map[string]interface{})["datasource"].(*datasource.Datasource)
  if !ok {
    err = fmt.Errorf("Could not get datasource from root object: %v\n", p.Info.RootValue)
  }
  return
}

func putInstance(p graphql.ResolveParams, ins *ybinds.Instance) (err error) {
  d, err := getDatasource(p)
  if err != nil {
    return
  }
  _, err = d.PutHead(ins)
  return
}

func getTopology(p graphql.ResolveParams) (val interface{}, err error) {
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  val = ins.Root.GetOrCreateTopology()
  return
}

func getNodes(p graphql.ResolveParams) (val interface{}, err error) {
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  nodes := make([]*ybinds.Node, 0, len(topology.Node))
  for _, node := range topology.Node {
    nodes = append(nodes, node)
  }
  val = &nodes
  return
}

func getNeighbors(p graphql.ResolveParams) (val interface{}, err error) {
  node, ok := p.Source.(*ybinds.Node)
  if !ok {
    err = fmt.Errorf("Expected *ybinds.Node, got %v\n", p.Source)
    return
  }
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  neighborMap := make(map[string]*ybinds.Node)
  for _, link := range topology.Link {
    if *node.Hostname == *link.EndpointA {
      if _, ok := neighborMap[*link.EndpointZ]; !ok {
        neighborMap[*link.EndpointZ] = topology.Node[*link.EndpointZ]
      }
    }
    if *node.Hostname == *link.EndpointZ {
      if _, ok := neighborMap[*link.EndpointA]; !ok {
        neighborMap[*link.EndpointA] = topology.Node[*link.EndpointA]
      }
    }
  }
  neighbors := make([]*ybinds.Node, 0, len(neighborMap))
  for _, node := range neighborMap {
    neighbors = append(neighbors, node)
  }
  val = &neighbors
  return
}

func addNode(p graphql.ResolveParams) (val interface{}, err error) {
  hostname := p.Args["hostname"].(string)
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  _, err = topology.NewNode(hostname)
  if err != nil {
    return
  }
  err = putInstance(p, ins)
  if err != nil {
    return
  }
  val = topology
  return
}

func delNode(p graphql.ResolveParams) (val interface{}, err error) {
  hostname := p.Args["hostname"].(string)
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  topology.DeleteNode(hostname)
  err = putInstance(p, ins)
  if err != nil {
    return
  }
  val = topology
  return
}

func getLinks(p graphql.ResolveParams) (val interface{}, err error) {
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  links := make([]*ybinds.Link, 0, len(topology.Link))
  for _, link := range topology.Link {
    links = append(links, link)
  }
  val = &links
  return
}

func getEndpoints(p graphql.ResolveParams) (val interface{}, err error) {
  link, ok := p.Source.(*ybinds.Link)
  if !ok {
    err = fmt.Errorf("Expected *ybinds.Link, got %v\n", p.Source)
    return
  }
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  endpoints := []*ybinds.Node{
    topology.Node[*link.EndpointA],
    topology.Node[*link.EndpointZ],
  }
  val = &endpoints
  return
}

func addLink(p graphql.ResolveParams) (val interface{}, err error) {
  prefix := p.Args["ipPrefix"].(string)
  endpointA := p.Args["endpointA"].(string)
  endpointZ := p.Args["endpointZ"].(string)
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  link, err := topology.NewLink(prefix)
  if err != nil {
    return
  }
  link.EndpointA = ygot.String(endpointA)
  link.EndpointZ = ygot.String(endpointZ)
  err = putInstance(p, ins)
  if err != nil {
    return
  }
  val = topology
  return
}

func delLink(p graphql.ResolveParams) (val interface{}, err error) {
  prefix := p.Args["ipPrefix"].(string)
  ins, err := getInstance(p)
  if err != nil {
    return
  }
  topology := ins.Root.GetOrCreateTopology()
  topology.DeleteLink(prefix)
  err = putInstance(p, ins)
  if err != nil {
    return
  }
  val = topology
  return
}
