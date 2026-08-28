// Plik niesie port rodziny control.* — przejęcia i oddania sterowania zleceniem — wraz z kształtem
// trzech żądań i trzech odpowiedzi oraz wpięciem tych komend do rejestru rdzenia.
package core

import (
	"context"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// ZadaniePrzejeciaSterowania jest żądaniem komendy przejęcia sterowania zleceniem od Koordynatora tury.
type ZadaniePrzejeciaSterowania struct {
	// WindowId — okno koordynatora prowadzące zlecenie.
	WindowId string `json:"windowId"`
	// Reason — powód przejęcia, wchodzi do śladu i niczego nie blokuje.
	Reason *string `json:"reason,omitempty"`
}

// WynikPrzejeciaSterowania jest odpowiedzią komendy przejęcia sterowania zleceniem przez Operatora rdzenia.
type WynikPrzejeciaSterowania struct {
	// Loop — bieg po przejęciu, zatrzymany, z dorobkiem Koordynatora nietkniętym.
	Loop shared.LoopState `json:"loop"`
	// Handover — zapis, który poszedł do trwałego śladu.
	Handover ZapisSterowania `json:"handover"`
}

// ZadanieOddaniaSterowania jest żądaniem komendy oddania sterowania zleceniem z powrotem Koordynatorowi.
type ZadanieOddaniaSterowania struct {
	// WindowId — okno koordynatora.
	WindowId string `json:"windowId"`
	// Note — notatka Operatora dla Koordynatora podejmującego bieg.
	Note *string `json:"note,omitempty"`
}

// WynikOddaniaSterowania jest odpowiedzią komendy oddania sterowania zleceniem z powrotem Koordynatorowi.
type WynikOddaniaSterowania struct {
	// Loop — bieg po oddaniu; podjęty, nie zaczęty od nowa.
	Loop shared.LoopState `json:"loop"`
	// Handover — zapis oddania.
	Handover ZapisSterowania `json:"handover"`
}

// ZadanieOdczytuSterowania jest żądaniem komendy odczytu bieżącego stanu sterowania zleceniem koordynatora.
type ZadanieOdczytuSterowania struct {
	// WindowId — okno koordynatora.
	WindowId string `json:"windowId"`
	// Limit — ile ostatnich zapisów historii; brak schodzi na wartość domyślną warstwy danych.
	Limit *int `json:"limit,omitempty"`
}

// WynikOdczytuSterowania jest odpowiedzią komendy odczytu bieżącego stanu sterowania zleceniem koordynatora.
type WynikOdczytuSterowania struct {
	// Loop — bieg zlecenia.
	Loop shared.LoopState `json:"loop"`
	// Controller — kto prowadzi zlecenie w tej chwili; pole stoi osobno, dopóki stan biegu go nie niesie.
	Controller Sterujacy `json:"controller"`
	// History — przejęcia i oddania, najnowsze pierwsze; wykaz pusty jest odpowiedzią, nie brakiem.
	History []ZapisSterowania `json:"history"`
}

// ZapisSterowania to jedno przejęcie albo jedno oddanie sterowania zleceniem w historii biegu koordynatora.
type ZapisSterowania struct {
	// WindowId — okno koordynatora, którego dotyczy.
	WindowId string `json:"windowId"`
	// Controller — kto steruje PO tym zapisie.
	Controller Sterujacy `json:"controller"`
	// At — chwila w milisekundach epoki.
	At int64 `json:"at"`
	// Actor — identyfikator klienta Operatora; pusty znaczy, że rdzeń nie potrafił go rozstrzygnąć.
	Actor *string `json:"actor,omitempty"`
	// Reason — powód przejęcia albo notatka przy oddaniu.
	Reason *string `json:"reason,omitempty"`
}

// PrzejecieSterowania jest portem rodziny control.* obsługującym przejęcie, oddanie i odczyt sterowania.
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

// NazwyPrzejeciaSterowania niesie trzy nazwy komend rodziny control.*, wstrzykiwane montażem, bo
// katalog nazw należy do kontraktu, którego ta rodzina jeszcze nie zna. Nazwa pominięta znaczy,
// że kontrakt tej komendy nie ogłasza, i komenda nie powstaje.
type NazwyPrzejeciaSterowania struct {
	Przejecie shared.MessageType
	Oddanie   shared.MessageType
	Odczyt    shared.MessageType
}

// zarejestrujPrzejecieSterowania wpina trzy komendy rodziny control.*. Port niewpięty nie jest
// ciszą: komendy odpowiadają błędem wewnętrznym z nazwą brakującego bytu, zamiast odpowiedzią
// nieznanej komendy.
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

// zarejestrujPrzejecieSterowaniaNiewpiete wpina odmowę rdzenia na miejsce brakujących obsługiwaczy komend.
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

// bladPrzejeciaSterowaniaNiewpietego składa odmowę rdzenia, gdy port rodziny control.* nie jest wpięty.
func bladPrzejeciaSterowaniaNiewpietego() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"przejęcie sterowania: rdzeń nie niesie portu przejęcia sterowania — "+
			"zlecenia nie ma jak przejąć ani oddać Koordynatorowi"))
}
