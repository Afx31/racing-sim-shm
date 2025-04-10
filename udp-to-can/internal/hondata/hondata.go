package hondata

type Frame660 struct {
	Rpm					uint16
	Speed				uint16
	Gear				uint8
	Voltage			uint8
}

type Frame661 struct {
	Iat					uint16
	Ect					uint16
	Mil					uint8
	Vts					uint8
	Cl					uint8
}

type Frame662 struct {
	Tps					uint16
	Map					uint16
}

type Frame663 struct {
	Inj					float64
	Ign					uint16
}

type Frame664 struct {
	LambdaRatio uint16
}

type Frame665 struct {
	Knock				uint16
}

type Frame666 struct {
	TargetCamAngle	float64
	ActualCamAngle	float64
}

type Frame667 struct {
	OilTemp			uint16
	OilPressure	uint16
	//Analog2			uint16
	//Analog3			uint16
}

// type Frame668 struct {
// 	Analog4			uint16
// 	Analog5			uint16
// 	Analog6			uint16
// 	Analog7			uint16
// }

type Frame669S300 struct {
	Frequency			uint8
	Duty					float64
	Content				float64
}
type Frame669KPRO struct {
	Frequency				uint8
	EthanolContent	float64
	FuelTemperature	uint16
}