bot:
	go build -o bin/expbot cmd/bot/*

service:
	go build -o bin/expservice cmd/service/*

all: bot service


initdb:
	if [ ! -d data ]; then \
		mkdir data; \
	fi
	if [ -f data/data.db ]; then \
		rm data/data.db; \
	fi
	sqlite3 data/data.db "CREATE TABLE raw_messages (id INT8 PRIMARY KEY, user_id INT8 NOT NULL, msg_value TEXT, created TEXT);"
	sqlite3 data/data.db "CREATE TABLE expenses (id INT8 PRIMARY KEY, exp_date TEXT, exp_value INT8, exp_description TEXT, created TEXT);"

build_clean: build initdb