# pg-cli

This is a CLI used to manage PostgreSQL databases for web applications.

<!-- vim-markdown-toc GFM -->

* [Why was this built?](#why-was-this-built)
* [Prerequisites](#prerequisites)
  * [Configuration](#configuration)
* [Installation](#installation)
* [Usage](#usage)
  * [Migrations](#migrations)
* [Development](#development)

<!-- vim-markdown-toc -->

## Why was this built?

In the past I've found myself relying on Makefiles to manage the PostgreSQL
databases associated with applications I've written. In Golang, unfortunately,
this has always felt like a bit of a kludge: at first these files would combine
the use of [Go and shell
script](https://github.com/kohrVid/calendar-api/blob/efdb530bd7a395134ad94b5e07cb2e97cccee1ab/Makefile)
but when this was refactored to only use Go, I found that I would have to rely
on [third-party packages to perform database migrations and application code
for everything
else](https://github.com/kohrVid/calendar-api/blob/8f116b4b5ed4fb5f866538c3f1a90d7bc77c276a/Makefile).

Though moving this functionality into a single package doesn't completely
remove the need for the Makefiles I've written, I think it would at least save
me having to re-write some of the application code for the database operations
in the work I'm planning.

This is still quite experimental but if others find it useful or have any ideas
for improvement then great!


##  Prerequisites

* Go v1.21+
* PostgreSQL

Optional for development:

* [gocov](https://github.com/axw/gocov#installation) (`go install github.com/axw/gocov/gocov@v1.0.0`)


### Configuration

This package is somewhat prescriptive in that your app must contain a YAML file
with configuration variables for each of the environments that you intend to
support. Environment names must be specified at the top level and each must
contain (or inherit) the following variables:

* `database_host` - the host name in the database URL
* `database_port` - port used to connect to the PostgreSQL database
* `database_user` - username used to connect to the database
* `database_name` - the name of the database

Optionally, an environment can also contain the following variables:

* `ssl_mode` - used to determine whether the database connection is made over
  SSL. By default this is set to `disabled` in pg-cli
* `data` - used to specify any records that should be seeded into the database

An example of this configuration can be found
[here](https://github.com/kohrVid/pg-cli/blob/master/example/env.yaml).
Note, because the CLI depends on
[viper](https://github.com/spf13/viper/issues/260) to parse YAML files, the
keys used in the configuration file are **NOT** case sensitive.

If your database requires a password, you can connect to it by assigning it to
an environment variable called `$DATABASE_PASSWORD` in your shell. Please note,
that database passwords should NOT be added to the configuration YAML as the
pg-cli will only look for the aforementioned environment variable. As this is
variable isn't parsed by the viper package, its name is **case sensitive** and
must be typed in all caps.


## Installation

To install this package, run:

    go install github.com/kohrVid/pg-cli@latest


## Usage

Create a new database:

    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT create

<sub>Note: if the `--config` flag isn't set, the example configuration in this
repo will be used to connect to Postgres.</sub>

Drop an existing database:

    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT drop

Seed the database with the records specified in the `ENVIRONMENT`.`data` value
of the configuration file:

    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT seed

Delete all rows in a database but maintain schema:

    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT clean


### Migrations

Migrations are run using files stored in the migration path. By default, the
pg-cli assumes that this `./migrations` but this can be set using the `--path`
(or `-p`) flag.

e.g.,

    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT migrate up --path "./custom/migrations/path"

As migrations depend on the [golang-migrate
package](https://github.com/golang-migrate), filenames should use the following
format:

    VERSION_migration_name.up.sql
    VERSION_migration_name.down.sql

The first migration file must generate a table (as opposed to an SQL function,
for example) with at least one field. Without this, the migration command will
fail silently. If you would like to create an SQL function in your first
migration an example work-around can be found
[here](https://github.com/kohrVid/pg-cli/blob/master/example/migrations/1_initialise_schema.up.sql).
Further details on how to write migration files can be found
on the [golang-migrate github repo](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md).

Apply all up migrations:

    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT migrate up

Apply all down migrations:

    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT migrate down

Step through migrations:

    # To go up a number of migrations, use a positive integer:
    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT migrate step -n 1

    # To go down a number of migrations, use a negative integer:
    pg-cli --config "./path/to/env.yaml" -e ENVIRONMENT migrate step -n -1


## Development

Clone the repo:

    git clone https://github.com/kohrVid/pg-cli.git
    cd pg-cli

Install dependencies:

    go mod vendor

Run tests:

    go test -v -count=1 ./...

To run just the tests for the `db` package:

    go test -v -count=1 ./db

To check coverage, you're advised to install gocov as per the [prerequisites
section](#prerequisites). To check the application's test coverage, run:

    gocov test ./... | gocov report
