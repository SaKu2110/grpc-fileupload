package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	filestream "github.com/SaKu2110/grpc/proto/gen/go/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type terget struct {
	host string
	path string
}

func request(src, dest string, client filestream.FileServiceClient) error {
	stream, err := client.Upload(context.Background())
	if err != nil {
		return err
	}

	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 1024)
	for {
		_, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		stream.Send(&filestream.UploadRequest{FilePath: dest, FileData: buf})
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)
	return nil
}

func main() {
	flag.Parse()
	var src, dest string = flag.Arg(0), flag.Arg(1)
	var t terget
	if src == "" || dest == "" {
		os.Exit(0)
	}

	port := ":50051"

	if arr := strings.Split(dest, ":"); len(arr) < 2 {
		os.Exit(0)
	} else {
		t = terget{host: arr[0], path: arr[1]}
	}

	conn, err := grpc.Dial(t.host+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := filestream.NewFileServiceClient(conn)

	if err := request(src, t.path, client); err != nil {
		panic(err)
	}
}
