# SSID - Sparsely Sequential ID

This algorithm generates sparsely sequencial 64-bit ID's, that are based in three main components: timestamp, worker/process number and a sequence number.
These are called sparsely sorted (or in mathematical terms [`k-sorted`](http://ci.nii.ac.jp/naid/110002673489/)). For example, keeping the `k` around the 10ms mark, between each ID generation batches, say A1 and A2, we guarantee that `ID(A1) < ID(A2)`, and also that both batches will be within 10ms in the ID space as well.

This project gets some of the foundational concepts out of the [Snowflake Project](https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake.html) the ID generator.

This ID generator is suitable to be used in a distributed environemnt, generating a 64-bit ID based in three components:

| Field | Bit lenght | Max value |Description |
| --- | --- | --- | --- |
| Timestamp | 40 | ~34.8 years | Machine's timestamp in miliseconds since a defined start date |
| GeneratorID | 8 | 256 | ID that identifies the instance that is generating the SSID |
| SequenceID | 15 | 32767 | The number of elements that can be generated with a single timestamp |

The only configurable ID is the `GeneratorID` that can represent a given machine, process, etc.

The base API is defined under the `ssid` package.

To use the API, two modalities are provided: a command line interface and a simple REST server.


## Build the image:

The simplest way to compile the project is by using [Docker](https://www.docker.com).

```
docker build -t ssid .
```

## (Option #1) Run the command line interface:

  Generate 1000 SSIDs
  ```
  docker run -ti ssid idgen generate -c 1000
  ```

  Example:
  ```bash
  INFO[0000] Configuration: {GeneratorID:0 StartTime:2016-07-04 00:00:00 +0000 UTC}
  INFO[0000] SSID: 268066528266551297
  INFO[0000] SSID: 268066528266551298
  INFO[0000] SSID: 268066528266551299
  INFO[0000] SSID: 268066528266551300
  INFO[0000] SSID: 268066528266551301
  INFO[0000] Total Generation Time: 12.685µs
  ```

  To ge the full options available by the command line run
  ```
  docker run -ti ssid  idgen generate --help
  ```
##  (Option #2) Run the REST server

  ```
  docker run -p 8080:8080  -t ssid
  ```
  This starts a REST webserver making it available in the por 8080.

  To reach the server hit the following endpoints:

###  1. Generate a single SSID:
  Running on  curl or browser:
  ```
  http://localhost:8080/ssid
  ```
  Example response:
  ```json
  {
    "config": {
      "GeneratorID": 0,
      "StartTime": "2016-07-04T00:00:00Z"
    },
    "ssids": [
      267965642873765889
    ]
  }
  ```

###  2. Generate a set of SSIDs:
  ```
  http://localhost:8080/ssid/<#SSID>/<GeneratorID>
  ```
  * `<#SSID> - number of SSIDs to generate`

  * `<GeneratorID> - generator ID`

  Running on curl or browser:
  ```
  http://localhost:8080/ssid/1000/7
  ```
  Example response:
  ```json
  "config": {
    "GeneratorID": 7,
    "StartTime": "2016-07-04T00:00:00Z"
  },
  "ssids": [
    267968567352262657,
    267968567352262658,
    267968567352262659,
    267968567352262660,
    267968567352262661,
    267968567352262662,
    267968567352262663,
    267968567352262664,
    267968567352262665,
    267968567352262666,
    267968567352262667,
    267968567352262668,
    267968567352262669,
    267968567352262670,
    267968567352262671,
    267968567352262672,
    267968567352262673,
    267968567352262674,
    267968567352262675,
    (...)
  ```

