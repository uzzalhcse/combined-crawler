# Use Ubuntu 22.04 as a base image
FROM ubuntu:22.04

# Install necessary tools
RUN apt-get update && \
    apt-get install -y wget

# Install Go
RUN wget https://golang.org/dl/go1.22.3.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
    rm go1.22.3.linux-amd64.tar.gz

# Set Go environment variables
ENV PATH="/usr/local/go/bin:${PATH}"

# Argument to specify which subdirectory to build
ARG PROJECT_DIR

# Set environment variable for project directory
ENV PROJECT_DIR=${PROJECT_DIR}
# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the application code
COPY . .

# Install Playwright CLI with the right version for later use
RUN PWGO_VER=$(grep -oE "playwright-go v\S+" /app/go.mod | sed 's/playwright-go //g') \
    && go install github.com/playwright-community/playwright-go/cmd/playwright@${PWGO_VER}

# Install dependencies and all browsers (or specify one)
RUN go run github.com/playwright-community/playwright-go/cmd/playwright@latest install --with-deps


# Build and run the Go application based on the argument
CMD ["sh", "-c", "cd /app/apps/${PROJECT_DIR} && go mod tidy && go run main.go"]
