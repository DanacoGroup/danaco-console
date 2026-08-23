import { elementIkony } from '../ikony/ikony';
import type { CzynnosciSesji } from './czynnosci-sesji';
import { utworzKarteSesji, type CzynnoscPowrotu } from './karta-sesji';
import { utworzStrefeZwijana } from './strefa-zwijana';
import type { WykazSrodowisk } from './wykaz-srodowisk';
import type { MigawkaSesji } from './zrodlo-sesji';

/**
 * Strefa sesji w tle.
 *
 * Pokazuje migawkę źródła sesji i zgłasza żądanie powrotu. Strefa nie zna
 * kanału ani nazwy żadnej komendy — dane przynosi `zrodlo-sesji`, a wiąże je
 * `wpiecie-sesji`.
 *
 * Trzy postacie, żadnych danych miejscowych:
 *   oczekiwanie — rdzeń jeszcze nie odpowiedział na `session.list`,
 *   pusto       — rdzeń odpowiedział i sesji w tle nie ma,
 *   błąd        — rdzeń odmówił; stan pusty nie udaje wtedy braku sesji.
 * Wykaz powstaje wyłącznie z wpisów rdzenia.
 */

const ETYKIETA = 'Sesje w tle';
const WYJASNIENIE =
  'Rozłączenie klienta nie kończy sesji ani jej procesów — tu widać sesje trwające na rdzeniu.';

/** Treści stanów pustych; opisy nazywają przyczynę, nie udają danych. */
const PUSTE = {
  oczekiwanie: {
    ikona: 'zegar',
    tytul: 'Oczekiwanie na rdzeń',
    opis: 'Wykaz sesji w tle pojawi się, gdy rdzeń odpowie na pytanie o sesje.',
  },
  pusto: {
    ikona: 'aktywnosc',
    tytul: 'Brak sesji w tle',
    opis: 'Żadna sesja poza bieżącą nie trwa na rdzeniu. Sesja rozłączona pojawi się tu sama.',
  },
  blad: {
    ikona: 'ostrzezenie',
    tytul: 'Rdzeń odmówił wykazu sesji',
    opis: '',
  },
} as const;

export interface StrefaSesji {
  /** Element `<details>` — strefa przywoływana, nie rysowana z urzędu. */
  element: HTMLDetailsElement;
  /** Przerysowuje strefę do postaci z migawki źródła. */
  ustawMigawke(migawka: MigawkaSesji): void;
  /** Nadaje czynność powrotu; `null` zdejmuje przyciski z wierszy. */
  ustawPowrot(czynnosc: CzynnoscPowrotu | null): void;
  /**
   * Nadaje czynności historii sesji. Wykaz pusty zdejmuje menu z wierszy —
   * montaż bez tożsamości klienta zostaje przy wykazie informacyjnym.
   */
  ustawCzynnosci(czynnosci: CzynnosciSesji): void;
  /** Dokłada wgląd w archiwum pod wykazem sesji bieżących; wolno raz. */
  osadzArchiwum(element: HTMLElement): void;
  /** Pokazuje treść odmowy czynności; znika przy następnej migawce. */
  zglosKomunikat(tekst: string): void;
}

/**
 * @param srodowiska Wykaz środowisk strony — nazwa środowiska w wierszu sesji
 *   pochodzi z tego samego bytu, co nazwa na karcie strefy pierwszej.
 */
export function utworzStrefeSesji(srodowiska: WykazSrodowisk): StrefaSesji {
  let powrot: CzynnoscPowrotu | null = null;
  let czynnosci: CzynnosciSesji = {};
  let migawka: MigawkaSesji = { stan: 'oczekiwanie', wpisy: [] };

  // Strefa jest zwinięta domyślnie: sesje w tle to wgląd w pracę już biegnącą,
  // czyli drugi plan wobec dróg wejścia w pracę nową. Zapowiedź niesie liczbę
  // sesji, więc zwinięcie nie ukrywa faktu, że coś trwa.
  const strefa = utworzStrefeZwijana({
    etykieta: ETYKIETA,
    wyjasnienie: WYJASNIENIE,
    klucz: 'strona.sesje',
    domyslnieRozwiniete: false,
  });
  const element = strefa.element;
  element.classList.add('dn-strona__strefa--sesje');
  element.setAttribute('aria-label', ETYKIETA);

  const komunikat = document.createElement('p');
  komunikat.className = 'dn-plakietka dn-plakietka--blad dn-strona__sesje-komunikat';
  komunikat.hidden = true;

  const tresc = document.createElement('div');
  tresc.className = 'dn-strona__sesje';
  // Wykaz zmienia się zdarzeniami rdzenia w trakcie pracy — czytnik ekranu
  // dostaje zmianę bez odebrania ogniska.
  tresc.setAttribute('aria-live', 'polite');

  strefa.cialo.append(komunikat, tresc);

  /**
   * Dopisek zapowiedzi — liczba sesji z migawki, nigdy liczba wymyślona.
   *
   * Przed odpowiedzią rdzenia i przy odmowie dopisek jest pusty: „0" znaczyłoby
   * „rdzeń odpowiedział i sesji nie ma", czyli co innego.
   */
  function odswiezDopisek(): void {
    if (migawka.stan === 'oczekiwanie' || migawka.stan === 'blad') {
      strefa.ustawDopisek('');
      return;
    }
    strefa.ustawDopisek(String(migawka.wpisy.length));
  }

  function przerysuj(): void {
    odswiezDopisek();
    komunikat.hidden = true;
    if (migawka.wpisy.length > 0) {
      tresc.replaceChildren(zbudujWykaz());
      return;
    }
    if (migawka.stan === 'blad') {
      tresc.replaceChildren(zbudujPustke('blad', migawka.blad));
      return;
    }
    tresc.replaceChildren(zbudujPustke(migawka.stan === 'oczekiwanie' ? 'oczekiwanie' : 'pusto'));
  }

  function zbudujWykaz(): HTMLUListElement {
    const wykaz = document.createElement('ul');
    wykaz.className = 'dn-strona__sesje-wykaz';
    for (const wpis of migawka.wpisy) {
      wykaz.append(
        utworzKarteSesji(wpis, { powrot, czynnosci, meldunek: zglos, srodowiska }).element,
      );
    }
    return wykaz;
  }

  /**
   * Meldunek czynności idzie tą samą drogą co odmowa powrotu: jedno pole nad
   * wykazem, znikające przy następnej migawce. Menu wiersza nie buduje
   * własnego miejsca na treść — inaczej ta sama odpowiedź rdzenia pojawiałaby
   * się w dwóch postaciach zależnie od tego, kto ją wywołał.
   */
  function zglos(tekst: string): void {
    komunikat.textContent = tekst;
    komunikat.hidden = false;
  }

  przerysuj();

  return {
    element,
    ustawMigawke(nowa) {
      migawka = nowa;
      przerysuj();
    },
    ustawPowrot(czynnosc) {
      powrot = czynnosc;
      przerysuj();
    },
    ustawCzynnosci(nowe) {
      czynnosci = nowe;
      przerysuj();
    },
    osadzArchiwum(archiwum) {
      // Archiwum stoi pod wykazem i poza `tresc`, bo `tresc` jest w całości
      // przerysowywana przy każdej migawce — wgląd rozwinięty przez Operatora
      // zwijałby się wtedy przy każdej zmianie w sesjach bieżących.
      strefa.cialo.append(archiwum);
    },
    zglosKomunikat: zglos,
  };
}

/** Stan pusty z biblioteki komponentów (`.dn-pusty-stan`, `komponenty/drobne.css`). */
function zbudujPustke(postac: keyof typeof PUSTE, opisBledu?: string): HTMLElement {
  const wzor = PUSTE[postac];
  const pustka = document.createElement('div');
  pustka.className = 'dn-pusty-stan dn-strona__sesje-pustka';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';
  tytul.textContent = wzor.tytul;

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';
  opis.textContent = postac === 'blad' ? (opisBledu ?? wzor.opis) : wzor.opis;

  pustka.append(elementIkony(wzor.ikona, { rozmiar: 24 }), tytul, opis);
  return pustka;
}
