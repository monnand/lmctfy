CC ?= gcc
AR ?= ar
CLMCTFY_ROOT = ../..
CFLAGS = -I$(CLMCTFY_ROOT)/clmctfy/include
CFLAGS += -I$(CLMCTFY_ROOT)/include
CFLAGS += -I$(CLMCTFY_ROOT)

OBJS = clmctfy.o mock.o
LIBCLMCTFY = libclmctfy.a
all: $(LIBCLMCTFY)

$(LIBCLMCTFY): $(OBJS)
	$(AR) $(ARFLAGS) $@ $^
%.o: %.c
	$(CC) -c $(CFLAGS) -o $@ $<
clean:
	rm -f *.o
	rm -f $(LIBCLMCTFY)
.PHONY: all clean
