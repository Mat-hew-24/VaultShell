#!/bin/bash

if [ "$1" = "exec" ]; then
    for arg in "$@"; do
        [[ "$arg" == -* ]] && continue
        [[ "$arg" == "exec" ]] && continue
        CONTAINER="$arg"
        break
    done
    /usr/local/bin/vaultshell enter-container "$CONTAINER"
    exit $?
fi

exec /usr/bin/docker "$@"
