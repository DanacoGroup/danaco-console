import type { PozycjaBezKomendy } from './panel-akcji';

/**
 * Teksty widoczne dla Operatora, wyjęte z plików budujących elementy.
 *
 * Osobny plik, bo to inna odpowiedzialność niż budowa widoku: zdania mówiące,
 * czego kontrakt nie niesie, dają się przeczytać w jednym miejscu i tam
 * poprawić, gdy komenda wejdzie do kontraktu. Każda pozycja nazywa czynność
 * okna i powód, dla którego nie ma za nią komendy.
 */

/** Czynności panelu Library Explorera bez komendy kontraktu. */
export const BEZ_KOMENDY_EXPLORER: readonly PozycjaBezKomendy[] = [
  {
    etykieta: 'Eksport',
    powod: 'kontrakt nie ma komendy eksportu zasobów biblioteki; nośnikiem miałby być panel akcji rdzenia',
  },
  {
    etykieta: 'Archiwizuj',
    powod: 'kontrakt nie ma komendy archiwizacji zasobu',
  },
  {
    etykieta: 'Przenieś',
    powod: 'kontrakt nie ma komendy przeniesienia zasobu w strukturze repozytorium',
  },
  {
    etykieta: 'Do kosza',
    powod: 'kontrakt nie ma komendy usunięcia odwracalnego ani kosza biblioteki',
  },
  {
    etykieta: 'Wykryj duplikaty',
    powod: 'wykrywanie duplikatów i niemal-duplikatów jest wymaganiem wobec rdzenia, nie komendą kontraktu',
  },
  {
    etykieta: 'Udostępnij linkiem',
    powod: 'kontrakt nie ma komendy udostępnienia zasobu odnośnikiem',
  },
];

/** Czynności panelu File Preview bez komendy kontraktu. */
export const BEZ_KOMENDY_PODGLAD: readonly PozycjaBezKomendy[] = [
  {
    etykieta: 'Porównaj z…',
    powod: 'porównywanie dokumentów zapisano jako wymaganie wobec rdzenia, nie jako komendę',
  },
  {
    etykieta: 'Pobierz',
    powod: 'kontrakt nie ma komendy pobrania treści pliku na urządzenie',
  },
];

/**
 * Czynności panelu Versioning Panel bez komendy kontraktu.
 *
 * Oznaczenia kamieniem milowym tu nie ma: `library.version.add` przyjmuje pole
 * `label`, więc etykieta wersji ma własne wejście (`dolozenie-wersji.ts`).
 */
export const BEZ_KOMENDY_WERSJE: readonly PozycjaBezKomendy[] = [
  {
    etykieta: 'Porównaj wersje',
    powod: 'kontrakt nie ma komendy porównania dwóch wersji zasobu biblioteki',
  },
  {
    etykieta: 'Archiwum historii',
    powod:
      'archiwum z treścią każdej wersji żąda bajtów wersji niebieżącej, których nie oddaje ' +
      'żadna komenda kontraktu; raport zmian bez treści składa się z historii pokazanej w oknie',
  },
];

/** Czynności panelu Tags & Collections bez komendy kontraktu. */
export const BEZ_KOMENDY_ETYKIETY: readonly PozycjaBezKomendy[] = [
  {
    etykieta: 'Kolekcja inteligentna',
    powod: 'kolekcja z regułą wymaga reguły po stronie rdzenia; library.collection.create przyjmuje wyłącznie nazwę i opis',
  },
  {
    etykieta: 'Połącz duplikaty etykiet',
    powod: 'kontrakt nie ma komendy łączenia etykiet',
  },
  {
    etykieta: 'Zmień kolor etykiety',
    powod: 'etykieta w kontrakcie jest samym napisem — nie ma pola koloru ani komendy jego zmiany',
  },
  {
    etykieta: 'Usuń nieużywaną etykietę',
    powod: 'kontrakt nie ma komendy usunięcia etykiety; library.tag.set zmienia wyłącznie etykiety pliku',
  },
  {
    etykieta: 'Zastosuj sugestię AI',
    powod: 'sugestia porządkująca pochodzi z paska narzędzi promptu modułu, którego kontrakt nie niesie',
  },
  {
    etykieta: 'Relacje tezaurusa',
    powod:
      'etykieta jest w kontrakcie samym napisem — nie ma pola relacji nadrzędny/podrzędny ' +
      'ani synonimu, więc wywóz SKOS wychodzi płaski i mówi to w swoim komentarzu',
  },
];

/**
 * Czynności panelu Metadata & Archive Panel bez komendy kontraktu.
 *
 * Wykaz jest długi, bo warstwa czwarta modułu jest obszerna, a kontrakt niesie
 * dla niej dokładnie jedną komendę: `knowledge.index` (przeliczenie wskaźnika
 * znaczenia, zakładka Higiena). Każda pozycja nazywa nie tylko brak, lecz jego
 * przyczynę — inaczej wykaz wyglądałby na listę rzeczy zapomnianych, a jest
 * pomiarem kontraktu.
 */
export const BEZ_KOMENDY_METADANE: readonly PozycjaBezKomendy[] = [
  {
    etykieta: 'Zapis metadanych Dublin Core',
    powod:
      'kontrakt nie ma komendy zapisu opisu zasobu biblioteki; library.tag.set zmienia ' +
      'wyłącznie etykiety i kolekcje pliku',
  },
  {
    etykieta: 'Metadane techniczne EXIF/IPTC/XMP/ID3',
    powod:
      'kontrakt nie ma komendy odczytu metadanych osadzonych w pliku biblioteki; media.inspect ' +
      'rozwiązuje assetId przez repozytorium zasobów modułu Design, więc identyfikator pliku ' +
      'biblioteki wraca stamtąd odmową. Materiał dźwiękowy i filmowy rozpoznaje się tą komendą ' +
      'w obszarze Archiwum, ale ze ścieżki na dysku Operatora — i o strumieniach, nie o EXIF',
  },
  {
    etykieta: 'Pole niestandardowe',
    powod:
      'definicja pola własnego żąda schematu metadanych po stronie rdzenia; kontrakt nie ma ' +
      'ani komendy definicji, ani miejsca na jej wartość przy zasobie',
  },
  {
    etykieta: 'Konwersja archiwalna PDF/A',
    powod:
      'document.convert rozwiązuje assetId przez repozytorium zasobów modułu Design, a plik ' +
      'biblioteki leży w innym rejestrze; profilu PDF/A kontrakt zresztą nie wymienia wśród ' +
      'formatów docelowych',
  },
  {
    etykieta: 'Pakiet archiwalny BagIt z manifestem sum kontrolnych',
    powod:
      'formatu BagIt kontrakt nie zna w żadnej komendzie. archive.pack pakuje zip, 7z albo ' +
      'tar.gz ze ścieżki na dysku Operatora i oddaje wynik jako zasób magazynu Designu — ' +
      'czynność stoi w obszarze Archiwum, ale pakietem archiwalnym biblioteki nie jest ' +
      'i tak się nie nazywa',
  },
  {
    etykieta: 'Metadane utrwalenia PREMIS/METS',
    powod: 'kontrakt nie ma komendy zapisu ani odczytu profilu utrwalenia zasobu',
  },
  {
    etykieta: 'Polityka retencji zasobów',
    powod:
      'reguła retencji należy do okna konfiguracji; moduł nie ma komendy, którą odczytałby ' +
      'ustawienie retencja_historii w swoim zasięgu',
  },
  {
    etykieta: 'Weryfikacja integralności',
    powod:
      'sprawdzenie sumy kontrolnej żąda przeliczenia jej z bajtów, a kontrakt nie niesie ' +
      'komendy pobierającej treść zasobu biblioteki; porównanie sumy z nią samą weryfikacją ' +
      'nie jest',
  },
  {
    etykieta: 'Niemal-duplikaty',
    powod:
      'wykrycie podobnej treści przy różnych nazwach żąda porównania treści całego zbioru; ' +
      'klient ma podgląd jednego pliku na wywołanie i żadnej komendy liczącej podobieństwo',
  },
  {
    etykieta: 'Scal duplikaty',
    powod: 'kontrakt nie ma komendy scalenia dwóch zasobów w jeden',
  },
  {
    etykieta: 'Normalizacja nazw plików',
    powod:
      'zmiana nazwy zasobu nie ma komendy; library.tag.set zmienia wyłącznie etykiety ' +
      'i kolekcje',
  },
  {
    etykieta: 'Dziennik audytu',
    powod: 'kontrakt nie ma komendy odczytu zdarzeń repozytorium biblioteki',
  },
  {
    etykieta: 'API repozytorium i webhooki',
    powod:
      'włączenie punktu końcowego i token dostępu należą do zakresu Integracje okna ' +
      'konfiguracji, nie do okna operacyjnego modułu',
  },
  {
    etykieta: 'Macierz izolacji dostępu do plików',
    powod:
      'osiem zakresów izolacji technicznej prowadzi okno punktów izolacji komendami ' +
      'isolation.technical.get i isolation.technical.set; druga macierz w tym module byłaby ' +
      'drugą prawdą o jednym stanie',
  },
];
