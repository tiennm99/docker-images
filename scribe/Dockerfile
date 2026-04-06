# ─────────────────────────────────────────────────────────────
# Stage 1: builder
#   Compiles Apache Thrift 0.9.1 (last version compatible with
#   the Scribe Thrift IDL) and then Scribe itself from source.
# ─────────────────────────────────────────────────────────────
FROM ubuntu:16.04 AS builder

# Avoid interactive prompts during package installation
ARG DEBIAN_FRONTEND=noninteractive

# Versions pinned for reproducibility
ENV THRIFT_VERSION=0.9.1
ENV SCRIBE_GIT_REF=master

# ── Build-time dependencies ──────────────────────────────────
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        autoconf \
        automake \
        libtool \
        pkg-config \
        git \
        wget \
        ca-certificates \
        # Thrift code-generation tools
        flex \
        bison \
        # Boost (filesystem, system, thread required by Scribe)
        libboost-dev \
        libboost-filesystem-dev \
        libboost-system-dev \
        libboost-thread-dev \
        # libevent (async I/O used by Thrift non-blocking server)
        libevent-dev \
        # OpenSSL (Thrift TLS transport)
        libssl-dev \
        zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

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
# Scribe's FB_INITIALIZE macro conflicts with AM_INIT_AUTOMAKE, causing
# automake to crash with "global options already processed". We patch
# automake to downgrade that fatal error to a warning.
RUN sed -i '/prog_error.*global options already processed/{N;d}' \
        /usr/share/automake-*/Automake/Options.pm

RUN set -eux; \
    git clone --depth 1 \
        https://github.com/facebookarchive/scribe.git /tmp/scribe; \
    cd /tmp/scribe; \
    autoreconf --force --install; \
    sed -i 's|BOOSTLIBDIR=`echo $BOOST_LDFLAGS.*|BOOSTLIBDIR=/usr/lib/x86_64-linux-gnu|g' configure; \
    ./configure LIBS="-lboost_system -lboost_filesystem"; \
    make -j"$(nproc)" -C src; \
    make -C src install; \
    rm -rf /tmp/scribe


# ─────────────────────────────────────────────────────────────
# Stage 2: runtime
#   Minimal image containing only the runtime libraries and the
#   compiled scribed binary.
# ─────────────────────────────────────────────────────────────
FROM ubuntu:16.04

LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"
LABEL description="Facebook Legacy Scribe – real-time log aggregation daemon"
LABEL version="2.2"

ARG DEBIAN_FRONTEND=noninteractive

# ── Runtime-only shared libraries ────────────────────────────
RUN apt-get update && apt-get install -y --no-install-recommends \
        libboost-filesystem1.58.0 \
        libboost-system1.58.0 \
        libboost-thread1.58.0 \
        libevent-2.0-5 \
        libssl1.0.0 \
        libstdc++6 \
    && rm -rf /var/lib/apt/lists/*

# Copy only what is needed from the builder stage
COPY --from=builder /usr/local/bin/scribed        /usr/local/bin/scribed
COPY --from=builder /usr/local/lib/libthrift*.so* /usr/local/lib/
RUN ldconfig

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
