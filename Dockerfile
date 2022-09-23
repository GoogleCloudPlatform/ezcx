# Borrowed from a Google example long, long ago.
# I'll need to figure out how to properly attribute this somehow!

FROM    golang:1.18-buster as builder
WORKDIR /app
COPY    . ./
RUN     go build -o service

FROM    debian:buster-slim
RUN     set -x && \
		apt-get update && \
		DEBIAN_FRONTEND=noninteractive apt-get install -y \
			ca-certificates && \
			rm -rf /var/lib/apt/lists/*
COPY    --from=builder /app/service /app/service

CMD     ["/app/service"]
