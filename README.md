# Steam Deck Stock Alerts

Monitor Steam Deck inventory and send ntfy notifications when stock changes.
BoltDB retains state across restarts to avoid repeated stock notifications.

## Run

Copy [.env.example](.env.example) to `.env` and set `NTFY_TOKEN`.
Edit [config.yaml](config.yaml) for packages, polling interval, country, and a reachable ntfy server.
The Compose file does not start ntfy.

The [Compose stack](docker-compose.yml) requires the external network `caddy-network`.
If it is absent, create it once with `docker network create caddy-network`.
Then start the service:

```sh
docker compose up -d
docker compose logs -f stock-checker
```

The stack builds the application locally and retains database and log volumes.
Keep `.env` out of Git.
Use [config.yaml](config.yaml) as the complete configuration reference and the [Makefile](Makefile) for deployment commands.
