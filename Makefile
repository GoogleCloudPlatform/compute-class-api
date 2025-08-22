.PHONY: gen-client-init
gen-clients-init:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.14.0

.PHONY: gen-client
gen-client:
	bash hack/gen-client.sh
