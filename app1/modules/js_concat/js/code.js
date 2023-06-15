function main(){
	var elem = find_element(nflow_data["elem"],payload)
    var field1 = ""
    if (nflow_data["field1"] != ""){
      field1 = find_element(nflow_data["field1"],elem)
    }
    var field2 = ""
    if (nflow_data["field2"] != ""){
      field2 = find_element(nflow_data["field2"],elem)
    }

  	payload[nflow_data["result"]] = field1 +nflow_data["connector"] + field2
}