// Odpowiedzialność pliku: droga poświadczenia od żądania do sejfu.
//
// Poświadczenie wchodzi żądaniem (`credential` przy koncie i przy punkcie
// dostępu) i nie wychodzi nigdy. Baza zna wyłącznie odwołanie — nazwę wpisu
// w sejfie albo ścieżkę profilu; kolumny na treść sekretu nie ma w schemacie
// w ogóle, a odpowiedź kontraktu niesie co najwyżej znacznik `hasCredential`
// albo `credentialRef`.
//
// Rdzeń nie jest sejfem i nie udaje sejfu. Zamiana sekretu na odwołanie należy
// do portu SejfPoswiadczen wypełnianego przy montażu. Gdy sejfu nie wpięto,
// poświadczenie jest odrzucane w tym samym wywołaniu: nie trafia do bazy, do
// dziennika ani do odpowiedzi, a byt powstaje bez odwołania i pracuje dalej.
// Cicha utrata sekretu byłaby gorsza od jawnego braku sejfu, więc warstwa wyżej
// widzi to po pustym odwołaniu i po `hasCredential = false`.
package core

import "context"

// SejfPoswiadczen jest portem magazynu sekretów. Bytem jest
// konto albo punkt dostępu wskazany swoim trwałym identyfikatorem.
type SejfPoswiadczen interface {
	// Zapisz umieszcza poświadczenie w sejfie i zwraca odwołanie do niego.
	Zapisz(ctx context.Context, byt, poswiadczenie string) (string, error)
	// Odczytaj zwraca poświadczenie bytu i znacznik, czy wpis istnieje. Brak
	// wpisu nie jest błędem — daje pusty sekret i false. Bytem jest tu
	// surowy klucz wpisu, bez przedrostka odwołania.
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
