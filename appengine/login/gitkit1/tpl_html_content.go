package gitkit1

const IDCardHTML = `

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
`

const UserInfoHTML = `
{{if not .User}}


{{else}}
  <p>You are logged in as {{.User.Name}}<br>
  ID: {{.User.ID}} <br>
  Email: {{.User.Email}}</p>

{{end}}`

const widgetHTML = `

<div id="gitkitWidget"></div>
<script type="text/javascript">

  var brandingURL = "{{ .BrandingURL }}";
  var faviconURL  = "{{ .FaviconURL  }}";
  var accountChooserConfig = {
    title:    'Sign in to {{.SiteName}}',
    favicon:  faviconURL,
    branding: brandingURL,
  };

  // http and relative urls are rejected as 'invalid domain'
  // Therefore fallback
  if (window.location.hostname == "localhost") {
    faviconURL = "";
    brandingURL = "";
    accountChooserConfig = {};
  }




  var signInSuccessUrl  = '{{.SignInSuccessURL}}';
  var signOutUrl = '{{.SignOutURL}}';

  // We erred in reversing the forward slash escaping, 
  // signInSuccessUrl = signInSuccessUrl.replace('\/', '/');
  // it was not necessary

  console.log("signInSuccessUrl: ",signInSuccessUrl, "  ---  ",signOutUrl)

  // siteName required for "manage account"
  google.identitytoolkit.start(
    '#gitkitWidget',
    {
      apiKey:           '{{.BrowserAPIKey}}',
      signInSuccessUrl: signInSuccessUrl,
      signOutUrl:       signOutUrl,
      oobActionUrl:     '{{.OOBActionURL}}',
      signInOptions:    ['google','facebook'],
      siteName:         '{{.SiteName}}',
      acUiConfig:        accountChooserConfig,      
      queryParameterForSignInSuccessUrl: 'red',
      callbacks: {
        signInSuccess: function(tokenString, accountInfo, opt_signInSuccessUrl) {
          console.log("opt_signInSuccessUrl",opt_signInSuccessUrl);
          console.log("accountInfo",accountInfo.email,accountInfo.providerId,accountInfo.displayName);

          // prevents redirect to signInSuccessUrl
          // return false; 
          return true;
        }
      },
      /*
        // https://developers.google.com/identity/toolkit/web/setup-frontend
         ajaxSender: function(url, data, completed) { },
      */      
    },
    '{{.POSTBody}}');

</script>

`
