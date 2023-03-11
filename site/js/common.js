function getElemenByNameOfChiled(elem, name) {
    console.log(elem)
    var buf = []
    for (var c in elem.children) {
        try {
            if (elem.children[c].getAttribute("name") == name) {
                buf.push(elem.children[c])
            }

        } catch (e) { continue }
    }
    return buf
}

function sortSelect(selElem) {
    var tmpAry = new Array();
    for (var i=0;i<selElem.options.length;i++) {
        tmpAry[i] = new Array();
        tmpAry[i][0] = selElem.options[i].text;
        tmpAry[i][1] = selElem.options[i].value;
    }
    tmpAry.sort();
    while (selElem.options.length > 0) {
        selElem.options[0] = null;
    }
    for (var i=0;i<tmpAry.length;i++) {
        var op = new Option(tmpAry[i][0], tmpAry[i][1]);
        selElem.options[i] = op;
    }
    return;
}