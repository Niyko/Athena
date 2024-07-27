![Athena logo](https://i.imgur.com/UQggP60.png)
Athena is a tool used to move change data capture from Microsoft SQL to Apache Kafka. Athena is written in Golang. Athena supports SASL authentication-based Kafka broker servers. Other tools like Debezium are available, but are a pain in the head to set up and manage them with MSSQL. Debezium gives you very little room for configuration while using it as a connector from services like confluent.io. Athena is very simple to set up and can be managed easily without any unwanted complications.

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

    // If you want to collect logs in clickhouse
    "clickHouse": true,
    "clickHouseHost": "<host>:<port>",
    "clickHouseUsername": "",
    "clickHousePassword": "",
    "clickHouseDatabase": "",
    "clickHouseTableName": "",
    "clickHouseTableTTL": 12,

    "kafkaHost": "",
    "kafkaSASLMechanisms": "PLAIN",
    "kafkaSecurityProtocol": "SASL_SSL",
    "kafkaSASLUsername": "",
    "kafkaSASLPassword": "",
    "kafkaTopic": "",

    "pollInterval": 10,
    "fetchLimit": 50,
    "skippedTables": []
}
`````

#### Create topic with name given in `config.json` without scheme and with required partitions in you Kafka broker.

#### Run the setup command in order to create the CDC in database and other required changes (Use athena.exe for Windows binaries).

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
Athena can be configured using the `config.json` file created on the root the Athena binary. Here are the details of the configuration keys and what they do in table format. Please not that MSSQL and Kafka connection options are not included on the table.

| Option | Description | Example |
| --- | --- | --- |
| `pollInterval` | Interval where next polling to the database is made. It's given in seconds format. | 10 |
| `fetchLimit` | Number of CDC changes rows that will be pulled from the table at once. | 50 |
| `skippedTables` | Array of tables that needs to skipped while taking CDC changes. | ["table1", "table2"] |
| `clickHouse` | Enable this if logging to Clickhouse is needed. | true/false |

## :hammer_and_wrench: How to build
You can build the binaries or do development of Athena by following the below steps. Athena is build fully on Golang. So you should install latest version of Go from [here](https://go.dev/doc/install). Do note that building binaries are managed with the [Goreleaser](https://goreleaser.com/).

* Clone that project from Github.
* Run `go mod download` command to install all mods.
* Run `SET GORUN=true` command to set gorun variable.
* Run `SET SENTRYDNS={sentry dns}` command to set sentry dns variable.
* Run the command `goreleaser release --snapshot --clean` for building the binaries.

## :page_with_curl: License
Athena is licensed under the [GNU GENERAL PUBLIC LICENSE](https://github.com/Niyko/Athena/blob/main/LICENSE).