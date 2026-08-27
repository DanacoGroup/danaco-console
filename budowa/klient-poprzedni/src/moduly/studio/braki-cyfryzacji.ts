/**
 * Narzędziownia cyfryzacji opisuje dwa różne braki: brak pozycji kontraktu jest brakiem produktu,
 * a brak składnika pakietu serwera jest usterką wdrożenia — dwóch nie wolno zlewać.
 */

/** Braki, które są usterką wdrożenia serwera, na przykład brak silnika rozpoznawania pisma albo rozpakowywacza archiwów, nie ograniczeniem produktu. */
export const SKLADNIKI_PAKIETU_SERWERA = {
  rozpoznanie:
    'Rozpoznanie tekstu liczy silnik rozpoznawania pisma, który JEST składnikiem pakietu serwera ' +
    '(Tesseract wraz z pakietem językowym polskim). Odmowa z jego powodu nie znaczy, że produkt ' +
    'tej funkcji nie ma — znaczy, że TEN serwer jest niekompletny i brakuje mu składnika pakietu. ' +
    'To usterka wdrożenia do zgłoszenia administratorowi serwera, nie brak do obejścia ' +
    'w oknie i nie rzecz, którą Operator sobie doinstalowuje.',
  archiwum:
    'Rozpakowanie archiwum liczy rozpakowywacz, który JEST składnikiem pakietu serwera (7-Zip). ' +
    'Odmowa znaczy niekompletny serwer, nie brak funkcji: wsad z archiwum jest zbudowany ' +
    'i pracuje, gdy pakiet serwera jest pełny. Archiwum idzie wprost do kolejki rdzenia — ' +
    'komenda studio.ingest.queue.add przyjmuje jego ścieżkę polem archivePath.',
  mowa:
    'Zapis głosem liczy silnik mowy, który JEST składnikiem pakietu serwera (piper wraz z głosem, ' +
    'espeak-ng, whisper). Nagranie nie opuszcza maszyny rdzenia. Odmowa znaczy niekompletny ' +
    'serwer — składnik do uzupełnienia w pakiecie.',
  urzadzenie:
    'Skaner i kamera należą do maszyny RDZENIA i wykaz oddaje je komenda ' +
    'studio.ingest.device.list. Pusty wykaz nie jest odmową ani usterką: znaczy, że ta maszyna nie ' +
    'ma podłączonego urządzenia wejściowego. Okno mówi to przed próbą, zamiast pozwolić ' +
    'Operatorowi nacisnąć i dostać błąd.',
} as const;

/**
 * Czym sterują nastawy rozpoznawania — zdania przy kontrolkach.
 *
 * Każde mówi, co nastawa naprawdę robi i którą komendą jedzie. Zdanie „kontrakt
 * tego nie niesie" nie stoi tu już ani raz, bo przestało być prawdą.
 */
export const NASTAWY_CYFRYZACJI = {
  silnik:
    'Silnik rozpoznawania: lokalny stoi na maszynie rdzenia i materiał jej nie opuszcza, chmurowy ' +
    'jest usługą zewnętrzną wołaną jawnym kluczem z okna konfiguracji. Brak wskazania bierze silnik ' +
    'z katalogu ustawień, a nie wybór okna.',
  jezyki:
    'Zestaw języków w oznaczeniu trzyliterowym, rozdzielony przecinkiem — pismo pisane po polsku ' +
    'i po angielsku naraz rozpoznaje się DWOMA językami, nie jednym. Brak wskazania bierze zestaw ' +
    'z katalogu ustawień.',
  prog:
    'Próg pewności: pozycja, której rozpoznanie wypadło poniżej progu, wraca w stan ponowienia ' +
    'zamiast wejść do edytora jako gotowa. Ponowienie jest stanem osobnym — pozycja nie jest ani ' +
    'gotowa, ani odmówiona.',
  czyszczenie:
    'Czyszczenie obrazu przed rozpoznaniem: prostowanie skosu, odszumianie, progowanie do dwóch ' +
    'poziomów i przycięcie marginesów. Wykonuje je rdzeń w ramach rozpoznania, na kopii — źródło ' +
    'zostaje nietknięte, więc próba nic nie kosztuje.',
  uklad:
    'Odtwarzanie układu: kolumny, tabele, nagłówki, stopki i przypisy wracają jako bloki wraz ' +
    'z położeniem na stronie. Bez tego wynikiem jest jeden ciąg tekstu, w którym dwie szpalty ' +
    'zlewają się w jedną.',
  korekta:
    'Poprawianie rozpoznanych słów idzie PRZED przyjęciem wyniku do edytora: każde słowo niesie ' +
    'położenie na stronie i pewność rozpoznania, więc pracuje się na obrazie obok tekstu. Poprawka ' +
    'wchodzi na warstwę tekstową pozycji kolejki, a nie do dokumentu — dokument powstaje dopiero ' +
    'przy przyjęciu.',
  przyjecie:
    'Przyjęcie wyniku zakłada dokument roboczy WRAZ Z PIERWSZĄ WERSJĄ w repozytorium sesji. Wiele ' +
    'pozycji wskazanych naraz składa się w jeden dokument w kolejności podania — tak wchodzi skan ' +
    'wielostronicowy rozłożony na pliki.',
  adres:
    'Pobranie strony sieciowej z oczyszczeniem z nawigacji i reklam idzie komendą studio.ingest.url ' +
    'i kończy się pozycją kolejki. Wniesienie fragmentu strony WPROST do dokumentu jest czynnością ' +
    'inną i idzie studio.insert.from.web wraz z zapisem pochodzenia.',
} as const;

/**
 * Zdania o czynnościach, które kontrakt niesie POZA rodziną `studio.ingest.*` —
 * żeby brak drogi z tego okna nie wyglądał na brak zdolności platformy.
 */
export const DROGI_CYFRYZACJI = {
  zalaczniki:
    'Załącznik wiadomości wchodzi tu jako zasób magazynu: komenda mail.message.get wciąga ' +
    'załączniki do magazynu rdzenia i oddaje ich identyfikatory. Skopiuj identyfikator do pola ' +
    'wskazania i wybierz źródło „zasób magazynu".',
  library:
    'Plik repozytorium wchodzi tą samą drogą — kolejka wczytywania przyjmuje zasób magazynu rdzenia ' +
    'zarówno z Designu, jak i z Library.',
  czyszczenie:
    'Poprawki obrazu wykonane osobno (komendy obszaru obrazu) oddają NOWY zasób, a pozycja kolejki ' +
    'zakłada się na nim. Źródło zostaje nietknięte. Czyszczenie w ramach samego rozpoznania idzie ' +
    'nastawami wyżej i nie wymaga osobnego zasobu.',
  archiwum:
    'Wsad z archiwum rozpakowuje rdzeń w ramach dołożenia do kolejki — komenda ' +
    'studio.ingest.queue.add przyjmuje ścieżkę archiwum i oddaje pozycje, które z niego powstały. ' +
    'Rozpakowywacz jest składnikiem pakietu serwera.',
} as const;
