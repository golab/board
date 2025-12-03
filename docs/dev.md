# Developer Guide

For additional documentation, see:
- The [design doc](design.md)

## Quickstart

1. Install golang

2. Install the local toolchain

```bash
$ make setup
```

3. Run the server

```bash
$ make run
```

4. Visit `http://localhost:8080` in your browser.

## Custom Setups

### Docker

A standalone Dockerfile has been provided to make containerization simple and straightforward.

1. Install docker

2. Build the docker container

```bash
$ make build-docker
```

This make command does this under the hood:
```
VERSION := $(shell git describe --tags)
docker build --build-arg VERSION=$(VERSION) -t board .
```
The intention is that the developer always has a way of checking which version (with respect to git tags) the container is running.

You can verify this on a running container with 
```bash
$ curl http://localhost:8080/api/version
```

3. Run the docker container, binding the container ports to your host ports

```bash
$ make run-docker
```

4. Visit `http://localhost:8080` in your browser

### Monitoring

The docker-compose file is to setup additional services for monitoring logs:
- Alloy reads log files (at `$LOG_PATH`, defaults to `./logs`) and sends to Loki
- Loki stores logs
- Grafana takes logs from Loki and displays in a UI (at `http://localhost:3000`)

1. Build and run the monitoring services (also builds the Board app itself in a container)

```bash
$ make monitoring-up
```

2. Visit grafana in your browser `http://localhost:3000` (user/password = admin/admin) 

3. Close down the monitoring services
```bash
$ make monitoring-down
```

## Pull Request Guidelines

The `main` branch is rarely updated (a couple times a month). Make new branches based off the `test` branch.

After making changes, run `make lint` to catch linting errors. `make fmt` will fix a lot of issues, but some you may have to fix by hand. Use `//nolint` sparingly.

Run `make test` to ensure all tests pass. Add new tests for new code.

Run `make test-race` to catch potential data races.

Submit pull request.
