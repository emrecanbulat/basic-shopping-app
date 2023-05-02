# **About The Project** 
**Language, libraries, tools, and IDEs:**
-   **Language**
    -   [Go v1.20](https://go.dev/dl/ "https://go.dev/dl/")        
-   **Libraries**
    -   [go-gorm/gorm](https://github.com/go-gorm/gorm "https://github.com/go-gorm/gorm") (ORM library for Golang)        
-   **Tools**    
    -   [PostgreSql](https://www.postgresql.org/download/ "https://www.postgresql.org/download/")
     -   [Docker](https://www.docker.com/products/docker-desktop/ "https://www.docker.com/products/docker-desktop/")           
    -   API platforms        
        -  [Postman](https://www.postman.com/downloads/ "https://www.postman.com/downloads/")             
-   **IDE’s**    
    -   [GoLand](https://www.jetbrains.com/go/download/#section=windows "https://www.jetbrains.com/go/download/#section=windows") (powerful code completion and nice [debugging feature](https://www.jetbrains.com/help/go/debugging-code.html "https://www.jetbrains.com/help/go/debugging-code.html")) or        
    -   [Visual Studio Code](https://code.visualstudio.com/ "https://code.visualstudio.com/")






 



## Installation

**1) Download the codebase**

```bash
    git clone https://github.com/emrecanbulat/ekinoks-shopping-app.git
```

**2) Fetch dependencies from** `go.mod`

```bash
  go mod download
```

**3) Set `.env` values**

*Run following command for generating a `.env` file from `.env.example`*

```bash
  cp .env.example .env
```

*You will see variables like app variables, PostgreSql credentials in `.env`. You must fill these variables before running the application*

**4) Run the Project**

```bash
  go run .\cmd\api .
```

*After running the following cURL command you should see `I’m OK.` message*
  
  `curl --request GET   --url localhost:8080/v1/healthcheck`

 ```bash
   {
    "message":"I'm OK.",
    "status":"available",
	"system_info":{"environment":"development","version":"1.0.0"}
    }
    
» This message means your Go server is up. You can also see some system information here
 ``` 
## Authorization
- **The first time you run the project, the seeder will automatically create an admin account.**
```bash
    email:      "admin@admin.com",
    password:   "password"
```
- **You can log in now. Get a response from the following cURL;**
```bash
    curl --request POST --location 'localhost:8080/v1/tokens/authentication' \
    --header 'Content-Type: application/json' \
    --data-raw '{
    "email": "admin@admin.com",
    "password": "password"
    }'
```
- **Get “token” from the response and set as the Bearer token in the HTTP HEADERS part. Example;**

```bash
{
  "Authorization" : "Bearer {{token}}"
}
```
## API Usage

- **PRODUCT**
#### 1) Create Product

```http
  curl --location 'localhost:8080/v1/products' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token}}' \
--data '{
      "title": "title",
      "description": "description",
      "price": 30899,
      "brand": "brand",
      "category": [
          "cat1",
          "cat2"
      ]
}'
```
  *Response;*
  ```
{
    "product": {
        "id": 22,
        "title": "title",
        "description": "description",
        "price": 30899,
        "brand": "brand",
        "category": [
            "cat1",
            "cat2"
        ]
    }
}
  ```


#### 2) Get Product

```http
curl --location 'localhost:8080/v1/products/${id}' \
--header 'Authorization: Bearer {{token}}'
```
  *Response;*
  ```
{
    "product": {
        "id": ${id},
        "title": "title",
        "description": "description",
        "price": 30899,
        "brand": "brand",
        "category": [
            "cat1",
            "cat2"
        ]
    }
}
  ```

#### 3) Get Product List (filter can be used)

```http
curl --location 'localhost:8080/v1/products?title=title&category=cat1&page=1&page_size=1&sort=-price&brand=brand' \
--header 'Authorization: Bearer {{token}}'
```
  *Response;*
  ```
{
    "meta": {
        "current_page": 1,
        "page_size": 1,
        "first_page": 1,
        "last_page": 22,
        "total_records": 22
    },
    "products": [
        {
            "id": 22,
            "title": "title",
            "description": "description",
            "price": 30899,
            "brand": "brand",
            "category": [
                "cat1",
                "cat2"
            ]
        }
    ]
}
  ```

#### 4) Update Product

```http
curl --location --request PUT 'localhost:8080/v1/products/${id}' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token}}' \
--data '{
      "title": "title (edit)",
      "description": "description (edit)",
      "price": 18899,
      "brand": "brand (edit)",
      "category": [
          "cat1 (edit)",
          "cat 2 (edit)"
      ]
}'

  ```
  *Response;*
  ```
{
    "product": {
        "id": ${id},
        "title": "title (edit)",
        "description": "description (edit)",
        "price": 18899,
        "brand": "brand (edit)",
        "category": [
            "cat1 (edit)",
            "cat 2 (edit)"
        ]
    }
}
  ```


#### 4) Delete Product

```http
curl --location --request DELETE 'localhost:8080/v1/products/${id}' \
--header 'Authorization: Bearer {{token}}'

  ```
  *Response;*
  ```
{
    "message": "product successfully deleted"
}
  ```

- **USER** 
#### 1) User Create & Register

```http
curl --location 'localhost:8080/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "full_name":"User",
    "email":"user@user.com",
    "password":"password",
    "address": "address",
    "phone":"0000000000"
}'

  ```
  *Response;*
  ```
{
    "access_token": {
        "expiry": "2023-05-03 15:20:23",
        "token": "{{token}}"
    },
    "user": {
        "id": ${id},
        "full_name": "User",
        "email": "user@user.com",
        "phone": "0000000000",
        "address": "address"
    }
}
  ```

- **ORDER** 
#### 1) New Order

```http
curl --location 'localhost:8080/v1/orders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer {{token}}' \
--data '{
      "product_id": ${id},
      "payment_type": "Credit_card",
      "amount_paid": 30899
}'

  ```
  *Response;*
  ```
{
    "order_details": {
        "amount_paid": 30899,
        "id": 10,
        "order_date": "2023-05-02 15:28:44",
        "payment_type": "Credit_card",
        "status": "processing"
    },
    "product": {
        "brand": "brand",
        "category": [
            "cat1",
            "cat2"
        ],
        "description": "description",
        "price": 30899,
        "title": "title"
    }
}
  ```

#### 2) Show Order

```http
curl --location 'localhost:8080/v1/orders/${id}' \
--header 'Authorization: Bearer {{token}}' \
--data ''

  ```
  *Response;*
  ```
{
    "order_details": {
        "amount_paid": 30899,
        "order_date": "2023-05-02 15:28:44",
        "payment_type": "Credit_card",
        "status": "processing"
    },
    "product": {
        "brand": "brand",
        "category": [
            "cat1",
            "cat2"
        ],
        "description": "description",
        "price": 30899,
        "title": "title"
    },
    "user": {
        "address": "address",
        "email": "admin@admin.com",
        "full_name": "Admin",
        "id": 1,
        "phone": "0000000000"
    }
}
  ```

#### 3) Order List

```http
curl --location 'localhost:8080/v1/orders?page_size=1' \
--header 'Authorization: Bearer {{token}}' \
--data ''

  ```
  *Response;*
  ```
{
    "meta": {
        "current_page": 1,
        "page_size": 1,
        "first_page": 1,
        "last_page": 10,
        "total_records": 10
    },
    "orders": [
        {
            "id": 1,
            "user": {
                "id": 1,
                "full_name": "Admin",
                "email": "admin@admin.com",
                "phone": "0000000000",
                "address": "address"
            },
            "product": {
                "id": 8,
                "title": "title (edit)",
                "description": "description (edit)",
                "price": 18899,
                "brand": "brand (edit)",
                "category": [
                    "cat1 (edit)",
                    "cat 2 (edit)"
                ]
            },
            "status": 0,
            "payment_type": "Cash",
            "amount_paid": 30899
        }
    ]
}
  ```