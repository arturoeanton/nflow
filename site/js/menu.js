function menu_action() {
    var menu = document.querySelector(".menu")
    var tab_btn = document.querySelector("body > div.menu > div.menu_tag > button")
    var tab = document.querySelector("body > div.menu > div.menu_tag ")
    if (menu.style.animationName != "hide_menu") {
        menu.style.animationName = "hide_menu"
        menu.style.left = `-${menu.offsetWidth}px`
        menu.boxShadow = "0px 0px 0px 0px";
        tab_btn.innerHTML = "☰"
        tab_btn.style.fontSize ="16px"

        tab.style.borderRadius= "0 6px 6px 0px";
        tab.style.boxShadow = "0px 1px 1px 0px";
        tab.style.right = "-30px"
        
        return
    }

    menu.style.animationName = "show_menu"
    menu.style.left = "0px"
    menu.boxShadow = "0px 2px 2px 2px";

    tab_btn.innerHTML = "&times;"
    tab_btn.style.fontSize ="20px"
        
    tab.style.boxShadow = "none";
    tab.style.right = "0px"

}