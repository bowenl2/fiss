#!/bin/bash
PWD="$(dirname $0)"
pushd $PWD
go-bindata -prefix="templates/" templates/
popd
