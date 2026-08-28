import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import './menu.css';

/**
 * Menu rozwijane to przycisk i lista, która wychodzi spod niego; nie zna treści listy ani nie zakłada bycia jedynym na scenie.
 */
export interface MenuRozwijane {
  element: HTMLElement;
  otwarte(): boolean;
  ustaw(otwarte: boolean): void;
  przelacz(): void;
  /** Podmienia zawartość listy; menu nie zna, co w niej stoi. */
  ustawTresc(dzieci: readonly HTMLElement[]): void;
  /** Zapala lub gasi znacznik „coś nowego" na uchwycie. */
  ustawZnacznik(widoczny: boolean): void;
  /** Zdejmuje nasłuchy dokumentu. Obowiązkowe. */
  zamknij(): void;
}

export interface OpcjeMenuRozwijanego {
  ikona: NazwaIkony;
  etykieta: string;
}

export function utworzMenuRozwijane(opcje: OpcjeMenuRozwijanego): MenuRozwijane {
  let otwarte = false;

  const element = document.createElement('div');
  element.className = 'dn-menu';

  const uchwyt = document.createElement('button');
  uchwyt.type = 'button';
  uchwyt.className = 'dn-btn-ikona dn-menu__uchwyt';
  uchwyt.setAttribute('aria-haspopup', 'menu');
  uchwyt.setAttribute('aria-label', opcje.etykieta);
  uchwyt.title = opcje.etykieta;
  uchwyt.append(elementIkony(opcje.ikona, { rozmiar: 16 }));

  const lista = document.createElement('div');
  lista.className = 'dn-menu__lista';
  lista.setAttribute('role', 'menu');
  lista.setAttribute('aria-label', opcje.etykieta);
  lista.hidden = true;

  element.append(uchwyt, lista);

  uchwyt.addEventListener('click', () => ustaw(!otwarte));
  uchwyt.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'ArrowDown' && zdarzenie.key !== 'ArrowUp') return;
    zdarzenie.preventDefault();
    ustaw(true);
    przesunOgnisko(zdarzenie.key === 'ArrowDown' ? 0 : -1, true);
  });

  lista.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key === 'ArrowDown') {
      zdarzenie.preventDefault();
      przesunOgnisko(1);
      return;
    }
    if (zdarzenie.key === 'ArrowUp') {
      zdarzenie.preventDefault();
      przesunOgnisko(-1);
    }
  });

  // Klawisz wyjścia zwija menu i wraca ogniskiem na uchwyt, jak w szufladzie sterowania.
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Escape' || !otwarte) return;
    zdarzenie.stopPropagation();
    ustaw(false);
    uchwyt.focus();
  });

  function pozycje(): HTMLElement[] {
    // Trzy role naraz, choć stawiana jest dziś tylko rola menuitem w tym pliku.
    const wybor = '[role="menuitem"], [role="menuitemcheckbox"], [role="menuitemradio"]';
    return [...lista.querySelectorAll<HTMLElement>(wybor)].filter((poz) => !poz.hidden);
  }

  /** Przesuwa ognisko o `krok` pozycji; `odNowa` liczy od krańca listy. */
  function przesunOgnisko(krok: number, odNowa = false): void {
    const wykaz = pozycje();
    if (wykaz.length === 0) return;
    if (odNowa) {
      const cel = krok >= 0 ? wykaz[0] : wykaz[wykaz.length - 1];
      cel.focus();
      return;
    }
    const teraz = wykaz.indexOf(document.activeElement as HTMLElement);
    const nastepna = (teraz + krok + wykaz.length) % wykaz.length;
    wykaz[nastepna].focus();
  }

  function poKlikuPoza(zdarzenie: Event): void {
    const cel = zdarzenie.target;
    if (cel instanceof Node && element.contains(cel)) return;
    ustaw(false);
  }

  function ustaw(nowe: boolean): void {
    if (nowe === otwarte) return;
    otwarte = nowe;
    lista.hidden = !otwarte;
    uchwyt.setAttribute('aria-expanded', String(otwarte));
    element.dataset.otwarte = otwarte ? 'tak' : 'nie';

    // Nasłuchy dokumentu żyją wyłącznie przy rozwiniętym menu.
    if (otwarte) document.addEventListener('pointerdown', poKlikuPoza, true);
    else document.removeEventListener('pointerdown', poKlikuPoza, true);
  }

  uchwyt.setAttribute('aria-expanded', 'false');
  element.dataset.otwarte = 'nie';

  return {
    element,
    otwarte: () => otwarte,
    ustaw,
    przelacz: () => ustaw(!otwarte),

    ustawTresc(dzieci) {
      lista.replaceChildren(...dzieci);
    },

    ustawZnacznik(widoczny) {
      if (widoczny) uchwyt.dataset.znacznik = 'tak';
      else delete uchwyt.dataset.znacznik;
    },

    zamknij() {
      document.removeEventListener('pointerdown', poKlikuPoza, true);
      otwarte = false;
      lista.hidden = true;
      uchwyt.setAttribute('aria-expanded', 'false');
      element.dataset.otwarte = 'nie';
    },
  };
}
