FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o hookah cmd/hookah/main.go

FROM alpine:3.21.3 AS prod
WORKDIR /app

ENV PORT=3000
ENV CONFIG_PATH=/etc/hookah/config.json
ENV TEMPLATES_PATH=/etc/hookah/templates

COPY --from=build /app/hookah /app/hookah

EXPOSE ${PORT}

CMD ["./hookah"]
