# builder image
FROM golang:1.18.3 as builder

RUN apt update && apt -y install libvips-dev

RUN mkdir /src

WORKDIR /src

ADD . .

RUN GOOS=linux go build -o ./build/main ./main.go


#deploying stage
FROM golang:1.18.3

RUN apt update && apt -y install libvips42

RUN mkdir -p /src

WORKDIR /src

COPY --from=builder /src/build/main .

#executable
CMD [ "./main" ]