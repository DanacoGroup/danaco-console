// Postać, w jakiej bramka trzyma hasło i PIN, wyłącznie biblioteką standardową Go, oraz
// wytworzenie tokenu sesji bramki, samoopisującym zapisem niosącym parametry wyprowadzenia.
package core

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

const (
	// nazwaFunkcjiSkrotu znakuje postać zapisu; zapis o innej nazwie sprawdzenie odmawia wprost, nie po cichu.
	nazwaFunkcjiSkrotu = "pbkdf2-sha256"
	// obrotySkrotu to koszt wyprowadzenia, wartość z zalecenia OWASP dla PBKDF2-HMAC-SHA256 z roku dwa tysiące dwudziestego trzeciego.
	obrotySkrotu = 600_000
	// dlugoscSoli i dlugoscSkrotu w bajtach; sól zgodna z RFC 8018 wymaga co najmniej ośmiu bajtów długości zapisu.
	dlugoscSoli   = 16
	dlugoscSkrotu = 32
	// dlugoscTokenu w bajtach; token jest poświadczeniem na okaziciela, jego jedyną obroną jest entropia losowości.
	dlugoscTokenu = 32
)

// zapisSekretu składa postać, w której sekret trafia do sejfu; nie da się jej odwrócić do napisu jawnego.
func zapisSekretu(sekret string) (string, error) {
	sol := make([]byte, dlugoscSoli)
	if _, err := rand.Read(sol); err != nil {
		return "", fmt.Errorf("core: bramka: brak losowości na sól hasła: %w", err)
	}
	skrot, err := pbkdf2.Key(sha256.New, sekret, sol, obrotySkrotu, dlugoscSkrotu)
	if err != nil {
		return "", fmt.Errorf("core: bramka: nie można wyprowadzić skrótu hasła: %w", err)
	}
	return strings.Join([]string{
		nazwaFunkcjiSkrotu,
		strconv.Itoa(obrotySkrotu),
		base64.RawStdEncoding.EncodeToString(sol),
		base64.RawStdEncoding.EncodeToString(skrot),
	}, "$"), nil
}

// sekretZgadzaSie sprawdza sekret względem zapisu z sejfu; parametry bierze z zapisu, nie
// ze stałych, bo hasło ustawione przy niższej liczbie obrotów ma dalej działać.
func sekretZgadzaSie(zapis, sekret string) (bool, error) {
	czesci := strings.Split(zapis, "$")
	if len(czesci) != 4 || czesci[0] != nazwaFunkcjiSkrotu {
		return false, fmt.Errorf("core: bramka: zapis sekretu w postaci nierozpoznanej")
	}
	obroty, err := strconv.Atoi(czesci[1])
	if err != nil || obroty <= 0 {
		return false, fmt.Errorf("core: bramka: zapis sekretu bez czytelnej liczby obrotów")
	}
	sol, err := base64.RawStdEncoding.DecodeString(czesci[2])
	if err != nil {
		return false, fmt.Errorf("core: bramka: zapis sekretu bez czytelnej soli: %w", err)
	}
	oczekiwany, err := base64.RawStdEncoding.DecodeString(czesci[3])
	if err != nil {
		return false, fmt.Errorf("core: bramka: zapis sekretu bez czytelnego skrótu: %w", err)
	}
	wyliczony, err := pbkdf2.Key(sha256.New, sekret, sol, obroty, len(oczekiwany))
	if err != nil {
		return false, fmt.Errorf("core: bramka: nie można wyprowadzić skrótu hasła: %w", err)
	}
	// Porównanie w czasie stałym: różnica czasu odpowiedzi zdradzałaby dopasowane bajty.
	return subtle.ConstantTimeCompare(wyliczony, oczekiwany) == 1, nil
}

// nowyTokenBramki wytwarza token sesji; awaria źródła losowości kończy czynność błędem,
// bo token przewidywalny byłby wpuszczeniem obcego do bramki.
func nowyTokenBramki() (string, error) {
	surowy := make([]byte, dlugoscTokenu)
	if _, err := rand.Read(surowy); err != nil {
		return "", fmt.Errorf("core: bramka: brak losowości na token sesji: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(surowy), nil
}

/*
ZNAKOW_KODU nie stoi tu przypadkiem: okno rejestracji ma sześć pól na kod
(`wejscie/ekrany/dostep.ts`, `ZNAKOW_DROGI`), a Operator przepisuje go ręcznie
z wiadomości. Token sesji ma 43 znaki i do tych pól nie wchodzi — dlatego kod
potwierdzenia jest osobnym bytem, krótkim i cyfrowym.

Sześć cyfr daje milion możliwości. Zgadywania pilnują cztery rzeczy naraz:
dławik prób kluczowany połączeniem, licznik pomyłek zamykający drogę po
`pulapProbDrogi` próbach, godzina ważności oraz zamknięcie drogi po pierwszym
użyciu — bez nich krótki kod byłby słabością, nie ułatwieniem.
*/
const znakowKoduPotwierdzenia = 6

// granicaLosowaniaCyfry — bajty od tej wartości w górę odpadają. Dwieście
// pięćdziesiąt sześć nie dzieli się przez dziesięć, więc reszta z dzielenia
// dawałaby cyfrom od zera do pięciu szansę większą niż pozostałym; kod tak
// krótki nie ma z czego oddać tej różnicy.
const granicaLosowaniaCyfry = 250

// nowyKodPotwierdzenia losuje kod przepisywany przez Operatora z wiadomości.
// Cyfry, nie litery: kod czyta się z ekranu telefonu i przepisuje na klawiaturze,
// a litery podobne do cyfr (O i zero, l i jedynka) mnożą pomyłki.
func nowyKodPotwierdzenia() (string, error) {
	cyfry := make([]byte, 0, znakowKoduPotwierdzenia)
	surowy := make([]byte, znakowKoduPotwierdzenia)
	for len(cyfry) < znakowKoduPotwierdzenia {
		if _, err := rand.Read(surowy); err != nil {
			return "", fmt.Errorf("core: brak losowości na kod potwierdzenia: %w", err)
		}
		for _, bajt := range surowy {
			if bajt >= granicaLosowaniaCyfry {
				continue
			}
			cyfry = append(cyfry, '0'+bajt%10)
			if len(cyfry) == znakowKoduPotwierdzenia {
				break
			}
		}
	}
	return string(cyfry), nil
}

// skrotTokenu zamienia token na jego rozpoznanie w bazie; token pochodzi ze źródła
// kryptograficznego, więc skrót bez soli i rozciągania wystarcza, inaczej niż przy haśle.
func skrotTokenu(token string) string {
	suma := sha256.Sum256([]byte(token))
	return hex.EncodeToString(suma[:])
}
