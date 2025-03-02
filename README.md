# Packs API
This is a simple application that calculates the optimal number of packs required to fulfill an order based on predefined criteria. The application consists of a frontend and backend, with the backend built using Golang and the frontend using React with TypeScript.

## Architecture

The application consists of three main components: frontend, backend, and database.

- **Backend** - **Golang**
- **Frontend** - **React**
- **Database** - **MongoDB**

# API Endpoints

## 1. Create an Order

- **POST** `/api/orders`
  - **Description**: Create a new order with the specified item quantity and pack sizes.
  - **Request Body**:
    ```json
    {
      "items": 251,
      "packSizes": [500, 250, 1000, 2000]
    }
    ```
  - **Response**:
    ```201 Created```

- On the **frontend**, users input the item quantity, pack sizes and click **Add Order**.
- A **request** is sent to the server, which validates the order:
    - Ensures the order quantity is greater than zero.
    - If the order is invalid (empty or erroneous), the server responds with an **error**.
- For valid orders:
    - The server calculates the **optimal number of packs** required to fulfill the order.
    - The order, along with shipping details, is saved to the **database**.
    - A **confirmation** with order details is returned to the frontend.

## 2. Retrieving Orders

- **GET** `/api/orders`
  - **Description**: Retrieve a list of all orders.
  - **Response**:
    ```200 OK```
    ```json
    {
      "data": [
        {
          "id": "67c1cb3195069c4f6f4d6fbc",
          "items": 251,
          "packSizes": [500, 250, 1000, 2000],
          "packQuantity": {
            "500": 1,
            "250": 1
          },
          "createdAt": "2025-02-28T14:41:53.722Z",
          "updatedAt": "2025-02-28T14:41:53.722Z"
        }
      ]
    }
    ```

- The **home page**  features a table listing all orders.

---

# How to Run the Code

This project utilizes **Docker Compose** to streamline tasks such as building and running the application.

## Getting Started

### 1. Start the Application
Run the following command to launch all required containers:

```sh
docker-compose up -d
```

### 2. Verify Running Containers
Check the status of the containers to ensure they are running properly:

```sh
docker-compose ps
```

### 3. Set Up the Database
- Connect to MongoDB.
- Create a database named packs-api.
- Add a collection called orders.

### 4. Access the Application
Visit the provided link to continue using the application.
http://localhost:3000

### 5. Stop the Application
When finished, shut down the application by running:

```sh
docker-compose down
```

### 6. Run Tests
To run tests, use this command:

```sh
make test
```

## Sample page 
<img width="800" alt="sample-ui" src="https://github.com/user-attachments/assets/8e93c27c-fa90-47a9-ab29-3249b71a1f1b" />

## To Do
- Add integration tests for repository layer
- Add authentication/authorization
