# 
# elder-peeler
#
# docker build . -t elder-peeler:1
# kubectl run ek --rm -it --image=elder-peeler:1 --restart=Never -- sh
#
#FROM golang:1.21.4-alpine3.17
FROM --platform=linux/amd64 golang:1.21.4-alpine3.17 as buildx
#
ENV SQS_URL="https://sqs.us-west-2.amazonaws.com/123456789012/guytest"
#
WORKDIR /app
#
COPY go.mod .
COPY go.sum .
RUN go mod download
#
COPY *.go ./
#
RUN go build -o /app/elder-peeler
#
ENTRYPOINT [ "/app/elder-peeler" ]
CMD ["bogus"]
#
