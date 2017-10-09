#!/bin/bash
set -euo pipefail
IFS=$'\n\t'

dir=$(dirname ${BASH_SOURCE[0]})


function production {
  ${dir}/migrate -database postgres://${POSTGRES_USER}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -path ${dir}/db up
}


if [[ ! ${POSTGRES_HOST:?Requires POSTGRES_HOST} \
   || ! ${POSTGRES_USER:?Requires POSTGRES_USER} \
   || ! ${POSTGRES_DB:?Requires POSTGRES_DB} \
   ]]; then
  exit 1
fi

production
