# Database deployment

The `zwischenton_identity_database` contains all users that have access to the ZwischentonCloud backend or
parts of it. The default user is called administrator with the password set to 'we4711'.
It also contains all API keys.

The `zwischenton_cloud_database` is used by the ZwischentonCloud for persistently storing zwischenton data.

## Server deployment

The [install script](../operation/install.sh) will install and secure the database.

### MYSQL cheatsheet

```bash
brew services start mysql
brew services restart mysql
brew services stop mysql
```

```mysql
SHOW DATABASES;
USE database;
SHOW TABLES;
SELECT * FROM table
```
