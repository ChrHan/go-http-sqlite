# SQLite Go Wrapper

A simple `Go` program to execute `CRUD` to `sqlite3` database.

Currently programmed only for this table specification:

    CREATE TABLE PRODUCTS (
      id  INT PRIMARY KEY,
      product_name  VARCHAR(20)
    )

# Prerequisite

Install `Go`

# Usage

1. Clone this repo
1. cd to `app/sqlite_wrapper`
1. Run `go install`
1. Put Environment Variables for custom sqlite3:

    export APP_FILENAME="{new_sqlite3_db_name}"
    export APP_LOG_LEVEL="{debug|info|warn|error|fatal|panic}"

1. Run program using `sqlite-wrapper {select|insert id product_name|update id product_name|delete id|delete_all}`

