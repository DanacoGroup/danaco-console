import type { StanWpisu } from './wpis-rozmowy';

/** Stała zawiera wszystkie napisy tekstowe warstwy rozmowy, czytane przez operatora w oknie komunikacji z modelem. */
export const NAPISY = {
  /** Nagłówek listy historii. */
  historia: 'Historia okna komunikacji',
  /** Etykieta pola wpisywania. */
  polePodpowiedz: 'Napisz do modelu. Enter wysyła, Shift+Enter dodaje wiersz.',
  /** Podpowiedź pod polem wpisywania, widoczna gdy otwarty jest wykaz pozycji po ukośniku. */
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
  /** Wykaz wytworów tury w trybie Streszczenie, nazwanych wedle wywołanych narzędzi. */
  wytwory: 'Wytwory tury — narzędzia',
  /** Stan pusty trybu „Streszczenie": tytuł. */
  streszczeniePustoTytul: 'Streszczenie powstaje z domkniętych tur',
  /** Opis stanu pustego trybu Streszczenie, wyświetlany zanim żadna tura się domknie. */
  streszczeniePustoOpis:
    'Streszczenie tury składa rdzeń i dokłada je do ostatniego fragmentu strumienia: '
    + 'czas, liczba tur wewnętrznych, koszt i konto, a obok nich nazwy narzędzi, '
    + 'które tura wywołała. W tym wątku żadna tura jeszcze się nie domknęła, więc '
    + 'nie ma z czego go złożyć. Wypowiedzi wątku nie zniknęły — czekają w trybie Zwykły.',
  /** Odpowiedź na Zatrzymaj, gdy rdzeń nie miał żadnej tury do przerwania. */
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

/** Funkcja zwraca nazwę stanu wpisu w postaci prezentowanej operatorowi obok treści wypowiedzi w oknie rozmowy. */
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

/** Stała zawiera pary nazwy pola prowenancji oraz jego etykiety prezentacyjnej wyświetlanej w bloku prowenancji. */
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
