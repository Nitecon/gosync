#!/bin/bash
for i in `go list -f '{{.Deps}}' | tr "[" " " | tr "]" " " | xargs go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' |grep  "\."`; do go get $i; done
go build main.go