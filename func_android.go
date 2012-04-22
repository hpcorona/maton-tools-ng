package main

import (
	"github.com/hpcorona/go-v8"
	)

var android = 
	`
function AndroidGenerator() {
	
}

AndroidGenerator.prototype = new Generator();

AndroidGenerator.prototype.workspace = function(_path) {
	fs.mkdir(path.join(_path, ".metadata"));
	var projsdir = path.join(_path, ".metadata", ".plugins", "org.eclipse.core.resources", ".projects");
	fs.mkdir(projsdir);
	
	fs.write(path.join(_path, ".metadata", "version.ini"), "org.eclipse.core.runtime=1");
	
	for (var i in Library.buildStack) {
		fs.mkdir(path.join(projsdir, Library.buildStack[i], "org.eclipse.jdt.core"));
	}
}

AndroidGenerator.prototype.genProject = function(_project) {
	var dot_classpath_tpl = '\
<?xml version="1.0" encoding="UTF-8"?>\n\
<classpath>\n\
	<classpathentry kind="src" path="java"/>\n\
	<classpathentry kind="src" path="gen"/>\n\
{{#spec.linked_sources}}\n\
	<classpathentry kind="src" path="{{name}}"/>\n\
{{/spec.linked_sources}}\n\
{{#spec.java}}\n\
	<classpathentry kind="src" path="code"/>\n\
{{/spec.java}}\n\
	<classpathentry kind="con" path="com.android.ide.eclipse.adt.ANDROID_FRAMEWORK"/>\n\
	<classpathentry kind="con" path="com.android.ide.eclipse.adt.LIBRARIES"/>\n\
{{#spec.jars}}\n\
	<classpathentry kind="lib" path="{{.}}"/>\n\
{{/spec.jars}}\n\
	<classpathentry kind="output" path="bin/classes"/>\n\
</classpath>\n\
';
	var dot_project_tpl = '\
<?xml version="1.0" encoding="UTF-8"?>\n\
<projectDescription>\n\
	<name>{{name}}</name>\n\
	<comment></comment>\n\
	<projects>\n\
	</projects>\n\
	<buildSpec>\n\
		<buildCommand>\n\
			<name>org.eclipse.cdt.managedbuilder.core.genmakebuilder</name>\n\
			<triggers>clean,full,incremental,</triggers>\n\
			<arguments>\n\
				<dictionary>\n\
					<key>?children?</key>\n\
					<value>?name?=outputEntries\\|?children?=?name?=entry\\\\\\\\\\\\\\|\\\\\\|\\||</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>?name?</key>\n\
					<value></value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.append_environment</key>\n\
					<value>true</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.buildArguments</key>\n\
					<value></value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.buildCommand</key>\n\
					<value>ndk-build</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.cleanBuildTarget</key>\n\
					<value>clean</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.contents</key>\n\
					<value>org.eclipse.cdt.make.core.activeConfigSettings</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.enableAutoBuild</key>\n\
					<value>false</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.enableCleanBuild</key>\n\
					<value>true</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.enableFullBuild</key>\n\
					<value>true</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.fullBuildTarget</key>\n\
					<value>V=1</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.stopOnError</key>\n\
					<value>true</value>\n\
				</dictionary>\n\
				<dictionary>\n\
					<key>org.eclipse.cdt.make.core.useDefaultBuildCmd</key>\n\
					<value>true</value>\n\
				</dictionary>\n\
			</arguments>\n\
		</buildCommand>\n\
		<buildCommand>\n\
			<name>com.android.ide.eclipse.adt.ResourceManagerBuilder</name>\n\
			<arguments>\n\
			</arguments>\n\
		</buildCommand>\n\
		<buildCommand>\n\
			<name>com.android.ide.eclipse.adt.PreCompilerBuilder</name>\n\
			<arguments>\n\
			</arguments>\n\
		</buildCommand>\n\
		<buildCommand>\n\
			<name>org.eclipse.jdt.core.javabuilder</name>\n\
			<arguments>\n\
			</arguments>\n\
		</buildCommand>\n\
		<buildCommand>\n\
			<name>com.android.ide.eclipse.adt.ApkBuilder</name>\n\
			<arguments>\n\
			</arguments>\n\
		</buildCommand>\n\
		<buildCommand>\n\
			<name>org.eclipse.cdt.managedbuilder.core.ScannerConfigBuilder</name>\n\
			<triggers>full,incremental,</triggers>\n\
			<arguments>\n\
			</arguments>\n\
		</buildCommand>\n\
	</buildSpec>\n\
	<natures>\n\
		<nature>com.android.ide.eclipse.adt.AndroidNature</nature>\n\
		<nature>org.eclipse.jdt.core.javanature</nature>\n\
		<nature>org.eclipse.cdt.core.cnature</nature>\n\
		<nature>org.eclipse.cdt.core.ccnature</nature>\n\
		<nature>org.eclipse.cdt.managedbuilder.core.managedBuildNature</nature>\n\
		<nature>org.eclipse.cdt.managedbuilder.core.ScannerConfigNature</nature>\n\
	</natures>\n\
	<linkedResources>\n\
{{#spec.linked_sources}}\n\
		<link>\n\
			<name>{{name}}</name>\n\
			<type>2</type>\n\
			<location>{{location}}</location>\n\
		</link>\n\
{{/spec.linked_sources}}\n\
{{#spec.java}}\n\
		<link>\n\
			<name>code</name>\n\
			<type>2</type>\n\
			<location>{{spec.java}}</location>\n\
		</link>\n\
{{/spec.java}}\n\
	</linkedResources>\n\
</projectDescription>\n\
';
	
	fs.mkdir(path.join(_project.outpath, "java"));
	fs.mkdir(path.join(_project.outpath, "gen"));
	fs.mkdir(path.join(_project.outpath, "bin/classes"));
	
	fs.write(path.join(_project.outpath, ".classpath"), Mustache.render(dot_classpath_tpl, _project));
	fs.write(path.join(_project.outpath, ".project"), Mustache.render(dot_project_tpl, _project));
}

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
	this.genProject(_project);
	
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
