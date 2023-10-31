#!/bin/sh

docker exec -ti warranty-redis sh
redis-cli
#select 0
#keys *
#DEL *