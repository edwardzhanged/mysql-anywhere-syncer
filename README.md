# MySQL-Anywhere-Syncer
Sync Mysql changes to any other database

---------------------------------------

## Features

* Real-time Sync: MySQL-Anywhere-Syncer captures changes in the MySQL database in real-time and syncs them to the target database.
* Configure real-time updates: The application can reload the configuration in response to changes in the configuration file without the need for a restart.
* Broad Compatibility: It can sync changes to any other database, providing a versatile solution for various database platforms. (Developing, support MongoDB only)
* Efficient and Reliable: It ensures data consistency and integrity during the synchronization process.

---------------------------------------
## Installation
### Build From Source
```bash
git clone https://github.com/edwardzhanged/mysql-anywhere-syncer.git
cd mysql-anywhere-syncer
make build
```
### Download Binary
Download From [Release](https://github.com/edwardzhanged/mysql-anywhere-syncer/releases/tag/v1.0.0)

## Configuration
The default configuration file is named app.yml. Also, you have the option to specify a different file name using the --config flag. 

An example configuration can be found within app.yml, which illustrates its primary purpose. You are also able to define multiple rules as needed.

### MongoDB
[MongoDB Configuration](docs/mongodb.md)

### Redis
TBD

## Run
If you only want to sync from latest MySQL position
```bash
./mysql-anywhere-syncer
```
If you want to dump the existing MySQL data to targets before auto sync
```bash
./mysql-anywhere-syncer --dump
```
Hung up
```bash
nohup ./mysql-anywhere-syncer > output.log 2>&1 
```

## Contribution
Implement the interface within the syncer package to support various different databases.
