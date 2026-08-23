// Odpowiedzialność pliku: silnik syntezy mowy modułu Translate — przeprowadzenie
// wybranego syntezatora do pliku dźwiękowego: plik tekstu, wiersz wywołania,
// sprawdzenie skutku i przekład odmów na kody kontraktu. Metody stoją na
// wspólnym `*adapterTlumaczenia` (typ i przedrostki deklaruje
// `adapter_modul_tlumaczenie.go`); `SyntezujMowe` — komenda, która tego silnika
// używa — mieszka w `adapter_modul_tlumaczenie_mowa.go`, a odpowiedź na pytanie
// „kogo w ogóle wołamy i jakim głosem" — w
// `adapter_modul_tlumaczenie_mowa_glos.go`. Nagłówek niesie rozstrzygnięcia
// wspólne obu plikom, żeby nie rozjechały się przy poprawce.
//
// Uruchomienia procesu tu nie ma: sekwencja startu (UruchomProces →
// PrzejmijDrzewo → pompy strumieni → oczekiwanie → ubicie drzewa) mieszka
// w pakiecie `zewnetrzne`, wspólnym z Terminalem, Developerem
// i `mowa/uruchomienie.go`. Ten plik woła `zewnetrzne.Wolaj` i o proces nie
// pyta.
//
// Dźwięk nie opuszcza maszyny Operatora: nie ma tu klienta HTTP ani adresu,
// a oba syntezatory są programami lokalnymi.
//
// Silniki są dwa, w tej kolejności pierwszeństwa:
//  1. `piper` — synteza neuronowa, głos, którego da się słuchać. Wchodzi
//     pierwszy, gdy stoi binarium oraz jest głos dla języka panelu. Głos leży
//     na maszynie jak każde inne binarium arsenału, więc jest zależnością
//     środowiska, a nie stanem produktu — rdzeń go nie pobiera, nie
//     wersjonuje i nie sprząta.
//  2. `espeak-ng` — synteza formantowa, głos mechaniczny, ale program jest
//     jednym plikiem bez stanu i bez modeli. Droga zapasowa: wchodzi, gdy
//     pipera nie ma albo nie ma dla tego języka głosu.
//
// Operator ma wiedzieć, którym silnikiem słucha: głos zapasowy brzmi inaczej
// niż dobry i nie ma być mylony z usterką nagrania. Kontrakt
// (`TranslateSpeechSynthesizeResponse`) niesie same `panelId` i `path`, bez
// pola na nazwę silnika, więc nazwa silnika idzie w nazwę pliku:
// `pan-…-piper-….wav` albo `pan-…-espeak-ng-….wav`. `path` jest polem kontraktu
// i niesie prawdę o tym, co powstało; ta sama nazwa ląduje w kolumnie
// `nagranie_odnosnik`, więc ślad również mówi, kto czytał. Gdy kontrakt dostanie
// pole `engine`, wystarczy je wypełnić wartością, którą ten plik już zna.
//
// Ścieżka programu i głosu bierze się z dwóch źródeł, w tej kolejności:
//
//   - zmienna środowiska (`DANACO_PIPER`, `DANACO_PIPER_GLOSY`,
//     `DANACO_ESPEAK`) — pierwszeństwo. Arsenał jest instalowany poza
//     produktem i bywa na każdej maszynie w innym miejscu; zmienna jest jedynym
//     wskazaniem, które działa bez przebudowy i bez migracji;
//   - wykrycie w miejscach typowych — nazwa goła w PATH, a przy jej braku
//     katalog arsenału (`/opt/danaco-arsenal/…`), dzięki czemu instalacja
//     typowa działa bez ustawiania czegokolwiek.
//
// Katalog ustawień produktu (`konfig`, wzór `mowa_model`) źródłem nie jest:
// wpisywałby położenie cudzego binarium do stanu produktu, a to jest fakt
// maszyny, nie nastawa Operatora. Gdy arsenał dostanie własną rodzinę ustawień,
// ten plik ma czytać ją, a nie dokładać trzecią drogę.
//
// Brak głosu to inna odmowa niż brak binarium: program stoi, więc `Stoi` mówi
// „jest", a czynność i tak nie wyjdzie. Rozróżnienie widać w treści odmowy, bo
// naprawy są różne — „zainstalować pipera" kontra „dołożyć plik głosu dla
// języka X do katalogu Y". Gdy zawiodą oba silniki, odmowa wymienia obie
// przyczyny osobno.
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
	// ładuje model sześćdziesięciomegabajtowy przy każdym wywołaniu i na słabszej
	// maszynie potrzebuje na to kilkunastu sekund — minuta zostawia zapas, a
	// zarazem nie pozwala programowi, który utknął, trzymać żądania bez końca
	// (`zewnetrzne.Wolaj` granicy wymaga).
	granicaSyntezy = 60 * time.Second
)

// ZSynteza wpina silnik syntezy mowy: uruchamiacz procesów (jedyna droga startu
// procesu w drzewie), dwa źródła izolacji — te same, którymi jadą
// Terminal, Developer i silnik rozpoznawania mowy — oraz katalog danych rdzenia,
// pod którym lądują nagrania.
//
// Zależność jest opcjonalna w tym sensie, że jej brak nie psuje pozostałych
// komend modułu; psuje wyłącznie `speech.synthesize`, która wtedy odmawia
// nazywając brak, zamiast oddać pustą ścieżkę udającą nagranie.
func (a *adapterTlumaczenia) ZSynteza(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy, katalogDanych string) *adapterTlumaczenia {

	a.uruchamiacz = uruchamiacz
	a.rozstrzygaczMowy = rozstrzygacz
	a.katalogIzolacji = katalog
	a.katalogDanych = katalogDanych
	return a
}

// zsyntezujDoPliku zamienia tekst na plik WAV i oddaje jego ścieżkę.
//
// Kolejność jest rozmyślna: najpierw pada pytanie, czym syntezować, i dopiero potem
// zakładamy cokolwiek na dysku. Katalog nagrań założony pod plik, który nigdy
// nie powstał, byłby śmieciem po odmowie.
func (a *adapterTlumaczenia) zsyntezujDoPliku(ctx context.Context,
	kodPanelu, jezyk, tekst string) (string, error) {

	if a.uruchamiacz == nil {
		return "", bladBrakuSyntezatora("rdzeń nie ma uruchamiacza procesów — nie ma czym wystartować " +
			"syntezatora mowy; naprawa: podpiąć warstwę kanału (injection) przy składaniu rdzenia")
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

	// Tekst idzie plikiem, nie argumentem. Piper czyta treść z pliku wskazanego
	// przełącznikiem `-i` albo ze standardowego wejścia, a port
	// `session.Uruchamiacz` wejścia procesu nie wystawia — więc plik jest jedyną
	// drogą, która nie wymaga rozszerzania portu. Znika przy tym granica
	// długości argumentu, o którą rozbija się treść dłuższego panelu. Espeak
	// dostaje ten sam plik przełącznikiem `-f`, żeby obie drogi różniły się
	// wyłącznie wierszem wywołania.
	sciezkaTekstu := sciezka + ".txt"
	if err := os.WriteFile(sciezkaTekstu, []byte(tekst), 0o600); err != nil {
		return "", bladSyntezyMowy("nie można zapisać tekstu do syntezy pod " + sciezkaTekstu + ": " + err.Error())
	}
	// Plik tekstowy jest rusztowaniem, nie wynikiem — znika niezależnie od tego,
	// czy synteza się udała. Zostawiony leżałby obok nagrania i wyglądał jak
	// część wyniku.
	defer func() { _ = os.Remove(sciezkaTekstu) }()

	okno, zasady, obszar := a.zasiegSyntezy()
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar,
		wybor.narzedzie, wybor.argumenty(sciezkaTekstu, sciezka), katalogPracySyntezy(obszar), granicaSyntezy)
	if err != nil {
		// Naruszenie izolacji znakujemy osobno — punkt izolacji Operatora to nie
		// jest usterka rdzenia (tak samo znakuje je Terminal i silnik mowy).
		if errors.Is(err, session.ErrIzolacja) {
			return "", bladIzolacjiSyntezy(err)
		}
		return "", bladSyntezyMowy(err.Error() + ogonSyntezatora(wynik.Diagnostyka))
	}

	// Plik ma istnieć i mieć rozmiar. Syntezator kończący się powodzeniem, ale
	// nie zostawiający nagrania, byłby odmową udającą sukces. Sprawdzany jest
	// więc skutek, a nie kod wyjścia procesu; piper potrafi wyjść zerem, gdy
	// tekst zwęzi się do samych znaków niewymawialnych.
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

// wyborSyntezatora niesie rozstrzygnięcie „czym czytamy": nazwę silnika (idzie
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
	// Espeak dostaje język w `-v`; głosu nie zgadujemy i nie mamy własnej tabelki
	// nazw — `espeak-ng` rozumie i kody (`pl`, `pl-PL`), i nazwy angielskie
	// (`polish`), więc tabelka w rdzeniu byłaby drugą prawdą o wykazie głosów.
	// Języka, którego syntezator nie zna, nie podmieniamy na domyślny: nagranie
	// polskiego zdania przeczytane po angielsku byłoby atrapą bez słowa
	// ostrzeżenia.
	return []string{"-v", w.glos, "-w", sciezkaNagrania, "-f", sciezkaTekstu}
}

// katalogNagran zakłada (gdy trzeba) katalog na nagrania pod katalogiem danych
// rdzenia. Pusty katalog danych jest odmową, nie powodem do wybrania czegoś
// z własnej głowy: nagranie zapisane w katalogu bieżącym procesu wylądowałoby
// tam, gdzie Operator go nie szuka, i nie zniknęłoby razem z resztą danych
// rdzenia.
func (a *adapterTlumaczenia) katalogNagran() (string, error) {
	podstawa := strings.TrimSpace(a.katalogDanych)
	if podstawa == "" {
		return "", bladBrakuSyntezatora("rdzeń nie zna katalogu danych — nie ma gdzie zapisać nagrania; " +
			"naprawa: wskazać katalog danych przełącznikiem -dane albo zmienną DANACO_KATALOG_DANYCH")
	}
	katalog := filepath.Join(podstawa, katalogSyntezyMowy)
	if err := os.MkdirAll(katalog, 0o755); err != nil {
		return "", bladSyntezyMowy("nie można założyć katalogu nagrań " + katalog + ": " + err.Error())
	}
	return katalog, nil
}

// nazwaNagrania składa nazwę pliku z kodu panelu, nazwy silnika i chwili
// syntezy. Silnik jest w nazwie, bo kontrakt nie ma pola na niego, a Operator
// ma wiedzieć, czy słucha głosu dobrego, czy zapasowego (nagłówek pliku).
// Chwila jest w nazwie, bo odsłuch bywa powtarzany po korekcie — nadpisywanie
// kasowałoby plik, do którego może już prowadzić wcześniejszy ślad.
func nazwaNagrania(kodPanelu, silnik string) string {
	return kodPanelu + "-" + silnik + "-" + strconv.FormatInt(time.Now().UTC().UnixMilli(), 10) + ".wav"
}

// zasiegSyntezy składa trójkę okno–zasady–obszar dla zasięgu platformy, tą samą
// drogą i z tego samego powodu, co `adapterMowy.zasiegPlatformy`: żądanie
// `speech.synthesize` niesie sam panel, a nie okno rozmowy, więc adresem jest
// najszerszy z ośmiu poziomów zasięgu, a nie podstawione po cichu
// puste struktury znaczące „izolacja wyłączona”.
//
// Okno dostaje `ExecutionEnvCore` wprost: nagranie ma powstać na dysku rdzenia,
// bo to rdzeń odda potem jego ścieżkę w odpowiedzi.
func (a *adapterTlumaczenia) zasiegSyntezy() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygaczMowy != nil {
		zasady = ZasadyIzolacji(a.rozstrzygaczMowy, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalogIzolacji != nil {
		obszar = ObszarOkna(a.katalogIzolacji.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// katalogPracySyntezy wskazuje katalog uruchomienia. Pusty jest odpowiedzią
// poprawną, nie brakiem: brama izolacji uzupełnia wtedy katalog własny zasięgu,
// a przy wyłączonym punkcie izolacji proces rusza w katalogu bieżącym rdzenia.
// Wpisanie tu czegoś z własnej głowy odbierałoby bramie rozstrzygnięcie, które
// należy do niej (wzór `mowa/silnik.go`, `katalogPracy`).
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

// bladBrakuSyntezatora znakuje zaplecze syntezy jako niedostępne: żądanie było
// poprawne, produkt nie jest zepsuty, brakuje czegoś w instalacji i komunikat
// mówi czego. Ten sam kod i to samo uzasadnienie, co przy braku pomocnika
// rozpoznawania (`bladSilnikaMowy`, `adapter_modul_mowa.go`) — dwie reguły dla
// jednego rodzaju braku byłyby rozjazdem.
func bladBrakuSyntezatora(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Translate: "+powod))
}

// bladSyntezyMowy znakuje syntezę, która ruszyła i się nie udała.
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
