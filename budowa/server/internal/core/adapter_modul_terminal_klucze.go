// Plik obsługuje komendy terminal.key.generate, terminal.key.import,
// terminal.key.list i terminal.key.remove — wykaz kluczy SSH znanych
// rdzeniowi, wytwarzanych biblioteką standardową bez pomocy ssh-keygen.
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

// przedrostekKlucza znakuje identyfikator wpisu wykazu kluczy, żeby kod klucza
// dało się odróżnić od kodów innych bytów terminala na pierwszy rzut oka.
const przedrostekKlucza = "tkey-"

// katalogKluczy jest podkatalogiem katalogu danych rdzenia, w którym leżą klucze
// wytworzone przez moduł — osobnym od katalogu SSH konta procesu, bo klucz
// produktu ma być usuwalny wraz z produktem.
const katalogKluczy = "klucze-ssh"

// dlugoscKluczaRSA jest długością klucza RSA w bitach. 3072 to dolna granica
// zalecana dziś dla kluczy RSA; 2048 nie jest już wyborem, który wolno zapisać
// w produkcie na następne lata.
const dlugoscKluczaRSA = 3072

// WytworzKlucz obsługuje komendę terminal.key.generate: wytwarza parę kluczy
// wskazanego rodzaju, zapisuje ją na dysku maszyny rdzenia i wciąga wpis do
// wykazu.
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
			"hasło klucza wskazuje sejf, a serwer nie ma jego czytnika; klucz z hasłem " +
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
		// Wiersza nie ma, więc plików też nie ma prawa zostać — byłyby materiałem,
		// o którym nikt już nie wie.
		_ = os.Remove(sciezka)
		_ = os.Remove(sciezka + ".pub")
		return shared.TerminalKeyGenerateResponse{}, err
	}
	return shared.TerminalKeyGenerateResponse{Key: kluczKontraktu(wiersz)}, nil
}

// WciagnijKlucz obsługuje komendę terminal.key.import: czyta klucz wskazany
// ścieżką, ustala jego rodzaj i odcisk, i wciąga wpis o nim do wykazu.
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
			"wciągnięcie klucza wymaga jego nazwy i ścieżki na maszynie serwera")
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
		// Klucz chroniony hasłem czyta się tylko z hasłem, którego rdzeń nie ma —
		// to nie jest powód odmowy.
		wiersz.Haslo = true
		wiersz.Rodzaj = shared.TerminalKeyTypeEd25519
		if chroniony.PublicKey != nil {
			wiersz.Odcisk = ssh.FingerprintSHA256(chroniony.PublicKey)
			wiersz.Rodzaj = rodzajKluczaZTypu(chroniony.PublicKey.Type())
			wiersz.KluczJawny = trescKluczaJawnego(chroniony.PublicKey, "")
		}
	case err != nil:
		return shared.TerminalKeyImportResponse{}, bladZadaniaTerminala(
			"plik " + sciezka + " nie jest kluczem prywatnym w postaci czytelnej dla serwera: " +
				err.Error())
	default:
		podpisujacy, err := ssh.NewSignerFromKey(prywatny)
		if err != nil {
			return shared.TerminalKeyImportResponse{}, bladZadaniaTerminala(
				"klucz " + sciezka + " ma rodzaj, którego serwer nie obsługuje: " + err.Error())
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

// WykazKluczy obsługuje komendę terminal.key.list: zwraca wykaz kluczy SSH
// zapisanych w dzienniku wyposażenia terminala.
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

// UsunKlucz obsługuje komendę terminal.key.remove: odpina wpisy książki hostów
// wskazujące klucz, zdejmuje go z wykazu i opcjonalnie usuwa jego pliki
// z dysku.
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
		// Niepowodzenie usunięcia pliku nie wycofuje zdjęcia z wykazu: wpis już nie
		// istnieje.
		_ = os.Remove(wiersz.Sciezka)
		_ = os.Remove(wiersz.Sciezka + ".pub")
	}
	odpowiedz := shared.TerminalKeyRemoveResponse{Removed: zdjety}
	if len(odpiete) > 0 {
		odpowiedz.DetachedHostIds = odpiete
	}
	return odpowiedz, nil
}

// paraKluczy wytwarza materiał klucza wskazanego rodzaju biblioteką
// standardową, bez wołania zewnętrznego programu ssh-keygen.
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
// pod prawami 0600, jedynymi, przy których ssh godzi się użyć klucza, i oddaje
// ścieżkę klucza prywatnego.
func (a *adapterTerminala) zapiszPlikiKlucza(kod string, prywatny, publiczny []byte) (string, error) {
	katalog := filepath.Join(a.katalogDanych, katalogKluczy)
	if strings.TrimSpace(a.katalogDanych) == "" {
		return "", bladWykonaniaTerminala("serwer nie zna własnego katalogu danych, " +
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

// trescKluczaJawnego składa wiersz klucza publicznego w postaci pliku
// authorized_keys, dołączając komentarz, gdy został podany.
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

// kluczKontraktu przekłada wiersz wykazu kluczy z magazynu na byt kontraktu
// TerminalSshKey do wydania w odpowiedzi.
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
