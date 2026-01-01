/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_icon_button } from '../common.js';
import { add_modal } from './base.js';
export {
    add_upload_modal
}

function add_upload_modal(state) {
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
    let m = new bootstrap.Modal(upload_modal);
    return upload_modal;

    // see upload() for where the upload_modal gets hidden after upload
    // a little indelicate, but it's another example of using the bootstrap object
    // to directly manipulate elements
}
