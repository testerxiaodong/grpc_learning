package sample

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpc_learning/pb"
)

func NewKeyBoard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}

func NewCpu() *pb.Cpu {
	brand := randomCpuBrand()
	name := randomCpuName(brand)
	numberCores := randomInt(2, 8)
	numberThreads := randomInt(numberCores, 12)
	minGHz := randomFloat64(2.0, 3.5)
	maxGHz := randomFloat64(minGHz, 5.0)

	return &pb.Cpu{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreads),
		MinGhz:        minGHz,
		MaxGhz:        maxGHz,
	}
}

func NewGpu() *pb.Gpu {
	brand := randomCpuBrand()
	name := randomCpuName(brand)
	minGHz := randomFloat64(2.0, 3.5)
	maxGHz := randomFloat64(minGHz, 5.0)
	memory := &pb.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit:  pb.Memory_GIGABYTE,
	}
	return &pb.Gpu{
		Brand:  brand,
		Name:   name,
		MinGhz: minGHz,
		MaxGhz: maxGHz,
		Memory: memory,
	}
}

func NewRAM() *pb.Memory {
	return &pb.Memory{
		Value: uint64(randomInt(4, 64)),
		Unit:  pb.Memory_GIGABYTE,
	}
}

func NewSSD() *pb.Storage {
	return &pb.Storage{
		Drive: pb.Storage_SSD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(128, 1024)),
			Unit:  pb.Memory_GIGABYTE,
		},
	}
}

func NewHDD() *pb.Storage {
	return &pb.Storage{
		Drive: pb.Storage_HDD,
		Memory: &pb.Memory{
			Value: uint64(randomInt(1, 6)),
			Unit:  pb.Memory_TERABYTE,
		},
	}
}

func NewScreen() *pb.Screen {
	height := randomInt(1080, 4320)
	width := height * 16 / 9
	return &pb.Screen{
		SizeInch: randomFloat32(13, 17),
		Resolution: &pb.Screen_Resolution{
			Height: uint32(height),
			Width:  uint32(width),
		},
		Panel:      randomScreenPanel(),
		Multitouch: randomBool(),
	}
}

func NewLaptop() *pb.Laptop {
	brand := randLaptopBrand()
	name := randomLaptopName(brand)
	laptop := &pb.Laptop{
		Id:       randomId(),
		Brand:    brand,
		Name:     name,
		Cpu:      NewCpu(),
		Memory:   NewRAM(),
		Gpus:     []*pb.Gpu{NewGpu()},
		Storages: []*pb.Storage{NewSSD(), NewHDD()},
		Screen:   NewScreen(),
		KeyBoard: NewKeyBoard(),
		Weight: &pb.Laptop_WeightKg{
			WeightKg: randomFloat64(1.0, 3.0),
		},
		PriceUsd:    randomFloat64(1500, 3000),
		ReleaseYear: uint32(randomInt(2015, 2019)),
		TimeStamp:   timestamppb.Now(),
	}
	return laptop
}

func RandomLaptopScore() float64 {
	return randomFloat64(1, 10)
}
