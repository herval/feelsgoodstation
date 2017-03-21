#!/bin/sh

GOARM=6 GOARCH=arm GOOS=linux go build -o station_pi main.go