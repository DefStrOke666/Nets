CC=g++
CFLAGS=-I.
DEPS = client.h clientClass.h exceptions.h
OBJ = client.o

%.o: %.c $(DEPS)
	$(CC) -c -o $@ $< $(CFLAGS)

client: client.o
	$(CC) -o $@ $^ $(CFLAGS)