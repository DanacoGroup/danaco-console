import { MessageRole } from '../../../shared/contract';
import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { godzinaWpisu, nazwaNadawcy, wpisModelu, type Wpis } from './wpis';

/**
 * Klasy wpisu według roli nadawcy: `user → czlowiek`, `assistant → inteligencja`,
 * `system → system`, `tool → system`. Wygląd niesie `komponenty/wpis.css`; okno
 * nie powtarza ani jednej z tych reguł.
 */
const KLASA_WPISU: Record<MessageRole, string> = {
  [MessageRole.User]: 'dn-wpis dn-wpis--czlowiek',
  [MessageRole.Assistant]: 'dn-wpis dn-wpis--inteligencja',
  [MessageRole.System]: 'dn-wpis dn-wpis--system',
  [MessageRole.Tool]: 'dn-wpis dn-wpis--system',
};

/**
 * Znak w medalionie wpisu. Rozróżnienie nadawcy niesie zawsze komplet trzech
 * kanałów — ikona, etykieta, kreska krawędzi — bo rola nigdy nie może stać na
 * samym kolorze (`komponenty/wpis.css`). Wynik narzędzia dzieli z systemem klasę
 * neutralną, ale nie ikonę: klasa mówi o randze wypowiedzi, ikona o jej źródle.
 */
const IKONA_WPISU: Record<MessageRole, NazwaIkony> = {
  [MessageRole.User]: 'uzytkownik',
  [MessageRole.Assistant]: 'agent',
  [MessageRole.System]: 'info',
  [MessageRole.Tool]: 'polecenie',
};

/** Obszar historii okna komunikacji. */
export interface Historia {
  /** Element montowany w oknie. */
  element: HTMLElement;
  /** Dokłada wpis na koniec historii. */
  dopisz(wpis: Wpis): void;
  /** Dokłada treść do ostatniego wpisu tej samej persony (strumień modelu). */
  dopiszDoOstatniego(persona: string, fragment: string): void;
  /** Liczba wpisów w historii. */
  liczbaWpisow(): number;
}

export function utworzHistorie(): Historia {
  const element = document.createElement('section');
  element.className = 'dc-okno__historia';
  element.setAttribute('aria-label', 'Historia okna komunikacji');

  let ostatniaPersona: string | null = null;
  let ostatniaTresc: HTMLElement | null = null;
  let liczba = 0;

  function przewinNaKoniec(): void {
    element.scrollTop = element.scrollHeight;
  }

  function dopisz(wpis: Wpis): void {
    const wezly = zbudujWpis(wpis);
    element.append(wezly.ramka);
    ostatniaPersona = wpis.persona;
    ostatniaTresc = wezly.tresc;
    liczba += 1;
    przewinNaKoniec();
  }

  function dopiszDoOstatniego(persona: string, fragment: string): void {
    if (ostatniaTresc === null || ostatniaPersona !== persona) {
      dopisz(wpisModelu(fragment, persona));
      return;
    }
    ostatniaTresc.textContent = `${ostatniaTresc.textContent ?? ''}${fragment}`;
    przewinNaKoniec();
  }

  return {
    element,
    dopisz,
    dopiszDoOstatniego,
    liczbaWpisow: () => liczba,
  };
}

/**
 * Etykieta tożsamości wpisu — rola, a po niej persona.
 *
 * Nadawcę podpisuje jedna etykieta (`.dn-wpis-nadawca`) w postaci „Model ·
 * Redaktor”. Persona równa nazwie roli zostaje pominięta, bo „Operator ·
 * Operator” nie niesie więcej niż „Operator”.
 */
function etykietaTozsamosci(wpis: Wpis): string {
  const rola = nazwaNadawcy(wpis.nadawca);
  return wpis.persona === rola ? rola : `${rola} · ${wpis.persona}`;
}

/**
 * Buduje węzły jednego wpisu: medalion, pasek tożsamości oraz treść.
 *
 * Medalion zajmuje pierwszą kolumnę dwukolumnowej siatki `.dn-wpis`, o szerokości
 * awatara. Bez niego pasek tożsamości wpadłby w tę kolumnę i został zgnieciony
 * do jej szerokości.
 */
function zbudujWpis(wpis: Wpis): { ramka: HTMLElement; tresc: HTMLElement } {
  const ramka = document.createElement('article');
  ramka.className = KLASA_WPISU[wpis.nadawca];

  const medalion = document.createElement('span');
  medalion.className = 'dn-wpis-medalion';
  medalion.append(elementIkony(IKONA_WPISU[wpis.nadawca], { rozmiar: 16 }));

  const tozsamosc = document.createElement('header');
  tozsamosc.className = 'dn-wpis-tozsamosc';

  const nadawca = document.createElement('span');
  nadawca.className = 'dn-wpis-nadawca';
  nadawca.textContent = etykietaTozsamosci(wpis);

  const godzina = document.createElement('time');
  godzina.className = 'dn-wpis-godzina';
  godzina.textContent = godzinaWpisu(wpis);

  tozsamosc.append(nadawca, godzina);

  const tresc = document.createElement('p');
  tresc.className = 'dn-wpis-tresc';
  tresc.textContent = wpis.tresc;

  ramka.append(medalion, tozsamosc, tresc);
  return { ramka, tresc };
}
