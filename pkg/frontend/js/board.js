/*
Copyright (c) 2026 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { Coord } from './common.js';

export function create_board(size) {
    let points = [];
    var i,j;
    for (i=0; i<size; i++) {
        let row = [];
        for (j=0; j<size; j++) {
            row.push(0);
        }
        points.push(row);
    }

    function clear() {
        for (let i=0; i<size; i++) {
            for (let j=0; j<size; j++) {
                points[i][j] = 0;
            }
        }
    }

    function set(start, color) {
        points[start.x][start.y] = color;
    }

    function get(start) {
        return points[start.x][start.y];
    }

    function remove(i, j) {
        // check to see if there's a stone there
        if (points[i][j] == 0) {
            return false;
        }

        let point = new Coord(i, j);

        // remove the stone on the board
        set(point, 0);

        return true;
    }
    return {
        clear,
        set,
        get,
        remove,
    };
}
