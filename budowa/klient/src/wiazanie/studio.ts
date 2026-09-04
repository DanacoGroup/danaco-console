// Wiązanie karty Studia z rdzeniem: słucha zdarzeń, woła komendy kontraktu
// i wypełnia węzły karty wzorami zdjętymi ze znacznika Właściciela.

import {
  ChangeKind,
  ChunkKind,
  Command,
  ReasoningEffort,
  EventType,
  ExecutionEnv,
  MessageRole,
  PermissionMode,
  WindowRole,
  WindowStatus,
  type Message,
  type PanelSection,
  type Session,
  type StreamChunkEvent,
  type StudioDocument,
} from '../../../shared/contract.ts';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen.ts';
import { zglosUchwyt } from '../polaczenie/rozdzielacz-zdarzen.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { wpiszPole } from './okno-modulu.ts';
import {
  kartyOkna,
  odczytajUkladPaneli,
  oknaRobocze,
  przypiszOknoKomunikacji,
  zapiszUkladPaneli,
} from './okna-robocze.ts';
import { zapewnijSesje } from './sesja-biezaca.ts';
import { zdejmijTrescPrzykladowaAutozapisu, zwiazAutozapis } from './studio-autozapis.ts';
import {
  zapiszDokumentZPostacia,
  zdejmijTrescPrzykladowaFormatu,
  zwiazFormatDokumentu,
} from './studio-dokument.ts';
import {
  zdejmijTrescPrzykladowaFormatowania,
  zwiazFormatowanie,
} from './studio-formatowanie.ts';
import { zdejmijTrescPrzykladowaStylow, zwiazStyle } from './studio-style.ts';
import { zdejmijTrescPrzykladowaStruktur, zwiazStruktury } from './studio-struktury.ts';
import {
  uruchomOperacjeDokumentu,
  zdejmijTrescPrzykladowaNarzedzi,
  zwiazNarzedzia,
} from './studio-narzedzia.ts';
import { zwiazKatalogOperacji } from './studio-katalog.ts';
import { zdejmijTrescPrzykladowaStanu, zwiazStanOkna } from './studio-stan-okna.ts';
import { zdejmijTrescPrzykladowaPlanu, zwiazPlan } from './studio-plan.ts';
import { zdejmijTrescPrzykladowaPlikow, zwiazPliki } from './studio-pliki.ts';
import { zdejmijTrescPrzykladowaPodgladu, zwiazPodglad } from './studio-podglad.ts';
import { zdejmijTrescPrzykladowaRepozytorium, zwiazRepozytorium } from './studio-repozytorium.ts';
import {
  porownajWersjeDokumentu,
  zdejmijTrescPrzykladowaRoznic,
  zwiazRoznice,
} from './studio-roznice.ts';

const KOD_MODULU_STUDIO = 'studio';

const PANEL_KART_STUDIA = 'st-karty';

interface WezlyStudia {
  historia: HTMLElement;
  formularz: HTMLFormElement;
  pole: HTMLTextAreaElement;
  zatrzymaj: HTMLElement | null;
  szyna: HTMLElement;
  nowyDokument: HTMLElement | null;
  kanwa: HTMLElement;
  status: HTMLElement;
  zapiszWersje: HTMLElement | null;
  porownajWersje: HTMLElement | null;
  uruchomOperacje: HTMLElement | null;
}

interface WzoryWpisow {
  czlowiek: HTMLElement | null;
  inteligencja: HTMLElement | null;
  system: HTMLElement | null;
}

interface FragmentStrumienia {
  tresc: StreamChunkEvent;
  ostatni: boolean;
}

interface StrumienOdpowiedzi {
  wpis: HTMLElement;
  tekst: string;
  nastepny: number;
  odlozone: Map<number, FragmentStrumienia>;
}

interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

// Każda karta ma własne okno i uchwyty; zamknięcie zwalnia wyłącznie jej wiązanie.
const WIAZANIA = new Map<string, WiazanieKarty>();

// Okno stojące podaje wołający; pustka znaczy okno z pierwszej wiadomości.
export function zwiazStudio(
  kanal: Kanal,
  nazwaSrodowiska: string,
  idOknaStojacego: string,
  wskazanieKorzenia: Element | string,
): boolean {
  const korzen = korzenKarty(wskazanieKorzenia);
  if (korzen === null) return false;
  const idKarty = korzen.getAttribute('data-karta') ?? '';
  if (idKarty === '') return false;
  if (WIAZANIA.get(idKarty)?.korzen === korzen) return false;
  const znalezione = zbierzWezly(korzen);
  if (znalezione === null) return false;
  zwolnijStudio(idKarty);
  const odlaczenia: Odsubskrybuj[] = [];
  // Nasłuchy karty schodzą razem z nią, zdjęte sterownikiem przerwania.
  const sterowanie = new AbortController();
  const przy = { signal: sterowanie.signal };
  odlaczenia.push(() => {
    sterowanie.abort();
  });
  WIAZANIA.set(idKarty, { korzen, odlaczenia });
  const wezly: WezlyStudia = znalezione;

  const wzory = zdejmijWzoryWpisow(wezly.historia);
  const wzorPozycji = zdejmijWzorPozycji(wezly.szyna);
  zdejmijTrescPrzykladowa(wezly, korzen);
  zdejmijCzynnosciWstazkiBezPokrycia(korzen);

  // Biblioteka wiąże pasmo kart raz, przy wczytaniu skryptu; karta Studia
  // wchodzi w dokument później, więc przełączanie panelu stoi tutaj.
  korzen.addEventListener('click', (zdarzenie) => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element)) return;
    const karta = cel.closest<HTMLElement>('.st-karty [role="tab"][data-karta]');
    if (karta !== null) przelaczKarte(korzen, karta);
  }, przy);

  const wpisy = new Map<string, HTMLElement>();
  const strumienie = new Map<string, StrumienOdpowiedzi>();
  let idOkna = '';
  let idSesji = '';
  let dokument: StudioDocument | null = null;

  const opiszDokument = (opisywany: StudioDocument | null): void => {
    opiszDokumentKarty(korzen, opisywany);
  };

  function wstawWpis(wiadomosc: Message): void {
    const wzor = wzorDlaRoli(wzory, wiadomosc.role);
    if (wzor === null) return;
    const wpis = wzor.cloneNode(true) as HTMLElement;
    wypelnijWpis(wpis, wiadomosc);
    wezly.historia.appendChild(wpis);
    wpisy.set(wiadomosc.id, wpis);
    wezly.historia.scrollTop = wezly.historia.scrollHeight;
  }

  function zmienWpis(zmiana: ChangeKind, wiadomosc: Message): void {
    const stojacy = wpisy.get(wiadomosc.id);
    if (zmiana === ChangeKind.Deleted) {
      stojacy?.remove();
      wpisy.delete(wiadomosc.id);
      strumienie.delete(wiadomosc.id);
      return;
    }
    if (stojacy === undefined) {
      wstawWpis(wiadomosc);
      return;
    }
    wypelnijWpis(stojacy, wiadomosc);
  }

  // Wiadomość modelu rdzeń rozgłasza po turze, więc godziny nadania jeszcze nie ma.
  function wstawWpisStrumienia(idWiadomosci: string): HTMLElement | null {
    const wzor = wzory.inteligencja;
    if (wzor === null) return null;
    const wpis = wzor.cloneNode(true) as HTMLElement;
    const nadawca = wpis.querySelector('.sta-wpis-nadawca');
    if (nadawca !== null) nadawca.textContent = nazwaNadawcy(MessageRole.Assistant);
    const godzina = wpis.querySelector('.sta-wpis-godzina');
    if (godzina !== null) godzina.textContent = '';
    for (const plakietka of wpis.querySelectorAll('.dn-plakietka')) plakietka.remove();
    wpiszTrescWpisu(wpis, '');
    wezly.historia.appendChild(wpis);
    wpisy.set(idWiadomosci, wpis);
    return wpis;
  }

  function przyjmijFragment(tresc: StreamChunkEvent, numer: number, ostatni: boolean): void {
    let strumien = strumienie.get(tresc.messageId);
    if (strumien === undefined) {
      const wpis = wpisy.get(tresc.messageId) ?? wstawWpisStrumienia(tresc.messageId);
      if (wpis === null) return;
      /* Porządek liczy się od fragmentu pierwszego, jaki to okno zobaczyło:
         liczenie od jedynki zatrzymałoby strumień zastany w połowie. */
      strumien = {
        wpis,
        tekst: trescWpisu(wpis),
        nastepny: numer === 0 ? 1 : numer,
        odlozone: new Map(),
      };
      strumienie.set(tresc.messageId, strumien);
    }
    if (numer > strumien.nastepny) {
      strumien.odlozone.set(numer, { tresc, ostatni });
      return;
    }
    if (numer !== 0 && numer < strumien.nastepny) return;
    naniesFragment(tresc.messageId, strumien, { tresc, ostatni });
  }

  function naniesFragment(
    idWiadomosci: string,
    strumien: StrumienOdpowiedzi,
    fragment: FragmentStrumienia,
  ): void {
    let biezacy: FragmentStrumienia | undefined = fragment;
    while (biezacy !== undefined) {
      strumien.tekst = zlozTresc(strumien.tekst, biezacy.tresc);
      wpiszTrescWpisu(strumien.wpis, strumien.tekst);
      wezly.historia.scrollTop = wezly.historia.scrollHeight;
      if (biezacy.ostatni) {
        strumienie.delete(idWiadomosci);
        return;
      }
      strumien.nastepny += 1;
      biezacy = strumien.odlozone.get(strumien.nastepny);
      strumien.odlozone.delete(strumien.nastepny);
    }
  }

  // Obszar tekstowy sam formularza nie zamyka: bez tego zostaje sam przycisk.
  wezly.pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    wezly.formularz.requestSubmit();
  }, przy);

  const zdejmijStanowisko = (): void => {
    idOkna = '';
  };

  /* Sesja powstaje pierwszą wiadomością, nie otwarciem okna ani wejściem
     w moduł; wyjątkiem jest okno wskazane przez wołającego. */
  const zapewnijStanowisko = async (): Promise<string> => {
    if (idOkna !== '') return idOkna;
    const stanowisko = await otworzStanowisko(
      kanal, wezly, wzorPozycji, wstawWpis, idOknaStojacego, idKarty, przy,
    );
    if (stanowisko === null) return '';
    idOkna = stanowisko.idOkna;
    idSesji = stanowisko.idSesji;
    // Karta bez okna komunikacji nie zamknie go w rdzeniu przy swoim zejściu.
    if (!przypiszOknoKomunikacji(idKarty, idOkna)) {
      oglos('Studio', 'Okno rozmowy stanęło, gdy jego karta zeszła z pasma — '
        + 'zamknięcie karty nie zamknie tego okna.', 'ostrzezenie');
    }
    zwiazUkladKart(kanal, idKarty, korzen, przy);
    zwiazPanele(kanal, idOkna, idKarty, korzen, odlaczenia);
    // Okno stojące prowadzi już dokument; nowy mnożyłby dokumenty przy wejściu.
    dokument = idOknaStojacego === ''
      ? await zalozDokument(kanal, idOkna, wezly)
      : await otworzDokumentOkna(kanal, idOkna, wezly);
    opiszDokument(dokument);
    return idOkna;
  };

  // Rdzeń odmawia wiadomości do okna w turze: pole czyści się po jej przyjęciu.
  const wyslijWiadomosc = async (tresc: string): Promise<void> => {
    const okno = await zapewnijStanowisko();
    if (okno === '') {
      oglos('Studio', 'Rdzeń nie założył stanowiska — wiadomość nie została wysłana.', 'blad');
      return;
    }
    const wynik = await wywolaj(kanal, Command.MessageSend, {
      windowId: okno,
      content: tresc,
      stream: true,
    });
    if (!wynik.udany) {
      oglos('Studio', wynik.blad?.message ?? 'Rdzeń odmówił przyjęcia wiadomości.', 'ostrzezenie');
      return;
    }
    // Pole czyści się tylko wtedy, gdy nadal niesie wysłany tekst.
    if (wezly.pole.value.trim() === tresc) wezly.pole.value = '';
  };

  wezly.formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    const tresc = wezly.pole.value.trim();
    if (tresc === '') return;
    void wyslijWiadomosc(tresc);
  }, przy);

  /* Zatrzymanie czynne niezależnie od stanu tury: rdzeń odpowiada na
     message.stop zawsze, a brak tury oddaje `stopped` równe fałszowi. */
  wezly.zatrzymaj?.addEventListener('click', () => {
    if (idOkna === '') {
      oglos('Studio', 'Okno nie prowadzi jeszcze rozmowy — nie ma czego zatrzymać.');
      return;
    }
    void wywolaj(kanal, Command.MessageStop, { windowId: idOkna });
  }, przy);

  const odlozWersje = (): void => {
    if (dokument === null) {
      oglos('Studio', 'Okno nie prowadzi dokumentu — nie ma czego odłożyć w repozytorium sesji.');
      return;
    }
    void zapiszDokument(kanal, dokument.id, wezly, (zapisany) => {
      dokument = zapisany;
      opiszDokument(zapisany);
    });
  };

  // Odłożenie wersji w repozytorium sesji ma mieć wyzwalacz.
  wezly.kanwa.addEventListener('keydown', (zdarzenie) => {
    if (!zdarzenie.ctrlKey || zdarzenie.key.toLowerCase() !== 's') return;
    zdarzenie.preventDefault();
    odlozWersje();
  }, przy);

  wezly.zapiszWersje?.addEventListener('click', () => {
    odlozWersje();
  }, przy);

  wezly.porownajWersje?.addEventListener('click', () => {
    void porownajWersjeDokumentu(idKarty).then((porownane) => {
      if (porownane) return;
      oglos('Studio', 'Okno nie prowadzi dokumentu — nie ma czego porównać.');
    });
  }, przy);

  // Wstążka jest drugim wyzwalaczem operacji, obok przycisku Tools Panel.
  wezly.uruchomOperacje?.addEventListener('click', () => {
    void uruchomOperacjeDokumentu(idKarty).then((uruchomiona) => {
      if (uruchomiona) return;
      oglos('Studio', 'Okno nie prowadzi jeszcze rozmowy — Tools Panel nie stoi.');
    });
  }, przy);

  wezly.kanwa.addEventListener('input', () => {
    odswiezPasStanu(wezly, dokument);
  }, przy);

  // Dokument stoi w oknie komunikacji, więc przycisk czeka na stanowisko.
  wezly.nowyDokument?.addEventListener('click', () => {
    if (idOkna === '') {
      oglos('Studio', 'Dokument powstaje w oknie rozmowy — wyślij pierwszą wiadomość.');
      return;
    }
    void zalozDokument(kanal, idOkna, wezly).then((zalozony) => {
      dokument = zalozony;
      opiszDokument(zalozony);
    });
  }, przy);

  // Uchwyt zgłoszony rozdzielaczowi liczy się jako odbiorca; cudze okno odpada.
  odlaczenia.push(
    zglosUchwyt(EventType.MessageChanged, (tresc) => {
      if (tresc.message.windowId !== idOkna) return;
      zmienWpis(tresc.change, tresc.message);
    }),
  );

  // Numer fragmentu i znacznik domknięcia stoją w kopercie, nie w treści.
  odlaczenia.push(
    zglosUchwyt(EventType.StreamChunk, (tresc, koperta) => {
      if (tresc.windowId !== idOkna) return;
      przyjmijFragment(tresc, koperta.seq ?? 0, koperta.done === true);
    }),
  );

  odlaczenia.push(
    zglosUchwyt(EventType.SessionChanged, (tresc) => {
      if (idSesji === '' || tresc.session.id !== idSesji) return;
      if (tresc.change === ChangeKind.Deleted) {
        zdejmijStanowisko();
        idSesji = '';
        oglos('Studio', 'Sesja tego okna została usunięta.', 'ostrzezenie');
        return;
      }
      const tytul = wezly.szyna.querySelector('.pt-pozycja[aria-current="true"] .pt-pozycja-tytul');
      if (tytul !== null) tytul.textContent = tresc.session.title ?? 'Sesja bez nazwy';
    }),
  );

  // Okno zamknięte poza wiązaniem: następna wiadomość zakłada stanowisko od nowa.
  odlaczenia.push(
    zglosUchwyt(EventType.WindowChanged, (tresc) => {
      if (idOkna === '' || tresc.window.id !== idOkna) return;
      if (tresc.change !== ChangeKind.Deleted && tresc.window.status !== WindowStatus.Closed) return;
      zdejmijStanowisko();
      oglos('Studio', 'Okno rozmowy zostało zamknięte w rdzeniu.', 'ostrzezenie');
    }),
  );

  odlaczenia.push(
    zglosUchwyt(EventType.StudioDocumentChanged, (tresc) => {
      if (dokument === null || tresc.document.id !== dokument.id) return;
      if (tresc.change === ChangeKind.Deleted) {
        dokument = null;
        opiszDokument(null);
        return;
      }
      dokument = tresc.document;
      opiszDokument(dokument);
    }),
  );

  opiszDokument(null);
  void opiszKanal(kanal, nazwaSrodowiska, korzen);
  void opiszWyborModelu(kanal, korzen);
  opiszWyborNakladu(korzen);
  // Okno wskazane stoi już w rdzeniu: historia wraca przy montażu.
  if (idOknaStojacego !== '') void zapewnijStanowisko();
  return true;
}

export function zwolnijStudio(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

// Rodzaj `final` zastępuje treść w całości, `error` domyka turę własnym zdaniem.
function zlozTresc(dotychczasowa: string, fragment: StreamChunkEvent): string {
  if (fragment.kind === ChunkKind.Text) return dotychczasowa + (fragment.text ?? '');
  if (fragment.kind === ChunkKind.Final || fragment.kind === ChunkKind.Error) {
    return fragment.text ?? dotychczasowa;
  }
  return dotychczasowa;
}

// Woła się przy montażu okna: odczyt z rdzenia idzie osobno, po stanowisku.
function zdejmijTrescPrzykladowa(wezly: WezlyStudia, korzen: Element): void {
  wezly.historia.replaceChildren();
  wezly.szyna.replaceChildren();
  wezly.kanwa.replaceChildren();
  zdejmijZnacznikiBezZrodla(korzen);
  zdejmijTrescPrzykladowaPlanu(korzen);
  zdejmijTrescPrzykladowaRoznic(korzen);
  zdejmijTrescPrzykladowaRepozytorium(korzen);
  zdejmijTrescPrzykladowaPlikow(korzen);
  zdejmijTrescPrzykladowaNarzedzi(korzen);
  zdejmijTrescPrzykladowaPodgladu(korzen);
  zdejmijTrescPrzykladowaStanu(korzen);
  zdejmijTrescPrzykladowaFormatowania(korzen);
  zdejmijTrescPrzykladowaStylow(korzen);
  zdejmijTrescPrzykladowaStruktur(korzen);
  zdejmijTrescPrzykladowaFormatu(korzen);
  zdejmijTrescPrzykladowaAutozapisu(korzen);
  // Pas stanu prototypu niesie miary cudzego dokumentu; okno bez dokumentu ma zero słów.
  odswiezPasStanu(wezly, null);
}

// Wstążka nie ma węzła na format wydania ani profil; prowadzi je panel podglądu.
function zdejmijCzynnosciWstazkiBezPokrycia(korzen: Element): void {
  korzen.querySelector('.st-wstazka [data-etykietka="Podgląd wydruku"]')?.remove();
  const drugorzedne = korzen.querySelectorAll('.st-wstazka-grupa--drugorzedna button');
  for (const [numer, czynnosc] of [...drugorzedne].entries()) {
    if (numer > 0) czynnosc.remove();
  }
}

// Znaczniki bez pokrycia: gałąź, ścieżka, miara różnicy, źródła, pas czytelności.
function zdejmijZnacznikiBezZrodla(korzen: Element): void {
  for (const znacznik of korzen.querySelectorAll('.sta-kontekst-akcji .sta-chip')) {
    if (znacznik.classList.contains('sta-chip--srodowisko')) continue;
    znacznik.remove();
  }
  for (const zrodlo of korzen.querySelectorAll('.sta-zrodlo')) zrodlo.remove();
  korzen.querySelector('.sta-kom-monitor')?.remove();
}

async function opiszWyborModelu(kanal: Kanal, korzen: Element): Promise<void> {
  const znak = korzen.querySelector('.sta-chip--model');
  const spis = korzen.querySelector('#pop-model');
  if (znak === null || spis === null) return;
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  if (!wynik.udany || wynik.wynik === undefined) return;
  const kanaly = wynik.wynik.channels;
  const wzor = spis.querySelector('.sta-popover-wiersz');
  if (wzor === null || kanaly.length === 0) return;
  const czysty = wzor.cloneNode(true) as HTMLElement;
  for (const wiersz of [...spis.querySelectorAll('.sta-popover-wiersz')]) wiersz.remove();
  for (const [numer, kanalModelu] of kanaly.entries()) {
    const wiersz = czysty.cloneNode(true) as HTMLElement;
    const etykieta = wiersz.querySelector('.sta-popover-etykieta');
    if (etykieta !== null) etykieta.textContent = kanalModelu.model ?? kanalModelu.name;
    const kolejnosc = wiersz.querySelector('.pt-mono');
    if (kolejnosc !== null) kolejnosc.textContent = String(numer + 1);
    spis.appendChild(wiersz);
  }
  const pierwszy = kanaly[0];
  if (pierwszy !== undefined) {
    znak.childNodes[0]?.replaceWith(pierwszy.model ?? pierwszy.name);
  }
}

// Prototyp niesie pięć stopni suwaka, a rdzeń zna trzy nakłady kontraktu.
function opiszWyborNakladu(korzen: Element): void {
  const spis = korzen.querySelector('#pop-wysilek');
  const znak = korzen.querySelector('[data-popover="pop-wysilek"]');
  if (spis === null || znak === null) return;
  const naklady = Object.values(ReasoningEffort);
  const suwak = spis.querySelector<HTMLInputElement>('.sta-suwak');
  const tytul = spis.querySelector('.sta-popover-tytul');
  const nazwij = (stopien: number): void => {
    const naklad = naklady[stopien] ?? naklady[0];
    if (naklad === undefined) return;
    znak.childNodes[0]?.replaceWith(naklad);
    if (tytul !== null) tytul.textContent = 'Wysiłek — ' + naklad;
  };
  if (suwak !== null) {
    suwak.min = '0';
    suwak.max = String(naklady.length - 1);
    suwak.value = String(naklady.length - 1);
    suwak.addEventListener('input', () => {
      nazwij(Number(suwak.value));
    });
  }
  nazwij(naklady.length - 1);
}

async function opiszKanal(kanal: Kanal, nazwaSrodowiska: string, korzen: Element): Promise<void> {
  const naglowek = korzen.querySelector('.sta-kom-naglowek');
  if (naglowek === null) return;
  const pola = [...naglowek.querySelectorAll('.sta-kom-pole')];
  wpiszPole(pola, 'Środowisko', nazwaSrodowiska);
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  const kanalModelu = wynik.udany ? wynik.wynik?.channels[0] : undefined;
  wpiszPole(pola, 'Model', kanalModelu?.model ?? kanalModelu?.name ?? '');
  wpiszPole(pola, 'Wysiłek', '');
}

function opiszDokumentKarty(korzen: Element, dokument: StudioDocument | null): void {
  const nazwa = dokument?.title ?? 'Dokument bez nazwy';
  const miejsca = '.st-wstazka-sesja span, #panel-editor .sta-okno-znacznik,'
    + ' .dn-karta--robocza .dn-karta-widoku-nazwa';
  for (const wezel of korzen.querySelectorAll(miejsca)) wezel.textContent = nazwa;
  // Znak dokumentu bieżącego niesie ikonę obok napisu, więc idzie w nim sam napis.
  const znak = korzen.querySelector('.sta-kom .sta-okno-belka .sta-chip[title="Dokument bieżący"]');
  const napis = znak === null ? undefined : [...znak.childNodes].find((w) => w.nodeType === Node.TEXT_NODE);
  if (napis !== undefined) napis.textContent = nazwa;
  wpiszWersje(korzen, dokument);
}

function wpiszWersje(korzen: Element, dokument: StudioDocument | null): void {
  for (const wezel of korzen.querySelectorAll('.st-wstazka-stan .st-miara')) {
    if (dokument?.versionId === undefined) wezel.remove();
    else wezel.textContent = `wersja ${dokument.versionId}`;
  }
}

interface Stanowisko {
  idOkna: string;
  idSesji: string;
}

// Okno podane przez wołającego stoi już w rejestrze, więc `window.create` nie pada.
async function otworzStanowisko(
  kanal: Kanal,
  wezly: WezlyStudia,
  wzorPozycji: HTMLElement | null,
  wstawWpis: (wiadomosc: Message) => void,
  idOknaStojacego: string,
  idKarty: string,
  przy: AddEventListenerOptions,
): Promise<Stanowisko | null> {
  const stanowisko = idOknaStojacego === ''
    ? await zalozStanowisko(kanal, idKarty)
    : await wskazStanowisko(kanal, idOknaStojacego);
  if (stanowisko === null) return null;

  const historia = await wywolaj(kanal, Command.MessageList, { windowId: stanowisko.idOkna });
  if (historia.udany && historia.wynik !== undefined) {
    for (const wiadomosc of historia.wynik.messages) wstawWpis(wiadomosc);
  }

  await wypelnijSzyne(kanal, wezly.szyna, wzorPozycji, stanowisko.idSesji, przy);
  return stanowisko;
}

// Okno stojące w sesji bez karty wraca do pracy; inaczej powstaje nowe.
async function zalozStanowisko(kanal: Kanal, idKarty: string): Promise<Stanowisko | null> {
  const idSesji = await zapewnijSesje(kanal, idKarty, 'Studio');
  if (idSesji === '') return null;

  const modul = await wskazModulStudia(kanal);
  if (modul === '') return null;

  // Wykaz okien idzie przed założeniem: zakładanie mnożyłoby okna rdzenia.
  const stojace = await wskazOknoWolne(kanal, idSesji, modul);
  if (stojace !== '') return { idOkna: stojace, idSesji };

  const kanalModelu = await wskazKanalModelu(kanal);
  if (kanalModelu === '') return null;

  const okno = await wywolaj(kanal, Command.WindowCreate, {
    sessionId: idSesji,
    moduleId: modul,
    modelChannelId: kanalModelu,
    workingDirs: [],
    executionEnv: ExecutionEnv.Local,
    permissionMode: PermissionMode.Manual,
    windowRole: WindowRole.Standalone,
  });
  if (!okno.udany || okno.wynik === undefined) return null;
  return { idOkna: okno.wynik.window.id, idSesji };
}

async function wskazStanowisko(kanal: Kanal, idOkna: string): Promise<Stanowisko | null> {
  const wynik = await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return { idOkna, idSesji: wynik.wynik.window.sessionId };
}

async function wskazOknoWolne(kanal: Kanal, idSesji: string, modul: string): Promise<string> {
  const wynik = await wywolaj(kanal, Command.WindowList, {
    sessionId: idSesji,
    status: WindowStatus.Open,
  });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  const zajete = oknaZajetePrzezKarty();
  const wolne = wynik.wynik.windows.find(
    (okno) => okno.moduleId === modul && !zajete.has(okno.id),
  );
  return wolne?.id ?? '';
}

function oknaZajetePrzezKarty(): ReadonlySet<string> {
  const zajete = new Set<string>();
  for (const okno of oknaRobocze()) {
    for (const karta of kartyOkna(okno)) {
      if (karta.idOknaKomunikacji !== '') zajete.add(karta.idOknaKomunikacji);
    }
  }
  return zajete;
}

async function wskazModulStudia(kanal: Kanal): Promise<string> {
  const wynik = await wywolaj(kanal, Command.ModuleList, {});
  if (!wynik.udany || wynik.wynik === undefined) return '';
  const modul = wynik.wynik.modules.find((pozycja) => pozycja.code === KOD_MODULU_STUDIO);
  return modul?.id ?? '';
}

// Wyboru kanału znacznik prototypu nie niesie, a `window.create` go wymaga.
async function wskazKanalModelu(kanal: Kanal): Promise<string> {
  const wynik = await wywolaj(kanal, Command.ChannelList, { enabledOnly: true });
  if (!wynik.udany || wynik.wynik === undefined) return '';
  return wynik.wynik.channels[0]?.id ?? '';
}

// Przełącza biblioteka `zakladki-paneli.js`, więc odczyt stanu idzie po jej nasłuchu.
function zwiazUkladKart(
  kanal: Kanal,
  idKarty: string,
  korzen: Element,
  przy: AddEventListenerOptions,
): void {
  if (idKarty === '') return;
  void odczytajUkladPaneli(kanal, idKarty, PANEL_KART_STUDIA).then((sekcje) => {
    if (sekcje !== null && sekcje.length > 0) naniesUkladKart(korzen, sekcje);
  });
  const naPrzelaczenie = (zdarzenie: Event): void => {
    const cel = zdarzenie.target;
    if (!(cel instanceof Element) || cel.closest('.st-karty [role="tab"][data-karta]') === null) return;
    globalThis.setTimeout(() => {
      void zapiszUkladPaneli(kanal, idKarty, PANEL_KART_STUDIA, zbierzUkladKart(korzen));
    }, 0);
  };
  korzen.addEventListener('click', naPrzelaczenie, przy);
}

function zbierzUkladKart(korzen: Element): PanelSection[] {
  const karty = korzen.querySelectorAll<HTMLElement>('.st-karty [role="tab"][data-karta]');
  return [...karty].map((karta, numer) => ({
    id: karta.dataset.karta ?? '',
    order: numer + 1,
    collapsed: false,
    hidden: karta.getAttribute('aria-selected') !== 'true',
  }));
}

function naniesUkladKart(korzen: Element, sekcje: PanelSection[]): void {
  const otwarta = [...sekcje]
    .sort((pierwsza, druga) => pierwsza.order - druga.order)
    .find((sekcja) => sekcja.hidden !== true);
  if (otwarta === undefined) return;
  const karta = korzen.querySelector<HTMLElement>(
    `.st-karty [role="tab"][data-karta="${otwarta.id}"]`,
  );
  if (karta === null || karta.getAttribute('aria-selected') === 'true') return;
  przelaczKarte(korzen, karta);
}

function przelaczKarte(korzen: Element, wybrana: HTMLElement): void {
  for (const karta of korzen.querySelectorAll<HTMLElement>('.st-karty [role="tab"][data-karta]')) {
    const czynna = karta === wybrana;
    karta.setAttribute('aria-selected', czynna ? 'true' : 'false');
    karta.tabIndex = czynna ? 0 : -1;
    const oznaczenie = karta.getAttribute('aria-controls') ?? '';
    const panel = oznaczenie === '' ? null : korzen.querySelector<HTMLElement>(`#${oznaczenie}`);
    if (panel !== null) panel.hidden = !czynna;
  }
}

async function wypelnijSzyne(
  kanal: Kanal,
  szyna: HTMLElement,
  wzor: HTMLElement | null,
  idBiezacej: string,
  przy: AddEventListenerOptions,
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
  }, przy);
}

function zbudujPozycje(wzor: HTMLElement, sesja: Session, biezaca: boolean): HTMLElement {
  const pozycja = wzor.cloneNode(true) as HTMLElement;
  const tytul = pozycja.querySelector('.pt-pozycja-tytul');
    // Identyfikator jest oznaczeniem magazynu; Operator czyta nazwę swojej pracy.
  if (tytul !== null) tytul.textContent = sesja.title ?? 'Sesja bez nazwy';
  if (biezaca) pozycja.setAttribute('aria-current', 'true');
  else {
    pozycja.removeAttribute('aria-current');
    pozycja.querySelector('.pt-tetno')?.remove();
  }
  return pozycja;
}

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

  return wczytajDokumentDoKanwy(kanal, dokument, wezly);
}

// Pustka znaczy okno bez dokumentu: kanwa zostaje wtedy w stanie pustym.
async function otworzDokumentOkna(
  kanal: Kanal,
  idOkna: string,
  wezly: WezlyStudia,
): Promise<StudioDocument | null> {
  const otwarty = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!otwarty.udany || otwarty.wynik === undefined) return null;
  return wczytajDokumentDoKanwy(kanal, otwarty.wynik.document, wezly);
}

async function wczytajDokumentDoKanwy(
  kanal: Kanal,
  dokument: StudioDocument,
  wezly: WezlyStudia,
): Promise<StudioDocument> {
  const tekst = await wywolaj(kanal, Command.StudioTextGet, { documentId: dokument.id });
  wezly.kanwa.textContent = tekst.udany && tekst.wynik !== undefined ? tekst.wynik.text : '';
  // Kanwa prototypu stoi zamknięta na edycję, bo niosła dokument przykładowy.
  wezly.kanwa.setAttribute('contenteditable', 'true');
  odswiezPasStanu(wezly, dokument);
  return dokument;
}

async function zapiszDokument(
  kanal: Kanal,
  idDokumentu: string,
  wezly: WezlyStudia,
  poZapisie: (dokument: StudioDocument) => void,
): Promise<void> {
  const zapisany = await zapiszDokumentZPostacia(
    kanal,
    idDokumentu,
    wezly.kanwa.textContent ?? '',
  );
  if (zapisany === null) return;
  poZapisie(zapisany);
  odswiezPasStanu(wezly, zapisany);
}

function odswiezPasStanu(wezly: WezlyStudia, dokument: StudioDocument | null): void {
  const pola = [...wezly.status.querySelectorAll('span')];
  const slowa = pola.find((pole) => pole.textContent?.startsWith('słów') === true);
  if (slowa !== undefined) slowa.textContent = `słów: ${policzSlowa(wezly.kanwa.textContent ?? '')}`;
  const wersja = pola.find((pole) => pole.textContent?.startsWith('wersja') === true);
  if (wersja !== undefined) {
    wersja.textContent = dokument?.versionId === undefined ? '' : `wersja ${dokument.versionId}`;
  }
}

function policzSlowa(tresc: string): number {
  const cialo = tresc.trim();
  return cialo === '' ? 0 : cialo.split(/\s+/u).length;
}

function zbierzWezly(korzen: Element): WezlyStudia | null {
  const kom = korzen.querySelector('.sta-kom');
  const historia = kom?.querySelector('.sta-kom-historia');
  const formularz = kom?.querySelector('form.sta-prompt');
  const pole = kom?.querySelector('textarea.sta-prompt-obszar');
  const szyna = korzen.querySelector('.st-szyna-lista');
  const kanwa = korzen.querySelector('.dn-kanwa');
  const status = korzen.querySelector('.st-status');
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
  const nowyDokument = korzen.querySelector('.st-szyna-naglowek button.dn-btn');
  return {
    historia,
    formularz,
    pole,
    zatrzymaj: wskazElement(kom?.querySelector('.sta-prompt-zatrzymaj') ?? null),
    szyna,
    nowyDokument: nowyDokument instanceof HTMLElement ? nowyDokument : null,
    kanwa,
    status,
    zapiszWersje: wskazElement(korzen.querySelector('.st-wstazka [data-etykietka="Zapisz wersję"]')),
    porownajWersje: wskazElement(
      korzen.querySelector('.st-wstazka [data-etykietka="Porównaj wersje"]'),
    ),
    uruchomOperacje: wskazElement(
      korzen.querySelector('.st-wstazka .st-wstazka-grupa--drugorzedna button'),
    ),
  };
}

function wskazElement(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? wezel : null;
}

function zdejmijWzoryWpisow(historia: HTMLElement): WzoryWpisow {
  return {
    czlowiek: sklonuj(historia.querySelector('.sta-wpis--czlowiek')),
    inteligencja: sklonuj(historia.querySelector('.sta-wpis--inteligencja')),
    system: sklonuj(historia.querySelector('.sta-wpis--system')),
  };
}

function zdejmijWzorPozycji(szyna: HTMLElement): HTMLElement | null {
  return sklonuj(szyna.querySelector('.pt-pozycja'));
}

function sklonuj(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? (wezel.cloneNode(true) as HTMLElement) : null;
}

function wzorDlaRoli(wzory: WzoryWpisow, rola: MessageRole): HTMLElement | null {
  if (rola === MessageRole.User) return wzory.czlowiek;
  if (rola === MessageRole.Assistant) return wzory.inteligencja;
  return wzory.system;
}

function nazwaNadawcy(rola: MessageRole): string {
  if (rola === MessageRole.User) return 'Operator';
  if (rola === MessageRole.Assistant) return 'Inteligencja';
  return 'System';
}

function wypelnijWpis(wpis: HTMLElement, wiadomosc: Message): void {
  const nadawca = wpis.querySelector('.sta-wpis-nadawca');
  if (nadawca !== null) nadawca.textContent = nazwaNadawcy(wiadomosc.role);
  const godzina = wpis.querySelector('.sta-wpis-godzina');
  if (godzina !== null) godzina.textContent = godzinaWpisu(wiadomosc.createdAt);
  for (const plakietka of wpis.querySelectorAll('.dn-plakietka')) plakietka.remove();
  wpiszTrescWpisu(wpis, wiadomosc.content);
}

// Pasek działań należy do biblioteki i przeżywa podmianę treści.
function wpiszTrescWpisu(wpis: HTMLElement, tekst: string): void {
  const tresc = wpis.querySelector('.sta-wpis-tresc');
  if (tresc === null) return;
  const akcje = tresc.querySelector('.sta-wpis-akcje');
  tresc.textContent = tekst;
  if (akcje !== null) tresc.appendChild(akcje);
}

function trescWpisu(wpis: HTMLElement): string {
  const tresc = wpis.querySelector('.sta-wpis-tresc');
  if (tresc === null) return '';
  const akcje = tresc.querySelector('.sta-wpis-akcje');
  return [...tresc.childNodes]
    .filter((wezel) => wezel !== akcje)
    .map((wezel) => wezel.textContent ?? '')
    .join('');
}

function godzinaWpisu(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

/* Panele wiąże się po założeniu okna komunikacji: przed nim nie ma o co pytać. */
function zwiazPanele(
  kanal: Kanal,
  idOkna: string,
  idKarty: string,
  korzen: Element,
  odlaczenia: Odsubskrybuj[],
): void {
  const zwiazane = [
    zwiazPlan(kanal, idOkna, korzen),
    zwiazRoznice(kanal, idOkna, idKarty, korzen),
    zwiazRepozytorium(kanal, idOkna, korzen),
    zwiazPliki(kanal, idOkna, korzen),
    zwiazNarzedzia(kanal, idOkna, idKarty, korzen),
    zwiazKatalogOperacji(kanal, idOkna, korzen),
    zwiazStanOkna(kanal, idOkna, korzen),
    zwiazPodglad(kanal, idOkna, korzen),
    zwiazFormatowanie(kanal, idOkna, korzen),
    zwiazStyle(kanal, idOkna, korzen),
    zwiazStruktury(kanal, idOkna, korzen),
    zwiazFormatDokumentu(kanal, idOkna, korzen),
    zwiazAutozapis(kanal, idOkna, korzen),
  ];
  for (const odlacz of zwiazane) {
    if (odlacz !== null) odlaczenia.push(odlacz);
  }
}
