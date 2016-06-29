#!/usr/bin/env bash

set -eu

rm -rf internal/*
makdir internal/rocksdb
curl -sL https://github.com/facebook/rocksdb/archive/v4.6.1.tar.gz | tar zxf - -C internal/rocksdb --strip-components=1
makdir internal/rocksdb
curl -sL https://github.com/google/snappy/archive/1.1.3.tar.gz | tar zxf - -C internal/snappy --strip-components=1

#git clean -dxf