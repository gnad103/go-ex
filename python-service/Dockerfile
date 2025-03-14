FROM python:3.9-slim

WORKDIR /app

# Copy requirements file
COPY python-service/requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Copy the source code
COPY proto/ ./proto/
COPY python-service/ ./python-service/

# Generate Python code from proto file
RUN python -m grpc_tools.protoc -I./proto --python_out=. --grpc_python_out=. ./proto/service.proto

# Expose port 50052 for gRPC
EXPOSE 50052

# Run the service
CMD ["python", "python-service/main.py"]