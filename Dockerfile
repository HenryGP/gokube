FROM golang:latest
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o gokube-app .
EXPOSE 8080 
CMD ["./gokube-app"]
