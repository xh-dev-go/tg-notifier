FROM golang:1.19-alpine3.17 as build
COPY . /app
WORKDIR /app
RUN go build -o /app/app .

FROM golang:1.19-alpine3.17
COPY --from=build /app/app /app/app
RUN chmod +x /app/app
WORKDIR /app
RUN ls -al
ENTRYPOINT ["/app/app"]
