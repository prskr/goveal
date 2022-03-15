FROM gcr.io/distroless/static:nonroot

USER nonroot:nonroot

COPY --chown=nonroot:nonroot goveal /app/goveal

EXPOSE 2233

ENTRYPOINT ["/app/goveal"]