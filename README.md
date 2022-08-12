# demo-grpc
Playground for grpc/protobuf

* script to generate test certificates for TLS : makeCert.sh, to be called from certif directory
* multiple services sharing the same connection,
* a service including triggering server stop,
* various authentication scheme with TLS ( none, server only, server and client ...) : unse -unsafe flag on cleint and/or server to test
* end2end test with : got test -v .

* (wip) autogenerate REST-to-grppc reverse proxy gateway - see // see https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#httprule

