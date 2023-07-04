package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc_learning/pb"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LaptopClient struct {
	service pb.LaptopServiceClient
}

func NewLaptopClient(cc *grpc.ClientConn) *LaptopClient {
	laptopClient := pb.NewLaptopServiceClient(cc)
	return &LaptopClient{laptopClient}
}

func (laptopClient *LaptopClient) CreateLaptop(laptop *pb.Laptop) {
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := laptopClient.service.CreateLaptop(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Println("laptop already exist")
		} else {
			log.Fatal("cannot creat laptop: ", err)
		}
		return
	}
	log.Printf("created laptop with id: %s", res.Id)
}

func (laptopClient *LaptopClient) SearchLaptop(filter *pb.Filter) {
	log.Println("search filter: ", filter)
	// 客户端超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 构建请求
	req := &pb.SearchLaptopRequest{Filter: filter}
	// 发送请求，获取响应流
	stream, err := laptopClient.service.SearchLaptop(ctx, req)
	if err != nil {
		log.Fatal("cannot search laptop: ", err)
	}
	// 从响应流中获取响应，直到出错或者EOF
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response: ", err)
		}
		laptop := res.GetLaptop()
		log.Println("- found: ", laptop.GetId())
		log.Println("- found: ", laptop.GetBrand())
		log.Println("- found: ", laptop.GetName())
		log.Println("- found: ", laptop.GetCpu().GetNumberCores())
		log.Println("- found: ", laptop.GetCpu().GetMinGhz())
		log.Println("- found: ", laptop.GetMemory().GetValue(), laptop.GetMemory().GetUnit())
		log.Println("- found: ", laptop.GetPriceUsd())
	}
}

func (laptopClient *LaptopClient) UploadImage(laptopID string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := laptopClient.service.UploadImage(ctx)
	if err != nil {
		log.Fatal("cannot upload image: ", err)
	}
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptopID,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}
	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info: ", err, stream.RecvMsg(nil))
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}
		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send image chunk data: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot recv res: ", err)
	}
	log.Printf("image upload success with id: %s, size: %d", res.GetId(), res.GetSize())
}

func (laptopClient *LaptopClient) RateLaptop(laptopIDs []string, scores []float64) error {
	// 防止下标越界
	if len(laptopIDs) < len(scores) {
		return errors.New("laptopIDs less than scores")
	}
	// 客户端超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := laptopClient.service.RateLaptop(ctx)
	if err != nil {
		return fmt.Errorf("cannot rate laptop: %v", err)
	}
	waitResponse := make(chan error)
	// go routine to receive response
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Println("no more data response")
				waitResponse <- nil
			}
			if err != nil {
				waitResponse <- fmt.Errorf("cannot receive stream response: %v", err)
				return
			}
			log.Println("received response: ", res)
		}
	}()
	// send request
	for i, laptopID := range laptopIDs {
		req := &pb.RateLaptopRequest{LaptopId: laptopID, Score: scores[i]}
		err := stream.Send(req)
		if err != nil {
			return fmt.Errorf("cannot send rate-laptop with id: %s, err: %v - %v", laptopID, err, stream.RecvMsg(nil))
		}
		log.Println("sent request: ", req)
	}
	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("cannot close send: %v", err)
	}
	err = <-waitResponse
	return err
}
