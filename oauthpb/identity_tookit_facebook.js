<script>
  // full example:
  // https://developers.facebook.com/docs/facebook-login/web

  window.fbAsyncInit = function() {
    FB.init({
      appId      : '942324259171809',
      cookie     : true,  // enable cookies for server access to session
      xfbml      : true,
      version    : 'v2.5'
    });
  };

  (function(d, s, id){
     var js, fjs = d.getElementsByTagName(s)[0];
     if (d.getElementById(id)) {return;}
     js = d.createElement(s); js.id = id;
     js.src = "//connect.facebook.net/en_US/sdk.js";
     fjs.parentNode.insertBefore(js, fjs);
   }(document, 'script', 'facebook-jssdk'));
</script>