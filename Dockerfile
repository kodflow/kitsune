FROM kodmain/builder:all as builder

USER root

COPY / /home/nobody/
RUN go build /home/nobody/test.go

FROM alpine as runner
COPY --from=builder /home/nobody/test /home/nobody/test
CMD [ "/home/nobody/test" ]