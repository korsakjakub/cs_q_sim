FROM golang:1.19 as build

WORKDIR /go/src/cs_q_sim

COPY . .
RUN go mod download; make build

FROM gcr.io/distroless/static-debian11

ENV CONFIG_PATH=config/
ENV CONFIG_NAME="config"
ENV CONFIG_TYPE="yaml"

COPY --from=build /go/bin/app /
CMD ["/qsim.out"]