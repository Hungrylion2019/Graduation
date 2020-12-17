.PHONY: all dev clean build env-up  run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@./start.sh
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@echo "Environment up"


##### RUN
run:
	@echo "Start cli ..."
	@docker-compose up -d cli

##### CLEAN
clean: 
	@echo "Stop environment ..."
	@echo "Environment down"
	@echo "Clean up ..."
	@./stop.sh
	@echo "Clean up done"
