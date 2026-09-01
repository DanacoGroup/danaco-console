/**
 * Wiązanie panelu narzędzi Studia z rdzeniem. Znacznik niesie biblioteka
 * Właściciela: ten plik nic nie buduje — powiela wzory wykazem dołożeń sesji.
 */

import { Command, EventType, type SessionTool } from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';

/** Wzory zdjęte z treści przykładowej: nagłówek grupy i wiersz pozycji. Klon zachowuje układ i klasy nadane przez bibliotekę. */
interface WzoryPanelu {
  grupa: HTMLElement | null;
  wiersz: HTMLElement | null;
}

/** Wzory zdejmuje pierwsze wiązanie: kolejne zastaje listę wypełnioną treścią rdzenia, z której wzoru odtworzyć się już nie da. */
let wzory: WzoryPanelu | null = null;

/** Odłączenia nasłuchów poprzedniego wiązania; bez nich wiązanie założone ponownie nanosiłoby jedno zdarzenie wielokrotnie. */
let odlaczenia: Array<() => void> = [];

/**
 * Zdejmuje treść przykładową panelu narzędzi, zabierając z niej wzory nagłówka
 * grupy i wiersza pozycji. Woła się przy montażu okna, przed powstaniem
 * stanowiska: wykaz narzędzi z prototypu opisuje cudzą sesję. Zwraca prawdę,
 * gdy panel stał w dokumencie.
 */
export function zdejmijTrescPrzykladowaNarzedzi(): boolean {
  return przygotujPanel() !== null;
}

/** Wskazuje listę panelu i opróżnia ją z treści przykładowej; pustka znaczy panel poza dokumentem. */
function przygotujPanel(): HTMLElement | null {
  const znaleziona = document.querySelector('#panel-tools .sta-okno-tresc.st-panel-lista');
  if (!(znaleziona instanceof HTMLElement)) return null;
  wzory ??= zdejmijWzory(znaleziona);
  zdejmijTrescPrzykladowa(znaleziona);
  return znaleziona;
}

/** Wiąże panel narzędzi z rdzeniem; panel pokazuje narzędzia dołożone do sesji okna. Prawda znaczy, że znacznik panelu stał i wiązanie stanęło. */
export function zwiazNarzedzia(kanal: Kanal, idOkna: string): boolean {
  const znaleziona = przygotujPanel();
  if (znaleziona === null) return false;
  const lista: HTMLElement = znaleziona;

  for (const odlacz of odlaczenia) odlacz();
  odlaczenia = [];

  let idSesji = '';
  let dolozenia: SessionTool[] = [];

  async function wczytaj(): Promise<void> {
    idSesji = await wskazSesjeOkna(kanal, idOkna);
    if (idSesji === '') return;
    dolozenia = await odczytajDolozenia(kanal, idSesji);
    wypelnij(lista, dolozenia);
  }

  /* Kliknięć nie wiążemy: panel nie ma węzła dołożenia ani zdjęcia — pola
     wyszukiwania wykazu tu nie ma, a jedyny węzeł wiersza poza nazwą niósł
     znak rozwinięcia bez pokrycia w kontrakcie. */
  odlaczenia.push(
    kanal.naZdarzenie(EventType.SessionToolAttached, (tresc) => {
      if (tresc.sessionId !== idSesji) return;
      if (dolozenia.some((pozycja) => pozycja.name === tresc.tool.name)) return;
      dolozenia = [...dolozenia, tresc.tool];
      wypelnij(lista, dolozenia);
    }),
    kanal.naZdarzenie(EventType.SessionToolDetached, (tresc) => {
      if (tresc.sessionId !== idSesji) return;
      dolozenia = dolozenia.filter((pozycja) => pozycja.name !== tresc.tool.name);
      wypelnij(lista, dolozenia);
    }),
  );

  void wczytaj();
  return true;
}

/** Zdejmuje wzory z treści przykładowej. Nagłówkiem grupy jest drugi napis listy — pierwszy niesie miarę zaznaczenia, której kontrakt nie oddaje. */
function zdejmijWzory(lista: HTMLElement): WzoryPanelu {
  const napisy = lista.querySelectorAll('.pt-etykieta');
  const wiersz = sklonuj(lista.querySelector('.st-panel-wiersz'));
  // Znak rozwinięcia schodzi ze wzoru przed powielaniem: dołożenie sesji nie ma stanu rozwinięcia, którym rdzeń wypełniłby ten znak.
  wiersz?.querySelector('.dn-meta')?.remove();
  return { grupa: sklonuj(napisy[1] ?? null), wiersz };
}

/** Zdejmuje treść przykładową listy wraz z przełącznikiem zakresu, miarą zaznaczenia i wywołaniem operacji — rodzina session.tool.* nic z tego nie niesie. */
function zdejmijTrescPrzykladowa(lista: HTMLElement): void {
  const bezPokrycia = '.dn-zakladki, .pt-etykieta, .st-panel-wiersz, .st-odsun-sekcja';
  for (const wezel of lista.querySelectorAll(bezPokrycia)) wezel.remove();
}

/** Nanosi dołożenia na listę: nagłówek grupy, pod nim jej pozycje, w kolejności, w jakiej rdzeń oddał wykaz. */
function wypelnij(lista: HTMLElement, dolozenia: SessionTool[]): void {
  lista.replaceChildren();
  let grupa: string | null = null;
  for (const narzedzie of dolozenia) {
    if (narzedzie.group !== grupa) {
      grupa = narzedzie.group;
      wstawNaglowek(lista, grupa);
    }
    wstawWiersz(lista, narzedzie);
  }
}

/** Wstawia nagłówek grupy powielony ze wzoru; grupa bez nazwy nie ma czego pokazać, więc nagłówek nie staje. */
function wstawNaglowek(lista: HTMLElement, grupa: string): void {
  const naglowek = sklonuj(wzory?.grupa ?? null);
  if (naglowek === null || grupa === '') return;
  naglowek.textContent = grupa;
  lista.appendChild(naglowek);
}

/** Wstawia wiersz pozycji powielony ze wzoru. Wiersz niesie nazwę skróconą, bo to ona jest tym, co Operator wpisuje po ukośniku. */
function wstawWiersz(lista: HTMLElement, narzedzie: SessionTool): void {
  const wiersz = sklonuj(wzory?.wiersz ?? null);
  if (wiersz === null) return;
  wiersz.textContent = narzedzie.shortName;
  lista.appendChild(wiersz);
}

/** Odczytuje sesję okna z rejestru okien; wykaz dołożeń idzie po sesji, a wiązanie dostaje identyfikator okna. */
async function wskazSesjeOkna(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  return wynik.wynik.window.sessionId;
}

/** Odczytuje narzędzia dołożone do sesji; odmowa rdzenia zostawia listę pustą. */
async function odczytajDolozenia(kanal: Kanal, idSesji: string): Promise<SessionTool[]> {
  const wynik = await wywolaj(kanal, Command.SessionToolList, { sessionId: idSesji });
  if (!wynik.udany || wynik.wynik === undefined) return [];
  return wynik.wynik.tools;
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}
