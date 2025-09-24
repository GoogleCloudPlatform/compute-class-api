FROM us-central1-docker.pkg.dev/gke-autoscaling-kokoro/kokoro/golang:1.23

ENV GOPATH /tmpfs/gopath/
ENV PATH $GOPATH/bin:$PATH

RUN go version
CMD ["/bin/bash"]
