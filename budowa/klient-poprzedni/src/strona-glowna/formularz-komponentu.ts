import type { Kanal } from '../protokol/kanal';
import { utworzPanelKomponentu } from './panel-komponentu';
import {
  RODZAJE_DO_ZALOZENIA,
  RODZAJ_BEZ_MAGAZYNU,
  utworzZalozenieKomponentu,
  type ZalozenieKomponentu,
} from './zalozenie-komponentu';

/** Formularz zakładania komponentu własnego jest wejściem do zakładania komponentu ze strony głównej, umieszczonym pod panelem komponentu już założonego. */
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
    // Wykaz komponentów dociąga strefa, pytając rdzeń na nowo.
    opcje.naZalozenie();
  }

  function powiedz(tresc: string, powodzenie: boolean): void {
    odpowiedz.textContent = tresc;
    odpowiedz.hidden = tresc === '';
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  // Panel komponentu założonego stoi pod formularzem i otwiera się razem z nim segmentem „Dodaj nowy”.
  const panel = utworzPanelKomponentu({
    kanal: opcje.kanal,
    naZmiane: opcje.naZalozenie,
  });

  const element = document.createElement('div');
  element.className = 'dn-strona__zaloz';
  element.append(formularz, panel.element);

  return {
    element,
    // Formularz otwiera segment „Dodaj nowy” belki strefy trzeciej.
    przelacz() {
      const otwarty = formularz.hidden;
      formularz.hidden = !otwarty;
      panel.przelacz();
      if (otwarty) nazwa.focus();
    },
  };
}
