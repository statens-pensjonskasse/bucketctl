# Siden vi mangler brannmuråpninger for å laste ned helm og kustomize direkte tar vi de fra et image som allerede har dem
FROM quay.io/argoproj/argocd:latest AS tools
FROM old-dockerhub.spk.no:5000/base-os/rockylinux9-minimal

COPY --from=tools /usr/local/bin/helm /usr/local/bin/
COPY --from=tools /usr/local/bin/kustomize /usr/local/bin/
COPY ./bin/bucketctl /usr/local/bin/

RUN microdnf install -y \
    git \
    make &&\
    microdnf clean all

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