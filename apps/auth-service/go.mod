module github.com/sg-evals/agent-blueprints-demo-monorepo/apps/auth-service

go 1.22

require (
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz v0.0.0
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/httpclient v0.0.0
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability v0.0.0
)

replace (
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/authz => ../../libs/authz
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/httpclient => ../../libs/httpclient
	github.com/sg-evals/agent-blueprints-demo-monorepo/libs/observability => ../../libs/observability
)
