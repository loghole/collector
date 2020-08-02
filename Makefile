run-alpha: db-image collector-image dashboard-image
	docker-compose -f docker-compose-alpha.yaml up clickhouse-db collector dashboard

db-image:
	cd clickhouse-db && $(MAKE) docker-image

collector-image:
	cd collector && $(MAKE) docker-image

dashboard-image:
	cd dashboard && $(MAKE) docker-image

go-lint:
	cd dashboard/backend && $(MAKE) go-lint
	cd collector && $(MAKE) go-lint
