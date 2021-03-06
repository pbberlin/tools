// ==UserScript==
// @name         rightclick-menu
// @description  include jQuery and make sure window.$ is the content page's jQuery version, and this.$ is our jQuery version. http://stackoverflow.com/questions/28264871/require-jquery-to-a-safe-variable-in-tampermonkey-script-and-console
// @namespace    http://your.homepage/
// @version      0.131
// @author       iche
// @downloadURL  http://localhost:8085/mnt01/tamper-monkey-rightclick-menu.js
// @updateURL    http://localhost:8085/mnt01/tamper-monkey-rightclick-menu.js //serving the head with possibly new version
// // https://developer.chrome.com/extensions/match_patterns
// @match        *://*.welt.de/*
// @match        *://*.handelsblatt.com/*
// @match        *://*.focus.de/*
// @exclude        *://*/*.png
// @exclude        *://*/*.jpg
// // @include     /^https?:\/\/www.flickr.com\/.*/

//  // @require      http://cdn.jsdelivr.net/jquery/2.1.3/jquery.min.js
// @require      https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js
// @grant        none
// @noframes
// @run-at      document-end
// ==/UserScript==


// fallback http://encosia.com/3-reasons-why-you-should-let-google-host-jquery-for-you/
if (typeof jQuery === 'undefined') {
    console.log("CDN blocked by Iran or China?");
    document.write(unescape('%3Cscript%20src%3D%22/path/to/your/scripts/jquery-2.1.4.min.js%22%3E%3C/script%3E'));
}

(function ($, undefined) {
    $(function () {
        //isolated jQuery start;
        console.log("about to add right click handler; " + $.fn.jquery + " Version");

        
        var menuHtml = "<div id='menu01' style='background-color:#ccc;padding:4px;z-index:1000;display:none'>";
        //menuHtml    += "  <li>item1</li>" ;
        menuHtml    += "  <li  id='menu01-item02'  style='height:64px;'>item2</li>" ;
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
            var obj = $(evt.target);
            var parAnchor = obj.closest("a");  // includes self

            
            if (parAnchor.length>0) {
                var domainX = document.domain;
                var href = obj.attr('href');
                if (href.indexOf('/') == 0) {
                    href = domainX + href
                }
                var text = obj.text()
                if (text.length > 100){
                    var textSh = "";
                    textSh  = text.substr(0,50); // start, end
                    textSh += " ... ";
                    textSh += text.substr(text.length-50,text.length-1);
                    text = textSh;
                }

                var formHtml = "";
                var prox01 = "http://localhost:8085/dedup";
                var prox02 = "https://libertarian-islands.appspot.com/dedup";

                formHtml += "<form action='"+prox01+"' method='post' ";
                formHtml += "    target='proxy-window' >";
                formHtml += "<input type='hidden'  name='url-x' value='"+href+"' >";
                formHtml += "<input type='submit'               value='subm'     >";
                formHtml += "</form>";

                var innerH = "";
                innerH += "<a target='proxy-window'  href='"+prox01 + "?url-x="+ href+"' >" + text + "</a>";
                innerH += formHtml;
                $('#menu01-item02').html(innerH);   
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