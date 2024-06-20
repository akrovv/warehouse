#!/bin/bash

show_menu() {
    clear
    echo "1. Products.Create"
    echo "2. Warehouses.Create"
    echo "3. Warehouses.GetLeftOvers"
    echo "4. Products.Reserve"
    echo "5. Products.CancelReservation"
    echo "6. Products.Add"
    echo "7. Products.Transfer" 
    echo "8. Products.Delete"
    echo "9. Выйти"
}

create_warehouse() {
    read -p "Название склада: " warehouse_name
    read -p "Доступность (true/false): " availability

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Warehouses.Create", "params": [[{"name": "'"$warehouse_name"'", "availability": '$availability'}]]}'

    curl_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}

create_product() {
    read -p "Имя товара: " name
    read -p "Размер товара: " size
    read -p "Код товара: " code
    read -p "Количество товара: " quantity

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Products.Create", "params": [[{"name": "'"$name"'", "size": "'"$size"'", "code": "'"$code"'", "quantity": '"$quantity"'}]]}'

    curl_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}

get_leftovers() {
    read -p "ID склада: " warehouse_id

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Warehouses.GetLeftOvers", "params": [{"warehouse_id": '$warehouse_id'}]}'

    curl_response=$(curl -s -X GET -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}

reserve_products() {
    read -p "ID склада: " warehouse_id
    read -p "Количество продуктов для резервирования: " quantity
    read -p "Код продукта для резервирования: " code

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Products.Reserve", "params": [[{"warehouse_id": '$warehouse_id', "quantity": '$quantity', "code": "'"$code"'"}]]}'

    curl_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}

cancel_reservation() {
    read -p "ID склада: " warehouse_id
    read -p "Количество продуктов для отмены резервирования: " quantity
    read -p "Код продукта для отмены резервирования: " code

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Products.CancelReservation", "params": [[{"warehouse_id": '$warehouse_id', "quantity": '$quantity', "code": "'"$code"'"}]]}'

    curl_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}

add_products() {
    read -p "ID склада: " warehouse_id
    read -p "Количество продуктов: " quantity
    read -p "Код продукта: " code

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Products.Add", "params": [[{"warehouse_id": '$warehouse_id', "code": "'"$code"'", "quantity": '$quantity'}]]}'

    curl_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}

transfer_products() {
    read -p "ID склада отправителя: " warehouse_from_id
    read -p "ID склада получателя: " warehouse_to_id
    read -p "Количество продуктов для перемещения: " quantity
    read -p "Код продукта: " code

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Products.Transfer", "params": [[{"warehouse_from_id": '$warehouse_from_id', "warehouse_to_id": '$warehouse_to_id', "code": "'"$code"'", "quantity": '$quantity'}]]}'

    curl_response=$(curl -s -X POST -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}

delete_products() {
    read -p "Код продукта: " code

    json_data='{"jsonrpc":"2.0", "id": 1, "method": "Products.Delete", "params": [[{"code": "'"$code"'"}]]}'

    curl_response=$(curl -s -X DELETE -H "Content-Type: application/json" -d "$json_data" http://localhost:8080/)
    echo "Ответ сервера:"
    echo "$curl_response"
    read -n 1 -s -r -p "enter, чтобы продолжить..."
}


handle_choice() {
    read -p ": " choice
    case $choice in
        1) create_product;;
        2) create_warehouse;;
        3) get_leftovers;;
        4) reserve_products;;
        5) cancel_reservation;;
        6) add_products;;
        7) transfer_products;;
        8) delete_products;;
        9) exit;;
        *) echo "Некорректный выбор";;
    esac
}

while true; do
    show_menu
    handle_choice
done