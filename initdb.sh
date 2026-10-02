#!/bin/sh

DB_DIR="data"

if [ ! -d "$DB_DIR" ]; then \
    mkdir "$DB_DIR"; \
fi
if [ ! -f "$DB_DIR"/data.db ]; then \
    sqlite3 "$DB_DIR"/data.db "CREATE TABLE raw_messages (id INT8 PRIMARY KEY, user_id INT8 NOT NULL, msg_value TEXT, created TEXT);" \

    sqlite3 "$DB_DIR"/data.db "CREATE TABLE expenses (id INT8 PRIMARY KEY, exp_date TEXT, exp_value INT8, exp_description TEXT, created TEXT);"
fi

exec "$@"  


