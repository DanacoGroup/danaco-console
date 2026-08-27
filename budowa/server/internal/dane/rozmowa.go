// Odpowiedzialność pliku: utrwalenie rozmowy okna komunikacji, zapis wiadomości i odpowiedzi modelu oraz odczyt historii okna po restarcie.
package dane

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"danacoconsole/shared"
)

// OpisOkna niesie parametry okna komunikacji potrzebne do założenia jego wiersza.
// Pola odpowiadają ustawieniom okna: własny moduł, własny kanał modelu, własna
// lista katalogów roboczych.
type OpisOkna struct {
	Id                  string
	IdSesji             string
	TytulSesji          string
	Projekt             string
	Modul               string
	KanalModelu         string
	Tytul               string
	KatalogiRobocze     []string
	SrodowiskoWykonania shared.ExecutionEnv
	TrybUprawnien       shared.PermissionMode
	RolaOkna            shared.WindowRole
	// Agent jest kodem eksperta nałożonego na kanał modelu okna. Pusty znaczy
	// model surowy.
	Agent string
}

// zrodloOpisuOkna podaje opis okna o wskazanym identyfikatorze rdzenia. Fałsz
// oznacza okno nieznane rejestrowi — wtedy wiersza nie da się założyć.
type zrodloOpisuOkna func(idOkna string) (OpisOkna, bool)

// UtrwalaczRozmowy zapisuje i odczytuje historię okien komunikacji. Pamięta
// odwzorowanie identyfikatora okna na wiersz, żeby nie pytać bazy przy każdej
// wiadomości.
type UtrwalaczRozmowy struct {
	zestaw *Zestaw
	zrodlo zrodloOpisuOkna

	mu   sync.Mutex
	okna map[string]int64
}

// NowyUtrwalaczRozmowy wiąże utrwalacz z zestawem repozytoriów i rejestrem okien
// rdzenia. Brak źródła opisu nie jest błędem konstrukcji — utrwalacz odnajdzie
// wtedy wyłącznie okna już zapisane w bazie.
func NowyUtrwalaczRozmowy(zestaw *Zestaw, zrodlo zrodloOpisuOkna) *UtrwalaczRozmowy {
	return &UtrwalaczRozmowy{zestaw: zestaw, zrodlo: zrodlo, okna: map[string]int64{}}
}

// Zapisz dopisuje wiadomość na koniec historii okna, nadając jej kolejny numer porządkowy w tej rozmowie.
func (u *UtrwalaczRozmowy) Zapisz(ctx context.Context, w shared.Message) error {
	oknoID, err := u.wierszOkna(ctx, w.WindowId)
	if err != nil {
		return err
	}
	wiersz, err := u.wierszWiadomosci(ctx, w, oknoID)
	if err != nil {
		return err
	}
	if _, err := u.zestaw.Wiadomosci.Dopisz(ctx, wiersz); err != nil {
		return err
	}
	return nil
}

// Zmien zapisuje treść i stan wiadomości już utrwalonej — tak domyka się
// odpowiedź modelu po zakończeniu strumienia.
func (u *UtrwalaczRozmowy) Zmien(ctx context.Context, w shared.Message) error {
	wiersz, err := u.zestaw.Wiadomosci.PoIdentyfikatorze(ctx, w.Id)
	if err != nil {
		return err
	}
	tresc := w.Content
	wiersz.Tresc = &tresc
	wiersz.Stan = w.Status
	// Zużycie tokenów domyka się w podsumowaniu ostatniej wiadomości; zero nie wymazuje wartości.
	meta := odczytajMetadaneWiadomosci(w.Metadata)
	if meta.TokenyWejscia > 0 {
		wiersz.TokenyWejscia = meta.TokenyWejscia
	}
	if meta.TokenyWyjscia > 0 {
		wiersz.TokenyWyjscia = meta.TokenyWyjscia
	}
	return u.zestaw.Wiadomosci.ZapiszWynik(ctx, wiersz)
}

// Historia zwraca całą zapisaną rozmowę okna, od najstarszej wiadomości do najnowszej, w oryginalnej kolejności.
func (u *UtrwalaczRozmowy) Historia(ctx context.Context, idOkna string) ([]shared.Message, error) {
	oknoID, err := u.wierszOkna(ctx, idOkna)
	if err != nil {
		return nil, err
	}
	idSesji, err := u.identyfikatorSesjiOkna(ctx, oknoID)
	if err != nil {
		return nil, err
	}
	wiersze, err := u.zestaw.Wiadomosci.ListaOkna(ctx, oknoID, 0)
	if err != nil {
		return nil, err
	}
	wykaz := make([]shared.Message, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, u.wiadomoscKontraktu(ctx, wiersz, idOkna, idSesji))
	}
	// Bloki nietekstowe tury wracają w metadanych wiadomości, żeby przeżyły restart jak sam tekst.
	return u.doklejBloki(ctx, idOkna, wykaz), nil
}

// identyfikatorSesjiOkna zwraca identyfikator rdzenia sesji, do której należy
// okno. Wiersz sesji założony poza rdzeniem go nie ma — wtedy wraca pusty napis,
// bo brak identyfikatora nie unieważnia historii.
func (u *UtrwalaczRozmowy) identyfikatorSesjiOkna(ctx context.Context, oknoID int64) (string, error) {
	okno, err := u.zestaw.Okna.Pobierz(ctx, oknoID)
	if err != nil {
		return "", err
	}
	sesja, err := u.zestaw.Sesje.Pobierz(ctx, okno.SesjaID)
	if err != nil {
		return "", err
	}
	if sesja.IdentyfikatorZewnetrzny == nil {
		return "", nil
	}
	return *sesja.IdentyfikatorZewnetrzny, nil
}

// metadaneWiadomosci to część obszaru Metadata kontraktu, której typ Message nie modeluje osobnymi polami: atrybucja i zużycie tokenów.
type metadaneWiadomosci struct {
	Persona        string `json:"persona,omitempty"`
	OknoZrodloweId string `json:"sourceWindowId,omitempty"`
	TokenyWejscia  int    `json:"inputTokens,omitempty"`
	TokenyWyjscia  int    `json:"outputTokens,omitempty"`
}

// odczytajMetadaneWiadomosci wydobywa pola atrybucji i zużycia z obszaru Metadata.
// Pusty albo nieczytelny obszar daje wartość zerową — brak metadanych nie jest
// błędem i nie przerywa zapisu.
func odczytajMetadaneWiadomosci(surowe json.RawMessage) metadaneWiadomosci {
	var meta metadaneWiadomosci
	if len(surowe) == 0 {
		return meta
	}
	_ = json.Unmarshal(surowe, &meta)
	return meta
}

// wierszWiadomosci przekłada wiadomość kontraktu na wiersz tabeli: treść, atrybucję i zużycie niesione w metadanych, oraz wykaz załączników.
func (u *UtrwalaczRozmowy) wierszWiadomosci(ctx context.Context, w shared.Message, oknoID int64) (Wiadomosc, error) {
	if w.Id == "" {
		return Wiadomosc{}, fmt.Errorf("dane: wiadomość bez identyfikatora nie da się utrwalić")
	}
	tresc, identyfikator := w.Content, w.Id
	wiersz := Wiadomosc{
		OknoID:                  oknoID,
		Rola:                    w.Role,
		RodzajTresci:            shared.ChunkKindText,
		Stan:                    w.Status,
		Tresc:                   &tresc,
		IdentyfikatorZewnetrzny: &identyfikator,
	}
	wiersz.Zalaczniki = zalacznikiDoKolumny(w.Attachments)
	meta := odczytajMetadaneWiadomosci(w.Metadata)
	if meta.Persona != "" {
		persona := meta.Persona
		wiersz.Persona = &persona
	}
	wiersz.TokenyWejscia = meta.TokenyWejscia
	wiersz.TokenyWyjscia = meta.TokenyWyjscia
	if meta.OknoZrodloweId != "" {
		// Nierozpoznane okno źródłowe zostawia atrybucję pustą; sama wiadomość trafia do bazy bez przypisania.
		if zrodloID, err := u.wierszOkna(ctx, meta.OknoZrodloweId); err == nil {
			wiersz.OknoZrodloweID = &zrodloID
		}
	}
	return wiersz, nil
}

// wiadomoscKontraktu przekłada wiersz tabeli na wiadomość kontraktu wraz
// z metadanymi atrybucji i zużycia odtworzonymi z kolumn.
func (u *UtrwalaczRozmowy) wiadomoscKontraktu(ctx context.Context, wiersz Wiadomosc, idOkna, idSesji string) shared.Message {
	w := shared.Message{
		Id:        strconv.FormatInt(wiersz.ID, 10),
		WindowId:  idOkna,
		SessionId: idSesji,
		Role:      wiersz.Rola,
		Status:    wiersz.Stan,
		CreatedAt: czasUtworzenia(wiersz.Utworzono),
	}
	if wiersz.IdentyfikatorZewnetrzny != nil {
		w.Id = *wiersz.IdentyfikatorZewnetrzny
	}
	if wiersz.Tresc != nil {
		w.Content = *wiersz.Tresc
	}
	w.Attachments = zalacznikiZKolumny(wiersz.Zalaczniki)
	w.Metadata = u.metadaneKontraktu(ctx, wiersz)
	return w
}

// zalacznikiDoKolumny składa wykaz odwołań w kolumnę JSON; wykaz pusty daje wartość pustą, nie tablicę.
func zalacznikiDoKolumny(odwolania []string) *string {
	if len(odwolania) == 0 {
		return nil
	}
	surowe, err := json.Marshal(odwolania)
	if err != nil {
		return nil
	}
	zapis := string(surowe)
	return &zapis
}

// zalacznikiZKolumny odtwarza wykaz odwołań z kolumny. Kolumna pusta albo
// nieczytelna daje wykaz pusty — historia ma się wyświetlić, a nie zniknąć przez
// jedno uszkodzone pole — ta sama zasada co przy czasie utworzenia.
func zalacznikiZKolumny(zapis *string) []string {
	if zapis == nil || *zapis == "" {
		return nil
	}
	var odwolania []string
	if err := json.Unmarshal([]byte(*zapis), &odwolania); err != nil {
		return nil
	}
	return odwolania
}

// metadaneKontraktu składa obszar Metadata z kolumn atrybucji i zużycia. Wiersz
// bez żadnej z tych wartości daje pusty obszar (nil) — kontrakt pomija go wtedy
// w wyjściu, zamiast nieść pusty obiekt.
func (u *UtrwalaczRozmowy) metadaneKontraktu(ctx context.Context, wiersz Wiadomosc) json.RawMessage {
	meta := metadaneWiadomosci{
		TokenyWejscia: wiersz.TokenyWejscia,
		TokenyWyjscia: wiersz.TokenyWyjscia,
	}
	if wiersz.Persona != nil {
		meta.Persona = *wiersz.Persona
	}
	if wiersz.OknoZrodloweID != nil {
		meta.OknoZrodloweId = u.identyfikatorOknaZrodlowego(ctx, *wiersz.OknoZrodloweID)
	}
	if meta == (metadaneWiadomosci{}) {
		return nil
	}
	surowe, err := json.Marshal(meta)
	if err != nil {
		return nil
	}
	return surowe
}

// identyfikatorOknaZrodlowego przekłada klucz wiersza okna źródłowego na identyfikator rdzenia znany warstwie wyższej.
func (u *UtrwalaczRozmowy) identyfikatorOknaZrodlowego(ctx context.Context, oknoID int64) string {
	okno, err := u.zestaw.Okna.Pobierz(ctx, oknoID)
	if err != nil || okno.IdentyfikatorZewnetrzny == nil {
		return ""
	}
	return *okno.IdentyfikatorZewnetrzny
}

// formatCzasuBazy odpowiada wyrażeniu strftime schematu: czas UTC w ISO 8601
// z częścią ułamkową sekundy.
const formatCzasuBazy = "2006-01-02T15:04:05.999Z"

// czasUtworzenia przekłada znacznik czasu bazy na milisekundy epoki kontraktu.
// Znacznik nieczytelny daje zero — historia ma się wyświetlić, a nie zniknąć.
func czasUtworzenia(znacznik string) int64 {
	chwila, err := time.Parse(formatCzasuBazy, znacznik)
	if err != nil {
		return 0
	}
	return chwila.UnixMilli()
}
