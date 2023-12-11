# Warehouse Service

**Warehouse Service** - бэкенд сервис, который позволяет резервироваться продукты для доставки. 
Программа предоставляет API для работы с приложением.

## :checkered_flag: Установка и запуск

### :package: Используя Docker

Для запуска **Warehouse Service** требуются следующие шаги:

1. Установить [Docker и Docker-compose](https://www.docker.com/) на локальную машину.
2. Склонировать репозиторий проекта.
3. Перейти в директорию `/backend`.
4. Собрать проект при помощи **Makefile:**
```bash
make -f Makefile
```
5. Использовать API, описанный ниже

## :bulb: API приложения

- ```GET /product``` - получение продуктов, возможных для резервирования, а также получение их общего количества и получение количества
каждого отдельного продукта. Пример получения списка продуктов:
```json
{
    "all_count": 33,
    "goods": [
        {
            "name": "milk",
            "size": "0.2m",
            "product_id": 1,
            "count": 15
        },
        {
            "name": "doors",
            "size": "1.8m",
            "product_id": 3,
            "count": 18
        }
    ]
}
```

- ```POST /product``` - добавление продуктов в таблицу продуктов. Пример добавления продукта:
```json
        {
            "name": "milk",
            "size": "0.2m",
            "count": 15
        }
```
- ```DELETE /product``` - удаление продукта из списка продуктов и из резерва. Пример работы удаления по следующему ID `[2]`: 
```json
[
    {
        "name": "tables",
        "size": "1.5m",
        "product_id": 2,
        "count": 20
    }
]
```

- ```POST /product/warehouse``` - резервирование продуктов для дальнейшей доставки. На вход подается массив уникальных ID продуктов.
Пример резервирования продуктов с ID `[1,2]`:
```json
[
    {
        "name": "milk",
        "size": "0.2m",
        "product_id": 1,
        "count": 15
    },
    {
        "name": "tables",
        "size": "1.5m",
        "product_id": 2,
        "count": 20
    }
]
```
- ```DELETE /product/warehouse``` - удаление продуктов из резерва. На вход подается массив уникальных ID продуктов. 
Пример удаления продуктов из резерва с ID `[1]`:
```json
[
    {
        "name": "milk",
        "size": "0.2m",
        "product_id": 1,
        "count": 15
    }
]
```

## Стартовые данные базы данных:

`Таблица products`

|  unique_code  |  product_name  |  size   |  count  |
|:-------------:|:--------------:|:-------:|:-------:|
|       1       |      milk      |  0.2m   |   15    |
|       2       |     tables     |  1.5m   |   20    |
|       3       |     doors      |  1.8m   |   18    |
|       4       |  microphones   |  0.15m  |    4    |

`Таблица warehouse`

| product_name | can_be_use |
| :----------: | :--------- |
 |             |            |


## Тесты
`Для использования тестов через postman`
1. Необходимо скачать коллекцию [postman-тестов](/postman_collection)
2. Запустить тесты из коллекции

`Для использования unit-тестов`
```bash
go test ./
```
Требования к тестам см. в [postman-test.md](/backend/doc/postman-test.md) и [unit-test.md](/backend/doc/unit-test.md)



