import type { AppComponent } from '../../../../shared/contract';

/**
 * Kanwa diagramu komponentów wraz z liniami zależności.
 *
 * Kanwa rysuje zbiór komponentów i nie zna kontraktu ani kanału — dostaje
 * komponenty i oddaje element.
 *
 * Linia zależności jest nazwana, nie narysowana: kontrakt niesie zależność jako
 * wykaz identyfikatorów (`dependsOn`), bez współrzędnych i bez kierunku
 * przepływu. Rysunek strzałek wymagałby danych, których nie ma.
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

/** Jeden komponent kanwy: nazwa, rodzaj, stos, zależności, kontrakt API. */
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

/** Wiersz opisu wewnątrz kafla komponentu. */
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
