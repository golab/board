/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/


import { Board } from './board.js';
import { BoardGraphics } from './boardgraphics.js';
import { TreeGraphics } from './treegraphics.js';

import { create_comments } from './comments.js';
import { create_layout } from './layout.js';
import { create_buttons } from './buttons.js';
import { create_modals } from './modals.js';

import { htmlencode, letterstocoord, coordtoid, opposite, Coord, prefer_dark_mode } from './common.js';

export {
    State
}

const FrameType = {
    DIFF: 0,
    FULL: 1,
}

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";

function b64_encode_arraybuffer(buffer) {
    let binary = '';
    const bytes = new Uint8Array(buffer);
    const len = bytes.byteLength;
    for (let i = 0; i < len; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
}

function b64_encode_unicode(str) {
    const text_encoder = new TextEncoder('utf-8');
    const encoded_data = text_encoder.encode(str);
    return btoa(String.fromCharCode(...encoded_data));
}

class State {
    constructor() {
        window.addEventListener("resize", (event) => this.resize(event));
        create_layout();
        this.compute_consts();
        this.color = 1;
        this.toggling = true;
        this.handicap = false;
        this.mark = "";
        this.input_buffer = 250;
        this.password = "";

        // pen variables
        this.pen_color = "#48BBFF";
        this.ispointerdown = false;
        this.penx = null;
        this.peny = null;

        this.keys_down = new Map();

        this.branch_jump = true;

        this.dark_mode = false;
        this.board = new Board(this.size);
        this.marks = new Map();
        this.pen = new Array();
        this.current = null;
        this.numbers = new Map();
        this.letters = new Array(26).fill(0);

        this.board_color = "light";
        this.textured_stones = true;
        this.board_graphics = new BoardGraphics(this);
        this.tree_graphics = new TreeGraphics(this);

        this.comments = create_comments(this);
        this.connected_users = {};

        this.board_graphics.draw_board();

        this.buttons = create_buttons(this);

        this.modals = create_modals(this);

        if (prefer_dark_mode()) {
            this.dark_mode_toggle();
        }

        this.resize();
        this.gameinfo = new Map();
        this.custom_label = "";
    }

    set_network_handler(handler) {
        this.network_handler = handler;
    }

    set_board_color(color) {
        this.board_color = color;
        this.board_graphics.draw_boardbg();
    }

    set_textured_stones(toggle) {
        this.textured_stones = toggle;
        this.board_graphics.draw_stones();
    }

    guest_nick(id) {
        return "Guest-" + id.substring(0, 4);
    }

    handle_current_users(users) {
        this.connected_users = {};
        for (let id in users) {
            let nick = users[id];
            if (nick == "") {
                nick = this.guest_nick(id);
            }
            this.connected_users[id] = nick;
        }
        this.modals.update_users_modal();
    }

    update_password(password) {
        this.password = password;
        this.modals.update_settings_modal();
    }

    update_settings(settings) {
        this.input_buffer = settings["buffer"];
        if (settings["size"] != this.size) {
            this.size = settings["size"];
            let review = document.getElementById("review");
            review.setAttribute("size", this.size);
            this.recompute_consts();
            this.board_graphics.reset_board();
            this.reset();
        }
        this.password = settings["password"];
        this.modals.update_settings_modal();
    }

    resize(event) {
        let content = document.getElementById("content");
        let arrows = document.getElementById("arrows");
        let h = arrows.offsetHeight*4.5;
        let new_width = Math.min(window.innerHeight*1.5 - h, window.innerWidth);
        content.style.width = new_width + "px";

        this.recompute_consts();
        this.board_graphics.resize();
        this.tree_graphics.resize();
        this.comments.resize();
        this.apply_pen();
        // it's a little hacky, but the buttons were being very annoying
        setTimeout(() => this.buttons.resize(), 100);
    }

    get_index_up() {
        let [x,y] = this.tree_graphics.current;

        while (true) {
            y--;
            if (y < 0) {
                return -1;
            }

            if (!this.tree_graphics.grid.has(y)) {
                continue;
            }

            let row = this.tree_graphics.grid.get(y);
            if (row.has(x)) {
                return row.get(x).index;
            }
        }
    }

    get_index_down() {
        let [x,y] = this.tree_graphics.current;

        while (true) {
            y++;
            if (!this.tree_graphics.grid.has(y)) {
                return -1;
            }

            let row = this.tree_graphics.grid.get(y);
            if (row.has(x)) {
                return row.get(x).index;
            }
        }

    }

    next_letter() {
        for (let i = 0; i < 26; i++) {
            if (this.letters[i] == 0) {
                return letters[i];
            }
        }
        return null;
    }

    free_letter(l) {
        let letter_index = l.charCodeAt(0)-65;
        this.letters[letter_index] = 0;
    }

    next_number() {
        let i = 1;
        while (true) {
            if (this.numbers.get(i) == null) {
                return i;
            }
            i++;
        }
    }

    free_number(i) {
        this.numbers.delete(i);
    }

    reset() {
        this.board_graphics.clear_and_remove();
        this.handicap = false;
        this.color = 1;
        this.toggling = true;
        this.mark = "";

        this.board = new Board(this.size);

        // update move number
        this.set_move_number(0);

        // update comments
        this.comments.clear();

        this.modals.update_modals();
    }

    get_gameinfo() {
        return this.gameinfo;
    }

    set_gameinfo(fields_object) {
        let fields = new Map();
        for (let field of fields_object) {
            let key = field["key"];
            let values = field["values"];
            fields.set(key, values);
        }
        /*
        let fields = new Map(Object.entries(fields_object));
        if (fields == null) {
            fields = new Map();
        }
        */
        let gameinfo = {};

        // currently doesn't play very nice with chinese characters

        if (fields.has("PB")) {
            let rank = "";
            if (fields.has("BR")) {
                rank = " [" + fields.get("BR") + "]";
            }
            gameinfo["Black"] = fields.get("PB") + rank;
        } else {
            gameinfo["Black"] = "Black";
        }

        document.getElementById("black-name").innerHTML = htmlencode(gameinfo["Black"]);

        if (fields.has("PW")) {
            let rank = "";
            if (fields.has("WR")) {
                rank = " [" + fields.get("WR") + "]";
            }
            gameinfo["White"] = fields.get("PW") + rank;
        } else {
            gameinfo["White"] = "White";
        }
        document.getElementById("white-name").innerHTML = htmlencode(gameinfo["White"]);

        if (fields.has("RE")) {
            gameinfo["Result"] = fields.get("RE");
        }

        if (fields.has("KM")) {
            gameinfo["Komi"] = fields.get("KM");
            let komi = document.getElementById("komi");
            let komi_float = parseFloat(gameinfo["Komi"]);
            if (komi_float == 375) {
                komi_float = 7.5;
            }
            komi.innerHTML = " (" + komi_float + ")";
        } else {
            komi.innerHTML = "";
        }

        if (fields.has("DT")) {
            gameinfo["Date"] = fields.get("DT");
        }

        if (fields.has("RU")) {
            gameinfo["Ruleset"] = fields.get("RU");
        }

        if (fields.has("AB")) {
            this.handicap = true;
        } else {
            this.handicap = false;
        }

        /*
        if (fields.has("PC")) {
            gameinfo["Place"] = fields.get("PC");
        }

        if (fields.has("SO")) {
            gameinfo["Source"] = fields.get("SO");
        }

        if (fields.has("EV")) {
            gameinfo["Event"] = fields.get("EV");
        }

        if (fields.has("N")) {
            gameinfo["Name"] = fields.get("N");
        }

        if (fields.has("GN")) {
            gameinfo["Game Name"] = fields.get("GN");
        }
        */

        this.gameinfo = gameinfo;
    }

    compute_consts() {
        let review = document.getElementById("review");
        let size = parseInt(review.getAttribute("size"));
        let arrows = document.getElementById("arrows");

        // this is the number of "squares" across the board, including margins
        let n = size+1;
        this.width = parseInt(review.offsetWidth) * (n-2)/n;
        this.size = size;
        this.side = this.width/(this.size-1);
        this.pad = this.side;

        // this is not very elegant
        let w = this.width + this.pad*2;
        review.style.height = w + "px";
        arrows.style.width = w + "px";
    }


    recompute_consts() {
        this.compute_consts();
        this.board_graphics.recompute_consts();
    }

    dark_mode_toggle() {
        let dark = "#1A1A1A";
        let light = "#F5F5F5";
        let bg_color = dark;
        let fg_color = light;
        let old_class = "btn-light";
        let new_class = "btn-dark";
        let new_setting = true;
        let old_icon = "bi-moon-fill";
        let new_icon = "bi-sun-fill";
        let old_black_stone = "bi-circle-fill";
        let new_black_stone = "bi-circle";
        if (this.dark_mode) {
            bg_color = light;
            fg_color = dark;
            old_class = "btn-dark";
            new_class = "btn-light";
            new_setting = false;
            old_icon = "bi-sun-fill";
            new_icon = "bi-moon-fill";
        }

        // change the setting
        this.dark_mode = new_setting;

        // change the background
        document.body.style.background = bg_color;

        // change the buttons
        let buttons = document.querySelectorAll("button");
        for (let button of buttons) {
            let cls = button.getAttribute("class");
            let new_cls = cls.replace(old_class, new_class);
            button.setAttribute("class", new_cls);
        }

        // change the black and white stone icons
        let black_stone_icon = document.getElementsByClassName(old_black_stone)[0];
        let white_stone_icon = document.getElementsByClassName(new_black_stone)[0];
        black_stone_icon.setAttribute("class", new_black_stone);
        white_stone_icon.setAttribute("class", old_black_stone);

        // change modals
        /*
        for (let modal of document.getElementsByClassName("modal-content")) {
            modal.style.background = bg_color;
            modal.style.color = fg_color;
        }
        */

        // change modal close labels
        /*
        for (let close of document.getElementsByClassName("btn-close")) {
            if (this.dark_mode) {
                close.classList.add("btn-close-white");
            } else {
                close.classList.remove("btn-close-white");
            }
        }
        */

        /*
        for (let bar of document.querySelectorAll(".form-control, .form-select, #upload-textarea")) {
            if (this.dark_mode) {
                bar.style.backgroundColor = "#444444";
                bar.style.color = "#FFFFFF";
            } else {
                bar.style.backgroundColor = "#FFFFFF";
                bar.style.color = "#000000";
            }
        }
        */

        if (this.dark_mode) {
            document.documentElement.setAttribute("data-bs-theme", "dark");
        } else {
            document.documentElement.setAttribute("data-bs-theme", "light");
        }
    }

    comments_toggle() {
        if (this.comments.hidden()) {
            this.comments.show();
            this.tree_graphics.schedule_render();
        } else {
            this.comments.hide();
            this.tree_graphics.schedule_render();
        }
    }

    set_move_number(d) {
        let num = document.getElementById("move-number");
        num.innerHTML = d;
    }

    get_move_number() {
        let num = document.getElementById("move-number");
        return num.innerHTML;
    }

    set_black() {
        this.color = 1;
        this.toggling = false;
        this.mark = "";
        this.board_graphics.clear_ghosts();
    }

    set_white() {
        this.color = 2;
        this.toggling = false;
        this.mark = "";
        this.board_graphics.clear_ghosts();
    }

    set_toggle() {
        this.toggling = true;
        this.update_color();
        this.mark = "";
        this.board_graphics.clear_ghosts();
    }

    set_eraser() {
        this.mark = "eraser";
        this.board_graphics.clear_ghosts();
    }

    set_pen() {
        this.mark = "pen";
        this.board_graphics.clear_ghosts();
    }

    draw_pen(x0, y0, x1, y1, pen_color) {
        // draw it
        this.board_graphics.draw_pen(x0, y0, x1, y1, pen_color);

        // save in the sgf

        if (x0 == null) {
            x0 = -1.0;
        }
        if (y0 == null) {
            y0 = -1.0;
        }
        let digs = 4;
        this.pen.push([x0, y0, x1, y1, pen_color]);
    }

    apply_pen() {
        for (let [x0, y0, x1, y1, color] of this.pen) {
            if (x0 == -1.0) {
                x0 = null;
            }
            if (y0 == -1.0) {
                y0 = null;
            }
            this.board_graphics.draw_pen(x0, y0, x1, y1, color);
        }
    }

    erase_pen() {
        this.board_graphics.clear_pen();
    }

    update_color(optional_color) {
        // if not toggling, do nothing
        if (!this.toggling) {
            return;
        }

        // if we have a current move we can check, do the opposite
        if (this.current != null) {
            if (this.board.get(this.current) == 1) {
                this.color = 2;
            } else if (this.board.get(this.current) == 2) {
                this.color = 1;
            }
            return;
        }

        // if we're on the initial move, change to initial color
        // (black for normal games, white for handicap)
        if (parseInt(this.get_move_number()) == 0) {
            if (this.handicap) {
                this.color = 2;
            } else {
                this.color = 1;
            }
            return;
        }

        // if we have an optional_color thrown in, do the opposite
        if (optional_color == 1 || optional_color == 2) {
            this.color = opposite(optional_color);
        }
    }

   set_triangle() {
        this.mark = "triangle";
        this.board_graphics.clear_ghosts();
    }

    set_square() {
        this.mark = "square";
        this.board_graphics.clear_ghosts();
    }

    set_letter() {
        this.mark = "letter";
        this.custom_label = "";
        this.board_graphics.clear_ghosts();
    }

    set_number() {
        this.mark = "number";
        this.board_graphics.clear_ghosts();
    }

    trigger_score() {
        this.mark = "score";
        this.network_handler.prepare_score()
    }

    upload() {
        let inp = document.getElementById("upload-sgf");
        inp.onchange = () => {
            // hide the upload modal
            this.modals.hide_modal("upload-modal");

            if (inp.files.length == 0) {
                // i guess do nothing
                return;
            } else if (inp.files.length == 1) {
                // if 1 file, it's easy
                let f = inp.files[0];

                // want to maintain same behavior if user uploads file named the same
                inp.value = "";

                let reader = new FileReader();
                //reader.readAsText(f);
                reader.readAsArrayBuffer(f);

                reader.addEventListener(
                    "load",
                    () => {
                        // encode unicode, and encode with base64
                        this.network_handler.prepare_upload(b64_encode_arraybuffer(reader.result));
                    },
                    false,
                );
            } else {
                // max out at 10, just in case
                let max = 10;
                let i = 0;
                // if multiple files, build promises
                let promises = [];
                for (let f of inp.files) {
                    if (i >= max) {
                        break;
                    }
                    i++;
                    promises.push(
                        new Promise((resolve, reject) => {
                            let reader = new FileReader();
                            reader.readAsArrayBuffer(f);
                            reader.addEventListener(
                                "load",
                                () => resolve(reader.result),
                                false,
                            );
                        })
                    );
                }
                // want to maintain same behavior for the next pass
                inp.value = "";
    
                // turn list of promises into 1 promise
                Promise.all(promises)
                    .then((values) => {
                        let payload = [];
                        for (let v of values) {
                            payload.push(b64_encode_arraybuffer(v));
                        }

                        // encode unicode, and encode with base64
                        //this.network_handler.prepare_upload(b64_encode_unicode(sgf));
                        this.network_handler.prepare_upload(payload);
                    }
                    );
            }

        }
    }

    paste() {
        let textarea = document.getElementById("upload-textarea");
        // get the textarea value
        let value = textarea.value;
        // hide the upload modal
        textarea.value = "";

        // because of promises, i have to wait until the data from the
        // url is fetched before i hide the upload modal
        // otherwise the upload modal is hidden before the data returns
        // and an erroneous stone shows up on the board for a split
        // second before the new sgf is loaded
        if (value.startsWith("http")) {
            this.network_handler.prepare_request(value);
            setTimeout(() => this.modals.hide_modal("upload-modal"), 0);
        } else {
            this.network_handler.prepare_upload(b64_encode_unicode(value));
            // i swear...
            // even though the timeout is set to 0, if we take it out
            // then there is a problem where the "click" event closes the modal
            // first and then a stone appears on the board
            // so just... leave this alone, even if it looks weird
            setTimeout(() => this.modals.hide_modal("upload-modal"), 0);
        }
    }

    link_ogs_game() {
        let textarea = document.getElementById("ogs-textarea");
        // get the textarea value
        let value = textarea.value;
        textarea.value = "";

        this.network_handler.prepare_link_ogs_game(value);

        // hide the upload modal
        this.modals.hide_modal("upload-modal");
    }


    get_sgf_link() {
        let href = window.location.href;
        return href + "/sgf";
    }

    get_link() {
        return window.location.href;
    }

    copy(text) {
        navigator.clipboard.writeText(text);
    }

    download() {
        // stolen from stack overflow
        var element = document.createElement('a');
        let href = window.location.href;
        element.setAttribute("href", href + "/sgf");
        //element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(text));
        let basename = href.split("/").pop();
        element.setAttribute('download', basename + ".sgf");

        element.style.display = 'none';
        document.body.appendChild(element);

        element.click();

        document.body.removeChild(element);
    }

    set_black_caps(n) {
        document.getElementById("black-caps").innerHTML = " - " + n;
    }

    set_white_caps(n) {
        document.getElementById("white-caps").innerHTML = " - " + n;
    }

    handle_frame(frame) {
        // clear all marks
        this.marks = new Map();
        this.current = null;
        this.board_graphics.clear_current();
        this.board_graphics.clear_ghosts();
        this.pen = new Array();
        this.numbers = new Map();
        this.letters = new Array(26).fill(0);
        this.board_graphics.remove_marks();

        // always clear comments
        this.comments.clear();

        // TODO: there's some weirdness to think about here:
        // essentially, we only handle metadata on a full frame
        // so there are things that are done on a full frame that AREN'T
        // done on a diff (for example, updating color)
        if (frame.metadata != null) {
            this.handle_metadata(frame.metadata);
        }

        if (frame.type == FrameType.DIFF && frame.diff != null) {
            this.apply_diff(frame.diff);
        } else if (frame.type == FrameType.FULL && frame.diff != null) {
            this.full_frame(frame.diff);
        }

        if (frame.marks != null) {
            this.handle_marks(frame.marks);
        }

        if (frame.comments != null) {
            this.handle_comments(frame.comments);
        }

        //////// scoring logic

        if (frame.black_caps != null) {
            this.set_black_caps(frame.black_caps);
        }

        if (frame.white_caps != null) {
            this.set_white_caps(frame.white_caps);
        }

        if (frame.black_area != null) {
            for (let coord of frame.black_area) {
                this.board_graphics.draw_black_area(coord.x, coord.y);
                let id = coordtoid(coord);
                this.marks.set(id, "score-black");
            }
        }

        if (frame.white_area != null) {
            for (let coord of frame.white_area) {
                this.board_graphics.draw_white_area(coord.x, coord.y);
                let id = coordtoid(coord);
                this.marks.set(id, "score-white");
            }

        }

        if (frame.dame != null) {
            //console.log(frame.dame);
        }

        if (frame.tree != null) {
            // save it
            if (frame.tree.nodes != null) {
                if (frame.tree.root != 0) {
                    // this is the node to attach onto
                    let up = frame.tree.up;
                    let graft = this.tree_graphics.nodes[up];
                    // graft it on
                    graft.down.push(frame.tree.root);
                    for (let index in frame.tree.nodes) {
                        let node = frame.tree.nodes[index];
                        this.tree_graphics.nodes[index] = node;
                    }
                } else {
                    this.tree_graphics.nodes = frame.tree.nodes;
                }
            }
 
            this.tree_graphics.handle_tree(frame.tree);

            // current should always exist
            let current_index = frame.tree.current;
          
            // setting move number must happen before setting the color
            let current_node = this.tree_graphics.nodes[current_index];
            let move_number = current_node.depth;
            this.set_move_number(move_number);

            // updating the color checks the current move
            let current_color = current_node.color;
            this.update_color(current_color);
        }
    
        this.tree_graphics.set_scroll();
        this.tree_graphics.schedule_render();
    }

    handle_comments(comments) {
        for (let cmt of comments) {
            cmt = cmt.trim();
            for (let cmt_line of cmt.split("\n")) {
                this.comments.update(cmt_line);
            }
        }
    }

    handle_marks(marks) {
        if ("current" in marks && marks.current != null) {
            let coord = marks.current;

            this.board_graphics._draw_current(
                coord.x,
                coord.y,
                opposite(this.board.get(coord)));
            this.current = coord;
        }

        if ("squares" in marks && marks.squares != null) {
            let squares = marks.squares;
            for (let coord of squares) {
                this.place_square(coord);
            }
        }

        if ("triangles" in marks && marks.triangles != null) {
            let triangles = marks.triangles;
            for (let coord of triangles) {
                this.place_triangle(coord);
            }
        }

        if ("labels" in marks && marks.labels != null) {
            let labels = marks.labels;
            for (let lb of labels) {
                this.place_label(lb);
            }
        }

        if ("pens" in marks && marks.pens != null) {
            let pens = marks.pens;
            for (let pen of pens) {
                let x0 = pen.x0;
                let y0 = pen.y0;
                if (pen.x0 == -1) {
                    x0 = null;
                }
                if (pen.y0 == -1) {
                    y0 = null;
                }
                this.draw_pen(x0, y0, pen.x1, pen.y1, pen.color);
            }
        }
    }

    full_frame(frame) {

        this.board.clear();
        this.board_graphics.clear_and_remove();
        for (let a of frame.add) {
            let col = a["color"];
            let coords = a["coords"];
            for (let coord of coords) {
                this.place_stone(coord.x, coord.y, col);
            }
        }
    }

    handle_metadata(metadata) {
        if (metadata.size != null && metadata.size != this.size) {
            let review = document.getElementById("review");
            review.setAttribute("size", metadata.size);
            this.recompute_consts();
            this.board_graphics.reset_board();
            this.reset();
        }
        this.set_gameinfo(metadata.fields);
        this.modals.update_modals();
    }

    apply_diff(diff) {
        for (let a of diff.add) {
            let col = a["color"];
            let coords = a["coords"];
            for (let coord of coords) {
                this.place_stone(coord.x, coord.y, col);
                this.board.set(coord, col);
            }
        }
        for (let r of diff.remove) {
            let coords = r["coords"];
            for (let coord of coords) {
                this.remove_stone(coord.x, coord.y);
                this.board.set(coord, 0);
            }
        }
    }

    place_stone(x, y, color) {
        // if out of bounds, just return
        if (x < 0 || x >= this.size || y < 0 || y >= this.size) {
            return;
        }

        let coord = new Coord(x, y);
        this.board.set(coord, color);

        this.board_graphics.draw_stone(x, y, color);

    }

    place_triangle(coord) {
        let color = 1;
        if (this.board.get(coord) == 1) {
            color = 2;
        }
        this.board_graphics._draw_triangle(coord.x, coord.y, color);

        let id = coordtoid(coord);
        this.marks.set(id, "triangle");

    }

    place_square(coord) {
        let color = 1;
        if (this.board.get(coord) == 1) {
            color = 2;
        }
        this.board_graphics._draw_square(coord.x, coord.y, color);
        let id = coordtoid(coord);
        this.marks.set(id, "square");

    }

    place_label(lb) {
        // each lb has a coord and a text
        
        let coord = lb.coord
        let id = coordtoid(coord);

        let i = parseInt(lb.text);
        if (Number.isInteger(i)) {
            this.marks.set(id, "number:" + lb.text);
            this.numbers.set(i, 1);
            this.board_graphics._draw_manual_number(coord.x, coord.y, lb.text);
        } else {
            this.marks.set(id, "letter:" + lb.text);
            let letter_index = lb.text.charCodeAt(0)-65;
            this.letters[letter_index] = 1;
            this.board_graphics._draw_manual_letter(coord.x, coord.y, lb.text);
        }
    }

    place_label_string(label) {
        let coord = label.coord;
        let id = coordtoid(coord);
        this.marks.set(id, "label:" + label.text);
        this.board_graphics.draw_custom_label(coord.x, coord.y, label.text);
    }

    remove_mark(x, y) {
        this.board_graphics.remove_mark(x, y);
    }

    remove_stone(x, y) {
        let erased = this.board.remove(x, y);
        // if there was no stone there, do nothing
        if (!erased) {
            return;
        }

        this.board_graphics.erase_stone(x, y);
    }
}
