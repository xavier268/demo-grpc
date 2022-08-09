#!/bin/bash

echo "Making a selfsigned certificate for the server"

echo "Generate rsa key for CA"
openssl genrsa -out ca.key 4096

echo "Generate self-signed certificate for CA"
openssl req -new -x509 -key ca.key -sha256 -subj "/C=FR/ST=PARIS/O=test organisation." -days 2 -out ca.cert

echo "Generate a key certificate for server"
openssl genrsa -out service.key 4096

echo "Create a signing request for server"
openssl req -new -key service.key -out service.csr -config certificate.conf

echo "sign the server certificate"
openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out service.pem -days 2 -sha256 -extfile certificate.conf -extensions req_ext

