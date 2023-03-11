function main(){
  	var storage_id =  nflow_data["storage_id"] +"-"
    if (nflow_data["storage_id"] == "") {
      storage_id = ""
    }
  	for (var key in form){
  		set_session("nflow_form",storage_id+key, form[key])
    }
}