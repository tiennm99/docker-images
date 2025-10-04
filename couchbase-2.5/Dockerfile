FROM ubuntu:12.04

LABEL maintainer="Tien Nguyen Minh <tiennm99@outlook.com>"

ARG CB_VERSION=2.5.2
ARG CB_RELEASE_URL=https://packages.couchbase.com/releases
ARG CB_PACKAGE=couchbase-server-enterprise_2.5.2_x86_64.deb

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH=$PATH:/opt/couchbase/bin:/opt/couchbase/bin/tools:/opt/couchbase/bin/install

RUN sed -i 's|http://archive.ubuntu.com/ubuntu/|http://old-releases.ubuntu.com/ubuntu/|g' /etc/apt/sources.list && \
    sed -i 's|http://security.ubuntu.com/ubuntu|http://old-releases.ubuntu.com/ubuntu|g' /etc/apt/sources.list && \
    apt-get update
RUN apt-get install -y wget
RUN wget --no-check-certificate $CB_RELEASE_URL/$CB_VERSION/$CB_PACKAGE
RUN apt-get install -y libssl0.9.8
RUN apt-get install -y librtmp0
RUN groupadd -g 1000 couchbase && useradd couchbase -u 1000 -g couchbase -M
RUN export INSTALL_DONT_START_SERVER=1 && \
    dpkg -i ./$CB_PACKAGE && rm -f ./$CB_PACKAGE

EXPOSE 8091 8092 11210 11211

CMD ["couchbase-server"]
