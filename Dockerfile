FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# Install Golang
RUN apt-get update && \
    apt-get install -y golang-go

# Install Java 
RUN apt-get install -y openjdk-11-jdk

# Install latest edge Nextflow
RUN apt-get install -y wget && \
    mkdir -p /nextflow && \
    wget -O /nextflow/nextflow https://github.com/nextflow-io/nextflow/releases/download/v23.07.0-edge/nextflow-23.07.0-edge-all && \
    chmod +x /nextflow/nextflow

# Set environment variables
ENV PATH="/usr/local/go/bin:${PATH}"  
ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk-amd64
ENV PATH="/usr/lib/jvm/java-11-openjdk-amd64/bin:${PATH}"
ENV NXF_VER=23.07.0-edge
ENV PATH="/nextflow:${PATH}"

# Create app directory
WORKDIR /app 

# Copy source code
COPY . /app

# Build application 
RUN go build -o main .

# Command to run app
CMD ["/app/main"]
