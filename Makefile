.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=darwin GOARCH=amd64 go build -o bin/admission-webhook-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o bin/admission-webhook-linux-amd64 .

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building go-webhook Docker image..."
	docker build -t go-webhook:1.0 . 

# From this point `kind` is required
.PHONY: cluster
cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	kind create cluster --config deploy/manifests/kind/kind.cluster.yaml

.PHONY: delete-cluster
delete-cluster:
	@echo "\n♻️  Deleting Kubernetes cluster..."
	kind delete cluster

.PHONY: push
push: docker-build
	@echo "\n📦 Pushing admission-webhook image into Kind's Docker daemon..."
	kind load docker-image go-webhook:1.0

.PHONY: deploy-config
deploy-config:
	@echo "\n⚙️  Applying cluster config..."
	kubectl apply -f deploy/manifests/cluster-config/

.PHONY: delete-config
delete-config:
	@echo "\n♻️  Deleting Kubernetes cluster config..."
	kubectl delete -f deploy/manifests/cluster-config/

.PHONY: deploy
deploy: push delete  deploy-config
	@echo "\n🚀 Deploying go-webhook..."
	kustomize build deploy/manifests/go-webhook/ | kubectl apply -f-


.PHONY: delete
delete:
	@echo "\n  Deleting go-webhook deployment if existing..."
	kubectl delete deploy go-webhook

.PHONY: logs
logs:
	@echo "\n🔍 Streaming simple-go-webhook logs..."
	sleep 8
	kubectl logs -l app=go-webhook -f

.PHONY: delete-all
delete-all: delete delete-config delete-pod delete-bad-pod
