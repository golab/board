/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { get_viewport, new_icon_button } from './common.js';

export function create_comments(_state) {
    let state = _state;
    let container = document.getElementById("comments");
    let comments = document.createElement("div");
    comments.style.textAlign = "left";
    comments.style.background = "#FEFEFE";
    let input_bar = document.createElement("input");
    input_bar.placeholder = "Comment...";

    input_bar.addEventListener("keypress", (event) => {
        if (event.key == "Enter") {
            let v = input_bar.value;
            input_bar.value = "";
            state.network_handler.prepare_comment(v);
        }
    });

    container.appendChild(comments);
    container.appendChild(input_bar);
    container.hidden = true;
    let _hidden = true;

    resize();

    function update(text) {
        let temp = document.createElement("div");
        temp.textContent = text;
        comments.innerHTML += temp.innerHTML + "<br>";
        temp.remove();
    }

    function store(text) {
        state.board.tree.current.add_field("C", text);
    }

    function clear() {
        comments.innerHTML = "";
    }

    function hidden() {
        return _hidden;
    }

    function show() {
        container.hidden = false;
        _hidden = false;
    }

    function hide() {
        container.hidden = true;
        _hidden = true;
    }

    function resize() {

        let vp = get_viewport();
        let new_width = 0;
        if (vp == "xs" || vp == "sm" || vp == "md") {
            let content = document.getElementById("content");
            new_width = content.offsetWidth;
        } else {
            let review = document.getElementById("review")
            new_width = window.innerWidth - review.offsetWidth - 100;
        }

        container.style.width = new_width + "px";
        comments.style.width = new_width + "px";
        input_bar.style.width = new_width + "px";
    }

    return {
        update,
        store,
        clear,
        hidden,
        hide,
        show,
        resize,
    };
}
