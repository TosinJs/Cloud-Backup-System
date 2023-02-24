# Cloud-Backup-System

## Description

This project was built in completion of a backend developer challenge.

## Requirements
Create an API that serves as a cloud backup system.

### General
- Users can create accounts with
  - Username
  - Password
  - Email
- Users can upload files up to 200mb
- Users can download uploaded files
- Users can create folders to hold files
- Create admin users for managing content 
- Admins can mark media files as unsafe
- Unsafe files get automatically deleted
- Users can stream audio and video

## Documentation
The documentation of this application can be found @:
<li>Postman @ https://github.com/TosinJs/Cloud-Backup-System/blob/master/Cloud%20Backup%20System.postman_collection.json </li>

## Run the Application Locally

```bash
# Clone the repository
$ git clone https://github.com/TosinJs/Cloud-Backup-System.git

# Install dependencies
$ cd 'Cloud Backup System'
$ go get ./..

# configuration 
# Create .env file in the root folder
$ touch .env

# populate the .env file with your files
$ DB_CONNECTION_STRING = "your postgresql connection string"
$ DB_CLIENT = "pg"
$ JWT_SECRET = "your JWT secret"
$ PORT = "3000"
$ DSN = "mysql connection string"
$ JWTSECRET = "jwt secret for auth token generation"
$ AWS_ID = "AWS ID"
$ AWS_SECRET = "AWS SECRET"
$ AWS_REGION = "AWS REGION"
$ AWS_TOKEN = "This can be left blank"

# Database Migrations (The migration files are in the internal/setup/database/ folder)
## Using the Makefile
$ make migrate-up 
$ make migrate-down

## Using the migrate cli tool directly
$ migrate -database ${your mysql DSN} -path internal/setup/database/migrations up 
$ migrate -database ${your mysql DSN} -path internal/setup/database/migrations down

# start
## Using the Makefile
$ make run

## Using go commands
$ go run main.go

```

## Application Flow
<p>The business logic is split into two main services: </p>
<li>Users</li>
<li>Files</li>

### Users
<p>The Users service contains all the logic for registeration and authentication of users </p>
<p>A JWT is retured to the user when they signup or login. The JWT is used to access the <strong>Files</strong> service</p>
<p><strong>Admin Credentials can be created to access the certain routes</strong>

![login flow](https://user-images.githubusercontent.com/68669102/211182773-d4f712ac-9c4f-4520-97c1-48a918b3a7eb.PNG)

### Files
<p>On signup each user has a folder assigned to their account </p>
<p>All the endpoints in the files service are protected endpoints </p>
<p>Users can Create, Download, and Delete files and folders from their folders </p>
<p>Admin users can flag files as offensive. After three flags the file is automatically deleted</p>
