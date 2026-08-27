// Pakiet models trzyma pojęcie kanału modelu i jednolitego strumienia fragmentów
// oraz rejestr kanałów budowany z danych, nie z kodu. Kanały echo, CLI i API
// mówią tym samym strumieniem i tym samym zestawem rodzajów fragmentów.
package models

import "danacoconsole/shared"

// Rodzaje diagnostyczne strumienia niosą prowenancję wywołania i metadane konta
// użytego przez kanał przy rotacji kont; obie wartości pochodzą z kontraktu,
// a pakiet ich nie definiuje samodzielnie.
const (
	RodzajProwenancji = shared.ChunkKindProvenance
	RodzajKonta       = shared.ChunkKindAccount
)
