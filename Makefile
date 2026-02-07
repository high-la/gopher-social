# 1. Load environment variables from .env file
include .env
export # This exports all variables from .env so they are available to shell commands

# 2. Configuration
MIGRATIONS_PATH = ./cmd/migrate/migrations

# 3. Targets
.PHONY: migrate/create
migrate/create:
	@echo "Creating migration: $(filter-out $@,$(MAKECMDGOALS))"
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate/up
migrate/up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(GOPHER_SOCIAL_DSN) up

.PHONY: migrate/down
migrate/down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(GOPHER_SOCIAL_DSN) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate/force
migrate/force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(GOPHER_SOCIAL_DSN) force $(filter-out $@,$(MAKECMDGOALS))

# 4. The "Catch-All" Target
# This is CRITICAL. It prevents Make from throwing an error like 
# "make: *** No rule to make target 'migration_name'.  Stop."
%:
	@:
