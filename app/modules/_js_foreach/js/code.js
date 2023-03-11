function main(){
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
}