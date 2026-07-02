# Trading Risk Management

A modern trading risk management web application built with Go. It helps traders calculate position size, risk per trade, leverage, risk/reward ratio, and manage trades through an intuitive dashboard.

## Features

- 🔐 JWT Authentication
- 👤 User Registration & Login
- 📊 Position Size Calculator
- 💰 Risk Management Calculator
- ⚖️ Leverage Support
- 📈 Risk/Reward Calculation
- 📅 Economic Calendar
- 📰 Market News
- 🎨 Responsive UI
- 🗄️ MySQL Database
- ⚡ Built with Go & GORM

---

## Tech Stack

- Go
- Gin
- GORM
- MySQL
- JWT Authentication
- HTML / CSS / JavaScript

---

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/Amir-shafiei/Trading-Risk_Management.git
cd Trading-Risk_Management
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Create a configuration file

Create a `.env` file in the project root.

Example:

```env
DBUSER=your_database_user
DBPASSWORD=your_database_password
DBHOST=localhost
DBPORT=3306
DBNAME=trading_risk_management

SERVERPORT=8080

JWT_SECRET=your_secret_key
```

> **Important:** The `.env` file is not included in this repository. You must create your own using your local database credentials.

### 4. Create the database

Create a MySQL database with the name specified in your `.env` file.

Example:

```sql
CREATE DATABASE trading_risk_management;
```

If AutoMigrate is enabled, the required tables will be created automatically when the application starts.

---

## Run the project

```bash
go run main.go
```

The application will start on:

```
http://localhost:8080
```

---

## Project Structure

```
Trading-Risk_Management/
│
├── config/
├── controllers/
├── middleware/
├── models/
├── routes/
├── services/
├── templates/
├── static/
├── utils/
├── main.go
└── .env
```

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| DBUSER | Database username |
| DBPASSWORD | Database password |
| DBHOST | Database host |
| DBPORT | Database port |
| DBNAME | Database name |
| SERVERPORT | Application port |
| JWT_SECRET | Secret key used for JWT authentication |


