(function ($) {
    $.fn.hovercard = function (options) {

        //Set defauls for the control
        var defaults = {
            width: 300,
            openOnLeft: false,
            openOnTop: false,
            cardImgSrc: "",
            detailsHTML: "",
            background: "#ffffff",
            delay: 0,
            onHoverIn: function () { },
            onHoverOut: function () { }
        };
        //Update unset options with defaults if needed
        var options = $.extend(defaults, options);

        //CSS for hover card. Change per your need, and move these styles to your stylesheet (recommended).
        if ($('#css-hovercard').length <= 0) {
            var hovercardTempCSS = '<style id="css-hovercard" type="text/css">' +
                                    '.hc-preview { position: relative; display:inline; }' +
                                    '.hc-name { font-weight:bold; position:relative; display:inline-block; }' +
                                    '.hc-details { left:-10px; margin-right:80px; text-align:left; '+
                                         'font-family:Sans-serif !important; font-size:12px !important; color:#666 !important; line-height:1.5em; '+
                                         'border:solid 1px #ddd; position:absolute; ' +
                                         '-moz-border-radius:3px;-webkit-border-radius:3px;border-radius:3px;top:-10px; '+
                                         '-moz-box-shadow:5px 5px 5px #888;-webkit-box-shadow:5px 5px 5px #888;box-shadow:5px 5px 5px #888;'+
                                         ' padding:2em 10px 10px; display:none;}' +
                                    '.hc-pic { width:70px; margin-top:-1em; float:right;  }' +
                                    '.hc-details-open-left { left: auto; right:-10px; text-align:right; margin-left:80px; margin-right:0; } ' +
                                    '.hc-details-open-left > .hc-pic { float:left; } ' +
                                    '.hc-details-open-top { bottom:-10px; top:auto; padding: 10px 10px 2em;} ' +
                                    '.hc-details-open-top > .hc-pic { margin-top:10px; float:right;  }' +
                                    '.hc-details .s-action{ position: absolute; top:8px; right:5px; } ' +
                                    '.hc-details .s-card-pad{ border-top: solid 1px #eee; margin-top:10px; padding-top:10px; overflow:hidden; } ' +
                                    '.hc-details-open-top .s-card-pad { border:none; border-bottom: solid 1px #eee; margin-top:0;padding-top:0; margin-bottom:10px;padding-bottom:10px; }' +
                                    '.hc-details .s-card .s-strong{ font-weight:bold; color: #555; } ' +
                                    '.hc-details .s-img{ float: left; margin-right: 10px; max-width: 70px;} ' +
                                    '.hc-details .s-name{ color:#222; font-weight:bold;} ' +
                                    '.hc-details .s-loc{ float:left;}' +
                                    '.hc-details-open-left .s-loc{ float:right;} ' +
                                    '.hc-details .s-href{ clear:both; float:left;} ' +
                                    '.hc-details .s-desc{ float:left; font-family: Georgia; font-style: italic; margin-top:5px;width:100%;} ' +
                                    '.hc-details .s-username{ text-decoration:none;} ' +
                                    '.hc-details .s-stats { display:block; float:left; margin-top:5px; clear:both; padding:0px;}' +
                                    '.hc-details ul.s-stats li{ list-style:none; float:left; display:block; padding:0px 10px !important; border-left:solid 1px #eaeaea;} ' +
                                    '.hc-details ul.s-stats li:first-child{ border:none; padding-left:0 !important;} ' +
                                    '.hc-details .s-count { font-weight: bold;} ' +
                                '.</style>")';

            $(hovercardTempCSS).appendTo('head');
        }
        //Executing functionality on all selected elements
        return this.each(function () {
            var obj = $(this);

            //wrap a parent span to the selected element
            obj.wrap('<div class="hc-preview" />');

            //add a relatively positioned class to the selected element
            obj.addClass("hc-name");

            //if card image src provided then generate the image elementk
            var hcImg = '';
            if (options.cardImgSrc.length > 0) {
                hcImg = '<img class="hc-pic" src="' + options.cardImgSrc + '" />';
            }

            //generate details span with html provided by the user
            var hcDetails = '<div class="hc-details" >' + hcImg + options.detailsHTML + '</div>';

            //append this detail after the selected element
            obj.after(hcDetails);
            obj.siblings(".hc-details").eq(0).css({ 'width': options.width, 'background': options.background });

            //toggle hover card details on hover
            obj.closest(".hc-preview").hover(function () {
                    
                var $this = $(this);
                adjustToViewPort($this);

                //Up the z index for the .hc-name to overlay on .hc-details
                $this.css("zIndex", "200");
                obj.css("zIndex", "100").find('.hc-details').css("zIndex", "50");

                var curHCDetails = $this.find(".hc-details").eq(0);
                curHCDetails.stop(true, true).delay(options.delay).fadeIn();


                //Default functionality on hoverin, and also allows callback
                if (typeof options.onHoverIn == 'function') {
                    //Callback function                    
                    options.onHoverIn.call(this);                    
                }


            }, function () {
                //Undo the z indices 
                $this = $(this);

                $this.find(".hc-details").eq(0).stop(true, true).fadeOut(300, function () {
                    $this.css("zIndex", "0");
                    obj.css("zIndex", "0").find('.hc-details').css("zIndex", "0");

                    if (typeof options.onHoverOut == 'function') {
                        options.onHoverOut.call(this);
                    }
                });
            });

            //Opening Directions adjustment
            function adjustToViewPort(hcPreview) {

                var hcDetails = hcPreview.find('.hc-details').eq(0);
                var hcPreviewRect = hcPreview[0].getBoundingClientRect();

                var hcdTop = hcPreviewRect.top - 20; //Subtracting 20px of padding;
                var hcdRight = hcPreviewRect.left + 35 + hcDetails.width(); //Adding 35px of padding;
                var hcdBottom = hcPreviewRect.top + 35 + hcDetails.height(); //Adding 35px of padding;
                var hcdLeft = hcPreviewRect.top - 10; //Subtracting 10px of padding;

                //Check for forced open directions, or if need to be autoadjusted
                if (options.openOnLeft || (hcdRight > window.innerWidth)) {
                    hcDetails.addClass("hc-details-open-left");
                } else {
                    hcDetails.removeClass("hc-details-open-left");
                }
                if (options.openOnTop || (hcdBottom > window.innerHeight)) {
                    hcDetails.addClass("hc-details-open-top");
                } else {
                    hcDetails.removeClass("hc-details-open-top");
                }
            };



        });

    };
})(jQuery);