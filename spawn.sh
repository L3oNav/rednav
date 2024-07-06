#!/bin/sh

set -e
tmpFile=$(mktemp)
go build -o "$tmpFile" main.go 
exec "$tmpFile" "$@"