# 1. Load environment variables from .env file
include .env
export # This exports all variables from .env so they are available to shell commands

# 2. Configuration
MIGRATIONS_PATH = ./cmd/migrate/migrations
SEED_PATH = ./cmd/migrate/seed/main.go

# 3. Targets

# ___________________________________________________________________________________________
#
# Run and Build section 
# ___________________________________________________________________________________________
# 

.PHONY: run
run: gen/docs
	@echo "Running application"
	@go run ./cmd/api
	

# ___________________________________________________________________________________________
#
# Migration section
# ___________________________________________________________________________________________
# 

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

#  
# lets say version 5(000005) failed, 
# Force the version back to 4:
# tell the database to act as if version 4 was the last successful one:
# make migrate/force 4

.PHONY: migrate/force
migrate/force:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(GOPHER_SOCIAL_DSN) force $(filter-out $@,$(MAKECMDGOALS))


# ___________________________________________________________________________________________
#
# Seed section
# ___________________________________________________________________________________________
# 
.PHONY: seed
seed:
	@go run $(SEED_PATH)


# ___________________________________________________________________________________________
#
# Swagger Section
# ___________________________________________________________________________________________
# 

.PHONY: gen/docs
gen/docs:
	@echo "Generating swagger docs...."
	@swag init -g ./api/main.go --parseDependency -d cmd,internal && swag fmt





    


# ___________________________________________________________________________________________
#
# Catch-All section
# ___________________________________________________________________________________________
# 

# docker command for running redis
# docker run -d --rm --name gophersocial-redis -p 6379:6379 redis:8.4 redis-server --loglevel warning

# 4. The "Catch-All" Target
# This is CRITICAL. It prevents Make from throwing an error like 
# "make: *** No rule to make target 'migration_name'.  Stop."
%:
	@:
