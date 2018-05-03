# DynamoDB backup Tool
This tool does two things:

* Backup of one existing AWS dynamodb table into a local folder. It will **NOT** backup the schema or other information (just the data).
* Restore of one existing backup to one single existing AWS dynamodb table.

Performance was not considered during development. Exports/Imports may take a while. Especially with large tables.

## Installation
Have a look at the [release-page.](https://github.com/michelvocks/dynamodb-backup-tool/releases)

Download and copy the single binary to the preferred place.

## How do I use it?
Export of table:
```bash
dynamodb-backup-tool -mode "export" -table "MyDynamoDBTable" -data "/path/to/data/folder"
```

Import of table:
```bash
dynamodb-backup-tool -mode "restore" -table "MyDynamoDBTable" -data "/path/to/data/folder"
```