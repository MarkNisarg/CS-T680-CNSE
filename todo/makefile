SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET> [ARGUMENTS]"
	@echo ""
	@echo "    Targets:"
	@echo "    build                                    Build the todo executable"
	@echo "    run                                      Run the todo program from code"
	@echo "    run-bin                                  Run the todo executable"
	@echo "    restore-db                               Restore the sample database (unix/mac)"
	@echo "    restore-db-windows                       Restore the sample database (windows)"
	@echo "    add-sample                               Add a sample row"
	@echo "    add item_json='<JSON GOES HERE>'         Add an item in the database"
	@echo "    list                                     List all items in the database"
	@echo "    query id=ITEM_ID                         Query an item in the database"
	@echo "    update item_json='<JSON GOES HERE>'      Update an item in the database"
	@echo "    delete id=ITEM_ID                        Delete an item from the database"
	@echo "    change-status id=ITEM_ID done=STATUS     Change the done status of an item in the database"


.PHONY: build
build:
	go build .

.PHONY: run
run:
	go run main.go

.PHONY: run-bin
run-bin:
	./todo.exe

.PHONY: restore-db
restore-db:
	cp ./data/todo.json.bak ./data/todo.json

.PHONY: restore-db-windows
restore-db-windows:
	copy .\data\todo.json.bak .\data\todo.json

.PHONY: add-sample
add-sample:
	./todo.exe -a '{"id":99, "title":"sample item", "done":true}'

.PHONY: add
add:
	./todo.exe -a '$(item_json)'

.PHONY: list
list:
	./todo.exe -l

.PHONY: query
query:
	./todo.exe -q $(id)

.PHONY: update
update:
	./todo.exe -u '$(item_json)'

.PHONY: delete
delete:
	./todo.exe -d $(id)

.PHONY: change-status
change-status:
	./todo.exe -q $(id) -s=$(done)
