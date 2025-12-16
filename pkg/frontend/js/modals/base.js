/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

export {
    add_modal
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
