package service

import (
	"bytes"
	"context"
	"errors"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc_learning/pb"
	"io"
	"log"
)

const maxImageSize = 1 << 20

type LaptopService struct {
	pb.UnimplementedLaptopServiceServer
	LaptopStore LaptopStore
	ImageStore  ImageStore
	RateStore   RateStore
}

func (l *LaptopService) CreateLaptop(ctx context.Context, request *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	// 获取请求对象中的字段做业务处理
	laptop := request.GetLaptop()
	log.Printf("receive a create-laptop request with id: %s", laptop.Id)

	if len(laptop.Id) > 0 {
		_, err := uuid.Parse(laptop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not valid UUID: %v", err)
		}
	} else {
		id := uuid.NewString()
		laptop.Id = id
	}
	// 超时控制与客户端取消处理，可抽离成一个函数进行代码复用
	//time.Sleep(7 * time.Second)
	err := ContextError(ctx)
	if err != nil {
		return nil, err
	}
	// save the laptop to store
	err = l.LaptopStore.Save(laptop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}
	log.Printf("saved laptop with id: %s", laptop.Id)
	return &pb.CreateLaptopResponse{
		Id: laptop.GetId(),
	}, nil
}

func (l *LaptopService) SearchLaptop(req *pb.SearchLaptopRequest, stream pb.LaptopService_SearchLaptopServer) error {
	filter := req.GetFilter()
	log.Printf("receive a search-laptop request with filter: %v", filter)
	err := l.LaptopStore.Search(stream.Context(), filter, func(laptop *pb.Laptop) error {
		res := &pb.SearchLaptopResponse{Laptop: laptop}
		err := stream.Send(res)
		if err != nil {
			return err
		}
		log.Printf("sent laptop with id: %s", laptop.GetId())
		return nil
	})
	if err != nil {
		return status.Errorf(codes.Internal, "unexpected error: %v", err)
	}
	return nil
}

func (l *LaptopService) UploadImage(stream pb.LaptopService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		return logError(status.Error(codes.Unknown, "cannot receive image info"))
	}
	laptopId := req.GetInfo().GetLaptopId()
	imageType := req.GetInfo().GetImageType()
	imageSize := 0
	imageData := bytes.Buffer{}
	for {
		if err := ContextError(stream.Context()); err != nil {
			return err
		}
		log.Println("waiting to receive more data")
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return logError(status.Error(codes.Unknown, "cannot receive image data"))
		}
		chunk := res.GetChunkData()
		size := len(chunk)
		log.Printf("received a chunk with size: %d", size)
		imageSize += size

		if imageSize > maxImageSize {
			return logError(status.Error(codes.InvalidArgument, "image too large"))
		}
		_, err = imageData.Write(chunk)
		if err != nil {
			return logError(status.Error(codes.Internal, "cannot write image data to byte buffer"))
		}
	}
	imageId, err := l.ImageStore.Save(laptopId, imageType, imageData)
	if err != nil {
		return logError(status.Error(codes.Internal, "cannot save image to imageStore"))
	}
	res := &pb.UploadImageResponse{Id: imageId, Size: uint32(imageSize)}
	err = stream.SendAndClose(res)
	if err != nil {
		return logError(status.Error(codes.Internal, "send res err"))
	}
	log.Printf("saved image with id: %s, size: %d", imageId, imageSize)
	return nil
}

func (l *LaptopService) RateLaptop(stream pb.LaptopService_RateLaptopServer) error {
	for {
		ctx := stream.Context()
		if err := ContextError(ctx); err != nil {
			return err
		}
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return logError(status.Error(codes.Unknown, "cannot receive client laptop rate message"))
		}
		laptopID := req.GetLaptopId()
		score := req.GetScore()
		log.Printf("receive a rate-laptop request: id = %s, score = %.2f", laptopID, score)
		// 查看laptop是否在laptopStore中
		find, err := l.LaptopStore.Find(laptopID)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot find laptop with id: %s", laptopID))
		}
		if find == nil {
			return status.Errorf(codes.NotFound, "laptop dont exist with id: %s", laptopID)
		}
		// 加入rateStore
		rate, err := l.RateStore.Add(laptopID, score)
		if err != nil {
			return logError(status.Errorf(codes.Internal, "cannot add to rateStore with id: %s", laptopID))
		}
		res := &pb.RateLaptopResponse{
			LaptopId:     laptopID,
			RateCount:    rate.Count,
			AverageScore: rate.SumScore / float64(rate.Count),
		}
		err = stream.Send(res)
		if err != nil {
			return logError(status.Errorf(codes.Unknown, "cannot send laptop rate response with id: %s", laptopID))
		}
	}
	return nil
}

func NewLaptopService(laptopStore LaptopStore, imageStore ImageStore, rateStore RateStore) *LaptopService {
	return &LaptopService{
		LaptopStore: laptopStore,
		ImageStore:  imageStore,
		RateStore:   rateStore,
	}
}

func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

func ContextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return logError(status.Error(codes.Canceled, "request is canceled"))
	case context.DeadlineExceeded:
		return logError(status.Error(codes.DeadlineExceeded, "deadline is exceeded"))
	default:
		return nil
	}
}
