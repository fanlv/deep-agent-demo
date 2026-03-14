# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
make build          # Build all three binaries (acp, cli, web)
make build-all      # go build ./...
make dev            # Start backend (port 8090) + frontend (port 5173) in parallel
make run-cli        # Run CLI: go run ./cmd/cli
make run-backend    # Run web backend only (port 8090)
make run-frontend   # Run frontend only (npm run dev, port 5173)
make web            # Full setup: kill existing, build, run backend + frontend
make web-stop       # Stop all services
```

Frontend (from `web/`):
```bash
npm run dev         # Vite dev server
npm run build       # TypeScript check + Vite build
npm run lint        # ESLint
```

Go tests: `go test ./...`

## Development Conventions

- Logging: use `pkg/logger` (`logger.Infof/Warnf/Errorf/...`). Avoid `log.Printf`, `fmt.Printf`, etc.
- Tests: during development, do not add unit tests unless explicitly requested.

### Code Layering

新增接口时必须严格分层，禁止在 handler 中直接操作文件/数据库：

1. **`types/model/`** — 定义 Request/Response 结构体，handler 和 service 共用
2. **`types/path/`** — 所有文件路径拼接逻辑统一放在此包（如 `PromptsDir()`、`ModelsConfigFile()`），其他包通过调用 path 包获取路径
3. **`repository/`** — 数据持久化层，负责文件/数据库的读写操作，调用 `types/path` 获取路径
4. **`services/`** — 业务逻辑层，调用 repository 完成数据操作
5. **`cmd/web/handler/`** — HTTP 层，只做参数校验和 JSON 响应，核心逻辑委托给 service

### API & Route Conventions

- 路由注册在 `cmd/web/api.go` 的 `registerRoutes` 中
- 接口统一使用 POST + JSON body 传参（包括查询类接口），不使用 query params
- 文件存储统一使用 `.md` 后缀

## Architecture

Multi-interface AI agent system with three entry points:

- **`cmd/web/`** — REST API server (Hertz framework, port 8090) with SSE streaming for real-time agent responses
- **`cmd/cli/`** — Interactive CLI interface
- **`cmd/acp/`** — ACP (Agent Communication Protocol) agent for external integration

### Core Packages

- **`services/agent/`** — Core agent service. `deepagent.go` wraps Eino (CloudWego AI framework). Sub-packages: `chatctx/` (chat context), `middlewares/` (planning, doc loading, tool wrapping), `sandbox/` (Docker sandbox), `prompt/` (system prompts)
- **`services/session/`** — Session/thread lifecycle management (in-memory + persistence)
- **`services/config/`** — Model configuration management
- **`pkg/modelbuilder/`** — LLM provider adapters: ARK, OpenAI, Claude, DeepSeek, Gemini, Ollama, Qwen. Configured via environment variables (see `pkg/modelbuilder/config.go`)
- **`repository/`** — Persistence layer (sessions, chat contexts, model configs stored in sandbox workspace)
- **`types/model/`** — Shared data types (Session, Message, Event, Request/Response)
- **`cmd/web/handler/`** — HTTP handlers: agent endpoints, thread CRUD, model config, SSE streaming

### Frontend (`web/src/`)

React 19 + TypeScript + Vite. Key files:
- `hooks/useAgentChat.ts` — SSE-based agent communication hook
- `utils/sse-client.ts` — SSE client wrapper
- `components/ChatBot.tsx` — Main chat interface
- `types/protocol.ts` — Event type enums matching backend

### Request Flow

1. Frontend sends POST `/api/v1/agent/run` with messages
2. Backend creates/resumes session, invokes agent via Eino Deep Agent framework
3. Agent executes with tools in Docker sandbox
4. Events stream back via SSE (thinking, tool calls, text chunks)

### API Endpoints

```
POST /api/v1/agent/init              — Create session
POST /api/v1/agent/run               — Run agent (SSE streaming)
GET  /api/v1/sessions                — List sessions
GET  /api/v1/sessions/:id/messages   — Get session messages
DELETE /api/v1/sessions/:id          — Delete session
POST /api/v1/prompt/get              — Get prompt by key
POST /api/v1/prompt/save             — Save prompt by key
GET  /api/v1/config/model/list       — List model configs
POST /api/v1/config/model/create     — Create model config
POST /api/v1/config/model/delete     — Delete model config
```

## LLM Provider Configuration

Each provider is configured via environment variables (see `pkg/modelbuilder/config.go`):
- ARK: `ARK_MODEL_ID`, `ARK_API_KEY`, `ARK_BASE_URL`
- OpenAI: `OPENAI_MODEL`, `OPENAI_API_KEY`, `OPENAI_BASE_URL`
- Claude: `CLAUDE_MODEL`, `CLAUDE_API_KEY`, `CLAUDE_BASE_URL`
- DeepSeek: `DEEPSEEK_MODEL`, `DEEPSEEK_API_KEY`
- Gemini: `GEMINI_MODEL`, `GEMINI_API_KEY`
- Ollama: `OLLAMA_MODEL`, `OLLAMA_BASE_URL`
- Qwen: `QWEN_MODEL`, `QWEN_API_KEY`
