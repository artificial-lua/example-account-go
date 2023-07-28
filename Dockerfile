FROM golang as go-builder

WORKDIR /builder

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build

FROM alpine

WORKDIR /app

COPY --from=go-builder /builder/example-account-go /app/example-account-go

COPY ./html /app/html

RUN apk add --no-cache bash

RUN apk add libc6-compat

CMD ["sh", "-c", "/app/example-account-go"]