import { IdentityMode, type IdentityCategory } from '../../../shared/contract';
import {
  poleWielowierszowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from './kontrolki-formularza';
import type { StanTozsamosci } from './stan-tozsamosci';
import {
  NAZWY_TRYBOW,
  TRYBY,
  nazwaWarstwy,
  ostrzezenieTrybu,
  trybZWyboru,
} from './warstwy-tozsamosci';

/**
 * Edytor jednej kategorii zasad — treść, tryb podania i ostrzeżenie o skutku.
 *
 * Przełącznik trybu rozstrzyga, czy prompt fabryczny producenta trafi do modelu,
 * więc ostrzeżenie stoi przy nim na stałe i zmienia się razem z nim.
 * Tryb bierze się kolejno z zapisu i z trybu proponowanego przez kategorię, a gdy
 * katalog nie podaje żadnego — z wartości domyślnej klucza `tozsamosc.tryb_domyslny`.
 * Kategoria bez zapisu na osi czynnej bierze treść z osi szerszej i edytor podaje
 * to wprost, zamiast pokazywać puste pole bez wyjaśnienia.
 */
export interface EdytorTozsamosci {
  /** Edytor osadzany w panelu tożsamości. */
  element: HTMLElement;
  /** Nanosi stan; treść wypełnia dopiero przy zmianie kategorii albo osi. */
  odswiez(): void;
}

export function utworzEdytorTozsamosci(stan: StanTozsamosci): EdytorTozsamosci {
  const tytul = document.createElement('h3');
  tytul.className = 'dm-edytor__tytul';

  const opis = document.createElement('p');
  opis.className = 'dn-pole-opis dm-edytor__opis';

  const tryb = poleWyboru(
    { etykieta: 'Tryb podania nakładki' },
    TRYBY.map((wartosc) => ({ wartosc, etykieta: NAZWY_TRYBOW[wartosc] })),
  );

  const ostrzezenie = document.createElement('p');
  ostrzezenie.className = 'dm-edytor__ostrzezenie';
  ostrzezenie.setAttribute('role', 'note');

  const TRESC_DEMONSTRACYJNA = 'Treść demo do późniejszej konfiguracji';

  const tresc = poleWielowierszowe(
    {
      etykieta: 'Treść kategorii',
      podpowiedz: 'treść, która wejdzie do nakładki tej kategorii',
    },
    16,
  );

  const pochodzenie = document.createElement('p');
  pochodzenie.className = 'dn-pole-opis dm-edytor__pochodzenie';

  const zapis = przycisk('Zapisz treść kategorii', 'dn-btn dn-btn--atrament');
  const zdejmij = przycisk('Zdejmij zapis z tej osi', 'dn-btn dn-btn--niebezpieczny');
  const odpowiedz = utworzWierszOdpowiedzi();

  const przyciski = document.createElement('div');
  przyciski.className = 'dm-edytor__przyciski';
  przyciski.append(zapis, zdejmij);

  const element = document.createElement('section');
  element.className = 'dm-edytor';
  element.append(
    tytul,
    opis,
    tryb.element,
    ostrzezenie,
    tresc.element,
    pochodzenie,
    przyciski,
    odpowiedz.element,
  );

  /**
   * Klucz wypełnienia: kategoria wraz z osią. Zmiana któregokolwiek członu
   * znaczy inną treść w edytorze; samo przeliczenie stanu — nie znaczy, więc
   * nie kasuje tego, co Operator właśnie pisze.
   */
  let wypelniony = '';

  function kluczWypelnienia(kategoria: IdentityCategory | null): string {
    const wskazanie = stan.wskazanie();
    return `${kategoria?.id ?? ''}|${wskazanie.os}|${wskazanie.bytOsi}`;
  }

  function ubierzOstrzezenie(): void {
    const wybrany = trybZWyboru(tryb.kontrolka.value);
    ostrzezenie.textContent = ostrzezenieTrybu(wybrany);
    ostrzezenie.dataset.tryb = wybrany;
    ostrzezenie.dataset.zastapienie = String(wybrany === IdentityMode.ZASTAP);
  }

  function wypelnij(kategoria: IdentityCategory | null): void {
    const zapisany = kategoria === null ? null : stan.dokument(kategoria.id);

    tytul.textContent =
      kategoria === null ? 'Kategorie zasad' : `${kategoria.name} · ${nazwaWarstwy(kategoria.layer)}`;
    opis.textContent =
      kategoria === null
        ? 'Wybierz kategorię z wykazu obok, aby zobaczyć i zmienić jej treść.'
        : (kategoria.description ?? '');
    opis.hidden = opis.textContent === '';

    tryb.kontrolka.value = zapisany?.mode ?? kategoria?.defaultMode ?? IdentityMode.ZASTAP;
    // Kategoria bez własnego zapisu na tej osi dostaje treść demonstracyjną
    // jako zawartość pola; w bazie nie ma jej, dopóki nie padnie „Zapisz treść
    // kategorii". Jedyną drogą zapisu pozostaje komenda `identity.document.set`.
    tresc.kontrolka.value = zapisany?.content ?? TRESC_DEMONSTRACYJNA;
    pochodzenie.textContent =
      kategoria === null
        ? ''
        : zapisany === null
          ? 'Ta oś nie ma własnego zapisu tej kategorii. Obowiązuje treść z osi szerszej; zapis tutaj przykryje ją dla wskazanego bytu.'
          : `Zapis własny tej osi. Odcisk treści: ${zapisany.contentHash ?? 'nie podano'}.`;

    tryb.element.hidden = kategoria === null;
    tresc.element.hidden = kategoria === null;
    przyciski.hidden = kategoria === null;
    zdejmij.hidden = zapisany === null;
    ostrzezenie.hidden = kategoria === null;

    odpowiedz.wyczysc();
    ubierzOstrzezenie();
  }

  tryb.kontrolka.addEventListener('change', ubierzOstrzezenie);

  zapis.addEventListener('click', () => {
    const kategoria = stan.wybrana();
    if (kategoria === null) return;
    void stan
      .zapisz(kategoria.id, tresc.kontrolka.value, trybZWyboru(tryb.kontrolka.value))
      .then((wynik) => {
        odpowiedz.pokaz(
          wynik.udany
            ? 'Treść zapisana. Podgląd złożonego promptu został przeliczony.'
            : `Rdzeń nie przyjął zapisu: ${wynik.blad?.message ?? 'brak treści błędu'}`,
          wynik.udany,
        );
      });
  });

  zdejmij.addEventListener('click', () => {
    const kategoria = stan.wybrana();
    const zapisany = kategoria === null ? null : stan.dokument(kategoria.id);
    if (zapisany === null) return;
    void stan.usun(zapisany.id).then((wynik) => {
      odpowiedz.pokaz(
        wynik.udany
          ? 'Zapis zdjęty. Obowiązuje treść tej kategorii z osi szerszej.'
          : `Rdzeń nie zdjął zapisu: ${wynik.blad?.message ?? 'brak treści błędu'}`,
        wynik.udany,
      );
    });
  });

  wypelnij(null);

  return {
    element,

    odswiez() {
      const kategoria = stan.wybrana();
      const klucz = kluczWypelnienia(kategoria);
      if (klucz !== wypelniony) {
        wypelniony = klucz;
        wypelnij(kategoria);
        return;
      }
      // Ta sama kategoria i ta sama oś: treści nie ruszamy, lecz przycisk
      // zdjęcia zapisu musi nadążyć za tym, czy zapis nadal istnieje.
      zdejmij.hidden = kategoria === null || stan.dokument(kategoria.id) === null;
    },
  };
}
