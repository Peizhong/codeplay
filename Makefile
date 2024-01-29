.PHONY: proto

proto: 
	protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        rpc/evaluator/evaluator.proto

	mockgen github.com/peizhong/codeplay/rpc/evaluator EvaluatorClient > ./gen/mock_evaluator/mock_evaluator.go

swag:
	swag init -d . --parseDependency --output ./gen/swagger/webapi


build_builder:
	sudo docker build -f builder.Dockerfile -t codeplay:v0.0.1-builder .

build_app:
	sudo docker build -f app.Dockerfile -t 10.10.10.1:5000/codeplay:v0.1.4 .
	sudo docker push 10.10.10.1:5000/codeplay:v0.1.4

sync_k8s:
	scp ./kubernetes/codeplay/*.yaml peizhong@10.10.10.1:~/source/repos/codeplay/kubernetes/codeplay/

sync_docker:
	scp ./docker/kafka/* peizhong@10.10.10.1:~/source/repos/codeplay/docker/kafka