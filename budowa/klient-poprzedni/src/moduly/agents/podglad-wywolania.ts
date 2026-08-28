import { Command } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';

/**
 * Podgląd wywołania podaje wiersz polecenia oraz prompt systemowy, z jakimi
 * ruszy proces modelu, zanim tura zostanie wysłana. Rdzeń liczy go tymi samymi
 * funkcjami, którymi wykonuje wysłanie wiadomości, wraz z nałożeniem eksperta.
 */
const ETYKIETA = 'Podgląd wywołania';
const WYJASNIENIE =
  'Wiersz polecenia i prompt systemowy, z jakimi ruszy proces modelu w tym oknie. ' +
  'Liczone tą samą drogą, którą idzie tura.';
const BEZ_OKNA =
  'Podaj identyfikator okna komunikacji. Wywołanie liczy się dla okna, a moduł ' +
  'Agents żadnego okna rozmowy nie zna — bez wskazania podgląd pokazywałby cudze.';

export interface PodgladWywolania {
  element: HTMLElement;
}

export function utworzPodgladWywolania(kanal: Kanal): PodgladWywolania {
  const element = document.createElement('section');
  element.className = 'da-agents__podglad';

  const naglowek = document.createElement('h3');
  naglowek.className = 'da-agents__naglowek';
  naglowek.textContent = ETYKIETA;

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'da-agents__wyjasnienie';
  wyjasnienie.textContent = WYJASNIENIE;

  const pasek = document.createElement('div');
  pasek.className = 'da-agents__pasek';

  const etykietaOkna = document.createElement('label');
  etykietaOkna.className = 'da-agents__etykieta';
  etykietaOkna.textContent = 'Okno komunikacji';

  const poleOkna = document.createElement('input');
  poleOkna.type = 'text';
  poleOkna.className = 'da-agents__pole';
  poleOkna.placeholder = 'okn_…';
  etykietaOkna.append(poleOkna);

  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  // Klasy przycisków modułu to `dn-btn` z odmianą; innych nazw arkusze nie
  // pokrywają.
  przycisk.className = 'dn-btn dn-btn--sm dn-btn--atrament';
  przycisk.textContent = 'Odczytaj wywołanie dla tego okna';

  pasek.append(etykietaOkna, przycisk);

  const stan = document.createElement('p');
  stan.className = 'da-agents__stan';
  stan.setAttribute('role', 'status');
  stan.textContent = BEZ_OKNA;

  const wiersz = document.createElement('pre');
  wiersz.className = 'da-agents__wiersz';
  wiersz.hidden = true;

  const prompt = document.createElement('pre');
  prompt.className = 'da-agents__prompt';
  prompt.hidden = true;

  element.append(naglowek, wyjasnienie, pasek, stan, wiersz, prompt);

  /** Chowa wynik poprzedniego odczytu — stara wartość pod nowym pytaniem kłamie. */
  function wyczysc(): void {
    wiersz.hidden = true;
    prompt.hidden = true;
    wiersz.textContent = '';
    prompt.textContent = '';
  }

  przycisk.addEventListener('click', () => {
    const idOkna = poleOkna.value.trim();
    if (idOkna === '') {
      wyczysc();
      stan.textContent = BEZ_OKNA;
      return;
    }
    wyczysc();
    stan.textContent = 'Odczyt z rdzenia…';
    kanal.wyslij(Command.ConfigExplainGet, { windowId: idOkna }, (wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Okno podaje treść odmowy rdzenia; pusty wiersz podany jako wywołanie
        // byłby odpowiedzią nieprawdziwą.
        stan.textContent = wynik.blad?.message ?? 'Rdzeń nie oddał wywołania dla tego okna.';
        return;
      }
      const argv = wynik.wynik.argv ?? [];
      const systemPrompt = wynik.wynik.systemPrompt ?? '';

      stan.textContent =
        `${argv.length} ${odmianaPozycji(argv.length)} wiersza polecenia · ` +
        `prompt systemowy ${systemPrompt.length} zn.`;

      wiersz.textContent = argv.join(' ');
      wiersz.hidden = argv.length === 0;

      // Pusty prompt jest odpowiedzią, nie brakiem odpowiedzi, więc pole
      // pozostaje widoczne.
      prompt.textContent =
        systemPrompt === '' ? 'Prompt systemowy pusty — żadna warstwa nie wniosła treści.' : systemPrompt;
      prompt.hidden = false;
    });
  });

  return { element };
}

/**
 * Odmiana rzeczownika liczonego według reguł polskiej liczby mnogiej: liczebnik
 * jeden bierze mianownik liczby pojedynczej, zakończenia od dwóch do czterech
 * poza nastką biorą mianownik liczby mnogiej, a pozostałe biorą dopełniacz.
 */
function odmianaPozycji(ile: number): string {
  if (ile === 1) return 'pozycja';
  const setki = ile % 100;
  const jednosci = ile % 10;
  const mnoga = jednosci >= 2 && jednosci <= 4 && (setki < 12 || setki > 14);
  return mnoga ? 'pozycje' : 'pozycji';
}
