import type {
  StudioChain,
  StudioChainStep,
  StudioOperation,
} from '../../../../shared/contract';
import { KATEGORIE_OPERACJI, nazwaOperacji } from './kategorie-operacji';
import { przycisk } from './zadania-wykaz';

/**
 * Warsztat łańcucha operacji — składanie sekwencji wykonywanej jednym
 * poleceniem: korekta → streszczenie → zmiana tonu.
 *
 * ── Skąd biorą się czynności ─────────────────────────────────────────────────
 * Z jednego wykazu, nie z dwóch. Katalog operacji kontekstowych — dwadzieścia
 * osiem czynności w siedmiu grupach — stoi w `kategorie-operacji.ts` i ten plik
 * go WYŁĄCZNIE czyta. Do tego dochodzą operacje własne Operatora z
 * `studio.operation.list`, odróżnione plakietką: Operator ma wiedzieć, czy
 * sięga po czynność fabryczną, czy po swoją.
 *
 * ── Czym warsztat NIE jest ──────────────────────────────────────────────────
 * Nie jest drugą maszynerią przebiegu. Warsztat składa łańcuch i zapisuje go
 * (`studio.chain.save`); przebieg prowadzi ta sama kolejka zadań, którą prowadzi
 * rozkład zlecenia. Kontrakt mówi wprost, że przebieg łańcucha prowadzi pętla
 * wykonawcza okna — więc prowadzi go pętla, a nie osobny licznik obok.
 */

/** Czynności warsztatu sięgające do rdzenia. */
export interface CzynnosciWarsztatuLancucha {
  /** Zapisuje łańcuch składany pod jego nazwą. */
  zapisz(): void;
  /** Uruchamia wskazany łańcuch na dokumencie czynnym. */
  uruchom(idLancucha: string): void;
  /** Wnosi łańcuch zapisany do warsztatu, żeby go zmienić. */
  wczytajDoZmiany(lancuch: StudioChain): void;
  /** Przestawia kroki składane. */
  ustawKroki(kroki: readonly StudioChainStep[]): void;
  /** Przestawia nazwę łańcucha składanego. */
  ustawNazwe(nazwa: string): void;
}

export interface WarsztatLancucha {
  element: HTMLElement;
  odswiez(
    kroki: readonly StudioChainStep[],
    nazwa: string,
    lancuchy: readonly StudioChain[],
    wlasne: readonly StudioOperation[],
  ): void;
}

export function utworzWarsztatLancucha(czynnosci: CzynnosciWarsztatuLancucha): WarsztatLancucha {
  const poleNazwy = document.createElement('input');
  poleNazwy.type = 'text';
  poleNazwy.className = 'dn-pole';
  poleNazwy.placeholder = 'Nazwa łańcucha, na przykład „redakcja pisma urzędowego"';
  poleNazwy.setAttribute('aria-label', 'Nazwa łańcucha operacji');
  poleNazwy.addEventListener('input', () => czynnosci.ustawNazwe(poleNazwy.value));

  const wyborCzynnosci = document.createElement('select');
  wyborCzynnosci.className = 'dn-wybor';
  wyborCzynnosci.setAttribute('aria-label', 'Czynność dokładana do łańcucha');

  const dolozGuzik = przycisk('Dołóż krok', () => {
    const idAkcji = wyborCzynnosci.value;
    if (idAkcji === '') return;
    czynnosci.ustawKroki([...biezaceKroki, { actionId: idAkcji }]);
  });

  const pasSkladania = document.createElement('div');
  pasSkladania.className = 'petla-warsztat__pas';
  pasSkladania.append(poleNazwy, wyborCzynnosci, dolozGuzik);

  const listaKrokow = document.createElement('ol');
  listaKrokow.className = 'petla-warsztat__kroki';
  listaKrokow.setAttribute('aria-label', 'Kroki łańcucha w kolejności wykonania');

  const pustkaKrokow = document.createElement('p');
  pustkaKrokow.className = 'dn-tekst-3';
  pustkaKrokow.textContent =
    'Łańcuch nie ma jeszcze ani jednego kroku. Wybierz czynność z wykazu i dołóż ją — ' +
    'kroki wykonują się w kolejności, w jakiej tu stoją.';

  const zapiszGuzik = przycisk('Zapisz łańcuch', () => czynnosci.zapisz());

  const naglowekZapisanych = document.createElement('h4');
  naglowekZapisanych.className = 'petla-podtytul';
  naglowekZapisanych.textContent = 'Łańcuchy zapisane';

  const listaZapisanych = document.createElement('ul');
  listaZapisanych.className = 'petla-warsztat__zapisane';

  const pustkaZapisanych = document.createElement('p');
  pustkaZapisanych.className = 'dn-tekst-3';
  pustkaZapisanych.textContent =
    'Rdzeń nie ma jeszcze ani jednego łańcucha dla tego zasięgu. Pusty wykaz nie jest ' +
    'usterką — łańcuch zakłada Operator, a fabrycznych łańcuchów ten produkt nie udaje.';

  const element = document.createElement('section');
  element.className = 'petla-warsztat';
  element.append(
    pasSkladania,
    pustkaKrokow,
    listaKrokow,
    zapiszGuzik,
    naglowekZapisanych,
    pustkaZapisanych,
    listaZapisanych,
  );

  let biezaceKroki: readonly StudioChainStep[] = [];

  /** Buduje wykaz czynności do wyboru: fabryczne z katalogu i własne Operatora. */
  function wypelnijWybor(wlasne: readonly StudioOperation[]): void {
    const wybrane = wyborCzynnosci.value;
    wyborCzynnosci.replaceChildren();
    const puste = document.createElement('option');
    puste.value = '';
    puste.textContent = '— wybierz czynność —';
    wyborCzynnosci.append(puste);

    for (const kategoria of KATEGORIE_OPERACJI) {
      const grupa = document.createElement('optgroup');
      grupa.label = kategoria.nazwa;
      for (const operacja of kategoria.operacje) {
        const pozycja = document.createElement('option');
        pozycja.value = operacja.id;
        pozycja.textContent = operacja.nazwa;
        grupa.append(pozycja);
      }
      wyborCzynnosci.append(grupa);
    }
    if (wlasne.length > 0) {
      const grupa = document.createElement('optgroup');
      grupa.label = 'Operacje własne Operatora';
      for (const operacja of wlasne) {
        const pozycja = document.createElement('option');
        pozycja.value = operacja.id;
        pozycja.textContent = `${operacja.name} (własna)`;
        grupa.append(pozycja);
      }
      wyborCzynnosci.append(grupa);
    }
    wyborCzynnosci.value = wybrane;
  }

  function wierszKroku(krok: StudioChainStep, numer: number): HTMLElement {
    const pozycja = document.createElement('li');
    pozycja.className = 'dn-karta petla-warsztat__krok';

    const nazwa = document.createElement('span');
    nazwa.className = 'petla-warsztat__krok-nazwa';
    nazwa.textContent = `${numer + 1}. ${nazwaOperacji(krok.actionId) ?? krok.actionId}`;
    pozycja.append(nazwa);

    // Przyjęcie wyniku bez decyzji Operatora jest nastawą KROKU, nie łańcucha:
    // korekta może wchodzić sama, a zmiana tonu wymagać spojrzenia.
    const przelacznik = document.createElement('label');
    przelacznik.className = 'petla-warsztat__krok-nastawa';
    const pole = document.createElement('input');
    pole.type = 'checkbox';
    pole.checked = krok.acceptAutomatically === true;
    pole.addEventListener('change', () => {
      const zmienione = biezaceKroki.map((wpis, indeks) =>
        indeks === numer ? { ...wpis, acceptAutomatically: pole.checked } : wpis,
      );
      czynnosci.ustawKroki(zmienione);
    });
    przelacznik.append(pole, document.createTextNode(' przyjmij wynik bez decyzji'));
    pozycja.append(przelacznik);

    const sterowanie = document.createElement('div');
    sterowanie.className = 'petla-warsztat__krok-sterowanie';
    if (numer > 0) {
      sterowanie.append(
        przycisk('W górę', () => czynnosci.ustawKroki(przestaw(biezaceKroki, numer, numer - 1))),
      );
    }
    if (numer < biezaceKroki.length - 1) {
      sterowanie.append(
        przycisk('W dół', () => czynnosci.ustawKroki(przestaw(biezaceKroki, numer, numer + 1))),
      );
    }
    sterowanie.append(
      przycisk('Usuń', () =>
        czynnosci.ustawKroki(biezaceKroki.filter((_, indeks) => indeks !== numer)),
      ),
    );
    pozycja.append(sterowanie);
    return pozycja;
  }

  function wierszZapisanego(lancuch: StudioChain): HTMLElement {
    const pozycja = document.createElement('li');
    pozycja.className = 'dn-karta petla-warsztat__zapisany';

    const nazwa = document.createElement('span');
    nazwa.textContent = `${lancuch.name} — ${lancuch.steps.length} kroków`;
    pozycja.append(nazwa);

    const zasieg = document.createElement('span');
    zasieg.className = 'dn-plakietka';
    zasieg.textContent = `zasięg ${lancuch.scope}`;
    pozycja.append(zasieg);

    pozycja.append(
      przycisk('Uruchom', () => czynnosci.uruchom(lancuch.id)),
      przycisk('Zmień', () => czynnosci.wczytajDoZmiany(lancuch)),
    );
    return pozycja;
  }

  return {
    element,

    odswiez(kroki, nazwa, lancuchy, wlasne) {
      biezaceKroki = kroki;
      wypelnijWybor(wlasne);
      if (poleNazwy.value !== nazwa) poleNazwy.value = nazwa;

      listaKrokow.replaceChildren();
      pustkaKrokow.hidden = kroki.length > 0;
      listaKrokow.hidden = kroki.length === 0;
      for (let numer = 0; numer < kroki.length; numer += 1) {
        const krok = kroki[numer];
        if (krok !== undefined) listaKrokow.append(wierszKroku(krok, numer));
      }
      // Zapis łańcucha bez nazwy albo bez kroków nie ma czego zapisać —
      // przycisk jest wtedy nieosiągalny, a nie milcząco bezskuteczny.
      zapiszGuzik.disabled = kroki.length === 0 || nazwa.trim() === '';

      listaZapisanych.replaceChildren();
      pustkaZapisanych.hidden = lancuchy.length > 0;
      listaZapisanych.hidden = lancuchy.length === 0;
      for (const lancuch of lancuchy) listaZapisanych.append(wierszZapisanego(lancuch));
    },
  };
}

/** Przestawia krok z jednego miejsca w drugie, zachowując pozostałe. */
function przestaw(
  kroki: readonly StudioChainStep[],
  skad: number,
  dokad: number,
): StudioChainStep[] {
  const kopia = [...kroki];
  const krok = kopia[skad];
  if (krok === undefined) return kopia;
  kopia.splice(skad, 1);
  kopia.splice(dokad, 0, krok);
  return kopia;
}
