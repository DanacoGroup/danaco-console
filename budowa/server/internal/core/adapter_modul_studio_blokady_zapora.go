// Plik obsługuje zaporę blokad — sprawdzenie blokad fragmentów wpięte
// w rejestr komend na drodze każdej komendy zmieniającej dokument Studia —
// oraz siatkę dopisującą ślad autora, gdy obsługiwacz go nie odłożył.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// zaporaZadanieZmiany jest tym, co zapora czyta z ładunku KAŻDEJ komendy:
// dokument, zakres (gdy jest) i podpis wykonawcy. Reszty ładunku nie zna
// i znać nie musi.
type zaporaZadanieZmiany struct {
	DocumentId string `json:"documentId"`
	RangeStart *int   `json:"rangeStart,omitempty"`
	RangeEnd   *int   `json:"rangeEnd,omitempty"`
	// Selection* to druga nazwa tego samego pojęcia w rodzinie komentarzy.
	SelectionStart *int `json:"selectionStart,omitempty"`
	SelectionEnd   *int `json:"selectionEnd,omitempty"`
	// Content niosą komendy przepisujące treść w całości; zapora uzgadnia je
	// przed wykonaniem.
	Content *string `json:"content,omitempty"`

	Author     *shared.StudioAuthor `json:"author,omitempty"`
	AgentId    *string              `json:"agentId,omitempty"`
	AgentName  *string              `json:"agentName,omitempty"`
	SubagentId *string              `json:"subagentId,omitempty"`
}

// zakres oddaje zakres żądania pod obiema nazwami kontraktu: rangeStart
// i rangeEnd albo selectionStart i selectionEnd.
func (z zaporaZadanieZmiany) zakres() (int, int, bool) {
	if z.RangeStart != nil && z.RangeEnd != nil {
		return *z.RangeStart, *z.RangeEnd, true
	}
	if z.SelectionStart != nil && z.SelectionEnd != nil {
		return *z.SelectionStart, *z.SelectionEnd, true
	}
	return 0, 0, false
}

// podpis oddaje pola tożsamości wykonawcy potrzebne rozpoznaniu, kto zawołał
// daną komendę zmieniającą.
func (z zaporaZadanieZmiany) podpis() kontrolaPodpisZadania {
	return kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	}
}

// zaporaKomendyBezZmianyTresci wymienia komendy rodziny studio.*, które
// dokumentu nie zmieniają, więc zapora ich nie dotyczy. Wykaz jest wykazem
// wyjątków, nie wykazem objętych, żeby komenda nowa była domyślnie sprawdzana.
var zaporaKomendyBezZmianyTresci = map[shared.MessageType]bool{
	shared.CommandStudioDocumentOpen:          true,
	shared.CommandStudioRepositoryList:        true,
	shared.CommandStudioDiffCompare:           true,
	shared.CommandStudioDiffVisual:            true,
	shared.CommandStudioDiffSource:            true,
	shared.CommandStudioDiffReportExport:      true,
	shared.CommandStudioPreviewRender:         true,
	shared.CommandStudioSearchSemantic:        true,
	shared.CommandStudioTrackingList:          true,
	shared.CommandStudioCommentList:           true,
	shared.CommandStudioAnnotationList:        true,
	shared.CommandStudioTemplateList:          true,
	shared.CommandStudioOperationList:         true,
	shared.CommandStudioChainList:             true,
	shared.CommandStudioBranchList:            true,
	shared.CommandStudioExportProfileList:     true,
	shared.CommandStudioRepositoryExport:      true,
	shared.CommandStudioPackageExport:         true,
	shared.CommandStudioIngestQueueList:       true,
	shared.CommandStudioIngestDeviceList:      true,
	shared.CommandStudioLockAdd:               true,
	shared.CommandStudioLockRemove:            true,
	shared.CommandStudioLockList:              true,
	shared.CommandStudioJournalList:           true,
	shared.CommandStudioJournalRevert:         true,
	shared.CommandStudioJournalRedo:           true,
	shared.CommandStudioModelChangesList:      true,
	shared.CommandStudioModelChangesRevert:    true,
	shared.CommandStudioModelChangesNavigate:  true,
	shared.CommandStudioBackupCreate:          true,
	shared.CommandStudioBackupList:            true,
	shared.CommandStudioBackupRestore:         true,
	shared.CommandStudioAutosaveGet:           true,
	shared.CommandStudioAutosaveSet:           true,
	shared.CommandStudioVersionSeriesList:     true,
	shared.CommandStudioVersionRestoreInitial: true,
	shared.CommandStudioMarkupList:            true,
	shared.CommandStudioMarkupTypeList:        true,
	shared.CommandStudioMarkupTypeSave:        true,
	shared.CommandStudioMarkupTypeDelete:      true,
	shared.CommandStudioAgentsClaim:           true,
	shared.CommandStudioAgentsRelease:         true,
	shared.CommandStudioAgentsSlotsList:       true,
	shared.CommandStudioAgentsConflictsList:   true,
	shared.CommandStudioAgentsSettingsGet:     true,
	shared.CommandStudioAgentsSettingsSet:     true,
	shared.CommandStudioDiffFormCompare:       true,

	// Czynności, które sprawdzają blokadę same, dokładniej niż zapora byłaby
	// w stanie.
	shared.CommandStudioMarkupAdd:      true,
	shared.CommandStudioMarkupRemove:   true,
	shared.CommandStudioClipboardCopy:  true,
	shared.CommandStudioClipboardPaste: true,
	shared.CommandStudioDiffHunkApply:  true,
	shared.CommandStudioProvenanceList: true,
	shared.CommandStudioViewGet:        true,
	shared.CommandStudioViewSet:        true,
}

// zaporaBlokadStudia owija w rejestrze każdą komendę rodziny `studio.*`, która
// może zmienić dokument. Wołane po `zarejestrujStudio` i po wpięciach odcinków
// postaci, wejścia i wydania — owija to, co w rejestrze już stoi.
func zaporaBlokadStudia(r *Rejestr, a *adapterStudia) {
	if r == nil || r.wpisy == nil || a == nil {
		return
	}
	for nazwa, obsluga := range r.wpisy {
		if !strings.HasPrefix(string(nazwa), "studio.") {
			continue
		}
		if zaporaKomendyBezZmianyTresci[nazwa] {
			continue
		}
		r.wpisy[nazwa] = a.zaporaOwin(nazwa, obsluga)
	}
}

// zaporaOwin składa obsługiwacza pilnowanego blokadami, wpinanego w miejsce
// obsługiwacza pierwotnego w rejestrze.
func (a *adapterStudia) zaporaOwin(nazwa shared.MessageType, obsluga Obsluga) Obsluga {
	return func(ctx context.Context, z protocol.Request) protocol.Odpowiedz {
		var zadanie zaporaZadanieZmiany
		// Ładunek innego kształtu nie jest usterką: komenda bez dokumentu nie
		// ma czego naruszyć.
		if err := json.Unmarshal(z.Ladunek, &zadanie); err != nil || zadanie.DocumentId == "" {
			return obsluga(ctx, z)
		}
		wykonawca := kontrolaRozpoznajWykonawce(ctx, zadanie.podpis())

		skladnica, err := a.kontrolaSkladnica()
		if err != nil {
			// Rdzeń bez tabel blokad ustępuje działaniu; brak nazywa się przy
			// pierwszej czynności blokad.
			return obsluga(ctx, z)
		}
		dokument, err := a.repozytorium.Dokument(ctx, zadanie.DocumentId)
		if err != nil {
			// Dokumentu nie ma albo odczyt się nie udał — rozstrzyga to
			// obsługiwacz własną odmową.
			return obsluga(ctx, z)
		}
		blokady, err := skladnica.BlokadyFragmentow(ctx, dokument.ID)
		if err != nil {
			return porazka(bladStudio(err))
		}
		wiazace := blokadaWiazaceDlaRak(blokady, wykonawca)
		if len(wiazace) == 0 {
			return obsluga(ctx, z)
		}

		// ── Droga pierwsza: zakres znany przed wykonaniem ────────────────────
		if od, do, jest := zadanie.zakres(); jest {
			if blokada, trafiona := blokadaNaZakresie(wiazace, od, do, wykonawca); trafiona {
				return porazka(kontrolaBladBlokady(blokada, wykonawca.nazwaWykonawcy()))
			}
			return obsluga(ctx, z)
		}

		// ── Droga druga: treść w żądaniu, uzgodnienie przed wykonaniem ───────
		if zadanie.Content != nil {
			uzgodnienie := blokadaUzgodnijTresc(wartoscTekstu(dokument.Tresc), *zadanie.Content,
				blokady, wykonawca)
			if uzgodnienie.CalkiemWBlokadzie {
				return porazka(kontrolaBladBlokady(uzgodnienie.PierwszaBlokada,
					wykonawca.nazwaWykonawcy()))
			}
			if len(uzgodnienie.Pominiete) == 0 {
				return obsluga(ctx, z)
			}
			przepisany, err := zaporaPrzepiszTresc(z.Ladunek, uzgodnienie.Tresc)
			if err != nil {
				return porazka(bladStudio(err))
			}
			z.Ladunek = przepisany
			odpowiedz := obsluga(ctx, z)
			return zaporaDopiszBilans(odpowiedz, blokadaBilans(uzgodnienie))
		}

		// ── Droga trzecia: czynność na całym dokumencie. Kopia zapasowa przed
		// wykonaniem, potem uzgodnienie.
		trescPrzed := wartoscTekstu(dokument.Tresc)
		if _, err := a.kopiaPrzedCzynnoscia(ctx, dokument,
			shared.StudioBackupReasonBeforeIrreversible); err != nil {

			return porazka(err)
		}
		odpowiedz := obsluga(ctx, z)
		if odpowiedz.Status != shared.EnvelopeStatusOk {
			return odpowiedz
		}
		po, err := a.repozytorium.Dokument(ctx, zadanie.DocumentId)
		if err != nil {
			return porazka(bladStudio(err))
		}
		uzgodnienie := blokadaUzgodnijTresc(trescPrzed, wartoscTekstu(po.Tresc),
			blokady, wykonawca)
		if len(uzgodnienie.Pominiete) == 0 {
			return odpowiedz
		}
		// Treść uzgodniona wraca do bazy: fragmenty zablokowane wracają,
		// reszta zmiany zostaje.
		po.Tresc = &uzgodnienie.Tresc
		if _, err := a.repozytorium.ZapiszDokument(ctx, po); err != nil {
			return porazka(bladStudio(err))
		}
		return zaporaDopiszBilans(odpowiedz, blokadaBilans(uzgodnienie))
	}
}

// zaporaPrzepiszTresc podmienia pole content w ładunku żądania przez mapę,
// zostawiając pozostałe pola nietknięte, bo struktura znana zaporze
// zgubiłaby pola jej nieznane.
func zaporaPrzepiszTresc(ladunek json.RawMessage, tresc string) (json.RawMessage, error) {
	pola := map[string]json.RawMessage{}
	if err := json.Unmarshal(ladunek, &pola); err != nil {
		return nil, err
	}
	zapis, err := json.Marshal(tresc)
	if err != nil {
		return nil, err
	}
	pola["content"] = zapis
	return json.Marshal(pola)
}

// zaporaDopiszBilans dopisuje bilans pominięć do odpowiedzi komendy. Bilans
// zastany nie jest nadpisywany, a scalany, bo obsługiwacz mógł już oddać
// bilans własny, o czym innym.
func zaporaDopiszBilans(odpowiedz protocol.Odpowiedz,
	bilans shared.StudioActionBalance) protocol.Odpowiedz {

	if odpowiedz.Status != shared.EnvelopeStatusOk {
		return odpowiedz
	}
	pola := map[string]json.RawMessage{}
	if len(odpowiedz.Wynik) > 0 {
		if err := json.Unmarshal(odpowiedz.Wynik, &pola); err != nil {
			// Wynik nie jest obiektem JSON — odpowiedzią jest odmowa
			// nazywająca, że bilansu nie dało się oddać.
			return porazka(kontrolaBladZaplecza(
				"czynność pominęła fragmenty zablokowane, ale jej odpowiedź nie ma " +
					"kształtu, w którym da się oddać bilans — pominięcia nie wolno przemilczeć"))
		}
	}
	if zastany, jest := pola["balance"]; jest {
		var poprzedni shared.StudioActionBalance
		if err := json.Unmarshal(zastany, &poprzedni); err == nil {
			bilans.Applied += poprzedni.Applied
			bilans.SkippedCount += poprzedni.SkippedCount
			bilans.Skipped = append(poprzedni.Skipped, bilans.Skipped...)
			if poprzedni.Note != nil && bilans.Note != nil {
				zdanie := *poprzedni.Note + " " + *bilans.Note
				bilans.Note = &zdanie
			}
		}
	}
	zapis, err := json.Marshal(bilans)
	if err != nil {
		return porazka(bladStudio(err))
	}
	pola["balance"] = zapis
	wynik, err := json.Marshal(pola)
	if err != nil {
		return porazka(bladStudio(err))
	}
	odpowiedz.Wynik = wynik
	return odpowiedz
}

// kopiaPrzedCzynnoscia zakłada kopię zapasową przed czynnością nieodwracalną:
// przyjęciem wszystkich zmian modelu, zamianą w całym dokumencie i zmianą
// formatu nośnika.
func (a *adapterStudia) kopiaPrzedCzynnoscia(ctx context.Context, dokument dane.DokumentStudia,
	powod shared.StudioBackupReason) (dane.KopiaZapasowaStudia, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return dane.KopiaZapasowaStudia{}, err
	}
	tresc := wartoscTekstu(dokument.Tresc)
	postacJSON := a.kopiaPostacJako(ctx, dokument.Kod)
	kopia, err := skladnica.ZapiszKopieZapasowa(ctx, dokument.ID, dane.KopiaZapasowaStudia{
		Kod:           nowyIdentyfikator(przedrostekKopiiStudia),
		Powod:         string(powod),
		Tresc:         &tresc,
		PostacJSON:    postacJSON,
		RozmiarBajtow: int64(len(tresc)),
		UdaloSie:      true,
	})
	if err != nil {
		return dane.KopiaZapasowaStudia{}, bladStudio(err)
	}
	return kopia, nil
}

// kopiaPostacJako zapisuje postać dokumentu jako ładunek kopii. Postaci
// nieczytelnej nie zamienia w brak — oddaje nil i kopia niesie samą treść, bo
// kopia treści jest lepsza niż brak kopii.
func (a *adapterStudia) kopiaPostacJako(ctx context.Context, kodDokumentu string) *string {
	postac, err := a.wejsciePostacDokumentu(ctx, kodDokumentu)
	if err != nil {
		return nil
	}
	zapis, err := json.Marshal(postac)
	if err != nil {
		return nil
	}
	tekst := string(zapis)
	return &tekst
}

// ── Siatka pod zapisem autora ───────────────────────────────────────────────

// sladWykonawcyStudia owija komendy Studia siatką, która dopisuje ślad autora
// tam, gdzie obsługiwacz go nie odłożył, mierząc skutek na treści, nie samo
// wywołanie.
func sladWykonawcyStudia(r *Rejestr, a *adapterStudia) {
	if r == nil || r.wpisy == nil || a == nil {
		return
	}
	for nazwa, obsluga := range r.wpisy {
		if !strings.HasPrefix(string(nazwa), "studio.") {
			continue
		}
		if zaporaKomendyBezZmianyTresci[nazwa] {
			continue
		}
		r.wpisy[nazwa] = a.sladOwin(obsluga)
	}
}

// sladOwin składa obsługiwacza z siatką śladu autora, wpinaną w miejsce
// obsługiwacza pierwotnego w rejestrze.
func (a *adapterStudia) sladOwin(obsluga Obsluga) Obsluga {
	return func(ctx context.Context, z protocol.Request) protocol.Odpowiedz {
		var zadanie zaporaZadanieZmiany
		if err := json.Unmarshal(z.Ladunek, &zadanie); err != nil || zadanie.DocumentId == "" {
			return obsluga(ctx, z)
		}
		wykonawca := kontrolaRozpoznajWykonawce(ctx, zadanie.podpis())
		if !wykonawca.czyWykonawca() {
			return obsluga(ctx, z)
		}
		przed, err := a.repozytorium.Dokument(ctx, zadanie.DocumentId)
		if err != nil {
			return obsluga(ctx, z)
		}
		trescPrzed := wartoscTekstu(przed.Tresc)
		sladowPrzed := a.sladIleZmianWykonawcy(ctx, przed.ID)

		odpowiedz := obsluga(ctx, z)
		if odpowiedz.Status != shared.EnvelopeStatusOk {
			return odpowiedz
		}

		po, err := a.repozytorium.Dokument(ctx, zadanie.DocumentId)
		if err != nil {
			return odpowiedz
		}
		trescPo := wartoscTekstu(po.Tresc)
		if trescPo == trescPrzed {
			// Treść ta sama: zmianę postaci odkłada droga postaci, siatka nie
			// ma czego tu dopisać.
			return odpowiedz
		}
		if a.sladIleZmianWykonawcy(ctx, po.ID) > sladowPrzed {
			// Obsługiwacz odłożył ślad sam — drugiego nie dokłada się.
			return odpowiedz
		}
		if err := a.sladOdlozBrakujacy(ctx, po, wykonawca, trescPrzed, trescPo); err != nil {
			// Skutek na dokumencie już zapadł, odpowiedź jest udana; odmowa
			// wraca, gdy śladu nie da się odłożyć.
			return porazka(err)
		}
		return odpowiedz
	}
}

// sladIleZmianWykonawcy liczy zmiany śledzone autora model w dokumencie,
// pomijając zmiany podpisane inaczej.
func (a *adapterStudia) sladIleZmianWykonawcy(ctx context.Context, dokumentID int64) int {
	zmiany, err := a.repozytorium.ZmianySledzone(ctx, dokumentID)
	if err != nil {
		return 0
	}
	ile := 0
	for _, zmiana := range zmiany {
		if zmiana.Autor == string(shared.StudioAuthorModel) {
			ile++
		}
	}
	return ile
}

// sladOdlozBrakujacy odkłada zmianę śledzoną i wpis dziennika dla zmiany
// treści, której obsługiwacz nie podpisał, licząc zakres z różnicy treści
// w znakach.
func (a *adapterStudia) sladOdlozBrakujacy(ctx context.Context, dokument dane.DokumentStudia,
	wykonawca kontrolaWykonawca, trescPrzed, trescPo string) error {

	od, doPrzed, doPo := sladZakresRoznicy(trescPrzed, trescPo)
	przedFragment := sladWycinek(trescPrzed, od, doPrzed)
	poFragment := sladWycinek(trescPo, od, doPo)

	rodzaj := shared.StudioChangeKindWstawienie
	if poFragment == "" {
		rodzaj = shared.StudioChangeKindUsuniecie
	}
	zmiana, err := a.repozytorium.ZapiszZmianeSledzona(ctx, dokument.ID, dane.ZmianaSledzona{
		Kod:        nowyIdentyfikator(przedrostekZmianyModeluStudia),
		Rodzaj:     string(rodzaj),
		Autor:      string(shared.StudioAuthorModel),
		ZakresOd:   int64(od),
		ZakresDo:   int64(doPrzed),
		TrescPrzed: kontrolaWskaznikTekstu(przedFragment),
		TrescPo:    kontrolaWskaznikTekstu(poFragment),
		Decyzja:    "oczekuje",
	})
	if err != nil {
		return bladStudio(err)
	}
	opis := "zmiana treści wniesiona przez " + wykonawca.nazwaWykonawcy() +
		" — ślad odłożony siatką rdzenia, bo czynność nie podpisała się sama"
	odKopia, doKopia := od, doPrzed
	if _, err := a.dziennikOdlozCzynnosc(ctx, dokument, wykonawca,
		shared.StudioActionKindTextEdit, opis, &odKopia, &doKopia, nil, nil,
		&zmiana.Kod); err != nil {

		return err
	}
	return nil
}

// sladZakresRoznicy oddaje początek różnicy oraz jej koniec po obu stronach,
// liczone w ZNAKACH — cięcie po bajtach rozcinałoby polskie litery dwubajtowe.
func sladZakresRoznicy(przed, po string) (int, int, int) {
	znakiPrzed, znakiPo := []rune(przed), []rune(po)
	od := 0
	for od < len(znakiPrzed) && od < len(znakiPo) && znakiPrzed[od] == znakiPo[od] {
		od++
	}
	sufiks := 0
	for sufiks < len(znakiPrzed)-od && sufiks < len(znakiPo)-od &&
		znakiPrzed[len(znakiPrzed)-1-sufiks] == znakiPo[len(znakiPo)-1-sufiks] {
		sufiks++
	}
	return od, len(znakiPrzed) - sufiks, len(znakiPo) - sufiks
}

// sladWycinek wycina fragment treści dokumentu w znakach; zakres poza treścią
// daje pustkę zamiast błędu.
func sladWycinek(tresc string, od, do int) string {
	znaki := []rune(tresc)
	od, do, poprawny := kontrolaZakresWTresci(znaki, od, do)
	if !poprawny {
		return ""
	}
	return string(znaki[od:do])
}
