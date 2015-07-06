// ==UserScript==
// @name         welt-context-menu
// @description  include jQuery and make sure window.$ is the content page's jQuery version, and this.$ is our jQuery version. http://stackoverflow.com/questions/28264871/require-jquery-to-a-safe-variable-in-tampermonkey-script-and-console
// @namespace    http://your.homepage/
// @version      0.1
// @author       iche
// @match        http://www.welt.de/*
// @match        https://www.welt.de/*
// @match        http://www.handelsblatt.com/*
// @match        https://www.handelsblatt.com/*
// @match        http://www.focus.de/*
// @require      http://cdn.jsdelivr.net/jquery/2.1.3/jquery.min.js
// @grant        none
// ==/UserScript==


(function ($, undefined) {
    $(function () {
        //isolated jQuery start;
        console.log("about to add right click handler");

        var menuHtml = "<div id='menu01' style='background-color:#ccc;padding:4px;z-index:1000'>";
        menuHtml    += "  <li>item1</li>" ;
        menuHtml    += "  <li  id='menu01-item02'>item2</li>" ;
        menuHtml    += "</div>" ;

        $("body").append(menuHtml);

        $('#menu01').on('click', function(e) {
            $('#menu01').hide();
        });

        $(document).on('click', function(e) {
            $('#menu01').hide();
        });


        function logX(arg01){
            console.log( "menu called " + arg01);
        }

        function showContextMenu(kind, evt){


            logX( kind + " 1 - x" + evt.pageX + " - y" + evt.pageY );


            var $lp = $(evt.target);
            var isAnchor = false

            for (i = 0; i < 10; i++) { 
                i++;
                isAnchor = $lp.is("A")     // .get(0).tagName
                if( isAnchor){
                    break;
                }
                $lp = $lp.parent(); // $( "html" ).parent() returns 
            }
            
            if (isAnchor) {
                var domainX = document.domain;
                var href = $lp.attr('href');
                if (href.indexOf('/') == 0) {
                    href = domainX + href
                }
                var text = $lp.text()
                if (text.length > 100){
                    var textSh = "";
                    textSh  = text.substr(0,50); // start, end
                    textSh += " ... ";
                    textSh += text.substr(text.length-50,text.length-1);
                    text = textSh;
                }

                $('#menu01-item02').html("<a href="+href+" >" + href +  " <br/>" + text + "</a>");   
            }

            var fromBottom =  evt.pageY - $('#menu01').height() - 8;
            $('#menu01').css({
                "position": "absolute", 
                "top":       fromBottom + "px",
                "left":      evt.pageX + "px",
            });            
            $('#menu01').show();
            

            logX( kind + " 2");    
        }    





        $(document).on('contextmenu', function(e){
            e.stopPropagation(); // before
            //e.preventDefault();  // prevent browser context menu from appear
            showContextMenu("allbrows",e);
        });


        console.log("right click handler added");



    //isolated jQuery end;
    });
})(window.jQuery.noConflict(true));

        
