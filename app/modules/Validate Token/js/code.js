function main(){
  next= "output_1"
  if (header["Authorization"] != "Bearer " + nflow_data["token"]){
    c.JSON(401,{"error":"Unauthorized"})
    next= "break"
  }
  
}