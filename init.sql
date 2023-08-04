CREATE TABLE IF NOT EXISTS "modules" (
  "id" SERIAL PRIMARY KEY,
  "name" text NOT NULL,
  "form" text NOT NULL,
  "mod" text NOT NULL,
  "code" text
);

CREATE TABLE IF NOT EXISTS "nflow" (
  "id" SERIAL PRIMARY KEY,
  "json" text NOT NULL,
  "name" text NOT NULL,
  "default_js" text
);

INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, '_end', '<hr>

<!-- add your form -->
    ', '{
    "title": "End",
    "in": 1,
    "out": null,
    "icon": "",
    "boxcolor": "darkorange",
    "custom": false,
    "param": {
        "type": "js",
        "script": "_end"
    }
}', 'function main(){
  c.JSON(200,{})
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, '_goroutine', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr/>', '{
    "title": "Parallel",
    "in": 1,
    "out": 2,
    "icon": "share",
    "custom": false,
    "editable":false,
    "boxcolor": "darkmagenta",
    "fontcolor": "white",
    "param": {
        "type": "gorutine",
        "name_box":"",
        "nflow_auth":"false"
    }
}', 'NULL');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, '_js_code', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr/>
    <textarea class="hide_component" name="code" df-code></textarea>
    <button class="btn-export" name="open" onclick=''open_coder(this)''> CODE </button>', '{
    "title": "JS Code",
    "icon": "source",
    "in": 1,
    "out": -1,
    "custom":false,
    "editable":false,
    "boxcolor": "darkred",
    "fontcolor": "white",
    "param": {
        "type":"js",
        "code":"",
        "name_box":"",
        "nflow_auth":"false"
    }
}', 'NULL');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, '_start_form', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<div class="nflow_auth">Auth:<input type="checkbox" df-nflow_auth="false" name="nflow_auth"></div>
<hr/>
<p>Url</p>
<input type="text"  name="urlpattern" df-urlpattern>
<button  class="btn-export" onclick=''setFieldNode();save();window.open(window.location+"/.."+getElemenByNameOfChiled(this.parentElement,"urlpattern")[0].value)''> OPEN </button>', '{
    "title": "Form Start",
    "in": 0,
    "out": 1,
    "icon": "video_settings",
    "custom": false,
    "editable":false,
    "boxcolor": "darkcyan",
    "fontcolor": "white",
    "param": {
        "type": "starter",
        "method": "ANY",
        "urlpattern":"",
        "name_box":"",
        "nflow_auth":"false"
    }
}', 'NULL');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'js_foreach', 'js_foreach', '{
    "title": "For each",
    "in": 1,
    "out": 2,
    "icon": "",
    "boxcolor": "darkred",
    "custom": false,
    "param": {
        "name_box": "",
        "nflow_auth": "false",
        "field": "",
        "item_name": "current_element",
        "type": "js",
        "script": "js_foreach"
    }
}', 'function main(){
  try{
    if (payload["__foreach"] == undefined || payload["__foreach"] == null){
      payload["__foreach"] = find_element(nflow_data["field"], payload)
    }
  	var current_list = payload["__foreach"]
    console.log("*>"+current_list.length)
    if (current_list.length == 0){
      next="output_1"
      return
    }
    next = "output_2"
    payload[nflow_data["item_name"]] = current_list.shift()
    payload["__foreach"]=current_list
	console.log("+>"+payload["__foreach"].length)
  }catch(e){
    next="output_1"
    return
  }
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, '_start_http', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<div class="nflow_auth">Auth:<input type="checkbox" df-nflow_auth="false" name="nflow_auth"></div>
<hr/>
<p>Url</p>
<input type="text"  name="urlpattern" df-urlpattern>

<p>Choose Method:</p>
<select name="method" df-method>
  <option value="GET">GET</option>
  <option value="POST">POST</option>
  <option value="PUT">PUT</option>
  <option value="PATCH">PATCH</option>
  <option value="DELETE">DELETE</option>
  <option value="ANY">ANY</option>
</select>
<button  class="btn-export" onclick=''setFieldNode();save();window.open(window.location+"/.."+getElemenByNameOfChiled(this.parentElement,"urlpattern")[0].value)''> OPEN </button>
</div>', '{
    "title": "Http Start",
    "in": 0,
    "out": 1,
    "icon": "http",
    "custom": false,
    "editable":false,
    "boxcolor": "darkcyan",
    "fontcolor": "white",
    "param": {
        "type": "starter",
        "method": "ANY",
        "urlpattern":"",
        "name_box":"",
        "nflow_auth":"false"
    }
}', 'NULL');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'form', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
<textarea class="hide_component" name="html_form" df-html_form="<div><input type=''input'' name=''n1''/><input type=''submit''/></div>"> </textarea>
<button class="btn-export" onclick="open_coder(this, `htmlmixed`)"> CODE</button>
<input type="input" df-template="form1.html" name="template">', '{
    "title": "Form",
    "in": 1,
    "out": 1,
    "icon": "book",
    "boxcolor": "darkorange",
    "custom": false,
    "param": {
        "name_box": "",
        "html_form": "<div><input type=''input'' name=''n1''/><input type=''submit''/></div>",
        "template": "form1.html",
        "type": "js",
        "script": "form"
    }
}', 'function main(){
  var action = "/nflow/node/run/"+__flow_name+"/"+__outputs["output_1"]
  var html_form = dromedary_data["html_form"]
  var template1 = dromedary_data["template"]  
  var code = template(html_form,{
        form_code:html_form, 
        action:action,
    	next:__outputs["output_1"],
    	payload:payload,
    	profile:get_profile()
  });
  var html = template(file_to_string("templete/" + template1), 
  {
        form_code:code, 
        action:action,
    	next:__outputs["output_1"],
    	payload:payload,
    	profile:get_profile()
  });
  c.HTML(200,html)
  
  if (payload == undefined || payload == null){
    payload = {}
  }
  payload["break"] = true
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'js_set_var', 'Description:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
Name:<input type="text" name="var_name" df-var_name="name1" class="nflow_field"> <br>
Value:<input type="text" name="var_value" df-var_value="" class="nflow_field"><br>
    ', '{
    "title": "Set Variable",
    "in": 1,
    "out": 1,
    "icon": "",
    "boxcolor": "darkred",
    "custom": false,
    "param": {
        "name_box": "",
        "var_name": "name1",
        "var_value": "",
        "type": "js",
        "script": "js_set_var"
    }
}', 'function main(){
  if (payload == undefined || payload == null ) {
    payload = {}
  }
  payload[nflow_data["var_name"]] = nflow_data["var_value"]
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'if_major', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>
<!-- add your form -->
Field:<input type="text" name="field" df-field="" class="nflow_field"><br>
Limit:<input type="number" name="limit" df-limit="0" class="nflow_field">    ', '{
    "title": "if higher/lower",
    "in": 1,
    "out": 2,
    "icon": "",
    "boxcolor": "darkred",
    "custom": false,
    "param": {
        "name_box": "",
        "field": "",
        "limit": "0",
        "type": "js",
        "script": "if_major"
    }
}', 'function main(){
  if (parseFloat(payload[nflow_data["field"]]) > parseFloat(nflow_data["limit"])){
    next= "output_1"
    return
  } 
  next="output_2"
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'if_equals', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>
<!-- add your form -->
Field:<input type="text" name="field" df-field="" class="nflow_field"><br>
Value:<input type="text" name="value" df-value="" class="nflow_field">    ', '{
    "title": "if equals",
    "in": 1,
    "out": 2,
    "icon": "",
    "boxcolor": "darkred",
    "custom": false,
    "param": {
        "name_box": "",
        "field": "",
        "value": "",
        "type": "js",
        "script": "if_equals"
    }
}', 'function main(){
  if (payload[nflow_data["field"]] == nflow_data["value"]){
    next= "output_1"
    return
  } 
  next="output_2"
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'json', '<input type="text" name="field" df-field="">
<input type="number" name="http_code" df-http_code="200">', '{
    "title": "JSON Render",
    "in": 1,
    "out": 0,
    "icon": "source",
    "boxcolor": "darkviolet",
    "custom": false,
    "param": {
        "field": "",
        "http_code": "200",
        "type": "js",
        "script": "json"
    }
}', 'function main(){
    var code = 200
    if (nflow_data["http_code"] != ""){
        code = parseInt(nflow_data["http_code"])
    }
    if (nflow_data["field"] == undefined || nflow_data["field"] =="" || nflow_data["field"] =="*"){
        c.JSON(code, payload)
        return
    }
    c.JSON(code,find_element(nflow_data["field"],payload))
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'parser_query_param', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
', '{
    "title": "Param to Payload",
    "in": 1,
    "out": 1,
    "icon": "filter_alt",
    "boxcolor": "darkred",
    "custom": false,
    "param": {
        "name_box": "",
        "nflow_auth": "false",
        "type": "js",
        "script": "parser_query_param"
    }
}', 'function main (){
    if (payload == undefined){
      payload = {}
    } 
    var params = url_values_to_map(c.QueryParams())
    for (key in params){
      payload[key] = c.QueryParam(key)
    }
  }');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'mustache', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
    
<textarea class="hide_component" name="template" df-template=""></textarea>
<button class="btn-export" onclick="open_coder(this, `htmlmixed`)"> CODE </button>', '{
    "title": "Mustache",
    "in": 1,
    "out": 0,
    "icon": "code",
    "boxcolor": "darkorange",
    "custom": false,
    "param": {
        "name_box": "",
        "nflow_auth": "false",
        "template": "",
        "type": "js",
        "script": "mustache"
    }
}', 'function main(){
    console.log(JSON.stringify(payload))
    c.HTML(200, mustache(dromedary_data["template"], payload))
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'Logout', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
    ', '{
    "title": "logout",
    "in": 1,
    "out": null,
    "icon": "",
    "boxcolor": "darkorange",
    "custom": false,
    "param": {
        "name_box": "",
        "nflow_auth": "false",
        "type": "js",
        "script": "Logout"
    }
}', 'function main(){
  delete_profile()
  c.HTML(200,"logout <a href=''/home'' >Home</a>")
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'template', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
<textarea class="hide_component" name="template" df-template=""></textarea>
<button class="btn-export" onclick="open_coder(this, `htmlmixed`)"> CODE </button>', '{
    "title": "Template",
    "in": 1,
    "out": 0,
    "icon": "code",
    "boxcolor": "darkorange",
    "custom": false,
    "param": {
        "name_box": "",
        "nflow_auth": "false",
        "template": "",
        "type": "js",
        "script": "template"
    }
}', 'function main(){
    console.log(JSON.stringify(payload))
    c.HTML(200, template(dromedary_data["template"], payload))
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'Validate Token', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
Token: <input type="input" name="token" df-token="test1">', '{
    "title": "",
    "in": 1,
    "out": 1,
    "icon": "",
    "boxcolor": "darkorange",
    "custom": false,
    "param": {
        "name_box": "",
        "token": "test1",
        "type": "js",
        "script": "Validate Token"
    }
}', 'function main(){
  next= "output_1"
  if (header["Authorization"] != "Bearer " + nflow_data["token"]){
    c.JSON(401,{"error":"Unauthorized"})
    next= "break"
  }
  
}');
INSERT INTO modules (id, name, form, mod, code) VALUES (DEFAULT, 'Validate Login', 'Name:<input type="text" name="name_box" df-name_box="" class="nflow_field">
<hr>

<!-- add your form -->
    ', '{
    "title": "Validate Login",
    "in": 1,
    "out": null,
    "icon": "",
    "boxcolor": "darkorange",
    "custom": false,
    "param": {
        "name_box": "",
        "type": "js",
        "script": "Validate Login"
    }
}', 'function main(){
  	var url_back = ""+get_session("auth-session","redirect_url")
  	var flag = false
  	if (payload["username"] == "user1") {
      if (payload["password"] == "1") {
      	flag = true
      }
    }
  
  	if (payload["username"] == "admin") {
      if (payload["password"] == "admin") {
      	flag = true
      }
    }
  	
  	if (flag) {
      	set_profile({"username":payload["username"]})
  		return c.HTML(200," <script>window.location.href = ''"+url_back+"''</script>")
    }
    return c.HTML(200," Error Login <br/>  <a href=''/home'' >Home</a>")

}');

INSERT INTO nflow (id, json, name, default_js) VALUES (DEFAULT, '{    "drawflow": {        "": {            "data": {}        },        "Home": {            "data": {                "090b19db-2cd3-4785-a68f-7c776b574783": {                    "class": "",                    "data": {                        "method": "ANY",                        "name_box": "",                        "nflow_auth": "false",                        "type": "starter",                        "urlpattern": "/login_flow1"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkcyan;\"></div><span class=\"material-icons\">http</span>Http Start</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"http\"></i> Http Start</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n<p>Url</p>\n<input type=\"text\"  name=\"urlpattern\" df-urlpattern>\n\n<p>Choose Method:</p>\n<select name=\"method\" df-method>\n  <option value=\"GET\">GET</option>\n  <option value=\"POST\">POST</option>\n  <option value=\"PUT\">PUT</option>\n  <option value=\"PATCH\">PATCH</option>\n  <option value=\"DELETE\">DELETE</option>\n  <option value=\"ANY\">ANY</option>\n</select>\n<button  class=\"btn-export\" onclick=''setFieldNode();save();window.open(window.location+\"/..\"+getElemenByNameOfChiled(this.parentElement,\"urlpattern\")[0].value)''> OPEN </button>\n</div>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "090b19db-2cd3-4785-a68f-7c776b574783",                    "inputs": {},                    "name": "",                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "fa94086f-91ee-4960-86df-43fa65f8475c",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 289,                    "pos_y": 268,                    "typenode": false                },                "10a96bc6-3816-4013-8931-1b6ff699c6d9": {                    "data": {                        "method": "ANY",                        "name_box": "Home",                        "nflow_auth": "true",                        "type": "starter",                        "urlpattern": "/home"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkcyan;\"></div><span class=\"material-icons\">video_settings</span>Form Start</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"video_settings\"></i> Form Start</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n<p>Url</p>\n<input type=\"text\"  name=\"urlpattern\" df-urlpattern>\n<button  class=\"btn-export\" onclick=''setFieldNode();save();window.open(window.location+\"/..\"+getElemenByNameOfChiled(this.parentElement,\"urlpattern\")[0].value)''> OPEN </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "10a96bc6-3816-4013-8931-1b6ff699c6d9",                    "inputs": {},                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "307699c8-243b-4518-8aaf-478931e33709",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 289,                    "pos_y": 156,                    "typenode": false                },                "19ad1116-1095-421b-870c-2b8808490039": {                    "class": "",                    "data": {                        "html_form": "{{index .profile \"username\"}}\n<div>\n  <input type=''input'' name=''n1''/>\n  <input type=''input'' name=''n3''/>\n  <input type=''submit''/>\n</div> ",                        "name_box": "",                        "nflow_auth": "false",                        "script": "form",                        "template": "form1.html",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div><span class=\"material-icons\">book</span>Form</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"book\"></i> Form</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr>\n\n<!-- add your form -->\n<textarea class=\"hide_component\" name=\"html_form\" df-html_form=\"<div><input type=''input'' name=''n1''/><input type=''submit''/></div>\"> </textarea>\n<button class=\"btn-export\" onclick=\"open_coder(this, `htmlmixed`)\"> CODE</button>\n<input type=\"input\" df-template=\"form1.html\" name=\"template\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "19ad1116-1095-421b-870c-2b8808490039",                    "inputs": {                        "input_1": {                            "connections": []                        }                    },                    "name": "",                    "outputs": {                        "output_1": {                            "connections": []                        }                    },                    "pos_x": -9477,                    "pos_y": -9819,                    "typenode": false                },                "307699c8-243b-4518-8aaf-478931e33709": {                    "data": {                        "html_form": "<div>\n  <input type=''input'' name=''n1''/>\n  <input type=''input'' name=''n3''/>\n  <input type=''submit''/>\n  Nuevo\n</div>",                        "name_box": "Form Datos1",                        "nflow_auth": "false",                        "script": "form",                        "template": "form1.html",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div><span class=\"material-icons\">book</span>Form</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"book\"></i> Form</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr>\n\n<!-- add your form -->\n<textarea class=\"hide_component\" name=\"html_form\" df-html_form=\"<div><input type=''input'' name=''n1''/><input type=''submit''/></div>\"> </textarea>\n<button class=\"btn-export\" onclick=\"open_coder(this, `htmlmixed`)\"> CODE</button>\n<input type=\"input\" df-template=\"form1.html\" name=\"template\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "307699c8-243b-4518-8aaf-478931e33709",                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "input": "output_1",                                    "node": "10a96bc6-3816-4013-8931-1b6ff699c6d9"                                }                            ]                        }                    },                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "c5221c80-92d0-41bc-98e0-c9d117a1b2e0",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 533,                    "pos_y": 145,                    "typenode": false                },                "5e9ff886-4991-4795-89c9-e02bb71bd6ed": {                    "class": "",                    "data": {                        "name_box": "",                        "nflow_auth": "false",                        "script": "Validate Login",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div>Validate Login</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"\"></i> Validate Login</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr>\n\n<!-- add your form -->\n    \n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "5e9ff886-4991-4795-89c9-e02bb71bd6ed",                    "inputs": {},                    "name": "",                    "outputs": {},                    "pos_x": -9141,                    "pos_y": -9371,                    "typenode": false                },                "66308384-c966-4dee-af14-b5a65177ad19": {                    "data": {                        "code": "function main(){\n\tc.JSON(200,{\"idp\":\"nflow2\", \"ss\":\"sss\"\n              })\n}",                        "name_box": "",                        "nflow_auth": "false",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkred;\"></div><span class=\"material-icons\">source</span>JS Code</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"source\"></i> JS Code</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n    <textarea class=\"hide_component\" name=\"code\" df-code></textarea>\n    <button class=\"btn-export\" name=\"open\" onclick=''open_coder(this)''> CODE </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "66308384-c966-4dee-af14-b5a65177ad19",                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "input": "output_1",                                    "node": "d94dae41-33bc-437a-86a9-d5951a538f7b"                                }                            ]                        }                    },                    "outputs": {},                    "pos_x": 605,                    "pos_y": 421,                    "typenode": false                },                "79db20d7-a350-4bd3-87f8-0b0468313e20": {                    "data": {                        "html_form": "<div><input type=''input'' name=''n2''/><input type=''submit''/></div>",                        "name_box": "Form Datos 2",                        "nflow_auth": "false",                        "script": "form",                        "template": "form1.html",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div><span class=\"material-icons\">book</span>Form</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"book\"></i> Form</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr>\n\n<!-- add your form -->\n<textarea class=\"hide_component\" name=\"html_form\" df-html_form=\"<div><input type=''input'' name=''n1''/><input type=''submit''/></div>\"> </textarea>\n<button class=\"btn-export\" onclick=\"open_coder(this, `htmlmixed`)\"> CODE</button>\n<input type=\"input\" df-template=\"form1.html\" name=\"template\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "79db20d7-a350-4bd3-87f8-0b0468313e20",                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "input": "output_1",                                    "node": "c5221c80-92d0-41bc-98e0-c9d117a1b2e0"                                }                            ]                        }                    },                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "7d42d888-3cbe-4033-8f40-3ef3f272bc58",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 1069,                    "pos_y": 135,                    "typenode": false                },                "7d42d888-3cbe-4033-8f40-3ef3f272bc58": {                    "data": {                        "field": "",                        "http_code": "200",                        "script": "json",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkviolet;\"></div><span class=\"material-icons\">source</span>JSON Render</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"source\"></i> JSON Render</div>\n\t\t<div class=\"box\">\n\t\t  <input type=\"text\" name=\"field\" df-field=\"\">\n<input type=\"number\" name=\"http_code\" df-http_code=\"200\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "7d42d888-3cbe-4033-8f40-3ef3f272bc58",                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "input": "output_1",                                    "node": "79db20d7-a350-4bd3-87f8-0b0468313e20"                                },                                {                                    "input": "output_2",                                    "node": "c5221c80-92d0-41bc-98e0-c9d117a1b2e0"                                }                            ]                        }                    },                    "outputs": {},                    "pos_x": 1379,                    "pos_y": 125,                    "typenode": false                },                "b913cba4-e239-40b6-84e9-7fcf3a5237f9": {                    "class": "",                    "data": {                        "code": "function main(){\n\tc.JSON(200,{\"idp\":\"nflow2\"})\n}",                        "name_box": "",                        "nflow_auth": "false",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkred;\"></div><span class=\"material-icons\">source</span>JS Code</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"source\"></i> JS Code</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n    <textarea class=\"hide_component\" name=\"code\" df-code></textarea>\n    <button class=\"btn-export\" name=\"open\" onclick=''open_coder(this)''> CODE </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "b913cba4-e239-40b6-84e9-7fcf3a5237f9",                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "input": "output_1",                                    "node": "fa94086f-91ee-4960-86df-43fa65f8475c"                                }                            ]                        }                    },                    "name": "",                    "outputs": {},                    "pos_x": 776,                    "pos_y": 259,                    "typenode": false                },                "c5221c80-92d0-41bc-98e0-c9d117a1b2e0": {                    "data": {                        "code": "function main(){\n  \tconsole.log(\"paso\")\n  \tif (payload[\"n1\"]  == \"jump\") {\n      next = \"output_2\"\n    }\n\tpayload[\"n1\"] = payload[\"n1\"] + \"_! modificado por js jeje\"\n  \tpayload[\"otro_campo\"] = \"agregado por js.\"\n  \n  \t\n  \n}",                        "name_box": "Validate 1",                        "nflow_auth": "false",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkred;\"></div><span class=\"material-icons\">source</span>JS Code</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"source\"></i> JS Code</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n    <textarea class=\"hide_component\" name=\"code\" df-code></textarea>\n    <button class=\"btn-export\" name=\"open\" onclick=''open_coder(this)''> CODE </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "c5221c80-92d0-41bc-98e0-c9d117a1b2e0",                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "input": "output_1",                                    "node": "307699c8-243b-4518-8aaf-478931e33709"                                }                            ]                        }                    },                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "79db20d7-a350-4bd3-87f8-0b0468313e20",                                    "output": "input_1"                                }                            ]                        },                        "output_2": {                            "connections": [                                {                                    "node": "7d42d888-3cbe-4033-8f40-3ef3f272bc58",                                    "output": "input_1",                                    "points": [                                        {                                            "pos_x": 1159,                                            "pos_y": 290.0000152587891                                        }                                    ]                                }                            ]                        }                    },                    "pos_x": 765,                    "pos_y": 139,                    "typenode": false                },                "d94dae41-33bc-437a-86a9-d5951a538f7b": {                    "class": "",                    "data": {                        "method": "ANY",                        "name_box": "",                        "nflow_auth": "false",                        "type": "starter",                        "urlpattern": "/login_flow21"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkcyan;\"></div><span class=\"material-icons\">http</span>Http Start</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"http\"></i> Http Start</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n<p>Url</p>\n<input type=\"text\"  name=\"urlpattern\" df-urlpattern>\n\n<p>Choose Method:</p>\n<select name=\"method\" df-method>\n  <option value=\"GET\">GET</option>\n  <option value=\"POST\">POST</option>\n  <option value=\"PUT\">PUT</option>\n  <option value=\"PATCH\">PATCH</option>\n  <option value=\"DELETE\">DELETE</option>\n  <option value=\"ANY\">ANY</option>\n</select>\n<button  class=\"btn-export\" onclick=''setFieldNode();save();window.open(window.location+\"/..\"+getElemenByNameOfChiled(this.parentElement,\"urlpattern\")[0].value)''> OPEN </button>\n</div>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "d94dae41-33bc-437a-86a9-d5951a538f7b",                    "inputs": {},                    "name": "",                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "66308384-c966-4dee-af14-b5a65177ad19",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 292,                    "pos_y": 414,                    "typenode": false                },                "de5bcff2-76b7-4cc9-b970-84b0dcfaeeee": {                    "class": "",                    "data": {                        "code": "function main(){\n  \tconsole.log(\"paso\")\n\tpayload[\"n1\"] = payload[\"n1\"] + \"_! modificado por js jeje\"\n  \tpayload[\"otro_campo\"] = \"agregado por js.\"\n  \n}",                        "name_box": "",                        "nflow_auth": "false",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkred;\"></div><span class=\"material-icons\">source</span>JS Code</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"source\"></i> JS Code</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n    <textarea class=\"hide_component\" name=\"code\" df-code></textarea>\n    <button class=\"btn-export\" name=\"open\" onclick=''open_coder(this)''> CODE </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "de5bcff2-76b7-4cc9-b970-84b0dcfaeeee",                    "inputs": {                        "input_1": {                            "connections": []                        }                    },                    "name": "",                    "outputs": {                        "output_1": {                            "connections": []                        }                    },                    "pos_x": -9240,                    "pos_y": -9821,                    "typenode": false                },                "f51a1499-7246-4ed0-84f9-4df8c9d61e81": {                    "class": "",                    "data": {                        "html_form": "{{index .profile \"username\"}}\n<div>\n  <input type=''input'' name=''n1''/>\n  <input type=''input'' name=''n3''/>\n  <input type=''submit''/>\n</div> ",                        "name_box": "",                        "nflow_auth": "false",                        "script": "form",                        "template": "form1.html",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div><span class=\"material-icons\">book</span>Form</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"book\"></i> Form</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr>\n\n<!-- add your form -->\n<textarea class=\"hide_component\" name=\"html_form\" df-html_form=\"<div><input type=''input'' name=''n1''/><input type=''submit''/></div>\"> </textarea>\n<button class=\"btn-export\" onclick=\"open_coder(this, `htmlmixed`)\"> CODE</button>\n<input type=\"input\" df-template=\"form1.html\" name=\"template\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "f51a1499-7246-4ed0-84f9-4df8c9d61e81",                    "inputs": {                        "input_1": {                            "connections": []                        }                    },                    "name": "",                    "outputs": {                        "output_1": {                            "connections": []                        }                    },                    "pos_x": -9230,                    "pos_y": -9818,                    "typenode": false                },                "fa94086f-91ee-4960-86df-43fa65f8475c": {                    "data": {                        "name_box": "",                        "nflow_auth": "false",                        "script": "Validate Token",                        "token": "test1",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div>Validate Token</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"\"></i> Validate Token</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr>\n\n<!-- add your form -->\nToken: <input type=\"input\" name=\"token\" df-token=\"test1\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "id": "fa94086f-91ee-4960-86df-43fa65f8475c",                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "input": "output_1",                                    "node": "090b19db-2cd3-4785-a68f-7c776b574783"                                }                            ]                        }                    },                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "b913cba4-e239-40b6-84e9-7fcf3a5237f9",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 531,                    "pos_y": 264,                    "typenode": false                },                "b0d21122-2e47-4bcc-9d7b-0fc94cc6bee2": {                    "id": "b0d21122-2e47-4bcc-9d7b-0fc94cc6bee2",                    "data": {                        "method": "ANY",                        "name_box": "/pepe sin login",                        "nflow_auth": "false",                        "type": "starter",                        "urlpattern": "/pepe11"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkcyan;\"></div><span class=\"material-icons\">video_settings</span>Form Start</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"video_settings\"></i> Form Start</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n<p>Url</p>\n<input type=\"text\"  name=\"urlpattern\" df-urlpattern>\n<button  class=\"btn-export\" onclick=''setFieldNode();save();window.open(window.location+\"/..\"+getElemenByNameOfChiled(this.parentElement,\"urlpattern\")[0].value)''> OPEN </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {},                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "21acf11d-a41a-4fb7-ae53-15d9754bddce",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 851,                    "pos_y": 417                },                "21acf11d-a41a-4fb7-ae53-15d9754bddce": {                    "id": "21acf11d-a41a-4fb7-ae53-15d9754bddce",                    "data": {                        "html_form": "<div><input type=''input'' name=''n1''/><input type=''submit''/></div>",                        "name_box": "",                        "script": "form",                        "template": "form1.html",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div><span class=\"material-icons\">book</span>Form</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"book\"></i> Form</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<hr>\n\n<!-- add your form -->\n<textarea class=\"hide_component\" name=\"html_form\" df-html_form=\"<div><input type=''input'' name=''n1''/><input type=''submit''/></div>\"> </textarea>\n<button class=\"btn-export\" onclick=\"open_coder(this, `htmlmixed`)\"> CODE</button>\n<input type=\"input\" df-template=\"form1.html\" name=\"template\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "node": "b0d21122-2e47-4bcc-9d7b-0fc94cc6bee2",                                    "input": "output_1"                                }                            ]                        }                    },                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "f9615b9a-5a82-4ea6-8b73-708ab877dcf6",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 1088,                    "pos_y": 416                },                "f9615b9a-5a82-4ea6-8b73-708ab877dcf6": {                    "id": "f9615b9a-5a82-4ea6-8b73-708ab877dcf6",                    "data": {                        "field": "",                        "http_code": "200",                        "script": "json",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkviolet;\"></div><span class=\"material-icons\">source</span>JSON Render</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"source\"></i> JSON Render</div>\n\t\t<div class=\"box\">\n\t\t  <input type=\"text\" name=\"field\" df-field=\"\">\n<input type=\"number\" name=\"http_code\" df-http_code=\"200\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "node": "21acf11d-a41a-4fb7-ae53-15d9754bddce",                                    "input": "output_1"                                }                            ]                        }                    },                    "outputs": {},                    "pos_x": 1359,                    "pos_y": 417                },                "09dfa380-1df9-4cc0-a5b1-a8eaef0ac9fc": {                    "id": "09dfa380-1df9-4cc0-a5b1-a8eaef0ac9fc",                    "data": {                        "name_box": "",                        "script": "Validate Login",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div>Validate Login</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"\"></i> Validate Login</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<hr>\n\n<!-- add your form -->\n    \n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "node": "3a4b785c-1b89-410f-82a5-e5127922105d",                                    "input": "output_1"                                }                            ]                        }                    },                    "outputs": {},                    "pos_x": 875,                    "pos_y": 596                },                "3a4b785c-1b89-410f-82a5-e5127922105d": {                    "id": "3a4b785c-1b89-410f-82a5-e5127922105d",                    "data": {                        "html_form": "<div>\n  User:<input type=''input'' name=''username''/><br/>\n  Passwrod:<input type=''password'' name=''password''/><br/>\n  <input type=''submit'' value=\"Login\"/>\n</div>",                        "name_box": "Form",                        "script": "form",                        "template": "login_template.html",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div><span class=\"material-icons\">book</span>Form</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"book\"></i> Form</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<hr>\n\n<!-- add your form -->\n<textarea class=\"hide_component\" name=\"html_form\" df-html_form=\"<div><input type=''input'' name=''n1''/><input type=''submit''/></div>\"> </textarea>\n<button class=\"btn-export\" onclick=\"open_coder(this, `htmlmixed`)\"> CODE</button>\n<input type=\"input\" df-template=\"form1.html\" name=\"template\">\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "node": "36504106-79a2-44bb-aa5d-51d95281b902",                                    "input": "output_1"                                }                            ]                        }                    },                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "09dfa380-1df9-4cc0-a5b1-a8eaef0ac9fc",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 591,                    "pos_y": 593                },                "36504106-79a2-44bb-aa5d-51d95281b902": {                    "id": "36504106-79a2-44bb-aa5d-51d95281b902",                    "data": {                        "method": "ANY",                        "name_box": "Login",                        "nflow_auth": "false",                        "type": "starter",                        "urlpattern": "/nflow_login"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkcyan;\"></div><span class=\"material-icons\">video_settings</span>Form Start</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"video_settings\"></i> Form Start</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n<p>Url</p>\n<input type=\"text\"  name=\"urlpattern\" df-urlpattern>\n<button  class=\"btn-export\" onclick=''setFieldNode();save();window.open(window.location+\"/..\"+getElemenByNameOfChiled(this.parentElement,\"urlpattern\")[0].value)''> OPEN </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {},                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "3a4b785c-1b89-410f-82a5-e5127922105d",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 296,                    "pos_y": 573                },                "c0251ac8-564b-4de4-8327-c42a1a5af3ba": {                    "id": "c0251ac8-564b-4de4-8327-c42a1a5af3ba",                    "data": {                        "name_box": "",                        "nflow_auth": "false",                        "script": "Logout",                        "type": "js"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkorange;\"></div>logout</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"\"></i> logout</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<hr>\n\n<!-- add your form -->\n    \n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {                        "input_1": {                            "connections": [                                {                                    "node": "59bb6ea9-e751-4073-a1f9-76d80b65129d",                                    "input": "output_1"                                }                            ]                        }                    },                    "outputs": {},                    "pos_x": 554,                    "pos_y": 760                },                "59bb6ea9-e751-4073-a1f9-76d80b65129d": {                    "id": "59bb6ea9-e751-4073-a1f9-76d80b65129d",                    "data": {                        "method": "ANY",                        "name_box": "Logout",                        "nflow_auth": "false",                        "type": "starter",                        "urlpattern": "/logout"                    },                    "html": "<div>\n  \n\n        <div class=\"title-box\">      <div class=\"mark_color_box\" style=\"background:darkcyan;\"></div><span class=\"material-icons\">video_settings</span>Form Start</div>\n            <div>\n                <input type=\"text\" value=\"\" name=\"input_box\" df-name_box=\"\" class=\"df_name_box\"\n                readOnly\n                />\n            </div>\n            </div>\n            <div class=\"form_prop_of_box\" style=\"display:none\">\n                <div>\n\t\t<div class=\"title-box\"><i class=\"video_settings\"></i> Form Start</div>\n\t\t<div class=\"box\">\n\t\t  Name:<input type=\"text\" name=\"name_box\" df-name_box=\"\" class=\"nflow_field\">\n<div class=\"nflow_auth\">Auth:<input type=\"checkbox\" df-nflow_auth=\"false\" name=\"nflow_auth\"></div>\n<hr/>\n<p>Url</p>\n<input type=\"text\"  name=\"urlpattern\" df-urlpattern>\n<button  class=\"btn-export\" onclick=''setFieldNode();save();window.open(window.location+\"/..\"+getElemenByNameOfChiled(this.parentElement,\"urlpattern\")[0].value)''> OPEN </button>\n\t\t</div>\n\t  </div>\n            </div>\n        ",                    "typenode": false,                    "inputs": {},                    "outputs": {                        "output_1": {                            "connections": [                                {                                    "node": "c0251ac8-564b-4de4-8327-c42a1a5af3ba",                                    "output": "input_1"                                }                            ]                        }                    },                    "pos_x": 287,                    "pos_y": 749                }            }        }    }}', 'app', 'function log () {
  profile = get_profile()
  console.log(
   ">>",
    "User:",profile["username"],
    "BoxID:", box_id,
    "BoxName:", box_name, 
    "BoxType:", box_type, 
    "Time:",duration_ms, "ms"
  )
}

function auth(){
    if (!exist_profile() ) {
        next = "login"
        return
    }
}');