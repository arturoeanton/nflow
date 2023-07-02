function main(){
  	var url_back = ""+get_session("auth-session","redirect_url")
  	var flag = false
  	if (payload["username"] == "user1") {
      if (payload["password"] == "1") {
      	flag = true
      }
    }
  
  	if (payload["username"] == "admin") {
      if (payload["password"] == "admin") {
      	flag = true
      }
    }
  	
  	if (flag) {
      	set_profile({"username":payload["username"]})
  		return c.HTML(200," <script>window.location.href = '"+url_back+"'</script>")
    }
    return c.HTML(200," Error Login <br/>  <a href='/home' >Home</a>")

}