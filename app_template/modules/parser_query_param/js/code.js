function main (){
    if (payload == undefined){
      payload = {}
    } 
    var params = url_values_to_map(c.QueryParams())
    for (key in params){
      payload[key] = c.QueryParam(key)
    }
  }