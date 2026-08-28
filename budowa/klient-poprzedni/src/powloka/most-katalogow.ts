import { invoke } from '@tauri-apps/api/core';
import { listen } from '@tauri-apps/api/event';

import { czyPowlokaNatywna } from './powloka-natywna';

/** Nazwa polecenia powłoki wywoływanego przy wskazywaniu katalogu; odpowiednik polecenia wybierz_katalog_roboczy w powłoce natywnej. */
const POLECENIE_WYBORU = 'wybierz_katalog_roboczy';

/** Zdarzenie powłoki niosące wybór katalogu dokonany z menu zasobnika systemowego; odpowiednik zdarzenia katalog roboczy. */
const ZDARZENIE_WYBORU = 'powloka:katalog-roboczy';

/** Sprawa, w której Operator wskazuje katalog za pomocą tego samego okna systemowego. Zmienia wyłącznie napis w belce, nie samo okno. */
export const Cel = {
  /** Gdzie powstają katalogi sesyjne i pliki robocze modelu. */
  KatalogRoboczy: 'katalog-roboczy',
  /** Do jakiego katalogu model dostaje wgląd. */
  PunktDostepu: 'punkt-dostepu',
} as const;

/** Jedna z wartości sprawy wskazania katalogu: katalog roboczy modelu albo punkt dostępu do katalogu wglądu. */
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
 * Otwiera natywne okno wyboru katalogu i zwraca wskazaną ścieżkę. Zwraca null przy rezygnacji
 * Operatora, przy braku powłoki natywnej oraz przy niepowodzeniu polecenia.
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

/** Podejmuje wybór dokonany z menu zasobnika powłoki i zwraca odsubskrybowanie działające natychmiast, także przed potwierdzeniem nasłuchu. */
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
