# Use the official Golang image
FROM golang:1.22.5

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use the official PostgreSQL image as the base image
# FROM postgres:14.5

# # Set the environment variables for the PostgreSQL container
# ENV POSTGRES_USER=postgres
# ENV POSTGRES_PASSWORD=postgres
# ENV POSTGRES_DB=users

# # Create a directory for the PostgreSQL data
# RUN mkdir /var/lib/postgresql/data

# # Copy the SQL file into the container
# COPY sql/users.sql /docker-entrypoint-initdb.d/create_tables.sql

# Expose the PostgreSQL port
# EXPOSE 5432

# Start the PostgreSQL service
# CMD ["postgres"]

# Expose port 80 to the outside world
EXPOSE 80

# Command to run the executable
CMD ["/app/main"]