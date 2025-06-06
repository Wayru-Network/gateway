# Stage 1: Build the application
FROM golang:1.24.3 AS build-stage
WORKDIR /app

# Install git and setup SSH
RUN apt-get update && apt-get install -y openssh-client git

# Setup SSH for private repos
RUN mkdir -p /root/.ssh
COPY .ssh/id_rsa /root/.ssh/
RUN chmod 600 /root/.ssh/id_rsa
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

# Configure git to use SSH for GitHub
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

# Set GOPRIVATE for private modules
ENV GOPRIVATE=github.com/Wayru-Network/serve

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
