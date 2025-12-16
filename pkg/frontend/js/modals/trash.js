/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { add_modal } from './base.js';
export {
    add_trash_modal
}

function add_trash_modal(state) {
    let id = "trash-modal";
    let paragraph = document.createElement("p");
    paragraph.setAttribute("id", id + "-paragraph");
    paragraph.innerHTML = "This will reset the entire board. Continue?";

    let title = document.createElement("h5");
    title.innerHTML = "Reset";

    let trash_modal = add_modal(id, title, paragraph, true, () => state.network_handler.prepare_trash());
    //trash_modal.addEventListener('hidden.bs.modal', () => modals_up.delete(id));
    let m = new bootstrap.Modal(trash_modal);
    return trash_modal;
}
