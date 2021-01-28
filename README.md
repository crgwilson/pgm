# PostgreSQL Schema Migrator

CLI tool to handle PostgreSQL database schema migrations using regular ol' SQL files.

Inspired by [goose](https://gitlab.com/blaskovicz/goose), and all the other
tools like this out there already.

## Quickstart

Create your SQL files using the following file naming scheme.

```console
<schema-version>.<up|down>.sql
```

...for example...

```console
001.up.sql
```

Initialize your database to work with the tool...

```console
pgm init
```

Run `pgm up` to apply all of your SQL files

```console
pgm up
```

## TODOs

* Use a pgpass file for connecting rather than command-line arguments
* Upgrade & downgrade to specific versions
* `list` subcommand to print out all found schema versions
* The CLI logger needs to be able to format strings properly
