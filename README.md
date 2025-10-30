# Board

[![Main](https://github.com/jarednogo/board/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/jarednogo/board/actions/workflows/main.yml)
[![Test](https://github.com/jarednogo/board/actions/workflows/test.yml/badge.svg)](https://github.com/jarednogo/board/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/jarednogo/board/blob/main/LICENSE)

This is the only online go board that allows for synchronous control across all participants.

This project is free and open-source and always will be. Feel free to contribute by submitting PRs, reporting bugs, suggesting features, or just sharing with friends.

[Main page](https://board.tripleko.com)

[Test page](https://board-test.tripleko.com)

[Discord](https://discord.gg/y4wGZyed3e)

## Developing

If you make a pull request, please use `test` as the target branch. The test domain (above) tracks the `test` branch while the main domain tracks the `main` branch.

### Running locally

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

### Running locally with docker

1. Install docker

2. Build the docker container

```bash
$ docker build . -t board
```

3. Run the docker container, binding the container ports to your host ports

```bash
$ docker run -p 8080:8080 board
```

4. Visit `http://localhost:8080` in your browser

### Making changes

1. Checkout a new local branch (base it off `test`)

2. After making changes, run `make lint` and `make test` to ensure code uniformity.

3. Make PRs against the `test` branch.
