/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_text_button, prefer_dark_mode, add_tooltip } from '../common.js';
import { add_modal } from './base.js';
export {
    add_settings_modal,
    set_password,
    remove_password
}

function remove_password() {
    let id = "settings-modal";

    let showremove_div = document.getElementById(id + "-showremove-div");
    showremove_div.hidden = true;

    let stars = document.getElementById(id + "-stars");
    stars.hidden = true;

    let password_bar = document.getElementById(id + "-password-bar");
    password_bar.value = "";

    let password_anchor = document.getElementById(id + "-password-anchor");
    password_anchor.hidden = false;
}

function set_password(password) {
    let id = "settings-modal";

    let password_anchor = document.getElementById(id + "-password-anchor");
    password_anchor.hidden = true;

    let stars = document.getElementById(id + "-stars");
    stars.innerHTML = "Password: " + "*".repeat(password.length);
    stars.hidden = false;

    let showremove_div = document.getElementById(id + "-showremove-div");
    showremove_div.hidden = false;

    let setcancel_div = document.getElementById(id + "-setcancel-div");
    setcancel_div.hidden = true;

    let password_bar = document.getElementById(id + "-password-bar");
    password_bar.hidden = true;
    password_bar.value = password;
}

function add_settings_modal(state) {
    function enter_password(password) {
        let id = "settings-modal";
        let m = bootstrap.Modal.getInstance("#"+id);
        m._element.focus();

        if (password == "") {
            cancel_password();
        } else {
            state.network_handler.prepare_settings(make_settings());
        }
    }

    function cancel_password() {
        let id = "settings-modal";
        let password_bar = document.getElementById(id + "-password-bar");
        let setcancel_div = document.getElementById(id + "-setcancel-div");
        let password_anchor = document.getElementById(id + "-password-anchor");

        password_bar.hidden = true;
        setcancel_div.hidden = true;
        password_anchor.hidden = false;

        let m = bootstrap.Modal.getInstance("#"+id);
        m._element.focus();
    }

    let id = "settings-modal";

    let title = document.createElement("h5");
    title.innerHTML = "Settings";

    let body = document.createElement("div");

    // nickname
    let nickname_div = document.createElement("div");
    let nickname_label = document.createElement("div");
    nickname_label.innerHTML = "Nickname";

    let nickname_bar = document.createElement("input")
    nickname_bar.setAttribute("class", "form-control");
    nickname_bar.id = id + "-nickname-bar";

    nickname_bar.value = localStorage.getItem("nickname") || "";

    nickname_bar.addEventListener(
        "keypress",
        (event) => {
            if (event.key == "Enter") {
                let m = bootstrap.Modal.getInstance("#"+id);
                m.hide();

                let nickname = nickname_bar.value;
                localStorage.setItem("nickname", nickname);
                state.network_handler.prepare_nickname(nickname);
            }
        }
    );
    nickname_div.appendChild(nickname_label);
    nickname_div.appendChild(nickname_bar);

    body.appendChild(nickname_div);
    body.appendChild(document.createElement("br"));

    // black player
    let black_player_div = document.createElement("div");
    let black_player_label = document.createElement("div");
    black_player_label.innerHTML = "Black";
    let black_player_bar = document.createElement("input");
    black_player_bar.setAttribute("class", "form-control");
    black_player_bar.id = id + "-blackplayer-bar";
    black_player_div.appendChild(black_player_label);
    black_player_div.appendChild(black_player_bar);
    body.appendChild(black_player_div)
    body.appendChild(document.createElement("br"));

    // white player
    let white_player_div = document.createElement("div");
    let white_player_label = document.createElement("div");
    white_player_label.innerHTML = "White";
    let white_player_bar = document.createElement("input");
    white_player_bar.setAttribute("class", "form-control");
    white_player_bar.id = id + "-whiteplayer-bar";
    white_player_div.appendChild(white_player_label);
    white_player_div.appendChild(white_player_bar);
    body.appendChild(white_player_div)
    body.appendChild(document.createElement("br"));

    // komi
    let komi_div = document.createElement("div");
    let komi_label = document.createElement("div");
    komi_label.innerHTML = "Komi";
    let komi_bar = document.createElement("input");
    komi_bar.setAttribute("class", "form-control");
    komi_bar.id = id + "-komi-bar";
    komi_div.appendChild(komi_label);
    komi_div.appendChild(komi_bar);
    body.appendChild(komi_div)
    body.appendChild(document.createElement("br"));

    // board size
    let size_element = document.createElement("div");

    let size_label = document.createElement("label");
    size_label.innerHTML = "Board size:&nbsp;";
    size_label.setAttribute("for", id + "-size-select");
    let size_select = document.createElement("select");
    size_select.setAttribute("id", id + "-size-select");
    size_select.setAttribute("class", "form-select");
    for (let x of ["9", "13", "19"]) {
        let opt = document.createElement("option");
        opt.setAttribute("value", x);
        opt.innerHTML = x;
        size_select.appendChild(opt);
    }
    size_select.value = state.size;

    size_element.appendChild(size_label);
    size_element.appendChild(size_select);

    body.appendChild(size_element);
    body.appendChild(document.createElement("br"));

    // up/down toggle
    let updown_element = document.createElement("div");
    updown_element.innerHTML = "Up/Down key behavior "

    let updown_anchor = document.createElement("a");
    updown_anchor.setAttribute("href", "/about#updown");
    updown_anchor.setAttribute("target", "_blank");

    let updown_icon = document.createElement("i");
    updown_icon.setAttribute("class", "bi-question-circle");
    add_tooltip(updown_icon, "This controls whether up/down moves up branches (like OGS) or modifies which branch to follow (like KGS)");
    updown_anchor.appendChild(updown_icon);
    updown_element.appendChild(updown_anchor);

    let d = document.createElement("div");
    d.setAttribute("class", "form-check form-switch");

    let inp = document.createElement("input");
    inp.setAttribute("class", "form-check-input");
    inp.setAttribute("type", "checkbox");
    inp.setAttribute("role", "switch");
    inp.checked = true;
    inp.setAttribute("id", "updown-switch");

    let label = document.createElement("label");
    label.setAttribute("class", "form-check-label");
    label.setAttribute("for", "updown-switch");
    label.innerHTML = "OGS";

    inp.onchange = function() {
        if (this.checked) {
            label.innerHTML = "OGS";
            state.branch_jump = true;
        } else {
            label.innerHTML = "KGS";
            state.branch_jump = false;
        }
    };

    d.appendChild(inp);
    d.appendChild(label);
    updown_element.append(d);

    body.appendChild(updown_element);
    body.appendChild(document.createElement("br"));

    // lightmode/darkmode toggle
    let darkmode_element = document.createElement("div");
    darkmode_element.innerHTML = "Darkmode"

    let d2 = document.createElement("div");
    d2.setAttribute("class", "form-check form-switch");

    let inp2 = document.createElement("input");
    inp2.setAttribute("class", "form-check-input");
    inp2.setAttribute("type", "checkbox");
    inp2.setAttribute("role", "switch");
    inp2.checked = prefer_dark_mode();
    inp2.setAttribute("id", "darkmode-switch");

    let label2 = document.createElement("label");
    label2.setAttribute("class", "form-check-label");
    label2.setAttribute("for", "darkmode-switch");
    let darkmode_icon = document.createElement("i");
    darkmode_icon.setAttribute("class", "bi-moon-fill");
    let lightmode_icon = document.createElement("i");
    lightmode_icon.setAttribute("class", "bi-sun-fill");
    label2.appendChild(darkmode_icon);

    inp2.onchange = function() {
        label2.innerHTML = "";
        if (this.checked) {
            label2.appendChild(darkmode_icon);
        } else {
            label2.appendChild(lightmode_icon);
        }
        state.dark_mode_toggle();
    };

    d2.appendChild(inp2);
    d2.appendChild(label2)
    darkmode_element.append(d2);

    body.appendChild(darkmode_element);
    body.appendChild(document.createElement("br"));

    // board color
    let board_color_element = document.createElement("div");

    let board_color_label = document.createElement("label");
    board_color_label.innerHTML = "Board color:&nbsp;";
    board_color_label.setAttribute("for", id + "-board-color-select");
    let board_color_select = document.createElement("select");
    board_color_select.setAttribute("id", id + "-board-color-select");
    board_color_select.setAttribute("class", "form-select");
    for (let x of ["light", "medium", "dark", "[custom]"]) {
        let opt = document.createElement("option");
        opt.setAttribute("value", x);
        opt.innerHTML = x;
        board_color_select.appendChild(opt);
    }
    let stored_board_color = localStorage.getItem("board_color") || "";
    if (stored_board_color != "") {
        state.board_color = stored_board_color;
        board_color_select.value = "";
    }
    board_color_select.value = state.board_color;

    let board_color_picker = document.createElement("input");
    board_color_picker.setAttribute("type", "color");
    board_color_picker.setAttribute("id", "board-color-picker");
    board_color_picker.style.display = "none";
    board_color_picker.onchange = () => {
        state.set_board_color(board_color_picker.value);
        localStorage.setItem("board_color", board_color_picker.value);
    };

    board_color_select.addEventListener(
        "change",
        () => {
            if (board_color_select.value == "[custom]") {
                board_color_picker.click();
                board_color_select.value = "";
            } else {
                state.set_board_color(board_color_select.value);
                localStorage.setItem("board_color", board_color_select.value);
            }
        }
    );

    board_color_element.appendChild(board_color_label);
    board_color_element.appendChild(board_color_select);

    body.appendChild(board_color_element);
    body.appendChild(document.createElement("br"));

    // textured stones
    let textured_stones_element = document.createElement("div");
    textured_stones_element.innerHTML = "Textured stones "

    let d_textured_stones = document.createElement("div");
    d_textured_stones.setAttribute("class", "form-check form-switch");

    let inp_textured_stones = document.createElement("input");
    inp_textured_stones.setAttribute("class", "form-check-input");
    inp_textured_stones.setAttribute("type", "checkbox");
    inp_textured_stones.setAttribute("role", "switch");
    inp_textured_stones.checked = true;
    inp_textured_stones.setAttribute("id", "textured-stones-switch");

    let label_textured_stones = document.createElement("label");
    label_textured_stones.setAttribute("class", "form-check-label");
    label_textured_stones.setAttribute("for", "textured-stones-switch");
    label_textured_stones.innerHTML = "On";
    let stored_textured_stones = localStorage.getItem("textured_stones") || "";
    if (stored_textured_stones == "true") {
        inp_textured_stones.checked = true;
        state.set_textured_stones(true);
    } else if (stored_textured_stones == "false") {
        inp_textured_stones.checked = false;
        state.set_textured_stones(false);
        label_textured_stones.innerHTML = "Off";
    }

    inp_textured_stones.onchange = function() {
        if (this.checked) {
            label_textured_stones.innerHTML = "On";
            state.set_textured_stones(true);
            localStorage.setItem("textured_stones", "true");
        } else {
            label_textured_stones.innerHTML = "Off";
            state.set_textured_stones(false);
            localStorage.setItem("textured_stones", "false");
        }
    };

    d_textured_stones.appendChild(inp_textured_stones);
    d_textured_stones.appendChild(label_textured_stones);
    textured_stones_element.append(d_textured_stones);

    body.appendChild(textured_stones_element);
    body.appendChild(document.createElement("br"));

    // password
    let password_anchor = document.createElement("a");
    let password_bar = document.createElement("input");
    password_bar.setAttribute("class", "form-control");

    let setcancel_div = document.createElement("div");
    setcancel_div.id = id + "-setcancel-div";

    let set_button = new_text_button(
        "Set",
        () => {
            enter_password(password_bar.value);
        }
    );

    let cancel_button = new_text_button(
        "Cancel",
        () => {
            cancel_password();
        }
    );

    cancel_button.setAttribute("class", "btn btn-secondary");
    set_button.setAttribute("class", "btn btn-primary");

    setcancel_div.appendChild(cancel_button);
    let nbsp = document.createElement("span");
    nbsp.innerHTML = "&nbsp;";
    setcancel_div.appendChild(nbsp);
    setcancel_div.appendChild(set_button);
    setcancel_div.hidden = true;

    password_anchor.style.cursor = "pointer";
    password_anchor.id = id + "-password-anchor";
    password_anchor.innerHTML = "Add password";
    password_anchor.onclick = () => {
        password_bar.hidden = false;
        password_bar.focus();
        setcancel_div.hidden = false;
        password_anchor.hidden = true;
    };

    let stars = document.createElement("div");
    stars.hidden = true;
    stars.id = id + "-stars";

    password_bar.hidden = true;
    password_bar.id = id + "-password-bar";
    password_bar.addEventListener("keypress", (event) => {
        if (event.key == "Enter") {
            enter_password(password_bar.value);
        }
    });

    let showremove_div = document.createElement("div");
    showremove_div.id = id + "-showremove-div";

    let remove_button = new_text_button(
        "Remove",
        () => {
            remove_password();
            state.network_handler.prepare_settings(make_settings());
        }
    );

    remove_button.setAttribute("class", "btn btn-primary");

    let show_button = new_text_button(
        "Show",
        () => {
            stars.textContent = "Password: " + password_bar.value;
        }
    );

    show_button.setAttribute("class", "btn btn-primary");

    showremove_div.appendChild(show_button);
    nbsp = document.createElement("span");
    nbsp.innerHTML = "&nbsp;";
    showremove_div.appendChild(nbsp);
    showremove_div.appendChild(remove_button);
    showremove_div.hidden = true;

    body.appendChild(password_anchor);

    body.appendChild(password_bar);
    body.appendChild(setcancel_div);

    body.appendChild(stars);

    body.appendChild(showremove_div);

    body.appendChild(document.createElement("br"));
    body.appendChild(document.createElement("br"));

    // input buffer

    let buffer_element = document.createElement("div");
    buffer_element.innerHTML = "Input buffer "

    let anchor = document.createElement("a");
    anchor.setAttribute("href", "/about#input-buffer");
    anchor.setAttribute("target", "_blank");

    let icon = document.createElement("i");
    icon.setAttribute("class", "bi-question-circle");
    add_tooltip(icon, "After a user changes the board in some way, the server will ignore everyone else for a given period of time");
    anchor.appendChild(icon);

    buffer_element.appendChild(anchor);

    let range = document.createElement("input");
    range.setAttribute("type", "range");
    range.setAttribute("class", "form-range");
    range.setAttribute("id", id + "-bufferrange");
    range.min = 0;
    range.max = 1000;
    range.step = 25;
    range.value = state.input_buffer;

    let output = document.createElement("output");
    output.setAttribute("id", id + "-bufferoutput");
    output.value = state.input_buffer + "ms";
    range.oninput = function() {output.value = this.value + "ms"};

    buffer_element.appendChild(range);
    buffer_element.appendChild(output);

    body.appendChild(buffer_element);

    let settings_modal = add_modal(
        id,
        title,
        body,
        true,
        () => state.network_handler.prepare_settings(make_settings())
    );
    settings_modal.addEventListener(
        'hidden.bs.modal',
        () => {
            stars.innerHTML = "Password: " + "*".repeat(password_bar.value.length);
        }
    );
    let m = new bootstrap.Modal(settings_modal);
    return settings_modal;
}

function make_settings() {
    let id = "settings-modal";
    let buffer = parseInt(document.getElementById(id + "-bufferrange").value);
    let size = parseInt(document.getElementById(id + "-size-select").value);
    let password = document.getElementById(id + "-password-bar").value;
    let nickname = document.getElementById(id + "-nickname-bar").value;

    let black = document.getElementById(id + "-blackplayer-bar").value;
    let white = document.getElementById(id + "-whiteplayer-bar").value;
    let komi = document.getElementById(id + "-komi-bar").value;
    return {
        "buffer": buffer,
        "size": size,
        "password": password,
        "nickname": nickname,
        "black": black,
        "white": white,
        "komi": komi};
}
