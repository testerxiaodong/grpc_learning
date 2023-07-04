package sample

import (
	"github.com/google/uuid"
	"grpc_learning/pb"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomCpuBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

func randomStringFromSet(a ...string) string {
	if len(a) == 0 {
		return ""
	}
	return a[rand.Intn(len(a))]
}

func randomCpuName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet(
			"Xeon E-2286M",
			"Core i9-9980Hk",
			"Core i7-9750H",
			"Core i5-9400F",
			"Core i3-1005G1")
	}
	return randomStringFromSet(
		"Ryzen 7 PRO 2700U",
		"Ryzen 5 PRO 3500U",
		"Ryzen 3 PRO 3200GE",
	)
}

func randLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("MacBook Air", "MacBook Pro")
	case "Dell":
		return randomStringFromSet("Latitude", "Vostro", "XPS", "Alienware")
	default:
		return randomStringFromSet("Thinkpad X1", "Thinkpad P1", "Thinkpad P53")
	}
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomId() string {
	return uuid.NewString()
}
