# OAM Admission controller

Admission controller is used for mutating and validating OAM component, trait and application configuration.
It is integrated with Kubernetes by [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).

Admission controller can provide web hook API for mutating and validating OAM spec, not only used by Rudr,
it can be used by all OAM implementation for validation.

## Validation List

1. checking component/trait existence for AppConfig
2. check component/trait really have parameters or properties mentioned in AppConfig
3. invalidate if a worker tries to bind to a port
4. check component/trait/AppConfig format when created
5. check application scope format when created
6. check application scope instance when they created by AppConfig


## Mutation List 
  
1. resolve and mutate `fromVariable()` func in AppConfig

## Validation Not Included

1. **schema require or not**: this could validate in CRD definition yaml. 

## How to install

### Generate the certificates and the CA

Admission Controller need you to prepare certificates and ca, **for none-production use**,
you cloud generate it by the shell script.

```shell
$ ./hack/gen_certs.sh
```

The script will generate the certificates in `./hack/_certs/`

```output
$ tree hack/_certs/
hack/_certs/
├── admission-crt.pem
├── admission-key.pem
├── admission.csr
├── ca.crt
├── ca.key
└── ca.srl
```

### Replace Helm Values with CA and certificates

1. use ca created here fill with CaBundle in charts/values.yaml

```shell
cat charts/admission/values_template.yaml | sed 's/_CaBundle_/'"$(cat hack/_certs/ca.crt | base64 | tr -d '\n')"'/g' > charts/admission/values.yaml
```

2. create certificates for Admission Controller as K8s secrets

_Notice: change namespace using the same namespace as helm installed_

```shell
kubectl create secret generic admission-oam -n default \
        --from-file=key.pem=hack/_certs/admission-key.pem \
        --from-file=cert.pem=hack/_certs/admission-crt.pem
```

3. install charts by helm tool.

```shell
helm install admission ./charts/admission
```

## How to use it ?

It will automatically check Application Configuration and other OAM resources when they are created or modified.

## Developing

1. make dependency to you vendor files build building docker image

```
go mod vendor
``` 

2. make sure `code-generator` in your `vendor` folder.

```shell script
mkdir -p vendor/k8s.io
cd vendor/k8s.io
git clone git@github.com:kubernetes/code-generator.git
```

