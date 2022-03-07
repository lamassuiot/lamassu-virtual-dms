FROM golang:1.16
WORKDIR /app
COPY . .
WORKDIR /app/cmd
ENV GOSUMDB=off
RUN go mod tidy
WORKDIR /app
RUN CGO_ENABLED=0 go build -o dms ./cmd/main.go

FROM alpine:3.14
COPY --from=0 /app/dms /
CMD ["/dms"]
