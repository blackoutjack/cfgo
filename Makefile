.PHONY: test

golze: *.go */*.go
	go build

test:
	python -m test
