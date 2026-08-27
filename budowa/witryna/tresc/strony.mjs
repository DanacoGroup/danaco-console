// Treść ośmiu stron pisanych (dziewiąta — „Pobierz" — powstaje z wykazu wydań).
//
// Co wolno tu napisać. Wyłącznie to, co produkt naprawdę robi. Witryna jest
// częścią produktu i podlega tej samej regule, co pasek uczciwości w aplikacji:
// zdanie obiecujące zdolność, której nie ma, jest cięższe niż brak zdania.
// Zdolności, których produkt nie ma bez dołożenia czegoś z zewnątrz, są nazwane
// po imieniu na stronie „Moduły" — bo Operator ma się o tym dowiedzieć tu,
// a nie po instalacji.
//
// SKĄD BIORĄ SIĘ OBRAZY NA TEJ WITRYNIE. Z gotowych makiet okien dostawy
// Właściciela: `design/05-okna/` — 29 makiet HTML (platformowe, moduły, przepływ,
// środowiska), wyrenderowanych do plików PNG przeglądarką w trybie bez okna
// (`chromium --headless --screenshot`, szerokość 1600 px, makiety mobilne 520 px)
// i złożonych w `witryna/media/okna/`. Żaden obraz na tej stronie nie jest
// rysunkiem poglądowym ani wizją tego, jak produkt mógłby wyglądać — każdy jest
// zdjęciem makiety, którą Właściciel dostarczył. Dlatego podpis pod obrazem
// mówi, z której makiety pochodzi: czytający ma wiedzieć, na co patrzy.
//
// Makieta to nie zrzut z działającego programu i strona nie ma tego zacierać —
// stąd zdanie o tym stoi przy pierwszym obrazie na stronie głównej.
const OKNA = 'media/okna';

/**
 * OBRAZ JEST WEJŚCIEM DO ŻYWEGO OKNA, NIE KOŃCEM DROGI.
 *
 * Makiety są interaktywne, więc obraz zawsze będzie od nich gorszy. Na stronach
 * pisanych stoi jednak obraz, a nie ramka — sześć ramek w jednej stronie znaczy
 * sześć wczytanych dokumentów z własnymi krojami, i strona o produkcie zaczyna
 * się ładować dłużej niż sam produkt. Rozwiązanie: obraz jest ODNOŚNIKIEM do
 * prototypu, a wszystkie okna naraz stoją na stronie „Prototypy", gdzie ramka
 * jest jedna i wczytuje się to, co wybrano.
 *
 * Ścieżka prototypu wyliczana jest z nazwy zrzutu (`moduly-studio` →
 * `05-okna/moduly/studio.html`), a nie podawana drugi raz obok — dwa zapisy
 * tej samej rzeczy rozjeżdżają się pierwszego dnia.
 */
function adresPrototypu(nazwa) {
  const rodziny = ['moduly', 'srodowiska', 'przeplyw', 'platformowe'];
  const rodzina = rodziny.find((r) => nazwa.startsWith(`${r}-`));
  if (!rodzina) return `prototypy/05-okna/WZORZEC-STANOWISKA.html`;
  return `prototypy/05-okna/${rodzina}/${nazwa.slice(rodzina.length + 1)}.html`;
}

/** Jeden obraz wraz z podpisem. `nazwa` to nazwa pliku bez rozszerzenia. */
function zrzut(nazwa, podpis) {
  const prototyp = adresPrototypu(nazwa);
  return `    <figure class="dc-zrzut">
      <a href="${prototyp}" title="Otwórz działający prototyp tego okna">
        <img src="${OKNA}/${nazwa}.png" alt="${podpis}" loading="lazy" decoding="async">
      </a>
      <figcaption>${podpis} <a href="${prototyp}">Otwórz prototyp ↗</a></figcaption>
    </figure>`;
}

/** Zestaw obrazów. `para: true` układa je w dwie kolumny na szerokim ekranie. */
function zrzuty(pozycje, { para = true } = {}) {
  return `  <div class="dc-zrzuty${para ? ' dc-zrzuty--para' : ''}">
${pozycje.map(([nazwa, podpis]) => zrzut(nazwa, podpis)).join('\n')}
  </div>`;
}

// Kolejność w tablicy jest kolejnością w nawigacji.
export const STRONY = [
  {
    adres: 'index.html',
    nazwa: 'Start',
    tytul: 'Danaco Console',
    opis: 'Stanowisko pracy z modelami AI — jedno okno, wiele modułów, praca na własnej maszynie.',
    tresc: `  <h1>Jedno stanowisko pracy z modelami AI</h1>
  <p class="lead">Danaco Console łączy rozmowę z modelem, pracę na plikach, agentów
  i automatyzacje w jednym oknie — na maszynie Operatora, bez oddawania treści pracy
  komukolwiek po drodze.</p>

${zrzuty([['przeplyw-03-centrum-dowodzenia', 'Centrum dowodzenia — przedpokój przed strefą roboczą: karty czterech środowisk, kafle wytworów, listwa ustawień.']], { para: false })}
  <p style="max-width:78ch"><small>Obrazy na tej witrynie to <strong>makiety okien z dostawy
  warstwy wizualnej</strong>, wyrenderowane bez retuszu — nie zrzuty z uruchomionego programu
  i nie wizje tego, jak produkt mógłby wyglądać. Podpis pod każdym obrazem mówi, którą makietę
  widać.</small></p>

  <h2>Cztery środowiska pracy</h2>
${zrzuty([
  ['srodowiska-talkin', 'TalkIn — rozmowa, wiedza i praca z treścią.'],
  ['srodowiska-workspace', 'WorkSpace — produktywność, organizacja i realizacja projektów.'],
  ['srodowiska-codestudio', 'CodeStudio — programowanie: edytor, terminal, kontrola wersji.'],
  ['srodowiska-multitaskingai', 'MultitaskingAI — panel orkiestracji pracy ciągłej.'],
])}

  <div class="grid grid-2">
    <div class="area">
      <h3>Wiele okien rozmowy naraz</h3>
      <p>Każde okno ma własny model, własną rolę i własny katalog roboczy.
      Koordynator i wykonawca mogą pracować obok siebie na jednym ekranie.</p>
    </div>
    <div class="area">
      <h3>Rdzeń u siebie</h3>
      <p>Rdzeń stoi na maszynie Operatora i domyślnie nasłuchuje wyłącznie na pętli
      zwrotnej. Rozmowa, historia i pliki zostają tam, gdzie powstały.</p>
    </div>
    <div class="area">
      <h3>Zero bramek po zalogowaniu</h3>
      <p>Jedyne miejsce, w którym produkt o cokolwiek pyta, to logowanie.
      Potem nie ma potwierdzeń, blokad ani okien „czy na pewno".</p>
    </div>
    <div class="area">
      <h3>Zapis, który nie znika</h3>
      <p>Zamknięcie okna niczego nie kasuje. Sesja trafia do bazy od założenia,
      a usunięcie jest osobną, jawną czynnością.</p>
    </div>
  </div>

  <div class="dc-uwaga">
    <p><strong>Produkt jest w budowie.</strong> Rdzeń i interfejs działają. Wydanie
    <strong>1.0.0</strong> z dnia 18.08.2026 jest zbudowane i <strong>da się je stąd
    pobrać</strong>: instalki Windows w dwóch architekturach — <strong>x64</strong>
    i <strong>ARM64</strong> — wraz z sumami kontrolnymi wykłada strona
    <a href="pobierz.html">Pobierz</a>. Osobno stoi tam pakiet serwera wdrożenia dla
    administratora: instalka niesie samo okno, a rdzeń pracuje na serwerze.
    Pobieranie jest <strong>chronione hasłem</strong> i na razie zastrzeżone — strona
    „Pobierz" mówi o tym przed przyciskami, żeby okienko z hasłem nie wyglądało na usterkę.
    Ta witryna mówi o stanie faktycznym, nie o zamierzonym.</p>
  </div>`,
  },
  {
    adres: 'mozliwosci.html',
    nazwa: 'Możliwości',
    tytul: 'Możliwości',
    opis: 'Co Danaco Console robi dziś: okna równoległe, agenci, kolejki, terminal, historia.',
    tresc: `  <h1>Co produkt robi</h1>
  <p class="lead">To, co produkt robi na uruchomionym rdzeniu — nie zapowiedzi.</p>

  <h2>Rozmowa i okna</h2>
  <p>Rozmowa z modelem toczy się w oknie osadzonym w module, nie zamiast niego —
  widok modułu i rozmowa stoją obok siebie. Okien równoległych może być kilka,
  każde z innym modelem i inną rolą, a scena pokazuje je jednocześnie.</p>

${zrzuty([
  ['moduly-studio', 'Studio — praca z modelem nad treścią: wersje, różnice, podgląd.'],
  ['srodowiska-mtai-okna-rol', 'Okna ról — koordynator i wykonawca obok siebie na jednej scenie.'],
])}

  <h2>Agenci i zespoły</h2>
  <p>Agent to zapisana tożsamość: model, rola, wskazówki, wersje. Agentów da się
  wersjonować, archiwizować i przywracać, a z kilku ułożyć zespół pracujący nad
  jednym zadaniem.</p>

  <h2>Automatyzacje i kolejki</h2>
  <p>Powtarzalna praca układa się w przebieg z krokami i zależnościami. Kolejka trzyma
  pozycje w ustalonej kolejności i przeprowadza je przez stany — <strong>ale posuwa je
  działanie Operatora</strong> (start, ponów, wstrzymaj, przerwij), a nie ona sama:
  po starcie rusza pozycja bieżąca, nie cały wykaz. Kolejka melduje <strong>swój</strong>
  stan i licznik obiegów; <strong>stanu poszczególnych pozycji dziś nie wystawia</strong>.</p>

${zrzuty([
  ['moduly-agents', 'Agents — agenci, ich wersje, archiwum i zespoły.'],
  ['moduly-automations', 'Automations — przebiegi, kroki, zależności, harmonogram.'],
])}

  <h2>Terminal i pliki</h2>
  <p>Wbudowany terminal i podgląd plików pracują w katalogu roboczym okna —
  tym samym, który widzi model.</p>

${zrzuty([
  ['moduly-terminal', 'Terminal — powłoka w katalogu roboczym okna.'],
  ['moduly-workspace', 'Workspace — przestrzeń robocza sesji: pliki, katalog, pamięć.'],
])}

  <h2>Historia i retencja</h2>
  <p>Historia rozmów trafia do bazy rdzenia. Zasada retencji rozstrzyga, co i jak długo
  ma zostać; przemiatanie wykonuje się przy starcie rdzenia i mówi w dzienniku, ile
  pozycji przycięło.</p>

  <h2>Dyktowanie</h2>
  <p>Mowa zamieniana na tekst <strong>lokalnie</strong> — dźwięk nie opuszcza maszyny.
  Rdzeń niesie pomocnik rozpoznawania mowy, ale <strong>samego silnika nie dostarcza</strong>:
  potrzebny jest Python z pakietem <code>faster-whisper</code>, zakładany osobno
  (<code>python -m pip install -r wymagania.txt</code>). Dopóki go nie ma, rdzeń mówi wprost,
  że dyktowanie jest niedostępne i czego brakuje — zamiast oddawać zmyślony zapis.</p>`,
  },
  {
    adres: 'moduly.html',
    nazwa: 'Moduły',
    tytul: 'Moduły',
    opis: 'Piętnaście modułów Danaco Console i uczciwy stan każdego z nich.',
    tresc: `  <h1>Moduły</h1>
  <p class="lead">Produkt niesie piętnaście modułów. Poniżej stan każdej rodziny —
  łącznie z tymi, które dziś są rejestrami bez zdolności wykonawczych.</p>

  <h2>Działają</h2>
  <table>
    <tbody>
      <tr><td><strong>Studio</strong></td><td>Praca z modelem nad treścią: wersje, różnice, podgląd.</td></tr>
      <tr><td><strong>Workspace</strong></td><td>Przestrzeń robocza sesji: pliki, katalog, pamięć.</td></tr>
      <tr><td><strong>Agents</strong></td><td>Agenci, ich wersje, archiwum i zespoły.</td></tr>
      <tr><td><strong>Multitasking</strong></td><td>Okna równoległe, obsada ról, plan etapów.</td></tr>
      <tr><td><strong>Terminal</strong></td><td>Powłoka w katalogu roboczym okna.</td></tr>
      <tr><td><strong>Automations</strong></td><td>Przebiegi, kroki, zależności, harmonogram.</td></tr>
      <tr><td><strong>Diagnostics</strong></td><td>Dziennik rdzenia i panel odmów.</td></tr>
      <tr><td><strong>Developer</strong></td><td>Wgląd w kontrakt, komendy i stan rdzenia.</td></tr>
      <tr><td><strong>Apps</strong></td><td>Prowadzi wdrożenie przez stany aż do wyniku i umie cofnąć je do wydania, które się powiodło. Nie hostuje produktu — adresu wdrożenia nie ma, bo rdzeń nie stawia serwera.</td></tr>
      <tr><td><strong>Translate</strong></td><td>Tłumaczy kanałem modelu, związany słownikiem Operatora: odpowiedniki terminów i nazwy wyłączone z przekładu obowiązują bezwzględnie. Kontrola jakości porównuje przekład ze źródłem — liczby, waluty, znaczniki podstawienia. Syntezuje też mowę: z panelu przekładu powstaje plik WAV na dysku Operatora — o ile w systemie stoi <code>piper</code> albo <code>espeak-ng</code>, bo instalka ich nie niesie; bez nich rdzeń mówi, czego brakuje.</td></tr>
      <tr><td><strong>Assistant</strong></td><td>Przyjmuje polecenie głosem lub tekstem, wykonuje je turą modelu i melduje wynik. Zlecenie da się wstrzymać, wznowić i ponowić.</td></tr>
      <tr><td><strong>Library</strong></td><td>Trzyma treść plików wiedzy na dysku obok bazy, pod sumą kontrolną, z wersjami. Wyszukiwanie dopasowuje <strong>nazwę pliku i jego treść</strong> (indeks pełnotekstowy) — dopasowania po ZNACZENIU dziś nie ma i nikt go nie udaje.</td></tr>
      <tr><td><strong>Browser</strong></td><td>Pobiera wskazaną stronę i robi z niej migawkę: tytuł, tekst, HTML. Nie wykonuje JavaScriptu strony i nie robi zrzutów ekranu — strona budowana w całości skryptem da migawkę szczątkową.</td></tr>
    </tbody>
  </table>

  <h2>Moduły w oknie</h2>
${zrzuty([
  ['moduly-library', 'Library — pliki wiedzy na dysku, pod sumą kontrolną, z wersjami.'],
  ['moduly-translate', 'Translate — przekład związany słownikiem Operatora, z kontrolą jakości.'],
  ['moduly-assistant', 'Assistant — polecenie głosem lub tekstem, wykonanie turą modelu i meldunek.'],
  ['moduly-diagnostics', 'Diagnostics — dziennik rdzenia i panel odmów.'],
  ['moduly-apps', 'Apps — wdrożenie prowadzone przez stany, z cofnięciem do wydania, które się powiodło.'],
  ['moduly-developer', 'Developer — wgląd w kontrakt, komendy i stan rdzenia.'],
])}

  <h2>Czego brakuje — stan na dziś</h2>
  <p>Poniżej to, czego produkt dziś <strong>nie zrobi bez dołożenia czegoś z zewnątrz</strong>.
  Piszemy to tutaj, żeby nikt nie dowiedział się o tym po instalacji:</p>
  <table>
    <tbody>
      <tr><td><strong>Design</strong></td><td>Prowadzi prompty strukturalne, zasoby i kompozycje: etykiety, filtrowanie, warstwy planszy. <strong>Obrazów nie wygeneruje, dopóki Operator nie założy kanału obrazowego</strong> — sam rdzeń umie z takim kanałem rozmawiać (adapter obrazowy i magazyn zasobów są na miejscu), brakuje wyłącznie wiersza w rejestrze kanałów: <code>channel.add</code> z <code>kind="api"</code> i <code>config.adapter="obrazy"</code>. Bez niego prompt składa się w gotowe polecenie i wraca Operatorowi, a odmowa mówi wprost, czego brakuje i jak to założyć.</td></tr>
      <tr><td><strong>Dyktowanie</strong></td><td>Rdzeń niesie pomocniki rozpoznawania mowy, ale <strong>silnika nie ma w instalce</strong>: Python z pakietem <code>faster-whisper</code> Operator zakłada sam. Do tego czasu rdzeń zgłasza dyktowanie jako niedostępne wraz z powodem — patrz <a href="wymagania.html">Wymagania</a>.</td></tr>
    </tbody>
  </table>`,
  },
  {
    // Treść tej strony powstaje z katalogu makiet w `prototypy.mjs`, nie tutaj —
    // wykaz okien czytany jest z dostawy przy składaniu, żeby strona nie mogła
    // obiecać okna, którego w dostawie nie ma.
    adres: 'prototypy.html',
    nazwa: 'Prototypy',
    tytul: 'Prototypy okien',
    opis: 'Działające prototypy okien Danaco Console: środowiska, moduły, przepływ, postać mobilna.',
    tresc: '',
  },
  {
    adres: 'bezpieczenstwo.html',
    nazwa: 'Bezpieczeństwo',
    tytul: 'Bezpieczeństwo',
    opis: 'Gdzie leżą dane, jak stoi rdzeń i co dokładnie chroni dostęp.',
    tresc: `  <h1>Bezpieczeństwo i prywatność</h1>
  <p class="lead">Produkt jest przeznaczony do pracy na treści, której nie wolno rozdawać.
  Stąd te rozstrzygnięcia.</p>

  <h2>Dane zostają u Operatora</h2>
  <p>Rdzeń, baza i pliki robocze leżą na maszynie Operatora. Rozmowy z modelem idą wyłącznie
  do modelu, który Operator wskazał — produkt nie pośredniczy w niczym więcej i nie odsyła
  treści pracy do wytwórcy.</p>

  <h2>Nasłuch domyślnie zamknięty</h2>
  <p>Rdzeń nasłuchuje domyślnie na pętli zwrotnej, czyli jest osiągalny wyłącznie z tej samej
  maszyny. Wystawienie szerzej jest możliwe, ale wymaga wskazania wprost i zostaje odnotowane
  w dzienniku wraz z tym, czego przy takim wystawieniu brakuje.</p>

  <h2>Jedna bramka, na wejściu</h2>
  <p>Hasło ustawia się przy pierwszym uruchomieniu; potem otwiera aplikację hasłem albo PIN-em.
  Przełącznik „nie wyloguj mnie" przedłuża wejście, żeby nie powtarzać go codziennie.
  Po zalogowaniu produkt nie pyta o nic więcej — brak blokad jest tu rozstrzygnięciem,
  nie niedopatrzeniem.</p>

${zrzuty([
  ['platformowe-ustawienia', 'Ustawienia — kanały modelu, poświadczenia, zasady retencji.'],
  ['platformowe-konfiguracja', 'Konfiguracja — ustawienia wdrożenia po stronie rdzenia.'],
])}

  <h2>Hasła i klucze</h2>
  <p>Hasło bramki nie leży nigdzie w postaci jawnej, a token wejścia zapisany jest w bazie
  wyłącznie jako skrót.</p>
  <p><strong>Inaczej jest z poświadczeniami do modeli.</strong> Rdzeń umie wskazywać je
  odwołaniem do sejfu, ale klucz podany wprost w konfiguracji kanału
  (<code>channel.add</code> z polem <code>apiKey</code>) <strong>trafia do bazy jawnym tekstem</strong>,
  razem z całą konfiguracją kanału. Baza leży na maszynie Operatora, więc nie wychodzi
  na zewnątrz — ale kto czyta plik bazy, czyta też ten klucz. Piszemy to tutaj, bo obietnica
  „klucze nie leżą w bazie" byłaby dziś nieprawdą.</p>

  <h2>Połączenie</h2>
  <p>Przy pracy na jednej maszynie połączenie nie opuszcza pętli zwrotnej. Przy wystawieniu
  poza nią rdzeń przyjmuje warstwę TLS wskazaną parą plików, a pochodzenie żądań jest
  sprawdzane wykazem — obcy adres nie otworzy gniazda do cudzego rdzenia.</p>`,
  },
  {
    adres: 'wymagania.html',
    nazwa: 'Wymagania',
    tytul: 'Wymagania',
    opis: 'Czego potrzebuje Danaco Console: system, sprzęt, dostęp do modeli.',
    tresc: `  <h1>Wymagania</h1>

  <h2>System</h2>
  <p>Stanowiskiem Operatora jest <strong>Windows 11</strong> — jedyna platforma produktu.
  Jedyny wybór, jaki zostaje pobierającemu, to architektura: <strong>x64</strong> dla
  procesorów Intel i AMD albo <strong>ARM64</strong>. Windows 10, Linux i macOS nie są
  obsługiwane. <strong>WebView2</strong> jest składnikiem Windows 11, więc instalka go nie
  niesie i nie zakłada.</p>

  <h2>Serwer wdrożenia</h2>
  <p>Instalka niesie samo okno wraz z powłoką — <strong>rdzenia w niej nie ma</strong>.
  Rdzeń i całe zaplecze stoją na serwerze wdrożenia, więc bez działającego serwera okno nie
  ma z czym rozmawiać, a praca bez łączności z nim nie jest przewidziana. Pakiet serwera
  zakłada administrator raz na maszynie serwerowej; jest osobną pozycją na stronie
  <a href="pobierz.html">Pobierz</a>.</p>

  <h2>Sprzęt</h2>
  <table>
    <tbody>
      <tr><td>Pamięć</td><td>8 GB wystarcza; 16 GB przy kilku oknach równoległych i dyktowaniu.</td></tr>
      <tr><td>Dysk</td><td>Około 250 MB na aplikację — zaplecze zostaje na serwerze wdrożenia.</td></tr>
      <tr><td>Sieć</td><td>Potrzebna wyłącznie do modeli działających zdalnie oraz do aktualizacji.</td></tr>
    </tbody>
  </table>

  <h2>Modele</h2>
  <p>Produkt nie dostarcza modeli. Operator wskazuje własne poświadczenia do modelu
  zdalnego albo model działający lokalnie.</p>

  <h2>Dyktowanie</h2>
  <p>Wymaga mikrofonu oraz <strong>osobno założonego silnika mowy</strong>. Instalka niesie
  wyłącznie pomocniki rozpoznawania (skrypty i wykaz zależności) — <strong>nie niesie ani
  Pythona, ani pakietu <code>faster-whisper</code>, ani modelu rozpoznawania</strong>. Operator
  zakłada je poleceniem <code>python -m pip install -r wymagania.txt</code>; do tej chwili rdzeń
  zgłasza dyktowanie jako niedostępne i mówi, czego brakuje.</p>`,
  },
  {
    adres: 'pobierz.html',
    nazwa: 'Pobierz',
    tytul: 'Pobierz',
    opis: 'Instalka Danaco Console i chronologia wszystkich wydań.',
    tresc: '',
  },
  {
    adres: 'pierwsze-kroki.html',
    nazwa: 'Pierwsze kroki',
    tytul: 'Pierwsze kroki',
    opis: 'Od instalacji do pierwszej rozmowy z modelem.',
    tresc: `  <h1>Pierwsze kroki</h1>
  <p class="lead">Od pobrania do pierwszej rozmowy — pięć czynności.</p>

  <h2>1. Instalacja</h2>
  <p>Pobierz plik ze strony <a href="pobierz.html">Pobierz</a> i uruchom go. Wybierz pozycję
  pasującą do architektury twojej maszyny: <strong>x64</strong> dla procesora Intel albo AMD,
  <strong>ARM64</strong> dla procesora ARM. Instalka niesie samo okno — po drugiej stronie
  musi stać serwer wdrożenia, założony przez administratora z osobnej pozycji na stronie
  „Pobierz".</p>
  <p><strong>Pobieranie jest dziś chronione hasłem</strong> — na czas wyłączności dostęp ma
  jedna osoba, a poświadczenie wydaje producent. W sprawie dostępu:
  <a href="mailto:support@danaco-group.pl">support@danaco-group.pl</a>.</p>

${zrzuty([
  ['przeplyw-01-okno-startowe', 'Okno startowe — pierwsze, co widzi Operator po uruchomieniu.'],
  ['przeplyw-02-rejestracja-logowanie', 'Rejestracja i logowanie — jedyna bramka w całym produkcie.'],
])}

  <h2>2. Ustawienie hasła</h2>
  <p>Pierwsze uruchomienie prosi o ustawienie hasła. To jedyna bramka w całym produkcie.
  Zaznacz „nie wyloguj mnie", jeśli pracujesz na własnej maszynie.</p>

  <h2>3. Wskazanie modelu</h2>
  <p>W Ustawieniach wskaż kanał modelu i poświadczenia. Bez tego okna rozmowy staną,
  ale reszta produktu działa.</p>

  <h2>4. Pierwsze okno</h2>
  <p>Ze strony głównej wejdź w środowisko i otwórz moduł. Rozmowa stoi po prawej stronie
  widoku modułu — nie zamiast niego.</p>

  <h2>5. Druga rozmowa obok pierwszej</h2>
  <p>Moduł Multitasking otwiera okna równoległe: każde z własnym modelem i rolą.
  Tak pracuje się parą koordynator–wykonawca.</p>`,
  },
  {
    adres: 'aktualizacje.html',
    nazwa: 'Aktualizacje',
    tytul: 'Aktualizacje',
    opis: 'Jak aktualizuje się Danaco Console: baner w oknie aplikacji i ponowne uruchomienie.',
    tresc: `  <h1>Aktualizacje</h1>
  <p class="lead">Aktualizacja jest jednym kliknięciem w oknie aplikacji.
  Ręczne pobieranie jest drogą zapasową, nie podstawową.</p>

  <h2>Baner „Aktualizuje"</h2>
  <p>Gdy na tej witrynie stanie wydanie nowsze niż zainstalowane, w oknie aplikacji pojawia
  się baner <strong>„Aktualizuje"</strong>. Baner działa jak przycisk: pobiera wydanie,
  zakłada je i uruchamia aplikację ponownie. Praca w toku zostaje zapisana wcześniej —
  sesje i historia leżą w bazie, nie w pamięci okna.</p>

  <h2>Czego baner nie robi</h2>
  <p>Nie aktualizuje sam z siebie i nie przerywa pracy. Dopóki nikt go nie kliknie,
  aplikacja działa dalej na wersji zainstalowanej. Nie ma też okna „czy na pewno" —
  kliknięcie banera JEST zgodą.</p>

  <h2>Wydania zostają na zawsze</h2>
  <p>Strona <a href="pobierz.html">Pobierz</a> trzyma pełną chronologię: każde wydanie
  z datą, także starsze. Powrót do wersji poprzedniej jest więc możliwy bez proszenia
  kogokolwiek o plik.</p>`,
  },
  {
    adres: 'wsparcie.html',
    nazwa: 'Wsparcie',
    tytul: 'Wsparcie',
    opis: 'Kontakt, zgłoszenia i informacje o producencie Danaco Console.',
    tresc: `  <h1>Wsparcie i kontakt</h1>

  <h2>Zgłoszenia</h2>
  <p>Usterki i pytania: <a href="mailto:support@danaco-group.pl">support@danaco-group.pl</a>.
  W zgłoszeniu warto podać wersję z okna Ustawień oraz to, co pokazał panel odmów
  w module Diagnostics — rdzeń zapisuje tam każdą odmowę wraz z powodem.</p>

  <h2>Producent</h2>
  <table>
    <tbody>
      <tr><td>Producent</td><td>Danaco Holding Group Sp. z o.o.</td></tr>
      <tr><td>Produkt</td><td>Danaco Console — platforma AI Workspace OS</td></tr>
      <tr><td>Stan</td><td>deweloperski</td></tr>
    </tbody>
  </table>

  <h2>Dane Operatora</h2>
  <p>Producent nie ma dostępu do treści pracy: rozmowy, pliki i historia leżą na maszynie
  Operatora. Zgłoszenie usterki nie wysyła niczego automatycznie — to, co ma trafić do
  wsparcia, Operator dołącza sam.</p>`,
  },
];
