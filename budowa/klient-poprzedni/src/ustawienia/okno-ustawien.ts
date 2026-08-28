import { oznaczOknoAplikacji } from '../komponenty/okno-aplikacji';
import './ustawienia.css';

import { elementIkony } from '../ikony/ikony';
import type { Kanal } from '../protokol/kanal';
import { REJESTR_SEKCJI_USTAWIEN, type SekcjaUstawien, type KodSekcjiUstawien } from './sekcje';
import { utworzStanUstawien, type StanUstawien } from './stan-ustawien';

/** Okno ustawień jest ramą: nawigacją sekcji i montażem sekcji czynnej, zbudowaną na natywnym oknie dialogowym tak jak okno konfiguracji i okno modeli. */
export interface OknoUstawien {
  /** Element `<dialog>` osadzany w dokumencie. */
  element: HTMLDialogElement;
  /** Otwiera okno i zleca odczyt sekcji czynnej. */
  otworz(): void;
  /** Zamyka okno; stan i subskrypcje zostają. */
  zamknij(): void;
  /** Odłącza subskrypcje wszystkich sekcji i usuwa okno z dokumentu. */
  rozlacz(): void;
}

export function utworzOknoUstawien(kanal: Kanal): OknoUstawien {
  const stan: StanUstawien = utworzStanUstawien();
  const sekcje = new Map<KodSekcjiUstawien, SekcjaUstawien>(
    REJESTR_SEKCJI_USTAWIEN.map((opis) => [opis.kod, opis.utworz(kanal)]),
  );

  function sekcjaCzynna(): SekcjaUstawien {
    const sekcja = sekcje.get(stan.czynna());
    if (sekcja === undefined) {
      throw new Error(`Okno Ustawień: brak sekcji zbudowanej dla kodu „${stan.czynna()}".`);
    }
    return sekcja;
  }

  const nawigacja = utworzNawigacjeSekcji((kod) => stan.ustawCzynna(kod));

  const cialo = document.createElement('div');
  cialo.className = 'du-okno__cialo';

  const wnetrzeCiala = document.createElement('div');
  wnetrzeCiala.className = 'du-okno__tresc';

  // Kolumna nawigacji, budowana zawsze, jest osadzana tylko, gdy rejestr niesie więcej niż jedną sekcję.
  if (REJESTR_SEKCJI_USTAWIEN.length > 1) {
    cialo.append(nawigacja.element, wnetrzeCiala);
  } else {
    cialo.append(wnetrzeCiala);
  }

  const element = oznaczOknoAplikacji({
    element: document.createElement('dialog'),
    kod: 'okno-ustawien',
    nazwa: 'Okno Ustawień',
  });
  element.className = 'dn-modal du-okno';
  element.dataset['okno'] = 'okno-ustawien';

  const odswiez = przycisk('Odczytaj sekcję ponownie', 'dn-btn dn-btn--zarys');
  const zamknij = przycisk('Zamknij', 'dn-btn dn-btn--atrament');

  const stopka = document.createElement('footer');
  stopka.className = 'dn-modal-stopka du-okno__stopka';
  stopka.append(odswiez, zamknij);

  element.append(naglowek(() => element.close()), cialo, stopka);

  function pokazSekcje(kod: KodSekcjiUstawien): void {
    wnetrzeCiala.replaceChildren(sekcjaCzynna().element);
    nawigacja.oznacz(kod);
    sekcjaCzynna().odswiez();
  }

  nawigacja.odswiez(REJESTR_SEKCJI_USTAWIEN, stan.czynna());
  stan.naZmiane((kod) => pokazSekcje(kod));

  odswiez.addEventListener('click', () => sekcjaCzynna().odswiez());
  zamknij.addEventListener('click', () => element.close());

  return {
    element,

    otworz() {
      if (!element.isConnected) document.body.append(element);
      if (!element.open) element.showModal();
      wnetrzeCiala.replaceChildren(sekcjaCzynna().element);
      nawigacja.oznacz(stan.czynna());
      sekcjaCzynna().odswiez();
    },

    zamknij: () => element.close(),

    rozlacz() {
      for (const sekcja of sekcje.values()) sekcja.rozlacz();
      element.remove();
    },
  };
}

/** Nawigacja sekcji: kolumna boczna zbudowana z rejestru stałego, z metodami odświeżenia wykazu i oznaczenia pozycji czynnej. */
interface NawigacjaSekcji {
  element: HTMLElement;
  odswiez(rejestr: readonly (typeof REJESTR_SEKCJI_USTAWIEN)[number][], czynna: KodSekcjiUstawien): void;
  oznacz(czynna: KodSekcjiUstawien): void;
}

function utworzNawigacjeSekcji(naWybor: (kod: KodSekcjiUstawien) => void): NawigacjaSekcji {
  const element = document.createElement('nav');
  element.className = 'dn-boczna du-nawigacja';
  element.setAttribute('aria-label', 'Sekcje Okna Ustawień');

  function oznacz(czynna: KodSekcjiUstawien): void {
    for (const pozycja of element.querySelectorAll<HTMLElement>('[data-sekcja]')) {
      const wybrana = pozycja.dataset['sekcja'] === czynna;
      if (wybrana) pozycja.setAttribute('aria-current', 'page');
      else pozycja.removeAttribute('aria-current');
    }
  }

  return {
    element,

    odswiez(rejestr, czynna) {
      const pozycje = rejestr.map((opis) => {
        const przycisk = document.createElement('button');
        przycisk.type = 'button';
        przycisk.className = 'dn-boczna-pozycja du-nawigacja__pozycja';
        przycisk.dataset['sekcja'] = opis.kod;

        const nazwa = document.createElement('span');
        nazwa.className = 'du-nawigacja__nazwa';
        nazwa.textContent = opis.nazwa;

        przycisk.append(elementIkony(opis.ikona, { rozmiar: 16 }), nazwa);
        przycisk.addEventListener('click', () => naWybor(opis.kod));
        return przycisk;
      });
      element.replaceChildren(...pozycje);
      oznacz(czynna);
    },

    oznacz,
  };
}

/** Nagłówek okna niesie ikonę, tytuł i przycisk zamknięcia, budując górny pasek ramy okna dialogowego ustawień. */
function naglowek(naZamkniecie: () => void): HTMLElement {
  const element = document.createElement('header');
  element.className = 'dn-modal-naglowek du-okno__naglowek';

  const tytul = document.createElement('h2');
  tytul.className = 'dn-modal-tytul';
  tytul.textContent = 'Ustawienia';

  const rozpychacz = document.createElement('span');
  rozpychacz.className = 'du-okno__rozpychacz';

  const zamkniecie = document.createElement('button');
  zamkniecie.type = 'button';
  zamkniecie.className = 'dn-btn-ikona';
  zamkniecie.setAttribute('aria-label', 'Zamknij okno ustawień');
  zamkniecie.title = 'Zamknij okno ustawień';
  zamkniecie.append(elementIkony('zamknij', { rozmiar: 18 }));
  zamkniecie.addEventListener('click', naZamkniecie);

  element.append(elementIkony('uzytkownik', { rozmiar: 20 }), tytul, rozpychacz, zamkniecie);
  return element;
}

function przycisk(tresc: string, klasa: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.textContent = tresc;
  return element;
}
