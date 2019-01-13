<p align="center">
  <img src="https://i.imgur.com/0cgLwMe.png" height="150" />
</p>


## Go CrowdEngine

Official golang implementation of the CrowdCompute engine

## Building the source

Clone the repository to your desired destination:

```
$ git clone https://github.com/crowdcompute/crowdengine
```
Build the `gocc` binary:

```
$ cd crowdengine
$ make build
```

You can now run the `gocc` binary in build/bin/gocc


## Running gocc

Requirements: docker

### Gocc CLI Flags
`gocc` can be supplied with command-line flags. Here is a list of some important flags:

  * `--datadir` Data directory to store `gocc` data
  * `--addr`  P2P listening interface
  * `--port` P2P listening port
  * `--maxpeers` Maximum number of peers to connect
  * `--rpc` Enable the RPC interface
  * `--rpcservices` List of rpc services allowed
  * `--rpcwhitelist` Allow IP addresses to access the RPC servers
  * `--socket` Enable IPC-RPC interface
  * `--socketpath` Path of the socker/pipe file
  * `--http` Enable the HTTP-RPC server
  * `--httpport` HTTP-RPC server listening port
  * `--httpaddr` HTTP-RPC server listening interface
  * `--httporigin` HTTP-RPC cross-origin value
  * `--ws` Enable the WS-RPC server
  * `--wsport` WS-RPC server listening port
  * `--wsaddr` WS-RPC server listening interface
  * `--wsorigin` WS-RPC cross-origin value


### Gocc Configuration

You can pass a `toml` configuration file to the binary instead of specifying each flag with the following command:

```
$ gocc --config /path/to/your_config.toml
```
### Node Communication and Interaction

CrowdEngine provides a Plain-text JSON RPC request–response protocol which can be used by decentralized applications(DApps) or clients to interact with a node. Requests and responses are formatted in JSON and transferred over HTTP, Websockets and IPC.

#### JSONRPC2.0 HTTP

CrowdEngine exposes an HTTP server which accepts JSONRPC requests. The following flags can be used to enable the JSONRPC2.0 HTTP interface:
```
$ gocc --rpc --http --httpport <portnumber> --httpaddr 0.0.0.0 --httporigin "*"
```

#### JSONRPC2.0 Websocket

CrowdEngine exposes a Websocket server which accepts JSONRPC requests. The following flags can be used to enable the JSONRPC2.0 Websocket interface:
```
$ gocc --rpc --ws --wsport <portnumber> --wsaddr 0.0.0.0 --wsorigin "*"
```

#### JSONRPC2.0 IPC
System programs can interact with `gocc` by writing/reading to a given socket/pipe file, which is created by default within the data directory. To enable socket/pipe use the following command
```
$ gocc --rpc --socket --socketpath /where/to/save/gocc.ipc
```
