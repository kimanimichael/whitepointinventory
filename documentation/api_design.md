# API Design

The system consists of one service with four main components at a high level:
- Users
- Farmers
- Purchases
- Payments

Currently, these four components interact directly with each other e.g. Purchases and Payments directly increment and decrement balances that are attributes to a farmer entry

Each of these components consists of various APIs

The components of this system and APIs relating to each of them are:

# Users

Exposed APIs are for creation, login, logout and retrieving of users from cookies

API endpoints:

- `POST /v1/users`: Create a new user account
- `GET /v1/users`: Retrieve a user
- `POST /v1/login`: Handle user login and issue cookie
- `POST /v1/logout`: Handle user logout and expire issued cookie

# Farmers

Exposed APIs are for creating, retrieving and deleting of farmer entries

API endpoints:

- `POST /v1/farmers`: Create a new farmer entry
- `GET /v1/farmers`: Retrieve single farmer by their name
- `GET /v1/farmer`: Retrieve all farmers
- `DELETE /v1/farmers/{farmer_id}`: Delete a farmer entry by ID

# Purchases

Exposed APIs are for creating, retrieving and deleting of purchase entries

API endpoints:

- `POST /v1/purchases`: Create a new purchase entry
- `GET /v1/purhcases`: Retrieve all purchases
- `GET /v1/purhcase`: Retrieve a single purchase by its ID
- `DELETE /v1/purchase/{purchase_id}`: Delete a purchase entry by its ID

# Payments

Exposed APIs are for creating, retrieving and deleting of payment entries

API endpoints:

- `POST /v1/payment`: Create a new payment entry
- `GET /v1/payments`: Retrieve all payment
- `GET /v1/payment`: Retrieve a single payment by its ID
- `DELETE /v1/payments/{payment_id}`: Delete a payment entry by its ID