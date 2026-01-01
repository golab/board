/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_icon_button } from '../common.js';
import { add_modal } from './base.js';
export {
    add_download_modal
}

function add_download_modal(state) {
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
    let m = new bootstrap.Modal(download_modal);
    return download_modal;
}
