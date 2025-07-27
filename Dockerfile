FROM cr.spk.no/docker-hub/alpine/helm:3.18.3 AS helm

FROM cr.spk.no/base/rockylinux:9.20250727014452-minimal@sha256:8a2692662f675d6ed7bf0c78c1e875a068ac6f021e4d3efb5dd938bfb1f50912

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
