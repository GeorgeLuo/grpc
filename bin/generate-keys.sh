#!/bin/sh
openssl req -newkey rsa:2048 \
  -new -nodes -x509 \
  -days 3650 \
  -out cert.pem \
  -keyout key.pem \
  -subj "/C=US/ST=New York/L=New York/O=GRPC Industries/OU=GRPC1/CN=localhost"
