# syntax=docker/dockerfile:1

# Both base images are pinned by digest; bump them deliberately, not by tag.
FROM golang:1.26.9@sha256:37cfcca63222ec499bec5e8ad7aad52dbf7529db36d5812b8eae2f0621aa92d6 AS build
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p /out && sh scripts/build.sh /out/docs-drift

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
# CI passes --build-arg VERSION=$(cat VERSION) so the image metadata matches the binary.
ARG VERSION=dev
LABEL org.opencontainers.image.version=$VERSION
COPY --from=build /out/docs-drift /usr/local/bin/docs-drift
ENTRYPOINT ["/usr/local/bin/docs-drift"]
