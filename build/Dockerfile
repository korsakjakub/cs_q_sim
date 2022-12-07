FROM golang:1.19 as build

WORKDIR /cs_q_sim
COPY . .
RUN CGO_ENABLED=0 make build

FROM gcr.io/distroless/static-debian11

ENV CONFIG_PATH=/go/

COPY --from=build /cs_q_sim/cmd/qsim.out /cs_q_sim/config/config.yaml /go/
CMD ["/go/qsim.out"]