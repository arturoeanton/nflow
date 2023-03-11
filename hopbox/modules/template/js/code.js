function main(){
    console.log(JSON.stringify(payload))
    c.HTML(200, template(dromedary_data["template"], payload))
}