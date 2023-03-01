#!/bin/sh
echo $AWS_LAMBDA_RUNTIME_API

if [ -z "${AWS_LAMBDA_RUNTIME_API}" ]; then
  exec /rie "$@"
else
  exec "$@"
fi