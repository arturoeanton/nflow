function main(){
  if (parseFloat(payload[nflow_data["field"]]) > parseFloat(nflow_data["limit"])){
    next= "output_1"
    return
  } 
  next="output_2"
}