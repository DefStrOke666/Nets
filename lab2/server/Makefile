CC=g++
CFLAGS=-I.
DEPS = server.h serverClass.h uploadThread.h exceptions.h
OBJ = server.o

%.o: %.c $(DEPS)
	$(CC) -c -o $@ $< $(CFLAGS)

server: server.o
	$(CC) -o $@ $^ $(CFLAGS) -lpthread