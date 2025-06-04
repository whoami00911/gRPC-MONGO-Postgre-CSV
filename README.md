## CRUD Приложение предоставляющее web API к данным, gRPC сервер считывающий данные по ссылке от клиента, добавляющий данные в mongoDB и выдающий отсортированный список клиенту (параметры сортировки задаются на клиенте)
### Стэк
- **Go**: go 1.23.5
- **Docker**

### Начало работы

Для запуска программы используйте следующую команду:

```bash
cd webAppCSV
make start
```
```bash
cd serverRPC
make start
```
### Использование web application

- **(Postman)Добавить продукт:**
Добавить сущность с уже существующим IP нельзя ни в кеш ни в БД
POST localhost:8085/products/
```
{
  "name": "name",
  "price": "price(decimal)"
}
```

- **(Postman)Получить все продукты из БД:**
```
GET localhost:8085/products/
```

- **(Postman)Получить один продукт из БД:**

```
GET localhost:8085/products/<id integer>
```

- **(Postman)Удалить все продукты**
```
DELETE localhost:8085/products/
```

- **(Postman)Удалить продукт по id:**
```
DELETE localhost:8085/products/<id integer>
```

- **(Postman)Обновить сущность по IP**
PUT localhost:8085/products/
```
{
  "name": "name",
  "price": "price(decimal)"
}
```

### Использование клиента
```
	list, err := client.List(ctx, &grpcPb.ListRequest{
		SortField:    grpcPb.ListRequest_name, // выбор поля для сортировки
		SortAsc:      1, // -1 по уменьшению, 1 по возврастани.
		PagingOffset: 0, // пропустить колличество записей
		PagingLimit:  10, // лимит записей
	})
```
