build:
	@go build -o  bin/adus


run: build
	@./bin/adus