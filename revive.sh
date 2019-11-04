#!/usr/bin/env bash

revive -config revive.toml  -formatter friendly ./...
golangci-link run