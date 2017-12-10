install:
	go install -v

fmt:
	go fmt
	cd ./lib && go fmt

test:
	cd ./lib && go test -v

image:
	docker build -t cirocosta/auto53 .

.PHONY: install fmt image test
