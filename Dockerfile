FROM golang:alpine

# Get build arguments
ARG BUILD_DATE
ARG BUILD_HOST
ARG GIT_URL
ARG BRANCH
ARG SHA
ARG VERSION
ARG PORT

# Install Tools and dependencies
RUN apk update; \
    apk add --update --no-cache openssl-dev musl-dev zlib-dev curl tzdata

# Set necessary environment variables needed for our image
ENV GO111MODULE=on

# Move to working directory /build
WORKDIR /build

# Copy and download dependencies using go mod
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -ldflags "\
    -X main.buildDate=$BUILD_DATE \
    -X main.buildHost=$BUILD_HOST \
    -X main.gitURL=$GIT_URL \
    -X main.branch=$BRANCH \
    -X main.sha=$SHA \
    -X main.version=$VERSION" \
    -o main ./cmd/svr/main.go;

# Move to / directory as the place for resulting binary folder
WORKDIR /

# Copy binary from build to main folder
RUN cp /build/main .

# Copy static files
RUN cp -r /build/swagge[r] ./swagger

# Clean up build folder
RUN rm -rf /build

# Export necessary port
EXPOSE $PORT

# Command to run when starting the container
CMD ["/main"]
