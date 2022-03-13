app:
	docker build -f build/package/Dockerfile -t crypto-watcher:1.0.0 --build-arg BUILD_REF=1.0.0 .

KIND_CLUSTER := crypto-watcher-cluster

kind-up:
	kind create cluster \
		--name $(KIND_CLUSTER) \
		--config deployments/k8s/kind/kind-config.yml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide

kind-load:
	kind load docker-image crypto-watcher:1.0.0 --name $(KIND_CLUSTER)

kind-apply:
	cat deployments/k8s/base/base-service.yml | kubectl apply -f -

kind-logs:
	kubectl logs -l app=crypto-watcher-service -f --tail=100 --all-containers=true --namespace=crypto-watcher

kind-restart:
	kubectl rollout restart deployment crypto-watcher --namespace=crypto-watcher

debug:
	expvarmon --ports=12000

build:
	go build -o cmd/crypto-watcher/crypto-watcher cmd/crypto-watcher/main.go

run:
	go run cmd/crypto-watcher/main.go