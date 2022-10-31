FROM golang:1.19 as builder
WORKDIR /
COPY . ./gocloudcamp
WORKDIR /gocloudcamp
RUN CGO_ENABLED=0 GOOS=linux go build -a -o cloud-app
#CMD [ "./cloud-app" ]

FROM alpine:3.16
WORKDIR /gocloudcamp
COPY --from=builder /gocloudcamp/cloud-app .
COPY --from=builder /gocloudcamp/migrations ./migrations
CMD [ "./cloud-app" ]