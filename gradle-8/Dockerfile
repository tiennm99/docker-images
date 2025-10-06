FROM ghcr.io/tiennm99/docker-images/oraclejdk-8:8u201

# Set environment variables
ENV GRADLE_HOME /opt/gradle
ENV PATH $GRADLE_HOME/bin:$PATH
ENV GRADLE_VERSION 8.14.3

# Install dependencies
RUN apt-get update && \
    apt-get install -y wget unzip && \
    rm -rf /var/lib/apt/lists/*

# Download and install Gradle
RUN wget -O gradle.zip "https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip" && \
    mkdir -p /opt && \
    unzip -d /opt gradle.zip && \
    mv /opt/gradle-${GRADLE_VERSION} /opt/gradle && \
    rm gradle.zip

# Create a non-root user
RUN groupadd -r gradle && useradd -r -g gradle gradle

# Set ownership
RUN chown -R gradle:gradle /opt/gradle

# Switch to the gradle user
USER gradle

# Set working directory
WORKDIR /home/gradle

# Default command
CMD ["gradle", "--version"]