package main

/*
#include <stdio.h>

void hello() {
    printf("Hello, C!\n");
}

*/
import "C"

func main() {
	C.hello()
}
