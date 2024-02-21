unsigned char secsLeft = 10;
bool on = true;

void setup() {
  pinMode(10, OUTPUT);
  
  Serial.begin(9600);
  Serial.println("Pomodoro Clock Started");
}

void loop() {
  digitalWrite(7, on);
  if (on) { on = --secsLeft != 0; }

  Serial.println(secsLeft);
  delay(1000);
}
