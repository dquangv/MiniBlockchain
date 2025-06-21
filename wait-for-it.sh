#!/bin/sh
# wait-for-it.sh from https://github.com/vishnubob/wait-for-it
# (rút gọn, đủ xài cho TCP check)
host="$1"
shift
port="$1"
shift
cmd="$@"

until nc -z "$host" "$port"; do
  echo "⏳ Waiting for $host:$port..."
  sleep 1
done

exec $cmd
