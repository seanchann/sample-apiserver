#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..


# Generates types_swagger_doc_generated file for the given group version.
# $1: Name of the group version
# $2: Path to the directory where types.go for that group version exists. This
# is the directory where the file will be generated.
kube_swagger_gen_types_swagger_doc() {
  local group_version=$1
  local gv_dir=$2
  local boilerplate_file=$3
  local TMPFILE="${TMPDIR:-/tmp}/types_swagger_doc_generated.$(date +%s).go"

  echo "Generating swagger type docs for ${group_version} at ${gv_dir}"

  go run $SCRIPT_ROOT/cmd/genswaggertypedocs/swagger_type_docs.go --type-src \
    "${gv_dir}/types.go" \
    --header-file "$boilerplate_file" \
    --package-name "${group_version##*/}" \
    --func-dest - \
    >>  "$TMPFILE"

  gofmt -w -s "$TMPFILE"
  mv "$TMPFILE" ""${gv_dir}"/types_swagger_doc_generated.go"
}



GROUP_VERSIONS=(sample/v1alpha1)
GROUP_VERSIONS_PKG_PATH=(pkg/apis/)

loop_i=0
for group_version in "${GROUP_VERSIONS[@]}"; do
  pkg_path=${GROUP_VERSIONS_PKG_PATH[$loop_i]}
  rm -rf $pkg_path/types_swagger_doc_generated.go
  kube_swagger_gen_types_swagger_doc "${group_version}" "$SCRIPT_ROOT/$pkg_path/${group_version}" "./boilerplate.go.txt"
  loop_i=`expr $loop_i + 1`
done

