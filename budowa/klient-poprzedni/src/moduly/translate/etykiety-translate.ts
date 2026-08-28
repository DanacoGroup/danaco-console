import { ExportFormat, TranslationStatus } from '../../../../shared/contract';
import type { ZdanieStanu } from './stan-okna-translate';

/**
 * Teksty i katalogi widoczne dla operatora w module Translate.
 */

/**
 * Podpowiedź języków docelowych tłumaczenia; pole pozostaje otwarte na kod języka spoza wykazu kontraktu.
 */
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

/**
 * Podpowiedź tonu tłumaczenia w panelu Translate; kontrakt przyjmuje w tym polu dowolny napis operatora.
 */
export const PODPOWIEDZ_TONU: readonly string[] = [
  'neutralny',
  'formalny',
  'nieformalny',
  'perswazyjny',
  'empatyczny',
  'techniczny',
];

/**
 * Formaty eksportu panelu Translate tworzą katalog zamknięty, pochodzący wprost ze słownika kontraktu.
 */
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
 * Ton panelu widoczny w nagłówku instancji jest jedną prawdą o nastawie tonu, oddawaną wprost przez rdzeń.
 */
export function opisTonuPanelu(ton: string | undefined): string {
  const nazwa = (ton ?? '').trim();
  return nazwa === '' ? 'ton domyślny rdzenia' : `ton: ${nazwa}`;
}

/**
 * Zdanie o języku źródłowym trzymanym przez rdzeń pokazuje rozbieżność ze stanem lokalnym bez cichego uzgadniania.
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

/**
 * Nazwa stanu panelu tłumaczenia widoczna dla operatora w nagłówku bieżącej instancji panelu Translate.
 */
export function nazwaStanuPanelu(stan: TranslationStatus): string {
  if (stan === TranslationStatus.Pending) return 'oczekuje';
  if (stan === TranslationStatus.Translating) return 'tłumaczenie w toku';
  if (stan === TranslationStatus.Ready) return 'gotowe';
  return 'błąd';
}

/**
 * Zdania stanów pustych panelu Translate mają jedną formę, różnicowaną wyłącznie treścią komunikatu operatora.
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
 * Odmowa braku okna sesji stoi poza wykazem stanów pustych, bo nie jest stanem oczekiwanym pierwszego użycia.
 */
export const BRAK_OKNA =
  'Rdzeń nie zwrócił żadnego okna tej sesji. Bez identyfikatora okna nie ma czym ' +
  'zaadresować translate.source.set ani translate.target.add.';

/**
 * Zdania objaśnień widoczne w dymkach pomocy przy elementach konfiguracji panelu modułu Translate tej budowy.
 */
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
