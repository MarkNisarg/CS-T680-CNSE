## ToDo CLI App

The ToDo App is a command-line interface (CLI) tool to manage a list of todo items. It uses a text file-based database to store and retrieve todo items. The default database file is `./data/todo.json`, but you can specify a different database file using the `-db` flag.

## Usage

To use the ToDo App, you can run the `todo` executable with various command-line options. Here are the available options:

```
Usage of ./todo:
  -a string
        Add an item to the database
  -d int
        Delete an item from the database
  -db string
        Name of the database file (default "./data/todo.json")
  -l    List all the items in the database
  -q int
        Query an item in the database
  -s    Change item 'done' status to true or false
  -u string
        Update an item in the database
```

### List all items

To list all items in the database, use the `-l` flag:

```
./todo -l
```

To list all items in the database, using the `make` command:
```
make list
```

Both commands will display all the items stored in the database.

### Query an item

To query a specific item by its ID, use the `-q` flag followed by the item ID:

```
./todo -q 2
```

To query a specific item by its ID, using the `make` command:
```
make query id=2
```

Both commands will retrieve and display the item with the specified ID from the database.

### Add an item

To add a new item to the database, use the `-a` flag followed by the item details in JSON format:

```
./todo -a '{"id":100, "title":"New item", "done":false}'
```

To add a new item to the database, using the `make` command:
```
make add '{"id":100, "title":"sample item", "done":false}'
```

Both commands will add a new item with the specified ID, title, and done status to the database.

### Update an item

To update an existing item in the database, use the `-u` flag followed by the item details in JSON format:

```
./todo -u '{"id":100, "title":"New item", "done":false}'
```

To update an existing item in the database, using the `make` command:
```
make update '{"id":100, "title":"sample item", "done":false}'
```

Both commands will update the item with the specified ID, modifying its title and done status in the database.

### Delete an item

To delete an item from the database, use the `-d` flag followed by the item ID:

```
./todo -d 2
```

To delete an item from the database, using the `make` command:

```
make delete id=2
```

Both commands will remove the item with the specified ID from the database.

### Change item status

To change the done status of an item in the database, use the `-q` followed by the item ID and `-s` flag followed by the new status (true or false):

```
./todo -q 2 -s=true
```

To change the done status of an item in the database, using the `make` command:

```
make change-status id=2 done=true
```

This command will update the done status of the item with the specified ID in the database.

### Makefile Commands

The provided Makefile includes several targets to automate common commands. Here are the available targets:

- `build`: Build the `todo` executable.
- `run`: Run the `todo` program from code.
- `run-bin`: Run the `todo` executable.
- `restore-db`: Restore the sample database (Unix/Mac).
- `restore-db-windows`: Restore the sample database (Windows).
- `add-sample`: Add a sample row to the database.
- `list`: List all items in the database.
- `query`: Query an item by ID.
- `add`: Add a new item to the database.
- `update`: Update an existing item in the database.
- `delete`: Delete an item from the database.
- `change-status`: Change the done status of an item in the database.

You can use these targets with the `make` command to execute the corresponding actions. For example, `make list` will list all items in the database, and `make add '{"id":100, "title":"New item", "done":false}'` will add a new item to the database.
