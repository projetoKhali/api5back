#!/bin/bash

set -euo pipefail

# shellcheck disable=SC1091
. .env

for prefix in DB DW; do
    echo -e "\n★ Dropping everything in ${prefix}\n"
    HOST=$(eval echo \$${prefix}_HOST)
    PORT=$(eval echo \$${prefix}_PORT)
    USER=$(eval echo \$${prefix}_USER)
    PASS=$(eval echo \$${prefix}_PASS)
    NAME=$(eval echo \$${prefix}_NAME)
    NAME_OR_EMPTY=""
    if [ -n "$NAME" ]; then
        NAME_OR_EMPTY="-d$NAME"
    fi
    PGPASSWORD="$PASS" psql -h "$HOST" -U "$USER" -p "$PORT" "$NAME_OR_EMPTY" -c \
        "drop schema public cascade; create schema public"
    echo -e "\n★ Done dropping everything in ${prefix}\n"
done
