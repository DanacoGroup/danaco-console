// Punkt wejścia narzędzia oglądania listów: składa postać graficzną listu
// aktywacji danymi przykładowymi i odkłada ją plikiem do obejrzenia
// przeglądarką.
//
// Wysyłka dołącza znaki marki częściami listu i wskazuje je odwołaniem `cid:`,
// którego przeglądarka nie zna. Podgląd wpisuje więc oba obrazy wprost w treść
// — jedyna różnica wobec listu, który dostaje Operator.
package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"danacoconsole/server/internal/core/listy"
)

func main() {
	sciezka := flag.String("plik", "", "plik, do którego odkłada się podgląd; bez niego treść idzie na wyjście")
	flag.Parse()

	wTresc := func(dane []byte) string {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(dane)
	}
	strona := listy.PodgladDanymiPrzykladowymi(wTresc(listy.ZnakJasny), wTresc(listy.ZnakCiemny))

	if *sciezka == "" {
		fmt.Print(strona)
		return
	}
	if err := os.WriteFile(*sciezka, []byte(strona), 0o644); err != nil {
		fmt.Fprintln(os.Stderr, "nie odłożono podglądu:", err)
		os.Exit(1)
	}
}
