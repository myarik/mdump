SHELL := /bin/bash # Use bash syntax

OS=linux
ARCH=amd64
APP_NAME=mdump

APP_VERSION?=0.0.1

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## build: builds an application
build:
	docker build -t ${APP_NAME} --build-arg APP_VERSION=${APP_VERSION} .
	docker tag ${APP_NAME}:latest ${APP_NAME}:${APP_VERSION}

#############################################################################
#
# If you require a different configuration from the defaults below, create a
# new file named "Makefile.local" in the same directory as this file and define
# your parameters there.
#############################################################################
-include Makefile.local

.PHONY: bin