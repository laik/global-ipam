#!/bin/bash

set -u -e

exit_with_error(){
  echo $1
  exit 1
}

GLOBAL_IPAM_BIN_SRC=/global-ipam
GLOBAL_IPAM_DST=/opt/cni/bin/global-ipam

yes | cp -f $GLOBAL_IPAM_BIN_SRC $GLOBAL_IPAM_DST || exit_with_error "Failed to copy $GLOBAL_IPAM_BIN_SRC to $GLOBAL_IPAM_DST"

echo "install global-ipam binary"

./cni-server