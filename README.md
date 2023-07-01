# eth tweet
Ethereum Tweet Service

The technical principle is to use Ethereum personal_sign to sign the tweet information, and then submit it to ipfs for storage. The front end can also read the information on ipfs to verify the signature to ensure the information is accurate.

Because the personal_sign signature is used, and it is a plaintext signature, it will not be used for malicious attacks.

This program provides http api for the front-end program to call to realize the separation of data and display.

Local data persistence supports sqlite and MySQL, but users can use it with a wallet without using the client.

online experience https://app.ethtweet.io/#/home

http api docs  [api](api.md)

## principle

Use the wallet to sign the tweet information, then submit it to ipfs, and then broadcast it to other nodes, and other nodes save and index it, and then submit it to ipfs.

The tweets of each address have an auto-incrementing nonce value, starting from 0, which is used to mark the order in which tweets are published, and there is also a sequence when nodes pull data.

You can use the web front end, and then call the rpc interface.

## config

By default, the `tweet.yaml` file is read, and the command-line arguments are read if it does not exist.

## mysql configuration

Modify the configuration file ```tweet.yaml``` to correctly configure the MySQL connection information.

In the window environment, if there is a MySQL 8.0 in the program running directory, it will automatically start MySQL, the required files and their paths
```
mysql\bin\mysqld.exe
mysql\bin\libprotobuf-lite.dll
```

start node
```
./EthTweet -config "./tweet.yaml"
```


## docker

Run the test, close the container data to automatically clear

```shell
docekr run  --rm -it -p 8080:8080 -p 4001:4001/udp -p 4001:4001/tcp chenjia404/ethtweet
```

save data run
```shell
docekr run -it -v ./databases:/databases -v ./keyStore:/keyStore -p 8080:8080 -p 4001:4001/udp -p 4001:4001/tcp chenjia404/ethtweet
```


## docker-compose

Under the docker-compose directory, an environment with MySQL is integrated, which can be started with one click.

## todo 
Add node statistics, record the longest online time, and connect every time you start

The function of encrypting private messages can generate encrypted private messages, and only the address holder can decrypt and read them.

app side

pc version web

Add attachment support: you can upload zip-like files and download them

Added attachment seeding: node reseeding