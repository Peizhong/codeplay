.PHONY: proto

proto:
	protoc -I=./rpc/ \
		--go_out=./rpc/ --go_opt=paths=source_relative \
    	--go-grpc_out=./rpc/ --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./rpc/ --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out ./gen/openapiv2 \
        rpc/evaluator/evaluator.proto

	# go install go.uber.org/mock/mockgen@latest
	mockgen github.com/peizhong/codeplay/rpc/evaluator EvaluatorClient > ./gen/mock_evaluator/mock_evaluator.go

swag:
	swag init -d . --parseDependency --output ./gen/swagger/webapi

docker_build:
	sudo docker build -t 10.10.10.1:5000/codeplay:v0.0.1 .
	# sudo docker push 10.10.10.1:5000/codeplay:v0.1.4

docker_run:
	sudo docker run --rm -it 10.10.10.1:5000/codeplay:v0.0.1 sh

local_run:
	CODEPLAY_FEATURE_GATES="{\"enable_gops\":true,\"enable_pprof\":true}" go run main.go web

sync_k8s:
	scp ./kubernetes/codeplay/*.yaml peizhong@10.10.10.1:~/source/repos/codeplay/kubernetes/codeplay/

sync_docker:
	scp ./docker/kafka/* peizhong@10.10.10.1:~/source/repos/codeplay/docker/kafka

fast_build:
	- sudo docker rmi registry.cn-shenzhen.aliyuncs.com/peizhong/codeplay:v0.0.2
	sudo docker build -f fast.Dockerfile -t registry.cn-shenzhen.aliyuncs.com/peizhong/codeplay:v0.0.3 .
	sudo docker push registry.cn-shenzhen.aliyuncs.com/peizhong/codeplay:v0.0.3