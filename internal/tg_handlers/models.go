package tg_handlers

const (
	characteristicsAnswerBody = `
Engine: %s,	
Drive Type: %s,	
Power: %s,	
Acceleration: %s,	
Top Speed: %s,	
Fuel Economy: %s,	
Transmission: %s,
						`
	characteristicsRequest = "опиши мне главные характеристики машины %s в виде одного json на английском, с полями: engine, power, acceleration, top_speed, fuel_economy, transmission, drive_type"
)
