# Always200
A node express API used as an example for deploying via the operator-sdk on openshift
```bash
curl http://localhost:8080/get
ok
curl -d "{data.json"} -X POST http://localhost:8080/post
ok
```
The operator creates the 
- Deployment (2 pods)
- Service 
- Route

## Usage
```bash
# deploy the resources
make deploy
# run the operator localy 
make run
# delete the deployment
make delete
```

Should work with any REST server provided its container is `EXPOSED 8080`
Change the image in the CR [here](https://github.com/austincunningham/always200/blob/master/deploy/crds/example.com_v1alpha1_always200_cr.yaml#L8)

## Container Images
[![Docker Repository on Quay](https://quay.io/repository/austincunningham/always200/status "Docker Repository on Quay")](https://quay.io/repository/austincunningham/always200)
[![Docker Repository on Quay](https://quay.io/repository/austincunningham/always200-operator/status "Docker Repository on Quay")](https://quay.io/repository/austincunningham/always200-operator)
