FROM golang:1.9

ENV LAST_FULL_REBUILD 2017-10-05

RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega
RUN go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
RUN go get -u -d github.com/dimus/smithwatr
RUN go build -tags 'postgres' -o ${GOPATH}/bin/migrate github.com/mattes/migrate/cli

RUN apt-get update && apt-get -yq install postgresql-client

WORKDIR ${GOPATH}/src/github.com/dimus/smithwatr
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ENTRYPOINT scripts/development.sh
