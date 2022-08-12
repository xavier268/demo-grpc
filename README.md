# demo-grpc
Playground for grpc/protobuf. 

Utility scripts are provided :

* **install.sh** : install needed extensions/plugins to **protoc**. The **protoc** utility itself should be installed before, typically using *dnf -y install protobuf-complier*.
* **makeCert.sh** (call from certif directory) : script to generate various test certificates and keys for TLS. See the *service.conf* and *client.conf* for the respective service and client configurations.

Grpc demonstrated features include :

* multiple services sharing the same connection,
* a service including triggering server stop,
* various authentication scheme with TLS ( none, server only, server and client ...) : use -unsafe flag on client and/or server to test

* autogenerate REST-to-grppc reverse proxy gateway - see // see https://cloud.google.com/endpoints/docs/grpc-service-config/reference/rpc/google.api#httprule
    + the gateway accept REST connection and ats as a gateway with a grpc server, then retuning the grpc response in REST format.
    + gateway CANNOT handle non-unray request (see documentation).

An end-to-end test is provided :

* Run : *go test -v .*
