package gql

import (
  "github.com/graphql-go/graphql"
)

func makeNodeType() (obj *graphql.Object) {
  obj = graphql.NewObject(
    graphql.ObjectConfig{
      Name: "Node",
      Fields: graphql.Fields{
        "hostname": &graphql.Field{
          Type: graphql.String,
        },
      },
    },
  )
  obj.AddFieldConfig(
    "neighbors", &graphql.Field{
      Type: graphql.NewList(obj),
      Resolve: getNeighbors,
    },
  )
  return
}
var nodeType = makeNodeType()

var linkType = graphql.NewObject(
  graphql.ObjectConfig{
    Name: "Link",
    Fields: graphql.Fields{
      "ipPrefix": &graphql.Field{
        Type: graphql.String,
      },
      "endpoints": &graphql.Field{
        Type: graphql.NewList(nodeType),
        Resolve: getEndpoints,
      },
    },
  },
)

var topologyType = graphql.NewObject(
  graphql.ObjectConfig{
    Name: "Topology",
    Fields: graphql.Fields{
      "nodes": &graphql.Field{
        Type: graphql.NewList(nodeType),
        Resolve: getNodes,
      },
      "links": &graphql.Field{
        Type: graphql.NewList(linkType),
        Resolve: getLinks,
      },
    },
  },
)
