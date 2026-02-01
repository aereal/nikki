#!/bin/bash

set -e

src_dir="$(dirname "$0")"
root="${src_dir}/../../.."
go tool gqlgen generate
