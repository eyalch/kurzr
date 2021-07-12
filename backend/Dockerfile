FROM golang:1.16.5-buster AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum Makefile ./
RUN make download

COPY . .

# CGO has to be disabled for alpine
RUN CGO_ENABLED=0 make build


FROM alpine:3.14.0

RUN apk add ca-certificates

COPY --from=builder /usr/src/app/backend /usr/share/bin/backend

CMD [ "/usr/share/bin/backend" ]
