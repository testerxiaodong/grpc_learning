package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc_learning/client"
	"grpc_learning/pb"
	"grpc_learning/sample"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const (
	username        = "admin"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}
	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCore:  4,
		MinCpuGhz:   2.5,
		MinRam:      &pb.Memory{Value: 8, Unit: pb.Memory_GIGABYTE},
	}
	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.GetId(), "tmp/laptop.jpg")
}

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)

	for i := 0; i < n; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}
	scores := make([]float64, n)
	for {
		fmt.Println("rate laptop (y/n)? ")
		var answer string
		_, err := fmt.Scan(&answer)
		if err != nil {
			break
		}
		if strings.ToLower(answer) != "y" {
			break
		}
		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}
		err = laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func authMethods() map[string]bool {
	const laptopServicePath = "/pb.LaptopService/"
	return map[string]bool{
		laptopServicePath + "CreateLaptop": true,
		laptopServicePath + "UploadImage":  true,
		laptopServicePath + "RateLaptop":   true,
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := ioutil.ReadFile("cert/ca-cert.pem")
	if err != nil {
		return nil, fmt.Errorf("cannot read ca cert: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, errors.New("cannot append ca pem to cert pool")
	}

	// load server's certificate and private key
	clientCert, err := tls.LoadX509KeyPair("cert/client-cert.pem", "cert/client-key.pem")
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal(err)
	}
	// 连接rpc服务端
	cc1, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		panic(err)
	}
	// 生成认证客户端
	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal(err)
	}
	// 注册认证客户端的拦截器
	cc2, err := grpc.Dial(
		*serverAddress,
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		panic(err)
	}
	// 调用rpc方法
	laptopClient := client.NewLaptopClient(cc2)
	testRateLaptop(laptopClient)
}
