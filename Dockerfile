# Stage 1: Build the React Frontend
FROM node:18-alpine AS frontend-builder
WORKDIR /app/client
COPY client/package*.json ./
RUN npm install
COPY client/ ./
RUN npm run build

# Stage 2: Build the Go Backend (UPDATED TO 1.24)
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the rest of the backend code
COPY . .
# Copy the built frontend from Stage 1 to the backend folder
COPY --from=frontend-builder /app/client/dist ./client/dist
# Build the Go binary
RUN go build -o main .

# Stage 3: Final Production Image
FROM alpine:latest
WORKDIR /root/
# Copy the binary and the static files
COPY --from=backend-builder /app/main .
COPY --from=backend-builder /app/client/dist ./client/dist

# Expose the port
EXPOSE 5000

# Run the app
CMD ["./main"]