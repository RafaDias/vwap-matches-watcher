# Crypto Watcher

Find vwap given subscriptions

[![CircleCI](https://circleci.com/gh/RafaDias/vwap-matches-watcher/tree/main.svg?style=shield)](https://circleci.com/gh/RafaDias/vwap-matches-watcher/tree/main)

## How It Works

![Tutorial](.static/getting-started.gif)

## How to config?
1. Clone the repository
2. Copy the .env.example in root directory with name .env

```console
git clone https://github.com/RafaDias/vwap-matches-watcher.git crypto-watcher
cd crypto-watcher
cp configs/.env.example .env
```

## How to run?

In the root directory, choose one of the environments below and run the commands:

<details>
<summary>Locally</summary>

```sh
make run
```
</details>

<details>
<summary>Docker</summary>

```sh
make build
docker run crypto-watcher:1.0.0
```
</details>

<details>
<summary>KIND (Kubernetes in docker)</summary>

```sh
make build      # creates a docker image for crypto-watcher 
make kind-up    # creates a cluster to simulate k8s
make kind-load  # load the crypto-watcher in envinroment
make kind-apply # Create a deployment with that image
make kind-logs  # Get logs from pods
```
</details>

## The Architecture
### Overview
![Alt text](.static/diagram.svg)


This project was developed following the principles of `clean architecture` and golang standard layout:
- [cmd/crypto-watcher](cmd/crypto-watcher) - It is the folder responsible for storing the main file and therefore the build itself;
- [configs](configs) - Just a folder for configurations, link environment variables;
- [deployments](deployments) - Packaging and Continuous Integration;
- [build](build) - Packaging (Docker);
- [internal](internal) - Project source files;
- [static](.images) - Static files, like images.

The source code is organized as follows:
- [application/usecase](internal/application/usecase) - Use cases orchestrate the flow of data to and from the entities, and direct those entities to use their Critical Business Rules to achieve the goals of the use case.
- [application/providers](internal/application/providers) - Integration with a support service.
- [domain](internal/domain) - The heart of the application. Contains the business rules.
- [infra](internal/infra) - It is the "dirtiest" layer of the project. Contains implementations of domain layer abstractions

## Development process
I started the project focused on the domain. I used unit tests to validate the concepts, and I thought of two structures: Price and Transaction.  
The transaction acts as an anti-corruption layer, so the outer layers do not use Price directly.  
This gives me more flexibility to change my domain without breaking contracts.

So I thought of the exchange service abstraction. I like working with interfaces because they make testing easier,
as I can inject test services and validate integrations.

I used a mock server to test WebSocket communication that just sends a message in the coinbase message format.

As it is not recommended to use external libraries only for testing and building, I did not use the decimal library, which is widely used to deal with floating points.
Working with default float can generate inaccurate results. I decided to keep it in float for simplicity, but I was could use big.Rat as well.

I even used https://www.piesocket.com/websocket-tester
to test the resiliency of the service when there is an unexpected shutdown.

That's it. Thanks for reading!

![](https://media.giphy.com/media/el7VG1XOOvi24oRXFt/giphy.gif)

## Author
- [Rafael Dias](https://www.linkedin.com/in/rafaeldiasmello/)