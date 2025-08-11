FROM golang:1.25rc3-bookworm

WORKDIR /app

COPY *.go /app

COPY . .

RUN go mod download 

RUN go build -o /github.com/zachdehooge/nexpull

CMD [ "/github.com/zachdehooge/nexpull" ]
