import grpc
import logger_pb2
import logger_pb2_grpc

def run():
    # Connect to the gRPC server
    channel = grpc.insecure_channel('localhost:8081')
    stub = logger_pb2_grpc.LoggerStub(channel)

    # Make a Write request
    write_response = stub.Write(logger_pb2.LogMessage(
        time="2023-01-01T00:00:00Z",
        severity=1,
        service="cobb",
        message="Key vault is a good servic e I Swear."
    ))
    print("Write response:", write_response)

    # Make a Read request
    read_response = stub.Read(logger_pb2.ReadRequest(
        service="test-service"
    ))
    print("Read response:", read_response)

if __name__ == '__main__':
    run()
