#!/bin/bash
PWD="$(dirname $0)"
pushd "${PWD}" > /dev/null
go-bindata -prefix="templates/" templates/
popd > /dev/null
