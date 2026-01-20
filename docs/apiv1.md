# V1 API

All API calls should be POST requests to `/api/v1/room/{board_id}`

The body of each request determines the execution of the event. In general, each event should have an `"event"` and an optional `"value"`:

- `{"event": "isprotected"}`

- `{"event": "checkpassword", "value": "mypass"}`

- `{"event": "debug"}`

- `{"event": "ping"}`

- `{"event": "upload_sgf", "value": "<base64 encoded sgf>"}`

- `{"event": "request_sgf", "value": "<url>"}`

- `{"event": "trash"}`

- `{"event": "graft", "value": "<string>"}`

- `{"event": "add_stone", "value": {"coords": [9, 9], "color": 1}}`

- `{"event": "pass", "value" 1}`

- `{"event": "remove_stone", "value": [9, 9]}`

- `{"event": "triangle", "value": [9, 9]}`

- `{"event": "square", "value": [9, 9]}`

- `{"event": "letter", "value": {"coords": [9, 9], "letter": "A"}}`

- `{"event": "number", "value": {"coords": [9, 9], "number": 7}}`

- `{"event": "label", "value": {"coords": [9, 9], "label": "O_o"}}`

- `{"event": "remove_mark", "value": [9, 9]}`

- `{"event": "cut"}`

- `{"event": "left"}`

- `{"event": "right"}`

- `{"event": "up"}`

- `{"event": "down"}`

- `{"event": "rewind"}`

- `{"event": "fastforward"}`

- `{"event": "goto_grid", "value": 23}`

- `{"event": "goto_coord", "value": [9, 9]}`

- `{"event": "comment", "value": "some comment"}`

- `{"event": "draw", "value": [0.1, 0.1, 0.2, 0.2, "#000000"]}`

- `{"event": "erase_pen"}`

- `{"event": "copy"}`

- `{"event": "clipboard"}`

- `{"event": "score"}`

- `{"event": "markdead", "value": [9, 9]}`

The response will contain `"success": true` or `"success": false` dependending on if the event was successfully received and interpreted by the server. Note: there might be some ambiguity about what constitutes a "failure." For example, if `"add_stone"` is attempted on a coordinate pair where a stone already exists, nothing will happen but the server will return a `"success": true` as the event was successfully processed (although the board state will remain unchanged).

If the event causes a board update, in addition to returning `"success": true`, the server will also return a `"frame"` event. If no frame event is present, the server will default to just returning the precipitating event.
