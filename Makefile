IMAGE = ghcr.io/rossgrat/steam-deck-stock-alerts
SERVER = potatoserver
REMOTE_DIR = ~/services/steam-deck-stock-alerts

.PHONY: build push deploy setup logs stop

build:
	docker build --platform linux/amd64 -t $(IMAGE):latest .

push: build
	docker push $(IMAGE):latest

deploy: push
	scp deploy/docker-compose.yml $(SERVER):$(REMOTE_DIR)/docker-compose.yml
	scp config.yaml $(SERVER):$(REMOTE_DIR)/config.yaml
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker pull $(IMAGE):latest && docker-compose up -d"

setup:
	ssh $(SERVER) "mkdir -p $(REMOTE_DIR)"
	scp deploy/docker-compose.yml $(SERVER):$(REMOTE_DIR)/docker-compose.yml
	scp config.yaml $(SERVER):$(REMOTE_DIR)/config.yaml
	@echo "Now SSH into $(SERVER) and create $(REMOTE_DIR)/.env with NTFY_TOKEN=<your_token>"

stop:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker-compose down"

logs:
	ssh $(SERVER) "cd $(REMOTE_DIR) && docker-compose logs -f"
