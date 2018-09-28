FROM golang:1.10

WORKDIR /go/src/spfloop
COPY . .

RUN go install -v
EXPOSE 8053
CMD ["spfloop"]