package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Sprawdziany autozapisu, kopii zapasowych i szeregów wersji.
//
// Szkody, które ten plik ma wykluczyć:
//  1. wskaźnik „zapisano" pokazany po zapisie NIEUDANYM — Operator zamknie okno
//     i straci pracę; to jest najgorszy możliwy błąd tego modułu;
//  2. zapisy samoczynne wchodzące do wykazu wersji Operatora, przez co historia
//     decyzji zamienia się w dziennik naciśnięć klawisza;
//  3. autor wersji podstawiony jako Operator tam, gdzie wiersz go nie niesie —
//     twierdzenie fałszywe dla każdej wersji zapisanej przez model.

// TestAutozapisNieudanyZapisJestWidocznyINazwany jest sednem uczciwości zapisu.
func TestAutozapisNieudanyZapisJestWidocznyINazwany(t *testing.T) {
	powod := "dysk odmówił zapisu: brak miejsca"
	nastawy := autozapisZlozNastawy(dane.NastawaPracyStudia{
		AutozapisCzynny:       true,
		AutozapisOdstepSekund: 120,
		KopieIleZachowac:      20,
		KopieWygasanieGodzin:  168,
		OstatniZapisNieudany:  true,
		OstatniPowodNiepowodz: &powod,
	})
	if nastawy.LastSaveFailed == nil || !*nastawy.LastSaveFailed {
		t.Fatalf("nieudany zapis nie wyszedł kontraktem — okno pokaże „zapisano”")
	}
	if nastawy.LastFailureReason == nil || *nastawy.LastFailureReason != powod {
		t.Errorf("powód niepowodzenia nie wyszedł kontraktem: %+v", nastawy.LastFailureReason)
	}
	if nastawy.LastSaveAt != nil {
		t.Errorf("nastawy po nieudanym zapisie oddały czas udanego zapisu: %+v",
			nastawy.LastSaveAt)
	}
}

// TestAutozapisKopiaNieudanaZostajeWidoczna mierzy, że kopia NIEUDANA nie znika
// razem z niepowodzeniem: jest jedynym śladem, że praca nie doszła na dysk.
func TestAutozapisKopiaNieudanaZostajeWidoczna(t *testing.T) {
	powod := "zapis samoczynny nie doszedł do skutku"
	kopia := autozapisZlozKopie(dane.KopiaZapasowaStudia{
		Kod:                "studio-kop-pierwsza",
		DokumentKod:        "studio-dok-pierwszy",
		Powod:              string(shared.StudioBackupReasonInterval),
		RozmiarBajtow:      4096,
		UdaloSie:           false,
		PowodNiepowodzenia: &powod,
		ZmianyNiezapisane:  true,
	})
	if kopia.Succeeded {
		t.Fatalf("kopia nieudana wyszła jako udana — Operator przywróciłby z niej pustkę")
	}
	if kopia.FailureReason == nil {
		t.Errorf("kopia nieudana nie nazywa powodu")
	}
	if kopia.UnsavedChanges == nil || !*kopia.UnsavedChanges {
		t.Errorf("kopia niosąca zmiany niezapisane nie jest tak oznaczona — Studio nie " +
			"zgłosi „mam niezapisany dokument z godziny X”")
	}
	if kopia.Bytes == nil || *kopia.Bytes != 4096 {
		t.Errorf("rozmiar kopii nie wyszedł kontraktem: %+v", kopia.Bytes)
	}
}

// TestAutozapisWersjaSzereguJestOdroznialna mierzy drugą szkodę: pozycja szeregu
// autozapisu daje się odróżnić w wykazie także w oderwaniu od licznika.
func TestAutozapisWersjaSzereguJestOdroznialna(t *testing.T) {
	tresc := "Treść zapisu samoczynnego."
	samoczynna := autozapisZlozWersje(dane.WersjaSzereguStudia{
		Kod: "studio-wer-samoczynna", DokumentKod: "studio-dok-pierwszy",
		Tresc: &tresc, Szereg: string(shared.StudioVersionSeriesAutosave),
	})
	if samoczynna.Label == nil || !strings.Contains(*samoczynna.Label, "samoczynny") {
		t.Errorf("wersja szeregu autozapisu nie jest odróżnialna w wykazie: %+v",
			samoczynna.Label)
	}

	// Etykieta nadana przez Operatora NIE jest nadpisywana — jego nazwa własna
	// jest ważniejsza niż podpowiedź rdzenia.
	wlasna := "wersja do wysłania"
	nazwana := autozapisZlozWersje(dane.WersjaSzereguStudia{
		Kod: "studio-wer-nazwana", Etykieta: &wlasna,
		Szereg: string(shared.StudioVersionSeriesAutosave),
	})
	if nazwana.Label == nil || *nazwana.Label != wlasna {
		t.Errorf("etykieta Operatora została nadpisana: %+v", nazwana.Label)
	}
}

// TestAutozapisAutorWersjiNiePodstawiaOperatora mierzy trzecią szkodę.
func TestAutozapisAutorWersjiNiePodstawiaOperatora(t *testing.T) {
	bezAutora := autozapisZlozWersje(dane.WersjaSzereguStudia{
		Kod: "studio-wer-stara", Szereg: string(shared.StudioVersionSeriesOperator),
	})
	if bezAutora.Author != nil {
		t.Errorf("wersja bez zapisanego autora dostała autora — brak wiedzy zamieniono "+
			"w twierdzenie: %q", *bezAutora.Author)
	}
	autorModelu := string(shared.StudioAuthorModel)
	kod := "agent-redaktor"
	zModelem := autozapisZlozWersje(dane.WersjaSzereguStudia{
		Kod: "studio-wer-modelu", Autor: &autorModelu, AutorAgentKod: &kod,
		Szereg: string(shared.StudioVersionSeriesOperator),
	})
	if zModelem.Author == nil || *zModelem.Author != shared.StudioAuthorModel {
		t.Errorf("autor `model` zgubił się w przekładzie: %+v", zModelem.Author)
	}
	if zModelem.AuthorAgentId == nil || *zModelem.AuthorAgentId != kod {
		t.Errorf("tożsamość wykonawcy wersji zgubiła się w przekładzie")
	}
}

// TestAutozapisNastawyWygasaniaSaJawne mierzy, że zasada wygasania kopii wychodzi
// kontraktem jako wartość Operatora, a nie jako stała rdzenia ukryta w kodzie.
func TestAutozapisNastawyWygasaniaSaJawne(t *testing.T) {
	nastawy := autozapisZlozNastawy(dane.NastawaPracyStudia{
		AutozapisCzynny:         true,
		AutozapisOdstepSekund:   45,
		AutozapisPrzyOdejsciu:   true,
		AutozapisPrzyZamknieciu: true,
		AutozapisPrzyPrzelacz:   false,
		KopieIleZachowac:        12,
		KopieWygasanieGodzin:    72,
	})
	if nastawy.IntervalSeconds == nil || *nastawy.IntervalSeconds != 45 {
		t.Errorf("odstęp autozapisu nie wyszedł kontraktem: %+v", nastawy.IntervalSeconds)
	}
	if nastawy.BackupRetentionCount == nil || *nastawy.BackupRetentionCount != 12 {
		t.Errorf("liczba zachowywanych kopii nie wyszła kontraktem")
	}
	if nastawy.BackupRetentionHours == nil || *nastawy.BackupRetentionHours != 72 {
		t.Errorf("czas wygasania kopii nie wyszedł kontraktem")
	}
	if nastawy.OnSwitch == nil || *nastawy.OnSwitch {
		t.Errorf("zapis przy przełączeniu wyłączony przez Operatora wyszedł jako włączony")
	}
}

// TestAutozapisSzeregiSaDwaINieRosna pilnuje, żeby nikt nie dołożył trzeciego
// pojęcia szeregu: rozróżnienie ma zostać dwuwartościowe, a wersję kluczową
// Studio ma już w `studio.version.label.set`.
func TestAutozapisSzeregiSaDwaINieRosna(t *testing.T) {
	szeregi := shared.WartosciStudioVersionSeries()
	if len(szeregi) != 2 {
		t.Fatalf("szeregów wersji ma być dwa, jest %d: %+v", len(szeregi), szeregi)
	}
	if string(szeregi[0]) != "operator" || string(szeregi[1]) != "autosave" {
		t.Errorf("słownik szeregów rozjechał się z migracją 367: %+v", szeregi)
	}
}
