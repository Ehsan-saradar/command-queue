# Command Queue Application

This is a simple client-server application written in Golang that allows clients to send commands to a server via an external message queue. The server reads messages from the queue and processes commands accordingly.

## Server

### Running the Server

To run the server, follow these steps:

1. Navigate to the `cmd/server` directory:
    ```bash
    cd cmd/server
    ```

2. Build the server executable:
    ```bash
    go build
    ```

3. Run the server binary:
    ```bash
    ./server -queue <queue_type> [additional_flags]
    ```

   Replace `<queue_type>` with either `rabbitmq` or `aws` depending on the type of message queue you want to use. You may also need to provide additional flags based on the type of queue selected. See below for details on the available flags.

### Server Flags

- `-queue`: Type of queue (required)
    - Options: `rabbitmq`, `aws`

- Additional flags based on the queue type:
    - RabbitMQ:
        - `-conn`: RabbitMQ connection string (required)
        - `-queueName`: Queue name (required)

    - AWS:
        - `-region`: AWS region (required)
        - `-queueURL`: AWS SQS queue URL (required)

## Client

### Running the Client

To run the client, follow these steps:

1. Navigate to the `cmd/client` directory:
    ```bash
    cd cmd/client
    ```

2. Build the client executable:
    ```bash
    go build
    ```

3. Run the client binary:
    ```bash
    ./client -queue <queue_type> [additional_flags]
    ```

   Replace `<queue_type>` with either `rabbitmq` or `aws` depending on the type of message queue you want to use. You may also need to provide additional flags based on the type of queue selected. See below for details on the available flags.

### Client Flags

- `-queue`: Type of queue (required)
    - Options: `rabbitmq`, `aws`

- Additional flags based on the queue type:
    - RabbitMQ:
        - None

    - AWS:
        - `-region`: AWS region (required)
        - `-queueURL`: AWS SQS queue URL (required)

- Additional flags for both types:
    - `-file`: Input file path (optional)
        - If provided, the client will read commands from the specified file.
        - If not provided, the client will read commands from standard input.

## Example Usage

### Running Server with RabbitMQ

```bash
./server -queue rabbitmq -conn <rabbitmq_connection_string> -queueName <queue_name>

```

### Running Server with AWS SQS
```bash
./server -queue aws -region <aws_region> -queueURL <sqs_queue_url>
```
