/**
 * Wiązanie okna Studia z rdzeniem. Znacznik niesie biblioteka Właściciela —
 * ten plik nic nie buduje: słucha zdarzeń, woła komendy kontraktu i wypełnia
 * stojące węzły, a wpisy i pozycje powiela z wzorów zdjętych ze znacznika.
 */

import {
  ChangeKind,
  Command,
  EventType,
  ExecutionEnv,
  MessageRole,
  PermissionMode,
  WindowRole,
  type Message,
  type Session,
  type StudioDocument,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { zwiazNarzedzia } from './studio-narzedzia.ts';
import { zwiazPlan } from './studio-plan.ts';
import { zwiazPliki } from './studio-pliki.ts';
import { zwiazPodglad } from './studio-podglad.ts';
import { zwiazRepozytorium } from './studio-repozytorium.ts';
import { zwiazRoznice } from './studio-roznice.ts';

/** Kod modułu Studia w rejestrze rdzenia; `module.list` oddaje po nim identyfikator, którego wymaga `window.create`. */
const KOD_MODULU_STUDIO = 'studio';

/** Węzły znacznika Właściciela, na których wiązanie pracuje. Brak któregokolwiek znaczy, że okno Studia nie stoi i wiązać nie ma czego. */
interface WezlyStudia {
  historia: HTMLElement;
  formularz: HTMLFormElement;
  pole: HTMLTextAreaElement;
  szyna: HTMLElement;
  nowyDokument: HTMLElement | null;
  kanwa: HTMLElement;
  status: HTMLElement;
}

/** Wzory wpisu historii zdjęte z treści przykładowej, po jednym na rodzaj nadawcy — kształt wpisu rdzenia bierze się stąd, nie z kodu. */
interface WzoryWpisow {
  czlowiek: HTMLElement | null;
  inteligencja: HTMLElement | null;
  system: HTMLElement | null;
}

/** Wiązanie stoi raz na dokument: `studio.js` i powłoka mogą wstawić okno ponownie, a podwójny nasłuch dawałby podwójne wpisy. */
let zwiazane = false;

/**
 * Wiąże okno Studia z rdzeniem. Kanał domyślnie z obiektu globalnego, bo
 * skrypty biblioteki są funkcjami domkniętymi i importu z nich nie ma.
 * Zwraca prawdę, gdy znacznik okna stał i wiązanie zostało założone.
 */
export function zwiazStudio(
  kanal: Kanal | undefined = globalThis.DanacoKanal,
  nazwaSrodowiska = '',
): boolean {
  if (zwiazane || kanal === undefined) return false;
  const znalezione = zbierzWezly();
  if (znalezione === null) return false;
  zwiazane = true;
  const wezly: WezlyStudia = znalezione;

  const wzory = zdejmijWzoryWpisow(wezly.historia);
  const wzorPozycji = zdejmijWzorPozycji(wezly.szyna);

  // Treść przykładowa znika, zanim padnie pierwsza odpowiedź rdzenia: pusty
  // węzeł jest uczciwy, znacznik z cudzym dokumentem — nie.
  wezly.historia.replaceChildren();
  wezly.szyna.replaceChildren();
  wezly.kanwa.replaceChildren();

  const wpisy = new Map<string, HTMLElement>();
  let idOkna = '';
  let dokument: StudioDocument | null = null;

  /** Wstawia wpis rdzenia na koniec historii, powielając wzór o kształcie zgodnym z rolą nadawcy. */
  function wstawWpis(wiadomosc: Message): void {
    const wzor = wzorDlaRoli(wzory, wiadomosc.role);
    if (wzor === null) return;
    const wpis = wzor.cloneNode(true) as HTMLElement;
    wypelnijWpis(wpis, wiadomosc);
    wezly.historia.appendChild(wpis);
    wpisy.set(wiadomosc.id, wpis);
    wezly.historia.scrollTop = wezly.historia.scrollHeight;
  }

  /** Nanosi zmianę wiadomości na wpis już stojący; wpis nieznany dopisuje, usunięty zdejmuje. */
  function zmienWpis(zmiana: ChangeKind, wiadomosc: Message): void {
    const stojacy = wpisy.get(wiadomosc.id);
    if (zmiana === ChangeKind.Deleted) {
      stojacy?.remove();
      wpisy.delete(wiadomosc.id);
      return;
    }
    if (stojacy === undefined) {
      wstawWpis(wiadomosc);
      return;
    }
    wypelnijWpis(stojacy, wiadomosc);
  }

  /* Enter wysyła, Shift+Enter przechodzi do nowego wiersza. Pole jest obszarem
     tekstowym, a ten sam z siebie formularza nie zamyka — bez tego jedyną drogą
     wysłania byłby przycisk. */
  wezly.pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    wezly.formularz.requestSubmit();
  });

  wezly.formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    const tresc = wezly.pole.value.trim();
    if (tresc === '' || idOkna === '') return;
    wezly.pole.value = '';
    void wywolaj(kanal, Command.MessageSend, { windowId: idOkna, content: tresc, stream: true });
  });

  // Ctrl+S zapisuje treść kanwy: pas stanu prototypu mówi o zapisie
  // automatycznym, ale odłożenie wersji w repozytorium sesji ma mieć wyzwalacz.
  wezly.kanwa.addEventListener('keydown', (zdarzenie) => {
    if (!zdarzenie.ctrlKey || zdarzenie.key.toLowerCase() !== 's') return;
    zdarzenie.preventDefault();
    if (dokument === null) return;
    void zapiszDokument(kanal, dokument.id, wezly, (zapisany) => {
      dokument = zapisany;
      opiszDokument(zapisany);
    });
  });

  wezly.kanwa.addEventListener('input', () => {
    odswiezPasStanu(wezly, dokument);
  });

  wezly.nowyDokument?.addEventListener('click', () => {
    if (idOkna === '') return;
    void zalozDokument(kanal, idOkna, wezly).then((zalozony) => {
      dokument = zalozony;
      opiszDokument(zalozony);
    });
  });

  kanal.naZdarzenie(EventType.MessageChanged, (tresc) => {
    if (tresc.message.windowId !== idOkna) return;
    zmienWpis(tresc.change, tresc.message);
  });

  zdejmijZnacznikiBezZrodla();
  void opiszKanal(kanal, nazwaSrodowiska);
  void otworzStanowisko(kanal, wezly, wzorPozycji, wstawWpis).then((okno) => {
    idOkna = okno;
    if (okno === '') return;
    zwiazPanele(kanal, okno);
    void zalozDokument(kanal, okno, wezly).then((zalozony) => {
      dokument = zalozony;
      opiszDokument(zalozony);
    });
  });

  return true;
}

/** Zdejmuje znaczniki kontekstu bez pokrycia w kontrakcie: nazwę gałęzi, ścieżkę repozytorium i miarę różnicy, których rdzeń nie oddaje. */
function zdejmijZnacznikiBezZrodla(): void {
  for (const znacznik of document.querySelectorAll('.sta-kontekst-akcji .sta-chip')) {
    if (znacznik.classList.contains('sta-chip--srodowisko')) continue;
    znacznik.remove();
  }
}

/** Wpisuje w nagłówek okna komunikacji środowisko wejścia i model kanału; pole wysiłku znika, bo kontrakt nie niesie jego wartości. */
async function opiszKanal(kanal: Kanal, nazwaSrodowiska: string): Promise<void> {
  const naglowek = document.querySelector('.sta-kom-naglowek');
  if (naglowek === null) return;
  const pola = [...naglowek.querySelectorAll('.sta-kom-pole')];
  wpiszPole(pola, 'Środowisko', nazwaSrodowiska);
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = wynik.udany ? wynik.wynik?.channels[0] : undefined;
  wpiszPole(pola, 'Model', kanalModelu?.model ?? kanalModelu?.name ?? '');
  wpiszPole(pola, 'Wysiłek', '');
}

/** Wpisuje wartość w pole nagłówka rozpoznane po jego podpisie; wartość pusta zdejmuje całe pole, bo pole bez wartości niczego nie mówi. */
function wpiszPole(pola: Element[], podpis: string, wartosc: string): void {
  const pole = pola.find((kandydat) => kandydat.textContent?.startsWith(podpis) === true);
  if (pole === undefined) return;
  if (wartosc === '') {
    pole.remove();
    return;
  }
  const dane = pole.querySelector('.dane');
  if (dane !== null) dane.textContent = wartosc;
}

/** Nadaje oknu i pasowi tytuł dokumentu; dokument bez nadanej nazwy dostaje nazwany stan pusty, nie własny identyfikator. */
function opiszDokument(dokument: StudioDocument | null): void {
  const nazwa = dokument?.title ?? 'Dokument bez nazwy';
  const miejsca = '.st-wstazka-sesja span, .sta-okno-znacznik, .dn-karta--robocza .dn-karta-widoku-nazwa';
  for (const wezel of document.querySelectorAll(miejsca)) wezel.textContent = nazwa;
  wpiszWersje(dokument);
}

/** Nanosi wersję dokumentu na miarę wstążki; dokument bez wersji w repozytorium sesji nie ma czego pokazać, więc miara znika. */
function wpiszWersje(dokument: StudioDocument | null): void {
  for (const wezel of document.querySelectorAll('.st-wstazka-stan .st-miara')) {
    if (dokument?.versionId === undefined) wezel.remove();
    else wezel.textContent = `wersja ${dokument.versionId}`;
  }
}

/** Zakłada sesję i okno komunikacji, wczytuje historię i szynę sesji; zwraca identyfikator okna albo pustkę, gdy rdzeń odmówił. */
async function otworzStanowisko(
  kanal: Kanal,
  wezly: WezlyStudia,
  wzorPozycji: HTMLElement | null,
  wstawWpis: (wiadomosc: Message) => void,
): Promise<string> {
  const sesja = await wywolaj(kanal, Command.SessionCreate, {});
  if (!sesja.udany || sesja.wynik === undefined) return '';
  const idSesji = sesja.wynik.session.id;

  const modul = await wskazModulStudia(kanal);
  const kanalModelu = await wskazKanalModelu(kanal);
  if (modul === '' || kanalModelu === '') return '';

  const okno = await wywolaj(kanal, Command.WindowCreate, {
    sessionId: idSesji,
    moduleId: modul,
    modelChannelId: kanalModelu,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    permissionMode: PermissionMode.Manual,
    windowRole: WindowRole.Standalone,
  });
  if (!okno.udany || okno.wynik === undefined) return '';
  const idOkna = okno.wynik.window.id;

  const historia = await wywolaj(kanal, Command.MessageList, { windowId: idOkna });
  if (historia.udany && historia.wynik !== undefined) {
    for (const wiadomosc of historia.wynik.messages) wstawWpis(wiadomosc);
  }

  await wypelnijSzyne(kanal, wezly.szyna, wzorPozycji, idSesji);
  return idOkna;
}

/** Odczytuje identyfikator modułu Studia z rejestru modułów; kod modułu jest stały, identyfikator nadaje rdzeń. */
async function wskazModulStudia(kanal: Kanal): Promise<string> {
  const wynik = await wywolaj(kanal, Command.ModuleList, {});
  if (!wynik.udany || wynik.wynik === undefined) return '';
  const modul = wynik.wynik.modules.find((pozycja) => pozycja.code === KOD_MODULU_STUDIO);
  return modul?.id ?? '';
}

/** Bierze pierwszy czynny kanał modelu z rejestru; wyboru kanału znacznik prototypu nie niesie, a `window.create` wymaga wskazania. */
async function wskazKanalModelu(kanal: Kanal): Promise<string> {
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  return wynik.wynik.channels[0]?.id ?? '';
}

/** Wypełnia szynę dokumentów pozycjami sesji rdzenia, powielając wzór pozycji zdjęty ze znacznika. */
async function wypelnijSzyne(
  kanal: Kanal,
  szyna: HTMLElement,
  wzor: HTMLElement | null,
  idBiezacej: string,
): Promise<void> {
  if (wzor === null) return;
  const wynik = await wywolaj(kanal, Command.SessionList, {});
  if (!wynik.udany || wynik.wynik === undefined) return;
  szyna.replaceChildren();
  for (const sesja of wynik.wynik.sessions) {
    szyna.appendChild(zbudujPozycje(wzor, sesja, sesja.id === idBiezacej));
  }
  szyna.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const pozycja = cel.closest('.pt-pozycja');
    if (pozycja === null) return;
    for (const inna of szyna.querySelectorAll('.pt-pozycja')) inna.removeAttribute('aria-current');
    pozycja.setAttribute('aria-current', 'true');
  });
}

/** Zwraca klon wzoru pozycji szyny opisany nazwą sesji; tętno zostaje wyłącznie przy pozycji bieżącej, bo tak niesie je znacznik. */
function zbudujPozycje(wzor: HTMLElement, sesja: Session, biezaca: boolean): HTMLElement {
  const pozycja = wzor.cloneNode(true) as HTMLElement;
  const tytul = pozycja.querySelector('.pt-pozycja-tytul');
  /* Sesja bez nadanej nazwy dostaje nazwany stan pusty, nie własny
     identyfikator: identyfikator jest oznaczeniem magazynu, a Operator czyta
     w szynie nazwę swojej pracy. */
  if (tytul !== null) tytul.textContent = sesja.title ?? 'Sesja bez nazwy';
  if (biezaca) pozycja.setAttribute('aria-current', 'true');
  else {
    pozycja.removeAttribute('aria-current');
    pozycja.querySelector('.pt-tetno')?.remove();
  }
  return pozycja;
}

/** Zakłada dokument w oknie, wczytuje go i wstawia jego treść w kanwę wraz z pasem stanu. */
async function zalozDokument(
  kanal: Kanal,
  idOkna: string,
  wezly: WezlyStudia,
): Promise<StudioDocument | null> {
  const zalozony = await wywolaj(kanal, Command.StudioDocumentCreate, { windowId: idOkna });
  if (!zalozony.udany || zalozony.wynik === undefined) return null;

  const otwarty = await wywolaj(kanal, Command.StudioDocumentOpen, {
    windowId: idOkna,
    documentId: zalozony.wynik.document.id,
  });
  const dokument = otwarty.udany && otwarty.wynik !== undefined
    ? otwarty.wynik.document
    : zalozony.wynik.document;

  const tekst = await wywolaj(kanal, Command.StudioTextGet, { documentId: dokument.id });
  wezly.kanwa.textContent = tekst.udany && tekst.wynik !== undefined ? tekst.wynik.text : '';
  // Kanwa prototypu stoi zamknięta na edycję, bo niosła dokument przykładowy;
  // z treścią rdzenia ma być polem pracy, więc atrybut idzie na otwarty.
  wezly.kanwa.setAttribute('contenteditable', 'true');
  odswiezPasStanu(wezly, dokument);
  return dokument;
}

/** Zapisuje treść kanwy jako nową wersję dokumentu i odświeża pas stanu wersją zwróconą przez rdzeń. */
async function zapiszDokument(
  kanal: Kanal,
  idDokumentu: string,
  wezly: WezlyStudia,
  poZapisie: (dokument: StudioDocument) => void,
): Promise<void> {
  const wynik = await wywolaj(kanal, Command.StudioDocumentSave, {
    documentId: idDokumentu,
    content: wezly.kanwa.textContent ?? '',
    createVersion: true,
  });
  if (!wynik.udany || wynik.wynik === undefined) return;
  poZapisie(wynik.wynik.document);
  odswiezPasStanu(wezly, wynik.wynik.document);
}

/** Nanosi na pas stanu liczbę słów policzoną z kanwy i wersję dokumentu; pola znajduje po ich własnych podpisach, nie po miejscu w rzędzie. */
function odswiezPasStanu(wezly: WezlyStudia, dokument: StudioDocument | null): void {
  const pola = [...wezly.status.querySelectorAll('span')];
  const slowa = pola.find((pole) => pole.textContent?.startsWith('słów') === true);
  if (slowa !== undefined) slowa.textContent = `słów: ${policzSlowa(wezly.kanwa.textContent ?? '')}`;
  const wersja = pola.find((pole) => pole.textContent?.startsWith('wersja') === true);
  if (wersja !== undefined) {
    wersja.textContent = dokument?.versionId === undefined ? '' : `wersja ${dokument.versionId}`;
  }
  /* Zapis idzie wyzwalaczem, nie zegarem — godzina zapisu samoczynnego nie ma
     w kontrakcie źródła, a godzina zmyślona mówi Operatorowi nieprawdę o tym,
     czy jego praca jest odłożona. */
  pola.find((pole) => pole.textContent?.includes('zapisano') === true)?.remove();
}

/** Liczba słów treści — ciągi znaków rozdzielone białymi znakami. */
function policzSlowa(tresc: string): number {
  const cialo = tresc.trim();
  return cialo === '' ? 0 : cialo.split(/\s+/u).length;
}

/** Wskazuje węzły okna Studia; pustka znaczy, że okno nie stoi w dokumencie. */
function zbierzWezly(): WezlyStudia | null {
  const kom = document.querySelector('.sta-kom');
  const historia = kom?.querySelector('.sta-kom-historia');
  const formularz = kom?.querySelector('form.sta-prompt');
  const pole = kom?.querySelector('textarea.sta-prompt-obszar');
  const szyna = document.querySelector('.st-szyna-lista');
  const kanwa = document.querySelector('.dn-kanwa');
  const status = document.querySelector('.st-status');
  if (
    !(historia instanceof HTMLElement) ||
    !(formularz instanceof HTMLFormElement) ||
    !(pole instanceof HTMLTextAreaElement) ||
    !(szyna instanceof HTMLElement) ||
    !(kanwa instanceof HTMLElement) ||
    !(status instanceof HTMLElement)
  ) {
    return null;
  }
  const nowyDokument = document.querySelector('.st-szyna-naglowek button.dn-btn');
  return {
    historia,
    formularz,
    pole,
    szyna,
    nowyDokument: nowyDokument instanceof HTMLElement ? nowyDokument : null,
    kanwa,
    status,
  };
}

/** Zdejmuje z historii po jednym wpisie każdego rodzaju jako wzór; klon zachowuje medalion, układ i klasy nadane przez bibliotekę. */
function zdejmijWzoryWpisow(historia: HTMLElement): WzoryWpisow {
  return {
    czlowiek: sklonuj(historia.querySelector('.sta-wpis--czlowiek')),
    inteligencja: sklonuj(historia.querySelector('.sta-wpis--inteligencja')),
    system: sklonuj(historia.querySelector('.sta-wpis--system')),
  };
}

/** Zdejmuje wzór pozycji szyny z pierwszej pozycji przykładowej. */
function zdejmijWzorPozycji(szyna: HTMLElement): HTMLElement | null {
  return sklonuj(szyna.querySelector('.pt-pozycja'));
}

/** Klon węzła wzorcowego, odporny na jego brak w znaczniku. */
function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

/** Wzór wpisu właściwy roli nadawcy; wynik narzędzia idzie kształtem systemowym, bo znacznik osobnego nie niesie. */
function wzorDlaRoli(wzory: WzoryWpisow, rola: MessageRole): HTMLElement | null {
  if (rola === MessageRole.User) return wzory.czlowiek;
  if (rola === MessageRole.Assistant) return wzory.inteligencja;
  return wzory.system;
}

/** Nazwa nadawcy w brzmieniu znacznika Właściciela — słownictwo pochodzi z prototypu, nie z tego pliku. */
function nazwaNadawcy(rola: MessageRole): string {
  if (rola === MessageRole.User) return 'Operator';
  if (rola === MessageRole.Assistant) return 'Inteligencja';
  return 'System';
}

/** Wypełnia klon wpisu treścią wiadomości rdzenia: nadawca, godzina, treść. Plakietki przykładowe znikają — rdzeń nie oddaje tego, co niosły. */
function wypelnijWpis(wpis: HTMLElement, wiadomosc: Message): void {
  const nadawca = wpis.querySelector('.sta-wpis-nadawca');
  if (nadawca !== null) nadawca.textContent = nazwaNadawcy(wiadomosc.role);
  const godzina = wpis.querySelector('.sta-wpis-godzina');
  if (godzina !== null) godzina.textContent = godzinaWpisu(wiadomosc.createdAt);
  for (const plakietka of wpis.querySelectorAll('.dn-plakietka')) plakietka.remove();
  const tresc = wpis.querySelector('.sta-wpis-tresc');
  if (tresc === null) return;
  // Pasek działań wpisu należy do biblioteki i przeżywa podmianę treści;
  // znika tylko zdanie przykładowe, nie przyciski, które znacznik niesie.
  const akcje = tresc.querySelector('.sta-wpis-akcje');
  tresc.textContent = wiadomosc.content;
  if (akcje !== null) tresc.appendChild(akcje);
}

/** Godzina wpisu w zapisie, którego używa znacznik historii — godziny i minuty. */
function godzinaWpisu(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

// Okno Studia wstawia w ramę powłoka, już po wczytaniu modułów; obserwator
// wiąże je w chwili, gdy znacznik stanie, i odłącza się po pierwszym wiązaniu.
if (!zwiazStudio()) {
  const obserwator = new MutationObserver(() => {
    if (zwiazStudio()) obserwator.disconnect();
  });
  obserwator.observe(document.documentElement, { childList: true, subtree: true });
}

/* Panele okna roboczego wiąże się dopiero po założeniu okna: każdy z nich pyta
   rdzeń o treść tego okna, a przed jego powstaniem nie ma o co pytać. Panel,
   którego znacznik nie stoi, zwraca fałsz i nie robi nic. */
function zwiazPanele(kanal: Kanal, idOkna: string): void {
  zwiazPlan(kanal, idOkna);
  zwiazRoznice(kanal, idOkna);
  zwiazRepozytorium(kanal, idOkna);
  zwiazPliki(kanal, idOkna);
  zwiazNarzedzia(kanal, idOkna);
  zwiazPodglad(kanal, idOkna);
}
