# Demo

These files create an ElasticSearch backend for OGRT and
some example programs, with which you can play around.

1. Install the client into the demo directory

    cd ../client
    ./vendorize
    ./configure --prefix=$(pwd)/../demo/ogrt
    make install

2. Spin up the ElasticSearch infrastructure (requires Docker)

    ./infra-init

3. Add an ElasticSearch output to the servers ogrt.conf

    [Outputs.ElasticSearch]
    Type = "JsonElasticSearch5"
    Workers = 2
    Params = "http:localhost:9200:ogrt"

4. Start ogrt-server

5. Generate a software tree

    ./generate-software

6. Enable OGRT

    . setup.sh

7. Push example jobs to ElasticSearch

    ./generate-jobs

8. Check OGRT Web Interface on http://localhost:8080/web/ to see if
   server received messages

9. Check Kibana on http://localhost:5601/ for the data. The index name
   is 'ogrt'


