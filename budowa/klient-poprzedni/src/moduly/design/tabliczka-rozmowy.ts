import { liczbaOkienRozmowy } from '../../okno-komunikacji/profil-modulu';
import { profilModulu } from '../../okno-komunikacji/rejestr-profilow';
import { KOD_MODULU, OKNO_DESIGN_BOARD } from './etykiety-designu';

/**
 * Tabliczka wejścia rozmowy — mówi, jak czat wiąże się z oknem roboczym, i liczy,
 * co tą drogą weszło.
 *
 * Postać rozmowy i okna obowiązkowe czyta z profilu modułu
 * (`okno-komunikacji/rejestr-profilow.ts`), więc zdanie o wiązaniu zmienia się
 * razem z profilem zamiast powtarzać go z pamięci.
 *
 * Liczby okien tabliczka nie podaje: profil Designu niesie `GRANICA_NIEPODANA`.
 * Pyta wyłącznie o prawo do rozmowy przez `liczbaOkienRozmowy`.
 *
 * Licznik obejmuje tylko wejścia rozmowy; zasoby z Prompt Buildera i Assets Panel
 * idą własną drogą.
 */
export interface TabliczkaRozmowy {
  /** Element osadzany w treści okna roboczego. */
  element: HTMLElement;
  /** Notuje jedno wejście rozmowy wraz z jego zdaniem. */
  zanotuj(zdanie: string): void;
}

/** Zdanie stanu do chwili pierwszego wejścia — brak drogi, nie brak zasobów. */
const BEZ_WEJSCIA =
  'Od otwarcia okna nie weszło tędy nic. To nie znaczy „brak zasobów": zasoby z Prompt ' +
  'Buildera i te przeniesione z Assets Panel idą własną drogą i tego licznika nie ruszają.';

export function utworzTabliczkeRozmowy(): TabliczkaRozmowy {
  let wejscia = 0;

  const naglowek = document.createElement('h4');
  naglowek.className = 'md-wejscie__naglowek';
  naglowek.textContent = 'Wejście rozmowy';

  const wiazanie = document.createElement('p');
  wiazanie.className = 'dn-pole-opis md-wejscie__wiazanie';
  wiazanie.textContent = zdanieOWiazaniu();

  const stan = document.createElement('p');
  stan.className = 'dn-pole-opis md-wejscie__stan';
  stan.textContent = BEZ_WEJSCIA;

  const element = document.createElement('section');
  element.className = 'md-wejscie';
  element.setAttribute('aria-label', 'Wejście rozmowy na okno robocze');
  element.append(naglowek, wiazanie, stan);

  return {
    element,

    zanotuj(zdanie) {
      wejscia += 1;
      stan.textContent = `Wejść z rozmowy od otwarcia okna: ${wejscia}. Ostatnie — ${zdanie}`;
    },
  };
}

/**
 * Zdanie o wiązaniu czatu z oknem roboczym — złożone z profilu modułu, więc
 * mówi to samo, co reszta produktu, i zmienia się razem z nim.
 */
export function zdanieOWiazaniu(): string {
  const profil = profilModulu(KOD_MODULU);
  const rola = profil.oknaObowiazkowe.includes(OKNO_DESIGN_BOARD.nazwa)
    ? `${OKNO_DESIGN_BOARD.nazwa} jest oknem roboczym modułu — profil wymienia je wśród okien ` +
      `obowiązkowych (${profil.oknaObowiazkowe.join(' · ')}).`
    : `UWAGA: profil modułu ${profil.nazwa} NIE wymienia okna ${OKNO_DESIGN_BOARD.nazwa} wśród ` +
      'obowiązkowych, a to okno zachowuje się jak robocze — rozejście się okna z profilem ' +
      'wymaga rozstrzygnięcia, nie ciszy.';

  if (liczbaOkienRozmowy(profil) === 0) {
    return (
      `${rola} Rozmowy moduł jednak NIE prowadzi — profil niesie postać „${profil.postacRozmowy}", ` +
      'więc wynik rozmowy nie ma tędy czym wejść i cała praca powstaje w oknach modułu.'
    );
  }

  if (profil.postacRozmowy === 'dymek-glosowy') {
    return (
      `${rola} Rozmowę prowadzi pływający awatar z oknem dymkowym, nie okno czatu. Jej wynik ` +
      'wchodzi tutaj zdarzeniem rdzenia design.asset.changed i kładzie się warstwą na kanwie.'
    );
  }

  return (
    `${rola} Rozmowę prowadzi Chat Window — jedno okno wspólne wszystkim modułom, przestawione ` +
    `na profil ${profil.nazwa}; stoi obok, na scenie sesji, a nie w tym module. Wynik rozmowy ` +
    'wchodzi tutaj zdarzeniem rdzenia design.asset.changed i sam kładzie się warstwą na kanwie; ' +
    'do rdzenia warstwa jedzie dopiero zapisem kompozycji.'
  );
}
