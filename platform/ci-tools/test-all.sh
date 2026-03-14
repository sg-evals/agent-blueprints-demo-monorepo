#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
FAILED=0

for mod in libs/authz libs/eventbus libs/retry libs/httpclient libs/featureflags libs/observability \
           apps/api-gateway apps/billing-service apps/auth-service apps/orders-service \
           apps/notifications-service apps/search-service apps/worker-ingest apps/worker-reconcile; do
    echo "=== Testing $mod ==="
    cd "$ROOT/$mod"
    if ! go test ./...; then
        FAILED=1
        echo "FAIL: $mod"
    fi
    echo ""
done

exit $FAILED
