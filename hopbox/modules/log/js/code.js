function main(){

    if (dromedary_data["field"] == undefined || dromedary_data["field"] ==""){
        ws_console_log(JSON.stringify(payload))
        return
    }

    ws_console_log(payload[dromedary_data["field"]])

}