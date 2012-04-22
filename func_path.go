package main

import (
  "strings"
	"github.com/hpcorona/go-v8"
	"path/filepath"
)

func load_path_functions(ctx *v8.V8Context) {
	// path functions
	ctx.AddFunc("_path_clean", path_clean)
	ctx.AddFunc("_path_ext", path_ext)
	ctx.AddFunc("_path_isAbs", path_isAbs)
	ctx.AddFunc("_path_join", path_join)
	ctx.AddFunc("_path_split", path_split)
	ctx.AddFunc("_path_splitList", path_splitList)

	ctx.Eval(`
    this.path = {
      "clean": function(args) { return _path_clean(args); },
      "ext": function(args) { return _path_ext(args); },
      "isAbs": function(args) { return _path_isAbs(args); },
      "join": function() { return _path_join.apply(this, arguments); },
      "split": function(args) { return _path_split(args); },
      "splitList": function(args) { return _path_splitList(args); }
    };
  `)
}

func path_clean(value ...interface{}) interface{} {
	paramCount(value, 1, "path.clean")

	return filepath.Clean(value[0].(string))
}

func path_ext(value ...interface{}) interface{} {
	paramCount(value, 1, "path.ext")

	return filepath.Ext(value[0].(string))
}

func path_isAbs(value ...interface{}) interface{} {
	paramCount(value, 1, "path.isAbs")

	return filepath.IsAbs(value[0].(string))
}

func path_join(value ...interface{}) interface{} {
	params := make([]string, len(value))
	for i := 0; i < len(value); i++ {
		params[i] = value[i].(string)
	}

	return filepath.Join(params...)
}

func path_split(value ...interface{}) interface{} {
	dir, file := filepath.Split(value[0].(string))

  idx := strings.LastIndex(file, ".")
  base := file
  ext := ""

  if idx >= 0 {
    base = file[0 : idx]
    ext = file[idx : ]
  }

	return []string{dir, file, base, ext}
}

func path_splitList(value ...interface{}) interface{} {
	return filepath.SplitList(value[0].(string))
}
