FROM centos:7

LABEL maintainer="miti99 <tiennm99@outlook.com>"

ENV LANG=en_US.UTF-8
ENV JAVA_HOME=/usr/java/jdk-8
ENV PATH=$JAVA_HOME/bin:$PATH
ENV JAVA_URL=https://repo.huaweicloud.com/java/jdk/8u201-b09/

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN set -eux; \
    ARCH="$(uname -m)" && \
    if [ "$ARCH" = "x86_64" ]; \
    then \
    JAVA_PKG="$JAVA_URL"jdk-8u201-linux-x64.tar.gz; \
    else \
    JAVA_PKG="$JAVA_URL"jdk-8u201-linux-i586.tar.gz; \
    fi && \
    JAVA_SHA256="$(curl -L "$JAVA_PKG".sha256)" ;  \
    curl -L --output /tmp/jdk.tar.gz "$JAVA_PKG" && \
    echo "$JAVA_SHA256" */tmp/jdk.tar.gz | sha256sum -c -; \
    mkdir -p "$JAVA_HOME"; \
    tar --extract --file /tmp/jdk.tar.gz --directory "$JAVA_HOME" --strip-components 1 && \
    rm -f /tmp/jdk.tar.gz; \
    ln -sfT "$JAVA_HOME" /usr/java/default; \
    ln -sfT "$JAVA_HOME" /usr/java/latest; \
    for bin in "$JAVA_HOME/bin/"*; do \
    base="$(basename "$bin")"; \
    [ ! -e "/usr/bin/$base" ]; \
    alternatives --install "/usr/bin/$base" "$base" "$bin" 20000; \
    done; \
    java -Xshare:dump;
