# MySQL-Anywhere-Syncer
Sync Mysql changes to any other database

---------------------------------------

## Features

* Real-time Sync: MySQL-Anywhere-Syncer captures changes in the MySQL database in real-time and syncs them to the target database.
* Configure real-time updates: The application can reload the configuration in response to changes in the configuration file without the need for a restart.
* Broad Compatibility: It can sync changes to any other database, providing a versatile solution for various database platforms. (Developing, support MongoDb only)
* Efficient and Reliable: It ensures data consistency and integrity during the synchronization process.

---------------------------------------
## Usage
#### Build
```bash
make build
```
#### Configuration
The default configuration file is named app.yml. Also, you have the option to specify a different file name using the --config flag. 

An example configuration can be found within app.yml, which illustrates its primary purpose. You are also able to define multiple rules as needed.

#### Run
If you only want to sync from latest MySQL position
```bash
./mysql-anywhere-syncer
```
If you want to dump the existing MySQL data to targets before auto sync
```bash
./mysql-anywhere-syncer --dump
```


## Benchmark


## Contribution
Implement the interface within the syncer package to support various different databases.


benchmark
mongodb connection没有close
自己总结点感悟吧
