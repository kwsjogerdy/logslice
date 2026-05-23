# logslice

A fast log filtering tool that supports structured and unstructured log formats.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

Filter logs by level, keyword, or time range:

```bash
# Filter by log level
logslice --level error app.log

# Filter by keyword
logslice --match "connection refused" app.log

# Filter structured JSON logs by field value
logslice --field status=500 app.log

# Pipe from stdin
tail -f app.log | logslice --level warn
```

### Flags

| Flag | Description |
|------|-------------|
| `--level` | Filter by log level (debug, info, warn, error) |
| `--match` | Filter lines containing a keyword or pattern |
| `--field` | Filter structured logs by key=value pair |
| `--since` | Show logs after a given timestamp |
| `--until` | Show logs before a given timestamp |

---

## Features

- Supports structured (JSON) and unstructured (plaintext) log formats
- Reads from files or stdin for easy pipeline integration
- Fast line-by-line streaming with minimal memory overhead

---

## License

MIT © 2024 yourusername