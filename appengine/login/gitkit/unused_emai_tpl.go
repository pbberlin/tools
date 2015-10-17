package gitkit

// Email templates for OOB action
// OOB is "out of band"
const (
	emailTemplateResetPassword = `<p>Dear user,</p>
<p>
Forgot your password?<br>
FavWeekday received a request to reset the password for your account <b>%[1]s</b>.<br>
To reset your password, click on the link below (or copy and paste the URL into your browser):<br>
<a href="%[2]s">%[2]s</a><br>
</p>
<p>FavWeekday Support</p>`

	emailTemplateChangeEmail = `<p>Dear user,</p>

<p>
Want to use another email address to sign into FavWeekday?<br>
FavWeekday received a request to change your account email address from %[1]s to <b>%[2]s</b>.<br>
To change your account email address, click on the link below (or copy and paste the URL into your browser):<br>
<a href="%[3]s">%[3]s</a><br>
</p>
<p>FavWeekday Support</p>`

	emailTemplateVerifyEmail = `Dear user,

<p>Thank you for creating an account on FavWeekday.</p>
<p>To verify your account email address, click on the link below (or copy and paste the URL into your browser):</p>
<p><a href="%[1]s">%[1]s</a></p>

<br>
<p>FavWeekday Support</p>`
)
