package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/thoas/go-funk"

	pb "github.com/DanielKnourek/CTC-tasks/trunk/task05/grpc/proto"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type serverImpl struct {
	pb.UnimplementedSEtcdServer
}

var etcd_client *clientv3.Client

func (s *serverImpl) HelloWorld(ctx context.Context, in *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloWorldResponse{Greeting: fmt.Sprintf("Hello %s", in.Name)}, nil
}

func (s *serverImpl) ListProducts(ctx context.Context, in *pb.Empty) (*pb.ProductListResponse, error) {
	log.Printf("Returning list of products")
	resp, err := etcd_client.Get(ctx, "product/", clientv3.WithPrefix())
	if err != nil {
		// handle error!
		log.Fatal(err)
	}
	// use the response
	// fmt.Printf("Get is done. Header is %v \n", resp.Kvs)

	products := funk.Map(resp.Kvs, func(kv *mvccpb.KeyValue) *pb.Product {
		amnt, err := strconv.ParseInt(string(kv.Value), 10, 32)
		if err != nil {
			log.Fatal(err)
		}
		return &pb.Product{
			Name:    strings.TrimPrefix(string(kv.Key), "product/"),
			Ammount: int32(amnt),
		}
	}).([]*pb.Product)
	return &pb.ProductListResponse{
		Products: products,
	}, nil
}

func (s *serverImpl) GetProduct(ctx context.Context, in *pb.ProductGetRequest) (*pb.ProductGetResponse, error) {
	log.Printf("Returning product %s", in.Id)
	resp, err := etcd_client.Get(ctx, fmt.Sprintf("product/%s", in.Id))
	if err != nil {
		// handle error!
		log.Fatal(err)
	}
	// use the response
	fmt.Printf("Get is done. Header is %v \n", resp.Kvs)

	amnt, err := strconv.ParseInt(string(resp.Kvs[0].Value), 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return &pb.ProductGetResponse{Product: &pb.Product{Name: in.Id, Ammount: int32(amnt)}}, nil
}

func (s *serverImpl) CreateProduct(ctx context.Context, in *pb.ProductPostRequest) (*pb.ProductPostResponse, error) {
	log.Printf("Creating product %s with ammount %d", in.Name, in.Ammount)
	resp, err := etcd_client.Put(ctx, fmt.Sprintf("product/%s", in.Name), fmt.Sprintf("%d", in.Ammount))
	if err != nil {
		// handle error!
		log.Fatal(err)
	}
	// use the response
	fmt.Printf("Put is done. Header is %v \n", resp)
	return &pb.ProductPostResponse{Product: &pb.Product{Name: in.Name, Ammount: in.Ammount}}, nil
}

func (s *serverImpl) DeleteProduct(ctx context.Context, in *pb.ProductDeleteRequest) (*pb.ProductDeleteResponse, error) {
	resp, err := etcd_client.Delete(ctx, fmt.Sprintf("product/%s", in.Id))
	if err != nil {
		// handle error!
		log.Fatal(err)
		return nil, err
	}
	// use the response
	fmt.Printf("Delete is done. Header is %v \n", resp)
	return &pb.ProductDeleteResponse{Id: in.Id}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func initEtcdConn() {
	ETCD_ENDPOINT_LIST := strings.Split(getEnv("ETCD_ENDPOINTS", "db_etcd:2379"), ",")
	log.Printf("Connecting to etcd on %v\n", ETCD_ENDPOINT_LIST)
	etcd_conn, err := clientv3.New(clientv3.Config{
		Endpoints:   ETCD_ENDPOINT_LIST,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("Error connecting to etcd: %v", err)
	}

	etcd_client = etcd_conn
}

func main() {
	GRPC_PORT := getEnv("GRPC_PORT", "50051")

	initEtcdConn()
	defer etcd_client.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPC_PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterSEtcdServer(server, &serverImpl{})
	log.Printf("Starting grpc server on port %s\n", GRPC_PORT)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
