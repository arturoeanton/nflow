function show_console(){


        if (bottom_panel.style.animationName != "show_bottom_panle") {
            bottom_panel.style.animationName = "show_bottom_panle"
            bottom_panel.style.bottom = "0px"
            return
        }

        if (bottom_panel.style.animationName != "hide_bottom_panle") {
            bottom_panel.style.animationName = "hide_bottom_panle"
            bottom_panel.style.bottom = "-270px"
            return
        }
    

}

function logout(){
    var scim_base = env[scim_base]
    window.location = scim_base+"/home";
}