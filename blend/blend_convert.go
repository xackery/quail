package blend

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func Convert(path string, shortname string, blenderPath string) error {

	wPy, err := os.Create(fmt.Sprintf("%s/cache/convert.py", path))
	if err != nil {
		return err
	}

	_, err = wPy.WriteString(convertPy)
	if err != nil {
		return fmt.Errorf("writeString: %w", err)
	}

	wPy.Close()

	if blenderPath == "" {
		blenderPath = blendPath()
	}
	if blenderPath == "" {
		return fmt.Errorf("blender not found, set it using --blender=<path>")
	}

	fullPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("filepath.Abs: %w", err)
	}
	//os.Chdir(fullPath)
	cmd := exec.Cmd{}
	//cmd.Env = append(cmd.Env, os.Environ()...)
	//cmd.Path = "/bin/bash"
	cmd.Path = blenderPath
	//cmd.Stdin = strings.NewReader(convertPy)
	//cmd.Args = []string{fmt.Sprintf("%s/%s.blend", fullPath, shortname)}
	//cmd.Path = "/usr/bin/open"

	//arg := fmt.Sprintf(`-c "%s --log-level -1 --background %s/%s.blend --python %s/cache/convert.py"`, blenderPath, path, shortname, path)
	//arg := "--python-console"
	//arg := fmt.Sprintf(`-c "%s --background %s/%s.blend --python-console"`, blenderPath, fullPath, shortname)
	cmd.Args = []string{fmt.Sprintf("--background %s/%s.blend", fullPath, shortname),
		//"--log-level -1",
		//"--python-console",
		"--python " + fmt.Sprintf("%s/cache/convert.py", fullPath),
		//		"--env-system-python /usr/local/bin/python",
		//"--python-use-system-env"
	}

	//fmt.Println("executing", arg)

	//--python-console with stdin?
	cmd.Dir = fullPath
	//fmt.Printf("%+v\n", cmd)

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf
	err = cmd.Start()
	if err != nil {
		fmt.Println("output:", outBuf.String())
		return fmt.Errorf("run: %w", err)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println("output:", outBuf.String())
		return fmt.Errorf("run: %w", err)
	}
	fmt.Println("output:", outBuf.String())
	os.Exit(0)
	return nil
}
