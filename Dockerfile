FROM golang:1.20.3-alpine

#set working dir to app
WORKDIR /app

#copy go.mod and go.sum files to the working dir
COPY go.mod go.sum ./

#download and install any required go deps
RUN go mod download

#copy the entire source code to the working dir
COPY . . 

# build the go application
RUN go build -o main .

# expose the port specified by the PORT env variable
EXPOSE 3000

#set the entry point of the container to the executable
CMD ["./main"]
