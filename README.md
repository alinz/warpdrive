# Warpdrive

### Postgres Setup

before running warpdrive, make sure you have a right role and database. you can run the follwoing sql in Postgres terminal

```bash
CREATE USER warpdrive WITH PASSWORD 'warpdrive';
CREATE DATABASE warpdrivedb;
```

and make sure to set the correct username, password and database in warpdrive.conf.