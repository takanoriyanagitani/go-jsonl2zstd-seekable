#!/bin/bash

wsm="./jsonl2zseek.wasm"

run_wasi() {
	cat /dev/stdin |
		wazero run "${wsm}"
}

input() {
	jq -c -n '{id:42, timestamp:634, severity:"INFO", body:"apt update done"}'
	jq -c -n '{id:43, timestamp:635, severity:"WARN", body:"apt update failure"}'
}

echo the original
input | jq -c
echo

echo decompressed
input | run_wasi | zstdcat | jq -c
echo

echo compressed
input | run_wasi | xxd
echo
