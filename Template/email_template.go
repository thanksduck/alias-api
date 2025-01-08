package template

import "fmt"

func VerifyEmailTemplate(name, magicLink string) (htmlBody, textBody string) {
	htmlBody = fmt.Sprintf(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
	<html dir="ltr" lang="en">
	  <head>
		<link rel="preload" as="image" href="https://static.20032003.xyz/1as/1as-logo.png" />
		<meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
		<meta name="x-apple-disable-message-reformatting" />
	  </head>
	  <body style='background-color:#ffffff;font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen-Sans,Ubuntu,Cantarell,"Helvetica Neue",sans-serif'>
		<table align="center" width="100%%" border="0" cellpadding="0" cellspacing="0" role="presentation" style='max-width:37.5em;margin:0 auto;padding:20px 25px 48px;background-image:url("/static/raycast-bg.png");background-position:bottom;background-repeat:no-repeat, no-repeat'>
		  <tbody>
			<tr style="width:100%%">
			  <td>
				<img alt="One Alias Service" height="48" src="https://static.20032003.xyz/1as/1as-logo.png" style="display:block;outline:none;border:none;text-decoration:none" width="73" />
				<h1 style="font-size:28px;font-weight:bold;margin-top:48px">🪄 Dear %s, Your magic verification link</h1>
				<table align="center" width="100%%" border="0" cellpadding="0" cellspacing="0" role="presentation" style="margin:24px 0">
				  <tbody>
					<tr>
					  <td>
						<p style="font-size:16px;line-height:26px;margin:16px 0">
						  <a href="%s" style="color:#FF6363;text-decoration-line:none" target="_blank">👉 Click here to Verify 👈</a>
						</p>
						<p style="font-size:16px;line-height:26px;margin:16px 0">If you didn't request this, please ignore this email.</p>
					  </td>
					</tr>
				  </tbody>
				</table>
				<p style="font-size:16px;line-height:26px;margin:16px 0">Best,<br />- One Alias Service Team</p>
				<hr style="width:100%%;border:none;border-top:1px solid #eaeaea;border-color:#dddddd;margin-top:48px" />
				<img height="32" alt="1@S" src="https://static.20032003.xyz/1as/1as-logo.png" style="display:block;outline:none;border:none;text-decoration:none;-webkit-filter:grayscale(100%%);filter:grayscale(100%%);margin:20px 0" width="48" />
				<p style="font-size:12px;line-height:24px;margin:16px 0;color:#8898aa;margin-left:4px">One Alias Service</p>
				<p style="font-size:12px;line-height:24px;margin:16px 0;color:#8898aa;margin-left:4px">Central Vista, Sector 17, Chandigarh, 160017</p>
			  </td>
			</tr>
		  </tbody>
		</table>
	  </body>
	</html>`, name, magicLink)
	textBody = fmt.Sprintf(`Dear %s,
	We hope this email finds you well. Thank you for choosing One Alias Service.
	
	Please verify your email address by clicking on the link below:
	
	%s
	
	If you didn't request this verification, please disregard this email. Your account security is important to us.
	
	Thank you for your trust in One Alias Service. We're here to assist you with any questions or concerns.
	
	Best regards,
	The One Alias Service Team
	`, name, magicLink)
	return htmlBody, textBody
}

func ForgetPasswordTemplate(name, resetLink string) (htmlBody, textBody string) {
	htmlBody = fmt.Sprintf(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html dir="ltr" lang="en">
  <head>
    <link
      rel="preload"
      as="image"
      href="https://static.20032003.xyz/1as/1as-logo.png"
    />
    <meta content="text/html; charset=UTF-8" http-equiv="Content-Type" />
    <meta name="x-apple-disable-message-reformatting" />
    <!--$-->
  </head>
  <div
    style="display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0"
  >
  </div>
  <body style="background-color:#f6f9fc;padding:10px 0">
    <table
      align="center"
      width="100%%"
      border="0"
      cellpadding="0"
      cellspacing="0"
      role="presentation"
      style="max-width:37.5em;background-color:#ffffff;border:1px solid #f0f0f0;padding:45px"
    >
      <tbody>
        <tr style="width:100%%">
          <td>
            <img
              alt="One Alias Service"
              height="48"
              src="https://static.20032003.xyz/1as/1as-logo.png"
              style="display:block;outline:none;border:none;text-decoration:none"
              width="73"
            />
            <table
              align="center"
              width="100%%"
              border="0"
              cellpadding="0"
              cellspacing="0"
              role="presentation"
            >
              <tbody>
                <tr>
                  <td>
                    <p
                      style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                    >
                      Hi
                      %s
                    </p>
                    <p
                      style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                    >
                      Someone recently requested a password change for your
                      One Alias Service account. If this was you, you can set a new
                      password here:
                    </p>
                    <a
                      href="%s"
                      style="line-height:100%%;text-decoration:none;display:block;max-width:100%%;mso-padding-alt:0px;background-color:#007ee6;border-radius:4px;color:#fff;font-family:&#x27;Open Sans&#x27;, &#x27;Helvetica Neue&#x27;, Arial;font-size:15px;text-align:center;width:210px;padding:14px 7px 14px 7px"
                      target="_blank"
                      ><span
                        ><!--[if mso
                          ]><i
                            style="mso-font-width:350%%;mso-text-raise:21"
                            hidden
                            >&#8202;</i
                          ><!
                        [endif]--></span
                      ><span
                        style="max-width:100%%;display:inline-block;line-height:120%%;mso-padding-alt:0px;mso-text-raise:10.5px"
                        >Reset password</span
                      ><span
                        ><!--[if mso
                          ]><i style="mso-font-width:350%%" hidden
                            >&#8202;&#8203;</i
                          ><!
                        [endif]--></span
                      ></a
                    >
                    <p
                      style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                    >
                      Link is valid for 10 minutes. If you don&#x27;t want to change your password or
                      didn&#x27;t request this, just ignore and delete this
                      message.
                    </p>
                    <p
                      style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                    >
                      To keep your account secure, please don&#x27;t forward
                      this email to anyone. See our Help Center for<!-- -->
                      <a
                        href="https://1as.in/support"
                        style="color:#067df7;text-decoration-line:none;text-decoration:underline"
                        target="_blank"
                        >more security tips.</a
                      >
                    </p>
                    <p
                      style="font-size:16px;line-height:26px;margin:16px 0;font-family:&#x27;Open Sans&#x27;, &#x27;HelveticaNeue-Light&#x27;, &#x27;Helvetica Neue Light&#x27;, &#x27;Helvetica Neue&#x27;, Helvetica, Arial, &#x27;Lucida Grande&#x27;, sans-serif;font-weight:300;color:#404040"
                    >
                      With Love,<br> One Alias Service
                    </p>
                  </td>
                </tr>
              </tbody>
            </table>
          </td>
        </tr>
      </tbody>
    </table>
    <!--/$-->
  </body>
</html>
`, name, resetLink)
	textBody = "Dear " + name + "\n\n" +
		"You have requested to reset your password. Please click on the link below to reset your password. This link is valid for 10 minutes.\n\n" +
		resetLink + "\n\n" +
		"Thank you,\n" +
		"One Alias Service Team"
	return htmlBody, textBody
}
