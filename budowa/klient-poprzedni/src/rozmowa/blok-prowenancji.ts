import { blokKodu, utworzBlokZwijany, wierszDanych, type BlokZwijany } from './blok-zwijany';
import { NAPISY, POLA_PROWENANCJI } from './etykiety-rozmowy';
import { wierszWywolania, type Prowenancja } from './prowenancja';

/** Blok pokazujący, co dokładnie poszło do modelu w tej turze, przy jednym wpisie rozmowy w tym samym oknie. */
export interface BlokProwenancji {
  /** Element montowany we wpisie. */
  element: HTMLElement;
  /** Wypełnia blok prowenancją albo go ukrywa, gdy jej brak. */
  aktualizuj(prowenancja: Prowenancja | null): void;
}

/** Blok pokazuje to, co rdzeń nadał pierwszym fragmentem strumienia w tej samej turze, zwinięty domyślnie. */
export function utworzBlokProwenancji(): BlokProwenancji {
  const blok: BlokZwijany = utworzBlokZwijany({
    tytul: NAPISY.prowenancja,
    ikona: 'oko',
    odmiana: 'akcent',
  });

  function aktualizuj(prowenancja: Prowenancja | null): void {
    if (prowenancja === null) {
      blok.pokaz(false);
      return;
    }
    blok.pokaz(true);
    blok.ustawPodtytul(podtytul(prowenancja));
    blok.tresc.replaceChildren(...zawartosc(prowenancja));
  }

  return { element: blok.element, aktualizuj };
}

/** Podtytuł widoczny bez rozwijania bloku: kanał wywołania, model, tryb uprawnień oraz numer tej próby. */
function podtytul(p: Prowenancja): string {
  // Kanał, obecny tylko dla kanałów sieciowych, idzie na początek, żeby Operator widział go od razu.
  const czesci = [p.kanal, p.model, p.trybUprawnien].filter((czesc) => czesc.length > 0);
  if (p.zrodloKonta.length > 0 && p.konto.length > 0) {
    // Skąd tożsamość: Operator widzi drogę nastawy, nie tylko wynik.
    czesci.push(`konto ${p.konto} (${p.zrodloKonta})`);
  }
  if (p.proba > 1) czesci.push(`próba ${p.proba}`);
  if (p.powodProby.length > 0) czesci.push(p.powodProby);
  return czesci.join(' · ');
}

/** Treść bloku prowenancji: kanał wywołania, argumenty, prompt systemowy, katalogi oraz parametry wywołania. */
function zawartosc(p: Prowenancja): HTMLElement[] {
  const czesci: HTMLElement[] = [];

  // Kanał wywołania — sekcja obecna tylko wtedy, gdy prowenancja niesie pola dodatkowe kanału.
  const kanal = polaKanalu(p)
    .filter(([, wartosc]) => wartosc.length > 0)
    .map(([nazwa, wartosc]) => wierszDanych(nazwa, wartosc));
  if (kanal.length > 0) {
    czesci.push(podnaglowek('Kanał wywołania'), ...kanal);
  }

  const wywolanie = wierszWywolania(p);
  if (wywolanie.length > 0) {
    czesci.push(podnaglowek('Wiersz wywołania'), blokKodu(wywolanie));
  }

  if (p.promptSystemowy.length > 0) {
    czesci.push(podnaglowek(naglowekPromptu(p)), blokKodu(p.promptSystemowy));
  }

  if (p.katalogi.length > 0) {
    czesci.push(podnaglowek('Katalogi robocze'), blokKodu(p.katalogi.join('\n')));
  }

  if (p.konfiguracjaMCP.length > 0) {
    czesci.push(podnaglowek('Konfiguracja MCP'), blokKodu(p.konfiguracjaMCP.join('\n')));
  }

  const dane = POLA_PROWENANCJI.map(([klucz, nazwa]) => [nazwa, String(p[klucz as keyof Prowenancja] ?? '')] as const)
    .filter(([, wartosc]) => wartosc.length > 0)
    .map(([nazwa, wartosc]) => wierszDanych(nazwa, wartosc));
  if (dane.length > 0) {
    czesci.push(podnaglowek('Parametry wywołania'), ...dane);
  }

  return czesci;
}

/** Pola tożsamości kanału sieciowego w kolejności prezentacji, wypełnione tylko dla kanałów sieciowych. */
function polaKanalu(p: Prowenancja): ReadonlyArray<readonly [string, string]> {
  return [
    ['Kanał', p.kanal],
    ['Adapter', p.adapter],
    ['Adres punktu końcowego', p.adres],
    ['Środowisko wykonania', p.srodowiskoWykonania],
  ];
}

/** Nagłówek promptu systemowego niesie skład warstw nakładki złożonych na ten sam prompt systemowy tury. */
function naglowekPromptu(p: Prowenancja): string {
  if (p.warstwyNakladki.length === 0) return 'Prompt systemowy';
  return `Prompt systemowy — warstwy: ${p.warstwyNakladki.join(' → ')}`;
}

/** Podnagłówek sekcji wewnątrz bloku prowenancji, wypisywany zawsze nad każdą kolejną sekcją tego bloku. */
function podnaglowek(tekst: string): HTMLElement {
  const element = document.createElement('h4');
  element.className = 'dc-blok__podnaglowek dn-etykieta-wersalikowa';
  element.textContent = tekst;
  return element;
}
