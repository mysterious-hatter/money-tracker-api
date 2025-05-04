# Some helpful information for CI/CD
## Open Postgesql in the Docker container
Go inside of the container
```console
docker exec -it <containerId> bash
```
Connect to the database
```console
psql -U <username> -d <database>
```
## View table schema
```sql
\d <tableName>
```
## View logs of the database
Once you are in the container, go to the directory with logs
```console
cd  /var/lib/postgresql/data/log
```
Type ```ls``` to see all the files in the directory
```console
ls
```
Then open the file you need
I use this pattern for the filenames: ```postgresql-%Y-%m-%d_%H%M%S.log```
```console
cat <your-log-file.log>
```
## Turn the logging on
Copy the configuration file from your databese's container to the local direcotry
```console
docker cp <containerName>s:var/lib/postgresql/data/postgresql.conf ./postgresql.conf
```
Open the configuration file ```postgresql.conf``` to edit
```console
nano postgresql.conf
```
Set the following settings
```ini
# Enable logging
logging_collector = on

# Log all queries
log_statement = 'all'

# uncomment this
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_destination = 'stderr'
```
Send the modified settings to the container
```console
docker cp ./postgresql.conf <containerName>:var/lib/postgresql/data/postgresql.conf
```
Restart the container
```console
docker restart <containerName>
```
