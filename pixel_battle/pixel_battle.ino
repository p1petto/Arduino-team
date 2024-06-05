#include "wifi.h"

#define LED_PIN D5
#define LED_NUM 288
#define COL 24
#define ROW 12

#include "FastLED.h"
CRGB leds[LED_NUM];

WiFiManager wifiManager;



bool is_correct_boundaries(int x, int y){
  if (x < 0 || x >= COL || y < 0 || y >= ROW)
    return false;
  return true;
}

bool is_correct_color(CRGB color){
  if (color[0] < 0 || color[0] > 255 || color[1] < 0 || color[1] > 255 || color[2] < 0 || color[2] > 255)
    return false;
  return true;
}

// Лента зигзаг, если x нечетный, то счёт идет в обратную сторону
int get_index(int x, int y){
  int index = 0;
  
  for (int col = 0; col < x; col++) {
    for (int row = 0; row < ROW; row++) {
      index++;
    }
  }
  if (x % 2 == 1){
    index += ROW - y - 1;
  }
  else 
    index += y;
  Serial.println(index);
  return index;
}

// Установка пикселя по координатам и цвету
void set_pixel(int x, int y, CRGB color) {
  if (is_correct_boundaries(x, y) && is_correct_color(color)){
    int index = get_index(x, y);
    leds[index] = color;
    FastLED.show();
  }
  else {
    Serial.println("Incorrect coordinates or color");
  }
}


void setup() {
  Serial.begin(9600);

  FastLED.addLeds<WS2812, LED_PIN, GRB>(leds, LED_NUM);
  FastLED.setBrightness(50);
  
// Очистка ленты вручную
  for (int i = 0; i < LED_NUM; i++) {
    leds[i] = CRGB::Black;
  }
  FastLED.show();

  //first parameter is name of access point, second is the password
  wifiManager.autoConnect("AP-NAME2");

  server.begin();

}

void loop() {
  // Варианты использования
  // set_pixel(0, 11, CRGB(15, 15, 15)); 
  // set_pixel(1, 5, CRGB::Red);

  if (Serial.available() > 0) {
    Serial.println("_______________________");
    int x = Serial.parseInt();  
    int y = Serial.parseInt(); 
    int r = Serial.parseInt();   
    int g = Serial.parseInt();   
    int b = Serial.parseInt();   
    
    while (Serial.available() > 0) {
      Serial.read();
    }

    Serial.print(x);
    Serial.print(" ");
    Serial.print(y);
    Serial.print(" RGB ");
    Serial.print(r);
    Serial.print(" ");
    Serial.print(g);
    Serial.print(" ");
    Serial.println(b);

    CRGB color = CRGB(r, g, b);

    set_pixel(x, y, color);
  }
}

