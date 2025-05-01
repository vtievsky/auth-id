APP=auth-id
APP_TAG=slaventius/${APP}:latest
# APP_TAG=slaventius/${APP}:v$(shell date +%s)

docker-build:
	@echo "building docker-image ${APP_TAG}"
	@docker build --no-cache --tag ${APP_TAG} --file build/app/Dockerfile .

up:  ## Поднять и запустить зависимости в контейнерах
	docker compose -f build/docker-compose.yml up -d --build

down:  ## Остановить и удалить зависимости
	docker compose -f build/docker-compose.yml down --remove-orphans
