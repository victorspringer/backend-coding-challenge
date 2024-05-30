lint-md:
	markdownlint-cli2 '**/*.md'
.PHONY: lint-md

run-db:
	docker-compose up -d mongo1 mongo2 mongo3

compose:
	docker-compose up
