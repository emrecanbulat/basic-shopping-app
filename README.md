# **About The Project** 
**Language, libraries, tools, and IDEs:**
-   **Language**
    -   [Go v1.20](https://go.dev/dl/ "https://go.dev/dl/")        
-   **Libraries**
    -   [go-gorm/gorm](https://github.com/go-gorm/gorm "https://github.com/go-gorm/gorm") (ORM library for Golang)
    -   [jwt.io](https://jwt.io/ "https://jwt.io/") (JSON Web Tokens)
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

  **!! Troubleshooting on the .env path**

  *if you get an error when you run the project that the values in your .env file cannot be read;

 Open `cmd/api/main.go` file.

 Please set `godotenv.Load(".env" )` to `godotenv.Load("../../.env")`

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
## » You can find more details about the project in the shoppingAppDoc.pdf file.
