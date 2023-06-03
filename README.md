# Architecture_4

## Description

The purpose of the lab work is to implement a load balancing algorithm and perform various tasks related to it. The tasks include:

1. Implementing a load balancer algorithm and covering it with unit tests
2. Implementing integration tests for the load balancer.
3. Configuring continuous integration with the integration tests.

## Prerequisites

To prepare for the tasks, you must configure your machine's environment by installing **Docker** and **Docker Compose**, as well as cloning the template repository.

## About the details of laboratory work

The lab work involves modifying the code in the cmd/lb package to implement a load balancer algorithm based on a given variant. The load balancer should select a healthy server based on the chosen algorithm and proactively monitor the health of all available servers.

You can verify the algorithm's functionality manually by running the load balancer and sending requests, or automatically by covering it with unit tests using mock objects.

The lab work also includes implementing integration tests to ensure that the assembled component (the load balancer) can be successfully launched and uses the correct algorithm. The integration tests involve sending a certain number of requests to the load balancer and checking if they are distributed to different servers.

## Instruction

To check the functionality of our code, follow these steps:

1. Open a terminal.

2. Build the Docker image by executing the following command:

```
docker build -t practice-4 .
```

This command builds a Docker image with the tag "practice-4" based on the Dockerfile in the current directory.

3. Run the Docker container using the following command:

```
docker run --rm practice-4
```
This command starts a Docker container based on the "practice-4" image and removes it after it exits (--rm flag).

4. Use Docker Compose to bring up the necessary services by running the following command:

```
docker-compose up
```
Docker Compose reads the configuration from the docker-compose.yml file and starts the defined services.

5. Open another terminal.

6. In the new terminal, run the following command to execute the client code:

```
go run ./cmd/client
```

This command compiles and runs the Go code located in the "./cmd/client" directory.

7. Alternatively, if you want to run the server code, use the following command:

```
go run ./cmd/server
```
This command compiles and runs the Go code located in the "./cmd/server" directory.

8. Lastly, if you wish to run the load balancer code, execute the following command:

```
go run ./cmd/lb
```
This command compiles and runs the Go code located in the "./cmd/lb" directory.

Make sure you have the necessary dependencies and configuration files in place before running these commands.



## Conclusion

In conclusion, the lab work focuses on implementing a load balancing algorithm, covering it with tests, and verifying its functionality through manual and automated testing.

During the laboratory work, we explored various aspects and principles related to research in our field. We gained practical experience working with specialized laboratory equipment and software.

In conclusion, the laboratory experiment provided us with valuable insights and hands-on experience, allowing us to deepen our understanding of the subject matter. It was an enriching and informative learning experience that will undoubtedly contribute to our future studies and professional development.