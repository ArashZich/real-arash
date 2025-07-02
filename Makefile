dev:
	ENVIRONMENT=development go run cmd/http/main.go

fresh-dev:
	ENVIRONMENT=development go run cmd/chef/main.go migrate:down && go run cmd/chef/main.go migrate:up && go run cmd/chef/main.go seed && go run cmd/http/main.go

migrate-dev:
	ENVIRONMENT=development go run cmd/chef/main.go migrate && go run cmd/http/main.go

prod:
	ENVIRONMENT=production go run cmd/http/main.go

fresh-prod:
	ENVIRONMENT=production go run cmd/chef/main.go migrate:down && go run cmd/chef/main.go migrate:up && go run cmd/chef/main.go seed && go run cmd/http/main.go

migrate-prod:
	ENVIRONMENT=production go run cmd/chef/main.go migrate && go run cmd/http/main.go
