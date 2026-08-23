// Port rodziny `control.*` — przejęcia i oddania sterowania zleceniem — wraz
// z kształtem jej trzech żądań i trzech odpowiedzi oraz wpięciem tych komend
// do rejestru rdzenia.
//
// Kształty stoją tutaj, a nie w `shared`, bo generowany kontrakt nie zna
// jeszcze żadnej z tych trzech nazw. Po wniesieniu komend do kontraktu typy
// stają się aliasami `shared.Control*`, a katalog wartości wraca do kontraktu.
//
// Nazwy komend są argumentem, a nie literałem: rejestr nie zawiera nazw
// własnych, wstrzykuje je montaż. Nazwa pusta niczego nie rejestruje
// (`Rejestr.Zarejestruj` ją pomija), więc dopóki kontrakt nie niesie tych
// komend, rdzeń nie ogłasza zdolności, której kontrakt nie zna.
//
// Rodzina nie bierze nadajnika i nie rozgłasza własnego zdarzenia. `LoopState`
// wychodzi już rodzinami `window.state.changed` i `progress.changed` oraz
// komendą `window.state.get`, a przejęcie zmienia właśnie `LoopState`:
// `Przejmij` woła `petla.Zatrzymaj`, `Oddaj` woła `petla.Wznow`, a każde z nich
// rozgłasza stan biegu obserwatorom pętli.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZadaniePrzejeciaSterowania jest żądaniem `control.takeover`.
type ZadaniePrzejeciaSterowania struct {
	// WindowId — okno koordynatora prowadzące zlecenie.
	WindowId string `json:"windowId"`
	// Reason — powód przejęcia. Wchodzi do śladu i niczego nie blokuje:
	// przycisk Operatora jest czynny także bez wyjaśnienia.
	Reason *string `json:"reason,omitempty"`
}

// WynikPrzejeciaSterowania jest odpowiedzią `control.takeover`.
type WynikPrzejeciaSterowania struct {
	// Loop — bieg po przejęciu. Zatrzymany, z dorobkiem Koordynatora
	// nietkniętym: liczba obiegów, ostatni wykonawca i powód tury zostają.
	Loop shared.LoopState `json:"loop"`
	// Handover — zapis, który poszedł do trwałego śladu.
	Handover ZapisSterowania `json:"handover"`
}

// ZadanieOddaniaSterowania jest żądaniem `control.release`.
type ZadanieOddaniaSterowania struct {
	// WindowId — okno koordynatora.
	WindowId string `json:"windowId"`
	// Note — notatka Operatora dla Koordynatora podejmującego bieg.
	Note *string `json:"note,omitempty"`
}

// WynikOddaniaSterowania jest odpowiedzią `control.release`.
type WynikOddaniaSterowania struct {
	// Loop — bieg po oddaniu; podjęty, nie zaczęty od nowa.
	Loop shared.LoopState `json:"loop"`
	// Handover — zapis oddania.
	Handover ZapisSterowania `json:"handover"`
}

// ZadanieOdczytuSterowania jest żądaniem `control.get`.
type ZadanieOdczytuSterowania struct {
	// WindowId — okno koordynatora.
	WindowId string `json:"windowId"`
	// Limit — ile ostatnich zapisów historii. Brak schodzi na wartość domyślną
	// warstwy danych; „bez granicy" nie jest żądaniem stawianym świadomie.
	Limit *int `json:"limit,omitempty"`
}

// WynikOdczytuSterowania jest odpowiedzią `control.get`.
type WynikOdczytuSterowania struct {
	// Loop — bieg zlecenia.
	Loop shared.LoopState `json:"loop"`
	// Controller — kto prowadzi zlecenie w tej chwili. Pole stoi osobno, dopóki
	// `LoopState` nie niesie sterującego: bez niego `control.get` na biegu
	// nigdy nie przejętym nie miałby czym odpowiedzieć.
	Controller Sterujacy `json:"controller"`
	// History — przejęcia i oddania, najnowsze pierwsze. Wykaz pusty znaczy
	// „biegu nie tknęła ręka człowieka" i jest odpowiedzią, nie brakiem.
	History []ZapisSterowania `json:"history"`
}

// ZapisSterowania to jedno przejęcie albo jedno oddanie sterowania.
type ZapisSterowania struct {
	// WindowId — okno koordynatora, którego dotyczy.
	WindowId string `json:"windowId"`
	// Controller — kto steruje PO tym zapisie.
	Controller Sterujacy `json:"controller"`
	// At — chwila w milisekundach epoki.
	At int64 `json:"at"`
	// Actor — identyfikator klienta Operatora. Pusty znaczy „rdzeń nie potrafił
	// tego rozstrzygnąć", tak samo jak `actorClientId` w zdarzeniach; tożsamość
	// klienta z gniazda ustala warstwa transportu.
	Actor *string `json:"actor,omitempty"`
	// Reason — powód przejęcia albo notatka przy oddaniu.
	Reason *string `json:"reason,omitempty"`
}

// PrzejecieSterowania jest portem rodziny `control.*`.
type PrzejecieSterowania interface {
	// Przejmij obsługuje `control.takeover`.
	Przejmij(ctx context.Context, z ZadaniePrzejeciaSterowania) (WynikPrzejeciaSterowania, error)
	// Oddaj obsługuje `control.release`.
	Oddaj(ctx context.Context, z ZadanieOddaniaSterowania) (WynikOddaniaSterowania, error)
	// Ster obsługuje `control.get`.
	Ster(ctx context.Context, z ZadanieOdczytuSterowania) (WynikOdczytuSterowania, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ PrzejecieSterowania = (*adapterPrzejeciaSterowania)(nil)

// NazwyPrzejeciaSterowania niesie trzy nazwy komend rodziny `control.*`.
//
// Katalog nazw należy do kontraktu, a kontrakt tej rodziny jeszcze nie zna —
// nazwy przychodzą z montażu. Nazwa pominięta znaczy „kontrakt tej komendy nie
// ogłasza" i komenda nie powstaje.
type NazwyPrzejeciaSterowania struct {
	Przejecie shared.MessageType
	Oddanie   shared.MessageType
	Odczyt    shared.MessageType
}

// zarejestrujPrzejecieSterowania wpina trzy komendy rodziny `control.*`.
//
// Port niewpięty nie jest ciszą: komendy zostają znane rdzeniowi i odpowiadają
// `internal_error` z nazwą brakującego bytu, zamiast `*.unknown` albo cichej
// zgody na przejęcie, które się nie odbyło.
func zarejestrujPrzejecieSterowania(r *Rejestr, p PrzejecieSterowania, n NazwyPrzejeciaSterowania) {
	if r == nil {
		return
	}
	if p == nil {
		zarejestrujPrzejecieSterowaniaNiewpiete(r, n)
		return
	}
	r.Zarejestruj(n.Przejecie, obsluz(p.Przejmij))
	r.Zarejestruj(n.Oddanie, obsluz(p.Oddaj))
	r.Zarejestruj(n.Odczyt, obsluz(p.Ster))
}

// zarejestrujPrzejecieSterowaniaNiewpiete wpina odmowę na miejsce obsługiwaczy.
func zarejestrujPrzejecieSterowaniaNiewpiete(r *Rejestr, n NazwyPrzejeciaSterowania) {
	r.Zarejestruj(n.Przejecie,
		obsluz(func(context.Context, ZadaniePrzejeciaSterowania) (WynikPrzejeciaSterowania, error) {
			return WynikPrzejeciaSterowania{}, bladPrzejeciaSterowaniaNiewpietego()
		}))
	r.Zarejestruj(n.Oddanie,
		obsluz(func(context.Context, ZadanieOddaniaSterowania) (WynikOddaniaSterowania, error) {
			return WynikOddaniaSterowania{}, bladPrzejeciaSterowaniaNiewpietego()
		}))
	r.Zarejestruj(n.Odczyt,
		obsluz(func(context.Context, ZadanieOdczytuSterowania) (WynikOdczytuSterowania, error) {
			return WynikOdczytuSterowania{}, bladPrzejeciaSterowaniaNiewpietego()
		}))
}

// bladPrzejeciaSterowaniaNiewpietego składa odmowę rdzenia bez portu `control.*`.
func bladPrzejeciaSterowaniaNiewpietego() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"przejęcie sterowania: rdzeń nie niesie portu przejęcia sterowania — "+
			"zlecenia nie ma jak przejąć ani oddać Koordynatorowi"))
}
