#!/usr/bin/sh

docker image ls | grep -P '(vault|form3tech|postgres)' | awk '{ print $3 }' | xargs docker image rmi -f 2>/dev/null
docker volume prune -f 2>/dev/null