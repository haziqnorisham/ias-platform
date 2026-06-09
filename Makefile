BINARY_NAME = ias_automation_v0.01
FRONTEND_REPO = https://github.com/haziqnorisham/ias_hc.git
FRONTEND_CACHE = .frontend
FRONTEND_DIST = frontend-dist

.PHONY: dev release frontend-pull frontend-rebuild clean clean-build run-dev

dev:
	@echo "==> Building for development (no embedded frontend)"
	@[ -f .env ] || (echo "ERROR: .env file missing. Copy example.env to .env and configure it." && exit 1)
	go build -o $(BINARY_NAME) .
	@echo "==> Build complete: $(BINARY_NAME)"

release: frontend-rebuild
	@echo "==> Building release with embedded frontend"
	go build -tags embed -o $(BINARY_NAME) .
	@echo "==> Release build complete: $(BINARY_NAME)"

frontend-pull:
	@if [ -d $(FRONTEND_CACHE) ]; then \
		echo "==> Pulling latest frontend changes..."; \
		cd $(FRONTEND_CACHE) && git pull; \
	else \
		echo "==> Cloning frontend repository..."; \
		git clone $(FRONTEND_REPO) $(FRONTEND_CACHE); \
	fi

frontend-rebuild: frontend-pull
	@echo "==> Installing frontend dependencies..."
	cd $(FRONTEND_CACHE) && npm ci
	@echo "==> Building frontend (Vite)..."
	cd $(FRONTEND_CACHE) && npm run build
	@rm -rf $(FRONTEND_DIST)
	cp -r $(FRONTEND_CACHE)/dist $(FRONTEND_DIST)
	@echo "==> Frontend copied to $(FRONTEND_DIST)/"

clean:
	rm -rf $(BINARY_NAME) $(FRONTEND_DIST) $(FRONTEND_CACHE)

clean-build:
	rm -f $(BINARY_NAME)

run-dev:
	@echo "==> Running in development mode (no embedded frontend)"
	go run .
