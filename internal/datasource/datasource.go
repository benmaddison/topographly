package datasource

import (
  "fmt"
  "os"
  "github.com/attic-labs/noms/go/config"
  "github.com/attic-labs/noms/go/datas"
  "github.com/attic-labs/noms/go/marshal"
  "github.com/benmaddison/topographly/internal/types"
)

type Datasource struct {
  db *datas.Database
  ds *datas.Dataset
}

func New(path string) (d *Datasource, err error) {
  cfg := config.NewResolver()
  db, ds, err := cfg.GetDataset(path)
  if err != nil {
    return
  }
  d = &Datasource{
    db: &db,
    ds: &ds,
  }
  return
}

func (d *Datasource) GetHead() (root *types.Root, err error) {
  head, ok := d.ds.MaybeHeadValue()
  if !ok {
    err = fmt.Errorf("No value at HEAD\n")
    return
  }
  root = types.NewRoot()
  err = marshal.Unmarshal(head, root)
  if err != nil {
    return
  }
  return
}

func (d *Datasource) PutHead(root *types.Root) (changed bool, err error) {
  rootValue, err := marshal.Marshal(*d.db, *root)
  if err != nil {
    return
  }
  head, ok := d.ds.MaybeHeadValue()
  if ok && head.Equals(rootValue) {
    return
  }
  *d.ds, err = (*d.db).CommitValue(*d.ds, rootValue)
  changed = true
  return
}

func (d *Datasource) Init() (err error) {
  _, err = d.GetHead()
  if err != nil {
    // commit an empty topology to the datasource
    fmt.Fprintf(os.Stdout, "No value at HEAD, initialising empty topology\n")
    newRoot := types.NewRoot()
    _, err = d.PutHead(newRoot)
  }
  return
}
