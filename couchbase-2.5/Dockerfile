FROM ubuntu:12.04

LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"

RUN sed -i 's|http://archive.ubuntu.com/ubuntu/|http://old-releases.ubuntu.com/ubuntu/|g' /etc/apt/sources.list && \
    sed -i 's|http://security.ubuntu.com/ubuntu|http://old-releases.ubuntu.com/ubuntu|g' /etc/apt/sources.list

ENV CB_VERSION=2.5.2
ENV CB_RELEASE_URL=https://packages.couchbase.com/releases
ENV CB_PACKAGE=couchbase-server-enterprise_2.5.2_x86_64.deb
ENV CB_SHA256=27a79a65758023c34ed900e8ef8c54bab4a65f4c84b7c94359cba910800a4b19

ENV PATH=$PATH:/opt/couchbase/bin:/opt/couchbase/bin/tools:/opt/couchbase/bin/install

RUN apt-get update && \
    apt-get install -yq wget libssl0.9.8 librtmp0 dos2unix && \
    apt-get autoremove && apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN export INSTALL_DONT_START_SERVER=1 && \
    wget -N --no-check-certificate $CB_RELEASE_URL/$CB_VERSION/$CB_PACKAGE && \
    echo "$CB_SHA256  $CB_PACKAGE" | sha256sum -c - && \
    dpkg -i ./$CB_PACKAGE && rm -f ./$CB_PACKAGE

ADD scripts/couchbase-start /usr/local/bin/
RUN dos2unix /usr/local/bin/couchbase-start && \
    chmod +x /usr/local/bin/couchbase-start

ENTRYPOINT ["/usr/local/bin/couchbase-start"]

EXPOSE 8091 8092 11210 11211
VOLUME /opt/couchbase/var
