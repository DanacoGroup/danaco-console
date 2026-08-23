import { oznaczOknoAplikacji } from '../komponenty/okno-aplikacji';
import './ustawienia.css';

import { elementIkony } from '../ikony/ikony';
import type { Kanal } from '../protokol/kanal';
import { REJESTR_SEKCJI_USTAWIEN, type SekcjaUstawien, type KodSekcjiUstawien } from './sekcje';
import { utworzStanUstawien, type StanUstawien } from './stan-ustawien';

/**
 * Okno Ustawień — rama: nawigacja sekcji i montaż sekcji czynnej.
 *
 * Okno stoi na natywnym `<dialog>`, tak jak okno konfiguracji i okno modeli:
 * warstwę tła, stos okien i zamknięcie klawiszem Esc daje przeglądarka, a nie
 * własna nakładka. Wygląd bierze z biblioteki (`dn-modal`, `dn-boczna`) —
 * plik nie zna ani jednej barwy.
 *
 * Rama nie niesie treści sekcji. Zna wyłącznie rejestr z `sekcje.ts` (kod,
 * nazwa, ikona, fabryka) i kontrakt `SekcjaUstawien` — co sekcja niesie
 * w środku, jest sprawą pliku sekcji, nie ramy.
 *
 * Rejestr niesie dziś jedną sekcję („uwierzytelnianie"); powód zwężenia
 * i warunek przywrócenia stoją w nagłówku `sekcje.ts`, a ta rama tylko na
 * niego wskazuje.
 *
 * Sekcje budują się razem, przy otwarciu okna. Rejestr sekcji jest stały (nie
 * przychodzi katalogiem rdzenia jak kategorie konfiguracji), więc nie ma
 * powodu budować ich leniwie przy pierwszym kliknięciu — a budowa od razu
 * utrzymuje stan każdej sekcji, na przykład wpisany, jeszcze niezapisany
 * formularz, przy przełączaniu zakładek, zamiast go gubić. Kod poniżej działa
 * dla jednej sekcji identycznie jak dla wielu: dołożenie drugiej pozycji do
 * `REJESTR_SEKCJI_USTAWIEN` nie wymaga zmiany tego pliku.
 *
 * Okno otwiera się natychmiast; każda sekcja sama rozstrzyga swój odczyt i swój
 * stan błędu. „Odczytaj ponownie" w stopce odświeża wyłącznie sekcję czynną —
 * sekcje niewidoczne nie ciągną rdzenia w tle bez powodu.
 */
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

  /**
   * Kolumna nawigacji jest budowana zawsze, ale osadzana w dokumencie tylko
   * wtedy, gdy rejestr niesie więcej niż jedną sekcję.
   *
   * Nawigacja do jednego miejsca jest nawigacją donikąd: kolumna
   * z pojedynczym, zawsze czynnym przyciskiem sugerowałaby Operatorowi wybór,
   * którego nie ma, i byłaby kłamstwem o kształcie okna.
   *
   * Próg jest samoczynny, a nie ręcznym przełącznikiem — warunek czyta długość
   * `REJESTR_SEKCJI_USTAWIEN` przy każdej budowie okna, więc po dołożeniu
   * drugiej sekcji kolumna wraca sama, bez zmiany w tym pliku. Gałąź „więcej
   * niż jedna" jest mechanizmem wzrostu okna, a nie martwym kodem.
   */
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

/** Nawigacja sekcji: kolumna boczna zbudowana z rejestru stałego. */
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

/** Nagłówek okna: ikona, tytuł, przycisk zamknięcia. */
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
