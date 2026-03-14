# Architecture Overview

## Summary

This monorepo contains 8 microservices and 6 shared libraries, organized as a Go workspace (`go.work`). All modules share a single dependency graph managed via `go work sync`.

## Services (apps/)

| Service                | Description                              |
|------------------------|------------------------------------------|
| api-gateway            | Public-facing HTTP gateway, routes to internal services |
| auth-service           | Authentication and token issuance        |
| billing-service        | Payment processing and invoicing         |
| orders-service         | Order lifecycle management               |
| notifications-service  | Email, SMS, and push notifications       |
| search-service         | Full-text search over domain entities    |
| worker-ingest          | Async ingestion of external data         |
| worker-reconcile       | Periodic reconciliation of billing data  |

## Shared Libraries (libs/)

| Library       | Purpose                                  |
|---------------|------------------------------------------|
| authz         | Authorization policies and middleware    |
| eventbus      | Publish/subscribe abstraction            |
| retry         | Retry logic with backoff strategies      |
| httpclient    | Standardized HTTP client with tracing    |
| featureflags  | Feature flag evaluation                  |
| observability | Logging, metrics, and tracing helpers    |

## Dependency Diagram

```
                        +-----------+
                        | api-gateway|
                        +-----+-----+
                              |
          +-------------------+-------------------+
          |                   |                   |
   +------+------+    +------+------+    +-------+------+
   | auth-service |    |orders-service|   |search-service|
   +------+------+    +------+------+    +-------+------+
          |                   |                   |
          |            +------+------+            |
          |            |billing-service|           |
          |            +------+------+            |
          |                   |                   |
          |   +---------------+--------+          |
          |   |notifications-service   |          |
          |   +------------------------+          |
          |                                       |
   +------+-------+                +--------------+
   |worker-ingest  |               |worker-reconcile|
   +--------------+                +--------------+

Shared library usage:

  api-gateway           -> authz, httpclient, observability, featureflags
  auth-service          -> authz, observability
  billing-service       -> retry, httpclient, eventbus, observability
  orders-service        -> authz, eventbus, httpclient, observability
  notifications-service -> eventbus, retry, httpclient, observability
  search-service        -> httpclient, observability, featureflags
  worker-ingest         -> eventbus, retry, observability
  worker-reconcile      -> retry, eventbus, observability
```
