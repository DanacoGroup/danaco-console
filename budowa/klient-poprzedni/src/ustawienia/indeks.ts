import type { Kanal } from '../protokol/kanal';
import { utworzOknoUstawien, type OknoUstawien } from './okno-ustawien';

/** Okno ustawień jest punktem zbiorczym katalogu: żyje między otwarciami, wraca do sekcji, na której operator skończył, i buduje się na nowo po zmianie kanału połączenia. */
let okno: OknoUstawien | null = null;

/** Kanał, na którym zbudowano to konkretne okno — podstawa rozpoznania zmiany połączenia z rdzeniem klienta. */
let osadzonyKanal: Kanal | null = null;

/** Otwiera okno ustawień i zleca odczyt sekcji czynnej; otwarcie nie czeka na rdzeń, okno pojawia się od razu. */
export function otworzOknoUstawien(kanal: Kanal): OknoUstawien {
  if (okno !== null && osadzonyKanal !== kanal) {
    okno.rozlacz();
    okno = null;
  }

  if (okno === null) {
    okno = utworzOknoUstawien(kanal);
    osadzonyKanal = kanal;
  }

  okno.otworz();
  return okno;
}

export type { OknoUstawien } from './okno-ustawien';
export type { SekcjaUstawien, KodSekcjiUstawien, OpisSekcjiUstawien } from './sekcje';
export { REJESTR_SEKCJI_USTAWIEN } from './sekcje';
export type { StanUstawien } from './stan-ustawien';

/** Most motywu wychodzi z katalogu osobno od okna, bo nastawa motywu obowiązuje także wtedy, gdy okno ustawień stoi zamknięte. */
export { podepnijMostMotywu, odepnijMostMotywu, KLUCZ_MOTYWU } from './most-motywu';
export type { MostMotywu, WyborMotywu } from './most-motywu';
