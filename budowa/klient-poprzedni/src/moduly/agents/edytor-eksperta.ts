import type { Agent, IdentityCategory, IdentityLayer } from '../../../../shared/contract';
import { utworzFormularzTozsamosci, type FormularzTozsamosci } from './formularz-tozsamosci';
import type { StanAgentow } from './stan-agentow';
import {
  KOLEJNOSC_WARSTW,
  NAZWY_WARSTW,
  utworzWarstwyPromptu,
  type WarstwyPromptu,
} from './warstwy-promptu';

/**
 * Widok edytora Agent Buildera w zakładce tożsamości, będący łącznikiem trzech
 * części: formularza tożsamości, edytora warstw promptu oraz wykazu kategorii
 * tożsamości pochodzącego z katalogu rdzenia.
 */
export interface EdytorEksperta {
  element: HTMLElement;
  /** Nanosi eksperta czynnego na formularz i na warstwy. */
  ustaw(ekspert: Agent | null): void;
  /** Podaje katalog kategorii tożsamości albo powód jego braku. */
  ustawKategorie(kategorie: readonly IdentityCategory[], powod: string): void;
  /** Wkleja treść na koniec instrukcji; zapis zostaje w rękach Operatora. */
  dopiszDoInstrukcji(tresc: string): void;
}

/**
 * Zależności edytora: wywołanie następujące po udanym zapisie eksperta, po
 * którym moduł odświeża bibliotekę. Edytor sam żadnego wywołania rdzenia nie
 * wykonuje, więc innych zależności nie ma.
 */
export interface OpcjeEdytora {
  /** Wywoływane po udanym zapisie — moduł odświeża wtedy bibliotekę. */
  naZapisie(ekspert: Agent): void;
}

export function utworzEdytorEksperta(stan: StanAgentow, opcje: OpcjeEdytora): EdytorEksperta {
  const tozsamosc: FormularzTozsamosci = utworzFormularzTozsamosci(stan, opcje);
  const warstwy: WarstwyPromptu = utworzWarstwyPromptu(stan);

  const katalog = document.createElement('ul');
  katalog.className = 'da-warstwy';

  const tytulKatalogu = document.createElement('h4');
  tytulKatalogu.className = 'da-panel__tytul';
  tytulKatalogu.textContent = 'Katalog kategorii tożsamości rdzenia';

  const element = document.createElement('section');
  element.className = 'da-panel da-edytor';
  element.append(
    tytul('Edytor — zakładka Tożsamość'),
    tozsamosc.element,
    warstwy.element,
    tytulKatalogu,
    katalog,
  );

  return {
    element,

    ustaw(ekspert) {
      tozsamosc.ustaw(ekspert);
      warstwy.ustaw(ekspert);
      element.dataset['ekspert'] = ekspert?.id ?? '';
    },

    ustawKategorie(kategorie, powod) {
      if (powod !== '') {
        katalog.replaceChildren(wierszWarstwy(powod, 'blad'));
        return;
      }
      if (kategorie.length === 0) {
        katalog.replaceChildren(
          wierszWarstwy(
            'Katalog kategorii tożsamości jest pusty — obowiązują same warstwy eksperta.',
            'puste',
          ),
        );
        return;
      }
      katalog.replaceChildren(...KOLEJNOSC_WARSTW.map((warstwa) => wierszKategorii(warstwa, kategorie)));
    },

    dopiszDoInstrukcji(tresc) {
      tozsamosc.dopiszDoInstrukcji(tresc);
    },
  };
}

function wierszKategorii(
  warstwa: IdentityLayer,
  kategorie: readonly IdentityCategory[],
): HTMLElement {
  const nalezace = kategorie.filter((kategoria) => kategoria.layer === warstwa);
  const tresc =
    nalezace.length === 0
      ? `${NAZWY_WARSTW[warstwa]}: katalog nie ma kategorii tej warstwy`
      : `${NAZWY_WARSTW[warstwa]}: ${nalezace.map((kategoria) => kategoria.name).join(' · ')}`;
  return wierszWarstwy(tresc, 'gotowe');
}

function wierszWarstwy(tresc: string, stan: string): HTMLElement {
  const element = document.createElement('li');
  element.className = 'da-warstwy__wiersz';
  element.dataset['stan'] = stan;
  element.textContent = tresc;
  return element;
}

function tytul(nazwa: string): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'da-panel__tytul';
  element.textContent = nazwa;
  return element;
}
