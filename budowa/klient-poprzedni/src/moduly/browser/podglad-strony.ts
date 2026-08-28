import type { BrowserSnapshot } from '../../../../shared/contract';

/**
 * Podgląd współdzielony — treść strony widoczna jednocześnie Operatorowi i modelowi
 * w oknie Browser Window. Pokazuje migawkę oraz obsługuje przewijanie treści
 * i zaznaczenie fragmentu, obie czynności miejscowo, na odczytanej migawce.
 */
export interface PodgladStrony {
  element: HTMLElement;
  /** Nanosi migawkę; `null` zostawia podgląd pusty. */
  odswiez(migawka: BrowserSnapshot | null): void;
  /** Tryb czytnika: sam tekst, bez źródła strony i metadanych. */
  ustawTrybCzytnika(wlaczony: boolean): void;
  /** Przewija treść o ekran: -1 w górę, 1 w dół. */
  przewin(kierunek: -1 | 1): void;
  /** Przewija treść do wskazanego wiersza; numer liczony od 1. */
  pokazWiersz(numer: number): void;
}

export function utworzPodgladStrony(naZaznaczenie: (fragment: string) => void): PodgladStrony {
  const naglowek = document.createElement('p');
  naglowek.className = 'mb-podglad__naglowek';

  const adres = document.createElement('p');
  adres.className = 'mb-podglad__adres';

  const tekst = document.createElement('pre');
  tekst.className = 'mb-podglad__tekst';
  tekst.tabIndex = 0;
  tekst.setAttribute('aria-label', 'Treść renderowana strony — obszar zaznaczenia');

  const zrodlo = document.createElement('pre');
  zrodlo.className = 'mb-podglad__zrodlo';
  zrodlo.hidden = true;

  const element = document.createElement('div');
  element.className = 'mb-podglad';
  element.dataset['czytnik'] = 'nie';
  element.append(naglowek, adres, tekst, zrodlo);

  let trybCzytnika = false;

  /** Zaznaczenie zbierane po każdym geście, który mógł je zmienić. */
  function zbierzZaznaczenie(): void {
    const wybor = document.getSelection();
    const fragment = wybor === null ? '' : wybor.toString().trim();
    if (fragment !== '' && !tekst.contains(wybor?.anchorNode ?? null)) return;
    naZaznaczenie(fragment);
  }

  tekst.addEventListener('mouseup', zbierzZaznaczenie);
  tekst.addEventListener('keyup', zbierzZaznaczenie);

  return {
    element,

    odswiez(migawka) {
      if (migawka === null) {
        naglowek.textContent = '';
        adres.textContent = '';
        tekst.textContent = '';
        zrodlo.textContent = '';
        zrodlo.hidden = true;
        return;
      }
      naglowek.textContent = (migawka.title ?? '').trim() === '' ? migawka.url : `${migawka.title}`;
      adres.textContent = opisMigawki(migawka);
      tekst.textContent = migawka.text ?? '';
      zrodlo.textContent = migawka.html ?? '';
      zrodlo.hidden = trybCzytnika || (migawka.html ?? '') === '';
      adres.hidden = trybCzytnika;
    },

    ustawTrybCzytnika(wlaczony) {
      trybCzytnika = wlaczony;
      element.dataset['czytnik'] = wlaczony ? 'tak' : 'nie';
      // Źródło strony znika z widoku, ale nie z migawki; powrót nie kosztuje odczytu.
      zrodlo.hidden = wlaczony || (zrodlo.textContent ?? '') === '';
      adres.hidden = wlaczony;
    },

    przewin(kierunek) {
      tekst.scrollBy({ top: kierunek * Math.max(tekst.clientHeight - 24, 120) });
    },

    // Wiersz wskazuje się udziałem w treści, a nie pomiarem wysokości linii.
    pokazWiersz(numer) {
      const wierszy = (tekst.textContent ?? '').split('\n').length;
      if (wierszy <= 1) return;
      const udzial = Math.min(Math.max(numer - 1, 0), wierszy - 1) / wierszy;
      tekst.scrollTo({ top: udzial * tekst.scrollHeight });
    },
  };
}

/**
 * Wiersz metadanych migawki: adres, chwila pobrania oraz to, co migawka niesie
 * poza samym tekstem. Wiersz mówi wprost o obecności źródła strony i zrzutu, więc
 * Operator wie, czego podgląd nie pokazuje, zamiast to wnioskować.
 */
function opisMigawki(migawka: BrowserSnapshot): string {
  const czesci = [migawka.url, new Date(migawka.capturedAt).toLocaleString('pl')];
  if ((migawka.html ?? '') !== '') czesci.push('ze źródłem strony');
  if ((migawka.screenshotRef ?? '') !== '') czesci.push(`zrzut: ${migawka.screenshotRef ?? ''}`);
  return czesci.join(' · ');
}
