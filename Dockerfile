FROM golang:1.23.4-bullseye AS builder

# RUN apt update && apt install libc6

WORKDIR /knit

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY . .

RUN scripts/build.sh 

FROM kcllang/kcl:v0.11.0

WORKDIR /knit

# COPY --from=builder /lib/x86_64-linux-gnu/libc.so /lib/x86_64-linux-gnu/libc.so
COPY --from=builder /knit/build/knit .

ENTRYPOINT ["/knit/knit"]
