package template

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/google/uuid"
)

func RenderNecklaceTemplate(id int, file_uri string, title string, uid uuid.UUID, organization_id int, width int, length int, category string) string {
	tmp, _ := template.New("html").Parse(NecklaceTemplate)
	var tpl bytes.Buffer
	_ = tmp.Execute(&tpl, map[string]interface{}{
		"id":              id,
		"file_uri":        file_uri,
		"title":           title,
		"uid":             uid,
		"organization_id": organization_id,
		"width":           width,
		"length":          length,
		"category":        category, // اضافه کردن کتگوری به داده‌های ارسالی به قالب

	})

	// Extract the JavaScript code
	jsRegex := regexp.MustCompile(`<script[^>]*type=["']text/javascript["'][^>]*>([\s\S]*?)<\/script>`)
	matches := jsRegex.FindAllSubmatch(tpl.Bytes(), -1)

	if matches == nil {
		return tpl.String() // No JS code found, return original template
	}

	var uglyTpl = tpl.Bytes()
	for _, match := range matches {
		jsCode := match[1]

		// Run uglify-js on the JavaScript code
		uglyJsCode, err := uglifyJs(jsCode)
		if err != nil {
			panic(err)
		}

		// Replace the original JavaScript code with the minified code in the template
		uglyTpl = bytes.Replace(uglyTpl, jsCode, uglyJsCode, 1)
	}

	return string(uglyTpl)
}

const NecklaceTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="/assets/images/favicon/favicon.ico" />
    <link
      rel="apple-touch-icon"
      sizes="180x180"
      href="/assets/images/favicon/apple-touch-icon.png"
    />
    <link
      rel="icon"
      type="image/png"
      sizes="32x32"
      href="/assets/images/favicon/favicon-32x32.png"
    />
    <link
      rel="icon"
      type="image/png"
      sizes="16x16"
      href="/assets/images/favicon/favicon-16x16.png"
    />
    <link rel="manifest" href="/assets/images/favicon/site.webmanifest" />
    <link
      rel="mask-icon"
      href="/assets/images/favicon/safari-pinned-tab.svg"
      color="#5bbad5"
    />
    <meta name="msapplication-TileColor" content="#603cba" />
    <meta name="theme-color" content="#ffffff" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />

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

    <!-- Terser -->
    <script src="https://cdn.jsdelivr.net/npm/terser/dist/bundle.min.js"></script>

    <!-- Include Axios CDN link -->
    <script src="https://cdn.jsdelivr.net/npm/axios/dist/axios.min.js"></script>

    <!-- INCLUDE WebAR.rocks.face MAIN SCRIPT -->
    <script src="/assets/dist/WebARRocksFace.js"></script>

    <!-- INCLUDE CANVAS2D HELPER -->
    <script src="/assets/helpers/WebARRocksFaceCanvas2DHelper.js"></script>

    <!-- INCLUDE RESIZER HELPER -->
    <script src="/assets/helpers/WebARRocksResizer.js"></script>
    <!-- INCLUDE MIRROR HELPER -->
    <script src="/assets/helpers/html2canvas.min.js"></script>

    <!-- INCLUDE LANDMARKS STABILIZER -->
    <script src="/assets/helpers/landmarksStabilizers/WebARRocksLMStabilizer.js"></script>
    <title>ARmo</title>


    <!-- Minify and replace JS scripts after the page is fully loaded -->
    <script>
      function minifyJS(jsCode) {
        const result = Terser.minify(jsCode);
        if (result.error) {
          console.error("Error minifying code:", result.error);
          return jsCode;
        }
        return result.code;
      }

      document.addEventListener("DOMContentLoaded", function() {
        const scriptTags = document.querySelectorAll("script[type='text/javascript']");
        scriptTags.forEach((scriptTag) => {
          const originalCode = scriptTag.innerText;
          const minifiedCode = minifyJS(originalCode);
          const newScriptTag = document.createElement("script");
          newScriptTag.type = "text/javascript";
          newScriptTag.innerText = minifiedCode;
          scriptTag.replaceWith(newScriptTag);
        });
      });
    </script>

    <!-- Original JavaScript code to be minified -->
    <script type="text/javascript">
    const _canvases = {
      face: null,
      overlay: null,
    };
    let _ctx = null,
      _earringImage = null;
    
    const _earringSettings = {
      image: "{{.file_uri}}",
      angleHide: 50, // head rotation angle in degrees from which we should hide the earrings
      angleHysteresis: 0.5, // add hysteresis to angleHide value, in degrees
      scale: 1, // width of the earring compared to the face width (1 -> 100% of the face width)
      pullUp: 0.05, // 0 -> earring are displayed at the bottom of the spotted position
      // 1 -> earring are displaed above the spotted position
      k: 0.7, // position is interpolated between 2 keypoints. this is the interpolation coefficient
      // 0-> earrings are at the bottom of the ear, 1-> earrings are further back
    };
    
    const _earringsVisibility = {
      right: false,
      left: false,
    };
    
    function start() {
      WebARRocksFaceCanvas2DHelper.init({
        spec: {
          NNCPath: "/assets/neuralNets/NN_EARS_4.json", // neural network model file
          canvas: _canvases.face,
        },
    
        callbackReady: function (err, spec) {
          // called when everything is ready
          if (err) {
            console.log("ERROR in demo.js: ");
            return;
          }
        },
    
        callbackTrack: function (detectState) {
          clear_canvas();
          if (detectState.isDetected) {
            //draw_faceCrop(detectState.faceCrop);
            draw_earrings(
              detectState.landmarks,
              detectState.faceWidth,
              detectState.ry,
              {{.width}}, 
              {{.length}},
              "{{.category}}"  // اضافه کردن کتگوری به تابع draw_earrings
            );
          } else {
            _earringsVisibility.right = true;
            _earringsVisibility.left = true;
          }
        },
      });
    }
    
    function mix_landmarks(posA, posB, k) {
      return [
        posA[0] * (1 - k) + posB[0] * k, // X
        posA[1] * (1 - k) + posB[1] * k, // Y
      ];
    }
    
    function draw_faceCrop(faceCrop) {
      _ctx.strokeStyle = "lime";
      _ctx.beginPath();
      _ctx.moveTo(faceCrop[0][0], faceCrop[0][1]);
      _ctx.lineTo(faceCrop[1][0], faceCrop[1][1]);
      _ctx.lineTo(faceCrop[2][0], faceCrop[2][1]);
      _ctx.lineTo(faceCrop[3][0], faceCrop[3][1]);
      _ctx.closePath();
      _ctx.stroke();
    }
    
    function cmToPixels(cm, dpi = 96) {
      const inches = cm / 2.54; // Convert centimeters to inches
      const pixels = inches * dpi; // Convert inches to pixels
      return pixels;
    }
    
    function draw_earrings(landmarks, faceWidth, ry, length, height, category) {
      // Variables to store positions of the earrings/necklace
      let leftPos, rightPos;
    
      // Check for the right ear visibility
      const rightEarringVisible =
        ry > -_earringSettings.angleHide - _earringSettings.angleHysteresis;
      if (rightEarringVisible) {
        rightPos = mix_landmarks(
          landmarks.rightEarBottom,
          landmarks.rightEarEarring,
          _earringSettings.k
        );
        _earringsVisibility.right = true;
      } else {
        _earringsVisibility.right = false;
      }
    
      // Check for the left ear visibility
      const leftEarringVisible =
        -ry > -_earringSettings.angleHide - _earringSettings.angleHysteresis;
      if (leftEarringVisible) {
        leftPos = mix_landmarks(
          landmarks.leftEarBottom,
          landmarks.leftEarEarring,
          _earringSettings.k
        );
        _earringsVisibility.left = true;
      } else {
        _earringsVisibility.left = false;
      }
    
      // Draw the necklace or bow tie if both earrings are visible
      if (
        _earringsVisibility.right &&
        _earringsVisibility.left &&
        rightPos &&
        leftPos
      ) {
        const necklacePos = [
          (leftPos[0] + rightPos[0]) / 2, // Average X position
          (leftPos[1] + rightPos[1]) / 2, // Average Y position
        ];
        draw_earring(necklacePos, length, height, category);
      }
    }
    
    function draw_earring(pos, width, height, category) {
      // Use width and height directly for sizing the earring
      const dWidth = cmToPixels(width);
      const dHeight = cmToPixels(height);
      const dx = pos[0] - dWidth / 2.0; // Center the necklace horizontally
    
      // Adjust vertical position based on category
      let dy;
      if (category === "bow_tie") {
        dy = pos[1] - dHeight * _earringSettings.pullUp + 200;
      } else {
        dy = pos[1] - dHeight * _earringSettings.pullUp + 220;
      }
    
      _ctx.drawImage(_earringImage, dx, dy, dWidth, dHeight);
    }
    
    function clear_canvas() {
      _ctx.clearRect(0, 0, _canvases.overlay.width, _canvases.overlay.height);
    }
    
    function main() {
      // Create earring image:
      _earringImage = new Image();
      _earringImage.src = _earringSettings.image;
    
      // Get canvas from the DOM:
      _canvases.face = document.getElementById("WebARRocksFaceCanvas");
      _canvases.overlay = document.getElementById("overlayCanvas");
    
      // Create 2D context for the overlay canvas (where the earring are drawn):
      _ctx = _canvases.overlay.getContext("2d");
    
      // Set the canvas to fullscreen
      // and add an event handler to capture window resize:
      WebARRocksResizer.size_canvas({
        isFullScreen: true,
        canvas: _canvases.face, // WebARRocksFace main canvas
        overlayCanvas: [_canvases.overlay], // other canvas which should be resized at the same size of the main canvas
        callback: start,
      });
    }
    
    window.addEventListener("load", main);
 
    // Function to generate a UUID v4
    function uuidv4() {
      return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function(c) {
        const r = (Math.random() * 16) | 0,
          v = c === "x" ? r : (r & 0x3) | 0x8;
        return v.toString(16);
      });
    }


    // Function to get user agent information
    function useUserAgent() {
      var userAgentInfo = {
        browser: "",
        operatingSystem: "",
        device: "",
      };

      var userAgent = window.navigator.userAgent;

      // Get browser name
      var browser =
        userAgent.match(
          /(opera|chrome|safari|firefox|msie|trident(?=\/))\/?\s*(\d+)/i
        ) || [];
      var browserName = browser[1] || "";
      if (/trident/i.test(browserName)) {
        var rv = userAgent.match(/\brv[ :]+(\d+)/g) || [];
        browserName = "IE " + (rv[1] || "");
      }
      if (browserName === "Chrome") {
        var edge = userAgent.match(/\b(OPR|Edge)\/(\d+)/);
        if (edge != null)
          browserName = edge.slice(1).join(" ").replace("OPR", "Opera");
      }

      // Get operating system
      let osName = "";
      if (/android/i.test(userAgent)) {
        osName = "Android";
      } else {
        const os = userAgent.match(
          /(windows nt|mac|linux|windows phone|iPad|iPhone|iPod)/i
        );
        osName = os ? os[0] : "";
      }

        // Get device model
        let device;
        if (/Mobile/i.test(userAgent) && !/Tablet/i.test(userAgent)) {
          device = "Mobile";
        } else if (/Tablet/i.test(userAgent)) {
          device = "Tablet";
        } else {
          device = "Desktop";
        }

      userAgentInfo = {
        browser: browserName,
        operatingSystem: osName,
        device: device,
      };

      return userAgentInfo;
    }

    // Example usage of useUserAgent function
    var { browser, device, operatingSystem } = useUserAgent();
    var uuid = uuidv4();

    // Function to make a GET request to get IP information
    async function getIPService() {
      try {
        var response = await axios.get("https://api.ipify.org/?format=json");
        return response.data;
      } catch (error) {
        console.error("Error fetching IP:", error);
        throw error;
      }
    }

    // Function to track user information
    async function trackUser() {
      var initialStates = {
        name: "",
        ip: "",
        browser_agent: "",
        operating_sys: "",
        device: "",
        is_ar: null,
        is_3d: null,
        is_vr: null,
        url: "",
        product_uid: null,
        organization_id: null,
        visit_uid: "",
      };

      var obj = {
        ...initialStates,
        name: "{{.title}}",
        product_uid: "{{.uid}}",
        browser_agent: browser,
        operating_sys: operatingSystem,
        device: device,
        is_3d: true,
        is_ar: true,
        is_vr: false,
        organization_id: {{.organization_id}},
        visit_uid: uuid,
      };

      try {
        // Get IP information
        var ipData = await getIPService();

        // Update obj with IP information
        obj = {
          ...obj,
          ip: ipData.ip,
        };
        // Make a POST request to track user information
        var response = await axios.post(
          "https://reality.armogroup.tech/api/v1/views",
          obj
        );
      } catch (error) {
        console.error("Error tracking user:", error);
      }
    }

    // Function to track the duration of the user's visit
      async function viewDuration(){
        var items = {
          visit_duration:10,
          product_uid: "{{.uid}}",
          organization_id: {{.organization_id}},
          visit_uid: uuid,
        }

        try {
          var res = await axios.post(
            "https://reality.armogroup.tech/api/v1/views/duration",
            items
          );
        } catch (error) {
          console.error("Error Duration:", error);
        }
      }

       // Call viewDuration every 10 seconds
        setInterval(() => {
          viewDuration()
        }, 10000);

        // Example usage of the trackUser function
        window.onload = function() {
          if (!window.trackUserCalled) {
            trackUser();
            window.trackUserCalled = true; // Set the flag to true to prevent future calls
          }
        };
  </script>
   

    <style type="text/css">
      body {
        margin: 0;
      }
      canvas {
        transform: scaleX(-1); /* flip horizontally (mirror mode) */
        position: absolute;
        top: 0;
        left: 0;
        max-width: 100vw;
        max-height: 100vh;
      }
      #overlayCanvas {
        z-index: 10;
      }
      #WebARRocksFaceCanvas {
        z-index: 1;
      }

      #divMain {
        position: relative;
        display: flex;
        justify-content: center;
        align-items: center;
      }

      #logoImage {
        height: 80px;
        width: 80px;
        z-index: 100;
      }

      #divLogo {
        display: flex;
        position: fixed;
        top: 28px;
        left: 28px;
        z-index: 100;
        width: 100vw;
        flex-direction: row;
      }
      #VTOButtons {
        display: flex;
        position: fixed;
        bottom: 48px;
        z-index: 100;
        flex-direction: row;
      }
    </style>
  </head>
  <body>
    <div id="divMain">
      <div id="divLogo">
        <img src="/assets/images/logo.png" alt="logo" id="logoImage" />
      </div>
      <!-- canvas where the earring will be displayed: -->
      <canvas width="600" height="600" id="overlayCanvas"></canvas>

      <!-- canvas where the video will be displayed: -->
      <canvas width="600" height="600" id="WebARRocksFaceCanvas"></canvas>
  
    </div>  
 
  </body>
</html>
`
