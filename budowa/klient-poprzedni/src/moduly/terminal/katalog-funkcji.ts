/**
 * Katalog funkcji modułu Terminal — inwentarz z dokumentacji projektowej wraz
 * ze stanem każdej pozycji w tej budowie.
 *
 * Po co on istnieje: moduł obiecuje zastąpić emulator terminala, multiplekser
 * sesji, klienta SSH i SFTP, harmonogram zadań, runner skryptów, menedżera
 * sekretów, konsolę szeregową i klienta kontenerów. Część tych zdolności ma
 * pokrycie w kontrakcie, część nie, a część wymaga programu, którego instalka
 * Danaco Console nie niesie. Bez wykazu Operator poznaje różnicę dopiero
 * w chwili, w której czegoś potrzebuje — a wtedy milczenie interfejsu czyta się
 * jak usterka platformy.
 *
 * Wykaz jest przepisem rozdziału „Katalog funkcji i narzędzi" dokumentacji
 * modułu, nie zbiorem pomysłów. Zmierzona liczba pozycji tego rozdziału wynosi
 * `LICZBA_POZYCJI` i tyle stoi niżej.
 *
 * Pozycja niesie trzy rzeczy: co funkcja robi, gdzie stoi w tej budowie albo
 * dlaczego jej nie ma, oraz od jakiego programu spoza instalki zależy. Zdania
 * o braku nazywają stronę braku: kontrakt bez komendy, rdzeń bez pomiaru, okno
 * poza tym złożeniem albo program, którego na maszynie rdzenia może nie być.
 */

/** Rozdział katalogu, do którego pozycja należy. */
export type RodzinaFunkcji =
  | 'Powłoki, sesje i emulator terminala'
  | 'Zdalne wykonanie i połączenia'
  | 'Skrypty, snippety i biblioteka'
  | 'Zadania, harmonogram i automatyzacja poleceniowa'
  | 'Sekrety, środowisko i uwierzytelnienie CLI'
  | 'Wynik, obserwowalność i rejestracja'
  | 'Narzędzia wiersza poleceń wbudowane w moduł';

/** Jedna pozycja katalogu funkcji. */
export interface PozycjaKatalogu {
  rodzina: RodzinaFunkcji;
  nazwa: string;
  /** Co funkcja robi — zdanie z dokumentacji modułu. */
  coRobi: string;
  /** Gdzie funkcja stoi w tej budowie albo dlaczego jej nie ma. */
  stan: string;
  /** Program spoza instalki, od którego funkcja zależy; pusty znaczy „bez zależności zewnętrznej”. */
  zaleznosc: string;
}

export const KATALOG_FUNKCJI: readonly PozycjaKatalogu[] = [
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Wielopowłokowy emulator',
    coRobi: 'Prowadzi sesję powłoki w karcie okna.',
    stan: 'Terminal Tabs otwiera kartę i wysyła polecenia. Sekwencji sterujących terminala okno nie interpretuje — wyjście stoi jako tekst wierszami, z rozróżnieniem strumienia zwykłego i diagnostycznego.',
    zaleznosc: 'powershell · cmd · bash · node · python · ssh',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Menedżer profili powłok',
    coRobi: 'Nazwane profile: typ powłoki, katalog startowy, zmienne, polecenie startowe.',
    stan: 'Kontrakt nie ma komendy przechowującej profil karty. Typ powłoki, katalog i nazwę podaje się przy otwarciu karty w Terminal Tabs; najbliższym odpowiednikiem profilu jest wpis książki hostów w Session Managerze — ten przeżywa restart rdzenia, bo leży w jego dzienniku.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Multiplekser paneli i układów',
    coRobi: 'Dzieli kartę na siatkę sąsiadujących paneli powłoki z zapisywanymi układami.',
    stan: 'Kontrakt nie zna panelu: karta jest najmniejszą jednostką, jaką rdzeń zakłada. Terminal Tabs dzieli układ WIDOKU karty na kolumny, co nie jest drugą powłoką.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Trwałe sesje nazwane',
    coRobi: 'Sesja żyje po stronie serwera po rozłączeniu klienta; podpięcie przywraca bufor i stan.',
    stan: 'Zbudowana. Session Manager i Terminal Tabs czytają wykaz kart z rdzenia (terminal.session.list), więc widać także karty otwarte przed rozłączeniem klienta — rdzeń odtwarza je przy starcie. Zamknięcie karty idzie do rdzenia (terminal.session.close) i wolno przy nim zakończyć procesy karty albo zostawić je biegnące.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Broadcast wpisu do wielu paneli',
    coRobi: 'Jedno polecenie trafia równolegle do zaznaczonych paneli i hostów.',
    stan: 'Kontrakt wiąże jedno wywołanie uruchomienia z jedną kartą. Rozgłoszenia nie ma; równoległość osiąga się osobnym uruchomieniem w każdej karcie.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Paleta poleceń',
    coRobi: 'Szybkie wywołanie czynności modułu wraz z wyszukiwaniem.',
    stan: 'Zbudowana. Paleta zna czynności wszystkich okien złożenia i ten katalog funkcji; otwiera ją skrót albo uchwyt w pasie modułu.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Bloki poleceń',
    coRobi: 'Każde polecenie i jego wynik jako osobny, składany blok z akcjami.',
    stan: 'Wynik wiąże się z procesem: Process Monitor pokazuje wyjście jednego uruchomienia osobno, a Output Console znakuje wiersze procesem. Składania i akcji na bloku w strumieniu okno nie prowadzi — granice poleceń wymagają integracji powłoki po stronie rdzenia.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Integracja powłoki',
    coRobi: 'Znaczniki katalogu roboczego i granic poleceń przekazywane automatycznie z powłoki.',
    stan: 'Brak po obu stronach: kontrakt nie przenosi tych znaczników, a rdzeń uruchamia powłokę bez czytania profilu użytkownika, więc nie ma gdzie osadzić skryptu integracyjnego.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Wyszukiwanie w buforze',
    coRobi: 'Wyszukiwanie wzorca w wyjściu, z wyrażeniami regularnymi i licznikiem trafień.',
    stan: 'Zbudowane w Output Console: wzorzec zwykły albo wyrażenie regularne, rozróżnianie wielkości liter, zawężenie do karty bieżącej.',
    zaleznosc: '',
  },
  {
    rodzina: 'Powłoki, sesje i emulator terminala',
    nazwa: 'Autouzupełnianie inline i podpowiedzi',
    coRobi: 'Uzupełnianie ścieżek, poleceń i flag oraz podpowiedź z historii.',
    stan: 'Kontrakt nie ma komendy podpowiedzi, a okno nie ma wglądu w system plików serwera. Historia poleceń karty ogranicza się do ostatniego polecenia, po które sięga ponowienie.',
    zaleznosc: '',
  },

  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Klient SSH',
    coRobi: 'Sesja zdalna jako karta powłoki.',
    stan: 'Zbudowany. Session Manager otwiera kartę powłoki zdalnej, podając adres celu zmienną środowiska karty — tak rdzeń przekazuje adres programowi ssh.',
    zaleznosc: 'ssh',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Książka hostów',
    coRobi: 'Katalog hostów z folderami, danymi połączenia i szybkim łączeniem.',
    stan: 'Zbudowana w Session Managerze nad dziennikiem rdzenia (terminal.host.save, .list, .remove), więc wpisy przeżywają odświeżenie strony i restart rdzenia. Czytanie i zapis pliku konfiguracyjnego OpenSSH zostaje jako droga wymiany z maszyną Operatora.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Przeskok przez host pośredni',
    coRobi: 'Łańcuch połączeń przez bastion.',
    stan: 'Wpis hosta niesie w kontrakcie wskazanie hosta pośredniego, ale rdzeń jeszcze go nie czyta. Do tego czasu łańcuch rozstrzyga konfiguracja OpenSSH maszyny, na której rdzeń uruchamia program ssh.',
    zaleznosc: 'ssh',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Przekierowanie portów',
    coRobi: 'Tunele lokalne, zdalne i dynamiczne wraz z panelem stanu.',
    stan: 'Zbudowana w Session Managerze (terminal.tunnel.open, .list, .close). Stan tunelu bierze się z procesu ssh, a nie z zapisu: przekierowanie, którego nie udało się założyć, wraca jako niepowodzenie wraz z powodem. Przepustowości okno nie pokazuje — rdzeń nie stoi w torze bajtów, a zero znaczyłoby nieprawdę.',
    zaleznosc: 'ssh',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Transfer plików',
    coRobi: 'Panel plików zdalnych, wysyłka i pobieranie, edycja w miejscu.',
    stan: 'Kontrakt nie ma komendy transferu, a moduł nie ma wglądu w system plików. Przenoszenie plików wykonuje się poleceniem w karcie, bez panelu.',
    zaleznosc: 'sftp · scp',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Synchronizacja katalogów',
    coRobi: 'Przyrostowa synchronizacja katalogów lokalnego i zdalnego.',
    stan: 'Bez własnego widoku: polecenie synchronizacji wykonuje się w karcie jak każde inne, a jego wynik idzie do Output Console.',
    zaleznosc: 'rsync',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Sesje na wielu hostach',
    coRobi: 'Jedna operacja na grupie hostów wraz z agregacją wyników.',
    stan: 'Kontrakt wiąże uruchomienie z jedną kartą. Grupa hostów w książce porządkuje wykaz, ale nie wykonuje jednej operacji na wielu wpisach.',
    zaleznosc: 'ssh',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Powłoka kontenera',
    coRobi: 'Karta powłoki wewnątrz kontenera wraz z wykazem kontenerów.',
    stan: 'Zbudowana. Terminal Tabs otwiera kartę powłoki kontenera, a wskazanie kontenera podaje się polem formularza karty. Rdzeń wykonuje polecenie klientem Dockera leżącym na jego maszynie; karta kontenerem nie zarządza — wykonuje polecenie w kontenerze, który już biegnie.',
    zaleznosc: 'docker',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Powłoka poda',
    coRobi: 'Powłoka w podzie, strumień dzienników, wybór kontekstu.',
    stan: 'Zbudowana. Terminal Tabs otwiera kartę powłoki poda, a nazwę poda podaje się polem formularza karty. Kontekst klastra i przestrzeń nazw bierze kubectl z konfiguracji maszyny rdzenia, o ile karta nie wskaże własnych.',
    zaleznosc: 'kubectl',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Konsola szeregowa',
    coRobi: 'Połączenie z portem szeregowym jako karta.',
    stan: 'Słownik powłok kontraktu ma już wartość dla portu szeregowego, a otwarcie karty przyjmuje urządzenie i prędkość transmisji. Brak jest po stronie rdzenia: nie otwiera on jeszcze urządzeń znakowych, więc okno tej powłoki nie oferuje.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Sesja Telnet',
    coRobi: 'Połączenie z urządzeniem sieciowym starszego typu.',
    stan: 'Zbudowana. Terminal Tabs otwiera kartę sesji Telnet z adresem celu i portem podanymi w formularzu karty. Telnet nie szyfruje ruchu — do maszyn, które mają SSH, właściwą kartą pozostaje karta zdalna.',
    zaleznosc: 'telnet',
  },
  {
    rodzina: 'Zdalne wykonanie i połączenia',
    nazwa: 'Odporność łącza',
    coRobi: 'Utrzymanie sesji przy zaniku łącza i zmianie adresu.',
    stan: 'Po stronie rdzenia: karty przeżywają rozłączenie klienta. Po stronie okna nie — bufor wyjścia ginie z połączeniem, a treść pamiętaną przez rdzeń odzyskuje się odczytem wyjścia procesu w Process Monitorze.',
    zaleznosc: '',
  },

  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Edytor skryptów',
    coRobi: 'Wieloliniowy edytor treści wraz z uruchomieniem.',
    stan: 'Zbudowany w Script Library. Treść idzie do rdzenia jako jedno polecenie, bo rdzeń podaje ją programowi powłoki pojedynczym argumentem — zapis do pliku nie jest potrzebny.',
    zaleznosc: '',
  },
  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Biblioteka skryptów',
    coRobi: 'Repozytorium skryptów z opisem, wersją i czasem ostatniego uruchomienia.',
    stan: 'Zbudowana w Script Library nad dziennikiem rdzenia (terminal.script.save, .list, .remove). Każdy zapis zakłada KOLEJNĄ wersję, więc poprzednie brzmienie treści zostaje; usunięcie pozycji zabiera ją wraz ze wszystkimi wersjami.',
    zaleznosc: '',
  },
  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Snippety poleceń',
    coRobi: 'Krótkie fragmenty do powtarzania, wstawiane z palety.',
    stan: 'Zbudowane w Script Library jako rodzaj pozycji biblioteki; uruchamiają się tą samą drogą co skrypt.',
    zaleznosc: '',
  },
  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Parametryzacja skryptów',
    coRobi: 'Deklaracja parametrów renderuje formularz wypełniany przed uruchomieniem.',
    stan: 'Zbudowana w Script Library. Deklaracja stoi w nagłówku komentarza treści, a podstawienie dzieje się w kliencie przed wysłaniem — parametr bez wartości zatrzymuje uruchomienie zamiast wejść pusty.',
    zaleznosc: '',
  },
  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Generowanie skryptu z języka naturalnego',
    coRobi: 'Model układa treść skryptu z opisu.',
    stan: 'Należy do okna rozmowy modułu, którego to złożenie nie osadza — okno rozmowy jest bytem sesji, wspólnym wszystkim modułom. Gotową treść wkleja się do edytora Script Library.',
    zaleznosc: '',
  },
  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Linter i formatowanie skryptów',
    coRobi: 'Statyczna analiza i formatowanie treści przed uruchomieniem.',
    stan: 'Zbudowana. Script Library prowadzi kontrolę wstępną w kliencie, sprawdzenie składni powłoki oraz pełną analizę na maszynie rdzenia (terminal.script.lint) wraz z formatowaniem treści. Odpowiedź mówi wprost, czy program analizy tam leży — pusty wykaz uwag przy jego braku znaczyłby fałszywie treść bez zastrzeżeń. Programów analizy instalka nie niesie.',
    zaleznosc: 'shellcheck · shfmt · PSScriptAnalyzer · ruff',
  },
  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Aliasy i funkcje powłoki',
    coRobi: 'Zarządzanie aliasami z interfejsu, synchronizowane między sesjami.',
    stan: 'Brak drogi: rdzeń uruchamia powłokę bez czytania profilu użytkownika, więc alias zapisany w profilu i tak nie obowiązywałby w karcie.',
    zaleznosc: '',
  },
  {
    rodzina: 'Skrypty, snippety i biblioteka',
    nazwa: 'Pliki konfiguracyjne powłoki',
    coRobi: 'Wersjonowany zestaw plików konfiguracyjnych, przenośny między hostami.',
    stan: 'Kontrakt nie ma komendy zapisu pliku na serwerze. Wersjonowanie repozytorium prowadzi moduł Developer.',
    zaleznosc: 'git',
  },

  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Runner zadań',
    coRobi: 'Wykrywa i uruchamia zadania z manifestów projektu.',
    stan: 'Zbudowany w Task & Schedule. Manifest wypisuje polecenie powłoki, a rozbiór jego treści dzieje się w kliencie. Odczyt pliku ma już komendę kontraktu — po jej wpięciu wypisywanie plikiem powłoki zostanie drogą zapasową dla powłok, w których odczyt nie zadziała.',
    zaleznosc: 'npm · make · task · just · program wypisujący plik',
  },
  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Harmonogram zadań',
    coRobi: 'Uruchomienia cykliczne wraz z kalendarzem najbliższych terminów.',
    stan: 'Zbudowany w Task & Schedule w obie strony: odczyt komendą harmonogramu, zapis dwoma krokami rodziny automatyk — zapisem automatyki o jednym kroku wołającym polecenie powłoki i dokładeniem jej wyrażenia cron. Rodzina komend zostaje jedna, bo kontrakt wiąże cykliczność z automatyką, a nie z kartą powłoki. Zapis jest planem wraz z terminem, nie budzikiem: rdzeń nie ma czym odpalić automatyki samodzielnie, więc uruchomienie prowadzi Operator z okien modułu Automations.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Wyzwalacze plikowe',
    coRobi: 'Uruchomienie polecenia przy zmianie plików i katalogów.',
    stan: 'Zbudowana w oknie Task & Schedule (terminal.watch.start, .stop, .list). Obserwacja wiąże się z KARTĄ powłoki, więc wyzwolone polecenie wykonuje się w jej katalogu i jej powłoką, tą samą drogą co polecenie wydane ręcznie. Wykaz niesie licznik wyzwoleń i powód niepowodzenia.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Potoki i sekwencje poleceń',
    coRobi: 'Łańcuch kroków z warunkami powodzenia i porażki.',
    stan: 'Zbudowany w Task & Schedule jako sekwencjonowanie po stronie klienta: okno wysyła krok, czeka na jego kod wyjścia i dopiero wtedy decyduje o następnym. Gałęzi równoległych i zmiennych potoku nie ma — to jest silnik rdzenia, którego kontrakt Terminalowi nie udostępnia.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Kolejka zadań i przebiegi',
    coRobi: 'Kolejkowanie oraz historia przebiegów ze statusem i wynikiem.',
    stan: 'Zbudowane w Task & Schedule. Kolejki silnika pętli obsługujące okno czyta i posuwa sześcioma działaniami słownika kontraktu; historia przebiegów liczy uruchomienia wydane z tego okna. Kolejki nie zakłada — komenda zakładająca żąda identyfikatora karty sesji, którego moduł nie zna.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Powiadomienia o zakończeniu',
    coRobi: 'Alert po zakończeniu długiego zadania oraz przy błędzie.',
    stan: 'Kontrakt nie ma komendy progu powiadomienia dla terminala. Czas trwania i kod wyjścia widać w Process Monitorze; kanał powiadomień jest własnością platformy, nie modułu.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Wyzwalacze zdarzeniowe',
    coRobi: 'Uruchomienie polecenia na zdarzenie zewnętrzne albo koniec innego zadania.',
    stan: 'Realizowane jawnym powiązaniem z modułem Automations, zgodnie z dokumentacją modułu. Terminal nie zakłada takiego powiązania sam.',
    zaleznosc: '',
  },
  {
    rodzina: 'Zadania, harmonogram i automatyzacja poleceniowa',
    nazwa: 'Bramy wykonania równoległego',
    coRobi: 'Uruchomienie wielu zadań z limitem współbieżności i agregacją wyników.',
    stan: 'Potok w Task & Schedule jest sekwencyjny. Limit współbieżności należy do silnika kolejek rdzenia i kontrakt nie oddaje go Terminalowi.',
    zaleznosc: '',
  },

  {
    rodzina: 'Sekrety, środowisko i uwierzytelnienie CLI',
    nazwa: 'Menedżer zmiennych środowiskowych',
    coRobi: 'Edytor zmiennych per sesja, projekt i host wraz z podglądem zestawu obowiązującego.',
    stan: 'Zmienne wchodzą wyłącznie przy otwarciu karty i tylko tą drogą — kontrakt nie ma komendy ich odczytu ani zmiany w karcie już otwartej. Warstw nadpisań moduł nie prowadzi.',
    zaleznosc: '',
  },
  {
    rodzina: 'Sekrety, środowisko i uwierzytelnienie CLI',
    nazwa: 'Ładowarka plików środowiska',
    coRobi: 'Automatyczne wczytanie pliku zmiennych po wejściu do katalogu.',
    stan: 'Brak: okno nie wie, kiedy zmienia się katalog roboczy, bo rdzeń nie przekazuje znaczników integracji powłoki, a każde polecenie jest osobnym procesem.',
    zaleznosc: '',
  },
  {
    rodzina: 'Sekrety, środowisko i uwierzytelnienie CLI',
    nazwa: 'Magazyn sekretów sesji',
    coRobi: 'Wstrzykiwanie sekretu do polecenia bez zapisu w historii.',
    stan: 'Kontrakt ma w zmiennej środowiska odwołanie do sejfu, ale karta terminala nie ma czytnika sejfu i rdzeń odmawia takiej zmiennej wprost, zamiast pominąć ją po cichu. Sekret podaje się więc wartością jawną albo wcale.',
    zaleznosc: '',
  },
  {
    rodzina: 'Sekrety, środowisko i uwierzytelnienie CLI',
    nazwa: 'Maskowanie wartości wrażliwych',
    coRobi: 'Zaciemnianie tokenów i kluczy w wyniku według wzorców.',
    stan: 'Brak. Wyjście idzie do okna w postaci, w jakiej wypisała je powłoka; maskowanie po stronie widoku dawałoby złudzenie ochrony, bo treść jawna i tak przeszła przez łącze i przez rejestr rdzenia.',
    zaleznosc: '',
  },
  {
    rodzina: 'Sekrety, środowisko i uwierzytelnienie CLI',
    nazwa: 'Menedżer kluczy SSH',
    coRobi: 'Generowanie, import i przypisanie klucza do hosta.',
    stan: 'Komendy wykazu kluczy stoją w kontrakcie, a część tajna nie przechodzi w nich przez łącze: klucz wytwarza się na maszynie rdzenia, a wciąga do wykazu wskazaniem ścieżki. Okno ich jeszcze nie wywołuje, więc do tego czasu klucze rozstrzyga wyłącznie konfiguracja maszyny, na której rdzeń uruchamia program ssh.',
    zaleznosc: 'ssh-keygen',
  },
  {
    rodzina: 'Sekrety, środowisko i uwierzytelnienie CLI',
    nazwa: 'Profile poświadczeń narzędzi wiersza poleceń',
    coRobi: 'Przełączanie profilu poświadczeń per karta.',
    stan: 'Osiągalne wyłącznie zmienną środowiska podaną przy otwarciu karty; osobnego widoku profili moduł nie ma.',
    zaleznosc: 'programy właściwe dostawcy chmury',
  },

  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Rejestracja sesji',
    coRobi: 'Nagranie przebiegu jako plik odtwarzalny z osią czasu.',
    stan: 'Kontrakt nie ma komendy rejestracji. Output Console eksportuje treść wyjścia, ale bez osi czasu odtwarzalnej poza platformą.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Odtwarzanie',
    coRobi: 'Przewijanie przebiegu krok po kroku.',
    stan: 'Zbudowane w Output Console jako przesuwanie granicy widocznych wierszy bufora widoku. To odtwarzanie tego, co okno zebrało w tym połączeniu, a nie nagrania z rdzenia.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Eksport transkryptu i dziennika',
    coRobi: 'Zapis karty albo zakresu wyjścia do pliku.',
    stan: 'Zbudowany w Terminal Tabs i Output Console jako zapis tekstu. Postaci z zachowaniem barw okno nie tworzy, bo nie interpretuje sekwencji sterujących terminala.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Parser strukturalny wyniku',
    coRobi: 'Rozpoznanie struktury w wyniku i prezentacja jako tabela albo drzewo.',
    stan: 'Poza tym złożeniem jako widok wyjścia. Rozbiór treści prowadzi Task & Schedule wyłącznie dla manifestów projektu, bo tam struktura jest znana z góry.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Wykresy zasobów procesu',
    coRobi: 'Przebieg zużycia procesora, pamięci i operacji wejścia-wyjścia.',
    stan: 'Kontrakt ma pola zużycia, ale rdzeń ich nie mierzy i nie przysyła. Process Monitor ich nie pokazuje — wpisanie zera znaczyłoby „nic nie zużywa", co byłoby nieprawdą.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Rejestr i drzewo procesów',
    coRobi: 'Pełny rejestr uruchomień z hierarchią potomków, filtrami i akcjami.',
    stan: 'Rejestr, filtry, grupowanie i zakończenie procesu zbudowane w Process Monitorze. Drzewa potomków nie ma: kontrakt niesie wyłącznie identyfikator procesu nadrzędnego, a wykaz obejmuje procesy rejestru rdzenia, nie wszystkie procesy maszyny.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Wykrywanie wzorców błędów',
    coRobi: 'Oznaczenie wiersza błędu plakietką wraz z czynnością.',
    stan: 'Częściowo: wiersz wyjścia diagnostycznego jest odróżniony, a proces zakończony błędem ma własny stan w wykazie. Rozpoznawania wzorców w treści okno nie prowadzi.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Różnicowanie wyników',
    coRobi: 'Porównanie wyjścia dwóch uruchomień obok siebie.',
    stan: 'Brak. Output Console pokazuje dwa widoki obok siebie po podziale, ale nie zestawia ich znak po znaku.',
    zaleznosc: '',
  },
  {
    rodzina: 'Wynik, obserwowalność i rejestracja',
    nazwa: 'Dziennik pętli wykonawczej',
    coRobi: 'Zapis zlecenia, zadań i komunikatów sterujących pętli.',
    stan: 'Należy do okna pętli wykonawczej, które jest bytem sesji wspólnym modułom i nie wchodzi do tego złożenia.',
    zaleznosc: '',
  },

  {
    rodzina: 'Narzędzia wiersza poleceń wbudowane w moduł',
    nazwa: 'Klient żądań sieciowych w karcie',
    coRobi: 'Wykonanie i podgląd żądania wraz z historią.',
    stan: 'Bez własnego widoku: żądanie wykonuje się poleceniem w karcie, a odpowiedź idzie do Output Console jak każde inne wyjście.',
    zaleznosc: 'curl',
  },
  {
    rodzina: 'Narzędzia wiersza poleceń wbudowane w moduł',
    nazwa: 'Podgląd i edycja plików strukturalnych',
    coRobi: 'Podgląd i zmiana treści plików danych wraz z kontrolą poprawności.',
    stan: 'Odczyt pliku ma komendę kontraktu, zapisu nie ma wcale. Podgląd wykonuje się dziś poleceniem wypisującym plik, a edycję prowadzi moduł Developer.',
    zaleznosc: '',
  },
  {
    rodzina: 'Narzędzia wiersza poleceń wbudowane w moduł',
    nazwa: 'Pulpit menedżerów pakietów',
    coRobi: 'Jednolity widok instalacji, aktualizacji i przeglądu zależności.',
    stan: 'Bez własnego widoku: polecenia menedżera pakietów wykonują się w karcie, a zadania z manifestu projektu wykrywa Task & Schedule.',
    zaleznosc: 'npm · pip · cargo · go · winget · brew',
  },
  {
    rodzina: 'Narzędzia wiersza poleceń wbudowane w moduł',
    nazwa: 'Konwertery i kodowanie',
    coRobi: 'Przekształcenia postaci danych i skróty kryptograficzne.',
    stan: 'Bez własnego widoku: przekształcenie wykonuje się poleceniem w karcie.',
    zaleznosc: 'programy właściwe przekształceniu',
  },
  {
    rodzina: 'Narzędzia wiersza poleceń wbudowane w moduł',
    nazwa: 'Zmienne robocze i wyrażenia',
    coRobi: 'Wynik poprzedniego polecenia dostępny jako wartość w następnym.',
    stan: 'Brak: każde polecenie jest odrębnym procesem powłoki, więc stan nie przechodzi między uruchomieniami. Ciąg zależnych kroków zapisuje się jako jeden skrypt w Script Library albo jako potok w Task & Schedule.',
    zaleznosc: '',
  },
  {
    rodzina: 'Narzędzia wiersza poleceń wbudowane w moduł',
    nazwa: 'Menedżer plików dwupanelowy',
    coRobi: 'Kopiowanie, przenoszenie i podgląd między dwoma katalogami.',
    stan: 'Kontrakt ma wyłącznie odczyt pojedynczego pliku — ani przeglądania katalogu, ani kopiowania, przenoszenia i usuwania. Operacje na plikach wykonuje się poleceniami w karcie.',
    zaleznosc: '',
  },
];

/** Zmierzona liczba pozycji katalogu — liczona z wykazu, nie wpisana. */
export const LICZBA_POZYCJI = KATALOG_FUNKCJI.length;
