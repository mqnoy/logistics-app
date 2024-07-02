# Build target
APP_NAME=core
TAG ?= latest

build-image:
	@latest_tag=$$(git describe --tags --abbrev=0 | sed 's/^v//'); \
	echo "Building docker images with version $$latest_tag"; \
	docker build -t $APP_NAME:$$latest_tag .

run:
	@latest_tag=$$(git describe --tags --abbrev=0 | sed 's/^v//'); \
	TAG=$$latest_tag  docker compose up

push-image:
	@latest_tag=$$(git describe --tags --abbrev=0 | sed 's/^v//'); \
	docker push registry.example.com/repo-name:$$latest_tag