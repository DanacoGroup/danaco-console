import './karty-sesji.css';

import { ProgressStatus } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { nowyIdentyfikator } from '../protokol/identyfikator';
import { utworzKarteSesji, type KartaSesji } from './karta-sesji';
import { utworzPostacPasa } from './postac-pasa-kart';
import { przeniesWedrowke } from './wedrowka-kart';
import { zwiazPasKartZRdzeniem, type ZamiaryKart } from './wpiecie-kart-sesji';
import type { MigawkaKart } from './zrodlo-kart-sesji';

/**
 * Pas kart sesji powłoki — poziome karty o mechanice zakładek.
 *
 * Jedna odpowiedzialność: wykaz kart i wybór jednej z nich. Wygląd karty
 * należy do `karta-sesji.ts`, klawiatura do `wedrowka-kart.ts`, stan pusty
 * i komunikat do `postac-pasa-kart.ts`, a droga do rdzenia do
 * `wpiecie-kart-sesji.ts`.
 *
 * Karta sesji jest sesją rdzenia. Gdy aplikacja udostępniła drogę do rdzenia,
 * pas przestaje być własnym wykazem: karty biorą się z `session.list` i zdarzeń
 * `session.changed`, ＋ zakłada sesję (`session.create`), zamknięcie karty
 * zamyka sesję (`session.close`), kosz usuwa ją trwale (`session.delete`, po
 * potwierdzeniu), a wybór przenosi ognisko klienta (`session.focus`).
 *
 * Zamknięcie i usunięcie to dwie czynności, nie dwie nazwy jednej. Krzyżyk
 * zamyka sesję i zostawia jej zapis w całości; kosz zdejmuje sesję z historii,
 * a zapis rdzeń kasuje trwale po terminie kosza — potwierdzenie pokazuje wykaz
 * tego, co znika (`potwierdzenie-usuniecia.ts`). Skład pasa rozstrzyga wtedy
 * wyłącznie rdzeń i pas nie zakłada ani jednej karty miejscowej. Bez rdzenia
 * pas jest samym widokiem.
 *
 * Karta sesji to nie okno komunikacji. Liczbą okien na scenie rządzi
 * przełącznik „Okna komunikacji: 1 2 3" układu okien równoległych, więc ＋ nie
 * wprowadza okna na scenę, a zamknięcie karty nie zdejmuje okna ze sceny.
 * Mechanika zakładek idzie wzorcem ARIA: ogniskuje się wyłącznie karta czynna,
 * a ＋ stoi poza pasem zakładek — zakłada kartę, nie wybiera istniejącej.
 */

/** Słuchacz zdarzenia dotyczącego jednej karty. */
export type SluchaczKarty = (karta: KartaSesji) => void;

export interface KartySesji {
  /** Pas montowany jako drugi wiersz powłoki. */
  element: HTMLElement;
  /** Zakłada kartę miejscową; w pasie rdzenia karta nie wchodzi do pasa. */
  dodaj(tytul: string, stan?: ProgressStatus): KartaSesji;
  /** Zamyka kartę o wskazanym identyfikatorze. */
  zamknij(id: string): void;
  /** Czyni kartę czynną — zaznaczenie w pasie, bez komendy kontraktu. */
  wybierz(id: string): void;
  /** Karta czynna albo brak, gdy pas jest pusty. */
  czynna(): KartaSesji | undefined;
  /** Wykaz kart w kolejności wyświetlania. */
  wykaz(): readonly KartaSesji[];
  /** Nadaje karcie czynnej nowy tytuł; w pasie rdzenia tytuł ma rdzeń. */
  przemianujCzynna(tytul: string): void;
  /** Przerysowuje pas do postaci z migawki rdzenia. */
  ustawMigawke(migawka: MigawkaKart): void;
  /** Pokazuje odpowiedź rdzenia przy pasie; następna migawka ją zdejmuje. */
  zglosKomunikat(tekst: string, waga?: 'blad' | 'info'): void;
  naWybor(sluchacz: SluchaczKarty): void;
  naZamkniecie(sluchacz: SluchaczKarty): void;
  naNowa(sluchacz: SluchaczKarty): void;
}

/** Ustawienia pasa: tytuł karty zakładanej przyciskiem ＋. */
export interface OpcjeKart {
  tytulNowej?(): string;
}

export function utworzKartySesji(opcje: OpcjeKart = {}): KartySesji {
  const karty: KartaSesji[] = [];
  const sluchaczeWyboru: SluchaczKarty[] = [];
  const sluchaczeZamkniecia: SluchaczKarty[] = [];
  const sluchaczeNowej: SluchaczKarty[] = [];
  let zamiary: ZamiaryKart | null = null;

  const element = document.createElement('div');
  element.className = 'dn-pasek dn-pasek--sesje dn-sesje';

  const lista = document.createElement('div');
  lista.className = 'dn-sesje__lista';
  lista.setAttribute('role', 'tablist');
  lista.setAttribute('aria-label', 'Karty sesji');

  const postac = utworzPostacPasa();

  const dodajKarte = document.createElement('button');
  dodajKarte.type = 'button';
  dodajKarte.className = 'dn-btn-ikona dn-sesje__dodaj';
  dodajKarte.setAttribute('aria-label', 'Załóż sesję');
  dodajKarte.title = 'Załóż sesję';
  dodajKarte.append(elementIkony('plus', { rozmiar: 18 }));
  dodajKarte.addEventListener('click', () => {
    const tytul = opcje.tytulNowej?.() ?? 'Nowa sesja';
    if (zamiary !== null) zamiary.przyNowej(tytul);
    else powiadom(sluchaczeNowej, dodaj(tytul));
  });

  element.append(lista, postac.pustka, dodajKarte, postac.komunikat);

  // Naciśnięcie karty jest w pasie rdzenia zamiarem (rozstrzyga go rdzeń),
  // a poza nim zmianą miejscową — to jedyna różnica obu postaci pasa.
  function zamierzWybor(id: string): void {
    if (zamiary !== null) zamiary.przyWyborze(id);
    else wybierz(id);
  }

  function zamierzZamkniecie(id: string): void {
    if (zamiary !== null) zamiary.przyZamknieciu(id);
    else zamknij(id);
  }

  /**
   * Usunięcie trwałe ma sens wyłącznie w pasie rdzenia: karta miejscowa nie ma
   * zapisu, który dałoby się skasować. Poza rdzeniem pas nie usuwa karty na
   * pocieszenie — mówi wprost, że nie ma czego usuwać.
   */
  function zamierzUsuniecie(id: string): void {
    const wskazana = karty.find((karta) => karta.id === id);
    if (zamiary === null || wskazana === undefined) {
      postac.zglos('Usunięcie trwałe wymaga rdzenia — bez niego karta nie ma zapisu.', 'blad');
      return;
    }
    zamiary.przyUsunieciu(id, wskazana.tytul());
  }

  const obsluga = {
    przyWyborze: zamierzWybor,
    przyZamknieciu: zamierzZamkniecie,
    przyUsunieciu: zamierzUsuniecie,
    przyKlawiszu: (id: string, zdarzenie: KeyboardEvent): void =>
      przeniesWedrowke(karty, id, zdarzenie, { wybierz: zamierzWybor, zamknij: zamierzZamkniecie }),
  };

  function dodaj(tytul: string, stan: ProgressStatus = ProgressStatus.Pending): KartaSesji {
    const karta = utworzKarteSesji({ id: nowyIdentyfikator('karta'), tytul, stan }, obsluga);
    // Pas związany z rdzeniem nie przyjmuje kart miejscowych: karta powstaje
    // dopiero z odpowiedzi rdzenia, więc zwrócona karta zostaje poza pasem.
    if (zamiary !== null) return karta;
    karty.push(karta);
    lista.append(karta.element);
    wybierz(karta.id);
    return karta;
  }

  function wybierz(id: string): void {
    const wskazana = karty.find((karta) => karta.id === id);
    if (wskazana === undefined) return;
    for (const karta of karty) karta.ustawCzynna(karta.id === id);
    powiadom(sluchaczeWyboru, wskazana);
  }

  function zamknij(id: string): void {
    const miejsce = karty.findIndex((karta) => karta.id === id);
    if (miejsce < 0) return;

    const [zamknieta] = karty.splice(miejsce, 1);
    if (zamknieta === undefined) return;
    const byla = zamknieta.czyCzynna();
    zamknieta.element.remove();

    // Po zamknięciu karty czynnej pas nie zostaje bez wskazania: wybór
    // przechodzi na sąsiada z prawej, a przy jego braku z lewej.
    const nastepna = karty[Math.min(miejsce, karty.length - 1)];
    if (byla && nastepna !== undefined) {
      wybierz(nastepna.id);
      nastepna.ustawOgnisko();
    }
    powiadom(sluchaczeZamkniecia, zamknieta);
  }

  /**
   * Uzgodnienie pasa z migawką rdzenia.
   *
   * Karta spoza wykazu znika bez zdarzenia zamknięcia: zamknięcie sesji
   * w rdzeniu nie jest zamknięciem karty przez Operatora i nie ma zdejmować
   * okna ze sceny.
   */
  function ustawMigawke(migawka: MigawkaKart): void {
    postac.schowajKomunikat();
    const wykazane = new Set(migawka.wpisy.map((wpis) => wpis.id));
    for (const karta of karty) if (!wykazane.has(karta.id)) karta.element.remove();

    const uporzadkowane = migawka.wpisy.map((wpis) => {
      const istniejaca = karty.find((karta) => karta.id === wpis.id);
      if (istniejaca === undefined) {
        return utworzKarteSesji({ id: wpis.id, tytul: wpis.tytul, stan: wpis.stan }, obsluga);
      }
      istniejaca.ustawTytul(wpis.tytul);
      istniejaca.ustawStan(wpis.stan);
      return istniejaca;
    });
    karty.length = 0;
    karty.push(...uporzadkowane);
    for (const karta of karty) lista.append(karta.element);

    const czynna = karty.find((karta) => karta.id === migawka.ogniskowana);
    if (czynna !== undefined && !czynna.czyCzynna()) wybierz(czynna.id);
    else for (const karta of karty) karta.ustawCzynna(karta.id === migawka.ogniskowana);

    postac.ustawPustke(migawka, karty.length);
  }

  const pas: KartySesji = {
    element,
    dodaj,
    zamknij,
    wybierz,
    czynna: () => karty.find((karta) => karta.czyCzynna()),
    wykaz: () => karty,
    ustawMigawke,

    przemianujCzynna(tytul) {
      // W pasie rdzenia tytuł karty jest nazwą sesji, nie nazwą modułu: moduł
      // widać w kontekście paska górnego, a kontrakt nie zna komendy zmiany
      // nazwy sesji — nadpisanie miejscowe rozjechałoby pas z rdzeniem przy
      // pierwszej migawce.
      if (zamiary !== null) return;
      karty.find((karta) => karta.czyCzynna())?.ustawTytul(tytul);
    },

    zglosKomunikat: (tekst, waga = 'blad') => postac.zglos(tekst, waga),

    naWybor: (sluchacz) => void sluchaczeWyboru.push(sluchacz),
    naZamkniecie: (sluchacz) => void sluchaczeZamkniecia.push(sluchacz),
    naNowa: (sluchacz) => void sluchaczeNowej.push(sluchacz),
  };

  // Wpięcie rdzenia na końcu: zamiary muszą zastać pas gotowy do przyjęcia
  // pierwszej migawki.
  zamiary = zwiazPasKartZRdzeniem(pas);

  return pas;
}

/** Rozesłanie zdarzenia do wszystkich słuchaczy. */
function powiadom(sluchacze: readonly SluchaczKarty[], karta: KartaSesji): void {
  for (const sluchacz of sluchacze) sluchacz(karta);
}

export type { KartaSesji };
