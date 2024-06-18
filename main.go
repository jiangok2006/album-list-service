package main

import (
	"context"
	"log"
	"net/http"

	"net"

	gwpb "buf.build/gen/go/jiangok/buf-hello/grpc-ecosystem/gateway/v2/album_list_service/v1/album_list_service/album_list_servicev1gateway"
	gpb "buf.build/gen/go/jiangok/buf-hello/grpc/go/album_list_service/v1/album_list_servicev1grpc"
	pb "buf.build/gen/go/jiangok/buf-hello/protocolbuffers/go/album_list_service/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	gpb.UnimplementedAlbumListServiceServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) GetAlbumList(ctx context.Context, in *pb.GetAlbumListRequest) (*pb.GetAlbumListResponse, error) {
	log.Printf("Received request")
	return &pb.GetAlbumListResponse{Id: "100", Title: "mytitle", Price: 10.00}, nil
}

// curl -k http://localhost:8090/v1/album_list_service
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	gpb.RegisterAlbumListServiceServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	conn, err := grpc.NewClient(
		"0.0.0.0:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = gwpb.RegisterAlbumListServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
