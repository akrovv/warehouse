## Таска
```
Реализовать API для работы с товарами на складе
```
API работает на базе реляционной СУБД **PostgreSQL** + **jsonRPC**.

## Инструкция по запуску
```
Запуск приложения: make up  
Запуск тестов: make test
Запуск linter'a: make lint
```

## cURL  
Вместо cURL можно использовать - **curs.bash**.  

## Стек
Backend: Golang, PostgreSQL, jsonRPC, docker, docker-compose

## Документация

### База данных
Всего реализовано 3 таблицы: **warehouses**, **products**, **warehouse_products**.  

**warehouses** и **products** соединены отношением **MANY-TO-MANY** через таблицу **warehouse_products**.   

Вставка в таблицу **warehouse_products** осуществляется через процедуру **insertWarehouseProducts**.  

**insertWarehouseProducts**  
На вход: **warehouse_id** (integer), **product_code** (UUID)  

При операции INSERT или UPDATE в таблицу **warehouse_products** срабататывает триггер - **tr_wareproducts_availability**, который поднимает исключение, в случае, если склад недоступен.  

### Склад

### Создать склад - POST Warehouses.Create
Принимает на вход массив json с параметрами склада.  

**Параметры**  
* name (string) - наименование склада
* availability (bool) - доступность склада

Пример json:
```json
[
    {
        "name": "Warehouse 1",
        "availability": true,
    },
    {
        "name": "Warehouse 2",
        "availability": false,
    }
]
```

Пример возможного запроса:
```bash
curl -v \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Warehouses.Create", "params": [[
        {"name": "Warehouse 1", "availability": true}, 
        {"name": "Warehouse 2", "availability": false}
    ]]}' \
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id":1,
    "result":
    [
        {
            "name": "Warehouse 1",
            "availability": true},
        {
            "name": "Warehouse 2",
            "availability": false
        }
    ],
    "error": null
}
```

### Получение оставшихся товаров со склада - GET Warehouses.GetLeftOvers
Принимает на вход json с id склада.  

**Параметры**  
* warehouse_id (integer) - id склада

Пример json:
```json
{
    "warehouse_id": 1,
}
```

Пример возможного запроса:
```bash
curl -v \
    -X GET \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Warehouses.GetLeftOvers", 
    "params": [
        {"warehouse_id": 1}
    ]}' \
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id": 1,
    "result": [
        {
            "name": "Product 2",
            "size": "30x120",
            "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11",
            "quantity": 5
        },
        {
            "name": "Product 1",
            "size": "50x50",
            "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "quantity": 10
        }
    ],
    "error":null
}
```

Пример возможной ошибки:
```json
{
    "id": 1,
    "result": null,
    "error": "service.GetLeftOvers returned: sql: no rows in result set"
}
```

## Товар

### Создать товар - POST Products.Create
Принимает на вход массив json с параметрами товара.  

**Параметры**  
* name (string) - наименование товара
* size (string) - размер товара
* code (string) - уникальный код (uuid)
* quantity (integer) - количество

Пример json:
```json
[
    {
        "name": "Product 1",
        "size": "50x50",
        "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "quantity": 20
    },
    {
        "name": "Product 2",
        "size": "30x120",
        "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11",
        "quantity": 5
    }
]
```

Пример возможного запроса:
```bash
curl -v \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Products.Create", 
    "params": 
    [[
        {"name": "Product 1", "size": "50x50", "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "quantity": 20}, 
        {"name": "Product 2", "size": "30x120", "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11", "quantity": 5}
    ]]
    }' \
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id":1,
    "result":
    [
        {
            "name": "Product 1",
            "size": "50x50",
            "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "quantity": 20
        },
        {
            "name": "Product 2",
            "size": "30x120",
            "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11",
            "quantity": 5
        }
    ],
    "error": null
}
```

Пример возможной ошибки:
```json
{
    "id": 1,
    "result": null,
    "error": "all calls returned: db.Exec with command INSERT to products returned: pq: duplicate key value violates unique constraint \"products_code_key\""
}
```

### Зарезервировать товар - POST Products.Reserve
Принимает на вход массив json с параметрами резерва.  

**Параметры**  
* warehouse_id (string) - id склада
* quantity (integer) - количество
* code (string) - уникальный код (uuid)

Пример json:
```json
[
    {
        "warehouse_id": 1, 
        "quantity": 10, 
        "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" 
    },
    {
        "warehouse_id": 1, 
        "quantity": 5,
        "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11" 
    }
]
```

Пример возможного запроса:
```bash
curl -v \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Products.Reserve", 
    "params": 
    [[
        {"warehouse_id": 1, "quantity": 10, "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" },
        {"warehouse_id": 1, "quantity": 5, "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11" }
    ]]}' \
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id": 1,
    "result":
    [
        {
            "warehouse_id": 1,
            "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "quantity": 10,
            "status": "reserved"
        },
        {
            "warehouse_id": 1,
            "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11",
            "quantity": 5,
            "status": "reserved"
            }
    ],
    "error": null
}
```

Пример возможной ошибки:
```json
{
    "id": 1,
    "result": null,
    "error": "all calls returned: db.Exec with command UPDATE to warehouse_products returned: pq: new row for relation \"warehouse_products\" violates check constraint \"warehouse_products_available_quantity_check\""
}
```


### Отменить резерв товар - POST Products.CancelReservation
Принимает на вход массив json с параметрами отмены резерва.  

**Параметры**  
* warehouse_id (string) - id склада
* quantity (integer) - количество
* code (string) - уникальный код (uuid)

Пример json:
```json
[
    {
        "warehouse_id": 1, 
        "quantity": 10, 
        "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" 
    },
    {
        "warehouse_id": 1, 
        "quantity": 5, 
        "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11" 
    }
]
```

Пример возможного запроса:
```bash
curl -v \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Products.CancelReservation", 
    "params": [[
        {"warehouse_id": 1, "quantity": 10, "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" },
        {"warehouse_id": 1, "quantity": 5, "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11" }
        ]]}' \
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id": 1,
    "result": [
        {
            "warehouse_id": 1,
            "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "quantity": 10,
            "status":"canceled"
        },
        {
            "warehouse_id": 1,
            "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11",
            "quantity": 5,
            "status":"canceled"
        }
    ],
    "error":null
}
```

Пример возможной ошибки:
```json
{
    "id": 1,
    "result": null,
    "error": "all calls returned: db.Exec with command UPDATE to warehouse_products returned: pq: new row for relation \"warehouse_products\" violates check constraint \"warehouse_products_reserved_quantity_check\""
}
```

### Пополнить товар на складе - POST Products.Add
Принимает на вход массив json с параметрами добавления товара.  

**Параметры**  
* warehouse_id (string) - id склада
* quantity (integer) - количество
* code (string) - уникальный код (uuid)

Пример json:
```json
{
    "warehouse_id": 1, 
    "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
    "quantity": 10
}
```

Пример возможного запроса:
```bash
curl -v \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Products.Add", 
    "params": [[
        {"warehouse_id": 1, "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "quantity": 10}
        ]]}' \
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id": 1,
    "result":
    [
        {
            "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "quantity": 10,
            "warehouse_id": 1
        }
    ],
    "error": null
}
```

Пример возможной ошибки:
```json
{
    "id": 1,
    "result": null,
    "error": "all calls returned: db.Exec with command UPDATE to warehouse_products returned: pq: insertion in no available warehouse"
}
```

### Перевести товары на другой склад - POST Products.Transfer
Принимает на вход массив json с параметрами для перевоза.  

**Параметры**  
* warehouse_from_id (string) - id склада, с которого перевод осуществляется
* warehouse_to_id (string) - id склада куда перевод осуществляется
* code (string) - уникальный код (uuid)
* quantity (integer) - количество

Пример json:
```json
{
    "warehouse_from_id": 1,
    "warehouse_to_id": 2, 
    "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", 
    "quantity": 10
}
```

Пример возможного запроса:
```bash
curl -v \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Products.Transfer", 
    "params": [
        [{"warehouse_from_id": 1, "warehouse_to_id": 2, "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", "quantity": 10}
    ]]}'\
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id": 1,
    "result":
    [
        {
            "warehouse_from_id": 1,
            "warehouse_to_id": 2,
            "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "quantity": 10
        }
    ],
    "error": null
}
```

Пример возможной ошибки:
```json
{
    "id": 1,
    "result": null,
    "error": "all calls returned: not enough quantity: 10, in warehouse: 1. available: 0"
}
```

### Утилизировать товар - DELETE Products.Delete
Принимает на вход массив json с параметрами для утилизации.  

**Параметры**  
* code (string) - уникальный код (uuid)

Пример json:
```json
[
    {
        "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" 
    },
    {
        "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11" 
    }
]
```

Пример возможного запроса:
```bash
curl -v \
    -X DELETE \
    -H "Content-Type: application/json" \
    -d '{"jsonrpc":"2.0", "id": 1, "method": "Products.Delete", "params": 
    [[
        {"code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11" },
        {"code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11" }
    ]]}' \
    http://localhost:8080/
```

Пример успешного ответа:
```json
{
    "id":1,
    "result":
    [
        {
            "name": "Product 1",
            "size": "50x50",
            "code": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
            "quantity": 60
        },
        {
            "name": "Product 2",
            "size": "30x120",
            "code": "a0eebc49-9c0b-3ef8-bb6d-6bb9bd380a11",
            "quantity": 5
        }
    ],
    "error":null
}
```

Пример возможной ошибки:
```json
{
    "id": 1,
    "result": null,
    "error": "all calls returned: db.Exec with command DELETE to products returned: sql: no rows in result set"
}
```










