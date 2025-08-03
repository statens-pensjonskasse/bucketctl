FROM cr.spk.no/docker-hub/alpine/helm:3.18.3 AS helm

FROM cr.spk.no/base/rockylinux:9.20250803014448-minimal@sha256:61d353c4805307db6adff8e6aa851981cb041f9cac5ef1f864b71634cade7b89

# Installerer helm ved å kopiere fra alpine/helm image
COPY --from=helm /usr/bin/helm /usr/bin/helm

# Installerer nødvendige CLI verktøy (trengs av bitbucket-config-applier og Jenkins jobb for bitbucket-config)
RUN microdnf install -y \
    binutils \
    curl \
    git \
    make \
    && microdnf clean all

COPY ./bin/bucketctl /usr/local/bin/

RUN adduser app
USER app
WORKDIR /home/app

RUN mkdir -p /home/app/.config/bucketctl &&\
    touch /home/app/.config/bucketctl/config.yaml &&\
    chmod 600 /home/app/.config/bucketctl/config.yaml &&\
    echo "base-url: https://git.spk.no" >> /home/app/.config/bucketctl/config.yaml &&\
    echo "git-url: ssh://git@git.spk.no:7999" >> /home/app/.config/bucketctl/config.yaml &&\
    echo "limit: 9001" >> /home/app/.config/bucketctl/config.yaml

ENTRYPOINT ["/usr/local/bin/bucketctl"]
