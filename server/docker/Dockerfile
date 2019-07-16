FROM centos:centos7
MAINTAINER Jordi Blasco <jordi.blasco@hpcnow.com>
ENV OGRT_VERSION 0.5.1
RUN yum -y install \
    curl \
    tar \
    bzip2 \
    wget && \
    yum clean all && \
    cd /usr/local && \
    wget -q https://github.com/georg-rath/ogrt/releases/download/v${OGRT_VERSION}/ogrt-server-v${OGRT_VERSION}.tar.bz2 && \
    tar -jxvf ogrt-server-v${OGRT_VERSION}.tar.bz2 && \
    rm -f ogrt-server-v${OGRT_VERSION}.tar.bz2 && \
    mv ogrt-server-v${OGRT_VERSION} ogrt && \ 
    cd ogrt && \ 
    mv /usr/local/ogrt/ogrt.conf /usr/local/ogrt/ogrt.conf.orig
COPY ogrt.conf /usr/local/ogrt/ogrt.conf
COPY entrypoint.sh /usr/local/ogrt/entrypoint.sh
EXPOSE 7971
ENTRYPOINT ["/usr/local/ogrt/entrypoint.sh"]
