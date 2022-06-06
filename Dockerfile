FROM golang:1.17.8

WORKDIR /tmp

RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
RUN chmod 700 get_helm.sh
RUN ./get_helm.sh

RUN helm plugin install https://github.com/hypnoglow/helm-s3.git


WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o ./warehouse-controller ./cmd/main.go

RUN ["chmod", "+x", "/app/scripts/entrypoint.sh"]

ENTRYPOINT ["/app/scripts/entrypoint.sh"]