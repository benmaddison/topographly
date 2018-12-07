package gql

import (
  "github.com/graphql-go/graphql"
)

func makeSchema() (schema graphql.Schema, err error) {
  schemaConfig := graphql.SchemaConfig{
    Query: queryType,
    Mutation: mutationType,
  }
  schema, err = graphql.NewSchema(schemaConfig)
  return
}
