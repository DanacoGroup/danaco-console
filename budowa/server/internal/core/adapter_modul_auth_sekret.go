// Postać, w jakiej bramka trzyma hasło i PIN, oraz wytworzenie tokenu sesji
// bramki.
//
// Sejf poświadczeń (`dane.SejfPlikowy`) jest schowkiem, nie funkcją skrótu:
// kładzie napis i oddaje napis, bo poświadczenie kanału modelu (klucz API) musi
// wyjść z powrotem w postaci użytecznej. Hasło bramki jest czymś odwrotnym —
// nie ma prawa wyjść ani jawnie, ani odwracalnie. Postać zapisu składa więc ten
// plik, wyłącznie z biblioteki standardowej Go:
//
//	crypto/pbkdf2   — PBKDF2 z RFC 8018, w wydaniu standardowym Go 1.24+;
//	crypto/sha256   — funkcja skrótu pod HMAC;
//	crypto/rand     — sól i token;
//	crypto/subtle   — porównanie w czasie stałym.
//
// Argon2id byłby doborem lepszym, ale mieszka w `golang.org/x/crypto`, którego
// `go.mod` rdzenia nie zaciąga.
//
// Zapis jest samoopisujący: napis kładziony w sejfie niesie nazwę funkcji,
// liczbę obrotów, sól i skrót. Dzięki temu podniesienie liczby obrotów albo
// zmiana funkcji nie unieważnia haseł już ustawionych — sprawdzenie czyta
// parametry z zapisu, a nie ze stałej.
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
	// nazwaFunkcjiSkrotu znakuje postać zapisu. Zapis o innej nazwie nie jest
	// odrzucany po cichu — sprawdzenie odmawia wprost, bo cisza znaczyłaby
	// „hasło się nie zgadza" tam, gdzie prawdą jest „zapisu nie rozumiem".
	nazwaFunkcjiSkrotu = "pbkdf2-sha256"
	// obrotySkrotu — koszt wyprowadzenia. Wartość z zalecenia OWASP dla
	// PBKDF2-HMAC-SHA256 (2023): 600 000 obrotów.
	obrotySkrotu = 600_000
	// dlugoscSoli i dlugoscSkrotu w bajtach; sól z RFC 8018 wymaga co najmniej
	// ośmiu, bierzemy szesnaście.
	dlugoscSoli   = 16
	dlugoscSkrotu = 32
	// dlugoscTokenu w bajtach. Token jest poświadczeniem na okaziciela, więc
	// jego jedyną obroną jest entropia — 256 bitów ze źródła kryptograficznego.
	dlugoscTokenu = 32
)

// zapisSekretu składa postać, w której sekret trafia do sejfu. Sekretu w tej
// postaci nie da się odwrócić: wraca z niej wyłącznie odpowiedź „zgadza się
// albo nie".
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

// sekretZgadzaSie sprawdza sekret względem zapisu z sejfu. Parametry bierze
// z zapisu, nie ze stałych — hasło ustawione przy niższej liczbie obrotów ma
// dalej działać. Zapis nieczytelny daje błąd, nie ciche „nie zgadza się".
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
	// Porównanie w czasie stałym: różnica czasu odpowiedzi zdradzałaby, ile
	// pierwszych bajtów zgadło się przy próbie.
	return subtle.ConstantTimeCompare(wyliczony, oczekiwany) == 1, nil
}

// nowyTokenBramki wytwarza token sesji. Awaria źródła losowości kończy
// czynność błędem: token przewidywalny byłby wpuszczeniem obcego do bramki.
// Przy identyfikatorach komunikatów, gdzie wystarcza licznik, brak losowości
// przechodzi dalej — tu nie.
func nowyTokenBramki() (string, error) {
	surowy := make([]byte, dlugoscTokenu)
	if _, err := rand.Read(surowy); err != nil {
		return "", fmt.Errorf("core: bramka: brak losowości na token sesji: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(surowy), nil
}

// skrotTokenu zamienia token na jego rozpoznanie w bazie. Token pochodzi ze
// źródła kryptograficznego, więc materiału do zgadywania nie ma i skrót bez
// soli oraz bez rozciągania wystarcza — inaczej niż przy haśle, które wymyśla
// człowiek.
func skrotTokenu(token string) string {
	suma := sha256.Sum256([]byte(token))
	return hex.EncodeToString(suma[:])
}
