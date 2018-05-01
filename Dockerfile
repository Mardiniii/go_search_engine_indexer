FROM iron/go:dev
WORKDIR /app

ENV SRC_DIR=/go/src/github.com/Mardiniii/go_search_engine_indexer/
COPY . $SRC_DIR

RUN cd $SRC_DIR; go build -o go_search_engine_indexer; cp go_search_engine_indexer /app/
ENTRYPOINT ["./go_search_engine_indexer"]
