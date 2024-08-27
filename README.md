# Backend Template

Backend Template Go

## Getting Started


### local development

* install go
```shell
go mod download
# .env file used for local development
docker compose -f docker-compose.database.yml up -d
# custom env-file
docker compose -f docker-compose.database.yml --env-file {custom-env-file} up -d
go run .
```
* load sql-dump to setup database


### Deployment

* edit env-file (.env)

```shell
docker compose -f docker-compose.prod.yml up -d
```

