/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { new_icon_button, add_tooltip, prefer_dark_mode } from './common.js';

export function create_buttons(_state) {
    const state = _state;

    // from observation, it seems the max number of buttons per row is 9
    // with 10 (on my device anyway) the buttons won't shrink by width anymore
    // and the overflow doesn't look good
    let review = document.getElementById("review");
    let button_row1 = document.getElementById("buttons-row1");
    let button_row2 = document.getElementById("buttons-row2");
    let button_row3 = document.getElementById("buttons-row3");

    let style = "";
    button_row1.style.margin = "auto";
    button_row2.style.margin = "auto";
    button_row3.style.margin = "auto";
    button_row1.style.display = "flex";
    button_row2.style.display = "flex";
    button_row3.style.display = "flex";

    let arrows = document.getElementById("arrows");

    // row 1

    let ghost_button = new_icon_button("bi-app");
    button_row1.appendChild(ghost_button);
    let h = button_row1.offsetHeight;
    ghost_button.style.display = "none";

    // toggle
    let toggle_button = new_icon_button("bi-circle-half", () => state.set_toggle());
    add_tooltip(toggle_button, "Place alternating stones (1)");
    button_row1.appendChild(toggle_button);

    // black stones
    let black_stone_button = new_icon_button("bi-circle-fill", () => state.set_black());
    add_tooltip(black_stone_button, "Place black stones (2)");
    button_row1.appendChild(black_stone_button);

    // white stones
    let white_stone_button = new_icon_button("bi-circle", () => state.set_white());
    add_tooltip(white_stone_button, "Place white stones (3)");
    button_row1.appendChild(white_stone_button);

    // pass
    let pass_button = new_icon_button("bi-arrow-repeat", () => state.network_handler.prepare_pass());
    add_tooltip(pass_button, "Pass (4)");
    button_row1.appendChild(pass_button);

    /*
    // remove stone
    let remove_stone_button = new_icon_button("bi-x-circle-fill", () => state.set_eraser());
    add_tooltip(remove_stone_button, "Remove stones");
    button_row1.appendChild(remove_stone_button);
    */

    // triangle
    let triangle_button = new_icon_button("bi-triangle", () => state.set_triangle());
    add_tooltip(triangle_button, "Place triangles (5)");
    button_row1.appendChild(triangle_button);

    // square
    let square_button = new_icon_button("bi-square", () => state.set_square());
    add_tooltip(square_button, "Place squares (6)");
    button_row1.appendChild(square_button);


    // letters
    let letter_button = new_icon_button("bi-alphabet-uppercase", () => state.set_letter());
    add_tooltip(letter_button, "Place letters (7)");
    button_row1.appendChild(letter_button);

    // numbers
    let number_button = new_icon_button("bi-123", () => state.set_number());
    add_tooltip(number_button, "Place numbers (8)");
    button_row1.appendChild(number_button);

    // row 2

    // pen
    let pen_button = new_icon_button("bi-pen", () => state.set_pen());
    add_tooltip(pen_button, "Draw with a pen (9)");
    button_row2.appendChild(pen_button);

    // color picker
    let color_picker = document.createElement("input");
    color_picker.setAttribute("type", "color");
    color_picker.setAttribute("id", "color-picker");

    let cls = prefer_dark_mode() ? "btn-dark" : "btn-light";
    color_picker.setAttribute("class", "btn " + cls + " wide-button");
    add_tooltip(color_picker, "Select pen color");
    color_picker.setAttribute("value", state.pen_color);

    // necessary because otherwise color picker doesn't end up being as
    // tall as the other buttons
    color_picker.style.height = h + "px";

    color_picker.onchange = function() {state.pen_color = this.value};
    button_row2.appendChild(color_picker);

    // eraser
    let eraser_button = new_icon_button("bi-eraser-fill", () => state.network_handler.prepare_erase_pen());
    add_tooltip(eraser_button, "Erase pen marks (0)");
    button_row2.appendChild(eraser_button);

    // scissors button
    let scissors_button = new_icon_button("bi-scissors", () => state.modals.show_modal("scissors-modal"));
    add_tooltip(scissors_button, "Cut branch (Ctrl+X / Cmd+X)");
    button_row2.appendChild(scissors_button);

    // copy button
    let copy_button = new_icon_button("bi-copy", () => state.network_handler.prepare_copy());
    add_tooltip(copy_button, "Copy branch (Ctrl+C / Cmd+C)");
    button_row2.appendChild(copy_button);

    // clipboard button
    let clipboard_button = new_icon_button("bi-clipboard", () => state.network_handler.prepare_clipboard());
    add_tooltip(clipboard_button, "Paste branch (Ctrl+V / Cmd+V)");
    button_row2.appendChild(clipboard_button);

    // score button
    let score_button = new_icon_button("bi-calculator", () => state.trigger_score());
    add_tooltip(score_button, "Score");
    button_row2.appendChild(score_button);

    // trash everything
    let trash_button = new_icon_button("bi-trash", () => state.modals.show_modal("trash-modal"));
    add_tooltip(trash_button, "Reset board");
    trash_button.setAttribute("class", "btn btn-danger wide-button");
    //trash_button.setAttribute("data-bs-toggle", "modal");
    //trash_button.setAttribute("data-bs-target", "#trash-modal");
    button_row2.appendChild(trash_button);

    // row 3

    // upload button
    let upload_button = new_icon_button("bi-upload", () => state.modals.show_modal("upload-modal"));
    add_tooltip(upload_button, "Upload SGF");
    //upload_button.setAttribute("data-bs-toggle", "modal");
    //upload_button.setAttribute("data-bs-target", "#upload-modal");
    button_row3.appendChild(upload_button);

    // download button
    let download_button = new_icon_button("bi-download", () => state.modals.show_modal("download-modal"));
    add_tooltip(download_button, "Download SGF");
    //download_button.setAttribute("data-bs-toggle", "modal");
    //download_button.setAttribute("data-bs-target", "#download-modal");
    button_row3.appendChild(download_button);

    // info button
    let info_button = new_icon_button("bi-info-circle", () => state.modals.show_modal("gameinfo-modal"));
    add_tooltip(info_button, "Game info");
    //info_button.setAttribute("data-bs-toggle", "modal");
    //info_button.setAttribute("data-bs-target", "#info-modal");
    button_row3.appendChild(info_button);

    // comments
    let comments_button = new_icon_button("bi-chat-left-text", () => state.comments_toggle());
    add_tooltip(comments_button, "Toggle comments");
    button_row3.appendChild(comments_button);

    // users
    let users_button = new_icon_button(
        "bi-people-fill",
        () => state.modals.show_modal("users-modal")
    );
    add_tooltip(users_button, "See connected users");
    button_row3.appendChild(users_button);

    // dark mode
    //let dark_mode_button = new_icon_button("bi-moon-fill", () => state.dark_mode_toggle());
    //add_tooltip(dark_mode_button, "Toggle dark mode");
    //button_row3.appendChild(dark_mode_button);

    // settings button
    let settings_button = new_icon_button("bi-gear", () => state.modals.show_modal("settings-modal"));
    add_tooltip(settings_button, "Settings");
    button_row3.appendChild(settings_button);

    // arrows

    // rewind
    let rewind_button = new_icon_button("bi-rewind-fill", () => state.network_handler.prepare_rewind());
    add_tooltip(rewind_button, "Rewind to beginning");
    arrows.appendChild(rewind_button);

    // go back one move
    let back_button = new_icon_button("bi-caret-left-fill", () => state.network_handler.prepare_left());
    add_tooltip(back_button, "Go back one move");
    arrows.appendChild(back_button);

    // move number
    let num = document.createElement("button");
    num.setAttribute("class", "btn " + cls + " disabled");
    num.setAttribute("id", "move-number");
    num.innerHTML = "0";
    arrows.appendChild(num);

    // go forward one move
    let forward_button = new_icon_button("bi-caret-right-fill", () => state.network_handler.prepare_right());
    add_tooltip(forward_button, "Go forward one move");
    arrows.appendChild(forward_button);

    let fastforward_button = new_icon_button("bi-fast-forward-fill", () => state.network_handler.prepare_fastforward());
    add_tooltip(fastforward_button, "Fast forward to end");
    arrows.appendChild(fastforward_button);

    // TODO: rethink this please
    let w = (state.width + state.pad*2);
    arrows.style.width = w + "px";
    arrows.style.margin = "auto";
    arrows.style.display = "flex";

    review.style.margin = "auto";
    review.style.display = "flex";
    review.style.height = w + "px";

    // name cards
    let namecards = document.getElementById("namecards");
    namecards.style.margin = "auto";
    namecards.style.display = "flex";

    let black = document.createElement("div");
    black.setAttribute("id", "black-namecard");
    black.classList.add("w-100", "text-white", "bg-dark", "justify-content-center");
    let black_name = document.createElement("span");
    black_name.setAttribute("id", "black-name");
    black.appendChild(black_name);
    let black_caps = document.createElement("span");
    black_caps.setAttribute("id", "black-caps");
    black.appendChild(black_caps);

    //add_tooltip(black, "Captures:");
    namecards.appendChild(black);

    let white = document.createElement("div");
    white.setAttribute("id", "white-namecard");
    white.classList.add("w-100", "text-black", "bg-light", "justify-content-center");
    let white_name = document.createElement("span");
    white_name.setAttribute("id", "white-name");
    white.appendChild(white_name);
    let white_caps = document.createElement("span");
    white_caps.setAttribute("id", "white-caps");
    white.appendChild(white_caps);

    let komi = document.createElement("komi");
    komi.setAttribute("id", "komi");
    white.appendChild(komi);

    namecards.appendChild(white);

    return {};
}
