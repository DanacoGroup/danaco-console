import type { AppComponent } from '../../../../shared/contract';

/**
 * Kanwa diagramu komponentów wraz z ich zależnościami. Kanwa dostaje zbiór
 * komponentów i oddaje element, nie znając ani kontraktu, ani kanału. Zależność
 * jest nazwana, a nie narysowana, ponieważ kontrakt niesie ją jako wykaz
 * identyfikatorów.
 */
export interface KanwaKomponentow {
  element: HTMLElement;
  /** Nanosi zbiór komponentów; pusty zbiór zostawia kanwę pustą. */
  nanies(komponenty: readonly AppComponent[]): void;
}

export function utworzKanweKomponentow(usun: (idKomponentu: string) => void): KanwaKomponentow {
  const lista = document.createElement('ul');
  lista.className = 'mp-kanwa';

  return {
    element: lista,
    nanies(komponenty) {
      lista.replaceChildren(...komponenty.map((komponent) => kafel(komponent, komponenty, usun)));
    },
  };
}

/**
 * Jeden kafel kanwy: nazwa komponentu, plakietka rodzaju, przycisk zdjęcia
 * z kanwy oraz wiersze stosu, opisu, zależności i kontraktu API. Wiersz
 * nieobowiązkowy powstaje tylko wtedy, gdy komponent niesie jego wartość.
 */
function kafel(
  komponent: AppComponent,
  wszystkie: readonly AppComponent[],
  usun: (idKomponentu: string) => void,
): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mp-kanwa__nazwa';
  nazwa.textContent = komponent.name;

  const rodzaj = document.createElement('span');
  rodzaj.className = 'dn-plakietka dn-plakietka--rola';
  rodzaj.textContent = komponent.kind;

  const usuniecie = document.createElement('button');
  usuniecie.type = 'button';
  usuniecie.className = 'dn-btn dn-btn--duch dn-btn--sm';
  usuniecie.textContent = 'Usuń z kanwy';
  usuniecie.addEventListener('click', () => usun(komponent.id));

  const glowa = document.createElement('div');
  glowa.className = 'mp-kanwa__glowa';
  glowa.append(nazwa, rodzaj, usuniecie);

  const element = document.createElement('li');
  element.className = 'mp-kanwa__kafel';
  element.dataset['komponent'] = komponent.id;
  element.append(glowa);

  if (komponent.stack !== undefined) element.append(wiersz('Stos', komponent.stack));
  if (komponent.description !== undefined) element.append(wiersz('Opis', komponent.description));
  element.append(wiersz('Zależy od', opisZaleznosci(komponent, wszystkie)));
  if (komponent.apiContract !== undefined) {
    element.append(wiersz('Kontrakt API', komponent.apiContract));
  }
  return element;
}

/**
 * Wiersz opisu wewnątrz kafla komponentu: etykieta oraz treść w dwóch osobnych
 * elementach, żeby arkusz stylów rozróżniał je bez wnikania w tekst.
 */
function wiersz(etykieta: string, tresc: string): HTMLElement {
  const nazwa = document.createElement('span');
  nazwa.className = 'mp-kanwa__etykieta';
  nazwa.textContent = etykieta;

  const wartosc = document.createElement('span');
  wartosc.className = 'mp-kanwa__wartosc';
  wartosc.textContent = tresc;

  const element = document.createElement('p');
  element.className = 'mp-kanwa__wiersz';
  element.append(nazwa, wartosc);
  return element;
}

/**
 * Zależności komponentu opisane nazwami, nie identyfikatorami.
 *
 * Identyfikator wskazujący komponent spoza kanwy zostaje w treści wprost, żeby
 * niespójność układu była widoczna.
 */
function opisZaleznosci(komponent: AppComponent, wszystkie: readonly AppComponent[]): string {
  const wskazane = komponent.dependsOn ?? [];
  if (wskazane.length === 0) return 'brak — komponent samodzielny';
  return wskazane
    .map((identyfikator) => {
      const cel = wszystkie.find((inny) => inny.id === identyfikator);
      return cel === undefined ? `${identyfikator} (poza kanwą)` : cel.name;
    })
    .join(' · ');
}
