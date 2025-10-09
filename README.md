# Board

This is the only online go board that allows for synchronous control across all participants.

This project is free and open-source and always will be. Feel free to contribute by submitting PRs, reporting bugs, suggesting features, or just sharing with friends.

[Main page](https://board.tripleko.com)

[Test page](https://board-test.tripleko.com)

[Discord](https://discord.gg/y4wGZyed3e)

## Developing

If you make a pull request, please use `test` as the target branch. The test domain (above) tracks the `test` branch while the main domain tracks the `main` branch.

### Running locally

1. Install golang

2. Run the server

```bash
$ make run
```

3. Visit `http://localhost:8080` in your browser.

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

### Developing

After installing golang, install the local toolchain:

```bash
$ make setup
```

Create a branch based on the `test` branch.

Use `make run-live` to build and run with `air` (which provides live reloading). 

After making changes, run `make test` and `make lint` to ensure code uniformity.

Make PRs against the `test` branch.
