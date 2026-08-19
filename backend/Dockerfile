# --- Build stage ---------------------------------------------------------
FROM golang:1.26.6 AS build

WORKDIR /src

# Cache dependencies first (this layer only changes when go.mod/go.sum change).
COPY go.mod ./
RUN go mod download

COPY . .

# Build a static, stripped binary so it runs on a minimal base image.
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# --- Runtime stage -------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/server /server

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/server"]
