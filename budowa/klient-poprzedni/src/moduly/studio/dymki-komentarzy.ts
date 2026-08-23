import {
  StudioChangeDecision,
  type StudioComment,
  type StudioTrackedChange,
} from '../../../../shared/contract';
import type { PropozycjaZmiany } from './pola-stanu';
import { NAZWY_RODZAJOW, RodzajZnakowania } from './przybornik-wykaz';
import { nazwaAutora, opiszZmiane, zmianyOczekujace } from './zmiany-modelu';

/**
 * Margines dokumentu — TRZY RÓŻNE BYTY, nie jeden.
 *
 * ── Rozstrzygnięcie Właściciela ─────────────────────────────────────────────
 * Na marginesie stają trzy rzeczy i Operator ma **po wyglądzie** wiedzieć, na
 * którą patrzy:
 *
 *   — **komentarz** — mówi o fragmencie i NIE niesie brzmienia; treści nie
 *     zmienia ani teraz, ani po rozwiązaniu wątku;
 *   — **propozycja zmiany** — niesie brzmienie fragmentu, ale NIE weszła
 *     w treść; Operator ją przyjmuje, odrzuca albo poprawia;
 *   — **zmiana śledzona** — jest JUŻ w treści i czeka na decyzję.
 *
 * Zlanie ich w jedną kartę „uwaga modelu" zabrałoby Operatorowi rozróżnienie, po
 * którym poznaje, czy dokument już się zmienił, czy jeszcze nie. Dlatego każda
 * karta niesie `data-rodzaj`, własny nagłówek i własne zdanie o skutku decyzji,
 * a arkusz `przybornik-znakowania.css` daje im trzy różne obramowania.
 *
 * ── Dlaczego dymek, a nie sama lista ────────────────────────────────────────
 * Komentarz ma stać przy miejscu, nie w wykazie. Kotwica w treści
 * (`powierzchnia-dokumentu.ts`) wskazuje fragment, a karta ustawia się na jej
 * wysokości. Spis do przejścia stoi osobno, w przyborniku znakowania — jedno nie
 * zastępuje drugiego: dymek mówi „tu", spis mówi „ile jeszcze".
 *
 * ── Czego kontrakt nie ma ───────────────────────────────────────────────────
 * Usunięcia komentarza. `studio.comment.*` niesie dodanie, wykaz i rozwiązanie;
 * komendy usuwającej nie ma. Przycisk „Usuń" nie jest więc udawany — jego
 * miejsce zajmuje „Rozwiąż wątek" wraz ze zdaniem o tej różnicy.
 */

/** Czynności marginesu zlecane oknu. */
export interface CzynnosciDymkow {
  /** Odpowiada w wątku wskazanego komentarza. */
  naOdpowiedz(idWatku: string, tresc: string): void;
  /** Przestawia rozwiązanie wątku. */
  naRozwiazanie(idKomentarza: string, rozwiazany: boolean): void;
  /** Przenosi kursor do fragmentu, do którego znakowanie jest przypięte. */
  naKotwice(idKomentarza: string): void;
  /** Rozstrzyga propozycję stojącą na marginesie — w całości. */
  naDecyzjePropozycji(przyjmij: boolean): void;
  /** Rozstrzyga jedną zmianę śledzoną. */
  naDecyzjeZmiany(kodZmiany: string, przyjmij: boolean): void;
}

/** Co stoi na marginesie w tej chwili. */
export interface WpisyMarginesu {
  komentarze: readonly StudioComment[];
  /** Propozycja brzmienia czekająca na decyzję; `null`, gdy żadnej nie ma. */
  propozycja: PropozycjaZmiany | null;
  /** Zakres, którego propozycja dotyczy — propozycja nie nosi go sama. */
  zakresPropozycji: { poczatek: number; koniec: number } | null;
  /** Zmiany śledzone dokumentu; margines pokazuje oczekujące decyzji. */
  zmiany: readonly StudioTrackedChange[];
}

/** Kolumna marginesu wraz z jej odświeżeniem. */
export interface DymkiKomentarzy {
  element: HTMLElement;
  /**
   * Przerysowuje margines.
   *
   * `polozenie` oddaje wysokość kotwicy w punktach powierzchni; kotwica bez
   * położenia (znakowanie dotyczące całego dokumentu albo fragmentu, którego już
   * nie ma) stawia kartę w kolejności wykazu, a nie na zgadniętej wysokości.
   */
  pokaz(wpisy: WpisyMarginesu, polozenie: (kod: string) => number | null): void;
  /** Liczba wątków otwartych — do zdania paska stanu. */
  ile(): number;
}

export function utworzDymkiKomentarzy(czynnosci: CzynnosciDymkow): DymkiKomentarzy {
  let watki = 0;

  const element = document.createElement('div');
  element.className = 'ms-dymki';
  element.setAttribute('aria-label', 'Margines dokumentu — komentarze, propozycje i zmiany');

  function pokaz(wpisy: WpisyMarginesu, polozenie: (kod: string) => number | null): void {
    const nadrzedne = wpisy.komentarze.filter(
      (komentarz) => komentarz.parentCommentId === undefined,
    );
    watki = nadrzedne.filter((komentarz) => komentarz.resolved !== true).length;
    const oczekujace = zmianyOczekujace(wpisy.zmiany);

    if (nadrzedne.length === 0 && wpisy.propozycja === null && oczekujace.length === 0) {
      const puste = document.createElement('p');
      puste.className = 'dn-pole-opis';
      puste.textContent =
        'Margines jest pusty: nie ma ani komentarza, ani propozycji brzmienia, ani zmiany ' +
        'śledzonej czekającej na decyzję. Zaznacz fragment — pływak przy zaznaczeniu i przybornik ' +
        'znakowania zakładają wszystkie trzy.';
      element.replaceChildren(puste);
      return;
    }

    const karty: HTMLElement[] = [];

    for (const komentarz of nadrzedne) {
      const odpowiedzi = wpisy.komentarze.filter(
        (pozycja) => pozycja.parentCommentId === komentarz.id,
      );
      const karta = utworzDymek(komentarz, odpowiedzi, czynnosci);
      ustawWysokosc(karta, polozenie(komentarz.id));
      karty.push(karta);
    }

    if (wpisy.propozycja !== null && wpisy.propozycja.tresc !== '') {
      const karta = utworzKartePropozycji(wpisy.propozycja, wpisy.zakresPropozycji, czynnosci);
      ustawWysokosc(karta, null);
      karty.push(karta);
    }

    for (const zmiana of oczekujace) {
      const karta = utworzKarteZmiany(zmiana, czynnosci);
      ustawWysokosc(karta, polozenie(zmiana.id));
      karty.push(karta);
    }

    element.replaceChildren(...karty);
  }

  return { element, pokaz, ile: () => watki };
}

/** Stawia kartę na wysokości kotwicy albo w kolejności wykazu. */
function ustawWysokosc(karta: HTMLElement, wysokosc: number | null): void {
  if (wysokosc !== null) karta.style.top = `${wysokosc}px`;
  karta.dataset['zakotwiczony'] = wysokosc === null ? 'nie' : 'tak';
}

/** Buduje kartę komentarza wraz z wątkiem odpowiedzi. */
function utworzDymek(
  komentarz: StudioComment,
  odpowiedzi: readonly StudioComment[],
  czynnosci: CzynnosciDymkow,
): HTMLElement {
  const dymek = document.createElement('article');
  dymek.className = 'ms-dymek-komentarza';
  dymek.dataset['komentarz'] = komentarz.id;
  dymek.dataset['rodzaj'] = RodzajZnakowania.Komentarz;
  dymek.dataset['rozwiazany'] = komentarz.resolved === true ? 'tak' : 'nie';

  dymek.append(
    naglowekKarty(NAZWY_RODZAJOW.komentarz.nazwa, NAZWY_RODZAJOW.komentarz.czym),
    glowaWpisu(komentarz),
    trescWpisu(komentarz.body),
  );

  for (const odpowiedz of odpowiedzi) {
    const wpis = document.createElement('div');
    wpis.className = 'ms-dymek-komentarza__odpowiedz';
    wpis.append(glowaWpisu(odpowiedz), trescWpisu(odpowiedz.body));
    dymek.append(wpis);
  }

  const pole = document.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole-kontrolka';
  pole.placeholder = 'odpowiedz w wątku';
  pole.setAttribute('aria-label', `Odpowiedź w wątku komentarza ${komentarz.id}`);

  const odpowiedz = document.createElement('button');
  odpowiedz.type = 'button';
  odpowiedz.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  odpowiedz.textContent = 'Odpowiedz';
  odpowiedz.addEventListener('click', () => {
    if (pole.value.trim() === '') return;
    czynnosci.naOdpowiedz(komentarz.id, pole.value.trim());
    pole.value = '';
  });

  const rozwiaz = document.createElement('button');
  rozwiaz.type = 'button';
  rozwiaz.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  rozwiaz.textContent = komentarz.resolved === true ? 'Otwórz wątek' : 'Rozwiąż wątek';
  rozwiaz.title =
    'Komendy usuwającej komentarz kontrakt nie niesie — studio.comment.* ma dodanie, wykaz ' +
    'i rozwiązanie. Rozwiązanie zamyka wątek bez kasowania go i da się je cofnąć.';
  rozwiaz.addEventListener('click', () =>
    czynnosci.naRozwiazanie(komentarz.id, komentarz.resolved !== true),
  );

  const doKotwicy = document.createElement('button');
  doKotwicy.type = 'button';
  doKotwicy.className = 'dn-btn dn-btn--sm dn-btn--duch';
  doKotwicy.textContent = 'Pokaż fragment';
  doKotwicy.addEventListener('click', () => czynnosci.naKotwice(komentarz.id));

  const pas = document.createElement('div');
  pas.className = 'ms-dymek-komentarza__pas';
  pas.append(pole, odpowiedz, rozwiaz, doKotwicy);
  dymek.append(pas);
  return dymek;
}

/**
 * Karta propozycji brzmienia.
 *
 * Odróżniona od komentarza tym, że **pokazuje brzmienie** wraz z jego długością,
 * i od zmiany śledzonej tym, że mówi wprost: w treści tego jeszcze nie ma.
 * Decyzja idzie `studio.proposal.decide`, a nie `studio.tracking.decide` — to
 * dwie różne komendy dla dwóch różnych bytów.
 */
function utworzKartePropozycji(
  propozycja: PropozycjaZmiany,
  zakres: { poczatek: number; koniec: number } | null,
  czynnosci: CzynnosciDymkow,
): HTMLElement {
  const karta = document.createElement('article');
  karta.className = 'ms-dymek-propozycji';
  karta.dataset['rodzaj'] = RodzajZnakowania.Propozycja;
  karta.dataset['propozycja'] =
    propozycja.idPropozycji === '' ? 'bez-odwolania' : propozycja.idPropozycji;

  const glowa = document.createElement('p');
  glowa.className = 'ms-dymek-komentarza__glowa';
  glowa.textContent =
    `model · operacja ${propozycja.idAkcji} · ${propozycja.tresc.length} znaków · ` +
    (zakres === null
      ? 'dotyczy całego dokumentu'
      : `znaki ${zakres.poczatek}–${zakres.koniec}`);

  const przyjmij = document.createElement('button');
  przyjmij.type = 'button';
  przyjmij.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
  przyjmij.textContent = 'Przyjmij brzmienie';
  przyjmij.dataset['czynnosc'] = 'przyjmij-propozycje';
  przyjmij.title =
    'Dopiero przyjęcie wnosi brzmienie do treści — komenda studio.proposal.decide. Do tej chwili ' +
    'dokument jest nietknięty.';
  przyjmij.addEventListener('click', () => czynnosci.naDecyzjePropozycji(true));

  const odrzuc = document.createElement('button');
  odrzuc.type = 'button';
  odrzuc.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  odrzuc.textContent = 'Odrzuć brzmienie';
  odrzuc.dataset['czynnosc'] = 'odrzuc-propozycje';
  odrzuc.addEventListener('click', () => czynnosci.naDecyzjePropozycji(false));

  const pas = document.createElement('div');
  pas.className = 'ms-dymek-komentarza__pas';
  pas.append(przyjmij, odrzuc);

  karta.append(
    naglowekKarty(NAZWY_RODZAJOW.propozycja.nazwa, NAZWY_RODZAJOW.propozycja.czym),
    glowa,
    trescWpisu(propozycja.tresc),
    pas,
  );
  return karta;
}

/**
 * Karta zmiany śledzonej.
 *
 * Stoi na marginesie obok oznaczenia w treści, bo oznaczenie w treści mówi
 * „gdzie", a karta mówi „co było przed" — i to drugie jest tym, czego Operator
 * potrzebuje do decyzji.
 */
function utworzKarteZmiany(
  zmiana: StudioTrackedChange,
  czynnosci: CzynnosciDymkow,
): HTMLElement {
  const karta = document.createElement('article');
  karta.className = 'ms-dymek-zmiany';
  karta.dataset['rodzaj'] = RodzajZnakowania.Zmiana;
  karta.dataset['zmiana'] = zmiana.id;
  karta.dataset['oczekuje'] = zmiana.decision === StudioChangeDecision.Oczekuje ? 'tak' : 'nie';

  const glowa = document.createElement('p');
  glowa.className = 'ms-dymek-komentarza__glowa';
  glowa.textContent = opiszZmiane(zmiana);

  const przed = document.createElement('p');
  przed.className = 'ms-dymek-zmiany__przed';
  przed.textContent =
    zmiana.before === undefined || zmiana.before === ''
      ? 'Przed zmianą nie stało tu nic — to dopisanie.'
      : `Przed: ${zmiana.before}`;

  const przyjmij = document.createElement('button');
  przyjmij.type = 'button';
  przyjmij.className = 'dn-btn dn-btn--sm dn-btn--sygnal';
  przyjmij.textContent = 'Przyjmij zmianę';
  przyjmij.title =
    'Zmiana jest JUŻ w treści — przyjęcie ją tam zostawia, a odrzucenie wycofuje. Komenda ' +
    'studio.tracking.decide.';
  przyjmij.addEventListener('click', () => czynnosci.naDecyzjeZmiany(zmiana.id, true));

  const odrzuc = document.createElement('button');
  odrzuc.type = 'button';
  odrzuc.className = 'dn-btn dn-btn--sm dn-btn--zarys';
  odrzuc.textContent = 'Odrzuć zmianę';
  odrzuc.addEventListener('click', () => czynnosci.naDecyzjeZmiany(zmiana.id, false));

  const doKotwicy = document.createElement('button');
  doKotwicy.type = 'button';
  doKotwicy.className = 'dn-btn dn-btn--sm dn-btn--duch';
  doKotwicy.textContent = 'Pokaż miejsce';
  doKotwicy.addEventListener('click', () => czynnosci.naKotwice(zmiana.id));

  const pas = document.createElement('div');
  pas.className = 'ms-dymek-komentarza__pas';
  pas.append(przyjmij, odrzuc, doKotwicy);

  karta.append(
    naglowekKarty(NAZWY_RODZAJOW.zmiana.nazwa, NAZWY_RODZAJOW.zmiana.czym),
    glowa,
    przed,
    pas,
  );
  return karta;
}

/** Nagłówek karty: nazwa rodzaju i zdanie o tym, czym się on różni. */
function naglowekKarty(nazwa: string, czym: string): HTMLElement {
  const podpis = document.createElement('p');
  podpis.className = 'ms-dymek__rodzaj';
  podpis.textContent = nazwa;
  podpis.title = czym;

  const naglowek = document.createElement('header');
  naglowek.className = 'ms-dymek__naglowek';
  naglowek.append(podpis);
  return naglowek;
}

/** Głowa wpisu: autor, czas i zakres, o ile komentarz jest przypięty. */
function glowaWpisu(komentarz: StudioComment): HTMLElement {
  const glowa = document.createElement('p');
  glowa.className = 'ms-dymek-komentarza__glowa';
  const czas = new Date(komentarz.createdAt).toLocaleString('pl-PL');
  const zakres =
    komentarz.selectionStart === undefined
      ? 'dotyczy całego dokumentu'
      : `znaki ${komentarz.selectionStart}–${komentarz.selectionEnd ?? komentarz.selectionStart}`;
  glowa.textContent = `${nazwaAutora(komentarz.author)} · ${czas} · ${zakres}`;
  return glowa;
}

/** Treść wpisu karty. */
function trescWpisu(tresc: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'ms-dymek-komentarza__tresc';
  element.textContent = tresc;
  return element;
}
