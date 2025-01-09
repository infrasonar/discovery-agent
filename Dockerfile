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

# Environment variable
#   ENV NETWORK=<network>
#   ENV TOKEN=<token>
#   ENV DAEMON=0
#   ENV CONFIG_PATH=/etc/infrasonar
#   ENV ASSET_NAME=<asset name>
#   ENV ASSET_ID=<asset id>
#   ENV API_URI=https://api.infrasonar.com
#   ENV SKIP_VERIFY=<0 or 1>
#   ENV CHECK_NMAP_INTERVAL=14400

# Run
CMD ["/discovery-agent"]
