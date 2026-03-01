# Steam Deck Stock Alerts

A Go service that monitors Steam Deck OLED inventory and sends push notifications via [ntfy](https://ntfy.sh) when stock status changes.

## Features

- Polls the Steam API at a configurable interval for multiple Steam Deck variants
- Sends notifications only on state transitions (in-stock / out-of-stock) to avoid spam
- Tracks state across restarts with an embedded BoltDB database
- API health monitoring with error/recovery alerts
- Structured JSON logging with rotation
- Multi-stage Docker build, deployed via Docker Compose

## Quick Start

1. Copy `.env.example` to `.env` and set your `NTFY_TOKEN`.
2. Edit `config.yaml` to configure packages, polling interval, and ntfy settings.
3. Run with Docker Compose:

```sh
docker compose up -d
```

## Configuration

**`config.yaml`**

| Key                | Description                        | Default                  |
| ------------------ | ---------------------------------- | ------------------------ |
| `polling_interval` | How often to check stock           | `1m`                     |
| `ntfy.url`         | ntfy server URL                    | `http://ntfy:80`         |
| `ntfy.topic`       | ntfy topic name                    | `steam-deck-alerts`      |
| `packages`         | List of `{id, name}` to monitor    | 512GB OLED, 1TB OLED     |
| `country_code`     | Region for the Steam API           | `US`                     |
| `log.path`         | Log file path                      | `/var/log/steam-deck-stock-alerts/app.log` |
| `db.path`          | BoltDB file path                   | `/data/stock.db`         |

**Environment variables:** `NTFY_TOKEN` (bearer token for ntfy auth).

## Deployment

The `Makefile` provides targets for deploying to a remote server:

```sh
make deploy   # Push config and pull latest image on the remote
make stop     # Stop the service
make logs     # Tail remote logs
```

## CI/CD

Pushes to `main` trigger a GitHub Actions workflow that:

1. Determines the next semantic version via [svu](https://github.com/caarlos0/svu)
2. Tags a release
3. Builds and pushes a Docker image to `ghcr.io/rossgrat/steam-deck-stock-alerts`
