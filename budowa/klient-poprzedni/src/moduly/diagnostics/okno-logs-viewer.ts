import {
  LogLevel,
  type DiagnosticsLogQueryRequest,
  type LogEntry,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  przelacznik,
  przelacznikWidoku,
  przestaw,
  przyciskAkcji as przycisk,
  przyciskBezKomendy,
  wiersz,
  wykaz,
} from '../../modele/kontrolki-formularza-braki';
import { utworzWyborZMenu, type WyborZMenu } from '../apps/wybor-z-menu';
import { powodBezKomendy } from './braki-kontraktu';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import type { StanDiagnostyki } from './stan-diagnostyki';
import { utworzStanTresci, type StanTresci } from './stany-okna';
import type { ZrodloDiagnostics } from './zrodlo-diagnostics';

/**
 * Logs Viewer to okno monitora modułu Diagnostics: przeszukuje dziennik rdzenia po wzorcu,
 * poziomie i źródle, w tym wpisy będące jedynym śladem odmowy wykonania komendy.
 */
export interface OknoLogsViewera {
  element: HTMLElement;
  odswiez(): void;
  /** Zamyka nasłuch zmiany zakresu wspólnego stanu modułu. */
  zamknij(): void;
}

export function utworzOknoLogsViewera(
  zrodlo: ZrodloDiagnostics,
  stan: StanDiagnostyki,
): OknoLogsViewera {
  const rama = utworzRameOkna({
    tytul: 'Logs Viewer',
    rola: 'monitor',
    kod: 'logs-viewer',
    przeznaczenie:
      'Przeszukiwanie dziennika rdzenia po wzorcu, poziomie i źródle — jedyny ślad odmowy komendy, gdy zapis wiersza błędu się nie powiedzie.',
    modul: 'Diagnostics',
    przedrostek: 'dg',
  });
  const tresc = utworzStanTresci();
  const powierzchnia = zlozPowierzchnieLogow(rama, tresc.element);
  const przypiete = new Set<string>();
  let ostatnieWpisy: readonly LogEntry[] = [];

  function odczytaj(): void {
    const zadanie = zbudujZadanie(powierzchnia, stan.zakres());
    if (zadanie === null) {
      // Odmowa okna sprząta jak odmowa rdzenia: wpisy przestają być bieżące, potwierdzenie gaśnie.
      ostatnieWpisy = [];
      tresc.potwierdzenie('', true);
      tresc.blad('Wzorzec nie jest poprawnym wyrażeniem regularnym — okno zatrzymało wyszukanie i nie pytało rdzenia.');
      return;
    }
    tresc.ladowanie('Przeszukiwanie dziennika…');
    void zrodlo.przeszukajDziennik(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Odmowa nie zostawia wpisów poprzedniego odczytu jako aktualnych — eksport oddałby zakres odrzucony.
        ostatnieWpisy = [];
        tresc.blad(zdanieNiepowodzenia('przeszukania dziennika', wynik.powod), wynik.blad);
        return;
      }
      ostatnieWpisy = wynik.wynik.entries;
      rysuj();
      // Wynik przycięty granicą bufora nie jest sukcesem — widoczne to w tonie potwierdzenia, nie w zdaniu.
      tresc.potwierdzenie(zdaniePodsumowania(wynik.wynik), wynik.wynik.truncated !== true);
      if (powierzchnia.autoPrzewijanie.dataset['wlaczony'] === 'true') przewinDoKonca(tresc);
    });
  }

  function rysuj(): void {
    const podzielony = powierzchnia.podzielWidok.dataset['wlaczony'] === 'true';
    const widoczne = podzielony ? ostatnieWpisy.filter((wpis) => przypiete.has(wpis.id)) : ostatnieWpisy;
    if (widoczne.length === 0) {
      tresc.pusto(
        podzielony
          ? 'Podział widoku włączony, ale żaden wpis nie jest przypięty — przypnij wpis albo wyłącz podział.'
          : 'Żaden wpis dziennika nie spełnia obecnych warunków wyszukania.',
      );
      return;
    }
    tresc.tresc().append(listaWpisow(widoczne, przypiete, przelaczPrzypiecie));
  }

  /** Przerysowanie klienckie nie jest nowym odczytem, więc potwierdzenie poprzedniego odczytu gaśnie. */
  function rysujKlienckie(): void {
    tresc.potwierdzenie('', true);
    rysuj();
  }

  function przelaczPrzypiecie(id: string): void {
    if (przypiete.has(id)) przypiete.delete(id);
    else przypiete.add(id);
    rysujKlienckie();
  }

  function eksportuj(): void {
    if (ostatnieWpisy.length === 0) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — wykaz wpisów jest pusty.', false);
      return;
    }
    pobierzPlik('dziennik-zakres.md', raportWpisow(ostatnieWpisy), 'text/markdown');
    tresc.potwierdzenie('Zakres dziennika pobrany jako plik Markdown.', true);
  }

  podepnijAkcje(powierzchnia, { odczytaj, eksportuj, rysuj: rysujKlienckie });
  const odsubskrybujZakres = stan.naZmiane(odczytaj);

  return { element: rama.element, odswiez: odczytaj, zamknij: odsubskrybujZakres };
}

/** Zdanie podsumowania odczytu dziennika: podaje liczbę oddanych wpisów oraz pola total i truncated wprost, bez przemilczenia. */
function zdaniePodsumowania(wynik: { entries: LogEntry[]; total?: number; truncated?: boolean }): string {
  const oddane = `Odczytano ${wynik.entries.length} wpis(y/ów)`;
  const calosc = wynik.total === undefined ? '' : ` z ${wynik.total} spełniających warunki`;
  const przyciecie =
    wynik.truncated === true
      ? ' — WYNIK PRZYCIĘTY granicą bufora rdzenia, zawęź wyszukanie, by zobaczyć resztę.'
      : '.';
  return `${oddane}${calosc}${przyciecie}`;
}

/** Kontrolki paska akcji Logs Viewera: przycisk szukania, przycisk eksportu oraz dwa przełączniki widoku. */
interface AkcjeLogow {
  szukajPrzycisk: HTMLButtonElement;
  eksportujPrzycisk: HTMLButtonElement;
  autoPrzewijanie: HTMLButtonElement;
  podzielWidok: HTMLButtonElement;
}

/**
 * Składa pasek akcji: wyszukanie, eksport zakresu i dwa przełączniki czysto klienckie —
 * auto-przewijanie oraz podział widoku. Dwie pozycje bez pokrycia w kontrakcie stoją jako
 * przycisk bez komendy.
 */
function zlozAkcjeLogow(gospodarz: HTMLElement): AkcjeLogow {
  const szukajPrzycisk = przycisk('Szukaj w dzienniku', 'dn-btn dn-btn--atrament');
  const eksportujPrzycisk = przycisk('Eksportuj zakres');
  const autoPrzewijanie = przelacznikWidoku('Auto-przewijanie', false);
  const podzielWidok = przelacznikWidoku('Podziel widok (wpisy przypięte)', false);

  gospodarz.append(
    szukajPrzycisk,
    eksportujPrzycisk,
    autoPrzewijanie,
    podzielWidok,
    przyciskBezKomendy(
      'Przekaż do Chat Window',
      powodBezKomendy('Podanie fragmentu dziennika oknu rozmowy wymagałoby komendy przekazania.'),
    ),
    przyciskBezKomendy(
      'Przekaż do Errors Panel',
      powodBezKomendy('Podanie wpisu dziennika Errors Panelowi wymagałoby komendy przekazania.'),
    ),
  );
  return { szukajPrzycisk, eksportujPrzycisk, autoPrzewijanie, podzielWidok };
}

/** Kontrolki filtra dziennika: wzorzec wyszukiwania, przełącznik wyrażenia regularnego, poziom, źródło, deduplikacja i górna granica liczby wpisów. */
interface FiltrLogow {
  wzorzec: HTMLInputElement;
  regex: HTMLInputElement;
  poziom: WyborZMenu;
  zrodlo: HTMLInputElement;
  deduplikuj: HTMLInputElement;
  granica: HTMLInputElement;
}

/** Kontrolki okna Logs Viewer: pasek akcji połączony z paskiem filtra dziennika w jedną powierzchnię sterowania. */
interface PowierzchniaLogow extends AkcjeLogow, FiltrLogow {}

/** Pozycje wyliczenia poziomu wpisu dziennika w polu wyboru; pusta wartość oznacza brak zawężenia do jednego poziomu. */
const POZIOMY = [
  { wartosc: '', etykieta: 'Wszystkie poziomy' },
  { wartosc: LogLevel.Error, etykieta: 'Błąd' },
  { wartosc: LogLevel.Warn, etykieta: 'Ostrzeżenie' },
  { wartosc: LogLevel.Info, etykieta: 'Informacja' },
  { wartosc: LogLevel.Debug, etykieta: 'Diagnostyczny' },
] as const;

/**
 * Składa kontrolki okna i osadza pasek narzędzi filtra w ramie. Wzorzec
 * i źródło stoją na pierwszym miejscu — to droga do odnalezienia odrzuconej
 * komendy, nie wygoda wyszukiwania.
 */
function zlozPowierzchnieLogow(
  rama: { akcje: HTMLElement; narzedzia: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
): PowierzchniaLogow {
  const akcje = zlozAkcjeLogow(rama.akcje);
  const wzorzec = pole('Wzorzec wyszukiwania', 'tekst albo wyrażenie regularne');
  const regex = przelacznik('Wzorzec jest wyrażeniem regularnym');
  // Rozwijanie z biblioteki menu-drzewo, nie natywny element select, przez obsadę modułu apps.
  const poziom = utworzWyborZMenu('Poziom wpisu', POZIOMY);
  const zrodloPole = pole('Źródło wpisu', 'np. library.collection.create');
  const deduplikuj = przelacznik('Scal wpisy identyczne licznikiem');
  const granica = poleLiczbowe('Górna granica liczby wpisów', 'domyślnie wg rdzenia');

  rama.narzedzia.append(
    wiersz('Wzorzec', wzorzec, { klasa: 'dg-wiersz', objasnienie: 'Puste pole nie zawęża wyszukania.' }),
    wiersz('Regex', regex, { klasa: 'dg-wiersz', objasnienie: 'Wzorzec niepoprawny jako wyrażenie zatrzymuje wyszukanie z odmową widoczną w oknie.' }),
    wiersz('Poziom', poziom.element, { klasa: 'dg-wiersz' }),
    wiersz('Źródło', zrodloPole, {
      klasa: 'dg-wiersz',
      objasnienie: 'Nazwa komendy odrzuconej przez rdzeń — jedyny trop odmowy bez wiersza błędu.',
    }),
    wiersz('Deduplikacja', deduplikuj, { klasa: 'dg-wiersz', objasnienie: 'Licznik wystąpień pokazuje się przy wpisie scalonym.' }),
    wiersz('Granica', granica, { klasa: 'dg-wiersz' }),
  );
  rama.cialo.append(stanTresci);

  return { ...akcje, wzorzec, regex, poziom, zrodlo: zrodloPole, deduplikuj, granica };
}

/** Podpina pasek akcji do czynności okna: wyszukania, eksportu oraz obu przełączników widoku klienckiego. */
function podepnijAkcje(
  powierzchnia: PowierzchniaLogow,
  obsluga: { odczytaj: () => void; eksportuj: () => void; rysuj: () => void },
): void {
  powierzchnia.szukajPrzycisk.addEventListener('click', obsluga.odczytaj);
  powierzchnia.eksportujPrzycisk.addEventListener('click', obsluga.eksportuj);
  powierzchnia.autoPrzewijanie.addEventListener('click', () => przestaw(powierzchnia.autoPrzewijanie));
  powierzchnia.podzielWidok.addEventListener('click', () => {
    przestaw(powierzchnia.podzielWidok);
    obsluga.rysuj();
  });
}

/**
 * Zadanie odczytu dziennika. Oddaje `null`, gdy wzorzec regex jest niepoprawną
 * składnią — wywołujący ma wtedy zatrzymać się z odmową widoczną, zamiast
 * wysyłać żądanie, które rdzeń odrzuci ciszej po swojej stronie.
 */
function zbudujZadanie(
  powierzchnia: FiltrLogow,
  zakres: { od?: number; do?: number },
): DiagnosticsLogQueryRequest | null {
  const wzorzec = powierzchnia.wzorzec.value.trim();
  const regex = powierzchnia.regex.checked;
  if (regex && wzorzec !== '') {
    try {
      // eslint-disable-next-line no-new -- sprawdzenie składni, wynik nie jest używany
      new RegExp(wzorzec);
    } catch {
      return null;
    }
  }
  const zadanie: DiagnosticsLogQueryRequest = {};
  if (wzorzec !== '') zadanie.pattern = wzorzec;
  if (regex) zadanie.regex = true;
  if (powierzchnia.poziom.wartosc() !== '') {
    zadanie.level = powierzchnia.poziom.wartosc() as LogLevel;
  }
  const zrodloWartosc = powierzchnia.zrodlo.value.trim();
  if (zrodloWartosc !== '') zadanie.source = zrodloWartosc;
  if (zakres.od !== undefined) zadanie.fromTime = zakres.od;
  if (zakres.do !== undefined) zadanie.toTime = zakres.do;
  if (powierzchnia.deduplikuj.checked) zadanie.deduplicate = true;
  const granica = Number.parseInt(powierzchnia.granica.value, 10);
  if (Number.isInteger(granica) && granica > 0) zadanie.limit = granica;
  return zadanie;
}

/** Przewija miejsce treści okna do ostatniego wpisu dziennika; obsługuje włączone auto-przewijanie klienckie. */
function przewinDoKonca(tresc: StanTresci): void {
  const miejsce = tresc.element.querySelector('.dg-tresc');
  if (miejsce !== null) miejsce.scrollTop = miejsce.scrollHeight;
}

/**
 * Buduje wykaz wpisów dziennika — wpisy przypięte pierwsze, w kolejności
 * bez zmiany, reszta chronologicznie od najnowszego. Czysta funkcja danych:
 * zna wpisy, przypięcia i wywołanie zwrotne kliknięcia, nie zna źródła.
 */
function listaWpisow(
  wpisy: readonly LogEntry[],
  przypiete: ReadonlySet<string>,
  przelacz: (id: string) => void,
): HTMLElement {
  const posortowane = [...wpisy].sort((a, b) => {
    const aPrzypiety = przypiete.has(a.id) ? 0 : 1;
    const bPrzypiety = przypiete.has(b.id) ? 0 : 1;
    if (aPrzypiety !== bPrzypiety) return aPrzypiety - bPrzypiety;
    return b.timestamp - a.timestamp;
  });
  const lista = wykaz('Wpisy dziennika', 'dg-wykaz');
  for (const wpis of posortowane) lista.append(wierszWpisu(wpis, przypiete.has(wpis.id), przelacz));
  return lista;
}

/** Buduje jeden wiersz wpisu dziennika wraz z jego licznikiem deduplikacji oraz przyciskiem przypięcia. */
function wierszWpisu(wpis: LogEntry, przypiety: boolean, przelacz: (id: string) => void): HTMLElement {
  const element = document.createElement('li');
  element.className = 'dg-wpis';
  element.dataset['poziom'] = wpis.level;
  element.dataset['przypiety'] = String(przypiety);

  const czas = document.createElement('span');
  czas.className = 'dg-pozycja__czas';
  czas.textContent = new Date(wpis.timestamp).toLocaleString('pl-PL');

  const tresc = document.createElement('span');
  tresc.className = 'dg-wpis__tresc';
  const zrodlo = wpis.source === undefined ? '' : `[${wpis.source}] `;
  tresc.textContent = `${zrodlo}${wpis.message}`;

  element.append(czas, tresc);

  if (wpis.repeatCount !== undefined && wpis.repeatCount > 1) {
    const licznik = document.createElement('span');
    licznik.className = 'dg-wpis__licznik';
    licznik.textContent = `×${wpis.repeatCount}`;
    element.append(licznik);
  }

  const akcje = document.createElement('span');
  akcje.className = 'dg-pozycja__akcje';
  const pinPrzycisk = przycisk(przypiety ? 'Odepnij' : 'Przypnij', 'dn-btn dn-btn--zarys');
  pinPrzycisk.addEventListener('click', () => przelacz(wpis.id));
  akcje.append(pinPrzycisk);
  element.append(akcje);

  return element;
}

/** Raport dziennika w formacie Markdown — treść pliku pobieranego przy wywołaniu eksportu zakresu wpisów. */
function raportWpisow(wpisy: readonly LogEntry[]): string {
  const wiersze = ['# Raport dziennika Diagnostics', ''];
  for (const wpis of wpisy) {
    const zrodlo = wpis.source === undefined ? 'nieznane' : wpis.source;
    const licznik = wpis.repeatCount === undefined ? '' : ` (×${wpis.repeatCount})`;
    wiersze.push(
      `- **${new Date(wpis.timestamp).toLocaleString('pl-PL')}** [${wpis.level}] \`${zrodlo}\`${licznik} — ${wpis.message}`,
    );
  }
  return wiersze.join('\n');
}
