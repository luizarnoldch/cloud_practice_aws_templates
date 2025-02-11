.PHONY: test test-local build run upload-cf deploy deploy-cf destroy

SAM_TEMPLATE_FILE := deployment/sam.yml
CF_TEMPLATE_FILE := deployment/cloud_formation.yml
S3_BUCKET := documentos-pdf
S3_PREFIX := deployments
CF_S3_TEMPLATE_PATH := s3://$(S3_BUCKET)/$(S3_PREFIX)/cloud_formation.yml
LAMBDA_PACKAGE := go_function.zip
LAMBDA_S3_PATH := s3://$(S3_BUCKET)/$(S3_PREFIX)/$(LAMBDA_PACKAGE)
STACK_NAME := go-lambda-practice

test:
	go test -v ./...
test-local:
	sam local invoke HelloWorldFunction --template $(SAM_TEMPLATE_FILE) --docker-volume-basedir $(PWD)
build:
	./build_go.sh
run:
	./run_lambda.sh
upload-cf: build
	aws s3 cp $(CF_TEMPLATE_FILE) $(CF_S3_TEMPLATE_PATH)
	aws s3 cp $(LAMBDA_PACKAGE) $(LAMBDA_S3_PATH)
deploy:
	sam deploy \
		--template-file $(SAM_TEMPLATE_FILE) \
		--stack-name $(STACK_NAME) \
		--capabilities CAPABILITY_NAMED_IAM \
		--resolve-s3
deploy-cf: upload-cf
	aws cloudformation create-stack \
		--stack-name $(STACK_NAME) \
		--template-body file://$(CF_TEMPLATE_FILE) \
		--capabilities CAPABILITY_NAMED_IAM 
destroy:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
