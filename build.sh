#!/bin/bash

BuildDocker() {
  go build -o ysl .
}

BuildExecutableFile() {
  ehco "build file"
}

if [ "$1" = "docker" ]; then
    BuildDocker
elif [ "$1" = "file" ]; then
    BuildExecutableFile
else
  echo -e "Parameter error"
fi