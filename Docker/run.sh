#!/bin/bash
pwd

sed -i -e "s;%DB_USERNAME%;$DB_USER;g" ./.env
sed -i -e "s;%DB_PASSWORD%;$DB_PASSWORD;g" ./.env
sed -i -e "s;%DB_NAME%;$DB_DATABASE;g" ./.env
sed -i -e "s;%DB_HOST%;$DB_HOST;g" ./.env
sed -i -e "s;%DB_PORT%;$DB_PORT;g" ./.env
sed -i -e "s;%SERVER_PORT%;$SERVER_PORT;g" ./.env
sed -i -e "s;%REDIS_SERVER%;$REDIS_HOST;g" ./.env
sed -i -e "s;%REDIS_SERVER_PORT%;$REDIS_PORT;g" ./.env
sed -i -e "s;%REDIS_DB%;$REDIS_DB;g" ./.env

cat ./.env
/dist/main