# Variables
CC = go run
MODE = -mode=debug
FRONTEND_DIR = frontend
BACKEND_DIR = backend

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

# dep-web: 
# 	cd backend && go build  -o ../builds/web ./cmd/web
# 	scp -P 27322 builds/web git@dc.chimege.com:/home/docscan/__web
# 	scp -P 27322 ./backend/confs/web.yaml git@dc.chimege.com:/home/docscan/web.yaml
# 	ssh -p 27322 -tt git@dc.chimege.com "cd /home/docscan && mv __web web"
# 	ssh -p 27322 -tt git@dc.chimege.com "supervisorctl restart docscan"
# 	rm -rf builds/web

# dep-ui: 
# 	cd $(FRONTEND_DIR) && yarn generate
# 	cd $(FRONTEND_DIR)/dist && zip dist.zip -r *
# 	mv $(FRONTEND_DIR)/dist/dist.zip .
# 	scp -P 27322 dist.zip git@dc.chimege.com:/var/www/html
# 	ssh -p 27322 -tt git@dc.chimege.com "cd /var/www/html && unzip -o dist.zip && rm -rf dist.zip"
# 	rm dist.zip
