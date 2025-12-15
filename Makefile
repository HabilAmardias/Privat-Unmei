build:
	docker compose up --build -d
start:
	docker compose up database redis backend frontend cron proxy postgres-exporter loki promtail prometheus grafana -d
down:
	docker compose down
datadown:
	docker volume prune --force && \
	docker volume rm privat-unmei_privat-unmei-data && \
	docker volume rm privat-unmei_redis-data