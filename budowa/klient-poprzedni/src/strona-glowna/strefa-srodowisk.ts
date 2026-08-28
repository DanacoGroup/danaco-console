import { utworzKarteSrodowiska, type KartaSrodowiska } from './karta-srodowiska';
import { utworzNaglowekStrefy } from './naglowek-strefy';
import { POZYCJE_SRODOWISK, type KodSrodowiska, type PozycjaSrodowiska } from './pozycje-srodowisk';
import { utworzSygnalWyboru, type SluchaczWyboru } from './sygnal-wyboru';

/** Strefa pierwsza pokazuje siatkę kart środowisk ze stanem czynnym i stanowi środek ciężkości strony głównej: największą powierzchnię i jedyny krój nagłówkowy na ekranie. */
const ETYKIETA = 'Strefa 1 · Karty środowisk';
const WYJASNIENIE = 'Wybór środowiska otwiera przestrzeń roboczą i przywraca jej karty sesji.';

export interface StrefaSrodowisk {
  element: HTMLElement;
  naWybor(sluchacz: SluchaczWyboru<PozycjaSrodowiska>): void;
  /** Przerysowuje strefę wykazem z rdzenia; stała zastana służy tylko jako treść przed odpowiedzią. */
  ustawWykaz(pozycje: readonly PozycjaSrodowiska[]): void;
  /** Oznacza środowisko czynne; `null` zdejmuje oznaczenie ze wszystkich. */
  ustawCzynne(kod: KodSrodowiska | null): void;
}

export function utworzStrefeSrodowisk(): StrefaSrodowisk {
  const sygnal = utworzSygnalWyboru<PozycjaSrodowiska>();

  const element = document.createElement('section');
  element.className = 'dn-strona__strefa dn-strona__strefa--srodowiska';
  element.setAttribute('aria-label', ETYKIETA);

  const siatka = document.createElement('div');
  siatka.className = 'dn-strona__siatka dn-strona__siatka--srodowiska';

  let karty: KartaSrodowiska[] = [];
  let czynne: string | null = null;

  /** Rysuje karty od nowa z podanego wykazu, zachowując oznaczenie czynnego. */
  function narysuj(pozycje: readonly PozycjaSrodowiska[]): void {
    karty = pozycje.map((pozycja) =>
      utworzKarteSrodowiska(pozycja, (wybrana) => sygnal.nadaj(wybrana)),
    );
    siatka.replaceChildren(...karty.map((karta) => karta.element));
    for (const karta of karty) karta.oznaczCzynna(karta.kod === czynne);
  }

  narysuj(POZYCJE_SRODOWISK);

  element.append(utworzNaglowekStrefy(ETYKIETA, WYJASNIENIE), siatka);

  return {
    element,
    naWybor: sygnal.sluchaj,
    ustawWykaz: narysuj,

    ustawCzynne(kod) {
      czynne = kod;
      for (const karta of karty) {
        karta.oznaczCzynna(karta.kod === kod);
      }
    },
  };
}
