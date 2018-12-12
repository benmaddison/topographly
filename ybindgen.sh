#!/usr/bin/env bash

# This should be run from the root of the repo.
# No attempt is made to fix up relative paths.

GEN_PATH="vendor/github.com/openconfig/ygot/generator/generator.go"
GEN_ARGS="-path=yang -generate_delete -generate_getters -generate_fakeroot -fakeroot_name=root -compress_paths"

# find matching yang schemas
for F in yang/topology-v*.yang; do
  # get the major version identifier
  V=$(echo $F | sed -E 's/^.*topology-(v.+)\.yang$/\1/')
  # set the output path for the generated bindings
  OUTFILE="internal/ybinds/$V/generated.go"
  # set up the generator command
  GEN_EXEC="$GEN_PATH $GEN_ARGS -package_name=$V -output_file=$OUTFILE $F"
  # create the output path and generate the bindings
  mkdir -p $OUTDIR && go run $GEN_EXEC && echo "[OK] $GEN_EXEC"
done
