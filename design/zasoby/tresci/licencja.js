/* Treść umowy licencyjnej jest fragmentem wstawianym do pola przewijanego w kroku instalacji i nie jest samodzielnym dokumentem do edycji.
   TREŚĆ UMOWY LICENCYJNEJ — wynik `zbuduj-licencje.py` ze źródła
   `zasoby/tresci/licencja-2.1.md`. Fragment, nie dokument: wstawia go
   składnik `dokument` do pola przewijanego w kroku 2.

   Nie edytuj tu niczego — zmiana wchodzi w źródle i przechodzi przez
   generator. Treść stoi w skrypcie, a nie w osobnym pliku do pobrania, bo
   okno bywa otwierane wprost z dysku, a przeglądarka blokuje wtedy
   pobieranie plików towarzyszących.
   ============================================================================ */
window.DanacoTresci = window.DanacoTresci || {};
window.DanacoTresci.licencja = `
<h2 id="1-postanowienia-wstępne">1. Postanowienia wstępne</h2>
<h3 id="11-charakter-dokumentu">1.1 Charakter dokumentu</h3>
<p>Niniejszy dokument („<strong>Licencja</strong>”) określa warunki, na jakich Danaco Holding Group Sp. z o.o. udostępnia oprogramowanie Danaco Console — Platforma AI Workspace OS („<strong>Oprogramowanie</strong>”) do korzystania. Licencja reguluje zakres uprawnień przyznanych korzystającemu, granice dozwolonego użycia, zasady dotyczące komponentów osób trzecich wbudowanych w Oprogramowanie, sposób postępowania z danymi wytworzonymi w toku korzystania, a także zasady odpowiedzialności Licencjodawcy odpowiadające rzeczywistemu, potwierdzonemu stanowi wykonania produktu.</p>
<p>Oprogramowanie jest <strong>produktem własnościowym</strong>. Nie jest oprogramowaniem otwartoźródłowym, nie jest oprogramowaniem darmowym w rozumieniu swobody redystrybucji i nie podlega żadnej z powszechnie stosowanych licencji wolnego oprogramowania. Fakt, że Oprogramowanie zawiera komponenty osób trzecich rozpowszechniane na licencjach otwartych (rozdz. 7), nie zmienia własnościowego charakteru samego Oprogramowania ani nie rozciąga warunków tych licencji na kod autorski Licencjodawcy — z zastrzeżeniem obowiązków wynikających wprost z licencji poszczególnych komponentów, opisanych w rozdz. 7.</p>
<p>Licencja odnosi się do wersji <strong>v2.0 o statusie Deweloperskim</strong>. Status ten nie jest formułą marketingową ani zastrzeżeniem ostrożnościowym: wynika z przyjętej decyzji architektonicznej („v2.0 do pierwszej publikacji; zakaz wersjonowania w trakcie budowy; status Deweloperski”) i odpowiada faktycznemu stanowi wykonania produktu, ustalonemu w drodze pełnego audytu repozytorium przeprowadzonego pierwotnie w stanie sprzed prac naprawczych, a następnie zweryfikowanemu ponownie i uzupełnionemu po zamknięciu prac naprawczych do stanu bieżącego stanu repozytorium tej samej gałęzi — rewizji, z której zbudowano i uruchomiono doręczany Egzemplarz (Załącznik C.1, C.3 pkt 4). Konsekwencje statusu Deweloperskiego dla rękojmi i odpowiedzialności opisuje rozdz. 10; konsekwencje dla zakresu funkcjonalnego, do którego Operator ma prawo mieć zaufanie, opisuje rozdz. 3.1 oraz rozdz. 10.3.</p>
<h3 id="12-strony-stosunku-licencyjnego">1.2 Strony stosunku licencyjnego</h3>
<p>Stronami stosunku licencyjnego są:</p>
<ol type="1">
<li><strong>Licencjodawca</strong> — Danaco Holding Group Sp. z o.o., podmiot, któremu przysługują autorskie prawa majątkowe do Oprogramowania, występujący w dokumentacji projektu również jako <strong>Producent</strong>;</li>
<li><strong>Licencjobiorca</strong> — podmiot, któremu Licencjodawca udostępnił egzemplarz Oprogramowania na warunkach niniejszej Licencji.</li>
</ol>
<p>Dokumentacja techniczna produktu posługuje się konsekwentnie pojęciem <strong>Operator</strong> na oznaczenie osoby, która faktycznie obsługuje Oprogramowanie: konfiguruje kanały modelu, konta, tożsamość modelu, punkty dostępu i nadania, prowadzi sesje i okna komunikacji oraz odpowiada za wynik pracy wykonanej z udziałem modelu. Pojęcia <strong>Licencjobiorca</strong> i <strong>Operator</strong> nie są tożsame: Licencjobiorcą może być osoba prawna, podczas gdy Operatorem jest zawsze osoba fizyczna zasiadająca przy Oprogramowaniu. Relacja między tymi rolami — a w szczególności to, czy Licencjobiorcą jest osoba fizyczna korzystająca z Oprogramowania, czy podmiot, na którego rzecz Operatorzy pracują — pozostaje kwestią modelu licencjonowania.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> Licencjobiorcą jest <strong>podmiot</strong> — osoba prawna, jednostka organizacyjna albo osoba fizyczna, w tym prowadząca działalność gospodarczą — któremu Licencjodawca udostępnił Egzemplarz, <strong>z prawem wskazania Operatorów</strong> w liczbie wynikającej z rozdz. 4.3. Licencjobiorca odpowiada za działania i zaniechania wskazanych Operatorów jak za własne. Jeżeli Egzemplarz udostępniono osobie fizycznej bez wskazania podmiotu, Licencjobiorcą jest ta osoba, a Operatorem — ona sama.</p>
<h3 id="13-zasada-zgodności-ze-stanem-faktycznym-produktu">1.3 Zasada zgodności ze stanem faktycznym produktu</h3>
<p>Licencja opisuje produkt rzeczywisty, nie zamierzony. Wszędzie tam, gdzie konieczne jest odniesienie do funkcji Oprogramowania — w szczególności w rozdz. 3 (przedmiot licencji), rozdz. 9 (dane Operatora) i rozdz. 10 (odpowiedzialność) — przywołany jest stan wykonania potwierdzony audytem repozytorium, a nie stan zaprojektowany w dokumentacji koncepcyjnej.</p>
<p>Reguła ta ma znaczenie prawne, nie redakcyjne: zakres, w jakim Licencjodawca mógłby ponosić odpowiedzialność za niezgodność Oprogramowania z opisem, wyznacza opis zawarty w niniejszej Licencji i w dokumentacji produktu przekazanej Licencjobiorcy — to jest w zamkniętym zestawie dokumentów wymienionym w rozdz. 1.4 pkt 3 i rozdz. 3.5 — a nie materiały koncepcyjne, prezentacyjne ani rejestry decyzji architektonicznych opisujące zamierzenia. Funkcje przewidziane koncepcją, lecz w wersji v2.0 niezaimplementowane albo niepodłączone do interfejsu, oznaczono w treści Licencji jednoznacznie — nie stanowią one przedmiotu świadczenia.</p>
<h3 id="14-hierarchia-dokumentów-projektu">1.4 Hierarchia dokumentów projektu</h3>
<p>W razie rozbieżności między dokumentami dotyczącymi Oprogramowania stosuje się następującą kolejność pierwszeństwa:</p>
<ol type="1">
<li>niniejsza Licencja — w zakresie uprawnień, ograniczeń, odpowiedzialności i zasad dotyczących komponentów osób trzecich;</li>
<li>odrębna umowa zawarta na piśmie między Licencjodawcą a Licencjobiorcą, jeżeli została zawarta i jeżeli wprost modyfikuje postanowienia Licencji;</li>
<li><strong>dokumentacja produktu</strong> przekazana Licencjobiorcy wraz z Egzemplarzem, obejmująca — poza niniejszym plikiem <code>LICENSE.md</code> — trzy dokumenty oznaczone nazwami plików:
<ul>
<li><strong><code>README.md</code></strong> — karta techniczno-produktowa Danaco Console: koncepcja, architektura trzech warstw, stos technologiczny, model katalogów oraz zakres funkcjonalny wraz z oznaczeniem stanu wykonania poszczególnych funkcji;</li>
<li><strong><code>INSTRUKCJA-UZYTKOWANIA.md</code></strong> — opis pracy z gotową Instalacją: strona główna, Środowiska, Okno komunikacji, okna równoległe, okno konfiguracji, punkty dostępu, konta i tożsamość modelu;</li>
<li><strong><code>INSTALACJA-I-KONFIGURACJA.md</code></strong> — procedura instalacji, budowania warstw i doprowadzenia świeżej Instalacji do stanu zdatnego do pracy, wykaz zmiennych środowiska i argumentów wywołania oraz postępowanie w razie niepowodzeń;</li>
</ul></li>
<li>dokumentacja projektowa i wewnętrzna Licencjodawcy — <strong>nieprzekazywana Licencjobiorcy i nieskładająca się na treść zobowiązania</strong> — w szczególności zapis decyzji architektonicznych, wykaz luk, brief projektowy, specyfikacje fal budowy, pakiety warstwy wizualnej oraz opracowania klasy Specyfikacja docelowa katalogów <code>architektura/</code>, <code>moduly/</code>, <code>interfejs-uzytkownika/</code>, <code>srodowiska/</code>, <code>specyfikacje/</code> i <code>funkcje-globalne/</code> — dokumentacja techniczna towarzysząca kodowi repozytorium budowy, odrębna od czterech dokumentów produktu wymienionych w pkt 3.</li>
</ol>
<p>Wykaz z pkt 3 jest zamknięty: dokumentem produktu jest wyłącznie plik wprost w nim wymieniony. Dokument nienależący do tego wykazu — choćby dotyczył Oprogramowania i pochodził od Licencjodawcy — nie tworzy zobowiązania co do zakresu funkcjonalnego (rozdz. 14.5).</p>
<h3 id="15-zastrzeżenie-o-przeglądzie-prawnym">1.5 Zastrzeżenie o przeglądzie prawnym</h3>
<p>Niniejszy dokument został przygotowany na podstawie dokumentacji technicznej i stanu faktycznego repozytorium, a kwestie handlowo-prawne pozostawione wcześniej do decyzji Operatora zostały rozstrzygnięte i wpisane do treści (Załącznik B). W tej postaci dokument stanowi <strong>wersję ustaloną</strong>, zdatną do dołączenia do dystrybucji Wydań o statusie Deweloperskim.</p>
<p><strong>Przed rozpoczęciem dystrybucji komercyjnej — w szczególności przed wprowadzeniem odpłatności (rozdz. 4.5) — zaleca się przegląd dokumentu przez radcę prawnego lub adwokata.</strong> Przegląd powinien objąć w szczególności:</p>
<ul>
<li>skuteczność wyłączenia rękojmi i ograniczenia odpowiedzialności wobec konsumentów oraz wobec przedsiębiorców, w świetle prawa właściwego wybranego zgodnie z rozdz. 14.1;</li>
<li>zgodność klauzul dotyczących danych osobowych z rozporządzeniem (UE) 2016/679 (RODO) oraz ustalenie ról administratora i podmiotu przetwarzającego (rozdz. 9.6);</li>
<li>kompletność i poprawność zestawienia komponentów osób trzecich wraz z wymaganymi notami licencyjnymi (rozdz. 7 i Załącznik A), w tym pełne zbadanie zamknięcia zależności warstwy Rust, którego wykaz w Załączniku A jest zweryfikowany częściowo;</li>
<li>dopuszczalność ograniczeń korzystania w świetle bezwzględnie obowiązujących przepisów o dozwolonym użytku programów komputerowych;</li>
<li>zgodność sposobu korzystania z zewnętrznych usług modeli językowych z warunkami umownymi ich dostawców (rozdz. 3.4 i 5.6);</li>
<li>ocenę dopuszczalności i skuteczności klauzuli zakazu budowy produktu konkurencyjnego oraz zakazu badań bezpieczeństwa (rozdz. 5.2 pkt 8 i 9), a także oznaczenia terminu naprawczego (rozdz. 11.3 pkt 2) i okresu poufności następczej (rozdz. 13.3);</li>
<li>sprawdzenie zgodności niniejszego dokumentu z pozostałymi dokumentami produktu wymienionymi w rozdz. 1.4 pkt 3 — w szczególności co do opisu stanu wykonania funkcji i co do danych liczbowych warstwy wizualnej;</li>
<li>weryfikację rozstrzygnięć Operatora zebranych w Załączniku B pod kątem skuteczności w świetle prawa właściwego.</li>
</ul>
<h3 id="16-sposób-oznaczania-rozstrzygnięć">1.6 Sposób oznaczania rozstrzygnięć</h3>
<p>Kwestie prawne i biznesowe, których dokumentacja projektu nie rozstrzygała, zostały rozstrzygnięte decyzją Operatora z dnia <strong>2026-08-26</strong> i wpisane do treści z oznaczeniem <strong>Rozstrzygnięcie (2026-08-26)</strong>. Komplet rozstrzygnięć, wraz ze wskazaniem przyjętego wariantu, zbiera Załącznik B. Zmiana któregokolwiek rozstrzygnięcia następuje wyłącznie w trybie rozdz. 8.6 albo 14.3.</p>
<hr />
<h2 id="2-definicje">2. Definicje</h2>
<p>Pojęcia pisane w Licencji wielką literą mają znaczenie nadane im w niniejszym rozdziale. Terminologia funkcjonalna odpowiada terminologii obowiązującej w dokumentacji technicznej Danaco Console; rozbieżność między znaczeniem potocznym a znaczeniem zdefiniowanym rozstrzyga się na rzecz znaczenia zdefiniowanego.</p>
<h3 id="21-podmioty">2.1 Podmioty</h3>
<p><strong>Licencjodawca</strong> — Danaco Holding Group Sp. z o.o., podmiot uprawniony z autorskich praw majątkowych do Oprogramowania; w dokumentacji technicznej występujący jako Producent.</p>
<p><strong>Twórca</strong> — Dariusz Naharnowicz, autor koncepcji produktu i twórca Oprogramowania w rozumieniu prawa autorskiego, wskazany jako twórca w metryce każdego dokumentu projektu.</p>
<p><strong>Licencjobiorca</strong> — podmiot, któremu udzielono licencji na korzystanie z Oprogramowania na warunkach niniejszego dokumentu.</p>
<p><strong>Operator</strong> — osoba fizyczna faktycznie obsługująca Oprogramowanie: konfigurująca je, prowadząca sesje pracy z modelem oraz odpowiadająca za polecenia wydawane modelowi i za wykorzystanie wyników jego pracy. Operator działa w imieniu i na rzecz Licencjobiorcy.</p>
<p><strong>Dostawca modelu</strong> — podmiot trzeci świadczący usługę udostępniania modelu sztucznej inteligencji albo dostarczający program kliencki umożliwiający korzystanie z takiego modelu, w szczególności dostawca programu <code>claude</code> (Claude Code CLI) wykorzystywanego przez kanał główny Oprogramowania.</p>
<p><strong>Osoba trzecia</strong> — podmiot inny niż Licencjodawca i Licencjobiorca, w szczególności autor lub uprawniony z praw do Komponentu osoby trzeciej.</p>
<h3 id="22-oprogramowanie-i-jego-warstwy">2.2 Oprogramowanie i jego warstwy</h3>
<p><strong>Oprogramowanie</strong> (Danaco Console) — całość programu komputerowego objętego niniejszą Licencją, obejmująca Rdzeń, Klienta, Powłokę, Kontrakt, schemat trwałości wraz z migracjami, warstwę wizualną oraz zasoby wbudowane, w postaci wykonywalnej dostarczonej Licencjobiorcy.</p>
<p><strong>Rdzeń</strong> — serwerowa część Oprogramowania napisana w języku Go (moduł <code>danacoconsole</code>). Rdzeń prowadzi trwałość, obsługuje komendy Kontraktu, zarządza sesjami i oknami komunikacji, rozstrzyga konfigurację, uruchamia proces Kanału modelu i przekazuje strumień jego odpowiedzi.</p>
<p><strong>Klient</strong> (interfejs) — część Oprogramowania napisana w języku TypeScript i renderowana przez komponent webview systemu operacyjnego. Klient odpowiada za prezentację, przyjmowanie poleceń Operatora i komunikację z Rdzeniem po Kontrakcie.</p>
<p><strong>Powłoka</strong> — natywna warstwa okna zbudowana w technologii Tauri (język Rust), uruchamiająca Rdzeń w tle, otwierająca okno interfejsu i osadzająca ikonę w zasobniku systemowym.</p>
<p><strong>Kontrakt</strong> — zbiór definicji nazw komend, zdarzeń, strumieni, narzędzi i kształtu koperty komunikatu, stanowiący jedyne źródło prawdy nazewnictwa w komunikacji Klienta z Rdzeniem. Kontrakt v2.0 obejmuje <strong>53 komendy, 31 zdarzeń i 39 deklaracji narzędzi modelu</strong>, z których generowane są wiązania dla obu warstw.</p>
<p>Koperta komunikatu jest jedna dla żądania, odpowiedzi i fragmentu strumienia; przesyłana jest w formacie JSON po połączeniu WebSocket. Jej pola <strong>wspólne</strong> to: <code>type</code> (wymagane — typ komunikatu w notacji <code>obszar.zasob.akcja</code>), <code>id</code> (wymagane — identyfikator zadania powtarzany w odpowiedzi i we fragmentach strumienia), <code>timestamp</code> (wymagane — czas nadania w milisekundach epoki), a ponadto <code>sessionId</code> (opcjonalne — sesja, której komunikat dotyczy; puste dla powitania połączenia) oraz <code>payload</code> (opcjonalne — treść właściwa, której kształt wyznacza typ komunikatu). Wyłącznie w odpowiedzi wypełniane są pola <code>status</code> i — przy statusie błędu — <code>error</code>; wyłącznie we fragmentach strumienia pola <code>seq</code> (numer fragmentu, liczony od jedności) i <code>done</code> (prawda w ostatnim fragmencie).</p>
<p>Skrót <code>{type, sessionId, payload}</code>, którym posługują się Instrukcja użytkowania rozdz. 1.4 oraz Instalacja i konfiguracja, nazywa wyłącznie trzy z powyższych pól — te przywoływane w tamtych dokumentach najczęściej — i nie zastępuje wykazu pełnego podanego tutaj.</p>
<p>Brzmienie powyższych nazw jest wiążące zgodnie z regułą wykładni z rozdz. 2.6 pkt 5; źródłem prawdy pozostaje plik Kontraktu, a przytoczony wykaz odpowiada jego stanowi w bieżącym stanie repozytorium (liczba komend, zdarzeń i deklaracji narzędzi nie uległa zmianie względem stanu sprzed prac naprawczych, na której przeprowadzono audyt pierwotny — zgodność potwierdzona odczytem <code>budowa/shared/contract.json</code>, źródła normatywnego Kontraktu).</p>
<p><strong>Katalog danych</strong> — katalog systemu plików, w którym Rdzeń przechowuje trwałość Oprogramowania. Domyślnym katalogiem danych jest <code>%LOCALAPPDATA%\\DanacoConsole</code>; katalog może zostać wskazany zmienną środowiska <code>DANACO_KATALOG_DANYCH</code> albo argumentem wywołania <code>--dane</code>.</p>
<p><strong>Baza</strong> — jeden plik bazy danych SQLite o nazwie <code>danaco-console.db</code> umieszczony w Katalogu danych, stanowiący całość trwałości Rdzenia. Schemat Bazy powstaje i jest aktualizowany przez migracje wkompilowane w Rdzeń, stosowane samoczynnie przy starcie.</p>
<h3 id="23-pojęcia-funkcjonalne-produktu">2.3 Pojęcia funkcjonalne produktu</h3>
<p><strong>Środowisko</strong> — profil widoczności modułów w nawigacji Oprogramowania. Zdefiniowane są cztery Środowiska: <strong>TalkIn</strong> (wiedza), <strong>WorkSpace</strong> (praca), <strong>CodeStudio</strong> (technologia) i <strong>MultitaskingAI</strong> (orkiestracja inteligencji). Środowisko nie jest pojemnikiem na maszyny ani katalogi.</p>
<p><strong>Moduł</strong> — jednostka funkcjonalna widoczna w nawigacji Środowiska.</p>
<p><strong>Okno operacyjne</strong> — wyspecjalizowany widok roboczy Modułu.</p>
<p><strong>Sesja</strong> — nadrzędna jednostka pracy Operatora, wspólna dla plików, pamięci, projektu i agentów; prezentowana jako <strong>karta sesji</strong>.</p>
<p><strong>Okno komunikacji</strong> — byt pośredni między Sesją a wiadomością, niosący moduł, kanał modelu, listę katalogów roboczych, środowisko wykonania, tryb uprawnień i rolę. Okno komunikacji jest najwęższym, ósmym poziomem zasięgu konfiguracji.</p>
<p><strong>Kanał modelu</strong> — wiersz rejestru sterowanego danymi opisujący sposób wywołania modelu. Kolumna <code>rodzaj_kanalu</code> schematu trwałości dopuszcza cztery wartości: <code>cli</code> (kanał główny, uruchamiający zewnętrzny program wiersza poleceń), <code>api</code> (kanał generyczny, wywołujący usługę sieciową opisaną wierszem rejestru), <code>sdk</code> oraz <code>lokalny</code>. Adapter kanału <strong>nie jest</strong> tożsamy z wartością kolumny <code>rodzaj_kanalu</code> — dobierany jest kluczem <code>adapter</code> z parametrów wiersza, wartością danych, nie nazwą typu. Kanał testowy <strong>echo</strong> (bez sieci i bez procesu, odsyłający treść zapytania) zakłada się jako wiersz rodzaju <code>lokalny</code> z parametrem <code>{"adapter":"echo"}</code> — „echo” nie jest samodzielną wartością kolumny <code>rodzaj_kanalu</code>.</p>
<p><strong>Konto</strong> — wpis opisujący tożsamość używaną przez Kanał modelu. Dla kanału rodzaju <code>cli</code> Konto wskazuje katalog konfiguracji przekazywany procesowi modelu zmienną <code>CLAUDE_CONFIG_DIR</code>. Konto nie przechowuje treści poświadczenia — przechowuje wyłącznie odwołanie do niego.</p>
<p><strong>Pula kont</strong> — uporządkowany zbiór Kont jednego rodzaju, w obrębie którego Rdzeń dokonuje rotacji przy rozpoznaniu wyczerpania limitu danego Konta.</p>
<p><strong>Nakładka tożsamości</strong> — trójwarstwowy mechanizm kształtowania zachowania modelu w kolejności krytyczności: <strong>konstytucja → profil → ekspertyza</strong> przekazywany procesowi modelu w trybie <code>ZASTAP</code> albo <code>DOLACZ</code>.</p>
<p><strong>Punkt dostępu</strong> — wpis opisujący maszynę lub katalog, do którego model ma mieć wgląd. <strong>Nadanie</strong> — przypisanie Punktu dostępu do Okna komunikacji wraz z trybem (<code>read</code> albo <code>write</code>) i podzbiorem korzeni.</p>
<p><strong>Prowenancja</strong> — zbiór informacji o tym, co faktycznie zostało przekazane modelowi w danym wywołaniu (argumenty wywołania, treść nakładki tożsamości, ustawienia, sumy kontrolne). Prowenancja pełni funkcję przejrzystości, nie kontroli dostępu. <strong>Prowenancja nie jest utrwalana</strong>: powstaje przed uruchomieniem procesu Tury, żyje wyłącznie w pamięci tego procesu i jest przekazywana do interfejsu na czas jej trwania. Schemat trwałości nie zawiera tabeli ani kolumny prowenancji, w szczególności nie zawiera jej tabela wiadomości. Po zamknięciu Tury odtworzenie z Bazy, co faktycznie zostało przekazane modelowi, nie jest możliwe (rozdz. 10.3 wykaz ograniczeń, pkt 20).</p>
<p><strong>Tura</strong> — pojedyncze wywołanie modelu w Oknie komunikacji, od przyjęcia wiadomości Operatora do zamknięcia strumienia odpowiedzi.</p>
<h3 id="24-pojęcia-dotyczące-danych">2.4 Pojęcia dotyczące danych</h3>
<p><strong>Dane Operatora</strong> — wszelkie dane wprowadzone do Oprogramowania przez Operatora albo powstałe w wyniku korzystania z Oprogramowania po stronie Licencjobiorcy, w szczególności: treść wiadomości i odpowiedzi modelu zapisana w Bazie, nazwy i struktura sesji, okien i kart, wartości ustawień, definicje kanałów modelu i kont, definicje punktów dostępu i nadań, dokumenty tożsamości modelu, a także pliki wytworzone przez model w katalogach roboczych.</p>
<p><strong>Treści Operatora</strong> — wytwór intelektualny powstały przy użyciu Oprogramowania, w szczególności teksty, dokumenty, kod źródłowy i inne materiały wygenerowane albo przetworzone w toku pracy z modelem.</p>
<p><strong>Poświadczenie</strong> — wartość umożliwiająca uwierzytelnienie wobec Dostawcy modelu albo wobec maszyny wskazanej Punktem dostępu (klucz, token, hasło, klucz prywatny). Poświadczenie nie jest przechowywane w Bazie — Baza przechowuje wyłącznie odwołanie do niego.</p>
<p><strong>Katalog roboczy</strong> — katalog systemu plików wskazany Oknu komunikacji, w którym model zostawia i odczytuje pliki pracy.</p>
<h3 id="25-pojęcia-dotyczące-dystrybucji-i-wersji">2.5 Pojęcia dotyczące dystrybucji i wersji</h3>
<p><strong>Egzemplarz</strong> — pojedyncza kopia Oprogramowania udostępniona Licencjobiorcy w postaci wykonywalnej.</p>
<p><strong>Instalacja</strong> — czynność doprowadzenia Egzemplarza do stanu zdatnego do uruchomienia na urządzeniu, wraz z utworzeniem Katalogu danych i Bazy.</p>
<p><strong>Wydanie</strong> — oznaczony zestaw plików wykonywalnych i zasobów Oprogramowania udostępniony przez Licencjodawcę jako całość.</p>
<p><strong>Aktualizacja</strong> — Wydanie następujące po Wydaniu posiadanym przez Licencjobiorcę, zastępujące je w całości albo w części.</p>
<p><strong>Status Deweloperski</strong> — status wersji v2.0 wynikający z przyjętej zasady wersjonowania, oznaczający, że produkt znajduje się w fazie budowy, a jego zakres funkcjonalny nie jest domknięty. Konsekwencje statusu opisuje rozdz. 10.</p>
<p><strong>Rdzeń lokalny</strong> i <strong>rdzeń serwerowy</strong> — dwa umiejscowienia Rdzenia przewidziane topologią docelową: w fazie budowy Rdzeń działa lokalnie na urządzeniu Operatora, docelowo zaś przewidziane jest jego przeniesienie na serwer <code>danaco-system</code> wraz z dystrybucją instalatorów zawierających Klienta i agenta lokalnego na urządzenia użytkowników.</p>
<h3 id="26-reguły-wykładni">2.6 Reguły wykładni</h3>
<ol type="1">
<li>Tytuły rozdziałów i podrozdziałów służą wyłącznie orientacji i nie wpływają na wykładnię postanowień.</li>
<li>Wyliczenia poprzedzone zwrotem „w szczególności” mają charakter przykładowy i nie wyczerpują zakresu pojęcia.</li>
<li>Odwołanie do rozdziału obejmuje wszystkie jego podrozdziały.</li>
<li>Liczba pojedyncza obejmuje liczbę mnogą i odwrotnie, o ile z kontekstu nie wynika inaczej.</li>
<li>Terminy techniczne zapisane czcionką o stałej szerokości (<code>w ten sposób</code>) oznaczają dosłowne nazwy elementów Oprogramowania: zmiennych środowiska, argumentów wywołania, plików, tabel, komend Kontraktu i kluczy konfiguracji. Nazwy te są wiążące co do brzmienia.</li>
</ol>
<hr />
<h2 id="3-przedmiot-licencji">3. Przedmiot licencji</h2>
<h3 id="31-zakres-przedmiotowy">3.1 Zakres przedmiotowy</h3>
<p>Przedmiotem Licencji jest korzystanie z Oprogramowania Danaco Console v2.0 w postaci wykonywalnej, obejmującego:</p>
<ol type="1">
<li><strong>Rdzeń</strong> — serwer w języku Go realizujący obsługę Kontraktu, trwałość w pliku SQLite, rozstrzyganie konfiguracji, prowadzenie sesji i okien komunikacji, uruchamianie procesu kanału modelu oraz przekazywanie strumienia odpowiedzi;</li>
<li><strong>Klienta</strong> — interfejs w języku TypeScript wraz z warstwą wizualną (żetony motywu, dwa motywy: jasny i ciemny, zestaw ikon, kroje pisma);</li>
<li><strong>Powłokę</strong> — natywne okno aplikacji zbudowane w technologii Tauri wraz z osadzeniem ikony w zasobniku systemowym;</li>
<li><strong>Kontrakt</strong> wraz z wygenerowanymi z niego wiązaniami obu warstw;</li>
<li><strong>schemat trwałości</strong> wraz z kompletem migracji wkompilowanych w Rdzeń;</li>
<li><strong>zasoby wbudowane</strong> — ikony, logotypy i sygnety marki, pliki krojów pisma;</li>
<li><strong>dokumentację produktu</strong> przekazaną Licencjobiorcy wraz z Egzemplarzem — pliki <code>README.md</code>, <code>INSTRUKCJA-UZYTKOWANIA.md</code> i <code>INSTALACJA-I-KONFIGURACJA.md</code> wraz z niniejszym plikiem <code>LICENSE.md</code> (rozdz. 1.4 pkt 3 i rozdz. 3.5).</li>
</ol>
<p>Przedmiot Licencji obejmuje wyłącznie <strong>postać wykonywalną</strong> Oprogramowania wraz z dokumentacją produktu wymienioną w pkt 7. Kod źródłowy, repozytorium projektu, dokumentacja projektowa i wewnętrzna Licencjodawcy (w tym zapis decyzji architektonicznych, wykaz znanych luk oraz dokumentacja techniczna repozytorium budowy), materiały koncepcyjne, pakiety projektowe warstwy wizualnej oraz narzędzia budowania <strong>nie stanowią przedmiotu Licencji</strong> i nie są udostępniane Licencjobiorcy, chyba że co innego wynika z odrębnej umowy zawartej na piśmie.</p>
<p>Rozgraniczenie to biegnie po wykazie z rozdz. 1.4 pkt 3, nie po miejscu powstania dokumentu: plik <code>README.md</code> dostarczany w katalogu produktu jest <strong>dokumentacją produktu objętą Licencją</strong>, mimo że pełni jednocześnie funkcję karty technicznej projektu. Nie jest natomiast objęta Licencją dokumentacja projektowa pozostająca w repozytorium budowy i nieprzekazywana Licencjobiorcy.</p>
<h3 id="32-postać-dystrybucyjna-i-topologia">3.2 Postać dystrybucyjna i topologia</h3>
<p>Oprogramowanie jest przeznaczone do pracy jako aplikacja desktopowa systemu Windows z Rdzeniem działającym lokalnie. Powłoka uruchamia Rdzeń jako proces w tle i otwiera okno interfejsu; komunikacja Klienta z Rdzeniem odbywa się po połączeniu WebSocket, domyślnie na porcie <strong>17870</strong>. Konfiguracja Rdzenia jest przyjmowana warstwowo — wartość domyślna, następnie zmienna środowiska, następnie argument wywołania — a zmiennymi rozpoznawanymi przez Rdzeń są wyłącznie: <code>DANACO_ROLA</code>, <code>DANACO_PORT</code>, <code>DANACO_KATALOG_DANYCH</code>, <code>DANACO_KATALOG_KLIENTA</code> oraz <code>DANACO_KATALOG_PROFILI</code>.</p>
<p>Konfiguracja pakietu instalacyjnego Powłoki wskazuje nazwę produktu „Danaco Console”, identyfikator <code>pl.danaco.console</code>, wydawcę „Danaco Holding Group Sp. z o.o.”, notę „© 2026 Danaco Holding Group Sp. z o.o.” oraz format instalatora NSIS dla systemu Windows.</p>
<p>Topologia docelowa przewiduje przeniesienie Rdzenia na serwer <code>danaco-system</code> oraz dystrybucję instalatorów zawierających Klienta i agenta lokalnego na urządzenia użytkowników. <strong>W wersji v2.0 topologia docelowa nie jest wdrożona</strong>: przewidziana jest koncepcją, a Oprogramowanie w bieżącej wersji pracuje z Rdzeniem lokalnym. Zmiana umiejscowienia Rdzenia nie zmienia zakresu niniejszej Licencji; zasady licencjonowania korzystania z Rdzenia serwerowego przez wielu Operatorów jednocześnie określa rozdz. 4.3 akapit ostatni.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> przyjmuje się <strong>model mieszany</strong> dystrybucji: instalator Windows udostępniany wskazanym Licencjobiorcom do pracy indywidualnej (Rdzeń lokalny) oraz — po odrębnym uzgodnieniu — dostęp do Rdzenia serwerowego dla pracy zespołowej. Jednostką licencjonowania jest w obu postaciach <strong>Operator</strong> (rozdz. 4.3).</p>
<h3 id="33-elementy-wyłączone-z-przedmiotu-licencji">3.3 Elementy wyłączone z przedmiotu licencji</h3>
<p>Przedmiotem Licencji <strong>nie są</strong>:</p>
<ol type="1">
<li><strong>modele sztucznej inteligencji</strong> jakiegokolwiek dostawcy — Oprogramowanie nie zawiera modelu, nie dostarcza wag modelu ani nie udziela dostępu do usługi modelu;</li>
<li><strong>program <code>claude</code> (Claude Code CLI)</strong> ani żaden inny zewnętrzny program uruchamiany przez Kanał modelu rodzaju <code>cli</code>; program taki Operator zapewnia we własnym zakresie (rozdz. 3.4);</li>
<li><strong>usługi sieciowe</strong> wywoływane przez Kanał modelu rodzaju <code>api</code> — adres, poświadczenie i warunki korzystania z takiej usługi pochodzą od jej dostawcy;</li>
<li><strong>serwery MCP</strong> i maszyny wskazane Punktami dostępu, wraz z oprogramowaniem na nich zainstalowanym;</li>
<li><strong>Komponenty osób trzecich</strong> w zakresie, w jakim ich własne licencje przyznają uprawnienia szersze niż niniejsza Licencja albo nakładają obowiązki odrębne — do komponentów tych stosuje się ich własne licencje (rozdz. 7);</li>
<li><strong>kod źródłowy</strong> Oprogramowania, repozytorium projektu oraz dokumentacja projektowa i wewnętrzna Licencjodawcy wymieniona w rozdz. 1.4 pkt 4; wyłączenie to <strong>nie obejmuje</strong> dokumentacji produktu z rozdz. 1.4 pkt 3, która jest Licencją objęta (rozdz. 3.1 pkt 7 i rozdz. 3.5);</li>
<li><strong>znaki towarowe, logotypy i oznaczenia</strong> Licencjodawcy, poza zakresem niezbędnym do zwykłego korzystania z Oprogramowania (rozdz. 6.4).</li>
</ol>
<h3 id="34-zależności-zewnętrzne-warunkujące-działanie">3.4 Zależności zewnętrzne warunkujące działanie</h3>
<p>Operator przyjmuje do wiadomości, że pełne wykorzystanie Oprogramowania wymaga zależności pozostających poza zakresem Licencji:</p>
<ol type="1">
<li><p><strong>Program modelu dla kanału głównego.</strong> Kanał modelu rodzaju <code>cli</code> uruchamia zewnętrzny program <code>claude</code> (Claude Code CLI), którego ścieżka pochodzi z wiersza rejestru kanału albo — w braku wskazania — z domyślnej nazwy <code>claude</code> odnajdywanej w zmiennej <code>PATH</code> systemu operacyjnego. Program ten nie jest dostarczany z Oprogramowaniem, a korzystanie z niego podlega warunkom jego dostawcy, w tym warunkom korzystania z konta u tego dostawcy. Uwaga: klucz konfiguracji <code>harness.program_claude</code> istnieje w katalogu ustawień Oprogramowania, lecz <strong>w wersji v2.0 nie jest odczytywany</strong> — wskazanie ścieżki programu innej niż domyślna następuje wyłącznie przez wiersz rejestru kanału.</p></li>
<li><p><strong>Konto Dostawcy modelu.</strong> Kanał główny wymaga co najmniej jednego Konta rodzaju <code>cli</code> wskazującego katalog konfiguracji przekazywany procesowi modelu. Uzyskanie i utrzymanie takiego konta, w tym pokrycie kosztów i przestrzeganie limitów, obciąża Licencjobiorcę.</p></li>
<li><p><strong>Doprowadzenie świeżej instalacji do pracy.</strong> W wersji v2.0 świeża Instalacja <strong>nie jest zdolna do wykonania wywołania modelu bez czynności przygotowawczych</strong>: rejestr kanałów modelu nie posiada zaczynu w migracjach, a pula kont pozostaje pusta. Doprowadzenie Instalacji do pracy wymaga utworzenia wiersza kanału głównego (rodzaj <code>cli</code>) oraz co najmniej jednego Konta rodzaju <code>cli</code>. Zakres dostępności tych dwóch czynności jest różny i wymaga rozróżnienia:</p>
<ul>
<li><strong>kanał modelu</strong> — <strong>interfejs zakładania kanału nie istnieje w wersji v2.0</strong>. Widoki aplikacji wywołują wyłącznie <code>channel.list</code>; utworzenie wiersza rejestru możliwe jest jedynie komendą Kontraktu <code>channel.add</code> wydaną poza interfejsem aplikacji. Komendą tą posługuje się samodzielna strona podglądu warstwy rozmowy (<code>podglad-rozmowy</code>), która przy braku kanału głównego dokłada jego wiersz; strona ta jest narzędziem podglądu warstwy, nie oknem aplikacji, i nie stanowi interfejsu zakładania kanałów w rozumieniu opisu przedmiotu świadczenia;</li>
<li><strong>konto</strong> — pełny cykl życia Konta (<code>account.list</code>, <code>account.add</code>, <code>account.update</code>, <code>account.remove</code>) jest dostępny <strong>z poziomu interfejsu</strong>, w oknie „Modele, konta i tożsamość” (rozdz. 10.3 wykaz funkcji działających, pkt 7). Świeża Instalacja wymaga zatem założenia Konta czynnością Operatora w interfejsie, nie zaś komendą wydaną poza nim.</li>
</ul>
<p>Okoliczność ta jest istotnym elementem opisu przedmiotu świadczenia i została ujawniona wprost, aby wyłączyć wątpliwość co do zgodności Oprogramowania z opisem.</p></li>
<li><p><strong>Środowisko systemowe.</strong> Oprogramowanie wymaga systemu operacyjnego Windows z dostępnym komponentem webview wykorzystywanym przez Powłokę oraz uprawnieniami pozwalającymi na zapis w Katalogu danych i nasłuch na porcie lokalnym.</p></li>
</ol>
<h3 id="35-dokumentacja-objęta-licencją">3.5 Dokumentacja objęta licencją</h3>
<p>Licencją objęta jest dokumentacja produktu przekazana Licencjobiorcy wraz z Egzemplarzem, w zakresie niezbędnym do korzystania z Oprogramowania. Jest to zamknięty zestaw czterech plików umieszczanych w katalogu produktu obok plików wykonywalnych:</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Plik</th>
<th>Rola</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>LICENSE.md</code></td>
<td>niniejszy dokument — warunki licencyjne</td>
</tr>
<tr>
<td><code>README.md</code></td>
<td>karta techniczno-produktowa: koncepcja, architektura, stos, zakres funkcjonalny wraz ze stanem wykonania</td>
</tr>
<tr>
<td><code>INSTRUKCJA-UZYTKOWANIA.md</code></td>
<td>praca z gotową Instalacją — widoki, Okno komunikacji, konfiguracja, konta i tożsamość modelu</td>
</tr>
<tr>
<td><code>INSTALACJA-I-KONFIGURACJA.md</code></td>
<td>instalacja, budowanie warstw, doprowadzenie świeżej Instalacji do pracy, zmienne środowiska i argumenty wywołania</td>
</tr>
</tbody>
</table></div>
<p>Wszędzie, gdzie Licencja posługuje się zwrotem „dokumentacja przekazana Licencjobiorcy” albo „dokumentacja produktu”, rozumie się przez to wyłącznie powyższy zestaw. Zwroty opisowe („dokumentacja wdrożeniowa”, „opis funkcjonalny”, „instrukcja Operatora”) nie są nazwami odrębnych dokumentów; odpowiadające im treści mieszczą się odpowiednio w plikach <code>INSTALACJA-I-KONFIGURACJA.md</code>, <code>README.md</code> i <code>INSTRUKCJA-UZYTKOWANIA.md</code>.</p>
<p>Licencjobiorca może sporządzać kopie dokumentacji na własne potrzeby i udostępniać ją swoim Operatorom. Publikowanie dokumentacji, jej fragmentów ani opracowań na jej podstawie poza organizacją Licencjobiorcy wymaga uprzedniej zgody Licencjodawcy udzielonej na piśmie.</p>
<hr />
<h2 id="4-udzielenie-licencji-i-zakres-uprawnień">4. Udzielenie licencji i zakres uprawnień</h2>
<h3 id="41-udzielenie-licencji">4.1 Udzielenie licencji</h3>
<p>Licencjodawca udziela Licencjobiorcy licencji <strong>niewyłącznej</strong>, <strong>nieprzenoszalnej</strong>, <strong>bez prawa udzielania sublicencji</strong>, na korzystanie z Oprogramowania w postaci wykonywalnej, wyłącznie na warunkach i w granicach określonych niniejszym dokumentem.</p>
<p>Licencja obejmuje wyłącznie pola eksploatacji wskazane wprost w rozdz. 4.2. Uprawnienia niewymienione nie są udzielane; w szczególności nie jest udzielane prawo zwielokrotniania Oprogramowania w celu wprowadzenia do obrotu, prawo najmu ani prawo użyczenia Egzemplarza.</p>
<p>Licencja nie przenosi na Licencjobiorcę żadnych autorskich praw majątkowych ani praw zależnych do Oprogramowania. Prawa te pozostają przy Licencjodawcy.</p>
<h3 id="42-uprawnienia-przyznane">4.2 Uprawnienia przyznane</h3>
<p>W ramach Licencji Licencjobiorca jest uprawniony do:</p>
<ol type="1">
<li><strong>instalowania i uruchamiania</strong> Oprogramowania na urządzeniach objętych zakresem podmiotowym Licencji (rozdz. 4.3), w celach wewnętrznej działalności Licencjobiorcy;</li>
<li><strong>korzystania z pełnej funkcjonalności</strong> udostępnionej w danym Wydaniu, w granicach opisanych w niniejszej Licencji i dokumentacji;</li>
<li><strong>konfigurowania</strong> Oprogramowania w zakresie przewidzianym mechanizmami samego produktu: definiowania kanałów modelu, kont i puli kont, punktów dostępu i nadań, dokumentów tożsamości modelu, katalogów roboczych, wartości ustawień na przewidzianych poziomach zasięgu i osiach rozstrzygania — przy czym uprawnienie to obejmuje <strong>zapis i rozstrzyganie</strong> wartości ustawień, nie zaś zapewnienie, że każda zapisana wartość zostanie zastosowana przy wywołaniu modelu; zakres wartości faktycznie stosowanych oraz ograniczenia zakładania kanałów i punktów dostępu ujawnia rozdz. 10.3 (wykaz ograniczeń, pkt 21 i 22) oraz rozdz. 3.4 pkt 3;</li>
<li><strong>sporządzenia kopii zapasowej</strong> Egzemplarza w liczbie niezbędnej do zabezpieczenia przed utratą, z zachowaniem wszystkich oznaczeń praw autorskich; kopia zapasowa nie może być używana równolegle z Egzemplarzem podstawowym;</li>
<li><strong>sporządzania kopii zapasowych Katalogu danych i Bazy</strong> — bez ograniczeń liczbowych, przy czym Licencjobiorca odpowiada za zabezpieczenie tych kopii stosownie do wrażliwości zawartych w nich Danych Operatora;</li>
<li><strong>korzystania z wytworów pracy</strong> — Treści Operatora powstałe przy użyciu Oprogramowania należą do Licencjobiorcy w granicach opisanych w rozdz. 6.5;</li>
<li><strong>wskazywania Operatorów</strong> uprawnionych do obsługi Oprogramowania w imieniu Licencjobiorcy, w liczbie wynikającej z rozdz. 4.3, przy czym Licencjobiorca odpowiada za ich działania i zaniechania jak za własne.</li>
</ol>
<h3 id="43-zakres-podmiotowy--liczba-stanowisk-i-użytkowników">4.3 Zakres podmiotowy — liczba stanowisk i użytkowników</h3>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> Licencja jest <strong>licencją imienną, liczoną według Operatorów</strong>. Jeden uprawniony Operator może korzystać z Oprogramowania na dowolnej liczbie urządzeń pozostających w dyspozycji Licencjobiorcy, przy zakazie korzystania równoczesnego przez tę samą osobę na więcej niż jednym urządzeniu. Liczbę uprawnionych Operatorów określa uzgodnienie towarzyszące udostępnieniu Egzemplarza (zamówienie, potwierdzenie udostępnienia albo odrębna umowa); w braku takiego uzgodnienia liczba uprawnionych Operatorów wynosi <strong>jeden</strong>. W topologii z Rdzeniem serwerowym jednostką liczenia pozostaje Operator: jedna Instalacja Rdzenia może obsługiwać wyłącznie Operatorów mieszczących się w uzgodnionej liczbie.</p>
<p>Niezależnie od wybranego wariantu obowiązują zasady następujące:</p>
<ol type="1">
<li>Licencjobiorca prowadzi ewidencję Instalacji na potrzeby wykazania zgodności korzystania z zakresem Licencji.</li>
<li>Udostępnienie Oprogramowania podmiotowi spoza organizacji Licencjobiorcy — w tym w modelu świadczenia usług na rzecz osób trzecich przy użyciu Oprogramowania — wymaga uprzedniej zgody Licencjodawcy udzielonej na piśmie.</li>
<li>Korzystanie przez podmioty powiązane z Licencjobiorcą nie jest objęte Licencją, chyba że wynika to wprost z odrębnej umowy.</li>
</ol>
<h3 id="44-zakres-terytorialny-i-czasowy">4.4 Zakres terytorialny i czasowy</h3>
<p>Licencja jest udzielana bez ograniczeń terytorialnych, z zastrzeżeniem, że Licencjobiorca zapewnia zgodność korzystania z przepisami obowiązującymi w miejscu korzystania, w szczególności z przepisami o ochronie danych osobowych, o kontroli eksportu oraz z regulacjami dotyczącymi systemów sztucznej inteligencji.</p>
<p>Okres obowiązywania Licencji reguluje rozdz. 11.</p>
<h3 id="45-odpłatność">4.5 Odpłatność</h3>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> Licencja na Wydania o statusie <strong>Deweloperskim jest nieodpłatna</strong>. Odpłatność — wraz z modelem rozliczeń, zasadami fakturowania i skutkami opóźnienia w zapłacie — zostanie wprowadzona najwcześniej wraz z pierwszym Wydaniem o statusie innym niż Deweloperski, nową wersją dokumentu Licencji w trybie rozdz. 8.6. Nieodpłatność fazy Deweloperskiej nie stanowi zrzeczenia się prawa do wynagrodzenia za Wydania przyszłe.</p>
<h3 id="46-brak-przeniesienia-praw">4.6 Brak przeniesienia praw</h3>
<p>Licencja nie stanowi sprzedaży Oprogramowania. Licencjobiorca nabywa wyłącznie prawo do korzystania w zakresie opisanym w rozdz. 4.2. Wszelkie prawa nieudzielone wprost pozostają zastrzeżone na rzecz Licencjodawcy.</p>
<p>Licencjobiorca nie może przenieść praw ani obowiązków wynikających z Licencji na osobę trzecią — w drodze czynności prawnej, wniesienia aportem, połączenia, podziału ani przejęcia — bez uprzedniej zgody Licencjodawcy udzielonej na piśmie pod rygorem nieważności.</p>
<h3 id="47-zakres-wykonania-modelu-a-zakres-licencji">4.7 Zakres wykonania modelu a zakres licencji</h3>
<p>Zgodnie z przyjętą architekturą zakres wykonania modelu (<code>local</code>, <code>core</code>, <code>remote</code>) jest parametrem Okna komunikacji i pozostaje niezależny od umiejscowienia Rdzenia. Licencja nie ogranicza Operatora co do wyboru maszyn, na których model pracuje, z zastrzeżeniem, że:</p>
<ol type="1">
<li>Licencjobiorca musi posiadać tytuł prawny do korzystania z maszyn i zasobów udostępnianych modelowi przez Punkty dostępu i Katalogi robocze;</li>
<li>udostępnienie modelowi zasobów osoby trzeciej wymaga jej zgody;</li>
<li>w wersji v2.0 wybór zakresu wykonania <strong>jest wyłącznie etykietą w prowenancji i nie zmienia miejsca uruchomienia procesu</strong> — proces modelu startuje zawsze na maszynie Rdzenia, a kanał połączenia zdalnego (SSH) jako kanał modelu nie istnieje. Operator nie może zatem opierać się na wartości zakresu wykonania jako na mechanizmie ograniczającym; okoliczność ta ma bezpośredni wpływ na ocenę ryzyka udostępniania zasobów i została ujawniona wprost.</li>
</ol>
<hr />
<h2 id="5-zasady-korzystania-i-ograniczenia">5. Zasady korzystania i ograniczenia</h2>
<h3 id="51-zasada-ogólna">5.1 Zasada ogólna</h3>
<p>Licencjobiorca korzysta z Oprogramowania zgodnie z jego przeznaczeniem, dokumentacją i niniejszą Licencją, w sposób nienaruszający praw Licencjodawcy, praw osób trzecich oraz przepisów prawa. Wątpliwość co do zakresu dozwolonego korzystania rozstrzyga się przez zapytanie skierowane do Licencjodawcy, nie przez samodzielne rozszerzającą wykładnię uprawnień.</p>
<h3 id="52-ograniczenia-dotyczące-samego-oprogramowania">5.2 Ograniczenia dotyczące samego Oprogramowania</h3>
<p>Licencjobiorca zobowiązuje się <strong>nie</strong>:</p>
<ol type="1">
<li>zwielokrotniać Oprogramowania w zakresie wykraczającym poza czynności niezbędne do jego zainstalowania, uruchomienia i sporządzenia kopii zapasowej;</li>
<li>rozpowszechniać Oprogramowania ani jego części — odpłatnie ani nieodpłatnie — w szczególności przez udostępnienie plików wykonywalnych, obrazów instalacyjnych, kontenerów ani obrazów maszyn wirtualnych zawierających Oprogramowanie;</li>
<li>najmować, użyczać, oddawać w leasing ani udostępniać Oprogramowania osobom trzecim w modelu współdzielenia dostępu, hostingu ani świadczenia usług przy jego użyciu na rzecz osób trzecich;</li>
<li>dokonywać dekompilacji, dezasemblacji ani innej postaci odtwarzania kodu źródłowego, z zastrzeżeniem rozdz. 5.7;</li>
<li>modyfikować Oprogramowania, tworzyć jego opracowań, tłumaczeń ani adaptacji, w tym modyfikować plików wykonywalnych, zasobów wbudowanych, migracji schematu ani warstwy wizualnej;</li>
<li>usuwać, zasłaniać ani zmieniać oznaczeń praw autorskich, znaków towarowych, not licencyjnych, oznaczeń wersji ani statusu wersji, zarówno w interfejsie Oprogramowania, jak i w plikach dystrybucji;</li>
<li>obchodzić ani wyłączać mechanizmów Oprogramowania służących identyfikacji wersji, prowenancji wywołania modelu ani zapisowi historii;</li>
<li>odwzorowywać chronionych elementów Oprogramowania wskazanych w rozdz. 6.1 — w szczególności Kontraktu, modelu danych i warstwy wizualnej — w produkcie własnym albo osoby trzeciej; zakaz ten nie obejmuje tworzenia produktów tej samej kategorii bez wykorzystania elementów chronionych;</li>
<li>przeprowadzać testów penetracyjnych, testów obciążeniowych ani badań bezpieczeństwa <strong>cudzych instalacji</strong> Oprogramowania; badania we własnym środowisku Licencjobiorcy są dozwolone, a ustalone podatności podlegają zgłoszeniu wyłącznie Licencjodawcy w trybie rozdz. 13.4.</li>
</ol>
<p><strong>Rozstrzygnięcie (2026-08-26) do pkt 8:</strong> przyjęto <strong>zakaz wąski</strong> — ochronie podlega odwzorowanie chronionych elementów wskazanych w rozdz. 6.1; tworzenie produktów tej samej kategorii bez wykorzystania tych elementów nie jest ograniczone. <strong>Rozstrzygnięcie do pkt 9:</strong> badania bezpieczeństwa i testy obciążeniowe <strong>we własnym środowisku Licencjobiorcy są dozwolone</strong>; zakaz obejmuje wyłącznie badanie cudzych instalacji, a ustalenia podlegają zgłoszeniu w trybie rozdz. 13.4. Rozstrzygnięcie usuwa sprzeczność między zakazem badań a obowiązkiem zgłaszania podatności.</p>
<h3 id="53-ograniczenia-wynikające-ze-statusu-deweloperskiego">5.3 Ograniczenia wynikające ze statusu Deweloperskiego</h3>
<p>Ze względu na status Deweloperski wersji v2.0 Licencjobiorca zobowiązuje się <strong>nie stosować</strong> Oprogramowania:</p>
<ol type="1">
<li>w procesach, w których błąd, przerwanie pracy albo utrata danych mogą spowodować zagrożenie życia, zdrowia albo bezpieczeństwa osób;</li>
<li>w systemach sterowania infrastrukturą krytyczną, urządzeniami medycznymi, środkami transportu ani w innych zastosowaniach o podwyższonym ryzyku;</li>
<li>jako jedynego miejsca przechowywania danych o znaczeniu istotnym dla działalności Licencjobiorcy — bez niezależnej kopii zapasowej prowadzonej poza Katalogiem danych;</li>
<li>jako narzędzia podejmującego automatycznie decyzje wywołujące skutki prawne wobec osób fizycznych albo w podobny sposób istotnie na nie wpływające, bez udziału człowieka.</li>
</ol>
<p>Ograniczenia te wynikają z rzeczywistego stanu wykonania produktu opisanego w rozdz. 10.3 i pozostają w mocy do czasu wydania przez Licencjodawcę wersji o statusie innym niż Deweloperski.</p>
<h3 id="54-zasady-korzystania-z-modelu-i-odpowiedzialność-za-polecenia">5.4 Zasady korzystania z modelu i odpowiedzialność za polecenia</h3>
<ol type="1">
<li>Oprogramowanie jest narzędziem pośredniczącym: przekazuje modelowi treść wskazaną przez Operatora wraz z Nakładką tożsamości, ustawieniami Okna komunikacji i konfiguracją mostów, po czym zwraca strumień odpowiedzi. <strong>Licencjodawca nie tworzy, nie weryfikuje ani nie zatwierdza treści wytwarzanych przez model.</strong></li>
<li>Operator odpowiada za treść poleceń wydawanych modelowi oraz za wykorzystanie wyników jego pracy, w tym za weryfikację ich poprawności przed użyciem.</li>
<li>Operator odpowiada za zakres uprawnień nadanych modelowi: za wybór trybu uprawnień, listę Katalogów roboczych, Punkty dostępu i Nadania (<code>read</code> albo <code>write</code>) oraz za konfigurację mostów MCP. Operator przyjmuje do wiadomości, że nadanie trybu <code>write</code> oznacza możliwość dokonywania przez model zmian w plikach na maszynach objętych nadaniem.</li>
<li>Operator przyjmuje do wiadomości, że <strong>w wersji v2.0 mechanizmy izolacji opisane dokumentacją (jedenaście ustawień izolacji rozstrzyganych na ośmiu poziomach zasięgu) są zdefiniowane i rozstrzygane, lecz nie mają egzekutorów — nie wpływają na sposób uruchomienia procesu modelu.</strong> Ustawienia izolacji nie mogą być traktowane jako zabezpieczenie techniczne.</li>
<li>Operator przyjmuje do wiadomości, że <strong>w wersji v2.0 model nie zachowuje ciągłości rozmowy między turami</strong>: każda tura uruchamia nowy proces bez historii poprzednich wiadomości Okna komunikacji. Okoliczność ta ma wpływ na sposób formułowania poleceń i na wynik pracy.</li>
</ol>
<h3 id="55-obowiązki-licencjobiorcy-w-zakresie-bezpieczeństwa">5.5 Obowiązki Licencjobiorcy w zakresie bezpieczeństwa</h3>
<ol type="1">
<li>Licencjobiorca zabezpiecza urządzenie, na którym działa Rdzeń, w sposób adekwatny do wrażliwości Danych Operatora zapisanych w Bazie.</li>
<li>Licencjobiorca odpowiada za ochronę Poświadczeń wskazywanych Oprogramowaniu przez odwołania (klucze API, klucze prywatne SSH, katalogi konfiguracji kont). Oprogramowanie nie przechowuje treści Poświadczeń w Bazie — ich bezpieczeństwo pozostaje po stronie Licencjobiorcy.</li>
<li>Licencjobiorca przyjmuje do wiadomości, że mechanizm uwierzytelniania jest <strong>w fazie budowy wyłączony</strong>, a rozstrzyganie dostępu do Rdzenia opiera się wyłącznie na zabezpieczeniu środowiska, w którym Rdzeń działa. W konsekwencji <strong>każdy podmiot mający dostęp sieciowy do portu nasłuchu Rdzenia ma dostęp do pełnej funkcjonalności Oprogramowania</strong>. Ograniczenie dostępu do tego portu należy do Licencjobiorcy.</li>
<li>Licencjobiorca nie udostępnia portu nasłuchu Rdzenia w sieci publicznej.</li>
</ol>
<h3 id="56-zgodność-z-warunkami-dostawców-modeli">5.6 Zgodność z warunkami Dostawców modeli</h3>
<p>Korzystanie z Oprogramowania w sposób prowadzący do wywołania modelu podlega równolegle warunkom umownym Dostawcy modelu. Licencjobiorca zobowiązuje się:</p>
<ol type="1">
<li>przestrzegać regulaminów i limitów Dostawcy modelu, w tym zakazów dotyczących współdzielenia kont, automatyzacji i obchodzenia limitów;</li>
<li>nie wykorzystywać mechanizmu Puli kont i rotacji Kont w sposób sprzeczny z warunkami Dostawcy modelu; mechanizm ten służy uporządkowanej pracy przy wielu uprawnionych kontach, nie zaś obchodzeniu ograniczeń;</li>
<li>zapewnić, że treści przekazywane modelowi mogą być temu modelowi przekazane zgodnie z prawem i z zobowiązaniami Licencjobiorcy wobec osób trzecich.</li>
</ol>
<p>Licencjodawca nie jest stroną stosunku prawnego między Licencjobiorcą a Dostawcą modelu i nie odpowiada za jego treść, wykonanie ani skutki zakończenia.</p>
<h3 id="57-uprawnienia-wynikające-z-przepisów-bezwzględnie-obowiązujących">5.7 Uprawnienia wynikające z przepisów bezwzględnie obowiązujących</h3>
<p>Ograniczenia opisane w rozdz. 5.2 pkt 4 i 5 nie naruszają uprawnień przysługujących Licencjobiorcy na podstawie bezwzględnie obowiązujących przepisów prawa właściwego, w szczególności — jeżeli prawem właściwym będzie prawo polskie — uprawnień do zwielokrotniania kodu i tłumaczenia jego formy w zakresie niezbędnym do uzyskania współdziałania z innymi programami, na zasadach i przy zachowaniu warunków określonych w art. 75 ustawy z dnia 4 lutego 1994 r. o prawie autorskim i prawach pokrewnych. Licencjobiorca zamierzający skorzystać z tych uprawnień zwraca się uprzednio do Licencjodawcy o udostępnienie informacji niezbędnych do współdziałania; Licencjodawca udostępnia je w rozsądnym terminie, o ile pozostaje to w jego możliwościach.</p>
<h3 id="58-skutki-naruszenia-ograniczeń">5.8 Skutki naruszenia ograniczeń</h3>
<p>Naruszenie ograniczeń opisanych w rozdz. 5.2, 5.3 albo 5.6 stanowi istotne naruszenie Licencji i uprawnia Licencjodawcę do jej wypowiedzenia w trybie rozdz. 11.3, niezależnie od dalej idących roszczeń.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> Licencja <strong>nie przewiduje kar umownych</strong>. Odpowiedzialność za naruszenie ograniczeń kształtują zasady ogólne prawa właściwego oraz uprawnienie do wypowiedzenia z rozdz. 11.3. Rozstrzygnięcie odpowiada nieodpłatnemu charakterowi fazy Deweloperskiej (rozdz. 4.5) i może ulec zmianie wraz z wprowadzeniem odpłatności, w trybie rozdz. 8.6.</p>
<hr />
<h2 id="6-własność-intelektualna">6. Własność intelektualna</h2>
<h3 id="61-prawa-do-oprogramowania">6.1 Prawa do Oprogramowania</h3>
<p>Oprogramowanie stanowi utwór w rozumieniu prawa autorskiego. Autorskie prawa majątkowe do Oprogramowania przysługują Licencjodawcy — Danaco Holding Group Sp. z o.o. Autorskie prawa osobiste przysługują Twórcy — Dariuszowi Naharnowiczowi — i nie podlegają zrzeczeniu ani przeniesieniu.</p>
<p>Ochronie podlega całość Oprogramowania oraz jego elementy dające się wyodrębnić, w szczególności:</p>
<ol type="1">
<li>kod źródłowy i kod wynikowy wszystkich trzech warstw (Rdzeń, Klient, Powłoka);</li>
<li><strong>Kontrakt</strong> — dobór, nazewnictwo i struktura komend, zdarzeń, strumieni i deklaracji narzędzi, wraz z konwencją generowania wiązań;</li>
<li><strong>model danych</strong> — struktura schematu trwałości, dobór i nazewnictwo tabel, kolumn oraz łańcuch więzów sesja → karta → okno → wiadomość;</li>
<li><strong>warstwa wizualna</strong> — zestaw żetonów motywu, dwa motywy (jasny i ciemny), siatka i zasady zestawu ikon, kompozycja widoków;</li>
<li><strong>oznaczenia marki</strong> — logotypy, sygnety i ich warianty;</li>
<li>dokumentacja produktu.</li>
</ol>
<p>Struktura katalogów, nazewnictwo modułów oraz konwencje architektoniczne Oprogramowania stanowią wyraz autorskiego doboru i również podlegają ochronie w zakresie, w jakim mają charakter twórczy.</p>
<h3 id="62-zastrzeżenie-praw">6.2 Zastrzeżenie praw</h3>
<p>Wszelkie prawa nieudzielone Licencjobiorcy wprost w rozdz. 4.2 pozostają zastrzeżone na rzecz Licencjodawcy. Milczenie Licencji co do określonego pola eksploatacji oznacza brak zgody, nie zaś zgodę dorozumianą.</p>
<h3 id="63-zakaz-działań-naruszających">6.3 Zakaz działań naruszających</h3>
<p>Licencjobiorca zobowiązuje się nie podejmować działań zmierzających do podważenia praw Licencjodawcy, w szczególności nie rejestrować na swoją rzecz oznaczeń zbieżnych z oznaczeniami Oprogramowania, nazw domen zawierających oznaczenie „Danaco Console” ani nie zgłaszać do ochrony rozwiązań odwzorowujących architekturę albo model danych Oprogramowania.</p>
<h3 id="64-znaki-towarowe-i-oznaczenia">6.4 Znaki towarowe i oznaczenia</h3>
<p>Oznaczenia „Danaco”, „Danaco Console”, „Danaco Holding Group”, nazwy Środowisk (TalkIn, WorkSpace, CodeStudio, MultitaskingAI) oraz logotypy i sygnety produktu stanowią oznaczenia Licencjodawcy. Licencja <strong>nie przyznaje</strong> prawa do używania tych oznaczeń poza zakresem niezbędnym do zwykłego korzystania z Oprogramowania i wewnętrznego odwoływania się do niego.</p>
<p>Użycie oznaczeń w materiałach marketingowych, w komunikacji publicznej, w nazwach produktów Licencjobiorcy albo w sposób sugerujący powiązanie gospodarcze z Licencjodawcą wymaga uprzedniej zgody udzielonej na piśmie.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> do czasu ewentualnej rejestracji znaków oznaczenia „Danaco Console” oraz oznaczenia graficzne produktu podlegają ochronie jako <strong>oznaczenia nierejestrowane</strong> — na podstawie prawa autorskiego do logotypów, przepisów o zwalczaniu nieuczciwej konkurencji oraz ochrony firmy Licencjodawcy. Uzyskanie praw ochronnych na znaki zostanie odzwierciedlone w kolejnej wersji dokumentu numerami praw, w trybie rozdz. 8.6, bez zmiany zakresu zobowiązań Licencjobiorcy.</p>
<h3 id="65-prawa-do-treści-operatora">6.5 Prawa do Treści Operatora</h3>
<ol type="1">
<li>Licencjodawca nie nabywa praw do Treści Operatora ani do Danych Operatora. Prawa te — w zakresie, w jakim powstają — przysługują Licencjobiorcy albo podmiotom uprawnionym zgodnie z prawem właściwym.</li>
<li>Licencjodawca nie uzyskuje dostępu do Danych Operatora w związku z samym korzystaniem przez Licencjobiorcę z Oprogramowania. Dane te pozostają w Katalogu danych na urządzeniu Licencjobiorcy (rozdz. 9).</li>
<li>Licencjodawca nie składa żadnego oświadczenia co do tego, czy i w jakim zakresie wytwór wygenerowany przez model podlega ochronie prawnoautorskiej, ani co do tego, komu przysługują do niego prawa. Ocena ta zależy od prawa właściwego, sposobu powstania wytworu i wkładu twórczego człowieka i pozostaje po stronie Licencjobiorcy.</li>
<li>Licencjobiorca odpowiada za to, aby korzystanie z wytworów pracy modelu nie naruszało praw osób trzecich.</li>
</ol>
<h3 id="66-informacje-zwrotne">6.6 Informacje zwrotne</h3>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> przekazując Licencjodawcy uwagi, zgłoszenia błędów albo propozycje funkcji („<strong>Informacje zwrotne</strong>”), Licencjobiorca udziela Licencjodawcy nieodpłatnej, niewyłącznej, nieograniczonej czasowo ani terytorialnie licencji na korzystanie z Informacji zwrotnych w celu rozwoju, poprawy i utrzymania Oprogramowania, bez obowiązku wskazywania źródła. Przekazanie Informacji zwrotnych jest dobrowolne; Licencjobiorca nie przekazuje w tym trybie treści, do których prawa przysługują osobom trzecim, ani rozwiązań, które zamierza zastrzec — w takim wypadku zakres korzystania wymaga odrębnego uzgodnienia.</p>
<hr />
<h2 id="7-komponenty-i-rozszerzenia">7. Komponenty i rozszerzenia</h2>
<h3 id="71-zasada-nadrzędna">7.1 Zasada nadrzędna</h3>
<p>Oprogramowanie zawiera oraz wykorzystuje komponenty osób trzecich („<strong>Komponenty</strong>”), rozpowszechniane na własnych licencjach. Do każdego Komponentu stosuje się <strong>jego własną licencję</strong>; niniejsza Licencja nie ogranicza uprawnień przyznanych Licencjobiorcy przez te licencje ani nie rozszerza ich na kod autorski Licencjodawcy.</p>
<p>W razie sprzeczności między postanowieniem niniejszej Licencji a warunkiem licencji Komponentu — w zakresie dotyczącym tego Komponentu — pierwszeństwo ma licencja Komponentu.</p>
<p>Zestawienie Komponentów w niniejszym rozdziale sporządzono na podstawie faktycznej zawartości repozytorium w stanie bieżącym: plików deklaracji zależności (<code>budowa/go.mod</code>, <code>budowa/client/package.json</code>, <code>budowa/desktop/src-tauri/Cargo.toml</code>) oraz plików zamknięcia zależności (<code>budowa/client/package-lock.json</code>, <code>budowa/desktop/src-tauri/Cargo.lock</code>). Zestawienie warstwy Rdzenia jest niezmienione względem rewizji pierwotnej stanu sprzed prac naprawczych; zestawienie warstwy Klienta uwzględnia dodanie zależności narzędziowej <code>jsdom</code> (środowisko DOM podkładane pod sprawdziany widoków, prace nad sprawdzianami widoków) — patrz rozdz. 7.3, pozostałe pozycje warstwy Klienta są niezmienione; zestawienie warstwy Powłoki uwzględnia usunięcie nieużywanej zależności bezpośredniej <code>serde_json</code> z <code>Cargo.toml</code>, dokonane w pracach nad skryptami wydania (prace nad skryptami wydania) — patrz rozdz. 7.4. Zestawienie zbiorcze zawiera Załącznik A.</p>
<h3 id="72-komponenty-warstwy-rdzenia-go">7.2 Komponenty warstwy Rdzenia (Go)</h3>
<p>Rdzeń jest budowany jako moduł Go <code>danacoconsole</code> (wymagana wersja języka: <code>go 1.26</code>). Zależności bezpośrednie i pośrednie deklarowane w <code>budowa/go.mod</code> oraz rodzaje ich licencji ustalone z treści plików licencyjnych w lokalnej pamięci podręcznej modułów:</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Komponent</th>
<th>Wersja</th>
<th>Rodzaj licencji</th>
<th>Rola w produkcie</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>github.com/coder/websocket</code></td>
<td>v1.8.15</td>
<td>ISC</td>
<td>transport WebSocket Rdzenia</td>
</tr>
<tr>
<td><code>golang.org/x/sys</code></td>
<td>v0.47.0</td>
<td>BSD 3-Clause</td>
<td>dostęp do funkcji systemowych</td>
</tr>
<tr>
<td><code>modernc.org/sqlite</code></td>
<td>v1.56.0</td>
<td>BSD 3-Clause</td>
<td>silnik SQLite bez zależności natywnych</td>
</tr>
<tr>
<td><code>github.com/dustin/go-humanize</code></td>
<td>v1.0.1</td>
<td>MIT</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
<tr>
<td><code>github.com/google/uuid</code></td>
<td>v1.6.0</td>
<td>BSD 3-Clause</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
<tr>
<td><code>github.com/mattn/go-isatty</code></td>
<td>v0.0.24</td>
<td>MIT</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
<tr>
<td><code>github.com/ncruces/go-strftime</code></td>
<td>v1.0.0</td>
<td>MIT</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
<tr>
<td><code>github.com/remyoudompheng/bigfft</code></td>
<td>v0.0.0-20230129092748-24d4a6f8daec</td>
<td>BSD 3-Clause</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
<tr>
<td><code>modernc.org/libc</code></td>
<td>v1.74.4</td>
<td>BSD 3-Clause</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
<tr>
<td><code>modernc.org/mathutil</code></td>
<td>v1.7.1</td>
<td>BSD 3-Clause</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
<tr>
<td><code>modernc.org/memory</code></td>
<td>v1.11.0</td>
<td>BSD 3-Clause</td>
<td>zależność pośrednia silnika SQLite</td>
</tr>
</tbody>
</table></div>
<p>Wszystkie wymienione licencje są licencjami zezwalającymi, dopuszczającymi rozpowszechnianie w postaci wynikowej pod warunkiem zachowania not o prawach autorskich i tekstu licencji. Licencjodawca spełnia ten warunek przez dołączenie zestawienia z Załącznika A wraz z pełnymi tekstami licencji do Wydania (rozdz. 7.9).</p>
<h3 id="73-komponenty-warstwy-klienta-typescript">7.3 Komponenty warstwy Klienta (TypeScript)</h3>
<p>Zależnością <strong>produkcyjną</strong> Klienta — to znaczy taką, której kod trafia do pakietu dostarczanego Licencjobiorcy — jest wyłącznie:</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Komponent</th>
<th>Wersja</th>
<th>Rodzaj licencji</th>
<th>Rola w produkcie</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>@tauri-apps/api</code></td>
<td>2.5.0</td>
<td>Apache-2.0 <strong>lub</strong> MIT (do wyboru)</td>
<td>wywołanie poleceń Powłoki z poziomu interfejsu</td>
</tr>
</tbody>
</table></div>
<p>Pozostałe zależności Klienta mają charakter <strong>narzędzi budowania i testowania</strong> i nie wchodzą do pakietu dostarczanego Licencjobiorcy:</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Komponent</th>
<th>Wersja</th>
<th>Rodzaj licencji</th>
<th>Rola</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>typescript</code></td>
<td>5.8.3</td>
<td>Apache-2.0</td>
<td>kontrola typów (<code>tsc --noEmit</code>)</td>
</tr>
<tr>
<td><code>vite</code></td>
<td>6.3.5</td>
<td>MIT</td>
<td>budowanie pakietu interfejsu</td>
</tr>
<tr>
<td><code>vitest</code></td>
<td>3.2.7</td>
<td>MIT</td>
<td>sprawdziany jednostkowe interfejsu</td>
</tr>
<tr>
<td><code>jsdom</code></td>
<td>^29.1.1</td>
<td>MIT</td>
<td>środowisko DOM dla sprawdzianów widoków (prace nad sprawdzianami widoków)</td>
</tr>
</tbody>
</table></div>
<p>Zamknięcie zależności zapisane w <code>package-lock.json</code> obejmuje <strong>141 pozycji</strong>. Rozkład rodzajów licencji w tym zamknięciu, odczytany z pól <code>license</code> pliku zamknięcia: MIT — 126 pozycji, Apache-2.0 — 3, ISC — 3, MIT-0 — 2, BSD-2-Clause — 2, BSD-3-Clause — 2, „Apache-2.0 OR MIT” — 1, BlueOak-1.0.0 — 1, CC0-1.0 — 1. Zestawienie nie ujawnia licencji o charakterze wzajemnym (copyleft) w tej warstwie.</p>
<h3 id="74-komponenty-warstwy-powłoki-rust--tauri">7.4 Komponenty warstwy Powłoki (Rust / Tauri)</h3>
<p>Zależności zadeklarowane w <code>budowa/desktop/src-tauri/Cargo.toml</code>, z rozróżnieniem sekcji <code>[dependencies]</code> (zależność produkcyjna, wchodząca do pliku wynikowego) i <code>[build-dependencies]</code> (narzędzie etapu budowania, nieobecne w pliku wynikowym) — rozróżnienie prowadzone konsekwentnie z rozdz. 7.3 dla warstwy Klienta:</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Komponent</th>
<th>Sekcja <code>Cargo.toml</code></th>
<th>Wersja wymagana</th>
<th>Rodzaj licencji</th>
<th>Rola w produkcie</th>
</tr>
</thead>
<tbody>
<tr>
<td><code>tauri</code></td>
<td><code>[dependencies]</code></td>
<td>2 (rozwiązana: 2.11.5)</td>
<td>Apache-2.0 <strong>lub</strong> MIT</td>
<td>powłoka okna, ikona zasobnika, obsługa obrazów ikon</td>
</tr>
<tr>
<td><code>tauri-plugin-dialog</code></td>
<td><code>[dependencies]</code></td>
<td>2 (rozwiązana: 2.7.2)</td>
<td>Apache-2.0 <strong>lub</strong> MIT</td>
<td>okno wyboru katalogu roboczego</td>
</tr>
<tr>
<td><code>serde</code></td>
<td><code>[dependencies]</code></td>
<td>1 (rozwiązana: 1.0.229)</td>
<td>MIT <strong>lub</strong> Apache-2.0</td>
<td>serializacja danych poleceń powłoki</td>
</tr>
<tr>
<td><code>tauri-build</code></td>
<td><strong><code>[build-dependencies]</code></strong> — zależność budowania</td>
<td>2 (rozwiązana: 2.6.3)</td>
<td>Apache-2.0 <strong>lub</strong> MIT</td>
<td>krok budowania powłoki; nie wchodzi do pliku wynikowego</td>
</tr>
</tbody>
</table></div>
<p><strong>Zmiana względem stanu sprzed prac naprawczych.</strong> Zależność bezpośrednia <code>serde_json</code> została usunięta z <code>Cargo.toml</code> w pracach nad ścieżką wydania i bramami mierzącymi wywołania jako nieużywana. <code>serde_json</code> w wersji 1.0.151 nadal występuje w zamknięciu <code>Cargo.lock</code> jako zależność <strong>pośrednia</strong>, wciągana przez <code>tauri</code>; nota licencyjna MIT <strong>lub</strong> Apache-2.0 pozostaje należna z tego tytułu, lecz komponent nie jest już zależnością bezpośrednią warstwy Powłoki i nie figuruje w tabeli powyżej.</p>
<p>Zamknięcie zależności zapisane w <code>Cargo.lock</code> obejmuje <strong>447 pakietów</strong> — jest to konsekwencja architektury Tauri, która wciąga zależności warstwy systemowej dla wielu platform. Rozkład rodzajów licencji ustalony przez odczyt pola <code>license</code> z manifestów pakietów obecnych w lokalnej pamięci podręcznej rejestru (<strong>274 z 447 pakietów zweryfikowane lokalnie; 173 pakiety niezweryfikowane</strong> — w większości pakiety właściwe dla platform innych niż Windows oraz pakiet własny Powłoki):</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Zadeklarowana licencja</th>
<th>Liczba pakietów</th>
</tr>
</thead>
<tbody>
<tr>
<td>MIT OR Apache-2.0 (w różnych zapisach)</td>
<td>144</td>
</tr>
<tr>
<td>MIT</td>
<td>47</td>
</tr>
<tr>
<td>Apache-2.0 OR MIT</td>
<td>31</td>
</tr>
<tr>
<td>Unicode-3.0</td>
<td>18</td>
</tr>
<tr>
<td>Unlicense OR MIT (w różnych zapisach)</td>
<td>11</td>
</tr>
<tr>
<td><strong>MPL-2.0</strong></td>
<td><strong>5</strong></td>
</tr>
<tr>
<td>BSD-3-Clause i warianty z BSD-3-Clause</td>
<td>6</td>
</tr>
<tr>
<td>Zlib i warianty z Zlib</td>
<td>5</td>
</tr>
<tr>
<td>Apache-2.0</td>
<td>2</td>
</tr>
<tr>
<td>pozostałe warianty złożone</td>
<td>5</td>
</tr>
</tbody>
</table></div>
<p><strong>Uwaga o licencjach wzajemnych.</strong> W zamknięciu zależności Powłoki występuje pięć pakietów na licencji <strong>MPL-2.0</strong> (Mozilla Public License 2.0): <code>cssparser</code> 0.36.0, <code>cssparser-macros</code> 0.6.1, <code>dtoa-short</code> 0.3.5, <code>option-ext</code> 0.2.0 oraz <code>selectors</code> 0.36.1. MPL-2.0 jest licencją wzajemną o zasięgu plikowym: obowiązek udostępnienia kodu źródłowego dotyczy plików objętych tą licencją i ich modyfikacji, nie zaś całości produktu, z którym zostały połączone. Licencjodawca nie modyfikuje tych pakietów, w związku z czym obowiązek sprowadza się do wskazania miejsca uzyskania ich kodu źródłowego i zachowania not licencyjnych.</p>
<p><strong>Zastrzeżenie o kompletności.</strong> Wykaz warstwy Rust jest zweryfikowany częściowo. Przed publikacją Wydania konieczne jest sporządzenie pełnego zestawienia licencji zamknięcia zależności właściwego dla budowy docelowej (narzędziem inwentaryzującym uruchomionym na maszynie budującej) i dołączenie go do Wydania.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> inwentarz licencji Komponentów jest <strong>generowany automatycznie przy każdym Wydaniu</strong> narzędziem inwentaryzującym uruchamianym na maszynie budującej i dołączany do Wydania jako plik <code>NOTICE</code>. Zestawienie Załącznika A dokumentuje stan na dzień wydania niniejszej wersji Licencji; w razie rozbieżności między Załącznikiem A a plikiem <code>NOTICE</code> Wydania rozstrzyga plik <code>NOTICE</code> jako zestawienie generowane z rzeczywistego stanu zależności.</p>
<h3 id="75-kroje-pisma">7.5 Kroje pisma</h3>
<p>Warstwa wizualna Oprogramowania opiera się na trzech krojach pisma, których pliki w formacie WOFF2 są <strong>wbudowane w Oprogramowanie i rozpowszechniane wraz z nim</strong>:</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Krój</th>
<th>Rola w produkcie</th>
<th>Licencja</th>
<th>Uprawniony</th>
</tr>
</thead>
<tbody>
<tr>
<td><strong>Space Grotesk</strong></td>
<td>nagłówki, tytuły środowisk, logotyp</td>
<td>SIL Open Font License 1.1</td>
<td>Copyright 2020 The Space Grotesk Project Authors</td>
</tr>
<tr>
<td><strong>IBM Plex Sans</strong></td>
<td>interfejs i treść</td>
<td>SIL Open Font License 1.1</td>
<td>Copyright 2019 IBM Corp.</td>
</tr>
<tr>
<td><strong>IBM Plex Mono</strong></td>
<td>dane techniczne, identyfikatory, terminal</td>
<td>SIL Open Font License 1.1</td>
<td>Copyright 2019 IBM Corp.</td>
</tr>
</tbody>
</table></div>
<p>Kroje osadzono w podzbiorach <code>latin</code> oraz <code>latin-ext</code>, obejmujących pełny zestaw polskich znaków diakrytycznych. Pełne teksty licencji obu rodzin znajdują się w plikach dołączonych do zasobów Oprogramowania (<code>LICENCJA-space-grotesk.txt</code>, <code>LICENCJA-ibm-plex.txt</code>).</p>
<p>Z licencji SIL OFL 1.1 wynikają obowiązki, które Licencjobiorca przyjmuje do wiadomości i zobowiązuje się respektować:</p>
<ol type="1">
<li>pliki krojów mogą być rozpowszechniane wyłącznie razem z Oprogramowaniem albo samodzielnie, <strong>zawsze wraz z tekstem licencji i notą o prawach autorskich</strong>;</li>
<li><strong>zakazana jest sprzedaż samych plików krojów</strong> jako odrębnego towaru;</li>
<li>zmodyfikowana wersja kroju nie może być rozpowszechniana pod nazwą zastrzeżoną (Reserved Font Name) przez autora oryginału;</li>
<li>licencja OFL dotyczy plików krojów, nie zaś dokumentów ani grafik złożonych przy ich użyciu — wytwory Operatora nie stają się przez to objęte OFL.</li>
</ol>
<p>Nieobowiązujące — wyłączone decyzją projektową i niewystępujące w Oprogramowaniu — są kroje Cormorant Garamond, Inter oraz JetBrains Mono.</p>
<h3 id="76-zestaw-ikon-i-zasoby-marki">7.6 Zestaw ikon i zasoby marki</h3>
<p>Zestaw ikon interfejsu opisany jest manifestem warstwy wizualnej w wersji v2.0 (<code>ikony/manifest.json</code>) i obejmuje <strong>82 pozycje</strong> na jednolitej siatce 24×24 z obrysem 1,75 i barwą dziedziczoną z kontekstu (<code>currentColor</code>). Liczba ta wynika wprost z manifestu (pole <code>liczba-ikon</code> o wartości 82 oraz 82 wpisy wykazu) i odpowiada wiązaniom zestawu w kodzie Klienta (13 + 16 + 11 + 14 + 24 + 4 = 82 pozycje w sześciu grupach zastosowań). Źródła pozycji, ustalone z pola <code>zrodlo</code> każdego wpisu manifestu:</p>
<ol type="1">
<li><strong>78 pozycji</strong> wywiedzionych z biblioteki <strong>Lucide</strong>, rozpowszechnianej na licencji <strong>ISC</strong>, z obrysem ujednoliconym do wartości przyjętej w zestawie;</li>
<li><strong>4 pozycje własne</strong> — emblematy Środowisk (<code>srodowisko-talkin</code>, <code>srodowisko-workspace</code>, <code>srodowisko-codestudio</code>, <code>srodowisko-multitaskingai</code>) wykonane w siatce i kresce zestawu; stanowią utwór Licencjodawcy.</li>
</ol>
<p>Licencja ISC jest licencją zezwalającą; jej warunkiem jest zachowanie noty o prawach autorskich i treści zezwolenia we wszystkich kopiach.</p>
<p><strong>Uzgodnienie z liczbą 83 pojawiającą się w katalogu plików zestawu.</strong> Katalog plików SVG zestawu (<code>client/src/ikony/svg/</code>) zawiera <strong>83 pliki</strong>, lecz jednym z nich jest <code>logo-danaco.svg</code> — znak marki wczytywany w kodzie odrębnie i <strong>nienależący do zestawu ikon</strong> ani do wykazu manifestu. Dla celów licencyjnych rozstrzyga podział rzeczowy, nie liczba plików katalogu: <strong>82 pozycje zestawu ikon</strong> (78 na licencji ISC, 4 własne) oraz <strong>znak marki</strong> podlegający rozdz. 6.4 wraz z pozostałymi zasobami marki.</p>
<p>Redakcja niniejszej Licencji zweryfikowała bezpośrednio treść <code>README.md</code> i <code>INSTRUKCJA-UZYTKOWANIA.md</code> dołączonych do Egzemplarza: oba dokumenty posługują się liczbą <strong>82 pozycji zestawu ikon</strong> i wprost objaśniają czytelnikowi ten sam podział rzeczowy — 82 pozycje wykazu manifestu wobec 83 plików katalogu z uwzględnieniem znaku marki. Wcześniejsza redakcja niniejszego rozdziału odnotowywała rozbieżność liczby 83 przypisywaną obu tym dokumentom; rozbieżność ta nie istnieje w treści dokumentów na dzień sporządzenia niniejszej wersji Licencji i została usunięta z treści rozdziału jako nieaktualna. Kwestia jednolitej liczby pozycji zestawu ikon w dokumentacji produktu nie wymaga już rozstrzygnięcia Operatora — wszystkie trzy dokumenty (niniejsza Licencja, <code>README.md</code>, <code>INSTRUKCJA-UZYTKOWANIA.md</code>) zgodnie posługują się liczbą wynikającą z manifestu jako źródła prawdy warstwy wizualnej, a odpowiadająca temu pozycja została usunięta z Załącznika B.</p>
<p><strong>Zasoby marki</strong> — logotypy, sygnety i ich warianty dla motywu jasnego i ciemnego (dwanaście plików katalogu zasobów marki) oraz znak <code>logo-danaco.svg</code> — stanowią utwór Licencjodawcy i nie podlegają licencjom komponentów. Ich użycie reguluje rozdz. 6.4.</p>
<h3 id="77-program-zewnętrzny-kanału-głównego">7.7 Program zewnętrzny kanału głównego</h3>
<p>Kanał modelu rodzaju <code>cli</code> uruchamia zewnętrzny program <code>claude</code> (Claude Code CLI). Program ten <strong>nie jest częścią Oprogramowania i nie jest rozpowszechniany przez Licencjodawcę</strong>. Oprogramowanie:</p>
<ol type="1">
<li>buduje argumenty wywołania warunkowo, zależnie od konfiguracji Okna komunikacji i Nakładki tożsamości (<code>--system-prompt</code> albo <code>--append-system-prompt</code>, <code>--permission-mode</code>, <code>--add-dir</code>, <code>--mcp-config</code>);</li>
<li>ustawia dla procesu zmienną <code>CLAUDE_CONFIG_DIR</code> wskazującą katalog konfiguracji Konta;</li>
<li>odczytuje strumień odpowiedzi w formacie <code>stream-json</code> i przekazuje go Klientowi w postaci fragmentów.</li>
</ol>
<p>Korzystanie z programu <code>claude</code> oraz z usług Dostawcy modelu podlega wyłącznie warunkom umownym tego Dostawcy. Licencjobiorca zapewnia legalność korzystania z programu i konta.</p>
<h3 id="78-rozszerzenia">7.8 Rozszerzenia</h3>
<p>Przyjęta decyzja architektoniczna przewiduje katalog rozszerzeń (przeglądanie, instalacja, zarządzanie, kategorie) oparty na tabeli <code>rozszerzenie</code> i rodzinie komend <code>extension.*</code>, przy czym marketplace jako model wydawniczy (publikowanie i dystrybucja rozszerzeń do innych operatorów) został odłożony do odrębnego opracowania.</p>
<p><strong>W wersji v2.0 mechanizm rozszerzeń nie istnieje</strong>: rodzina komend <code>extension.*</code> nie występuje w Kontrakcie wykonawczym, a Oprogramowanie nie udostępnia drogi instalowania ani uruchamiania rozszerzeń. Funkcjonalność ta jest przewidziana koncepcją i w bieżącej wersji niezaimplementowana.</p>
<p>W konsekwencji Licencja <strong>nie reguluje</strong> zasad tworzenia, dystrybucji ani korzystania z rozszerzeń.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> do czasu udostępnienia mechanizmu rozszerzeń wraz z odrębnymi zasadami ich licencjonowania rozszerzenia mogą pochodzić <strong>wyłącznie od Licencjodawcy</strong>, a instalowanie w Oprogramowaniu komponentów rozszerzających pochodzących od osób trzecich pozostaje niedozwolone na podstawie rozdz. 5.2 pkt 5. Zasady dotyczące rozszerzeń osób trzecich — w tym licencja API, odpowiedzialność za rozszerzenie i tryb weryfikacji — zostaną wprowadzone nową wersją dokumentu w trybie rozdz. 8.6 przed udostępnieniem mechanizmu.</p>
<h3 id="79-obowiązki-dokumentacyjne-przy-dystrybucji">7.9 Obowiązki dokumentacyjne przy dystrybucji</h3>
<p>Licencjodawca zobowiązuje się dołączać do każdego Wydania Oprogramowania zestawienie Komponentów wraz z pełnymi tekstami ich licencji oraz notami o prawach autorskich, w postaci pliku not licencyjnych dostępnego Licencjobiorcy wraz z Egzemplarzem.</p>
<p>Licencjobiorca, rozpowszechniając Oprogramowanie w przypadkach, w których rozpowszechnianie jest dopuszczalne na podstawie odrębnej zgody, zobowiązany jest przekazać ten plik w postaci niezmienionej.</p>
<hr />
<h2 id="8-aktualizacje-i-wersjonowanie">8. Aktualizacje i wersjonowanie</h2>
<h3 id="81-zasada-wersjonowania">8.1 Zasada wersjonowania</h3>
<p>Oprogramowanie podlega zasadzie wersjonowania ustalonej decyzją architektoniczną: <strong>oznaczenie v2.0 obowiązuje do pierwszej publikacji, w trakcie budowy zakazane jest podnoszenie numeru wersji, a status wersji brzmi „Deweloperski”</strong>. W konsekwencji numer wersji Oprogramowania <strong>nie jest miernikiem zakresu funkcjonalnego</strong> ani stopnia dojrzałości: dwa Wydania oznaczone tym samym numerem v2.0 mogą różnić się zakresem działających funkcji.</p>
<p>Licencjobiorca przyjmuje do wiadomości, że identyfikacja Wydania następuje przez oznaczenie Wydania nadane przez Licencjodawcę, nie przez numer wersji produktu.</p>
<p>Numer <code>1.0.0</code> występujący w metadanych pakietu Powłoki, w manifeście pakietu Klienta oraz w manifeście pakietu Rust jest odwzorowaniem wersji produktu w rozumieniu powyższej zasady i podlega tej samej uwadze.</p>
<h3 id="82-charakter-aktualizacji">8.2 Charakter Aktualizacji</h3>
<ol type="1">
<li>Aktualizacja może obejmować poprawki błędów, zmiany zakresu funkcjonalnego, zmiany interfejsu, zmiany Kontraktu oraz zmiany schematu trwałości.</li>
<li>Ze względu na status Deweloperski Licencjodawca <strong>zastrzega prawo do zmian niezachowujących zgodności wstecznej</strong>, w szczególności do zmiany nazw komend i zdarzeń Kontraktu, zmiany kluczy konfiguracji oraz zmiany struktury tabel.</li>
<li>Aktualizacja może usunąć funkcję dostępną w Wydaniu poprzednim, jeżeli funkcja ta okazała się niezgodna z przyjętą architekturą albo nie została doprowadzona do stanu użytecznego.</li>
</ol>
<h3 id="83-migracje-schematu-trwałości">8.3 Migracje schematu trwałości</h3>
<ol type="1">
<li>Schemat Bazy jest tworzony i aktualizowany przez migracje wkompilowane w Rdzeń i stosowane samoczynnie przy jego starcie, w jednej transakcji, z kontrolą sumy kontrolnej kroku.</li>
<li><strong>Migracje działają wyłącznie w kierunku naprzód.</strong> Oprogramowanie nie udostępnia mechanizmu cofnięcia migracji ani przywrócenia schematu do stanu poprzedniego. Po uruchomieniu nowszego Wydania na istniejącym Katalogu danych powrót do Wydania wcześniejszego może być niemożliwy.</li>
<li>Przed instalacją Aktualizacji Licencjobiorca zobowiązany jest wykonać kopię zapasową Katalogu danych. Zaniechanie tej czynności obciąża Licencjobiorcę; Licencjodawca nie odpowiada za skutki niemożności powrotu do Wydania wcześniejszego.</li>
</ol>
<h3 id="84-uprawnienie-do-aktualizacji-i-sposób-dostarczania">8.4 Uprawnienie do Aktualizacji i sposób dostarczania</h3>
<p>Oprogramowanie <strong>nie zawiera mechanizmu samoczynnej aktualizacji</strong>: Powłoka nie korzysta z komponentu aktualizatora, a Rdzeń nie nawiązuje połączeń z serwerami Licencjodawcy w celu sprawdzenia dostępności nowego Wydania. Aktualizacja następuje wyłącznie przez świadome działanie Licencjobiorcy — pobranie i zainstalowanie Wydania udostępnionego przez Licencjodawcę.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> w okresie statusu Deweloperskiego Aktualizacje udostępnione przez Licencjodawcę wchodzą w zakres Licencji <strong>nieodpłatnie</strong>, przy czym Licencjodawca <strong>nie zobowiązuje się</strong> do ich wydawania, do utrzymywania Wydań wcześniejszych ani do usuwania podatności w oznaczonym terminie (rozdz. 8.5). Wsparcie techniczne świadczone jest w miarę możliwości, drogą poczty elektronicznej na adres wskazany w rozdz. 14.7, bez gwarantowanych czasów reakcji. Zobowiązania serwisowe mogą zostać wprowadzone wraz z odpłatnością (rozdz. 4.5), nową wersją dokumentu w trybie rozdz. 8.6.</p>
<h3 id="85-brak-zobowiązania-do-rozwoju">8.5 Brak zobowiązania do rozwoju</h3>
<p>Licencjodawca <strong>nie zaciąga zobowiązania</strong> do wydawania Aktualizacji, rozwijania Oprogramowania, doprowadzania do stanu działającego funkcji przewidzianych koncepcją ani do świadczenia wsparcia technicznego wykraczającego poza zakres opisany w rozdz. 8.4. Postanowienie to jest bezpośrednią konsekwencją statusu Deweloperskiego i nieodpłatnego charakteru udostępnienia w tej fazie (rozdz. 4.5); może ulec zmianie wraz z wprowadzeniem odpłatności, w trybie rozdz. 8.6.</p>
<h3 id="86-zmiany-treści-licencji">8.6 Zmiany treści Licencji</h3>
<p>Licencjodawca może wydać nową wersję dokumentu Licencji wraz z nowym Wydaniem Oprogramowania. Nowa wersja Licencji obowiązuje wobec Wydań udostępnionych po jej wydaniu. Do Wydania już posiadanego przez Licencjobiorcę stosuje się wersję Licencji dostarczoną wraz z tym Wydaniem, chyba że strony postanowią inaczej.</p>
<p>Zmiana Licencji nie może pozbawić Licencjobiorcy uprawnień nabytych w stosunku do Wydania już posiadanego.</p>
<hr />
<h2 id="9-dane-i-treści-operatora">9. Dane i treści Operatora</h2>
<h3 id="91-miejsce-przechowywania-danych">9.1 Miejsce przechowywania danych</h3>
<p>Całość trwałości Oprogramowania mieści się w <strong>jednym pliku bazy SQLite</strong> o nazwie <code>danaco-console.db</code>, umieszczonym w Katalogu danych. Katalogiem danych jest domyślnie <code>%LOCALAPPDATA%\\DanacoConsole</code>; wskazanie innego katalogu następuje zmienną środowiska <code>DANACO_KATALOG_DANYCH</code> albo argumentem wywołania <code>--dane</code>. Zmiana Katalogu danych przenosi całość trwałości — nie istnieje drugie, niezależne miejsce zapisu.</p>
<p>Dane Operatora pozostają na urządzeniu, na którym działa Rdzeń. <strong>Licencjodawca nie otrzymuje kopii Danych Operatora w związku z korzystaniem z Oprogramowania.</strong></p>
<h3 id="92-zakres-danych-zapisywanych-przez-oprogramowanie">9.2 Zakres danych zapisywanych przez Oprogramowanie</h3>
<p>W Bazie zapisywane są w szczególności:</p>
<ol type="1">
<li>sesje, karty sesji i okna komunikacji wraz z ich ustawieniami (moduł, kanał modelu, katalogi robocze, środowisko wykonania, tryb uprawnień, rola). <strong>Moment powstania wiersza wymaga zastrzeżenia:</strong> wiersz sesji i wiersz okna komunikacji <strong>nie powstają w chwili utworzenia sesji ani okna</strong> (<code>session.create</code>, <code>window.create</code>), lecz dopiero przy zapisie pierwszej wiadomości w danym oknie; ponadto zmiana ustawień okna (<code>window.update</code>) <strong>nie jest utrwalana</strong> — obowiązuje wyłącznie w pamięci Rdzenia do czasu jego zatrzymania (rozdz. 10.3 wykaz ograniczeń, pkt 17). Sesja i okno, które nie doszły do pierwszej wiadomości, nie zostaną odtworzone po restarcie Rdzenia i nie znajdą się w kopii zapasowej Katalogu danych (rozdz. 9.5, rozdz. 11.4 pkt 3);</li>
<li><strong>treść wiadomości Operatora oraz treść odpowiedzi modelu</strong> — w tabeli wiadomości powiązanej łańcuchem więzów sesja → karta → okno → wiadomość;</li>
<li>wartości ustawień konfiguracji na przewidzianych poziomach zasięgu i osiach rozstrzygania — z zastrzeżeniem, że zapis wartości nie jest równoznaczny z jej stosowaniem przy wywołaniu modelu (rozdz. 10.3 wykaz ograniczeń, pkt 21);</li>
<li>definicje kanałów modelu, kont (bez treści Poświadczeń), punktów dostępu i nadań, dokumentów tożsamości modelu.</li>
</ol>
<p><strong>Czego Baza nie zapisuje.</strong> Poza treścią Poświadczeń (rozdz. 9.3) Baza nie zapisuje <strong>danych prowenancji wywołań modelu</strong>: struktura prowenancji powstaje przed uruchomieniem procesu Tury i jest przekazywana do interfejsu, lecz nie ma dla niej ani tabeli, ani kolumny w schemacie trwałości. Po zamknięciu Tury nie da się z Bazy odtworzyć, co faktycznie zostało przekazane modelowi — dostępna pozostaje wyłącznie treść wiadomości i odpowiedzi (rozdz. 2.3, rozdz. 10.3 wykaz ograniczeń, pkt 20). Okoliczność ta ma znaczenie dla rozliczalności pracy wykonanej z udziałem modelu i dla odtwarzania przebiegu zdarzeń po fakcie.</p>
<p><strong>Uwaga o treściach obszernych.</strong> Przyjęta decyzja architektoniczna przewiduje przechowywanie treści obszernych jako plików poza Bazą, z zapisem odwołania w Bazie. Mechanizm ten jest przewidziany koncepcją i <strong>w wersji v2.0 niezaimplementowany</strong> — pole odwołania do treści nie jest wypełniane, a treść zapisywana jest w całości w Bazie. Okoliczność ta ma znaczenie dla planowania rozmiaru Katalogu danych i dla procedur usuwania danych.</p>
<h3 id="93-postępowanie-z-poświadczeniami">9.3 Postępowanie z Poświadczeniami</h3>
<p>Zgodnie z przyjętą architekturą Oprogramowanie <strong>nie przechowuje treści Poświadczeń w Bazie</strong>. Baza przechowuje wyłącznie odwołania — nazwę wpisu w magazynie sekretów, ścieżkę katalogu konfiguracji Konta albo ścieżkę klucza. Kontrakt przyjmuje treść poświadczenia wyłącznie w żądaniach dodania i aktualizacji Konta oraz Punktu dostępu; żadna odpowiedź ani żadne zdarzenie nie zwraca tej treści — zwracana jest wyłącznie informacja o obecności poświadczenia albo jego odwołanie.</p>
<p>Odpowiedzialność za zabezpieczenie miejsc, w których faktyczne Poświadczenia się znajdują (katalogi konfiguracji kont, pliki kluczy, magazyn sekretów systemu operacyjnego), spoczywa na Licencjobiorcy.</p>
<h3 id="94-przekazywanie-danych-na-zewnątrz">9.4 Przekazywanie danych na zewnątrz</h3>
<p>Oprogramowanie przekazuje dane poza urządzenie wyłącznie w następstwie czynności Operatora i wyłącznie w następujących sytuacjach:</p>
<ol type="1">
<li><strong>wywołanie modelu kanałem rodzaju <code>cli</code></strong> — treść wiadomości, Nakładka tożsamości i parametry wywołania są przekazywane uruchamianemu lokalnie programowi zewnętrznemu, który we własnym zakresie komunikuje się z usługą Dostawcy modelu;</li>
<li><strong>wywołanie modelu kanałem rodzaju <code>api</code></strong> — treść zapytania jest wysyłana na adres wskazany wierszem rejestru kanału, przy użyciu poświadczenia wskazanego konfiguracją. Droga ta jest w wersji v2.0 <strong>praktycznie nieosiągalna</strong>: Oprogramowanie nie udostępnia interfejsu zakładania kanałów (rozdz. 3.4 pkt 3), a Konta rodzaju <code>api</code> nie mają drogi wykonawczej (rozdz. 10.3 wykaz ograniczeń, pkt 19). Jej uruchomienie wymagałoby ręcznego zasiania tabeli <code>kanal_modelu</code> poza Oprogramowaniem albo wydania komendy <code>channel.add</code> poza interfejsem aplikacji. Ujawnienie tej drogi ma charakter pełnego opisu możliwych wyjść danych, nie zaś opisu funkcji działającej;</li>
<li><strong>most MCP</strong> — konfiguracja mostu przekazana procesowi modelu może powodować nawiązanie połączenia z maszyną wskazaną Punktem dostępu, w tym połączenia z użyciem programu <code>ssh</code> obecnego w systemie operacyjnym.</li>
</ol>
<p>Poza powyższymi przypadkami Oprogramowanie nie nawiązuje połączeń wychodzących. <strong>Oprogramowanie nie zawiera telemetrii przekazywanej Licencjodawcy, nie raportuje zdarzeń użycia, nie sprawdza dostępności aktualizacji i nie wysyła zgłoszeń o błędach.</strong></p>
<p>Operator przyjmuje do wiadomości, że treść przekazana modelowi opuszcza kontrolę Oprogramowania i podlega dalej zasadom przetwarzania stosowanym przez Dostawcę modelu. Ocena dopuszczalności przekazania określonej kategorii danych (w tym danych osobowych, tajemnicy przedsiębiorstwa, informacji poufnych osób trzecich) należy do Licencjobiorcy i jest dokonywana <strong>przed</strong> przekazaniem.</p>
<h3 id="95-kopie-zapasowe-i-usuwanie-danych">9.5 Kopie zapasowe i usuwanie danych</h3>
<ol type="1">
<li>Licencjobiorca odpowiada za wykonywanie kopii zapasowych Katalogu danych. Oprogramowanie <strong>nie wykonuje kopii zapasowych samoczynnie</strong> i nie zawiera mechanizmu przywracania danych z kopii. <strong>Zakres kopii wyznacza zakres utrwalenia:</strong> kopia Katalogu danych obejmuje wyłącznie to, co zostało zapisane w Bazie. Nie obejmuje zatem sesji i okien, które nie doszły do pierwszej wiadomości, zmian ustawień okna dokonanych komendą <code>window.update</code>, ani danych prowenancji wywołań modelu — żaden z tych elementów nie jest utrwalany (rozdz. 9.2 pkt 1, rozdz. 10.3 wykaz ograniczeń, pkt 17 i 20).</li>
<li>Usunięcie Katalogu danych powoduje nieodwracalną utratę całości trwałości — sesji, okien, historii wiadomości, ustawień, definicji kanałów, kont, punktów dostępu i nadań.</li>
<li>Odinstalowanie Oprogramowania nie usuwa Katalogu danych; jego usunięcie wymaga odrębnej czynności Licencjobiorcy.</li>
<li>Oprogramowanie nie udostępnia w wersji v2.0 funkcji eksportu całości Danych Operatora do formatu przenośnego. Dostęp do danych poza Oprogramowaniem możliwy jest przez odczyt pliku Bazy narzędziami zgodnymi z SQLite.</li>
</ol>
<h3 id="96-dane-osobowe">9.6 Dane osobowe</h3>
<p>Jeżeli Dane Operatora obejmują dane osobowe, administratorem tych danych jest Licencjobiorca. Licencjodawca nie przetwarza tych danych w związku z korzystaniem przez Licencjobiorcę z Oprogramowania, ponieważ nie uzyskuje do nich dostępu (rozdz. 9.1 i 9.4).</p>
<p>Licencjobiorca zobowiązuje się:</p>
<ol type="1">
<li>przetwarzać dane osobowe przy użyciu Oprogramowania zgodnie z przepisami, w szczególności z rozporządzeniem (UE) 2016/679;</li>
<li>przed przekazaniem danych osobowych modelowi ocenić podstawę prawną tego przekazania oraz warunki przetwarzania stosowane przez Dostawcę modelu, w tym miejsce przetwarzania i ewentualne przekazanie poza Europejski Obszar Gospodarczy;</li>
<li>wdrożyć adekwatne środki techniczne i organizacyjne zabezpieczające urządzenie, na którym działa Rdzeń, z uwzględnieniem faktu, że uwierzytelnianie dostępu do Rdzenia jest w fazie budowy wyłączone (rozdz. 5.5 pkt 3).</li>
</ol>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> postanowienia rozdz. 9.6 dotyczą instalacji z <strong>Rdzeniem lokalnym</strong>. Uruchomienie topologii z Rdzeniem serwerowym prowadzonym przez Licencjodawcę albo podmiot z nim powiązany wymaga — <strong>przed rozpoczęciem przetwarzania</strong> — zawarcia odrębnej umowy obejmującej powierzenie przetwarzania danych osobowych w rozumieniu art. 28 RODO, wraz z określeniem prowadzącego serwer, lokalizacji przetwarzania, okresu retencji i zasad usuwania danych. Do czasu zawarcia takiej umowy udostępnienie Rdzenia serwerowego danemu Licencjobiorcy nie następuje.</p>
<h3 id="97-dostęp-licencjodawcy-do-danych-przy-czynnościach-wsparcia">9.7 Dostęp Licencjodawcy do danych przy czynnościach wsparcia</h3>
<p>Jeżeli Licencjobiorca zwróci się o wsparcie techniczne, przekazanie Licencjodawcy plików diagnostycznych, wycinków Bazy albo zrzutów interfejsu następuje wyłącznie z inicjatywy Licencjobiorcy i na jego odpowiedzialność. Licencjobiorca zobowiązany jest usunąć z przekazywanych materiałów dane, których przekazanie byłoby niedopuszczalne. Licencjodawca wykorzystuje przekazane materiały wyłącznie w celu udzielenia wsparcia i usuwa je po jego zakończeniu.</p>
<hr />
<h2 id="10-ograniczenie-rękojmi-i-odpowiedzialności">10. Ograniczenie rękojmi i odpowiedzialności</h2>
<h3 id="101-charakter-świadczenia">10.1 Charakter świadczenia</h3>
<p>Oprogramowanie jest udostępniane <strong>w stanie, w jakim się znajduje</strong> („as is”), z zastrzeżeniem statusu Deweloperskiego wersji v2.0. Licencjodawca udostępnia produkt w toku budowy, o niedomkniętym zakresie funkcjonalnym, którego rzeczywisty stan wykonania został ustalony audytem i ujawniony w niniejszym dokumencie.</p>
<h3 id="102-wyłączenie-rękojmi">10.2 Wyłączenie rękojmi</h3>
<p>W najszerszym zakresie dopuszczalnym przez prawo właściwe Licencjodawca wyłącza rękojmię za wady Oprogramowania oraz nie udziela żadnych gwarancji, w szczególności nie zapewnia, że:</p>
<ol type="1">
<li>Oprogramowanie będzie działać nieprzerwanie, bez błędów i bez przestojów;</li>
<li>Oprogramowanie będzie przydatne do określonego celu zamierzonego przez Licencjobiorcę;</li>
<li>wyniki pracy modelu uzyskane przy użyciu Oprogramowania będą poprawne, kompletne, aktualne albo nadające się do wykorzystania bez weryfikacji;</li>
<li>dane zapisane w Bazie nie ulegną utracie ani uszkodzeniu;</li>
<li>funkcje przewidziane koncepcją produktu zostaną doprowadzone do stanu działającego.</li>
</ol>
<p>Wyłączenie rękojmi nie dotyczy przypadków, w których przepis bezwzględnie obowiązujący nie dopuszcza takiego wyłączenia, ani szkody wyrządzonej umyślnie.</p>
<h3 id="103-ujawnienie-rzeczywistego-stanu-wykonania">10.3 Ujawnienie rzeczywistego stanu wykonania</h3>
<p>Licencjodawca ujawnia poniżej stan wykonania wersji v2.0 w zakresie istotnym dla oceny przydatności Oprogramowania. Ujawnienie ma charakter oświadczenia o stanie przedmiotu świadczenia i wyłącza możliwość powołania się przez Licencjobiorcę na nieświadomość tych ograniczeń.</p>
<p><strong>Funkcje potwierdzone jako działające:</strong></p>
<ol type="1">
<li>transport komunikatów po WebSocket wraz z kopertą Kontraktu, odpornością na nieznaną komendę i przetrwaniem pracy Rdzenia mimo rozłączenia Klienta;</li>
<li>obsługa wszystkich komend Kontraktu po stronie Rdzenia;</li>
<li>trwałość w pliku SQLite wraz z samoczynnym stosowaniem migracji i odtwarzaniem sesji oraz okien po restarcie Rdzenia;</li>
<li>kanał główny modelu (uruchomienie programu zewnętrznego, katalog konfiguracji per Konto, odczyt strumienia <code>stream-json</code>, rozpoznanie wyczerpania limitu, rotacja Kont w obrębie puli, nadanie prowenancji przed startem procesu);</li>
<li>rejestr kanałów modelu sterowany danymi wraz z odświeżeniem po zmianie;</li>
<li>okno konfiguracji ze sterowanym danymi katalogiem kategorii, definicji i opcji oraz zapisem wartości na poziomach zasięgu i osiach — <strong>z istotnym zastrzeżeniem</strong>: działa katalog, rozstrzyganie i zapis wartości, natomiast <strong>większość zapisanych kluczy nie jest odczytywana przy wywołaniu modelu</strong> (wykaz ograniczeń, pkt 21). Zapis wartości nie jest równoznaczny z jej stosowaniem;</li>
<li>zarządzanie kontami, punktami dostępu, nadaniami i tożsamością modelu z poziomu interfejsu — <strong>z zastrzeżeniem</strong>, że nie da się zapisać Punktu dostępu rodzaju <code>localDirectory</code> (wykaz ograniczeń, pkt 22); funkcja działa w zakresie kont, tożsamości modelu, punktów dostępu rodzaju <code>mcpBridge</code> i nadań;</li>
<li>Nakładka tożsamości w trybach <code>ZASTAP</code> i <code>DOLACZ</code> przekazywana procesowi modelu;</li>
<li>zapis historii rozmowy do Bazy;</li>
<li>warstwa wizualna: dwa motywy, zwarta gęstość, zestaw ikon;</li>
<li>zatrzymanie bieżącej tury.</li>
</ol>
<p><strong>Ograniczenia istotne — funkcje przewidziane koncepcją, w wersji v2.0 niedziałające albo działające częściowo:</strong></p>
<ol type="1">
<li><strong>Świeża instalacja nie wykona wywołania modelu bez czynności przygotowawczych</strong> — brak zaczynu rejestru kanałów, pusta pula kont odmawia wywołania, a Klient przekazuje rodzaj kanału zamiast identyfikatora wiersza rejestru (rozdz. 3.4 pkt 3).</li>
<li><strong>Brak ciągłości rozmowy</strong> — każda tura uruchamia nowy proces bez historii; model nie widzi poprzednich wiadomości Okna komunikacji.</li>
<li><strong>Historia rozmowy nie jest odtwarzana w interfejsie</strong> — mimo zapisu do Bazy Klient nie pobiera zapisanych wiadomości; po odświeżeniu widoku historia znika z ekranu.</li>
<li><strong>Okna operacyjne modułów nie istnieją</strong> — z siedemdziesięciu dziewięciu wierszy katalogu okien operacyjnych w Bazie, będących liczbą obowiązującą w całej dokumentacji produktu (<code>INSTRUKCJA-UZYTKOWANIA.md</code> rozdz. 7), zbudowany jest <strong>jeden widok</strong> — okno komunikacji, któremu odpowiada piętnaście wierszy katalogu, po jednym na moduł; wybór modułu w nawigacji nie przeładowuje przestrzeni roboczej. Wszystkie piętnaście modułów (Studio, Workspace, Browser, Research, Library, Translate, Roundtable, Design, Assistant, Terminal, Developer, Diagnostics, Apps, Agents, Automations) występuje wyłącznie jako pozycje nawigacji.</li>
<li><strong>Środowiska TalkIn, WorkSpace, CodeStudio i MultitaskingAI</strong> mają nawigację i opis, lecz nie mają realizacji funkcjonalnej; środowisko MultitaskingAI nie posiada panelu orkiestracji, ról, kolejek ról ani monitora.</li>
<li><strong>Pulpit dowodzenia (Mission Control)</strong> — dane przykładowe wbudowane w interfejs (plik <code>dane-przykladowe.ts</code>) zostały usunięte; pulpit rysuje się obecnie na podstawie stanu Rdzenia, odczytywanego komendami <code>session.list</code>, <code>window.list</code> i <code>channel.list</code> oraz podtrzymywanego subskrypcją zdarzeń <code>session.changed</code>, <code>window.changed</code>, <code>queue.changed</code> i <code>progress.changed</code>, z uczciwymi stanami pustymi tam, gdzie Rdzeń nie zwrócił jeszcze danych. Zmiana ta nie obejmuje zamiarów: zamiary wyrażone na pulpicie (kolejka, utworzenie zasobu, decyzja) nadal <strong>nie wykonują komend Kontraktu</strong> — zapisują się wyłącznie do lokalnego paska działań pulpitu, bez wywołania kanału.</li>
<li><strong>Kolejki</strong> — komendy kolejek zmieniają wyłącznie stan wiersza; pozycje kolejki nie powstają, silnik wykonania nie istnieje.</li>
<li><strong>Przekazanie zlecenia koordynator → wykonawca w interfejsie</strong> jest wyłącznie animacją lokalną, bez wywołania komendy Kontraktu. Sama pętla koordynator–wykonawca istnieje w Rdzeniu i uruchamia się wraz z turą.</li>
<li><strong>Wybór modelu, modelu zapasowego i nakładu rozumowania na poziomie Okna komunikacji</strong> — wartości zapisują się pod kluczami nierozpoznawanymi przez Rdzeń i nie wpływają na wywołanie.</li>
<li><strong>Zakres wykonania i host wykonania</strong> — wyłącznie etykieta w prowenancji; proces modelu startuje zawsze na maszynie Rdzenia (rozdz. 4.7 pkt 3).</li>
<li><strong>Izolacja</strong> — jedenaście ustawień rozstrzyganych na ośmiu poziomach zasięgu, bez egzekutorów (rozdz. 5.4 pkt 4).</li>
<li><strong>Sterowanie platformą przez narzędzia modelu</strong> — trzydzieści dziewięć deklaracji narzędzi jest generowanych z Kontraktu, lecz nie ma konsumenta; model nie steruje platformą.</li>
<li><strong>Nawigacja platformy przez Kontrakt</strong> — komendy strony głównej, środowisk, modułów i przestrzeni roboczej są zaimplementowane po obu stronach, lecz w większości nie są wywoływane przez żaden widok. <strong>Wyjątek</strong>: <code>session.list</code>, <code>session.bind</code> i <code>session.focus</code> są wywoływane przez stronę główną w strefie „Sesje w tle” — odczyt sesji trwających na Rdzeniu, powrót do sesji przez powiązanie połączenia (<code>session.bind</code>) i przeniesienie ogniska karty czynnej (<code>session.focus</code>), dostępne wyłącznie po ustaleniu tożsamości klienta z powitania połączenia.</li>
<li><strong>Pamięć wielopoziomowa</strong> — istnieje warstwa danych, brak komend Kontraktu i konsumentów.</li>
<li><strong>Rozszerzenia i katalog rozszerzeń</strong> — nie występują w Kontrakcie wykonawczym (rozdz. 7.8).</li>
<li><strong>Telemetria postępu w interfejsie</strong> — zdarzenie <code>progress.changed</code> <strong>ma już konsumenta</strong>: źródło danych pulpitu dowodzenia subskrybuje je wraz z <code>session.changed</code>, <code>window.changed</code> i <code>queue.changed</code> i wprowadza jego treść do składania kompletu danych pulpitu. Zapytanie o stan okna po stronie Rdzenia nadal zwraca zawsze stan oczekujący, dopóki dla tego okna nie istnieje proces Tury — ograniczenie to dotyczy odczytu stanu okna niezależnie od tego, czy zdarzenie postępu ma subskrybenta po stronie Klienta.</li>
<li><strong>Wiersz sesji i wiersz okna powstają dopiero przy pierwszej wiadomości</strong>, nie w chwili <code>session.create</code> ani <code>window.create</code>; <strong>zmiana ustawień okna</strong> (<code>window.update</code>) nie jest utrwalana w Bazie, lecz obowiązuje wyłącznie w pamięci Rdzenia do jego zatrzymania; sesje trafiają zawsze do pierwszego środowiska. Skutki dla zakresu utrwalenia opisuje rozdz. 9.2 pkt 1, dla kopii zapasowych — rozdz. 9.5, dla stanu po zakończeniu Licencji — rozdz. 11.4 pkt 3.</li>
<li><strong>Instalator NSIS Powłoki nie powstaje.</strong> Na rewizji audytu pierwotnego stanu sprzed prac naprawczych konfiguracja Powłoki (<code>tauri.conf.json</code>) nie zawierała kroku poprzedzającego budowanie (<code>beforeBuildCommand</code>) ani adresu trybu rozwojowego (<code>devUrl</code>), a repozytorium nie zawierało skryptu składającego katalog produktu — złożenie katalogu <code>C:\\DanacoConsole_App</code> było wówczas czynnością ręczną. <strong>Stan ten uległ zmianie w pracach nad ścieżką wydania</strong>: <code>tauri.conf.json</code> zawiera obecnie <code>beforeBuildCommand</code> i <code>devUrl</code>; repozytorium zawiera <code>budowa/scripts/wydanie.sh</code> i <code>budowa/scripts/pakowanie.sh</code>, którymi jednym poleceniem składa się produkt (rdzeń <code>danaco-console.exe</code>, powłoka <code>Danaco Console.exe</code>, <code>client/dist</code>) do katalogu <code>C:\\DanacoConsole_App</code> — z tej właśnie ścieżki zbudowano i uruchomiono Egzemplarz, do którego dołączona jest niniejsza Licencja. Punkt niniejszy, dawniej ujawniający brak jakiejkolwiek automatyzacji budowania, pozostaje aktualny w zakresie węższym: <strong>wygenerowanie instalatora Windows w formacie NSIS nadal nie następuje</strong> — krok <code>cargo tauri build</code> wytwarzający instalator wymaga obecności narzędzia NSIS na maszynie budującej; ścieżka wydania pomija budowę Powłoki automatycznie przy braku <code>cargo</code> lub <code>tauri-cli</code>, a sam instalator pozostaje krokiem odrębnym, udokumentowanym w skrypcie wydania jako wymagający NSIS, którego na maszynie audytowanej brak. Punkt niniejszy podlega ponownej weryfikacji przy każdym Wydaniu (Załącznik C.4).</li>
<li><strong>Zmiany definicji Kont wymagają ponownego uruchomienia Rdzenia</strong> — pula rotacji budowana jest jednorazowo przy starcie; stan wyczerpania limitu nie jest utrwalany. Konta rodzaju <code>api</code> i <code>sdk</code> nie mają drogi wykonawczej (skutki dla kanału rodzaju <code>api</code> opisuje rozdz. 9.4 pkt 2).</li>
<li><strong>Prowenancja nie jest utrwalana.</strong> Struktura prowenancji powstaje przed uruchomieniem procesu Tury i jest przekazywana do interfejsu, lecz w całym schemacie trwałości nie ma dla niej tabeli ani kolumny — w szczególności nie ma jej tabela wiadomości, której polami treściowymi są: rola, persona, okno źródłowe, rodzaj treści, stan, treść, odwołanie do treści, liczniki tokenów wejścia i wyjścia, kolejność i czas utworzenia — obok klucza głównego <code>id</code> i klucza obcego <code>okno_komunikacji_id</code> wiążącego wiadomość z Oknem komunikacji (<code>budowa/server/internal/store/migracja_002_okna.sql</code>). Wyliczenie to obejmuje pola treściowe tabeli, nie stanowi pełnego wykazu jej schematu. <strong>Po zamknięciu Tury nie da się odtworzyć, co faktycznie poszło do modelu</strong> — jakie argumenty wywołania, jaka treść Nakładki tożsamości i jakie ustawienia zostały zastosowane. Ograniczenie to jest istotne wszędzie tam, gdzie Licencjobiorca opiera się na rozliczalności pracy wykonanej z udziałem modelu; wyklucza ono również posługiwanie się prowenancją jako dowodem po fakcie (rozdz. 2.3, rozdz. 9.2).</li>
<li><strong>Większość kluczy katalogu ustawień zapisuje się, lecz nie jest odczytywana przy wywołaniu modelu.</strong> W wykonaniu Tury stosowane są wyłącznie: Nakładka tożsamości (tryb <code>ZASTAP</code> albo <code>DOLACZ</code>), tryb uprawnień, katalogi robocze, katalog roboczy sesji oraz mosty MCP. Pozostałe klucze katalogu ustawień — w tym <strong>cała rodzina <code>harness.*</code></strong> (<code>harness.program_claude</code>, <code>harness.plik_ustawien</code>, <code>harness.konfiguracja_mcp</code>) — zapisują się w Bazie <strong>bez wpływu na wywołanie modelu</strong>. Ujawnione wcześniej przypadki szczegółowe (rozdz. 3.4 pkt 1 — program kanału głównego; pkt 9 powyżej — model, model zapasowy i nakład rozumowania; pkt 10 — zakres wykonania; pkt 11 — izolacja) są przykładami tej reguły ogólnej, nie wyjątkami od niej. W konsekwencji uprawnienie do konfigurowania Oprogramowania przyznane w rozdz. 4.2 pkt 3 obejmuje <strong>zapis i rozstrzyganie wartości</strong>, nie zaś zapewnienie, że każda zapisana wartość wpłynie na zachowanie modelu.</li>
<li><strong>Punktu dostępu rodzaju <code>localDirectory</code> nie da się zapisać.</strong> Schemat trwałości wymaga dla tego rodzaju powiązania z wierszem urządzenia (kolumna <code>urzadzenie_id</code> wraz z więzem sprawdzającym <code>rodzaj &lt;&gt; 'localDirectory' OR urzadzenie_id IS NOT NULL</code>), a wiersz w tabeli <code>urzadzenie</code> nie powstaje w żaden sposób produkcyjny: nie istnieje repozytorium urządzeń w warstwie danych Rdzenia, migracje nie niosą zaczynu tej tabeli, a jedyne jej zasilenia występują w plikach sprawdzianów. W wersji v2.0 zarządzanie punktami dostępu z poziomu interfejsu obejmuje zatem wyłącznie punkty rodzaju <code>mcpBridge</code>.</li>
</ol>
<p>Wykaz powyższy odzwierciedla stan w bieżącym stanie repozytorium — audyt pierwotny przeprowadzono w stanie sprzed prac naprawczych, a punkty 6, 16 i 18 zaktualizowano po weryfikacji zmian wprowadzonych pracami naprawczymi, poprzedzającymi bieżący stan repozytorium (Załącznik C.1, C.3 pkt 4). Punkty nieoznaczone jako zaktualizowane pozostają prawdziwe na obu rewizjach. Wykaz nie stanowi zobowiązania do usunięcia wymienionych ograniczeń.</p>
<h3 id="104-ograniczenie-odpowiedzialności">10.4 Ograniczenie odpowiedzialności</h3>
<p>W najszerszym zakresie dopuszczalnym przez prawo właściwe Licencjodawca nie ponosi odpowiedzialności za:</p>
<ol type="1">
<li>utracone korzyści, utratę przychodu, utratę spodziewanych oszczędności, utratę renomy ani szkody pośrednie;</li>
<li>utratę, uszkodzenie albo ujawnienie Danych Operatora oraz Treści Operatora;</li>
<li>skutki decyzji podjętych na podstawie wyników pracy modelu;</li>
<li>działania modelu w zakresie zasobów udostępnionych mu przez Operatora, w tym za zmiany dokonane w plikach na maszynach objętych nadaniem trybu <code>write</code>;</li>
<li>przerwy w działaniu albo zmianę warunków świadczenia usług przez Dostawcę modelu, w tym za wyczerpanie limitów, zawieszenie konta albo zaprzestanie udostępniania programu zewnętrznego;</li>
<li>niezgodność korzystania z Oprogramowania z warunkami Dostawcy modelu albo z przepisami obowiązującymi Licencjobiorcę;</li>
<li>szkody wynikłe z korzystania z Oprogramowania w sposób sprzeczny z rozdz. 5.3.</li>
</ol>
<p>Ograniczenia nie dotyczą szkody wyrządzonej umyślnie ani odpowiedzialności, której wyłączyć nie można na podstawie przepisu bezwzględnie obowiązującego.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> w okresie nieodpłatnego udostępniania (rozdz. 4.5) odpowiedzialność Licencjodawcy jest ograniczona do szkody wyrządzonej <strong>z winy umyślnej</strong>, w najszerszym zakresie dopuszczalnym przez prawo właściwe. Z chwilą wprowadzenia odpłatności limit odpowiedzialności zostanie oznaczony jako równowartość wynagrodzenia zapłaconego przez Licencjobiorcę w okresie dwunastu miesięcy poprzedzających zdarzenie szkodzące, nową wersją dokumentu w trybie rozdz. 8.6. Ograniczenia nie uchybiają odpowiedzialności, której wyłączyć nie można na podstawie przepisu bezwzględnie obowiązującego.</p>
<h3 id="105-rozkład-ryzyka">10.5 Rozkład ryzyka</h3>
<p>Strony przyjmują, że opisany rozkład odpowiedzialności odpowiada charakterowi świadczenia: Oprogramowanie o statusie Deweloperskim udostępniane jest z ujawnieniem rzeczywistych ograniczeń (rozdz. 10.3), a Licencjobiorca podejmuje decyzję o korzystaniu z pełną wiedzą o tych ograniczeniach.</p>
<h3 id="106-siła-wyższa">10.6 Siła wyższa</h3>
<p>Żadna ze stron nie odpowiada za niewykonanie zobowiązań spowodowane okolicznościami pozostającymi poza jej rozsądną kontrolą, w szczególności działaniem siły wyższej, awarią infrastruktury teleinformatycznej niezależną od strony, decyzją organu władzy publicznej ani zaprzestaniem świadczenia usług przez Dostawcę modelu.</p>
<hr />
<h2 id="11-okres-obowiązywania-i-zakończenie-licencji">11. Okres obowiązywania i zakończenie licencji</h2>
<h3 id="111-wejście-w-życie">11.1 Wejście w życie</h3>
<p>Licencja wiąże z chwilą pierwszego z następujących zdarzeń: zainstalowania Oprogramowania, pierwszego uruchomienia Oprogramowania albo zaakceptowania warunków Licencji w toku instalacji. Przystąpienie do korzystania z Oprogramowania oznacza przyjęcie warunków Licencji w całości.</p>
<p>Podmiot, który nie akceptuje warunków Licencji, zobowiązany jest odstąpić od instalacji, a Egzemplarz — usunąć.</p>
<h3 id="112-okres-obowiązywania">11.2 Okres obowiązywania</h3>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> Licencja na <strong>Wydanie posiadane</strong> przez Licencjobiorcę jest udzielana <strong>bezterminowo</strong>, z zastrzeżeniem uprawnień do wypowiedzenia z rozdz. 11.3. Uprawnienie do Aktualizacji reguluje odrębnie rozdz. 8.4; wygaśnięcie albo zmiana tego uprawnienia nie pozbawia Licencjobiorcy prawa do korzystania z Wydania już posiadanego.</p>
<p>Niezależnie od wybranego wariantu Licencja wygasa najpóźniej z chwilą zaprzestania korzystania z Oprogramowania przez Licencjobiorcę połączonego z usunięciem wszystkich Egzemplarzy.</p>
<h3 id="113-wypowiedzenie-i-rozwiązanie">11.3 Wypowiedzenie i rozwiązanie</h3>
<ol type="1">
<li><strong>Wypowiedzenie przez Licencjobiorcę.</strong> Licencjobiorca może w każdym czasie zakończyć Licencję przez zaprzestanie korzystania z Oprogramowania i usunięcie wszystkich Egzemplarzy oraz kopii zapasowych Oprogramowania. Zakończenie nie rodzi po stronie Licencjodawcy obowiązku zwrotu świadczeń; przy nieodpłatnej Licencji fazy Deweloperskiej (rozdz. 4.5) świadczenia podlegające zwrotowi nie powstają.</li>
<li><strong>Wypowiedzenie przez Licencjodawcę z przyczyny naruszenia.</strong> Licencjodawca może wypowiedzieć Licencję ze skutkiem natychmiastowym, jeżeli Licencjobiorca istotnie narusza jej postanowienia, w szczególności ograniczenia z rozdz. 5.2, 5.3 albo 5.6, i nie zaprzestaje naruszenia w terminie wyznaczonym w wezwaniu skierowanym na piśmie albo drogą elektroniczną. <strong>Rozstrzygnięcie (2026-08-26):</strong> termin ten wynosi <strong>trzydzieści dni</strong> dla naruszeń usuwalnych (w szczególności przekroczenia liczby Operatorów albo braku ewidencji); przy naruszeniach nieusuwalnych — w szczególności rozpowszechnieniu Oprogramowania albo obejściu oznaczeń — wypowiedzenie może nastąpić <strong>bez wyznaczania terminu</strong>.</li>
<li><strong>Wypowiedzenie bez przyczyny.</strong> Wobec bezterminowego charakteru Licencji na Wydanie posiadane (rozdz. 11.2) Licencjodawcy <strong>nie przysługuje</strong> uprawnienie do wypowiedzenia Licencji bez wskazania przyczyny w odniesieniu do Wydania już posiadanego przez Licencjobiorcę; nie ogranicza to swobody Licencjodawcy co do udostępniania Wydań przyszłych.</li>
<li><strong>Rozwiązanie za porozumieniem.</strong> Strony mogą zakończyć Licencję w każdym czasie za obopólnym porozumieniem wyrażonym na piśmie.</li>
</ol>
<h3 id="114-skutki-zakończenia">11.4 Skutki zakończenia</h3>
<p>Z chwilą zakończenia Licencji:</p>
<ol type="1">
<li>wygasają wszystkie uprawnienia Licencjobiorcy opisane w rozdz. 4.2;</li>
<li>Licencjobiorca zobowiązany jest zaprzestać korzystania z Oprogramowania, odinstalować je ze wszystkich urządzeń i usunąć wszystkie Egzemplarze oraz kopie zapasowe Oprogramowania;</li>
<li>Licencjobiorca zachowuje prawo do Danych Operatora i Treści Operatora oraz prawo do zachowania kopii Katalogu danych — zakończenie Licencji <strong>nie zobowiązuje do usunięcia własnych danych</strong>. Zakres danych możliwych do zachowania wyznacza jednak zakres ich faktycznego utrwalenia w Bazie (rozdz. 9.2 pkt 1 oraz rozdz. 9.5 pkt 1); Licencjobiorca zamierzający zabezpieczyć stan pracy przed zakończeniem Licencji powinien wziąć to pod uwagę, tym bardziej że Oprogramowanie nie udostępnia funkcji eksportu całości Danych Operatora do formatu przenośnego (rozdz. 9.5 pkt 4);</li>
<li>Licencjobiorca zachowuje prawo korzystania z wytworów pracy powstałych przy użyciu Oprogramowania w okresie obowiązywania Licencji;</li>
<li>na żądanie Licencjodawcy Licencjobiorca potwierdza wykonanie obowiązków z pkt 2 oświadczeniem złożonym na piśmie albo drogą elektroniczną.</li>
</ol>
<h3 id="115-postanowienia-trwające">11.5 Postanowienia trwające</h3>
<p>Zakończenie Licencji nie wpływa na moc obowiązującą postanowień, które z natury mają obowiązywać dłużej, w szczególności rozdz. 6 (własność intelektualna), rozdz. 10 (odpowiedzialność), rozdz. 13 (poufność) oraz rozdz. 14 (postanowienia końcowe).</p>
<hr />
<h2 id="12-zgodność-licencyjna-i-kontrola-korzystania">12. Zgodność licencyjna i kontrola korzystania</h2>
<h3 id="121-obowiązek-zgodności">12.1 Obowiązek zgodności</h3>
<p>Licencjobiorca korzysta z Oprogramowania w granicach zakresu podmiotowego Licencji (rozdz. 4.3) i utrzymuje ewidencję pozwalającą wykazać zgodność korzystania z tym zakresem. Przedmiotem ewidencji — stosownie do rozstrzygnięcia rozdz. 4.3 — jest wykaz uprawnionych Operatorów oraz Instalacji, z których korzystają. Obowiązek ewidencyjny z rozdz. 4.3 pkt 1 i obowiązek niniejszy są jednym i tym samym obowiązkiem opisanym w dwóch miejscach dokumentu ze względu na układ rozdziałów.</p>
<h3 id="122-weryfikacja">12.2 Weryfikacja</h3>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> Licencjodawca może żądać od Licencjobiorcy <strong>oświadczenia o liczbie Instalacji i uprawnionych Operatorów</strong>, składanego w postaci dokumentowej, nie częściej niż <strong>raz w roku kalendarzowym</strong>. Licencjodawcy nie przysługuje prawo audytu w siedzibie Licencjobiorcy ani dostęp do treści zapisanych w Bazie.</p>
<p><strong>Oprogramowanie nie zawiera mechanizmu technicznej kontroli licencji</strong>: nie weryfikuje klucza licencyjnego, nie łączy się z serwerem aktywacji ani nie raportuje faktu uruchomienia. Weryfikacja zgodności może zatem odbywać się wyłącznie środkami organizacyjnymi.</p>
<h3 id="123-skutki-stwierdzenia-niezgodności">12.3 Skutki stwierdzenia niezgodności</h3>
<p>W razie stwierdzenia korzystania wykraczającego poza zakres Licencji strony w pierwszej kolejności ustalają zakres niezgodności i sposób jej usunięcia. Nieusunięcie niezgodności w uzgodnionym terminie stanowi istotne naruszenie Licencji w rozumieniu rozdz. 11.3 pkt 2.</p>
<hr />
<h2 id="13-poufność">13. Poufność</h2>
<h3 id="131-informacje-poufne">13.1 Informacje poufne</h3>
<p>Za informacje poufne uważa się nieujawnione publicznie informacje dotyczące Oprogramowania, przekazane Licencjobiorcy przez Licencjodawcę albo powzięte w związku z korzystaniem z Oprogramowania, w szczególności: dokumentację techniczną nieprzeznaczoną do publikacji, opisy architektury wewnętrznej, struktury Kontraktu i modelu danych, informacje o niedziałających funkcjach i podatnościach, a także informacje o planach rozwoju produktu.</p>
<p>Za informacje poufne <strong>nie uważa się</strong> informacji, które: są publicznie dostępne bez naruszenia zobowiązania do poufności; były znane Licencjobiorcy przed ich przekazaniem; zostały uzyskane zgodnie z prawem od osoby trzeciej nieobjętej zobowiązaniem do poufności; albo zostały opracowane samodzielnie bez wykorzystania informacji poufnych.</p>
<h3 id="132-obowiązki">13.2 Obowiązki</h3>
<p>Licencjobiorca zobowiązuje się zachować informacje poufne w tajemnicy, nie ujawniać ich osobom trzecim bez zgody Licencjodawcy udzielonej na piśmie oraz udostępniać je własnym Operatorom i współpracownikom wyłącznie w zakresie niezbędnym do korzystania z Oprogramowania, po zobowiązaniu ich do poufności.</p>
<p>Obowiązek zachowania poufności nie stoi na przeszkodzie ujawnieniu informacji na żądanie uprawnionego organu, w zakresie i w trybie wynikającym z przepisów; o takim żądaniu Licencjobiorca niezwłocznie zawiadamia Licencjodawcę, o ile nie jest to zabronione.</p>
<h3 id="133-okres-obowiązywania">13.3 Okres obowiązywania</h3>
<p>Zobowiązanie do poufności obowiązuje w okresie trwania Licencji oraz przez <strong>trzy lata</strong> po jej zakończeniu (<strong>Rozstrzygnięcie 2026-08-26</strong>), a w odniesieniu do informacji stanowiących tajemnicę przedsiębiorstwa: przez cały okres, w którym zachowują one taki charakter. Ta druga część zobowiązania nie jest ograniczona terminem.</p>
<h3 id="134-ujawnienie-podatności">13.4 Ujawnienie podatności</h3>
<p>Licencjobiorca, który poweźmie informację o podatności bezpieczeństwa Oprogramowania, zawiadamia o niej wyłącznie Licencjodawcę i powstrzymuje się od jej publikacji przez okres uzgodniony z Licencjodawcą, nie dłuższy jednak niż konieczny do wydania poprawki.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> zgłoszenia podatności kieruje się na adres poczty elektronicznej <strong>security@danaco-group.pl</strong>. Strony stosują zasadę skoordynowanego ujawniania: Licencjobiorca powstrzymuje się od publikacji przez okres uzgodniony z Licencjodawcą, nie dłuższy niż <strong>dziewięćdziesiąt dni</strong> od zgłoszenia, chyba że poprawka zostanie wydana wcześniej albo strony uzgodnią inaczej.</p>
<hr />
<h2 id="14-postanowienia-końcowe">14. Postanowienia końcowe</h2>
<h3 id="141-prawo-właściwe">14.1 Prawo właściwe</h3>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> prawem właściwym dla Licencji jest <strong>prawo polskie</strong>, z wyłączeniem norm kolizyjnych. Jeżeli Licencjobiorcą jest konsument, wybór prawa nie pozbawia go ochrony przyznanej przepisami bezwzględnie obowiązującymi państwa jego zwykłego pobytu.</p>
<h3 id="142-rozstrzyganie-sporów">14.2 Rozstrzyganie sporów</h3>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> spory wynikłe z Licencji rozstrzyga <strong>sąd powszechny właściwy miejscowo dla siedziby Licencjodawcy</strong>. Zastrzeżenie to nie wiąże Licencjobiorcy będącego konsumentem — w sporach z konsumentem właściwość sądu ustala się według przepisów ogólnych.</p>
<p>Niezależnie od wybranego wariantu strony podejmą w dobrej wierze próbę polubownego rozwiązania sporu przed skierowaniem sprawy na drogę postępowania sądowego.</p>
<h3 id="143-forma-zmian">14.3 Forma zmian</h3>
<p>Zmiany Licencji w stosunku dwustronnym wymagają formy pisemnej albo dokumentowej pod rygorem nieważności, z zastrzeżeniem rozdz. 8.6, który dotyczy wydawania nowych wersji dokumentu Licencji wraz z nowymi Wydaniami.</p>
<h3 id="144-klauzula-salwatoryjna">14.4 Klauzula salwatoryjna</h3>
<p>Nieważność albo bezskuteczność któregokolwiek z postanowień Licencji nie wpływa na ważność pozostałych postanowień. W miejsce postanowienia nieważnego albo bezskutecznego strony stosują postanowienie skuteczne, najbliższe celowi gospodarczemu postanowienia zastępowanego. Dotyczy to w szczególności postanowień rozdz. 5.2, 10.2 i 10.4, których zakres podlega ograniczeniu do granic dopuszczalnych przez prawo właściwe.</p>
<h3 id="145-całość-porozumienia">14.5 Całość porozumienia</h3>
<p>Licencja wraz z załącznikami stanowi całość porozumienia stron w zakresie korzystania z Oprogramowania i zastępuje wcześniejsze ustalenia, oświadczenia i zapewnienia dotyczące tego przedmiotu, z zastrzeżeniem rozdz. 1.4 pkt 2. Materiały prezentacyjne, koncepcyjne i poglądowe dotyczące produktu nie stanowią części Licencji i nie tworzą zobowiązania co do zakresu funkcjonalnego.</p>
<h3 id="146-język-dokumentu">14.6 Język dokumentu</h3>
<p>Językiem Licencji jest język polski. W razie sporządzenia tłumaczenia rozstrzygające znaczenie ma wersja polska.</p>
<h3 id="147-dane-kontaktowe-i-doręczenia">14.7 Dane kontaktowe i doręczenia</h3>
<p>Oświadczenia i zawiadomienia związane z Licencją składa się na adres siedziby Licencjodawcy albo na wskazany przez niego adres poczty elektronicznej.</p>
<p><strong>Rozstrzygnięcie (2026-08-26):</strong> dane kontaktowe Licencjodawcy:</p>
<div class="dn-zestawienie-pole"><table>
<tbody>
<tr>
<td><strong>Licencjodawca</strong></td>
<td>Danaco Holding Group Sp. z o.o.</td>
</tr>
<tr>
<td><strong>Adres siedziby</strong></td>
<td><em>[do uzupełnienia danymi rejestrowymi]</em></td>
</tr>
<tr>
<td><strong>KRS / NIP</strong></td>
<td><em>[do uzupełnienia danymi rejestrowymi]</em></td>
</tr>
<tr>
<td><strong>Korespondencja i wsparcie</strong></td>
<td>support@danaco-group.pl</td>
</tr>
<tr>
<td><strong>Zgłoszenia podatności</strong></td>
<td>security@danaco-group.pl (rozdz. 13.4)</td>
</tr>
</tbody>
</table></div>
<p>Uzupełnienie pól rejestrowych następuje przed dystrybucją zewnętrzną i nie stanowi zmiany Licencji w rozumieniu rozdz. 14.3.</p>
<h3 id="148-nagłówki-i-załączniki">14.8 Nagłówki i załączniki</h3>
<p>Załączniki A, B i C stanowią integralną część Licencji. W razie rozbieżności między treścią rozdziału a treścią załącznika rozstrzyga treść rozdziału, z wyjątkiem Załącznika A w zakresie zestawienia Komponentów, gdzie rozstrzyga treść załącznika jako zestawienia szczegółowego.</p>
<hr />
<h2 id="załącznik-a--zestawienie-komponentów-osób-trzecich">Załącznik A — zestawienie komponentów osób trzecich</h2>
<p>Zestawienie sporządzono na podstawie bieżącego stanu repozytorium (audyt pierwotny w stanie sprzed prac naprawczych, zaktualizowany po zamknięciu prac naprawczych — Załącznik C.3 pkt 4). Zestawienie <strong>nie zastępuje</strong> pełnych tekstów licencji, które Licencjodawca dołącza do Wydania zgodnie z rozdz. 7.9.</p>
<h3 id="a1-warstwa-rdzenia--moduł-go-danacoconsole-wymagane-go-126">A.1. Warstwa Rdzenia — moduł Go <code>danacoconsole</code> (wymagane <code>go 1.26</code>)</h3>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Lp.</th>
<th>Komponent</th>
<th>Wersja</th>
<th>Licencja</th>
<th>Zależność</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td><code>github.com/coder/websocket</code></td>
<td>v1.8.15</td>
<td>ISC</td>
<td>bezpośrednia</td>
</tr>
<tr>
<td>2</td>
<td><code>golang.org/x/sys</code></td>
<td>v0.47.0</td>
<td>BSD 3-Clause</td>
<td>bezpośrednia</td>
</tr>
<tr>
<td>3</td>
<td><code>modernc.org/sqlite</code></td>
<td>v1.56.0</td>
<td>BSD 3-Clause</td>
<td>bezpośrednia</td>
</tr>
<tr>
<td>4</td>
<td><code>github.com/dustin/go-humanize</code></td>
<td>v1.0.1</td>
<td>MIT</td>
<td>pośrednia</td>
</tr>
<tr>
<td>5</td>
<td><code>github.com/google/uuid</code></td>
<td>v1.6.0</td>
<td>BSD 3-Clause</td>
<td>pośrednia</td>
</tr>
<tr>
<td>6</td>
<td><code>github.com/mattn/go-isatty</code></td>
<td>v0.0.24</td>
<td>MIT</td>
<td>pośrednia</td>
</tr>
<tr>
<td>7</td>
<td><code>github.com/ncruces/go-strftime</code></td>
<td>v1.0.0</td>
<td>MIT</td>
<td>pośrednia</td>
</tr>
<tr>
<td>8</td>
<td><code>github.com/remyoudompheng/bigfft</code></td>
<td>v0.0.0-20230129092748-24d4a6f8daec</td>
<td>BSD 3-Clause</td>
<td>pośrednia</td>
</tr>
<tr>
<td>9</td>
<td><code>modernc.org/libc</code></td>
<td>v1.74.4</td>
<td>BSD 3-Clause</td>
<td>pośrednia</td>
</tr>
<tr>
<td>10</td>
<td><code>modernc.org/mathutil</code></td>
<td>v1.7.1</td>
<td>BSD 3-Clause</td>
<td>pośrednia</td>
</tr>
<tr>
<td>11</td>
<td><code>modernc.org/memory</code></td>
<td>v1.11.0</td>
<td>BSD 3-Clause</td>
<td>pośrednia</td>
</tr>
</tbody>
</table></div>
<h3 id="a2-warstwa-klienta--pakiet-npm-danaco-console-client">A.2. Warstwa Klienta — pakiet npm <code>danaco-console-client</code></h3>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Lp.</th>
<th>Komponent</th>
<th>Wersja</th>
<th>Licencja</th>
<th>Charakter</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td><code>@tauri-apps/api</code></td>
<td>2.5.0</td>
<td>Apache-2.0 OR MIT</td>
<td>produkcyjna</td>
</tr>
<tr>
<td>2</td>
<td><code>typescript</code></td>
<td>5.8.3</td>
<td>Apache-2.0</td>
<td>narzędzie budowania</td>
</tr>
<tr>
<td>3</td>
<td><code>vite</code></td>
<td>6.3.5</td>
<td>MIT</td>
<td>narzędzie budowania</td>
</tr>
<tr>
<td>4</td>
<td><code>vitest</code></td>
<td>3.2.7</td>
<td>MIT</td>
<td>narzędzie testowe</td>
</tr>
<tr>
<td>5</td>
<td><code>jsdom</code></td>
<td>^29.1.1</td>
<td>MIT</td>
<td>narzędzie testowe (środowisko DOM, prace nad sprawdzianami widoków)</td>
</tr>
</tbody>
</table></div>
<p>Zamknięcie zależności: 141 pozycji; rozkład licencji — MIT 126, Apache-2.0 3, ISC 3, MIT-0 2, BSD-2-Clause 2, BSD-3-Clause 2, „Apache-2.0 OR MIT” 1, BlueOak-1.0.0 1, CC0-1.0 1. Brak licencji wzajemnych.</p>
<h3 id="a3-warstwa-powłoki--pakiet-cargo-danaco-console-powloka">A.3. Warstwa Powłoki — pakiet Cargo <code>danaco-console-powloka</code></h3>
<p>Zależności bezpośrednie (bieżący stan repozytorium; kolumna „Sekcja” wskazuje, czy zależność pochodzi z <code>[dependencies]</code> — wchodzi do pliku wynikowego — czy z <code>[build-dependencies]</code> — narzędzie etapu budowania):</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Lp.</th>
<th>Komponent</th>
<th>Sekcja</th>
<th>Wersja rozwiązana</th>
<th>Licencja</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td><code>tauri</code></td>
<td><code>[dependencies]</code></td>
<td>2.11.5</td>
<td>Apache-2.0 OR MIT</td>
</tr>
<tr>
<td>2</td>
<td><code>tauri-plugin-dialog</code></td>
<td><code>[dependencies]</code></td>
<td>2.7.2</td>
<td>Apache-2.0 OR MIT</td>
</tr>
<tr>
<td>3</td>
<td><code>serde</code></td>
<td><code>[dependencies]</code></td>
<td>1.0.229</td>
<td>MIT OR Apache-2.0</td>
</tr>
<tr>
<td>4</td>
<td><code>tauri-build</code></td>
<td><code>[build-dependencies]</code></td>
<td>2.6.3</td>
<td>Apache-2.0 OR MIT</td>
</tr>
</tbody>
</table></div>
<p><code>serde_json</code> (rozwiązana 1.0.151, MIT OR Apache-2.0) figurowała jako zależność bezpośrednia w stanie sprzed prac naprawczych; usunięta z <code>Cargo.toml</code> jako nieużywana w pracach nad skryptami wydania, pozostaje w zamknięciu <code>Cargo.lock</code> wyłącznie jako zależność pośrednia wciągana przez <code>tauri</code> — patrz rozdz. 7.4.</p>
<p>Zamknięcie zależności: 447 pakietów, z czego 274 zweryfikowane co do zadeklarowanej licencji. Komponenty na licencji wzajemnej MPL-2.0: <code>cssparser</code> 0.36.0, <code>cssparser-macros</code> 0.6.1, <code>dtoa-short</code> 0.3.5, <code>option-ext</code> 0.2.0, <code>selectors</code> 0.36.1.</p>
<p><strong>Zestawienie warstwy Powłoki wymaga uzupełnienia przed publikacją</strong> — patrz rozdz. 7.4 (zastrzeżenie o kompletności).</p>
<h3 id="a4-kroje-pisma">A.4. Kroje pisma</h3>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Lp.</th>
<th>Krój</th>
<th>Warianty</th>
<th>Licencja</th>
<th>Plik licencji</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>Space Grotesk</td>
<td>500, 600, 700 (latin, latin-ext)</td>
<td>SIL OFL 1.1</td>
<td><code>LICENCJA-space-grotesk.txt</code></td>
</tr>
<tr>
<td>2</td>
<td>IBM Plex Sans</td>
<td>400, 500, 600, 700 (latin, latin-ext)</td>
<td>SIL OFL 1.1</td>
<td><code>LICENCJA-ibm-plex.txt</code></td>
</tr>
<tr>
<td>3</td>
<td>IBM Plex Mono</td>
<td>400, 500, 600 (latin, latin-ext)</td>
<td>SIL OFL 1.1</td>
<td><code>LICENCJA-ibm-plex.txt</code></td>
</tr>
</tbody>
</table></div>
<h3 id="a5-zestaw-ikon-i-zasoby-marki">A.5. Zestaw ikon i zasoby marki</h3>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Lp.</th>
<th>Zasób</th>
<th>Liczba</th>
<th>Licencja / status</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>Ikony wywiedzione z biblioteki Lucide</td>
<td>78</td>
<td>ISC</td>
</tr>
<tr>
<td>2</td>
<td>Emblematy Środowisk (własne)</td>
<td>4</td>
<td>utwór Licencjodawcy</td>
</tr>
<tr>
<td>3</td>
<td><strong>Zestaw ikon łącznie — wykaz manifestu <code>ikony/manifest.json</code></strong></td>
<td><strong>82</strong></td>
<td>78 ISC + 4 własne</td>
</tr>
<tr>
<td>4</td>
<td>Znak marki <code>logo-danaco.svg</code> (poza zestawem ikon, w katalogu plików zestawu)</td>
<td>1 plik</td>
<td>utwór Licencjodawcy</td>
</tr>
<tr>
<td>5</td>
<td>Logotypy, sygnety i ich warianty (katalog zasobów marki)</td>
<td>12 plików</td>
<td>utwór Licencjodawcy</td>
</tr>
</tbody>
</table></div>
<p>Katalog plików SVG zestawu zawiera 83 pliki: 82 pozycje zestawu ikon oraz znak marki z poz. 4. <code>README.md</code> i <code>INSTRUKCJA-UZYTKOWANIA.md</code> posługują się zgodnie liczbą 82 pozycji zestawu ikon i objaśniają ten sam podział rzeczowy wobec 83 plików katalogu — rozdz. 7.6 opisuje weryfikację tej zgodności; kwestia nie figuruje już w Załączniku B jako nierozstrzygnięta.</p>
<h3 id="a6-zależności-zewnętrzne-nieobjęte-dystrybucją">A.6. Zależności zewnętrzne nieobjęte dystrybucją</h3>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Lp.</th>
<th>Zależność</th>
<th>Status</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>Program <code>claude</code> (Claude Code CLI)</td>
<td>nie jest dystrybuowany; zapewnia Operator; warunki Dostawcy modelu</td>
</tr>
<tr>
<td>2</td>
<td>Usługa modelu wywoływana kanałem rodzaju <code>api</code></td>
<td>nie jest dystrybuowana; warunki dostawcy usługi</td>
</tr>
<tr>
<td>3</td>
<td>Program <code>ssh</code> systemu operacyjnego</td>
<td>wykorzystywany przez most MCP; element systemu operacyjnego</td>
</tr>
<tr>
<td>4</td>
<td>Komponent webview systemu Windows</td>
<td>element systemu operacyjnego</td>
</tr>
</tbody>
</table></div>
<hr />
<h2 id="załącznik-b--wykaz-rozstrzygnięć-operatora">Załącznik B — wykaz rozstrzygnięć Operatora</h2>
<p>Kwestie wymienione poniżej były we wcześniejszej redakcji dokumentu pozostawione do decyzji Operatora. Zostały rozstrzygnięte decyzją Operatora z dnia <strong>2026-08-26</strong> i wpisane do treści rozdziałów wskazanych w tabeli. W razie rozbieżności między tabelą a treścią rozdziału rozstrzyga treść rozdziału (rozdz. 14.8).</p>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Lp.</th>
<th>Rozdział</th>
<th>Kwestia</th>
<th>Rozstrzygnięcie</th>
</tr>
</thead>
<tbody>
<tr>
<td>1</td>
<td>1.2</td>
<td>Definicja Licencjobiorcy</td>
<td>podmiot (w tym osoba fizyczna) z prawem wskazania Operatorów</td>
</tr>
<tr>
<td>2</td>
<td>3.2</td>
<td>Model dystrybucji końcowej</td>
<td>model mieszany: instalator lokalny + Rdzeń serwerowy po odrębnym uzgodnieniu; jednostka: Operator</td>
</tr>
<tr>
<td>3</td>
<td>4.3</td>
<td>Liczba stanowisk i jednostka liczenia</td>
<td>licencja imienna per Operator; wiele urządzeń, zakaz korzystania równoczesnego; domyślnie 1 Operator</td>
</tr>
<tr>
<td>4</td>
<td>4.5</td>
<td>Odpłatność</td>
<td>nieodpłatna dla statusu Deweloperskiego; odpłatność od pierwszego Wydania produkcyjnego, w trybie rozdz. 8.6</td>
</tr>
<tr>
<td>5</td>
<td>5.2 pkt 8</td>
<td>Zakaz budowy produktu konkurencyjnego</td>
<td>zakaz wąski — wyłącznie odwzorowanie elementów chronionych z rozdz. 6.1</td>
</tr>
<tr>
<td>6</td>
<td>5.2 pkt 9</td>
<td>Badania bezpieczeństwa</td>
<td>dozwolone we własnym środowisku; zakaz badania cudzych instalacji; zgłoszenia trybem rozdz. 13.4</td>
</tr>
<tr>
<td>7</td>
<td>5.8</td>
<td>Kary umowne</td>
<td>brak kar umownych; zasady ogólne</td>
</tr>
<tr>
<td>8</td>
<td>6.4</td>
<td>Status ochronny znaków</td>
<td>ochrona oznaczeń nierejestrowanych; rejestracja odzwierciedlana w trybie rozdz. 8.6</td>
</tr>
<tr>
<td>9</td>
<td>6.6</td>
<td>Informacje zwrotne</td>
<td>nieodpłatna, niewyłączna licencja dla Licencjodawcy na korzystanie w celu rozwoju produktu</td>
</tr>
<tr>
<td>10</td>
<td>7.4</td>
<td>Inwentarz licencji Komponentów</td>
<td>plik <code>NOTICE</code> generowany automatycznie przy każdym Wydaniu; rozstrzyga <code>NOTICE</code></td>
</tr>
<tr>
<td>11</td>
<td>7.8</td>
<td>Rozszerzenia</td>
<td>wyłącznie od Licencjodawcy; rozszerzenia osób trzecich niedozwolone do odrębnej regulacji</td>
</tr>
<tr>
<td>12</td>
<td>8.4</td>
<td>Aktualizacje i wsparcie</td>
<td>Aktualizacje w ramach Licencji nieodpłatnie, bez zobowiązania do wydawania; wsparcie e-mail bez czasów reakcji</td>
</tr>
<tr>
<td>13</td>
<td>9.6</td>
<td>Role RODO przy Rdzeniu serwerowym</td>
<td>rozdz. 9.6 dotyczy Rdzenia lokalnego; topologia serwerowa wymaga uprzedniej umowy powierzenia (art. 28 RODO)</td>
</tr>
<tr>
<td>14</td>
<td>10.4</td>
<td>Górna granica odpowiedzialności</td>
<td>w fazie nieodpłatnej: wina umyślna; po wprowadzeniu odpłatności: limit 12-miesięcznego wynagrodzenia</td>
</tr>
<tr>
<td>15</td>
<td>11.2</td>
<td>Okres obowiązywania</td>
<td>licencja bezterminowa na Wydanie posiadane; Aktualizacje odrębnie (rozdz. 8.4)</td>
</tr>
<tr>
<td>16</td>
<td>11.3 pkt 2</td>
<td>Termin usunięcia naruszenia</td>
<td>30 dni dla naruszeń usuwalnych; bez terminu dla nieusuwalnych</td>
</tr>
<tr>
<td>17</td>
<td>12.2</td>
<td>Kontrola zgodności</td>
<td>oświadczenie o liczbie Instalacji i Operatorów, nie częściej niż raz w roku; bez audytu i dostępu do Bazy</td>
</tr>
<tr>
<td>18</td>
<td>13.3</td>
<td>Okres poufności następczej</td>
<td>trzy lata; tajemnica przedsiębiorstwa bez terminu</td>
</tr>
<tr>
<td>19</td>
<td>13.4</td>
<td>Kanał zgłaszania podatności</td>
<td>security@danaco-group.pl; skoordynowane ujawnianie do 90 dni</td>
</tr>
<tr>
<td>20</td>
<td>14.1</td>
<td>Prawo właściwe</td>
<td>prawo polskie, z zachowaniem bezwzględnej ochrony konsumenta</td>
</tr>
<tr>
<td>21</td>
<td>14.2</td>
<td>Właściwość sądu</td>
<td>sąd siedziby Licencjodawcy; wobec konsumentów — przepisy ogólne</td>
</tr>
<tr>
<td>22</td>
<td>14.7</td>
<td>Dane identyfikacyjne Licencjodawcy</td>
<td>wpisane adresy e-mail; pola rejestrowe (adres, KRS, NIP) do uzupełnienia przed dystrybucją zewnętrzną</td>
</tr>
<tr>
<td>23</td>
<td>1.5</td>
<td>Przegląd prawny</td>
<td>wersja ustalona dla fazy Deweloperskiej; przegląd prawny zalecany przed dystrybucją komercyjną</td>
</tr>
</tbody>
</table></div>
<p>Poprzednia redakcja niniejszego załącznika zawierała dodatkową pozycję 11 (rozdz. 7.6 — jednolita liczba pozycji zestawu ikon w dokumentach produktu). Pozycja została usunięta po weryfikacji bezpośredniej treści <code>README.md</code> i <code>INSTRUKCJA-UZYTKOWANIA.md</code>, która wykazała, że oba dokumenty posługują się zgodnie liczbą 82 pozycji zestawu ikon — rozbieżność, dla której pozycja istniała, nie występuje już w treści dokumentów produktu (rozdz. 7.6). Numeracja pozostałych pozycji została odpowiednio przesunięta; poniższa uwaga posługuje się numeracją bieżącą.</p>
<p><strong>Uwaga o czterech pozycjach: poz. 5, 6, 16 i 18.</strong> Pozycje te we wcześniejszej redakcji dokumentu zawierały wartości przyjęte bez oznaczenia (odpowiednio: zakaz konkurencji, zakaz badań bezpieczeństwa, termin czternastu dni, pięcioletni okres poufności), następnie oznaczone jako wymagające decyzji. Rozstrzygnięcia z dnia 2026-08-26 zastępują te wartości: zakaz wąski, badania własnego środowiska dozwolone, termin trzydziestu dni dla naruszeń usuwalnych, trzyletni okres poufności następczej.</p>
<hr />
<h2 id="załącznik-c--metryka-weryfikacji-faktograficznej">Załącznik C — metryka weryfikacji faktograficznej</h2>
<h3 id="c1-podstawa-faktograficzna">C.1. Podstawa faktograficzna</h3>
<p>Twierdzenia o funkcjach, ograniczeniach, zależnościach i sposobie działania Oprogramowania zawarte w niniejszym dokumencie oparto na:</p>
<ol type="1">
<li>skonsolidowanym wyniku audytu repozytorium przeprowadzonego pierwotnie na stanu sprzed prac naprawczych, uzupełnionym ponowną weryfikacją punktową w bieżącym stanie repozytorium tej samej gałęzi po fali naprawczych — stanu, z którego zbudowano i uruchomiono doręczany Egzemplarz (C.3 pkt 4);</li>
<li>dokumencie nadrzędnym projektu (README) oraz zapisie decyzji architektonicznych;</li>
<li>dokumentacji technicznej towarzyszącej kodowi (kontrakt komunikacyjny — Kontrakty komunikacji; model danych — Model danych; mechanizm izolacji — Izolacja i zależności; katalog modułów i środowisk — <code>moduly/</code>, <code>srodowiska/</code>; katalog okien operacyjnych w bazie i wykaz luk — rozdz. 18 i 20 <code>README.md</code>);</li>
<li>bezpośrednim odczycie plików repozytorium: deklaracji zależności, plików zamknięcia zależności, konfiguracji pakietu Powłoki, wzorca konfiguracji środowiska, migracji schematu, manifestu warstwy wizualnej oraz plików licencyjnych krojów pisma.</li>
</ol>
<h3 id="c2-elementy-zweryfikowane-bezpośrednio-w-kodzie">C.2. Elementy zweryfikowane bezpośrednio w kodzie</h3>
<div class="dn-zestawienie-pole"><table>
<thead>
<tr>
<th>Twierdzenie dokumentu</th>
<th>Źródło weryfikacji</th>
</tr>
</thead>
<tbody>
<tr>
<td>Domyślny port nasłuchu 17870</td>
<td>wzorzec konfiguracji środowiska, konfiguracja Rdzenia</td>
</tr>
<tr>
<td>Wykaz zmiennych środowiska (pięć pozycji)</td>
<td>wzorzec konfiguracji środowiska</td>
</tr>
<tr>
<td>Katalog danych <code>%LOCALAPPDATA%\\DanacoConsole</code></td>
<td>moduł katalogu danych Rdzenia</td>
</tr>
<tr>
<td>Nazwa pliku Bazy <code>danaco-console.db</code></td>
<td>moduł pliku bazy Rdzenia</td>
</tr>
<tr>
<td>Nazwa produktu, identyfikator, wydawca, format instalatora</td>
<td>konfiguracja pakietu Powłoki</td>
</tr>
<tr>
<td>Brak komponentu samoczynnej aktualizacji</td>
<td>konfiguracja i manifest Powłoki</td>
</tr>
<tr>
<td>Zależności i wersje trzech warstw</td>
<td>pliki deklaracji i zamknięcia zależności</td>
</tr>
<tr>
<td>Rodzaje licencji Komponentów warstwy Go</td>
<td>pliki licencyjne w lokalnej pamięci podręcznej modułów</td>
</tr>
<tr>
<td>Licencje krojów pisma (OFL 1.1)</td>
<td>pliki licencyjne dołączone do zasobów</td>
</tr>
<tr>
<td>Liczba i źródło pozycji zestawu ikon</td>
<td>manifest warstwy wizualnej</td>
</tr>
<tr>
<td>Przechowywanie odwołań zamiast treści poświadczeń</td>
<td>schemat tabel kont i punktów dostępu</td>
</tr>
<tr>
<td>Brak połączeń wychodzących poza kanał modelu i most</td>
<td>przegląd użyć warstwy sieciowej w Rdzeniu</td>
</tr>
<tr>
<td>Pola koperty Kontraktu: wymagane <code>type</code>, <code>id</code>, <code>timestamp</code>; opcjonalne <code>sessionId</code>, <code>payload</code>; pola odpowiedzi <code>status</code>, <code>error</code>; pola strumienia <code>seq</code>, <code>done</code></td>
<td>plik Kontraktu, sekcja koperty</td>
</tr>
<tr>
<td>Wersja rozwiązana <code>serde</code> 1.0.229 (zależność bezpośrednia) i <code>serde_json</code> 1.0.151 (zależność pośrednia po usunięciu z <code>Cargo.toml</code> w pracach nad skryptami wydania)</td>
<td>plik zamknięcia zależności warstwy Rust w bieżącym stanie repozytorium</td>
</tr>
<tr>
<td><code>tauri-build</code> jest zależnością sekcji <code>[build-dependencies]</code>, nie <code>[dependencies]</code></td>
<td><code>budowa/desktop/src-tauri/Cargo.toml</code></td>
</tr>
<tr>
<td>Zestaw ikon: 82 pozycje wykazu manifestu (78 wywiedzionych z Lucide, 4 własne) przy 83 plikach katalogu zestawu, z których jeden jest znakiem marki; zgodność <code>README.md</code> i <code>INSTRUKCJA-UZYTKOWANIA.md</code> z liczbą 82</td>
<td>manifest warstwy wizualnej, wiązania zestawu w kodzie Klienta, wykaz plików katalogu, odczyt bezpośredni obu dokumentów</td>
</tr>
<tr>
<td>Zdarzenie <code>progress.changed</code> ma konsumenta; dane przykładowe pulpitu dowodzenia usunięte, pulpit odczytuje <code>session.list</code>/<code>window.list</code>/<code>channel.list</code></td>
<td><code>budowa/client/src/mission-control/zrodlo-pulpitu.ts</code></td>
</tr>
<tr>
<td><code>session.list</code>, <code>session.bind</code>, <code>session.focus</code> wywoływane przez stronę główną (strefa „Sesje w tle”)</td>
<td><code>budowa/client/src/strona-glowna/wpiecie-sesji.ts</code>, <code>zrodlo-sesji.ts</code></td>
</tr>
<tr>
<td>Kolumny CHECK <code>rodzaj_kanalu IN ('cli','api','sdk','lokalny')</code>; kanał <code>echo</code> zakłada się jako rodzaj <code>lokalny</code> z parametrem <code>{"adapter":"echo"}</code></td>
<td><code>budowa/server/internal/store/migracja_001_fundament.sql</code>; <code>budowa/server/internal/models/kanal.go</code></td>
</tr>
<tr>
<td>Istnienie <code>budowa/scripts/wydanie.sh</code> i <code>budowa/scripts/pakowanie.sh</code>; katalog <code>C:\\DanacoConsole_App</code> zbudowany tą ścieżką i uruchomiony</td>
<td>zawartość skryptów; zawartość katalogu produktu</td>
</tr>
<tr>
<td>Brak utrwalenia prowenancji — brak tabeli i kolumny w schemacie; wykaz pól tabeli wiadomości</td>
<td>migracje schematu trwałości; moduł prowenancji Rdzenia</td>
</tr>
<tr>
<td>Niemożność zapisu Punktu dostępu rodzaju <code>localDirectory</code> — więz wymagający wiersza urządzenia i brak produkcyjnego źródła tego wiersza</td>
<td>migracja punktów dostępu; przegląd warstwy danych Rdzenia i migracji pod kątem tabeli urządzeń</td>
</tr>
<tr>
<td>Zakres ustawień stosowanych przy wywołaniu modelu; brak odczytu rodziny <code>harness.*</code></td>
<td>katalog ustawień w migracji; przegląd składania wywołania w Rdzeniu</td>
</tr>
<tr>
<td>Dostępność cyklu życia Konta z poziomu interfejsu przy braku interfejsu zakładania kanału</td>
<td>warstwa kont i formularz konta w Kliencie; wywołania <code>channel.list</code> w widokach</td>
</tr>
<tr>
<td>Nazwy plików dokumentacji produktu</td>
<td>zawartość katalogu produktu; wykaz dokumentów w <code>README.md</code></td>
</tr>
</tbody>
</table></div>
<h3 id="c3-ograniczenia-weryfikacji">C.3. Ograniczenia weryfikacji</h3>
<ol type="1">
<li>Zestawienie licencji zamknięcia zależności warstwy Rust jest niepełne — zweryfikowano 274 z 447 pakietów (rozdz. 7.4).</li>
<li>Weryfikacja miała charakter statyczny: nie uruchamiano procesu budowania ani samego Oprogramowania w celu potwierdzenia zachowania w czasie działania. Twierdzenia o działaniu funkcji pochodzą z wyniku audytu repozytorium.</li>
<li>Dokument opisuje stan na dwóch oznaczonych rewizjach tej samej gałęzi: rewizji audytu pierwotnego stanu sprzed prac naprawczych oraz bieżącego stanu repozytorium, na której punktowo zweryfikowano i zaktualizowano twierdzenia unieważnione falą naprawczą (niniejszy punkt 4). Twierdzenie nieoznaczone żadną z tych rewizji z osobna obowiązuje na obu. Każda kolejna zmiana zależności, warstwy wizualnej albo zakresu funkcjonalnego wymaga ponownej weryfikacji rozdz. 7, 10.3 oraz Załącznika A na rewizji faktycznie wydawanej.</li>
<li><strong>Rewizja stanu sprzed prac naprawczych a bieżący stan gałęzi roboczej.</strong> Między rewizją audytu pierwotnego stanu sprzed prac naprawczych a HEAD gałęzi roboczej (bieżącego stanu repozytorium) leżą trzy rewizje <strong>zatwierdzone</strong>, nie niezatwierdzone zmiany robocze: prac nad skryptami wydania (ścieżka wydania — <code>tauri.conf.json</code> z <code>beforeBuildCommand</code> i <code>devUrl</code>, <code>budowa/scripts/wydanie.sh</code> i <code>pakowanie.sh</code>, brama mierząca wywołania, usunięcie <code>serde_json</code> z zależności bezpośrednich Powłoki), prac nad zasilaniem widoków z rdzenia (domknięcie warstwy wizualnej v2.0, pulpit dowodzenia zasilany z Rdzenia zamiast danych przykładowych, testy widoków) oraz porządek dokumentacyjny. Egzemplarz zbudowano z zatwierdzonej bieżącego stanu repozytorium; niniejszy opis dotyczy tej rewizji, a nie późniejszego stanu katalogu roboczego repozytorium, w którym w chwili sporządzania tego zestawienia równolegle trwa kolejna, odrębna fala prac nieobjęta niniejszym dokumentem. Różnica wobec rewizji stanu sprzed prac naprawczych nie ogranicza się do konfiguracji budowy Powłoki: obejmuje również pulpit dowodzenia i konsumenta zdarzenia <code>progress.changed</code> (rozdz. 10.3 wykaz ograniczeń, pkt 6 i 16), usunięcie <code>serde_json</code> jako zależności bezpośredniej (rozdz. 7.4, A.3) oraz stronę główną wywołującą <code>session.list</code>/<code>session.bind</code>/<code>session.focus</code> (rozdz. 10.3 pkt 13). Twierdzenia niniejszego dokumentu, poza punktami zaktualizowanymi wprost, odnoszą się do stanu potwierdzonego w bieżącym stanie repozytorium; różnicę wobec stanu sprzed prac naprawczych oznaczono w każdym z punktów, których dotyczy (rozdz. 10.3 wykaz ograniczeń, pkt 6, 13, 16, 18). Przed każdym kolejnym Wydaniem należy powtórzyć weryfikację na rewizji wówczas faktycznie wydawanej.</li>
<li>Twierdzenia o zgodności niniejszego dokumentu z pozostałymi dokumentami produktu zweryfikowano przez bezpośredni odczyt tych dokumentów w katalogu produktu. Rozbieżność co do liczby pozycji zestawu ikon, opisywana we wcześniejszej redakcji niniejszej Licencji, nie występuje w treści <code>README.md</code> ani <code>INSTRUKCJA-UZYTKOWANIA.md</code> na dzień tej weryfikacji — oba dokumenty posługują się liczbą 82 zgodnie z manifestem warstwy wizualnej; odpowiadająca temu pozycja Załącznika B została usunięta (rozdz. 7.6).</li>
</ol>
<h3 id="c4-zakres-wymagający-ponownego-przeglądu-przy-każdym-wydaniu">C.4. Zakres wymagający ponownego przeglądu przy każdym Wydaniu</h3>
<ol type="1">
<li>rozdz. 7 i Załącznik A — zestawienie Komponentów i ich licencji, w tym wersje rozwiązane zależności bezpośrednich warstwy Powłoki oraz liczba i źródła pozycji zestawu ikon;</li>
<li>rozdz. 10.3 — ujawnienie rzeczywistego stanu wykonania, ze szczególnym uwzględnieniem pkt 6 (źródło danych pulpitu dowodzenia), pkt 13 (wyjątek nawigacji platformy — sesje w tle), pkt 16 (konsument <code>progress.changed</code>), pkt 18 (instalator NSIS Powłoki), pkt 21 (zakres ustawień stosowanych przy wywołaniu) i pkt 22 (punkty dostępu rodzaju <code>localDirectory</code>);</li>
<li>rozdz. 3.4 — zależności zewnętrzne warunkujące działanie, w tym podział na czynności dostępne z interfejsu i czynności wymagające komendy Kontraktu;</li>
<li>rozdz. 9.2 i 9.4 — zakres danych zapisywanych i przekazywanych na zewnątrz, w tym moment powstania wiersza sesji i okna oraz brak utrwalenia prowenancji;</li>
<li>rozdz. 2.2 — zestaw pól koperty Kontraktu, wiążący co do brzmienia na podstawie rozdz. 2.6 pkt 5;</li>
<li>rozdz. 1.4 pkt 3 i rozdz. 3.5 — wykaz plików dokumentacji produktu objętej Licencją, wraz z ich zgodnością co do stanu opisywanego produktu.</li>
</ol>
<p><strong>Niniejszy dokument stanowi wersję ustaloną: rozstrzygnięcia Operatora z dnia 2026-08-26 zostały wpisane do treści (Załącznik B). Przed dystrybucją komercyjną zalecany jest przegląd prawny (rozdz. 1.5); przed dystrybucją zewnętrzną wymagane jest uzupełnienie pól rejestrowych w rozdz. 14.7.</strong></p>
<hr />
<p><em>Koniec dokumentu. Licencja produktu — Akt prawny, wersja 2.1, 2026-08-26.</em></p>
<hr />
<p><em>Danaco Console — Platforma AI Workspace OS · v2.0 · status Deweloperski</em> <em>© 2026 Danaco Holding Group Sp. z o.o. Wszelkie prawa zastrzeżone — Dariusz Naharnowicz.</em></p>
`;
