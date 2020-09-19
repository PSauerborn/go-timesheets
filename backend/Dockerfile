
FROM golang:latest as build

ENV GO111MODULE=on

WORKDIR /app/server

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build

FROM alpine:latest as server

WORKDIR /app/server

COPY --from=build /app/server/go-timesheets ./

RUN chmod +x ./go-timesheets

EXPOSE 10091

CMD [ "./go-timesheets" ]