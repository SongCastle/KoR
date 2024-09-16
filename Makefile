.PHONY: back
.PHONY: front

# container
init:
	cp back/.env.sample back/.env
	cp db/.env.sample db/.env
	cp db/password.txt.sample db/password.txt
	cp front/.env.sample front/.env
	@make build
	@make up
	docker compose exec back ash -c "chmod +x serve.sh && ./serve.sh"

build:
	docker compose build --no-cache

up:
	docker compose up -d
	docker compose exec back ash -c "chmod +x serve.sh && ./serve.sh"

stop:
	docker compose stop

down:
	docker compose down

destroy:
	docker compose down --rmi all --volumes --remove-orphans

ps:
	docker compose ps

# backend
back:
	docker compose exec back ash

# frontend
front:
	docker compose exec front bash