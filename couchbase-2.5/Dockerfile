# Docs: https://docs.couchbase.com/couchbase-manual-2.5/cb-install

FROM ubuntu:12.04

LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"

COPY *.deb /tmp/

RUN dpkg -i /tmp/libgpg-error0_*.deb && \
    dpkg -i /tmp/libgcrypt11_*.deb && \
    dpkg -i /tmp/libp11-kit0_*.deb && \
    dpkg -i /tmp/libtasn1-3_*.deb && \
    dpkg -i /tmp/libgnutls26_*.deb && \
    dpkg -i /tmp/librtmp0_*.deb && \
    dpkg -i /tmp/couchbase-server-enterprise_*.deb && \
    rm /tmp/*.deb

# Expose ports
# 8091: Web Administration Port
# 8092: Couchbase API Port
# 11210: Internal/External Bucket Port
# 11211: Client interface (proxy)
EXPOSE 8091 8092 11210 11211

VOLUME /opt/couchbase/var

CMD ["/opt/couchbase/bin/couchbase-server"]
