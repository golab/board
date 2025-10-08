/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_text_button, new_icon_button, add_tooltip, prefer_dark_mode } from './common.js';

function make_settings() {
    let id = "settings-modal";
    let buffer = parseInt(document.getElementById(id + "-bufferrange").value);
    let size = parseInt(document.getElementById(id + "-size-select").value);
    let password = document.getElementById(id + "-password-bar").value;
    let nickname = document.getElementById(id + "-nickname-bar").value;
    return {"buffer": buffer, "size": size, "password": password, "nickname": nickname};
}

function get_nickname() {
    let id = "settings-modal";
    return document.getElementById(id + "-nickname-bar").value;
}

export function create_modals(_state) {
    // apparently necessary to make the tooltips work properly
    const state = _state;
    const modal_ids = [];
    const modals_up = new Map();

    // don't forget to add new modals here
    // otherwise clicking the button throws an error
    add_trash_modal();
    add_scissors_modal();
    add_gameinfo_modal();
    add_download_modal();
    add_upload_modal();
    add_error_modal();
    add_prompt_modal();
    add_info_modal();
    add_settings_modal();
    add_users_modal();
    enable_tooltips();
    add_toast();

    function add_toast() {
        let toast_div = document.createElement("div");
        toast_div.id = "toasts";
        toast_div.setAttribute(
            "class",
            "toast-container position-fixed top-0 end-0 p-3"
        );
        document.body.appendChild(toast_div);
    }

    function show_toast(body_text) {
        let container = document.getElementById("toasts");
    
        let toast = document.createElement("div");

        toast.setAttribute("class", "toast");
        toast.setAttribute("data-bs-autohide", false);
        toast.role = "alert";

        let flex = document.createElement("div")
        flex.setAttribute("class", "d-flex");

        let body = document.createElement("div");
        body.setAttribute("class", "toast-body");
        body.innerHTML = body_text;
        flex.appendChild(body);
    
        let close = document.createElement("button");
        close.setAttribute("class", "btn-close me-2 m-auto");
        close.setAttribute("data-bs-dismiss", "toast");
        close.setAttribute("aria-label", "Close");
        flex.appendChild(close);

        toast.appendChild(flex);

        container.appendChild(toast);
    
        let t = bootstrap.Toast.getOrCreateInstance(toast);

        t.show();
    }

    function get_prompt_bar() {
        let id = "prompt-modal-input";
        let bar = document.getElementById(id);
        return bar.value;
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

    function enter_password(password) {
        let id = "settings-modal";
        let m = bootstrap.Modal.getInstance("#"+id);
        m._element.focus();

        if (password == "") {
            cancel_password();
        } else {
            //set_password(password);
            state.network_handler.prepare_settings(make_settings());
        }
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

    function enable_tooltips() {
        var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
        var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {return new bootstrap.Tooltip(tooltipTriggerEl)});
    }
 
    function update_modals() {
        update_gameinfo_modal();
        update_settings_modal();
    }

    function add_settings_modal() {
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
        nickname_bar.addEventListener(
            "keypress",
            (event) => {
                if (event.key == "Enter") {
                    hide_modal(id);
                    let nickname = nickname_bar.value;
                    state.network_handler.prepare_nickname(nickname);
                }
            }
        );
        nickname_div.appendChild(nickname_label);
        nickname_div.appendChild(nickname_bar);

        body.appendChild(nickname_div);
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
        d.appendChild(label)
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
                modals_up.delete(id);
                stars.innerHTML = "Password: " + "*".repeat(password_bar.value.length);
            }
        );
        modal_ids.push(id);
        let m = new bootstrap.Modal(settings_modal);
    }

    function update_settings_modal() {
        let id = "settings-modal";

        let range = document.getElementById(id + "-bufferrange");
        range.value = state.input_buffer;
        let output = document.getElementById(id + "-bufferoutput");
        output.value = state.input_buffer + "ms";
        
        let select = document.getElementById(id + "-size-select");
        select.value = state.size;

        let password = state.password;
        if (password == "") {
            remove_password();
        } else {
            set_password(password);
        }
    }

    function add_error_modal() {
        let id = "error-modal";
        let paragraph = document.createElement("p");
        paragraph.setAttribute("id", id + "-paragraph");
        let button = new_icon_button("bi-exclamation-triangle-fill");
        button.setAttribute("class", "btn btn-danger");
        paragraph.appendChild(button);
        let span1 = document.createElement("span");
        span1.setAttribute("id", id + "-message");
        paragraph.appendChild(span1);

        let title = document.createElement("h5");
        title.innerHTML = "Error";

        let error_modal = add_modal(id, title, paragraph, false);
        error_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        error_modal.addEventListener('shown.bs.modal', () => modals_up.set(id, true));
        modal_ids.push(id);

        // here i have to manually create the bootstrap modal object
        // because there isn't anywhere else that automatically triggers
        // its creation (not just the document element error_modal)
        // but later i want to show it
        let m = new bootstrap.Modal(error_modal);
    }

    function add_prompt_modal() {
        let id = "prompt-modal";
        let paragraph = document.createElement("p");
        paragraph.setAttribute("id", id + "-paragraph");
        let button = new_icon_button("bi-info-square");
        button.setAttribute("class", "btn btn-primary");
        paragraph.appendChild(button);
        let span1 = document.createElement("span");
        span1.setAttribute("id", id + "-message");
        paragraph.appendChild(span1);

        paragraph.appendChild(document.createElement("br"));
        paragraph.appendChild(document.createElement("br"));

        let prompt_bar = document.createElement("input");
        prompt_bar.setAttribute("class", "form-control");
        prompt_bar.id = id + "-input";
        prompt_bar.addEventListener("keypress", (event) => {
            if (event.key == "Enter") {
                let button = document.getElementById(id + "-ok");
                button.click();
            }
        });

        paragraph.appendChild(prompt_bar);

        let title = document.createElement("h5");
        title.innerHTML = "Prompt";

        let prompt_modal = add_modal(id, title, paragraph, true);
        prompt_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        prompt_modal.addEventListener(
            'shown.bs.modal',
            () => {
                modals_up.set(id, true);
                prompt_bar.focus();
            }
        );
        modal_ids.push(id);

        let m = new bootstrap.Modal(prompt_modal);
    }

    function add_info_modal() {
        let id = "info-modal";
        let paragraph = document.createElement("p");
        paragraph.setAttribute("id", id + "-paragraph");
        let button = new_icon_button("bi-info-square");
        button.setAttribute("class", "btn btn-primary");
        paragraph.appendChild(button);
        let span1 = document.createElement("span");
        span1.setAttribute("id", id + "-message");
        paragraph.appendChild(span1);

        let title = document.createElement("h5");
        title.innerHTML = "Info";

        let info_modal = add_modal(id, title, paragraph, false);
        info_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        info_modal.addEventListener('shown.bs.modal', () => modals_up.set(id, true));
        modal_ids.push(id);

        let m = new bootstrap.Modal(info_modal);
    }

    function show_modal(id) {
        let m = bootstrap.Modal.getInstance("#"+id);
        m.show();
        modals_up.set(id, true);
    }

    function hide_modal(id) {
        let m = bootstrap.Modal.getInstance("#"+id);
        m.hide();
        // this next line might be redundant
        modals_up.delete(id);
    }

    function show_error_modal(message) {
        let span = document.getElementById("error-modal-message");
        span.innerHTML = "&nbsp;" + message;
        span.appendChild(document.createElement("br"));
        span.innerHTML += "Consider making a bug report on ";

        let anchor1 = document.createElement("a");
        anchor1.setAttribute("href", "https://discord.gg/tUKtEXPE");
        anchor1.setAttribute("target", "_blank");
        anchor1.innerHTML = "discord";

        span.appendChild(anchor1);

        span.appendChild(document.createElement("br"));
        span.innerHTML += "It may be useful to attach the current ";

        let anchor2 = document.createElement("a");
        anchor2.setAttribute("href", window.location.href + "/debug");
        anchor2.setAttribute("target", "_blank");
        anchor2.innerHTML = "server state";

        span.appendChild(anchor2);

        let m = bootstrap.Modal.getInstance("#error-modal");
        m.show();
        modals_up.set("error-modal", true);
    }

    function show_info_modal(message) {
        let span = document.getElementById("info-modal-message");
        span.innerHTML = "&nbsp;" + message;
        let m = bootstrap.Modal.getInstance("#info-modal");
        m.show();
        modals_up.set("info-modal", true);
    }

    function show_prompt_modal(message, handler) {
        let span = document.getElementById("prompt-modal-message");
        span.innerHTML = "&nbsp;" + message;
        let m = bootstrap.Modal.getInstance("#prompt-modal");
        let button = document.getElementById("prompt-modal-ok");
        button.onclick = handler;
        m.show();
        modals_up.set("prompt-modal", true);
    }

    function add_scissors_modal() {
        let id = "scissors-modal";
        let paragraph = document.createElement("p");
        paragraph.setAttribute("id", id + "-paragraph");
        paragraph.innerHTML = "Are you sure you want to remove this entire move branch?";

        let title = document.createElement("h5");
        title.innerHTML = "Delete Branch";

        let scissors_modal = add_modal(id, title, paragraph, true, () => state.network_handler.prepare_cut());
        scissors_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        modal_ids.push(id);

        let m = new bootstrap.Modal(scissors_modal);
    }

    function add_trash_modal() {
        let id = "trash-modal";
        let paragraph = document.createElement("p");
        paragraph.setAttribute("id", id + "-paragraph");
        paragraph.innerHTML = "This will reset the entire board. Continue?";

        let title = document.createElement("h5");
        title.innerHTML = "Reset";

        let trash_modal = add_modal(id, title, paragraph, true, () => state.network_handler.prepare_trash());
        trash_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        modal_ids.push(id);
        let m = new bootstrap.Modal(trash_modal);
    }

    function add_download_modal() {
        let id = "download-modal";
        let body_element = document.createElement("div");

        let save_button = new_icon_button("bi-floppy", () => state.download());
        save_button.setAttribute("class", "btn btn-primary");
        save_button.setAttribute("id", "save-sgf");
        let line1 = document.createElement("p");
        line1.appendChild(save_button);
        let span1 = document.createElement("span");
        span1.style.cursor = "pointer";
        span1.innerHTML = "&nbsp;Save SGF to disk";
        span1.onclick = () => document.getElementById("save-sgf").click();
        line1.appendChild(span1);
        body_element.appendChild(line1);

        let copy_button = new_icon_button("bi-copy", () => state.copy(state.get_sgf_link()));
        copy_button.setAttribute("id", "copy-sgf");
        copy_button.setAttribute("type", "button");
        copy_button.setAttribute("class", "btn btn-primary");

        copy_button.setAttribute("data-bs-toggle", "tooltip");
        copy_button.setAttribute("data-bs-placement", "bottom");
        copy_button.setAttribute("data-bs-trigger", "click");
        copy_button.setAttribute("data-bs-title", "Copied!");

        // an example of using the "bootstrap" object
        // haven't needed to use it until this tooltip
        copy_button.addEventListener("shown.bs.tooltip", (e) => {
            setTimeout(() => {
                const tooltip = bootstrap.Tooltip.getInstance(e.target);
                if (tooltip) {tooltip.hide();}
            }, 1500);
        });

        let line2 = document.createElement("p");
        line2.appendChild(copy_button);
        let span2 = document.createElement("span");
        span2.style.cursor = "pointer";
        span2.innerHTML = "&nbsp;Copy link to raw SGF";
        span2.onclick = () => document.getElementById("copy-sgf").click();
        line2.appendChild(span2);
        body_element.appendChild(line2);

        let title = document.createElement("h5");
        title.innerHTML = "Download SGF";

        let download_modal = add_modal(id, title, body_element);
        download_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        modal_ids.push(id);
        let m = new bootstrap.Modal(download_modal);
    }

    function add_upload_modal() {
        let id = "upload-modal";
        let body_element = document.createElement("div");

        // upload file from disk

        // dummy element which acts as upload operation
        let inp = document.createElement("input");
        // we CAN accept multiple files
        // ... but should we?
        inp.setAttribute("multiple", "true");
        inp.id = "upload-sgf";
        inp.style.display = "none";
        inp.setAttribute("type", "file");
        inp.onclick = () => state.upload();
 
        let disk_button = new_icon_button("bi-hdd", () => document.getElementById("upload-sgf").click());
        disk_button.setAttribute("class", "btn btn-primary");
        let line1 = document.createElement("p");
        line1.appendChild(inp);
        line1.appendChild(disk_button);
        let span1 = document.createElement("span");
        span1.style.cursor = "pointer";
        span1.innerHTML = "&nbsp;Upload from disk";
        span1.onclick = () => document.getElementById("upload-sgf").click();
        line1.appendChild(span1);
        body_element.appendChild(line1);

        // upload via url/sgf paste

        let paste_button = new_icon_button("bi-box-arrow-in-right", () => state.paste());
        paste_button.setAttribute("type", "button");
        paste_button.setAttribute("class", "btn btn-primary");

        let line2 = document.createElement("p");
        line2.appendChild(paste_button);
        let span2 = document.createElement("span");
        span2.innerHTML = "&nbsp;";
        line2.appendChild(span2);

        let input_text = document.createElement("input");
        input_text.setAttribute("type", "text");
        input_text.setAttribute("id", "upload-textarea");
        input_text.setAttribute("placeholder", "Paste SGF or URL");

        input_text.addEventListener("keypress", (event) => {
            if (event.key == "Enter") {
                state.paste();
            }
        });

        line2.appendChild(input_text);
        body_element.appendChild(line2);

        // finish

        let title = document.createElement("h5");
        title.innerHTML = "Upload SGF";

        let upload_modal = add_modal(id, title, body_element);
        upload_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        modal_ids.push(id);
        let m = new bootstrap.Modal(upload_modal);

        // see upload() for where the upload_modal gets hidden after upload
        // a little indelicate, but it's another example of using the bootstrap object
        // to directly manipulate elements
    }

    function add_gameinfo_modal() {
        let id = "gameinfo-modal";
        let body_element = document.createElement("div");

        let link_button = new_icon_button("bi-link", () => state.copy(state.get_link()));
        link_button.setAttribute("type", "button");
        link_button.setAttribute("class", "btn btn-primary");

        link_button.setAttribute("data-bs-toggle", "tooltip");
        link_button.setAttribute("data-bs-placement", "bottom");
        link_button.setAttribute("data-bs-trigger", "click");
        link_button.setAttribute("data-bs-title", "Copied!");

        // an example of using the "bootstrap" object
        // haven't needed to use it until this tooltip
        link_button.addEventListener("shown.bs.tooltip", (e) => {
            setTimeout(() => {
                const tooltip = bootstrap.Tooltip.getInstance(e.target);
                if (tooltip) {tooltip.hide();}
            }, 1500);
        });

        let line1 = document.createElement("p");
        line1.appendChild(link_button);
        let span1 = document.createElement("span");
        span1.innerHTML = "&nbsp;Share board";
        line1.appendChild(span1);

        let href = window.location.href;
        let input_text = document.createElement("input");
        input_text.setAttribute("class", "form-control");
        input_text.setAttribute("type", "text");
        input_text.setAttribute("disabled", true);
        input_text.setAttribute("value", href);
        line1.appendChild(input_text);

        let line2 = document.createElement("p");
        line2.setAttribute("id", "info-modal-gameinfo");

        body_element.appendChild(line1);
        body_element.appendChild(line2);

        let title = document.createElement("h5");
        title.innerHTML = "Game Information";

        let gameinfo_modal = add_modal(id, title, body_element, false);
        update_gameinfo_modal();
        // on close of modal
        gameinfo_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        modal_ids.push(id);
        let m = new bootstrap.Modal(gameinfo_modal);
    }

    function add_users_modal() {
        let id = "users-modal";
        let paragraph = document.createElement("p");
        paragraph.setAttribute("id", id + "-paragraph");
        /*
        let button = new_icon_button("bi-info-square");
        button.setAttribute("class", "btn btn-primary");
        paragraph.appendChild(button);
        */

        let span1 = document.createElement("span");
        span1.setAttribute("id", id + "-message");
        paragraph.appendChild(span1);

        let title = document.createElement("h5");
        title.innerHTML = "Users";

        let users_modal = add_modal(id, title, paragraph, false);
        users_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
        users_modal.addEventListener('shown.bs.modal', () => modals_up.set(id, true));
        modal_ids.push(id);

        let m = new bootstrap.Modal(users_modal);
    }

    function update_users_modal() {
        let body = document.getElementById("users-modal-message");
        body.innerHTML = "";
        body.appendChild(document.createElement("br"));
        for (let id in state.connected_users) {
            let nick = state.connected_users[id];
            let temp = document.createElement("div");
            temp.textContent = nick;
            body.innerHTML += temp.innerHTML + "<br>";
        }
    }

    function update_gameinfo_modal() {
        let body = document.getElementById("info-modal-gameinfo");
        body.innerHTML = "";
        let gameinfo = state.get_gameinfo();
        let result = "";
        for (let key in gameinfo) {
            let temp = document.createElement("div");
            temp.textContent = key + ": " + gameinfo[key];
            result += temp.innerHTML + "<br>";
        }
        body.innerHTML = result;
    }

    function add_modal(id, title_element, body_element, ok, handler) {
        let modals = document.getElementById("modals");
        let modal = document.createElement("div");
        modal.setAttribute("class", "modal fade");
        modal.setAttribute("id", id);
        modal.setAttribute("tabindex", "-1");
        modal.setAttribute("aria-labelledby", id + "-label");
        modal.setAttribute("aria-hidden", "true");

        let dialog = document.createElement("div");
        dialog.setAttribute("class", "modal-dialog");

        let content = document.createElement("div");
        content.setAttribute("class", "modal-content");

        let header = document.createElement("div");
        header.setAttribute("class", "modal-header");
        header.setAttribute("id", id + "-header");
        //header.style.backgroundColor = "#000000";

        title_element.setAttribute("class", "modal-title");
        title_element.setAttribute("id", id + "-label");

        let button_x = document.createElement("button");
        button_x.setAttribute("class", "btn-close");
        button_x.setAttribute("data-bs-dismiss", "modal");
        button_x.setAttribute("aria-label", "Close");

        header.appendChild(title_element);
        header.appendChild(button_x);

        let body = document.createElement("div");
        body.setAttribute("class", "modal-body");
        body.setAttribute("id", id + "-body");
        //body.style.backgroundColor = "#000000";

        body.appendChild(body_element);
        
        let footer = document.createElement("div");
        footer.setAttribute("class", "modal-footer");
        footer.setAttribute("id", id + "-footer");
        //footer.style.backgroundColor = "#000000";

        let button_close = document.createElement("button");
        button_close.setAttribute("class", "btn btn-secondary");
        button_close.setAttribute("data-bs-dismiss", "modal");
        button_close.innerHTML = "Close";
        footer.appendChild(button_close);

        if (ok) {
            let button_ok = document.createElement("button");
            button_ok.setAttribute("class", "btn btn-primary");
            button_ok.setAttribute("data-bs-dismiss", "modal");
            button_ok.onclick = handler;
            button_ok.innerHTML = "Ok";
            button_ok.id = id + "-ok";
            footer.appendChild(button_ok);
        }

        content.appendChild(header);
        content.appendChild(body);
        content.appendChild(footer);

        dialog.appendChild(content);

        modal.appendChild(dialog);
        modals.appendChild(modal);
        return modal;
    }

    return {
        get_prompt_bar,
        update_modals,
        modals_up,
        show_modal,
        show_error_modal,
        show_info_modal,
        show_prompt_modal,
        update_settings_modal,
        update_gameinfo_modal,
        update_users_modal,
        hide_modal,
        show_toast,
        get_nickname,
    };

}
