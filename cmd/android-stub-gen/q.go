package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// StubTemplate 是生成的C++ stub代码的模板
const StubTemplate = `
#include <android/log.h>
#include <dlfcn.h>
#include <stdlib.h>

typedef void goMainFunc_t();  

int main(int argc, char** argv) {
    __android_log_print(ANDROID_LOG_VERBOSE, "miqt_stub", "Starting up");
        
    void* handle = dlopen("{{.SourceSOFile}}", RTLD_LAZY);
    if (handle == NULL) {
        __android_log_print(ANDROID_LOG_VERBOSE, "miqt_stub", "miqt_stub: null handle opening so: %s", dlerror());
        exit(1);
    }
    
    void* goMain = dlsym(handle, "{{.FunctionName}}");
    if (goMain == NULL) {
        __android_log_print(ANDROID_LOG_VERBOSE, "miqt_stub", "miqt_stub: null handle looking for function: %s", dlerror());
        exit(1);        
    }
    
    __android_log_print(ANDROID_LOG_VERBOSE, "miqt_stub", "miqt_stub: Found target, calling");
    
    // Cast to function pointer and call
    goMainFunc_t* f = (goMainFunc_t*)goMain;
    f();
    
    __android_log_print(ANDROID_LOG_VERBOSE, "miqt_stub", "miqt_stub: Target function returned");
    return 0;
}
`

// Config 用于填充模板的结构体
type Config struct {
	SourceSOFile string
	FunctionName string
	DestSOFile   string
	QTPath       string
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s src.so function-name dest.so\n", os.Args[0])
		os.Exit(1)
	}

	argSourceSOFile := os.Args[1]
	argFunctionName := os.Args[2]
	argDestSOFile := os.Args[3]
	qtPath := "D:\\Qt\\6.9.0\\android_arm64_v8a\\" // 假设Qt路径

	config := Config{
		SourceSOFile: filepath.Base(argSourceSOFile),
		FunctionName: argFunctionName,
		DestSOFile:   argDestSOFile,
		QTPath:       qtPath,
	}

	// 创建临时目录
	tmpdir, err := os.MkdirTemp("", "android-gen-stub")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	fmt.Printf("- Using temporary directory: %s\n", tmpdir)
	fmt.Printf("- Found Qt path: %s\n", qtPath)

	// 生成C++ stub代码
	stubFilePath := filepath.Join(tmpdir, "miqtstub.cpp")
	stubFile, err := os.Create(stubFilePath)
	if err != nil {
		log.Fatalf("Failed to create stub file: %v", err)
	}
	defer stubFile.Close()

	tmpl, err := template.New("stub").Parse(StubTemplate)
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	err = tmpl.Execute(stubFile, config)
	if err != nil {
		log.Fatalf("Failed to execute template: %v", err)
	}

	// 编译C++代码
	cxx := "g++"
	args := []string{
		"-shared",
		"-ldl",
		"-llog",
		filepath.Join(qtPath, "plugins/platforms/libplugins_platforms_qtforandroid_arm64-v8a.so"),
		filepath.Join(qtPath, "lib/libQt5Widgets_arm64-v8a.so"),
		filepath.Join(qtPath, "lib/libQt5Gui_arm64-v8a.so"),
		filepath.Join(qtPath, "lib/libQt5Core_arm64-v8a.so"),
		filepath.Join(qtPath, "lib/libQt5Svg_arm64-v8a.so"),
		filepath.Join(qtPath, "lib/libQt5AndroidExtras_arm64-v8a.so"),
		"-fPIC",
		"-DQT_WIDGETS_LIB",
		"-I" + filepath.Join(qtPath, "include/QtWidgets"),
		"-I" + filepath.Join(qtPath, "include/"),
		"-I" + filepath.Join(qtPath, "include/QtCore"),
		"-DQT_GUI_LIB",
		"-I" + filepath.Join(qtPath, "include/QtGui"),
		"-DQT_CORE_LIB",
		stubFilePath,
		"-Wl,-soname," + filepath.Base(argDestSOFile),
		"-o", argDestSOFile,
	}

	cmd := exec.Command(cxx, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to compile stub: %v", err)
	}

	fmt.Println("Done.")
}
