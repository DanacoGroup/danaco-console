// Pakiet models trzyma pojęcie kanału modelu i jednolitego strumienia fragmentów
// oraz rejestr kanałów budowany z danych, nie z kodu.
// Wszystkie kanały — echo, CLI, API — mówią tym samym strumieniem i tym samym
// zestawem rodzajów fragmentów; odbiorca nie zna dostawcy.
package models

import "danacoconsole/shared"

// Rodzaje diagnostyczne strumienia — prowenancja wywołania i metadane konta
// użytego przez kanał przy rotacji kont. Obie wartości pochodzą z kontraktu;
// pakiet niczego nie definiuje samodzielnie.
//
// W kontrakcie są to wartości przelotowe: nie mają odpowiednika w kolumnie
// wiadomosc.rodzaj_tresci, więc nie wchodzą do WartosciBazyChunkKind i nie
// trafiają do bazy, której CHECK dopuszcza wyłącznie rodzaje treści rozmowy.
//
// Rodzaje treści rozmowy nie mają tu drugiej nazwy: bierze się je wprost ze
// `shared`, a odwzorowanie na kolumnę bazy robi pakiet dane.
const (
	RodzajProwenancji = shared.ChunkKindProvenance
	RodzajKonta       = shared.ChunkKindAccount
)
