package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/bobbydeveaux/go-blog-frontend/app/common"
	pb "github.com/bobbydeveaux/go-blog-proto/post"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func getPosts(w http.ResponseWriter, r *http.Request) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.

	rpc, err := c.GetPosts(context.Background(), &pb.PostRequest{Name: ""})
	if err != nil {
		fmt.Println("could not greet: %v", err)
	}

	b, err := json.Marshal(rpc.Message)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

func main() {
	flag.Parse()
	router := mux.NewRouter()
	http.Handle("/", httpInterceptor(router))

	router.HandleFunc("/posts", getPosts).Methods("GET")

	http.ListenAndServe(":8181", nil)
}

func httpInterceptor(router http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		if origin := req.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here if its Preflighted OPTIONS request
		if req.Method == "OPTIONS" {
			return
		}
		startTime := time.Now()

		router.ServeHTTP(w, req)

		finishTime := time.Now()
		elapsedTime := finishTime.Sub(startTime)

		switch req.Method {
		case "GET":
			// We may not always want to StatusOK, but for the sake of
			// this example we will
			common.LogAccess(w, req, elapsedTime)
		case "POST":
			// here we might use http.StatusCreated
		}

	})
}
