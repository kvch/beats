FROM debian:jessie
MAINTAINER Simão Martins "simao.martins@tecnico.ulisboa.pt"

EXPOSE 749 88

ENV DEBIAN_FRONTEND noninteractive
# The -qq implies --yes
RUN apt-get -qq update
RUN apt-get -qq install locales krb5-kdc krb5-admin-server
RUN apt-get -qq clean

RUN locale-gen "en_US.UTF-8"
RUN echo "LC_ALL=\"en_US.UTF-8\"" >> /etc/default/locale

# Kerberos user information
ENV REALM ELASTIC
ENV LIBBEAT_USER_PRINCIPAL elastic
ENV ES_SERVICE_PRINCIPAL HTTP/elasticsearch_kerberos
ENV PASSWORD changeme
ENV PRINCIPALS "$LIBBEAT_USER_PRINCIPAL $ES_SERVICE_PRINCIPAL"
ENV KRB5_TRACE /dev/stdout

COPY init-script.sh /tmp/
CMD /tmp/init-script.sh
