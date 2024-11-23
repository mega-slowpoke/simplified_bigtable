#!/bin/bash
SRC_DIR="."
TAR_DIR="."

protoc -I=$SRC_DIR --go_out=$TAR_DIR external_api.proto