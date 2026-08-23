import type { AccessGrant } from '../../../shared/contract';
import type { StanDostepow } from './stan-dostepow';
import { utworzWierszNadania } from './wiersz-nadania';

/**
 * Zbiór nadań okna rozmowy — druga połowa sekcji dostępów.
 *
 * Okno ma zbiór nadań, nie jedno nadanie. Lista pokazuje ten zbiór
 * w kolejności, z jawnym oznaczeniem nadania głównego: kolejność rozstrzyga,
 * w jakiej postaci nadania trafią do konfiguracji mostów, a główne wskazuje
 * punkt, od którego model zaczyna.
 *
 * Lista przebudowuje się przy każdej zmianie zbioru. To świadomy wybór: zmiana
 * kolejności albo oznaczenia głównego przestawia wszystkie wiersze naraz, więc
 * nanoszenie wartości na wiersze już zbudowane byłoby trudniejsze i mniej
 * wierne niż zbudowanie ich od nowa z odpowiedzi rdzenia.
 */
export interface ListaNadan {
  /** Element osadzany w sekcji dostępów. */
  element: HTMLElement;
  /** Nanosi zbiór nadań okna czynnego. */
  odswiez(): void;
}

export function utworzListeNadan(stan: StanDostepow): ListaNadan {
  const tytul = document.createElement('h3');
  tytul.className = 'dd-nadania__tytul';
  tytul.textContent = 'Dostępy nadane temu oknu';

  const opis = document.createElement('p');
  opis.className = 'dd-nadania__opis';

  const naglowek = document.createElement('header');
  naglowek.className = 'dd-nadania__naglowek';
  naglowek.append(tytul, opis);

  const uwaga = document.createElement('p');
  uwaga.className = 'dd-nadania__uwaga';
  uwaga.setAttribute('role', 'status');
  uwaga.hidden = true;

  const wykaz = document.createElement('ul');
  wykaz.className = 'dd-nadania__wykaz';

  const element = document.createElement('section');
  element.className = 'dd-nadania';
  element.append(naglowek, uwaga, wykaz);

  let podpis = '';

  return {
    element,

    odswiez() {
      const nadania = stan.nadania();
      const okno = stan.oknoID();
      const nowyPodpis = `${okno}#${podpisZbioru(nadania)}`;
      opis.textContent = zdanieOpisu(okno, nadania.length);

      if (nowyPodpis === podpis) return;
      podpis = nowyPodpis;

      if (okno === '') {
        uwaga.hidden = true;
        wykaz.replaceChildren(pusty(KOMUNIKAT_BEZ_OKNA));
        return;
      }

      if (nadania.length === 0) {
        uwaga.hidden = true;
        wykaz.replaceChildren(pusty(KOMUNIKAT_BEZ_NADAN));
        return;
      }

      uwaga.hidden = nadania.some((nadanie) => nadanie.primary);
      uwaga.textContent = KOMUNIKAT_BEZ_GLOWNEGO;

      wykaz.replaceChildren(
        ...nadania.map((nadanie, numer) =>
          utworzWierszNadania({
            nadanie,
            punkt: stan.punkt(nadanie.accessPointId),
            stan,
            pozycja: numer + 1,
            liczba: nadania.length,
          }).element,
        ),
      );
    },
  };
}

/**
 * Podpis zbioru — zmiana któregokolwiek członu przebudowuje listę. Człony
 * odpowiadają dokładnie temu, co wiersz rysuje.
 */
function podpisZbioru(nadania: readonly AccessGrant[]): string {
  return nadania
    .map(
      (nadanie) =>
        `${nadanie.id}:${nadanie.order}:${nadanie.mode}:${String(nadanie.primary)}:${(
          nadanie.roots ?? []
        ).join(',')}`,
    )
    .join('|');
}

/** Zdanie nagłówka: które okno i ile nadań. */
function zdanieOpisu(okno: string, liczba: number): string {
  if (okno === '') return 'Sekcja nie jest związana z oknem rozmowy.';
  return `Okno ${okno} · nadań: ${liczba}. Nadanie żyje per okno rozmowy, nie per sesja — drugie okno tej samej sesji ma własny zbiór.`;
}

const KOMUNIKAT_BEZ_OKNA =
  'Sekcja czeka na okno rozmowy. Nadanie dostępu jest wiązaniem okna z punktem, więc bez okna nie ma czego nadać.';

const KOMUNIKAT_BEZ_NADAN =
  'Okno nie ma ani jednego nadania — model nie sięgnie z niego do żadnej maszyny ani katalogu. Nadaj dostęp z wykazu punktów obok.';

const KOMUNIKAT_BEZ_GLOWNEGO =
  'Żadne nadanie nie jest oznaczone jako główne. Oznacz jedno — od niego model zaczyna, gdy nie wskazano punktu wprost.';

/** Zdanie zamiast pustej listy — Operator ma wiedzieć, dlaczego nic nie ma. */
function pusty(tresc: string): HTMLLIElement {
  const element = document.createElement('li');
  element.className = 'dd-nadania__pusty';
  element.textContent = tresc;
  return element;
}
