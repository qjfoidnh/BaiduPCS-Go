function s(j, r) {
    var a = [];
    var p = [];
    var o = "";
    var v = j.length;
    for (var q = 0; q < 256; q++) {
        a[q] = j.substr((q % v), 1).charCodeAt(0);
        // console.log(q, q%v);
        p[q] = q;
        // console.log(q, a[q], p[q]);
    }
    for (var u = q = 0; q < 256; q++) {
        u = (u + p[q] + a[q]) % 256;
        var t = p[q];
        p[q] = p[u];
        p[u] = t;
    }
    // console.log(p)
    for (var i = u = q = 0; q < r.length; q++) {
        i = (i + 1) % 256;
        u = (u + p[i]) % 256;
        var t = p[i];
        p[i] = p[u];
        p[u] = t;
        k = p[((p[i] + p[u]) % 256)];
        o += String.fromCharCode(r.charCodeAt(q) ^ k);
    }
    return o;
};

function base64encode(t) {
    var r, e, a, o, n, i, s = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
    for (a = t.length,
    e = 0,
    r = ""; a > e; ) {
        if (o = 255 & t.charCodeAt(e++),
        e == a) {
            r += s.charAt(o >> 2),
            r += s.charAt((3 & o) << 4),
            r += "==";
            break
        }
        if (n = t.charCodeAt(e++),
        e == a) {
            r += s.charAt(o >> 2),
            r += s.charAt((3 & o) << 4 | (240 & n) >> 4),
            r += s.charAt((15 & n) << 2),
            r += "=";
            break
        }
        i = t.charCodeAt(e++),
        r += s.charAt(o >> 2),
        r += s.charAt((3 & o) << 4 | (240 & n) >> 4),
        r += s.charAt((15 & n) << 2 | (192 & i) >> 6),
        r += s.charAt(63 & i)
    }
    return r
}

var s1 = s("e8c7d729eea7b54551aa594f942decbe", "37dbe07ade9359c1aa70807e847f768c13360ad2");
console.log(s1);
console.log(base64encode(s1));