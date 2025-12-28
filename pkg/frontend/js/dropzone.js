/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { get_dims } from './common.js';

function dropzone(fn) {
    let [width, height] = get_dims();
    let zone = document.createElement("div");
    zone.id = "drop-upload";
    zone.classList.add("text-center");
    zone.style.position = "fixed";
    zone.style.width = width + "px";
    zone.style.height = height + "px";
    zone.style.left = "0";
    zone.style.top = "0";
    zone.style.pointerEvents = "none";
    // need to make sure it's on top of the board graphics
    zone.style.zIndex = 2000;

    let inp = document.createElement("input");
    inp.setAttribute("type", "file");
    inp.multiple = true;
    inp.hidden = true;
    zone.appendChild(inp);

    let overlay = document.createElement("div");
    overlay.id = "drop-upload-overlay";
    overlay.hidden = true;
    overlay.style.position = "relative";
    overlay.style.top = "40%";
    overlay.style.color = "#48BBFF";

    let icon = document.createElement("i");
    icon.classList.add("bi-cloud-arrow-up-fill");
    icon.style.fontSize = "400%";
    overlay.appendChild(icon);
    let text = document.createElement("div");
    text.innerHTML = "Drag and drop to upload";
    text.style.fontSize = "200%";
    overlay.appendChild(text);
    zone.appendChild(overlay);

    function prevent_defaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    function highlight(e) {
        overlay.hidden = false;
        zone.style.background = "#000000cc";
    }

    function unhighlight(e) {
        overlay.hidden = true;
        zone.style.background = "";
    }

    function handle_drop(e) {
        const dt = e.dataTransfer;
        const files = dt.files;
        fn(files);
    }

    // prevent defaults
    for (let evt of ["dragenter", "dragover", "dragleave", "drop"]) {
        window.addEventListener(evt, prevent_defaults, false);
    }

    // highlight during dragenter or dragover
    for (let evt of ["dragenter", "dragover"]) {
        window.addEventListener(evt, highlight, false);
    }

    // unhighlight on dragleave or drop
    for (let evt of ["dragleave", "drop"]) {
        window.addEventListener(evt, unhighlight, false);
    }

    // handle drop
    window.addEventListener("drop", handle_drop, false);

    return zone;
}

export function create_dropzone(fn) {
    let body = document.body;
    body.appendChild(dropzone(fn));
    return {};
}

export function resize_dropzone() {
    let dropzone = document.getElementById("drop-upload");
    let [width, height] = get_dims();
    dropzone.style.width = width + "px";
    dropzone.style.height = height + "px";
}
