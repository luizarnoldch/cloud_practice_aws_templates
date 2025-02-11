TEMPLATE_FILE := deploy.yml
STACK_NAME := go-lambda-practice

test:
		echo "test makefile"
build:
	./build_go.sh
run:
	./run_lambda.sh
deploy:
	sam deploy --template-file $(TEMPLATE_FILE) --stack-name $(STACK_NAME) --capabilities CAPABILITY_NAMED_IAM --resolve-s3
destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
