#!/bin/sh
# Skrypt dowodzi, że wpięcie nałożone na kopie dwóch plików pakietu models przechodzi test pakietu zewnetrzne, a oryginalne pliki models pozostają nietknięte, co skrypt sprawdza sumą kontrolną na końcu.
set -e

KAT=server/internal/zewnetrzne/.dowod-wpiecia
ROBOCZY=$(mktemp -d)
MODELS=server/internal/models

PRZED=$(md5sum $MODELS/adapter_api.go $MODELS/zadanie_api.go)

cp $MODELS/adapter_api.go $MODELS/zadanie_api.go "$ROBOCZY/"
cp $KAT/zastosuj.py "$ROBOCZY/"
python3 "$ROBOCZY/zastosuj.py"

cat > "$ROBOCZY/overlay.json" <<EOF
{"Replace":{
  "$PWD/$MODELS/adapter_api.go": "$ROBOCZY/adapter_api.go",
  "$PWD/$MODELS/zadanie_api.go": "$ROBOCZY/zadanie_api.go"
}}
EOF

echo
echo "===== BEZ WPIECIA — musi byc CZERWONE (mierzy cos, czego jeszcze nie ma) ====="
if go test -count=1 -tags wpiecie_s5 ./server/internal/zewnetrzne/... >/dev/null 2>&1; then
	echo "UWAGA: zielone bez wpiecia — albo wpiecie juz nalozono, albo dowod przestal mierzyc."
else
	echo "czerwone, zgodnie z oczekiwaniem"
fi

echo
echo "===== Z WPIECIEM (overlay) — musi byc ZIELONE ====="
go test -count=1 -tags wpiecie_s5 -overlay="$ROBOCZY/overlay.json" ./server/internal/zewnetrzne/...

echo
PO=$(md5sum $MODELS/adapter_api.go $MODELS/zadanie_api.go)
if [ "$PRZED" = "$PO" ]; then
	echo "===== ORYGINALY models/ NIETKNIETE (md5 bez zmiany) ====="
else
	echo "===== BLAD: pliki models/ ZMIENILY SIE! ====="
	exit 1
fi
rm -rf "$ROBOCZY"
