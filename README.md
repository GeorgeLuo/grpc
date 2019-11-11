### Setup development environment
Build dev environment
```
sudo docker build -t grpc --file Dockerfile.dev .
```
Generate key and cert for mutual TLS. These default locations will place the key/cert local to where the command is ran. Changing the location will require configuration in the server and client to reference correctly.
```
chmod +x ./bin/generate-keys.sh
```
```
./bin/generate-keys.sh
```
Initialize dev environment
```
sudo docker run --volume "$(pwd)":/go/src/github.com/GeorgeLuo/grpc --interactive --tty --publish 8443:8443 grpc
```
Run go server, exposed to port 8443
```
go run *.go
```

### Use grpc server with client
use client.go as cli-like process. To start

```
go run grpc-client/* start -cert cert.pem -key key.pem -command ./resources/test_process.sh -host localhost
```
to stop
```
go run grpc-client/* stop -cert cert.pem -key key.pem -task_id "987f769fca40-3635" -host localhost
```
to get status
```
go run grpc-client/* status -cert cert.pem -key key.pem -task_id "987f769fca40-3635" -host localhost
```

The status subcommand also supports multiple task_ids in param. This is done by providing more than one task_id parameter.

### Use grpc with curl

Following the above instructions to generate a key and cert with a self-signed cert

test a message to run a process using POST
```
curl   -X POST   --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem  https://localhost:8443/start   -H 'Content-Type: application/json'   -d '{"command":"./resources/test_process.sh"}'
```
test a message to stop a process using POST
```
curl   -X POST  --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem  https://localhost:8443/stop   -H 'Content-Type: application/json'   -d '{"task_id":"987f769fca40-3635"}'
```
test a message to get status by task_id
```
curl   -X POST   --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem https://localhost:8443/status   -H 'Content-Type: application/json'   -d '{"task_id":"987f769fca40-3635"}'
```

## Aliased Requests

An alias can be provided to be used instead of a task id for record. Pass an alias to start as a field of the body.

```
curl   -X POST   --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem  https://localhost:8443/start   -H 'Content-Type: application/json'   -d '{"command":"./resources/test_process.sh", "alias":"test"}'
```

Retrieve the status using the same aliases
```
curl   -X POST   --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem https://localhost:8443/status   -H 'Content-Type: application/json'   -d '{"alias":"test"}'
```

Or using the cli, to start using an alias
```
go run grpc-client/* start -cert cert.pem -key key.pem -command ./resources/test_process.sh -alias test_proc -host localhost
```
to get status with alias
```
go run grpc-client/* status -cert cert.pem -key key.pem -alias test_proc -host localhost
```

Note in the case a task id AND an alias is provided (to status or stop endpoint), the alias will take priority in evaluation. If the alias is not mapped, the task id will NOT resolve. This is due to future consideration where alias will encapsulate multiple processes and will provide the more complex output.

## Execute Jobs

Jobs are collections of processes to execute under a single job alias. A job alias may not take the same name as a previously taken alias naming a task (and vice-versa). Tasks are formed within the request in the same format as a standard task (can still be aliased). Naming a task within the job with the same alias will lead to unexpected behavior (to be determined what that is).

A task that fails under a job's execution will trigger a bailout of the previously fired tasks. Logically, identical to sending a stop process on each of the task_ids within a job. This action should eventually be configured as an option.

A sample start:
```
curl   -X POST   --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem https://localhost:8443/jobs/start   -H 'Content-Type: application/json'   -d '{"alias":"test","tasks":[{"command":"echo 1"},{"command":"echo 2"},{"command":"echo 3"},{"command":"echo 4"},{"command":"echo 5"}]}'
```

status:
```
curl   -X POST   --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem https://localhost:8443/jobs/status   -H 'Content-Type: application/json'   -d '{"alias":"test"}'
```

stop:
```
curl   -X POST   --cert ./cert.pem   --key ./key.pem   --cacert ./cert.pem https://localhost:8443/jobs/stop   -H 'Content-Type: application/json'   -d '{"alias":"test"}'
```

## Remote Usage With Docker

To generate a set of cert and key, run the ./bin/generate-keys-remote.sh executable, and make the additions to your openssl.cnf file:

```
[ req ]
...
req_extensions          = san_reqext

[ san_reqext ]
subjectAltName      = @alt_names

...

[ alt_names ]
IP.0            = XXX.XXX.XXX.XX
```

Once the files are generated, deliver them to the remote host. On the remote host

```
sudo docker build -t grpc .
sudo docker run --volume /keys:/keys --interactive --tty --publish 8443:8443 grpc
```

The first “/keys” directory is the directory from the host running the docker container, which should contain the cert and key file. In this usage, the docker run will point to the key and cert in the volume. Remember you must now provide the remote host information through the client:

```
grpc-client start -host XXX.XXX.XXX.XX -command some_command
```

## TODO

- Support alias mapping for multiple commands
- Scheduling job
- proxy jobs between hosts
- pass in json files to define jobs of alias mapped to commands
- server side configs for location of scripts/processes
- client side configs for location of job definitions
