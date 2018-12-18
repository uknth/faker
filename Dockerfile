FROM golang:latest

WORKDIR /go/src/github.com/uknth/faker

COPY . .

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN dep ensure -v

RUN go install -v github.com/uknth/faker/...

RUN export PATH="$PATH:/go/bin"

CMD ["faker-server"]