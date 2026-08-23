// Komendy `terminal.key.generate`, `terminal.key.import`, `terminal.key.list`
// i `terminal.key.remove` — wykaz kluczy SSH znanych rdzeniowi.
//
// ── Biblioteka wkompilowana, nie `ssh-keygen` ───────────────────────────────
// Para kluczy powstaje w Go: `crypto/ed25519`, `crypto/rsa` i `crypto/ecdsa`
// wytwarzają materiał, a `golang.org/x/crypto/ssh` zapisuje go w postaci OpenSSH
// i liczy odcisk. `ssh-keygen` byłby tu programem spoza instalki wołanym po to,
// żeby zrobić rzecz, którą biblioteka standardowa robi w kilku wierszach — a
// wytworzenie klucza jest czynnością, bez której książka hostów przestaje mieć
// czym się łączyć. Ta sama decyzja co przy git (go-git), PDF (pdfcpu)
// i wyszukiwaniu (regexp).
//
// ── Czego rdzeń nie robi z kluczem ──────────────────────────────────────────
// Klucz prywatny nie opuszcza dysku maszyny rdzenia w żadną stronę: wytworzenie
// oddaje sam odcisk i klucz publiczny, wciągnięcie do wykazu bierze ŚCIEŻKĘ, nie
// treść, a wykaz nie ma pola, w którym materiał tajny mógłby się znaleźć.
//
// ── Hasło klucza ────────────────────────────────────────────────────────────
// Kontrakt każe podawać hasło ODWOŁANIEM do sejfu, nigdy treścią — i to jest
// właściwe. Rdzeń nie ma dziś czytnika sejfu (ten sam brak, który zmienna
// środowiska karty nazywa w `adapter_modul_terminal_powloki.go`), więc żądanie
// z odwołaniem kończy się odmową NAZYWAJĄCĄ ten brak. Wytworzenie klucza bez
// hasła w odpowiedzi na prośbę o klucz z hasłem byłoby cichym obniżeniem
// ochrony — a to jest gorsze niż odmowa.
package core

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekKlucza znakuje identyfikator wpisu wykazu kluczy.
const przedrostekKlucza = "tkey-"

// katalogKluczy jest podkatalogiem katalogu danych rdzenia, w którym leżą klucze
// wytworzone przez moduł. Osobny katalog, a nie `~/.ssh` konta procesu: klucz
// wytworzony przez produkt ma być rozpoznawalny i usuwalny wraz z produktem,
// a konfiguracja OpenSSH Operatora nie ma prawa zmienić się sama.
const katalogKluczy = "klucze-ssh"

// dlugoscKluczaRSA jest długością klucza RSA w bitach. 3072 to dolna granica
// zalecana dziś dla kluczy RSA; 2048 nie jest już wyborem, który wolno zapisać
// w produkcie na następne lata.
const dlugoscKluczaRSA = 3072

// WytworzKlucz obsługuje `terminal.key.generate`.
func (a *adapterTerminala) WytworzKlucz(ctx context.Context,
	z shared.TerminalKeyGenerateRequest) (shared.TerminalKeyGenerateResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalKeyGenerateResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.TerminalKeyGenerateResponse{}, bladZadaniaTerminala(
			"wytworzenie klucza wymaga jego nazwy w wykazie")
	}
	if strings.TrimSpace(wartoscTekstu(z.PassphraseRef)) != "" {
		return shared.TerminalKeyGenerateResponse{}, bladZadaniaTerminala(
			"hasło klucza wskazuje sejf, a rdzeń nie ma jego czytnika; klucz z hasłem " +
				"trzeba dziś wytworzyć poza produktem i wciągnąć do wykazu komendą " +
				"terminal.key.import — wytworzenie klucza BEZ hasła w odpowiedzi na tę prośbę " +
				"byłoby cichym obniżeniem ochrony")
	}

	prywatny, publiczny, err := paraKluczy(z.KeyType)
	if err != nil {
		return shared.TerminalKeyGenerateResponse{}, err
	}
	blok, err := ssh.MarshalPrivateKey(prywatny, strings.TrimSpace(wartoscTekstu(z.Comment)))
	if err != nil {
		return shared.TerminalKeyGenerateResponse{}, bladWykonaniaTerminala(
			"nie można zapisać klucza w postaci OpenSSH: " + err.Error())
	}
	kod := nowyIdentyfikator(przedrostekKlucza)
	jawny := trescKluczaJawnego(publiczny, wartoscTekstu(z.Comment))
	sciezka, err := a.zapiszPlikiKlucza(kod, pem.EncodeToMemory(blok), []byte(jawny+"\n"))
	if err != nil {
		return shared.TerminalKeyGenerateResponse{}, err
	}

	wiersz := dane.KluczTerminala{
		Kod:        kod,
		Nazwa:      nazwa,
		Rodzaj:     z.KeyType,
		Odcisk:     ssh.FingerprintSHA256(publiczny),
		KluczJawny: jawny,
		Sciezka:    sciezka,
	}
	if err := dziennik.ZapiszKlucz(ctx, wiersz); err != nil {
		// Wiersza nie ma, więc plików też nie ma prawa zostać: klucz prywatny
		// leżący poza wykazem jest materiałem, o którym nikt już nie wie.
		_ = os.Remove(sciezka)
		_ = os.Remove(sciezka + ".pub")
		return shared.TerminalKeyGenerateResponse{}, err
	}
	return shared.TerminalKeyGenerateResponse{Key: kluczKontraktu(wiersz)}, nil
}

// WciagnijKlucz obsługuje `terminal.key.import`.
//
// Klucz wskazuje się ścieżką, więc rdzeń go CZYTA, żeby powiedzieć o nim prawdę:
// jakiego jest rodzaju, jaki ma odcisk i czy jest chroniony hasłem. Wpis
// przepisany z samego żądania byłby wpisem o pliku, którego nikt nie otworzył.
func (a *adapterTerminala) WciagnijKlucz(ctx context.Context,
	z shared.TerminalKeyImportRequest) (shared.TerminalKeyImportResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalKeyImportResponse{}, err
	}
	nazwa := strings.TrimSpace(z.Name)
	sciezka := strings.TrimSpace(z.Path)
	if nazwa == "" || sciezka == "" {
		return shared.TerminalKeyImportResponse{}, bladZadaniaTerminala(
			"wciągnięcie klucza wymaga jego nazwy i ścieżki na maszynie rdzenia")
	}
	tresc, err := os.ReadFile(sciezka)
	if err != nil {
		return shared.TerminalKeyImportResponse{}, bladBrakuZasobuTerminala(
			"klucz prywatny " + sciezka)
	}

	wiersz := dane.KluczTerminala{
		Kod:     nowyIdentyfikator(przedrostekKlucza),
		Nazwa:   nazwa,
		Sciezka: sciezka,
	}
	prywatny, err := ssh.ParseRawPrivateKey(tresc)
	var chroniony *ssh.PassphraseMissingError
	switch {
	case errors.As(err, &chroniony):
		// Klucz chroniony hasłem czyta się tylko z hasłem, którego rdzeń nie ma.
		// To NIE jest powód odmowy: wykaz ma nieść taki klucz, a `ssh` odczyta go
		// sam przy połączeniu. Odcisk bierze się wtedy z klucza publicznego —
		// z pliku obok albo z części publicznej, którą niesie sam błąd.
		wiersz.Haslo = true
		wiersz.Rodzaj = shared.TerminalKeyTypeEd25519
		if chroniony.PublicKey != nil {
			wiersz.Odcisk = ssh.FingerprintSHA256(chroniony.PublicKey)
			wiersz.Rodzaj = rodzajKluczaZTypu(chroniony.PublicKey.Type())
			wiersz.KluczJawny = trescKluczaJawnego(chroniony.PublicKey, "")
		}
	case err != nil:
		return shared.TerminalKeyImportResponse{}, bladZadaniaTerminala(
			"plik " + sciezka + " nie jest kluczem prywatnym w postaci czytelnej dla rdzenia: " +
				err.Error())
	default:
		podpisujacy, err := ssh.NewSignerFromKey(prywatny)
		if err != nil {
			return shared.TerminalKeyImportResponse{}, bladZadaniaTerminala(
				"klucz " + sciezka + " ma rodzaj, którego rdzeń nie obsługuje: " + err.Error())
		}
		publiczny := podpisujacy.PublicKey()
		wiersz.Odcisk = ssh.FingerprintSHA256(publiczny)
		wiersz.Rodzaj = rodzajKluczaZTypu(publiczny.Type())
		wiersz.KluczJawny = trescKluczaJawnego(publiczny, "")
	}
	if wiersz.KluczJawny == "" {
		if jawny, err := os.ReadFile(sciezka + ".pub"); err == nil {
			wiersz.KluczJawny = strings.TrimSpace(string(jawny))
		}
	}
	if err := dziennik.ZapiszKlucz(ctx, wiersz); err != nil {
		return shared.TerminalKeyImportResponse{}, err
	}
	return shared.TerminalKeyImportResponse{Key: kluczKontraktu(wiersz)}, nil
}

// WykazKluczy obsługuje `terminal.key.list`.
func (a *adapterTerminala) WykazKluczy(ctx context.Context,
	_ shared.TerminalKeyListRequest) (shared.TerminalKeyListResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalKeyListResponse{}, err
	}
	wiersze, err := dziennik.Klucze(ctx)
	if err != nil {
		return shared.TerminalKeyListResponse{}, err
	}
	wykaz := make([]shared.TerminalSshKey, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, kluczKontraktu(wiersz))
	}
	return shared.TerminalKeyListResponse{Keys: wykaz, Total: len(wykaz)}, nil
}

// UsunKlucz obsługuje `terminal.key.remove`.
//
// Odpięcie wpisów książki hostów idzie PRZED usunięciem klucza i idzie zawsze,
// także wtedy, gdy plików nie usuwamy: wpis wskazujący klucz zdjęty z wykazu
// wskazywałby na nic, a kontrakt każe te wpisy wymienić w odpowiedzi.
func (a *adapterTerminala) UsunKlucz(ctx context.Context,
	z shared.TerminalKeyRemoveRequest) (shared.TerminalKeyRemoveResponse, error) {

	dziennik, err := a.dziennikWyposazenia()
	if err != nil {
		return shared.TerminalKeyRemoveResponse{}, err
	}
	kod := strings.TrimSpace(z.KeyId)
	if kod == "" {
		return shared.TerminalKeyRemoveResponse{}, bladZadaniaTerminala(
			"zdjęcie klucza wymaga jego wskazania")
	}
	wiersz, err := dziennik.Klucz(ctx, kod)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.TerminalKeyRemoveResponse{Removed: false}, nil
	}
	if err != nil {
		return shared.TerminalKeyRemoveResponse{}, err
	}
	odpiete, err := dziennik.OdepnijKlucz(ctx, kod)
	if err != nil {
		return shared.TerminalKeyRemoveResponse{}, err
	}
	zdjety, err := dziennik.UsunKlucz(ctx, kod)
	if err != nil {
		return shared.TerminalKeyRemoveResponse{}, err
	}
	if z.DeleteFiles != nil && *z.DeleteFiles && wiersz.Sciezka != "" {
		// Niepowodzenie usunięcia pliku nie wycofuje zdjęcia z wykazu: wpis już
		// nie istnieje, a plik zostaje na dysku, gdzie Operator go widzi.
		_ = os.Remove(wiersz.Sciezka)
		_ = os.Remove(wiersz.Sciezka + ".pub")
	}
	odpowiedz := shared.TerminalKeyRemoveResponse{Removed: zdjety}
	if len(odpiete) > 0 {
		odpowiedz.DetachedHostIds = odpiete
	}
	return odpowiedz, nil
}

// paraKluczy wytwarza materiał klucza wskazanego rodzaju.
func paraKluczy(rodzaj shared.TerminalKeyType) (any, ssh.PublicKey, error) {
	switch rodzaj {
	case shared.TerminalKeyTypeEd25519:
		jawny, tajny, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, bladWykonaniaTerminala("nie można wytworzyć klucza ed25519: " + err.Error())
		}
		publiczny, err := ssh.NewPublicKey(jawny)
		if err != nil {
			return nil, nil, bladWykonaniaTerminala(err.Error())
		}
		// Wskaźnik, nie wartość: postać OpenSSH wymaga typu ed25519.PrivateKey
		// przekazanego wskaźnikiem.
		return &tajny, publiczny, nil
	case shared.TerminalKeyTypeRsa:
		tajny, err := rsa.GenerateKey(rand.Reader, dlugoscKluczaRSA)
		if err != nil {
			return nil, nil, bladWykonaniaTerminala("nie można wytworzyć klucza RSA: " + err.Error())
		}
		publiczny, err := ssh.NewPublicKey(&tajny.PublicKey)
		if err != nil {
			return nil, nil, bladWykonaniaTerminala(err.Error())
		}
		return tajny, publiczny, nil
	case shared.TerminalKeyTypeEcdsa:
		tajny, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, bladWykonaniaTerminala("nie można wytworzyć klucza ECDSA: " + err.Error())
		}
		publiczny, err := ssh.NewPublicKey(&tajny.PublicKey)
		if err != nil {
			return nil, nil, bladWykonaniaTerminala(err.Error())
		}
		return tajny, publiczny, nil
	default:
		return nil, nil, bladZadaniaTerminala("rodzaj klucza " + string(rodzaj) +
			" nie należy do słownika kontraktu")
	}
}

// zapiszPlikiKlucza odkłada klucz prywatny i publiczny na dysk maszyny rdzenia
// i oddaje ścieżkę klucza prywatnego.
//
// Prawa pliku prywatnego to 0600 — wyłącznie konto procesu rdzenia. `ssh` sam
// odmawia użycia klucza o prawach szerszych, więc plik zapisany inaczej byłby
// plikiem, którym nie da się połączyć.
func (a *adapterTerminala) zapiszPlikiKlucza(kod string, prywatny, publiczny []byte) (string, error) {
	katalog := filepath.Join(a.katalogDanych, katalogKluczy)
	if strings.TrimSpace(a.katalogDanych) == "" {
		return "", bladWykonaniaTerminala("rdzeń nie zna własnego katalogu danych, " +
			"więc nie ma gdzie położyć wytworzonego klucza")
	}
	if err := os.MkdirAll(katalog, 0o700); err != nil {
		return "", bladWykonaniaTerminala("nie można założyć katalogu kluczy: " + err.Error())
	}
	sciezka := filepath.Join(katalog, kod)
	if err := os.WriteFile(sciezka, prywatny, 0o600); err != nil {
		return "", bladWykonaniaTerminala("nie można zapisać klucza prywatnego: " + err.Error())
	}
	if err := os.WriteFile(sciezka+".pub", publiczny, 0o644); err != nil {
		_ = os.Remove(sciezka)
		return "", bladWykonaniaTerminala("nie można zapisać klucza publicznego: " + err.Error())
	}
	return sciezka, nil
}

// trescKluczaJawnego składa wiersz klucza publicznego w postaci `authorized_keys`.
func trescKluczaJawnego(publiczny ssh.PublicKey, komentarz string) string {
	if publiczny == nil {
		return ""
	}
	wiersz := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(publiczny)))
	if komentarz = strings.TrimSpace(komentarz); komentarz != "" {
		wiersz += " " + komentarz
	}
	return wiersz
}

// rodzajKluczaZTypu przekłada nazwę algorytmu OpenSSH na słownik kontraktu.
// Nazwa nierozpoznana schodzi na `ed25519` — wykaz kontraktu ma trzy pozycje
// i nie ma w nim wartości „inny”, a wpis ma powstać mimo to.
func rodzajKluczaZTypu(typ string) shared.TerminalKeyType {
	switch {
	case strings.HasPrefix(typ, "ssh-rsa"), strings.HasPrefix(typ, "rsa-"):
		return shared.TerminalKeyTypeRsa
	case strings.HasPrefix(typ, "ecdsa-"):
		return shared.TerminalKeyTypeEcdsa
	default:
		return shared.TerminalKeyTypeEd25519
	}
}

// kluczKontraktu przekłada wiersz wykazu kluczy na byt kontraktu.
func kluczKontraktu(w dane.KluczTerminala) shared.TerminalSshKey {
	return shared.TerminalSshKey{
		Id:            w.Kod,
		Name:          w.Nazwa,
		KeyType:       w.Rodzaj,
		Fingerprint:   w.Odcisk,
		PublicKey:     wskaznikTekstu(w.KluczJawny),
		Path:          wskaznikTekstu(w.Sciezka),
		HasPassphrase: w.Haslo,
		CreatedAt:     chwilaBazy(w.Utworzono),
	}
}
