#!/usr/bin/env bash

# This should be run from the root of the repo.
# No attempt is made to fix up relative paths.

YANG_DIR="yang"

# Set up a venv to run the yang linter
VENV_PATH=".venv"
echo "Creating venv at $VENV_PATH"
/usr/bin/env virtualenv $VENV_PATH
echo "Installing pyang"
$VENV_PATH/bin/pip install -U pip pyang
echo "Getting openconfig pyang plugin"
/usr/bin/env wget -O $VENV_PATH/oc-pyang.zip https://github.com/openconfig/oc-pyang/archive/master.zip
/usr/bin/env unzip -u $VENV_PATH/oc-pyang.zip -d $VENV_PATH/oc-pyang/

echo "Running oc-pyang linter"
if ! $VENV_PATH/bin/pyang --plugindir=$VENV_PATH/oc-pyang/oc-pyang-master/openconfig_pyang/plugins/ --lint $YANG_DIR/topology-v*.yang; then
  exit 1
fi

GEN_PATH="vendor/github.com/openconfig/ygot/generator/generator.go"
GEN_ARGS="-path=yang -generate_delete -generate_getters -generate_fakeroot -fakeroot_name=root -compress_paths"

# find matching yang schemas
echo "Generating Go bindings"
for F in $YANG_DIR/topology-v*.yang; do
  # get the major version identifier
  V=$(echo $F | sed -E 's/^.*topology-(v.+)\.yang$/\1/')
  # set the output path for the generated bindings
  OUTDIR="internal/ybinds/$V"
  OUTFILE="$OUTDIR/generated.go"
  # set up the generator command
  GEN_EXEC="$GEN_PATH $GEN_ARGS -package_name=$V -output_file=$OUTFILE $F"
  # create the output path and generate the bindings
  mkdir -p $OUTDIR && go run $GEN_EXEC && echo "[OK] $GEN_EXEC"
done
