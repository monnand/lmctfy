CC ?= gcc
AR ?= ar
CLMCTFY_ROOT = ../..
CFLAGS = -I$(CLMCTFY_ROOT)/clmctfy/include
CFLAGS += -I$(CLMCTFY_ROOT)/include
CFLAGS += -I$(CLMCTFY_ROOT)

OBJS = clmctfy_raw.o clmctfy.o clmctfy_mock.o
LIBCLMCTFY = libclmctfy.a
all: $(LIBCLMCTFY)

$(LIBCLMCTFY): $(OBJS)
	$(AR) $(ARFLAGS) $@ $^
%.o: %.c clmctfy-mock.h macros.h
	$(CC) -c $(CFLAGS) -o $@ $<
clean:
	rm -f *.o
	rm -f $(LIBCLMCTFY)
.PHONY: all clean
