.PHONY: all
all: epithet-oidc-policy-lambda.zip				## run tests and build binaries

epithet-oidc-policy-lambda:					## build linux binary for lambda
	GOOS=linux GOARCH=amd64 go build

epithet-oidc-policy-lambda.zip: epithet-oidc-policy-lambda		## build lambda zip
	zip -X epithet-oidc-policy-lambda.zip epithet-oidc-policy-lambda

.PHONY: clean
clean:							## clean all local resources
	go clean
	rm -f epithet-*

.PHONY: help
help:							## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
