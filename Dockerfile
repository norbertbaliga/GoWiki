FROM golang:1.16.11-bullseye
 
RUN mkdir -p /app
 
WORKDIR /app
 
COPY . /app
 
RUN go build ./wiki.go
 
CMD ["./wiki"]