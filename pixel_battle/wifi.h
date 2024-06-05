#include <ESP8266WiFi.h>
#include <DNSServer.h>
#include <ESP8266WebServer.h>
#include <WiFiManager.h>         // https://github.com/tzapu/WiFiManager

// Set web server port number to 80
WiFiServer server(80);

// Variable to store the HTTP request
String header;