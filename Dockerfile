FROM gcr.io/distroless/static:nonroot

COPY ./out/app /bin/wp-easy-sync

ENTRYPOINT [ "wp-easy-sync" ]
