#!/bin/bash

## NOT TO BE USED FOR REAL SECURITY


######### CA

echo "Generate rsa key for CA"
openssl genrsa -out ca.key 4096

echo "Generate self-signed certificate for CA"
openssl req -new -x509 -key ca.key -sha256 -subj "/C=FR/ST=PARIS/O=test organisation." -days 2 -out ca.cert


######### Service/Server 

echo "Generate a key certificate for server"
openssl genrsa -out service.key 4096

echo "Create a signing request for server" # for real, modify the service.conf file
openssl req -new -key service.key -out service.csr -config service.conf

echo "sign the server certificate" # for real, need to use create srl file (serial), not recreate each time.
openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out service.pem -days 2 -sha256 -extfile service.conf -extensions req_ext

echo "Verify server certificate"
openssl x509 -in service.pem -text -noout

######### Client 

echo "Generate rsa key for client"
openssl genrsa -out client.key 4096

echo "Create a signing request for client" # for real, modify the service.conf file
openssl req -new -key client.key -out client.csr -config client.conf 

echo "sign the client certificate" # for real, need to use create srl file (serial), not recreate each time.
openssl x509 -req -in client.csr -CA ca.cert -CAkey ca.key -CAcreateserial -out client.pem -days 2 -sha256 -extfile client.conf -extensions v3_req

echo "Verify client certificate"
openssl x509 -in client.pem -text -noout


