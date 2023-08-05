#! /bin/sh

set -a
source .env
set +a

goose -dir sql/schema "$@"
