# Simple gRPC micro-services application for reading and saving records from file to DB

Features:

* gRPC endpoints
* Incremental reading of CSV file
* Graceful shutdown
* Creates a new record in a database, or updates the existing one


## How to launch the application

Start PostgreSQL database server. In the PostgreSQL database named `postgres` (it is a default one) execute the SQL statements given in the file `migrate/migration.sql`.
Database connection information (see file `config/app.yaml`):
* server address: `127.0.0.1` (localhost)
* server port: `5432`
* database name: `postgres`
* username: `postgres`
* password: `postgres`

Install the application from the Terminal:
```shell
go get github.com/nettyrnp/go-grpc
```

Run the Persistence Service from the Terminal #1:
```shell
go run services/persistor/persistor.go
```

Run the Ingestor Service from the Terminal #2:
```shell
go run services/ingestor/ingestor.go
```

Now in the Terminals you can observe the logs describing the communication between the three services (including Postgres):

```shell
# Terminal #1:
2018/08/07 01:49:43 sql: INSERT INTO public.people(id, name, email, mobile_number) VALUES (33, 'Damian', 'dolor@cursus.com', '(+44)01699955892')
2018/08/07 01:49:43 sql: UPDATE public.people SET name='Damien', email='dolor@cursus.com', mobile_number='(+44)01699955892' WHERE id=33
...
2018/08/07 01:49:44 sql: INSERT INTO public.people(id, name, email, mobile_number) VALUES (88, 'Lev', 'porttitor.vulputate@velitegetlaoreet.ca
', '(+44)0138796288')
2018/08/07 01:49:44 sql: INSERT INTO public.people(id, name, email, mobile_number) VALUES (89, 'Clark', 'commodo.at@sagittisDuisgravida.net',
'(+44)01168557827')
...
# should print a list of executed sql queries

```shell
# Terminal #2:
2018/08/07 01:49:41 Reading: loaded 39 non-duplicate lines
2018/08/07 01:49:41 Saving: created 38, updated 1 records in DB
2018/08/07 01:49:42 Reading: loaded 41 non-duplicate lines
2018/08/07 01:49:42 Saving: created 40, updated 1 records in DB
2018/08/07 01:49:43 Reading: loaded 22 non-duplicate lines
2018/08/07 01:49:43 Saving: created 22, updated 0 records in DB
...
# should print the stats on Reading the CSV file & Saving to DB operations
