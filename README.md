# Hotel Reservation Backend Project

## Project Overview
This project provides a backend system for managing hotel reservations. Users can book rooms, and admins can manage bookings through a set of CRUD APIs. The project incorporates authentication and authorization using JWT tokens to ensure secure access to the system.

## Features
- **User Booking**: Users can book a room at a hotel.
- **Admin Booking Management**: Admins can check and manage bookings.
- **Authentication & Authorization**: Secure access using JWT tokens.
- **CRUD APIs**:
  - Hotels: Create, Read, Update, and Delete operations with JSON responses.
  - Rooms: Create, Read, Update, and Delete operations with JSON responses.
- **Database Management**:
  - Scripts for database migration and seeding.



## Project environment variables
---

HTTP_LISTEN_ADDRESS=:3000
JWT_SECRET=thejwtsecretofallsecrets
MONGO_DB_NAME=hotel-reservation
MONGO_DB_URL=mongodb://localhost:27017
MONGO_DB_URL_TEST=mongodb://localhost:27017
---
