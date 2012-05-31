package main

import (
	"github.com/hpcorona/go-v8"
	"fmt"
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

AndroidGenerator.prototype.genManifest = function(_project) {
	if (_project.spec.android == null) return;
	
	var manifest_tpl = '\
<?xml version="1.0" encoding="utf-8"?>\n\
<manifest xmlns:android="http://schemas.android.com/apk/res/android"\n\
    package="{{spec.android.package}}"\n\
		android:installLocation="preferExternal"\n\
    android:versionCode="{{spec.app.versionCode}}"\n\
    android:versionName="{{spec.app.versionName}}">\n\
{{#spec.android.permissions}}\n\
    <uses-permission android:name="{{.}}" />\n\
{{/spec.android.permissions}}\n\
    <uses-sdk android:minSdkVersion="{{spec.android.version.min}}" android:targetSdkVersion="{{spec.android.version.target}}"/>\n\
{{#spec.android.gles2}}\n\
    <uses-feature android:glEsVersion="0x00020000" android:required="true" />\n\
{{/spec.android.gles2}}\n\
		<supports-screens android:resizeable="false"\n\
			android:smallScreens="true"\n\
			android:normalScreens="true"\n\
			android:largeScreens="true"\n\
			android:xlargeScreens="true"\n\
			android:anyDensity="true" />\n\
\n\
    <application\n\
        android:icon="@drawable/ic_launcher"\n\
        android:label="@string/app_label" android:allowClearUserData="true">\n\
        <activity\n\
            android:name=".{{spec.android.activity}}"\n\
            android:screenOrientation="portrait"\n\
            android:label="@string/app_name" >\n\
            <intent-filter>\n\
                <action android:name="android.intent.action.MAIN" />\n\
                <category android:name="android.intent.category.LAUNCHER" />\n\
            </intent-filter>\n\
        </activity>\n\
    </application>\n\
</manifest>\n\
';
	fs.write(path.join(_project.outpath, "AndroidManifest.xml"), Mustache.render(manifest_tpl, _project));
	
	var res_ldpi = path.join(_project.outpath, "res", "drawable-ldpi");
	var res_mdpi = path.join(_project.outpath, "res", "drawable-mdpi");
	var res_hdpi = path.join(_project.outpath, "res", "drawable-hdpi");
	var layout = path.join(_project.outpath, "res", "layout");
	var values = path.join(_project.outpath, "res", "values");
	var assets = path.join(_project.outpath, "assets");
	
	fs.mkdir(res_ldpi);
	fs.mkdir(res_mdpi);
	fs.mkdir(res_hdpi);
	fs.mkdir(layout);
	fs.mkdir(values);
	fs.mkdir(assets);
	
	if (_project.spec.android.icon != undefined) {
		if (_project.spec.android.icon.ldpi != null) {
			fs.cp(_project.spec.android.icon.ldpi, path.join(res_ldpi, "ic_launcher.png"));
		}
		if (_project.spec.android.icon.mdpi != null) {
			fs.cp(_project.spec.android.icon.mdpi, path.join(res_mdpi, "ic_launcher.png"));
		}
		if (_project.spec.android.icon.hdpi != null) {
			fs.cp(_project.spec.android.icon.hdpi, path.join(res_hdpi, "ic_launcher.png"));
		}
	}
	
	var strings_tpl = '\
<?xml version="1.0" encoding="utf-8"?>\n\
<resources>\n\
  <string name="app_name">{{spec.app.name}}</string>\n\
  <string name="app_label">{{spec.app.label}}</string>\n\
</resources>\
';

	fs.write(path.join(values, "strings.xml"), Mustache.render(strings_tpl, _project));
	fs.write(path.join(_project.outpath, "lint.xml"), ['<?xml version="1.0" encoding="UTF-8"?>',"<lint>","</lint>"]);
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
	</linkedResources>\n\
</projectDescription>\n\
';

	var source_dirs = [];
	
	for (var i in _project.spec.linked_sources) {
		fs.symlink(_project.spec.linked_sources[i].location, path.join(_project.outpath, _project.spec.linked_sources[i].name));
		source_dirs.push(_project.spec.linked_sources[i].name);
	}
	
	var libsdir = path.join(_project.outpath, "libs");
	fs.mkdir(libsdir);
	for (var i in _project.spec.jars) {
		fs.cpt(_project.spec.jars[i], libsdir);
	}
	_project.spec.jars = fs.ls(libsdir, "*");	
	
	if (_project.spec.java != undefined) {
		fs.symlink(_project.spec.java, path.join(_project.outpath, "java"));
	} else {
		fs.mkdir(path.join(_project.outpath, "java"));
	}
	source_dirs.push("java");
	
	fs.mkdir(path.join(_project.outpath, "gen"));
	fs.mkdir(path.join(_project.outpath, "bin/classes"));
	
	fs.write(path.join(_project.outpath, ".classpath"), Mustache.render(dot_classpath_tpl, _project));
	fs.write(path.join(_project.outpath, ".project"), Mustache.render(dot_project_tpl, _project));
	
	if (_project.spec.android != undefined) {
		var ant_tpl = "\
source.dir={{sources}}\n\
{{#keystore}}\n\
key.store={{keystore}}\n\
{{/keystore}}\n\
{{#keyalias}}\n\
key.alias={{keyalias}}\n\
{{/keyalias}}\n\
";

		var ant_props = {
				sources: source_dirs.join(":"),
				keystore: _project.spec.android.keystore,
				keyalias: _project.spec.android.keyalias
			};
	
		fs.write(path.join(_project.outpath, "ant.properties"), Mustache.render(ant_tpl, ant_props));
	}
	
	this.genManifest(_project);
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
		flags: [].concat(_project.flags),
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
{{#flags}}\n\
LOCAL_CFLAGS += {{.}}\n\
{{/flags}}\n\
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
APP_ABI := armeabi armeabi-v7a x86\n\
{{#modules}}\n\
APP_MODULES += {{.}}\n\
{{/modules}}\n\
APP_STL := gnustl_static\n\
APP_OPTIM := {{configuration}}\n\
APP_EXPORT_CFLAGS += -g\n\
APP_GNUSTL_FORCE_CPP_FEATURES := exceptions rtti\n\
";
	
	fs.write(path.join(android.path, "jni", "Application.mk"), Mustache.render(application_mk_tpl, android));
	
	var project_properties_tpl = "\
target=android-{{spec.android.version.target}}\
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
	if (_project.type != "application") return;
	
	var ANDROID_CMD = os.findCmd("android");
	os.run(ANDROID_CMD, "update", "project", "-p", _project.outpath);
	
	os.log("Fixing APK name...");
	var build = path.join(_project.outpath, "build.xml");
	var content = fs.read(build);
	content = content.replace(_project.spec.android.activity, _project.spec.android.apk);
	fs.write(build, content);
}
`
	
func load_android_functions(ctx *v8.V8Context) {
	_,err := ctx.Eval(android)
	if err != nil {
    fmt.Printf("=====\nERROR\n=====\n%s:%s", "internal android", err.Error())
    return
  }
}
