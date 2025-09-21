package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"io"
	"time"
	"gopkg.in/ini.v1"
	"os/exec"
)

var defaultConfigFile = "config.ini"

func main() {

	configFile := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Controllo argomento -config e creo il file di configurazione di default se non esiste
	if configFile == nil || strings.TrimSpace(*configFile) == "" {
		if _, err := os.Stat(defaultConfigFile); err == nil {
			fmt.Println("Specificare le configurazioni con -config")
		}else {
			fmt.Println("Creazione del file di configurazione predefinito...")
			file, err := os.Create("config.ini")
			if err != nil {
				fmt.Println("Errore creazione file di configurazione:", err)
				return
			}
			defer file.Close()
			_, err = file.WriteString(defaultConfigContent2)
			if err != nil {
				fmt.Println("Errore scrittura file di configurazione:", err)
				return
			}
			fmt.Println("File di configurazione creato: config.ini")
		
		}

		return
	}

	// Controllo che il file di configurazione esista
	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		fmt.Println("Il file di configurazione specificato non esiste.")
		return
	}

	// Leggo il file di configurazione
	cfg, err := ini.Load(*configFile)
	if err != nil {
		fmt.Println("Errore nel caricamento del file di configurazione:", err)
		return
	}

	// Verifico che le configurazioni obbligatorie siano presenti
	var missing []string
	if cfg.Section("prefisso").Key("value").String() == "" {
		missing = append(missing, "prefisso.value")
	}
	// if cfg.Section("cifre").Key("value").String() == "" {
	// 	missing = append(missing, "cifre.value")
	// }
	if cfg.Section("sorgente").Key("value").String() == "" {
		missing = append(missing, "sorgente.value")
	}
	if cfg.Section("destinazione").Key("value").String() == "" {
		missing = append(missing, "destinazione.value")
	}
	if cfg.Section("elaborati").Key("value").String() == "" {
		missing = append(missing, "elaborati.value")
	}

	if len(missing) > 0 {
		fmt.Println("Le seguenti configurazioni sono mancanti nel file di configurazione:", strings.Join(missing, ", "))
		return
	}

	// Parsing delle configurazioni
	var config Config
	config.Prefisso = cfg.Section("prefisso").Key("value").String()
	config.Cifre, err = cfg.Section("cifre").Key("value").Int()
	if err != nil {
		// fmt.Println("Errore nella conversione di 'cifre.value' in intero:", err)
		// return
		config.Cifre = 4
	}
	config.Sorgente = cfg.Section("sorgente").Key("value").String()
	config.Destinazione = cfg.Section("destinazione").Key("value").String()
	config.Elaborati = cfg.Section("elaborati").Key("value").String()
	config.Comando = cfg.Section("comando").Key("value").String()
	config.Verbose, err = cfg.Section("verbose").Key("value").Bool()
	if err != nil {
		fmt.Println("Errore nella conversione di 'verbose.value' in booleano:", err)
		return
	}

	// Carico il file di indice
	var counter int
	if _, err := os.Stat(config.counterFilePath()); os.IsNotExist(err) {
		counter = 0
		file, err := os.Create(config.counterFilePath())
		if err != nil {
			fmt.Println("Errore creazione file di contatore:", err)
			return
		}
		defer file.Close()
		_, err = file.WriteString("0")
		if err != nil {
			fmt.Println("Errore scrittura file di contatore:", err)
			return
		}
	} else {
		data, err := os.ReadFile(config.counterFilePath())
		if err != nil {
			fmt.Println("Errore lettura file di contatore:", err)
			return
		}
		counter, err = strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			fmt.Println("Errore conversione contatore:", err)
			return
		}
	}

	// Se non esiste creo il file di mappatura
	if _, err := os.Stat(config.mapFilePath()); os.IsNotExist(err) {
		file, err := os.Create(config.mapFilePath())
		if err != nil {
			fmt.Println("Errore creazione file di mappatura:", err)
			return
		}
		defer file.Close()
		_, err = file.WriteString("Indice,Originale,Rinominato\n")
		if err != nil {
			fmt.Println("Errore scrittura file di mappatura:", err)
			return
		}
	}

	fmt.Printf("Configurazione caricata correttamente: %+v\n", config)

	files, err := os.ReadDir(config.Sorgente)
	if err != nil {
		fmt.Println("Errore lettura directory sorgente:", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		originalPath := filepath.Join(config.Sorgente, file.Name())
		fileExtension := filepath.Ext(file.Name())
		newFilePath := config.formatPath(counter, fileExtension)
		processedFilePath := config.processedPath(file.Name())

		args, err := buildCommandArgs(config, counter, originalPath, newFilePath, *configFile)
		if err != nil {
			fmt.Println("Errore preparazione argomenti comando:", err)
			continue
		}

		// Copio il file nella destinazione con il nuovo nome
		src, err := os.Open(originalPath)
		if err != nil {
			fmt.Println("Errore apertura file sorgente:", err)
			continue
		}
		defer src.Close()

		dst, err := os.Create(newFilePath)
		if err != nil {
			fmt.Println("Errore creazione file destinazione:", err)
			continue
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			fmt.Println("Errore copia file:", err)
			continue
		}

		// Sposto il file originale nella cartella elaborati
		err = os.Rename(originalPath, processedFilePath)
		if err != nil {
			fmt.Println("Errore spostamento file originale:", err)
			continue
		}

		// Aggiorno il file di mappatura
		mapFile, err := os.OpenFile(config.mapFilePath(), os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Errore apertura file di mappatura:", err)
			continue
		}
		defer mapFile.Close()

		_, err = mapFile.WriteString(fmt.Sprintf("%d,%s,%s\n", counter, file.Name(), filepath.Base(newFilePath)))
		if err != nil {
			fmt.Println("Errore scrittura file di mappatura:", err)
			continue
		}

		// Incremento il contatore e aggiorno il file di contatore
		counter++
		err = os.WriteFile(config.counterFilePath(), []byte(strconv.Itoa(counter)), 0644)
		if err != nil {
			fmt.Println("Errore aggiornamento file di contatore:", err)
			continue
		}

		// Eseguo il comando opzionale se specificato TODO
		// fmt.Println("Eseguo il comando")
		
		cmd := renderCommand(config.Comando, args)

		err = executeConfiguredCommand(cmd, args)
		if err != nil {
			fmt.Println("Errore esecuzione comando:", err)
			continue
		}

	}

}

func buildCommandArgs(cfg Config, counter int, srcPath, dstPath, cfgPath string) (commandArgs, error) {
    srcAbs, err := filepath.Abs(srcPath)
    if err != nil {
        return commandArgs{}, fmt.Errorf("absolute src path: %w", err)
    }
    dstAbs, err := filepath.Abs(dstPath)
    if err != nil {
        return commandArgs{}, fmt.Errorf("absolute dst path: %w", err)
    }
    cfgAbs, err := filepath.Abs(cfgPath)
    if err != nil {
        return commandArgs{}, fmt.Errorf("absolute config path: %w", err)
    }

	srcRootAbs, err := filepath.Abs(cfg.Sorgente)
    if err != nil {
        return commandArgs{}, fmt.Errorf("absolute source root: %w", err)
    }
    dstRootAbs, err := filepath.Abs(cfg.Destinazione)
    if err != nil {
        return commandArgs{}, fmt.Errorf("absolute destination root: %w", err)
    }

    srcRelFile, err := filepath.Rel(srcRootAbs, srcAbs)
    if err != nil {
        srcRelFile = filepath.Base(srcAbs)
    }
    srcRelDir, err := filepath.Rel(srcRootAbs, filepath.Dir(srcAbs))
    if err != nil {
        srcRelDir = "."
    }
    dstRelFile, err := filepath.Rel(dstRootAbs, dstAbs)
    if err != nil {
        dstRelFile = filepath.Base(dstAbs)
    }
    dstRelDir, err := filepath.Rel(dstRootAbs, filepath.Dir(dstAbs))
    if err != nil {
        dstRelDir = "."
    }

    srcFileName := filepath.Base(srcAbs)
    dstFileName := filepath.Base(dstAbs)
    srcExt := strings.TrimPrefix(filepath.Ext(srcFileName), ".")
    dstExt := strings.TrimPrefix(filepath.Ext(dstFileName), ".")

    return commandArgs{
        SrcFileName:    srcFileName,
        SrcFileBase:    strings.TrimSuffix(srcFileName, filepath.Ext(srcFileName)),
        SrcFileExt:     srcExt,
        SrcFilePath:    srcAbs,
        SrcDirPath:     filepath.Dir(srcAbs),
        SrcRelFilePath: srcRelFile,
        SrcRelDirPath:  srcRelDir,
        DstDirPath:     filepath.Dir(dstAbs),
        DstFileName:    dstFileName,
        DstFilePath:    dstAbs,
        DstRelFilePath: dstRelFile,
        DstRelDirPath:  dstRelDir,
        Counter:       fmt.Sprintf("%d", counter),
        ConfigDir:     filepath.Dir(cfgAbs),
        Timestamp:     time.Now().Format("20060102_150405"),
        ExtensionCase: dstExt, // adatta qui se devi forzare upper/lower case
    }, nil
}

func renderCommand(template string, args commandArgs) string {
    replacements := []string{
        "$srcFileName", args.SrcFileName,
        "$srcFileBase", args.SrcFileBase,
        "$srcFileExt", args.SrcFileExt,
        "$srcFilePath", args.SrcFilePath,
        "$srcDirPath", args.SrcDirPath,
        "$dstDirPath", args.DstDirPath,
        "$dstFileName", args.DstFileName,
        "$dstFilePath", args.DstFilePath,
        "$counter", args.Counter,
        "$configDir", args.ConfigDir,
        "$timestamp", args.Timestamp,
        "$extensionCase", args.ExtensionCase,

		// Alias legacy per compatibilità
    "$fName", args.SrcFileName,
    "$relFilePath", args.SrcRelFilePath,
    "$relPath", args.SrcRelDirPath,
    "$absFilePath", args.SrcFilePath,
    "$absPath", args.SrcDirPath,

    "$fNameComp", args.DstFileName,
    "$relFilePathComp", args.DstRelFilePath,
    "$relPathComp", args.DstRelDirPath,
    "$absFilePathComp", args.DstFilePath,
    "$absPathComp", args.DstDirPath,
    }
    replacer := strings.NewReplacer(replacements...)
    return replacer.Replace(template)
}

func executeConfiguredCommand(template string, args commandArgs) error {
    rendered := strings.TrimSpace(renderCommand(template, args))
    if rendered == "" {
        return nil
    }

    cmd := exec.Command("bash", "-lc", rendered)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    return cmd.Run()
}