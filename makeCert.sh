#!/bin/bash

echo "Making a selfsigned certificate for the server"

echo "Generate rsa key for CA"
openssl genrsa -out ca.key 4096

echo "Generate certificate signed with generated CA key"
openssl req -new -x509 -key ca.key -sha256 -subj "/C=FR/ST=PARIS/O=test organisation." -days 2 -out ca.cert