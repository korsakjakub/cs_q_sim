FROM golang:1.19 as build

WORKDIR /cs_q_sim
COPY . .
RUN CGO_ENABLED=0 go build -o cmd/web.out cmd/web/main.go

FROM gcr.io/distroless/static-debian11

ENV CONFIG_PATH=/go/

WORKDIR /go
COPY --from=build /cs_q_sim/cmd/web.out /cs_q_sim/config/config.yaml /cs_q_sim/web /go/
COPY --from=build /cs_q_sim/web /go/web
ENTRYPOINT ["/go/web.out"]