FROM golang:1.16.5-buster AS builder

RUN apk add --no-cache git

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN make download

COPY . .
RUN make build


FROM alpine:3.14.0

RUN apk add ca-certificates

COPY --from=builder /usr/src/app/backend /usr/share/bin/backend

CMD [ "/usr/share/bin/backend" ]
