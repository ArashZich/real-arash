package templates

import (
	"bytes"
	"html/template"
)

// RenderRedirectTemplate renders the redirect template
func RenderRedirectTemplate(method, url string, data map[string]string) (string, error) {

	tmp, _ := template.New("html").Parse(redirectTemplate)
	var tpl bytes.Buffer
	_ = tmp.Execute(&tpl, map[string]interface{}{
		"method": method,
		"payURL": url,
		"data":   data,
	})
	return tpl.String(), nil
}

const redirectTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <title>بازگشت به سایت</title>
    <style>
      @font-face {
        font-family: "IranSans";
        src: url("https://armogroup.storage.iran.liara.space/fonts/IranSans/IRANSansWeb%28FaNum%29_Medium.woff") format("woff");
      }

      html {
        font-family: "IranSans", Times, serif;
      }

      .text-center {
        text-align: center;
      }

      .mt-2 {
        margin-top: 2em;
      }

      .spinner {
        display: flex;
        justify-content: center;
        align-items: center;
        margin-top: 100px;
      }

      .dot {
        width: 15px;
        height: 15px;
        margin: 0 5px;
        background-color: #5e35b1;
        border-radius: 50%;
        animation: bounce 1.5s infinite ease-in-out;
      }

      .dot:nth-child(1) {
        animation-delay: -0.3s;
      }

      .dot:nth-child(2) {
        animation-delay: -0.15s;
      }

      @keyframes bounce {
        0%, 80%, 100% {
          transform: scale(0);
        }
        40% {
          transform: scale(1);
        }
      }
    </style>
  </head>
  <body>
    <div class="spinner">
      <div class="dot"></div>
      <div class="dot"></div>
      <div class="dot"></div>
    </div>
    <form
      class="text-center mt-2"
      method="{{ .method }}"
      action="{{ .payURL }}"
    >
      <p>در بازگشت به سایت پذیرنده.</p>
      <p>
        اگر تا 
        <span id="countdown">10</span>
        ثانیه به صورت خودکار به درگاه پرداخت منتقل نشدید، روی دکمه زیر کلیک کنید.
      </p>
      {{ range $key, $value := .data }}
      <input type="hidden" name="{{ $key }}" value="{{ $value }}" />
      {{ end }}

      <button type="submit">بازگشت به سایت</button>
    </form>
    <script>
      // Total seconds to wait
      let seconds = 10;

      function submitForm() {
        document.forms[0].submit();
      }

      function countdown() {
        seconds = seconds - 1;
        if (seconds <= 0) {
          // submit the form
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
    </script>
  </body>
</html>
`
