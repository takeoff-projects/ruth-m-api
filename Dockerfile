FROM golang:1.17.0-alpine3.13
RUN mkdir /app
ADD . /app
WORKDIR /app
ENV GOOGLE_CLOUD_PROJECT=roi-takeoff-user47
RUN go build -o main .
CMD ["/app/main"]