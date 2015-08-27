package fetch_rss

var economistHomepage = []byte(`<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml"
      xml:lang="en" lang="en" dir="ltr"
      xmlns:og="http://ogp.me/ns#"
      xmlns:fb="https://www.facebook.com/2008/fbml">

<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
  <title>The Economist - World News, Politics, Economics, Business & Finance</title>
  <link rel="shortcut icon" href="http://cdn.static-economist.com/sites/default/files/econfinal_favicon.ico" type="image/x-icon" />
<script language="javascript">
  var gs_channels = "default";
</script>
<script language="javascript" type="text/javascript"
src="//economist.gscontxt.net/main/channels.cgi?url=http%3A%2F%2Fwww.economist.com%2Fnode%2F21555491">
</script>
<script type="text/javascript">
  if (typeof gs_channels === "undefined") { gs_channels = "default"; }
  function GPTGetCookie(cname) {
      var name = cname + "=";
      var ca = document.cookie.split(';');
      for(var i=0; i < ca.length; i++) {
          var c = ca[i];
          while (c.charAt(0) === ' ') c = c.substring(1);
          if (c.indexOf(name) === 0) return decodeURIComponent(c.substring(name.length, c.length));
      }
      return false;
  }
  var googletag = googletag || {};
  var econgpt = econgpt || {};
  googletag.cmd = googletag.cmd || [];
  (function() {    var gads = document.createElement("script");
    gads.async = true;
    gads.type = "text/javascript";
    var useSSL = "https:" == document.location.protocol;
    gads.src = (useSSL ? "https:" : "http:") + "//www.googletagservices.com/tag/js/gpt.js";
    var node = document.getElementsByTagName("script")[0];
    node.parentNode.insertBefore(gads, node);  })();
</script><script type="text/javascript">
  googletag.cmd.push(function() {  econgpt.gpt_leaderboard_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [[728,90], [970,90], [970,250]], 'gpt_leaderboard_ad').addService(googletag.pubads()).setTargeting('dcopt', 'ist').setTargeting('pos', 'ldr_top').setTargeting('tile', '1');  econgpt.gpt_mpu_no_wrapper_top_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [20,20], 'gpt_mpu_no_wrapper_top_ad').addService(googletag.pubads()).setTargeting('pos', 'mpu_no_wrapper_top').setTargeting('tile', '6');  econgpt.gpt_top_mpu_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [[300,250], [350,900], [350,600], [300,600]], 'gpt_top_mpu_ad').addService(googletag.pubads()).setTargeting('pos', 'mpu_top').setTargeting('tile', '7');  econgpt.gpt_ribbon_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [351,49], 'gpt_ribbon_ad').addService(googletag.pubads()).setTargeting('pos', 'absribbon').setTargeting('tile', '8');  econgpt.gpt_pencil_slug_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [970,35], 'gpt_pencil_slug_ad').addService(googletag.pubads()).setTargeting('pos', 'pslug_top').setTargeting('tile', '3');  econgpt.gpt_subscription_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [223,90], 'gpt_subscription_ad').addService(googletag.pubads()).setTargeting('pos', 'sub_top').setTargeting('tile', '2');  econgpt.gpt_slider_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [1,1], 'gpt_slider_ad').addService(googletag.pubads()).setTargeting('pos', 'slider').setTargeting('tile', '11');  econgpt.gpt_adcast = googletag.defineSlot("/5605/teg.fmsq/j2ek", [250,1000], 'gpt_adcast').addService(googletag.pubads()).setTargeting('pos', 'adcast').setTargeting('tile', '12');  econgpt.gpt_button_ad_1 = googletag.defineSlot("/5605/teg.fmsq/j2ek", [125,125], 'gpt_button_ad_1').addService(googletag.pubads()).setTargeting('pos', 'button1');  econgpt.gpt_button_ad_2 = googletag.defineSlot("/5605/teg.fmsq/j2ek", [125,125], 'gpt_button_ad_2').addService(googletag.pubads()).setTargeting('pos', 'button2');  econgpt.gpt_button_ad_3 = googletag.defineSlot("/5605/teg.fmsq/j2ek", [125,125], 'gpt_button_ad_3').addService(googletag.pubads()).setTargeting('pos', 'button3');  econgpt.gpt_button_ad_4 = googletag.defineSlot("/5605/teg.fmsq/j2ek", [125,125], 'gpt_button_ad_4').addService(googletag.pubads()).setTargeting('pos', 'button4');  econgpt.gpt_button_ad_5 = googletag.defineSlot("/5605/teg.fmsq/j2ek", [125,125], 'gpt_button_ad_5').addService(googletag.pubads()).setTargeting('pos', 'button5');  econgpt.gpt_button_ad_6 = googletag.defineSlot("/5605/teg.fmsq/j2ek", [125,125], 'gpt_button_ad_6').addService(googletag.pubads()).setTargeting('pos', 'button6');  econgpt.gpt_bottom_mpu_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [[300,250], [50,50], [300,600]], 'gpt_bottom_mpu_ad').addService(googletag.pubads()).setTargeting('pos', 'mpu_bottom').setTargeting('tile', '9');  econgpt.gpt_bottom_right_mpu_ad = googletag.defineSlot("/5605/teg.fmsq/j2ek", [[300,250], [45,45]], 'gpt_bottom_right_mpu_ad').addService(googletag.pubads()).setTargeting('pos', 'mpu_bottom_right').setTargeting('tile', '10');    googletag.pubads().setTargeting('subs', 'n').setTargeting('sdn', 'n');
    var subs = 'n';
    // The ec_omniture_user_sub cookie has info about the user's subscription.
    // Set subs=[y/n] based on the cookie.
    var sub_arr_regex = /((print|digital)-subscriber|ent-product)/;
    if (sub_arr_regex.test(document.cookie) === true) {
      subs = 'y';
    }
    googletag.pubads().setTargeting("subs", subs);    // Set etear variable based on the cookie
    var etear_cookie_val = GPTGetCookie('ec_ads_etear');
    if (etear_cookie_val) {
      etear = etear_cookie_val;
      googletag.pubads().setTargeting("etear", etear);
    }
    googletag.pubads().enableSingleRequest();
    googletag.pubads().setTargeting("gs_cat", gs_channels);
    googletag.pubads().collapseEmptyDivs();    googletag.enableServices();
  });
</script><script type="text/javascript">/* global Econ */
/**
 * Add event listeners to GPT ad rendering.
 *
 * @param context
 */
window.Econ = typeof Econ !== 'undefined' ? Econ : {};
Econ.gpt = Econ.gpt || {};
Econ.gpt.attach = Econ.gpt.attach || {};(function(document){
  "use strict";  var deltaPrefix = 'block-ec_ads-',
      deltaSuffix = '_ad';
  Econ.gpt.wide_leader = false;
  Econ.gpt.house_ad_leader = false;
  /**
   * Implements Econ.gpt.attach.HOOK
   * @param googletag_ref
   */
  Econ.gpt.attach.ecAds = function(googletag_ref) {
    googletag_ref.pubads().addEventListener('slotRenderEnded', function(event) {
      var pos = event.slot.getTargeting("pos")[0],
          adBlock,
          adLabel,
          domSlotID = event.slot.getSlotId().getDomId(),
          domSlot = document.getElementById(domSlotID),
          domStyles,
          adSize = Number(event.size[0]),
          domParents,
          billBoardWidth = 970;      // When ads are rendered in the ldr_top slot,
      // we check the size so we know if we need to show the pencil slug or not.
      if (pos === "ldr_top") {
        // If it is fully wide, then remove the subscription ad.
        if (adSize === billBoardWidth) {
          Econ.gpt.wide_leader = true;
          // We need to call this behavior because if we are in async mode, the ads are called AFTER document-ready
          // and the behavior may have already been called before the ad finished rendering.
          Drupal.behaviors.ecGPTPencilSlug(document);
        }
        else {
          // Make the top ad sticky if it's a narrower leader ad and it is not a house ad.
          if (!Econ.gpt.house_ad_leader) {    
            Econ.gpt.stickyAds();
          }  
        }
      }      // Remove empty ad slots. We stay in sync with the "collapsed div" motif. Google adds a display: none to ads
      // that are empty, so we just put our wrapper insync with that.
      domStyles = domSlot !== null ? window.getComputedStyle(domSlot) : '';
      domParents = domSlot !== null ? domSlot.parentNode.parentNode.parentNode : '';
      if (domSlot !== null &&
         domParents.getAttribute('class').indexOf('block-ec_ads') !== -1) {
        if (domStyles.getPropertyValue('display') === 'none' || adSize === 1) {
          domParents.style.display = 'none';
        }
        else {
          domParents.style.display = '';
        }
      }      // For special reports sponsor ads. The POS for these is "undefined".
      if (domSlotID.substr(0, 21) === 'gpt_sponsored_sr_logo') {
        domSlot.style.float = 'right';
      }      // For mobile, we are displaying ads using the refresh method,
      // which means we need to restore the Advertisemnt label.
      if (pos.substr(0,10) === 'mobile_mpu') {
        adBlock = document.getElementById(deltaPrefix + pos + deltaSuffix);
        if (adBlock.className.indexOf('processed') === -1) {
          adBlock.insertBefore(adLabel, adBlock.firstChild);
          adBlock.className += ' processed';
        }
      }
    });
  };
})(window.document);
</script>
<link rel="canonical" href="http://www.economist.com/" />
<meta name="description" content="The Economist offers authoritative insight and opinion on international news, politics, business, finance, science, technology and the connections between them." />
<meta name="pubdate" content="20120516" />
<meta name="revisit-after" content="1 day" />
<meta name="twitter:card" content="summary_large_image" />
<meta name="twitter:site" content="@TheEconomist" />
<meta property="fb:app_id" content="173277756049645" />
<meta property="og:image" content="http://cdn.static-economist.com/sites/default/files/the-economist-logo.gif" />
<meta property="og:site_name" content="The Economist" />
  <link rel="apple-touch-icon" sizes="144x144" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-144x144.png" /><link rel="apple-touch-icon" sizes="120x120" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-120x120.png" /><link rel="apple-touch-icon" sizes="114x114" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-114x114.png" /><link rel="apple-touch-icon" sizes="72x72" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-72x72.png" /><link rel="apple-touch-icon" href="http://cdn.static-economist.com/sites/default/files/images/icons/touch-icon-iphone.png" /><link rel="apple-touch-icon-precomposed" sizes="144x144" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-144x144-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="120x120" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-120x120-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="114x114" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-114x114-precomposed.png"><link rel="apple-touch-icon-precomposed" sizes="72x72" href="http://cdn.static-economist.com/sites/default/files/images/icons/apple-touch-icon-72x72-precomposed.png"><link rel="apple-touch-icon-precomposed" href="http://cdn.static-economist.com/sites/default/files/images/icons/touch-icon-iphone-precomposed.png">  <link type="text/css" rel="stylesheet" media="all" href="http://cdn.static-economist.com/sites/default/files/css/css_5b7f88bc4d30779c79a733cfcf63a4c8.css" />
  <!--[if IE 9]>
<link rel="stylesheet" href="/sites/all/themes/econfinal/styles/ie9.css?D" />
<![endif]-->
<!--[if IE 8]>
<link rel="stylesheet" href="/sites/all/themes/econfinal/styles/ie8.css?D" />
<script src="http://html5shiv.googlecode.com/svn/trunk/html5.js"></script>
<![endif]-->  <link rel="publisher" href="https://plus.google.com/100470681032489535736" />
  <script type="text/javascript" src="http://cdn.static-economist.com/sites/default/files/js/js_951d9d70881ff14d49a42e64fdef5717.js"></script>
<script type="text/javascript">
<!--//--><![CDATA[//><!--
jQuery.extend(Drupal.settings, {"basePath":"\/","offline":{"now":1440682027},"insert":{"widgets":{"imagefield_widget":{"fields":{"use_original_size":"input:checked[name$=\"[use_original_size]\"]","image_link":"input:textfield[name$=\"[image_link]\"]"}}}},"jcarousel":{"ajaxPath":"\/jcarousel\/ajax\/views"},"ec_user_settings":{"auto_select_country":true},"omniture":{"edge_server":"economistcomprod","click_tracking":[{"name":"hp_cover_share_facebook","selector":"data-ec-omniture","link_track_vars":"events","events":"event32","event":"click","edge_server":"economistcomprod","link_track_events":"event32"},{"name":"hp_cover_share_twitter","selector":"data-ec-omniture","link_track_vars":"events","events":"event31","event":"click","edge_server":"economistcomprod","link_track_events":"event31"},{"name":"hp_cover_share_linkedin","selector":"data-ec-omniture","link_track_vars":"events","events":"event34","event":"click","edge_server":"economistcomprod","link_track_events":"event34"},{"name":"hp_cover_share_plusone","selector":"data-ec-omniture","link_track_vars":"events","events":"event35","event":"click","edge_server":"economistcomprod","link_track_events":"event35"},{"name":"hp_cover_share_pinterest","selector":"data-ec-omniture","link_track_vars":"events","events":"event37","event":"click","edge_server":"economistcomprod","link_track_events":"event37"},{"name":"home|minimap","selector":"#footer p a","link_track_vars":"prop15","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|facebook","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|twitter","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|linked-in","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|plusone","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|tumblr","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|instagram","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|youtube","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|rss","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"rightrail|social_share|newsletters","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"home|touts|issue|cover","selector":"data-ec-omniture","link_track_vars":"prop15","event":"click","edge_server":"economistcomprod"},{"name":"home|touts|issue|subscribe","selector":"data-ec-omniture","link_track_vars":"prop15","event":"click","edge_server":"economistcomprod"},{"name":"home|touts|issue|printedition","selector":"data-ec-omniture","link_track_vars":"prop15","event":"click","edge_server":"economistcomprod"},{"name":"recommend_comment","selector":"data-ec-omniture","events":"event36","eVar31":"comment","link_track_events":"event36","prop18":null,"link_track_vars":"prop18,eVar31,events","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|contact_us","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|help","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|my_account","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|subscribe","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|print_edition","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|digital_editions","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|events","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|jobs_economist_com","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|timekeeper_saved_articles","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|united_states","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|britain","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|europe","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|china","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|asia","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|americas","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|middle_east_africa","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|international","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|business_finance","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|economics","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|markets_data","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|science_technology","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|special_reports","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|culture","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|multimedia_library","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_debates","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|letters_to_the_editor","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_quiz","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|buttonwoods_notebook","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|democracy_in_america","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|erasmus","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|free_exchange","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|game_theory","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|graphic_detail","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|gulliver","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|prospero","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_explains","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|topics","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|economics_a_z","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|style_guide","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_world_in_2015","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|which_mba_","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|mba_services","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_gmat_tutor","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|executive_education_navigator","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|reprints_and_permissions","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_group","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_intelligence_unit","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_intelligence_unit_store","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_corporate_network","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|ideas_people_media","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|intelligent_life","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|roll_call","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|cq","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|eurofinance","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|the_economist_store","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"mini_map_home|view_complete_site_index_","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|contact_us","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|help","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|about_us","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|advertise_with_us","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|editorial_staff","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|staff_books","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|careers","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|site_index","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|accessibility","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|privacy_policy","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|cookies_info","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"},{"name":"footer_home|terms_of_use","selector":"data-ec-omniture","event":"click","edge_server":"economistcomprod"}],"click_tracking_frames":["click_tracking_top_fb","click_tracking_footer_fb"],"click_tracking_top_fb":[{"name":"top_fb","selector":"data-ec-omniture-frame","ud_process_prop":"prop18","ud_process_val":"username","link_track_vars":"prop18,prop5,eVar24,eVar31,events","prop18":"","prop5":"Published August 27, 2015 13:26:54","eVar24":"Published August 27, 2015 13:26:54","eVar31":"homepage","events":"event32","event":"click","edge_server":"economistcomprod","link_track_events":"event32"}],"click_tracking_footer_fb":[{"name":"footer_fb","selector":"data-ec-omniture-frame","ud_process_prop":"prop18","ud_process_val":"username","link_track_vars":"prop18,prop5,eVar24,eVar31,events","prop18":"","prop5":"Published August 27, 2015 13:26:54","eVar24":"Published August 27, 2015 13:26:54","eVar31":"homepage","events":"event32","event":"click","edge_server":"economistcomprod","link_track_events":"event32"}]},"socialShare":{"shareCoverLink":"\u003cdiv class=\"link-social\"\u003eShare\u003c\/div\u003e","issueNID":"21661557","shareHTML":"  \u003cdiv class=\"qtip-click-content\"\u003e\n    \u003cul class=\"share-icons\"\u003e\n      \u003cli class=\"share-icon-facebook\"\u003e\u003ca href=\"http:\/\/www.facebook.com\/dialog\/feed?redirect_uri=http%3A%2F%2Fwww.economist.com%2Fclose-me.html\u0026amp;app_id=173277756049645\u0026amp;link=http%3A%2F%2Fwww.economist.com%2Fprintedition%2F2015-08-22%3Ffsrc%3Dscn%252Ffb_ec%252FThe%2520Economist%2520-%2520Aug%252022nd%25202015\u0026amp;name=The%20Economist%20-%20Aug%2022nd%202015\u0026amp;description=This%20week%20in%20The%20Economist%3A%20Genetic%20engineering%20%7C%20NGOs%20in%20China%20%7C%20Britain%E2%80%99s%20Labour%20Party%20%7C%20Higher%20education%20in%20America%20%7C%20The%20Ashley%20Madison%20hack\" target=\"_blank\" data-ec-omniture=\"hp_cover_share_facebook\"\u003eFacebook\u003c\/a\u003e\u003c\/li\u003e\n      \u003cli class=\"share-icon-twitter\"\u003e\u003ca href=\"https:\/\/twitter.com\/intent\/tweet?status=This%20week%20in%20The%20Economist%3A%20Genetic%20engineering%20%7C%20NGOs%20in%20China%20%7C%20Britain%E2%80%99s%20Labour%20Party%20http%3A%2F%2Fecon.st%2F1Kxi3Sq\" target=\"_blank\" data-ec-omniture=\"hp_cover_share_twitter\"\u003eTwitter\u003c\/a\u003e\u003c\/li\u003e \u003c!-- 550 pixels wide and 450 pixels high --\u003e\n      \u003cli class=\"share-icon-linkedin\"\u003e\u003ca href=\"http:\/\/www.linkedin.com\/shareArticle?mini=true\u0026amp;ro=false\u0026amp;summary=This%20week%20in%20The%20Economist%3A%20Genetic%20engineering%20%7C%20NGOs%20in%20China%20%7C%20Britain%E2%80%99s%20Labour%20Party%20%7C%20Higher%20education%20in%20America%20%7C%20The%20Ashley%20Madison%20hack\u0026amp;source=The%20Economist\u0026amp;url=http%3A%2F%2Fwww.economist.com%2Fprintedition%2F2015-08-22%3Ffsrc%3Dscn%252Fln_ec%252FThe%2520Economist%2520-%2520Aug%252022nd%25202015\u0026amp;title=The%20Economist%20-%20Aug%2022nd%202015\u0026amp;description=This%20week%20in%20The%20Economist%3A%20Genetic%20engineering%20%7C%20NGOs%20in%20China%20%7C%20Britain%E2%80%99s%20Labour%20Party%20%7C%20Higher%20education%20in%20America%20%7C%20The%20Ashley%20Madison%20hack\" target=\"_blank\" data-ec-omniture=\"hp_cover_share_linkedin\"\u003eLinkedin\u003c\/a\u003e\u003c\/li\u003e\n      \u003cli class=\"share-icon-google\"\u003e\u003ca href=\"https:\/\/plus.google.com\/share?url=http%3A%2F%2Fwww.economist.com%2Fprintedition%2F2015-08-22%3Ffsrc%3Dscn%252Fgn_ec%252FThe%2520Economist%2520-%2520Aug%252022nd%25202015\u0026amp;t=The%20Economist%20-%20Aug%2022nd%202015\u0026amp;description=This%20week%20in%20The%20Economist%3A%20Genetic%20engineering%20%7C%20NGOs%20in%20China%20%7C%20Britain%E2%80%99s%20Labour%20Party%20%7C%20Higher%20education%20in%20America%20%7C%20The%20Ashley%20Madison%20hack\" target=\"_blank\" data-ec-omniture=\"hp_cover_share_plusone\"\u003eGoogle+\u003c\/a\u003e\u003c\/li\u003e\n      \u003cli class=\"share-icon-pinterest\"\u003e\u003ca href=\"http:\/\/pinterest.com\/pin\/create\/button\/?url=http%3A%2F%2Fwww.economist.com%2Fprintedition%2F2015-08-22%3Ffsrc%3Dscn%252Fpn_ec%252FThe%2520Economist%2520-%2520Aug%252022nd%25202015\u0026amp;description=The%20Economist%20-%20Aug%2022nd%202015\" target=\"_blank\" data-ec-omniture=\"hp_cover_share_pinterest\"\u003ePinterest\u003c\/a\u003e\u003c\/li\u003e\n    \u003c\/ul\u003e\n  \u003c\/div\u003e"},"cdn":{"cdnPath":"http:\/\/cdn.static-economist.com","cdnPathPNG":"http:\/\/cdn.static-economist.com"},"ecGPT":{"isAsync":true},"ecAds":{"ad_site":"teg","ad_zone":"j2ek","sitecode":"fmsq","ad_call_type":"gpt","subs":"n","sdn":"n","sz_nojs":"20x20"},"ec_mostx_active_tab":"commented","ecHomepage":{"rotato_initial_delay":5,"rotato_interval":4,"isFront":true},"EconStyles":[{"title":"Pull quote","description":"Better reading","outputTemplate":"\u003cdiv class=\"pullquote\"\u003e{$selection}\u003c\/div\u003e"}],"omnitureOverride":{"overriddenVars":{"s.pageName":["homepage","homepage"],"s.channel":["home","home"],"s.prop1":["homepage","homepage"],"s.prop2":["",""]},"defaultVars":{"s.pageName":["homepage","homepage"],"s.channel":["home","home"],"s.prop1":["",""],"s.prop2":["",""]},"end":["1-01-1970","1-01-1970"]},"bluekai":{"siteID":18452},"ecCookieMessage":{"cid":"ec_cookie_message_","version":0},"ecFloodlight":{"regions":{"UK":{"anon":"\/\/2382269.fls.doubleclick.net\/activityi;src=2382269;type=mains756;cat=2013-091","reg":"\/\/2382269.fls.doubleclick.net\/activityi;src=2382269;type=mains756;cat=2013-998","sub":"\/\/2382269.fls.doubleclick.net\/activityi;src=2382269;type=mains756;cat=2013-079"},"NA":{"anon":"\/\/1113415.fls.doubleclick.net\/activityi;src=1113415;type=econo358;cat=econo673","reg":"\/\/1113415.fls.doubleclick.net\/activityi;src=1113415;type=econo358;cat=econo673","sub":"\/\/1113415.fls.doubleclick.net\/activityi;src=1113415;type=econo358;cat=econo673"}},"countries":{"AU":{"anon":"\/\/4122168.fls.doubleclick.net\/activityi;src=4122168;type=subsc565;cat=theec313","reg":"\/\/4122168.fls.doubleclick.net\/activityi;src=4122168;type=subsc565;cat=theec313","sub":""},"HK":{"anon":"\/\/2533238.fls.doubleclick.net\/activityi;src=2533238;type=dsp20753;cat=dsp_e721","reg":"\/\/2533238.fls.doubleclick.net\/activityi;src=2533238;type=dsp20753;cat=dsp_e721","sub":""},"SG":{"anon":"\/\/2533238.fls.doubleclick.net\/activityi;src=2533238;type=dsp20753;cat=dsp_e721","reg":"\/\/2533238.fls.doubleclick.net\/activityi;src=2533238;type=dsp20753;cat=dsp_e721","sub":""}}},"ecRegOverlayShowOnLoad":true,"EcRegOverlayEnabled":true,"ecPixelTracking":{"adroit":{"countries":["US"],"scriptUrl":"\/\/www.imiclk.com\/cgi\/r.cgi?m=3\u0026mid=Ay3XMfb1\u0026ptid=HOME"},"adsrvr":{"countries":["IN"],"regions":["AP"],"scriptUrl":"\/\/insight.adsrvr.org\/track\/evnt\/?adv=5gu6big\u0026ct=0:2qbfzvew\u0026fmt=3"},"fosina":{"countries":["US"],"scriptUrl":"\/\/secure.fastclick.net\/w\/roitrack.cgi?aid=1000045446"},"media_math":{"countries":["IN"],"regions":["AP"],"imageUrl":"\/\/insight.adsrvr.org\/track\/evnt\/?adv=5gu6big\u0026ct=0:zx1ys9qy\u0026fmt=3"},"the_trade_desk":{"countries":["IN"],"regions":["AP"],"scriptUrl":"\/\/pixel.mathtag.com\/event\/js?mt_id=413291\u0026mt_adid=120100\u0026v1=\u0026v2=\u0026v3=\u0026s1=\u0026s2=\u0026s3="},"tribalFusion":{"countries":["US"],"scriptUrl":"http:\/\/a.tribalfusion.com\/i.cid?c=548473\u0026d=30\u0026page=landingPage"},"turn":{"countries":["IN"],"regions":["AP"],"imageUrl":"\/\/r.turn.com\/r\/beacon?b2=R_QfAflH-jzroqgoR-7UpyoPLd4h6zhRDu90-2mSJwKPFE1jUk1i_ZKaX90YPbtpfum-xQNkFnpSs2wXvdPFSg\u0026cid=\u0026bprice="},"polar":{"scriptUrl":"\/\/plugin.mediavoice.com\/customers\/economist\/mediavoice.js"}},"adsStickyTiming":3000,"q":"node\/21555491","alias":"node\/21555491","nid":"21555491","uid":0,"regional_covers":{"A":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","AP":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","E":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","EU":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","LA":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","ME":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","NA":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","UK":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg","default":"http:\/\/cdn.static-economist.com\/sites\/default\/files\/imagecache\/print-cover-thumbnail\/print-covers\/20150822_cuk400.jpg"},"ec_google_cse":{"cx":"013751040265774567329:pqjb-wvrj-q","resultsUrl":"http:\/\/www.economist.com\/search\/gcs","gname":"economist-search","queryParameterName":"ss","fromurl":"\/","noResultsString":"Your query returned no results. Please try a different search term. (Did you check your spelling? You can also try rephrasing your query or using more general search terms.)","isResultsPage":false,"omnitureTrackerId":"gsc-i-id1"},"loginform":{"action":"https:\/\/www.economist.com\/user\/login?destination=node%2F21555491"},"formvalidation":{"formname":"#user-login"},"popups":{"originalCSS":{"\/modules\/node\/node.css":1,"\/modules\/poll\/poll.css":1,"\/modules\/system\/defaults.css":1,"\/modules\/system\/system.css":1,"\/modules\/system\/system-menus.css":1,"\/modules\/user\/user.css":1,"\/sites\/all\/modules\/cck\/theme\/content-module.css":1,"\/sites\/all\/modules\/contrib\/oembed\/oembed.css":1,"\/sites\/all\/modules\/ctools\/css\/ctools.css":1,"\/sites\/all\/modules\/custom\/ec_cookie_message\/css\/ec_cookie_message.css":1,"\/sites\/all\/modules\/custom\/ec_messaging\/css\/ec_messaging.css":1,"\/sites\/all\/modules\/date\/date.css":1,"\/sites\/all\/modules\/date\/date_popup\/themes\/datepicker.1.7.css":1,"\/sites\/all\/modules\/date\/date_popup\/themes\/jquery.timeentry.css":1,"\/sites\/all\/modules\/ec_components\/css\/ec_tip.css":1,"\/sites\/all\/modules\/ec_offline\/styles\/ec_offline_online.css":1,"\/sites\/all\/modules\/ec_vote\/css\/ec_vote.css":1,"\/sites\/all\/modules\/filefield\/filefield.css":1,"\/sites\/all\/modules\/logintoboggan\/logintoboggan.css":1,"\/sites\/all\/modules\/mollom\/mollom.css":1,"\/sites\/all\/modules\/og\/theme\/og.css":1,"\/sites\/all\/modules\/video_filter\/video_filter.css":1,"\/sites\/all\/modules\/wysiwyg_fontcolor\/tinymce\/fontcolor\/css\/fontcolor.css":1,"\/sites\/all\/modules\/ec_comments\/css\/ec_comments.css":1,"\/sites\/all\/modules\/cck\/modules\/fieldgroup\/fieldgroup.css":1,"\/sites\/all\/modules\/views\/css\/views.css":1,"\/sites\/all\/modules\/ec_omniture\/ec_omniture_link_tracking.css":1,"\/sites\/all\/modules\/custom\/ec_social\/ec_social.css":1,"\/sites\/all\/modules\/cck\/modules\/content_multigroup\/content_multigroup.css":1,"\/sites\/all\/modules\/ec_components\/css\/ec_components.css":1,"\/sites\/all\/modules\/ec_ads\/css\/ec_ads.css":1,"\/sites\/all\/modules\/ec_ads\/css\/ec_ads_topic_ribbon.css":1,"\/sites\/all\/modules\/custom\/ec_social\/ec-social-share-follow.css":1,"\/sites\/all\/modules\/ec_ads\/css\/ec_ads_slider.css":1,"\/sites\/all\/modules\/ec_ads\/css\/ec_ads_overlay.css":1,"\/sites\/all\/modules\/custom\/wysiwyg_econstyles\/tinymce\/econstyles\/css\/econstyles_editor.css":1,"\/sites\/all\/modules\/ec_offers\/css\/ec-offers.css":1,"\/sites\/all\/modules\/wysiwyg_fontcolor\/tinymce\/fontcolor\/css\/fontcolor_editor.css":1,"\/sites\/all\/modules\/ec_ads\/ec_gpt\/css\/ec_gpt.css":1,"\/sites\/all\/modules\/custom\/ec_google_cse\/css\/ec_google_cse.css":1,"\/sites\/all\/themes\/econfinal\/styles\/reset.css":1,"\/sites\/all\/themes\/econfinal\/styles\/grid.css":1,"\/sites\/all\/themes\/econfinal\/styles\/style.css":1,"\/sites\/all\/themes\/econfinal\/styles\/components.css":1,"\/sites\/all\/themes\/econfinal\/styles\/typography.css":1,"\/sites\/all\/themes\/econfinal\/styles\/skins.css":1,"\/sites\/all\/themes\/econfinal\/styles\/user-profile-comments.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-guest_contribution-tpl.css":1,"\/sites\/all\/themes\/econfinal\/styles\/page-header.css":1,"\/sites\/all\/themes\/econfinal\/styles\/page-navigation.css":1,"\/sites\/all\/themes\/econfinal\/styles\/page-ads.css":1,"\/sites\/all\/themes\/econfinal\/styles\/page-footer.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-homepage.css":1,"\/sites\/all\/themes\/econfinal\/styles\/page-taxonomy-term.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-issue.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-page-tools.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-print-cover.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-print-sub-ad.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-poll-article.css":1,"\/sites\/all\/themes\/econfinal\/styles\/captcha.css":1,"\/sites\/all\/themes\/econfinal\/styles\/site-index.css":1,"\/sites\/all\/themes\/econfinal\/styles\/inform-component.css":1,"\/sites\/all\/themes\/econfinal\/styles\/pager.css":1,"\/sites\/all\/themes\/econfinal\/styles\/block-blogs.css":1,"\/sites\/all\/themes\/econfinal\/styles\/block-products-events.css":1,"\/sites\/all\/themes\/econfinal\/styles\/block-radio.css":1,"\/sites\/all\/themes\/econfinal\/styles\/block-most-popular.css":1,"\/sites\/all\/themes\/econfinal\/styles\/block-brightcove_player-playlist.css":1,"\/sites\/all\/themes\/econfinal\/styles\/topic-page.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-rss-page.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-page-tpl.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-promotions-page-tpl.css":1,"\/sites\/all\/themes\/econfinal\/styles\/media-directory.css":1,"\/sites\/all\/themes\/econfinal\/styles\/ec-user.css":1,"\/sites\/all\/themes\/econfinal\/styles\/ec-user-block.css":1,"\/sites\/all\/themes\/econfinal\/styles\/ec-highlight-block.css":1,"\/sites\/all\/themes\/econfinal\/styles\/ec-glossary.css":1,"\/sites\/all\/themes\/econfinal\/styles\/ec-charts.css":1,"\/sites\/all\/themes\/econfinal\/styles\/ec-sections.css":1,"\/sites\/all\/themes\/econfinal\/styles\/node-event.css":1,"\/sites\/all\/themes\/econfinal\/styles\/interactive-map.css":1,"\/sites\/all\/themes\/econfinal\/styles\/vendor-prexifed.css":1},"originalJS":{"\/sites\/all\/modules\/ec_ads\/ec_gpt\/js\/ec_gpt.js":1,"\/sites\/all\/modules\/ec_ajax_form\/js\/ec_ajax_form.js":1,"\/sites\/all\/modules\/ec_base\/js\/ec_base.js":1,"\/sites\/all\/modules\/ec_base\/js\/ec_console.js":1,"\/sites\/all\/modules\/ec_base\/js\/ec_analytics.js":1,"\/sites\/all\/modules\/ec_base\/js\/plugins\/jquery.ec_timeago.js":1,"\/sites\/all\/modules\/ec_base\/js\/plugins\/jquery.cookie.js":1,"\/sites\/all\/modules\/ec_base\/js\/plugins\/jquery.json.js":1,"\/sites\/all\/libraries\/qtip\/jquery.qtip-1.0.0-rc3.min.js":1,"\/sites\/all\/modules\/ec_components\/js\/ec_tip.js":1,"\/sites\/all\/modules\/ec_vote\/js\/ec_vote.js":1,"\/sites\/all\/modules\/jquery_aop\/misc\/jquery.aop.js":1,"\/sites\/all\/modules\/jsonify\/jsonify.js":1,"\/sites\/all\/modules\/mollom\/mollom.js":1,"\/sites\/all\/modules\/og\/og.js":1,"\/sites\/all\/modules\/swftools\/shared\/swfobject2\/swfobject.js":1,"\/sites\/all\/modules\/ec_user\/js\/ec_user.js":1,"\/sites\/all\/libraries\/simplemodal\/jquery.simplemodal.1.4.1.min.js":1,"\/sites\/all\/libraries\/jquery_expander\/jquery.expander.min.js":1,"\/sites\/all\/modules\/ec_comments\/js\/ec_comments.js":1,"\/sites\/all\/modules\/ec_homepage\/ec_homepage.js":1,"\/sites\/all\/modules\/ec_base\/js\/ec_overlays.js":1,"\/sites\/all\/modules\/ec_ads\/js\/stepcarousel.js":1,"\/sites\/all\/modules\/ec_mostx\/js\/ec_mostx.js":1,"\/sites\/all\/modules\/ec_ads\/js\/ec_ads_slider.js":1,"\/sites\/all\/modules\/custom\/bluekai\/bluekai.js":1,"\/sites\/all\/modules\/custom\/ec_floodlight\/ec_floodlight.js":1,"\/sites\/all\/modules\/custom\/ec_overlay\/js\/ec-overlay.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_adroit\/js\/ec_pixel_tracking_adroit.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_adsrvr\/js\/ec_pixel_tracking_adsrvr.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_fosina\/js\/ec_pixel_tracking_fosina.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_media_math\/js\/ec_pixel_tracking_media_math.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_tapad\/ec_pixel_tracking_tapad.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_the_trade_desk\/js\/ec_pixel_tracking_the_trade_desk.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_tribal_fusion\/js\/ec_pixel_tracking_tribal_fusion.js":1,"\/sites\/all\/modules\/custom\/ec_pixel_tracking\/ec_pixel_tracking_turn\/js\/ec_pixel_tracking_turn.js":1,"\/sites\/all\/modules\/custom\/ec_polar\/ec_polar.js":1,"\/sites\/all\/modules\/ec_ads\/js\/ec_ads.js":1,"\/sites\/all\/modules\/custom\/ec_google_cse\/js\/ec_google_cse.js":1,"\/sites\/all\/modules\/custom\/ec_google_cse\/js\/non_require_js_loader.js":1,"\/sites\/all\/modules\/custom\/ec_social\/ec_social.js":1,"\/sites\/all\/modules\/ajax\/jquery\/jquery.a_form.packed.js":1,"\/sites\/all\/modules\/ajax\/ajax.js":1,"\/sites\/all\/modules\/custom\/ec_brightcove\/js\/ec_brightcove.base.js":1,"\/sites\/all\/modules\/custom\/ec_cookie_message\/js\/ec_cookie_message.js":1,"\/sites\/all\/modules\/custom\/ec_messaging\/js\/ec_messaging.js":1,"\/sites\/all\/modules\/ec_components\/js\/ec_components_plugins.js":1,"\/sites\/all\/modules\/ec_offline\/js\/ec_offline_redirect.js":1,"\/sites\/all\/modules\/ec_offline\/js\/ec_offline.js":1,"\/sites\/all\/modules\/ec_omniture\/ec_omniture_link_tracking.js":1,"\/sites\/all\/modules\/ec_ads\/js\/ec_ads_banner.js":1,"\/sites\/all\/modules\/ec_blogs\/js\/ec_blogs.js":1,"\/sites\/all\/themes\/econfinal\/js\/page-header.js":1,"\/sites\/all\/themes\/econfinal\/js\/node-page-tools.js":1,"\/sites\/all\/themes\/econfinal\/js\/node-poll-article.js":1,"\/sites\/all\/themes\/econfinal\/js\/site-index.js":1,"\/sites\/all\/themes\/econfinal\/js\/node-guest_contribution.js":1,"\/sites\/all\/themes\/econfinal\/js\/AC_OETags.js":1,"\/sites\/all\/themes\/econfinal\/js\/ec-blogs.js":1,"\/sites\/all\/themes\/econfinal\/js\/ec-messages.js":1,"\/sites\/all\/themes\/econfinal\/js\/interactive-map.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/drag.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/jquery.color.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/jquery.form.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/jquery.autogrow.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/jquery.timer.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/jquery.hoverIntent.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/modernizr_2.8.3.min.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/jquery.flash.min.js":1,"\/sites\/all\/themes\/econfinal\/js\/plugins\/jcarousellite_1.0.1.min.js":1,"\/sites\/all\/modules\/ec_ads\/ec_gpt\/js\/ec_gpt_sticky.js":1}}});
//--><!]]>
</script>
<script type="text/javascript">
<!--//--><![CDATA[//><!--
document.write(unescape('%3Cscript src="//service.maxymiser.net/cdn/economist/js/mmcore.js" type="text/javascript"%3E%3C/script%3E'))
//--><!]]>
</script>
<script type="text/javascript">
<!--//--><![CDATA[//><!--

      function homepagePlayerLoaded(id) {
        player = brightcove.api.getExperience(id);
        if (player) {
          experienceModule = player.getModule(APIModules.EXPERIENCE);
          menuModule = player.getModule(APIModules.MENU);
          experienceModule.addEventListener(BCExperienceEvent.TEMPLATE_READY, function(result) {
            var video = player.getModule(APIModules.VIDEO_PLAYER);
            video.addEventListener(BCMediaEvent.PLAY, function(e) {
              menuModule.closeMenuPage();
            });
          });
        }
      }
  
//--><!]]>
</script>
<script type="text/javascript">
<!--//--><![CDATA[//><!--
/**
 * The main function: wireUpBrightcoveExperienceOnLoaded
 * is the function that we have set in brightcove_player to be the default
 * function to be called when a BC player is called.
 * @ref http://docs.brightcove.com/en/video-cloud/smart-player-api/samples/set-ad-policy.html
 */
/*
* Custom Code: Brightcove Smart Analytics v2.0
*/

// Setup a new object.
ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d = {};
ECOmnitureBCLogger = function(msg) {
  Economist.console.log(msg);
}


ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.omnitureBCTemplateLoaded = function (experienceID) {
  // Smart vs Flash players use slightly different APIs.
  var player = Econ.bcIsFlashPlayer() ? brightcove.getExperience(experienceID) : brightcove.api.getExperience(experienceID);
  var modVP = player.getModule(brightcove.api.modules.APIModules.VIDEO_PLAYER);
  var modExp = player.getModule(brightcove.api.modules.APIModules.EXPERIENCE);
  var modCon = player.getModule(brightcove.api.modules.APIModules.CONTENT);
  var modSocial = player.getModule(APIModules.SOCIAL);
  var ads = player.getModule(APIModules.ADVERTISING);

  // Set the variables to the function to localize them across the other events.
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.player = player;
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modVP = modVP;
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modExp = modExp;
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modCon = modCon;
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modSocial = modSocial;
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.mediaPlayerName = experienceID;
  // Add the player ad keys.
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.ads = ads;

  // Setup the handlers.
  modExp.addEventListener(brightcove.api.events.ExperienceEvent.TEMPLATE_READY, ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onTemplateReady);
  if (player.type === "flash") {
    modExp.addEventListener(BCExperienceEvent.CONTENT_LOAD, ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.bcpOnContentLoad);
    modVP.addEventListener(BCVideoEvent.VIDEO_CHANGE, ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.bcpOnVideoChange);
    $('#' + experienceID).attr('data-bc-mode', 'flash');
  }

  // We assume only one mobile hero player.
  // @see ec_gallery_mobile.js
  if ($('#' + experienceID + '.BrightcoveExperience').length == 1 &&
    typeof Drupal.settings !== 'undefined' &&
    typeof Drupal.settings.node_type !== 'undefined' &&
    Drupal.settings.node_type === 'homepage') {
    player.getModule(brightcove.api.modules.APIModules.VIDEO_PLAYER)._canPlayWithoutUserInteraction = true;
    EconMobileHeroSetup(player);
  }
}

// Begin handlers for linking and sharing.
ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.bcpOnContentLoad = function(evt) {
  // Initialize share link.
  var tabBar = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.tabBar;
  var modVP  = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modVP;
  var playlistId = (tabBar) ? tabBar.getSelectedData().id : modVP.getCurrentVideo().lineupId;
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.bcpUpdateLink(modVP.getCurrentVideo().id, playlistId);
}

ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.bcpOnVideoChange = function(evt) {
  var tabBar = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.tabBar;
  var modVP  = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modVP;
  var modExp = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modExp;
  if (modExp.getReady()) { // If template is Ready.
    // Because TemplateReady has already fired we can now access the currentVideo and currentPlaylist from the tabBar module.
    var playlistId = (tabBar) ? tabBar.getSelectedData().id : modVP.getCurrentVideo().lineupId;
    ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.bcpUpdateLink(modVP.getCurrentVideo().id, playlistId);
  }
}

ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.bcpUpdateLink = function(videoId, playlistId) {
  var playlistKey = "bclid";
  var videoKey = "bctid";
  var modSocial = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modSocial;
  var modVP  = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modVP;
  var tabBar = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.tabBar;
  // If we are on the /multimedia page, then we use the playlist id, etc. else
  // we just shorten the url that we are on.
  var info = shortUrlAPI = null;
  var newLink = document.URL;

  if (!window.location.origin) {
    window.location.origin = window.location.protocol+"//"+window.location.host;
  }
  var baseURL = window.location.origin + "/multimedia";
  newLink = baseURL + "?" + playlistKey + "=" + playlistId + "&" + videoKey + "=" + videoId;

  // Shorten the URL if we can.
  shortUrlAPI = Drupal.settings.basePath + 'brightcove-player/short-url?url=' + encodeURIComponent(newLink);
  // Fetch the short-url in synchronous manner.
  $.ajax({
    url: shortUrlAPI,
    async: false,
    dataType: "json",
    success: function(result) {
      info = result;
    }
  });
  if (info != null) {
    modSocial.setLink(info.short_url);
  }
  else {
    modSocial.setLink(newLink);
  }
}

// **************** End handlers for linking and sharing. *************

/**
 * Even though the BC documentation says to set the policyHandler and set the keys,
 * it seems that this does not work and thus this function is not used in the onTemplateReady.
 * This should be looked into.
 * @ref: http://docs.brightcove.com/en/video-cloud/smart-player-api/samples/set-ad-policy.html
 */
ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.getAdPolicyHandler = function(adPolicy) {}

// Function used on templateReady.
// @see omnitureBCTemplateLoaded
ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onTemplateReady = function(evt) {
  try {
    var player = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.player;

    var modAds = player.getModule(APIModules.ADVERTISING);
    // We don't use the adPolicyHandler anymore.
    var adPolicyHandler = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.getAdPolicyHandler;
    var adKey = Drupal.settings.ecAds.ad_site + '.' + Drupal.settings.ecAds.sitecode + '/' + Drupal.settings.ecAds.ad_zone;
    var adPolicy = modAds.getAdPolicy(adPolicyHandler);

    var adServerURL = "https:" == document.location.protocol ? "https:" : "http:";
    adServerURL += "//ad.doubleclick.net/pfadx/";
    adServerURL += Drupal.settings.ecAds.ad_site + '.' + Drupal.settings.ecAds.sitecode + '/' + Drupal.settings.ecAds.ad_zone
    adPolicy.adServerURL = adServerURL;
    // Update the adPolicy for the ad-call for v1do.
    modAds.setAdPolicy(adPolicy);

    var modVP = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modVP;

    modVP.addEventListener(brightcove.api.events.MediaEvent.PLAY, ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onPlay);
    modVP.addEventListener(brightcove.api.events.MediaEvent.STOP, ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onStop);
    modVP.addEventListener(brightcove.api.events.MediaEvent.PROGRESS, ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onProgress);

    // Sharing and linking
    if (evt.target.experience) {
      return;
    }
    // Load up a library item.
    var modExp = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.modExp;
    var tabBar = modExp.getElementByID("playlistTabs");
    ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.tabBar = tabBar;
    // If we have a channel passed in, then set focus to that tab.
    if (typeof multimedia_channel !== 'undefined') {
      var tabs = tabBar.getData();
      for (var i=0, len=tabs.length; i<len; i++) {
        if (multimedia_channel == tabs[i].displayName.toLowerCase().replace(/[^a-z0-9]/g, '')) {
          // Set the tab, and then play the first video in that tab.
          tabBar.setSelectedIndex(i);
          modVP.loadVideo(tabs[i].videoIds[0]);
        }
      }
    }
  }
  catch (err) {
    ECOmnitureBCLogger(err);
  }
};

ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onPlay = function(evt) {
  var mediaLength = evt.duration; //Required video duration
  ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.mediaLength = mediaLength; // Needed in another function.
  var mediaOffset = Math.floor(evt.position); //Required video position
  var mediaID = (evt.media.id).toString(); //Required video id
  var mediaFriendly = evt.media.displayName; //Required video title
  var mediaName = mediaID + ":" + mediaFriendly; //Required Format video name
  var mediaRefID = evt.media.referenceId; //Optional reference id
  var mediaPlayerType = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.player.type; //Optional player type
  var mediaTagsArray = evt.media.tags; //Optional tags
  var mediaPlaylist = ""; //Optional Playlist
  var mediaTagsArray2 = new Array();
  var mediaPlayerName = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.mediaPlayerName;
  // HTML5 Players have their tags flattend to a 1-dimensional array already.
  if (mediaPlayerType == 'html') {
    mediaTagsArray2 = mediaTagsArray;
  }
  else {
    for (var i = 0; i < mediaTagsArray.length; i++) {
      mediaTagsArray2.push(mediaTagsArray[i]['name']);
    }
  }
  if (mediaPla</lenpe == "flash") { //Optional playlist id
    if (evt.media.lineupId != null) {
      mediaPlaylist = (evt.media.lineupId).toString();
    }

  } else {
    if (evt.media.playlistID !== null) {
      mediaPlaylist = (evt.media.playlistID).toString();
    }
  }
  /* Check for start of video */
  if (mediaOffset == 0) {
    /* These data points are optional. If using SC14, change context data variables to hard
    coded variable names and change trackVars above. */
    s.contextData['bc_tags'] = mediaTagsArray2.toString(); //Optional returns list of tags for current video. Flash only.
    s.contextData['bc_refid'] = mediaRefID; //Optional retursn reference id
    s.contextData['bc_player'] = mediaPlayerName; //Optional player name is currently hard coded. Will be dynamic in later releases.
    s.contextData['bc_playertype'] = mediaPlayerType; //Optional returns flash or html
    s.contextData['bc_playlist'] = mediaPlaylist; //Optional returns playlist number for current video.
    s.Media.open(mediaName, mediaLength, mediaPlayerName);
    s.Media.play(mediaName, mediaOffset);
  } else {
    s.Media.play(mediaName, mediaOffset);
  }
}

ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onStop = function (evt) {
  try {
    var mediaOffset = Math.floor(evt.position);
    var mediaLength = ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.mediaLength;
    if (mediaOffset == mediaLength) {
      s.Media.stop(me</lenme, mediaOffset);
      s.Media.close(mediaName);
    } else {
      s.Media.stop(mediaName, mediaOffset);
    }
  }
  catch(err) {
    ECOmnitureBCLogger(err);
  }
}

ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.onProgress = function(evt) {
  s.Media.monitor = function(s, media) {
    if (media.event == "MILESTONE") {
      /* Use to set additional data points during milestone calls */
      //s.Media.track(media.name); Uncomment if setting extra milestone data.
    }
  }
}

//--><!]]>
</script>
    <script type="text/javascript">
    $('html').addClass('js');
  </script>
  </head>

<body class="front not-logged-in page-node node-type-homepage one-sidebar sidebar-right path-node-21555491 world-menu business-menu economics-menu printedition-menu science-technology-menu culture-menu">
      
  <!-- Begin Tealium Integration -->
  <script type="text/javascript">
    <!--//--><![CDATA[//><!--
    var utag_data = {"node_id":"21555491","node_title":"Published August 27, 2015 13:26:54","node_type":"homepage","node_created":"1337152064","node_changed":"1440682014","node_author":"1930114","user_status":"anonymous","dfp_zone":"j2ek","dfp_site":"FMSQ","dfp_targeting":"5605\/teg.fmsq\/j2ek","news_keywords":"","section":"","issue_date":"","content_source":"web","content_type":"homepage","story_title":"","subsection":""};
    //--><!]]>
  </script>

      <!-- Loading script asynchronously -->
    <script type="text/javascript">
      (function(a,b,c,d){
        a='//tags.tiqcdn.com/utag/teg/main/prod/utag.js';
        b=document;c='script';d=b.createElement(c);d.src=a;d.type='text/java'+c;d.async=true;
        a=b.getElementsByTagName(c)[0];a.parentNode.insertBefore(d,a);
      })();
    </script>
  
  <!-- End Tealium Integration -->
        <div id="fb-root"></div>
<script>
  window.fbAsyncInit = function() {
    FB.init({
      status: true,
      cookie: true,
      xfbml: true,
      appId: '173277756049645',
      oauth: true,
      channelUrl: 'http://www.economist.com/fb/fb_channel.html'
    });
  };

  (function() {
    var e = document.createElement('script');
    e.src = document.location.protocol + '//connect.facebook.net/en_US/all.js';
    e.async = true;
    document.getElementById('fb-root').appendChild(e);
  }());
</script>                <div id="leaderboard" class="clearfix">
    <div id="block-ec_ads-leaderboard_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div id="leaderboard-ad"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_leaderboard_ad" data-cb-ad-id="Leaderboard ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_leaderboard_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26dcopt%3Dist%26pos%3Dldr_top&sz=728x90|970x90|970x250&tile=1&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26dcopt%3Dist%26pos%3Dldr_top&sz=728x90|970x90|970x250&tile=1&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>

<div id="block-ec_ads-subscription_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div id="subslug-ad"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_subscription_ad" data-cb-ad-id="Subscription ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_subscription_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dsub_top&sz=223x90&tile=2&c=328904146">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dsub_top&sz=223x90&tile=2&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>

<div id="block-ec_ads-pencil_slug_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_pencil_slug_ad" data-cb-ad-id="Pencil slug ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_pencil_slug_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dpslug_top&sz=970x35&tile=3&c=328904146">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dpslug_top&sz=970x35&tile=3&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

  </div>
<header class="mh">
  <!-- 
  To use as a standalone component please wrap it up into a <header class="mh"></header> 
-->
<div class="mh-stripe">
  <div class="mh-stripe-wrap">
    <ul class="mh-user-menu"><li class="first"><span>More from The Economist</span><ul><li class="first"><a href="/digital">The Economist digital editions</a></li>
<li><a href="/newsletters">Newsletters</a></li>
<li><a href="/events">Events</a></li>
<li><a href="http://jobs.economist.com">Jobs.Economist.com</a></li>
<li><a href="http://store.economist.com/">The Economist Store</a></li>
<li class="last"><a href="/bookmarks" data-ec-bookmark-click="20|9408299|bookmark &gt; more from the economist &gt; BM saved" data-ec-omniture="masthead|act_prod|bookmarks">Timekeeper reading list</a></li>
</ul></li>
<li class="even"><span>My Subscription</span><ul><li class="first"><a href="/products/subscribe">Subscribe to The Economist</a></li>
<li><a href="/activate">Activate my digital subscription</a></li>
<li><a href="/user">Manage my subscription</a></li>
<li class="last"><a href="/products/renew">Renew</a></li>
</ul></li>
<li class="masthead-user"><a href="/user/login?destination=node%2F21555491" class="show-login">Log in or register</a></li>
<li class="masthead-subscribe even last"><a href="https://subscriptions.economist.com/GLB/MAST/T1" class="show-subscribe">Subscribe</a></li>
</ul>
          <div class="mh-search">
        <form action="http://google.com/search"  accept-charset="UTF-8" method="GET" id="search-theme-form">
<div><div id="search" class="container-inline">
  <div class="form-item clearfix" id="edit-search-theme-form-1-wrapper">
 <label for="edit-search-theme-form-1">Search this site:</label>
<input type="text" maxlength="128" name="query" id="edit-search-theme-form-1" size="15" value="" title="Enter the terms you wish to search for." autocorrect="off" class="form-text search-field" />
</div>
<input type="submit" name="op" id="edit-submit" value="Search"  class="form-submit" />
<input type="hidden" name="form_id" id="edit-search-theme-form" value="search_theme_form"  />
<input type="hidden" name="sitesearch" id="edit-sitesearch" value="economist.com"  />
</div>

</div></form>
      </div>
      </div> <!-- /.mh-stripe-wrap -->
</div> <!-- /.mh-stripe -->  <div class="mh-nav mh-big">
    <div class="mh-nav-wrap">
              <h1 class="svg-logo"><a href="/" class="active"><img class="mh-logo" width="170" height="85" src="//cdn.static-economist.com/sites/all/themes/econfinal/images/svg/logo.svg" alt="The Economist" /></a></h1>
                    <nav>
          <ul class="mh-nav-links"><li class="first"><a href="/content/politics-this-week" title="" class="sub-menu-link">World politics</a><ul class="mh-subnav"><li class="first"><a href="/content/politics-this-week">Politics this week</a></li>
<li class="even"><a href="/sections/united-states">United States</a></li>
<li><a href="/sections/britain">Britain</a></li>
<li class="even"><a href="/sections/europe">Europe</a></li>
<li><a href="/sections/china">China</a></li>
<li class="even"><a href="/sections/asia">Asia</a></li>
<li><a href="/sections/americas">Americas</a></li>
<li class="even"><a href="/sections/middle-east-africa">Middle East &amp; Africa</a></li>
<li class="last"><a href="/sections/international">International</a></li>
</ul></li>
<li class="even"><a href="/sections/business-finance" class="sub-menu-link">Business &amp; finance</a><ul class="mh-subnav"><li class="first"><a href="/sections/business-finance">All Business &amp; finance</a></li>
<li class="even last"><a href="/whichmba">Which MBA?</a></li>
</ul></li>
<li class=""><a href="/sections/economics" class="sub-menu-link">Economics</a><ul class="mh-subnav"><li class="first"><a href="/sections/economics">All Economics</a></li>
<li class="even"><a href="/economics-a-to-z">Economics A-Z</a></li>
<li><a href="/markets-data">Markets &amp; data</a></li>
<li class="even last"><a href="/indicators">Indicators</a></li>
</ul></li>
<li class="even"><a href="/sections/science-technology" class="sub-menu-link">Science &amp; technology</a><ul class="mh-subnav"><li class="first"><a href="/sections/science-technology">All Science &amp; technology</a></li>
<li class="even last"><a href="/technology-quarterly" title="Technology Quarterly">Technology Quarterly</a></li>
</ul></li>
<li class=""><a href="/sections/culture" class="sub-menu-link">Culture</a><ul class="mh-subnav"><li class="first"><a href="/sections/culture">All Culture</a></li>
<li class="even"><a href="http://moreintelligentlife.com/">More Intelligent Life</a></li>
<li><a href="/styleguide/introduction">Style guide</a></li>
<li class="even last"><a href="/economist-quiz">The Economist Quiz</a></li>
</ul></li>
<li class="even"><a href="/blogs" class="sub-menu-link">Blogs</a><ul class="mh-subnav"><li class="first"><a href="/blogs">Latest updates</a></li>
<li class="even"><a href="/blogs/buttonwood" title="Financial markets">Buttonwood&#039;s notebook</a></li>
<li><a href="/blogs/democracyinamerica" title="American politics">Democracy in America</a></li>
<li class="even"><a href="/blogs/erasmus">Erasmus</a></li>
<li><a href="/blogs/freeexchange" title="Economics">Free exchange</a></li>
<li class="even"><a href="/blogs/gametheory" title="Sports">Game theory</a></li>
<li><a href="/blogs/graphicdetail" title="Charts, maps and infographics">Graphic detail</a></li>
<li class="even"><a href="/blogs/gulliver" title="Business travel">Gulliver</a></li>
<li><a href="/blogs/prospero" title="Books, arts and culture">Prospero</a></li>
<li class="even last"><a href="/blogs/economist-explains">The Economist explains</a></li>
</ul></li>
<li class=""><a href="http://www.economist.com/debate" class="sub-menu-link">Debate</a><ul class="mh-subnav"><li class="first"><a href="http://www.economist.com/debate">Economist debates</a></li>
<li class="even last"><a href="/content/letters-to-the-editor" title="">Letters to the editor</a></li>
</ul></li>
<li class="even"><a href="/multimedia" class="sub-menu-link">Multimedia</a><ul class="mh-subnav"><li class="first"><a href="/films">Economist Films</a></li>
<li class="even"><a href="http://radio.economist.com">Economist Radio</a></li>
<li><a href="/multimedia">Multimedia library</a></li>
<li class="even last"><a href="/audio-edition">The Economist in audio</a></li>
</ul></li>
<li class="last"><a href="/printedition" title="" class="sub-menu-link">Print edition</a><ul class="mh-subnav"><li class="first"><a href="/printedition/">Current issue</a></li>
<li class="even"><a href="/printedition/covers">Previous issues</a></li>
<li><a href="/printedition/specialreports">Special reports</a></li>
<li class="even"><a href="/content/politics-this-week">Politics this week</a></li>
<li><a href="/content/business-this-week">Business this week</a></li>
<li class="even"><a href="/sections/leaders">Leaders</a></li>
<li><a href="/printedition/kallery">KAL&#039;s cartoon</a></li>
<li class="even last"><a href="/sections/obituary">Obituaries</a></li>
</ul></li>
</ul>                  </nav>
          </div>
  </div> <!-- /.mh-nav -->
</header> <!-- /header -->
  <div id="page" class="container">
    <a name="top" id="navigation-top"></a>

            
    <div id="columns" class="clearfix">
                        <div id="leadspot" class="grid-16 clearfix">
        <div id="block-ec_homepage-ec_homepage_superhero" class="block block-ec_homepage 
">
    <div class="content clearfix">
        <div id="superhero" class="clearfix">
      <div class="hero-superhero"><ul id="hero" class="hero-multiple"><li class="selected"><div class="hero-item hero-item-1"><div class="hero-comment"><a href="/node/21662544/comments#comments" title="Comments" class="comment-icon"><span>5</span></a></div><a href="/news/leaders/21662544-fear-about-chinas-economy-can-be-overdone-investors-are-right-be-nervous-great-fall" class="hero-tab"><h2 class="fly-title">Financial markets</h2><p class="headline">The Great Fall of China</p></a></div><div class="hero-media"><a href="/news/leaders/21662544-fear-about-chinas-economy-can-be-overdone-investors-are-right-be-nervous-great-fall" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/superhero/20150829_LDP001_473.jpg" alt="The Great Fall of China" title="The Great Fall of China"  class="imagecache imagecache-superhero" width="473" height="266" /></a></div></li><li><div class="hero-item hero-item-2"><a href="/news/finance-and-economics/21662584-ukraines-deal-its-creditors-less-impressive-it-appears-tinkering" class="hero-tab"><h2 class="fly-title">Ukraine’s debt restructuring</h2><p class="headline">Tinkering at the edges</p></a></div><div class="hero-media" style="visibility:hidden;opacity:0"><a href="/news/finance-and-economics/21662584-ukraines-deal-its-creditors-less-impressive-it-appears-tinkering" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/superhero/20150829_FNP503_473.jpg" alt="Tinkering at the edges" title="Tinkering at the edges"  class="imagecache imagecache-superhero" width="473" height="266" /></a></div></li><li><div class="hero-item hero-item-3"><a href="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft" class="hero-tab"><h2 class="fly-title">Flight safety</h2><p class="headline">Perils of private planes</p></a></div><div class="hero-media" style="visibility:hidden;opacity:0"><a href="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/original-size/20150829_WOC584_473.png" alt="Perils of private planes" title="Perils of private planes"  class="imagecache imagecache-original-size" width="946" height="532" /></a></div></li><li><div class="hero-item hero-item-4"><div class="hero-comment"><a href="/blogs/economist-explains/2015/08/economist-explains-20#comments" title="Comments" class="comment-icon"><span>6</span></a></div><a href="/blogs/economist-explains/2015/08/economist-explains-20" class="hero-tab"><h2 class="fly-title">Black-hole theory</h2><p class="headline">Hawking a new idea</p></a></div><div class="hero-media" style="visibility:hidden;opacity:0"><a href="/blogs/economist-explains/2015/08/economist-explains-20" class="hero-image" rel="nofollow"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/superhero/20150829_BLP509_473.jpg" alt="Hawking a new idea" title="Hawking a new idea"  class="imagecache imagecache-superhero" width="473" height="266" /></a></div></li></ul></div>
      <div class="cover-image-container"><a href="/printedition" class="cover-image" data-ec-omniture="home|touts|issue|cover"></a><ul><li><a href="/printedition" data-ec-omniture="home|touts|issue|printedition">Full contents</a></li><li class="last"><a href="/products/subscribe" data-ec-omniture="home|touts|issue|subscribe">Subscribe</a></li></ul></div>
    </div>  </div>
</div>

      </div>
      
      <div id="column-content" class="grid-10 grid-first clearfix">
                                <!-- Create left column on search pages -->
                                                  <!-- DoubleClick Floodlight Tag: Please do not remove -->
<script type="text/javascript">
document.write('<iframe id="floodlight" src="" width="1" height="1" frameborder="0" style="display:none"></iframe>');
</script>
<!-- End of DoubleClick Floodlight Tag: Please do not remove -->
<!-- Homepage node -->
<div class="grid-7 grid-first push-3">
  <div id="homepage-center-inner">
    <section class="news-package typog-package">

  <h1 class="fly-title">Shia, not shale</h1>

    <article>
  <a href="/news/finance-economics/21662570-kingdom-can-stand-more-pain-it-will-take-much-cheaper-oil-saudi-arabia-take-action" >
          <div>
          <h2 class="headline">It will take much cheaper oil for Saudi Arabia to take action</h2>
      <p class="rubric">
        The kingdom can stand more pain        <span data-href-redirect="/node/21662570/comments#comments" class="comment-icon"><span>0</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/world/middle-east" title="Middle East &amp;amp; Africa">More in Middle East &amp; Africa &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">Murder and politics in Northern Ireland</h1>

    <article>
  <a href="/news/britain/21662500-murder-former-ira-man-causes-political-tremors-consequences-killing" >
          <div>
          <h2 class="headline">The consequences of a killing</h2>
      <p class="rubric">
        The murder of a former IRA man causes political tremors        <span data-href-redirect="/node/21662500/comments#comments" class="comment-icon"><span>1</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/britain" title="Britain">More in Britain &raquo;</a></div>
  </section><section class="ec-homepage-player news-package typog-package">
      <h1 class="fly-title">Latest audio and video</h1>
    <div style="height:422px;"><object id="ec-homepage-video" class="BrightcoveExperience"><param name="bgcolor" value="#FFFFFF" /><param name="isUI" value="true" /><param name="isVid" value="true" /><param name="dynamicStreaming" value="true" /><param name="autoStart" value="false" /><param name="wmode" value="opaque" /><param name="includeAPI" value="true" /><param name="linkBaseURL" value="http://www.economist.com/multimedia" /><param name="playerID" value="1545427201001" /><param name="playerKey" value="AQ~~,AAABDH-R__E~,dB4S9tmhdOrAcjB6eqWZCo1XXp-OU2vB" /><param name="width" value="402" /><param name="height" value="422" /><param name="templateLoadHandler" value="ECOmnitureBrightCoveHandlersf2f697c8f0d23fecb2b4ce9cc3be8e5d.omnitureBCTemplateLoaded" /><param name="@videoPlayer" value="" /><param name="labels" value="http://cdn.static-economist.com/sites/all/modules/custom/ec_brightcove/EcBcLables.xml" /></object></div>
  <p class="package-more">
    <a href="/multimedia">More in Multimedia &raquo;</a>  </p>
</section><section class="news-package typog-package">

  <h1 class="fly-title">Economics</h1>

    <article>
  <a href="/blogs/freeexchange/2015/08/american-economy" >
          <div>
          <h2 class="headline">The American economy</h2>
      <p class="rubric">
        How exposed are American households to the stock market?        <span data-href-redirect="/blogs/freeexchange/2015/08/american-economy#comments" class="comment-icon"><span>20</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/business-finance" title="Business &amp;amp; finance">More in Business &amp; finance &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">European mobile telecoms</h1>

    <article>
  <a href="/news/business/21661660-eus-new-competition-chief-will-have-rule-wave-mergers-together-we-stand" >
          <div>
          <h2 class="headline">Together we stand</h2>
      <p class="rubric">
        The EU’s new competition chief will have to rule on a wave of mergers        <span data-href-redirect="/node/21661660/comments#comments" class="comment-icon"><span>6</span></span>      </p>
    </div>
      </a>
</article>
  
    <ul><li class="first last"><div class="">
  <a href="https://espresso.economist.com/385877ed6e1207dc4a965ffe024e7862" class="headline">Shedding some light on Europe’s most cut-throat mobile market</a>  </div>
</li>
</ul>  
    <div class="package-more"><a href="/sections/business-finance" title="Business &amp;amp; finance">More in Business &amp; finance &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">Singapore’s democracy</h1>

    <article>
  <a href="/news/asia/21662410-fifty-years-singapores-ruling-party-looks-secure-unequal-contest" >
          <div>
          <h2 class="headline">The first election since the death of Lee Kuan Yew</h2>
      <p class="rubric">
        Fifty years on, Singapore’s ruling party looks secure        <span data-href-redirect="/node/21662410/comments#comments" class="comment-icon"><span>23</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/asia" title="Asia">More in Asia &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">China&#039;s economy</h1>

    <article>
  <a href="/blogs/freeexchange/2015/08/chinas-stockmarket" >
          <div>
          <h2 class="headline">The government giveth and taketh away</h2>
      <p class="rubric">
        China reduces its direct interventions in wobbly markets, but the central bank aims a monetary boost at the broader economy        <span data-href-redirect="/blogs/freeexchange/2015/08/chinas-stockmarket#comments" class="comment-icon"><span>66</span></span>      </p>
    </div>
      </a>
</article>
  
    <ul><li class="first"><div class="">
  <a href="https://espresso.economist.com/376187c6d721478403b2d6a3f4aedd05?utm_content=buffer355f4&amp;utm_medium=social&amp;utm_source=twitter.com&amp;utm_campaign=buffer" class="headline">Spreading the pain: where a Chinese slowdown really hurts</a>  </div>
</li>
<li class="even last"><div class="">
  <a href="/blogs/graphicdetail/2015/08/daily-chart-9" class="headline">Daily chart: The gravity of China’s great fall</a>  <a href="/blogs/graphicdetail/2015/08/daily-chart-9#comments" title="Comments" class="comment-icon"><span>4</span></a></div>
</li>
</ul>  
    <div class="package-more"><a href="/sections/china" title="China">More in China &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">American oil</h1>

    <article>
  <a href="/news/finance-and-economics/21661673-long-overdue-easing-protectionist-export-ban-nafta-naphtha" >
          <div>
          <h2 class="headline">Nafta naphtha</h2>
      <p class="rubric">
        A long-overdue easing of a protectionist export ban        <span data-href-redirect="/node/21661673/comments#comments" class="comment-icon"><span>7</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/sections/business-finance" title="Business &amp;amp; finance">More in Business &amp; finance &raquo;</a></div>
  </section><section class="news-package typog-package">

  <h1 class="fly-title">Africa&#039;s worst war</h1>

    <article>
  <a href="/news/middle-east-and-africa/21662478-disagreeing-agree-south-sudan-agrees-peace-deal-unlikely-last" >
          <div>
          <h2 class="headline">Disagreeing to agree</h2>
      <p class="rubric">
        South Sudan signs a peace deal that is unlikely to last        <span data-href-redirect="/node/21662478/comments#comments" class="comment-icon"><span>1</span></span>      </p>
    </div>
      </a>
</article>
  
  
    <div class="package-more"><a href="/world/middle-east" title="Middle East &amp;amp; Africa">More in Middle East &amp; Africa &raquo;</a></div>
  </section>      </div>
</div>

<div class="side-box multiple-box grey-palette typog-highlights grid-3 pull-7">
  <div id="homepage-highlight-1"><article >
  <a href="/news/britain/21662591-net-migration-britain-has-never-been-higher-immigration-breaks-record">
    <div>
                    <h1 class="fly-title">Immigration breaks a record</h1>
                          <p class="rubric">David Cameron promised to reduce net migration “from the hundreds of thousands to the tens of thousands”. Thankfully, he failed.
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_brc577.png" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="136" />      </a>
</article></div><div id="homepage-highlight-2"><article >
  <a href="/blogs/prospero/2015/08/british-theatre">
    <div>
                    <h1 class="fly-title">Adrian Noble</h1>
                          <p class="rubric">A chat about Oscar Wilde and British theatre
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_BKP503_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-3"><article >
  <a href="/blogs/erasmus/2015/08/europes-religious-war">
    <div>
                    <h1 class="fly-title">Religion and debt</h1>
                          <p class="rubric">A French minister is the latest to explain the euro-crisis in religious terms
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_BLP508_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-4"><article >
  <a href="/blogs/democracyinamerica/2015/08/down-syndrome">
    <div>
                    <h1 class="fly-title">Down syndrome</h1>
                          <p class="rubric">A controversial abortion bill puts moderate Republican presidential candidate John Kasich in an awkward spot
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_USP501_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-5"><article >
  <a href="/news/science-and-technology/21662365-scientists-are-developing-jab-might-only-need-be-given-once-lifetime-why-universal">
    <div>
                    <h1 class="fly-title">The &quot;universal&quot; flu vaccine</h1>
                          <p class="rubric">Scientists are developing a jab that might only need to be given once in a lifetime
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150829_BLP507_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-6"><article >
  <a href="/news/united-states/21661815-obscure-dispute-about-rubbish-could-reshape-agency-working-whos-boss">
    <div>
                    <h1 class="fly-title">Companies and employment</h1>
                          <p class="rubric">An obscure dispute about rubbish could reshape agency working
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150822_USD001_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div><div id="homepage-highlight-7"><article >
  <a href="http://www.economist.com/sciencebriefs">
    <div>
                    <h1 class="fly-title">Science briefs</h1>
                          <p class="rubric">The missing 95% of the universe
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/sci_320.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="119" />      </a>
</article></div><div id="homepage-highlight-8"><article >
  <a href="/news/middle-east-and-africa/21661826-costly-valuable-lessons-guerrilla-army-once-fought">
    <div>
                    <h1 class="fly-title">Hizbullah’s learning curve</h1>
                          <p class="rubric">Costly but valuable lessons for a guerrilla army that once fought in the shadows
</p>
                </div>
          <img src="http://cdn.static-economist.com/sites/default/files/imagecache/homepage_highlight/20150822_MAP003_168.jpg" alt="" title=""  class="imagecache imagecache-homepage_highlight" width="168" height="95" />      </a>
</article></div></div>
<!-- End homepage node -->              </div> <!-- /#main-area -->

              <div id="column-right" class="grid-6 clearfix">
          <div id="homepage-touts">
  <div class="grid-3 grid-first">
    <div class="tout tout-1">
          </div>

    <div class="tout tout-2">
          </div>
  </div> <!-- /#homepage-touts -->

  <div class="grid-3">
    <div class="tout tout-issue">
      <div class="cover-image-container">
              </div>
    </div>
  </div> <!-- /#homepage-issue -->
</div>
          <div id="block-ec_ads-mpu_no_wrapper_top_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty ec-ads-remove-wrapper"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_mpu_no_wrapper_top_ad" data-cb-ad-id="Mpu no wrapper top ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_mpu_no_wrapper_top_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_no_wrapper_top&sz=20x20&tile=6&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_no_wrapper_top&sz=20x20&tile=6&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>

<div id="block-ec_ads-top_mpu_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty"><p class="ec-ads-label">Advertisement</p><!-- Site: Web.  Zone: Home |  --> <div id="gpt_top_mpu_ad" data-cb-ad-id="Top mpu ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_top_mpu_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_top&sz=300x250|350x900|350x600|300x600&tile=7&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_top&sz=300x250|350x900|350x600|300x600&tile=7&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>

<div id="block-ec_ads-ribbon_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec_topic_ribbon ec-ads-remove-if-empty"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_ribbon_ad" data-cb-ad-id="Ribbon ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_ribbon_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dabsribbon&sz=351x49&tile=8&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dabsribbon&sz=351x49&tile=8&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>

<div id="block-ec_social-right_rail_social_share_buttons" class="block block-ec_social 
">
    <div class="content clearfix">
    <div id="social-share-buttons-block">
  <div class="title">Follow <cite>The Economist</cite></div>
  <div class="social-share-buttons">
    <ul class="clearfix">
      <li class="facebook">
        <a data-ec-omniture="rightrail|social_share|facebook" href="http://www.facebook.com/TheEconomist" title="Facebook" target="_blank">Facebook</a>
      </li>
      <li class="twitter">
        <a data-ec-omniture="rightrail|social_share|twitter" href="http://twitter.com/TheEconomist" title="Twitter" target="_blank">Twitter</a>
      </li>
      <li class="linked-in">
        <a data-ec-omniture="rightrail|social_share|linked-in" href="http://www.linkedin.com/groups/Economist-official-group-Economist-newspaper-3056216" title="Linked in" target="_blank">Linked in</a>
      </li>
      <li class="google-plus">
        <a data-ec-omniture="rightrail|social_share|plusone" href="https://plus.google.com/100470681032489535736/posts" title="Google plus" target="_blank">Google plus</a>
      </li>
      <li class="tumblr">
        <a data-ec-omniture="rightrail|social_share|tumblr" href="http://theeconomist.tumblr.com/" title="Tumblr" target="_blank">Tumblr</a>
      </li>
      <li class="instagram">
        <a data-ec-omniture="rightrail|social_share|instagram" href="http://instagram.com/theeconomist/" title="Instagram" target="_blank">Instagram</a>
      </li>
      <li class="youtube">
        <a data-ec-omniture="rightrail|social_share|youtube" href="http://www.youtube.com/user/economistmagazine" title="YouTube" target="_blank">YouTube</a>
      </li>
      <li class="rss">
        <a data-ec-omniture="rightrail|social_share|rss" href="/rss" title="RSS" target="_blank">RSS</a>
      </li>
      <li class="newsletters">
        <a data-ec-omniture="rightrail|social_share|newsletters" href="/newsletters" title="Newsletters" target="_blank">Newsletters</a>
      </li>
    </ul>
  </div>
</div>  </div>
</div>

<div id="block-ec_blogs-ec_blogs_block_recent" class="block block-ec_blogs 
">
    <div class="content clearfix">
    <div class="title">
                 <h6><a href="/latest-updates">Latest updates &raquo;</a></h6>
                 </div><div id="latest-updates"><article id="node-21662570" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/news/finance-economics/21662570-kingdom-can-stand-more-pain-it-will-take-much-cheaper-oil-saudi-arabia-take-action">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/articles/main/20150829_blp511.jpg" alt="Shia, not shale" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/news/finance-economics/21662570-kingdom-can-stand-more-pain-it-will-take-much-cheaper-oil-saudi-arabia-take-action"><span class="latest-updates-fly-title">Shia, not shale</span>: It will take much cheaper oil for Saudi Arabia to take...</a></p>
  <p class="dateline">
    <span class="section">Finance & Economics</span>
    <span class="timestamp" title="2015-08-27T13:21:14+00:00"> 37 mins ago</span>
  </p>
</article>

<article id="node-21662583" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/blogs/graphic-detail/20150829_woc583.png" alt="American and British flight safety: airlines v light aircraft" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/american-and-british-flight-safety-airlines-v-light-aircraft"><span class="latest-updates-fly-title">American and British flight safety: airlines v light aircraft</span>: The perils...</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T12:46:21+00:00"> 1 hrs 22 mins ago</span>
  </p>
</article>

<article id="node-21657486" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/daily-dispatches-0">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/blogs/graphic-detail/china.png" alt="Daily dispatches" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/daily-dispatches-0"><span class="latest-updates-fly-title">Daily dispatches</span>: China crisis</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T12:12:48+00:00"> 1 hrs 55 mins ago</span>
  </p>
</article>

<article id="node-21662580" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/news/markets-and-data/21662580-retail-sales-producer-prices-wages-and-exchange-rates">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/articles/main/20150829_int600.png" alt="" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/news/markets-and-data/21662580-retail-sales-producer-prices-wages-and-exchange-rates"><span class="latest-updates-fly-title"></span>: Retail sales, producer prices, wages and exchange rates</a></p>
  <p class="dateline">
    <span class="section">Markets and data</span>
    <span class="timestamp" title="2015-08-27T10:15:56+00:00"> 3 hrs 47 mins ago</span>
  </p>
</article>

<article id="node-21662579" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/news/markets-and-data/21662579-foreign-reserves">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/articles/main/20150829_int500.png" alt="" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/news/markets-and-data/21662579-foreign-reserves"><span class="latest-updates-fly-title"></span>: Foreign reserves</a></p>
  <p class="dateline">
    <span class="section">Markets and data</span>
    <span class="timestamp" title="2015-08-27T10:15:09+00:00"> 3 hrs 34 mins ago</span>
  </p>
</article>

<article id="node-21624322" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/ebola-graphics">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/07/blogs/graphic-detail/20150829_wom999.png" alt="Ebola in graphics" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/ebola-graphics"><span class="latest-updates-fly-title">Ebola in graphics</span>: The toll of a tragedy</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T09:23:55+00:00"> August 27th, 9:23</span>
  </p>
</article>

<article id="node-21662502" class="ec-sections-latest-updates-block clearfix" data-href-redirect="/blogs/graphicdetail/2015/08/daily-chart-14">
      <img src="http://cdn.static-economist.com/sites/default/files/imagecache/50_by_50/images/2015/08/blogs/graphic-detail/20150829_woc578_0.png" alt="Daily chart" title=""  class="imagecache imagecache-50_by_50" width="50" height="50" />    <p><a href="/blogs/graphicdetail/2015/08/daily-chart-14"><span class="latest-updates-fly-title">Daily chart</span>: Who wants to live forever?</a></p>
  <p class="dateline">
    <span class="section">Graphic detail</span>
    <span class="timestamp" title="2015-08-27T00:01:10+00:00"> August 27th, 0:01</span>
  </p>
</article>

</div><div class="more-latest-updates"><a href="/latest-updates" class="more">More latest updates &raquo;</a></div>  </div>
</div>

<div id="block-ec_mostx-mostpopular" class="block block-ec_mostx 
">
    <div class="content clearfix">
    <div id="most-lists" class="block">
  <div id="block-title"><p>Most commented</p></div>
  <div class="list-wrapper">
          <ul id="commented-list" class="show">
        <li class="mostx-first"><a href="/news/europe/21661941-wanting-burden-shared-germany-eu-country-which-takes-most-asylum-seekers-straining?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21661941"><img src="http://cdn.static-economist.com/sites/default/files/imagecache/mostx_block/images/2015/08/articles/main/20150822_eup503.jpg" alt="Germany, the EU country which takes the most asylum seekers, is straining" title=""  class="mostx-image-first" width="168" height="95" /><span class="mostx-list">1</span><span class="mostx-text"><span class="mostx-fly-title">Refugees in Europe</span>Germany, the EU country which takes the most asylum seekers, is straining</span></a></li><li><a href="/news/europe/21662019-it-may-have-been-frances-latest-islamist-attack-time-no-one-was-killed-attempted-murder?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21662019"><span class="mostx-list">2</span><span class="mostx-text"><span class="mostx-fly-title">Terrorism in France</span>: Attempted murder on the Paris express</span></a></li><li><a href="/news/europe/21661810-journey-capital-hinterland-shows-how-grim-life-has-become-and-how-russians?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21661810"><span class="mostx-list">3</span><span class="mostx-text"><span class="mostx-fly-title">Russia’s economy</span>: The path to penury</span></a></li><li><a href="/news/international/21661812-islamic-states-revival-slavery-extreme-though-it-finds-disquieting-echoes-across?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21661812"><span class="mostx-list">4</span><span class="mostx-text"><span class="mostx-fly-title">Islam and slavery</span>: The persistence of history</span></a></li><li><a href="/news/business-and-finance/21662092-china-sneezing-rest-world-rightly-nervous-causes-and-consequences-chinas?spc=scode&amp;spv=xm&amp;ah=9d7f7ab945510a56fa6d37c30b6f1709" id="track-commented-21662092"><span class="mostx-list">5</span><span class="mostx-text"><span class="mostx-fly-title">Market turmoil</span>: The causes and consequences of China's market crash</span></a></li>      </ul>
    
  </div>
</div>
  </div>
</div>

<div id="block-ec_ads-bottom_mpu_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty"><p class="ec-ads-label">Advertisement</p><!-- Site: Web.  Zone: Home |  --> <div id="gpt_bottom_mpu_ad" data-cb-ad-id="Bottom mpu ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_bottom_mpu_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_bottom&sz=300x250|50x50|300x600&tile=9&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_bottom&sz=300x250|50x50|300x600&tile=9&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>

<div id="block-ec_pixel_tracking_onscroll-bottom_onscroll" class="block block-ec_pixel_tracking_onscroll 
">
    <div class="content clearfix">
    <div id="onscroll-ad-holder-mpu2"></div>  </div>
</div>

<div id="block-block-1" class="block block-block 
">
    <div class="content clearfix">
    <div id="product-events">
<div class="title">Products and events</div>
<div class="product-first linked">
<p class="products-events-section"><strong><a  class="social-link" href="/economist-quiz" target="_blank">Test your EQ</a></strong><br />
Take our weekly news quiz to stay on top of the headlines
</p>
</div>
<div class="product-section product-section-last linked">
<p class="products-events-section"><strong><a  class="social-link" href="http://econ.st/R7pQMy" target="_blank">Want more from <i>The Economist?</i></a></strong><br />
Visit The Economist e-store and you’ll find a range of carefully selected products for business and pleasure, Economist books and diaries, and much more</p>
</div>
</div>
  </div>
</div>

<div id="block-ec_ads-bottom_right_mpu_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads ec-ads-remove-if-empty"><p class="ec-ads-label">Advertisement</p><!-- Site: Web.  Zone: Home |  --> <div id="gpt_bottom_right_mpu_ad" data-cb-ad-id="Bottom right mpu ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_bottom_right_mpu_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_bottom_right&sz=300x250|45x45&tile=10&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dmpu_bottom_right&sz=300x250|45x45&tile=10&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>

<div id="block-ec_pixel_tracking_onscroll-bottom_right_onscroll" class="block block-ec_pixel_tracking_onscroll 
">
    <div class="content clearfix">
    <div id="onscroll-ad-holder"></div>  </div>
</div>

        </div> <!-- /#side-bar -->
      
      
    </div> <!-- /#columns -->
          <div id="footer-classifieds" class="clearfix">
        <div class="title">Classified ads</div>
        <div id="block-ec_ads-button_ad_1" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified first">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_1" data-cb-ad-id="Button ad 1"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_button_ad_1");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton1&sz=125x125&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton1&sz=125x125&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

<div id="block-ec_ads-button_ad_2" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_2" data-cb-ad-id="Button ad 2"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_button_ad_2");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton2&sz=125x125&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton2&sz=125x125&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

<div id="block-ec_ads-button_ad_3" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_3" data-cb-ad-id="Button ad 3"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_button_ad_3");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton3&sz=125x125&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton3&sz=125x125&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

<div id="block-ec_ads-button_ad_4" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_4" data-cb-ad-id="Button ad 4"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_button_ad_4");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton4&sz=125x125&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton4&sz=125x125&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

<div id="block-ec_ads-button_ad_5" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_5" data-cb-ad-id="Button ad 5"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_button_ad_5");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton5&sz=125x125&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton5&sz=125x125&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

<div id="block-ec_ads-button_ad_6" class="block block-ec_ads 
ec-ads-gpt ec-classified-box ec-ads-classified last">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_button_ad_6" data-cb-ad-id="Button ad 6"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_button_ad_6");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton6&sz=125x125&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dbutton6&sz=125x125&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

      </div>
        <aside class="site-index">
      <div>
  <div class="svg-logo"><img class="mh-logo" width="170" height="85" src="//cdn.static-economist.com/sites/all/themes/econfinal/images/svg/logo.svg" alt="The Economist" /></div>
      <ul class="site-index-1">
              <li><a href="/contact-info" data-ec-omniture="mini_map_home|contact_us">Contact us</a></li>
              <li><a href="/help/home" data-ec-omniture="mini_map_home|help">Help</a></li>
              <li><a href="/user" data-ec-omniture="mini_map_home|my_account">My account</a></li>
              <li><a href="/products/subscribe" data-ec-omniture="mini_map_home|subscribe">Subscribe</a></li>
              <li><a href="/printedition" data-ec-omniture="mini_map_home|print_edition">Print edition</a></li>
              <li><a href="/digital" data-ec-omniture="mini_map_home|digital_editions">Digital editions</a></li>
              <li><a href="/events-conferences" data-ec-omniture="mini_map_home|events">Events</a></li>
              <li><a href="http://jobs.economist.com/" data-ec-omniture="mini_map_home|jobs_economist_com">Jobs.Economist.com</a></li>
              <li><a href="/bookmarks" data-ec-omniture="mini_map_home|timekeeper_saved_articles">Timekeeper saved articles</a></li>
          </ul>
  </div>
<div>
  <h6>Sections</h6>
  <ul class="footer-index-site-2-content">
                <li><a href="/sections/united-states" data-ec-omniture="mini_map_home|united_states">United States</a></li>
                <li><a href="/sections/britain" data-ec-omniture="mini_map_home|britain">Britain</a></li>
                <li><a href="/sections/europe" data-ec-omniture="mini_map_home|europe">Europe</a></li>
                <li><a href="/sections/china" data-ec-omniture="mini_map_home|china">China</a></li>
                <li><a href="/sections/asia" data-ec-omniture="mini_map_home|asia">Asia</a></li>
                <li><a href="/sections/americas" data-ec-omniture="mini_map_home|americas">Americas</a></li>
                <li><a href="/sections/middle-east-africa" data-ec-omniture="mini_map_home|middle_east_africa">Middle East &amp; Africa</a></li>
                <li><a href="/sections/international" data-ec-omniture="mini_map_home|international">International</a></li>
                <li><a href="/sections/business-finance" data-ec-omniture="mini_map_home|business_finance">Business &amp; finance</a></li>
                <li><a href="/sections/economics" data-ec-omniture="mini_map_home|economics">Economics</a></li>
                <li><a href="/markets-data" data-ec-omniture="mini_map_home|markets_data">Markets &amp; data</a></li>
                <li><a href="/sections/science-technology" data-ec-omniture="mini_map_home|science_technology">Science &amp; technology</a></li>
                <li><a href="http://www.economist.com/specialreports" data-ec-omniture="mini_map_home|special_reports">Special reports</a></li>
                <li><a href="/sections/culture" data-ec-omniture="mini_map_home|culture">Culture</a></li>
                <li><a href="/multimedia" data-ec-omniture="mini_map_home|multimedia_library">Multimedia library</a></li>
            </ul>
      <h6>Debate and discussion</h6>
    <ul>
            <li><a href="http://www.economist.com/debate" data-ec-omniture="mini_map_home|the_economist_debates">The Economist debates</a></li>
                <li><a href="http://www.economist.com/letters" data-ec-omniture="mini_map_home|letters_to_the_editor">Letters to the editor</a></li>
                <li><a href="/economist-quiz" data-ec-omniture="mini_map_home|the_economist_quiz">The Economist Quiz</a></li>
      </ul>
</div>
  <div>
    <h6>Blogs</h6>
    <ul>

              <li><a href="/blogs/buttonwood" data-ec-omniture="mini_map_home|buttonwoods_notebook">Buttonwood&#039;s notebook</a></li>
              <li><a href="/blogs/democracyinamerica" data-ec-omniture="mini_map_home|democracy_in_america">Democracy in America</a></li>
              <li><a href="/blogs/erasmus" data-ec-omniture="mini_map_home|erasmus">Erasmus</a></li>
              <li><a href="/blogs/freeexchange" data-ec-omniture="mini_map_home|free_exchange">Free exchange</a></li>
              <li><a href="/blogs/gametheory" data-ec-omniture="mini_map_home|game_theory">Game theory</a></li>
              <li><a href="/blogs/graphicdetail" data-ec-omniture="mini_map_home|graphic_detail">Graphic detail</a></li>
              <li><a href="/blogs/gulliver" data-ec-omniture="mini_map_home|gulliver">Gulliver</a></li>
              <li><a href="/blogs/prospero" data-ec-omniture="mini_map_home|prospero">Prospero</a></li>
              <li><a href="/blogs/economist-explains" data-ec-omniture="mini_map_home|the_economist_explains">The Economist explains</a></li>
          </ul>
  </div>
  <div>
    <h6>Research and insights</h6>
    <ul>
                              <li><a href="/topics" data-ec-omniture="mini_map_home|topics">Topics</a></li>
                                        <li><a href="/economics-a-to-z" data-ec-omniture="mini_map_home|economics_a_z">Economics A-Z</a></li>
                                        <li><a href="/styleguide/introduction" data-ec-omniture="mini_map_home|style_guide">Style guide</a></li>
                                        <li><a href="http://www.theworldin.com/" data-ec-omniture="mini_map_home|the_world_in_2015">The World in 2015</a></li>
                                        <li><a href="/whichmba" data-ec-omniture="mini_map_home|which_mba_">Which MBA?</a></li>
                                        <li><a href="https://success.economist.com/?fsrc=econfooter" data-ec-omniture="mini_map_home|mba_services">MBA Services</a></li>
                                        <li><a href="https://gmat.economist.com/?gsrc=economist_footer&amp;c3ch=Economist&amp;c3nid=site footer" data-ec-omniture="mini_map_home|the_economist_gmat_tutor">The Economist GMAT Tutor</a></li>
                                        <li><a href="https://execed.economist.com/?fsrc=econ-foot" data-ec-omniture="mini_map_home|executive_education_navigator">Executive Education Navigator</a></li>
                                        <li><a href="http://www.economist.com/rights" data-ec-omniture="mini_map_home|reprints_and_permissions">Reprints and permissions</a></li>
                                        </ul>
              <h6><a href="http://www.economistgroup.com" data-ec-omniture="mini_map_home|the_economist_group">The Economist Group &raquo;</a></h6>
            <ul>
                                        <li><a href="http://www.eiu.com" data-ec-omniture="mini_map_home|the_economist_intelligence_unit">The Economist Intelligence Unit</a></li>
                                        <li><a href="http://store.eiu.com" data-ec-omniture="mini_map_home|the_economist_intelligence_unit_store">The Economist Intelligence Unit Store</a></li>
                                        <li><a href="http://www.corporatenetwork.com" data-ec-omniture="mini_map_home|the_economist_corporate_network">The Economist Corporate Network</a></li>
                                        <li><a href="http://ideaspeoplemedia.com/" data-ec-omniture="mini_map_home|ideas_people_media">Ideas People Media</a></li>
                                        <li><a href="http://www.moreintelligentlife.com" data-ec-omniture="mini_map_home|intelligent_life">Intelligent Life</a></li>
                                        <li><a href="http://www.rollcall.com/?t=0506EC&amp;p=econ&amp;s=econ" data-ec-omniture="mini_map_home|roll_call">Roll Call</a></li>
                                        <li><a href="http://www.cq.com/news.do" data-ec-omniture="mini_map_home|cq">CQ</a></li>
                                        <li><a href="http://www.eurofinance.com" data-ec-omniture="mini_map_home|eurofinance">EuroFinance</a></li>
                                        <li><a href="http://store.economist.com" data-ec-omniture="mini_map_home|the_economist_store">The Economist Store</a></li>
                                        </ul>
              <h6 class="minimap-site-index"><a href="/content/site-index" data-ec-omniture="mini_map_home|view_complete_site_index_">View complete site index »</a></h6>
            <ul>
                      </ul>
  </div>


    </aside>
    <div id="block-ec_ads-slider_ad" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <!-- Site: Web.  Zone: Home |  --> <div id="gpt_slider_ad" data-cb-ad-id="Slider ad"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_slider_ad");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dslider&sz=1x1&tile=11&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dslider&sz=1x1&tile=11&c=328904146" />
  </a>
</noscript></div>  </div>
</div>

<div id="block-ec_ads-adcast" class="block block-ec_ads 
ec-ads-gpt">
    <div class="content clearfix">
    <div class="ec-ads-remove-if-empty"><!-- Site: Web.  Zone: Home |  --> <div id="gpt_adcast" data-cb-ad-id="Adcast"><script type="text/javascript">  googletag.cmd.push(function() {      googletag.display("gpt_adcast");  });
</script><noscript>
  <a href="//pubads.g.doubleclick.net/gampad/jump?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dadcast&sz=250x1000&tile=12&c=328904146" target="_blank">
   <img src="//pubads.g.doubleclick.net/gampad/ad?iu=/5605/teg.fmsq/j2ek&t=subs%3Dn%26sdn%3Dn%26pos%3Dadcast&sz=250x1000&tile=12&c=328904146" />
  </a>
</noscript></div></div>  </div>
</div>


    <footer>
        <div class="footer-stripe-top">
    <ul>
                                <li><a href="/contact-info" data-ec-omniture="footer_home|contact_us">Contact us</a></li>
                          <li><a href="http://www.economist.com/help" data-ec-omniture="footer_home|help">Help</a></li>
                          <li><a href="http://www.economist.com/help/about-us#About_Economistcom" data-ec-omniture="footer_home|about_us">About us</a></li>
                          <li><a href="http://www.economistgroupmedia.com" data-ec-omniture="footer_home|advertise_with_us">Advertise with us</a></li>
                          <li><a href="/mediadirectory" data-ec-omniture="footer_home|editorial_staff">Editorial Staff</a></li>
                          <li><a href="/mediadirectory/books" data-ec-omniture="footer_home|staff_books">Staff Books</a></li>
                          <li><a href="http://www.economistgroupcareers.com" data-ec-omniture="footer_home|careers">Careers</a></li>
                          <li><a href="/content/site-index" data-ec-omniture="footer_home|site_index">Site index</a></li>
                          </ul>
  </div>

<div class="footer-stripe-bottom">
  <ul>
    <li>Copyright &copy; The Economist Newspaper Limited 2015. All rights reserved.</li>
                        <li><a href="/help/accessibilitypolicy" data-ec-omniture="footer_home|accessibility">Accessibility</a></li>
                  <li><a href="http://www.economistgroup.com/results_and_governance/governance/privacy" data-ec-omniture="footer_home|privacy_policy">Privacy policy</a></li>
                  <li><a href="/cookies-info" data-ec-omniture="footer_home|cookies_info">Cookies info</a></li>
                  <li><a href="/legal/terms-of-use" data-ec-omniture="footer_home|terms_of_use">Terms of use</a></li>
                            <li class="cookie-pref"><div id="teconsent"><script type="text/javascript" src="http://consent-st.truste.com/get?name=notice.js&domain=economist.com&c=teconsent"></script></div></li>
        </ul>
</div>
    </footer> <!-- /footer -->
  </div>

  <script type="text/javascript" src="http://admin.brightcove.com/js/BrightcoveExperiences_all.js"></script>
<script type="text/javascript" src="http://admin.brightcove.com/js/APIModules_all.js"></script>
 <!-- Google Code for Economist Home Page Audience Remarketing List -->
<script type="text/javascript">
/* <![CDATA[ */
var google_conversion_id = 1021426701;
var google_conversion_language = "en";
var google_conversion_format = "3";
var google_conversion_color = "666666";
var google_conversion_label = "8uJbCJ_f8gEQjfiG5wM";
var google_conversion_value = 0;
/* ]]> */
</script>
<script type="text/javascript" src="//www.googleadservices.com/pagead/conversion.js">
</script>
<noscript>
<div style="display:inline;">
<img height="1" width="1" style="border-style:none;" alt="" src="//www.googleadservices.com/pagead/conversion/1021426701/?label=8uJbCJ_f8gEQjfiG5wM&amp;guid=ON&amp;script=0"/>
</div>
</noscript>
<!-- Google Code for Visitors (Non-Subscriber) Remarketing List -->
<script type="text/javascript">
/* <![CDATA[ */
var google_conversion_id = 1012729503;
var google_conversion_language = "en";
var google_conversion_format = "3";
var google_conversion_color = "666666";
var google_conversion_label = "LLuqCOn_0QEQn4304gM";
var google_conversion_value = 0;
/* ]]> */
</script>
<script type="text/javascript"
src="//www.googleadservices.com/pagead/conversion.js">
</script>
<noscript>
<div style="display:inline;">
<img height="1" width="1" style="border-style:none;" alt=""
src="//www.googleadservices.com/pagead/conversion/1012729503/?label=LLu
qCOn_0QEQn4304gM&amp;guid=ON&amp;script=0"/>
</div>
</noscript>

<!-- Quantcast Tag -->
<script type="text/javascript">
var _qevents = _qevents || [];

(function() {
var elem = document.createElement('script');
elem.src = (document.location.protocol == "https:" ? "https://secure" : "http://edge") + ".quantserve.com/quant.js";
elem.async = true;
elem.type = "text/javascript";
var scpt = document.getElementsByTagName('script')[0];
scpt.parentNode.insertBefore(elem, scpt);
})();

_qevents.push({
qacct:"p-a8GHW19EK4IzY"
});
</script>

<noscript>
<div style="display:none;">
<img src="//pixel.quantserve.com/pixel/p-a8GHW19EK4IzY.gif" border="0" height="1" width="1" alt="Quantcast"/>
</div>
</noscript>
<!-- End Quantcast tag -->

<!--
    SiteCatalyst code version: H.13    Copyright 1997-2009 Omniture, Inc.
    More info available at http://www.omniture.com
-->
<script language="javascript" type="text/javascript">var s_account = 'economistcomprod';
</script><script language="javascript" src="//cdn.static-economist.com/sites/default/files/external/ec_omniture/3_5/ec_omniture_s_code.js" type="text/javascript"></script>
<script language="javascript" type="text/javascript">
  s.pageName="homepage";
s.pageType="";
s.server="economist.com";
s.channel="home";
s.contextData.subsection="";
s.prop1="homepage";
s.prop2="";
s.prop3="web";
s.prop4="homepage";
s.prop5="";
s.prop6="";
s.prop7="";
s.prop11="not_logged_in";
s.prop12="";
s.prop13="anonymous";
s.prop14="";
s.prop15="";
s.prop26="";
s.prop28="";
s.prop31="2013|08|14";
s.prop34="econfinal";
s.prop40="0";
s.prop41="fmsq";
s.prop42="j2ek";
s.hier1="";
s.state="";
s.zip="";
  // For grape shot.
  if (typeof s !== "undefined" && typeof window.gs_channels !== "undefined") {
    s.prop6 = window.gs_channels;
  }
  if (typeof s !== "undefined" && typeof Economist !== "undefined" && typeof Economist.userData !== "undefined" && typeof Economist.userData.data !== "undefined") {
    s.prop11 = "logged_in";

    if (typeof Economist.userData.data.username !== "undefined" && typeof Economist.analytics !== 'undefined' && typeof Economist.analytics.omnitureTools !== 'undefined') {
      s.prop12 = Economist.analytics.omnitureTools.clean(Economist.userData.data.username);
    }

    if (typeof Economist.userData.data.uid !== "undefined") {
      s.prop40 = Economist.userData.data.uid;
    }

    if (typeof s.pageName !== "undefined" && s.pageName === "pay_barrier|registration_offer") {
      s.pageName = "pay_barrier|subscription_offer";
      s.events = "event2,event23";
    }
  }

  if (typeof s !== "undefined" && typeof Economist !== "undefined" && typeof Economist.userData === "undefined" && $.cookie("ec_old_user_data") !== null) {
    // Omniture is around, but we dont have userData structure in the Economist associative array.
    // Check to see if we have a saved cookie with the last logged in user details within, if so use that.
    try {
      Economist.oldUserData = $.cookie("ec_old_user_data").split('|');
    }
    catch (e) {
      // If we get a error its probably due to the cookie containing invalid data,
      // so unset the oldUserData.
      delete Economist.oldUserData;
    }

    if (typeof Economist.oldUserData !== "undefined" && Economist.oldUserData instanceof Array && typeof Economist.analytics !== 'undefined' && typeof Economist.analytics.omnitureTools !== 'undefined') {
      s.prop12 = Economist.analytics.omnitureTools.clean(Economist.oldUserData[0]);
      s.prop40 = Economist.oldUserData[1];
    }
  }

  // Process any queued events if we are not using sessions.
  if ($.cookie('ec_omniture_queue')) {
    try {
      var processQueue = (function() {
        var queued_events, queued_event, event_name;
        queued_events = JSON.parse($.cookie('ec_omniture_queue'));
        for (queued_event in queued_events) {
          for (event_name in queued_events[queued_event]) {
            if (queued_events[queued_event].hasOwnProperty(event_name)) {
              if (typeof s[event_name] !== "undefined" && s[event_name].length > 0) {
                s[event_name] += ',' + queued_events[queued_event][event_name];
              } else {
                s[event_name] = queued_events[queued_event][event_name];
              }
            }
          }
        }
      }());
    } catch (e) {
      // do nothing but catch the error so it doesn't propagate -
      // we ensure the cookie is deleted in the finally block.
    } finally {
      $.cookie('ec_omniture_queue', '', { path: '/', expires: -5 });
    }
  }
  var s_code=s.t();if(s_code)document.write(s_code)
</script>
<script language="JavaScript" type="text/javascript">
   if(navigator.appVersion.indexOf('MSIE')>=0)document.write(unescape('%3C')+'\!-'+'-');
</script>
<noscript>
  <a href="http://www.omniture.com" title="Web Analytics">
    <img src="/stats.economist.com/b/ss/economistcomprod/1/H.20.3--NS/0" height="1" width="1" border="0" alt="" />
  </a>
</noscript>
<!-- End SiteCatalyst code version: H.13 -->
<script type="text/javascript" src="http://cdn.static-economist.com/sites/default/files/js/js_880c1d9b8688156fe6b261195971ed0d.js"></script>
      <!--
START of Interactive Avenue Pixel Tag: Please do not remove!!
-->
<script type="text/javascript">
  /* <![CDATA[ */
  var google_conversion_id = 1035071974;
  var google_custom_params = window.google_tag_params;
  var google_remarketing_only = true;
  /* ]]> */
</script>
<script type="text/javascript"
        src="//www.googleadservices.com/pagead/conversion.js">
</script>
<noscript>
  <div style="display:inline;">
    <img height="1" width="1" style="border-style:none;" alt=""
         src="//googleads.g.doubleclick.net/pagead/viewthroughconversion/1035071974/?value=0&amp;guid=ON&amp;script=0"/>
  </div>
</noscript>

<script type="text/javascript">
  /* <![CDATA[ */
  var google_conversion_id = 1047392724;
  var google_custom_params = window.google_tag_params;
  var google_remarketing_only = true;
  /* ]]> */
</script>
<script type="text/javascript"
        src="//www.googleadservices.com/pagead/conversion.js">
</script>
<noscript>
  <div style="display:inline;">
    <img height="1" width="1" style="border-style:none;" alt=""
         src="//googleads.g.doubleclick.net/pagead/viewthroughconversion/1047392724/?value=0&amp;guid=ON&amp;script=0"/>
  </div>
</noscript>
<!--
END of Interactive Avenue Pixel Tag: Please do not remove!!
-->
<script type='text/javascript'>
var ebRand = Math.random()+'';
ebRand = ebRand * 1000000;
document.write('<scr'+'ipt src="//bs.serving-sys.com/BurstingPipe/ActivityServer.bs?cn=as&ActivityID=326761&rnd=' + ebRand + '"></scr' + 'ipt>');
</script>
<noscript>
<img width="1" height="1" style="border:0" src="//bs.serving-sys.com/BurstingPipe/ActivityServer.bs?cn=as&amp;ActivityID=326761&;ns=1"/>
</noscript>
<!--
Start of DoubleClick Floodlight Tag: Please do not remove
-->
<script type="text/javascript">
  var axel = Math.random() + "";
  var a = axel * 10000000000000;
  var d = document.createElement('iframe');
  d.setAttribute('height', '1');
  d.setAttribute('width', '1');
  d.setAttribute('style', 'display:none');
  d.setAttribute('src', '//4396156.fls.doubleclick.net/activityi;src=4396156;type=Econo014;cat=Econo0;ord=1;num=' + a + '?');
  document.body.appendChild(d);</script>
<noscript>
  <iframe src="//4396156.fls.doubleclick.net/activityi;src=4396156;type=Econo014;cat=Econo0;ord=1;num=1?" width="1" height="1" frameborder="0" style="display:none"></iframe>
</noscript>
<!-- End of DoubleClick Floodlight Tag: Please do not remove -->
<!--
START of Universal McCann Tag: Please do not remove!!
-->
<script src="//d21j20wsoewvjq.cloudfront.net/tt-8c244b370747c1930a4e0967254778ddbb69f6a409e62beebe5f92191a09a3a1.js" async>
</script>
<!--
END of Universal McCann Tag: Please do not remove!!
-->
<script src="//tag.yieldoptimizer.com/ps/ps?t=s&p=1814&pg=Page"></script>          <script type="text/javascript">
  (function() {
    var em = document.createElement('script'); em.type = 'text/javascript'; em.async = true;
    em.src = ('https:' == document.location.protocol ? 'https://us-ssl' : 'http://us-cdn') + '.effectivemeasure.net/em.js';
    var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(em, s);
  })();
</script>
<noscript>
<img src="//us.effectivemeasure.net/em_image" alt="" style="position:absolute; left:-5px;" />
</noscript>          <!-- Begin BlueKai Tag -->
<script type="text/javascript">
window.bk_async = function() {
	// INSERT DATA HERE IN THE FORM:
	// bk_addPageCtx("<<SOMEKEY>>", "<<SOMEVALUE>>");
	
	
	try{
		if(typeof Econ.user.user.country !== 'undefined'){
			var bkCountry = (typeof Econ.user.user.country.iso != 'undefined' ? Econ.user.user.country.iso : '');
			var bkRegion = (typeof Econ.user.user.country.region != 'undefined' ? Econ.user.user.country.region : '');
			bk_addPageCtx("ctry", bkCountry);
			bk_addPageCtx("reg", bkRegion);
		}
	}	
	catch(e){}
	
	try{	
		if(typeof s !='undefined'){
			var bkChannel = (typeof s.channel != 'undefined' ? s.channel : '');
			var bkPageName = (typeof s.pageName != 'undefined' ? s.pageName : '');

			var bkprop1 = (typeof s.prop1 != 'undefined' ? s.prop1 : '');
			var bkprop3 = (typeof s.prop3 != 'undefined' ? s.prop3 : '');
			var bkprop4 = (typeof s.prop4 != 'undefined' ? s.prop4 : '');
			var bkprop5 = (typeof s.prop5 != 'undefined' ? s.prop5 : '');
			var bkprop6 = (typeof s.prop6 != 'undefined' ? s.prop6 : '');
			var bkprop8 = (typeof s.prop8 != 'undefined' ? s.prop8 : '');
			var bkprop9 = (typeof s.prop9 != 'undefined' ? s.prop9 : '');
			var bkprop10 = (typeof s.prop10 != 'undefined' ? s.prop10 : '');
			var bkprop11 = (typeof s.prop11 != 'undefined' ? s.prop11 : '');
			var bkprop13 = (typeof s.prop13 != 'undefined' ? s.prop13 : '');
			var bkprop16 = (typeof s.prop16 != 'undefined' ? s.prop16 : '');
			var bkprop40 = (typeof s.prop40 != 'undefined' ? s.prop40 : '');
			var bkprop41 = (typeof s.prop41 != 'undefined' ? s.prop41 : '');
			var bkprop42 = (typeof s.prop42 != 'undefined' ? s.prop42 : '');
			var bkprop53 = (typeof s.prop53 != 'undefined' ? s.prop53 : '');
			var bkprop54 = (typeof s.prop54 != 'undefined' ? s.prop54 : '');

			var bkevents = (typeof s.events != 'undefined' ? s.events : '');

			bk_addPageCtx("chan", bkChannel);
			bk_addPageCtx("page", bkPageName);
			bk_addPageCtx("p1", bkprop1);
			bk_addPageCtx("p3", bkprop3);
			bk_addPageCtx("p4", bkprop4);
			bk_addPageCtx("p5", bkprop5);
			bk_addPageCtx("p6", bkprop6);
			bk_addPageCtx("p8", bkprop8);
			bk_addPageCtx("p9", bkprop9);
			bk_addPageCtx("p10", bkprop10);
			bk_addPageCtx("p11", bkprop11);
			bk_addPageCtx("p13", bkprop13);
			bk_addPageCtx("p16", bkprop16);
			bk_addPageCtx("did", bkprop40);
			bk_addPageCtx("p41", bkprop41);
			bk_addPageCtx("p42", bkprop42);
			bk_addPageCtx("p53", bkprop53);
			bk_addPageCtx("p54", bkprop54);
			bk_addPageCtx("events", bkevents);
		}
	}
	catch(e){}	
  bk_doJSTag(18452, 10);

};
(function() {
  var scripts=document.getElementsByTagName('script')[0];
  var s=document.createElement('script');
  s.async = true;
  s.src = '//tags.bkrtx.com/js/bk-coretag.js';
  scripts.parentNode.insertBefore(s, scripts);
}());
</script>
<!-- End BlueKai Tag -->
    </body>
</html>`)
