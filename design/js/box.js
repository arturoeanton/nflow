var modules = {}

function show_icon() {
  elem = document.getElementById("box_icon")
  span_show_icon.innerHTML = "<i class='" + elem.value + "'></i>"
}

function create_box(module) {
  if (module == undefined) {
    module = ""
  }


  Swal.fire({
    title: `Box: ${module == "" ? "new" : module}`,
    html:
      `
        <div class="w3-bar" style="background:#19191a">
        <button class="w3-bar-item w3-button tablink w3-sand" onclick="openTab(event,'Manifest')">Manifest</button>
        <button class="w3-bar-item w3-button tablink" onclick="openTab(event,'HTML')">HTML</button>
        <button class="w3-bar-item w3-button tablink" onclick="openTab(event,'Code')">Code</button>
        </div>
        <div id="Manifest" class="w3-container open">
            <br/>
            Title: <input type="text" id="box_title" style="width: 80px;"> 
            Box Color: <input type="text" id="box_color" value="darkorange"  style="width: 80px;"> 
            Icon: <input type="text" id="box_icon" style="width: 80px;">
            In:  <input type="number" id="box_in" style="width: 25px;"> 
            Out: <input type="number" id="box_out" style="width: 25px;">
            Custom:  <input type="checkbox" id="box_custom" >
            <br/>
            <br/>

            <textarea id='code_manifest'></textarea>


        </div>
        <div id="HTML" class="w3-container open" style="display:none">
            <br/>
            <textarea id='code_html'></textarea>
        </div>
        <div id="Code" class="w3-container open" style="display:none">
            <br/>
            <textarea id='code_code'></textarea>
        </div>
        
        `,
    focusConfirm: false,
    onClose: () => {
      cm.getWrapperElement().parentNode.removeChild(cm.getWrapperElement());
      ch.getWrapperElement().parentNode.removeChild(ch.getWrapperElement());
      cj.getWrapperElement().parentNode.removeChild(cj.getWrapperElement());

      console.log('close!');

    },
    preConfirm: () => {
      if (module == "" || module == undefined || module == null) {
        module = prompt("Please enter the name of Box", "");
        if (module == "" || module == null || module == undefined) {
          return false
        }
      }


      var json = {}
      json["title"] = box_title.value
      json["in"] = parseInt(box_in.value)
      json["out"] = parseInt(box_out.value)
      json["icon"] = box_icon.value
      json["boxcolor"] = box_color.value

      json["custom"] = box_custom.checked

      json["param"] = {}
      var param = JSON.parse(cm.getDoc().getValue())
      for (var i in param) {
        json["param"][i] = param[i]
      }
      json["param"]["type"] = "js"
      json["param"]["script"] = module

      var elementTemp = document.createElement("p")
      elementTemp.innerHTML = ch.getDoc().getValue()
      getParamAndSetNames(elementTemp, json["param"])
      ch.getDoc().setValue(elementTemp.innerHTML)
      elementTemp.remove()


      fetch("nflow/module/manifest/" + module, {
        method: 'post',
        body: JSON.stringify(json, 4, 4)
      }).then(response => {
        fetch("nflow/module/box/" + module, {
          method: 'post',
          body: ch.getDoc().getValue()
        }).then(response => {
          fetch("nflow/module/code/" + module, {
            method: 'post',
            body: cj.getDoc().getValue()
          }).then(response => {
            filter_box()
            cm=null;
            ch=null;
            cj=null;
            Swal.fire({
              position: 'top-end',
              icon: 'success',
              title: 'Your work has been saved',
              showConfirmButton: false,
              timer: 1500
            })
          })
        })
      })

      return
    }
  })


  cm = codemirror_new(code_manifest, "application/ld+json", true)
  ch = codemirror_new(code_html, "htmlmixed")
  cj = codemirror_new(code_code, "javascript")

  if (module != "") {
    fetch("nflow/module/manifest/" + module).then(response => response.text()).then(data => {
      var json = JSON.parse(data)
      box_title.value = json["title"]
      box_in.value = json["in"]
      box_out.value = json["out"]
      box_icon.value = json["icon"]
      box_custom.checked = json["custom"]
      box_color.value = json["boxcolor"]





    })
    fetch("nflow/module/box/" + module).then(response => response.text()).then(data => {
      var elementTemp = document.createElement("p")
      elementTemp.innerHTML = data
      var param = {}
      getParamAndSetNames(elementTemp, param)
      ch.getDoc().setValue(elementTemp.innerHTML)
      elementTemp.remove()
      cm.getDoc().setValue(JSON.stringify(param, 4, 4));
    })
    fetch("nflow/module/code/" + module).then(response => response.text()).then(data => { cj.getDoc().setValue(data); })
  } else {
    cm.getDoc().setValue(`{
}`);
    ch.getDoc().setValue(`Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<div class="nflow_auth">Auth:<input type="checkbox" df-nflow_auth="false" name="nflow_auth"></div>
<hr/>

<!-- add your form -->
    `);
    cj.getDoc().setValue(`function main(){
}`);
  }
}


function openTab(evt, tabName) {
  var i;
  var x = document.getElementsByClassName("open");
  for (i = 0; i < x.length; i++) {
    x[i].style.display = "none";
  }

  tablinks = document.getElementsByClassName("tablink");
  for (i = 0; i < x.length; i++) {
    tablinks[i].className = tablinks[i].className.replace(" w3-sand", "");
  }
  evt.currentTarget.className += " w3-sand";

  var elementTemp = document.createElement("p")
  elementTemp.innerHTML = ch.getDoc().getValue()
  var param = {}
  getParamAndSetNames(elementTemp, param)
  ch.getDoc().setValue(elementTemp.innerHTML)
  cm.getDoc().setValue(JSON.stringify(param, 4, 4));
  elementTemp.remove()

  document.getElementById(tabName).style.display = "block";
  cm.refresh()
  ch.refresh()
  cj.refresh()

}

menu_input_search.addEventListener("keyup", function (event) {
  filter_box()
  // Number 13 is the "Enter" key on the keyboard
  if (event.keyCode === 13) {
    // Cancel the default action, if needed
    event.preventDefault();
    // Trigger the button element with a click
    filter_box()
  }
});

function filter_box() {
  var f = menu_input_search.value
  fetch('nflow/modules')
    .then(response => response.json())
    .then(function (data) {
      modules_bar = document.querySelector(".menu > .menu_items")
      modules_bar.innerHTML = ""

      for (const key in data) {
        data.name = key
        data.html = ""
        modules[key] = data[key]
        f = f.toUpperCase()
        re = new RegExp(`${f}.*`)
        if (data[key].title.toUpperCase().match(re)) {
          var editable = ""
          if (data[key].editable) {
            editable = `
                        <span class="material-icons" style="left: 230px;    cursor: pointer;" onclick="deleteModule('${key}')" >delete</span>
                        <span class="material-icons" style="left: 250px;  cursor: pointer;" onclick="create_box('${key}')">edit</span>
                        `
          }
          style = ""
          if (data[key]["boxcolor"] != "" && data[key]["boxcolor"] != undefined) {
            style = "background:" + data[key]["boxcolor"] + ";"
          }
          modules_bar.innerHTML += ` 
                        <div class="drag-drawflow menu_item" draggable="true" ondragstart="drag(event)" data-node="${key}">
                        <div class="mark_menu_color_box" style="${style}"></div>
                        <span style="margin:10px; text-overflow: ellipsis;white-space: nowrap;">${editable}${data[key].title}</span>
                        </div>`

          modules[key].html = data[key]["HTMLForm"]
        }
      }
    }).catch(function (err) {
      console.warn('Load list modules', err);
    });

}

function deleteModule(module) {

  Swal.fire({
    title: 'Do you want to delete Box?',
    showDenyButton: true,
    showCancelButton: true,
    confirmButtonText: `Delete`,
    denyButtonText: `Don't Delete`,
  }).then((result) => {
    /* Read more about isConfirmed, isDenied below */
    if (result.isConfirmed) {
      fetch("nflow/module/" + module, {
        method: 'delete'
      }).then(response => {
        filter_box()
      })
    } else if (result.isDenied) {

    }
  })



}


function getParamAndSetNames(element, param) {
  if (param == undefined) param = {}
  for (var i in element.children) {
    var elem = element.children[i]
    for (var j in element.children) {
      param = getParamAndSetNames(elem, param)
    }
    try {
      var attributes = elem.getAttributeNames()
    } catch (e) {
      continue
    }
    for (var j in attributes) {
      var attr = attributes[j]
      if (attr.startsWith("df-")) {
        param[attr.substr(3)] = elem.getAttribute(attr)
        elem.setAttribute("name", attr.substr(3))
      }
    }
  }
  return param
}
menu_input_search.value = ""
filter_box()