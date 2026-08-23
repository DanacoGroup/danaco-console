// Odpowiedzialność pliku: AUTOZAPIS, KOPIE ZAPASOWE i SZEREGI WERSJI —
// `studio.autosave.get/set/run`, `studio.backup.create/list/restore`,
// `studio.version.series.list` i `studio.version.restore.initial`.
//
// ── Dlaczego autozapis idzie OSOBNYM szeregiem wersji ───────────────────────
// Wersje nazwane i kluczowe zakłada Operator; to jest historia jego decyzji.
// Gdyby zapisy samoczynne wchodziły do tego samego wykazu, po godzinie pracy
// historia przestałaby być historią decyzji i stała się dziennikiem naciśnięć
// klawisza — Operator nie znalazłby w niej własnej wersji nazwanej. Rozróżnienie
// idzie kolumną `szereg` (migracja 367) i pojęcia wersji kluczowej NIE zakłada
// się drugiego: Studio ma je w `studio.version.label.set`.
//
// ── Dlaczego kopia zapasowa jest osobna od wersji ────────────────────────────
// Kopia ma przetrwać awarię procesu I AWARIĘ ZAPISU. Wersja leży w repozytorium
// i zakłada się ją tym samym zapisem, który właśnie się nie udał — więc wersja
// nie ochroni pracy przed nieudanym zapisem. Kopia zakłada się niezależnie.
//
// ── Uczciwość zapisu, wymóg bezwzględny ─────────────────────────────────────
// Wskaźnik „zapisano" pokazany, gdy zapis się nie udał, jest najgorszym możliwym
// błędem tego modułu: Operator zamknie okno i straci pracę. Dlatego nieudany
// zapis samoczynny NIE wraca odmową, która przepadnie w logu — wraca odpowiedzią
// niosącą `saved: false`, NAZWANY powód i KOPIĘ, w której praca została. Wiersz
// nastaw zapamiętuje ten skutek, więc następne `autosave.get` powie prawdę także
// wtedy, gdy okno tymczasem się przeładowało.
package core

import (
	"context"
	"encoding/json"
	"errors"

	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// NastawyAutozapisu obsługuje `studio.autosave.get`.
func (a *adapterStudia) NastawyAutozapisu(ctx context.Context,
	z shared.StudioAutosaveGetRequest) (shared.StudioAutosaveGetResponse, error) {

	nastawa, _, err := a.autozapisNastawa(ctx, z.WindowId, z.DocumentId)
	if err != nil {
		return shared.StudioAutosaveGetResponse{}, err
	}
	return shared.StudioAutosaveGetResponse{Settings: autozapisZlozNastawy(nastawa)}, nil
}

// UstawAutozapis obsługuje `studio.autosave.set`.
//
// Zasada wygasania kopii jest tu JAWNYM, ODWRACALNYM ustawieniem Operatora —
// nie stałą rdzenia. Pola pominięte zostają w brzmieniu zastanym: żądanie
// przestawiające sam odstęp nie ma zerować liczby zachowywanych kopii.
func (a *adapterStudia) UstawAutozapis(ctx context.Context,
	z shared.StudioAutosaveSetRequest) (shared.StudioAutosaveSetResponse, error) {

	nastawa, _, err := a.autozapisNastawa(ctx, z.WindowId, z.DocumentId)
	if err != nil {
		return shared.StudioAutosaveSetResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAutosaveSetResponse{}, err
	}

	nastawa.AutozapisCzynny = z.Enabled
	if z.IntervalSeconds != nil {
		if *z.IntervalSeconds < 5 {
			return shared.StudioAutosaveSetResponse{}, bladWskazaniaStudio(
				"odstęp autozapisu krótszy niż pięć sekund — zapis co ułamek sekundy " +
					"zamieniłby historię w dziennik naciśnięć klawisza")
		}
		nastawa.AutozapisOdstepSekund = int64(*z.IntervalSeconds)
	}
	if z.OnBlur != nil {
		nastawa.AutozapisPrzyOdejsciu = *z.OnBlur
	}
	if z.OnClose != nil {
		nastawa.AutozapisPrzyZamknieciu = *z.OnClose
	}
	if z.OnSwitch != nil {
		nastawa.AutozapisPrzyPrzelacz = *z.OnSwitch
	}
	if z.BackupRetentionCount != nil {
		if *z.BackupRetentionCount < 1 {
			return shared.StudioAutosaveSetResponse{}, bladWskazaniaStudio(
				"zasada wygasania kopii zachowująca zero kopii nie jest zasadą " +
					"wygasania, a wyłączeniem kopii zapasowych")
		}
		nastawa.KopieIleZachowac = int64(*z.BackupRetentionCount)
	}
	if z.BackupRetentionHours != nil {
		if *z.BackupRetentionHours < 1 {
			return shared.StudioAutosaveSetResponse{}, bladWskazaniaStudio(
				"kopia wygasająca po mniej niż godzinie nie przetrwałaby awarii procesu")
		}
		nastawa.KopieWygasanieGodzin = int64(*z.BackupRetentionHours)
	}
	if err := skladnica.ZapiszNastaweAutozapisu(ctx, nastawa); err != nil {
		return shared.StudioAutosaveSetResponse{}, bladStudio(err)
	}
	// Odpowiedź czyta stan z bazy, a nie powtarza żądania: ma mówić, jak jest.
	zapisana, _, err := a.autozapisNastawa(ctx, z.WindowId, z.DocumentId)
	if err != nil {
		return shared.StudioAutosaveSetResponse{}, err
	}
	return shared.StudioAutosaveSetResponse{Settings: autozapisZlozNastawy(zapisana)}, nil
}

// WykonajAutozapis obsługuje `studio.autosave.run`.
func (a *adapterStudia) WykonajAutozapis(ctx context.Context,
	z shared.StudioAutosaveRunRequest) (shared.StudioAutosaveRunResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioAutosaveRunResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioAutosaveRunResponse{}, err
	}
	nastawa, err := skladnica.NastawaPracy(ctx, stan.dokument.Okno, &stan.dokument.ID)
	if err != nil {
		return shared.StudioAutosaveRunResponse{}, bladStudio(err)
	}

	var powod shared.StudioBackupReason = shared.StudioBackupReasonInterval
	if z.Trigger != nil && *z.Trigger != "" {
		powod = *z.Trigger
	}

	// Postać z żądania wchodzi PRZED treścią: treść dokumentu wylicza się z
	// drzewa (`postacZapisz`), więc odwrotna kolejność nadpisałaby ją drzewem
	// sprzed edycji.
	if len(z.Form) > 0 {
		var forma shared.StudioDocumentForm
		if err := json.Unmarshal(z.Form, &forma); err != nil {
			return shared.StudioAutosaveRunResponse{}, bladWskazaniaStudio(
				"postać podana do zapisu samoczynnego jest nieczytelna: " + err.Error())
		}
		forma.DocumentId = stan.dokument.Kod
		pominiete, err := a.autozapisPrzywrocPostac(ctx, stan, &forma)
		if err != nil {
			return shared.StudioAutosaveRunResponse{}, err
		}
		_ = pominiete
	}
	if z.Content != nil {
		postacUzgodnijZTrescia(&stan.forma, *z.Content)
	}
	tresc := postacTekstFormy(&stan.forma)

	// KOPIA ZAPASOWA IDZIE PIERWSZA i niesie zmiany niezapisane. Gdyby zapis się
	// nie udał, praca zostaje w niej — i to jest cała jej rola.
	postacJSON, bladZapisuPostaci := json.Marshal(stan.forma)
	kopiaWiersz := dane.KopiaZapasowaStudia{
		Kod:               nowyIdentyfikator(przedrostekKopiiStudia),
		Powod:             string(powod),
		Tresc:             &tresc,
		RozmiarBajtow:     int64(len(tresc)),
		UdaloSie:          true,
		ZmianyNiezapisane: true,
	}
	if bladZapisuPostaci == nil {
		zapis := string(postacJSON)
		kopiaWiersz.PostacJSON = &zapis
	}
	kopia, err := skladnica.ZapiszKopieZapasowa(ctx, stan.dokument.ID, kopiaWiersz)
	if err != nil {
		return shared.StudioAutosaveRunResponse{}, bladStudio(err)
	}
	zlozonaKopia := autozapisZlozKopie(kopia)

	// Zapis właściwy. Niepowodzenie NIE wraca odmową: wraca odpowiedzią mówiącą
	// wprost, że się nie udało, wraz z kopią, w której praca została.
	if err := a.postacZapisz(ctx, stan); err != nil {
		powodNiepowodzenia := "zapis samoczynny dokumentu " + stan.dokument.Kod +
			" nie doszedł do skutku: " + err.Error() +
			". Praca została w kopii zapasowej " + kopia.Kod +
			" — przywróć ją przez studio.backup.restore, zanim zamkniesz okno."
		if bladSkutku := skladnica.ZapiszSkutekAutozapisu(ctx, nastawa.ID, nil, true,
			&powodNiepowodzenia); bladSkutku != nil {

			return shared.StudioAutosaveRunResponse{}, bladStudio(bladSkutku)
		}
		return shared.StudioAutosaveRunResponse{
			Saved:         false,
			FailureReason: &powodNiepowodzenia,
			Backup:        &zlozonaKopia,
		}, nil
	}

	// Zapis się udał, więc kopia przestaje nieść pracę, której w dokumencie nie
	// ma. Wiersz zostaje — jest częścią wykazu kopii — ale zakłada się nowy
	// z wyzerowanym wskaźnikiem, bo tabela nie przestawia tej kolumny po fakcie,
	// a nadpisanie wiersza kopii zatarłoby ślad, że kopia była zakładana przed
	// zapisem.
	chwila := time.Now().UTC().Format(time.RFC3339Nano)
	if err := skladnica.ZapiszSkutekAutozapisu(ctx, nastawa.ID, &chwila, false, nil); err != nil {
		return shared.StudioAutosaveRunResponse{}, bladStudio(err)
	}

	// Wersja idzie w SZEREGU AUTOZAPISU, nie w szeregu Operatora.
	autor := string(shared.StudioAuthorUzytkownik)
	wersjaWiersz := dane.WersjaSzereguStudia{
		Kod:    nowyIdentyfikator(przedrostekWersjiStudio),
		Tresc:  &tresc,
		Autor:  &autor,
		Szereg: string(shared.StudioVersionSeriesAutosave),
	}
	if bladZapisuPostaci == nil {
		zapis := string(postacJSON)
		wersjaWiersz.PostacJSON = &zapis
	}
	wersja, err := skladnica.ZapiszWersjeSzeregu(ctx, stan.dokument.ID, wersjaWiersz)
	if err != nil {
		return shared.StudioAutosaveRunResponse{}, bladStudio(err)
	}
	zlozonaWersja := autozapisZlozWersje(wersja)

	// Przemiatanie wedle nastawy Operatora. Szeregu Operatora nie tyka.
	if _, err := skladnica.PrzemiecWersjeAutozapisu(ctx, stan.dokument.ID,
		nastawa.KopieIleZachowac); err != nil {

		return shared.StudioAutosaveRunResponse{}, bladStudio(err)
	}
	if _, err := skladnica.PrzemiecKopieZapasowe(ctx, stan.dokument.ID,
		nastawa.KopieIleZachowac, nastawa.KopieWygasanieGodzin); err != nil {

		return shared.StudioAutosaveRunResponse{}, bladStudio(err)
	}

	zapisano := chwilaBazy(chwila)
	return shared.StudioAutosaveRunResponse{
		Saved:   true,
		Version: &zlozonaWersja,
		Backup:  &zlozonaKopia,
		SavedAt: &zapisano,
	}, nil
}

// ZalozKopieZapasowa obsługuje `studio.backup.create`.
func (a *adapterStudia) ZalozKopieZapasowa(ctx context.Context,
	z shared.StudioBackupCreateRequest) (shared.StudioBackupCreateResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioBackupCreateResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioBackupCreateResponse{}, err
	}
	var powod shared.StudioBackupReason = shared.StudioBackupReasonManual
	if z.Reason != nil && *z.Reason != "" {
		powod = *z.Reason
	}
	tresc := wartoscTekstu(dokument.Tresc)
	niezapisane := false
	if z.Content != nil {
		niezapisane = *z.Content != tresc
		tresc = *z.Content
	}
	wiersz := dane.KopiaZapasowaStudia{
		Kod:               nowyIdentyfikator(przedrostekKopiiStudia),
		Powod:             string(powod),
		Tresc:             &tresc,
		RozmiarBajtow:     int64(len(tresc)),
		UdaloSie:          true,
		ZmianyNiezapisane: niezapisane,
	}
	if len(z.Form) > 0 {
		zapis := string(z.Form)
		wiersz.PostacJSON = &zapis
	} else {
		stan, err := a.postacWczytaj(ctx, dokument.Kod)
		if err != nil {
			return shared.StudioBackupCreateResponse{}, err
		}
		if postacJSON, err := json.Marshal(stan.forma); err == nil {
			zapis := string(postacJSON)
			wiersz.PostacJSON = &zapis
		}
	}
	kopia, err := skladnica.ZapiszKopieZapasowa(ctx, dokument.ID, wiersz)
	if err != nil {
		return shared.StudioBackupCreateResponse{}, bladStudio(err)
	}
	return shared.StudioBackupCreateResponse{Backup: autozapisZlozKopie(kopia)}, nil
}

// KopieDokumentu obsługuje `studio.backup.list`.
//
// Bez wskazania dokumentu oddaje kopie niosące zmiany NIEZAPISANE po wszystkich
// dokumentach okna — tym Studio samo zgłasza „mam niezapisany dokument z godziny
// X, przywrócić?", zamiast czekać, aż Operator się domyśli.
func (a *adapterStudia) KopieDokumentu(ctx context.Context,
	z shared.StudioBackupListRequest) (shared.StudioBackupListResponse, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioBackupListResponse{}, err
	}
	var wiersze []dane.KopiaZapasowaStudia
	switch {
	case kontrolaTekstNiepusty(z.DocumentId):
		dokument, err := a.dokumentDoCzynnosci(ctx, *z.DocumentId)
		if err != nil {
			return shared.StudioBackupListResponse{}, err
		}
		if wiersze, err = skladnica.KopieZapasowe(ctx, dokument.ID); err != nil {
			return shared.StudioBackupListResponse{}, bladStudio(err)
		}
	default:
		if wiersze, err = skladnica.KopieNiezapisane(ctx, wartoscTekstu(z.WindowId)); err != nil {
			return shared.StudioBackupListResponse{}, bladStudio(err)
		}
	}

	tylkoNiezapisane := z.UnsavedOnly != nil && *z.UnsavedOnly
	odpowiedz := shared.StudioBackupListResponse{Backups: []shared.StudioDocumentBackup{}}
	for _, wiersz := range wiersze {
		if tylkoNiezapisane && !wiersz.ZmianyNiezapisane {
			continue
		}
		odpowiedz.Backups = append(odpowiedz.Backups, autozapisZlozKopie(wiersz))
		if wiersz.ZmianyNiezapisane {
			odpowiedz.UnsavedCount++
		}
	}
	return odpowiedz, nil
}

// PrzywrocKopie obsługuje `studio.backup.restore`.
//
// Przywrócenie DO NOWEGO DOKUMENTU jest osobną drogą, bo przywracanie samo nie
// ma kasować tego, co jest: Operator, który po awarii nie jest pewien, która
// wersja jest lepsza, ma dostać obie.
func (a *adapterStudia) PrzywrocKopie(ctx context.Context,
	z shared.StudioBackupRestoreRequest) (shared.StudioBackupRestoreResponse, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioBackupRestoreResponse{}, err
	}
	if strings.TrimSpace(z.BackupId) == "" {
		return shared.StudioBackupRestoreResponse{},
			bladWskazaniaStudio("przywrócenie bez wskazania kopii zapasowej")
	}
	kopia, err := skladnica.KopiaZapasowa(ctx, strings.TrimSpace(z.BackupId))
	if err != nil {
		return shared.StudioBackupRestoreResponse{}, dziennikBrakWiersza(err,
			"kopia zapasowa nie istnieje: "+z.BackupId)
	}
	if !kopia.UdaloSie {
		// Kopia NIEUDANA nie niesie pracy — niesie ślad, że praca nie doszła na
		// dysk. Przywrócenie z niej wstawiłoby do dokumentu pustkę i nazwało to
		// przywróceniem.
		powod := "kopia " + kopia.Kod + " jest zapisem NIEUDANYM"
		if kopia.PowodNiepowodzenia != nil {
			powod += " (" + *kopia.PowodNiepowodzenia + ")"
		}
		return shared.StudioBackupRestoreResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Studio: "+powod+
				" — nie niesie treści do przywrócenia, a jest śladem, że zapis się nie udał"))
	}
	if kopia.Tresc == nil {
		return shared.StudioBackupRestoreResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeConflict, "moduł Studio: kopia "+kopia.Kod+
				" nie niesie treści — przywrócenie z niej zostawiłoby dokument pusty"))
	}

	zrodlowy, err := a.dokumentDoCzynnosci(ctx, kopia.DokumentKod)
	if err != nil {
		return shared.StudioBackupRestoreResponse{}, err
	}

	docelowyKod := zrodlowy.Kod
	if z.AsNewDocument != nil && *z.AsNewDocument {
		tytul := "kopia dokumentu " + zrodlowy.Kod
		if zrodlowy.Tytul != nil && *zrodlowy.Tytul != "" {
			tytul = *zrodlowy.Tytul + " (przywrócony z kopii)"
		}
		if kontrolaTekstNiepusty(z.Title) {
			tytul = *z.Title
		}
		okno := zrodlowy.Okno
		if kontrolaTekstNiepusty(z.WindowId) {
			okno = *z.WindowId
		}
		nowy, err := a.repozytorium.ZapiszDokument(ctx, dane.DokumentStudia{
			Kod:    nowyIdentyfikator(przedrostekDokumentuStudio),
			Okno:   okno,
			Tytul:  &tytul,
			Format: zrodlowy.Format,
			Tresc:  kopia.Tresc,
		})
		if err != nil {
			return shared.StudioBackupRestoreResponse{}, bladStudio(err)
		}
		docelowyKod = nowy.Kod
	}

	stan, err := a.postacWczytaj(ctx, docelowyKod)
	if err != nil {
		return shared.StudioBackupRestoreResponse{}, err
	}
	if kopia.PostacJSON != nil && strings.TrimSpace(*kopia.PostacJSON) != "" {
		var forma shared.StudioDocumentForm
		if err := json.Unmarshal([]byte(*kopia.PostacJSON), &forma); err != nil {
			return shared.StudioBackupRestoreResponse{}, kontrolaBladZaplecza(
				"kopia " + kopia.Kod + " niesie nieczytelną postać: " + err.Error())
		}
		forma.DocumentId = stan.dokument.Kod
		if _, err := a.autozapisPrzywrocPostac(ctx, stan, &forma); err != nil {
			return shared.StudioBackupRestoreResponse{}, err
		}
	}
	postacUzgodnijZTrescia(&stan.forma, *kopia.Tresc)
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioBackupRestoreResponse{}, err
	}
	return shared.StudioBackupRestoreResponse{
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
	}, nil
}

// WersjeWSzeregach obsługuje `studio.version.series.list`.
func (a *adapterStudia) WersjeWSzeregach(ctx context.Context,
	z shared.StudioVersionSeriesListRequest) (shared.StudioVersionSeriesListResponse, error) {

	dokument, err := a.dokumentDoCzynnosci(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioVersionSeriesListResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioVersionSeriesListResponse{}, err
	}
	wiersze, err := skladnica.WersjeSzeregow(ctx, dokument.ID)
	if err != nil {
		return shared.StudioVersionSeriesListResponse{}, bladStudio(err)
	}

	odpowiedz := shared.StudioVersionSeriesListResponse{Versions: []shared.StudioVersion{}}
	// Licznik obu szeregów liczy się z CAŁEGO wykazu, nie z tego, co zostało po
	// zawężeniu: przełącznik „pokaż także zapisy samoczynne" ma pokazywać, ile
	// ich jest, jeszcze przed włączeniem.
	for _, wiersz := range wiersze {
		if wiersz.Szereg == string(shared.StudioVersionSeriesAutosave) {
			odpowiedz.AutosaveCount++
		} else {
			odpowiedz.OperatorCount++
		}
	}
	for _, wiersz := range wiersze {
		if z.Series != nil && *z.Series != "" && wiersz.Szereg != string(*z.Series) {
			continue
		}
		if z.MilestonesOnly != nil && *z.MilestonesOnly && !wiersz.KamienMilowy {
			continue
		}
		if z.Limit != nil && *z.Limit > 0 && len(odpowiedz.Versions) >= *z.Limit {
			break
		}
		odpowiedz.Versions = append(odpowiedz.Versions, autozapisZlozWersje(wiersz))
	}
	zalozycielska, err := skladnica.WersjaZalozycielska(ctx, dokument.ID)
	switch {
	case err == nil:
		odpowiedz.InitialVersionId = &zalozycielska.Kod
	case errors.Is(err, dane.ErrBrakWiersza):
		// Dokument bez ani jednej wersji Operatora nie ma wersji założycielskiej
		// i to jest prawda o nim, nie brak wiedzy rdzenia.
	default:
		return shared.StudioVersionSeriesListResponse{}, bladStudio(err)
	}
	return odpowiedz, nil
}

// PrzywrocWersjeZalozycielska obsługuje `studio.version.restore.initial`.
//
// Wersje nowsze ZOSTAJĄ — tak samo jak przy `studio.repository.restore` — więc
// samo cofnięcie do stanu pierwotnego jest odwracalne. Wersja stanu bieżącego
// zakłada się PRZED powrotem, bo inaczej praca sprzed powrotu nie miałaby
// w historii ani jednego punktu, do którego Operator mógłby wrócić.
func (a *adapterStudia) PrzywrocWersjeZalozycielska(ctx context.Context,
	z shared.StudioVersionRestoreInitialRequest) (shared.StudioVersionRestoreInitialResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioVersionRestoreInitialResponse{}, err
	}
	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return shared.StudioVersionRestoreInitialResponse{}, err
	}
	zalozycielska, err := skladnica.WersjaZalozycielska(ctx, stan.dokument.ID)
	if err != nil {
		return shared.StudioVersionRestoreInitialResponse{}, dziennikBrakWiersza(err,
			"dokument "+stan.dokument.Kod+" nie ma wersji założycielskiej — historia "+
				"szeregu Operatora jest pusta, więc nie ma do czego wracać")
	}
	if zalozycielska.Tresc == nil {
		return shared.StudioVersionRestoreInitialResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeConflict, "moduł Studio: wersja "+
				"założycielska "+zalozycielska.Kod+" nie niesie treści — powrót do niej "+
				"zostawiłby dokument pusty"))
	}

	if z.CreateVersion == nil || *z.CreateVersion {
		if _, err := a.zalozWersjeDokumentu(ctx, stan.dokument,
			wartoscTekstu(stan.dokument.Tresc), shared.StudioAuthorUzytkownik, nil); err != nil {

			return shared.StudioVersionRestoreInitialResponse{}, err
		}
		// Dokument przeczytany ponownie, bo założenie wersji przestawiło
		// wskaźnik wersji bieżącej.
		if stan, err = a.postacWczytaj(ctx, z.DocumentId); err != nil {
			return shared.StudioVersionRestoreInitialResponse{}, err
		}
	}

	// Kopia przed czynnością nieodwracalną — powrót do stanu pierwotnego jest
	// dokładnie taką czynnością.
	if _, err := a.kopiaPrzedCzynnoscia(ctx, stan.dokument,
		shared.StudioBackupReasonBeforeIrreversible); err != nil {

		return shared.StudioVersionRestoreInitialResponse{}, err
	}

	if zalozycielska.PostacJSON != nil && strings.TrimSpace(*zalozycielska.PostacJSON) != "" {
		var forma shared.StudioDocumentForm
		if err := json.Unmarshal([]byte(*zalozycielska.PostacJSON), &forma); err != nil {
			return shared.StudioVersionRestoreInitialResponse{}, kontrolaBladZaplecza(
				"wersja założycielska niesie nieczytelną postać: " + err.Error())
		}
		forma.DocumentId = stan.dokument.Kod
		if _, err := a.autozapisPrzywrocPostac(ctx, stan, &forma); err != nil {
			return shared.StudioVersionRestoreInitialResponse{}, err
		}
	}
	postacUzgodnijZTrescia(&stan.forma, *zalozycielska.Tresc)
	if err := a.postacZapisz(ctx, stan); err != nil {
		return shared.StudioVersionRestoreInitialResponse{}, err
	}
	drzewo, err := dziennikZapisDrzewa(stan.forma)
	if err != nil {
		return shared.StudioVersionRestoreInitialResponse{}, err
	}
	if _, err := a.dziennikOdlozCzynnosc(ctx, stan.dokument,
		kontrolaWykonawca{Rodzaj: shared.StudioAuthorUzytkownik},
		shared.StudioActionKindTextEdit,
		"powrót dokumentu do wersji założycielskiej "+zalozycielska.Kod,
		nil, nil, drzewo, drzewo, nil); err != nil {

		return shared.StudioVersionRestoreInitialResponse{}, err
	}
	return shared.StudioVersionRestoreInitialResponse{
		Document: a.zlozDokument(stan.dokument),
		Form:     stan.forma,
		Version:  autozapisZlozWersje(zalozycielska),
	}, nil
}

// ── Wspólne ─────────────────────────────────────────────────────────────────

// autozapisNastawa odczytuje wiersz nastaw pracy dla pary okno-dokument albo dla
// samego okna, zakładając go, gdy Operator nigdy nastaw nie ruszał.
//
// Okno bierze się ze wskazania, a gdy go nie ma — z dokumentu. Nastawy bez okna
// nie istnieją: nastawa okna jest wierszem tej samej tabeli i musi wiedzieć,
// czyja jest.
func (a *adapterStudia) autozapisNastawa(ctx context.Context, oknoZadania,
	dokumentZadania *string) (dane.NastawaPracyStudia, dane.DokumentStudia, error) {

	skladnica, err := a.kontrolaSkladnica()
	if err != nil {
		return dane.NastawaPracyStudia{}, dane.DokumentStudia{}, err
	}
	var dokument dane.DokumentStudia
	var wskazanie *int64
	okno := wartoscTekstu(oknoZadania)
	if kontrolaTekstNiepusty(dokumentZadania) {
		if dokument, err = a.dokumentDoCzynnosci(ctx, *dokumentZadania); err != nil {
			return dane.NastawaPracyStudia{}, dane.DokumentStudia{}, err
		}
		wskazanie = &dokument.ID
		if okno == "" {
			okno = dokument.Okno
		}
	}
	if okno == "" {
		return dane.NastawaPracyStudia{}, dane.DokumentStudia{}, bladWskazaniaStudio(
			"nastawy bez wskazania okna ani dokumentu — nastawa musi wiedzieć, czyja jest")
	}
	nastawa, err := skladnica.NastawaPracy(ctx, okno, wskazanie)
	if err != nil {
		return dane.NastawaPracyStudia{}, dane.DokumentStudia{}, bladStudio(err)
	}
	return nastawa, dokument, nil
}

// autozapisPrzywrocPostac nakłada drzewo docelowe na stan bieżący, wraz z bytami
// o własnych wierszach. Rachunek jest ten sam, co przy cofaniu czynności —
// dlatego woła się go, a nie pisze drugi raz.
func (a *adapterStudia) autozapisPrzywrocPostac(ctx context.Context, stan *stanPostaci,
	cel *shared.StudioDocumentForm) ([]shared.StudioSkippedItem, error) {

	biezaca := stan.forma
	stan.forma.Blocks = cel.Blocks
	if cel.PageSetup != nil {
		stan.forma.PageSetup = cel.PageSetup
	}
	return a.dziennikPrzywrocByty(ctx, stan, &biezaca, cel)
}

// autozapisZlozNastawy składa nastawy autozapisu kontraktu z wiersza.
//
// Czas ostatniego zapisu i wskaźnik niepowodzenia idą z bazy, nie z pamięci
// procesu: okno przeładowane po awarii ma zobaczyć prawdę o ostatnim zapisie,
// a nie stan czysty.
func autozapisZlozNastawy(wiersz dane.NastawaPracyStudia) shared.StudioAutosaveSettings {
	odstep := int(wiersz.AutozapisOdstepSekund)
	ileKopii := int(wiersz.KopieIleZachowac)
	godzin := int(wiersz.KopieWygasanieGodzin)
	nastawy := shared.StudioAutosaveSettings{
		Enabled:              wiersz.AutozapisCzynny,
		IntervalSeconds:      &odstep,
		OnBlur:               wskaznikLogiczny(wiersz.AutozapisPrzyOdejsciu),
		OnClose:              wskaznikLogiczny(wiersz.AutozapisPrzyZamknieciu),
		OnSwitch:             wskaznikLogiczny(wiersz.AutozapisPrzyPrzelacz),
		BackupRetentionCount: &ileKopii,
		BackupRetentionHours: &godzin,
		LastSaveFailed:       wskaznikLogiczny(wiersz.OstatniZapisNieudany),
		LastFailureReason:    wiersz.OstatniPowodNiepowodz,
	}
	if wiersz.OstatniZapis != nil && *wiersz.OstatniZapis != "" {
		chwila := chwilaBazy(*wiersz.OstatniZapis)
		nastawy.LastSaveAt = &chwila
	}
	return nastawy
}

// autozapisZlozKopie składa kopię zapasową kontraktu z wiersza.
func autozapisZlozKopie(wiersz dane.KopiaZapasowaStudia) shared.StudioDocumentBackup {
	bajtow := wiersz.RozmiarBajtow
	return shared.StudioDocumentBackup{
		Id:             wiersz.Kod,
		DocumentId:     wiersz.DokumentKod,
		Reason:         shared.StudioBackupReason(wiersz.Powod),
		Bytes:          &bajtow,
		Succeeded:      wiersz.UdaloSie,
		FailureReason:  wiersz.PowodNiepowodzenia,
		UnsavedChanges: wskaznikLogiczny(wiersz.ZmianyNiezapisane),
		CreatedAt:      chwilaBazy(wiersz.Utworzono),
	}
}

// autozapisZlozWersje składa wersję kontraktu z wiersza szeregu.
//
// Autor wychodzi WYŁĄCZNIE wtedy, gdy wiersz go niesie: wersje założone przed
// dobudową tego pola autora nie mają, a podstawienie tu Operatora zamieniłoby
// brak wiedzy w twierdzenie fałszywe dla każdej wersji zapisanej przez model.
func autozapisZlozWersje(wiersz dane.WersjaSzereguStudia) shared.StudioVersion {
	wersja := shared.StudioVersion{
		Id:              wiersz.Kod,
		DocumentId:      wiersz.DokumentKod,
		Label:           wiersz.Etykieta,
		Summary:         wiersz.Podsumowanie,
		ContentHash:     wiersz.SkrotTresci,
		CreatedAt:       chwilaBazy(wiersz.Utworzono),
		Milestone:       wskaznikLogiczny(wiersz.KamienMilowy),
		BranchId:        wiersz.GalazKod,
		ProposalId:      wiersz.PropozycjaKod,
		AuthorAgentId:   wiersz.AutorAgentKod,
		AuthorAgentName: wiersz.AutorAgentNazwa,
	}
	if wiersz.Autor != nil && *wiersz.Autor != "" {
		autor := shared.StudioAuthor(*wiersz.Autor)
		wersja.Author = &autor
	}
	// Szereg wychodzi opisem w podsumowaniu wyłącznie wtedy, gdy wiersz go nie
	// niesie inaczej: kontrakt wersji nie ma pola szeregu, a wykaz rozdziela je
	// licznikami odpowiedzi. Etykieta zapisu samoczynnego mówi to Operatorowi
	// wprost, żeby pozycja wykazu była odróżnialna także w oderwaniu od licznika.
	if wiersz.Szereg == string(shared.StudioVersionSeriesAutosave) && wersja.Label == nil {
		etykieta := "zapis samoczynny"
		wersja.Label = &etykieta
	}
	return wersja
}
