# hook-replay

A CLI tool to intercept, capture, and replay webhook requests. Useful for local webhook development and debugging.

## About

This project was built for personal use and to practice Go. Many technical decisions — such as using the standard library’s flag package instead of a CLI framework, or adopting a simple layered architecture — were made intentionally for learning and experimentation rather than as "production recommendations".

## Installation

```sh
go install github.com/mathiasdonoso/hook-replay/cmd/hr@latest
```

## Quick Start

Start a proxy server that captures incoming webhooks and forwards them to your local app:

```sh
hr serve --port 3000 --forward localhost:8080
```

View the captured requests:

```sh
hr log
```

Replay the most recent request:

```sh
hr replay --last
```

Replay a specific request by ID, three times with a 500ms delay between each:

```sh
hr replay abc12345 --times 3 --delay 500
```

## Commands

### `hr serve`

Starts an HTTP proxy server that captures incoming webhook requests and forwards them to a target.

| Flag        | Type   | Default | Description                          |
|-------------|--------|---------|--------------------------------------|
| `--port`    | uint   | `3000`  | Port to listen on                    |
| `--forward` | string | —       | Target URL to forward requests to (required) |

```sh
hr serve --port 9000 --forward localhost:8080
```

---

### `hr log`

Displays the last 20 captured requests in a formatted table.

```sh
hr log
```

Output columns: `Id`, `Source`, `Path`, `Method`, `Body`, `ReceivedAt`

---

### `hr replay [id]`

Replays a previously captured request. Accepts the first characters of the event ID as shown in `hr log`.

| Flag       | Type   | Default | Description                                      |
|------------|--------|---------|--------------------------------------------------|
| `--last`   | bool   | `false` | Replay the most recent captured event            |
| `--times`  | uint   | `1`     | Number of times to replay the request            |
| `--delay`  | uint   | `0`     | Delay in milliseconds between each replay        |
| `--target` | string | —       | Override the original forward target             |

```sh
# Replay a specific event
hr replay abc12345

# Replay the last event to a different target
hr replay --last --target localhost:9000

# Replay an event 5 times with a 1 second delay
hr replay abc12345 --times 5 --delay 1000
```

## Data Storage

Captured events are stored in a local SQLite database at `~/.hook-replay/events.db`. The database and schema are created automatically on first run.
