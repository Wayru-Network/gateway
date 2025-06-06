# Stage 1: Build the application
FROM golang:1.24.3 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o gateway ./cmd

# Stage 2: Copy the application to a distroless image
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build-stage /app/gateway .
EXPOSE 4050
CMD ["./gateway"]
