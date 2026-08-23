// Odpowiedzialność pliku: ZAPORA BLOKAD — sprawdzenie blokad fragmentów wpięte
// w rejestr komend, na drodze KAŻDEJ komendy zmieniającej dokument Studia.
//
// ── Dlaczego w rejestrze, a nie w obsługiwaczach ─────────────────────────────
// Wymaganie Właściciela mówi: sprawdzenie stoi na drodze każdej komendy
// zmieniającej dokument, po stronie serwera, PRZED dotknięciem treści. Komend
// zmieniających dokument Studio ma dziś ponad setkę i pisze je czterech
// wykonawców naraz. Wywołanie sprawdzenia w każdym obsługiwaczu z osobna
// znaczyłoby: sto miejsc do pominięcia przez pomyłkę, a każde pominięcie to
// cicha dziura w blokadzie. Zapora wpięta w rejestr obejmuje wszystkie te
// komendy JEDNYM warunkiem — i obejmuje też te, których jeszcze nikt nie
// napisał, bo działa po nazwie rodziny, nie po wykazie obsługiwaczy.
//
// Rejestr jest jedynym miejscem, w którym rdzeń rozstrzyga „co wykonać"
// (`rejestr.go`), więc jest też jedynym miejscem, przez które przechodzi
// KAŻDE wywołanie — także wywołanie modelu, bo model woła komendy tą samą
// drogą co klient. Owinięcie wpisu rejestru jest zatem tym samym, co postawienie
// straży w drzwiach, a nie przy każdym stoliku.
//
// ── Dwie drogi sprawdzenia, bo dwa kształty komend ──────────────────────────
//
//  1. KOMENDA NA FRAGMENCIE — żądanie niesie `rangeStart` i `rangeEnd`. Zakres
//     znany PRZED wykonaniem, więc sprawdzenie jest czyste: zakres stykający się
//     z blokadą wiążącą kończy się odmową NAZWANĄ i obsługiwacz nie rusza. Tak
//     idzie przygniatająca większość czynności postaci i treści.
//
//  2. KOMENDA NA CAŁYM DOKUMENCIE — żądanie zakresu nie niesie (zamiana
//     w całym dokumencie, przyjęcie wszystkich zmian, przestawienie formatu).
//     Tu zakresu przed wykonaniem NIE MA: powstaje dopiero z rachunku
//     obsługiwacza. Odmowa całości byłaby nieproporcjonalna — Właściciel mówi to
//     wprost — więc zapora robi trzy rzeczy: zakłada kopię zapasową (czynność
//     nieodwracalna i tak jej wymaga), puszcza obsługiwacza, a potem UZGADNIA
//     wynik z blokadami: fragmenty zablokowane wracają do brzmienia zastanego,
//     a odpowiedź dostaje BILANS pominięć. Blokada nie zostaje przy tym naruszona
//     na zewnątrz: żaden inny wołający nie widzi stanu przejściowego, bo
//     uzgodnienie zamyka się w tym samym wywołaniu, a odpowiedź niesie już stan
//     uzgodniony.
//
// ── Dlaczego bilans dopisuje się do odpowiedzi tutaj ────────────────────────
// Bilans musi wyjść odpowiedzią TEJ komendy, którą Operator albo model zawołał —
// inaczej pominięcie zostałoby przemilczane, a to jest zakazane. Odpowiedzi
// większości tych komend należą do innych odcinków i ich kształtów ten plik nie
// zmienia. Dopisuje więc pole `balance` do gotowego ładunku JSON (`wynik` jest
// `json.RawMessage`), zamiast żądać od czterech wykonawców, żeby każdy dołożył
// u siebie to samo pole i pamiętał o nim w każdej nowej komendzie.
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
	// Content niosą komendy przepisujące treść w całości. Gdy jest, zapora
	// uzgadnia je PRZED wykonaniem — to jest sprawdzenie najczystsze z możliwych.
	Content *string `json:"content,omitempty"`

	Author     *shared.StudioAuthor `json:"author,omitempty"`
	AgentId    *string              `json:"agentId,omitempty"`
	AgentName  *string              `json:"agentName,omitempty"`
	SubagentId *string              `json:"subagentId,omitempty"`
}

// zakres oddaje zakres żądania pod obiema nazwami kontraktu.
func (z zaporaZadanieZmiany) zakres() (int, int, bool) {
	if z.RangeStart != nil && z.RangeEnd != nil {
		return *z.RangeStart, *z.RangeEnd, true
	}
	if z.SelectionStart != nil && z.SelectionEnd != nil {
		return *z.SelectionStart, *z.SelectionEnd, true
	}
	return 0, 0, false
}

// podpis oddaje pola tożsamości wykonawcy.
func (z zaporaZadanieZmiany) podpis() kontrolaPodpisZadania {
	return kontrolaPodpisZadania{
		Author: z.Author, AgentId: z.AgentId, AgentName: z.AgentName, SubagentId: z.SubagentId,
	}
}

// zaporaKomendyBezZmianyTresci wymienia komendy rodziny `studio.*`, które
// dokumentu NIE zmieniają, więc zapora ich nie dotyczy.
//
// Wykaz jest wykazem WYJĄTKÓW, nie wykazem objętych — i to jest zamierzone.
// Gdyby zapora obejmowała wykaz komend zmieniających, komenda dopisana przez
// innego wykonawcę i niewpisana do wykazu byłaby cichą dziurą w blokadzie.
// Tak ułożony wykaz myli się w drugą stronę: komenda nowa jest domyślnie
// sprawdzana, a najgorsze, co może z tego wyjść, to sprawdzenie zbędne przy
// odczycie — widoczne od razu i nieszkodliwe dla dokumentu.
//
// Odczyty i wykazy stoją tu wszystkie, bo odczyt nie ma czego naruszyć.
// Czynności samych blokad i dziennika stoją tu, bo są narzędziem Operatora
// NAD blokadą i sprawdzają jego prawa same, dokładniej (blokadę zdejmuje
// wyłącznie Operator — sprawdza to `ZdejmijBlokade`).
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

	// Czynności, które sprawdzają blokadę SAME, i to dokładniej niż zapora.
	// Wpis tutaj nie jest osłabieniem: obsługiwacz robi to samo sprawdzenie
	// przed dotknięciem treści, a zapora zrobiłaby je na zakresie o innym
	// znaczeniu i przez to albo za szeroko, albo za wąsko.
	//
	//   markup.add — propozycja na marginesie jest JEDYNĄ drogą wykonawcy do
	//     fragmentu pod blokadą i musi przechodzić; wyróżnienie barwą sprawdza
	//     blokadę samo (`postacOdcinkiDozwolone`).
	//   markup.remove — zdejmuje barwę z zakresu SWOJEGO znakowania, którego
	//     zapora z żądania nie zna.
	//   clipboard.copy — odczyt fragmentu, także zablokowanego; wycięcie sprawdza
	//     blokadę samo.
	//   clipboard.paste — miejsce wklejenia jedzie polem `offset`, nie zakresem.
	//   diff.hunk.apply — zakres w żądaniu dotyczy WERSJI ŹRÓDŁOWEJ, nie treści
	//     bieżącej; uzgodnienie liczy się na treści bieżącej.
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

// zaporaOwin składa obsługiwacza pilnowanego blokadami.
func (a *adapterStudia) zaporaOwin(nazwa shared.MessageType, obsluga Obsluga) Obsluga {
	return func(ctx context.Context, z protocol.Request) protocol.Odpowiedz {
		var zadanie zaporaZadanieZmiany
		// Ładunek innego kształtu NIE jest tu usterką: nie każda komenda Studia
		// niesie dokument, a te, które go nie niosą, nie mają czego naruszyć.
		// Odmowa z powodu niedopasowania kształtu zablokowałaby komendę,
		// o której zapora nie ma nic do powiedzenia.
		if err := json.Unmarshal(z.Ladunek, &zadanie); err != nil || zadanie.DocumentId == "" {
			return obsluga(ctx, z)
		}
		wykonawca := kontrolaRozpoznajWykonawce(ctx, zadanie.podpis())

		skladnica, err := a.kontrolaSkladnica()
		if err != nil {
			// Rdzeń złożony bez tabel blokad nie ma czym sprawdzić, czy zmiana
			// jest wolna. Pilnowanie ustępuje wtedy DZIAŁANIU, a nie odwrotnie:
			// zamiana pracy całego modułu na odmowę z powodu braku montażu
			// byłaby szkodą większą niż brak sprawdzenia. Brak nazywa się
			// wprost przy pierwszej czynności samych blokad.
			return obsluga(ctx, z)
		}
		dokument, err := a.repozytorium.Dokument(ctx, zadanie.DocumentId)
		if err != nil {
			// Dokumentu nie ma albo odczyt się nie udał — rozstrzygnięcie tego
			// należy do obsługiwacza, który powie o tym własną treścią odmowy.
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

		// ── Droga trzecia: czynność na całym dokumencie ──────────────────────
		// Zakresu przed wykonaniem nie ma. Kopia zapasowa PRZED (czynność
		// nieodwracalna i tak jej wymaga), wykonanie, uzgodnienie wyniku.
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
		// Treść uzgodniona wraca do bazy: fragmenty zablokowane odzyskują
		// brzmienie zastane, reszta zmiany zostaje.
		po.Tresc = &uzgodnienie.Tresc
		if _, err := a.repozytorium.ZapiszDokument(ctx, po); err != nil {
			return porazka(bladStudio(err))
		}
		return zaporaDopiszBilans(odpowiedz, blokadaBilans(uzgodnienie))
	}
}

// zaporaPrzepiszTresc podmienia pole `content` w ładunku żądania, zostawiając
// wszystkie pozostałe pola nietknięte.
//
// Przez mapę, nie przez strukturę: zapora nie zna kształtów żądań innych
// odcinków, a złożenie ładunku ze znanej jej struktury zgubiłoby każde pole,
// o którym nie wie.
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

// zaporaDopiszBilans dopisuje bilans pominięć do odpowiedzi komendy.
//
// Bilans zastany NIE jest nadpisywany, a scalany: obsługiwacz mógł już oddać
// własny bilans (pominięcia formatu, cechy nieprzeniesione), a bilans blokad
// jest o czym innym. Nadpisanie zamieniłoby jedno przemilczenie na drugie.
func zaporaDopiszBilans(odpowiedz protocol.Odpowiedz,
	bilans shared.StudioActionBalance) protocol.Odpowiedz {

	if odpowiedz.Status != shared.EnvelopeStatusOk {
		return odpowiedz
	}
	pola := map[string]json.RawMessage{}
	if len(odpowiedz.Wynik) > 0 {
		if err := json.Unmarshal(odpowiedz.Wynik, &pola); err != nil {
			// Wynik nie jest obiektem JSON — nie ma gdzie dopisać pola.
			// Milczenie byłoby tu przemilczeniem pominięcia, więc odpowiedzią
			// jest odmowa nazywająca, że bilansu nie dało się oddać.
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

// kopiaPrzedCzynnoscia zakłada kopię zapasową przed czynnością nieodwracalną.
//
// Wymóg Właściciela wymienia trzy takie czynności wprost: przyjęcie wszystkich
// zmian modelu, zamianę w całym dokumencie i zmianę formatu nośnika. Zapora
// zakłada kopię przed KAŻDĄ czynnością na całym dokumencie, bo wszystkie trzy
// wchodzą tą drogą, a dołożenie do wykazu czwartej nie powinno wymagać
// pamiętania o kopii.
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

// sladWykonawcyStudia owija komendy Studia siatką, która dopisuje ŚLAD AUTORA
// tam, gdzie obsługiwacz go nie odłożył.
//
// ── Dlaczego siatka, a nie zaufanie obsługiwaczom ───────────────────────────
// Przełącznik „pokaż wszystko, co zrobił model" stoi na jednym założeniu: każda
// zmiana wykonawcy ma ślad podpisany wykonawcą. Zmierzone: `studio.document.save`
// zawołane przez wykonawcę zmieniało treść dokumentu i NIE odkładało ani zmiany
// śledzonej, ani wpisu dziennika — czyli praca modelu wchodziła do pisma
// niewidzialna dla przełącznika. Właściciel nazwał taką drogę wprost USTERKĄ do
// naprawy, nie ograniczeniem do zgłoszenia.
//
// Naprawa nie może stać w obsługiwaczu tej jednej komendy: komend zmieniających
// treść jest w Studiu ponad setka i pisze je czterech wykonawców, a każda nowa
// mogłaby przeoczyć ślad tak samo. Siatka mierzy więc SKUTEK — czy treść się
// zmieniła — i dopisuje ślad tylko wtedy, gdy obsługiwacz go nie odłożył.
// Obsługiwacz, który odkłada ślad sam (czynności postaci), nie dostaje drugiego.
//
// Siatka NIE zastępuje śladu odkładanego przez obsługiwacza i nie ma zastąpić:
// tamten zna zakres i rodzaj zmiany dokładnie, a siatka zna tylko to, że treść
// jest inna. Dlatego zakres siatki to zakres RÓŻNICY, a nie zgadywane miejsce.
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

// sladOwin składa obsługiwacza z siatką śladu autora.
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
			// Treść ta sama. Zmianę samej POSTACI odkłada droga postaci i tam
			// ślad już jest — siatka nie ma czego tu dopisać, a dopisanie wpisu
			// „coś się stało" bez skutku na treści byłoby zaśmieceniem wykazu.
			return odpowiedz
		}
		if a.sladIleZmianWykonawcy(ctx, po.ID) > sladowPrzed {
			// Obsługiwacz odłożył ślad sam — drugiego nie dokładamy.
			return odpowiedz
		}
		if err := a.sladOdlozBrakujacy(ctx, po, wykonawca, trescPrzed, trescPo); err != nil {
			// Ślad jest tu warunkiem kontroli Operatora nad pracą modelu, więc
			// jego brak nie może przejść ciszą. Skutek na dokumencie już zapadł
			// i odmowa byłaby nieprawdą — dlatego odpowiedź zostaje udana, a brak
			// śladu wraca odmową dopiero wtedy, gdy nie da się go odłożyć wcale.
			return porazka(err)
		}
		return odpowiedz
	}
}

// sladIleZmianWykonawcy liczy zmiany śledzone autora `model` w dokumencie.
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

// sladOdlozBrakujacy odkłada zmianę śledzoną i wpis dziennika dla zmiany treści,
// której obsługiwacz nie podpisał.
//
// Zakres bierze się z RÓŻNICY treści — wspólny przedrostek i wspólny sufiks
// przycięte w ZNAKACH — bo tyle da się o zmianie powiedzieć uczciwie. Zakres
// „cały dokument" byłby wygodniejszy i nieprawdziwy: podświetlenie zaznaczyłoby
// wtedy pismo w całości i przestałoby cokolwiek pokazywać.
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

// sladWycinek wycina fragment treści w znakach; zakres poza treścią daje pustkę.
func sladWycinek(tresc string, od, do int) string {
	znaki := []rune(tresc)
	od, do, poprawny := kontrolaZakresWTresci(znaki, od, do)
	if !poprawny {
		return ""
	}
	return string(znaki[od:do])
}
