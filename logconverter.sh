#!/bin/bash

LOGS_DIR="${1:-logs}"
OUT_DIR="${2:-logstxt}"

mkdir -p "$OUT_DIR"

inotifywait -m -e close_write --format '%f' "$LOGS_DIR" | while read filename; do
    input="$LOGS_DIR/$filename"
    output="$OUT_DIR/${filename%.*}.txt"
    sleep 1
    ansi2txt < "$input" > "$output"
    echo "Processed: $filename -> $output"
done
