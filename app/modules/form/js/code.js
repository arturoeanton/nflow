function main(){
  var action = "/nflow/node/run/"+__flow_name+"/"+__outputs["output_1"]
  var html_form = dromedary_data["html_form"]
  console.log(JSON.stringify(dromedary_data))
  c.HTML(200,`<form method='post' id="form_main" action='${action}'> 
  	${html_form}
  </form>`)
  if (payload == undefined || payload == null){
    payload = {}
  }
  payload["break"] = true
}