.PHONY: build

build:
	sam build -s cmd

local:
	sam build -s cmd --debug
	sam local start-api -p 3030 --debug --skip-pull-image --warm-containers EAGER

deploy:
	sam build -s cmd
	sam deploy --profile lambda-ci

migrate:
	go run db/migration.go