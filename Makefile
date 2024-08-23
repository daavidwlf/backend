build: 
	docker build

start:
	docker compose up --build --detach

run logs:
	docker compose up --build

stop:
	docker compose down

list:
	docker ps