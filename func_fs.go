package main

import (
	"fmt"
	"github.com/hpcorona/go-v8"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"bufio"
)

func load_fs_functions(ctx *v8.V8Context) {
	// file system functions
	ctx.AddFunc("_fs_cd", fs_cd)
	ctx.AddFunc("_fs_cp", fs_cp)
  ctx.AddFunc("_fs_cpt", fs_cpt)
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
  ctx.AddFunc("_fs_lsn", fs_lsn)
  ctx.AddFunc("_fs_wd", fs_wd)
	ctx.AddFunc("_fs_read", fs_read)
	ctx.AddFunc("_fs_write", fs_write)

	ctx.Eval(`
  this.fs =
    {
      "cd": function(args) { _fs_cd(args); },
      "cp": function(v0, v1) { _fs_cp(v0, v1); },
      "cpt": function(v0, v1) { _fs_cpt(v0, v1); },
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
      "ls": function() { return _fs_ls.apply(this, arguments); },
      "lsn": function() { return _fs_lsn.apply(this, arguments); },
      "wd": function() { return _fs_wd(); },
			"read": function(v0) { return _fs_read(v0); },
			"write": function() { return _fs_write.apply(this, arguments); }
    };
  `)
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

	_, err = io.Copy(df, sf)
	if err != nil {
		panic(err.Error())
	}
  df.Close()

  finfo, err := os.Stat(src)
  if err != nil {
    panic(err.Error())
  }
  dof, err := os.Open(dst)
  if err != nil {
    panic(err.Error())
  }
  os.Chmod(dst, finfo.Mode())
  dof.Close()

	return true
}

func fs_cpt(value ...interface{}) interface{} {
  paramCount(value, 2, "fs.cpt")

  src := value[0].(string)
  dst := value[1].(string)

  _, f := filepath.Split(src)
  dstf := filepath.Join(dst, f)

  return fs_cp(src, dstf)
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

	err := os.MkdirAll(dst, 2147484141)
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
  paramMin(value, 1, "fs.ls")
  paramMax(value, 3, "fs.ls")

  path := value[0].(string)
  matchp := ""
  if len(value) > 1 {
    matchp = value[1].(string)
  }
  cs := true
  if len(value) > 2 {
    cs = value[2].(bool)
  }

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}

	filess := make([]string, 0, len(files))

	for i := 0; i < len(files); i++ {
    name := files[i].Name()
    if matchp == "" {
      filess = append(filess, filepath.Join(path, name))
    } else {
      if match(matchp, name, cs) {
        filess = append(filess, filepath.Join(path, name))
      }
    }
	}

	return filess
}

func listDir(mp MatchPattern, cs bool, path string, prefix string) []string {
  files, err := ioutil.ReadDir(path)
  if err != nil {
    panic(err.Error())
  }

  filess := make([]string, 0, len(files))

  for i := 0; i < len(files); i++ {
    name := filepath.Join(prefix, files[i].Name())
    if mp.Match(name, cs) {
      filess = append(filess, name)
    }

    if files[i].IsDir() {
      newpath := filepath.Join(path, files[i].Name())
      filess = append(filess, listDir(mp, cs, newpath, name)...)
    }
  }

  return filess
}

func fs_lsn(value ...interface{}) interface{} {
  paramMin(value, 1, "fs.lsn")
  paramMax(value, 3, "fs.lsn")

  path := value[0].(string)
  var matchp MatchPattern
  if len(value) > 1 {
    matchp = NewMatchPattern(value[1].(string))
  }
  cs := true
  if len(value) > 2 {
    cs = value[2].(bool)
  }

  if path == "" {
    return listDir(matchp, cs, ".", "")
  }
  return listDir(matchp, cs, path, "")
}

func fs_wd(value ...interface{}) interface{} {
  paramCount(value, 0, "fs.wd")

  path, _ := os.Getwd()

  return path
}

func fs_read(value ...interface{}) interface{} {
	paramCount(value, 1, "fs.read")
	
	data, err := ioutil.ReadFile(value[0].(string))
	if err != nil {
		panic(err.Error())
	}
	
	return string(data)
}

func fs_write(value ...interface{}) interface{} {
	paramMin(value, 2, "fs.write")
	
	file, err := os.Create(value[0].(string))
	if err != nil {
		panic(err.Error())
	}
	
	defer file.Close()
	
	last := false
	writer := bufio.NewWriter(file)
	for i := 1; i < len(value); i++ {
		last = i == (len(value) - 1)
    intd, ok := value[i].([]interface{})
    if ok {
      for j := 0; j < len(intd); j++ {
				last_this := last && j == (len(intd) - 1)
				if !last_this {
					writer.WriteString(fmt.Sprintf("%s\n", intd[j].(string)))
				} else {
					writer.WriteString(intd[j].(string))
				}
      }
    } else {
			if !last {
				writer.WriteString(fmt.Sprintf("%s\n", value[i].(string)))
			} else {
				writer.WriteString(value[i].(string))
			}
    }
	}
	
	writer.Flush()
	
	return true
}
