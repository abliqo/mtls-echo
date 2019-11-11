#!/bin/bash
#
# Creates certificates for testing
#
DEST_DIR=$(pwd)/certs
DAYS=3650

mkdir -p $DEST_DIR

echo "Creating root CA"
openssl req -new \
    -newkey rsa:2048 -nodes -keyout $DEST_DIR/rootca.key \
    -x509 -sha256 -days $DAYS \
    -subj "/C=US/ST=teststate/L=testloc/O=testorg/OU=testunit/CN=rootCA" \
    -out $DEST_DIR/rootca.pem

echo "Creating server cert"
openssl req -new \
    -newkey rsa:2048 -nodes -keyout $DEST_DIR/server.key \
    -sha256 -days $DAYS \
    -subj "/C=US/ST=teststate/L=testloc/O=testorg/OU=testunit/CN=localhost" \
    -out $DEST_DIR/server.csr

openssl x509 -req \
    -in $DEST_DIR/server.csr \
    -CA $DEST_DIR/rootca.pem \
    -CAkey $DEST_DIR/rootca.key \
    -CAcreateserial \
    -sha256 \
    -extfile server.ext \
    -out $DEST_DIR/server.pem

echo "Creating client cert"
openssl req -new \
    -newkey rsa:2048 -nodes -keyout $DEST_DIR/client.key \
    -sha256 -days $DAYS \
    -subj "/C=US/ST=teststate/L=testloc/O=testorg/OU=testunit/CN=localhost" \
    -out $DEST_DIR/client.csr

openssl x509 -req \
    -in $DEST_DIR/client.csr \
    -CA $DEST_DIR/rootca.pem \
    -CAkey $DEST_DIR/rootca.key \
    -CAcreateserial \
    -sha256 \
    -extfile client.ext \
    -out $DEST_DIR/client.pem

rm -f \
    $DEST_DIR/*.csr \
    $DEST_DIR/*.srl

echo "Done"