.PHONY: build build-all build-acp build-cli build-web clean run-cli run-web run-frontend run-backend dev web web-backend

BACKEND_PORT := 8090
FRONTEND_PORT := 5173

build-all:
	@echo "Building all applications..."
	go build ./...
	@echo "All applications built successfully!"

build:
	@mkdir -p bin
	@echo "Building acp..."
	go build -o bin/deepagent-acp ./cmd/acp
	@echo "Building cli..."
	go build -o bin/deepagent-cli ./cmd/cli
	@echo "Building web..."
	go build -o bin/deepagent-web ./cmd/web
	@echo "All binaries built to bin/"

build-acp:
	@mkdir -p bin
	go build -o bin/deepagent-acp ./cmd/acp

build-cli:
	@mkdir -p bin
	go build -o bin/deepagent-cli ./cmd/cli

build-web:
	@mkdir -p bin
	go build -o bin/deepagent-web ./cmd/web

run-cli:
	go run ./cmd/cli

run-backend:
	go run ./cmd/web/main.go

run-frontend:
	cd web && npm run dev

dev:
	@echo "Starting backend and frontend..."
	@trap 'kill 0' EXIT; \
	$(MAKE) run-backend & \
	$(MAKE) run-frontend & \
	wait

web:
	@if [ -z "$$LOCAL_MEMORY" ]; then \
		echo "❌ LOCAL_MEMORY environment variable is not set. Please set it first."; \
		echo "   Example: export LOCAL_MEMORY=/path/to/local_memory"; \
		exit 1; \
	fi; \
	mkdir -p "$$LOCAL_MEMORY/workspace" "$$LOCAL_MEMORY/agent"; \
	chmod 755 "$$LOCAL_MEMORY/workspace" "$$LOCAL_MEMORY/agent"; \
	echo "✅ LOCAL_MEMORY directories ready: $$LOCAL_MEMORY/{workspace,agent}"
	@echo "🚀 Starting web services..."
	@backend_pid=$$(lsof -ti:$(BACKEND_PORT) 2>/dev/null); \
	frontend_pid=$$(lsof -ti:$(FRONTEND_PORT) 2>/dev/null); \
	if [ -n "$$backend_pid" ]; then \
		echo "🔄 Stopping existing backend (pid: $$backend_pid)..."; \
		kill $$backend_pid 2>/dev/null || true; \
		sleep 1; \
	fi; \
	if [ -n "$$frontend_pid" ]; then \
		echo "🔄 Stopping existing frontend (pid: $$frontend_pid)..."; \
		kill $$frontend_pid 2>/dev/null || true; \
		sleep 1; \
	fi; \
	echo "📦 Building backend..."; \
	mkdir -p bin; \
	go build -o bin/deepagent-web ./cmd/web || exit 1; \
	echo "✅ Backend built successfully"; \
	echo "🌐 Starting backend on port $(BACKEND_PORT)..."; \
	nohup ./bin/deepagent-web > /tmp/deepagent-backend.log 2>&1 & \
	echo "🎨 Starting frontend on port $(FRONTEND_PORT)..."; \
	cd web && ( \
		if [ ! -x node_modules/.bin/vite ]; then \
			echo "📦 Installing frontend dependencies (first run)..."; \
			npm ci || npm install || exit 1; \
		fi; \
		nohup npm run dev > /tmp/deepagent-frontend.log 2>&1 & \
	); \
	sleep 2; \
	echo ""; \
	echo "================================================"; \
	echo "  Backend:  http://localhost:$(BACKEND_PORT)"; \
	echo "  Frontend: http://localhost:$(FRONTEND_PORT)"; \
	echo "================================================"; \
	echo ""; \
	echo "Stop: make web-stop"; \
	echo ""; \
	echo "Logs (backend):  tail -f /tmp/deepagent-backend.log"; \
	echo "Logs (frontend): tail -f /tmp/deepagent-frontend.log"; \
	echo ""; \
	echo "📜 Following backend logs (Ctrl+C stops log follow only)...";

web-backend:
	@echo "🔄 Restarting web backend..."
	@backend_pid=$$(lsof -ti:$(BACKEND_PORT) 2>/dev/null); \
	if [ -n "$$backend_pid" ]; then \
		echo "Stopping existing backend (pid: $$backend_pid)..."; \
		kill $$backend_pid 2>/dev/null || true; \
		sleep 1; \
	fi; \
	echo "📦 Building backend..."; \
	mkdir -p bin; \
	go build -o bin/deepagent-web ./cmd/web || exit 1; \
	echo "✅ Backend built successfully"; \
	echo "🌐 Starting backend on port $(BACKEND_PORT)..."; \
	nohup sh -c 'trap "" INT; exec ./bin/deepagent-web' > /tmp/deepagent-backend.log 2>&1 & \
	sleep 1; \
	echo "✅ Backend restarted: http://localhost:$(BACKEND_PORT)"

web-stop:
	@echo "🛑 Stopping web services..."
	@backend_pid=$$(lsof -ti:$(BACKEND_PORT) 2>/dev/null); \
	frontend_pid=$$(lsof -ti:$(FRONTEND_PORT) 2>/dev/null); \
	if [ -n "$$backend_pid" ]; then \
		echo "Stopping backend (pid: $$backend_pid)..."; \
		kill $$backend_pid 2>/dev/null || true; \
	fi; \
	if [ -n "$$frontend_pid" ]; then \
		echo "Stopping frontend (pid: $$frontend_pid)..."; \
		kill $$frontend_pid 2>/dev/null || true; \
	fi; \
	echo "✅ All services stopped"

web-status:
	@echo "📊 Web services status:"
	@backend_pid=$$(lsof -ti:$(BACKEND_PORT) 2>/dev/null); \
	frontend_pid=$$(lsof -ti:$(FRONTEND_PORT) 2>/dev/null); \
	if [ -n "$$backend_pid" ]; then \
		echo "  Backend:  ✅ Running (pid: $$backend_pid, port: $(BACKEND_PORT))"; \
	else \
		echo "  Backend:  ❌ Not running"; \
	fi; \
	if [ -n "$$frontend_pid" ]; then \
		echo "  Frontend: ✅ Running (pid: $$frontend_pid, port: $(FRONTEND_PORT))"; \
	else \
		echo "  Frontend: ❌ Not running"; \
	fi

clean:
	rm -rf bin
