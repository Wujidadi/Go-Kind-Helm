FROM wujidadi/ubuntu-tuned:20250315

ENV GO_VERSION=1.24.1
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

RUN echo "\033[38;2;255;215;0m更新 apt 套件庫 ...\033[0m" && \
    apt-get update && apt-get install -y --no-install-recommends apt-utils && apt-get upgrade -y; \
    apt-get install -y --no-install-recommends software-properties-common; \
    echo "" && \
    echo "\033[38;2;255;215;0m安裝 Golang\033[0m" && \
    curl -LO https://go.dev/dl/go${GO_VERSION}.linux-arm64.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-arm64.tar.gz && \
    rm go${GO_VERSION}.linux-arm64.tar.gz && \
    go version && \
    echo "\033[38;2;255;215;0m安裝 PostgreSQL client\033[0m" && \
    apt-get install -y postgresql-client && \
    echo "" && \
    echo "\033[38;2;255;215;0m清理 apt 套件庫\033[0m" && \
    apt-get autoremove -y; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*; \
    echo "" && \
    echo "\033[38;2;255;215;0m更改預設文字編輯器為 Vim\033[0m" && \
    echo 'export EDITOR=vim' >> /root/.zshrc; \
    echo "" && \
    echo "\033[38;2;255;215;0m建立 Zsh 歷史紀錄目錄\033[0m" && \
    touch /root/.zsh_history

WORKDIR /app

COPY src/app/go.mod src/app/go.sum ./
RUN go mod download && \
    echo "\033[38;2;0;255;0m拉取 Redis SDK...\033[0m" && \
    go get github.com/redis/go-redis/v9 && \
    go mod tidy

COPY src/app/. ./
RUN go build -o /bin/go-server main.go
