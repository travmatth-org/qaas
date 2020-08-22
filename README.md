# Fortunes As A Service (FAAS)

A ridiculously overengineered cloud native application, managing all of your fortune-telling needs.

## Development

The project is managed inside of a Docker container. If you're using VSCode, as I am, you'll be able to remotely connect to this container automatically, otherwise you're able to build a container with an image based on `.dev/container/Dockerfile`. There are several development related commands available in the Makefile:

```
# compile the `faas` executable
make build
# compile the `faas` executable and move all website assets into /dist folder for deployment 
make build.all
# compile & run `faas` with development options, for local testing over localhost:8080
make run
# remove the compiled binary and associated cache files
make clean
# run linter on project files
make lint
# run `go vet` on project files
make vet
# run tests on project files
make test
# lint, vet, test
make cicd
# generate & open html report on code coverage
make coverage.view
```

## Deployment