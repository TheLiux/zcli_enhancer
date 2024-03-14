#!/bin/bash

# Il percorso dell'eseguibile buildato da Go
EXECUTABLE="./zcli_enhancer"

# Il percorso di destinazione in /usr/local/bin
DESTINATION="/usr/local/bin/zcli_enhancer"

# Copia l'eseguibile in /usr/local/bin
sudo cp $EXECUTABLE $DESTINATION

# Imposta una variabile d'ambiente MYTOOL_HOME puntando al /usr/local/bin
echo 'export MYTOOL_HOME=/usr/local/bin' >> ~/.bash_profile

# Ricarica .bash_profile per applicare le modifiche
source ~/.bash_profile

echo "zcli_enhancer installed with success"
