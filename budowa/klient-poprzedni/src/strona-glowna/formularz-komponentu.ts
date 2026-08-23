import type { Kanal } from '../protokol/kanal';
import { utworzPanelKomponentu } from './panel-komponentu';
import {
  RODZAJE_DO_ZALOZENIA,
  RODZAJ_BEZ_MAGAZYNU,
  utworzZalozenieKomponentu,
  type ZalozenieKomponentu,
} from './zalozenie-komponentu';

/**
 * Formularz zakładania komponentu własnego — wejście do `component.create`
 * ze strony głównej.
 *
 * Formularz stoi zwinięty: środkiem ciężkości strony głównej są karty
 * środowisk, a formularz stale rozłożony odbierałby im pas ekranu przy
 * czynności wykonywanej rzadko. Rozwija go segment „Dodaj nowy" belki strefy
 * trzeciej.
 *
 * Wykaz rodzajów bierze się z `zalozenie-komponentu.ts`: `assistant` odmawia,
 * bo platforma nie ma magazynu profili. Powód stoi na ekranie, nie tylko
 * w komentarzu — inaczej zostałaby sama krótsza lista bez wyjaśnienia.
 *
 * Pod formularzem stoi panel komponentu już założonego (`panel-komponentu.ts`)
 * — zmiana i przypisanie. Zakładanie i zmiana idą jedną szufladą, bo są jedną
 * pracą w dwóch krokach, a osobny uchwyt do drugiego kroku kazałby Operatorowi
 * szukać go po założeniu komponentu.
 */
export interface FormularzKomponentu {
  element: HTMLElement;
  /** Rozwija albo zwija formularz — wywołuje go segment „Dodaj nowy". */
  przelacz(): void;
}

export interface OpcjeFormularza {
  kanal: Kanal;
  /** Komponent założony — strefa ma dociągnąć wykaz z rdzenia na nowo. */
  naZalozenie(): void;
}

export function utworzFormularzKomponentu(opcje: OpcjeFormularza): FormularzKomponentu {
  const zrodlo: ZalozenieKomponentu = utworzZalozenieKomponentu(opcje.kanal);

  const rodzaj = document.createElement('select');
  rodzaj.className = 'dn-pole-kontrolka dn-strona__zaloz-rodzaj';
  rodzaj.setAttribute('aria-label', 'Rodzaj komponentu');
  for (const pozycja of RODZAJE_DO_ZALOZENIA) {
    const opcja = document.createElement('option');
    opcja.value = pozycja.kod;
    opcja.textContent = pozycja.nazwa;
    rodzaj.append(opcja);
  }

  const nazwa = document.createElement('input');
  nazwa.type = 'text';
  nazwa.className = 'dn-pole-kontrolka dn-strona__zaloz-nazwa';
  nazwa.placeholder = 'nazwa komponentu';
  nazwa.setAttribute('aria-label', 'Nazwa komponentu');

  const opis = document.createElement('input');
  opis.type = 'text';
  opis.className = 'dn-pole-kontrolka dn-strona__zaloz-opis';
  opis.placeholder = 'opis — pole opcjonalne';
  opis.setAttribute('aria-label', 'Opis komponentu');

  const zaloz = document.createElement('button');
  zaloz.type = 'submit';
  zaloz.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  zaloz.textContent = 'Utwórz';

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dn-strona__zaloz-odpowiedz';
  odpowiedz.hidden = true;

  const granica = document.createElement('p');
  granica.className = 'dn-strona__zaloz-granica';
  granica.textContent = RODZAJ_BEZ_MAGAZYNU.powod;

  const formularz = document.createElement('form');
  formularz.className = 'dn-strona__zaloz-formularz';
  formularz.hidden = true;
  formularz.append(rodzaj, nazwa, opis, zaloz, odpowiedz, granica);

  formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    void wyslij();
  });

  async function wyslij(): Promise<void> {
    powiedz('Zakładanie komponentu…', true);
    const wynik = await zrodlo.zaloz(
      rodzaj.value as (typeof RODZAJE_DO_ZALOZENIA)[number]['kod'],
      nazwa.value,
      opis.value,
    );
    powiedz(wynik.zdanie, wynik.udane);
    if (!wynik.udane) return;
    nazwa.value = '';
    opis.value = '';
    // Wykaz dociąga strefa, pytając rdzeń na nowo. Doklejenie kafla z odpowiedzi
    // pokazałoby stan, którego rdzeń nie potwierdził drugim odczytem.
    opcje.naZalozenie();
  }

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  // Panel komponentu założonego stoi pod formularzem zakładania i otwiera się
  // razem z nim: „Dodaj nowy" jest jedynym uchwytem przybornika, a czynności nad
  // komponentem już założonym (`component.update`, `component.assign`) nie mają
  // w strefie drugiej innego wejścia — kafel prowadzi do modułu, nie do zmiany.
  // Drugi uchwyt do jednej szuflady byłby drugą drogą do jednego bytu.
  const panel = utworzPanelKomponentu({
    kanal: opcje.kanal,
    naZmiane: opcje.naZalozenie,
  });

  const element = document.createElement('div');
  element.className = 'dn-strona__zaloz';
  element.append(formularz, panel.element);

  return {
    element,
    // Formularz nie ma własnego uchwytu — otwiera go segment „Dodaj nowy"
    // belki strefy trzeciej. Drugi przycisk do tego samego formularza byłby
    // drugą drogą do jednego bytu.
    przelacz() {
      const otwarty = formularz.hidden;
      formularz.hidden = !otwarty;
      panel.przelacz();
      if (otwarty) nazwa.focus();
    },
  };
}
