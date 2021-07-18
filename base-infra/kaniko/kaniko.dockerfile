FROM gcr.io/kaniko-project/executor:debug AS kaniko

FROM centos:7

USER root

RUN yum update -y

RUN yum -y install wget

RUN wget https://dl.google.com/go/go1.16.4.linux-amd64.tar.gz && mkdir /golang && tar -xzf go1.16.4.linux-amd64.tar.gz --directory /golang && mv /golang/go /usr/local

COPY server.go .

COPY --from=kaniko /kaniko .

COPY --from=kaniko /kaniko/warmer /kaniko/warmer
COPY --from=kaniko /kaniko/docker-credential-gcr /kaniko/docker-credential-gcr
COPY --from=kaniko /kaniko/docker-credential-ecr-login /kaniko/docker-credential-ecr-login
COPY --from=kaniko /kaniko/docker-credential-acr /kaniko/docker-credential-acr

COPY --from=kaniko /kaniko/ssl/certs/ /kaniko/ssl/certs/
COPY --from=kaniko /etc/nsswitch.conf /etc/nsswitch.conf

ENV HOME /root
ENV USER root
# ENV PATH $PATH:/usr/local/bin:/kaniko:/busybox
ENV SSL_CERT_DIR=/kaniko/ssl/certs
ENV DOCKER_CONFIG /kaniko/.docker/
ENV DOCKER_CREDENTIAL_GCR_CONFIG /kaniko/.config/gcloud/docker_credential_gcr_config.json
WORKDIR /workspace
RUN ["/kaniko/docker-credential-gcr", "config", "--token-source=env"]
WORKDIR /
ENTRYPOINT [ "/usr/local/go/bin/go", "run", "server.go" ]