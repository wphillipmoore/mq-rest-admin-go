#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
mq_dev_env="${MQ_DEV_ENV_PATH:-${repo_root}/../mq-rest-admin-dev-environment}"

if [ ! -d "$mq_dev_env" ]; then
  echo "mq-rest-admin-dev-environment not found at: $mq_dev_env" >&2
  echo "Clone it as a sibling directory or set MQ_DEV_ENV_PATH." >&2
  exit 1
fi

export COMPOSE_PROJECT_NAME=mqrest-go
export QM1_REST_PORT=9463
export QM2_REST_PORT=9464
export QM1_MQ_PORT=1434
export QM2_MQ_PORT=1435

cd "$mq_dev_env"
exec scripts/mq_stop.sh
