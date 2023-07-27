FROM golang as go-builder

WORKDIR /builder

COPY ./* ./

RUN go build

CMD ["./example-account-go"]