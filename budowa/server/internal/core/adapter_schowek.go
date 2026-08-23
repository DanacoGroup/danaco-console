// Odpowiedzialność pliku: rodzina `clipboard.*` — trwała historia schowka
// Operatora.
//
// ── Czyj jest schowek i czego rdzeń NIE robi ────────────────────────────────
// Schowek należy do maszyny Operatora. Rdzeń stoi na serwerze i schowka tej
// maszyny nie widzi — ani go nie czyta, ani do niego nie pisze. Ta rodzina nie
// udaje inaczej: nie ma tu ani jednej ścieżki, która sięgałaby po cudzy
// schowek, i nie ma komendy „wklej", bo wklejenie jest czynnością okna.
//
// Podział ról jest taki:
//   - Operator kopiuje u siebie; okno oddaje skopiowaną treść rdzeniowi
//     komendą `clipboard.push`. Rdzeń dostaje treść, a nie dostęp.
//   - Rdzeń daje tej treści trwałość: historia przestaje ginąć razem z kartą
//     i jest ta sama na każdej maszynie tego samego Operatora.
//   - Operator wybiera wpis z wykazu (`clipboard.list`); okno wstawia go
//     u siebie — do pola, do schowka systemowego, gdziekolwiek.
//
// Dzięki temu podziałowi komenda robi dokładnie to, co obiecuje jej nazwa,
// i ani kroku więcej.
//
// ── Powtórzenie nie mnoży wpisów ────────────────────────────────────────────
// Ta sama treść skopiowana drugi raz podnosi wpis zastany na czoło wykazu.
// Rozstrzyga to warunek UNIQUE na odcisku treści, a nie odczyt-i-zapis
// w adapterze: dwa okna kopiujące naraz rozjechałyby się na odczycie.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	przedrostekWpisuSchowka = "schowek-"

	// dlugoscZajawkiSchowka domyka podgląd wpisu pokazywany w wykazie. Wpis
	// bywa całym dokumentem, a lista ma się otworzyć od razu — pełna treść
	// każdego wpisu w wykazie zamieniłaby odczyt historii w pobranie archiwum.
	dlugoscZajawkiSchowka = 200
)

// adapterSchowka wypełnia port `Schowek`.
type adapterSchowka struct {
	repozytorium dane.RepozytoriumSchowka
}

// nowyAdapterSchowka wiąże port z historią schowka.
func nowyAdapterSchowka(repozytorium dane.RepozytoriumSchowka) *adapterSchowka {
	return &adapterSchowka{repozytorium: repozytorium}
}

// WykazSchowka obsługuje `clipboard.list`.
func (a *adapterSchowka) WykazSchowka(ctx context.Context,
	z shared.ClipboardListRequest) (shared.ClipboardListResponse, error) {

	if a.repozytorium == nil {
		return shared.ClipboardListResponse{}, bladZapleczaSchowka()
	}
	sito := dane.SitoSchowka{
		Fraza:        wartoscTekstu(z.Query),
		TylkoPrzypie: z.PinnedOnly != nil && *z.PinnedOnly,
		Granica:      wartoscLiczby(z.Limit),
		Przesuniecie: wartoscLiczby(z.Offset),
	}
	if z.Kind != nil {
		sito.Rodzaj = string(*z.Kind)
	}
	wpisy, wszystkich, err := a.repozytorium.WpisySchowka(ctx, sito)
	if err != nil {
		return shared.ClipboardListResponse{}, bladMagazynuSchowka(err)
	}
	wykaz := make([]shared.ClipboardEntry, 0, len(wpisy))
	for _, wpis := range wpisy {
		wykaz = append(wykaz, wpisSchowkaKontraktu(wpis))
	}
	return shared.ClipboardListResponse{Entries: wykaz, Total: wszystkich}, nil
}

// DopiszDoSchowka obsługuje `clipboard.push`.
func (a *adapterSchowka) DopiszDoSchowka(ctx context.Context,
	z shared.ClipboardPushRequest) (shared.ClipboardPushResponse, error) {

	if a.repozytorium == nil {
		return shared.ClipboardPushResponse{}, bladZapleczaSchowka()
	}
	if z.Content == "" {
		return shared.ClipboardPushResponse{}, bladWskazaniaSchowka(
			"wpis bez treści — nie ma czego zapamiętać")
	}
	rodzaj := shared.ClipboardEntryKindText
	if z.Kind != nil && strings.TrimSpace(string(*z.Kind)) != "" {
		rodzaj = string(*z.Kind)
	}
	if err := sprawdzRodzajWpisuSchowka(shared.ClipboardEntryKind(rodzaj)); err != nil {
		return shared.ClipboardPushResponse{}, err
	}

	// Odcisk liczymy z rodzaju i treści razem: ta sama treść wniesiona raz jako
	// tekst, a raz jako plik, jest dwoma różnymi wpisami — inaczej wkleja się
	// jedno i drugie.
	suma := sha256.Sum256([]byte(rodzaj + "\x00" + z.Content))
	zajawka := zajawkaSchowka(z.Content, shared.ClipboardEntryKind(rodzaj))

	wpis := dane.WpisSchowka{
		Kod:           nowyIdentyfikator(przedrostekWpisuSchowka),
		Rodzaj:        rodzaj,
		Tresc:         z.Content,
		Odcisk:        hex.EncodeToString(suma[:]),
		Zajawka:       &zajawka,
		RozmiarBajtow: int64(len(z.Content)),
		Wrazliwy:      z.Sensitive != nil && *z.Sensitive,
		OknoZrodlowe:  z.SourceWindowId,
		Utworzono:     time.Now().UnixMilli(),
	}
	zapisany, bylo, err := a.repozytorium.DopiszWpisSchowka(ctx, wpis)
	if err != nil {
		return shared.ClipboardPushResponse{}, bladMagazynuSchowka(err)
	}
	return shared.ClipboardPushResponse{
		Entry: wpisSchowkaKontraktu(zapisany), AlreadyPresent: bylo,
	}, nil
}

// PrzypnijWpisSchowka obsługuje `clipboard.pin`.
func (a *adapterSchowka) PrzypnijWpisSchowka(ctx context.Context,
	z shared.ClipboardPinRequest) (shared.ClipboardPinResponse, error) {

	if a.repozytorium == nil {
		return shared.ClipboardPinResponse{}, bladZapleczaSchowka()
	}
	wpis, err := a.repozytorium.PrzypnijWpisSchowka(ctx, strings.TrimSpace(z.EntryId), z.Pinned)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.ClipboardPinResponse{}, bladNieznanegoWpisuSchowka(z.EntryId)
		}
		return shared.ClipboardPinResponse{}, bladMagazynuSchowka(err)
	}
	return shared.ClipboardPinResponse{Entry: wpisSchowkaKontraktu(wpis)}, nil
}

// UsunZeSchowka obsługuje `clipboard.delete`.
func (a *adapterSchowka) UsunZeSchowka(ctx context.Context,
	z shared.ClipboardDeleteRequest) (shared.ClipboardDeleteResponse, error) {

	if a.repozytorium == nil {
		return shared.ClipboardDeleteResponse{}, bladZapleczaSchowka()
	}
	zPrzypietymi := z.IncludePinned != nil && *z.IncludePinned
	usunietych, err := a.repozytorium.UsunWpisySchowka(ctx,
		strings.TrimSpace(wartoscTekstu(z.EntryId)), zPrzypietymi)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.ClipboardDeleteResponse{}, bladNieznanegoWpisuSchowka(wartoscTekstu(z.EntryId))
		}
		return shared.ClipboardDeleteResponse{}, bladMagazynuSchowka(err)
	}
	return shared.ClipboardDeleteResponse{Deleted: usunietych}, nil
}

// zajawkaSchowka składa podgląd wpisu.
//
// Wpis nietekstowy niesie w treści bajty w postaci base64, więc jego zajawką
// jest opis, a nie początek zapisu: pierwsze dwieście znaków base64 nie mówi
// czytelnikowi nic poza tym, że to base64.
func zajawkaSchowka(tresc string, rodzaj shared.ClipboardEntryKind) string {
	if rodzaj != shared.ClipboardEntryKindText {
		return string(rodzaj) + ", " + zapisRozmiaruSchowka(len(tresc))
	}
	skrocona := strings.TrimSpace(tresc)
	if utf8.RuneCountInString(skrocona) <= dlugoscZajawkiSchowka {
		return skrocona
	}
	runy := []rune(skrocona)
	return string(runy[:dlugoscZajawkiSchowka]) + "…"
}

// zapisRozmiaruSchowka podaje wielkość wpisu w postaci czytelnej dla człowieka.
func zapisRozmiaruSchowka(bajtow int) string {
	switch {
	case bajtow >= 1024*1024:
		return zapisLiczbyMiary(float64(bajtow)/(1024*1024)) + " MB"
	case bajtow >= 1024:
		return zapisLiczbyMiary(float64(bajtow)/1024) + " kB"
	default:
		return zapisLiczbyMiary(float64(bajtow)) + " B"
	}
}

// sprawdzRodzajWpisuSchowka odbija rodzaj spoza wyliczenia kontraktu.
func sprawdzRodzajWpisuSchowka(rodzaj shared.ClipboardEntryKind) error {
	switch rodzaj {
	case shared.ClipboardEntryKindText, shared.ClipboardEntryKindImage,
		shared.ClipboardEntryKindFile:
		return nil
	default:
		return bladWskazaniaSchowka("nie znam rodzaju wpisu „" + string(rodzaj) +
			"” — rdzeń zna: text, image, file")
	}
}

// wpisSchowkaKontraktu przekłada wiersz na wpis kontraktu.
func wpisSchowkaKontraktu(w dane.WpisSchowka) shared.ClipboardEntry {
	return shared.ClipboardEntry{
		Id: w.Kod, Kind: shared.ClipboardEntryKind(w.Rodzaj), Content: w.Tresc,
		Preview: w.Zajawka, SizeBytes: int(w.RozmiarBajtow), Pinned: w.Przypiety,
		Sensitive: w.Wrazliwy, SourceWindowId: w.OknoZrodlowe,
		CreatedAt: w.Utworzono, UsedAt: w.Uzyto,
	}
}

// bladZapleczaSchowka nazywa brak historii schowka po stronie rdzenia.
func bladZapleczaSchowka() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"schowek: rdzeń nie ma wpiętej historii schowka — naprawa: podpiąć "+
			"repozytorium schowka przy składaniu rdzenia"))
}

// bladWskazaniaSchowka nazywa niepoprawne żądanie.
func bladWskazaniaSchowka(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "schowek: "+powod))
}

// bladNieznanegoWpisuSchowka nazywa wskazanie wpisu, którego nie ma.
func bladNieznanegoWpisuSchowka(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"schowek: nie ma wpisu o identyfikatorze "+strings.TrimSpace(kod)))
}

// bladMagazynuSchowka nazywa niepowodzenie zapisu albo odczytu.
func bladMagazynuSchowka(err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"schowek: "+err.Error()))
}
