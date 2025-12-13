/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { Coord } from './common.js';

export {
    Board,
}

class Group {
    constructor(coords, rep, color, libs) {
        this.coords = coords;
        this.rep = rep;
        this.color = color;
        this.libs = libs;
    }
}

class Board {
    constructor(size, init_root=true) {

        this.size = size;
        this.points = [];
        var i,j;
        for (i=0; i<size; i++) {
            let row = [];
            for (j=0; j<size; j++) {
                row.push(0);
            }
            this.points.push(row);
        }
    }

    clear() {
        for (let i=0; i<this.size; i++) {
            for (let j=0; j<this.size; j++) {
                this.points[i][j] = 0;
            }
        }
    }

    set(start, color) {
        this.points[start.x][start.y] = color;
    }

    get(start) {
        return this.points[start.x][start.y];
    }

    remove(i, j) {
        // check to see if there's a stone there
        if (this.points[i][j] == 0) {
            return false;
        }

        // we need to remember what color it was
        let color = this.points[i][j];
        
        let point = new Coord(i, j);

        // remove the stone on the board
        this.set(point, 0);

        return true;
    }
}
