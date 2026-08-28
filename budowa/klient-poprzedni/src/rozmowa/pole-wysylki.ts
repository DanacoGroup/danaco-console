import { elementIkony } from '../ikony/ikony';
import { NAPISY } from './etykiety-rozmowy';
import { opisStanu, type StanRozmowy } from './stan-rozmowy';
import type { WykazPoUkosniku } from './wykaz-po-ukosniku';

/**
 * Zdarzenia zgłaszane przez pole wysyłki wywołującemu widokowi, obejmujące wysłanie i
 * przerwanie tury.
 */
export interface ObslugaWysylki {
  /** Operator zatwierdził wypowiedź. */
  naWyslanie(tekst: string): void;
  /** Operator zażądał przerwania tury. */
  naPrzerwanie(): void;
}

/**
 * Ustawienia pola wysyłki, których obsługa wysyłania i przerywania tury nie potrafi
 * wyprowadzić samodzielnie.
 */
export interface OpcjePolaWysylki {
  // Pasek zlecenia pod polem wypowiedzi, opcjonalny — powstaje dopiero w powłoce produktu.
  pasekZlecenia?: HTMLElement;
  // Wykaz po ukośniku nad polem wypowiedzi; polem sterują wyłącznie cztery ruchy klawiatury.
  wykazUkosnika?: WykazPoUkosniku;
}

/**
 * Pole wpisywania wypowiedzi Operatora wraz z przyciskiem wysyłki oraz przyciskiem
 * zatrzymania tury odpowiedzi.
 */
export interface PoleWysylki {
  /** Element montowany w widoku rozmowy. */
  element: HTMLElement;
  /** Odświeża wskaźnik stanu wysyłania. */
  pokazStan(stan: StanRozmowy): void;
  /** Ustawia ognisko na polu wpisywania. */
  ustawOgnisko(): void;
  /** Dokłada gotowe polecenie do treści wypowiedzi; nie wysyła go. */
  wstaw(tekst: string): void;
  // Odbitka wypowiedzi z innego połączenia, pokazywana w polu na krótko i znikająca bez
  // wysłania.
  odbijWypowiedzZZewnatrz(tekst: string): void;
}

/**
 * Czas, przez jaki odbitka cudzej wypowiedzi stoi w polu, zanim zniknie bez śladu w treści
 * wpisywania.
 */
const CZAS_ODBITKI_MS = 1200;

/**
 * Buduje pole wypowiedzi Operatora wraz z przyciskiem wysyłki i zawsze klikalnym
 * przyciskiem zatrzymania.
 */
export function utworzPoleWysylki(
  obsluga: ObslugaWysylki,
  opcje: OpcjePolaWysylki = {},
): PoleWysylki {
  const element = document.createElement('form');
  // Wiersz dolny okna: wyściółka, tło panelu, bez kreski górnej.
  element.className = 'dn-rozmowa-dol dc-wysylka';

  // Powłoka pola niesie ramkę, promień i ognisko, nie obszar tekstu.
  const powloka = document.createElement('div');
  powloka.className = 'dn-prompt';

  // Grot — sygnatura wejścia, znak dekoracyjny poza drzewem dostępności.
  const grot = document.createElement('span');
  grot.className = 'dn-prompt-grot';
  grot.setAttribute('aria-hidden', 'true');
  grot.textContent = '❯';

  const pole = document.createElement('textarea');
  // Obszar tekstu bez ramki, rosnący między dwoma progami wysokości wewnątrz powłoki.
  pole.className = 'dn-prompt-obszar dc-wysylka__pole';
  pole.rows = 1;
  pole.placeholder = NAPISY.polePodpowiedz;
  pole.setAttribute('aria-label', NAPISY.polePodpowiedz);

  const pasek = document.createElement('div');
  // Pasek akcji pod polem: odstęp 4 px, odsunięcie od pola 8 px.
  pasek.className = 'dn-rozmowa-akcje dc-wysylka__pasek';

  const wskaznik = document.createElement('span');
  wskaznik.className = 'dc-wysylka__wskaznik';
  wskaznik.setAttribute('role', 'status');

  // Podpowiedź wykazu po ukośniku w rzędzie akcji, nie w placeholderze pola.
  const podpowiedzWykazu = document.createElement('span');
  podpowiedzWykazu.className = 'dc-wysylka__podpowiedz';
  podpowiedzWykazu.setAttribute('role', 'status');
  podpowiedzWykazu.hidden = true;
  podpowiedzWykazu.textContent = NAPISY.wykazPodpowiedz;

  const przerwij = przycisk(NAPISY.przerwij, 'zatrzymaj', 'dn-btn dn-btn--zarys');
  const wyslij = przycisk(NAPISY.wyslij, 'wyslij', 'dn-btn dn-btn--atrament');
  wyslij.type = 'submit';

  // Odbitka cudzej wypowiedzi w osobnej warstwie nad polem, poza jego treścią i drzewem
  // wejścia.
  const odbitka = document.createElement('output');
  odbitka.className = 'dc-wysylka__odbitka';
  odbitka.hidden = true;
  odbitka.setAttribute('aria-live', 'polite');
  let zegarOdbitki = 0;

  powloka.append(grot, pole, odbitka, wyslij);
  // Pasek zlecenia stoi pierwszy w rzędzie, przed wskaźnikiem stanu i przyciskiem
  // zatrzymania.
  if (opcje.pasekZlecenia !== undefined) pasek.append(opcje.pasekZlecenia);
  pasek.append(podpowiedzWykazu, wskaznik, przerwij);
  // Wykaz po ukośniku poprzedza powłokę pola, bo otwiera się w górę nad rozmową.
  const wykaz = opcje.wykazUkosnika;
  if (wykaz !== undefined) element.append(wykaz.element);
  element.append(powloka, pasek);

  function zatwierdz(): void {
    const tresc = pole.value;
    if (tresc.trim().length === 0) return;
    pole.value = '';
    obsluga.naWyslanie(tresc);
  }

  /* Wykaz po ukośniku — wyzwalacz z klawiatury.

     Wykaz pojawia się natychmiast po wpisaniu „/", nie po dopisaniu liter.
     Wyzwalaczem jest zdarzenie `input`, a nie `keydown`: znak jest już w polu,
     więc warunek „ukośnik stoi na początku treści" da się sprawdzić na treści,
     a nie zgadywać z klawisza (martwe klawisze, wklejenie, IME).

     Ukośnik w środku zdania nie otwiera niczego — ścieżka `/home/ubuntu`
     wpisana w wypowiedź jest tekstem, nie komendą. */

  /** Czy treść pola jest wpisem po ukośniku (i jaka jest jego fraza). */
  function frazaUkosnika(tresc: string): string | null {
    if (!tresc.startsWith('/')) return null;
    const reszta = tresc.slice(1);
    // Nowy wiersz kończy komendę: to już wypowiedź wielowierszowa, nie wpis.
    return reszta.includes('\n') ? null : reszta;
  }

  function odswiezPodpowiedz(): void {
    podpowiedzWykazu.hidden = wykaz === undefined || !wykaz.otwarty();
  }

  if (wykaz !== undefined) {
    wykaz.naZmiane(() => odswiezPodpowiedz());
    pole.addEventListener('input', () => {
      const fraza = frazaUkosnika(pole.value);
      if (fraza === null) {
        // Skasowanie ukośnika zwija wykaz — tak samo jak Escape.
        if (wykaz.otwarty()) wykaz.zamknij();
        odswiezPodpowiedz();
        return;
      }
      if (!wykaz.otwarty()) wykaz.otworz();
      wykaz.ustawFraze(fraza);
      odswiezPodpowiedz();
    });
  }

  element.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    zatwierdz();
  });

  przerwij.addEventListener('click', () => obsluga.naPrzerwanie());

  // Enter wysyła, Shift+Enter dokłada wiersz; przy otwartym wykazie te same klawisze sterują
  // nim.
  pole.addEventListener('keydown', (zdarzenie) => {
    if (wykaz !== undefined && wykaz.otwarty()) {
      if (zdarzenie.key === 'ArrowDown' || zdarzenie.key === 'ArrowUp') {
        zdarzenie.preventDefault();
        wykaz.przesun(zdarzenie.key === 'ArrowDown' ? 1 : -1);
        return;
      }
      if (zdarzenie.key === 'Escape') {
        zdarzenie.preventDefault();
        zdarzenie.stopPropagation();
        wykaz.zamknij();
        odswiezPodpowiedz();
        return;
      }
      if (zdarzenie.key === 'Enter' && !zdarzenie.shiftKey) {
        zdarzenie.preventDefault();
        // Pole czyści się po zatwierdzeniu pozycji wykazu, bo wybór nie jest wypowiedzią do
        // modelu.
        const wybrano = wykaz.zatwierdz();
        odswiezPodpowiedz();
        if (wybrano) {
          pole.value = '';
          return;
        }
        // Wykaz przycięty do zera treścią pola idzie zwykłą drogą wysyłki.
        zatwierdz();
        return;
      }
    }
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    zatwierdz();
  });

  return {
    element,
    pokazStan(stan) {
      wskaznik.textContent = opisStanu(stan);
      // `aria-busy` niesie stan pracy bez odbierania klikalności.
      wyslij.setAttribute('aria-busy', String(stan.wysyla));
      element.dataset['wysyla'] = String(stan.wysyla);
    },
    ustawOgnisko: () => pole.focus(),

    odbijWypowiedzZZewnatrz(tekst) {
      const tresc = tekst.trim();
      if (tresc.length === 0) return;
      // Operator pisze — pole należy do niego i nikt go z niego nie wyrywa.
      if (pole.value.length > 0 || document.activeElement === pole) return;

      // Odbitka mieszka w warstwie nad polem, nie w jego treści, by uniknąć podwójnego wysłania.
      odbitka.textContent = tresc;
      odbitka.hidden = false;
      powloka.dataset['odbitka'] = 'tak';
      window.clearTimeout(zegarOdbitki);
      zegarOdbitki = window.setTimeout(() => {
        odbitka.hidden = true;
        odbitka.textContent = '';
        delete powloka.dataset['odbitka'];
      }, CZAS_ODBITKI_MS);
    },

    // Narzędzia paska i akcje modułu wstawiają polecenie do pola, nie wysyłają go od razu.
    wstaw(tekst) {
      const zastana = pole.value;
      const rozdzielnik = zastana.length > 0 && !zastana.endsWith('\n') ? '\n' : '';
      pole.value = `${zastana}${rozdzielnik}${tekst}`;
      pole.focus();
      pole.setSelectionRange(pole.value.length, pole.value.length);
    },
  };
}

/**
 * Przycisk paska wysyłki złożony z ikony oraz etykiety słownej, osadzony w rzędzie
 * sterowania oknem.
 */
function przycisk(etykieta: string, ikona: 'wyslij' | 'zatrzymaj', klasa: string): HTMLButtonElement {
  const element = document.createElement('button');
  element.type = 'button';
  element.className = klasa;
  element.append(elementIkony(ikona, { rozmiar: 16, klasa: 'dn-ikona' }));

  const napis = document.createElement('span');
  napis.textContent = etykieta;
  element.append(napis);
  return element;
}
