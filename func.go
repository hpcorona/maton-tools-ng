package main

import (
	"fmt"
	"github.com/hpcorona/go-v8"
)

func loadFunctions(ctx *v8.V8Context) {
	// operating system functions
	load_os_functions(ctx)

	// file system functions
	load_fs_functions(ctx)
	
	// path functions
	load_path_functions(ctx)

  // ng functions
	load_ng_functions(ctx)
	
	// mustache functions
	load_mustache_functions(ctx)
	
	// project management functions
	load_library_functions(ctx)
	
	// android generator
	load_android_functions(ctx)
}

func paramCount(value []interface{}, pcount int, fname string) {
  if value == nil && pcount == 0 {
    return
  }

	if len(value) != pcount {
		panic(fmt.Sprintf("%s only supports %d parameter(s); %d passed\n{%v}\n", fname, pcount, len(value), value))
	}
}

func paramMin(value []interface{}, pmin int, fname string) {
  if value == nil && pmin == 0 {
    return
  }

	if len(value) < pmin {
		panic(fmt.Sprintf("%s must have a minimum of %d parameter(s); %d passed\n{%v}\n", fname, pmin, len(value), value))
	}
}

func paramMax(value []interface{}, pmax int, fname string) {
  if value == nil && pmax == 0 {
    return
  }

	if len(value) > pmax {
		panic(fmt.Sprintf("%s must have a maximum of %d parameter(s); %d passed\n{%v}\n", fname, pmax, len(value), value))
	}
}
