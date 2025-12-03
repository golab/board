# Design

## Project layout

For the most part, source code is stored under `/pkg`.

The notable exceptions would be:
- `/cmd` which contains only `main.go` and is the entry point to the application.
- `/internal` which contains a few subpackages exclusively for testing.
- `/integration` which contains integration and benchmark tests.

## `/pkg`

To avoid circular imports, we will be rather strict about subpackages, and import direction.  To this end, subpackages which will fall into three main types:
- `/core`, which contains types and functions that will need to be used ubiquitously throughout the rest of the application. (Examples: `color.go`, `board.go`, `parser.go`).
- Standalone packages which do not import from any other subpackage and are also not imported in `/core` (Examples: `/zip`, `/logx`, `config`).
- The "main" packages (i.e. `/state`, `/room`, `/hub`, and `/app')

### `core`

It is expected that the `core` package contains code that will be needed in basically all the other packages that are related to the application's network and game logic. For example, anytime there is a reason to make reference to black and white (the colors), we would use `core.Black` and `core.White`.

The `core` package should be standalone: it should not import from any other subpackage. This is one of the ways we can avoid circular dependencies. On the other hand, new subpackages should not necessarily be chucked into `core` by default. In fact, it is desirable to keep `core` as lean as possible. Only packages which cannot be easily decoupled from the rest of `core` or packages which are part of the fundamental business logic of the application should live here.

### `state` -> `room` -> `hub` -> `app`

The `state` package can import from `core` or any of the standalone packages.

The `room` package can import from `core`, `state`, or any of the standalone packages.

The `hub` package can import from `core`, `state`, `room`, or any of the standalone packages.

The `app` package can import from `core`, `state`, `room`, `hub`, or any of the standalone packages.
