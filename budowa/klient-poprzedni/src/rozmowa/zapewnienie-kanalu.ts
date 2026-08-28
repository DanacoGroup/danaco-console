import { Command, KnownChannelKinds } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';

/**
 * Rodzaj kanału głównego — pierwsza pozycja katalogu kontraktu, wykorzystywana przy
 * zakładaniu rejestru kanałów.
 */
const RODZAJ_GLOWNY = KnownChannelKinds[0];

/**
 * Nazwa wiersza zakładanego w rejestrze kanałów przy braku kanału głównego, widoczna dla
 * Operatora w interfejsie.
 */
const NAZWA_GLOWNEGO = 'Kanał główny — Claude Code CLI';

/**
 * Wynik zapewnienia kanału głównego: identyfikator założonego wiersza rejestru albo opis
 * napotkanej przeszkody.
 */
export interface WynikZapewnienia {
  /** Identyfikator kanału z rejestru; pusty, gdy zapewnienie się nie udało. */
  idKanalu: string;
  /** Skąd wziął się kanał — do komunikatu dla Operatora. */
  pochodzenie: 'zastany' | 'zalozony' | 'brak';
  /** Opis przeszkody, gdy `pochodzenie` równa się `brak`. */
  przeszkoda: string;
}

/**
 * Zapewnia w rejestrze kanałów wiersz kanału głównego, dokładając go komendą channel.add,
 * gdy rejestr go nie ma.
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
/**
 * Dokłada nowy wiersz kanału głównego do rejestru kanałów, dokładnie tak samo, jak
 * zrobiłby to sam Operator ręcznie w panelu.
 */
/**
 * Dokłada nowy wiersz kanału głównego do rejestru kanałów prowadzonego przez rdzeń tej
 * aplikacji.
 */
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

/**
 * Opis przeszkody przeznaczony dla Operatora; pusty komunikat zwrócony przez rdzeń
 * zastępuje zdanie zastępcze.
 */
function opis(komunikat: string | undefined): string {
  return komunikat === undefined || komunikat.length === 0
    ? 'rdzeń nie podał przyczyny'
    : komunikat;
}
