install:
	go install -v

fmt:
	go fmt
	cd ./lib && go fmt

.PHONY: install fmt
