FROM golang:latest

WORKDIR /go/src/github.com/uknth/faker

COPY . .

RUN go mod tidy
RUN go install -v github.com/uknth/faker/...

RUN export PATH="$PATH:/go/bin"

CMD ["faker"]