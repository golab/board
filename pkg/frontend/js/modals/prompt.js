/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_icon_button } from '../common.js';
import { add_modal } from './base.js';
export {
    add_prompt_modal
}

function add_prompt_modal(state) {
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
    prompt_modal.addEventListener(
        'shown.bs.modal',
        () => {
            prompt_bar.focus();
        }
    );

    let m = new bootstrap.Modal(prompt_modal);
    return prompt_modal;
}
