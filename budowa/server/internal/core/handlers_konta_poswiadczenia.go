// Plik prowadzi drogę poświadczenia od żądania do sejfu; baza zna wyłącznie
// odwołanie do sekretu, a rdzeń nie przechowuje treści sekretu w żadnej
// postaci.
package core

import "context"

// SejfPoswiadczen jest portem magazynu sekretów. Bytem jest
// konto albo punkt dostępu wskazany swoim trwałym identyfikatorem.
type SejfPoswiadczen interface {
	// Zapisz umieszcza poświadczenie w sejfie i zwraca odwołanie do niego.
	Zapisz(ctx context.Context, byt, poswiadczenie string) (string, error)
	// Odczytaj zwraca poświadczenie bytu i znacznik istnienia wpisu; brak
	// wpisu nie jest błędem.
	Odczytaj(ctx context.Context, byt string) (string, bool)
	// Usun kasuje poświadczenie bytu. Brak wpisu nie jest błędem.
	Usun(ctx context.Context, byt string) error
}

// odwolaniePoswiadczenia zamienia poświadczenie żądania na odwołanie zapisywane
// w bazie. Brak poświadczenia w żądaniu znaczy „bez zmiany" i daje brak
// odwołania; brak sejfu daje to samo, bo sekretu nie ma gdzie odłożyć.
func odwolaniePoswiadczenia(ctx context.Context, sejf SejfPoswiadczen,
	byt string, poswiadczenie *string) (*string, error) {

	if sejf == nil || poswiadczenie == nil || *poswiadczenie == "" {
		return nil, nil
	}
	odwolanie, err := sejf.Zapisz(ctx, byt, *poswiadczenie)
	if err != nil {
		return nil, err
	}
	if odwolanie == "" {
		return nil, nil
	}
	return &odwolanie, nil
}

// usunPoswiadczenie kasuje sekret bytu, który przestał istnieć. Niepowodzenie
// kasowania nie unieważnia usunięcia bytu — wiersz już zniknął, a sejf zgłosi
// swój błąd osobno.
func usunPoswiadczenie(ctx context.Context, sejf SejfPoswiadczen, byt string) {
	if sejf == nil || byt == "" {
		return
	}
	_ = sejf.Usun(ctx, byt)
}
