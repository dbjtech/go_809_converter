package terminal

var (
	FirmwareType = map[string]string{
		"F": "ZJ210W",
		"H": "ZJ220",
		"I": "ZJ220S",
		"E": "ZJ210",
		"K": "ZJ210B",
		"S": "ZJ210L",
		"T": "ZJ211W",
		"U": "ZJ211",
		"Y": "ZJ300",
		"V": "ZJ300WL",
		"W": "ZJ300L",
		"X": "ZJ300W",
		"J": "IV100",
	}

	ResFirmwareType = map[string]string{
		"ZJ210W":  "F",
		"ZJ220":   "H",
		"ZJ220S":  "I",
		"ZJ210":   "E",
		"ZJ210B":  "K",
		"ZJ210L":  "S",
		"ZJ211W":  "T",
		"ZJ211":   "U",
		"ZJ300":   "Y",
		"ZJ300WL": "V",
		"ZJ300L":  "W",
		"ZJ300W":  "X",
		"IV100":   "J",
	}

	SwitchStatus = map[string]int{
		"sent":     0,
		"sending":  1,
		"worked":   2,
		"unworked": 3,
		"downlink": -1,
		"ignore":   0,
		"success":  1,
		"failed":   2,
	}
)

var AlarmStatus = struct {
	LongStop   int
	LowVoltage int
	OverSpeed  int
}{
	// 长时间停留
	LongStop: 1 << 0,
	// 电瓶低电
	LowVoltage: 1 << 1,
	//超速
	OverSpeed: 1 << 2,
}
