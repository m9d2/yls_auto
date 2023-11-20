#!/bin/bash

find static -type f -not -name 'web.go' -exec rm {} +
cd ./web && npm install && npm run build && cp -r ./dist/. ../static/
cd ../
go build -o ysl 
