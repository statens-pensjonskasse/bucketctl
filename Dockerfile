FROM cr.spk.no/base/k8s-tools:20250209014905@sha256:ac12417845b9100a4f69541b737903c981c8fb4469fb10b70d10b6f37b06f84b

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
