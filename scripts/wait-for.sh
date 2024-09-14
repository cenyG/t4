#!/bin/sh

host="$1"
shift
cmd="$@"

until nc -z "$host"; do
  >&2 echo "Waiting for $host..."
  sleep 1
done

>&2 echo "$host is available - executing command"
exec $cmd