#!/bin/bash

# Rewrite references from github.com/clear-street/gogen-avro to gopkg.in/clear-street/gogen-avro.<version>

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <version>"
  exit 1
fi
 
GITHUB_REPO="github.com/clear-street/gogen-avro"
VERSION="$1"
GOPKG_REPO="gopkg.in/clear-street/gogen-avro.$VERSION"

sed -i "s|$GITHUB_REPO|$GOPKG_REPO|" container/*.go generator/*.go types/*.go gogen-avro/main.go example/*/*.go test.sh test/*/*.go
