/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { uuid, opposite, Coord } from './common.js';
export {
    NetworkHandler
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";

// marks that can be placed on many cells with a click-and-drag
const DRAG_MARKS = new Set(["triangle", "square", "letter", "number"]);

// specifically into 4 bytes
function pack_int(n) {
    let byte_array = new Uint8Array([0, 0, 0, 0]);
    for (let i=0; i<4; i++) {
        let each_byte = n & 0xff;
        byte_array[i] = each_byte;
        n = (n-each_byte)/256;
    }
    return byte_array;
}

function wait_for_open(ws) {
  if (ws.readyState === WebSocket.OPEN) {
    return Promise.resolve();
  }

  return new Promise((resolve, reject) => {
    const onOpen = () => {
      ws.removeEventListener("open", onOpen);
      ws.removeEventListener("error", onError);
      resolve();
    };

    const onError = (evt) => {
      ws.removeEventListener("open", onOpen);
      ws.removeEventListener("error", onError);
      reject(evt);
    };

    ws.addEventListener("open", onOpen);
    ws.addEventListener("error", onError);
  });
}

// it's really both a network handler and event handler
class NetworkHandler {
    constructor(url) {
        this.url = url;
        this.connect();

        this.add_listeners();

        this.blocked = null;

        // click-and-drag marks: while the pointer is down the cells are only
        // previewed locally, then sent as a single "marks" event on release.
        // null when no gesture is in progress
        this.drag = null;

        // set when a drag was just sent (or cancelled partway), so the
        // click event that follows the release doesn't also act on the
        // last cell
        this.drag_committed = false;

        // in a var for now so we can do exponential backoff in the future
        this.backoff = 0;

        // necessary because of cloudfront's auto-timeout max of 60sec
        setInterval(() => this.keep_warm(), 30000)
    }

    connect() {
        this.socket = new WebSocket(this.url);
        this.socket.onmessage = (event) => this.onmessage(event);
        this.socket.onopen = (event) => this.onopen(event);
        this.socket.onclose = (event) => this.reconnect(event);
    }

    debug(message) {
        let payload = {"event": "debug", "value": message};
        this.send(payload);
    }

    onopen(event) {
        console.log("connected!");

        // hide the info modal
        //this.state.modals.hide_modal("info-modal");
        this.state.modals.hide_broken_connection_icon();

        // check to see if the room is protected by a password
        this.prepare_isprotected();

        // send in saved nickname, if previously entered
        let nickname = this.state.modals.get_nickname();
        if (nickname != "") {
            this.prepare_nickname(this.state.modals.get_nickname());
        }

        // reset backoff
        this.backoff = 0;
    }

    reconnect(event) {
        this.state.reset();
        this.state.modals.show_broken_connection_icon();
        /*
        if (!this.state.modals.modals_up.has("info-modal")) {
            this.state.modals.show_info_modal("Reconnecting...");
        }
        */
        if (document.hasFocus() && this.socket.readyState == WebSocket.CLOSED) {
            console.log("reconnecting in", this.backoff, "ms");
            setTimeout(() => this.connect(), this.backoff);
            this.backoff *= 2;
            if (this.backoff == 0) {
                this.backoff = 500;
            } else if (this.backoff > 120000) {
                this.backoff = 0;
            }

        }
    }

    // evidently necessary because unfocused tabs don't hide modals
    focus() {
        // clear keys down
        this.state.keys_down = new Map();

        // hide info modal
        if (this.socket.readyState == WebSocket.OPEN) {
            this.state.modals.hide_broken_connection_icon();
            //this.state.modals.hide_modal("info-modal");
        } else if (this.socket.readyState == WebSocket.CLOSED) {
            this.reconnect();
        }
    }

    add_listeners() {
        document.addEventListener("click", (event) => this.click(event));
        document.addEventListener("pointermove", (event) => this.pointermove(event));
        document.addEventListener("pointerdown", (event) => this.pointerdown(event));
        document.addEventListener("pointerup", (event) => this.pointerup(event));
        document.addEventListener("mouseout", (event) => this.mouseout(event));

        document.addEventListener("touchstart", (event) => this.touchstart(event));
        document.addEventListener("touchend", (event) => this.touchend(event));
        document.addEventListener("touchcancel", (event) => this.touchcancel(event));
        // the passive arg is necessary to disable regular touch events
        document.addEventListener("touchmove", (event) => this.touchmove(event), {passive: false});

        document.addEventListener("keydown", (event) => this.keydown(event));
        document.addEventListener("keyup", (event) => this.keyup(event));

        window.addEventListener("focus", () => this.focus());
    }

    set_state(state) {
        this.state = state;
    }

    keep_warm() {
        let payload = {"event": "ping"};
        this.prepare(payload);
    }

    fromserver(payload) {
        let evt = payload["event"];
        let id = payload["eventid"];
        if (this.blocked == id) {
            this.blocked = null;
        }
        var coords;
        var index;
        var value;
        var userid;
        var label;
        switch (evt) {
            case "frame":
                this.state.handle_frame(payload["value"]);
                if (payload["value"].tree != null) {
                    // the board moved to another node mid-drag
                    this.drag_cancel();
                } else {
                    this.drag_refresh();
                }
                break;
            case "triangle":
                coords = new Coord(payload["value"][0], payload["value"][1]);
                this.state.place_triangle(coords);
                break;
            case "square":
                coords = new Coord(payload["value"][0], payload["value"][1]);
                this.state.place_square(coords);
                break;
            case "letter":
                coords = new Coord(payload["value"]["coords"][0], payload["value"]["coords"][1]);
                label = new Object();
                label.coord = coords;
                label.text = payload["value"]["letter"];
                this.state.place_label(label);
                break;
            case "number":
                coords = new Coord(payload["value"]["coords"][0], payload["value"]["coords"][1]);
                label = new Object();
                label.coord = coords;
                label.text = payload["value"]["number"];
                this.state.place_label(label);

                //this.state.place_number(coords[0], coords[1], payload["value"]["number"]);
                break;
            case "label":
                coords = new Coord(payload["value"]["coords"][0], payload["value"]["coords"][1]);
                label = new Object();
                label.coord = coords;
                label.text = payload["value"]["label"];
                this.state.place_label_string(label);
                break;
            case "remove_mark":
                coords = payload["value"];
                this.state.remove_mark(coords[0], coords[1]);
                break;
            case "trash":
                this.state.reset();
                break;
            case "update_settings":
                value = payload["value"];
                this.state.update_settings(value);
                break;
            case "draw":
                let [x0,y0,x1,y1, pen_color] = payload["value"];
                this.state.draw_pen(x0, y0, x1, y1, pen_color);
                break;
            case "erase_pen":
                this.state.erase_pen();
                break;
            case "comment":
                this.state.comments.update(payload["value"]);
                //this.state.comments.store(payload["value"]);
                break;
            case "error":
                value = payload["value"];
                this.state.modals.show_error_modal(value);
                break;
            case "isprotected":
                if (payload["value"]) {
                    if (this.state.password != "") {
                        // if we already have a password saved, send it
                        this.prepare_checkpassword(this.state.password);
                    } else {
                        let handler = () => {
                            let v = this.state.modals.get_prompt_bar();
                            this.prepare_checkpassword(v);
                        }
                        this.state.modals.show_prompt_modal(
                            "Enter password:",
                            handler
                        );
                    }
                }
                break;
            case "checkpassword":
                if (payload["value"] != "") {
                    this.state.update_password(payload["value"]);
                } else {
                    this.state.modals.show_toast(
                        "Wrong password. You can observe, but not edit"
                    );
                }
                break;
            case "global":
                value = payload["value"];
                this.state.modals.show_toast(value);
                break;
            case "connected_users":
                value = payload["value"];
                this.state.handle_current_users(value);
                break;
        }
    }

    prepare_pass() {
        let payload = {"event":"pass", "value":this.state.color};
        this.prepare(payload);
    }

    prepare_rewind() {
        let payload = {"event":"rewind"};
        this.prepare(payload);
    }

    prepare_fastforward() {
        let payload = {"event":"fastforward"};
        this.prepare(payload);
    }

    prepare_left() {
        let payload = {"event":"left"};
        this.prepare(payload);
    }

    prepare_right() {
        let payload = {"event":"right"};
        this.prepare(payload);
    }

    prepare_up() {
        let payload = {"event":"up"};
        this.prepare(payload);
    }

    prepare_down() {
        let payload = {"event":"down"};
        this.prepare(payload);
    }

    async prepare_upload(data) {
        let payload = {"event":"upload_sgf", "value": data};

        await wait_for_open(this.socket);
        this.prepare(payload);
    }

    prepare_request(url) {
        let payload = {"event":"request_sgf", "value": url};
        this.prepare(payload);
    }

    prepare_link_ogs_game(url) {
        let payload = {"event":"link_ogs_game", "value": url};
        this.prepare(payload);
    }

    prepare_trash() {
        let payload = {"event":"trash", "value": "all"};
        this.prepare(payload);
    }

    prepare_cut() {
        let payload = {"event": "cut"};
        this.prepare(payload);
    }

    prepare_copy() {
        let payload = {"event": "copy"};
        this.prepare(payload);
    }

    prepare_clipboard() {
        let payload = {"event": "clipboard"};
        this.prepare(payload);
    }

    prepare_score() {
        let payload = {"event": "score"};
        this.prepare(payload);
    }

    prepare_markdead(coord) {
        let payload = {"event": "markdead", "value": coord};
        this.prepare_payload(payload);
    }

    prepare_settings(settings) {
        let payload = {"event": "update_settings", "value": settings};
        this.prepare(payload);
    }

    prepare_erase_pen() {
        let payload = {"event":"erase_pen"};
        this.prepare(payload);
    }

    prepare_comment(text) {
        let payload = {"event": "comment", "value": text};
        this.prepare(payload);
    }

    prepare_isprotected() {
        let payload = {"event": "isprotected"};
        this.prepare(payload);
    }

    prepare_checkpassword(text) {
        let payload = {"event": "checkpassword", "value": text};
        this.prepare(payload);
    }

    prepare_nickname(text) {
        let payload = {"event": "update_nickname", "value": text};
        this.prepare(payload);
    }

    prepare(payload) {
        // before anything check for integrity of socket
        if (this.socket.readyState != WebSocket.OPEN) {
            return;
        }

        let evt = payload["event"];

        // if a modal is up, then the events we allow are:
        // "label"
        // "trash"
        // "update_settings"
        // "scissors"
        // "upload_sgf"
        // "request_sgf"
        // "link_ogs_game"
        // "checkpassword"
        // "isprotected"
        // "update_nickname"
        if (
            this.state.modals.modals_up.size > 0 &&
            evt != "label" &&
            evt != "trash" &&
            evt != "update_settings" &&
            evt != "cut" &&
            evt != "upload_sgf" &&
            evt != "request_sgf" &&
            evt != "link_ogs_game" &&
            evt != "checkpassword" &&
            evt != "isprotected" &&
            evt != "update_nickname") {
            return;
        }

        // pointer movements don't get shared...
        // ...unless the pointer is down or the mark is "pen"
        if (evt == "pointermove") {
            let coords = payload["value"];
            if (this.state.mark != "") {
                this.state.board_graphics.draw_ghost_mark(coords[0], coords[1]);
            } else {
                this.state.board_graphics.draw_ghost_stone(coords[0], coords[1], this.state.color);
            }
            return;
        }

        this.send(payload);
    }

    send(payload) {
        // blocking logic
        if (this.blocked != null) {
            return;
        }
        let id = uuid();
        payload["eventid"] = id;

        // only block under certain circumstances
        let type = payload["event"];
        if (type == "add_stone" || type == "goto_grid") {
            this.blocked = id;
        }

        //console.log("sending:", payload);
        
        // first create the json payload
        let json_payload = JSON.stringify(payload);

        // find the length in bytes
        let encoder = new TextEncoder();
        let payload_bytes = encoder.encode(json_payload);

        // then send the length to the socket
        let length = pack_int(payload_bytes.length);
        this.socket.send(length);

        // then send the payload
        this.socket.send(payload_bytes);
    }

    onmessage(event) {
        //console.log("receiving:", event.data);
        let payload = JSON.parse(event.data);
        this.fromserver(payload);
    }

    keydown(event) {
        let payload = {"event": "keydown", "value": event.key};
        let shift = this.state.keys_down.has("Shift");
        let ctrl = this.state.keys_down.has("Control");
        let alt = this.state.keys_down.has("Alt");
        let meta = this.state.keys_down.has("Meta");
        // logical xor
        let jump = this.state.branch_jump != shift
        let on_input_bar = document.activeElement.tagName == "INPUT";
        if (on_input_bar) {
            return;
        }

        switch(event.key) {
            case "ArrowUp":
                if (jump){
                    let index = this.state.get_index_up();
                    if (index == -1) {
                        return;
                    }
                    payload = {"event": "goto_grid", "value": index};
                    this.prepare(payload);
                } else {
                    this.prepare_up();
                }
                break;

            case "ArrowDown":
                if (jump){
                    let index = this.state.get_index_down();
                    if (index == -1) {
                        return;
                    }
                    payload = {"event": "goto_grid", "value": index};
                    this.prepare(payload);
                } else {
                    this.prepare_down();
                }
                break;

            case "ArrowLeft":
                if (shift || ctrl || alt || meta) {break;}
                this.prepare_left();
                break;
            case "ArrowRight":
                if (shift || ctrl || alt || meta) {break;}
                this.prepare_right();
                break;
            
            case "1":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_toggle();
                break;
            case "2":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_black();
                break;
            case "3":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_white();
                break;
            case "4":
                if (shift || ctrl || alt || meta) {break;}
                this.prepare_pass();
                break;
            case "5":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_triangle();
                break;
            case "6":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_square();
                break;
            case "7":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_letter();
                break;
            case "8":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_number();
                break;
            case "9":
                if (shift || ctrl || alt || meta) {break;}
                this.state.set_pen();
                break;
            case "0":
                if (shift || ctrl || alt || meta) {break;}
                this.prepare_erase_pen();
                break;
            case "x":
                if (ctrl) {
                    this.prepare_cut();
                }
                break;
            case "c":
                if (ctrl) {
                    this.prepare_copy();
                }
                break;
            case "v":
                if (ctrl) {
                    this.prepare_clipboard();
                }
                break;
            case "Enter":
                if (ctrl) {
                    this.state.trigger_score();
                }
                break;
            default:
                this.state.keys_down.set(event.key, true);
        }
    }

    keyup(event) {
        if (this.state.keys_down.has(event.key)) {
            this.state.keys_down.delete(event.key);
        }
    }

    // begin a (potential) click-and-drag of marks. nothing is drawn or
    // sent until the pointer reaches a second cell, so a plain click
    // still goes through click()
    drag_start(clientX, clientY, target) {
        this.drag = null;
        this.drag_committed = false;
        this.state.board_graphics.clear_drag_marks();
        if (!DRAG_MARKS.has(this.state.mark)) {
            return;
        }
        if (this.state.modals.modals_up.size > 0) {
            return;
        }
        let review = document.getElementById("review");
        if (review == null || !review.contains(target)) {
            return;
        }
        let coords = this.state.board_graphics.pos_to_coord(clientX, clientY);
        if (!this.in_bounds(coords)) {
            return;
        }
        this.drag = {
            mark: this.state.mark,
            start: coords,
            last: coords,
            active: false,
            seen: new Set(),
            cells: [],
        };
    }

    in_bounds(coords) {
        let size = this.state.size;
        return coords[0] >= 0 && coords[1] >= 0 && coords[0] < size && coords[1] < size;
    }

    drag_add(x, y) {
        let id = x.toString() + "-" + y.toString();
        if (this.drag.seen.has(id)) {
            return;
        }
        this.drag.seen.add(id);
        // cells that already have a mark are left alone
        // (the server skips them too)
        if (this.state.marks.has(id)) {
            return;
        }
        this.drag.cells.push([x, y]);
    }

    // returns true if a drag is underway (so the caller can skip ghosts)
    drag_move(clientX, clientY) {
        if (this.drag == null) {
            return false;
        }
        if (this.state.mark != this.drag.mark) {
            // tool changed mid-drag
            this.drag_cancel();
            return false;
        }
        let coords = this.state.board_graphics.pos_to_coord(clientX, clientY);
        if (!this.in_bounds(coords)) {
            return this.drag.active;
        }
        let [x0, y0] = this.drag.last;
        let [x1, y1] = coords;
        if (x0 == x1 && y0 == y1) {
            return this.drag.active;
        }
        if (!this.drag.active) {
            this.drag.active = true;
            this.drag_add(x0, y0);
        }
        // a fast drag skips cells, so fill in the line from the last one
        let steps = Math.max(Math.abs(x1 - x0), Math.abs(y1 - y0));
        for (let i = 1; i <= steps; i++) {
            this.drag_add(
                Math.round(x0 + (x1 - x0)*i/steps),
                Math.round(y0 + (y1 - y0)*i/steps));
        }
        this.drag.last = coords;
        this.state.board_graphics.clear_ghosts();
        this.state.board_graphics.draw_drag_marks(this.drag.cells);
        return true;
    }

    // pointer released: send the whole drag as one event
    drag_end() {
        let drag = this.drag;
        this.drag = null;
        if (drag == null || !drag.active) {
            return;
        }
        this.state.board_graphics.clear_drag_marks();
        this.drag_committed = true;
        if (drag.cells.length == 0 || this.state.mark != drag.mark) {
            return;
        }
        this.prepare(this.marks_payload(drag.mark, drag.cells));
    }

    // the "marks" event places one kind of mark on a list of cells (a
    // single click is just a list of one). the server skips cells that
    // are taken and hands out the letters/numbers itself, so two users
    // can't end up with the same one
    marks_payload(mark, cells) {
        let value = {"mark": mark, "coords": cells};
        if (mark == "letter" && this.state.custom_label != "") {
            value["mark"] = "label";
            value["label"] = this.state.custom_label;
        }
        // lets the server drop the event if someone navigated away
        // before it arrived
        if (this.state.current_index != null) {
            value["node"] = this.state.current_index;
        }
        return {"event": "marks", "value": value};
    }

    drag_cancel() {
        if (this.drag == null) {
            return;
        }
        if (this.drag.active) {
            // the release is still coming; it shouldn't count as a click
            this.drag_committed = true;
        }
        this.drag = null;
        this.state.board_graphics.clear_drag_marks();
    }

    // marks changed under an in-progress drag (someone else placed some):
    // drop the cells that are now taken
    drag_refresh() {
        if (this.drag == null || !this.drag.active) {
            return;
        }
        this.drag.cells = this.drag.cells.filter(
            ([x, y]) => !this.state.marks.has(x.toString() + "-" + y.toString()));
        this.state.board_graphics.draw_drag_marks(this.drag.cells);
    }

    pointerdown(event) {
        if (event.pointerType == "mouse") {
            this.state.ispointerdown = true;
            if (event.button == 0) {
                this.drag_start(event.clientX, event.clientY, event.target);
            }
            //let payload = {"event": "pointerdown"};
            //this.prepare(payload);
        }
    }

    pointerup(event) {
        if (event.pointerType == "mouse") {
            this.state.ispointerdown = false;
            this.state.penx = null;
            this.state.peny = null;
            this.drag_end();

            //let payload = {"event": "pointerup"};
            //this.prepare(payload);
        }
    }

    mouseout(event) {
        this.state.board_graphics.clear_ghosts();
    }

    pointermove(event) {
        if (event.pointerType == "mouse") {
            if (this.state.mark == "pen" && this.state.ispointerdown) {
                let [x,y,inside] = this.state.board_graphics.board_relative_coords(event.clientX, event.clientY);
                let x0 = this.state.penx;
                if (x0 == null) {
                    x0 = x;
                }
                let y0 = this.state.peny;
                if (y0 == null) {
                    y0 = y;
                }
                let payload = {"event": "draw", "value": [x0, y0, x, y, this.state.pen_color]};
                this.state.penx = x;
                this.state.peny = y;
                this.prepare(payload);
            } else {
                if (this.drag_move(event.clientX, event.clientY)) {
                    // the drag preview replaces the ghost
                    return;
                }
                let coords = this.state.board_graphics.pos_to_coord(event.clientX, event.clientY);
                let payload = {"event": "pointermove", "value": coords};
                this.prepare(payload);
            }
        }
        // mobile doesn't get ghosts
    }

    touchstart(event) {
        this.state.ispointerdown = true;
        if (event.touches.length == 1) {
            let touch = event.touches[0];
            this.drag_start(touch.clientX, touch.clientY, event.target);
        } else {
            // a second finger means pinch/scroll, not a drag
            this.drag_cancel();
        }
    }

    touchend(event) {
        this.state.ispointerdown = false;
        this.state.penx = null;
        this.state.peny = null;
        this.drag_end();
    }

    touchcancel(event) {
        this.state.ispointerdown = false;
        this.state.penx = null;
        this.state.peny = null;
        this.drag_cancel();
    }

    touchmove(event) {
        if (this.state.mark == "pen" && this.state.ispointerdown) {
            let touch = event.touches[0];
            let [x,y,inside] = this.state.board_graphics.board_relative_coords(touch.clientX, touch.clientY);
            if (inside) {
                // necessary for capturing touch inside svg
                // for example, disables "pull refresh" on android
                event.preventDefault();
            }
            let x0 = this.state.penx;
            if (x0 == null) {
                x0 = x;
            }
            let y0 = this.state.peny;
            if (y0 == null) {
                y0 = y;
            }

            let payload = {"event": "draw", "value": [x0, y0, x,y, this.state.pen_color]};
            this.state.penx = x;
            this.state.peny = y;

            this.prepare(payload);
        } else if (this.drag != null) {
            let touch = event.touches[0];
            let [x,y,inside] = this.state.board_graphics.board_relative_coords(touch.clientX, touch.clientY);
            if (inside) {
                // keep the page from scrolling while dragging marks
                event.preventDefault();
            }
            this.drag_move(touch.clientX, touch.clientY);
        }
    }

    click(event) {
        // check to see if it's within explorer
        let x = event.clientX;
        let y = event.clientY;
        let flip = this.state.keys_down.has("Shift");

        // this click is just the release of a drag
        if (this.drag_committed) {
            this.drag_committed = false;
            return;
        }

        // only clicks that land directly on the explorer or the board count
        // otherwise a click on something overlapping them (e.g. the
        // collapsed button dropdown) would also place a stone underneath
        if (this.state.tree_graphics.container.contains(event.target)) {
            //let explorer_node = this.state.tree_graphics.capture_mouse(x, y);
            let index = this.state.tree_graphics.capture_mouse(x, y);
            if (index != -1) {
                let payload = {"event": "goto_grid", "value": index};
                this.prepare(payload);
            }
            return;
        }

        let review = document.getElementById("review");
        if (review == null || !review.contains(event.target)) {
            return;
        }

        let coords = this.state.board_graphics.pos_to_coord(x, y);

        // don't need to share if click is outside the board
        if (coords[0] < 0 || coords[1] < 0 || coords[0] >= this.state.size || coords[1] >= this.state.size) {
            return;
        }

        // also don't need to share if there is already a stone there
        let has_child = false;
        let stone_there = this.state.board.get(new Coord(coords[0], coords[1])) != 0;

        let payload = {};

        if (this.state.mark == "score") {
            payload = {"event": "markdead", value: coords};
        } else if (this.state.mark != "") {
            if (this.state.mark == "pen") {
                return;
            }

            let id = coords[0].toString() + "-" + coords[1].toString();
            if (this.state.marks.has(id)) {
                payload = {"event": "remove_mark", value: coords};
            } else {
                if (!DRAG_MARKS.has(this.state.mark)) {
                    return;
                }
                if (this.state.mark == "letter") {
                    if (flip) {
                        // shift+click: ask for a custom label
                        // (sent as "label" since that's what's allowed
                        // through while the prompt is still up)
                        let handler = () => {
                            let label = this.state.modals.get_prompt_bar();
                            this.state.modals.clear_prompt_bar();
                            this.state.custom_label = label;
                            payload = {"event": "label"};
                            payload["value"] = {"coords": coords, "label": label};
                            this.prepare(payload);
                        }
                        this.state.modals.show_prompt_modal(
                            "Enter up to three characters:",
                            handler
                        );
                        return;
                    }
                    if (this.state.custom_label == "" && this.state.next_letter() == null) {
                        // out of letters
                        return;
                    }
                }
                // a click is a drag over a single cell
                payload = this.marks_payload(this.state.mark, [coords]);
            }
        } else if (this.state.toggling) {
            if (this.state.keys_down.has("Shift")) {
                payload = {"event": "goto_coord", "value": coords};
            } else {
                if (stone_there) {
                    return;
                }
                //payload = {"event": "stone-toggle", "value": coords, "color": this.state.color};
                payload = {"event": "add_stone", "value": {"coords": coords, "color": this.state.color}};
            }
        } else {
            if (stone_there) {
                // if there is a stone there already and we are in manual mode
                // then remove the stone
                payload = {"event": "remove_stone", "value": coords};
            } else {
                let color = this.state.color;
                if (flip) {
                    color = opposite(color);
                }
                //payload = {"event": "stone-manual", "value": coords, "color": color};
                payload = {"event": "add_stone", "value": {"coords": coords, "color": color}};
            }
        }
        this.prepare(payload);
    }
}
