package cwgrpc

import (
	"CW_DB_v2/internal/domain/models"
	"CW_DB_v2/internal/services/cw"
	"context"
	"errors"
	cwv1 "github.com/ykwais/CW_GO_protos/gen/go/cw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type CW interface {
	Login(ctx context.Context, login, password string) (userId int64, err error)
	Register(ctx context.Context, login, password, email, realName string) (userID int64, err error)

	ListPhotos() ([]models.Photo, error)
	PhotosForMainScreen(ctx context.Context, dateStart string, dateEnd string) (photos []models.BetterPhoto, err error)
	PhotosOfAutomobile(id int64) (photos []models.Photo, err error)
}

type serverAPI struct {
	cwv1.UnimplementedServiceServer
	*slog.Logger
	cw CW
}

func RegisterServerAPI(gRPC *grpc.Server, logger *slog.Logger, service CW) {
	cwv1.RegisterServiceServer(gRPC, &serverAPI{
		UnimplementedServiceServer: cwv1.UnimplementedServiceServer{},
		Logger:                     logger,
		cw:                         service,
	}) // тут происходит связывание сервера и его реализации

}

/*
	ниже представлены обработчики запросов, поступающие на сервер. Каждый из обработчиков отвечает за валидацию, вызов реализации метода - Login например
	и посылку ответа на клиент
*/

func (s *serverAPI) PhotosOfAutomobile(req *cwv1.PhotosOfAutomobileRequest, res grpc.ServerStreamingServer[cwv1.PhotosOfAutomobileResponse]) error {
	s.Logger.Info("photos of current automobile start")
	combinePathData, err := s.cw.PhotosOfAutomobile(req.Id)
	if err != nil {
		return err
	}

	for _, photo := range combinePathData {
		chunkSize := 1024 * 1024
		data := photo.Data
		for i := 0; i < len(data); i += chunkSize {
			end := i + chunkSize
			if end > len(data) {
				end = len(data)
			}

			response := &cwv1.PhotosOfAutomobileResponse{
				PhotoName: photo.Name,
				Chunk:     data[i:end],
			}

			if err := res.Send(response); err != nil {
				s.Logger.Error("failed to send photo chunk", err)
				return err
			}
		}
	}

	return nil
}

func (s *serverAPI) PhotosForMainScreen(req *cwv1.PhotosForMainScreenRequest, res grpc.ServerStreamingServer[cwv1.PhotosForMainScreenResponse]) error {
	s.Logger.Info("start PhotosForMainScreen")
	combinePhotoDatas, err := s.cw.PhotosForMainScreen(context.Background(), req.DateBegin, req.DateEnd)
	if err != nil {
		s.Logger.Error("failed to get photos for main screen", err)
		return err
	}

	for _, photo := range combinePhotoDatas {
		chunkSize := 1024 * 1024
		data := photo.Data
		for i := 0; i < len(data); i += chunkSize {
			end := i + chunkSize
			if end > len(data) {
				end = len(data)
			}

			response := &cwv1.PhotosForMainScreenResponse{
				Chunk:       data[i:end],
				Brand:       photo.Brand,
				Model:       photo.Model,
				VehicleId:   photo.VehicleId,
				PricePerDay: photo.TotalCost,
			}

			if err := res.Send(response); err != nil {
				s.Logger.Error("failed to send photo chunk", err)
				return err
			}
		}
	}

	s.Logger.Info("end PhotosForMainScreen")
	return nil

}

func (s *serverAPI) ListPhotos(req *cwv1.EmptyRequest, stream cwv1.Service_ListPhotosServer) error {
	s.Logger.Info("start ListPhotos")
	photos, err := s.cw.ListPhotos()
	if err != nil {
		s.Logger.Error("failed to list photos", err)
		return err
	}

	for _, photo := range photos {

		chunkSize := 1024 * 1024
		data := photo.Data
		for i := 0; i < len(data); i += chunkSize {
			end := i + chunkSize
			if end > len(data) {
				end = len(data)
			}

			response := &cwv1.ListPhotosResponse{
				PhotoName: photo.Name,
				Chunk:     data[i:end],
			}

			if err := stream.Send(response); err != nil {
				s.Logger.Error("failed to send photo chunk", err)
				return err
			}
		}
	}

	//for _, photo := range photos {
	//	response := &cwv1.ListPhotosResponse{
	//		PhotoName: photo.Name,
	//		Chunk:     photo.Data,
	//	}
	//
	//	if err := stream.Send(response); err != nil {
	//		s.Logger.Error("failed to send photo", err)
	//		return err
	//	}
	//}

	s.Logger.Info("all photos sent successfully")
	return nil
}

func (s *serverAPI) Login(ctx context.Context, req *cwv1.LoginRequest) (*cwv1.LoginResponse, error) {
	s.Logger.Info("on LOGIN request get: ", slog.String("login", req.Login), slog.String("password", req.Password))
	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty login or password")
	}

	userID, err := s.cw.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, cw.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid login or password")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &cwv1.LoginResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *cwv1.RegisterRequest) (*cwv1.RegisterResponse, error) {

	if req.GetLogin() == "" || req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty login or password")
	}

	resultCh := make(chan *cwv1.RegisterResponse, 1)

	go func() {
		userID, err := s.cw.Register(ctx, req.GetLogin(), req.GetPassword(), req.GetEmail(), req.GetRealName())
		if err != nil {
			resultCh <- nil
			return
		}
		resultCh <- &cwv1.RegisterResponse{
			UserId: userID,
		}
	}()

	registerResponse := <-resultCh

	if registerResponse == nil {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	return registerResponse, nil

	//userID, err := s.cw.Register(ctx, req.GetLogin(), req.GetPassword())
	//if err != nil {
	//	if errors.Is(err, cw.ErrUserExists) {
	//		return nil, status.Error(codes.AlreadyExists, "user already exists")
	//	}
	//	return nil, status.Error(codes.Internal, "internal error")
	//}

	//return &cwv1.RegisterResponse{
	//	UserId: userID,
	//}, nil
}

//func (s *serverAPI) isAdmin(ctx context.Context, req *cwv1.IsAdminRequest) (*cwv1.IsAdminResponse, error) {
//	if req.GetUserId() == 0 {
//		return nil, status.Error(codes.InvalidArgument, "user ID cannot be 0")
//	}
//
//	isAdmin, err := s.cw.IsAdmin(ctx, req.GetUserId())
//	if err != nil {
//		if errors.Is(err, storage.ErrUserNotFound) {
//			return nil, status.Error(codes.NotFound, "user not found")
//		}
//		return nil, status.Error(codes.Internal, "internal error")
//	}
//
//	return &cwv1.IsAdminResponse{IsAdmin: isAdmin}, nil
//}
