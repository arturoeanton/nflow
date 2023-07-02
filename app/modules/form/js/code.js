function main(){
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
}