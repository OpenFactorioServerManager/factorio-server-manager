# CLAUDE.md — Factorio Server Manager

This file documents the codebase structure, build workflows, and conventions for AI assistants working on this repository.

---

## Project Overview

**Factorio Server Manager (FSM)** is a self-hosted web application for managing a Factorio game server. It runs as a background service on the same machine as the Factorio binary and provides:

- Start/stop/kill the Factorio server process
- Manage save files (upload, download, delete, create)
- Manage mods (install from mod portal, toggle, update, delete)
- Manage mod packs (grouped mod configurations)
- View game logs and console output in real-time
- Edit server settings and game settings
- User management with session-based authentication

The application is a **Go backend** (REST API + WebSocket) serving a **React frontend** as static files.

---

## Repository Structure

```
factorio-server-manager/
├── app/                        # Compiled frontend output (served statically by Go)
│   └── index.html              # Static HTML shell
├── src/                        # Go backend source (module root)
│   ├── main.go                 # Application entry point
│   ├── go.mod                  # Go module: github.com/OpenFactorioServerManager/factorio-server-manager
│   ├── go.sum
│   ├── bootstrap/
│   │   └── config.go           # CLI flags, JSON config loading, singleton Config
│   ├── factorio/               # Core Factorio domain logic
│   │   ├── server.go           # Server struct, lifecycle (Run/Stop/Kill)
│   │   ├── server_linux.go     # Linux-specific process management
│   │   ├── server_windows.go   # Windows-specific process management
│   │   ├── saves.go            # Save file listing/management
│   │   ├── save.go             # Save file binary format parsing
│   │   ├── save_test.go        # Tests for save parsing across Factorio versions
│   │   ├── mods.go             # Mod file operations (delete all, startup)
│   │   ├── mod_Mods.go         # Mods collection type
│   │   ├── mod_modInfo.go      # Mod metadata (info.json parsing)
│   │   ├── mod_modSimple.go    # Simplified mod representation
│   │   ├── mod_modpack.go      # Mod pack read/write operations
│   │   ├── mod_portal.go       # Factorio mod portal HTTP client
│   │   ├── config.go           # config.ini parsing
│   │   ├── credentials.go      # Factorio account credentials (mod portal login)
│   │   ├── gamelog.go          # Factorio log file tailing
│   │   ├── rcon.go             # RCON console command support
│   │   └── version.go          # Version type [4]uint with binary/text marshalling
│   ├── api/                    # HTTP API layer
│   │   ├── routes.go           # Route table + gorilla/mux router setup
│   │   ├── handlers.go         # Core handlers: saves, server control, users, settings
│   │   ├── auth.go             # Session auth middleware, SQLite user store (GORM)
│   │   ├── mods_handler.go     # Installed mods CRUD handlers
│   │   ├── mod_modpack_handler.go  # Mod pack handlers
│   │   ├── mod_portal_handler.go   # Mod portal proxy handlers
│   │   ├── *_test.go           # Integration tests using httptest
│   │   └── websocket/
│   │       ├── wshub.go        # WebSocket hub: rooms, broadcast, control handlers
│   │       └── wsclient.go     # WebSocket client connection management
│   ├── lockfile/
│   │   └── lockfile.go         # Read/write file locking utility
│   └── factorio_testfiles/     # Factorio save file fixtures for version tests
│       ├── test_0_13.zip ... test_1_1_14.zip
│       └── belt-balancer_*.zip (mod fixtures)
├── ui/                         # React frontend source
│   ├── index.js                # Webpack entry: mounts <App /> to #app
│   ├── index.scss              # Global styles (imports Tailwind)
│   ├── notifications.js        # window.flash() implementation
│   ├── api/
│   │   ├── client.js           # Axios instance (withCredentials, error interceptor)
│   │   ├── socket.js           # Native WebSocket + EventEmitter bus, auto-reconnect
│   │   └── resources/          # Per-domain API functions (thin wrappers over client.js)
│   │       ├── user.js
│   │       ├── server.js
│   │       ├── saves.js
│   │       ├── mods.js
│   │       ├── settings.js
│   │       └── log.js
│   └── App/
│       ├── App.jsx             # Root component: auth state, serverStatus, routing
│       ├── views/              # Full-page route components
│       │   ├── Controls.jsx    # Server start/stop controls
│       │   ├── Login.jsx
│       │   ├── Saves/          # Save file management
│       │   ├── Mods/           # Mod management + mod portal
│       │   ├── ServerSettings.jsx
│       │   ├── GameSettings.jsx
│       │   ├── Console.jsx     # In-game console (WebSocket)
│       │   ├── Logs.jsx        # Log viewer (WebSocket)
│       │   ├── UserManagement/
│       │   └── Help.jsx
│       └── components/         # Reusable UI primitives
│           ├── Layout.jsx      # App shell (nav sidebar)
│           ├── Button.jsx, ButtonLink.jsx
│           ├── Input.jsx, InputPassword.jsx, Select.jsx, Checkbox.jsx
│           ├── Modal.jsx, ConfirmDialog.jsx
│           ├── Panel.jsx, Label.jsx
│           ├── Flash.jsx       # Toast notification display
│           ├── Error.jsx
│           └── Tabs/
├── docker/                     # Docker deployment
│   ├── Dockerfile              # Production image (downloads release from GitHub)
│   ├── Dockerfile-build        # Multi-stage build image
│   ├── Dockerfile-local        # Local dev image
│   ├── docker-compose.yaml     # Full stack with Traefik reverse proxy
│   ├── docker-compose.simple.yaml
│   └── entrypoint.sh
├── Makefile                    # Build orchestration
├── webpack.config.js           # Frontend bundler config
├── tailwind.config.js          # Custom Tailwind theme
├── package.json                # Node dependencies
└── conf.json.example           # Example runtime configuration
```

---

## Build System

### Full build (frontend + backend + zip archive)

```bash
make build                          # Detects OS, builds for Linux or Windows
make gen_release                    # Build both Linux and Windows releases
```

The `make build` process:
1. Runs `npm install && npm run build` to bundle the frontend into `app/`
2. Compiles the Go binary
3. Packages everything into `build/factorio-server-manager-<os>.zip`

### Backend only

```bash
cd src

# Linux
CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o ../factorio-server-manager/factorio-server-manager .

# Windows (requires mingw cross-compiler)
GO111MODULE=on GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
  CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc \
  go build -ldflags="-extldflags=-static" -o ../factorio-server-manager/factorio-server-manager.exe .
```

### Frontend only

```bash
npm install
npm run development    # One-time dev build
npm run watch          # Watch mode (dev)
npm run build          # Production build (minified)
```

Webpack compiles `ui/index.js` → `app/bundle.js` and `ui/index.scss` → `app/style.css`.

### Clean

```bash
make clean    # Removes build/, compiled frontend assets, node_modules/
```

---

## Running Locally

The binary reads from a `conf.json` file. Copy the example first:

```bash
cp conf.json.example conf.json
```

Then run the binary with flags pointing at a local Factorio installation:

```bash
./factorio-server-manager \
  --dir /path/to/factorio \
  --conf ./conf.json \
  --port 8080
```

The web UI is served at `http://localhost:8080`. On first start, an `admin` user is created with a randomly generated password printed to stdout.

### CLI flags / environment variables

All flags have `FSM_*` environment variable equivalents (see `src/bootstrap/config.go`):

| Flag | Env var | Default | Description |
|------|---------|---------|-------------|
| `--conf` | `FSM_CONF` | `./conf.json` | Config file path |
| `--dir` | `FSM_DIR` | `./` | Factorio installation directory |
| `--host` | `FSM_SERVER_IP` | `0.0.0.0` | Web server bind address |
| `--port` | `FSM_PORT` | `80` | Web server port |
| `--bin` | `FSM_BINARY` | `bin/x64/factorio` | Factorio binary (relative to `--dir`) |
| `--mod-pack-dir` | `FSM_MODPACK_DIR` | `./mod_packs` | Mod packs storage directory |
| `--rcon-port` | `FSM_RCON_PORT` | random 40000–45000 | RCON port |
| `--max-upload` | `FSM_MAX_UPLOAD` | `20` | Max upload size in MB |
| `--autostart` | `FSM_AUTOSTART` | `false` | Auto-start Factorio on FSM start |

---

## Configuration File (`conf.json`)

The JSON config file (default `./conf.json`) provides persistent settings. CLI flags are applied first, then the JSON file overrides. Key fields:

```json
{
  "sq_lite_database_file": "sqlite.db",
  "cookie_encryption_key": "",       // Base64 32-byte key; auto-generated if empty
  "rcon_pass": "",                   // Auto-generated if empty
  "settings_file": "server-settings.json",
  "log_file": "factorio-server-manager.log",
  "console_cache_size": 25,          // Lines of console output to cache in memory
  "secure": true                     // Set false for HTTP (no HTTPS); affects cookie SameSite
}
```

On startup, FSM auto-generates `cookie_encryption_key` and `rcon_pass` if missing/default and writes them back to `conf.json`. It also migrates any legacy LevelDB database to SQLite.

---

## Testing

### Go tests

Tests live alongside their packages in `*_test.go` files.

```bash
cd src
go test ./...              # Run all tests
go test ./factorio/...     # Run factorio package tests only
go test ./api/...          # Run API integration tests
```

**API integration tests** (`src/api/*_test.go`) require a `.env` file at `src/.env` with:

```
FSM_DIR=/path/to/factorio
FSM_CONF=/path/to/conf.json
FSM_MODPACK_DIR=/path/to/mod_packs
mod_dir=/path/to/mods
factorio_username=your_factorio_username
factorio_password=your_factorio_password
```

Without `.env`, a warning is logged and tests that require Factorio credentials will fail.

Tests use `github.com/stretchr/testify` for assertions and `net/http/httptest` for HTTP handler testing. The `TestMain` function in `src/api/mods_handler_test.go` sets up the shared test environment.

**Save format tests** (`src/factorio/save_test.go`) parse binary save files from `src/factorio_testfiles/` and verify the header structure across Factorio versions (0.13 through 1.1.14).

### Frontend tests

There are no frontend tests. `npm test` exits with an error by design.

---

## API Reference

All endpoints require authentication (session cookie) except `POST /api/login`. Routes prefixed with `ServerOff: true` return `HTTP 423 Locked` while the Factorio server is running.

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/login` | Login (`{username, password}`) → sets session cookie |
| `GET` | `/api/logout` | Logout, clears session |
| `GET` | `/api/user/status` | Get current logged-in user |
| `GET` | `/api/user/list` | List all users |
| `POST` | `/api/user/add` | Add user `{username, password, role}` |
| `POST` | `/api/user/remove` | Remove user `{username}` |
| `POST` | `/api/user/password` | Change own password `{old_password, new_password, new_password_confirmation}` |

### Server Control

| Method | Path | Server must be off? | Description |
|--------|------|---------------------|-------------|
| `POST` | `/api/server/start` | Yes | Start server `{savefile, bindip, port, latency}` |
| `GET` | `/api/server/stop` | No | Graceful stop |
| `GET` | `/api/server/kill` | No | Force kill |
| `GET` | `/api/server/status` | No | Returns Server struct as JSON |
| `GET` | `/api/server/facVersion` | No | Returns `{version, base_mod_version}` |

### Saves

| Method | Path | Server must be off? | Description |
|--------|------|---------------------|-------------|
| `GET` | `/api/saves/list` | No | List saves; `?latest=true` appends symlink entry |
| `GET` | `/api/saves/dl/{save}` | No | Download save file |
| `POST` | `/api/saves/upload` | No | Upload save (multipart, field `savefile`, `.zip` only) |
| `GET` | `/api/saves/rm/{save}` | No | Delete save |
| `GET` | `/api/saves/create/{save}` | Yes | Create new save via Factorio binary `--create` |
| `POST` | `/api/saves/mods` | Yes | Load mods from save file |

### Mods

| Method | Path | Server must be off? | Description |
|--------|------|---------------------|-------------|
| `GET` | `/api/mods/list` | No | List installed mods |
| `POST` | `/api/mods/toggle` | Yes | Toggle mod enabled/disabled |
| `POST` | `/api/mods/delete` | Yes | Delete a mod |
| `POST` | `/api/mods/delete/all` | Yes | Delete all mods |
| `POST` | `/api/mods/update` | Yes | Update a mod |
| `POST` | `/api/mods/upload` | Yes | Upload mod file (multipart) |
| `GET` | `/api/mods/download` | No | Download all mods as zip |

### Mod Portal

| Method | Path | Server must be off? | Description |
|--------|------|---------------------|-------------|
| `GET` | `/api/mods/portal/list` | No | List all mods from portal |
| `GET` | `/api/mods/portal/info/{mod}` | No | Get mod details |
| `POST` | `/api/mods/portal/install` | Yes | Install single mod from portal |
| `POST` | `/api/mods/portal/install/multiple` | Yes | Install multiple mods |
| `POST` | `/api/mods/portal/login` | No | Login to Factorio account |
| `GET` | `/api/mods/portal/loginstatus` | No | Check portal login status |
| `GET` | `/api/mods/portal/logout` | No | Logout from portal |

### Mod Packs

All mod pack routes follow `/api/mods/packs/{modpack}/...` — see `src/api/routes.go` for the full list.

### Settings & Config

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/settings` | Get server-settings.json contents |
| `POST` | `/api/settings/update` | Update and save server settings |
| `GET` | `/api/config` | Get config.ini contents |
| `GET` | `/api/log/tail` | Get recent Factorio log lines |

### WebSocket (`/ws`)

Requires authentication. Messages are JSON with shape:

```json
{ "room_name": "string", "message": "...", "controls": { "type": "string", "value": "string" } }
```

**Control types** (sent from client):
- `subscribe` — join a room (value: room name)
- `unsubscribe` — leave a room
- `command` — send RCON command to Factorio server

**Rooms**:
- `gamelog` — real-time Factorio console output; replays cached lines on subscribe
- `server_status` — JSON-serialized Server struct on status change

---

## Backend Conventions (Go)

### Configuration singleton
`bootstrap.GetConfig()` returns the global `Config` struct. It is set once by `bootstrap.NewConfig()` in `main.go`.

### Factorio Server singleton
`factorio.GetFactorioServer()` returns a pointer to the global `Server` instance. Always use this rather than creating new instances.

### HTTP handler pattern
Every handler uses the deferred response pattern:

```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
    var resp interface{}
    defer func() {
        WriteResponse(w, resp)
    }()
    // ... set resp, call w.WriteHeader() for errors
}
```

`WriteResponse` JSON-encodes `resp` to the writer. Set `w.WriteHeader(statusCode)` before returning for non-200 responses.

### Route middleware
- `AuthMiddleware` — validates session cookie; redirects to `/login` if missing
- `ServerOffMiddleware` — returns `HTTP 423` if Factorio server is currently running

Routes marked `ServerOff: true` in the route table are automatically wrapped with `ServerOffMiddleware`.

### Authentication
- Passwords are bcrypt-hashed then base64-encoded before storage in SQLite
- Sessions use gorilla/sessions cookie store with the `CookieEncryptionKey` from config
- The `Secure` config flag controls the `Secure` attribute on session cookies (set `false` for plain HTTP)

### WebSocket hub
`websocket.WebsocketHub` is a package-level singleton initialized in `init()`. Use `WebsocketHub.GetRoom(name)` to obtain or create a room, then `room.Send(message)` to broadcast to subscribers. The `gamelog` room has special caching behavior (keeps the last `ConsoleCacheSize` lines).

### Version type
`factorio.Version` is `[4]uint` (major, minor, patch, build). Use `v.Greater(other)`, `v.Equals(other)` for comparisons. Implements `encoding.TextMarshaller`/`TextUnmarshaller` and `encoding.BinaryReader` for save file parsing.

---

## Frontend Conventions (React)

### Authentication and routing
`App.jsx` holds `isAuthenticated` state. The `ProtectedRoute` component redirects unauthenticated users to `/login`. After login, `serverStatus` is fetched and subscribed via WebSocket, then passed down to all views as a prop.

### API layer
- `ui/api/client.js` — Axios instance; all requests use `withCredentials: true`. A response interceptor shows a flash notification for 502 errors and silently rejects 401s.
- `ui/api/resources/*.js` — thin wrappers that call `client.get/post` and return data. Import these in views rather than using `client` directly.
- `ui/api/socket.js` — Exports an EventEmitter `bus`. Use `socket.emit('log subscribe')`, `socket.on('gamelog', handler)` etc. The WebSocket auto-reconnects every 5 seconds on close.

### Flash notifications
`window.flash(message, color)` is available globally (defined in `ui/notifications.js`, rendered by `<Flash />`). The Axios interceptor calls this automatically for network errors.

### Styling
- Tailwind CSS with a **custom color palette** — do not use default Tailwind color names. Use the project's colors:
  - `gray-dark`, `gray-medium`, `gray-light`
  - `white`, `dirty-white`
  - `green`, `green-light`, `blue`, `blue-light`
  - `red`, `red-light`, `orange`, `black`
- Additional Tailwind utilities: extended `width`/`margin` values (`72`, `80`, `88`, `96`)
- Global SCSS is in `ui/index.scss`

### Component library
Reusable primitives live in `ui/App/components/`. Prefer these over raw HTML:
- `<Button>`, `<ButtonLink>` — styled action buttons
- `<Input>`, `<InputPassword>`, `<Select>`, `<Checkbox>` — form controls
- `<Modal>`, `<ConfirmDialog>` — dialogs
- `<Panel>` — content card/container
- `<Label>` — form label
- `<Tabs>` — tabbed interface

---

## Docker Deployment

```bash
# Simple (HTTP, no Traefik):
docker compose -f docker/docker-compose.simple.yaml up -d

# Full stack with Traefik (HTTPS):
DOMAIN_NAME=fsm.example.com EMAIL_ADDRESS=you@example.com \
  docker compose -f docker/docker-compose.yaml up -d
```

The production image pulls a release zip from GitHub. For local builds, use `Dockerfile-local` or `Dockerfile-build`.

Volumes:
- `/opt/fsm-data` — FSM config, SQLite database, credentials
- `/opt/factorio/saves` — save files
- `/opt/factorio/mods` — installed mods
- `/opt/factorio/config` — Factorio config directory
- `/opt/fsm/mod_packs` — mod pack definitions

Ports:
- `80/tcp` (or `443/tcp` with Traefik) — web UI
- `34197/udp` — Factorio game server

---

## Contributing

- Branch from `develop`; open PRs targeting `develop`
- Update `CHANGELOG.md` with a human-readable entry for each change
- Go: follow standard Go formatting (`gofmt`); no lint tooling is currently configured
- Frontend: no linter configured; follow existing patterns
