/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

import { letterstocoord } from './common.js';

export {
    Parser,
    merge,
}

class Expr {
    constructor(type, value) {
        this.type = type;
        this.value = value;
    }
}

class SGFNode {
    constructor(fields, index) {
        this.fields = fields;
        this.down = [];
        this.index = index;
    }

    to_sgf(root=false) {
        let result = "";
        if (root) {
            result += "(";
        }
        result += ";";
        for (let [field, values] of this.fields) {
            result += field;
            for (let value of values) {
                result += "[";
                result += value.replaceAll("]", "\\]");
                result += "]";
            }
        }

        for (let d of this.down) {
            if (this.down.length > 1) {
                result += "(" + d.to_sgf() + ")";
            } else {
                result += d.to_sgf();
            }
        }
        if (root) {
            result += ")";
        }
        return result;
    }
}

let x = 0;

function iswhitespace(c) {
    return (c == "\n" || c == " " || c == "\t" || c == "\r");
}

class Parser {
    constructor(text) {
        this.text = text;
        this.index = 0;
    }

    parse() {
        this.skip_whitespace();
        let c = this.read();
        if (c == "(") {
            return this.parse_branch();
        } else {
            return new Expr("error", "unexpected " + c);
        }
    }

    skip_whitespace() {
        while (true) {
            if (iswhitespace(this.peek(0))) {
                this.read();
            } else {
                break;
            }
        }
        return new Expr("whitespace", "");
    }

    parse_key() {
        let s = "";
        while (true) {
            let c = this.peek(0);
            if (c == "\0") {
                return new Expr("error", "bad key");
            } else if (c < "A" || c > "Z") {
                break;
            }
            s += this.read();
        }
        return new Expr("key", s);
    }

    parse_field() {
        let s = "";
        while (true) {
            let t = this.read();
            if (t == "\0") {
                return new Expr("error", "bad field");
            } else if (t == "]") {
                break;
            } else if (t == "\\" && this.peek(0) == "]") {
                t = this.read();
            }
            s += t;
        }
        return new Expr("field", s);
    }

    parse_nodes() {
        let n = this.parse_node();
        if (n.type == "error") {
            return n;
        }
        let root = n.value;
        let cur = root;
        while (true) {
            let c = this.peek(0);
            if (c == ";") {
                this.read();
                let m = this.parse_node();
                    if (m.type == "error") {
                        return m;
                    }
                    let next = m.value;
                    cur.down.push(next);
                    cur = next;
            } else {
                break;
            }
        }
        return new Expr("nodes", [root, cur]);
    }

    parse_node() {
        var result;
        let fields = new Map();
        let index = 0;
        while (true) {
            this.skip_whitespace();
            let c = this.peek(0);
            if (c == "(" || c == ";" || c == ")") {
                break;
            }
            if (c < "A" || c > "Z") {
                return new Expr("error", "bad node (expected key)" + c);
            }
            result = this.parse_key();
            if (result.type == "error") {
                return result;
            }
            let key = result.value;

            let multifield = [];
            this.skip_whitespace();
            if (this.read() != "[") {
                return new Expr("error", "bad node (expected field) " + this.read());
            }
            result = this.parse_field();
            if (result.type == "error") {
                return result;
            }
            multifield.push(result.value);

            while (true) {
                this.skip_whitespace();
                if (this.peek(0) == "[") {
                    this.read();
                    result = this.parse_field();
                    if (result.type == "error") {
                        return result;
                    }
                    multifield.push(result.value);
                } else {
                    break;
                }
            }

            this.skip_whitespace();
            fields.set(key, multifield);
        }
        // need to add actual indexing
        let n = new SGFNode(fields, index);
        return new Expr("node", n);
    }

    parse_branch() {
        let root = null;
        let current = null;
        while (true) {
            let c = this.read();
            if (c == "\0") {
                return new Expr("error", "unfinished branch, expected ')'");
            } else if (c == ";") {
                let result = this.parse_nodes();
                if (result.type == "error") {
                    return result;
                }
                let node = result.value[0];
                let cur = result.value[1];
                if (root == null) {
                    root = node;
                    current = cur;
                } else {
                    current.down.push(node);
                    current = cur;
                }
            } else if (c == "(") {
                let result = this.parse_branch();
                if (result.type == "error") {
                    return result;
                }
                let new_branch = result.value;
                if (root == null) {
                    root = new_branch;
                    current = new_branch;
                } else {
                    current.down.push(new_branch);
                }
            } else if (c == ")") {
                break;
            }
        }
        return new Expr("branch", root);
    }

    read() {
        if (this.index >= this.text.length) {
            return "\0";
        }
        let result = this.text[this.index];
        this.index++;
        return result;
    }

    unread() {
        if (this.index == 0) {
            return;
        }
        this.index--;
    }

    peek(n) {
        if (this.index+n >= this.text.length) {
            return "\0";
        }
        return this.text[this.index+n];
    }
}

function merge(sgfs) {
    if (sgfs.length == 0) {
        return "";
    } else if (sgfs.length == 1) {
        return sgfs[0];
    }
    let size = 0;
    let fields = new Map();
    fields.set("GM", ["1"]);
    fields.set("FF", ["4"]);
    fields.set("CA", ["UTF-8"]);
    fields.set("PB", ["Black"]);
    fields.set("PW", ["White"]);
    fields.set("RU", ["Japanese"]);
    fields.set("KM", ["6.5"]);
    let new_root = new SGFNode(fields, 0);
    for (let sgf of sgfs) {
        let p = new Parser(sgf);
        let root = p.parse().value;
        let sizes = root.fields.get("SZ") || [];

        // if SZ is not provided, assume 19
        let _size = 19;
        if (sizes.length > 0) {
            _size = sizes[0];
        }

        // if we haven't set the (assumed) same size yet, set it
        if (size == 0) {
            size = _size;
        }

        // if not all the sgfs are the same size, just return the first one
        if (_size != size) {
            return sgfs[0];
        }

        if (root.fields.has("B") ||
            root.fields.has("W") ||
            root.fields.has("AB") ||
            root.fields.has("AW")) {
            // strip fields and save the node
            for (let f of ["RU", "SZ", "KM", "TM", "OT"]) {
                root.fields.delete(f);
            }
            new_root.down.push(root);
        } else {
            // otherwise save all the children
            for (let d of root.down) {
                new_root.down.push(d);
            }
        }
    }

    new_root.fields.set("SZ", [size.toString()]);
    return new_root.to_sgf(true);

    /*
    // if we exit the loop, all the sgfs are the same size
    // so make a new root, and set all the sgfs as children
    let result = "(;GM[1]FF[4]CA[UTF-8]PB[Black]PW[White]RU[Japanese]KM[6.5]SZ["
        + size.toString() + "]";
    for (let sgf of sgfs) {
        //console.log(sgf);
        result += "(" + sgf + ")";
    }
    result += ")";
    return result;
    */

}

function test() {
    let data = `(;GM[1]FF[4]CA[UTF-8]AP[CGoban:3]ST[2]
RU[Japanese]SZ[19]KM[6.50]
PW[ alice ]PB[bob ]
(;B[pd]
(;W[qf]
;B[nc]
(;W[qc]
;B[qd]C[comment [some comment\\]])
(;W[qd]
;B[qc]
;W[rc]TR[qd]
;B[qe]
;
;W[rd]
;B[pe]))
(;W[qc]
;B[qd]
;W[pc]TR[qc][pd][qd]
;B[od]LB[pc:D][qc:B][pd:A][qd:C])
(;W[oc]
;B[pc]
;W[mc]))
(;B[qg]))`

    let p = new Parser(data);
    let result = p.parse();
    if (result.type == "error") {
        console.log(result);
        return;
    }
    console.log(result.value);
}

//test();



