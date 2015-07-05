// ==UserScript==
// @name         welt-context-menu
// @namespace    http://your.homepage/
// @version      0.1
// @description  enter something useful
// @author       You
// @match        http://www.welt.de/*
// @match        https://www.welt.de/*
// @match        http://www.handelsblatt.com/*
// @match        https://www.handelsblatt.com/*
// @match        http://www.focus.de/*
// @grant        none
// @require http://code.jquery.com/jquery-latest.js
// ==/UserScript==


// inspired by http://stackoverflow.com/questions/4909167/how-to-add-a-custom-right-click-menu-to-a-webpage
// // @run-at context-menu
//
console.log("add01");

$("body").append( "<div id='menu01' style='background-color:#ccc;padding:4px;'>  <li>item1</li><li>item2</li>   </div>" );

function logX(arg01){
    console.log( "menu called " + arg01);
}

function func2(kind, evt){


    logX( kind + " 1 - x" + evt.pageX + " - y" + evt.pageY );

    var isAnchor = $(evt.target).is("A")     // .get(0).tagName
    
    if (true) {
        $('#menu01').append("<li>"+ evt.target + "</li>");   
        var fromBottom =  evt.pageY - $('#menu01').height() - 8;
        $('#menu01').css({
            "position": "absolute", 
            "top":       fromBottom + "px",
            "left":      evt.pageX + "px",
        });            
        $('#menu01').show();
    }
    

    logX( kind + " 2");    
}    



$('#menu01').on('click', function(e) {
    $('#menu01').hide();
});


$(document).on('contextmenu', function(e){
    e.stopPropagation(); // before
    //e.preventDefault();
    func2("allbrows",e);
});


console.log("add02");


