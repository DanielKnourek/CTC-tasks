package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"

	pb "github.com/DanielKnourek/CTC-tasks/trunk/task05/grpc/proto"
	"google.golang.org/grpc"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

var conn *grpc.ClientConn

func initConnection() {
	GRPC_PORT := getEnv("GRPC_PORT", "50051")
	GRPC_ADDRESS := getEnv("GRPC_ADDRESS", "localhost")

	connection, err := grpc.Dial(fmt.Sprintf("%s:%s", GRPC_ADDRESS, GRPC_PORT), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	conn = connection
}

func hello(w http.ResponseWriter, _ *http.Request) {
	ctx, ctx_cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctx_cancel()
	client := pb.NewSEtcdClient(conn)
	res, err := client.HelloWorld(ctx, &pb.HelloWorldRequest{Name: "Daniel"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return
	}
	log.Printf("Greeting: %s", res.Greeting)

	data, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("could not marshal: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}

func products(w http.ResponseWriter, _ *http.Request) {
	ctx, ctx_cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctx_cancel()
	client := pb.NewSEtcdClient(conn)
	res, err := client.ListProducts(ctx, &pb.Empty{})
	if err != nil {
		log.Fatalf("could not list products: %v", err)
		return
	}
	log.Printf("Products: %s", res.Products)

	data, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("could not marshal: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	ctx, ctx_cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctx_cancel()
	client := pb.NewSEtcdClient(conn)
	res, err := client.GetProduct(ctx, &pb.ProductGetRequest{Id: productId})
	if err != nil {
		log.Fatalf("could not get product: %v", err)
		return
	}
	log.Printf("Product: %s", res.Product)

	data, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("could not marshal: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}

func delProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productId := vars["id"]

	ctx, ctx_cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctx_cancel()
	client := pb.NewSEtcdClient(conn)
	res, err := client.DeleteProduct(ctx, &pb.ProductDeleteRequest{Id: productId})
	if err != nil {
		log.Fatalf("could not delete product: %v", err)
		return
	}
	log.Printf("Product: %s", res.Id)

	data, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("could not marshal: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}

func putProduct(w http.ResponseWriter, r *http.Request) {
	ctx, ctx_cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer ctx_cancel()

	var Product_req pb.Product
	err := json.NewDecoder(r.Body).Decode(&Product_req)
	if err != nil {
		log.Fatalf("could not decode: %v", err)
		return
	}

	client := pb.NewSEtcdClient(conn)
	res, err := client.CreateProduct(ctx, &pb.ProductPostRequest{Name: Product_req.Name, Ammount: Product_req.Ammount})
	if err != nil {
		log.Fatalf("could not put product: %v", err)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("could not marshal: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", data)
}

func main() {
	initConnection()
	defer conn.Close()

	CLIENT_PORT := getEnv("CLIENT_PORT", "8080")

	// ctx, ctx_cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer ctx_cancel()

	router := mux.NewRouter()

	router_api := router.PathPrefix("/api/").Subrouter()
	router_api.HandleFunc("/hello", hello).Methods(http.MethodGet, http.MethodPost)
	router_api.HandleFunc("/product", products).Methods(http.MethodGet)
	router_api.HandleFunc("/product", putProduct).Methods(http.MethodPut)

	router_prod := router_api.PathPrefix("/product/").Subrouter()
	router_prod.HandleFunc("/{id}", getProduct).Methods(http.MethodGet)
	router_prod.HandleFunc("/{id}", delProduct).Methods(http.MethodDelete)
	router_prod.NotFoundHandler = http.HandlerFunc(products)

	fileServer := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fileServer))

	panic(http.ListenAndServe(fmt.Sprintf(":%s", CLIENT_PORT), router))

	// client := pb.NewSEtcdClient(conn)

	// r, err := client.ListProducts(ctx, &pb.Empty{})
	// if err != nil {
	// 	log.Fatalf("could not list products: %v", err)
	// 	return
	// }
	// log.Printf("Products: %v", r.Products.Products)
}
