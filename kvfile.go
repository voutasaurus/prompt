package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type kvfile struct {
	Path string
	m    map[string]string
}

func (kv *kvfile) Read(key string) (value string, ok bool) {
	kv.readAll(kv.Path)
	value, ok = kv.m[key]
	return value, ok
}

func (kv *kvfile) readAll(path string) error {
	kv.m = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// done reading
			return nil
		}
		return err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for i := 1; scan.Scan(); i++ {
		if err := kv.readline(scan.Text()); err != nil {
			return fmt.Errorf("read error in %q on line %d: %v", path, i, err)
		}
	}
	return scan.Err()
}

func (kv *kvfile) readline(line string) error {
	// skip empty lines
	if strings.Trim(line, " \t\n") == "" {
		return nil
	}

	ss := strings.Split(line, "=")
	key, value := ss[0], strings.Join(ss[1:], "=")

	if _, duplicate := kv.m[key]; duplicate {
		return fmt.Errorf("%q is duplicated", key)
	}

	kv.m[key] = value

	return nil
}

func (kv *kvfile) WriteNew(key, value string) error {
	if err := kv.readAll(kv.Path); err != nil {
		return fmt.Errorf("cannot read from %q as it is malformed: %v", kv.Path, err)
	}

	if _, duplicate := kv.m[key]; duplicate {
		return fmt.Errorf("cannot write to %q as it already contains a value for %q", kv.Path, key)
	}

	kv.m[key] = value
	if err := kv.writeAll(kv.Path); err != nil {
		return fmt.Errorf("cannot write to %q as it is malformed: %v", kv.Path, err)
	}
	return nil
}

func (kv *kvfile) Write(key, value string) error {
	if err := kv.readAll(kv.Path); err != nil {
		return fmt.Errorf("cannot read from %q as it is malformed: %v", kv.Path, err)
	}

	kv.m[key] = value
	if err := kv.writeAll(kv.Path); err != nil {
		return fmt.Errorf("cannot write to %q as it is malformed: %v", kv.Path, err)
	}
	return nil
}

func (kv *kvfile) writeAll(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for k, v := range kv.m {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}

	return nil
}
