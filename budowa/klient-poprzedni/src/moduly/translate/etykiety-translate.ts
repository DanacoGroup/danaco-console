import { ExportFormat, TranslationStatus } from '../../../../shared/contract';
import type { ZdanieStanu } from './stan-okna-translate';

/**
 * Teksty i katalogi widoczne dla Operatora w module Translate.
 *
 * Wyjęte z plików budujących elementy, żeby zmiana zdania nie była zmianą
 * widoku — tak samo jak `sterowanie/etykiety-sterowania.ts`.
 *
 * Dwa rodzaje katalogu są rozdzielone:
 *
 *  — katalog zamknięty pochodzi z kontraktu i jest listą wyboru: formaty
 *    eksportu (`ExportFormat`) i stany panelu (`TranslationStatus`) mają
 *    w kontrakcie skończony zbiór wartości, więc lista wyboru jest tu prawdą;
 *  — katalog otwarty jest wyłącznie podpowiedzią: kontrakt przyjmuje język
 *    i ton jako dowolny napis i nie ma komendy zwracającej ich wykaz. Pole
 *    zostaje edytowalne, a podpowiedź nie udaje katalogu rdzenia.
 */

/** Podpowiedź języków docelowych — pole pozostaje otwarte na kod spoza wykazu. */
export const PODPOWIEDZ_JEZYKOW: readonly string[] = [
  'pl',
  'en',
  'de',
  'fr',
  'es',
  'it',
  'uk',
  'cs',
];

/** Podpowiedź tonu tłumaczenia — kontrakt przyjmuje dowolny napis. */
export const PODPOWIEDZ_TONU: readonly string[] = [
  'neutralny',
  'formalny',
  'nieformalny',
  'perswazyjny',
  'empatyczny',
  'techniczny',
];

/** Formaty eksportu panelu — katalog zamknięty kontraktu (`ExportFormat`). */
export const FORMATY_EKSPORTU: readonly { wartosc: ExportFormat; etykieta: string }[] = [
  { wartosc: ExportFormat.Docx, etykieta: 'DOCX' },
  { wartosc: ExportFormat.Txt, etykieta: 'Tekst zwykły' },
  { wartosc: ExportFormat.Markdown, etykieta: 'Markdown' },
  { wartosc: ExportFormat.Pdf, etykieta: 'PDF' },
  { wartosc: ExportFormat.Html, etykieta: 'HTML' },
];

/**
 * Wartość steru kanału, gdy kanał nie został wskazany.
 *
 * Na uchwycie steru stoi wartość, nie nazwa nastawy — a wartością „brak
 * wskazania" jest to, co robi wtedy rdzeń: bierze kanał czynny okna
 * (`TranslateTargetAddRequest.channelId`).
 */
export const WARTOSC_KANALU_OKNA = 'kanał czynny okna';

/**
 * Ton panelu widoczny w nagłówku instancji — jedna prawda o nastawie.
 *
 * Rdzeń oddaje `panel.tone` w odpowiedziach `target.add`, `translation.set`
 * i `panel.tone.set`; nagłówek pokazuje tę właśnie wartość, a nie treść pola
 * wejściowego.
 *
 * Pusty ton jest stanem, nie brakiem danych: kontrakt nie wymaga tonu, a panel
 * bez niego przekłada się tonem domyślnym rdzenia. Zdanie mówi to wprost,
 * zamiast pokazywać pustkę albo kreskę.
 */
export function opisTonuPanelu(ton: string | undefined): string {
  const nazwa = (ton ?? '').trim();
  return nazwa === '' ? 'ton domyślny rdzenia' : `ton: ${nazwa}`;
}

/**
 * Zdanie o języku źródłowym, który trzyma rdzeń — druga strona jednej prawdy.
 *
 * Pole „Język źródłowy" jest wejściem: `source.detect` wpisuje tam rozpoznanie,
 * którego rdzeń u siebie jeszcze nie ma. Gdyby to samo pole przerysowywać
 * wartością rdzenia, każde ogłoszenie stanu kasowałoby świeże rozpoznanie.
 * Wartość bieżącą pokazuje więc zdanie obok pola, czytając `stan.jezykZrodlowy()`.
 *
 * Rozbieżność jest nazwana, nie uzgodniona po cichu: okno nic samo nie wysyła
 * ani nie nadpisuje, mówi tylko, że obie wartości się różnią i co je uzgodni
 * (zapis źródła).
 */
export function zdanieJezykaRdzenia(
  tekstRdzenia: string,
  jezykRdzenia: string,
  wPolu: string,
): string {
  const rdzen = jezykRdzenia.trim();
  const pole = wPolu.trim();
  if (tekstRdzenia === '') {
    return 'Rdzeń nie ma jeszcze tekstu źródłowego tego okna, więc nie ma też jego języka.';
  }
  if (rdzen === '') {
    return pole === ''
      ? 'Rdzeń trzyma tekst źródłowy BEZ języka — pole języka po jego stronie jest puste.'
      : `Rdzeń trzyma tekst źródłowy BEZ języka, a w polu wyżej stoi „${pole}". ` +
          'Zapis źródła uzgodni obie wartości.';
  }
  if (pole === '' || pole === rdzen) return `Rdzeń trzyma dziś język źródłowy „${rdzen}".`;
  return (
    `Rdzeń trzyma język źródłowy „${rdzen}", a w polu wyżej stoi „${pole}" — to jeszcze nie ` +
    'jest zapisane. Zapis źródła uzgodni obie wartości.'
  );
}

/** Nazwa stanu panelu widoczna dla Operatora. */
export function nazwaStanuPanelu(stan: TranslationStatus): string {
  if (stan === TranslationStatus.Pending) return 'oczekuje';
  if (stan === TranslationStatus.Translating) return 'tłumaczenie w toku';
  if (stan === TranslationStatus.Ready) return 'gotowe';
  return 'błąd';
}

/**
 * Zdania stanów pustych — jedna forma, różnicowana wyłącznie treścią.
 *
 * Każde z nich tłumaczy, czym okno jest i jak je zapełnić: pustka jest stanem
 * oczekiwanym, więc zdanie prowadzi do pierwszej czynności, zamiast nazywać brak.
 *
 * Tytuł i opis są rozdzielone, bo nośnik stanu pustego ma dwa stopnie pisma:
 * tytuł 13 px półgrubym i opis 12 px.
 */
export const PUSTE: Record<
  'zrodlo' | 'panele' | 'glosariusz' | 'pamiec' | 'formaty' | 'jakosc',
  ZdanieStanu
> = {
  zrodlo: {
    tytul: 'Bez tekstu źródłowego',
    opis:
      'To okno przyjmuje tekst do przełożenia. Wpisz albo wklej treść i zapisz źródło — ' +
      'zapis zakłada podział na segmenty i rozsyła je do wszystkich paneli naraz.',
  },
  panele: {
    tytul: 'Bez języka docelowego',
    opis:
      'Każdy język docelowy to osobny panel przekładu. Dodaj pierwszy formularzem ' +
      '„+ Dodaj język" — panele aktualizują się potem równolegle.',
  },
  glosariusz: {
    tytul: 'Glosariusz pusty w tej sesji',
    opis:
      'Definiuj terminy formularzem wyżej. Kontrakt nie ma komendy odczytu glosariusza ' +
      '(translate.glossary.list ani .get), więc widać tu wyłącznie terminy zapisane stąd.',
  },
  pamiec: {
    tytul: 'Pamięć tłumaczeń nieodpytana',
    opis:
      'To okno pyta pamięć o podpowiedź dla jednego segmentu i jednego panelu. Wskaż panel, ' +
      'wybierz segment źródłowy albo wpisz fragment i uruchom szukanie.',
  },
  formaty: {
    tytul: 'Bez dokumentu i bez panelu',
    opis:
      'To okno wydobywa tekst z dokumentu do Source Panel i wydaje panel w formacie pliku. ' +
      'Wskaż ścieżkę dokumentu albo dodaj pierwszy język docelowy.',
  },
  jakosc: {
    tytul: 'Kontroli jakości jeszcze nie uruchomiono',
    opis:
      'To okno prowadzi kontrolę jakości wszystkich paneli naraz i zbiera ich zastrzeżenia ' +
      'w jednym wykazie. Uruchom kontrolę przyciskiem wyżej.',
  },
};

/**
 * Odmowa, nie pustka — dlatego zdanie stoi poza wykazem stanów pustych.
 *
 * Brak okna sesji nie jest stanem oczekiwanym pierwszego użycia: moduł nie ma
 * wtedy czym zaadresować ani jednego żądania i mówi to jako powód niewykonania
 * czynności (`okno.blad`, wiersz odpowiedzi).
 */
export const BRAK_OKNA =
  'Rdzeń nie zwrócił żadnego okna tej sesji. Bez identyfikatora okna nie ma czym ' +
  'zaadresować translate.source.set ani translate.target.add.';

/** Zdania objaśnień [?] przy elementach konfiguracji modułu. */
export const OBJASNIENIA = {
  jezykZrodlowy:
    'Puste pole zostawia rozpoznanie językowi rdzenia — kontrakt dopuszcza żądanie bez ' +
    'języka źródłowego i wtedy rdzeń rozpoznaje go sam.',
  ponownaSegmentacja:
    'Zaznaczone żąda ponownego podziału tekstu na segmenty przy zapisie źródła. ' +
    'Segmentacja jest po stronie rdzenia; klient jej nie liczy.',
  jezykDocelowy:
    'Kod języka panelu. Kontrakt przyjmuje dowolny napis i nie zwraca wykazu języków, ' +
    'więc lista obok pola jest podpowiedzią, a nie katalogiem.',
  tonPanelu:
    'Ton tłumaczenia obowiązujący ten jeden panel. Pole otwarte — kontrakt nie ' +
    'ogranicza zbioru tonów. W instancji panelu pole niesie ton, który trzyma rdzeń: ' +
    'po zapisie pokazuje to samo, co nagłówek panelu.',
  kanalPrzekladu:
    'Kanał modelu, którym rdzeń wykona przekład. Brak wskazania bierze kanał czynny okna — ' +
    'tak, jak mówi opis pola channelId w kontrakcie. Wykaz obejmuje wyłącznie kanały czynne.',
  formatEksportu:
    'Format pliku wyniku. Wykaz pochodzi z kontraktu (ExportFormat); XLIFF, wymieniony ' +
    'w wykazie okien, nie jest jego wartością.',
  nieTlumacz:
    'Termin oznaczony jako nietłumaczony zostaje w tłumaczeniach w postaci źródłowej — ' +
    'nazwy własne, oznaczenia handlowe, symbole.',
  sciezkaGlosariusza:
    'Ścieżka pliku TBX albo CSV po stronie rdzenia. Klient plików nie czyta ani nie ' +
    'zapisuje — robi to rdzeń pod wskazaną ścieżką.',
} as const;
