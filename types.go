package main

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	// Nome
	Prefisso string
	Cifre    int

	// Percorsi
	Sorgente     string
	Destinazione string
	Elaborati    string

	// Esecuzione
	Comando string
	Verbose bool
}

func (config Config) counterFileName() string {
	return "." + config.Prefisso + ".counter"
}

func (config Config) counterFilePath() string {
	return filepath.Join(config.Destinazione, config.counterFileName())
}

func (config Config) mapFileName() string {
	return config.Prefisso + ".csv"
}

func (config Config) mapFilePath() string {
	return filepath.Join(config.Destinazione, config.mapFileName())
}

func (config Config) formatName(idx int, ext string) string {
	return fmt.Sprintf("%s%d%s", config.Prefisso, idx, ext)
}

func (config Config) formatPath(idx int, ext string) string {
	return filepath.Join(config.Destinazione, config.formatName(idx, ext))
}

func (config Config) processedPath(originalName string) string {
	return filepath.Join(config.Elaborati, originalName)
}

type commandArgs struct {
	SrcFileName   string
	SrcFileBase   string
	SrcFileExt    string
	SrcFilePath   string
	SrcDirPath    string
	DstDirPath    string
	DstFileName   string
	DstFilePath   string
	Counter       string
	ConfigDir     string
	Timestamp     string
	ExtensionCase string
	SrcRelFilePath string
	SrcRelDirPath string
	DstRelFilePath string
	DstRelDirPath string
}
