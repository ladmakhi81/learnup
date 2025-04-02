# Learning Management System (LMS)

A feature-rich **Learning Management System (LMS)** built with **Golang**, **GORM**, and **Tus protocol** for efficient learning management. This system supports **students, teachers, supporters, admins, and super admins** with a wide range of functionalities.

## 🚀 Features

### 🔹 Super Admin

- Manage entire system, including:
    - **Tickets, Forums, Courses, Teachers, Students, Admins**
    - **Permissions & Roles, Banners, Categories, Orders**
    - **Transactions, Notifications, Carts, Course Videos & Participants**

### 🔹 Admin

- Verify new admins and courses
- Set course fees
- Answer support tickets
- Manage forums

### 🔹 Teachers

- Create and manage courses
- Verify discount coupons (created by admin)
- Upload course videos
- Manage forums and respond to questions/comments
- Manage enrolled students

### 🔹 Students

- Comment and like courses
- Participate in webinars
- Purchase courses
- Create transactions via payment gateway

## 🛠 Tech Stack

- **Backend:** Golang
- **ORM:** GORM
- **File Uploads:** Tus Protocol
- **Database:** PostgreSQL
- **Authentication:** JWT-based authentication
- **Messaging & Notifications:** WebSockets & Email notifications

## 📌 Installation

### Prerequisites

- Golang installed (v1.20+ recommended)
- PostgreSQL/MySQL database setup
- Docker (optional for containerized deployment)

### Steps

```sh
# Clone the repository
git clone https://github.com/yourusername/lms.git
cd lms

# Install dependencies
go mod tidy

# Setup environment variables
cp .env.example .env
# Edit .env file with your configurations

# Run database migrations
go run main.go migrate

# Start the server
go run main.go
```