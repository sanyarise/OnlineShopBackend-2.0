password := 123
username := shopteam
db_name := shop
init_dir := $(CURDIR)/storage/init_db

start_db:
	(docker stop postgres && docker rm postgres) || true 
	docker run \
	--name postgres \
	-e POSTGRESQL_USERNAME=$(username) \
	-e POSTGRESQL_DATABASE=$(db_name) \
	-e POSTGRESQL_PASSWORD=$(password) \
	-p 5432:5432 \
	-v $(init_dir):/docker-entrypoint-initdb.d \
	 bitnami/postgresql:latest