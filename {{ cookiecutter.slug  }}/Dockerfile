FROM golang:1.21-alpine as builder
WORKDIR /app/
COPY . .
RUN go build -o bin/{{ cookiecutter.slug }} main.go

FROM alpine:3.19
WORKDIR /app/
COPY --from=builder /app/bin/{{ cookiecutter.slug }} .
ENTRYPOINT ["./{{ cookiecutter.slug }}"]
