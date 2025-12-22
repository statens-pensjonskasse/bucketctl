FROM cr.spk.no/docker-hub/alpine/helm:3.18.3 AS helm

FROM ghcr.io/statens-pensjonskasse/rockylinux:9-minimal@sha256:236cb1b370a3e075a9200c446017023e03da9a4e90c04adcf2fda009da49bc9e

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
