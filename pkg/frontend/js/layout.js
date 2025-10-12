/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { get_dims } from './common.js';

function create_review() {
    let review_container = document.createElement("div");
    let review = document.createElement("div");
    review.id = "review";
    review.setAttribute("size", "19");
    review_container.appendChild(review);
    return review_container;
}

function create_buttons1() {
    let buttons1 = document.createElement("div");
    buttons1.id = "buttons-row1";

    let buttons1_container = document.createElement("div");
    buttons1_container.appendChild(buttons1);
    buttons1_container.style.flex = "1";
    return buttons1_container;
}

function create_buttons2() {
    let buttons2 = document.createElement("div");
    buttons2.id = "buttons-row2";

    let buttons2_container = document.createElement("div");
    buttons2_container.appendChild(buttons2);
    buttons2_container.style.flex = "1";

    return buttons2_container;
}

function create_buttons3() {
    let buttons3 = document.createElement("div");
    buttons3.id = "buttons-row3";

    let buttons3_container = document.createElement("div");
    buttons3_container.appendChild(buttons3);

    return buttons3_container;
}

function create_explorer() {
    let explorer_container = document.createElement("div");
    let container = document.createElement("div");
    container.id = "explorer_container";
    container.style.overflowX = "auto";
    container.style.overflowY = "auto";
    container.style.resize = "vertical";
    let explorer = document.createElement("div");
    explorer.id = "explorer";
    explorer.style.position = "relative";
    container.appendChild(explorer);
    explorer_container.appendChild(container);
    return explorer_container;
}

function create_arrows() {
    let arrows = document.createElement("div");
    arrows.id = "arrows";

    let arrows_container = document.createElement("div");
    arrows_container.appendChild(arrows);

    return arrows_container;
}

function create_comments() {
    let comments = document.createElement("div");
    comments.id = "comments";

    let comments_container = document.createElement("div");
    comments_container.appendChild(comments);

    return comments_container;
}

function _layout(height) {
    let container_fluid = document.createElement("div");
    container_fluid.classList.add("container-fluid");
    container_fluid.classList.add("text-center");

    let cols = document.createElement("div");
    cols.classList.add("d-flex");
    cols.classList.add("gap-3");
    cols.classList.add("flex-column", "flex-md-row");

    let col1 = document.createElement("div");
    col1.classList.add("d-flex");
    col1.classList.add("flex-column");
    col1.style.width = "65%";

    let col2 = document.createElement("div");
    col2.classList.add("d-flex");
    col2.classList.add("flex-column");
    col2.classList.add("flex-fill");

    let button_row = document.createElement("div");
    button_row.classList.add("d-flex");
    button_row.appendChild(create_buttons1());
    button_row.appendChild(create_buttons2());

    col1.appendChild(button_row);
    col1.appendChild(create_review());
    col1.appendChild(create_arrows());
    col1.appendChild(create_buttons3());

    col2.appendChild(create_explorer());
    col2.appendChild(create_comments());

    cols.appendChild(col1);
    cols.appendChild(col2);

    container_fluid.appendChild(cols);

    return container_fluid;
}

function create_div(cl, id) {
    let d = document.createElement("div");
    if (cl != null && cl != "") {
        d.setAttribute("class", cl);
    }
    if (id != null && id != "") {
        d.id = id;
    }
    return d;
}

function layout() {
    let container_fluid = document.createElement("div");
    container_fluid.classList.add("container-fluid");
    container_fluid.classList.add("text-center");

    let row1 = create_div("row");
    let a = create_div("col-lg-4 col-sm-12 gx-0");
    let a1 = create_div("", "buttons-row1");
    a.appendChild(a1);
    let b = create_div("col-lg-4 col-sm-12 gx-0");
    let b1 = create_div("", "buttons-row2");
    b.appendChild(b1);
    let c = create_div("col-lg-4 gx-0 ps-lg-4");
    let namecards = create_div("w-100", "namecards");
    c.appendChild(namecards);

    row1.appendChild(a);
    row1.appendChild(b);
    row1.appendChild(c);

    let row2 = create_div("row");
    let d = create_div("col-lg-8 col-sm-12 gx-0");
    let r = create_div("", "review");
    r.setAttribute("size", "19");
    let arrows = create_div("", "arrows");
    let b3 = create_div("", "buttons-row3");
    d.appendChild(r);
    d.appendChild(arrows);
    d.appendChild(b3);

    let e = create_div("col-lg-4 gx-0 ps-lg-4");
    let exp_container = create_div("", "explorer_container");
    exp_container.style.overflowX = "auto";
    exp_container.style.overflowY = "auto";
    exp_container.style.resize = "vertical";
    let explorer = create_div("", "explorer");
    explorer.style.position = "relative";
    exp_container.appendChild(explorer);
    let comments = create_div("", "comments");
    e.appendChild(exp_container);
    e.appendChild(comments);

    row2.appendChild(d);
    row2.appendChild(e);

    //let row3 = create_div("row");
    //let f = create_div("col-lg-8 col-sm-12 gx-0");
    //let b3 = create_div("", "buttons-row3");
    //f.appendChild(b3);

    //row3.appendChild(f);

    container_fluid.appendChild(row1);
    container_fluid.appendChild(row2);
    //container_fluid.appendChild(row3);

    return container_fluid;
}

export function create_layout() {
    let [width, height] = get_dims();
    //console.log(width, height);
    let content = document.getElementById("content");
    content.appendChild(layout());

    return {};
}
