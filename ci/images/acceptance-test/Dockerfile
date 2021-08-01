FROM ubuntu

# Grafana
RUN  apt-get update \
  && apt-get install --yes --quiet --no-install-recommends ca-certificates curl gnupg gcc \
  && echo deb https://packages.grafana.com/oss/deb stable main > /etc/apt/sources.list.d/grafana.list \
  && curl https://packages.grafana.com/gpg.key | apt-key add - \
  && apt-get update \
  && apt-get install --yes --quiet --no-install-recommends grafana \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Go
ENV GOLANG_VERSION 1.16.6
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 be333ef18b3016e9d7cb7b1ff1fdb0cac800ca0be4cf2290fe613b3d069dfe0d
RUN    curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
    && echo "$GOLANG_DOWNLOAD_SHA256 golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz
ENV PATH /usr/local/go/bin:$PATH

# Ginkgo
ENV GOBIN /usr/local/bin
RUN go get github.com/onsi/gomega && go install github.com/onsi/ginkgo/ginkgo@latest
