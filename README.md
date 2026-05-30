![Athena logo](https://i.imgur.com/UQggP60.png)
Athena is a lightweight Change Data Capture (CDC) solution that streams changes from Microsoft SQL Server to Apache Kafka. Built in Golang, it supports SASL-authenticated Kafka brokers and provides a straightforward setup experience. Unlike alternatives such as Debezium, which can be complex to configure and manage. Athena offers greater simplicity and operational ease. It automatically manages CDC setups, publishes database changes to a single Kafka topic, and delivers a clean, intuitive event format that is easy for downstream consumers to understand and process.

## :cyclone: How things work

- Creates a message for changes like `create`, `update`, `delete` for rows in MSSQL database tables to a single Kafka topic.
- Athena only creates messages for all new table changes. Existing ones are ignored.
- Kafka topic have to be created before hand. Unlike Debezium, Athena will not create the topic own it's own.
- By default, Athena will poll for changes for all tables, you can use the `skippedTables` option in the `config.json` to ignore any tables.

## :cyclone: Simple Installation
You can download the pre-compiled binaries from the Github [releases](https://github.com/Niyko/Athena/releases) page and copy them to the desired location. After that you can follow the below steps in order.

#### Create a `config.json` file in the root folder where you but the Athena binary. Here is the format of the JSON file. Fill all the credentials also.
You can find more details about the paramters in config file in below sections.
`````json
{
    "dbHost": "",
    "dbPort": 1433,
    "dbUser": "",
    "dbPassword": "",
    "dbName": "",

    "kafkaHost": "",
    "kafkaEnableTLS": false,
    "kafkaTopic": "",

    "kafkaSASLMechanisms": "NONE",
    "kafkaSASLUsername": "",
    "kafkaSASLPassword": "",

    "pollInterval": 10,
    "fetchLimit": 50,
    "skippedTables": [],

    // If you want to collect logs in clickhouse
    "clickHouse": true,
    "clickHouseHost": "<host>:<port>",
    "clickHouseUsername": "",
    "clickHousePassword": "",
    "clickHouseDatabase": "",
    "clickHouseTableName": "",
    "clickHouseTableTTL": 12
}
`````

#### Create topic with name given in `config.json` without scheme and with required partitions in you Kafka broker.

#### Run the setup command in order to create the CDC in database and other required setup.

`````bash
./athena setup
`````

#### Setup a service for running Athena in the background. Setting this up will different for Windows and Linux. Below given are the steps to create them on a Linux distro.

#### Create a service file called `athena_mssql_kafka.service` in the directory `/etc/systemd/system` using the following commands.

`````bash
cd /etc/systemd/system
nano athena_mssql_kafka.service
`````

#### Copy and paste the below contents to the above created service file `athena_mssql_kafka.service`.

`````s
[Unit]
Description=Athena MSSQL Kafka Service
After=network.target

[Service]
Type=simple
ExecStart=athena run

[Install]
WantedBy=multi-user.target
`````

> Please note that path in `ExecStart` needs to change while creating the service file.

#### Now you can start the service and also check the status of the service.

`````bash
systemctl start athena_mmsql_kafka.service
systemctl status athena_mmsql_kafka.service
`````

## :gear: Configuring Athena
Athena can be configured using the `config.json` file created on the root the Athena binary. Here are the details of the configuration keys and what they do in table format.

| Option | Description | Example |
| --- | --- | --- |
| `dbHost` | aaaaaaaaaaa | 127.0.0.1 |
| `dbPort` | aaaaaaaaaaa | 1433 |
| `dbUser` | aaaaaaaaaaa |  |
| `dbPassword` | aaaaaaaaaaa |  |
| `dbName` | aaaaaaaaaaa |  |
| `kafkaHost` | aaaaaaaaaaa |  |
| `kafkaTopic` | aaaaaaaaaaa |  |
| `kafkaEnableTLS` | aaaaaaaaaaa | `true`, `false` |
| `kafkaSASLMechanisms` | aaaaaaaaaaa | `NONE`, `SASL-PLAIN`, `SASL-SCRAM-SHA-256`, `SASL-SCRAM-SHA-512` |
| `kafkaSASLUsername` | aaaaaaaaaaa |  |
| `pollInterval` | Interval where next polling to the database is made. It's given in seconds format. | 10 |
| `fetchLimit` | Number of CDC changes rows that will be pulled from the table at once. | 50 |
| `skippedTables` | Array of tables that needs to skipped while taking CDC changes. | ["table1", "table2"] |
| `clickHouse` | aaaaaaaaaaa | 10 |
| `clickHouseHost` | aaaaaaaaaaa | 10 |
| `clickHouseUsername` | aaaaaaaaaaa | 10 |
| `clickHousePassword` | aaaaaaaaaaa | 10 |
| `clickHouseDatabase` | aaaaaaaaaaa | 10 |
| `clickHouseTableName` | aaaaaaaaaaa | 10 |
| `clickHouseTableTTL` | aaaaaaaaaaa | 10 |

## :hammer_and_wrench: How to build
You can build the binaries or do development of Athena by following the below steps. Athena is build fully on Golang. So you should install latest version of Go from [here](https://go.dev/doc/install). Do note that building binaries are managed with the [Goreleaser](https://goreleaser.com/).

* Clone that project from Github.
* Run `go mod download` command to install all mods.
* Run `SET GORUN=true` command to set gorun variable.
* Run the command `goreleaser release --snapshot --clean` for building the binaries.

## :page_with_curl: License
Athena is licensed under the [MIT License](https://github.com/Niyko/Athena/blob/main/LICENSE).