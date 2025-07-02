package template

import (
	"bytes"
	"regexp"
	"text/template"

	"github.com/google/uuid"
)

func RenderShoesTemplate(id int, file_uri string, title string, uid uuid.UUID, organization_id int) string {
	tmp, _ := template.New("html").Parse(shoesTemplate)
	var tpl bytes.Buffer
	_ = tmp.Execute(&tpl, map[string]interface{}{
		"id":              id,
		"file_uri":        file_uri,
		"title":           title,
		"uid":             uid,
		"organization_id": organization_id,
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

const shoesTemplate = `
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

    <!-- INCLUDE WEBARROCKSHAND SCRIPT -->
    <script src="/assets/dist/WebARRocksHand.js"></script>

    <!-- THREE.JS - REPLACE IT BY three.min.js FOR PROD -->
    <script src="/assets/library/three/three.js"></script>

    <!-- THREE.JS HELPERS -->
    <script src="/assets/library/three/GLTFLoader.js"></script>

    <!-- WEBARROCKSHAND THREEJS VTO HELPER -->
    <script src="/assets/helpers/HandTrackerThreeHelper.js"></script>

    <!-- INCLUDE LANDMARKS STABILIZER -->
    <script src="/assets/helpers/landmarksStabilizers/OneEuroLMStabilizer.js"></script>

    <!-- POSEFLIP FILTER -->
    <script src="/assets/helpers/PoseFlipFilter.js"></script>


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
    const _settings = {
      threshold: 0.8, // detection sensitivity, between 0 and 1
      NNVersion: 28, // best: 28
      
      // CONVERSES SHOES:
      // 3D models:
      shoeRightPath: "{{.file_uri}}",
      isModelLightMapped: false,
      occluderPath: "/assets/file/occluder.glb",
      
      // pose settings:
      scale: 0.95,
      translation: [0, -0.02, 0], // Z -> verical, Y+ -> front way
      
      // debug flags:
      debugCube: false, // Add a cube
      debugDisplayLandmarks: true,
      };
      
      const _three = {
      loadingManager: null,
      };
      
      const _states = {
      notLoaded: -1,
      loading: 0,
      running: 1,
      busy: 2,
      };
      let _state = _states.notLoaded;
      let _isSelfieCam = false;
      
      function setFullScreen(cv) {
      cv.width = window.innerWidth;
      cv.height = window.innerHeight;
      }
      
      // entry point:
      function main() {
      _state = _states.loading;
      
      const handTrackerCanvas = document.getElementById("handTrackerCanvas");
      const VTOCanvas = document.getElementById("ARCanvas");
      
      setFullScreen(handTrackerCanvas);
      setFullScreen(VTOCanvas);
      
      HandTrackerThreeHelper.init({
        poseLandmarksLabels: [
          "ankleBack",
          "ankleOut",
          "ankleIn",
          "ankleFront",
          "heelBackOut",
          "heelBackIn",
          "pinkyToeBaseTop",
          "middleToeBaseTop",
          "bigToeBaseTop",
        ],
        enableFlipObject: true, //true,
        cameraZoom: 1,
        freeZRot: false,
        threshold: _settings.threshold,
        scanSettings: {
          multiDetectionSearchSlotsRate: 0.5,
          multiDetectionMaxOverlap: 0.3,
          multiDetectionOverlapScaleXY: [0.5, 1],
          multiDetectionEqualizeSearchSlotScale: true,
          multiDetectionForceSearchOnOtherSide: true,
          multiDetectionForceChirality: 1,
          disableIsRightHandNNEval: true,
          overlapFactors: [1.0, 1.0, 1.0],
          translationScalingFactors: [0.3, 0.3, 1],
          nScaleLevels: 2, // in the higher scale level, the size of the detection window is the smallest video dimension
          scale0Factor: 0.5,
        },
        VTOCanvas: VTOCanvas,
        handTrackerCanvas: handTrackerCanvas,
        debugDisplayLandmarks: false,
        NNsPaths: ["/assets/neuralNets/NN_FOOT_28.json"],
        maxHandsDetected: 2,
        stabilizationSettings: {
          //qualityFactorRange: [0.4, 0.7],
          NNSwitchMask: {
            isRightHand: false,
            isFlipped: false,
          },
        },
        landmarksStabilizerSpec: {
          minCutOff: 0.001,
          beta: 3, // lower => more stabilized
        },
      })
        .then(function (three) {
        handTrackerCanvas.style.zIndex = 3; // fix a weird bug on iOS15 / safari
        start(three);
        })
        .catch(function (err) {
        console.log("INFO in main.js: an error happens ", err);
        });
      }
      
      function start(three) {
      three.loadingManager.onLoad = function () {
        console.log("INFO in main.js: Everything is loaded");
        _state = _states.running;
      };
      
      // set tonemapping:
      three.renderer.toneMapping = THREE.ACESFilmicToneMapping;
      three.renderer.outputEncoding = THREE.sRGBEncoding;
      
      // set lighting:
      if (!_settings.isModelLightMapped) {
        const pointLight = new THREE.PointLight(0xffffff, 2);
        const ambientLight = new THREE.AmbientLight(0xffffff, 0.8);
        three.scene.add(pointLight, ambientLight);
      }
      
      // add a debug cube:
      if (_settings.debugCube) {
        const s = 1;
        const cubeGeom = new THREE.BoxGeometry(s, s, s);
        const cubeMesh = new THREE.Mesh(cubeGeom, new THREE.MeshNormalMaterial());
        HandTrackerThreeHelper.add_threeObject(cubeMesh);
      }
      
      function transform(threeObject) {
        threeObject.scale.multiplyScalar(_settings.scale);
        threeObject.position.add(
        new THREE.Vector3().fromArray(_settings.translation)
        );
      }
      
      // load the shoes 3D model:
      new THREE.GLTFLoader().load(_settings.shoeRightPath, function (gltf) {
        const shoe = gltf.scene;
        transform(shoe);
        HandTrackerThreeHelper.add_threeObject(shoe);
      });
      
      new THREE.GLTFLoader(three.loadingManager).load(
        _settings.occluderPath,
        function (gltf) {
        const occluder = gltf.scene.children[0];
        transform(occluder);
        HandTrackerThreeHelper.add_threeOccluder(occluder);
        }
      );
      } //end start()
      
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
          is_3d: false,
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
        overflow: hidden;
      }
      canvas {
        position: fixed;
        top: 0;
        left: 0;
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
      #handTrackerCanvas {
        z-index: 1;
      }
      #ARCanvas {
        z-index: 4;
      }
      canvas {
        /* mirror the canvas - useful if camera is in selfie mode only: */
        /*transform: rotateY(180deg);*/
      }
    </style>
  </head>
  <body>
    <div id="divMain">
      <div id="divLogo">
        <img src="/assets/images/logo.png" alt="logo" id="logoImage" />
      </div>
	  <div id="canvases">
        <!-- Canvas used by the hand tracker and displaying the video -->
        <canvas id="handTrackerCanvas"></canvas>

        <!-- Canvas displayed above, with the THREE.js scene -->
        <canvas id="ARCanvas"></canvas>
      </div>
    </div>

  </body>
</html>
`
