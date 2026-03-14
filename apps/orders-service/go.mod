module github.com/sg-evals/agent-blueprints-demo-monorepo/apps/orders-service

go 1.22

require (
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz v0.0.0
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus v0.0.0
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/retry v0.0.0
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability v0.0.0
)

replace (
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz => ../../libs/authz
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/eventbus => ../../libs/eventbus
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/retry => ../../libs/retry
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability => ../../libs/observability
)
