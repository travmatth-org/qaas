# Fortunes As A Service (FAAS)

A ridiculously overengineered cloud native application, managing all of your fortune-telling needs.

## Development

The project is managed inside of a Docker container. If you're using VSCode, as I am, you'll be able to remotely connect to this container automatically, otherwise you're able to build a container with an image based on `.dev/container/Dockerfile`. There are several commands available in the Makefile:
```
make build # compile the `faas` executable
make run # compile & run `faas` with development options
make clean # remove the compiled binary and associated cache files
make lint # run linter on project files
make vet # run `go vet` on project files
make test # run tests on project files
make check # lint, vet, test
