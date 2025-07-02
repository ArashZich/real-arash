package template

import (
	"bytes"
	"html/template"

	"github.com/ARmo-BigBang/kit/exp"
)

// RenderRedirectTemplate renders the redirect template
func RenderRedirectTemplate(successFull bool, msg string, targetWebsite string) string {
	tmp, _ := template.New("html").Parse(redirectTemplate)
	var tpl bytes.Buffer
	title := exp.TerIf(successFull, "پرداخت شما با موفقیت انجام شد", "پرداخت ناموفق")
	_ = tmp.Execute(&tpl, map[string]interface{}{
		"successFull":   successFull,
		"targetWebsite": targetWebsite,
		"title":         title,
		"message":       msg,
	})
	return tpl.String()
}

const redirectTemplate = `
<!DOCTYPE html>
<html lang="fa" dir="rtl">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>ARmo</title>
</head>
<style>
  @font-face {
    font-family: "IranSans";
    src: url("https://armogroup.storage.iran.liara.space/fonts/IranSans/IRANSansWeb%28FaNum%29_Medium.woff") format("woff");
    /* Add more formats if needed (e.g., woff2, ttf, etc.) */
  }

  html {
    font-family: "IranSans", Times, serif;
  }

  body {
    height: 95vh;
  }

  .container {
    display: flex;
    flex-direction: column;
    height: 100%;
    justify-content: center;
    align-items: center;
  }

  .body_header--success {
    color: #118d57;
  }

  .body_header--error {
    color: #b71d18;
  }

  .body_div_footer {
    margin-top: 24px;
    text-align: center;
  }

  .body_div_footer_button {
    width: 400px;
    padding: 12px 0;
    position: relative;
    border-radius: 8px;
    border-color: #5e35b1;
    color: #5e35b1;
    font-size: 15px;
    background-color: #fff;
    cursor: pointer;
    overflow: hidden;
    margin-top: 24px;
    font-family: "IranSans";
  }

  .body_div_footer:hover .body_div_footer_button {
    color: white;
    background-color: #5e35b1;
  }
</style>
<body>
  <div class="container">
  <div class="{{if .successFull}}body_header--success{{else}}body_header--error{{end}}">
  <h2 class="body_header_title">{{ .message }}</h2>
    </div>
    <div class="body_div_footer">
      <p>در حال بازگشت به سایت آرمو</p>
      <p>اگر تا <span id="countdown">10</span> ثانیه دیگر به صورت خودکار به سایت آرمو منتقل نشدید، روی دکمه زیر کلیک کنید.</p>
      <button type="button" onclick="redirect()" class="body_div_footer_button">بازگشت به سایت</button>
    </div>
  </div>
  <script>
    // Total seconds to wait
    let seconds = 10;

    function submitForm() {
      // Redirect to the target website
      window.location.href = "{{ .targetWebsite }}";
    }

    function countdown() {
      seconds = seconds - 1;
      if (seconds <= 0) {
        // Redirect to the target website
        submitForm();
      } else {
        // Update remaining seconds
        document.getElementById("countdown").innerHTML = seconds;
        // Count down using javascript
        window.setTimeout("countdown()", 1000);
      }
    }

    // Run countdown function
    countdown();

    // Function to redirect to the specified URL
    function redirect() {
      window.location.href = "{{ .targetWebsite }}";
    }
  </script>
</body>
</html>
`
