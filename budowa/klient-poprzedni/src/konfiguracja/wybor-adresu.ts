import {
  ConfigAxis,
  ConfigScope,
  type SettingDefinition,
} from '../../../shared/contract';
import type { AdresUstawienia } from './adres-ustawienia';
import {
  BYTY_OSI,
  BYTY_ZASIEGOW,
  NAZWY_OSI,
  NAZWY_ZASIEGOW,
  OSIE_OD_NAJWEZSZEJ,
  ZASIEGI_OD_NAJWEZSZEGO,
  osWymagaBytu,
  zasiegWymagaBytu,
} from './zasiegi';

/**
 * Wybór adresu w przestrzeni konfiguracji: poziom zasięgu, byt poziomu, oś,
 * byt osi.
 *
 * Ta sama kontrolka obsługuje dwie role okna: pasek u góry ustawia punkt
 * widzenia, względem którego liczone jest dziedziczenie pól, a panel przy polu
 * ustawia adres zapisu wartości nadpisującej. Różni je wyłącznie wykaz
 * dopuszczalnych poziomów i osi — pasek podaje wszystkie, panel pola tylko te
 * z katalogu (`allowedScopes`, `allowedAxes`).
 *
 * Pola bytu znikają tam, gdzie poziom lub oś bytu nie mają (poziom globalny,
 * oś platformy).
 */
export interface WyborAdresu {
  /** Wiersz kontrolek osadzany w pasku albo w panelu pola. */
  element: HTMLElement;
  /** Adres wskazany kontrolkami. */
  adres(): AdresUstawienia;
  /** Nanosi adres na kontrolki; człony spoza wykazu zostają pominięte. */
  ustaw(adres: AdresUstawienia): void;
  /** Subskrypcja zmiany adresu. */
  naZmiane(sluchacz: () => void): void;
}

export function utworzWyborAdresu(
  poziomy: readonly ConfigScope[],
  osie: readonly ConfigAxis[],
): WyborAdresu {
  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  const zasieg = listaWyboru(
    poziomy.map((poziom) => ({ wartosc: poziom, etykieta: NAZWY_ZASIEGOW[poziom] })),
    'Poziom zasięgu',
  );
  const os = listaWyboru(
    osie.map((wartosc) => ({ wartosc, etykieta: NAZWY_OSI[wartosc] })),
    'Oś rozstrzygania',
  );
  const bytZasiegu = poleBytu('Byt poziomu');
  const bytOsi = poleBytu('Byt osi');

  const koszykBytuZasiegu = opisany('Byt poziomu', bytZasiegu);
  const koszykBytuOsi = opisany('Byt osi', bytOsi);

  function ubierz(): void {
    const wybranyZasieg = zasieg.value as ConfigScope;
    const wybranaOs = os.value as ConfigAxis;
    koszykBytuZasiegu.hidden = !zasiegWymagaBytu(wybranyZasieg);
    bytZasiegu.placeholder = BYTY_ZASIEGOW[wybranyZasieg] ?? '';
    koszykBytuOsi.hidden = !osWymagaBytu(wybranaOs);
    bytOsi.placeholder = BYTY_OSI[wybranaOs] ?? '';
  }

  for (const kontrolka of [zasieg, os, bytZasiegu, bytOsi]) {
    kontrolka.addEventListener('change', () => {
      ubierz();
      oglos();
    });
  }

  const element = document.createElement('div');
  element.className = 'dk-adres';
  element.append(
    opisany('Poziom', zasieg),
    koszykBytuZasiegu,
    opisany('Oś', os),
    koszykBytuOsi,
  );

  ubierz();

  return {
    element,

    adres: () => ({
      zasieg: zasieg.value as ConfigScope,
      bytZasiegu: koszykBytuZasiegu.hidden ? '' : bytZasiegu.value.trim(),
      os: os.value as ConfigAxis,
      bytOsi: koszykBytuOsi.hidden ? '' : bytOsi.value.trim(),
    }),

    ustaw(adres) {
      if (poziomy.includes(adres.zasieg)) zasieg.value = adres.zasieg;
      if (osie.includes(adres.os)) os.value = adres.os;
      bytZasiegu.value = adres.bytZasiegu;
      bytOsi.value = adres.bytOsi;
      ubierz();
    },

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),
  };
}

/**
 * Poziomy dopuszczone przez katalog. Lista pusta nie ma w kontrakcie znaczenia
 * własnego, więc brak wskazania otwiera wszystkie poziomy, a rozstrzygnięcie
 * o dopuszczalności zapisu zostaje przy rdzeniu.
 */
export function dopuszczalneZasiegi(definicja: SettingDefinition): ConfigScope[] {
  const wskazane = definicja.allowedScopes ?? [];
  const wybrane = ZASIEGI_OD_NAJWEZSZEGO.filter((poziom) => wskazane.includes(poziom));
  return wybrane.length > 0 ? wybrane : [...ZASIEGI_OD_NAJWEZSZEGO];
}

/**
 * Osie dopuszczone przez katalog. Lista pusta ma tu w kontrakcie znaczenie
 * własne — „wyłącznie platforma" — więc pozycja bez wskazanych osi dostaje
 * jedną pozycję wyboru zamiast wszystkich osi.
 */
export function dopuszczalneOsie(definicja: SettingDefinition): ConfigAxis[] {
  const wskazane = definicja.allowedAxes ?? [];
  const wybrane = OSIE_OD_NAJWEZSZEJ.filter((os) => wskazane.includes(os));
  return wybrane.length > 0 ? wybrane : [ConfigAxis.Platform];
}

interface PozycjaWyboru {
  wartosc: string;
  etykieta: string;
}

function listaWyboru(pozycje: readonly PozycjaWyboru[], opis: string): HTMLSelectElement {
  const wybor = document.createElement('select');
  wybor.className = 'dn-pole-kontrolka dk-adres__lista';
  wybor.setAttribute('aria-label', opis);
  for (const pozycja of pozycje) {
    const element = document.createElement('option');
    element.value = pozycja.wartosc;
    element.textContent = pozycja.etykieta;
    wybor.append(element);
  }
  return wybor;
}

function poleBytu(opis: string): HTMLInputElement {
  const pole = document.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole-kontrolka dk-adres__byt';
  pole.setAttribute('aria-label', opis);
  return pole;
}

/** Kontrolka wraz z etykietą nad nią. */
function opisany(etykieta: string, kontrolka: HTMLElement): HTMLElement {
  const koszyk = document.createElement('label');
  koszyk.className = 'dn-pole dk-adres__pole';

  const nazwa = document.createElement('span');
  nazwa.className = 'dn-pole-etykieta';
  nazwa.textContent = etykieta;

  koszyk.append(nazwa, kontrolka);
  return koszyk;
}
