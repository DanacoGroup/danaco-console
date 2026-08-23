import { Command, type Channel } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';

/**
 * Wykaz kanałów modelu czytany z rejestru rdzenia.
 *
 * Rejestr jest sterowany danymi: nowy kanał to nowy wiersz, nie nowy typ
 * w kodzie. Sterowania modelu głównego i zapasowego czytają wykaz stąd,
 * zamiast prowadzić własną listę nazw.
 *
 * Wykaz jest katalogiem wyboru wspólnym dla całego klienta, nie ustawieniem
 * okna — jeden egzemplarz obsługuje dowolną liczbę okien i żadne z nich nie
 * zapisuje w nim swojego stanu. Pusty wykaz nie wyłącza sterowania: pole
 * pokazuje wartość bieżącą okna i przyjmuje wpis operatora.
 */
export interface RejestrKanalow {
  /** Kanały znane w tej chwili; pusta lista, dopóki rdzeń nie odpowie. */
  kanaly(): Channel[];
  /**
   * Czy rdzeń odpowiedział na `channel.list` choć raz.
   *
   * Pusty wykaz znaczy dwie różne rzeczy — „jeszcze nie wiem" i „rejestr jest
   * pusty" — a widok musi je rozróżnić: pierwsza to wskaźnik ładowania, druga
   * to stan pusty ze zdaniem o pustym rejestrze. Bez tej flagi wskaźnik
   * ładowania nie zgasłby nigdy na rdzeniu z autentycznie pustym rejestrem.
   */
  odpowiedzOtrzymana(): boolean;
  /** Zamawia wykaz z rdzenia. */
  odswiez(): void;
  /** Subskrypcja zmian wykazu. */
  naZmiane(sluchacz: (kanaly: Channel[]) => void): Odsubskrybuj;
}

export function utworzRejestrKanalow(kanal: Kanal): RejestrKanalow {
  const zmiany = utworzMagistrale<Channel[]>();
  let wykaz: Channel[] = [];
  let odpowiedziano = false;

  return {
    kanaly: () => wykaz,

    odpowiedzOtrzymana: () => odpowiedziano,

    odswiez() {
      kanal.wyslij(Command.ChannelList, {}, (wynik) => {
        // Odpowiedź odnotowujemy także wtedy, gdy rdzeń odmówił albo nie podał
        // wykazu: pytanie zostało rozstrzygnięte, więc wskaźnik ładowania ma
        // zgasnąć, a widok pokazać stan pusty zamiast wiecznego czekania.
        odpowiedziano = true;
        const odebrane = wynik.wynik?.channels;
        wykaz = wynik.udany && odebrane !== undefined ? odebrane : wykaz;
        zmiany.oglos(wykaz);
      });
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/** Nazwa kanału pokazywana operatorowi: nazwa własna, rodzaj i model. */
export function nazwaKanalu(kanal: Channel): string {
  const model = kanal.model !== undefined && kanal.model.length > 0 ? ` · ${kanal.model}` : '';
  const czynny = kanal.enabled ? '' : ' · nieczynny';
  return `${kanal.name} (${kanal.kind})${model}${czynny}`;
}
