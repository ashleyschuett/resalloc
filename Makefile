all: setup migrate

fresh: clean setup migrate

setup:
	sqlite3 resalloc.db ""

migrate:
	cat migrations/1445043850_users_table_migration.sql | sqlite3 resalloc.db

clean:
	rm resalloc.db

test:
	docker-compose build
	docker-compose up -d
