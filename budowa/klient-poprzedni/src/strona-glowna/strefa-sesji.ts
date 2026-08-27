import { elementIkony } from '../ikony/ikony';
import type { CzynnosciSesji } from './czynnosci-sesji';
import { utworzKarteSesji, type CzynnoscPowrotu } from './karta-sesji';
import { utworzStrefeZwijana } from './strefa-zwijana';
import type { WykazSrodowisk } from './wykaz-srodowisk';
import type { MigawkaSesji } from './zrodlo-sesji';

/** Strefa sesji w tle pokazuje migawkę źródła sesji i zgłasza żądanie powrotu, rozróżniając osobnymi napisami oczekiwanie, brak sesji i odmowę rdzenia. */
const ETYKIETA = 'Sesje w tle';
const WYJASNIENIE =
  'Rozłączenie klienta nie kończy sesji ani jej procesów — tu widać sesje trwające na rdzeniu.';

/** Treści stanów pustych strefy sesji w tle: opisy nazywają przyczynę braku wykazu, nie udają danych, których nie ma. */
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
  /** Nadaje czynności historii sesji; wykaz pusty zdejmuje menu z wierszy sesji. */
  ustawCzynnosci(czynnosci: CzynnosciSesji): void;
  /** Dokłada wgląd w archiwum pod wykazem sesji bieżących; wolno raz. */
  osadzArchiwum(element: HTMLElement): void;
  /** Pokazuje treść odmowy czynności; znika przy następnej migawce. */
  zglosKomunikat(tekst: string): void;
}

/** Buduje strefę sesji dla podanego wykazu środowisk, z którego wiersz sesji czerpie nazwę środowiska tak samo jak karta strefy pierwszej. */
export function utworzStrefeSesji(srodowiska: WykazSrodowisk): StrefaSesji {
  let powrot: CzynnoscPowrotu | null = null;
  let czynnosci: CzynnosciSesji = {};
  let migawka: MigawkaSesji = { stan: 'oczekiwanie', wpisy: [] };

  // Strefa jest zwinięta domyślnie: sesje w tle to wgląd w pracę już biegnącą.
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
  // Wykaz zmienia się zdarzeniami rdzenia w trakcie pracy; czytnik dostaje zmianę bez odebrania ogniska.
  tresc.setAttribute('aria-live', 'polite');

  strefa.cialo.append(komunikat, tresc);

  // Dopisek jest pusty przed odpowiedzią rdzenia i przy odmowie, bo zero znaczyłoby co innego.
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

  // Meldunek czynności idzie drogą odmowy powrotu; menu wiersza nie buduje własnego miejsca na treść.
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
      // Archiwum stoi poza treścią przerysowywaną, żeby wgląd rozwinięty nie zwijał się przy zmianie sesji.
      strefa.cialo.append(archiwum);
    },
    zglosKomunikat: zglos,
  };
}

/** Buduje jeden ze stanów pustych strefy z gotowego wzoru biblioteki komponentów, dobranego po postaci przekazanej wywołaniu. */
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
