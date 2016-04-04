#!/bin/bash

mkdir -p binaries
rm binaries/*
cd binaries

GOOS=windows GOARCH=386 go build github.com/ScottBrooks/modmerge/modmerge
zip modmerge-win32.zip modmerge.exe
rm modmerge.exe

GOOS=linux GOARCH=386 go build github.com/ScottBrooks/modmerge/modmerge
gzip modmerge
mv modmerge.gz modmerge-linux.gz

GOOS=darwin go build github.com/ScottBrooks/modmerge/modmerge
gzip modmerge
mv modmerge.gz modmerge-osx.gz


