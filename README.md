
# User product service
This is service for storing and working with customer products. 
Service shceme:
![image](https://github.com/AnisaForWork/user_orders/assets/138590484/63f4a983-a967-4a86-9c65-d273a949a7c1)

## Requirements

- Relational database for storing info about user and their products;
- Use provided auth server;
- Generate check in PDF format using provided pattern;
- No hardcoded configurations;
- Use docker compose;
- Only owner can edit, view and delete their products;
- RESTful.

## Requests
Documented via swagger 2.
- `POST /auth/reg` Register user - info about user added to DB;
- `POST /auth/auth` Authenticate user - user provides password plus login and receives jwt token;
- `POST /product/` Create product 
- `DELETE /product/:barcode` Delete product  
- `GET /product/all` View all user products;
- `GET /product/:barcode` View user product by id;
- `GET /product/:barcode/check` Generate product check.
- `GET /product/check/:checkName` Get product check file.

## Used technologies
- DB - MySQL;
- Router with Gin;
- swagger with swaggo
    
