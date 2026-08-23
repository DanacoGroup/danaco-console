import type { NazwaIkony } from '../ikony/ikony';
import { nazwaRoli, nazwaSrodowiska } from '../okno-komunikacji/etykiety-okna';
import {
  nazwaModulu,
  nazwaNakladu,
  nazwaTrybuUprawnien,
} from '../sterowanie/etykiety-sterowania';
import { nazwaKanalu, type RejestrKanalow } from '../sterowanie/rejestr-kanalow';
import type { MigawkaSterowania } from '../sterowanie/stan-sterowania';

/**
 * Osiem ustawień okna sprowadzonych do wierszy podsumowania.
 *
 * Plik nie zna ani jednej nazwy prezentacyjnej z osobna — wszystkie bierze
 * z istniejących słowników: `okno-komunikacji/etykiety-okna` (środowisko,
 * rola), `sterowanie/etykiety-sterowania` (moduł, nakład, uprawnienia)
 * i `sterowanie/rejestr-kanalow` (nazwa kanału modelu). Jedyne, co dokłada,
 * to złożenie wartości w wiersz czytelny bez otwierania kontrolki.
 *
 * Host wykonania nie jest osobnym wierszem, lecz dopiskiem przy środowisku
 * wykonania — uszczegóławia zasięg, a nie stanowi odrębnego ustawienia. Nazwa
 * hosta zostaje przez to widoczna bez dziewiątej pozycji.
 */
export interface PozycjaPodsumowania {
  /** Klucz porządkowy wiersza; zarazem wartość `data-ustawienie`. */
  klucz: string;
  /** Etykieta wersalikowa ustawienia. */
  etykieta: string;
  /** Bieżąca wartość zapisana słowami operatora. */
  wartosc: string;
  /** Ikona wiersza; stan nigdy nie jest niesiony samym kolorem. */
  ikona: NazwaIkony;
  /** Objaśnienie kontekstowe ustawienia — treść dymka pomocy. */
  objasnienie: string;
  /** Wiersz czeka na dane z rdzenia — wskaźnik ładowania obok wartości. */
  ladowanie?: boolean;
}

/** Osiem wierszy podsumowania w kolejności, w jakiej stoją w widoku. */
export function pozycjePodsumowania(
  migawka: MigawkaSterowania,
  rejestr: RejestrKanalow,
): PozycjaPodsumowania[] {
  const { okno, ustawienia } = migawka;

  // Pusty wykaz przed odpowiedzią rdzenia: nazwy kanałów są jeszcze w drodze.
  // Rejestr nie odróżnia braku odpowiedzi od pustego rejestru (zgłoszone),
  // więc do tego czasu wiersz kanału niesie wskaźnik ładowania, nie wyrok.
  const wykazWDrodze = rejestr.kanaly().length === 0;

  return [
    {
      klucz: 'srodowisko',
      etykieta: 'Środowisko wykonania',
      wartosc: opisSrodowiska(migawka),
      ikona: 'globus',
      objasnienie:
        'Gdzie wykonują się polecenia okna: lokalnie albo na wskazanym hoście. Dopisek po kropce to host wykonania.',
    },
    {
      klucz: 'katalogi',
      etykieta: 'Katalogi robocze',
      wartosc: opisKatalogow(okno.workingDirs),
      ikona: 'folder',
      objasnienie:
        'Katalogi, w których okno czyta i zapisuje pliki. Brak katalogu nie blokuje pracy okna.',
    },
    {
      klucz: 'model',
      etykieta: 'Model',
      wartosc: opisKanalu(rejestr, okno.modelChannelId, 'Bez wskazania'),
      ikona: 'gwiazdka',
      objasnienie:
        'Kanał modelu prowadzący rozmowę tego okna. Wykaz kanałów pochodzi z rejestru rdzenia.',
      ladowanie: wykazWDrodze && okno.modelChannelId.length > 0,
    },
    {
      klucz: 'model-zapasowy',
      etykieta: 'Model zapasowy',
      wartosc: opisKanalu(rejestr, ustawienia.kanalZapasowy, 'Bez modelu zapasowego'),
      ikona: 'odswiez',
      objasnienie:
        'Kanał przejmujący rozmowę, gdy model główny odmówi albo wyczerpie limit.',
      ladowanie: wykazWDrodze && ustawienia.kanalZapasowy.length > 0,
    },
    {
      klucz: 'naklad',
      etykieta: 'Nakład rozumowania',
      wartosc: nazwaNakladu(ustawienia.nakladRozumowania),
      ikona: 'zegar',
      objasnienie:
        'Ile rozumowania model poświęca na odpowiedź — od zwięzłego do pogłębionego.',
    },
    {
      klucz: 'uprawnienia',
      etykieta: 'Tryb uprawnień',
      wartosc: nazwaTrybuUprawnien(okno.permissionMode),
      ikona: 'tarcza',
      objasnienie:
        'Zakres czynności, które okno wykonuje bez pytania Operatora o zgodę.',
    },
    {
      klucz: 'rola',
      etykieta: 'Rola okna',
      wartosc: nazwaRoli(okno.windowRole),
      ikona: 'uzytkownik',
      objasnienie:
        'Miejsce okna w pętli koordynator–wykonawca: planuje, wykonuje albo pracuje samodzielnie.',
    },
    {
      klucz: 'modul',
      etykieta: 'Moduł',
      wartosc: nazwaModulu(okno.moduleId),
      ikona: 'menu',
      objasnienie: 'Zestaw narzędzi i kontekstu, z którym pracuje to okno.',
    },
  ];
}

/** Zasięg wykonania wraz z nazwą hosta, jeżeli operator ją wskazał. */
function opisSrodowiska(migawka: MigawkaSterowania): string {
  const zasieg = nazwaSrodowiska(migawka.okno.executionEnv);
  const host = migawka.ustawienia.hostWykonania.trim();
  return host.length > 0 ? `${zasieg} · ${host}` : zasieg;
}

/**
 * Liczba katalogów roboczych po polsku.
 *
 * Pusta lista nie jest brakiem gotowości ani błędem — okno pracuje bez
 * wskazanego katalogu, więc wiersz mówi to wprost.
 */
function opisKatalogow(katalogi: string[]): string {
  const ile = katalogi.length;
  if (ile === 0) return 'Brak — okno pracuje bez wskazanego katalogu';
  if (ile === 1) return `1 katalog · ${katalogi[0]}`;
  return `${ile} ${odmianaKatalogu(ile)}`;
}

/** Odmiana rzeczownika „katalog" przy liczebniku większym od jedności. */
function odmianaKatalogu(ile: number): string {
  const dziesiatki = ile % 100;
  const jednosci = ile % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (dziesiatki < 12 || dziesiatki > 14);
  return mnoga ? 'katalogi' : 'katalogów';
}

/**
 * Nazwa kanału modelu z rejestru rdzenia.
 *
 * Identyfikator spoza wykazu nie znika i nie zostaje podmieniony — widok
 * pokazuje prawdę o oknie, a nie pozycję pierwszą z brzegu. Dopóki wykaz
 * jest pusty (odpowiedź rdzenia w drodze), dopisek „spoza wykazu" byłby
 * przedwczesny — wiersz pokazuje sam identyfikator ze wskaźnikiem ładowania.
 */
function opisKanalu(
  rejestr: RejestrKanalow,
  identyfikator: string,
  gdyPusty: string,
): string {
  if (identyfikator.length === 0) return gdyPusty;
  const wykaz = rejestr.kanaly();
  if (wykaz.length === 0) return identyfikator;
  const kanal = wykaz.find((pozycja) => pozycja.id === identyfikator);
  return kanal !== undefined ? nazwaKanalu(kanal) : `${identyfikator} (spoza wykazu)`;
}
