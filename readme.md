## Queryhw

A program for benchmarking time_bucket("1 minute") queries in timescaledb
with a configurable number of parallel worker threads.

Queryhw is written in Go.

Usage of ./queryhw:

    -d string
        database connection string for timescaledb, see docs for lib/pq
        (default "postgres://postgres:xxx@db/homework?sslmode=disable")
    -f string
        the path to a CSV file containing the queries to run 
        (default "-" read CSV from STDIN)
    -n int
        the number of concurrent workers to run 
        (default GOMAXPROCS - number of hardware threads on the machine)
    -v
        print more verbose output as the program runs

## How to run queryhw

### Prerequisites

You'll need docker and docker-compose. 
These instructions are tested on Linux and Mac OS with docker=20.10.12 (community) and docker-compose 1.29.2
See https://docs.docker.com/engine/install/ and https://docs.docker.com/compose/install/
for how to install docker and docker-compose on your platform.

### Build

Clone the git repo and cd into it.

    git clone git@github.com:eloff/queryhw.git
    cd queryhw

Create a pgdata directory. **Important do not skip this step.**
Unfortunately git doesn't allow an empty directory and Postgres doesn't
like a directory containing any files, including hidden dot-prefixed files
like .keep
    
    mkdir pgdata

Then start the Go container and connect to an interactive session:

    docker-compose run --rm app

List the directory. You should see main.go, querytool/, this readme, etc.

Run the tests:

    go test ./...
    ?       github.com/eloff/queryhw        [no test files]
    ok      github.com/eloff/queryhw/querytool      0.002s

Build the program:

    go build

Check the usage:

    ./queryhw -h

Run the program:

    ./queryhw < data/query_params.csv

Experiment by running queryhw with different
input sources and values for -n (number of workers.)

### Troubleshooting

Because nothing ever seems to work quite like it's supposed to.

If you get an error like this, the database is still initializing, give it more time:

    error running query: dial tcp db:5432: connect: connection refused

Another cause of a similar error was using an outdated version of docker.
Updating docker, rebooting, and then following the instructions below
to recreate the docker containers from scratch solved it.

If you get an error like:

    error running query: dial tcp: lookup db: Temporary failure in name resolution

I solved this by following the instructions to recreate the docker container
from scratch below, and then running:

    sudo service docker stop
    sudo service docker start

If troubleshooting, you can verify the database was started and setup correctly by running

    psql -h 127.0.0.1 -p 5438 -U postgres -d homework
    select count(*) from cpu_usage;

This requires psql installed on your local machine. There should be 345600 rows.
It's also possible to use docker to attach an interactive shell to the db service
and use psql from there (but I expect anyone at Timescale has it installed!)

If something goes wrong you may need to recreate the docker containers from scratch.

You can do this by running:
    
    docker-compose down
    docker-compose rm
    sudo rm -rf pgdata
    mkdir pgdata
