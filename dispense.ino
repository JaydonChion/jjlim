#include <ESP8266WiFi.h>
#include <ESP8266HTTPClient.h>
#include <Servo.h>
#define led_pin 5
Servo myservo; // create servo object to control a servo

const char* ssid = "ssid"; // Your ssid
const char* password = "password"; // Your password
 
void setup () {
 
Serial.begin(115200);
WiFi.begin(ssid, password);
myservo.attach(16);
myservo.write(0);
pinMode(led_pin,OUTPUT);
 
while (WiFi.status() != WL_CONNECTED) {
 
delay(1000);
Serial.print("Connecting..");
 
}
 
}
 
void loop() {

if (WiFi.status() == WL_CONNECTED) { //Check WiFi connection status
  Serial.println("connected");
HTTPClient http;  //Declare an object of class HTTPClient
 
http.begin("http://your_web_server_ip_address/machineRead/1");  //Specify request destination
delay(100);
int httpCode = http.GET();        //Send the request
 
if (httpCode > 0) { //Check the returning code
 
String payload = http.getString();   //Get the request response payload
Serial.println(payload);                     //Print the response payload
if (payload == "dispense") {
   myservo.write(90); // sets the servo position according to the scaled value
   digitalWrite(led_pin,HIGH);
   delay(1000);
   myservo.write(0); // sets the servo position according to the scaled value
   digitalWrite(led_pin,LOW);
   delay(100);
  
  }
}
Serial.flush();
 
http.end();   //Close connection
 
}
 
delay(500);    //Send a request every 30 seconds
 
}
