BINARY_NAME = ias_automation_v0.01

.PHONY: dev clean clean-build run-dev

dev:
	@echo "==> Building for development (no embedded frontend)"
	@[ -f .env ] || (echo "ERROR: .env file missing. Copy example.env to .env and configure it." && exit 1)
	CGO_ENABLED=0 go build -o $(BINARY_NAME) .
	@echo "==> Build complete: $(BINARY_NAME)"

clean:
	rm -f $(BINARY_NAME)

clean-build:
	rm -f $(BINARY_NAME)

run-dev:
	@echo "==> Running in development mode (no embedded frontend)"
	go run .
