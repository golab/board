/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_icon_button } from '../common.js';
import { add_modal } from './base.js';
export {
    add_gameinfo_modal
}

function add_gameinfo_modal(state) {
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
    // on close of modal
    let m = new bootstrap.Modal(gameinfo_modal);
    return gameinfo_modal;
}
