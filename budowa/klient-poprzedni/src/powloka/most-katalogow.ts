import { invoke } from '@tauri-apps/api/core';
import { listen } from '@tauri-apps/api/event';

import { czyPowlokaNatywna } from './powloka-natywna';

/**
 * Most do natywnego okna wyboru katalogu — jedyny konsument polecenia powłoki
 * `wybierz_katalog_roboczy` (`desktop/src-tauri/src/dialog_katalogu.rs`).
 *
 * Okno systemu operacyjnego zwraca ścieżkę istniejącą i rozwiniętą, czego pole
 * tekstowe nie zapewnia. Katalog wskazuje się w dwóch sprawach — katalog roboczy
 * modelu (`Cel.KatalogRoboczy`, ustawienie `katalog.roboczy.podstawa`) i zakres
 * wglądu modelu (`Cel.PunktDostepu`, punkt dostępu rodzaju `localDirectory`) —
 * ale czynność systemu jest w obu ta sama, więc okno jest jedno, a cel zmienia
 * wyłącznie napis w belce.
 *
 * Poza powłoką natywną wskazanie wraca wartością `null` tak samo jak rezygnacja
 * z wyboru, a widok zostawia drogę wpisania ścieżki ręcznie. Żadna ścieżka
 * wykonania nie rzuca wyjątkiem i nie odrzuca obietnicy.
 */

/** Nazwa polecenia powłoki; odpowiednik `polecenia::wybierz_katalog_roboczy`. */
const POLECENIE_WYBORU = 'wybierz_katalog_roboczy';

/** Zdarzenie powłoki niosące wybór z zasobnika; odpowiednik `ZDARZENIE_KATALOG`. */
const ZDARZENIE_WYBORU = 'powloka:katalog-roboczy';

/** Sprawa, w której Operator wskazuje katalog. Zmienia napis, nie okno. */
export const Cel = {
  /** Gdzie powstają katalogi sesyjne i pliki robocze modelu. */
  KatalogRoboczy: 'katalog-roboczy',
  /** Do jakiego katalogu model dostaje wgląd. */
  PunktDostepu: 'punkt-dostepu',
} as const;

/** Jedna z wartości `Cel`. */
export type Cel = (typeof Cel)[keyof typeof Cel];

/**
 * Napisy w belce okna. Rejestr sterowany danymi: nowy cel to nowy wiersz
 * tutaj, nie nowa gałąź w kodzie ani nowe okno.
 */
const TYTULY: Record<Cel, string> = {
  [Cel.KatalogRoboczy]: 'Wskaż katalog roboczy',
  [Cel.PunktDostepu]: 'Wskaż katalog, do którego model ma mieć dostęp',
};

/**
 * Czy interfejs działa w powłoce natywnej i ma dostęp do okna systemowego.
 *
 * Rozpoznanie stoi w `powloka-natywna.ts`, bo pyta o nie także most rdzenia.
 * Wyjście powtórzone tutaj oszczędza wywołującym (`dostepy/dialog-katalogu.ts`)
 * drugiego importu.
 */
export { czyPowlokaNatywna };

/**
 * Otwiera natywne okno wyboru katalogu i zwraca wskazaną ścieżkę.
 *
 * Zwraca `null` w trzech przypadkach: przy rezygnacji Operatora, przy braku
 * powłoki natywnej oraz przy niepowodzeniu polecenia. Rozróżnienie ich nie jest
 * tu potrzebne, bo brak wyboru niczego nie nadpisuje. Wywołujący sprawdza
 * `czyPowlokaNatywna`, jeżeli chce powiedzieć Operatorowi, dlaczego okna nie było.
 */
export async function wskazKatalog(cel: Cel = Cel.KatalogRoboczy): Promise<string | null> {
  if (!czyPowlokaNatywna()) return null;
  try {
    const wybor = await invoke<string | null>(POLECENIE_WYBORU, { tytul: TYTULY[cel] });
    return typeof wybor === 'string' && wybor !== '' ? wybor : null;
  } catch (blad) {
    console.warn('[powłoka] okno wyboru katalogu nie otworzyło się', blad);
    return null;
  }
}

/**
 * Podejmuje wybór dokonany z menu zasobnika powłoki.
 *
 * Zwraca odsubskrybowanie działające natychmiast — także wtedy, gdy powłoka
 * nie zdążyła jeszcze potwierdzić nasłuchu: znacznik przerwania pilnuje, żeby
 * spóźniona odpowiedź nie ożywiła zdjętej subskrypcji.
 */
export function naWskazanieZZasobnika(sluchacz: (sciezka: string) => void): () => void {
  if (!czyPowlokaNatywna()) return () => undefined;

  let czynna = true;
  let zdejmij: (() => void) | null = null;

  void listen<string>(ZDARZENIE_WYBORU, (zdarzenie) => {
    if (czynna && typeof zdarzenie.payload === 'string' && zdarzenie.payload !== '') {
      sluchacz(zdarzenie.payload);
    }
  })
    .then((odsubskrybuj) => {
      if (czynna) zdejmij = odsubskrybuj;
      else odsubskrybuj();
    })
    .catch((blad: unknown) => {
      console.warn('[powłoka] nasłuch wyboru katalogu z zasobnika nie ruszył', blad);
    });

  return () => {
    czynna = false;
    if (zdejmij !== null) zdejmij();
    zdejmij = null;
  };
}
