// Sprawdzian odczytów warstwy urządzeń: przekładu odpowiedzi obcych warstw
// (SANE, WIA, CUPS) na kształt rdzenia oraz znakowania wartości wchodzących do
// skryptu PowerShella.
//
// Sprawdzian nie uruchamia ani jednego programu — sprawdza to, co da się
// sprawdzić bez skanera i bez drukarki: czy rdzeń odczytuje odpowiedź, i czy
// odmawia tam, gdzie odpowiedzi nie zrozumiał, zamiast oddać pusty wykaz.
package core

import (
	"testing"

	"danacoconsole/shared"
)

func TestOdczytDrukarekCupsRozpoznajeDomyslnaIWstrzymana(t *testing.T) {
	wyjscie := "printer HP-42 is idle.  enabled since Mon 18 Aug 2026\n" +
		"printer Etykiety disabled since Mon 18 Aug 2026 -\n" +
		"system default destination: HP-42\n"

	drukarki := odczytajDrukarkiCups(wyjscie)
	if len(drukarki) != 2 {
		t.Fatalf("oczekiwano dwóch kolejek, jest %d: %+v", len(drukarki), drukarki)
	}
	if drukarki[0].Nazwa != "HP-42" || !drukarki[0].Domyslna || !drukarki[0].Gotowa {
		t.Errorf("kolejka domyślna odczytana błędnie: %+v", drukarki[0])
	}
	if drukarki[1].Nazwa != "Etykiety" || drukarki[1].Domyslna || drukarki[1].Gotowa {
		t.Errorf("kolejka wstrzymana odczytana błędnie: %+v", drukarki[1])
	}
}

func TestOdczytUrzadzenWiaOdmawiaZamiastPustegoWykazu(t *testing.T) {
	// Pusta odpowiedź nie jest wykazem pustym: warstwa nie odpowiedziała wcale.
	if _, err := odczytajUrzadzeniaWia("   \n"); err == nil {
		t.Error("brak odpowiedzi warstwy WIA oddany jako pusty wykaz — to twierdzenie " +
			"„szukałem i nic nie ma”, którego rdzeń nie ma prawa postawić")
	}
	// Odpowiedź niezrozumiała też jest odmową, a nie pustką.
	if _, err := odczytajUrzadzeniaWia("Get-Printer : polecenie nieznane"); err == nil {
		t.Error("niezrozumiała odpowiedź warstwy WIA oddana jako pusty wykaz")
	}
	// Wykaz pusty ODDANY przez warstwę jest wynikiem prawdziwym.
	urzadzenia, err := odczytajUrzadzeniaWia("[]")
	if err != nil {
		t.Fatalf("pusty wykaz oddany przez warstwę potraktowany jako usterka: %v", err)
	}
	if len(urzadzenia) != 0 {
		t.Errorf("z pustego wykazu powstały urządzenia: %+v", urzadzenia)
	}

	urzadzenia, err = odczytajUrzadzeniaWia(
		`[{"id":"\\\\.\\Usbscan0","name":"Canon LiDE 400"},{"id":"WIA-2","name":""}]`)
	if err != nil {
		t.Fatalf("wykazu WIA nie da się odczytać: %v", err)
	}
	if len(urzadzenia) != 2 {
		t.Fatalf("oczekiwano dwóch urządzeń, jest %d", len(urzadzenia))
	}
	if urzadzenia[0].Name != "Canon LiDE 400" ||
		urzadzenia[0].Kind != shared.StudioInputDeviceKindSkaner {
		t.Errorf("urządzenie odczytane błędnie: %+v", urzadzenia[0])
	}
	if urzadzenia[1].Name != "WIA-2" {
		t.Errorf("urządzenie bez nazwy ma wystąpić pod swoim identyfikatorem: %+v", urzadzenia[1])
	}
}

func TestNapisPowerShellaNieWypuszczaWartosciWPolecenie(t *testing.T) {
	if wynik := napisPowerShell("skaner'; Remove-Item C:\\ -Recurse; '"); wynik !=
		`'skaner''; Remove-Item C:\ -Recurse; '''` {
		t.Errorf("apostrof nie został podwojony, wartość wychodzi z napisu: %s", wynik)
	}
}
