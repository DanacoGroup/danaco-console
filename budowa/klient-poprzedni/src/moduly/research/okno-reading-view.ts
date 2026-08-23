import type { LibraryPreview } from '../../../../shared/contract';
import { przycisk, utworzWierszOdpowiedzi } from '../../modele/kontrolki-formularza';
import { AKCJE_LEKTURY } from './akcje-okien';
import {
  ustawStanLektury,
  wczytajZrodlo,
  wykonajAkcjeLektury,
  wskazaneZrodlo,
  zapiszWypis,
  type KontekstLektury,
} from './czynnosci-lektury';
import { KODY_OKIEN } from './kody-okien';
import { utworzRameBadania } from './rama-badania';
import type { StanBadania } from './stan-badania';
import { utworzStanOknaBadania } from './stan-okna-badania';

/**
 * Reading View — okno lektury i ekstrakcji.
 *
 * Okno szóste opracowania (rozdz. 3.6), którego katalog okien rdzenia jeszcze
 * nie zna — patrz `kody-okien.ts`. Zbudowane, bo dwie jego czynności mają dziś
 * pokrycie: wczytanie treści dokumentu repozytorium i zamiana zaznaczonego
 * fragmentu w ustalenie z powiązaniem do czytanego źródła.
 *
 * Czego tu jeszcze nie ma i dlaczego: podświetlenia trwałe, notatki na
 * marginesie i wypisy zbiorcze mają już w kontrakcie własne komendy i własny byt
 * wraz z kotwicą pozycji, ale rdzeń nie ma dla nich uchwytu. Okno ich nie udaje —
 * zaznaczenie jest zaznaczeniem przeglądarki, a trwałym staje się dopiero jako
 * ustalenie. Pozostałe operacje na źródle stoją w panelu akcji pod nazwami
 * swoich komend i wracają odmową rdzenia, zamiast znikać z okna.
 *
 * Plik składa widok; zachowanie po naciśnięciu leży w `czynnosci-lektury`.
 */
export interface OknoReadingView {
  element: HTMLElement;
  odswiez(): void;
}

export function utworzOknoReadingView(
  stan: StanBadania,
  przejdz: (kodOkna: string) => void,
): OknoReadingView {
  let strona = 1;
  let podglad: LibraryPreview | null = null;
  /** Źródło, dla którego treść jest już wczytana — zapora przed pętlą odczytu. */
  let wczytane = '';

  const naglowekMaterialu = document.createElement('p');
  naglowekMaterialu.className = 'mr-lektura__material';

  const tresc = document.createElement('div');
  tresc.className = 'mr-lektura__tresc';
  tresc.tabIndex = 0;
  tresc.setAttribute('role', 'document');
  tresc.setAttribute('aria-label', 'Treść czytanego źródła');

  const kontekst: KontekstLektury = {
    stan,
    okno: utworzStanOknaBadania(),
    odpowiedz: utworzWierszOdpowiedzi(),
    tresc,
    strona: () => strona,
    ustawStrone: (numer) => {
      strona = numer < 1 ? 1 : numer;
      void wczytajZrodlo(kontekst);
    },
    wchlonPodglad: (nowy) => {
      podglad = nowy;
      przerysujTresc();
      odswiez();
    },
    przejdz,
  };

  const doUstalenia = przycisk('→ ustalenie', 'dn-btn dn-btn--sm dn-btn--atrament');
  doUstalenia.addEventListener('click', () => void zapiszWypis(kontekst));

  const poprzednia = przycisk('◄ Strona poprzednia', 'dn-btn dn-btn--sm dn-btn--duch');
  poprzednia.addEventListener('click', () => kontekst.ustawStrone(strona - 1));

  const nastepna = przycisk('Strona następna ►', 'dn-btn dn-btn--sm dn-btn--duch');
  nastepna.addEventListener('click', () => kontekst.ustawStrone(strona + 1));

  const nawigacja = document.createElement('div');
  nawigacja.className = 'mr-lektura__nawigacja';
  nawigacja.append(poprzednia, nastepna);

  kontekst.okno.tresc.append(
    naglowekMaterialu,
    nawigacja,
    tresc,
    doUstalenia,
    kontekst.odpowiedz.element,
  );

  const rama = utworzRameBadania(
    KODY_OKIEN.lektura,
    'Reading View',
    'pomocnicze',
    AKCJE_LEKTURY,
    (akcja) => void wykonajAkcjeLektury(kontekst, akcja),
  );
  rama.cialo.append(kontekst.okno.element);

  /** Treść czytana wraz z odnośnikiem do podglądu graficznego, gdy rdzeń go oddał. */
  function przerysujTresc(): void {
    if (podglad === null) {
      tresc.replaceChildren();
      return;
    }
    const akapit = document.createElement('p');
    akapit.className = 'mr-lektura__akapit';
    akapit.textContent = podglad.text ?? '';

    const czesci: HTMLElement[] = [akapit];
    if ((podglad.imageRef ?? '') !== '') czesci.push(odnosnikGraficzny(podglad.imageRef ?? ''));
    tresc.replaceChildren(...czesci);
  }

  function odswiez(): void {
    const zrodlo = wskazaneZrodlo(stan);
    naglowekMaterialu.textContent = opisMaterialu(zrodlo?.title ?? '', strona, podglad);

    // Wskazanie nowego źródła zeruje stronę i pociąga jego treść. Warunek na
    // `wczytane` jest zaporą: `odswiez` biegnie z każdego ogłoszenia stanu,
    // a wczytanie ogłasza stan ponownie.
    const wskazane = zrodlo?.id ?? '';
    if (wskazane !== '' && wskazane !== wczytane) {
      wczytane = wskazane;
      strona = 1;
      void wczytajZrodlo(kontekst);
      return;
    }
    if (wskazane === '') {
      wczytane = '';
      podglad = null;
      przerysujTresc();
    }
    ustawStanLektury(kontekst, (podglad?.text ?? '') !== '');
  }

  odswiez();
  return { element: rama.element, odswiez };
}

/** Odnośnik do podglądu graficznego — rdzeń oddaje wskazanie, nie obraz. */
function odnosnikGraficzny(wskazanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis mr-lektura__odnosnik';
  element.textContent = `Podgląd graficzny strony: ${wskazanie}. Rdzeń oddaje wskazanie zasobu, a nie samą grafikę — z tej postaci cytatu nie da się zaznaczyć.`;
  return element;
}

/** Zdanie nagłówka: co jest czytane i na której stronie. */
function opisMaterialu(tytul: string, strona: number, podglad: LibraryPreview | null): string {
  if (tytul === '') return 'Nie wskazano materiału do lektury.';
  const stron = podglad?.pageCount;
  return stron === undefined
    ? `Materiał: ${tytul} — strona ${String(strona)}.`
    : `Materiał: ${tytul} — strona ${String(strona)} z ${String(stron)}.`;
}
