var app = new Vue({
  el: '#app',
  data: {
    scanner: null,
    activeCameraId: null,
    cameras: [],
    scans: []
  },
  mounted: function () {
    var self = this;
    self.scanner = new Instascan.Scanner({ video: document.getElementById('preview'), scanPeriod: 5 });
    self.scanner.addListener('scan', function (content, image) {
      self.scans.unshift({ date: +(Date.now()), content: content });
      var xhttp = new XMLHttpRequest();
      var resultText = "";
        xhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {    
            resultText = xhttp.responseText;

            if(resultText == "fail"){
              alert("Something wrong! Please call us for help")

            }else{
              alert("Item dispensing!")
            }
            
                          }
        };
        var url = "http://127.0.0.1/item/"+content;
        xhttp.open("GET",url, true);
        xhttp.send();
        
    });
    Instascan.Camera.getCameras().then(function (cameras) {
      self.cameras = cameras;
      if (cameras.length > 0) {
        self.activeCameraId = cameras[0].id;
        self.scanner.start(cameras[0]);
      } else {
        console.error('No cameras found.');
      }
    }).catch(function (e) {
      console.error(e);
    });
  },
  methods: {
    formatName: function (name) {
      return name || '(unknown)';
    },
    selectCamera: function (camera) {
      this.activeCameraId = camera.id;
      this.scanner.start(camera);
    }
  }
});




