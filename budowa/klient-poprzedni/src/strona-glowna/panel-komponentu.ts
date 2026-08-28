import { ConfigScope, type Component } from '../../../shared/contract';
import type { Kanal } from '../protokol/kanal';
import {
  poziomOKodzie,
  POZIOMY_ODRZUCANE,
  POZIOMY_PRZYPISANIA,
  utworzZmianeKomponentu,
  type BytPoziomu,
  type ZmianaKomponentu,
} from './zmiana-komponentu';

/** Panel komponentu założonego stoi w przyborniku strefy drugiej, pod formularzem zakładania, i udostępnia zmianę pól oraz przypisanie komponentu do poziomu zasięgu. */
export interface PanelKomponentu {
  element: HTMLElement;
  /** Rozwija albo zwija panel. */
  przelacz(): void;
  /** Odczytuje wykaz komponentów i byty poziomów z rdzenia na nowo. */
  odswiez(): void;
}

export interface OpcjePanelu {
  kanal: Kanal;
  /** Komponent zmieniony — strefa ma dociągnąć kafle z rdzenia na nowo. */
  naZmiane(): void;
  /** Źródło komend; wstrzykiwane dla sprawdzianu, domyślnie z kanału. */
  zrodlo?: ZmianaKomponentu;
}

export function utworzPanelKomponentu(opcje: OpcjePanelu): PanelKomponentu {
  const zrodlo = opcje.zrodlo ?? utworzZmianeKomponentu(opcje.kanal);
  let komponenty: readonly Component[] = [];
  let byty: readonly BytPoziomu[] = [];

  const komponent = document.createElement('select');
  komponent.className = 'dn-pole-kontrolka';
  komponent.dataset['ster'] = 'komponent';
  komponent.setAttribute('aria-label', 'Komponent własny do zmiany');
  komponent.addEventListener('change', () => {
    wypelnijZWybranego();
    opiszWiazanie();
  });

  const nazwa = document.createElement('input');
  nazwa.type = 'text';
  nazwa.className = 'dn-pole-kontrolka';
  nazwa.placeholder = 'nowa nazwa — puste zostawia bez zmian';
  nazwa.setAttribute('aria-label', 'Nowa nazwa komponentu');

  const opis = document.createElement('input');
  opis.type = 'text';
  opis.className = 'dn-pole-kontrolka';
  opis.placeholder = 'nowy opis — puste zostawia bez zmian';
  opis.setAttribute('aria-label', 'Nowy opis komponentu');

  const czynny = document.createElement('input');
  czynny.type = 'checkbox';
  czynny.className = 'dn-pole-kontrolka';
  czynny.id = 'dn-komponent-czynny';

  const czynnyPodpis = document.createElement('label');
  czynnyPodpis.className = 'dn-pole-etykieta';
  czynnyPodpis.htmlFor = czynny.id;
  czynnyPodpis.textContent = 'Komponent czynny';

  const czynnyWiersz = document.createElement('div');
  czynnyWiersz.className = 'dn-pole';
  czynnyWiersz.append(czynny, czynnyPodpis);

  const zmien = document.createElement('button');
  zmien.type = 'button';
  zmien.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  zmien.textContent = 'Zapisz zmianę';
  zmien.addEventListener('click', () => void wyslijZmiane());

  const poziom = document.createElement('select');
  poziom.className = 'dn-pole-kontrolka';
  poziom.dataset['ster'] = 'poziom';
  poziom.setAttribute('aria-label', 'Poziom zasięgu przypisania');
  for (const pozycja of POZIOMY_PRZYPISANIA) {
    const opcja = document.createElement('option');
    opcja.value = pozycja.kod;
    opcja.textContent = pozycja.nazwa;
    poziom.append(opcja);
  }
  poziom.addEventListener('change', () => void przestawPoziom());

  const byt = document.createElement('select');
  byt.className = 'dn-pole-kontrolka';
  byt.dataset['ster'] = 'byt';
  byt.setAttribute('aria-label', 'Byt poziomu, z którym wiąże się komponent');
  byt.addEventListener('change', () => opiszWiazanie());

  const bytWpisywany = document.createElement('input');
  bytWpisywany.type = 'text';
  bytWpisywany.className = 'dn-pole-kontrolka';
  bytWpisywany.placeholder = 'identyfikator bytu poziomu';
  bytWpisywany.setAttribute('aria-label', 'Identyfikator bytu poziomu');
  bytWpisywany.addEventListener('input', () => opiszWiazanie());

  const wiazanie = document.createElement('p');
  wiazanie.className = 'dn-pole-opis';
  wiazanie.setAttribute('aria-live', 'polite');

  const przypisz = document.createElement('button');
  przypisz.type = 'button';
  przypisz.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  przypisz.textContent = 'Przypisz komponent';
  przypisz.addEventListener('click', () => void wyslijPrzypisanie());

  const odpowiedz = document.createElement('p');
  odpowiedz.className = 'dn-strona__zaloz-odpowiedz';
  odpowiedz.hidden = true;

  const granica = document.createElement('p');
  granica.className = 'dn-strona__zaloz-granica';
  granica.textContent =
    'Przypisanie zapamiętuje poziom, na którym komponent obowiązuje — nie przenosi bytu ' +
    'modułowego, nie nadaje uprawnień i nie włącza komponentu do żadnej pętli wykonania. ' +
    'Kontrakt zaznacza przy tej komendzie, że znaczenia przypisania Właściciel nie ' +
    `rozstrzygnął. ${POZIOMY_ODRZUCANE}`;

  const element = document.createElement('div');
  element.className = 'dn-strona__zaloz-formularz';
  element.dataset['panel'] = 'komponent';
  element.hidden = true;
  element.append(
    komponent,
    nazwa,
    opis,
    czynnyWiersz,
    zmien,
    poziom,
    byt,
    bytWpisywany,
    wiazanie,
    przypisz,
    odpowiedz,
    granica,
  );

  /** Komponent wskazany w wykazie; `null`, gdy wykaz jest pusty. */
  function wybrany(): Component | null {
    return komponenty.find((pozycja) => pozycja.id === komponent.value) ?? null;
  }

  // Wartości pól bierze się z komponentu wskazanego jako podpowiedź miejsca zastanego, nie jako wysyłkę.
  function wypelnijZWybranego(): void {
    const pozycja = wybrany();
    nazwa.value = '';
    opis.value = '';
    czynny.checked = pozycja?.enabled ?? true;
  }

  /** Zdanie o wiązaniu — nazwami, nie identyfikatorami, i przed czynnością. */
  function opiszWiazanie(): void {
    const pozycja = wybrany();
    const wybranyPoziom = poziomOKodzie(poziom.value);
    if (pozycja === null) {
      wiazanie.textContent =
        'Nie ma czego wiązać: rdzeń nie oddał ani jednego komponentu własnego. Załóż komponent ' +
        'formularzem powyżej.';
      przypisz.disabled = true;
      return;
    }
    if (!wybranyPoziom.zBytem) {
      wiazanie.textContent =
        `Wiążesz komponent „${pozycja.name}" z poziomem globalnym — bez bytu, bo ten poziom ` +
        'identyfikatora bytu nie przyjmuje.';
      przypisz.disabled = false;
      return;
    }
    const nazwaBytu = nazwaWybranegoBytu(wybranyPoziom.zrodloBytow === 'wpisywany');
    if (nazwaBytu === '') {
      wiazanie.textContent =
        `Wskaż byt poziomu „${wybranyPoziom.nazwa}" — bez niego rdzeń odmówi, bo ten poziom ` +
        'wymaga identyfikatora bytu.';
      przypisz.disabled = true;
      return;
    }
    wiazanie.textContent =
      `Wiążesz komponent „${pozycja.name}" z: ${nazwaBytu} — poziom „${wybranyPoziom.nazwa}".`;
    przypisz.disabled = false;
  }

  /** Nazwa bytu wskazanego: z wykazu albo wpisana wprost. */
  function nazwaWybranegoBytu(wpisywany: boolean): string {
    if (wpisywany) {
      const wpisany = bytWpisywany.value.trim();
      return wpisany === '' ? '' : wpisany;
    }
    return byty.find((pozycja) => pozycja.id === byt.value)?.nazwa ?? '';
  }

  /** Identyfikator bytu do żądania; pusty znaczy „bez bytu". */
  function idBytu(): string {
    const wybranyPoziom = poziomOKodzie(poziom.value);
    if (!wybranyPoziom.zBytem) return '';
    return wybranyPoziom.zrodloBytow === 'wpisywany' ? bytWpisywany.value.trim() : byt.value;
  }

  /** Przestawienie poziomu dociąga byty tego poziomu i przestawia kontrolki. */
  async function przestawPoziom(): Promise<void> {
    const wybranyPoziom = poziomOKodzie(poziom.value);
    byt.hidden = !wybranyPoziom.zBytem || wybranyPoziom.zrodloBytow === 'wpisywany';
    bytWpisywany.hidden = !wybranyPoziom.zBytem || wybranyPoziom.zrodloBytow !== 'wpisywany';

    if (wybranyPoziom.zrodloBytow === 'srodowiska') byty = await zrodlo.srodowiska();
    else if (wybranyPoziom.zrodloBytow === 'okna') byty = await zrodlo.okna();
    else byty = [];

    byt.replaceChildren(
      ...byty.map((pozycja) => {
        const opcja = document.createElement('option');
        opcja.value = pozycja.id;
        opcja.textContent = pozycja.nazwa;
        return opcja;
      }),
    );
    // Wykaz pusty przy poziomie żądającym bytu jest brakiem po stronie rdzenia i musi to powiedzieć.
    if (byty.length === 0 && !byt.hidden) {
      const opcja = document.createElement('option');
      opcja.value = '';
      opcja.textContent = 'rdzeń nie oddał ani jednego bytu tego poziomu';
      byt.append(opcja);
    }
    opiszWiazanie();
  }

  async function wyslijZmiane(): Promise<void> {
    const pozycja = wybrany();
    if (pozycja === null) {
      pokaz('Nie ma czego zmieniać — rdzeń nie oddał ani jednego komponentu własnego.', false);
      return;
    }
    const nowaNazwa = nazwa.value.trim();
    const nowyOpis = opis.value.trim();
    const stanCzynnosci = czynny.checked;
    if (nowaNazwa === '' && nowyOpis === '' && stanCzynnosci === pozycja.enabled) {
      // Żądanie bez ani jednego pola zmienionego byłoby wywołaniem bez treści udającym zmianę.
      pokaz(
        'Nic nie zmieniono: nazwa i opis są puste, a stan czynności taki jak w rdzeniu. ' +
          'Pole puste zostawia wartość bez zmian — to nie żądanie pustej nazwy.',
        false,
      );
      return;
    }
    const wynik = await zrodlo.zmien({
      componentId: pozycja.id,
      ...(nowaNazwa === '' ? {} : { name: nowaNazwa }),
      ...(nowyOpis === '' ? {} : { description: nowyOpis }),
      ...(stanCzynnosci === pozycja.enabled ? {} : { enabled: stanCzynnosci }),
    });
    pokaz(wynik.zdanie, wynik.udane);
    if (wynik.udane) {
      opcje.naZmiane();
      odswiez();
    }
  }

  async function wyslijPrzypisanie(): Promise<void> {
    const pozycja = wybrany();
    if (pozycja === null) {
      pokaz('Nie ma czego przypisać — rdzeń nie oddał ani jednego komponentu własnego.', false);
      return;
    }
    const wybranyPoziom = poziomOKodzie(poziom.value);
    const identyfikator = idBytu();
    if (wybranyPoziom.zBytem && identyfikator === '') {
      pokaz(
        `Poziom „${wybranyPoziom.nazwa}" wymaga identyfikatora bytu — bez niego rdzeń odmówi.`,
        false,
      );
      return;
    }
    const wynik = await zrodlo.przypisz({
      componentId: pozycja.id,
      scope: wybranyPoziom.kod,
      ...(identyfikator === '' ? {} : { scopeId: identyfikator }),
    });
    pokaz(wynik.zdanie, wynik.udane);
    if (wynik.udane) opcje.naZmiane();
  }

  function pokaz(zdanie: string, powodzenie: boolean): void {
    odpowiedz.hidden = false;
    odpowiedz.textContent = zdanie;
    odpowiedz.dataset['powodzenie'] = String(powodzenie);
  }

  function odswiez(): void {
    void (async () => {
      const wynik = await zrodlo.wykaz();
      komponenty = wynik.komponenty;
      komponent.replaceChildren(
        ...komponenty.map((pozycja) => {
          const opcja = document.createElement('option');
          opcja.value = pozycja.id;
          opcja.textContent = pozycja.enabled
            ? pozycja.name
            : `${pozycja.name} — wyłączony`;
          return opcja;
        }),
      );
      if (komponenty.length === 0) {
        const opcja = document.createElement('option');
        opcja.value = '';
        opcja.textContent =
          wynik.zdanie === ''
            ? 'rdzeń nie ma ani jednego komponentu własnego'
            : 'wykazu komponentów rdzeń nie oddał';
        komponent.append(opcja);
        if (wynik.zdanie !== '') pokaz(wynik.zdanie, false);
      }
      wypelnijZWybranego();
      await przestawPoziom();
    })();
  }

  return {
    element,

    przelacz() {
      element.hidden = !element.hidden;
      if (!element.hidden) odswiez();
    },

    odswiez,
  };
}

/** Poziom globalny jest wartością początkową selektora poziomu, dopóki operator nie wybierze innego poziomu. */
export const POZIOM_POCZATKOWY: ConfigScope = ConfigScope.Global;
