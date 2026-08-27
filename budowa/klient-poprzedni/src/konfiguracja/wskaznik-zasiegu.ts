import type { ConfigEntry, SettingDefinition } from '../../../shared/contract';
import { adresWpisu, opisAdresu, type AdresUstawienia } from './adres-ustawienia';
import { utworzLancuchZasiegow } from './lancuch-zasiegow';
import type { PunktWidzenia, Rozstrzygniecie } from './rozstrzygniecie';
import { dopuszczalneOsie, dopuszczalneZasiegi, utworzWyborAdresu } from './wybor-adresu';

/**
 * Wskaźnik zasięgu przy jednym ustawieniu. Plakietka mówi, skąd pochodzi wartość,
 * a rozwinięcie pokazuje pełne dziedziczenie oraz wybór poziomu, na którym wartość
 * ma zostać nadpisana.
 */
export interface WskaznikZasiegu {
  /** Plakietka pochodzenia wraz z rozwijanym panelem. */
  element: HTMLElement;
  /** Adres, pod który pójdzie najbliższy zapis tego pola. */
  adres(): AdresUstawienia;
  /** Nanosi rozstrzygnięcie oraz — przy pierwszym wywołaniu — punkt widzenia. */
  odswiez(rozstrzygniecie: Rozstrzygniecie, punkt: PunktWidzenia): void;
  /** Subskrypcja zmiany adresu zapisu. */
  naZmianeCelu(sluchacz: () => void): void;
}

/**
 * Zależności wskaźnika: pozycja katalogu ustawień oraz sposób zdjęcia zapisu. Zdjęcie
 * idzie komendą `config.reset`, a wskaźnik sam do rdzenia nie sięga, więc czynność
 * przychodzi z zewnątrz jako wywołanie zwrotne.
 */
export interface ZaleznosciWskaznika {
  definicja: SettingDefinition;
  /** Zdejmuje wskazany zapis komendą `config.reset`. */
  naPrzywrocenie(wpis: ConfigEntry): void;
}

export function utworzWskaznikZasiegu(zaleznosci: ZaleznosciWskaznika): WskaznikZasiegu {
  const lancuch = utworzLancuchZasiegow(zaleznosci.naPrzywrocenie);
  const cel = utworzWyborAdresu(
    dopuszczalneZasiegi(zaleznosci.definicja),
    dopuszczalneOsie(zaleznosci.definicja),
  );

  const plakietka = document.createElement('button');
  plakietka.type = 'button';
  plakietka.className = 'dn-plakietka dk-wskaznik__plakietka';
  plakietka.setAttribute('aria-expanded', 'false');

  const panel = document.createElement('div');
  panel.className = 'dk-wskaznik__panel';
  panel.hidden = true;
  panel.append(naglowekDziedziczenia(), lancuch.element, naglowekCelu(), cel.element);

  plakietka.addEventListener('click', () => {
    panel.hidden = !panel.hidden;
    plakietka.setAttribute('aria-expanded', String(!panel.hidden));
  });

  const element = document.createElement('div');
  element.className = 'dk-wskaznik';
  element.append(plakietka, panel);

  /** Punkt widzenia nanosi się na wybór celu tylko raz, przy pierwszym odświeżeniu. */
  let punktNaniesiony = false;

  return {
    element,

    adres: cel.adres,

    odswiez(rozstrzygniecie, punkt) {
      if (!punktNaniesiony) {
        cel.ustaw(punkt);
        punktNaniesiony = true;
      }
      lancuch.odswiez(rozstrzygniecie);
      ubierzPlakietke(plakietka, rozstrzygniecie);
    },

    naZmianeCelu: cel.naZmiane,
  };
}

/**
 * Plakietka mówi pochodzenie wartości jednym zwrotem, bez rozwijania panelu. Wartość
 * domyślna katalogu jest tu nazwana wprost, więc klucz bez ani jednego zapisu nie
 * wygląda na klucz bez odpowiedzi.
 */
function ubierzPlakietke(
  plakietka: HTMLButtonElement,
  rozstrzygniecie: Rozstrzygniecie,
): void {
  const zrodlo = rozstrzygniecie.zrodlo;
  const opis = zrodlo === null ? 'wartość domyślna' : opisAdresu(adresWpisu(zrodlo));

  plakietka.textContent = opis;
  plakietka.dataset.domyslna = String(rozstrzygniecie.domyslna);
  plakietka.className = rozstrzygniecie.domyslna
    ? 'dn-plakietka dk-wskaznik__plakietka'
    : 'dn-plakietka dn-plakietka--sygnal dk-wskaznik__plakietka';
  plakietka.title = `Wartość pochodzi z: ${opis}. Naciśnij, aby zobaczyć dziedziczenie i wybrać poziom zapisu.`;
}

function naglowekDziedziczenia(): HTMLElement {
  return naglowek(
    'Dziedziczenie',
    'Zapisy tego klucza od najwęższego poziomu. Wiersz oznaczony sygnałem obowiązuje w punkcie widzenia okna.',
  );
}

function naglowekCelu(): HTMLElement {
  return naglowek(
    'Poziom zapisu',
    'Zatwierdzenie wartości w polu powyżej zapisze ją pod tym adresem.',
  );
}

function naglowek(tytul: string, opis: string): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dk-wskaznik__naglowek';

  const nazwa = document.createElement('h4');
  nazwa.className = 'dk-wskaznik__tytul';
  nazwa.textContent = tytul;

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'dn-pole-opis';
  wyjasnienie.textContent = opis;

  element.append(nazwa, wyjasnienie);
  return element;
}
