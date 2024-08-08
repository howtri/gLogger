# Idea

Server in Golang that writes to and reads from log files in a standardized format like
time severity service message

Use gRPC and then have python clients that communicate with the server over gRPC

Plan
- Write the protobuf for the log message [X]
- generate the server gRPC code [X]
- in main initialize the gRPC server so it is listening [X]
- create a function handleRequest that is passed to a goRoutine along with a mutex for writes []
- have a function that handles the write path and one for the read path [X]
- maintain another function that handles emitting metrics through the use of channels []
- create the grafana dashboards []
- Create Docker Images []
- Deploy on K8s []


# Notes
    route_guide.pb.go, which contains all the protocol buffer code to populate, serialize, and retrieve request and response message types.
    route_guide_grpc.pb.go, which contains the following:
        An interface type (or stub) for clients to call with the methods defined in the RouteGuide service.
        An interface type for servers to implement, also with the methods defined in the RouteGuide service.

# Commands

python3 -m grpc_tools.protoc -I=app/protocol --python_out=client/ --grpc_python_out=client/ app/protocol/logger.proto
export PATH=$PATH:$(go env GOPATH)/bin

# gLogger service start
sudo docker build -t glogger .
sudo docker run -p 8081:8081 -p 2112:2112 glogger

# gLogger, Prometheus collector, Grafana start
docker-compose up --build

# Stop services
docker-compose down

# Get K3s
curl -sfL https://get.k3s.io | sh -

# Push glogger to docker hub
docker push howtri/glogger:latest