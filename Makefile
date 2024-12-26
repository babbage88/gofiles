DOCKER_HUB:=jtrahan88/gofiles:
SHELL := /bin/bash

provisionbuilder:
	docker buildx create --name gfilesbuilder --use
	docker buildx inspect --bootstrap

tailwind:
	./tailwindcss -i templates/app.css -o static/tailwind.css

buildandpushimage: tailwind
	docker buildx use gfilesbuilder
	docker buildx build --platform linux/amd64 -t $(DOCKER_HUB_TEST)$(tag) . --push

run-local: tailwind
	go run .