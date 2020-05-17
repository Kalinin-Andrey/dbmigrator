# example
# docker run --rm -it --name dbmigrator -p=5432:5432 -v ${pwd}/migration:/migration -e dsn="host=localhost port=5432 dbname=postgres user=postgres password=postgres sslmode=disable" gomigrator gomigrator --help
FROM golang:1.14

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

# for mounting volume with migrations
ENV dir="/migration"

CMD ["dbmigrator --help"]

