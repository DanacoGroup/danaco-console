import { Command, type Channel } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';

/** Wykaz kanałów modelu czytany z rejestru rdzenia, wspólny dla całego klienta i niezależny od pojedynczego okna. */
export interface RejestrKanalow {
  /** Kanały znane w tej chwili; pusta lista, dopóki rdzeń nie odpowie. */
  kanaly(): Channel[];
  // Czy rdzeń odpowiedział na wykaz kanałów choć raz; pusty wykaz bywa dwuznaczny.
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
        // Odpowiedź odnotowujemy także po odmowie, żeby wskaźnik ładowania zgasł.
        odpowiedziano = true;
        const odebrane = wynik.wynik?.channels;
        wykaz = wynik.udany && odebrane !== undefined ? odebrane : wykaz;
        zmiany.oglos(wykaz);
      });
    },

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
  };
}

/** Nazwa kanału pokazywana operatorowi: nazwa własna kanału, jego rodzaj oraz model, jeśli jest już znany. */
export function nazwaKanalu(kanal: Channel): string {
  const model = kanal.model !== undefined && kanal.model.length > 0 ? ` · ${kanal.model}` : '';
  const czynny = kanal.enabled ? '' : ' · nieczynny';
  return `${kanal.name} (${kanal.kind})${model}${czynny}`;
}
