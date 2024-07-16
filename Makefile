test:
	bash -c "set -m; bash '$(CURDIR)/scripts/test.sh'"

run:
	bash -c "set -m; bash '$(CURDIR)/scripts/run.sh'"

db-console:
	docker exec -it uservice-subscriptions-postgres-subscriptions-1 \
		bash -c "PGPASSWORD=postgres psql -U postgres -d postgres"

PHONY: test run
