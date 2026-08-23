import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { przycisk } from '../../modele/kontrolki-formularza';
import {
  BARWY_ADNOTACJI,
  KLASY_DYMKA,
  NARZEDZIA_ADNOTACJI,
  OBJASNIENIA,
} from './etykiety-browser';
import type { NarzedzieAdnotacji } from './slady-adnotacji';

/**
 * Pasek narzędzi adnotacji — pływa na warstwie płótna, przy dolnej krawędzi.
 *
 * Kontrolki trybu adnotacji. Rysowanie jest w `plotno-adnotacji.ts`, a rozmowa
 * z rdzeniem w `warstwa-adnotacji.ts`. Pasek leży na rysunku i nie zabiera
 * wysokości podglądowi, więc pozycjonowanie należy do warstwy płótna
 * (`adnotacja.css`), a nie do układu okna.
 *
 * Wybór narzędzia i barwy niesie `aria-pressed`, nie klasa: stan wypowiedziany
 * atrybutem czyta czytnik ekranu, widzi sprawdzian i sięga po niego arkusz.
 *
 * Żaden przycisk nie gaśnie. „Dodaj do rozmowy" przy pustym płótnie pozostaje
 * naciskalny i odpowiada zdaniem o tym, że nie ma czego wysłać.
 */
export interface PasekAdnotacji {
  element: HTMLElement;
  /** Nanosi bieżący wybór narzędzia i barwy na przyciski paska. */
  odswiez(): void;
}

/** Czego pasek potrzebuje od warstwy: stan płótna i cztery czynności. */
export interface UjsciaAdnotacji {
  narzedzie(): NarzedzieAdnotacji;
  ustawNarzedzie(kod: NarzedzieAdnotacji): void;
  zeton(): string;
  ustawZeton(zeton: string): void;
  /** Treść stawiana narzędziem „tekst". */
  ustawNapis(tresc: string): void;
  wyczysc(): void;
  dodajDoRozmowy(): void;
  zamknij(): void;
}

export function utworzPasekAdnotacji(ujscia: UjsciaAdnotacji): PasekAdnotacji {
  const narzedzia = NARZEDZIA_ADNOTACJI.map((pozycja) => {
    const kontrolka = przycisk(pozycja.nazwa, 'dn-btn dn-btn--sm dn-btn--zarys');
    kontrolka.addEventListener('click', () => {
      ujscia.ustawNarzedzie(pozycja.kod);
      odswiez();
    });
    return { kod: pozycja.kod, kontrolka };
  });

  const barwy = BARWY_ADNOTACJI.map((pozycja) => {
    const kontrolka = przycisk('', `dn-btn-ikona mb-adnotacja__barwa ${pozycja.klasa}`);
    // Próbka jest plamą barwy, więc jej nazwa musi paść w atrybucie: bez tego
    // wybór byłby sygnalizowany samym kolorem, czego zabrania zasada
    // dostępności żetonów stanu (`motyw/stany.css`).
    kontrolka.setAttribute('aria-label', `Barwa adnotacji: ${pozycja.nazwa}`);
    kontrolka.addEventListener('click', () => {
      ujscia.ustawZeton(pozycja.zeton);
      odswiez();
    });
    return { zeton: pozycja.zeton, kontrolka };
  });

  const napis = document.createElement('input');
  napis.type = 'text';
  napis.className = 'dn-pole-kontrolka mb-adnotacja__napis';
  napis.placeholder = 'Treść napisu';
  napis.setAttribute('aria-label', 'Treść stawiana narzędziem „tekst"');
  napis.addEventListener('input', () => ujscia.ustawNapis(napis.value));

  const wyczysc = przycisk('Wyczyść adnotacje', 'dn-btn dn-btn--sm dn-btn--zarys');
  wyczysc.addEventListener('click', () => ujscia.wyczysc());

  const doRozmowy = przycisk('Dodaj do rozmowy', 'dn-btn dn-btn--sm dn-btn--atrament');
  doRozmowy.addEventListener('click', () => ujscia.dodajDoRozmowy());

  const zamknij = przycisk('Zamknij tryb adnotacji', 'dn-btn dn-btn--sm dn-btn--zarys');
  zamknij.addEventListener('click', () => ujscia.zamknij());

  const element = document.createElement('div');
  element.className = 'mb-adnotacja__pasek';
  element.setAttribute('role', 'toolbar');
  element.setAttribute('aria-label', 'Narzędzia adnotacji na płótnie');
  element.append(
    grupa('Narzędzie', narzedzia.map((pozycja) => pozycja.kontrolka)),
    napis,
    grupa('Barwa', barwy.map((pozycja) => pozycja.kontrolka)),
    wyczysc,
    doRozmowy,
    zamknij,
    utworzDymekObjasnienia(OBJASNIENIA.adnotacja, KLASY_DYMKA),
  );

  function odswiez(): void {
    const wybrane = ujscia.narzedzie();
    for (const pozycja of narzedzia) {
      pozycja.kontrolka.setAttribute('aria-pressed', String(pozycja.kod === wybrane));
    }
    const zeton = ujscia.zeton();
    for (const pozycja of barwy) {
      pozycja.kontrolka.setAttribute('aria-pressed', String(pozycja.zeton === zeton));
    }
    // Pole napisu ma znaczenie wyłącznie przy narzędziu „tekst" — pozostaje
    // widoczne i edytowalne zawsze, a różnicę niesie znacznik dla arkusza.
    element.dataset['narzedzie'] = wybrane;
  }

  odswiez();
  return { element, odswiez };
}

/** Grupa przycisków paska wraz z nazwą czytaną przez technologie wspomagające. */
function grupa(nazwa: string, kontrolki: readonly HTMLElement[]): HTMLElement {
  const element = document.createElement('div');
  element.className = 'mb-adnotacja__grupa';
  element.setAttribute('role', 'group');
  element.setAttribute('aria-label', nazwa);
  element.append(...kontrolki);
  return element;
}
