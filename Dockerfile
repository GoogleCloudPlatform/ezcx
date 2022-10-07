# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
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
