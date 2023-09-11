package templates

type Purchase struct {
	Description string
	Type        string
}
type ReceiptTemplate struct {
	Name      string
	Purchases []Purchase
	Total     string
}

var ReceiptEmail = "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\">\n<html data-editor-version=\"2\" class=\"sg-campaigns\" xmlns=\"http://www.w3.org/1999/xhtml\">\n    <head>\n      <meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\">\n      <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, minimum-scale=1, maximum-scale=1\">\n      <!--[if !mso]><!-->\n      <meta http-equiv=\"X-UA-Compatible\" content=\"IE=Edge\">\n      <!--<![endif]-->\n      <!--[if (gte mso 9)|(IE)]>\n      <xml>\n        <o:OfficeDocumentSettings>\n          <o:AllowPNG/>\n          <o:PixelsPerInch>96</o:PixelsPerInch>\n        </o:OfficeDocumentSettings>\n      </xml>\n      <![endif]-->\n      <!--[if (gte mso 9)|(IE)]>\n  <style type=\"text/css\">\n    body {width: 600px;margin: 0 auto;}\n    table {border-collapse: collapse;}\n    table, td {mso-table-lspace: 0pt;mso-table-rspace: 0pt;}\n    img {-ms-interpolation-mode: bicubic;}\n  </style>\n<![endif]-->\n      <style type=\"text/css\">\n    body, p, div {\n      font-family: arial,helvetica,sans-serif;\n      font-size: 14px;\n    }\n    body {\n      color: #000000;\n    }\n    body a {\n      color: #1188E6;\n      text-decoration: none;\n    }\n    p { margin: 0; padding: 0; }\n    table.wrapper {\n      width:100% !important;\n      table-layout: fixed;\n      -webkit-font-smoothing: antialiased;\n      -webkit-text-size-adjust: 100%;\n      -moz-text-size-adjust: 100%;\n      -ms-text-size-adjust: 100%;\n    }\n    img.max-width {\n      max-width: 100% !important;\n    }\n    .column.of-2 {\n      width: 50%;\n    }\n    .column.of-3 {\n      width: 33.333%;\n    }\n    .column.of-4 {\n      width: 25%;\n    }\n    @media screen and (max-width:480px) {\n      .preheader .rightColumnContent,\n      .footer .rightColumnContent {\n        text-align: left !important;\n      }\n      .preheader .rightColumnContent div,\n      .preheader .rightColumnContent span,\n      .footer .rightColumnContent div,\n      .footer .rightColumnContent span {\n        text-align: left !important;\n      }\n      .preheader .rightColumnContent,\n      .preheader .leftColumnContent {\n        font-size: 80% !important;\n        padding: 5px 0;\n      }\n      table.wrapper-mobile {\n        width: 100% !important;\n        table-layout: fixed;\n      }\n      img.max-width {\n        height: auto !important;\n        max-width: 100% !important;\n      }\n      a.bulletproof-button {\n        display: block !important;\n        width: auto !important;\n        font-size: 80%;\n        padding-left: 0 !important;\n        padding-right: 0 !important;\n      }\n      .columns {\n        width: 100% !important;\n      }\n      .column {\n        display: block !important;\n        width: 100% !important;\n        padding-left: 0 !important;\n        padding-right: 0 !important;\n        margin-left: 0 !important;\n        margin-right: 0 !important;\n      }\n      .social-icon-column {\n        display: inline-block !important;\n      }\n    }\n  </style>\n      <!--user entered Head Start--><!--End Head user entered-->\n    </head>\n    <body>\n      <center class=\"wrapper\" data-link-color=\"#1188E6\" data-body-style=\"font-size:14px; font-family:arial,helvetica,sans-serif; color:#000000; background-color:#f1f2ff;\">\n        <div class=\"webkit\">\n          <table cellpadding=\"0\" cellspacing=\"0\" border=\"0\" width=\"100%\" class=\"wrapper\" bgcolor=\"#f1f2ff\">\n            <tr>\n              <td valign=\"top\" bgcolor=\"#f1f2ff\" width=\"100%\">\n                <table width=\"100%\" role=\"content-container\" class=\"outer\" align=\"center\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\">\n                  <tr>\n                    <td width=\"100%\">\n                      <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\">\n                        <tr>\n                          <td>\n                            <!--[if mso]>\n    <center>\n    <table><tr><td width=\"600\">\n  <![endif]-->\n                                    <table width=\"100%\" cellpadding=\"0\" cellspacing=\"0\" border=\"0\" style=\"width:100%; max-width:600px;\" align=\"center\">\n                                      <tr>\n                                        <td role=\"modules-container\" style=\"padding:0px 0px 0px 0px; color:#000000; text-align:left;\" bgcolor=\"#f1f2ff\" width=\"100%\" align=\"left\"><table class=\"module preheader preheader-hide\" role=\"module\" data-type=\"preheader\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"display: none !important; mso-hide: all; visibility: hidden; opacity: 0; color: transparent; height: 0; width: 0;\">\n    <tr>\n      <td role=\"module-content\">\n        <p></p>\n      </td>\n    </tr>\n  </table><table class=\"module\" role=\"module\" data-type=\"text\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"2ad54b32-f159-49b6-8c60-bcb9c8c13777\" data-mc-module-version=\"2019-10-22\">\n    <tbody>\n      <tr>\n        <td style=\"padding:18px 0px 18px 24px; line-height:22px; text-align:inherit; background-color:#febc5d;\" height=\"100%\" valign=\"top\" bgcolor=\"#febc5d\" role=\"module-content\"><div><div style=\"font-family: inherit; text-align: inherit\"><span style=\"color: #ffffff; font-size: 24px\">EQUIPHUNTER</span></div><div></div></div></td>\n      </tr>\n    </tbody>\n  </table><table class=\"module\" role=\"module\" data-type=\"text\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"377a9cf5-6e13-454b-a86e-5a79f72cc92a\" data-mc-module-version=\"2019-10-22\">\n    <tbody>\n      <tr>\n        <td style=\"padding:40px 0px 10px 0px; line-height:36px; text-align:inherit;\" height=\"100%\" valign=\"top\" bgcolor=\"\" role=\"module-content\"><div><div style=\"font-family: inherit; text-align: center\"><span style=\"font-style: normal; font-variant-caps: normal; font-weight: normal; letter-spacing: normal; text-align: start; text-indent: 0px; text-transform: none; word-spacing: 0px; -webkit-text-stroke-width: 0px; text-decoration: none; caret-color: rgb(0, 0, 0); color: #000000; white-space: pre-wrap; font-size: 36px\">Hi, </span><span style=\"font-style: normal; font-variant-caps: normal; font-weight: normal; letter-spacing: normal; text-align: start; text-indent: 0px; text-transform: none; word-spacing: 0px; -webkit-text-stroke-width: 0px; text-decoration: none; caret-color: rgb(0, 0, 0); color: #000000; white-space: pre-wrap; font-size: 36px\"><strong>{{.Name}}. Thank you for your purchase</strong></span></div><div></div></div></td>\n      </tr>\n    </tbody>\n  </table><table class=\"wrapper\" role=\"module\" data-type=\"image\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"63f7b8a6-68c9-4f59-8c01-8a9a58420573\">\n    <tbody>\n      <tr>\n        <td style=\"font-size:6px; line-height:10px; padding:20px 0px 20px 0px;\" valign=\"top\" align=\"center\">\n          <img class=\"max-width\" border=\"0\" style=\"display:block; color:#000000; text-decoration:none; font-family:Helvetica, arial, sans-serif; font-size:16px; max-width:50% !important; width:50%; height:auto !important;\" width=\"300\" alt=\"\" data-proportionally-constrained=\"true\" data-responsive=\"true\" src=\"https://vendor-proposals.s3.us-east-2.amazonaws.com/email-images/1329x513.png\">\n        </td>\n      </tr>\n    </tbody>\n  </table><table class=\"module\" role=\"module\" data-type=\"text\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"377a9cf5-6e13-454b-a86e-5a79f72cc92a.2\" data-mc-module-version=\"2019-10-22\">\n    <tbody>\n      <tr>\n        <td style=\"padding:20px 0px 20px 0px; line-height:22px; text-align:inherit;\" height=\"100%\" valign=\"top\" bgcolor=\"\" role=\"module-content\"><div><div style=\"font-family: inherit; text-align: center\"><span style=\"font-style: normal; font-variant-caps: normal; font-weight: normal; letter-spacing: normal; text-align: start; text-indent: 0px; text-transform: none; word-spacing: 0px; -webkit-text-stroke-width: 0px; text-decoration: none; caret-color: rgb(0, 0, 0); color: #000000; white-space: pre-wrap; font-size: 20px\">This is the detail of your purchase:</span></div><div></div></div></td>\n      </tr>\n    </tbody>\n  </table><table class=\"module\" role=\"module\" data-type=\"text\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"377a9cf5-6e13-454b-a86e-5a79f72cc92a.1\" data-mc-module-version=\"2019-10-22\">\n    <tbody>\n      <tr>\n        <td style=\"padding:10px 0px 10px 0px; line-height:18px; text-align:inherit;\" height=\"100%\" valign=\"top\" bgcolor=\"\" role=\"module-content\">\n          <div>\n            {{with .Purchases}}\n              {{range .}}\n              <div style=\"font-family: inherit; text-align: center\"><span style=\"font-size: 18px\"><strong>- {{.Type}}</strong></span><span style=\"font-size: 18px\"> {{.Description}}</span></div>          \n              {{ end }}\n            {{end}}\n            \n          </div>\n        </td>\n      </tr>\n    </tbody>\n  </table><table class=\"module\" role=\"module\" data-type=\"text\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"377a9cf5-6e13-454b-a86e-5a79f72cc92a.1.2\" data-mc-module-version=\"2019-10-22\">\n    <tbody>\n      <tr>\n        <td style=\"padding:10px 0px 10px 0px; line-height:18px; text-align:inherit;\" height=\"100%\" valign=\"top\" bgcolor=\"\" role=\"module-content\"><div><div style=\"font-family: inherit; text-align: center\"><span style=\"font-size: 24px\"><strong>Total:</strong></span><span style=\"font-size: 24px\"> {{.Total}} USD</span></div><div></div></div></td>\n      </tr>\n    </tbody>\n  </table><table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"module\" data-role=\"module-button\" data-type=\"button\" role=\"module\" style=\"table-layout:fixed;\" width=\"100%\" data-muid=\"77b3cc70-85aa-41d0-9bbb-024b97bc04b2\">\n      <tbody>\n        <tr>\n          <td align=\"center\" bgcolor=\"\" class=\"outer-td\" style=\"padding:50px 0px 50px 0px;\">\n            <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" class=\"wrapper-mobile\" style=\"text-align:center;\">\n              <tbody>\n                <tr>\n                <td align=\"center\" bgcolor=\"#febc5d\" class=\"inner-td\" style=\"border-radius:6px; font-size:16px; text-align:center; background-color:inherit;\">\n                  <a href=\"https://equiphunter.com/login\" style=\"background-color:#febc5d; border:1px solid #333333; border-color:#333333; border-radius:20px; border-width:1px; color:#ffffff; display:inline-block; font-size:14px; font-weight:400; letter-spacing:0px; line-height:14px; padding:12px 18px 12px 18px; text-align:center; text-decoration:none; border-style:solid; width:216px;\" target=\"_blank\">TAKE A LOOK!</a>\n                </td>\n                </tr>\n              </tbody>\n            </table>\n          </td>\n        </tr>\n      </tbody>\n    </table><table class=\"module\" role=\"module\" data-type=\"divider\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"4c9bbf96-a686-4d52-908e-1fe9dd88ea4c\">\n    <tbody>\n      <tr>\n        <td style=\"padding:0px 0px 0px 0px;\" role=\"module-content\" height=\"100%\" valign=\"top\" bgcolor=\"\">\n          <table border=\"0\" cellpadding=\"0\" cellspacing=\"0\" align=\"center\" width=\"100%\" height=\"1px\" style=\"line-height:1px; font-size:1px;\">\n            <tbody>\n              <tr>\n                <td style=\"padding:0px 0px 1px 0px;\" bgcolor=\"#dcdcdc\"></td>\n              </tr>\n            </tbody>\n          </table>\n        </td>\n      </tr>\n    </tbody>\n  </table><table class=\"module\" role=\"module\" data-type=\"text\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" width=\"100%\" style=\"table-layout: fixed;\" data-muid=\"377a9cf5-6e13-454b-a86e-5a79f72cc92a.1.1.1\" data-mc-module-version=\"2019-10-22\">\n    <tbody>\n      <tr>\n        <td style=\"padding:40px 0px 20px 0px; line-height:22px; text-align:inherit;\" height=\"100%\" valign=\"top\" bgcolor=\"\" role=\"module-content\"><div><div style=\"font-family: inherit; text-align: center\"><span style=\"font-style: normal; font-variant-caps: normal; font-weight: normal; letter-spacing: normal; text-align: start; text-indent: 0px; text-transform: none; word-spacing: 0px; -webkit-text-stroke-width: 0px; text-decoration: none; caret-color: rgb(0, 0, 0); color: #000000; white-space: pre-wrap\">© 2020 - Equiphunter. Copyright all rights reserved.</span></div><div></div></div></td>\n      </tr>\n    </tbody>\n  </table></td>\n                                      </tr>\n                                    </table>\n                                    <!--[if mso]>\n                                  </td>\n                                </tr>\n                              </table>\n                            </center>\n                            <![endif]-->\n                          </td>\n                        </tr>\n                      </table>\n                    </td>\n                  </tr>\n                </table>\n              </td>\n            </tr>\n          </table>\n        </div>\n      </center>\n    </body>\n  </html>"
