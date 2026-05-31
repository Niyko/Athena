![Athena logo](https://i.imgur.com/UQggP60.png)
Athena is a lightweight Change Data Capture (CDC) solution that streams changes from Microsoft SQL Server to Apache Kafka. Built in Golang, it supports SASL-authenticated Kafka brokers and provides a straightforward setup experience. Unlike alternatives such as Debezium, which can be complex to configure and manage. Athena offers greater simplicity and operational ease. It automatically manages CDC setups, publishes database changes to a single Kafka topic, and delivers a clean, intuitive event format that is easy for downstream consumers to understand and process.

## :zap: How things work

- Creates a message for changes like `create`, `update`, `delete` for rows in MSSQL database tables to a single Kafka topic.
- Athena only creates messages for all new table changes. Existing ones are ignored.
- Kafka topic have to be created before hand. Unlike Debezium, Athena will not create the topic own it's own.
- All CDC setups in MSSQL is automatically done by Athena when `setup` command is run.
- By default, Athena will poll for changes for all tables, you can use the `skippedTables` option in the `config.json` to ignore any tables.

## :cyclone: Simple Installation
You can download the pre-compiled binaries from the Github [releases](https://github.com/Niyko/Athena/releases) page and copy them to the desired location. After that you can follow the below steps in order.

#### Create a `config.json` file in the root folder where you but the Athena binary. Here is the format of the JSON file. Fill all the credentials also.
You can find more details about the paramters in config file in below sections.
`````json
{
    "dbHost": "127.0.0.1",
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
    "clickHouse": false,
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
| `dbHost` | Database host of MSSQL | 127.0.0.1 |
| `dbPort` | Database port of MSSQL | 1433 |
| `dbUser` | Username for the MSSQL database |  |
| `dbPassword` | Password for the MSSQL database |  |
| `dbName` | Database name of MSSQL |  |
| `kafkaHost` | Host with port for the Kafka server  | 127.0.0.1:9092 |
| `kafkaTopic` | Kafka topic that you created for table changes to show |  |
| `kafkaEnableTLS` | Enables TLS for Kafka connection | `true`, `false` |
| `kafkaSASLMechanisms` | SASL mechanism that need to be used for Kafka connection | `NONE`, `SASL-PLAIN`, `SASL-SCRAM-SHA-256`, `SASL-SCRAM-SHA-512` |
| `kafkaSASLUsername` | SASL user name of the Kafka server |  |
| `kafkaSASLPassword` | SASL password of the Kafka server |  |
| `pollInterval` | Interval where next polling to the database is made. It's given in seconds format. | 10 |
| `fetchLimit` | Number of CDC changes rows that will be pulled from the table at once. | 50 |
| `skippedTables` | Array of tables that needs to skipped while taking CDC changes. | ["table1", "table2"] |
| `clickHouse` | Enable Clickhouse logs. Table and struture for Clickhouse is automatically created by Athena when `setup` command is run | `true`, `false` |
| `clickHouseHost` | Host with port for the Clickhouse server | 127.0.0.1:8123 |
| `clickHouseUsername` | User name of Clickhouse server |  |
| `clickHousePassword` | Password of Clickhouse server |  |
| `clickHouseDatabase` | Dasebase name of Clickhouse server |  |
| `clickHouseTableName` | Table name of Clickhouse server |  |
| `clickHouseTableTTL` | Time to live for each record in hours | 24 |

## :mushroom: Helper options in Athena
Athena executable have some other helper functions apart from `setup` or `run` which are explained below. These can be run like eg: `./athena uninstall`

| Option | Description |
| --- | --- |
| `uninstall` | Will disable CDC in MSSQL database and remove the SQlite database |
| `add-cdc` | Will run CDC setup in the MSSQL database |
| `remove-cdc` | Will disable CDC in MSSQL database |
| `clear-cdc-history` | Clear CDC history or changes that Athena didn't process yet from the MSSQL database |
| `recreate-clickhouse` | Rerun the Clickhouse migration |
| `recreate-sqlite` | Recreate the SQlite database and rerun the migration |
| `help` | To view all the options available |

## :triangular_ruler: Development
For setting up development environment, there is a docker file in the folder `dev`. It will create all necessary services like MSSQL with sample database, Kafka etc. This same environment can be used for running integration tests.

* Install latest version of Go from [here](https://go.dev/doc/install).
* Clone that project from Github.
* Run `go mod download` command to install all mods.
* Then run the commands below as needed.

`````bash
cd dev
docker composer up -d
cd ..
set GORUN=true # Used for identifying if script is run from go run command to choose correct path for reading config.json or db.sqlite
go run . setup
go run . run
`````

> Please note that the `docker-compose.yml` in the `dev` folder should only be used for development purpose.

## :cactus: How to run tests
Before running the tests, make sure you have setup the development environment and also `config.json` is setup correctly.

`````bash
cd tests
go test -v -run TestIntegration
`````

## :hammer_and_wrench: How to build
You can build the binaries or do development of Athena by following the below steps. Athena is build fully on Golang. So you should install latest version of Go from [here](https://go.dev/doc/install). Do note that building binaries are managed with the [Goreleaser](https://goreleaser.com/).

* Clone that project from Github.
* Run `go mod download` command to install all mods.
* Run `SET GORUN=true` command to set gorun variable.
* Run the command `goreleaser release --snapshot --clean` for building the binaries.

## :page_with_curl: License
Athena is licensed under the [MIT License](https://github.com/Niyko/Athena/blob/main/LICENSE).