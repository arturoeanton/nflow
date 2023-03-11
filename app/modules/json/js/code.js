function main(){
    var code = 200
    if (nflow_data["http_code"] != ""){
        code = parseInt(nflow_data["http_code"])
    }
    if (nflow_data["field"] == undefined || nflow_data["field"] =="" || nflow_data["field"] =="*"){
        c.JSON(code, payload)
        return
    }
    c.JSON(code,find_element(nflow_data["field"],payload))
}