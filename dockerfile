# syntax=docker/dockerfile:1

FROM golang:1.20
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
COPY .env .env
ENV TZ=Europe/Istanbul
RUN  CGO_ENABLED=0 GOOS=linux go build -o /bin/app .

FROM scratch
ENV TZ=Europe/Istanbul
COPY --from=0 /bin/app /bin/app
COPY .env .env

ENTRYPOINT [ "/bin/app" ]

