package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"grpc_learning/pb"
	"grpc_learning/service"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	secretKey      = "secretKey"
	tokenDuration  = 15 * time.Minute
	serverCertFile = "cert/server-cert.pem"
	serverKeyFile  = "cert/server-key.pem"
	caCertFile     = "cert/ca-cert.pem"
)

func seedUsers(userStore service.UserStore) error {
	err := createUser(userStore, "admin", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")
}

func createUser(userStore service.UserStore, username string, password string, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func accessibleRoles() map[string][]string {
	const laptopServicePath = "/pb.LaptopService/"
	return map[string][]string{
		laptopServicePath + "CreateLaptop": {"admin"},
		laptopServicePath + "UploadImage":  {"admin"},
		laptopServicePath + "RateLaptop":   {"admin", "user"},
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		return nil, err
	}

	pemClientCA, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read ca cert: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, errors.New("cannot append ca pem to cert pool")
	}

	// create credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	return credentials.NewTLS(config), nil
}

func runGRPCServer(
	authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	jwtManager *service.JwtManager,
	enableTLS bool,
	listener net.Listener,
) error {
	// 身份验证拦截器
	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	// grpc选项
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}
	// TLS验证
	if enableTLS {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			return err
		}
		serverOptions = append(serverOptions, grpc.Creds(tlsCredentials))
	}
	// 注册拦截器：一元拦截器和流拦截器
	grpcServer := grpc.NewServer(serverOptions...)
	// 注册服务
	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)
	// 开启grpc反射
	reflection.Register(grpcServer)
	// 启动服务
	err := grpcServer.Serve(listener)
	if err != nil {
		return err
	}
	log.Printf("Start GRPC server on port %s, TLS = %t", listener.Addr().String(), enableTLS)
	return nil
}

func runRESTServer(
	authServer pb.AuthServiceServer,
	laptopServer pb.LaptopServiceServer,
	enableTLS bool,
	listener net.Listener,
) error {
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authServer)
	if err != nil {
		return err
	}
	err = pb.RegisterLaptopServiceHandlerServer(ctx, mux, laptopServer)
	if err != nil {
		return err
	}
	log.Printf("Start REST server on port %s, TLS = %t", listener.Addr().String(), enableTLS)
	if enableTLS {
		return http.ServeTLS(listener, mux, serverCertFile, serverKeyFile)
	}
	return http.Serve(listener, mux)
}

func main() {
	fmt.Println("hello world from server...")
	port := flag.Int("port", 0, "the server port")
	flag.Parse()
	log.Printf("start server on port %d", *port)

	userStore := service.NewInMemoryUserStore()
	err := seedUsers(userStore)
	if err != nil {
		log.Fatal("cannot seed user")
	}

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	jwtManager := service.NewJwtManager(secretKey, tokenDuration)
	authService := service.NewAuthService(userStore, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("img")
	rateStore := service.NewInMemoryRateStore()

	laptopServer := service.NewLaptopService(laptopStore, imageStore, rateStore)

	err = runRESTServer(authService, laptopServer, false, listener)
	if err != nil {
		panic(err)
	}
}
