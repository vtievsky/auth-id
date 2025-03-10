APP=auth-id
APP_TAG=slaventius/${APP}:latest

docker-build:
	@echo "building docker-image ${APP_TAG}"
	@docker build --no-cache --tag ${APP_TAG} --file build/app/Dockerfile .