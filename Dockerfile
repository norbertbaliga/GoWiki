FROM golang:1.16.5-buster
 
RUN mkdir -p /app
 
WORKDIR /app
 
COPY . /app
 
RUN go build ./wiki.go
 
CMD ["./wiki"]