FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

RUN apt-get update && apt-get install -y nmap

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /discovery-agent

# Run
CMD ["/discovery-agent"]
