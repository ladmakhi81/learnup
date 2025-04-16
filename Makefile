dev:
	@docker compose up --build -d && docker ps
stop:
	@docker compose stop
update-doc:
	@./generate_doc.sh && make dev