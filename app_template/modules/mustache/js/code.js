function main(){
    console.log(JSON.stringify(payload))
    c.HTML(200, mustache(dromedary_data["template"], payload))
}