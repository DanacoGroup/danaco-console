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
 * Reading View jest oknem lektury i ekstrakcji, wczytującym treść dokumentu repozytorium i zamieniającym zaznaczony fragment w ustalenie powiązane ze źródłem.
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

    // Wskazanie nowego źródła zeruje stronę i pociąga treść; warunek chroni przed pętlą odczytu.
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

/** Funkcja tworzy odnośnik do podglądu graficznego strony, ponieważ rdzeń oddaje wyłącznie wskazanie zasobu, a nie samą grafikę. */
function odnosnikGraficzny(wskazanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis mr-lektura__odnosnik';
  element.textContent = `Podgląd graficzny strony: ${wskazanie}. Rdzeń oddaje wskazanie zasobu, a nie samą grafikę — z tej postaci cytatu nie da się zaznaczyć.`;
  return element;
}

/** Funkcja układa zdanie nagłówka opisujące tytuł czytanego materiału oraz numer bieżącej i łącznej liczby stron. */
function opisMaterialu(tytul: string, strona: number, podglad: LibraryPreview | null): string {
  if (tytul === '') return 'Nie wskazano materiału do lektury.';
  const stron = podglad?.pageCount;
  return stron === undefined
    ? `Materiał: ${tytul} — strona ${String(strona)}.`
    : `Materiał: ${tytul} — strona ${String(strona)} z ${String(stron)}.`;
}
