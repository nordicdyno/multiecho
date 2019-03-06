TAG ?= test
.PHONY: build
build:
	docker build -t nordicdyno/multiecho:$(TAG) .

.PHONY: install
install:
	go install .

.PHONY: run
run:
	docker run --rm -it nordicdyno/multiecho:$(TAG)
