package main

import (
	"io"
	"io/fs"
	"net"
	"os"

	filestream "github.com/SaKu2110/grpc/proto/gen/go/v1"
	"google.golang.org/grpc"
)

const (
	FILE_PERMISSION fs.FileMode = 0644
)

var _ filestream.FileServiceServer = &Server{}

type Server struct {
	filestream.UnimplementedFileServiceServer
}

func (s *Server) Upload(stream filestream.FileService_UploadServer) error {
	var file *os.File = nil

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if file == nil {
			if file, err = os.OpenFile(resp.FilePath, os.O_RDWR|os.O_CREATE, FILE_PERMISSION); err != nil {
				return err
			}
			defer file.Close()
		}
		file.Write(resp.FileData)
	}

	if err := stream.SendAndClose(&filestream.UploadResponse{
		Message: "Complate!",
	}); err != nil {
		return err
	}
	return nil
}

func main() {
	port := ":50051"
	serv := grpc.NewServer()
	filestream.RegisterFileServiceServer(serv, &Server{})

	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	if err := serv.Serve(lis); err != nil {
		panic(err)
	}
}
