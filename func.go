package main

import (
	"fmt"
	"github.com/hpcorona/go-v8/v8"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func loadFunctions(ctx *v8.V8Context) {
	// operating system functions
	ctx.AddFunc("_os_log", os_log)
	ctx.AddFunc("_os_env", os_env)
	ctx.AddFunc("_os_exit", os_exit)
	ctx.AddFunc("_os_run", os_run)
	ctx.AddFunc("_os_findCmd", os_findCmd)

	ctx.Eval(`
  this.os = 
    {
      "log": function() { _os_log.apply(this, arguments); },
      "env": function(args) { return _os_env(args); },
      "exit": function(args) { _os_exit(args); },
      "run": function() { return _os_run.apply(this, arguments); },
      "findCmd": function(args) { return _os_findCmd(args); }
    };
  `)

	// file system functions
	ctx.AddFunc("_fs_cd", fs_cd)
	ctx.AddFunc("_fs_cp", fs_cp)
	ctx.AddFunc("_fs_mv", fs_mv)
	ctx.AddFunc("_fs_touch", fs_touch)
	ctx.AddFunc("_fs_rm", fs_rm)
	ctx.AddFunc("_fs_mkdir", fs_mkdir)
	ctx.AddFunc("_fs_rmdir", fs_rmdir)
	ctx.AddFunc("_fs_tempdir", fs_tempdir)
	ctx.AddFunc("_fs_symlink", fs_symlink)
	ctx.AddFunc("_fs_link", fs_link)
	ctx.AddFunc("_fs_truncate", fs_truncate)
	ctx.AddFunc("_fs_isDir", fs_isDir)
	ctx.AddFunc("_fs_ls", fs_ls)

	ctx.Eval(`
  this.fs =
    {
      "cd": function(args) { _fs_cd(args); },
      "cp": function(v0, v1) { _fs_cp(v0, v1); },
      "mv": function(v0, v1) { _fs_mv(v0, v1); },
      "touch": function(args) { _fs_touch(args); },
      "rm": function(args) { _fs_rm(args); },
      "mkdir": function(args) { _fs_mkdir(args); },
      "rmdir": function(args) { _fs_rmdir(args); },
      "tempdir": function() { return _fs_tempdir(); },
      "symlink": function(v0, v1) { _fs_symlink(v0, v1); },
      "link": function(v0, v1) { _fs_link(v0, v1); },
      "truncate": function(v0, v1) { _fs_truncate(v0, v1); },
      "isDir": function(args) { return _fs_isDir(args); },
      "ls": function(args) { return _fs_ls(args); }
    };
  `)

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

	return []string{dir, file}
}

func path_splitList(value ...interface{}) interface{} {
	return filepath.SplitList(value[0].(string))
}

func os_log(value ...interface{}) interface{} {
	for i := 0; i < len(value); i++ {
		fmt.Printf("%v", value[i])
	}

	fmt.Printf("\n")

	return true
}

func os_findCmd(value ...interface{}) interface{} {
	paramCount(value, 1, "os.findCmd")

	path, err := exec.LookPath(value[0].(string))
	if err != nil {
		panic(err.Error())
	}

	return path
}

func os_run(value ...interface{}) interface{} {
	paramMin(value, 1, "os.run")

  fmt.Printf("Run '%s'", value[0].(string))

  strs := make([]string, len(value) - 1)
  for i := 0; i < len(value) - 1; i++ {
    strs[i] = value[i+1].(string)
    fmt.Printf(" %s", strs[i])
  }
  fmt.Println()

	cmd := exec.Command(value[0].(string), strs...)
	err := cmd.Run()
	if err != nil {
		panic(err.Error())
	}

	out, err := cmd.CombinedOutput()
	if err == nil {
		fmt.Printf("%s\n", string(out))
	}

  return true
}

func os_env(value ...interface{}) interface{} {
	paramCount(value, 1, "os.env")

	return os.Getenv(value[0].(string))
}

func os_exit(value ...interface{}) interface{} {
	paramCount(value, 1, "os.exit")

  fmt.Printf("System exit '%d'.\n", value[0].(int))

	os.Exit(value[0].(int))

	return true
}

func fs_cd(value ...interface{}) interface{} {
	paramCount(value, 1, "fs.cd")

  fmt.Printf("Change current directory to '%s'.", value[0].(string))

	os.Chdir(value[0].(string))

	return true
}

func fs_cp(value ...interface{}) interface{} {
	paramCount(value, 2, "fs.cp")

	src := value[0].(string)
	dst := value[1].(string)

  fmt.Printf("Copy '%s' to '%s'.\n", src, dst)

	sf, err := os.Open(src)
	if err != nil {
		panic(err.Error())
	}
	defer sf.Close()
	df, err := os.Create(dst)
	if err != nil {
		panic(err.Error())
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_mv(value ...interface{}) interface{} {
	paramCount(value, 2, "fs.mv")

	src := value[0].(string)
	dst := value[1].(string)

  fmt.Printf("Move '%s' to '%s'.\n", src, dst)

	err := os.Rename(src, dst)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_touch(value ...interface{}) interface{} {
	paramCount(value, 1, "fs.touch")

	dst := value[0].(string)

  fmt.Printf("Touch file '%s'.\n", dst)

	df, err := os.Create(dst)
	if err != nil {
		panic(err.Error())
	}
	defer df.Close()

	return true
}

func fs_rm(value ...interface{}) interface{} {
	paramCount(value, 1, "fs.rm")

	dst := value[0].(string)

  fmt.Printf("Remove file '%s'.\n", dst)

	err := os.Remove(dst)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_rmdir(value ...interface{}) interface{} {
	paramCount(value, 1, "fs.rmdir")

	dst := value[0].(string)

  fmt.Printf("Remove directory '%s'.\n", dst)

	err := os.RemoveAll(dst)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_mkdir(value ...interface{}) interface{} {
	paramCount(value, 1, "fs.mkdir")

	dst := value[0].(string)

  fmt.Printf("Create directory '%s'.\n", dst)

	err := os.MkdirAll(dst, 444)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_tempdir(value ...interface{}) interface{} {
	paramCount(value, 0, "fs.tempdir")

	return os.TempDir()
}

func fs_symlink(value ...interface{}) interface{} {
	paramCount(value, 2, "fs.symlink")

	old := value[0].(string)
	neu := value[1].(string)

  fmt.Printf("Symlink '%s' from file '%s'.\n", neu, old)

	err := os.Symlink(old, neu)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_link(value ...interface{}) interface{} {
	paramCount(value, 2, "fs.link")

	old := value[0].(string)
	neu := value[1].(string)

  fmt.Printf("Hard link '%s' from file '%s'.\n", neu, old)

	err := os.Link(old, neu)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_truncate(value ...interface{}) interface{} {
	paramCount(value, 2, "fs.truncate")

	file := value[0].(string)
	size := value[1].(int64)

  fmt.Printf("Limit '%s' to %d bytes in size.\n", file, size)

	err := os.Truncate(file, size)
	if err != nil {
		panic(err.Error())
	}

	return true
}

func fs_isDir(value ...interface{}) interface{} {
	paramCount(value, 1, "fs.isDir")

	fi, err := os.Stat(value[0].(string))
	if err != nil {
		panic(err.Error())
	}

	return fi.IsDir()
}

func fs_ls(value ...interface{}) interface{} {
	files, err := ioutil.ReadDir(value[0].(string))
	if err != nil {
		panic(err.Error())
	}

	filess := make([]string, len(files))

	for i := 0; i < len(files); i++ {
		filess[i] = files[i].Name()
	}

	return filess
}

