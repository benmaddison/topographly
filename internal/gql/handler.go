package gql

import (
  "context"
  "net/http"
  "github.com/graphql-go/handler"
  "github.com/benmaddison/topographly/internal/datasource"
)

func GetHandler(d *datasource.Datasource) (h http.Handler, err error) {
  rootObjectFn := func(ctx context.Context, r *http.Request) (obj map[string]interface{}) {
    ins, err := d.GetHead()
    if err != nil {
      panic(err)
    }
    obj = map[string]interface{}{
      "datasource": d,
      "instance": ins,
    }
    return
  }

  schema, err := makeSchema()
  if err != nil {
    return
  }

  handlerConfig := handler.Config{
    Schema: &schema,
    Pretty: true,
    GraphiQL: true,
    RootObjectFn: rootObjectFn,
  }
  h = handler.New(&handlerConfig)
  return
}
