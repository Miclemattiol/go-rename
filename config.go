package main

const defaultConfigContent = `
[nome]
prefisso=       # Prefisso utilizzato per formare il nome dei file
cifre=4         # Numero di cifre per il contatore progressivo (default: 4)

[percorsi]
sorgente=       # Percorso della cartella sorgente dei file da elaborare
destinazione=   # Percorso della cartella dove verranno salvati i file rinominati
elaborati=      # Percorso della cartella dove spostare i file originali dopo l'elaborazione

[esecuzione]
verbose=false   # true per abilitare i messaggi di debug, false per disabilitarli
comando=        # Comando facoltativo eseguito per ogni file.
                # Sostituzioni disponibili:
                # $srcFileName -> nome originale con estensione
                # $srcFileBase -> nome originale senza estensione
                # $srcFileExt -> estensione originale senza punto
                # $srcFilePath -> percorso assoluto del file sorgente
                # $srcDirPath -> directory del file sorgente
                # $dstDirPath -> directory di destinazione
                # $dstFileName -> nuovo nome con estensione
                # $dstFilePath -> percorso assoluto del file rinominato
                # $counter -> indice formattato secondo le cifre configurate
                # $configDir -> directory del file di configurazione
                # $timestamp -> data/ora corrente (YYYYMMDD_HHmmss)
                # $extensionCase -> estensione finale come verrà scritta`

const defaultConfigContent2 = `
[prefisso]
value=
# Prefisso utilizzato per formare il nome dei file

[sorgente]
value=
# Percorso assoluto alla cartella sorgente dei file da elaborare

[destinazione]
value=
# Percorso assoluto alla cartella dove verranno salvati i file rinominati

[elaborati]
value=
# Percorso assoluto alla cartella dove spostare i file originali dopo l'elaborazione

[comando]
value=
# Comando facoltativo eseguito per ogni file.
# Puoi usare: $fName, $relFilePath, $relPath, $absFilePath, $absPath,
# oppure le stesse variabili con 'Comp' per riferirti al file rinominato.
# Esempio: echo "Elaborato $fName in $absFilePathComp"
# Altre sostituzioni disponibili:
# $srcFileName -> nome originale con estensione
# $srcFileBase -> nome originale senza estensione
# $srcFileExt -> estensione originale senza punto
# $srcFilePath -> percorso assoluto del file sorgente
# $srcDirPath -> directory del file sorgente
# $dstDirPath -> directory di destinazione
# $dstFileName -> nuovo nome con estensione
# $dstFilePath -> percorso assoluto del file rinominato
# $counter -> indice formattato secondo le cifre configurate
# $configDir -> directory del file di configurazione
# $timestamp -> data/ora corrente (YYYYMMDD_HHmmss)
# $extensionCase -> estensione finale come verrà scritta

[verbose]
value=false
# true per abilitare i messaggi di debug, false per disabilitarli`