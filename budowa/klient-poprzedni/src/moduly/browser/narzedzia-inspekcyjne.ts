import type { BrowserSnapshot } from '../../../../shared/contract';
import { utworzDymekObjasnienia } from '../../komponenty/dymek';
import { KLASY_DYMKA, OBJASNIENIA } from './etykiety-browser';
import { KLASA_PRZYCISKU, przyciskCzynnosci } from './przyciski-browser';
import type { StanPrzegladania } from './stan-przegladania';

/**
 * Narzędzia inspekcyjne Browser Window — panel warstwy czwartej, otwierany
 * skrótem albo z paska kontekstu w trybie administracyjnym.
 */
export interface NarzedziaInspekcyjne {
  element: HTMLElement;
}

export function utworzNarzedziaInspekcyjne(
  stan: StanPrzegladania,
  powiedz: (tresc: string, powodzenie: boolean) => void,
): NarzedziaInspekcyjne {
  /** Migawka odłożona do porównania różnic; wartość pusta znaczy „nie ma z czym porównywać”. */
  let odlozona: BrowserSnapshot | null = null;

  const naglowek = document.createElement('h4');
  naglowek.className = 'mb-panel__tytul';
  naglowek.textContent = 'Narzędzia inspekcyjne';

  const wynik = document.createElement('pre');
  wynik.className = 'mb-inspekcja__wynik';

  async function pokazZrodlo(): Promise<void> {
    powiedz('Odczyt źródła strony z rdzenia (migawka z polem includeHtml)…', true);
    await stan.zaciagnijMigawke(true);
    const migawka = stan.migawka();
    if (migawka === null) {
      wynik.textContent = '';
      powiedz(stan.powodMigawki(), false);
      return;
    }
    const zrodlo = migawka.html ?? '';
    wynik.textContent = zrodlo === '' ? '' : zrodlo;
    powiedz(
      zrodlo === ''
        ? 'Rdzeń oddał migawkę BEZ źródła strony (pole html puste) — nie ma czego pokazać.'
        : `Źródło strony ${migawka.url}: ${zrodlo.length} znaków.`,
      zrodlo !== '',
    );
  }

  function odloz(): void {
    const migawka = stan.migawka();
    if (migawka === null) {
      powiedz('Nie ma czego odłożyć — najpierw przejdź do strony w Browser Window.', false);
      return;
    }
    odlozona = migawka;
    powiedz(
      `Migawka ${migawka.id} odłożona do porównania (${(migawka.text ?? '').length} znaków treści).`,
      true,
    );
  }

  function porownaj(): void {
    const biezaca = stan.migawka();
    if (odlozona === null || biezaca === null) {
      powiedz('Porównanie potrzebuje dwóch migawek — odłóż jedną, potem odśwież stronę.', false);
      return;
    }
    const zestawienie = zestawWierszy(odlozona, biezaca);
    wynik.textContent = zestawienie.opis;
    powiedz(zestawienie.zdanie, true);
  }

  const pasek = document.createElement('div');
  pasek.className = 'mb-panel__pasek';
  pasek.append(
    przyciskCzynnosci('Pokaż źródło strony', KLASA_PRZYCISKU.zarys, () => void pokazZrodlo()),
    przyciskCzynnosci('Odłóż migawkę do porównania', KLASA_PRZYCISKU.zarys, odloz),
    przyciskCzynnosci('Porównaj z bieżącą', KLASA_PRZYCISKU.zarys, porownaj),
    utworzDymekObjasnienia(OBJASNIENIA.inspekcja, KLASY_DYMKA),
  );
  // Cztery narzędzia mają uchwyty w rdzeniu; czynności stoją w panelu rodzin, przy właściwej stronie.

  const element = document.createElement('section');
  element.className = 'mb-panel mb-inspekcja';
  element.setAttribute('aria-label', 'Narzędzia inspekcyjne strony — warstwa czwarta');
  element.append(naglowek, pasek, wynik);

  return { element };
}

/** Wynik zestawienia dwóch migawek strony: zdanie dla wiersza odpowiedzi wraz z pełną treścią panelu różnic. */
interface ZestawienieMigawek {
  zdanie: string;
  opis: string;
}

/**
 * Zestawienie treści dwóch migawek po wierszach. Miary są trzy: pierwszy wiersz
 * różniący się wskazuje miejsce zmiany, liczba wierszy dodanych i zdjętych
 * mówi o jej rozmiarze, a równość treści rozstrzyga, czy zmiana w ogóle była.
 */
function zestawWierszy(
  wczesniejsza: BrowserSnapshot,
  biezaca: BrowserSnapshot,
): ZestawienieMigawek {
  const przed = (wczesniejsza.text ?? '').split('\n');
  const po = (biezaca.text ?? '').split('\n');
  if (przed.join('\n') === po.join('\n')) {
    return {
      zdanie: `Treść bez zmian względem migawki ${wczesniejsza.id}.`,
      opis: `Migawki ${wczesniejsza.id} i ${biezaca.id} niosą tę samą treść (${przed.length} wierszy).`,
    };
  }

  const pierwszaRoznica = numerPierwszejRoznicy(przed, po);
  const zdjete = przed.filter((linia) => !po.includes(linia));
  const dodane = po.filter((linia) => !przed.includes(linia));
  return {
    zdanie:
      `Treść zmieniona: pierwsza różnica w wierszu ${pierwszaRoznica}, ` +
      `wierszy dodanych ${dodane.length}, zdjętych ${zdjete.length}.`,
    opis: [
      `Migawka wcześniejsza: ${wczesniejsza.id} (${przed.length} wierszy)`,
      `Migawka bieżąca: ${biezaca.id} (${po.length} wierszy)`,
      '',
      ...dodane.slice(0, GRANICA_WIERSZY).map((linia) => `+ ${linia}`),
      ...zdjete.slice(0, GRANICA_WIERSZY).map((linia) => `- ${linia}`),
    ].join('\n'),
  };
}

/**
 * Ile wierszy różnicy panel wypisuje.
 *
 * Zestawienie ma pokazać, co się zmieniło, a nie przepisać stronę: przy zmianie
 * obejmującej cały dokument pełny wypis byłby drugą kopią treści w oknie.
 */
const GRANICA_WIERSZY = 40;

/** Numer pierwszego wiersza, w którym treści migawek zaczynają się rozchodzić; liczony od jedynki, nie od zera. */
function numerPierwszejRoznicy(przed: readonly string[], po: readonly string[]): number {
  const wspolne = Math.min(przed.length, po.length);
  for (let numer = 0; numer < wspolne; numer += 1) {
    if (przed[numer] !== po[numer]) return numer + 1;
  }
  return wspolne + 1;
}
