#!/bin/sh

function join_by { local IFS="$1"; shift; echo "$*"; }

IFS=',' read -r -a rolemap <<< "$APIKEYROLEMAP"

for i in {1..9}; do
    var="\$APIKEYROLEMAP_${i}_KEY"
    key=$(eval echo "$var")
    var="\$APIKEYROLEMAP_${i}_ROLE"
    role=$(eval echo "$var")
    if [ ! -z "$key" ]; then
        rolemap+=("$key:$role")
    fi
done

export APIKEYROLEMAP=$(join_by , "${rolemap[@]}")

/$GO_PROJECT_NAME
