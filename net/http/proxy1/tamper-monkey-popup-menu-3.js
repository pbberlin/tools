// ==UserScript==
// @name         popup-menu
// @description  include jQuery and make sure window.$ is the content page's jQuery version, and this.$ is our jQuery version. 
// @description  http://stackoverflow.com/questions/28264871/require-jquery-to-a-safe-variable-in-tampermonkey-script-and-console
// @namespace    http://your.homepage/
// @version      0.122
// @author       iche
// @downloadURL  http://localhost:8085/mnt01/tamper-monkey-popup-menu.js
// @updateURL    http://localhost:8085/mnt01/tamper-monkey-popup-menu.js //serving the head with possibly new version
// // https://developer.chrome.com/extensions/match_patterns
// @match        *://www.welt.de/*
// @match        *://www.handelsblatt.com/*  
// @match        *://www.focus.de/*
// // @include     /^https?:\/\/www.flickr.com\/.*/
// // @require      http://cdn.jsdelivr.net/jquery/2.1.3/jquery.min.js
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
        console.log("about to add hover popups; " + $.fn.jquery + " Version");


var vPaddingOffs = 25 ;
var vShadowOffs  = 6 ;
var hOffs  = 12 ;  // horizontal offset

function PopupContent(obj){

          
    var parAnchor = obj.closest("a");  // includes self

    
    if (parAnchor.length<=0) {
        return "";
    }

    var domainY = location.host;
    var protocolY = location.protocol;

    var href = parAnchor.attr('href');
    if (href.indexOf('/') === 0) {
        href = domainY + href;
    }

    var text = obj.text();
    if (text.length > 100){
        var textSh = "";
        textSh  = text.substr(0,50); // start, end
        textSh += " ... ";
        textSh += text.substr(text.length-50,text.length-1);
        text = textSh;
    }
    
    var formHtml = "";
    var prox01 = "http://localhost:8085/weedout";
    var prox02 = "https://libertarian-islands.appspot.com/weedout";

    formHtml += "<form action='"+prox01+"' method='post' ";
    formHtml += "    target='proxy-window' >";
    formHtml += "<input type='hidden'  name='prot'  value='"+protocolY+"' >";
    formHtml += "<input type='hidden'  name='url-x' value='"+href+"' >";
    formHtml += "<input type='submit'               value='subm'     >";
    formHtml += "</form>";

    var html = "";
    html += "<a target='proxy-window'  href='"+prox01 + "?url-x="+ href+"&prot="+protocolY+"' >" + text + "</a>";
    html += formHtml;

    return html;

}


function AddCSS(){

    if ($('#css-hover-popup').length <= 0) {

        var s =  '';
        s += '<style type="text/css"  id="css-hover-popup" >';
        s += '';
        s += '.hc-b2{';
        s += '    position:absolute;';
        s += '    z-index: 10;';
        s += '    left:-10px;   top:-10px;';
        s += '    width:  240px;  ';
        s += '    /* dont fix the height - use jQuery.outerHeight() for computations */ ';
        s += '    height: auto ';
        s += '';
        s += '    /* alternating upon top-bottom, adapt vPaddingOffs accordingly */ ';
        // s += '    border-left: solid 2px #d22; ';
        // s += '    border-right:solid 2px #d22; ';
        s += '    margin:         0px !important; ' ;
        s += '    padding:        0px; ';
        s += '    padding-left:   0px !important; ';
        s += '    padding-right:  0px !important;' ;
        s += '';
        s += '    display:none;';
        s += '}';
        s += '';
        s += '.hc-b3{';
        s += '    text-align:left; ';
        s += '    font-family:Sans-serif !important; font-size:12px !important; ';
        s += '    color:#666 !important; line-height:1.5em;';
        s += '    background:url("arrow-desc.png") no-repeat scroll center 0 transparent;';
        s += '    background-image: none;';
        s += '    background-color: #ddd !important;  ';
        s += '    border:solid 1px #ddd; ';
        s += '    border-radius:4px;';
        s += '    box-shadow:3px 3px 3px #888;';
        s += '    /* skipping the old stuff';
        s += '      -moz-border-radius:3px;-webkit-border-radius:3px;';
        s += '      -moz-box-shadow:5px 5px 5px #888;-webkit-box-shadow:5px 5px 5px #888;  */';
        s += '    margin:         0px !important; ' ;
        s += '    padding:        5px 8px !important;  ';
        s += '}';
        s += '';


        s += '</style>';

        $(s).appendTo('head');


        var popupUpScaffold = '<div id="popup1" class="hc-b2" ><div id="pop1inner" class="hc-b3" >content popup1</div></div>'; 
        $('body').after(popupUpScaffold);  //Append after selected element

    }

}





function HidePop(obj){
    var popup = $('#popup1');
    // popup.hide();
    popup.stop(true, true).fadeOut(300, function () {
        console.log("fadeout-callback");
    });
}

// 
// hovercard.js made the anchor overlay the popup with
//          wrapper.css("zIndex", "200");
//          anchor.css("zIndex", "100");
//          popup.css("zIndex", "50");
// And resetting the all to zero on mouseleave.
// This would only look good with align right-or-left.
function ShowPop(obj){
    

    var popup = $('#popup1');
    var inner = $('#pop1inner');

    var html = PopupContent(obj);
    inner.html(html); // setting content => force sizing


    var objAbsTL = obj.offset();

    var scrollTop  = $('body').scrollTop();   // How many pixel is document scrolled in the viewport
    var scrollLeft = $('body').scrollLeft();   

    var viewportTop  = objAbsTL.top  - scrollTop;
    var viewportLeft = objAbsTL.left + obj.outerWidth() - scrollLeft;

    var vert = "lower";
    if (viewportTop < $(window).height() / 2) {
        vert = "upper";
    }
    var horz = "right";
    if (viewportLeft < $(window).width() * 0.75) {
        horz = "left";
    }

    console.log("quadrant", vert, horz);


    if ( vert == "upper" ) {
        // show below
        popup.css({
            'padding-top':    vPaddingOffs + 'px',
            'padding-bottom': '0px',
            'background-position': 'center top',
        });
        popup.css({
            left: objAbsTL.left - hOffs + 'px',
            top:  -20 + objAbsTL.top   +  obj.outerHeight()  + 'px',
        });

    } else {
        // show above
        popup.css({
            'padding-top':    '0px',
            'padding-bottom':  vPaddingOffs  +  'px',
            'background-position': 'center bottom',
        });
        popup.css({
            left: objAbsTL.left  - hOffs + 'px',
            top:  objAbsTL.top   - popup.outerHeight() + vPaddingOffs - + vShadowOffs  + 'px',
        });


        // popup.css({
        //     top:     'auto',
        //     bottom:  $( document ).height() -  objAbsTL.top  - vPaddingOffs + 'px',
        // });

    }

    inner.css({
        'text-align':  'left',
    });
    if ( horz == "right" ) {
        popup.css({
            left:  objAbsTL.left + obj.outerWidth() - popup.outerWidth() + hOffs +  'px',
            right: 'auto',
        });
        inner.css({
            'text-align':  'right',
        });
        console.log("righted",objAbsTL.left + obj.outerWidth());
    }


    if ( obj.outerWidth() <= 240 ) {
        popup.css({
            width: 240 + 'px',
        });
    } else if ( obj.outerWidth() > 480 ) {
        popup.css({
            width: 480 + 'px',
        });
    } else {
        popup.css({
            width: obj.outerWidth() + 2*hOffs + 'px',
        });
    }

    inner.show();
    popup.stop(true, true).delay(1).fadeIn(500);


}


$( document ).ready(function() {

    var wh = $( window ).height();            // Returns height of browser viewport
    var dh = $( document ).height();          // Returns height of HTML document
    var scrollVert = $('body').scrollTop();   // How many pixel is document scrolled in the viewport
    var obj = $('body').find('*').last();
    var os1 = obj.offset();                   // Our object's position relative to document root
    var os2 = obj[0].getBoundingClientRect(); // Our object's position relative to document root - MINUS scroll pos
    console.log("height window - doc:",wh,dh," sctp:", scrollVert , " offs1",os1.top, " offs2", os2.top);
    // console.log("offs1-l", os1.left, " offs2-l",  os2.left);



    AddCSS();


    $( 'a' ).on( "mouseenter", function(evt) {
        var obj = $(evt.target);
        ShowPop(obj);
    });

    $( '#popup1' ).on( "mouseleave", function(evt) {
        var obj = $(evt.target);
        HidePop(obj);
    });



    $( 'a' ).on( "focusin", function(evt) {
        var obj = $(evt.target);
        ShowPop(obj);
    });

    $( '#popup1' ).on( "focusout", function(evt) {
        var obj = $(evt.target);
        HidePop(obj);
    });

    console.log( "document ready completed" );
});

        

        console.log("hover popups handler added");
        //isolated jQuery end;
    });
})(window.jQuery.noConflict(true));

        
