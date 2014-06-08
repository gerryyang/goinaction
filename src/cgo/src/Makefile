# Makefile

CC     =  gcc
CXX    =  g++
RANLIB =  ranlib
AR     =  ar
AROPT  =  -scurv

INCLUDE =  -I../inc
LIBPATH =  -L../lib
CFLAGS  =  -Wall -Os -g -pipe -D_REENTRANT

LIBS = -lhello

OBJS_HELLO =  hello.o
BIN_DEMO   =  demo
OBJS_DEMO  =  main.o

LIB_HELLO_STATIC  =  ../lib/libhello.a
LIB_HELLO_DYNAMIC =  ../lib/libhello.so

LIBRARY = STATIC
ifeq ($(LIBRARY), STATIC)
LIB_HELLO = $(LIB_HELLO_STATIC)
else
CFLAGS += -fPIC
LIB_HELLO = $(LIB_HELLO_DYNAMIC)
endif

LIB_API  =  $(LIB_HELLO)
OBJS_API =  $(OBJS_HELLO)

all: $(LIB_API)
bin: $(BIN_DEMO)

$(LIB_HELLO_STATIC): $(OBJS_HELLO)
	$(AR) $(AROPT) $@ $? 
	$(RANLIB) $@

$(LIB_HELLO_DYNAMIC): $(OBJS_HELLO)
	$(CC) -shared -o $@ $?

$(BIN_DEMO): $(OBJS_DEMO) $(LIB_HELLO)
	$(CXX) -o $@ $^ $(LIBPATH) $(LIBS)
	@echo "build $(BIN_DEMO) OK!"


install:
	@echo "nothing to install"
		
.PHONY: clean
clean:
	rm -f $(OBJS_API) $(LIB_API) 
	rm -f $(OBJS_DEMO) $(BIN_DEMO)



%.o: %.cpp 
	$(CXX) $(CFLAGS) $(INCLUDE) -c $<

%.o: %.c
	$(CC) $(CFLAGS) $(INCLUDE) -c $<


