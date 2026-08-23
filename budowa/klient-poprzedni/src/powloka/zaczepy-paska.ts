import type { PozycjaUstawienia } from '../strona-glowna/pozycje-ustawien';

/**
 * Zaczepy paska — czynności, których pasek nie zna, podane mu z zewnątrz.
 *
 * Kształt zaczepów jest potrzebny naraz czterem miejscom powłoki
 * (`pasek-gorny.ts`, `akcje-paska.ts`, `menu-profilu.ts`,
 * `wyszukiwanie-globalne.ts`), a te importują się nawzajem. Wspólny kształt
 * w pliku bez zależności jest jedyną postacią, która nie zawiązuje cyklu
 * importów.
 *
 * Okna Always On Display, Mobile, Modeli czy Konfiguracji potrzebują kanału do
 * rdzenia, którego powłoka nie ma. Pasek dostaje więc czynność, nie drogę.
 *
 * Zaczep czyta się w chwili naciśnięcia, nie przy montażu: powłoka powstaje bez
 * połączenia z rdzeniem, więc funkcja zapamiętana przy montażu byłaby pusta,
 * a zaczep podłączony później nie doszedłby do kontrolek. Stąd wszystkie
 * odczyty `zaczepy.otworzUstawienie` stoją w obsłudze naciśnięcia.
 *
 * Brak zaczepu nie jest brakiem okna. Kontrolka bez podanej czynności zostaje
 * klikalna i mówi, że ten pasek nie dostał drogi do okna — nie że okna nie ma,
 * bo okna są zbudowane i osiągalne z listwy strony głównej.
 */
export interface ZaczepyPaska {
  /**
   * Otwiera okno platformy odpowiadające pozycji ustawień.
   *
   * Jedna czynność na wszystkie siedem pozycji, bo po drugiej stronie i tak
   * stoi jedno rozgałęzienie (`aplikacja/akcje-ustawien.ts`). Siedem osobnych
   * zaczepów byłoby siedmioma drogami do tego samego rozgałęzienia.
   */
  otworzUstawienie?(pozycja: PozycjaUstawienia): void;

  /**
   * Otwiera albo zamyka centrum powiadomień.
   *
   * Osobny zaczep, a nie pozycja ustawień: centrum nie jest oknem platformy
   * z listwy strony głównej, tylko kolumną powłoki środowiska, a jego
   * wyzwalaczem jest plakietka dzwonka — element warstwy 1, jedyna spoczynkowa
   * reprezentacja mechanizmu (katalog komponentów, rozdz. 11.6).
   */
  otworzPowiadomienia?(): void;
}

/** Gdzie postawić zdanie; bez odbiorcy zdanie przepada świadomie. */
export type OdbiorcaZdania = (tytul: string, tresc: string) => void;

/**
 * Zdanie o braku drogi — jedno na cały pasek.
 *
 * Wypowiadają je trzy kontrolki (przycisk obecności, pozycja menu profilu,
 * wynik wyszukiwania). Trzy kopie rozjechałyby się przy pierwszej korekcie,
 * a to zdanie odróżnia brak okna od braku drogi do okna w tym pasku.
 */
export function zdanieBezDrogi(nazwaOkna: string): string {
  return (
    `Ten pasek nie dostał drogi do okna „${nazwaOkna}". Okno jest zbudowane `
    + 'i osiągalne z listwy strony głównej — brakuje wyłącznie zaczepu '
    + 'w widoku środowiska, nie samego okna.'
  );
}

/**
 * Wykonuje zaczep albo mówi, dlaczego go nie ma. Odczyt zaczepu jest tutaj —
 * czyli w chwili naciśnięcia, nie w chwili montażu kontrolki.
 */
export function wykonajZaczep(
  zaczepy: ZaczepyPaska,
  pozycja: PozycjaUstawienia,
  naKomunikat?: OdbiorcaZdania,
): void {
  const otworz = zaczepy.otworzUstawienie;
  if (otworz === undefined) {
    naKomunikat?.(pozycja.nazwa, zdanieBezDrogi(pozycja.nazwa));
    return;
  }
  otworz(pozycja);
}
