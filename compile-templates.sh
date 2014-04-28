#!/bin/bash
PWD="$(dirname $0)"
pushd $PWD
go-bindata -debug -prefix="templates/" templates/
popd