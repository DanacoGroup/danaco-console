import type { SettingCategory } from '../../../shared/contract';
import { utworzPoleUstawienia, type PoleUstawienia } from './pole-ustawienia';
import type { StanKonfiguracji } from './stan-konfiguracji';
import { polePozostajeWidoczne } from './widocznosc-pol';

/**
 * Prawa kolumna okna: formularz jednej kategorii zbudowany z katalogu. Panel
 * nie zna pola z góry — dostaje kategorię, pyta stan o jej pozycje i buduje
 * wiersz na pozycję. Kategoria bez pozycji mówi wprost, że katalog nie ma dla
 * niej wierszy.
 */
export interface PanelKategorii {
  /** Kolumna treści osadzana w ciele okna. */
  element: HTMLElement;
  /** Przebudowuje formularz na wskazaną kategorię. */
  pokaz(kategoria: SettingCategory | null): void;
  /** Nanosi wartości i pochodzenia ze stanu na zbudowane pola. */
  odswiez(): void;
}

/**
 * Kod kategorii katalogu rozwiniętej we własnym oknie punktów izolacji —
 * jedyna nazwa zakresu, którą ten plik zna; punkt styku dwóch okien, nie
 * lista pól zaszyta w kodzie interfejsu.
 */
const KATEGORIA_IZOLACJI = 'izolacja';

export function utworzPanelKategorii(
  stan: StanKonfiguracji,
  naOtwarcieIzolacji: () => void,
): PanelKategorii {
  const naglowek = document.createElement('header');
  naglowek.className = 'dk-tresc__naglowek';

  const tytul = document.createElement('h3');
  tytul.className = 'dk-tresc__tytul';

  const opis = document.createElement('p');
  opis.className = 'dk-tresc__opis';

  naglowek.append(tytul, opis);

  const formularz = document.createElement('div');
  formularz.className = 'dk-formularz';

  const element = document.createElement('section');
  element.className = 'dk-tresc';
  element.append(naglowek, formularz);

  let pola: PoleUstawienia[] = [];

  function odswiez(): void {
    for (const pole of pola) {
      pole.odswiez();
      const definicja = stan.definicjaKlucza(pole.klucz);
      pole.ustawWidocznosc(definicja === null || polePozostajeWidoczne(definicja, stan));
    }
  }

  return {
    element,

    pokaz(kategoria) {
      if (kategoria === null) {
        tytul.textContent = 'Konfiguracja';
        opis.textContent = '';
        opis.hidden = true;
        pola = [];
        formularz.replaceChildren(stanPusty(TYTUL_BEZ_KATEGORII, zdanieBezKategorii(stan)));
        return;
      }

      tytul.textContent = kategoria.name;
      opis.textContent = kategoria.description ?? '';
      opis.hidden = kategoria.description === undefined;

      const definicje = stan.definicjeKategorii(kategoria.id);
      pola = definicje.map((definicja) => utworzPoleUstawienia({ definicja, stan }));

      const przejscia =
        kategoria.id === KATEGORIA_IZOLACJI ? [przejscieDoIzolacji(naOtwarcieIzolacji)] : [];

      formularz.replaceChildren(
        ...przejscia,
        ...(pola.length > 0
          ? pola.map((pole) => pole.element)
          : [
              stanPusty(
                'Kategoria bez pozycji',
                `Katalog nie ma ani jednej pozycji w kategorii „${kategoria.name}".`,
              ),
            ]),
      );

      odswiez();
    },

    odswiez,
  };
}

/**
 * Przejście do okna punktów izolacji nad formularzem zakresu „Izolacja".
 * Zdanie mówi, czego formularz poniżej nie daje, zamiast samego „otwórz".
 * Przycisk nie zastępuje formularza: obie drogi pozostają czynne.
 */
function przejscieDoIzolacji(naOtwarcie: () => void): HTMLElement {
  const zdanie = document.createElement('p');
  zdanie.className = 'dn-pole-opis';
  zdanie.textContent =
    'Zakres „Izolacja" jest rozwinięty w osobnym oknie o trzech panelach: selektor zasięgu, ' +
    'macierz jedenastu punktów i podgląd polityki efektywnej po dziedziczeniu. Wiersze poniżej ' +
    'to te same klucze, ale zawsze pod jednym punktem widzenia okna — bez wyboru poziomu, na ' +
    'którym reguła ma obowiązywać, i bez podglądu wyniku.';

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = 'dn-btn dn-btn--zarys';
  przycisk.textContent = 'Otwórz okno konfiguracji punktów izolacji →';
  przycisk.addEventListener('click', naOtwarcie);

  const element = document.createElement('div');
  element.className = 'dk-formularz__przejscie';
  element.append(zdanie, przycisk);
  return element;
}

const TYTUL_BEZ_KATEGORII = 'Brak kategorii ustawień';

/**
 * Zdanie stanu pustego dobrane do fazy odczytu: ten sam pusty katalog znaczy
 * co innego w trakcie pytania rdzenia, co innego po odmowie, a co innego,
 * gdy kategorii naprawdę nie ma.
 */
function zdanieBezKategorii(stan: StanKonfiguracji): string {
  switch (stan.faza()) {
    case 'spoczynek':
    case 'odczyt':
      return 'Katalog kategorii jedzie z rdzenia. Okno pozostaje otwarte i czynne.';
    case 'blad':
      return 'Katalog kategorii nie dotarł z rdzenia — powód nad stopką okna. Odczyt można powtórzyć.';
    case 'gotowe':
      return 'Rdzeń odpowiedział, ale nie zna ani jednej kategorii ustawień. Katalog kategorii jest pusty.';
  }
}

/**
 * Stan pusty formularza — rodzina `.dn-pusty-stan`. Mówi, że brak treści jest
 * oczekiwany, nie błędem ładowania. Eksportowany, bo panel obszarów sesji
 * potrzebuje tego samego stanu pustego.
 */
export function stanPusty(
  tytul: string,
  zdanie: string,
  klasa = 'dk-formularz__pusto',
): HTMLElement {
  const element = document.createElement('div');
  element.className = `dn-pusty-stan ${klasa}`;

  const naglowek = document.createElement('p');
  naglowek.className = 'dn-pusty-stan-tytul';
  naglowek.textContent = tytul;

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';
  opis.textContent = zdanie;

  element.append(naglowek, opis);
  return element;
}
