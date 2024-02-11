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

    "kafkaHost": "",
    "kafkaSASLMechanisms": "PLAIN",
    "kafkaSecurityProtocol": "SASL_SSL",
    "kafkaSASLUsername": "",
    "kafkaSASLPassword": "",
    "kafkaTopic": "",

    "pollInterval": 10,
    "fetchLimit": 50,
    "mssqlCDCRetentionPeriod": 1440,
    "skippedTables": []
}
`````

#### Create topic with name given in `config.json` without scheme and with required partitions in you Kafka broker.

#### Run the setup command in order to create the CDC in database and other required changes (Use athena.exe for Windows binaries).

`````bash
athena setup
`````

