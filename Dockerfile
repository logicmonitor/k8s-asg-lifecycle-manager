FROM golang:1.8

ENV KUBERNETES_VERSION="v1.7.5"
RUN curl -L "https://storage.googleapis.com/kubernetes-release/release/${KUBERNETES_VERSION}/bin/linux/amd64/kubectl" -o /usr/local/bin/kubectl \
    && chmod +x /usr/local/bin/kubectl


# This dockerfile assumes you've already compiled the
# scheduler binary in pwd

COPY k8s-asg-lifecycle-manager /k8s-asg-lifecycle-manager

EXPOSE 8080
ENTRYPOINT ["/k8s-asg-lifecycle-manager"]
