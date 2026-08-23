import type { StudioProvenance } from '../../../../shared/contract';
import {
  poleLogiczne,
  poleWielowierszowe,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';

/**
 * Panel treści fragmentu, postaci w całości i pochodzenia.
 *
 * ── Dlaczego te pięć czynności stoi razem ───────────────────────────────────
 * `studio.text.get`, `studio.text.edit`, `studio.document.form.get`,
 * `studio.document.form.save` i `studio.provenance.list` dotyczą jednego: treści
 * i postaci jako CAŁOŚCI, a nie pojedynczej cechy. Rozsypane po trzech panelach
 * nastaw byłyby nie do znalezienia.
 *
 * ── Dlaczego zapis postaci nie żąda treści ──────────────────────────────────
 * `studio.document.form.save` przyjmuje treść jako pole opcjonalne: brak znaczy
 * „bez zmiany treści". Ten przycisk odczytuje postać i zapisuje ją z powrotem,
 * zakładając wersję — czyli utrwala to, co nastawy strony i style zmieniły, bez
 * dotykania ani jednej litery. To jest droga, którą postać przestaje ginąć nawet
 * wtedy, gdy Operator nie pisał.
 *
 * ── Dlaczego brzmienie fragmentu ma pole wielowierszowe ─────────────────────
 * Bo poprawa fragmentu bez przepisywania całości jest osią zamówienia, a fragment
 * bywa akapitem, nie wyrazem. Pole jednowierszowe wymuszałoby wklejanie akapitu
 * w linijkę wysokości jednego wiersza.
 *
 * Panel nie woła rdzenia i nie zna dokumentu — zleca czynności warstwie wyżej.
 */

/** Czynności panelu treści i postaci. */
export interface CzynnosciTresciPanelu {
  /** Odczyt treści zaznaczonego fragmentu wraz z jego postacią i blokadami. */
  naOdczytFragmentu(): void;
  /** Nowe brzmienie zaznaczonego fragmentu. */
  naZmianeBrzmienia(brzmienie: string, zachowajPostac: boolean): void;
  /** Odczyt pełnej postaci dokumentu. */
  naOdczytPostaci(): void;
  /** Zapis postaci wraz z założeniem wersji, bez zmiany treści. */
  naZapisPostaci(): void;
  /** Wykaz pochodzenia fragmentów dokumentu. */
  naPochodzenia(): void;
}

/** Panel treści i postaci wraz z jego sterowaniem. */
export interface StronaPanelTresci {
  element: HTMLElement;
  przestawWidocznosc(): void;
  widoczny(): boolean;
  /** Treść odczytanego fragmentu wchodzi w pole brzmienia jako punkt wyjścia. */
  pokazFragment(tresc: string, zdanieOBlokadach: string): void;
  /** Zestawienie postaci odczytanej — czego dokument ma ile. */
  pokazPostac(zdanie: string): void;
  /** Wykaz pochodzenia fragmentów. */
  pokazPochodzenia(zapisy: readonly StudioProvenance[]): void;
  pokazOdpowiedz(tresc: string, powodzenie: boolean): void;
}

export function utworzStronePanelTresci(
  czynnosci: CzynnosciTresciPanelu,
): StronaPanelTresci {
  const odpowiedz = utworzWierszOdpowiedzi();

  const brzmienie = poleWielowierszowe(
    {
      etykieta: 'Brzmienie zaznaczonego fragmentu',
      opis:
        'Odczyt wpisuje tu treść fragmentu; zmiana wysyła ją z powrotem. Pole puste znaczy ' +
        'USUNIĘCIE fragmentu — tak stanowi kontrakt i tak ma zostać, bo usunięcie fragmentu jest ' +
        'czynnością, nie brakiem wskazania.',
    },
    6,
  );

  const zachowajPostac = poleLogiczne({
    etykieta: 'Zachowaj postać znaku fragmentu zastanego',
    opis: 'Wyłączone znaczy, że nowe brzmienie bierze postać domyślną miejsca.',
  });
  zachowajPostac.kontrolka.checked = true;

  const zdanieBlokad = document.createElement('p');
  zdanieBlokad.className = 'dn-pole-opis';

  const odczytaj = przycisk('Odczytaj zaznaczony fragment', 'dn-btn dn-btn--sm dn-btn--zarys');
  odczytaj.dataset['czynnosc'] = 'odczytaj-fragment';
  odczytaj.title =
    'Idzie komendą studio.text.get — oddaje treść fragmentu wraz z jego postacią znaku, postacią ' +
    'akapitu i blokadami, które go obejmują.';
  odczytaj.addEventListener('click', () => czynnosci.naOdczytFragmentu());

  const zmien = przycisk('Zmień brzmienie fragmentu', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zmien.dataset['czynnosc'] = 'zmien-brzmienie';
  zmien.title =
    'Idzie komendą studio.text.edit. Blokada fragmentu jest sprawdzana W RDZENIU, przed ' +
    'dotknięciem treści — okno jej nie pilnuje i pilnować nie może.';
  zmien.addEventListener('click', () => {
    czynnosci.naZmianeBrzmienia(brzmienie.kontrolka.value, zachowajPostac.kontrolka.checked);
  });

  const zdaniePostaci = document.createElement('p');
  zdaniePostaci.className = 'dn-pole-opis';

  const odczytajPostac = przycisk('Odczytaj postać dokumentu', 'dn-btn dn-btn--sm dn-btn--zarys');
  odczytajPostac.dataset['czynnosc'] = 'odczytaj-postac';
  odczytajPostac.title =
    'Idzie komendą studio.document.form.get — arkusz stylów, nastawy strony, sekcje, bloki, ' +
    'tabele, obiekty, aparat dokumentu, pola i blokady.';
  odczytajPostac.addEventListener('click', () => czynnosci.naOdczytPostaci());

  const zapiszPostac = przycisk(
    'Utrwal postać i załóż wersję',
    'dn-btn dn-btn--sm dn-btn--sygnal',
  );
  zapiszPostac.dataset['czynnosc'] = 'utrwal-postac';
  zapiszPostac.title =
    'Idzie komendą studio.document.form.save bez pola treści — utrwala to, co nastawy strony ' +
    'i style zmieniły, nie dotykając ani jednej litery. To jest droga, którą postać przestaje ' +
    'ginąć także wtedy, gdy Operator nie pisał.';
  zapiszPostac.addEventListener('click', () => czynnosci.naZapisPostaci());

  const wykazPochodzen = document.createElement('ul');
  wykazPochodzen.className = 'ms-postac__wykaz';

  const pochodzenia = przycisk('Pochodzenie fragmentów', 'dn-btn dn-btn--sm dn-btn--zarys');
  pochodzenia.dataset['czynnosc'] = 'odczytaj-pochodzenia';
  pochodzenia.title =
    'Idzie komendą studio.provenance.list — skąd pochodzi każdy wniesiony fragment. To jest ' +
    'podstawa pod bibliografię i pod podobieństwa w panelu Redaktora.';
  pochodzenia.addEventListener('click', () => czynnosci.naPochodzenia());

  const element = document.createElement('section');
  element.className = 'ms-postac ms-postac--tresc';
  element.dataset['panel'] = 'tresc-i-postac';
  element.hidden = true;
  element.setAttribute('aria-label', 'Treść fragmentu, postać dokumentu i pochodzenie');

  element.append(
    grupa('Fragment', [odczytaj, zdanieBlokad, brzmienie.element, zachowajPostac.element, zmien]),
    grupa('Postać dokumentu', [odczytajPostac, zdaniePostaci, zapiszPostac]),
    grupa('Pochodzenie', [pochodzenia, wykazPochodzen]),
    odpowiedz.element,
  );

  return {
    element,

    przestawWidocznosc() {
      element.hidden = !element.hidden;
    },

    widoczny: () => !element.hidden,

    pokazFragment(tresc, zdanieOBlokadach) {
      brzmienie.kontrolka.value = tresc;
      zdanieBlokad.textContent = zdanieOBlokadach;
    },

    pokazPostac(zdanie) {
      zdaniePostaci.textContent = zdanie;
    },

    pokazPochodzenia(zapisy) {
      wykazPochodzen.replaceChildren(
        ...zapisy.map((zapis) => {
          const pozycja = document.createElement('li');
          pozycja.className = 'ms-postac__zasada';
          pozycja.dataset['pochodzenie'] = zapis.id;
          const skad =
            zapis.sourceUrl ?? zapis.libraryFileId ?? zapis.sourceTitle ?? 'źródło nienazwane';
          pozycja.textContent =
            `znaki ${zapis.rangeStart}–${zapis.rangeEnd} · ${zapis.kind} · ${skad}` +
            `${zapis.sourceVersion === undefined ? '' : ` · wersja ${zapis.sourceVersion}`}` +
            `${zapis.author === undefined ? '' : ` · wniósł ${zapis.author}`}`;
          return pozycja;
        }),
      );
      if (zapisy.length === 0) {
        const puste = document.createElement('li');
        puste.className = 'dn-pole-opis';
        puste.textContent =
          'Rdzeń nie ma ani jednego zapisu pochodzenia dla tego dokumentu. Zapisy powstają przy ' +
          'wnoszeniu fragmentów ze stron i z Biblioteki — dokument pisany od zera ich nie ma i to ' +
          'jest stan prawidłowy, nie brak.';
        wykazPochodzen.append(puste);
      }
    },

    pokazOdpowiedz: (tresc, powodzenie) => odpowiedz.pokaz(tresc, powodzenie),
  };
}

/** Grupa pól panelu. */
function grupa(tytul: string, zawartosc: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'ms-postac__tytul';
  naglowek.textContent = tytul;

  const element = document.createElement('div');
  element.className = 'ms-postac__grupa';
  element.append(naglowek, ...zawartosc);
  return element;
}
