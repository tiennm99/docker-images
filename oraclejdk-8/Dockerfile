FROM ubuntu:18.04

# Set environment variables
ENV JAVA_HOME /opt/jdk1.8.0_201
ENV PATH $JAVA_HOME/bin:$PATH
ENV JAVA_VERSION 8u201
ENV JAVA_URL https://repo.huaweicloud.com/java/jdk/8u201-b09/jdk-8u201-linux-x64.tar.gz

# Install dependencies
RUN apt-get update && \
    apt-get install -y wget curl && \
    rm -rf /var/lib/apt/lists/*

# Download and install Oracle JDK 8 u201 from Huawei Cloud repository
RUN mkdir -p /tmp/jdk && \
    cd /tmp/jdk && \
    wget -O jdk-8u201-linux-x64.tar.gz $JAVA_URL && \
    mkdir -p /opt && \
    cd /opt && \
    tar -xzf /tmp/jdk/jdk-8u201-linux-x64.tar.gz && \
    rm -rf /tmp/jdk

# Set up alternatives for java commands
RUN update-alternatives --install /usr/bin/java java $JAVA_HOME/bin/java 100 && \
    update-alternatives --install /usr/bin/javac javac $JAVA_HOME/bin/javac 100

# Verify installation
RUN java -version && javac -version

# Clean up
RUN apt-get clean

# Set working directory
WORKDIR /app

# Default command
CMD ["bash"]