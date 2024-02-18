
# Command Queue Application

This is a Golang implementation of a Client-Server application that utilizes an external queue to manage and execute commands. The server reads commands from the queue and executes them, while the client sends commands to the queue. This application is designed to be robust, scalable, and easy to read and maintain.

## Overview

This application consists of a server component (`cmd/server/server.go`) and a client component (`cmd/client/client.go`). The server implements an ordered map data structure in memory and executes commands received from an external queue in parallel. The client accepts commands from the standard input or a file and sends them to the external queue.

## Specifications

The application implements the following features:

- **Server**:
  - Implements an ordered map data structure in memory.
  - Reads messages (commands) from an external queue.
  - Supports adding, removing, and retrieving items from the data structure.
  - Executes commands in parallel as much as possible.

- **Client**:
  - Can be configured from the command line or a file.
  - Sends messages (commands) to the external queue.
  - Supports running multiple clients in parallel.

- **External Queue**:
  - Can be Amazon Simple Queue Service (SQS) or RabbitMQ.

- **Client and Server Messages**:
  - Messages represent commands that the server should execute.
  - Commands include addItem, deleteItem, getItem, and getAllItems.

## Usage

### Server

To run the server, execute the following command:

```bash  
go run cmd/server/server.go -queue <queue_type> [additional_options]
```  

Supported `queue_type` values: `rabbitmq` or `aws`.

Additional options:
- `region`: AWS region (required for aws).
- `queueURL`: AWS SQS queue URL (required for aws).
- `connectionString`: RabbitMQ connection string (required for rabbitmq).
- `queueName`: Queue name (required for rabbitmq).

### Client
To run the client, execute the following command:
```bash  
go run cmd/client/client.go -queue <queue_type> [additional_options]
```  

Supported `queue_type` values: `rabbitmq` or `aws`.

Additional options:
- `region`: AWS region (required for aws).
- `queueURL`: AWS SQS queue URL (required for aws).
- `connectionString`: RabbitMQ connection string (required for rabbitmq).
- `file`: Input file path (optional).

### Dependencies
- AWS SDK for Go (for aws queue type)
- RabbitMQ Go Client (for rabbitmq queue type)

## Assumptions
-   The application assumes that the external queue is configured and accessible.
-   The ordered map data structure is implemented in-memory without using external packages.
-   Error handling for network failures or invalid configurations is not extensively covered in this version of the code.
-   Due to parallel processing of items from the external queue, it is possible that the order of items in the external queue does not match the order in which items are processed by the server. Therefore, the order of items in the external queue may not correspond to the order in which they are processed and stored in the in-memory queue.