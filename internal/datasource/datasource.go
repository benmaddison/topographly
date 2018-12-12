package datasource

import (
  "fmt"
  "os"
  "time"
  "github.com/attic-labs/noms/go/config"
  "github.com/attic-labs/noms/go/datas"
  "github.com/attic-labs/noms/go/marshal"
  nomstypes "github.com/attic-labs/noms/go/types"
  "github.com/benmaddison/topographly/internal/ybinds"
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

func (d *Datasource) GetHead() (ins *ybinds.Instance, err error) {
  head, ok := d.ds.MaybeHead()
  if !ok {
    err = fmt.Errorf("No value at HEAD\n")
    return
  }
  ins, err = ybinds.NewInstance()
  if err != nil {
    return
  }
  err = marshal.Unmarshal(head, ins)
  if err != nil {
    return
  }
  return
}

func (d *Datasource) PutHead(ins *ybinds.Instance) (changed bool, err error) {
  rootValue, err := marshal.Marshal(*d.db, *ins)
  if err != nil {
    return
  }
  head, ok := d.ds.MaybeHeadValue()
  if ok && head.Equals(rootValue) {
    return
  }
  opts := commitOptions(ins)
  *d.ds, err = (*d.db).Commit(*d.ds, rootValue, opts)
  changed = true
  return
}

func (d *Datasource) Init() (err error) {
  _, err = d.GetHead()
  if err != nil {
    // commit an empty topology to the datasource
    fmt.Fprintf(os.Stdout, "No value at HEAD, initialising empty topology\n")
    ins, lerr := ybinds.NewInstance()
    if lerr != nil {
      err = lerr
      return
    }
    _, err = d.PutHead(ins)
  }
  return
}

func commitOptions(ins *ybinds.Instance) (opts datas.CommitOptions) {
  ts := time.Now().Unix()
  opts = datas.CommitOptions{
    Parents: nomstypes.Set{},
    Meta: nomstypes.NewStruct(
      "Meta",
      nomstypes.StructData{
        "timestamp": nomstypes.Number(ts),
        "schema": nomstypes.String(ins.SchemaName),
      },
    ),
    Policy: nil,
  }
  return
}
