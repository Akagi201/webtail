#!/usr/bin/env bash

echo "go-bindata template files..."
go-bindata data/...
goimports -w bindata.go
