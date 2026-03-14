# agent-blueprints-demo-monorepo

Semi-synthetic microservice monorepo for demonstrating background coding agents.

## Quick Start

```bash
go work sync
./platform/ci-tools/test-all.sh
```

## Structure

```
apps/                  # Microservices
  api-gateway/         # Public-facing HTTP gateway
  auth-service/        # Authentication and token issuance
  billing-service/     # Payment processing and invoicing
  orders-service/      # Order lifecycle management
  notifications-service/ # Email, SMS, and push notifications
  search-service/      # Full-text search
  worker-ingest/       # Async data ingestion
  worker-reconcile/    # Billing reconciliation

libs/                  # Shared libraries
  authz/               # Authorization policies and middleware
  eventbus/            # Publish/subscribe abstraction
  retry/               # Retry logic with backoff
  httpclient/          # HTTP client with tracing
  featureflags/        # Feature flag evaluation
  observability/       # Logging, metrics, and tracing

docs/                  # Documentation
  architecture/        # Architecture overview and diagrams
  runbooks/            # Operational runbooks
  service-catalog/     # Service metadata

tools/                 # Development and build tools
  gen_repo/            # Deterministic repo generator CLI

platform/              # Platform and CI tooling
  ci-tools/            # CI scripts (test-all.sh)
```
