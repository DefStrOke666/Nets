## Usage

```bash
$ cd server/
$ make
$ ./server -p [port]
```
```bash
$ cd client/
$ make
$ ./client -a [server address] -p [server port] -f [file]
```

## Working example

```bash
$ ./client -a 46.243.142.237 -p 8888 -f [file]
```
```bash
$ ssh lab2@46.243.142.237
password: lab2
$ cd lab2/server/
$ cat log.txt
$ cd uploads/
```