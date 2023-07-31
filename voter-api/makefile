SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <:TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "     build              Build the voter-api executable"
	@echo "     run                Run the voter-api program from code"
	@echo "     run-bin            Run the voter-api executable"
	@echo "     get-all            Get all voters with all voter history"
	@echo "     get-voter          Get a voter by passing id=<:voterId> on command line"
	@echo "     add-voter          Add a voter record by passing id=<:voterId>, firstName="<:firstName>", & lastName="<:lastName>" on command line"
	@echo "     update-voter       Update a voter record by passing id=<:voterId>, firstName="<:firstName>", & lastName="<:lastName>" on command line"
	@echo "     delete-all         Delete all voters"
	@echo "     delete-voter       Delete a voter by passing id=<:voterId> on command line"
	@echo "     get-voter-history  Get voter history for the voter with id=<:voterId>"
	@echo "     get-voter-poll     Get single voter poll data with id=<:voterId> and pollid=<:pollId>"
	@echo "     add-voter-poll     Add a voter poll record for the voter with id=<:voterId>"
	@echo "     update-voter-poll  Update a voter poll record with id=<:voterId> and pollid=<:pollId>"
	@echo "     delete-voter-poll  Delete a voter poll record with id=<:voterId> and pollid=<:pollId>"
	@echo "     health-check       Get the health status of the voter API"

.PHONY: build
build:
	go build .

.PHONY: run
run:
	go run main.go

.PHONY: run-bin
run-bin:
	./voter-api.exe

.PHONY: get-all
get-all:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters

.PHONY: get-voter
get-voter:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)

.PHONY: add-voter
add-voter:
	curl -w "\nHTTP Status: %{http_code}\n" -d '{ "firstName": "$(firstName)", "lastName": "$(lastName)" }' -H "Content-Type: application/json" -X POST http://localhost:1080/voters/$(id)

.PHONY: update-voter
update-voter:
	curl -w "\nHTTP Status: %{http_code}\n" -d '{ "firstName": "$(firstName)", "lastName": "$(lastName)" }' -H "Content-Type: application/json" -X PUT http://localhost:1080/voters/$(id)

.PHONY: delete-all
delete-all:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters

.PHONY: delete-voter
delete-voter:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters/$(id)

.PHONY: get-voter-history
get-voter-history:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)/polls

.PHONY: get-voter-poll
get-voter-poll:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)/polls/$(pollid)

.PHONY: add-voter-poll
add-voter-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:1080/voters/$(id)/polls/$(pollid)

.PHONY: update-voter-poll
update-voter-poll:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X PUT http://localhost:1080/voters/$(id)/polls/$(pollid)

.PHONY: delete-voter-poll
delete-voter-poll:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters/$(id)/polls/$(pollid)

.PHONY: health-check
health-check:
	curl -w "\nHTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/health
