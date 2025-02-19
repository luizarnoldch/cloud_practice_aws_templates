SAM_TEMPLATE_FILE := deployment/sam.yml
STACK_NAME := go-lambda-dynamodb-practice

test:
		go test ./... -v
build:
	./build_go.sh
docker-up:
	docker-compose up -d
docker-down:
	docker-compose down
docker-down-all:
	docker-compose down -v
run:
	./run_lambda.sh
deploy:
	sam deploy --template-file $(TEMPLATE_FILE) --stack-name $(STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
