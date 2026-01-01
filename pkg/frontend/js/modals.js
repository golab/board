/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_text_button, new_icon_button, add_tooltip, prefer_dark_mode } from './common.js';

import { add_modal } from './modals/base.js';
import { add_trash_modal } from './modals/trash.js';
import { add_scissors_modal } from './modals/scissors.js';
import { add_gameinfo_modal } from './modals/gameinfo.js';
import { add_download_modal } from './modals/download.js';
import { add_upload_modal } from './modals/upload.js';
import { add_error_modal } from './modals/error.js';
import { add_prompt_modal } from './modals/prompt.js';
import { add_info_modal } from './modals/info.js';
import { add_settings_modal, set_password, remove_password } from './modals/settings.js';
import { add_users_modal } from './modals/users.js';

function get_nickname() {
    let id = "settings-modal";
    return document.getElementById(id + "-nickname-bar").value;
}

function add_broken_connection_icon() {
    let icon = document.createElement("div");
    let obj = document.createElement("i");
    obj.setAttribute("class", "bi-wifi-off");
    icon.appendChild(obj);
    icon.id = "broken-connection-icon";
    icon.style.position = "fixed";
    icon.style.bottom = "12px";
    icon.style.right = "12px";
    icon.style.color = "red";
    icon.style.fontSize = "30px";
    icon.hidden = true;
    document.body.appendChild(icon);
}

function add_toast() {
    let toast_div = document.createElement("div");
    toast_div.id = "toasts";
    toast_div.setAttribute(
        "class",
        "toast-container position-fixed top-0 end-0 p-3"
    );
    document.body.appendChild(toast_div);
}

export function create_modals(_state) {
    // apparently necessary to make the tooltips work properly
    const state = _state;
    const modal_ids = [];
    const modals_up = new Map();

    function register_modal(modal) {
        modal.addEventListener('hidden.bs.modal', () => modals_up.delete(modal.id));
        modal_ids.push(modal.id);
    }

    // don't forget to add new modals here
    // otherwise clicking the button throws an error
    register_modal(add_trash_modal(state));
    register_modal(add_scissors_modal(state));
    register_modal(add_gameinfo_modal(state));
    update_gameinfo_modal();

    register_modal(add_download_modal(state));
    register_modal(add_upload_modal(state));
    register_modal(add_error_modal(state));
    register_modal(add_prompt_modal(state));
    register_modal(add_info_modal());
    register_modal(add_settings_modal(state));
    register_modal(add_users_modal());

    enable_tooltips();
    add_toast();
    add_broken_connection_icon();

    function show_broken_connection_icon() {
        let icon = document.getElementById("broken-connection-icon");
        icon.hidden = false;
    }

    function hide_broken_connection_icon() {
        let icon = document.getElementById("broken-connection-icon");
        icon.hidden = true;
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

    function clear_prompt_bar() {
        let id = "prompt-modal-input";
        let bar = document.getElementById(id);
        bar.value = "";
    }


    function enable_tooltips() {
        var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
        var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {return new bootstrap.Tooltip(tooltipTriggerEl)});
    }
 
    function update_modals() {
        update_gameinfo_modal();
        update_settings_modal();
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

    function show_modal(id) {
        let m = bootstrap.Modal.getInstance("#"+id);
        m.show();
        modals_up.set(id, true);
    }

    function hide_modal(id) {
        let m = bootstrap.Modal.getInstance("#"+id);
        m.hide();
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

    return {
        get_prompt_bar,
        clear_prompt_bar,
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
        show_broken_connection_icon,
        hide_broken_connection_icon,
    };
}
