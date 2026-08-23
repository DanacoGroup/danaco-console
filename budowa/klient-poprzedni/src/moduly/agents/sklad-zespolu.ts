import type { Agent } from '../../../../shared/contract';

/**
 * Skład zespołu — wybór ekspertów biblioteki, w kolejności zaznaczania.
 *
 * Kolejność jest treścią, nie porządkiem wyświetlania: pole `agentIds` niesie
 * ekspertów w kolejności nadanej przez Operatora. Składanie tablicy przez
 * przejście po bibliotece brałoby kolejność z wykazu ekspertów, więc zespół
 * wracałby po zapisie w innym porządku niż złożony — wybór trzyma zatem własną
 * listę i dopisuje na jej koniec.
 *
 * Zespół wczytany może wskazywać eksperta spoza bieżącego wykazu: zarchiwizowanego
 * albo odciętego frazą zawężającą. Taki wpis zostaje w składzie i pokazuje się
 * wierszem nazywającym brak, żeby ponowny zapis nie okroił zespołu po cichu.
 */
export interface SkladZespolu {
  element: HTMLElement;
  /** Skład bieżący w kolejności nadanej przez Operatora. */
  wybrani(): readonly string[];
  /** Nanosi skład wprost — używane przy wczytaniu zespołu z rdzenia. */
  ustawSklad(idEkspertow: readonly string[]): void;
  /** Nanosi wykaz ekspertów, z których da się wybierać. */
  ustawBiblioteke(eksperci: readonly Agent[]): void;
}

export function utworzSkladZespolu(naZmiane: () => void): SkladZespolu {
  let biblioteka: readonly Agent[] = [];
  let sklad: string[] = [];

  const lista = document.createElement('ul');
  lista.className = 'da-sklad';

  const element = document.createElement('div');
  element.className = 'da-panel';

  const tytul = document.createElement('h5');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Skład zespołu';

  const opis = document.createElement('p');
  opis.className = 'dn-pole-opis';
  opis.textContent =
    'Zaznaczenie dopisuje eksperta na koniec składu — kolejność zaznaczania jest kolejnością zespołu.';

  element.append(tytul, opis, lista);

  function przelacz(idEksperta: string): void {
    sklad = sklad.includes(idEksperta)
      ? sklad.filter((wpis) => wpis !== idEksperta)
      : [...sklad, idEksperta];
    przerysuj();
    naZmiane();
  }

  function wiersz(idEksperta: string, nazwa: string, pozycja: number): HTMLElement {
    const pole = document.createElement('input');
    pole.type = 'checkbox';
    pole.className = 'dn-przelacznik';
    pole.checked = pozycja > 0;
    pole.addEventListener('change', () => przelacz(idEksperta));

    const etykieta = document.createElement('label');
    etykieta.className = 'da-sklad__nazwa';
    etykieta.textContent = pozycja > 0 ? `${pozycja}. ${nazwa}` : nazwa;
    etykieta.prepend(pole);

    const wpis = document.createElement('li');
    wpis.className = 'da-sklad__wiersz';
    wpis.dataset['ekspert'] = idEksperta;
    wpis.append(etykieta);
    return wpis;
  }

  function przerysuj(): void {
    const wiersze: HTMLElement[] = [];
    for (const ekspert of biblioteka) {
      wiersze.push(wiersz(ekspert.id, ekspert.name, sklad.indexOf(ekspert.id) + 1));
    }
    // Eksperci składu spoza bieżącego wykazu — nazwani brakiem, nie pominięci.
    const znani = new Set(biblioteka.map((ekspert) => ekspert.id));
    for (const idEksperta of sklad) {
      if (znani.has(idEksperta)) continue;
      wiersze.push(
        wiersz(idEksperta, `${idEksperta} — poza bieżącym wykazem`, sklad.indexOf(idEksperta) + 1),
      );
    }
    lista.replaceChildren(...wiersze);
  }

  przerysuj();

  return {
    element,
    wybrani: () => sklad,

    ustawSklad(idEkspertow) {
      sklad = [...idEkspertow];
      przerysuj();
    },

    ustawBiblioteke(eksperci) {
      biblioteka = eksperci;
      przerysuj();
    },
  };
}
