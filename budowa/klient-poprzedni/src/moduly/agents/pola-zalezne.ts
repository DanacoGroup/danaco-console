import {
  CapabilitySupport,
  type AdapterCapabilities,
  type SessionConfigFieldCapability,
} from '../../../../shared/contract';

/**
 * Zestaw pól zależnych od kanału — wykaz budowany z deklaracji zdolności adaptera
 * dostawcy (`config.capabilities.get`), a nie z listy wpisanej w kodzie. Nowa
 * wersja programu dostawcy zmienia wykaz bez zmiany klienta.
 */
export interface PolaZalezne {
  element: HTMLElement;
  /** Nanosi deklarację zdolności adaptera. */
  ustaw(zdolnosci: AdapterCapabilities): void;
  /** Pokazuje trwający odczyt deklaracji. */
  ladowanie(): void;
  /** Pokazuje odmowę odczytu deklaracji. */
  blad(powod: string): void;
}

const NAZWY_OBSLUGI: Record<string, string> = {
  [CapabilitySupport.Supported]: 'obsługiwane',
  [CapabilitySupport.Partial]: 'częściowo',
  [CapabilitySupport.Unsupported]: 'nieobsługiwane',
};

export function utworzPolaZalezne(): PolaZalezne {
  const tytul = document.createElement('h4');
  tytul.className = 'da-panel__tytul';
  tytul.textContent = 'Pola zależne od kanału';

  const zrodlo = document.createElement('p');
  zrodlo.className = 'dn-pole-opis';

  const lista = document.createElement('ul');
  lista.className = 'da-zdolnosci';

  const element = document.createElement('section');
  element.className = 'da-panel da-zdolnosci__panel';
  element.append(tytul, zrodlo, lista);

  function komunikat(tresc: string, stan: string): void {
    zrodlo.textContent = tresc;
    lista.replaceChildren();
    element.dataset['stan'] = stan;
  }

  return {
    element,

    ustaw(zdolnosci) {
      const pola = zdolnosci.fields ?? [];
      element.dataset['stan'] = pola.length === 0 ? 'puste' : 'gotowe';
      zrodlo.textContent = opisZrodla(zdolnosci, pola.length);
      lista.replaceChildren(...pola.map(wiersz));
    },

    ladowanie: () => komunikat('Odczyt deklaracji zdolności adaptera w toku…', 'ladowanie'),

    blad: (powod) => komunikat(powod, 'blad'),
  };
}

function opisZrodla(zdolnosci: AdapterCapabilities, liczba: number): string {
  if (liczba === 0) {
    const powod = (zdolnosci.reason ?? '').trim();
    return powod === ''
      ? 'Adapter nie zadeklarował ani jednego pola obszaru model.'
      : `Adapter nie zadeklarował pól obszaru model: ${powod}`;
  }
  const wersja = (zdolnosci.providerVersion ?? '').trim();
  const czesci = [`adapter ${zdolnosci.adapterId}`, `transport ${zdolnosci.transport}`];
  if (wersja !== '') czesci.push(`wersja dostawcy ${wersja}`);
  return `${czesci.join(' · ')} — ${liczba} pól obszaru model.`;
}

function wiersz(pole: SessionConfigFieldCapability): HTMLElement {
  const sciezka = document.createElement('span');
  sciezka.className = 'da-zdolnosci__pole';
  sciezka.textContent = pole.fieldPath;

  const plakietka = document.createElement('span');
  plakietka.className = 'dn-plakietka da-zdolnosci__obsluga';
  plakietka.dataset['obsluga'] = pole.support;
  plakietka.textContent = NAZWY_OBSLUGI[pole.support] ?? pole.support;

  const element = document.createElement('li');
  element.className = 'da-zdolnosci__wiersz';
  element.dataset['pole'] = pole.fieldPath;
  element.append(sciezka, plakietka);

  const powod = (pole.reason ?? '').trim();
  if (powod !== '') {
    const wyjasnienie = document.createElement('span');
    wyjasnienie.className = 'da-zdolnosci__powod';
    wyjasnienie.textContent = powod;
    element.append(wyjasnienie);
  }
  return element;
}
