build:
	sudo docker compose up --build -d
start:
	sudo docker compose up database redis backend frontend cron proxy postgres-exporter loki promtail prometheus grafana -d
down:
	sudo docker compose down
datadown:
	sudo docker volume prune --force && \
	sudo docker volume rm privat-unmei_privat-unmei-data && \
	sudo docker volume rm privat-unmei_redis-data
