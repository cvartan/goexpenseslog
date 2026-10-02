FROM debian:bookworm-slim


RUN apt-get update && apt-get install -y \
    curl \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*



RUN mkdir -p /opt/goexpenseslog/bin

RUN mkdir -p /opt/goexpenseslog/data

RUN mkdir /opt/goexpenseslog/config

WORKDIR /opt/goexpenseslog

COPY /bin/expbot bin/

COPY /config/config.yml config/

COPY --chmod=+x initdb.sh .

ENTRYPOINT ["./initdb.sh"]

CMD ["bin/expbot","-config","config/config.yml"]

