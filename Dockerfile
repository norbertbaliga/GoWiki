FROM golang:alpine as builder
WORKDIR /app 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" wiki.go

FROM scratch
WORKDIR /app
ARG pages_folder="./pages"
ENV GOWIKI_PAGES_PATH ${pages_folder}
COPY --from=builder /app/wiki .
COPY --from=builder /app/css ./css
COPY --from=builder /app/js ./js
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/pages ${GOWIKI_PAGES_PATH}
ENTRYPOINT ["/app/wiki"]