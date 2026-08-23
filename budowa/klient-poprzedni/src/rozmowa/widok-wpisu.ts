import './wpis.css';

import { elementIkony } from '../ikony/ikony';
import { utworzBlokNarzedzi } from './blok-narzedzi';
import { utworzBlokProwenancji } from './blok-prowenancji';
import { utworzBlokZwijany } from './blok-zwijany';
import { NAPISY, nazwaStanuWpisu } from './etykiety-rozmowy';
import { RodzajNadawcy, znakiNadawcy, type KlasaNadawcy } from './nadawca';
import { czesciStopki, wierszBledu } from './stopka-wpisu';
import {
  warstwyZapisu,
  WIDOK_ZAPISU_DOMYSLNY,
  type WarstwyZapisu,
  type WidokZapisu,
} from './widok-zapisu';
import { godzinaWpisu, type WpisRozmowy } from './wpis-rozmowy';

/** Widok jednego wpisu historii. */
export interface WidokWpisu {
  /** Element montowany w liście. */
  element: HTMLElement;
  /** Odświeża widok tym samym wpisem po zmianie. */
  aktualizuj(wpis: WpisRozmowy): void;
  /** Przestawia wpis na inny tryb widoku transkryptu; treść zostaje w pamięci. */
  ustawWidokZapisu(widok: WidokZapisu): void;
}

/**
 * Widok wpisu rozmowy.
 *
 * Wpis jest siatką dwukolumnową `24px 1fr` (`komponenty/wpis.css`):
 * w pierwszej kolumnie medalion nadawcy, w drugiej tożsamość i wszystko, co pod
 * nią idzie. Bez medalionu pasek tożsamości wpadłby w kolumnę awatara i został
 * zgnieciony do 24 px.
 *
 * Kolejność warstw odpowiada kolejności strumienia: prowenancja stoi przed
 * treścią, bo rdzeń nadaje ją przed jakimkolwiek tekstem i przed startem
 * procesu. Dalej idzie podgląd pracy modelu, potem odpowiedź, potem narzędzia,
 * błędy i podsumowanie tury.
 *
 * Rozróżnienie nadawcy niesie ikona medalionu, klasa semantyczna (trzy barwy,
 * nie dziewięć — `nadawca.ts`) i etykieta słowna — nigdy sama barwa tła.
 *
 * Klasy `dc-wpis` i `dc-wpis--<nadawca>` nie mają własnych reguł; stoją jako
 * uchwyty, po których wpis i jego nadawca dają się rozpoznać z zewnątrz.
 *
 * Widok transkryptu jest filtrem nad tym wpisem, nie drugim widokiem. Cztery
 * tryby (`widok-zapisu.ts`) sterują wyłącznie tym, które warstwy są rysowane;
 * wpis zostaje tym samym obiektem na tej samej pozycji listy. Przełączenie nie
 * kasuje niczego, nie woła rdzenia i nie blokuje — schowana warstwa wraca
 * w całości po powrocie do trybu, który ją pokazuje.
 */
export function utworzWidokWpisu(
  wpis: WpisRozmowy,
  widok: WidokZapisu = WIDOK_ZAPISU_DOMYSLNY,
): WidokWpisu {
  const znaki = znakiNadawcy(wpis.nadawca);
  let warstwy: WarstwyZapisu = warstwyZapisu(widok);
  // Ostatni wpis trzymany miejscowo: zmiana trybu przerysowuje wpis z tego, co
  // okno już ma, bez pytania rdzenia o powtórzenie tury.
  let ostatni: WpisRozmowy = wpis;

  const element = document.createElement('article');

  const medalion = document.createElement('span');
  medalion.className = 'dn-wpis-medalion';
  medalion.append(
    elementIkony(znaki.ikona, { rozmiar: 16, etykieta: znaki.etykieta }),
  );

  const tozsamosc = document.createElement('header');
  tozsamosc.className = 'dn-wpis-tozsamosc';

  // Persona mieszka w tekście nadawcy („Model · Redaktor”), a nie w osobnym
  // elemencie o własnym stopniu i grubości.
  const etykieta = document.createElement('span');
  etykieta.className = 'dn-wpis-nadawca';

  // Plakietka roli przy wyniku narzędzia — niesie nazwę narzędzia.
  const rola = document.createElement('span');
  rola.className = 'dn-plakietka dn-plakietka--rola';
  rola.hidden = true;

  const godzina = document.createElement('time');
  godzina.className = 'dn-wpis-godzina';

  const stan = document.createElement('span');
  // Odmianę stanu niesie klasa widoku `dc-wpis__stan`; biblioteka nie ma dla
  // niej osobnego wariantu plakietki.
  stan.className = 'dn-plakietka dc-wpis__stan';

  tozsamosc.append(etykieta, rola, godzina, stan);

  const prowenancja = utworzBlokProwenancji();
  const rozumowanie = utworzBlokZwijany({ tytul: NAPISY.rozumowanie, ikona: 'zegar' });
  const narzedzia = utworzBlokNarzedzi();

  const tresc = document.createElement('p');
  tresc.className = 'dn-wpis-tresc';

  const bledy = document.createElement('div');
  bledy.className = 'dc-wpis__bledy';

  const stopka = document.createElement('footer');
  stopka.className = 'dc-wpis__stopka';

  element.append(
    medalion,
    tozsamosc,
    prowenancja.element,
    rozumowanie.element,
    tresc,
    narzedzia.element,
    bledy,
    stopka,
  );

  function aktualizuj(nowy: WpisRozmowy): void {
    ostatni = nowy;
    element.className = klasyWpisu(nowy, znaki.klasa);
    etykieta.textContent = tekstTozsamosci(nowy, znaki.etykieta);
    godzina.textContent = godzinaWpisu(nowy);
    element.dataset['stan'] = nowy.stan;
    stan.textContent = opisStanuWpisu(nowy);
    ustawRoleNarzedzia(nowy);

    // Warstwa wyłączona trybem dostaje puste dane, nie ukryty element z treścią:
    // blok, który sam chowa się przy braku danych, po przełączeniu trybu nie
    // zostawia w drzewie napisu, którego Operator nie prosił o pokazanie.
    prowenancja.aktualizuj(warstwy.prowenancja ? nowy.prowenancja : null);
    aktualizujRozumowanie(nowy);
    tresc.textContent = warstwy.tresc ? nowy.tresc : '';
    tresc.hidden = !warstwy.tresc || nowy.tresc.length === 0;
    narzedzia.aktualizuj(warstwy.narzedzia ? nowy.narzedzia : []);
    bledy.replaceChildren(...nowy.bledy.map(wierszBledu));
    stopka.replaceChildren(...czesciStopki(nowy, warstwy));
  }

  /**
   * Przestawienie trybu: nowy rozkład warstw i jedno przerysowanie.
   *
   * Rozwinięcie bloku narzędzi jest ustawiane tylko tutaj, a nie przy każdym
   * odświeżeniu — inaczej w trybie `pelny` kolejny fragment strumienia
   * rozwijałby z powrotem blok zwinięty ręką Operatora.
   */
  function ustawWidokZapisu(nowyWidok: WidokZapisu): void {
    warstwy = warstwyZapisu(nowyWidok);
    element.dataset['widokZapisu'] = nowyWidok;
    narzedzia.ustawRozwiniecie(warstwy.narzedziaRozwiniete);
    aktualizuj(ostatni);
  }

  /** Nazwa narzędzia przy wpisie z jego wynikiem; bez niej plakietka znika. */
  function ustawRoleNarzedzia(nowy: WpisRozmowy): void {
    const nazwa = nowy.nadawca === RodzajNadawcy.Narzedzie ? nowy.narzedzia[0]?.nazwa ?? '' : '';
    rola.textContent = nazwa;
    rola.hidden = nazwa.length === 0;
  }

  /**
   * Podgląd pracy modelu — zwinięty, żeby nie zasłaniał odpowiedzi.
   *
   * W trybach `zwykly` i `streszczenie` blok nie powstaje wcale, bo tok
   * rozumowania jest tym, co odróżnia tryb `rozumowanie` od `zwyklego`.
   */
  function aktualizujRozumowanie(nowy: WpisRozmowy): void {
    const jest = warstwy.rozumowanie && nowy.rozumowanie.length > 0;
    rozumowanie.pokaz(jest);
    if (!jest) return;
    rozumowanie.ustawPodtytul(`znaków: ${nowy.rozumowanie.length}`);
    rozumowanie.tresc.textContent = nowy.rozumowanie;
  }

  element.dataset['widokZapisu'] = widok;
  narzedzia.ustawRozwiniecie(warstwy.narzedziaRozwiniete);
  aktualizuj(wpis);
  return { element, aktualizuj, ustawWidokZapisu };
}

/**
 * Komplet klas wpisu: budowa z biblioteki, klasa semantyczna nadawcy, stan pracy
 * i dwie nazwy `dc-*` jako uchwyty rozpoznania nadawcy.
 *
 * `dn-wpis--pracuje` niesie kropkę tętna przy nadawcy (`komponenty/wpis.css`) —
 * wpis, w którym tura jeszcze biegnie.
 */
function klasyWpisu(wpis: WpisRozmowy, klasa: KlasaNadawcy): string {
  const pracuje = wpis.stan === 'wysylanie' || wpis.stan === 'strumien';
  return [
    'dn-wpis',
    `dn-wpis--${klasa}`,
    ...(pracuje ? ['dn-wpis--pracuje'] : []),
    'dc-wpis',
    `dc-wpis--${wpis.nadawca}`,
  ].join(' ');
}

/** Tożsamość mówiącego: rodzaj nadawcy i persona w jednym tekście. */
function tekstTozsamosci(wpis: WpisRozmowy, etykieta: string): string {
  const persona = wpis.persona.trim();
  if (persona.length === 0 || persona === etykieta) return etykieta;
  return `${etykieta} · ${persona}`;
}

/** Stan wpisu wraz z licznikiem fragmentów strumienia. */
function opisStanuWpisu(wpis: WpisRozmowy): string {
  const nazwa = nazwaStanuWpisu(wpis.stan);
  if (wpis.fragmenty === 0) return nazwa;
  return `${nazwa} · fragmentów: ${wpis.fragmenty}`;
}
