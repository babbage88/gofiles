DOCKER_HUB:=jtrahan88/gofiles:
SHELL := /bin/bash

provisionbuilder:
	docker buildx create --name gfilesbuilder --use
	docker buildx inspect --bootstrap

tailwind:
	./tailwindcss -i templates/app.css -o static/tailwind.css

buildandpushlocalk3: tailwind
	docker buildx use gfilesbuilder
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_HUB_TEST)$(tag) . --push

run-local: tailwind
	go run .