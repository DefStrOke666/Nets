CC=g++
CFLAGS=-I.
DEPS = ipv4.h ipv6.h utils.h exceptions.h
OBJ = multicast.o

%.o: %.c $(DEPS)
	$(CC) -c -o $@ $< $(CFLAGS)

multicast: multicast.o
	$(CC) -o $@ $^ $(CFLAGS)
