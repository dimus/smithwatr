FROM golang:1.9-alpine

ENV LAST_FULL_REBUILD 2017-10-05
RUN apk update && apk add bash git postgresql-client && apk upgrade

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
RUN go get -u -d github.com/dimus/smithwatr
RUN go build -tags 'postgres' -o /go/bin/migrate github.com/mattes/migrate/cli


WORKDIR /go/src/github.com/dimus/smithwatr
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ENTRYPOINT scripts/development.sh
