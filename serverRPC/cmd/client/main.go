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

func TryConnect() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient("grpc-server:8889", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("1) %s", err)
		return nil, err
	}
	return conn, nil
}

func main() {
	time.Sleep(time.Second * 10)
	//Создать соединение с gRPC сервером
	conn, err := grpc.NewClient("grpc-server:8889", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("1) %s", err)
		for {
			conn, err = TryConnect()
			if err != nil {
				fmt.Printf("2) %s", err)
				time.Sleep(time.Second * 1) // Добавляем задержку между попытками
			} else {
				fmt.Println("Соединение с сервером установлено")
				break
			}
		}
	}
	fmt.Println("Соединение с сервером установлено")
	defer conn.Close()

	//Создать контекст подключения

	// Инициализировать нового клиента
	client := grpcPb.NewSortServiceClient(conn)

	for {
		//Отправить данные на сервер
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		resp, err := client.Fetch(ctx, &grpcPb.FetchRequest{
			Url: "http://web-app:8085/products/",
		})

		if err != nil {
			fmt.Println("Fetch error")
			fmt.Printf("3) %s", err)
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
			fmt.Printf("4) %s", err)
		}
		fmt.Println("Сортировка от сервера: ", list)

		time.Sleep(1 * time.Second)
		cancel()
	}
}
