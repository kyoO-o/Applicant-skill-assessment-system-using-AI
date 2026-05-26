# Variables
CC = go run
MODE = -mode=debug
FRONTEND_DIR = frontend
BACKEND_DIR = backend
SERVER = root@167.172.65.141
BUILDS_DIR = builds

# Determine the processor type
PROCESSOR_TYPE := $(shell uname -m)

web:
ifeq ($(PROCESSOR_TYPE), arm64)
	cd $(BACKEND_DIR)/cmd && $(CC) -tags=dynamic ./web $(MODE)
else
	cd $(BACKEND_DIR)/cmd && $(CC) ./web $(MODE)
endif

ui: 
	cd $(FRONTEND_DIR) && pnpm dev --dotenv ./env/.env --host

dep-web:
	mkdir -p $(BUILDS_DIR)
	cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o ../$(BUILDS_DIR)/api ./cmd/web
	scp $(BUILDS_DIR)/api $(SERVER):/home/deploy/api/__api
	rm -f $(BUILDS_DIR)/api
	ssh $(SERVER) "cd /home/deploy/api && mv api _api 2>/dev/null || true && mv __api api"
	ssh $(SERVER) "systemctl restart web-api"

dep-ui:
	cd $(FRONTEND_DIR) && NUXT_PUBLIC_API_BASE=https://skillz.works pnpm generate
	rsync -avz --delete $(FRONTEND_DIR)/.output/public/ $(SERVER):/var/www/html/
# 	cd frontend && NUXT_PUBLIC_API_BASE=http://167.172.65.141 pnpm generate
# 	rsync -avz --delete frontend/.output/public/ root@167.172.65.141:/var/www/html/
