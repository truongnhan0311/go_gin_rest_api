FROM golang:alpine

RUN apk add --no-cache --upgrade bash nano

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /Docker
COPY ./Docker .
RUN cat crontab.staging | crontab -
WORKDIR /api

# Copy and download dependency using go mod
COPY api/go.mod .
COPY api/go.sum .

RUN go mod download

# Copy the code into the container
COPY ./api .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN mkdir download
RUN mkdir config

RUN cp /api/main .
COPY /Docker/.env.docker ./.env
COPY /Docker/run.sh ./run.sh
RUN chmod 777 ./run.sh
RUN cp /api/config/arial.ttf ./config/arial.ttf
RUN cp /api/config/phieu-bao-hanh.pdf ./config/phieu-bao-hanh.pdf
# Export necessary port
EXPOSE 3001

# Command to run when starting the container
CMD ["./run.sh"]
#CMD ["./run.sh"]