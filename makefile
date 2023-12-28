run: build
	@./bin/pipcas
build:
	@go build -o ./bin/pipcas ./cmd 