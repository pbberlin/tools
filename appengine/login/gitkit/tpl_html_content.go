package gitkit

import (
	"fmt"

	"github.com/pbberlin/tools/net/http/routes"
)

// client side javascript to measure signin status
const signedInCode = `
  <script type=text/javascript>
    function getCookie(name) {
        name = name + "=";
        var ca = document.cookie.split(';');
        for(var i=0; i<ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0)==' ')   c = c.substring(1);
            if (c.indexOf(name) == 0)  return name + " is present";
        }
        return "";
    }
    console.log( "cukis:", getCookie("gtoken") );
  </script>
`

const home1 = `{{if not .User}}
<!-- 
  We could just as easy do automatic forwarding if not signed on:

  <script type=text/javascript>
    google.identitytoolkit.setConfig({  
      widgetUrl: '{{.WidgetURL}}'
    });
    // automatic forwarding to login:
    // google.identitytoolkit.signIn();
  </script>
  <a href="#" onclick="google.identitytoolkit.signIn()">Sign In</a><br><br>

-->


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
    </p>

  <p><a href="https://security.google.com/settings/security/permissions">Remove Google App Permissions</a></p>
  <p><a href="https://www.facebook.com/settings?tab=applications">Remove Facebook App Permissions</a></p>

	{{ .CookieDump }}

{{end}}`

const home2 = `{{if .User}}
  <p>Your favorite weekday:</p>
  <form method="POST" action="{{.UpdateWeekdayURL}}">
    <input type="hidden" name="xsrftoken" value="{{.UpdateWeekdayXSRFToken}}">
    <select name="favorite">
      {{range $index, $weekday := .Weekdays}}
        <option value="{{$index}}" {{if eq $.WeekdayIndex $index}}selected{{end}}>{{$weekday}}</option>{{end}}
    </select>
    <button type="submit">change</button>
  </form>
{{end}}`

const widget = fmt.Sprintf(`<div id="gitkitWidget"></div>
<script type="text/javascript">


  var brandingURL = "https://%v.appspot.com/auth/accountChooserBranding.html";
  var faviconURL  = "https://%v.appspot.com/favicon.ico";
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
`, routes.AppID(), routes.AppID())
