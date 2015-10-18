all: setup migrate

fresh: clean setup migrate

setup:
	sqlite3 resalloc.db ""

migrate:
	cat migrations/1445043850_users_table_migration.sql | sqlite3 resalloc.db
	cat migrations/1445127756_resources_table_migration.sql | sqlite3 resalloc.db
	cat migrations/1445139240_machines_table_migration.sql | sqlite3 resalloc.db
	cat migrations/1445143449_leases_table_migration.sql | sqlite3 resalloc.db

clean:
	rm resalloc.db || true

test:
	docker-compose build
	docker-compose up -d
	docker-compose logs resalloc
