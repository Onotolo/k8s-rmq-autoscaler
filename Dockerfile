FROM golang:buster as builder

ENV PROJECT github.com/medal-labs/k8s-rmq-autoscaler

WORKDIR /go/src/$PROJECT

COPY src/. /go/src/$PROJECT

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -o /k8s-rmq-autoscaler .

FROM alpine as release
COPY --from=builder /k8s-rmq-autoscaler /k8s-rmq-autoscaler

ENTRYPOINT ["/k8s-rmq-autoscaler"]
