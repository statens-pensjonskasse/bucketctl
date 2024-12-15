FROM cr.spk.no/base/k8s-tools:20241215014900@sha256:a5bb1a8aa1c1c5a4a3650c8e9d5db9ba5e2365a2f727237b1496b18f42d376e9

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
