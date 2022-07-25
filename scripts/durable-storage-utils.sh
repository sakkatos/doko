#!/bin/bash

init_psql_db_cluster() {
    if [ ! "$(ls -A $PGSQL_DATADIR)" ]; then
        echo -e "Initialising PostgreSQL DB Cluster..."
        rm -rf "${PGSQL_HOME}"/logfile
        "${PGSQL_BASEDIR}"/bin/pg_ctl -U "${PGSQL_ROOTUSER}" -D "${PGSQL_DATADIR}" initdb -o "-E=UTF8 --no-locale"
    fi
}

start_psql_db_cluster() {
    mkdir -p "${PGSQL_HOME}"
    mkdir -p "${PGSQL_DATADIR}"
    export PGSQLD_PID=$(ps -ax | grep -v " grep " | grep "${PGSQL_BASEDIR}"/bin/postgres | awk '{ print $1 }')
    if [ -z "${PGSQLD_PID}" ]; then
        init_psql_db_cluster
        echo -e "Starting up PostgreSQL DB Server..."
        "${PGSQL_BASEDIR}"/bin/pg_ctl -U "${PGSQL_ROOTUSER}" -D "${PGSQL_DATADIR}" -o "-p ${PGSQL_TCP_PORT} -k ${PGSQL_HOME}" -l "${PGSQL_HOME}"/logfile start
        export PGSQLD_PID=$!
    fi
}

create_doko_db_user() {
    DB_USER_NAME="${1}"
    if psql -h localhost -p "${PGSQL_TCP_PORT}" "${PGSQL_ROOTUSER}" -t -c '\du' | cut -d \| -f 1 | grep -qw "${DB_USER_NAME}"; then
        echo "PostgreSQL DB User ${DB_USER_NAME} already exists."
    else
        createuser -e -h localhost -p "${PGSQL_TCP_PORT}" -d -r -S "${DB_USER_NAME}"
    fi
}

create_doko_db() {
    DB_NAME="${1}"
    DB_OWNER_NAME="${2:-$PGSQL_ROOTUSER}"
    if psql -h localhost -p "${PGSQL_TCP_PORT}" "${PGSQL_ROOTUSER}" -lqt | cut -d \| -f 1 | grep -qw "${DB_NAME}"; then
        echo "PostgreSQL doko DB ${DB_NAME} already exists. Continuing..."
    else
        createdb -e -h localhost -p "${PGSQL_TCP_PORT}" -O "${DB_OWNER_NAME}" "${DB_NAME}"
    fi
}
