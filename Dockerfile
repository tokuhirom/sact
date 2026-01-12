FROM gcr.io/distroless/static-debian12:nonroot

COPY sact /usr/local/bin/sact

ENTRYPOINT ["/usr/local/bin/sact"]
