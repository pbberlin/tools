package gitkit

const UNUSEDwidgetRaw = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<title>Favorite Weekday</title>
<script type="text/javascript" src="//www.gstatic.com/authtoolkit/js/gitkit.js"></script>
<link type="text/css" rel="stylesheet" href="//www.gstatic.com/authtoolkit/css/gitkit.css">
</head>
<body>
<div id="gitkitWidget"></div>
<script type="text/javascript">


  var brandingURL = "https://tec-news.appspot.com/auth/accountChooserBranding.html";
  var faviconURL  = "https://tec-news.appspot.com/favicon.ico";
  var accountChooserConfig = {
              title:    'Sign in to tec-news',
              favicon:  faviconURL,
              branding: brandingURL,
         };

  if (window.location.hostname == "localhost") {
    faviconURL = "";
    brandingURL = "";
    accountChooserConfig = {};
  }


  google.identitytoolkit.start(
    '#gitkitWidget',
    {
      apiKey:           '{{.BrowserAPIKey}}',
      signInSuccessUrl: '{{.SignInSuccessUrl}}',
      signOutUrl:       '{{.SignOutURL}}',
      oobActionUrl:     '{{.OOBActionURL}}',
      signInOptions:    ['google','facebook'],
      siteName:         'tec-news',
      acUiConfig:        accountChooserConfig,      
    },


    '{{.POSTBody}}');
      // acUiConfig - 
      // is rejected because of "invalid domain"
      // favicon:  'http://tec-news.appspot.com/favicon.ico' 
      // favicon:  'http://localhost:8087/favicon.ico' 
      // branding: 'http://tec-news.appspot.com/account_choooser_branding'

  	  // "signInOptions": ["password","google","facebook"]

      // siteName required for "manage account"
</script>
</body>
</html>`

const UNUSEDhomeRaw = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
<title>Favorite Weekday</title>
<script type="text/javascript" src="//www.gstatic.com/authtoolkit/js/gitkit.js"></script>
<link type="text/css" rel="stylesheet" href="//www.gstatic.com/authtoolkit/css/gitkit.css">
</head>
<body>
{{if not .User}}
<script type=text/javascript>
  google.identitytoolkit.setConfig({
    widgetUrl: '{{.WidgetURL}}',
  });
  google.identitytoolkit.signIn();
</script>
{{else}}
<div id="navbar"></div>
<script type=text/javascript>
  google.identitytoolkit.signInButton(
      '#navbar',
      {
        widgetUrl:   '{{.WidgetURL}}',
        signOutUrl:  '{{.SignOutURL}}'
      }
  );
</script>
<br>
<p>Hello {{.User.Name}} ID: {{.User.ID}} Email: {{.User.Email}}</p>
<p>
  14119357422359180555 - peter.buchmann@web.de<br>
  14952300052240127534 - peter.buchmann.68@gmail.com<br>
  13036053464987178369 - peter.buchmann.68@gmail.com<br>
  </p>

<p><a href="https://security.google.com/settings/security/permissions">Remove Google App Permissions</a></p>
<p><a href="https://www.facebook.com/settings?tab=applications">Remove Facebook App Permissions</a></p>

<p>Your favorite weekday:</p>
<form method="POST" action="{{.UpdateWeekdayURL}}">
  <input type="hidden" name="xsrftoken" value="{{.UpdateWeekdayXSRFToken}}">
  <select name="favorite">
    {{range $index, $weekday := .Weekdays}}
      <option value="{{$index}}" {{if eq $.WeekdayIndex $index}}selected{{end}}>{{$weekday}}</option>{{end}}
  </select>
  <button type="submit">change</button>
</form>
<br>
<p>Tired of FavWeekday?</p>
<form method="POST" action="{{.DeleteAccountURL}}">
  <input type="hidden" name="xsrftoken" value="{{.DeleteAccountXSRFToken}}">
  <button type="submit">delete account</button>
</form>
{{end}}
</body>
</html>`
