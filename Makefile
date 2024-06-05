lint-md:
	markdownlint-cli2 '**/*.md'
.PHONY: lint-md

run-db:
	docker-compose up -d mongo1 mongo2 mongo3 redis

build-vendor:
	cd lib/context && go mod vendor && \
	cd ../image && go mod vendor && \
	cd ../log && go mod vendor && \
	cd ../../services/authentication && \
	go mod vendor && \
	cd ../user && go mod vendor && \
	cd ../rating && go mod vendor && \
	cd ../movie && go mod vendor

compose:
	make build-vendor && docker-compose up
