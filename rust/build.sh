#!/bin/bash

if [[ "$(uname)" =~ MSYS*|CYGWIN*|MINGW* ]]; then
  echo 'Using Rust GNU toolchain on Windows'
  rustup override set stable-gnu
fi

if ! command -v cbindgen &>/dev/null; then
  echo 'Downloading cbindgen...'
  cargo install cbindgen
fi

cbindgen --config cbindgen.toml --output include/gol.h
cargo build --release
