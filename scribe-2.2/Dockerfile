FROM centos:7 AS builder

ENV THRIFT_VERSION=0.9.1

# CentOS 7 is EOL, so yum must use the vault repositories.
RUN sed -i 's|^mirrorlist=|#mirrorlist=|g; s|^#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' \
        /etc/yum.repos.d/CentOS-Base.repo

RUN yum install -y \
        autoconf \
        automake \
        bison \
        boost-devel \
        flex \
        gcc \
        gcc-c++ \
        git \
        libevent-devel \
        libtool \
        make \
        openssl-devel \
        wget \
        zlib-devel \
    && yum clean all

# Scribe uses Facebook's old Thrift IDL and C++ runtime, which match Apache Thrift 0.9.x.
RUN set -eux; \
    wget -q "https://archive.apache.org/dist/thrift/${THRIFT_VERSION}/thrift-${THRIFT_VERSION}.tar.gz" \
         -O /tmp/thrift.tar.gz; \
    tar -xzf /tmp/thrift.tar.gz -C /tmp; \
    cd /tmp/thrift-${THRIFT_VERSION}; \
    ./configure \
        --without-java \
        --without-erlang \
        --without-perl \
        --without-php \
        --without-ruby \
        --without-haskell \
        --without-go \
        --without-python \
        --without-qt4 \
        --without-c_glib \
        --without-tests; \
    make -j"$(nproc)"; \
    make install; \
    ldconfig

RUN set -eux; \
    cd /tmp/thrift-${THRIFT_VERSION}/contrib/fb303; \
    ./bootstrap.sh; \
    ./configure --with-thriftpath=/usr/local; \
    make -j"$(nproc)"; \
    make install; \
    ldconfig; \
    rm -rf /tmp/thrift*

RUN set -eux; \
    git clone --depth 1 \
        https://github.com/facebookarchive/scribe.git /tmp/scribe; \
    cd /tmp/scribe; \
    autoreconf --force --install; \
    sed -i 's|BOOSTLIBDIR=`echo $BOOST_LDFLAGS.*|BOOSTLIBDIR=/usr/lib64|g' configure; \
    ./configure LIBS="-lboost_system -lboost_filesystem"; \
    make -j"$(nproc)" -C src; \
    make -C src install; \
    rm -rf /tmp/scribe

FROM centos:7

LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"

# CentOS 7 is EOL, so yum must use the vault repositories.
RUN sed -i 's|^mirrorlist=|#mirrorlist=|g; s|^#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' \
        /etc/yum.repos.d/CentOS-Base.repo

RUN yum install -y \
        boost-filesystem \
        boost-system \
        boost-thread \
        libevent \
        openssl \
    && yum clean all

COPY --from=builder /usr/local/bin/scribed        /usr/local/bin/scribed
COPY --from=builder /usr/local/lib/libthrift*.so* /usr/local/lib/
COPY --from=builder /usr/local/lib/libfb303*.so*  /usr/local/lib/
RUN echo "/usr/local/lib" > /etc/ld.so.conf.d/scribe.conf && ldconfig

RUN groupadd -r scribe \
    && useradd -r -g scribe -s /sbin/nologin scribe \
    && mkdir -p /var/log/scribe /etc/scribe \
    && chown -R scribe:scribe /var/log/scribe /etc/scribe

COPY config/scribe.conf       /etc/scribe/scribe.conf
COPY scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

VOLUME ["/var/log/scribe", "/etc/scribe"]

EXPOSE 1463

USER scribe

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
