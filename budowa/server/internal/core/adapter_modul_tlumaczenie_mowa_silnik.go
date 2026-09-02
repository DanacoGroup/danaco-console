// Odpowiedzialność pliku: silnik syntezy mowy modułu Translate — wybór
// syntezatora, wiersz wywołania, sprawdzenie skutku i przekład odmów na kody
// kontraktu. Uruchomienia procesu tu nie ma: idzie ono pakietem zewnetrzne.
package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// Nazwy silników. Wchodzą do nazwy pliku nagrania, więc Operator czyta
	// z `path`, którym silnikiem powstało to, czego słucha.
	silnikPiper  = "piper"
	silnikEspeak = "espeak-ng"

	// Wskazania Operatora. Zmienna środowiska ma pierwszeństwo przed wykryciem
	// w miejscach typowych — uzasadnienie w nagłówku pliku.
	zmiennaPipera      = "DANACO_PIPER"
	zmiennaGlosowPiper = "DANACO_PIPER_GLOSY"
	zmiennaEspeaka     = "DANACO_ESPEAK"

	// Miejsca typowe arsenału. Ścieżki zapasowe, nie jedyne — sprawdzane
	// dopiero, gdy goła nazwa nie stoi w PATH.
	piperArsenalu       = "/opt/danaco-arsenal/piper/bin/piper"
	glosyPiperaArsenalu = "/opt/danaco-arsenal/glosy"

	// katalogSyntezyMowy to podkatalog katalogu danych rdzenia, w którym lądują
	// nagrania. Osobny, żeby dało się je skasować hurtem bez ruszania bazy.
	katalogSyntezyMowy = "synteza-mowy"

	// granicaSyntezy to granica czasu jednego uruchomienia syntezatora. Piper
	// ładuje ciężki model przy każdym wywołaniu; minuta zostawia zapas, a
	// zarazem nie pozwala utkniętemu programowi trzymać żądania bez końca.
	granicaSyntezy = 60 * time.Second
)

// ZSynteza wpina silnik syntezy mowy: uruchamiacz procesów, dwa źródła
// izolacji i katalog danych rdzenia, pod którym lądują nagrania. Zależność
// opcjonalna: jej brak psuje wyłącznie `speech.synthesize`, która wtedy
// odmawia nazywając brak.
func (a *adapterTlumaczenia) ZSynteza(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy, katalogDanych string) *adapterTlumaczenia {

	a.uruchamiacz = uruchamiacz
	a.rozstrzygaczMowy = rozstrzygacz
	a.katalogIzolacji = katalog
	a.katalogDanych = katalogDanych
	return a
}

// zsyntezujDoPliku zamienia tekst na plik WAV i oddaje jego ścieżkę.
// Kolejność jest rozmyślna: najpierw pytanie, czym syntezować, dopiero potem
// zakłada się cokolwiek na dysku.
func (a *adapterTlumaczenia) zsyntezujDoPliku(ctx context.Context,
	kodPanelu, jezyk, tekst string) (string, error) {

	if a.uruchamiacz == nil {
		return "", bladBrakuSyntezatora("serwer nie ma uruchamiacza procesów — nie ma czym wystartować " +
			"syntezatora mowy; naprawa: podpiąć warstwę kanału (injection) przy składaniu serwera")
	}

	wybor, err := dobierzSyntezator(jezyk)
	if err != nil {
		return "", err
	}

	katalog, err := a.katalogNagran()
	if err != nil {
		return "", err
	}
	sciezka := filepath.Join(katalog, nazwaNagrania(kodPanelu, wybor.silnik))

	// Tekst idzie plikiem, nie argumentem — port procesu nie wystawia
	// wejścia, a argument ma granicę.
	sciezkaTekstu := sciezka + ".txt"
	if err := os.WriteFile(sciezkaTekstu, []byte(tekst), 0o600); err != nil {
		return "", bladSyntezyMowy("nie można zapisać tekstu do syntezy pod " + sciezkaTekstu + ": " + err.Error())
	}
	// Plik tekstowy jest rusztowaniem, nie wynikiem — znika niezależnie od
	// tego, czy synteza się udała.
	defer func() { _ = os.Remove(sciezkaTekstu) }()

	okno, zasady, obszar := a.zasiegProgramowTlumaczenia(ctx)
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		wybor.narzedzie, wybor.argumenty(sciezkaTekstu, sciezka), katalogPracySyntezy(obszar), granicaSyntezy)
	if err != nil {
		// Naruszenie izolacji jest znakowane osobno — punkt izolacji Operatora to
		// nie jest usterka rdzenia.
		if errors.Is(err, session.ErrIzolacja) {
			return "", bladIzolacjiSyntezy(err)
		}
		return "", bladSyntezyMowy(err.Error() + ogonSyntezatora(wynik.Diagnostyka))
	}

	// Plik ma istnieć i mieć rozmiar — liczy się skutek, nie kod wyjścia;
	// piper bywa zerem bez nagrania.
	opis, err := os.Stat(sciezka)
	if err != nil {
		return "", bladSyntezyMowy(wybor.silnik + " zakończył się powodzeniem, ale nie zostawił nagrania pod " +
			sciezka + " (" + err.Error() + ")" + ogonSyntezatora(wynik.Diagnostyka))
	}
	if opis.Size() == 0 {
		_ = os.Remove(sciezka)
		return "", bladSyntezyMowy(wybor.silnik + " zostawił nagranie puste (0 bajtów) pod " + sciezka +
			ogonSyntezatora(wynik.Diagnostyka))
	}
	return sciezka, nil
}

// wyborSyntezatora niesie rozstrzygnięcie „czym jest czytane": nazwę silnika (idzie
// w nazwę pliku), narzędzie dla `zewnetrzne.Wolaj` i głos, gdy silnik go używa.
type wyborSyntezatora struct {
	silnik    string
	narzedzie zewnetrzne.Narzedzie
	glos      string
}

// argumenty składa wiersz wywołania wybranego silnika. Jedno miejsce na oba
// wiersze, bo różnica między silnikami sprowadza się właśnie do nich — reszta
// drogi (plik tekstu, plik wyjścia, granica czasu, brama izolacji) jest wspólna.
func (w wyborSyntezatora) argumenty(sciezkaTekstu, sciezkaNagrania string) []string {
	if w.silnik == silnikPiper {
		return []string{"-m", w.glos, "-i", sciezkaTekstu, "-f", sciezkaNagrania}
	}
	// Espeak dostaje język w -v bez podmiany na domyślny, żeby nie czytać
	// cicho w złym języku.
	return []string{"-v", w.glos, "-w", sciezkaNagrania, "-f", sciezkaTekstu}
}

// katalogNagran zakłada, gdy trzeba, katalog na nagrania pod katalogiem
// danych rdzenia. Pusty katalog danych jest odmową: nagranie zapisane
// w katalogu bieżącym procesu wylądowałoby tam, gdzie Operator go nie szuka.
func (a *adapterTlumaczenia) katalogNagran() (string, error) {
	podstawa := strings.TrimSpace(a.katalogDanych)
	if podstawa == "" {
		return "", bladBrakuSyntezatora("serwer nie zna katalogu danych — nie ma gdzie zapisać nagrania; " +
			"naprawa: wskazać katalog danych przełącznikiem -dane albo zmienną DANACO_KATALOG_DANYCH")
	}
	katalog := filepath.Join(podstawa, katalogSyntezyMowy)
	if err := os.MkdirAll(katalog, 0o755); err != nil {
		return "", bladSyntezyMowy("nie można założyć katalogu nagrań " + katalog + ": " + err.Error())
	}
	return katalog, nil
}

// nazwaNagrania składa nazwę pliku z kodu panelu, nazwy silnika i chwili
// syntezy. Silnik jest w nazwie, bo kontrakt jej nie niesie. Chwila jest
// w nazwie, bo odsłuch bywa powtarzany po korekcie — nadpisanie skasowałoby
// plik wcześniejszego śladu.
func nazwaNagrania(kodPanelu, silnik string) string {
	return kodPanelu + "-" + silnik + "-" + strconv.FormatInt(time.Now().UTC().UnixMilli(), 10) + ".wav"
}

// zasiegProgramowTlumaczenia składa trójkę okno-zasady-obszar zasięgu
// platformy. Żądanie modułu niesie sam panel, nie okno rozmowy, więc adresem
// jest najszerszy poziom zasięgu. Okno dostaje `ExecutionEnvCore` wprost.
func (a *adapterTlumaczenia) zasiegProgramowTlumaczenia(ctx context.Context) (session.Okno,
	session.Zasady, session.Obszar) {

	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasieg := ZasiegKonta(ctx)
	zasady := session.Zasady{}
	if a.rozstrzygaczMowy != nil {
		zasady = ZasadyIzolacji(a.rozstrzygaczMowy, zasieg)
	}
	obszar := session.Obszar{}
	if a.katalogIzolacji != nil {
		obszar = ObszarOkna(a.katalogIzolacji.Ustal(zasieg, ""), "")
	}
	return okno, zasady, obszar
}

// katalogPracySyntezy wskazuje katalog uruchomienia. Pusty jest odpowiedzią
// poprawną, nie brakiem: brama izolacji uzupełnia wtedy katalog własny
// zasięgu, a przy wyłączonym punkcie izolacji proces rusza w katalogu
// bieżącym rdzenia.
func katalogPracySyntezy(obszar session.Obszar) string {
	return strings.TrimSpace(obszar.KatalogRoboczy)
}

// ogonSyntezatora dokłada do odmowy to, co powiedział sam syntezator — zwykle
// jedyne zdanie wskazujące naprawę (np. brak głosu dla wskazanego języka).
func ogonSyntezatora(diagnostyka string) string {
	if strings.TrimSpace(diagnostyka) == "" {
		return ""
	}
	return " — syntezator powiedział: " + strings.TrimSpace(diagnostyka)
}

// bladBrakuSyntezatora znakuje zaplecze syntezy jako niedostępne: żądanie
// było poprawne, produkt nie jest zepsuty, brakuje czegoś w instalacji
// i komunikat mówi czego. Ten sam kod niesie brak pomocnika rozpoznawania
// mowy.
func bladBrakuSyntezatora(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Translate: "+powod))
}

// bladSyntezyMowy znakuje syntezę, która ruszyła i się nie udała — proces
// wystartował, ale nagranie nie powstało albo jest puste.
func bladSyntezyMowy(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Translate: "+powod))
}

// bladIzolacjiSyntezy znakuje zatrzymanie przez punkt izolacji Operatora — tak
// samo znakuje je Terminal i silnik rozpoznawania mowy.
func bladIzolacjiSyntezy(err error) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"moduł Translate: punkt izolacji zatrzymał uruchomienie syntezatora mowy: "+err.Error()))
}
