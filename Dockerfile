FROM golang:1.23.0

# Install protoc and necessary tools
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        protobuf-compiler \
        git \
    && rm -rf /var/lib/apt/lists/*

# Install protoc-gen-go plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

ENV PATH="/go/bin:${PATH}"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate gRPC code
RUN protoc --go_out=. --go-grpc_out=. proto/pvz.proto

# Build application
RUN go build -o server ./cmd/server

CMD ["./server"]