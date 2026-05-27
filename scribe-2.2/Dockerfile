# ─────────────────────────────────────────────────────────────
# Stage 1: builder
#   CentOS 7 chosen because automake 1.13 predates the 1.14+
#   "global options already processed" bug, and OpenSSL 1.0.2
#   has BN_init (required by Thrift 0.9.1).
# ─────────────────────────────────────────────────────────────
FROM centos:7 AS builder

# Versions pinned for reproducibility
ENV THRIFT_VERSION=0.9.1

# ── Point yum to vault (CentOS 7 is EOL) ────────────────────
RUN sed -i 's|^mirrorlist=|#mirrorlist=|g; s|^#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' \
        /etc/yum.repos.d/CentOS-Base.repo

# ── Build-time dependencies ──────────────────────────────────
RUN yum install -y \
        gcc gcc-c++ make autoconf automake libtool \
        git wget \
        flex bison \
        boost-devel \
        libevent-devel \
        openssl-devel \
        zlib-devel \
    && yum clean all

# ── Build Apache Thrift ──────────────────────────────────────
# Scribe uses the old Facebook Thrift IDL and C++ runtime, which
# is API-compatible with Apache Thrift 0.9.x.
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

# ── Build fb303 (Scribe dependency) ─────────────────────────
# fb303 is Facebook's base service interface, bundled with Thrift.
RUN set -eux; \
    cd /tmp/thrift-${THRIFT_VERSION}/contrib/fb303; \
    ./bootstrap.sh; \
    ./configure --with-thriftpath=/usr/local; \
    make -j"$(nproc)"; \
    make install; \
    ldconfig; \
    rm -rf /tmp/thrift*

# ── Build Scribe ─────────────────────────────────────────────
# Facebook archived Scribe at https://github.com/facebookarchive/scribe
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


# ─────────────────────────────────────────────────────────────
# Stage 2: runtime
#   Minimal image containing only the runtime libraries and the
#   compiled scribed binary.
# ─────────────────────────────────────────────────────────────
FROM centos:7

LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"
LABEL description="Facebook Legacy Scribe – real-time log aggregation daemon"
LABEL version="2.2"

# ── Point yum to vault (CentOS 7 is EOL) ────────────────────
RUN sed -i 's|^mirrorlist=|#mirrorlist=|g; s|^#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' \
        /etc/yum.repos.d/CentOS-Base.repo

# ── Runtime-only shared libraries ────────────────────────────
RUN yum install -y \
        boost-filesystem \
        boost-system \
        boost-thread \
        libevent \
        openssl \
    && yum clean all

# Copy only what is needed from the builder stage
COPY --from=builder /usr/local/bin/scribed        /usr/local/bin/scribed
COPY --from=builder /usr/local/lib/libthrift*.so* /usr/local/lib/
COPY --from=builder /usr/local/lib/libfb303*.so*  /usr/local/lib/
RUN echo "/usr/local/lib" > /etc/ld.so.conf.d/scribe.conf && ldconfig

# ── Least-privilege user ─────────────────────────────────────
RUN groupadd -r scribe \
    && useradd -r -g scribe -s /sbin/nologin scribe \
    && mkdir -p /var/log/scribe /etc/scribe \
    && chown -R scribe:scribe /var/log/scribe /etc/scribe

# ── Default configuration and entrypoint ─────────────────────
COPY config/scribe.conf       /etc/scribe/scribe.conf
COPY scripts/docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# /var/log/scribe  – persistent log output (mount a volume here)
# /etc/scribe      – configuration (mount a custom scribe.conf here)
VOLUME ["/var/log/scribe", "/etc/scribe"]

# Scribe listens on this TCP port by default
EXPOSE 1463

USER scribe

ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
