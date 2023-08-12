var id = document.getElementById("drawflow");
const editor = new Drawflow(id);
var flag_running_save = false
var current_id = undefined;
var clipboard_node = undefined;
 
function init() {
    if (flag_running_save) {
        return
    }
    try {
        flag_running_save = true
        editor.reroute = true;
        // aca es donde esta el workflow
        editor.drawflow = {}
        fetch("nflow/app")
            .then(response => response.json())
            .then(data => {
                editor.drawflow = data
                editor.useuuid = true
                editor.start();
                tabs = document.querySelector(".top_panel > .panel_items")
                for (var name in data.drawflow) {
                    if (name == "Home" || name == "" || name == undefined) {
                        continue
                    }
                    tabs.innerHTML += `<button id="module_${name}" class="panel_item" onclick="editor.changeModule('${name}'); changeModule(event);" >${name}</button>`;
                }
            })
    } finally {
        flag_running_save = false
    }
}


function refresh() {
    if (flag_running_save) {
        return
    }
    try {
        flag_running_save = true

        fetch("nflow/app")
            .then(response => response.json())
            .then(data => {
                editor.clear()
                editor.drawflow = data
                editor.load()
            })

    } finally {
        flag_running_save = false
    }
}


// Events!
editor.on('nodeCreated', function (id) {
    console.log("Node created " + id);
    save()
})

editor.on('nodeRemoved', function (id) {
    current_id = undefined;
    console.log("Node removed " + id);
    id_box_in_prop = undefined
    setFieldNode()
    panel_prop.innerHTML = ""
    panel_prop.style.width = 0;
    panel_prop.style.padding = 0;
    id_box_in_prop = undefined;
    save()
})

function nodeSelectedCustom (id) {
    console.log(1)
    current_id = id;
    panel_prop.style.color = "white";
    panel_prop.style.padding = "20px";

    id_box_in_prop = id
    bottom_panel.style.animationName = "show_bottom_panle"
    bottom_panel.style.bottom = "0px"
    while(1){
        try{
            document.querySelector("#node-" + id_box_in_prop + " > script").remove()
        } catch(e){
            console.log("clear scripts")
            break
        }
    }

    var node_df = editor.getNodeFromId(id_box_in_prop)
    console.log(node_df)
    try {
        var form_prop = document.querySelector("#node-" + id_box_in_prop + " .form_prop_of_box")
        panel_prop.innerHTML = form_prop.innerHTML
    } catch (e) {
        panel_prop.innerHTML = node_df.html
    }

    while(1){
        try{
            document.querySelector("#node-" + id_box_in_prop + " > script").remove()
        } catch(e){
            console.log("clear scripts")
            break
        }
    }

    var inputs = document.querySelectorAll("#panel_prop  input")
    var textareas = document.querySelectorAll("#panel_prop  textarea")
    var selects = document.querySelectorAll("#panel_prop  select")


    for (var i = 0; i < inputs.length; i++) {
        var elem = inputs[i]
        if (elem.type == "checkbox") {
            elem.checked = false
            if (node_df.data[elem.getAttribute("name")] == "true") {
                elem.checked = true
            }
            continue
        }
        elem.value = node_df.data[elem.getAttribute("name")]
    }
    textarea_code.value =""
    for (var i = 0; i < textareas.length; i++) {
        var elem = textareas[i]
        elem.value = node_df.data[elem.getAttribute("name")]
        code = elem.value
        textarea_code.value = elem.value
        mode_code = elem.getAttribute("name")
    }

    for (var i = 0; i < selects.length; i++) {
        var elem = selects[i]
        elem.value = node_df.data[elem.getAttribute("name")]
    }




    console.log("Node selected " + id);
    var _in = document.getElementById("node-" + id)
    console.log(_in)
    var scripts = _in.getElementsByTagName("script")
    for (var i in scripts) {
        var scriptNode = document.createElement('script');
        scriptNode.innerHTML = scripts[i].text;
        _in.appendChild(scriptNode);
    }
    save()
}

editor.on('nodeSelected', nodeSelectedCustom)

var id_box_in_prop = undefined
editor.on('nodeUnselected', function (flag) {
    try {
        textarea_code.value =""
        console.log('nodeUnselected')
        setFieldNode()

        panel_prop.innerHTML = ""
        panel_prop.style.padding = 0;
        id_box_in_prop = undefined;

        save()
    } finally {
        current_id = undefined;
        bottom_panel.style.animationName = "hide_bottom_panle"
        bottom_panel.style.bottom = "-270px"
    }

})

editor.on('moduleCreated', function (name) {
    console.log("Module Created " + name);
})

editor.on('moduleChanged', function (name) {
    console.log("Module Changed " + name);
})

editor.on('connectionCreated', function (connection) {
    console.log('Connection created');
    console.log(connection);
    var node = editor.getNodeFromId(connection.output_id)
    for (var i in node.outputs) {
        var connections = node.outputs[i].connections
        if (connections.length > 1) {
            var c = connection
            editor.removeSingleConnection(c.output_id, c.input_id, c.output_class, c.input_class)
            // throw "Error length connection output"
        }
    }
    save()

})

editor.on('connectionRemoved', function (connection) {
    console.log('Connection removed');
    console.log(connection);
    save()
})

editor.on('mouseMove', function (position) {
    //console.log('Position mouse x:' + position.x + ' y:' + position.y);
})

editor.on('nodeMoved', function (id) {
    console.log("Node moved " + id);
    save()
})

editor.on('zoom', function (zoom) {
    console.log('Zoom level ' + zoom);
})

editor.on('translate', function (position) {
    console.log('Translate x:' + position.x + ' y:' + position.y);
})

editor.on('addReroute', function (id) {
    console.log("Reroute added " + id);
})

editor.on('removeReroute', function (id) {
    console.log("Reroute removed " + id);
})

editor.on('contextmenu', function (event) {
    console.log(event)
    if (current_id != undefined) {
        copy(current_id)
    }else if (clipboard_node != undefined){
        paste(event.offsetX, event.offsetY)
    }  
})

/* DRAG EVENT */

/* Mouse and Touch Actions */

var elements = document.getElementsByClassName('drag-drawflow');
for (var i = 0; i < elements.length; i++) {
    elements[i].addEventListener('touchend', drop, false);
    elements[i].addEventListener('touchmove', positionMobile, false);
    elements[i].addEventListener('touchstart', drag, false);
}

var mobile_item_selec = '';
var mobile_last_move = null;

function positionMobile(ev) {
    mobile_last_move = ev;
}

function allowDrop(ev) {
    ev.preventDefault();
}

function drag(ev) {
    if (ev.type === "touchstart") {
        mobile_item_selec = ev.target.closest(".drag-drawflow").getAttribute('data-node');
    } else {
        ev.dataTransfer.setData("node", ev.target.getAttribute('data-node'));
    }
}

function drop(ev) {
    if (ev.type === "touchend") {
        var parentdrawflow = document.elementFromPoint(mobile_last_move.touches[0].clientX, mobile_last_move.touches[0].clientY).closest("#drawflow");
        if (parentdrawflow != null) {
            addNodeToDrawFlow(mobile_item_selec, mobile_last_move.touches[0].clientX, mobile_last_move.touches[0].clientY);
        }
        mobile_item_selec = '';
    } else {
        ev.preventDefault();
        var data = ev.dataTransfer.getData("node");
        addNodeToDrawFlow(data, ev.clientX, ev.clientY);
    }

}

function copy (id){
    clipboard_node = editor.getNodeFromId(id)
}
function paste(x, y) {
    editor.addNode("", Object.keys(clipboard_node.inputs).length, Object.keys(clipboard_node.outputs).length, x, y, "", clipboard_node.data, clipboard_node.html)
}

function addNodeToDrawFlow(name, pos_x, pos_y) {
    if (editor.editor_mode === 'fixed') {
        return false;
    }
    pos_x = pos_x * (editor.precanvas.clientWidth / (editor.precanvas.clientWidth * editor.zoom)) - (editor.precanvas.getBoundingClientRect().x * (editor.precanvas.clientWidth / (editor.precanvas.clientWidth * editor.zoom)));
    pos_y = pos_y * (editor.precanvas.clientHeight / (editor.precanvas.clientHeight * editor.zoom)) - (editor.precanvas.getBoundingClientRect().y * (editor.precanvas.clientHeight / (editor.precanvas.clientHeight * editor.zoom)));


    item = modules[name]
    if (item != undefined) {
        var icon = ""
        if (item.icon != undefined) {
            if (item.icon != "" && item.icon.indexOf("fa-") < 0) {
                icon = `<span class="material-icons">${item.icon}</span>`
            }
        }

        var style = ""

        if (item["boxcolor"] != "" && item["boxcolor"] != undefined) {
            style += "background:" + item["boxcolor"] + ";"
        }

        var html_box = `<div>
  

        <div class="title-box">      <div class="mark_color_box" style="${style}"></div>${icon}${item.title}</div>
            <div>
                <input type="text" value="" name="input_box" df-name_box="" class="df_name_box"
                readOnly
                />
            </div>
            </div>
            <div class="form_prop_of_box" style="display:none">
                ${item.html}
            </div>
        `

        if (item.out < 0) {
            Swal.fire({
                title: 'How many outputs do you want?',
                input: 'number',
                inputAttributes: {
                    autocapitalize: 'off'
                },
                showCancelButton: true,
                confirmButtonText: 'Ok',
                showLoaderOnConfirm: true,
                preConfirm: (out) => {
                    if (out == undefined || out == null || out == "")
                        out = 0
                    out = parseInt(out)
                    item = modules[name]
                    id = editor.addNode(item.name, item.in, out, pos_x, pos_y, item.name, item.param, html_box)
                },
                allowOutsideClick: () => !Swal.isLoading()
            })
        } else {
            out = item.out
            if (out == undefined || out == null || out == "")
                out = 0
            out = parseInt(out)
            var id = editor.addNode(item.name, item.in, out, pos_x, pos_y, item.name, item.param, html_box)
            console.log(id)
        }

    }
}

function setFieldNode() {
    console.log(id_box_in_prop)
    if (id_box_in_prop != undefined) {
        var id = id_box_in_prop
        var inputs = document.querySelectorAll("#panel_prop  input")
        var textareas = document.querySelectorAll("#panel_prop  textarea")
        var selects = document.querySelectorAll("#panel_prop  select")

        var field_data = {}
        var old_data = editor.getNodeFromId(id).data
        for (var key in old_data) {
            field_data[key] = old_data[key]
        }

        for (var i = 0; i < inputs.length; i++) {
            var elem = inputs[i]

            if (elem.type == "checkbox") {
                field_data[elem.getAttribute("name")] = "false"
                if (elem.checked) {
                    field_data[elem.getAttribute("name")] = "true"
                }
                continue
            }

            field_data[elem.getAttribute("name")] = elem.value

        }
        for (var i = 0; i < textareas.length; i++) {
            var elem = textareas[i]
            field_data[elem.getAttribute("name")] = elem.value
        }

        for (var i = 0; i < selects.length; i++) {
            var elem = selects[i]
            field_data[elem.getAttribute("name")] = elem.value
        }




        editor.updateNodeDataFromId(id, field_data)
    }
}

var transform = '';

function showpopup(e) {
    e.target.closest(".drawflow-node").style.zIndex = "9999";
    e.target.children[0].style.display = "block";
    //document.getElementById("modalfix").style.display = "block";

    //e.target.children[0].style.transform = 'translate('+translate.x+'px, '+translate.y+'px)';
    transform = editor.precanvas.style.transform;
    editor.precanvas.style.transform = '';
    editor.precanvas.style.left = editor.canvas_x + 'px';
    editor.precanvas.style.top = editor.canvas_y + 'px';
    //console.log(transform);

    //e.target.children[0].style.top  =  -editor.canvas_y - editor.container.offsetTop +'px';
    //e.target.children[0].style.left  =  -editor.canvas_x  - editor.container.offsetLeft +'px';
    editor.editor_mode = "fixed";

}

function closemodal(e) {
    e.target.closest(".drawflow-node").style.zIndex = "2";
    e.target.parentElement.parentElement.style.display = "none";
    //document.getElementById("modalfix").style.display = "none";
    editor.precanvas.style.transform = transform;
    editor.precanvas.style.left = '0px';
    editor.precanvas.style.top = '0px';
    editor.editor_mode = "edit";
}

function changeModule(event) {
    var all = document.querySelectorAll(".menu ul li");
    for (var i = 0; i < all.length; i++) {
        all[i].classList.remove('selected');
    }
    event.target.classList.add('selected');
}

function changeMode(option) {

    //console.log(lock.id);
    if (option == 'lock') {
        lock.style.display = 'none';
        unlock.style.display = 'block';
    } else {
        lock.style.display = 'block';
        unlock.style.display = 'none';
    }

}



function open_coder(elem, mode) {
    if (mode == undefined) {
        mode = "javascript"
    }
    var code = elem.previousElementSibling.value
    textarea = elem.previousElementSibling

    Swal.fire({
        title: 'Coder',
        html: `<textarea id="code1" >${code}</textarea>`,
        width: '8000px',
        heigth: '8000px',
        preConfirm: () => {
            textarea = elem.previousElementSibling
            textarea.value = editor.getValue()
            console.log(textarea)
            save()
        }
    }).then((result) => {
        /* Read more about isConfirmed, isDenied below */

    })
    var editor = codemirror_new(code1, mode)

}

function save_template(name, code_template){

    var _datos = {
        name: name,
        content : code_template
    }

    fetch("/nflow/templates", {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(_datos),
    });

 
}

function delete_template(name){

   

    fetch("/nflow/templates/"+name, {
      method: "DELETE",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    }).then (response => response.json())
    .then(data =>{
        nodeSelectedCustom(id_box_in_prop);
    })

 
}



function create_template(name, code_template){

    var _datos = {
        name: name,
        content : code_template
    }

    fetch("/nflow/templates", {
      method: "PUT",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify(_datos),
    }).then (response => response.json())
    .then(data =>{
        console.log(data);
        nodeSelectedCustom(id_box_in_prop);
    })

 
}
function open_template(elem,name) {
    mode = "htmlmixed"
    

    fetch("/nflow/templates/"+name)
    .then(response => response.json())
    .then(template => {
        Swal.fire({
            title: 'Coder',
            html: `<textarea id="code1" >${ template.content}</textarea>`,
            width: '8000px',
            heigth: '8000px',
            preConfirm: () => {
                if (confirm("Are you sure to save the template?\nThis template will change for all boxes.") == true) {
                    console.log(editor.getValue())
                    save_template(name,editor.getValue())
                  } else {
                  }
                
            }
        }).then((result) => {
            /* Read more about isConfirmed, isDenied below */
    
        })
        var editor = codemirror_new(code1, mode)
        
    })
}



var words = undefined
fetch("nflow/ui/intellisense").then(response => response.json()).then(data => { words = data["js_words"] })

function codemirror_new(code1, mode, readOnly) {
    CodeMirror.hint.javascript = function (cm) {

        var list = words
        var cursor = editor.getCursor();
        var currentLine = editor.getLine(cursor.line);
        var start = cursor.ch;
        var end = start;
        while (end < currentLine.length && /[\w$]+/.test(currentLine.charAt(end))) ++end;
        while (start && /[\w$]+/.test(currentLine.charAt(start - 1))) --start;
        var curWord = start != end && currentLine.slice(start, end);
        var regex = new RegExp('^' + curWord, 'i');
        var result = {
            list: (!curWord ? list : list.filter(function (item) {
                return item.match(regex);
            })).sort(),
            from: CodeMirror.Pos(cursor.line, start),
            to: CodeMirror.Pos(cursor.line, end)
        };


        //var inner = { from: cm.getCursor(), to: cm.getCursor(), list: words_filter };





        return result;
    };

    var autocomplete = false
    if (mode == "javascript") {
        autocomplete = "autocomplete"
    }

    var editor = CodeMirror.fromTextArea(code1, {
        lineNumbers: true,
        lineWrapping: false,
        mode: mode,
        theme: "monokai",
        showFormatButton: true,
        matchBrackets: true,
        autoCloseBrackets: true,
        showCursorWhenSelecting: true,
        continueComments: "Enter",
        keyMap: "sublime",
        foldGutter: true,
        gutters: ["CodeMirror-linenumbers", "CodeMirror-foldgutter"],
        extraKeys: {
            "Ctrl-Space": autocomplete,
            "Ctrl-S": function (cm) {
                textarea.value = cm.getValue()
                console.log(textarea)
                save(true)
            },
            "Cmd-S": function (instance) {
                textarea.value = cm.getValue()
                console.log(textarea)
                save(true)
            },
            "Ctrl-7": "toggleComment",
            "Shift-Ctrl-F": function (cm) {
                cm.setOption("fullScreen", !cm.getOption("fullScreen"));
            }
        }
    });



    if (readOnly == true) {
        console.log("readOnly")
        editor.options.readOnly = true
    }

    editor.setSize("100%", "600");
    editor.setOption("fullScreen", true);

    setTimeout(() => {
        editor.setOption("fullScreen", false);
        document.getElementsByClassName("swal2-confirm swal2-styled")[0].onfocus = function (ev) {
            editor.focus(ev)
        }
    }, 500);
    return editor
}

function save(show,callfx) {
    console.log("setFieldNode")
    setFieldNode()
    if (flag_running_save) {
        return
    }
    try {
        flag_running_save = true
        var dflow = editor.export()

        elemens = Array.prototype.forEach.call(document.getElementsByClassName("modules"),
            module => {
                key_workspace = module.getAttribute('name')
                if (key_workspace == undefined) {
                    return
                }

                if (dflow.drawflow[key_workspace] == undefined) {
                    dflow.drawflow[key_workspace] = { "data": {} }
                }

                for (key in dflow.drawflow[key_workspace].data) {
                    try {
                        item = dflow.drawflow[key_workspace].data[key]
                        var elem = document.getElementById("node-" + key)
                        if (elem == null) {
                            continue
                        }
                        for (data_key in item.data) {

                            var sub_elems = null
                            try {
                                sub_elems = elem.querySelectorAll("[name='" + data_key + "']")
                                console.log(key)
                            } catch (e) {
                                console.log(key)
                                console.log(e)
                                continue
                            }
                            try {
                                if (sub_elems[0] instanceof Element) {
                                    item.data[data_key] = sub_elems[0].value;
                                }
                            } catch (e) {
                                console.log(sub_elems)
                                console.log(e)
                            }


                        }
                    } catch (e) {
                        console.log(e)
                        console.log(key)
                    }
                }
            }
        )



        var json_str = JSON.stringify(dflow, 4, 4)
        console.log(dflow)


        fetch('nflow/app', {
            method: 'POST',
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            body: json_str
        }).then(response => {
            return response.json()
        })
            .then(data => {
                if (data.msg == "ok") {

                    console.log("saved!")

                    if (callfx){
                        flag_running_save = false
                        callfx()
                    }

                    if (show) {
                        Swal.fire({
                            position: 'top-end',
                            icon: 'success',
                            title: 'Your work has been saved',
                            showConfirmButton: false,
                            timer: 1500
                        })
                    }
                } else {
                    Swal.fire({
                        title: 'Save',
                        html: `<pre>${data.msg}</pre>`
                    })
                }
            })


    } catch (e) {
        console.log(e)
    } finally {
        console.log("save flag_running_save = false")
        flag_running_save = false
    }


}

function open_data(text, mode, readOnly) {
    try {
        var j = JSON.parse(text)
        text = "<textarea id='codejson'>" + JSON.stringify(j, 0, 4) + "</textarea>"
    } catch (e) { }



    Swal.fire({
        title: 'Result',
        html: `<pre style="text-align: left;">${text}</pre>`
    })
    return codemirror_new(codejson, mode, readOnly)
}


function add_workspace() {
    Swal.fire({
        title: 'Name?',
        input: 'text',
        inputAttributes: {
            autocapitalize: 'off'
        },
        showCancelButton: true,
        confirmButtonText: 'Ok',
        showLoaderOnConfirm: true,
        preConfirm: (name) => {
            editor.addModule(name);
            tabs = document.querySelector(".top_panel > .panel_items")
            tabs.innerHTML += `<button id="module_${name}" class="panel_item" onclick="editor.changeModule('${name}'); changeModule(event);" >${name}</button>`;

            elem = document.getElementById(`module_${name}`)
            elem.click()
            setFieldNode()
            save()
        },
        allowOutsideClick: () => !Swal.isLoading()
    })

}

function remove_workspace() {
    var name = editor.module;
    editor.removeModule(name)
    document.getElementById(`module_${name}`).remove();
    save()
}

init()