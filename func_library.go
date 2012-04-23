package main

import (
	"github.com/hpcorona/go-v8"
	"fmt"
	)

var library = 
	`
function extend(from, to) {
  if (from == null || typeof from != "object") return from;
  if (from.constructor != Object && from.constructor != Array) return from;
  if (from.constructor == Date || from.constructor == RegExp || from.constructor == Function ||
  from.constructor == String || from.constructor == Number || from.constructor == Boolean)
  return new from.constructor(from);

  to = to || new from.constructor();

  for (var name in from) {
    to[name] = typeof to[name] == "undefined" ? extend(from[name], null) : to[name];
  }

  return to;
}

function Project(_name, _type, _basepath) {
	this.root = ng.wd();
  this.name = _name;
  this.type = _type;
  this.basepath = path.join(ng.wd(), _basepath);
  this.outpath = path.join(this.basepath, "build");
  this.platforms = [];
  this.configurations = [];
  this.configurators = [];
  this.includedirs = [];
  this.libdirs = [];
  this.ldflags = [];
	this.ldlibs = [];
  this.defines = [];
  this.flags = [];
	this.files = [];
  this.depends = [];
	this.spec = {};
	this.realDepends = null;
	this.module = null;
	this.includes = null;
	this.source = null;
	this.binary = null;
	this.platform = null;
	this.configuration = null;
	this.generator = null;
}

Project.prototype = {
    addPlatforms: function(_platform) {
        this.platforms = this.platforms.concat(_platform);
    },

    addConfigurations: function(_configuration) {
        this.configurations = this.configurations.concat(_configuration);
    },

		addSpec: function(_spec, _val) {
				this.spec[_spec] = _val;
		},

    addIncludeDirs: function(_include) {
        this.includedirs = this.includedirs.concat(_include);
    },

    addLibraryDirs: function(_libs) {
        this.libdirs = this.libdirs.concat(_libs);
    },

    addLdFlags: function(_ldflags) {
        this.ldflags = this.ldflags.concat(_ldflags);
    },

    addDefines: function(_defines) {
        this.defines = this.defines.concat(_defines);
    },

		addLdLibs: function(_ldlibs) {
				this.ldlibs = this.ldlibs.concat(_ldlibs);
		},

    addFlags: function(_flags) {
        this.flags = this.flags.concat(_flags);
    },

		addFiles: function(_files) {
				this.files = this.files.concat(_files);
		},

    addDependencies: function(_dependency) {
				this.depends = this.depends.concat(_dependency);
    },

    clone: function() {
			return extend(this);
    },

		addConfigurator: function(_platform, _config, _func) {
			var _platforms = [].concat(_platform);
			var _configs = [].concat(_config);
			
			if (_platforms.length == 0) {
				if (_configs.length == 0) {
					var selector = {func: _func};
					this.configurators.push(selector);
				} else {
					for (var c in _configs) {
						var selector = {config: _configs[c], func: _func};
						this.configurators.push(selector);
					}
				}
			} else {
				for (var p in _platforms) {
					if (_configs.length == 0) {
						var selector = {platform: _platforms[p], func: _func};
						this.configurators.push(selector);
					} else {
						for (var c in _configs) {
							var selector = {platform: _platforms[p], config: _configs[c], func: _func};
							this.configurators.push(selector);
						}
					}
				}
			}
		},
		
		target: function(_platform, _config, _outputbase) {
			var config = this.clone();
			
			for (var i in this.configurators) {
				var c = this.configurators[i];
				if (c.platform == _platform && (c.config == undefined || c.config == _config)) {
					c.func(config, _config, _platform);
				} else if (c.config == _config && (c.platform == undefined || c.platform == _platform)) {
					c.func(config, _config, _platform);
				} else if (c.config == undefined && c.platform == undefined) {
					c.func(config, _config, _platform);
				}
			}
			
			config.platforms = [ _platform ];
			config.configurations = [ _config ];
			config.outpath = path.join(_outputbase, config.name);
			config.platform = _platform;
			config.configuration = _config;
			
			return config;
		},
		
		configure: function() {
			this.generator.configure(this);
			
			if (this.postconfig != undefined) {
				this.postconfig();
			}
		},
		
		generate: function() {
			this.generator.generate(this);
			
			if (this.postgenerate != undefined) {
				this.postgenerate();
			}
		},
		
		finalize: function() {
			this.generator.finalize(this);
			
			if (this.postfinalize != undefined) {
				this.postfinalize();
			}
		}
}

var Library = (typeof module !== "undefined" && module.exports) || {};

(function(exports){
	
	exports.projects = {};
	exports.buildStack = [];
	exports.targetProjects = {};
	exports.generator = null;
	exports.workspace = null;
	
	exports.addProjects = function(_project) {
		var _projects = [].concat(_project);
		for (var i in _projects) {
			var name = _projects[i].name;
			exports.projects[name] = _projects[i];
		}
	}
	
	exports.load = function(_paths) {
		var _allpaths = [].concat(_paths);
		
		for (var i in _allpaths) {
			var tinc = path.join(_allpaths[i], "project.js");
			ng.include(tinc);
		}
	}
	
	exports.clearBuildStack = function() {
		exports.buildStack = [];
		exports.targetProjects = {};
		exports.generator = null;
		exports.workspace = null;
	}
	
	exports.target = function(_projname, _platform, _config, _outputbase) {
		if (exports.targetProjects[_projname] != null) {
			return true;
		}
		
		if (exports.generator == null) {
			if (_platform == "android") {
				exports.generator = new AndroidGenerator();
				exports.workspace = path.join(_outputbase, _platform + "_" + _config);
			} else {
				os.error("Unsupported generator: " + _platform);
				return false;
			}
		}
		
		var _proj = exports.projects[_projname];
		if (_proj == undefined) {
			os.error("The project " + _projname + " was not found.")
			return false;
		}
		
		for (var i in _proj.depends) {
			var _dprojname = _proj.depends[i];
			
			if (exports.targetProjects[_dprojname] == null) {
				exports.target(_dprojname, _platform, _config, exports.workspace);
			}
		}
		
		var t = _proj.target(_platform, _config, exports.workspace);
		exports.buildStack.push(t.name);
		exports.targetProjects[_projname] = t;
		
		return true;
	}
	
	exports.getProp = function(_projname, _prop) {
		var _proj = exports.targetProjects[_projname];
		if (_proj == undefined) {
			os.error("Count not get the property " + _prop + " because the project " + _projname + " does not exists");
			return undefined;
		}
		
		if (_prop == "INCLUDE") {
			return _proj.includes;
		} else if (_prop == "TARGET") {
			return _proj.target;
		} else if (_prop == "SOURCE") {
			return _proj.source;
		} else if (_prop == "BINARY") {
			return _proj.binary;
		} else if (_prop == "PLATFORM") {
			return _proj.platform;
		} else if (_prop == "CONFIGURATION") {
			return _proj.configuration;
		} else if (_prop == "OUTPUT") {
			return _proj.outpath;
		} else {
			os.error("Invalid property " + _prop + ". It must be one of INCLUDE, TARGET, SOURCE, BINARY, PLATFORM, CONFIGURATION, OUTPUT")
			return undefined;
		}
	}
	
	exports.getTarget = function(_projname) {
		return exports.targetProjects[_projname];
	}
	
	exports.configure = function() {
		exports.generator.workspace(exports.workspace);
		
		for (var i in exports.buildStack) {
			var t = Library.getTarget(exports.buildStack[i]);
			t.generator = exports.generator;
			t.configure();
		}
	}
	
	exports.generate = function() {
		for (var i in exports.buildStack) {
			var t = Library.getTarget(exports.buildStack[i]);
			t.generate();
		}
	}

	exports.finalize = function() {
		for (var i in exports.buildStack) {
			var t = Library.getTarget(exports.buildStack[i]);
			t.finalize();
		}
	}
	
	exports.build = function() {
		this.configure();
		this.generate();
		this.finalize();
	}
	
})(Library);

function Generator() {
	
}

Generator.prototype = {
	generate: function(_project) {
	},
	configure: function(_project) {
	},
	finalize: function(_project) {
	},
	workspace: function(_path) {
	}
}
`
	
func load_library_functions(ctx *v8.V8Context) {
	_,err := ctx.Eval(library)
	if err != nil {
    fmt.Printf("=====\nERROR\n=====\n%s:%s", "internal library", err.Error())
    return
  }
}
