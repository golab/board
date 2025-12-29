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
    let id = "settings-modal";

    let title = document.createElement("h5");
    title.innerHTML = "Settings";

    let body = document.createElement("div");

    let gameinfo_toggle = document.createElement("a");
    gameinfo_toggle.setAttribute("data-bs-toggle", "collapse");
    gameinfo_toggle.setAttribute("href", "#gameinfo-collapse");
    gameinfo_toggle.setAttribute("aria-expanded", "false");
    gameinfo_toggle.classList.add("text-decoration-none");
    gameinfo_toggle.classList.add("text-reset");
    gameinfo_toggle.classList.add("fw-bold");
    let gameinfo_label = document.createElement("span");
    gameinfo_label.innerHTML = "Game Info ";
    let gameinfo_chevron = document.createElement("i");
    gameinfo_chevron.classList.add("bi");
    gameinfo_chevron.classList.add("bi-chevron-right");
    gameinfo_toggle.appendChild(gameinfo_label);
    gameinfo_toggle.appendChild(gameinfo_chevron);

    body.appendChild(gameinfo_toggle);
    body.appendChild(gameinfo_settings(state));
    body.appendChild(document.createElement("br"));

    let appearance_toggle = document.createElement("a");
    appearance_toggle.setAttribute("data-bs-toggle", "collapse");
    appearance_toggle.setAttribute("href", "#appearance-collapse");
    appearance_toggle.setAttribute("aria-expanded", "false");
    appearance_toggle.classList.add("text-decoration-none");
    appearance_toggle.classList.add("text-reset");
    appearance_toggle.classList.add("fw-bold");
    let appearance_label = document.createElement("span");
    appearance_label.innerHTML = "Appearance ";
    let appearance_chevron = document.createElement("i");
    appearance_chevron.classList.add("bi");
    appearance_chevron.classList.add("bi-chevron-right");
    appearance_toggle.appendChild(appearance_label);
    appearance_toggle.appendChild(appearance_chevron);

    body.appendChild(appearance_toggle);
    body.appendChild(appearance_settings(state));
    body.appendChild(document.createElement("br"));

    let behavior_toggle = document.createElement("a");
    behavior_toggle.setAttribute("data-bs-toggle", "collapse");
    behavior_toggle.setAttribute("href", "#behavior-collapse");
    behavior_toggle.setAttribute("aria-expanded", "false");
    behavior_toggle.classList.add("text-decoration-none");
    behavior_toggle.classList.add("text-reset");
    behavior_toggle.classList.add("fw-bold");
    let behavior_label = document.createElement("span");
    behavior_label.innerHTML = "Behavior ";
    let behavior_chevron = document.createElement("i");
    behavior_chevron.classList.add("bi");
    behavior_chevron.classList.add("bi-chevron-right");
    behavior_toggle.appendChild(behavior_label);
    behavior_toggle.appendChild(behavior_chevron);

    body.appendChild(behavior_toggle);
    body.appendChild(behavior_settings(state));
    body.appendChild(document.createElement("br"));

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
            let stars = document.getElementById(id + "-stars");
            let password_bar = document.getElementById(id + "-password-bar");
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

function gameinfo_settings(state) {
    let id = "settings-modal";
    let gameinfo_collapse = document.createElement("div");
    gameinfo_collapse.classList.add("collapse");
    gameinfo_collapse.id = "gameinfo-collapse";

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

    gameinfo_collapse.appendChild(nickname_div);
    gameinfo_collapse.appendChild(document.createElement("br"));

    // black player
    let black_player_div = document.createElement("div");
    let black_player_label = document.createElement("div");
    black_player_label.innerHTML = "Black";
    let black_player_bar = document.createElement("input");
    black_player_bar.setAttribute("class", "form-control");
    black_player_bar.id = id + "-blackplayer-bar";
    black_player_div.appendChild(black_player_label);
    black_player_div.appendChild(black_player_bar);
    gameinfo_collapse.appendChild(black_player_div)
    gameinfo_collapse.appendChild(document.createElement("br"));

    // white player
    let white_player_div = document.createElement("div");
    let white_player_label = document.createElement("div");
    white_player_label.innerHTML = "White";
    let white_player_bar = document.createElement("input");
    white_player_bar.setAttribute("class", "form-control");
    white_player_bar.id = id + "-whiteplayer-bar";
    white_player_div.appendChild(white_player_label);
    white_player_div.appendChild(white_player_bar);
    gameinfo_collapse.appendChild(white_player_div)
    gameinfo_collapse.appendChild(document.createElement("br"));

    // komi
    let komi_div = document.createElement("div");
    let komi_label = document.createElement("div");
    komi_label.innerHTML = "Komi";
    let komi_bar = document.createElement("input");
    komi_bar.setAttribute("class", "form-control");
    komi_bar.id = id + "-komi-bar";
    komi_div.appendChild(komi_label);
    komi_div.appendChild(komi_bar);
    gameinfo_collapse.appendChild(komi_div)
    gameinfo_collapse.appendChild(document.createElement("br"));

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

    gameinfo_collapse.appendChild(size_element);

    return gameinfo_collapse;
}

function appearance_settings(state) {
    let id = "settings-modal";
    let appearance_collapse = document.createElement("div");
    appearance_collapse.classList.add("collapse");
    appearance_collapse.id = "appearance-collapse";

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

    appearance_collapse.appendChild(darkmode_element);
    appearance_collapse.appendChild(document.createElement("br"));

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

    let saved_board_color_picker = "#f2bc74";
    if (stored_board_color.startsWith("#")) {
        saved_board_color_picker = stored_board_color;
    }
    let board_color_picker = document.createElement("input");
    board_color_picker.value = saved_board_color_picker;
    board_color_picker.setAttribute("type", "color");
    board_color_picker.setAttribute("id", "board-color-picker");
    board_color_picker.style.display = "none";
    board_color_picker.onchange = () => {
        saved_board_color_picker = board_color_picker.value;
        state.set_board_color(board_color_picker.value);
        localStorage.setItem("board_color", board_color_picker.value);
    };

    board_color_select.addEventListener(
        "change",
        () => {
            if (board_color_select.value == "[custom]") {
                state.set_board_color(saved_board_color_picker);
                localStorage.setItem("board_color", saved_board_color_picker);
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

    appearance_collapse.appendChild(board_color_element);
    appearance_collapse.appendChild(document.createElement("br"));

    // stone colors
    // black stones
    let black_stone_color = document.createElement("div");
    let black_stone_color_label = document.createElement("div");
    black_stone_color_label.innerHTML = "Black stones:";

    let black_stone_color_select = document.createElement("select");
    black_stone_color_select.setAttribute("id", id + "-black-stone-color-select");
    black_stone_color_select.setAttribute("class", "form-select");
    for (let x of ["gradient", "[custom]"]) {
        let opt = document.createElement("option");
        opt.setAttribute("value", x);
        opt.innerHTML = x;
        black_stone_color_select.appendChild(opt);
    }

    let black_stored = localStorage.getItem("black_stone_color") || "gradient";
    let saved_b_color_picker = "#000000";
    if (black_stored.startsWith("#")) {
        saved_b_color_picker = black_stored;
    }
    let b_color_picker = document.createElement("input");
    b_color_picker.value = saved_b_color_picker;
    b_color_picker.setAttribute("type", "color");
    b_color_picker.setAttribute("id", "board-color-picker");
    b_color_picker.style.display = "none";
    b_color_picker.onchange = () => {
        saved_b_color_picker = b_color_picker.value;
        state.set_black_stone_color(b_color_picker.value);
        localStorage.setItem("black_stone_color", b_color_picker.value);
    };
    //b_color_picker.value = black_stored;
    black_stone_color_select.value = black_stored;
    state.set_black_stone_color(black_stored);
    black_stone_color.appendChild(black_stone_color_label);
    black_stone_color.appendChild(b_color_picker);

    black_stone_color_select.addEventListener(
        "change",
        () => {
            if (black_stone_color_select.value == "[custom]") {
                state.set_black_stone_color(saved_b_color_picker);
                localStorage.setItem("black_stone_color", saved_b_color_picker);
                b_color_picker.click();
                black_stone_color_select.value = "";
            } else {
                state.set_black_stone_color(black_stone_color_select.value);
                localStorage.setItem("black_stone_color", black_stone_color_select.value);
            }
        }
    );

    appearance_collapse.appendChild(black_stone_color);
    appearance_collapse.appendChild(black_stone_color_select);
    appearance_collapse.appendChild(document.createElement("br"));

    // white stones
    let white_stone_color = document.createElement("div");
    let white_stone_color_label = document.createElement("div");
    white_stone_color_label.innerHTML = "White stones:";

    let white_stone_color_select = document.createElement("select");
    white_stone_color_select.setAttribute("id", id + "-white-stone-color-select");
    white_stone_color_select.setAttribute("class", "form-select");
    for (let x of ["pattern", "gradient", "[custom]"]) {
        let opt = document.createElement("option");
        opt.setAttribute("value", x);
        opt.innerHTML = x;
        white_stone_color_select.appendChild(opt);
    }

    let white_stored = localStorage.getItem("white_stone_color") || "pattern";
    let saved_w_color_picker = "#FFFFFF";
    if (white_stored.startsWith("#")) {
        saved_w_color_picker = white_stored;
    }
    let w_color_picker = document.createElement("input");
    w_color_picker.value = saved_w_color_picker;
    w_color_picker.setAttribute("type", "color");
    w_color_picker.setAttribute("id", "board-color-picker");
    w_color_picker.style.display = "none";
    w_color_picker.onchange = () => {
        saved_w_color_picker = w_color_picker.value;
        state.set_white_stone_color(w_color_picker.value);
        localStorage.setItem("white_stone_color", w_color_picker.value);
    };
    //w_color_picker.value = white_stored;
    white_stone_color_select.value = white_stored;
    state.set_white_stone_color(white_stored);
    white_stone_color.appendChild(white_stone_color_label);
    white_stone_color.appendChild(w_color_picker);

    white_stone_color_select.addEventListener(
        "change",
        () => {
            if (white_stone_color_select.value == "[custom]") {
                state.set_white_stone_color(saved_w_color_picker);
                localStorage.setItem("white_stone_color", saved_w_color_picker);
                w_color_picker.click();
                white_stone_color_select.value = "";
            } else {
                state.set_white_stone_color(white_stone_color_select.value);
                localStorage.setItem("white_stone_color", white_stone_color_select.value);
            }
        }
    );

    appearance_collapse.appendChild(white_stone_color);
    appearance_collapse.appendChild(white_stone_color_select);
    appearance_collapse.appendChild(document.createElement("br"));


    // black stone outline
    let black_outline_color = document.createElement("div");
    let black_outline_color_label = document.createElement("div");
    black_outline_color_label.innerHTML = "Black stone border:";

    let black_outline_color_select = document.createElement("select");
    black_outline_color_select.setAttribute("id", id + "-black-outline-color-select");
    black_outline_color_select.setAttribute("class", "form-select");
    for (let x of ["none", "[custom]"]) {
        let opt = document.createElement("option");
        opt.setAttribute("value", x);
        opt.innerHTML = x;
        black_outline_color_select.appendChild(opt);
    }

    let black_outline_stored = localStorage.getItem("black_outline_color") || "none";
    let saved_b_outline_color_picker = "#FFFFFF";
    if (black_outline_stored.startsWith("#")) {
        saved_b_outline_color_picker = black_outline_stored;
    }
    let b_outline_color_picker = document.createElement("input");
    b_outline_color_picker.value = saved_b_outline_color_picker;
    b_outline_color_picker.setAttribute("type", "color");
    b_outline_color_picker.setAttribute("id", "board-color-picker");
    b_outline_color_picker.style.display = "none";
    b_outline_color_picker.onchange = () => {
        saved_b_outline_color_picker = b_outline_color_picker.value;
        state.set_black_outline_color(b_outline_color_picker.value);
        localStorage.setItem("black_outline_color", b_outline_color_picker.value);
    };
    black_outline_color_select.value = black_outline_stored;
    state.set_black_outline_color(black_outline_stored);
    black_outline_color.appendChild(black_outline_color_label);
    black_outline_color.appendChild(b_outline_color_picker);

    black_outline_color_select.addEventListener(
        "change",
        () => {
            if (black_outline_color_select.value == "[custom]") {
                state.set_black_outline_color(saved_b_outline_color_picker);
                localStorage.setItem("black_outline_color", saved_b_outline_color_picker);
                b_outline_color_picker.click();
                black_outline_color_select.value = "";
            } else {
                state.set_black_outline_color(black_outline_color_select.value);
                localStorage.setItem("black_outline_color", black_outline_color_select.value);
            }
        }
    );

    appearance_collapse.appendChild(black_outline_color);
    appearance_collapse.appendChild(black_outline_color_select);
    appearance_collapse.appendChild(document.createElement("br"));

    // white stone outline
    let white_outline_color = document.createElement("div");
    let white_outline_color_label = document.createElement("div");
    white_outline_color_label.innerHTML = "White stone border:";

    let white_outline_color_select = document.createElement("select");
    white_outline_color_select.setAttribute("id", id + "-white-outline-color-select");
    white_outline_color_select.setAttribute("class", "form-select");
    for (let x of ["none", "[custom]"]) {
        let opt = document.createElement("option");
        opt.setAttribute("value", x);
        opt.innerHTML = x;
        white_outline_color_select.appendChild(opt);
    }

    let white_outline_stored = localStorage.getItem("white_outline_color") || "none";
    let saved_w_outline_color_picker = "#000000";
    if (white_outline_stored.startsWith("#")) {
        saved_w_outline_color_picker = white_outline_stored;
    }
    let w_outline_color_picker = document.createElement("input");
    w_outline_color_picker.value = saved_w_outline_color_picker;
    w_outline_color_picker.setAttribute("type", "color");
    w_outline_color_picker.setAttribute("id", "board-color-picker");
    w_outline_color_picker.style.display = "none";
    w_outline_color_picker.onchange = () => {
        saved_w_outline_color_picker = w_outline_color_picker.value;
        state.set_white_outline_color(w_outline_color_picker.value);
        localStorage.setItem("white_outline_color", w_outline_color_picker.value);
    };
    white_outline_color_select.value = white_outline_stored;
    state.set_white_outline_color(white_outline_stored);
    white_outline_color.appendChild(white_outline_color_label);
    white_outline_color.appendChild(w_outline_color_picker);

    white_outline_color_select.addEventListener(
        "change",
        () => {
            if (white_outline_color_select.value == "[custom]") {
                state.set_white_outline_color(saved_w_outline_color_picker);
                localStorage.setItem("white_outline_color", saved_w_outline_color_picker);
                w_outline_color_picker.click();
                white_outline_color_select.value = "";
            } else {
                state.set_white_outline_color(white_outline_color_select.value);
                localStorage.setItem("white_outline_color", white_outline_color_select.value);
            }
        }
    );

    appearance_collapse.appendChild(white_outline_color);
    appearance_collapse.appendChild(white_outline_color_select);
    appearance_collapse.appendChild(document.createElement("br"));

    // shadows
    let shadow_element = document.createElement("div");
    shadow_element.innerHTML = "Shadows "

    let d_shadow = document.createElement("div");
    d_shadow.setAttribute("class", "form-check form-switch");

    let inp_shadow = document.createElement("input");
    inp_shadow.setAttribute("class", "form-check-input");
    inp_shadow.setAttribute("type", "checkbox");
    inp_shadow.setAttribute("role", "switch");
    inp_shadow.checked = true;
    inp_shadow.setAttribute("id", "textured-stones-switch");

    let label_shadow = document.createElement("label");
    label_shadow.setAttribute("class", "form-check-label");
    label_shadow.setAttribute("for", "textured-stones-switch");
    label_shadow.innerHTML = "On";
    let stored_shadow = localStorage.getItem("shadow") || "";
    if (stored_shadow == "true") {
        inp_shadow.checked = true;
        state.set_shadow(true);
    } else if (stored_shadow == "false") {
        inp_shadow.checked = false;
        state.set_shadow(false);
        label_shadow.innerHTML = "Off";
    }

    inp_shadow.onchange = function() {
        if (this.checked) {
            label_shadow.innerHTML = "On";
            state.set_shadow(true);
            localStorage.setItem("shadow", "true");
        } else {
            label_shadow.innerHTML = "Off";
            state.set_shadow(false);
            localStorage.setItem("shadow", "false");
        }
    };

    d_shadow.appendChild(inp_shadow);
    d_shadow.appendChild(label_shadow);
    shadow_element.append(d_shadow);

    appearance_collapse.appendChild(shadow_element);
    return appearance_collapse;
}

function behavior_settings(state) {
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
    let behavior_collapse = document.createElement("div");
    behavior_collapse.classList.add("collapse");
    behavior_collapse.id = "behavior-collapse";

    // up/down toggle
    let updown_element = document.createElement("div");
    updown_element.innerHTML = "Up/Down Keys "

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
            localStorage.setItem("updown", "true");
        } else {
            label.innerHTML = "KGS";
            state.branch_jump = false;
            localStorage.setItem("updown", "false");
        }
    };

    let stored_updown = localStorage.getItem("updown") || "";
    if (stored_updown == "true") {
        label.innerHTML = "OGS";
        inp.checked = true;
        state.branch_jump = true;
    } else if (stored_updown == "false") {
        label.innerHTML = "KGS";
        inp.checked = false;
        state.branch_jump = false;
    }

    d.appendChild(inp);
    d.appendChild(label);
    updown_element.append(d);

    behavior_collapse.appendChild(updown_element);
    behavior_collapse.appendChild(document.createElement("br"));


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
    password_anchor.innerHTML = "Add Password";
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

    behavior_collapse.appendChild(password_anchor);

    behavior_collapse.appendChild(password_bar);
    behavior_collapse.appendChild(setcancel_div);

    behavior_collapse.appendChild(stars);

    behavior_collapse.appendChild(showremove_div);

    behavior_collapse.appendChild(document.createElement("br"));
    behavior_collapse.appendChild(document.createElement("br"));

    // input buffer

    let buffer_element = document.createElement("div");
    buffer_element.innerHTML = "Input Buffer "

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

    behavior_collapse.appendChild(buffer_element);

    return behavior_collapse;
}
