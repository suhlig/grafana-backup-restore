FROM ubuntu

ARG golang_version
ARG golang_download_sha256

# Grafana
RUN  apt-get update \
  && apt-get install --yes --quiet --no-install-recommends \
       ca-certificates \
       curl \
       g++ \
       gcc \
       gnupg \
       libc6-dev \
       make \
       pkg-config \
  && echo deb https://packages.grafana.com/oss/deb stable main > /etc/apt/sources.list.d/grafana.list \
  && curl https://packages.grafana.com/gpg.key | apt-key add - \
  && apt-get update \
  && apt-get install --yes --quiet --no-install-recommends grafana \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Go
RUN    curl -fsSL "https://golang.org/dl/go$golang_version.linux-amd64.tar.gz" -o golang.tar.gz \
    && echo "$golang_download_sha256 golang.tar.gz" | sha256sum -c - \
    && tar -C /usr/local -xzf golang.tar.gz \
    && rm golang.tar.gz
ENV PATH /usr/local/go/bin:$PATH

# Ginkgo
ENV GOBIN /usr/local/bin
RUN go get github.com/onsi/gomega && go install github.com/onsi/ginkgo/ginkgo@latest
