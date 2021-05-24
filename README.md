# Always200

## Requirements
operator-sdk : v0.19.0

# Usage
A node express API used as an example for deploying via the operator-sdk on openshift
```bash
curl http://localhost:8080/get
ok
curl -d "{data.json"} -X POST http://localhost:8080/post
ok
```
The operator creates the 
- Deployment (will scale up or down to the spec.size specified in the always200 CR)
- Service 
- Route

## Deployment
Run locally
```bash
# deploy the resources
make deploy
# run the operator localy 
make run
# delete the deployment
make delete
```
Deploy on cluster
```bash
# deploy the resources
make deploy/image
# delete the deployment
make delete
```

Should work with any REST server provided its container is `EXPOSED 8080`
Change the image in the CR [here](https://github.com/austincunningham/always200/blob/master/deploy/crds/example.com_v1alpha1_always200_cr.yaml#L8)

## Container Images
- always200 rest container [![Docker Repository on Quay](https://quay.io/repository/austincunningham/always200/status "Docker Repository on Quay")](https://quay.io/repository/austincunningham/always200)
- always200 operator containter [![Docker Repository on Quay](https://quay.io/repository/austincunningham/always200-operator/status "Docker Repository on Quay")](https://quay.io/repository/austincunningham/always200-operator)

## TODO
- Write some unit tests