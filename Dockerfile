FROM gcr.io/distroless/static:nonroot

COPY ./out/app /bin/wvc-sync

ENTRYPOINT [ "wvc-sync" ]
