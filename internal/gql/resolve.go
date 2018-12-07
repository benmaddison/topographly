package gql

import (
  "fmt"
  "github.com/graphql-go/graphql"
  "github.com/openconfig/ygot/ygot"
  "github.com/benmaddison/topographly/internal/datasource"
  "github.com/benmaddison/topographly/internal/types"
)

func getDatasource(p graphql.ResolveParams) (d *datasource.Datasource, err error) {
  d, ok := p.Info.RootValue.(map[string]interface{})["datasource"].(*datasource.Datasource)
  if !ok {
    err = fmt.Errorf("Could not get datastore from root object: %v\n", p.Info.RootValue)
  }
  return
}

func getTopology(p graphql.ResolveParams) (val interface{}, err error) {
  d, err := getDatasource(p)
  if err != nil {
    return
  }
  root, err := d.GetHead()
  if err != nil {
    return
  }
  if root.Topology == nil {
    root.Topology = &types.Topology{}
  }
  val = root.Topology
  return
}

func putTopology(p graphql.ResolveParams, t *types.Topology) (err error) {
  d, err := getDatasource(p)
  if err != nil {
    return
  }
  root := &types.Root{
    Topology: t,
  }
  _, err = d.PutHead(root)
  return
}

func getNodes(p graphql.ResolveParams) (val interface{}, err error) {
  t, ok := p.Source.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", p.Source)
    return
  }
  nodes := make([]*types.Topology_Node, 0, len(t.Node))
  for _, node := range t.Node {
    nodes = append(nodes, node)
  }
  val = &nodes
  return
}

func getNeighbors(p graphql.ResolveParams) (val interface{}, err error) {
  n, ok := p.Source.(*types.Topology_Node)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology_Node, got %v\n", p.Source)
    return
  }
  topology, err := getTopology(p)
  if err != nil {
    return
  }
  t, ok := topology.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", t)
    return
  }
  neighborMap := make(map[string]*types.Topology_Node, 0)
  for _, l := range t.Link {
    if *n.Hostname == *l.EndpointA {
      if _, ok := neighborMap[*l.EndpointZ]; !ok {
        neighborMap[*l.EndpointZ] = t.Node[*l.EndpointZ]
      }
    }
    if *n.Hostname == *l.EndpointZ {
      if _, ok := neighborMap[*l.EndpointA]; !ok {
        neighborMap[*l.EndpointA] = t.Node[*l.EndpointA]
      }
    }
  }
  neighbors := make([]*types.Topology_Node, 0, len(neighborMap))
  for _, node := range neighborMap {
    neighbors = append(neighbors, node)
  }
  val = &neighbors
  return
}

func addNode(p graphql.ResolveParams) (val interface{}, err error) {
  hostname := p.Args["hostname"].(string)
  topology, err := getTopology(p)
  if err != nil {
    return
  }
  t, ok := topology.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", t)
    return
  }
  _, err = t.NewNode(hostname)
  if err != nil {
    return
  }
  err = putTopology(p, t)
  if err != nil {
    return
  }
  val = t
  return
}

func delNode(p graphql.ResolveParams) (val interface{}, err error) {
  hostname := p.Args["hostname"].(string)
  topology, err := getTopology(p)
  if err != nil {
    return
  }
  t, ok := topology.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", t)
    return
  }
  t.DeleteNode(hostname)
  err = putTopology(p, t)
  if err != nil {
    return
  }
  val = t
  return
}

func getLinks(p graphql.ResolveParams) (val interface{}, err error) {
  t, ok := p.Source.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", p.Source)
    return
  }
  links := make([]*types.Topology_Link, 0, len(t.Link))
  for _, link := range t.Link {
    links = append(links, link)
  }
  val = &links
  return
}

func getEndpoints(p graphql.ResolveParams) (val interface{}, err error) {
  l, ok := p.Source.(*types.Topology_Link)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology_Link, got %v\n", p.Source)
    return
  }
  topology, err := getTopology(p)
  if err != nil {
    return
  }
  t, ok := topology.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", t)
    return
  }
  endpoints := []*types.Topology_Node{
    t.Node[*l.EndpointA],
    t.Node[*l.EndpointZ],
  }
  val = &endpoints
  return
}

func addLink(p graphql.ResolveParams) (val interface{}, err error) {
  prefix := p.Args["ipPrefix"].(string)
  endpointA := p.Args["endpointA"].(string)
  endpointZ := p.Args["endpointZ"].(string)
  topology, err := getTopology(p)
  if err != nil {
    return
  }
  t, ok := topology.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", t)
    return
  }
  link, err := t.NewLink(prefix)
  if err != nil {
    return
  }
  link.EndpointA = ygot.String(endpointA)
  link.EndpointZ = ygot.String(endpointZ)
  err = putTopology(p, t)
  if err != nil {
    return
  }
  val = t
  return
}

func delLink(p graphql.ResolveParams) (val interface{}, err error) {
  prefix := p.Args["ipPrefix"].(string)
  topology, err := getTopology(p)
  if err != nil {
    return
  }
  t, ok := topology.(*types.Topology)
  if !ok {
    err = fmt.Errorf("Expected *types.Topology, got %v\n", t)
    return
  }
  t.DeleteLink(prefix)
  err = putTopology(p, t)
  if err != nil {
    return
  }
  val = t
  return
}
