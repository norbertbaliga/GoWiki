FROM golang:rc
 
RUN mkdir -p /app
 
WORKDIR /app
 
COPY . /app
 
RUN go build ./wiki.go
 
CMD ["./wiki"]