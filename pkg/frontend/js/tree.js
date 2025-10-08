/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { letterstocoord } from './common.js';
export {
    Tree
}

class Node {
    constructor(captured, index, up, fields, depth) {
        this.captured = captured;
        this.down = [];
        this.up = up;
        this.fields = fields;
        this.index = index;
        this.preferred_child = 0;
        this.depth = depth;
    }

    add_field(key, value) {
        if (this.fields == null) {
            this.fields = new Map();
        }
        if (!this.fields.has(key)) {
            this.fields.set(key, []);
        }
        this.fields.get(key).push(value);
    }

    remove_field(key, value) {
        if (this.fields == null) {
            this.fields = new Map();
            return;
        }
        if (!this.fields.has(key)) {
            return;
        }

        let index = -1;
        let values = this.fields.get(key);
        let i = 0;
        for (let v of values) {
            if (v == value) {
                index = i;
            }
            i++;
        }
        if (index == -1) {
            return;
        }

        values.splice(index, 1);

        if (values.length == 0) {
            this.fields.delete(key);
        }
    }

    // want to know color of node
    // if it's an erasing node, or if it's a pass node
    has_move() {
        return this.fields.has("B") || this.fields.has("W");
    }

    is_pass() {
        if (!this.has_move()) {
            return false;
        }
        if (this.fields.has("B") && this.fields.get("B")[0] == "") {
            return true;
        }
        if (this.fields.has("W") && this.fields.get("W")[0] == "") {
            return true;
        }
        return false;
    }

    color() {
        if (!this.has_move()) {
            return 0;
        }
        if (this.fields.has("B")) {
            return 1;
        }
        if (this.fields.has("W")) {
            return 2;
        }

        return 0;
    }

    coord() {
        if (!this.has_move() || this.is_pass()) {
            return null;
        }
        if (this.fields.has("B")) {
            return letterstocoord(this.fields.get("B")[0]);
        }
        if (this.fields.has("W")) {
            return letterstocoord(this.fields.get("W")[0]);
        }
    }

    a_stones() {
        let add = {};
        add[1] = [];
        add[2] = [];
        if (this.fields.has("AB")) {
            for (let c of this.fields.get("AB")) {
                add[1].push(c);
            }
        }
        if (this.fields.has("AW")) {
            for (let c of this.fields.get("AW")) {
                add[2].push(c);
            }
        }
        return add;
    }

    colors() {
        let a = this.a_stones();
        let result = new Map();
        let x = this.color();
        if (x != 0) {
            result.set(x, 1);
        }
        if (a[1].length > 0) {
            result.set(1, 1);
        }
        if (a[2].length > 0) {
            result.set(2,1);
        }
        return result;
    }

    has_child(coord, color) {
        for (let child of this.down) {
            // TODO: not sure if this functions as expected when we erase a stone
            // but it's fine for now
            if (child.coord() != null && child.coord().is_equal(coord) && child.color() == color) {
                return true;
            }
        }
        return false;
    }
}

function make_fields() {
    let fields = new Map();
    fields.set("GM", ["1"]);
    fields.set("FF", ["4"]);
    fields.set("CA", ["UTF-8"]);
    fields.set("PB", ["Black"]);
    fields.set("PW", ["White"]);
    fields.set("RU", ["Japanese"]);
    fields.set("KM", ["6.5"]);
    return fields;
}

class Tree {
    constructor(init_root) {
        this.next_index = 0;
        this.nodes = new Map();
        this.root = null;
        this.max_depth = 0;
        this.current_depth = 0;
        if (init_root) {
            this.root = new Node([], this.next_index, null, make_fields(), this.current_depth);
            this.max_depth = 1;
            this.current_depth = 1;
            this.next_index++;
            this.nodes.set(0, this.root);
        }
        this.current = this.root;
    }

    push(removed, set_index=-1, fields=null, force=false) {
        // to allow for reindexing on the fly
        let new_index = set_index;
        if (set_index == -1) {
            new_index = this.next_index;
        }

        let n = new Node(removed, new_index, this.current, fields, this.current_depth);
        if (this.root != null && set_index==-1 && !force) {
            // first check if it's already there
            for (let i=0; i<this.current.down.length; i++) {
                let node = this.current.down[i];
                let coord_old = node.coord();
                let coord_new = n.coord();
                if (coord_old != null && coord_new != null && coord_old.x == coord_new.x && coord_old.y == coord_new.y && node.color() == n.color()) {
                    this.current.preferred_child = i;
                    this.right();
                    return;
                }
            }
        }

        this.nodes.set(new_index, n);
        this.next_index = new_index+1;
        if (this.root == null) {
            this.root = n;
            this.current = n;
        } else {
            this.current.down.push(n);
            this.current.preferred_child = this.current.down.length-1
            let index = this.current.down.length - 1;
            this.current = this.current.down[index];
        }

        this.current_depth++;
        if (this.current_depth > this.max_depth) {
            this.max_depth = this.current_depth;
        }
    }

    push_pass(color) {
        let fields = new Map();
        let key = "B";
        if (color == 2) {
            key = "W";
        }
        fields.set(key, [""]);

        let removed = {};
        removed[1] = [];
        removed[2] = [];

        this.push(removed, -1, fields);
    }

    recompute_max_depth() {
        let start = this.root;
        let max_depth = 0;
        let stack = [start];
        while (stack.length > 0) {
            let cur = stack.pop();
            if (cur.depth > max_depth) {
                max_depth = cur.depth;
            }

            for (let child of cur.down) {
                stack.push(child);
            }
        }
        this.max_depth = max_depth+1;
    }

    cut(index) {
        // assume we are cutting a child of index
        let j = -1;
        for (let i=0; i<this.current.down.length; i++) {
            let node = this.current.down[i];
            if (node.index == index) {
                j = i;
                break;
            }
        }
        if (j == -1) {
            return;
        }

        // splice(x, y) removes y elements starting at index x
        this.current.down.splice(j, 1);
        this.nodes.delete(index);

        // adjust prefs
        if (this.current.preferred_child >= this.current.down.length) {
            this.current.preferred_child = 0;
        }

        // adust max_depth
        // TODO: figure this out
        this.recompute_max_depth();
    }

    left() {
        if (this.current.up == null) {
            return null;
        }
        let result = this.current;
        this.current = this.current.up;
        this.current_depth--;
        return result;
    }

    right() {
        if (this.current.down.length == 0) {
            return null;
        }
        let index = this.current.preferred_child;
        this.current = this.current.down[index];
        this.current_depth++;
        return this.current;
        //return [this.current.value, this.current.captured, this.current.color];
    }

    up() {
        if (this.current.down.length == 0) {
            return;
        }
        let c = this.current.preferred_child;
        let mod = this.current.down.length;
        this.current.preferred_child = (((c-1)%mod) + mod)%mod;
    }

    down() {
        if (this.current.down.length == 0) {
            return;
        }
        let c = this.current.preferred_child;
        let mod = this.current.down.length;
        this.current.preferred_child = (((c+1)%mod) + mod)%mod;
    }

    rewind() {
        this.current = this.root;
        this.current_depth = 1;
    }

    preferred() {
        // get a map of preferred indexes
        let pref = new Map();
        let cur = this.root;
        while (true) {
            pref.set(cur.index, true);
            if (cur.down.length == 0) {
                break;
            }
            cur = cur.down[cur.preferred_child];
        }
        return pref;
    }

    set_preferred(index) {
        let n = this.nodes.get(index);
        // set preferred path
        let cur = n;
        while(true) {
            let old_index = cur.index;
            cur = cur.up;
            if (cur == null) {
                break;
            }
            for (let i=0; i < cur.down.length; i++) {
                if (cur.down[i].index == old_index) {
                    cur.preferred_child = i;
                }
            }
        }
    }

    reset_prefs() {
        let stack = [this.root];
        while (stack.length > 0) {
            let cur = stack.pop();
            cur.preferred_child = 0;

            if (cur.down.length == 1) {
                stack.push(cur.down[0]);
            } else if (cur.down.length > 1) {
                // should go backward through array
                for (let i=cur.down.length-1; i >=0; i--) {
                    let n = cur.down[i];
                    stack.push(n);
                }
            }
        }
    }

    set_prefs(prefs) {
        let stack = [this.root];
        while (stack.length > 0) {
            let cur = stack.pop();
            let key = cur.index.toString();
            let p = prefs[key];
            cur.preferred_child = p;

            if (cur.down.length == 1) {
                stack.push(cur.down[0]);
            } else if (cur.down.length > 1) {
                // should go backward through array
                for (let i=cur.down.length-1; i >=0; i--) {
                    let n = cur.down[i];
                    stack.push(n);
                }
            }
        }
    }

    to_sgf() {
        let result = "";
        result += "(";
        if (this.root.fields == null) {
            /*
            let fields = new Map();
            fields.set("GM", ["1"]);
            fields.set("FF", ["4"]);
            fields.set("CA", ["UTF-8"]);
            // TODO: doesn't have to be 19
            fields.set("SZ", ["19"]);
            fields.set("PB", ["Clack"]);
            fields.set("PW", ["White"]);
            */
            this.root.fields = make_fields();
        }
        let stack = [this.root];

        while (stack.length > 0) {
            let cur = stack.pop();
            if (typeof cur == "string") {
                result += cur;
                continue;
            }
            result += ";";

            for (let [key, values] of cur.fields) {
                if (key == "IX") {
                    continue;
                }
                result += key;
                for (let v of values) {
                    result += "[";
                    result += v.replaceAll("]", "\\]");
                    result += "]";
                }
            }
            result += "IX" + "[" + cur.index + "]";

            if (cur.down.length == 1) {
                stack.push(cur.down[0]);
            } else if (cur.down.length > 1) {
                // should go backward through array
                for (let i=cur.down.length-1; i >=0; i--) {
                    let n = cur.down[i];
                    stack.push(")");
                    stack.push(n);
                    stack.push("(");
                }
            }
        }

        result += ")";
        return result;
    }
}
