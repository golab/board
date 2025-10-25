/*
Copyright (c) 2025 Jared Nishikawa

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package sgfsamples

const (
	Empty             = "(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black])"
	SimpleTwoBranches = "(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black](;B[pd];W[dd];B[pp];W[dp])(;B[dd];W[ee]))"
	SimpleWithComment = "(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black]C[comment1](;B[pd];W[dd];B[pp];W[dp])(;B[dd];W[ee]))"
	SimpleFourMoves   = "(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black];B[pd];W[dd];B[pp];W[dp];B[];W[])"
	SimpleEightMoves  = "(;PW[White]RU[Japanese]KM[6.5]GM[1]FF[4]CA[UTF-8]SZ[19]PB[Black];B[pd];W[dd];B[pp];W[dp];B[cc];W[dc];B[cd];W[de])"
	Scoring1          = "(;SZ[19]KM[6.5]FF[4]WRE[B+11.5]OT[3x30 byo-yomi]PB[black]PW[white]GM[1]DT[2025-09-14]TM[2400]RU[Japanese]CA[UTF-8];B[pd];W[dd];B[pp];W[dp];B[fq];W[qc];B[qd];W[pc];B[od];W[nb];B[dn];W[en];B[em];W[fn];B[cq];W[dm];B[dq];W[cn];B[cf];W[fc];B[cd];W[cc];B[bc];W[ce];B[bd];W[be];B[de];W[bf];B[dc];W[ed];B[db];W[cg];B[dl];W[cm];B[pj];W[ep];B[gq];W[mc];B[md];W[ld];B[nc];W[lb];B[le];W[ke];B[nd];W[rc];B[kf];W[kd];B[jf];W[lf];B[me];W[lg];B[kh];W[pi];B[oi];W[qj];B[ph];W[qi];B[oj];W[ql];B[df];W[dh];B[ff];W[hd];B[dg];W[ch];B[hf];W[gg];B[he];W[gd];B[fe];W[lj];B[lh];W[qh];B[dj];W[ei];B[fk];W[gj];B[gk];W[hj];B[hk];W[ej];B[cl];W[ek];B[el];W[gm];B[gl];W[ij];B[ik];W[im];B[jl];W[jk];B[il];W[jn];B[bi];W[bh];B[bk];W[kk];B[hn];W[hm];B[kl];W[ll];B[km];W[hq];B[cp];W[co];B[ji];W[jj];B[lm];W[ml];B[hp];W[eq];B[er];W[fr];B[gr];W[ip];B[io];W[gp];B[ho];W[dr];B[fs];W[cr];B[br];W[bq];B[bp];W[ar];B[aq];W[bs];B[as];W[jo];B[iq];W[ar];B[fj];W[fi];B[as];W[lo];B[mm];W[ar];B[ob];W[oc];B[as];W[nm];B[bm];W[no];B[nn];W[mn];B[on];W[oo];B[po];W[pn];B[om];W[nl];B[ol];W[nq];B[jp];W[kp];B[oq];W[nr];B[np];W[mp];B[op];W[mo];B[or];W[an];B[gi];W[hh];B[am];W[bn];B[dk];W[ar];B[fm];W[gn];B[as];W[aj];B[bj];W[ar];B[je];W[bq];B[fh];W[gh];B[br];W[fg];B[as];W[pm];B[pl];W[qm];B[oh];W[lq];B[ic];W[jb];B[jc];W[kc];B[ib];W[ia];B[ha];W[ja];B[hb];W[rp];B[ro];W[qo];B[rq];W[rn];B[qp];W[so];B[qr];W[ar];B[na];W[mb];B[as];W[bb];B[ab];W[ar];B[kb];W[ka];B[as];W[cb];B[ca];W[ar];B[la];W[ma];B[as];W[ai];B[ak];W[ar];B[qk];W[rk];B[as];W[bl];B[al];W[ar];B[kn];W[ko];B[as];W[jd];B[id];W[ar];B[mk];W[lk];B[as];W[ad];B[ba];W[ar];B[ii];W[hi];B[as];W[gf];B[ge];W[ar];B[qg];W[bq];B[rg];W[pk];B[ok];W[ap];B[rd];W[ki];B[jh];W[sc];B[eg];W[eh];B[kr];W[jq];B[jr];W[ip];B[hr];W[mi];B[mj];W[li];B[mh];W[rh];B[lr];W[ns];B[mr];W[mq];B[sh];W[sf];B[se];W[sd];B[re];W[si];B[sg];W[ee];B[ef];W[os];B[ps];W[nj];B[nk];W[ni];B[nh];W[ig];B[qk];W[rj];B[ah];W[ag];B[if];W[jg];B[kg];W[ms];B[ls];W[ai];B[ds];W[fp];B[jp];W[ci];B[cj];W[ip];B[pb];W[qb];B[jp];W[cs];B[es];W[ip];B[qa];W[ra];B[jp];W[kq];B[ip];W[sq];B[sr];W[sp];B[go];W[fo];B[pk];W[ae];B[fl];W[jm];B[aj];W[ah];B[];W[])"
)
