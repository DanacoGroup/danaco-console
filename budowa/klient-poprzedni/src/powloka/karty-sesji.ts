import './karty-sesji.css';

import { ProgressStatus } from '../../../shared/contract';
import { elementIkony } from '../ikony/ikony';
import { nowyIdentyfikator } from '../protokol/identyfikator';
import { utworzKarteSesji, type KartaSesji } from './karta-sesji';
import { utworzPostacPasa } from './postac-pasa-kart';
import { przeniesWedrowke } from './wedrowka-kart';
import { zwiazPasKartZRdzeniem, type ZamiaryKart } from './wpiecie-kart-sesji';
import type { MigawkaKart } from './zrodlo-kart-sesji';

// Pas kart sesji powłoki — poziome karty o mechanice zakładek, wykaz kart i wybór jednej z nich.

/** Słuchacz zdarzenia dotyczącego jednej karty, wywoływany przy każdej zmianie stanu pasa kart sesji powłoki. */
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

/** Ustawienia pasa: tytuł karty zakładanej przyciskiem plus, zanim rdzeń nada jej właściwą nazwę sesji. */
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

  // Naciśnięcie karty jest w pasie rdzenia zamiarem, poza nim zmianą miejscową — jedyna różnica pasa.
  function zamierzWybor(id: string): void {
    if (zamiary !== null) zamiary.przyWyborze(id);
    else wybierz(id);
  }

  function zamierzZamkniecie(id: string): void {
    if (zamiary !== null) zamiary.przyZamknieciu(id);
    else zamknij(id);
  }

  /** Usunięcie trwałe ma sens tylko w pasie rdzenia — karta miejscowa nie ma zapisu do skasowania. */
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
    // Pas związany z rdzeniem nie przyjmuje kart miejscowych: karta powstaje dopiero z odpowiedzi rdzenia.
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

    // Po zamknięciu karty czynnej pas nie zostaje bez wskazania: wybór przechodzi na sąsiada obok.
    const nastepna = karty[Math.min(miejsce, karty.length - 1)];
    if (byla && nastepna !== undefined) {
      wybierz(nastepna.id);
      nastepna.ustawOgnisko();
    }
    powiadom(sluchaczeZamkniecia, zamknieta);
  }

  /** Karta spoza wykazu znika bez zdarzenia zamknięcia — zamknięcie sesji to nie zamknięcie karty. */
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
      // W pasie rdzenia tytuł karty jest nazwą sesji, nie modułu, bo kontrakt nie zna zmiany nazwy sesji.
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

/** Rozesłanie zdarzenia do wszystkich słuchaczy zarejestrowanych na pasie kart sesji powłoki aplikacji. */
function powiadom(sluchacze: readonly SluchaczKarty[], karta: KartaSesji): void {
  for (const sluchacz of sluchacze) sluchacz(karta);
}

export type { KartaSesji };
