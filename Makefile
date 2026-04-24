.PHONY: frontend backend build dev-frontend dev-backend docker clean

frontend:
	cd src/frontend && pnpm install && pnpm build

backend:
	@mkdir -p bin
	cd src/backend && go build -ldflags "-s -w" -o ../bin/oversteplab ./cmd/server/

build: frontend backend

dev-frontend:
	cd src/frontend && pnpm dev

dev-backend:
	cd src/backend && go run ./cmd/server/main.go

docker:
	docker-compose up --build -d

clean:
	rm -rf bin/
	rm -rf src/backend/internal/web/dist/*
	cd src/frontend && rm -rf node_modules/
