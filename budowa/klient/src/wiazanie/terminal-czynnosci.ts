// Czynności okna Terminal: karty powłoki, wyjście procesu i strumień zbiorczy,
// książka gospodarzy, klucze SSH, biblioteka skryptów, odczyt pliku i tunel.
import {
  Command,
  EventType,
  TerminalKeyType,
  TerminalScriptKind,
  TerminalShell,
  TerminalTunnelKind,
  type TerminalOutputLine,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';

const NAGLOWEK = 'Terminal';

/*
zwiazCzynnosciTerminala stawia pas nad wykazem kart powłoki.

Pas stoi nad ciałem panelu, bo ciało jest wymieniane przy każdym odświeżeniu
wykazu; pas postawiony w nim znikałby razem z kartami.
*/
export function zwiazCzynnosciTerminala(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
  dolacz: (odsubskrybuj: Odsubskrybuj) => void,
  przy: AddEventListenerOptions,
): void {
  postawPas(korzen);
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const czynnosc = cel.closest<HTMLElement>('[data-terminal]')?.dataset.terminal;
    if (czynnosc === undefined) return;
    zdarzenie.stopPropagation();
    void wykonaj(kanal, korzen, czynnosc, idOkna, odswiez, dolacz);
  }, przy);
}

const PRZYCISKI: ReadonlyArray<readonly [string, string]> = [
  ['otworz', 'Otwórz kartę'],
  ['wyjscie', 'Odczytaj wyjście'],
  ['strumien', 'Śledź wyjście'],
  ['plik', 'Odczytaj plik'],
  ['obserwuj', 'Obserwuj pliki'],
  ['gospodarz', 'Zapisz gospodarza'],
  ['gospodarz-usun', 'Usuń gospodarza'],
  ['klucz', 'Wytwórz klucz'],
  ['klucz-wniesc', 'Wnieś klucz'],
  ['klucz-usun', 'Usuń klucz'],
  ['skrypt', 'Zapisz skrypt'],
  ['skrypt-sprawdz', 'Sprawdź skrypt'],
  ['skrypt-usun', 'Usuń skrypt'],
  ['tunel', 'Otwórz tunel'],
];

function postawPas(korzen: Element): void {
  const panel = korzen.querySelector('#panel-tabs');
  const cialo = panel?.querySelector('.sta-okno-tresc');
  if (panel === null || panel === undefined || cialo === null || cialo === undefined) return;
  if (panel.querySelector('[data-terminal-wpis]') !== null) return;
  const pas = korzen.ownerDocument.createElement('div');
  pas.className = 'dn-pas-dzialan';
  const pole = korzen.ownerDocument.createElement('input');
  pole.type = 'text';
  pole.className = 'dn-pole dn-pole--sm';
  pole.placeholder = 'Nazwa, ścieżka albo wzorzec';
  pole.setAttribute('aria-label', 'Nazwa, ścieżka albo wzorzec dla czynności terminala');
  pole.dataset.terminalWpis = '';
  pas.appendChild(pole);
  for (const [czynnosc, etykieta] of PRZYCISKI) pas.appendChild(przycisk(korzen, czynnosc, etykieta));
  panel.insertBefore(pas, cialo);
}

function przycisk(korzen: Element, czynnosc: string, etykieta: string): HTMLButtonElement {
  const wezel = korzen.ownerDocument.createElement('button');
  wezel.type = 'button';
  wezel.className = 'dn-btn dn-btn--duch dn-btn--sm';
  wezel.dataset.terminal = czynnosc;
  wezel.textContent = etykieta;
  return wezel;
}

function wpis(korzen: Element): string {
  const pole = korzen.querySelector<HTMLInputElement>('[data-terminal-wpis]');
  return (pole?.value ?? '').trim();
}

/* Drugie naciśnięcie potwierdza czynność nieodwracalną; napis wraca sam, więc
   uzbrojenie nie zostaje na przycisku na stałe. */
function potwierdzone(korzen: Element, czynnosc: string): boolean {
  const guzik = korzen.querySelector<HTMLElement>(`[data-terminal="${czynnosc}"]`);
  if (guzik === null) return true;
  if (guzik.dataset.uzbrojone === 'tak') {
    delete guzik.dataset.uzbrojone;
    return true;
  }
  const napis = guzik.textContent ?? '';
  guzik.dataset.uzbrojone = 'tak';
  guzik.textContent = 'Naciśnij ponownie';
  globalThis.setTimeout(() => {
    delete guzik.dataset.uzbrojone;
    guzik.textContent = napis;
  }, 5000);
  return false;
}

async function pierwszaKarta(kanal: Kanal, idOkna: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.TerminalSessionList, { windowId: idOkna });
  return wynik.wynik?.sessions[0]?.id ?? '';
}

async function wykonaj(
  kanal: Kanal,
  korzen: Element,
  czynnosc: string,
  idOkna: string,
  odswiez: () => void,
  dolacz: (odsubskrybuj: Odsubskrybuj) => void,
): Promise<void> {
  if (idOkna === '') {
    oglos(NAGLOWEK, 'Rdzeń nie dał okna terminala dla tej karty.', 'ostrzezenie');
    return;
  }
  if (czynnosc === 'otworz') return otworzKarte(kanal, korzen, idOkna, odswiez);
  if (czynnosc === 'wyjscie') return odczytajWyjscie(kanal, korzen, idOkna);
  if (czynnosc === 'strumien') return sledzWyjscie(kanal, korzen, idOkna, dolacz);
  if (czynnosc === 'plik') return odczytajPlik(kanal, korzen, idOkna);
  if (czynnosc === 'obserwuj') return obserwujPliki(kanal, korzen, idOkna);
  if (czynnosc === 'gospodarz') return zapiszGospodarza(kanal, korzen, odswiez);
  if (czynnosc === 'gospodarz-usun') return usunGospodarza(kanal, korzen, odswiez);
  if (czynnosc === 'klucz') return wytworzKlucz(kanal, korzen);
  if (czynnosc === 'klucz-wniesc') return wniesKlucz(kanal, korzen);
  if (czynnosc === 'klucz-usun') return usunKlucz(kanal, korzen);
  if (czynnosc === 'skrypt') return zapiszSkrypt(kanal, korzen, odswiez);
  if (czynnosc === 'skrypt-sprawdz') return sprawdzSkrypt(kanal, korzen);
  if (czynnosc === 'skrypt-usun') return usunSkrypt(kanal, korzen, odswiez);
  if (czynnosc === 'tunel') return otworzTunel(kanal, korzen, idOkna);
  /* Czynność spoza obsłużonych odmawia zamiast milczeć: przycisk bez gałęzi
     wyglądałby jak działający. */
  oglos(NAGLOWEK, `Czynność „${czynnosc}" nie jest prowadzona przez to okno.`, 'ostrzezenie');
}

async function otworzKarte(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  odswiez: () => void,
): Promise<void> {
  const nazwa = wpis(korzen);
  const wynik = await wywolaj(kanal, Command.TerminalSessionOpen, {
    windowId: idOkna,
    shell: TerminalShell.Bash,
    ...(nazwa === '' ? {} : { title: nazwa }),
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia karty powłoki.', 'ostrzezenie');
    return;
  }
  odswiez();
}

/* Wyjście czytane jest z procesu stojącego w wykazie najwyżej: okno nie
   prowadzi wskazania procesu, a kolejność wykazu jest jedynym wskazaniem,
   które Operator widzi. */
async function odczytajWyjscie(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const procesy = await wywolaj(kanal, Command.TerminalProcessList, { windowId: idOkna });
  const idProcesu = procesy.wynik?.processes[0]?.id ?? '';
  if (idProcesu === '') {
    oglos(NAGLOWEK, 'Żaden proces nie stoi w wykazie tego okna.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalOutputRead, { processId: idProcesu, tail: 60 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu wyjścia.', 'ostrzezenie');
    return;
  }
  const cialo = korzen.querySelector('#panel-output .sta-okno-tresc');
  const wiersze = [wynik.wynik.stdout, wynik.wynik.stderr]
    .join('\n')
    .split('\n')
    .filter((wiersz) => wiersz.trim() !== '');
  if (cialo !== null) cialo.replaceChildren(...wiersze.map((wiersz) => wierszWyjscia(cialo, wiersz)));
  if (wiersze.length === 0) {
    oglos(NAGLOWEK, 'Proces nie wypisał dotąd niczego.');
    return;
  }
  if (wynik.wynik.truncated) {
    oglos(NAGLOWEK, 'Wyjście przycięte granicą bufora — proces wypisał więcej.', 'ostrzezenie');
  }
}

function wierszWyjscia(cialo: Element, tekst: string): HTMLElement {
  const wiersz = cialo.ownerDocument.createElement('div');
  wiersz.className = 'oc-wiersz';
  wiersz.textContent = tekst;
  return wiersz;
}

/*
sledzWyjscie zapisuje okno na zbiorcze wyjście wszystkich kart powłoki.

Zapis jest subskrypcją, nie jednorazowym odczytem: wiersze przychodzą potem
zdarzeniem `stream.chunk`. Uchwyt idzie do zwolnień karty, inaczej zostałby
po jej zamknięciu i dopisywał do zdjętego już drzewa.
*/
async function sledzWyjscie(
  kanal: Kanal,
  korzen: Element,
  idOkna: string,
  dolacz: (odsubskrybuj: Odsubskrybuj) => void,
): Promise<void> {
  const panel = korzen.querySelector('#panel-output');
  const cialo = panel?.querySelector('.sta-okno-tresc') ?? null;
  if (panel === null || panel === undefined || cialo === null) return;
  if (panel.getAttribute('data-strumien') === 'tak') {
    oglos(NAGLOWEK, 'Okno już śledzi zbiorcze wyjście kart.');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalOutputStream, { windowId: idOkna, tail: 60 });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu na wyjście.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.subscribed) {
    oglos(NAGLOWEK, 'Rdzeń oddał ogon historii, ale obserwacji nie założył.', 'ostrzezenie');
    return;
  }
  panel.setAttribute('data-strumien', 'tak');
  cialo.replaceChildren(...wynik.wynik.lines.map((wiersz: TerminalOutputLine) =>
    wierszWyjscia(cialo, wiersz.text)));
  dolacz(zglosUchwyt(EventType.StreamChunk, (tresc) => {
    if (tresc.windowId !== idOkna) return;
    const tekst = (tresc.text ?? '').replace(/\n$/, '');
    if (tekst === '') return;
    cialo.appendChild(wierszWyjscia(cialo, tekst));
  }));
  oglos(NAGLOWEK, 'Okno śledzi zbiorcze wyjście kart powłoki.');
}

/* Treść pliku staje w panelu plików: panel stał dotąd pusty ze zdaniem o braku
   pokrycia, choć komenda odczytu czekała bez drogi. */
async function odczytajPlik(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const sciezka = wpis(korzen);
  if (sciezka === '') {
    oglos(NAGLOWEK, 'Odczyt pliku potrzebuje ścieżki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const idKarty = await pierwszaKarta(kanal, idOkna);
  if (idKarty === '') {
    oglos(NAGLOWEK, 'Okno nie ma karty powłoki, na której plik dałoby się odczytać.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalFileRead, {
    sessionId: idKarty,
    path: sciezka,
    tail: 200,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił odczytu pliku.', 'ostrzezenie');
    return;
  }
  const cialo = korzen.querySelector('#panel-pliki .sta-okno-tresc');
  if (cialo !== null) {
    cialo.replaceChildren(...wynik.wynik.content.split('\n')
      .map((wiersz) => wierszWyjscia(cialo, wiersz)));
  }
  oglos(NAGLOWEK, wynik.wynik.truncated
    ? `Plik ${wynik.wynik.path} przycięty granicą odczytu.`
    : `Plik ${wynik.wynik.path} odczytany w całości.`);
}

/*
obserwujPliki zakłada obserwację ścieżek karty powłoki.

Wzorzec i polecenie stoją w jednym polu rozdzielone znakiem pionowym, bo okno
nie ma dwóch pól, a komenda bez polecenia po zmianie jest bezczynna.
*/
async function obserwujPliki(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const [wzorzec = '', polecenie = ''] = wpis(korzen).split('|').map((czesc) => czesc.trim());
  if (wzorzec === '' || polecenie === '') {
    oglos(NAGLOWEK,
      'Obserwacja potrzebuje wzorca i polecenia po zmianie, na przykład „*.go | go build ./...".',
      'ostrzezenie');
    return;
  }
  const idKarty = await pierwszaKarta(kanal, idOkna);
  if (idKarty === '') {
    oglos(NAGLOWEK, 'Okno nie ma karty powłoki, w której obserwacja mogłaby stanąć.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalWatchStart, {
    sessionId: idKarty,
    pattern: wzorzec,
    command: polecenie,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił założenia obserwacji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Obserwacja stoi na wzorcu „${wzorzec}".`);
}

/* Wpis książki zakłada się nazwą i adresem celu rozdzielonymi spacją; pusty
   identyfikator znaczy dla rdzenia wpis nowy. */
async function zapiszGospodarza(kanal: Kanal, korzen: Element, odswiez: () => void): Promise<void> {
  const czesci = wpis(korzen).split(/\s+/);
  const nazwa = czesci[0] ?? '';
  const cel = czesci[1] ?? '';
  if (nazwa === '' || cel === '') {
    oglos(NAGLOWEK,
      'Wpis gospodarza potrzebuje nazwy i adresu celu, na przykład „serwer ubuntu@10.0.0.1".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalHostSave, {
    host: { id: '', name: nazwa, target: cel, createdAt: 0, updatedAt: 0 },
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu gospodarza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Gospodarz „${nazwa}" stoi w książce.`);
  odswiez();
}

async function usunGospodarza(kanal: Kanal, korzen: Element, odswiez: () => void): Promise<void> {
  const nazwa = wpis(korzen);
  const wykaz = await wywolaj(kanal, Command.TerminalHostList, {});
  const cel = wykaz.wynik?.hosts.find((gospodarz) => gospodarz.name === nazwa)?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'Książka gospodarzy nie ma wpisu o tej nazwie.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, 'gospodarz-usun')) return;
  const wynik = await wywolaj(kanal, Command.TerminalHostRemove, { hostId: cel });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia gospodarza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Gospodarz „${nazwa}" zdjęty z książki.`);
  odswiez();
}

/* Materiał klucza prywatnego nigdy nie idzie do ogłoszenia ani do schowka:
   okno nazywa oznaczenie i odcisk, bo to wystarcza do rozpoznania klucza. */
async function wytworzKlucz(kanal: Kanal, korzen: Element): Promise<void> {
  const nazwa = wpis(korzen);
  if (nazwa === '') {
    oglos(NAGLOWEK, 'Klucz potrzebuje nazwy wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalKeyGenerate, {
    name: nazwa,
    keyType: TerminalKeyType.Ed25519,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wytworzenia klucza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Klucz „${nazwa}" wytworzony; odcisk ${wynik.wynik.key.fingerprint}.`);
}

/* Wniesienie bierze klucz leżący już na maszynie rdzenia, więc pole podaje
   ścieżkę, a nie treść klucza — treść nie przechodzi tędy nigdy. */
async function wniesKlucz(kanal: Kanal, korzen: Element): Promise<void> {
  const czesci = wpis(korzen).split(/\s+/);
  const nazwa = czesci[0] ?? '';
  const sciezka = czesci[1] ?? '';
  if (nazwa === '' || sciezka === '') {
    oglos(NAGLOWEK,
      'Wniesienie klucza potrzebuje nazwy i ścieżki klucza na maszynie rdzenia.',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalKeyImport, { name: nazwa, path: sciezka });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił wniesienia klucza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Klucz „${nazwa}" wciągnięty do wykazu.`);
}

/* Pliki klucza zostają na dysku maszyny rdzenia: okno zdejmuje wyłącznie wpis
   wykazu, bo skasowania materiału nie da się cofnąć. */
async function usunKlucz(kanal: Kanal, korzen: Element): Promise<void> {
  const nazwa = wpis(korzen);
  const wykaz = await wywolaj(kanal, Command.TerminalKeyList, {});
  const cel = wykaz.wynik?.keys.find((klucz) => klucz.name === nazwa)?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'Wykaz kluczy nie ma pozycji o tej nazwie.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, 'klucz-usun')) return;
  const wynik = await wywolaj(kanal, Command.TerminalKeyRemove, { keyId: cel });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia klucza.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Klucz „${nazwa}" zdjęty z wykazu; pliki zostały na dysku rdzenia.`);
}

/* Nazwa i treść skryptu stoją w jednym polu rozdzielone znakiem pionowym;
   wersję pozycji nadaje rdzeń, więc okno wysyła zero. */
async function zapiszSkrypt(kanal: Kanal, korzen: Element, odswiez: () => void): Promise<void> {
  const [nazwa = '', zawartosc = ''] = wpis(korzen).split('|').map((czesc) => czesc.trim());
  if (nazwa === '' || zawartosc === '') {
    oglos(NAGLOWEK,
      'Pozycja biblioteki potrzebuje nazwy i treści, na przykład „budowa | go build ./...".',
      'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalScriptSave, {
    script: {
      id: '',
      name: nazwa,
      kind: TerminalScriptKind.Script,
      shell: TerminalShell.Bash,
      content: zawartosc,
      version: 0,
      createdAt: 0,
      updatedAt: 0,
    },
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił zapisu pozycji biblioteki.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja „${nazwa}" stoi w bibliotece.`);
  odswiez();
}

async function sprawdzSkrypt(kanal: Kanal, korzen: Element): Promise<void> {
  const zawartosc = wpis(korzen);
  if (zawartosc === '') {
    oglos(NAGLOWEK, 'Sprawdzenie potrzebuje treści skryptu wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalScriptLint, {
    content: zawartosc,
    shell: TerminalShell.Bash,
  });
  if (!wynik.udany || wynik.wynik === undefined) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił sprawdzenia skryptu.', 'ostrzezenie');
    return;
  }
  if (!wynik.wynik.analyzerAvailable) {
    oglos(NAGLOWEK,
      `Narzędzie ${wynik.wynik.analyzer} nie stoi na maszynie rdzenia — treści nikt nie sprawdził.`,
      'ostrzezenie');
    return;
  }
  const uwagi = wynik.wynik.findings.length;
  oglos(NAGLOWEK, uwagi === 0
    ? `Narzędzie ${wynik.wynik.analyzer} nie ma uwag do treści.`
    : `Narzędzie ${wynik.wynik.analyzer} zgłasza ${uwagi} uwag.`);
}

async function usunSkrypt(kanal: Kanal, korzen: Element, odswiez: () => void): Promise<void> {
  const nazwa = wpis(korzen);
  const wykaz = await wywolaj(kanal, Command.TerminalScriptList, {});
  const cel = wykaz.wynik?.scripts.find((pozycja) => pozycja.name === nazwa)?.id ?? '';
  if (cel === '') {
    oglos(NAGLOWEK, 'Biblioteka nie ma pozycji o tej nazwie.', 'ostrzezenie');
    return;
  }
  if (!potwierdzone(korzen, 'skrypt-usun')) return;
  const wynik = await wywolaj(kanal, Command.TerminalScriptRemove, { scriptId: cel });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił usunięcia pozycji.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Pozycja „${nazwa}" zdjęta z biblioteki wraz z wersjami.`);
  odswiez();
}

/* Tunel idzie przez wpis książki gospodarzy o nazwie z pola: bez gospodarza
   rdzeń nie ma przez co go poprowadzić. */
async function otworzTunel(kanal: Kanal, korzen: Element, idOkna: string): Promise<void> {
  const nazwa = wpis(korzen);
  const wykaz = await wywolaj(kanal, Command.TerminalHostList, {});
  const gospodarz = wykaz.wynik?.hosts.find((wpisKsiazki) => wpisKsiazki.name === nazwa)?.id ?? '';
  if (gospodarz === '') {
    oglos(NAGLOWEK, 'Tunel potrzebuje nazwy gospodarza z książki wpisanej w polu.', 'ostrzezenie');
    return;
  }
  const wynik = await wywolaj(kanal, Command.TerminalTunnelOpen, {
    windowId: idOkna,
    kind: TerminalTunnelKind.Dynamic,
    hostId: gospodarz,
  });
  if (!wynik.udany) {
    oglos(NAGLOWEK, wynik.blad?.message ?? 'Rdzeń odmówił otwarcia tunelu.', 'ostrzezenie');
    return;
  }
  oglos(NAGLOWEK, `Tunel przez gospodarza „${nazwa}" otwarty.`);
}
