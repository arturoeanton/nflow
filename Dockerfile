FROM golang:1.19 as development
# Add a work directory
WORKDIR /nflow
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
COPY config.docker.toml config.toml

RUN go mod tidy
RUN make build

# Expose port
EXPOSE 9090
# Start app
CMD /nflow/nFlow  -w /nflow/app/

# docker build -t hopbox/nflow .

