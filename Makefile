DOCKER_HUB:=jtrahan88/gofiles:
SHELL := /bin/bash

provisionbuilder:
	docker buildx create --name gfilesbuilder --use
	docker buildx inspect --bootstrap

tailwind:
	./tailwindcss -i templates/app.css -o static/tailwind.css

buildandpushimage: tailwind
	docker buildx use gfilesbuilder
	docker buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_HUB)$(tag) . --push

deploylocalk3: buildandpushimage
	kubectl --kubeconfig ~/.kube/config-local rollout restart deployment gofiles

deployprodk3: buildandpushimage
	kubectl rollout restart deployment gofiles

run-local: tailwind
	go run .
