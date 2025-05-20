package main

import (
	"context"
	"fmt"
	"gRPC-server/pkg/parseCSV/grpcPb"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//Создать соединение с gRPC сервером
	conn, err := grpc.NewClient("localhost:8889", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("1) %s", err)
	}
	defer conn.Close()

	//Создать контекст подключения
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	defer cancel()

	// Инициализировать нового клиента
	client := grpcPb.NewSortServiceClient(conn)

	//Отправить данные на сервер
	resp, err := client.Fetch(ctx, &grpcPb.FetchRequest{
		Url: "http://localhost:8085/products/",
	})

	if err != nil {
		log.Fatalf("2) %s", err)
	}

	log.Printf("Fetch success: %s", resp.Status)

	//отправить запрос на получение данных
	list, err := client.List(ctx, &grpcPb.ListRequest{
		SortField:    grpcPb.ListRequest_name,
		SortAsc:      1,
		PagingOffset: 0,
		PagingLimit:  10,
	})
	if err != nil {
		log.Fatalf("3) %s", err)
	}
	fmt.Println("Сортировка от сервера: ", list)
}
