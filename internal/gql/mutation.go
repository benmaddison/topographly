package gql

import (
  "github.com/graphql-go/graphql"
)

var mutationType = graphql.NewObject(
  graphql.ObjectConfig{
    Name: "Mutation",
    Fields: graphql.Fields{
      "addNode": &graphql.Field{
        Type: topologyType,
        Resolve: addNode,
        Args: graphql.FieldConfigArgument{
          "hostname": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
        },
      },
      "delNode": &graphql.Field{
        Type: topologyType,
        Resolve: delNode,
        Args: graphql.FieldConfigArgument{
          "hostname": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
        },
      },
      "addLink": &graphql.Field{
        Type: topologyType,
        Resolve: addLink,
        Args: graphql.FieldConfigArgument{
          "ipPrefix": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
          "endpointA": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
          "endpointZ": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
        },
      },
      "delLink": &graphql.Field{
        Type: topologyType,
        Resolve: delLink,
        Args: graphql.FieldConfigArgument{
          "ipPrefix": &graphql.ArgumentConfig{
            Type: graphql.NewNonNull(graphql.String),
          },
        },
      },
    },
  },
)
