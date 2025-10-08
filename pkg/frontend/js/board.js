/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { Tree } from './tree.js';
import { Parser } from './sgf.js';
import { opposite, Coord, ObjectSet, Result, letterstocoord } from './common.js';

export {
    Board,
    from_sgf
}

function from_sgf(b64_data) {
    let data = atob(b64_data);

    let p = new Parser(data);
    let result = p.parse();
    if (result.type == "error") { 
        console.log(result);
        return;
    }

    // result.value should be an SGFNode (the root)
    // now i need to turn that into a tree
    let root = result.value
    let size_field = root.fields.get("SZ");
    let size = 19;
    if (size_field != null) {
        if (size_field.length != 1) {
            console.log("error: SZ is a multifield");
            return;
        }
        size = parseInt(size_field[0]);
        if (size != 9 && size != 13 && size != 19) {
            console.log("error: non-standard size");
            return;
        }
    }
    let board = new Board(size, false);
    let stack = [root];
    
    while (stack.length > 0) {
        let cur = stack.pop();
        if (typeof cur == "string") {
            if (cur == "<") {
                // stole some of this from graphics.js (see left())
                let node = board.tree.left();
                if (node == null) {
                    continue;
                }

                let captured = node.captured;

                // redraw captured stones
                for (let c of [1,2]) {
                    for (let xy of captured[c]) {
                        board.set(xy, c);
                    }
                }

                let coord = node.coord();

                if (coord != null) {
                    // clear previous move
                    board.set(coord, 0);
                }

                // clear other stones
                for (let x of node.fields.get("AB") || []) {
                    let c = letterstocoord(x);
                    board.set(c, 0);
                }
                for (let x of node.fields.get("AW") || []) {
                    let c = letterstocoord(x);
                    board.set(c, 0);
                }
            }
            continue;
        }

        board.handle_sgfnode(cur);
        // push on children in reverse order
        for (let i=cur.down.length-1; i >=0; i--) {
            stack.push("<")
            stack.push(cur.down[i]);
        }
    }
    board.tree.current = board.tree.root;
    board.tree.reset_prefs();
    board.clear();
    
    return board;
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

        this.tree = new Tree(init_root);
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

    // get field value at current node
    get_field(key) {
        if (this.tree.current.fields != null && this.tree.current.fields.has(key)) {
            return this.tree.current.fields.get(key);
        }
        return [];
    }

    to_sgf() {
        return this.tree.to_sgf();
    }

    copy() {
        let b = new Board(this.size);
        var i,j;
        for (i=0; i<b.size; i++) {
            for (j=0; j<b.size; j++) {
                b.points[i][j] = this.points[i][j];
            }
        }
        return b;
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

        // push new node to the tree
        let removed = {};
        removed[1] = [];
        removed[2] = [];
        removed[color].push(point);

        let fields = new Map();
        fields.set("AE", [point.to_letters()]);

        this.tree.push(removed, -1, fields);
        return true;
    }

    handle_sgfnode(node) {
        let fields = node.fields;
        let index = -1;
        if (fields.has("IX")) {
            index = parseInt(fields.get("IX")[0]);
        }

        if (fields.has("B")) {
            let coord = letterstocoord(fields.get("B")[0]);
            // coord may be null, but that's ok
            this.place(coord, 1, index, fields, true);
            return;
        } else if (fields.has("W")) {
            let coord = letterstocoord(fields.get("W")[0]);
            // coord may be null, but that's ok
            this.place(coord, 2, index, fields, true);
            return;
        }

        // manual 
        let ae = [];
        let ab = []
        let aw = [];
        if (fields.has("AE")) {
            ae = fields.get("AE");
        }
        if (fields.has("AB")) {
            ab = fields.get("AB");
        }
        if (fields.has("AW")) {
            aw = fields.get("AW");
        }

        // it's super annoying if more than two of these are present
        
        // i'm going to be very opinionated about this and simply not do
        // any board calculations
        // 99% of the time, it shouldn't matter
        let removed = {};
        removed[1] = [];
        removed[2] = [];
        for (let c of ae) {
            let coord = letterstocoord(c);
            let color = this.get(coord);
            if (color != 0) {
                removed[color].push(coord);
            }
            this.set(coord, 0);
        }

        for (let c of ab) {
            let coord = letterstocoord(c);
            this.set(coord, 1);
        }

        for (let c of aw) {
            let coord = letterstocoord(c);
            this.set(coord, 2);
        }

        this.tree.push(removed, index, fields, true);
    }

    place(start, color, set_index=-1, fields=null, force=false) {
        if (start == null) {
            // start == null means fields has no "B" or "W"
            // that means no captured
            let removed = {};
            removed[1] = [];
            removed[2] = [];

            this.tree.push(removed, set_index, fields, force);
            return new Result(true, []);
        }

        let dead = {};
        dead[1] = [];
        dead[2] = []
        // check to see if move is illegal
        let results = this.legal(start, color);
        if (!results.ok) {
            return new Result(false, dead);
        }

        // put the stone on the board
        this.set(start, color);

        // remove dead group
        for (let coord of results.values) {
            this.set(coord, 0);
            dead[opposite(color)].push(coord);
        }

        if (fields == null) {
            fields = new Map();
            let key = "B";
            if (color == 2) {
                key = "W";
            }
            fields.set(key, [start.to_letters()]);
        }

        this.tree.push(dead, set_index, fields, force);
        return new Result(true, dead);
    }

    legal(start, color) {
        // if there's already a stone there, it's illegal
        if (this.get(start) != 0) {
            return new Result(false, []);
        }
        this.set(start, color);
        // if it has >0 libs, it's legal
        let gp = this.find_group(start);
        let enough_libs = false;
        if (gp.libs.size > 0) {
            enough_libs = true;
        }

        // remove any groups of opposite color with 0 libs
        // important: only check neighboring area
        let dead_set = new ObjectSet();
        let killed_something = false;
        let nbs = this.neighbors(start);
        for (let nb of nbs) {
            if (this.get(nb) == 0) {
                continue;
            }
            let gp = this.find_group(nb);
            if ((gp.libs.size == 0) && (gp.color == opposite(color))) {
                for (let coord of gp.coords) {
                    this.set(coord, 0);
                    dead_set.add(coord);
                    killed_something = true;
                }
            }
        }
        let ok = true;
        if (!(enough_libs || killed_something)) {
            this.set(start, 0);
            ok = false;
        }
        let dead = [];
        for (let d of dead_set) {
            dead.push(JSON.parse(d));
        }
        return new Result(ok, dead);
    }

    neighbors(start) {
        let nbs = [];
        var x,y;
        for (x=-1; x<=1; x++) {
            for (y=-1; y<=1; y++) {
                if ((x!=0 && y!=0) || (x==0 && y==0)) {
                    continue;
                }
                let new_x = start.x+x;
                let new_y = start.y+y;
                if (new_x < 0 || new_y < 0) {
                    continue;
                }
                if (new_x >= this.size || new_y >= this.size) {
                    continue;
                }
                nbs.push(new Coord(new_x, new_y));
            }
        }
        return nbs;
    }

    find_group(start) {
        let c = this.copy();
        let color = this.get(start);
        let stack = [start];
        c.set(start, 0);
        let group = [];
        let libs = new ObjectSet();
        if (color == 0) {
            return group;
        }
        let rep = start;
        var point;
        while (stack.length > 0) {
            point = stack.pop();
            group.push(point);
            if (point.x < rep.x) {
                rep = point;
            } else if ((point.x == rep.x) && (point.y < rep.y)) {
                rep = point;
            }
            let nbs = this.neighbors(point);
            var nb;
            for (nb of nbs) {
                if (c.get(nb) == color) {
                    stack.push(nb);
                }
                if (this.get(nb) == 0) {
                    libs.add([nb.x, nb.y]);
                }
                c.set(nb, 0);
            }
        }
        return new Group(group, rep, color, libs);
    }

    groups() {
        var i,j;
        let check = [];
        for (i=0; i<this.size; i++) {
            check.push([]);
            for(j=0; j<this.size; j++) {
                check[i].push(0);
            }
        }
        let groups = [];
        for (i=0; i<this.size; i++) {
            for(j=0; j<this.size; j++) {
                if (check[i][j] == 0 && this.points[i][j] != 0) {
                    let group = this.find_group(new Coord(i,j));
                    for (let c of group.coords) {
                        check[i][j] = 1;
                    }
                    groups.push(group);
                }
            }
        }
        return groups;
    }
}
