package main

import (
	"github.com/hpcorona/go-v8"
	)

var android = 
	`
function AndroidGenerator() {
	
}

AndroidGenerator.prototype = new Generator();
AndroidGenerator.prototype.generate = function(_project) {
	var android = {
		apppath: (_project.type == "application" ? _project.outpath : null),
		path: _project.outpath,
		source: _project.source,
		includes: this.genIncludes(_project),
		statics: this.genStaticLibs(_project),
		whole_statics: this.genWholeStaticLibs(_project),
		shared: this.genSharedLibs(_project),
		files: [].concat(_project.files),
		dependencies: this.genDepOutputs(_project),
		direct_dependencies: this.genDirectDepsOutputs(_project),
		ldlibs: [].concat(_project.ldlibs),
		ldflags: [].concat(_project.ldflags),
		defines: [].concat(_project.defines),
		module: _project.module,
		modules: this.genModules(),
		libtype: (_project.type == "application" || _project.type == "shared" ? "BUILD_SHARED_LIBRARY" : "BUILD_STATIC_LIBRARY"),
		platform: _project.platform,
		spec: _project.spec,
		configuration: _project.configuration
	};
	
	fs.rmdir(android.path);
	fs.mkdir(android.path);
	fs.symlink(_project.basepath, android.source);
	
	fs.mkdir(path.join(android.path, "jni"));
	
	var android_mk = path.join(android.path, "jni", "Android.mk");
	var android_mk_tpl = "\
{{#apppath}}\n\
APPLICATION_PATH := {{apppath}}\n\n\
{{/apppath}}\n\
{{#direct_dependencies}}\n\
include {{.}}/jni/Android.mk\n\
{{/direct_dependencies}}\n\
\n\
LOCAL_PATH := {{source}}\n\
\n\
include $(CLEAR_VARS)\n\
\n\
{{#defines}}\n\
LOCAL_CFLAGS += -D{{.}}\n\
{{/defines}}\n\
\n\
{{#includes}}\n\
LOCAL_C_INCLUDES +=	{{.}}\n\
{{/includes}}\n\
\n\
LOCAL_MODULE    := {{module}}\n\
LOCAL_SRC_FILES := \\\n\
{{#files}}\n\
	{{.}} \\\n\
{{/files}}\n\
\n\
{{#ldlibs}}\n\
LOCAL_LDLIBS += -l{{.}}\n\
{{/ldlibs}}\n\
{{#ldflags}}\n\
LOCAL_LDLIBS += {{.}}\n\
{{/ldflags}}\n\
{{#statics}}\n\
LOCAL_STATIC_LIBRARIES += {{.}}\n\
{{/statics}}\n\
{{#whole_statics}}\n\
LOCAL_WHOLE_STATIC_LIBRARIES += {{.}}\n\
{{/whole_statics}}\n\
{{#shared}}\n\
LOCAL_SHARED_LIBRARIES += {{.}}\n\
{{/shared}}\n\
\n\
include $({{libtype}})\n\
";

	fs.write(path.join(android.path, "jni", "Android.mk"), Mustache.render(android_mk_tpl, android));
	
	if (_project.type != "application") return;
	
	var application_mk_tpl = "\
{{#modules}}\n\
APP_MODULES += {{.}}\n\
{{/modules}}\n\
APP_STL := gnustl_static\n\
APP_OPTIM := {{configuration}}\n\
APP_EXPORT_CFLAGS += -g\n\
";
	
	fs.write(path.join(android.path, "jni", "Application.mk"), Mustache.render(application_mk_tpl, android));
	
	var project_properties_tpl = "\
target={{spec.android_version}}\
";

	fs.write(path.join(android.path, "project.properties"), Mustache.render(project_properties_tpl, android));
}

AndroidGenerator.prototype.genModules = function() {
	var modules = [];
	
	for (var i in Library.buildStack) {
		var t = Library.getTarget(Library.buildStack[i]);
		if (t.type == "application") {
			modules.push("game");
		} else {
			modules.push(t.name);
		}
	}
	
	return modules;
}

AndroidGenerator.prototype.genDirectDepsOutputs = function(_project) {
	var outputs = [];
	
	for (var i in _project.depends) {
		outputs.push(Library.getProp(_project.depends[i], "OUTPUT"));
	}
	
	return outputs;
}

AndroidGenerator.prototype.genDepOutputs = function(_project) {
	var outputs = [];
	
	for (var i in _project.realDepends) {
		outputs.push(Library.getProp(_project.realDepends[i], "OUTPUT"));
	}
	
	return outputs;
}

AndroidGenerator.prototype.genIncludes = function(_project) {
	var includes = [].concat(_project.includes);
	
	for (var i in _project.depends) {
		var projincs = Library.getProp(_project.depends[i], "INCLUDE");
		
		includes = includes.concat(projincs);
	}
	
	return includes;
}

AndroidGenerator.prototype.genStaticLibs = function(_project) {
	var libs = [];
	
	for (var i in _project.realDepends) {
		var target = Library.getTarget(_project.realDepends[i]);
		if (target == undefined) {
			os.error("Target project not found " + project.realDepends[i]);
		}
		
		if (target.type != "whole_static" && target.type != "static") continue;
		
		libs.push(_project.realDepends[i]);
	}
	
	return libs;
}

AndroidGenerator.prototype.genWholeStaticLibs = function(_project) {
	var libs = [];
	
	for (var i in _project.realDepends) {
		var target = Library.getTarget(_project.realDepends[i]);
		if (target == undefined) {
			os.error("Target project not found " + project.realDepends[i]);
		}
		
		if (target.type != "whole_static") continue;
		
		libs.push(_project.realDepends[i]);
	}
	
	return libs;
}

AndroidGenerator.prototype.genSharedLibs = function(_project) {
	var libs = [];
	
	for (var i in _project.realDepends) {
		var target = Library.getTarget(_project.realDepends[i]);
		if (target == undefined) {
			os.error("Target project not found " + project.realDepends[i]);
		}
		
		if (target.type != "shared") continue;
		
		libs.push(_project.realDepends[i]);
	}
	
	return libs;
}

AndroidGenerator.prototype.configure = function(_project) {
	if (_project.type == "application") {
		_project.module = "game";
	} else if (_project.type == "whole_static" || _project.type == "static" || _project.type == "shared") {
		_project.module = _project.name;
	} else {
		os.error("Invalid project type for android " + _project.type + ". Please specify one of application, static, whole_static, shared");
	}
	
	this.module = null;
	this.includes = null;
	this.source = null;
	this.binary = null;
	this.generator = null;
	this.realDepends = null;
	
	_project.binary = path.join(_project.outpath, "bin");
	_project.source = path.join(_project.outpath, "src");
	_project.java = path.join(_project.outpath, "java");
	_project.generator = this;
	
	_project.includes = [];
	for (var i in _project.includedirs) {
		_project.includes.push(path.join(_project.source, _project.includedirs[i]));
	}
	
	if (_project.type == "application") {
		_project.realDepends = this.genRealDepends(_project);
	} else {
		_project.realDepends = [].concat(_project.depends);
	}
}

AndroidGenerator.prototype.genRealDepends = function(_project) {
	var rdeps = [].concat(_project.depends);
	
	for (var i in Library.buildStack) {
		var t = Library.getTarget(Library.buildStack[i]);
		
		if (t.type == "whole_static" || t.type == "shared") {
			rdeps.push(t.name);
		}
	}
	
	var sorted_arr = rdeps.sort();
	var results = [];
	for (var i = 0; i < sorted_arr.length; i++) {
		if (results.length == 0 || sorted_arr[i] != results[results.length - 1]) {
			results.push(sorted_arr[i]);
		}
	}
	
	return results;
}

AndroidGenerator.prototype.finalize = function(_project) {
	
}
`
	
func load_android_functions(ctx *v8.V8Context) {
	ctx.Eval(android)
}
