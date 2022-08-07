#!/usr/bin/env bash

cd ..
go list -u -m -f '{{if .Update}}{{.}}{{end}}' all
