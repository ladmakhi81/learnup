dev:
	@docker compose up --build -d
stop:
	@docker compose stop
update-doc:
	@./generate_doc.sh && make dev