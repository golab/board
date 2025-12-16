/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { add_modal } from './base.js';
export {
    add_users_modal
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

    let m = new bootstrap.Modal(users_modal);
    return users_modal;
}
