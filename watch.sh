#!/bin/bash

LOGS_DIR="/var/lib/vaultshell/logs"
OUT_DIR="/var/lib/vaultshell/logstxt"

nohup bash /usr/local/bin/logconverter.sh "$LOGS_DIR" "$OUT_DIR" > "/var/lib/vaultshell/watch.log" 2>&1 &
