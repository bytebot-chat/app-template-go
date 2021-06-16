# App-template-go

This is a template app for a bytebot app written in go, fork it to create your own.

## Running locally
`go run .`

## Running within a docker container
`docker build . -t bytebot-template-app`
`docker run bytebot-template-app # Add --net=host if you are running the gateway locally with docker-compose`

## TODO
- github actions?
- k8s deployment?
