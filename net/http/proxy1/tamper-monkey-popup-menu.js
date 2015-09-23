// ==UserScript==
// @name         popup-menu
// @description  include jQuery and make sure window.$ is the content page's jQuery version, and this.$ is our jQuery version. 
// @description  http://stackoverflow.com/questions/28264871/require-jquery-to-a-safe-variable-in-tampermonkey-script-and-console
// @namespace    http://your.homepage/
// @version      0.12
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


var vPadOffs = 30; // padding offset; top & bottom are 10; one is 20
    vPadOffs += 2; // border with
var hPadSimple = 10; // simple horizontal padding offset

var winBrdrLeft = 4;                     // border to outer browser window
var winBrdrRight = winBrdrLeft+2+2;        // popup border 2*1, shadow

var vDist = 0;  // vertical distance to popup - cannot be set - any gap causes mouseleave event


function PopupContent(obj){

          
    var parAnchor = obj.closest("a");  // includes self

    
    if (parAnchor.length<=0) {
    	return "";
    }

    var domainX = document.domain;
    console.log("domainX", domainX);

    var href = obj.attr('href');
    if (href.indexOf('/') === 0) {
        href = domainX + href;
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
    formHtml += "<input type='hidden'  name='url-x' value='"+href+"' >";
    formHtml += "<input type='submit'               value='subm'     >";
    formHtml += "</form>";

    var html = "";
    html += "<a target='proxy-window'  href='"+prox01 + "?url-x="+ href+"' >" + text + "</a>";
    html += formHtml;

    return html;

}


function AddCSS(){

	if ($('#css-hover-popup').length <= 0) {

		var s =  '';
	    s += '<style type="text/css"  id="css-hover-popup" >';
	    s += '.hc-a{';
	    s += '    position: relative; ';
	    s += '    display:inline; ';
	    s += '}';
	    s += '.hc-b1{}';
	    s += '.hc-b2{';
	    s += '    position:absolute;';
	    s += '    z-index: 10;';
	    s += '    left:-10px;   top:-10px;';
	    s += '    /* fixed height required for positioning computations */';
	    s += '    width:150px;  height:36px;';
	    s += '    text-align:left; ';
	    s += '    font-family:Sans-serif !important; font-size:12px !important; ';
	    s += '    color:#666 !important; line-height:1.5em;';
	    s += '';
	    s += '    background:url("arrow-desc.png") no-repeat scroll center 0 transparent;';
	    s += '    background-image: none;';
	    s += '    background-color: #ddd !important;  ';
	    s += '    border:solid 1px #ddd; ';
	    s += '    border-radius:4px;';
	    s += '    box-shadow:3px 3px 3px #888;';
	    s += '    /* skipping the old stuff';
	    s += '      -moz-border-radius:3px;-webkit-border-radius:3px;';
	    s += '      -moz-box-shadow:5px 5px 5px #888;-webkit-box-shadow:5px 5px 5px #888;  */';
	    s += '    margin: 0px !important;' ;
	    s += '    padding:  10px; ';
	    s += '    padding-left:  10px !important; ';
	    s += '    padding-right: 10px !important;' ;
	    s += '    /* alternating upon top-bottom, adapt vPadOffs accordingly */ ';
	    s += '    aapadding-top: 20px;';
	    s += '    display:none;';
	    s += '}';
	    s += '</style>';

	    $(s).appendTo('head');

	}

}

function CreateAnchorWrappers(){

    var start = new Date().getTime();

    // Create parent wrappers
    // $( '.hovering-popup' ).each(function( index ) {
    $( 'a' ).each(function( index ) {
        // console.log( index + ": " + $( this ).text() );
        var obj = $(this);
        if (!obj.parent().hasClass("hc-a")){
            obj.wrap('<div class="hc-a" />');
            obj.addClass("hc-b1");
            var specialPop = '<div class="hc-b2" >' + 'specialpop' + '</div>'; // Potentially a lot of DOM insertions over time
            obj.after(specialPop);  //Append after selected element
            
            var w = obj.width();
            if (w < 320) {
                w = 320;
            }
            obj.siblings(".hc-b2").eq(0).css({ 'width': w +'px' });
        }
    });

    // Remember when we finished
    var end = new Date().getTime();
    console.log("creating wrappers took ",end - start);
}


function HidePop(obj){
    // var popup = $('.hc-b2');
    // popup.hide();
    obj.stop(true, true).fadeOut(300, function () {
        console.log("fadeout-callback");
    });
}

// A wrapping of a-tags is required,
// to bind the a-tag and the popup 
// into one single mouseover-mouseleave-context.
// 
// The wrapping makes positioning of the popup relative.
// At left:0, top:0; the popup overlays the original anchor-tag.
// 
// The wrapping must be done in advance, causing a lot of DOM overhead.
// Approx. 3 ms per link.
// Dynamic wrapping could be done like this:
// Reacting to a-tag-mouseover, 
// creating wrapper,
// calling parent.triggerHandler("mouseover"),
// calling it only once.
// 
// Horizonal positioning is centered.
// We could also make it align right-or-left   - with left: auto; right:-10px;
// 
// hovercard.js made the anchor overlay the popup with
//          wrapper.css("zIndex", "200");
//          anchor.css("zIndex", "100");
//          popup.css("zIndex", "50");
// And resetting the all to zero on mouseleave.
// This would only look good with align right-or-left.
function ShowPop(obj){
    
    // CreateAnchorWrappers()    // must be done in advance

    // var popup = $('.hc-b2');
    var popup = obj.siblings(".hc-b2").eq(0);
    if (popup.length<=0){
        popup = obj.find(".hc-b2").eq(0);
    }

    var objAbsTL = obj.offset();
    var objRelTL = obj[0].getBoundingClientRect();

    var new_left =  - (popup.width() / 2) + (obj.width() / 2) - hPadSimple;
    
    // prevent transgression to the left
    if (objRelTL.left + new_left < winBrdrLeft){
        new_left = -objRelTL.left + winBrdrLeft;
    }
    // prevent transgression to the right
    if (objRelTL.left + new_left + popup.width() + 2*hPadSimple > $(window).width() - winBrdrRight){
        new_left = -popup.width() - 2*hPadSimple - ( objRelTL.left - $(window).width() )  - winBrdrRight;
    }



    if (objRelTL.top < $(window).height() / 2) {
        popup.css({
            'padding-top':    '20px',
            'padding-bottom': '10px',
            left: new_left + 'px',
            top: obj.height() +  1*vDist + 'px',
            'background-position': 'center top',
            //'background-image': "url(arrow-desc.png)",
        });
    } else {
        console.log("ort>wh", objRelTL.top, $(window).height() / 2);
        popup.css({
            'padding-top':    '10px',
            'padding-bottom': '20px',
            left: new_left + 'px',
            top: - popup.height() - vPadOffs - 1*vDist + 'px',
            'background-position': 'center bottom',
            //'background-image':  "url(arrow-asc.png)",
        });
    }

    var anch = obj.siblings(".hc-b1").eq(0);
    if (anch.length<=0){
        anch = obj.find(".hc-b1").eq(0);
    }
    var html = PopupContent(anch);
    popup.html(html);

    // popup.show();
    popup.stop(true, true).delay(200).fadeIn(500);

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

    CreateAnchorWrappers();


    // Here popup would not be clickable.
    // It disappears as soon as cursor leaves obj.
    $( 'a' ).on( "mouseover", function() {
    });


    $(document).on( "click", "a" , function(e) {
    	// console.log("click1");
        // e.preventDefault();
    });

    // We have to latch the events onto the parent wrapper node.
    // Thus all parent nodes need to be created on documentload :(

    // mouseover fires again for child elements.
    // it prevents click events 
    // $( '.hc-a' ).on( "mouseover", function(e) {

    $( '.hc-a' ).on( "mouseenter", function(e) {
        var obj = $(this);
        ShowPop(obj);
    });
    $( '.hc-a' ).on( "mouseleave", function() {
        var obj = $(this).find(".hc-b2").eq(0);
        HidePop(obj);
    });


    $( '.hc-a' ).on( "focusin", function() {
        // var obj = $(this);
        // ShowPop(obj);
    });

    $( '.hc-a' ).on( "focusout", function() {
        // var obj = $(this).find(".hc-b2").eq(0);
        // HidePop(obj)
    });

    console.log( "document ready completed" );
});


        

        console.log("hover popups handler added");
        //isolated jQuery end;
    });
})(window.jQuery.noConflict(true));

        
