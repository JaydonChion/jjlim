var app = new Vue({
  el: "#app",
  data: {
    scanner: null,
    activeCameraId: null,
    cameras: [],
    scans: []
  },
  mounted: function() {
    var self = this;
    self.scanner = new Instascan.Scanner({
      video: document.getElementById("preview"),
      scanPeriod: 5
    });
    self.scanner.addListener("scan", function(content, image) {
      self.scans.unshift({ date: +Date.now(), content: content });
      var xhttp = new XMLHttpRequest();
      var resultText = "";
      xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
          resultText = xhttp.responseText;
          if (resultText == "fail") {
            document.getElementById("scan").innerHTML =
              "<div style='color:red; font-size:25px;'>Something went wrong! Please call us for help</div>";
            // alert("Something wrong! Please call us for help");
          } else {
            document.getElementById("scan").innerHTML =
              "<div style='color:lightgreen; font-size:25px;'>Sucessfully, Item will be dispensing in just a second!</div>";
            // alert("Item dispensing!");
          }
        }
      };
      var url = "http://192.168.43.108/item/" + content;
      xhttp.open("GET", url, true);
      xhttp.send();
    });
    Instascan.Camera.getCameras()
      .then(function(cameras) {
        self.cameras = cameras;
        if (cameras.length > 0) {
          self.activeCameraId = cameras[0].id;
          self.scanner.start(cameras[0]);
        } else {
          console.error("No cameras found.");
        }
      })
      .catch(function(e) {
        console.error(e);
      });
  },
  methods: {
    formatName: function(name) {
      return name || "(unknown)";
    },
    selectCamera: function(camera) {
      this.activeCameraId = camera.id;
      this.scanner.start(camera);
    }
  }
});
