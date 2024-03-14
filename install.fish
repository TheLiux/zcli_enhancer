#!/usr/bin/env fish

# Il percorso dell'eseguibile buildato da Go
set EXECUTABLE "./zcli_enhancer"

# Il percorso di destinazione in /usr/local/bin
set DESTINATION "/usr/local/bin/zcli_enhancer"

# Copia l'eseguibile in /usr/local/bin
sudo cp $EXECUTABLE $DESTINATION

# Imposta una variabile d'ambiente zcli_enhancer puntando al /usr/local/bin
set -Ux ZCLI_ENHANCER /usr/local/bin

echo "zcli_enhancer installed with success"
