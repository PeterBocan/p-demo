FROM alpine:latest 

RUN adduser -D app
USER app 

COPY --chown=app --chmod=700 ./demo /demo
ENV DEMO_PORT=6734

EXPOSE 6734

ENTRYPOINT [ "/demo" ]