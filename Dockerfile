FROM alpine

COPY ./out/app /bin/wp-easy-sync

ENTRYPOINT [ "wp-easy-sync" ]
