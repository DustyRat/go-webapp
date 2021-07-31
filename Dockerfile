FROM golang:alpine

# Set necessary environment variables needed for our image
ENV GO111MODULE=on

# Get build arguments
ARG SHORT_COMMIT
ARG NOW
ARG HOSTNAME
ARG GIT_URL
ARG ENVIRONMENT
ARG BRANCH_NAME

# Install Tools and dependencies
RUN apk add --update --no-cache openssl-dev musl-dev zlib-dev curl tzdata

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
    -X main.buildDate=$NOW \
    -X main.buildHost=$HOSTNAME \
    -X main.gitURL=$GIT_URL \
    -X main.environment=$ENVIRONMENT \
    -X main.branch=$BRANCH_NAME \
    -X main.version=$SHORT_COMMIT" \
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
EXPOSE 3000

# Command to run when starting the container
CMD ["/main"]
