export NS=always200
.PHONY: deploy
deploy:
	@oc new-project $(NS)
	@oc apply -f deploy/role.yaml
	@oc apply -f deploy/service_account.yaml
	@oc apply -f deploy/role_binding.yaml
	@oc apply -f deploy/crds/example.com_always200s_crd.yaml
	@oc apply -f deploy/crds/*_cr.yaml

.PHONY: delete
delete:

	@oc delete -f deploy/role.yaml
	@oc delete -f deploy/service_account.yaml
	@oc delete -f deploy/role_binding.yaml
	@oc delete -f deploy/crds/*_cr.yaml
	@oc delete -f deploy/crds/example.com_always200s_crd.yaml
	@oc delete project $(NS)

.PHONY: run
run:
	@operator-sdk run local --watch-namespace $(NS)

.PHONY: deploy/image
deploy/image:
	@oc new-project $(NS)
	@oc apply -f deploy/role.yaml
	@oc apply -f deploy/service_account.yaml
	@oc apply -f deploy/role_binding.yaml
	@oc apply -f deploy/crds/example.com_always200s_crd.yaml
	@oc apply -f deploy/operator.yaml