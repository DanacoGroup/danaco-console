import { pokazKomunikat } from '../../aplikacja/komunikaty';
import { oznaczWarstwe } from './warstwy-designu';

/**
 * Źródło wartości znaczników paska kontekstu kanwy Design Board: okno modułu,
 * liczba kanałów obrazowych rejestru, tryb narzędzia i powiększenie. Pasek
 * własnych selektorów nie stawia, tylko nazywa miejsce, w którym każdy stoi.
 */
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

/**
 * Jeden znacznik paska: przycisk osadzany w pasku oraz nastawa wartości
 * bieżącej wraz ze zdaniem pokazywanym po naciśnięciu znacznika.
 */
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

  // Srodowisko nie zmienia sie w cyklu zycia modulu, wiec znacznik stoi raz.
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
 * Składa znacznik kontekstowy jako pigułkę klikalną zawsze. Naciśnięcie
 * pokazuje zdanie komunikatem, a to samo zdanie idzie w atrybut opisu, żeby
 * technologie wspomagające czytały je bez naciskania.
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
