#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

mkdir -p hack/_certs/

# CREATE THE PRIVATE KEY FOR OUR CUSTOM CA
openssl genrsa -out hack/_certs/ca.key 2048

# GENERATE A CA CERT WITH THE PRIVATE KEY, THIS IS THE CA WE WILL NEED.
openssl req -new -x509 -key hack/_certs/ca.key -out hack/_certs/ca.crt -config hack/ca_config.txt

# CREATE THE PRIVATE KEY FOR OUR OAM ADMISSION CONTROLLER
openssl genrsa -out hack/_certs/admission-key.pem 2048

# CREATE A CSR FROM THE CONFIGURATION FILE AND OUR PRIVATE KEY
openssl req -new -key hack/_certs/admission-key.pem -subj "/CN=admission.default.svc" -out hack/_certs/admission.csr -config hack/admission_config.txt

# CREATE THE CERT SIGNING BY THE CSR AND THE CA CREATED BEFORE FOR OUR OAM ADMISSION CONTROLLER
openssl x509 -req -in hack/_certs/admission.csr -CA hack/_certs/ca.crt -CAkey hack/_certs/ca.key -CAcreateserial -out hack/_certs/admission-crt.pem