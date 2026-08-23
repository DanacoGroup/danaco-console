import { Command, KnownChannelKinds } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';

/** Rodzaj kanału głównego — pierwsza pozycja katalogu kontraktu. */
const RODZAJ_GLOWNY = KnownChannelKinds[0];

/** Nazwa wiersza zakładanego przy braku kanału głównego w rejestrze. */
const NAZWA_GLOWNEGO = 'Kanał główny — Claude Code CLI';

/** Wynik zapewnienia kanału: identyfikator wiersza albo opis przeszkody. */
export interface WynikZapewnienia {
  /** Identyfikator kanału z rejestru; pusty, gdy zapewnienie się nie udało. */
  idKanalu: string;
  /** Skąd wziął się kanał — do komunikatu dla Operatora. */
  pochodzenie: 'zastany' | 'zalozony' | 'brak';
  /** Opis przeszkody, gdy `pochodzenie` równa się `brak`. */
  przeszkoda: string;
}

/**
 * Zapewnia w rejestrze kanałów wiersz kanału głównego.
 *
 * Rejestr kanałów jest sterowany danymi: kanał istnieje wtedy, gdy istnieje
 * jego wiersz, a nie wtedy, gdy typ dopisano w kodzie. Świeża baza rdzenia nie
 * ma ani jednego wiersza, więc okno komunikacji nie miałoby czym rozmawiać —
 * warstwa dokłada wiersz komendą `channel.add`, tak samo jak zrobiłby to
 * Operator w panelu sterowania okna.
 *
 * Niepowodzenie nie przerywa niczego poza tym wywołaniem: wynik niesie opis
 * przeszkody, a okno pozostaje czynne.
 */
export function zapewnijKanalGlowny(
  kanal: Kanal,
  gotowe: (wynik: WynikZapewnienia) => void,
): void {
  kanal.wyslij(Command.ChannelList, { enabledOnly: true }, (odpowiedz) => {
    if (!odpowiedz.udany) {
      gotowe({ idKanalu: '', pochodzenie: 'brak', przeszkoda: opis(odpowiedz.blad?.message) });
      return;
    }
    const zastany = (odpowiedz.wynik?.channels ?? []).find(
      (pozycja) => pozycja.kind === RODZAJ_GLOWNY && pozycja.enabled,
    );
    if (zastany !== undefined) {
      gotowe({ idKanalu: zastany.id, pochodzenie: 'zastany', przeszkoda: '' });
      return;
    }
    zaloz(kanal, gotowe);
  });
}

/** Dokłada wiersz kanału głównego do rejestru. */
function zaloz(kanal: Kanal, gotowe: (wynik: WynikZapewnienia) => void): void {
  kanal.wyslij(
    Command.ChannelAdd,
    { name: NAZWA_GLOWNEGO, kind: RODZAJ_GLOWNY, enabled: true },
    (odpowiedz) => {
      const dodany = odpowiedz.wynik?.channel;
      if (!odpowiedz.udany || dodany === undefined) {
        gotowe({ idKanalu: '', pochodzenie: 'brak', przeszkoda: opis(odpowiedz.blad?.message) });
        return;
      }
      gotowe({ idKanalu: dodany.id, pochodzenie: 'zalozony', przeszkoda: '' });
    },
  );
}

/** Opis przeszkody dla Operatora; pusty komunikat rdzenia zastępuje zdanie zastępcze. */
function opis(komunikat: string | undefined): string {
  return komunikat === undefined || komunikat.length === 0
    ? 'rdzeń nie podał przyczyny'
    : komunikat;
}
