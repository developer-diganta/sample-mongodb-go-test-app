# Sample-mongodb-go-test-app

This is a sample application built with Go that interacts with a MongoDB database. It stores a person's name with a unique id.

## Prerequisites

Before running the application, make sure you have the following prerequisites installed on your system:

- Go

## Getting Started

To get started with the application, follow these steps:

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/sample-mongodb-go-test-app
   ```
  
 2. Navigate to the project directory:
    ```bash
    cd sample-mongodb-go-test-app
    ```
 3. Install the dependencies: 
    ```bash
    go mod download
    ```
 4. Set up environment variables:

    4.1 Create a .env file in the project root directory.
 
    4.2 Add the following environment variables to the file:
    ```bash
    MONGODB_USERNAME=your-mongodb-username
    MONGODB_PASSWORD=your-mongodb-password
    ```
    
5. Replace the ```cluster0.wyqvltm.mongodb.net ``` of variable  ```connectionString```  with your cluster address.
6. Build and run the application:
    ```bash
    go run main.go
    ```
 The server will start running on http://localhost:8080.

