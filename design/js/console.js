var loc = window.location;
var uri = 'ws:';

if (loc.protocol === 'https:') {
    uri = 'wss:';
}
var pos_init = (window.loc.protocol +"//"+  window.loc.host +"/").length
var pos_end = window.loc.href.indexOf("/nflow")
var sub = window.loc.href.substr(pos_init,pos_end -pos_init) 
uri += '//' + loc.host;
if (sub != ""){
  uri += '/'+ sub
}
uri += '/nflow/console/ws';

ws_console = new WebSocket(uri)

ws_console.onopen = function () {
    console.log('Connected')
}

ws_console.onmessage = function (evt) {
    //console.text += evt.data +"\n"
    if (evt.data == "") {
        return
    }
   
    if (evt.data.startsWith("__node_id:run:")){
      setNodeRunnig(evt.data.split(":")[2])
      return
    }
    if (evt.data.startsWith("__node_id:stop:")){
      setNodeStop(evt.data.split(":")[2])
      return
    }
    print_console (evt.data)
}

function setNodeRunnig (node_id){
  var node = document.getElementById("node-"+node_id)
  node.style.borderColor= "red"
  node.style.borderWidth= "3px"
}

function setNodeStop (node_id){
  var node = document.getElementById("node-"+node_id)
  node.style.borderColor= ""
  node.style.borderWidth= ""
}



function print_console(text) {
    var dt = new Date();
    log_time =`${
        dt.getDate().toString().padStart(2, '0')}/${
        (dt.getMonth()+1).toString().padStart(2, '0')}/${
        dt.getFullYear().toString().padStart(4, '0')} ${
        dt.getHours().toString().padStart(2, '0')}:${
        dt.getMinutes().toString().padStart(2, '0')}:${
        dt.getSeconds().toString().padStart(2, '0')}`


    textarea_console.value += log_time + " - " + text +"\n"
    textarea_console.scrollTop = textarea_console.scrollHeight;
}

input_cmd.addEventListener("keyup", function(event) {
    // Number 13 is the "Enter" key on the keyboard
    if (event.keyCode === 13) {
      // Cancel the default action, if needed
      event.preventDefault();
      // Trigger the button element with a click
      if (this.value == "clear"){
        textarea_console.value=">>" 
      }else{
        ws_console.send(this.value)
        textarea_console.value += ">>" + this.value +"\n"
        textarea_console.scrollTop = textarea_console.scrollHeight;
      }
      this.value=""

    }
  });