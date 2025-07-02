package template

import (
	"bytes"
	"text/template"
)

func RenderExpiredTemplate(message string) string {
	tmp, _ := template.New("html").Parse(ExpiredTemplate)
	var tpl bytes.Buffer
	_ = tmp.Execute(&tpl, map[string]interface{}{
		"message": message,
	})
	return tpl.String()
}

const ExpiredTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>ARMO</title>
      <meta
      name="description"
      content="Virtual Try On - Experience augmented reality through our interactive platform."
    />
    <meta
      name="keywords"
      content="AR, Augmented Reality, Virtual Try On, AR Applications, AR Technology, AR User Experience, AR Developers, Virtual Reality, AR Hardware, AR Benefits, AR Challenges"
    />
    <meta name="author" content="ARmo Group" />
    <meta property="og:title" content="Virtual Try On" />
    <meta
      property="og:description"
      content="Explore the latest in augmented reality with our Virtual Try On platform."
    />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://armogroup.tech" />
    <meta
      property="og:image"
      content="https://bytebase.armogroup.tech/api/v1/files/download/armo-logo-main.jpg
    "
    />
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
        margin: 0;
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100vh;
        background-color: black;
        color: white;
      }
    </style>
  </head>
  <body>
    <h1>{{.message}}</h1>
  </body>
</html>
`
