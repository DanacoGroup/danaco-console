import type { StanWpisu } from './wpis-rozmowy';

/**
 * Napisy warstwy rozmowy — jedno miejsce dla wszystkiego, co Operator czyta.
 *
 * Napis nigdy nie stoi obok samej barwy: każdy stan niesie słowo, bo stan
 * sygnalizowany wyłącznie kolorem jest niedopuszczalny.
 */
export const NAPISY = {
  /** Nagłówek listy historii. */
  historia: 'Historia okna komunikacji',
  /** Etykieta pola wpisywania. */
  polePodpowiedz: 'Napisz do modelu. Enter wysyła, Shift+Enter dodaje wiersz.',
  /**
   * Podpowiedź stojąca pod polem, gdy otwarty jest wykaz po ukośniku.
   *
   * Zdanie mówi wprost, że pole wpisywania jest zarazem filtrem wykazu, i
   * wymienia klawisze, które przy otwartym wykazie znaczą co innego niż zwykle.
   */
  wykazPodpowiedz: 'Wpisz, by zawęzić wykaz · ↑↓ wybór · Enter zatwierdza · Esc zamyka',
  /** Przycisk wysyłki. */
  wyslij: 'Wyślij',
  /** Przycisk zatrzymania. */
  przerwij: 'Zatrzymaj',
  /** Blok prowenancji. */
  prowenancja: 'Co poszło do modelu',
  /** Blok podglądu pracy modelu. */
  rozumowanie: 'Praca modelu',
  /** Blok wywołań narzędzi. */
  narzedzia: 'Narzędzia',
  /** Blok metadanych konta. */
  konto: 'Konto kanału',
  /** Nagłówek listy błędów tury. */
  bledy: 'Błąd kanału',
  /** Podpis stopki z podsumowaniem tury. */
  podsumowanie: 'Tura zamknięta',
  /** Wykaz typów zdarzeń przechwyconych w turze — tryb „Pełny". */
  zdarzenia: 'Zdarzenia tury',
  /**
   * Wykaz wytworów tury — tryb „Streszczenie".
   *
   * Napis mówi „narzędzia", a nie „pliki": rdzeń nie nadaje spisu wytworzonych
   * plików, więc wykaz wymienia nazwy wywołanych narzędzi.
   */
  wytwory: 'Wytwory tury — narzędzia',
  /** Stan pusty trybu „Streszczenie": tytuł. */
  streszczeniePustoTytul: 'Streszczenie powstaje z domkniętych tur',
  /**
   * Stan pusty trybu „Streszczenie": opis.
   *
   * Tłumaczy, czym streszczenie jest i skąd się bierze, zamiast meldować „brak
   * danych".
   */
  streszczeniePustoOpis:
    'Streszczenie tury składa rdzeń i dokłada je do ostatniego fragmentu strumienia: '
    + 'czas, liczba tur wewnętrznych, koszt i konto, a obok nich nazwy narzędzi, '
    + 'które tura wywołała. W tym wątku żadna tura jeszcze się nie domknęła, więc '
    + 'nie ma z czego go złożyć. Wypowiedzi wątku nie zniknęły — czekają w trybie Zwykły.',
  /**
   * Odpowiedź na kliknięcie „Zatrzymaj", kiedy rdzeń nie miał czego przerwać
   * — rdzeń oddaje wtedy `stopped` równe fałsz, a kliknięcie i tak musi zostawić
   * ślad na ekranie.
   */
  zatrzymanieBezTury:
    'Zatrzymanie: w tym oknie nie biegła żadna tura — rdzeń nie miał czego przerwać.',
  /** Zatrzymanie tury, o której historia okna nic nie wiedziała. */
  zatrzymanieBezWpisu:
    'Zatrzymanie: rdzeń przerwał turę tego okna, choć historia nie wskazywała żadnej w biegu.',
  /** Stan pustej historii. */
  pustoTytul: 'Rozmowa jeszcze się nie zaczęła',
  /** Opis pustej historii. */
  pustoOpis: 'Wypowiedź Operatora otwiera turę kanału głównego.',
} as const;

/** Nazwa stanu wpisu widoczna przy wypowiedzi. */
export function nazwaStanuWpisu(stan: StanWpisu): string {
  switch (stan) {
    case 'wysylanie':
      return 'wysyłanie';
    case 'strumien':
      return 'w strumieniu';
    case 'zakonczony':
      return 'zakończona';
    case 'pusty':
      return 'zamknięta bez odpowiedzi';
    case 'przerwany':
      return 'przerwana';
    case 'bledny':
      return 'błąd';
  }
}

/** Nazwa pola prowenancji w postaci prezentacyjnej. */
export const POLA_PROWENANCJI: ReadonlyArray<readonly [string, string]> = [
  ['model', 'Model'],
  ['modelZapasowy', 'Model zapasowy'],
  ['naklad', 'Nakład rozumowania'],
  ['trybUprawnien', 'Tryb uprawnień'],
  ['katalogRoboczy', 'Katalog roboczy'],
  ['plikUstawien', 'Plik ustawień'],
  ['wznowienie', 'Wznowienie rozmowy'],
  ['konto', 'Konto'],
  ['katalogKonta', 'Katalog konta'],
  ['trybNakladki', 'Tryb nakładki'],
  ['skrotNakladki', 'Skrót nakładki'],
  ['chwila', 'Chwila wywołania'],
];
