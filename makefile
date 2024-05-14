#Dockerfile vars
TAG=0.0.0

#vars
IMAGENAME=go-web-cache
REPO=docker.io/supporttools
IMAGEFULLNAME=${REPO}/${IMAGENAME}:${TAG}

.PHONY: help build push all

help:
	    @echo "Makefile arguments:"
	    @echo ""
	    @echo "tag - Docker Tag"
	    @echo ""
	    @echo "Makefile commands:"
	    @echo "build"
	    @echo "push"
		@echo "bump"
	    @echo "all"

.DEFAULT_GOAL := all

build:
		@docker buildx build --platform linux/amd64 --pull --build-arg VERSION=${TAG} --build-arg GIT_COMMIT=`git rev-parse HEAD` --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") --cache-from supporttools/go-web-cache:latest -t supporttools/go-web-cache:${TAG} -t supporttools/go-web-cache:latest --push -f Dockerfile .

push:
	    @docker push ${IMAGEFULLNAME}
		@docker tag ${IMAGEFULLNAME} ${REPO}/${IMAGENAME}:latest
		@docker push ${REPO}/${IMAGENAME}:latest

bump:
		@make build push

all: build push