import { pokazKomunikat } from '../../aplikacja/komunikaty';
import { oznaczWarstwe } from './warstwy-designu';

/**
 * Pasek kontekstu kanwy Design Board — znaczniki warstwy pierwszej.
 *
 * Opracowanie wymienia nad kanwą znaczniki projektu, środowiska, modelu, trybu
 * narzędzia i poziomu powiększenia, a przy znaczniku kontekstowym stanowi, że
 * kliknięcie otwiera selektor warstwy drugiej, który zwija się po wyborze.
 *
 * Selektorów pasek NIE dubluje. Wybór silnika stoi w Prompt Builderze, tryb
 * narzędzia w znaczniku trybu nad kanwą, powiększenie w przyborniku — a dwie
 * kontrolki nastawiające jedną wartość rozjeżdżają się przy pierwszej zmianie
 * i Operator przestaje wiedzieć, która mówi prawdę. Kliknięcie znacznika mówi
 * więc, GDZIE jego selektor stoi, i tam Operatora odsyła.
 *
 * Znacznik środowiska jest tu jedynym, który nie ma czego pokazać. Środowisko
 * jest własnością powłoki i moduł nie dostaje go ani parametrem, ani żadną
 * komendą swojego obszaru; znacznik mówi to wprost, zamiast wpisywać nazwę
 * wziętą z niczego.
 */

/** Skąd znacznik bierze wartość i gdzie stoi jego selektor. */
export interface ZrodloZnacznika {
  /** Okno modułu ustalone z rdzenia; puste, gdy rdzeń go nie wskazał. */
  okno(): string;
  /** Zdanie o oknie modułu — z odczytu stanu okna albo o jego braku. */
  opisOkna(): string;
  /** Liczba kanałów obrazowych w rejestrze — kandydatów na silnik. */
  silnikow(): number;
  /** Nazwa trybu narzędzia kanwy. */
  tryb(): string;
  /** Powiększenie kanwy; 1 znaczy skalę naturalną. */
  powiekszenie(): number;
}

export interface PasekKontekstu {
  element: HTMLElement;
  odswiez(): void;
}

/** Jeden znacznik paska: nazwa, wartość bieżąca i zdanie po naciśnięciu. */
interface Znacznik {
  element: HTMLButtonElement;
  ustaw(wartosc: string, zdanie: string): void;
}

export function utworzPasekKontekstu(zrodlo: ZrodloZnacznika): PasekKontekstu {
  const projekt = znacznik('Projekt');
  const srodowisko = znacznik('Środowisko');
  const silnik = znacznik('Silnik');
  const tryb = znacznik('Tryb');
  const powiekszenie = znacznik('Powiększenie');

  const element = document.createElement('div');
  element.className = 'md-kontekst';
  element.setAttribute('aria-label', 'Pasek kontekstu kanwy');
  oznaczWarstwe(element, 1);
  element.append(
    projekt.element,
    srodowisko.element,
    silnik.element,
    tryb.element,
    powiekszenie.element,
  );

  // Środowisko nie zmienia się w cyklu życia modułu i nie zależy od stanu, więc
  // stoi raz — odświeżanie go co ramkę powtarzałoby tę samą prawdę.
  srodowisko.ustaw(
    'nie podane',
    'Moduł Design jest dostępny w środowiskach WorkSpace i CodeStudio, ale KTÓRE z nich ' +
      'jest bieżące, moduł wie wyłącznie od powłoki. Ani jedna komenda obszaru design nie ' +
      'niesie środowiska, a okno modułu w rdzeniu wskazuje moduł i sesję, nie środowisko. ' +
      'Znacznik nie wpisuje tu nazwy, której nie ma skąd wziąć.',
  );

  return {
    element,

    odswiez() {
      const okno = zrodlo.okno();
      projekt.ustaw(
        okno === '' ? 'okno nieustalone' : okno,
        `Znacznik pokazuje okno modułu w rdzeniu — to ono jest przedmiotem komend obszaru ` +
          `i bez niego zapis kompozycji oraz generowanie odmawiają. ${zrodlo.opisOkna()}`,
      );

      const kanalow = zrodlo.silnikow();
      silnik.ustaw(
        kanalow === 0 ? 'brak kanału obrazowego' : `${String(kanalow)} do wyboru`,
        kanalow === 0
          ? 'W rejestrze kanałów modelu nie ma ani jednego kanału obrazowego, więc generowanie ' +
            'odmówi, nazywając ten brak. Kanały zakłada się poza modułem, w konfiguracji kont.'
          : 'Selektor silnika stoi w oknie Prompt Builder, przy polach promptu. Pasek go nie ' +
            'dubluje: dwie kontrolki nastawiające jeden kanał rozeszłyby się przy pierwszej ' +
            'zmianie.',
      );

      tryb.ustaw(
        zrodlo.tryb(),
        'Selektor trybu stoi nad kanwą, w znaczniku „Tryb narzędzia kanwy". Tam też stoi ' +
          'zdanie o tym, co wybrany tryb w tej budowie robi.',
      );

      powiekszenie.ustaw(
        `${String(Math.round(zrodlo.powiekszenie() * 100))}%`,
        'Powiększenie zmienia się przyciskami przybornika kanwy: powiększ, pomniejsz, widok ' +
          'naturalny. Znacznik pokazuje nastawę, nie zmienia jej.',
      );
    },
  };
}

/**
 * Znacznik kontekstowy — pigułka klikalna zawsze.
 *
 * Przycisk, nie napis: opracowanie stanowi, że znacznik kontekstowy otwiera się
 * kliknięciem, a element nieklikalny nie miałby jak tego spełnić. Zdanie idzie
 * dymkiem, bo znacznik ma zostać pigułką, a nie rozrosnąć się w akapit.
 */
function znacznik(nazwa: string): Znacznik {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = 'dn-plakietka md-kontekst__znacznik';
  element.dataset['znacznik'] = nazwa.toLowerCase();

  let zdanie = '';
  element.addEventListener('click', () => {
    pokazKomunikat({ tytul: `${nazwa} — znacznik kontekstu`, tresc: zdanie, waga: 'info' });
  });

  return {
    element,
    ustaw(wartosc, tresc) {
      zdanie = tresc;
      element.textContent = `${nazwa}: ${wartosc}`;
      // Powód czytany przez technologie wspomagające bez naciskania.
      element.setAttribute('aria-description', tresc);
    },
  };
}
