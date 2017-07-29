#!/bin/bash

if [ -z "${LISTEN_PORT}" ]; then
	LISTEN_PORT=6379;
fi

if [ -z "${MASTER_NAME}" ]; then
	echo "Missing master name, assuming \"master\"";
	MASTER_NAME="master";
fi

if [ -z "${SENTINEL_ADDR}" ]; then
	SENTINEL_ADDR="redis-sentinel:26379";
fi

exec "$@ -listen :{$LISTEN_PORT} -sentinel ${SENTINEL_ADDR} -master ${MASTER_NAME}"