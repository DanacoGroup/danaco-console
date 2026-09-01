/**
 * Wiązanie karty Studia z rdzeniem. Znacznik niesie biblioteka Właściciela —
 * ten plik nic nie buduje: słucha zdarzeń, woła komendy kontraktu i wypełnia
 * węzły karty, a wpisy i pozycje powiela z wzorów zdjętych ze znacznika.
 * Karty stoją w płótnie obok siebie: wiązanie idzie od korzenia karty i trzyma
 * stan po jej identyfikatorze.
 */

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
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { wpiszPole } from './okno-modulu.ts';
import {
  odczytajUkladPaneli,
  przypiszOknoKomunikacji,
  zapiszUkladPaneli,
} from './okna-robocze.ts';
import { zapewnijSesje } from './sesja-biezaca.ts';
import { zglosUchwyt } from './zdarzenia.ts';
import { zdejmijTrescPrzykladowaNarzedzi, zwiazNarzedzia } from './studio-narzedzia.ts';
import { zdejmijTrescPrzykladowaPlanu, zwiazPlan } from './studio-plan.ts';
import { zdejmijTrescPrzykladowaPlikow, zwiazPliki } from './studio-pliki.ts';
import { zdejmijTrescPrzykladowaPodgladu, zwiazPodglad } from './studio-podglad.ts';
import { zdejmijTrescPrzykladowaRepozytorium, zwiazRepozytorium } from './studio-repozytorium.ts';
import {
  porownajWersjeDokumentu,
  zdejmijTrescPrzykladowaRoznic,
  zwiazRoznice,
} from './studio-roznice.ts';

/** Kod modułu Studia w rejestrze rdzenia; `module.list` oddaje po nim identyfikator, którego wymaga `window.create`. */
const KOD_MODULU_STUDIO = 'studio';

/** Panel, pod którym rdzeń trzyma układ kart pasma okna roboczego Studia; sekcja to karta pasma, ukryta znaczy nieotwartą. */
const PANEL_KART_STUDIA = 'st-karty';

/** Węzły znacznika Właściciela, na których wiązanie pracuje. Brak któregokolwiek znaczy, że okno Studia nie stoi i wiązać nie ma czego. */
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
}

/** Wzory wpisu historii zdjęte z treści przykładowej, po jednym na rodzaj nadawcy — kształt wpisu rdzenia bierze się stąd, nie z kodu. */
interface WzoryWpisow {
  czlowiek: HTMLElement | null;
  inteligencja: HTMLElement | null;
  system: HTMLElement | null;
}

/** Fragment strumienia wraz z tym, co niesie o nim koperta: numerem w strumieniu i znacznikiem domknięcia. */
interface FragmentStrumienia {
  tresc: StreamChunkEvent;
  ostatni: boolean;
}

/** Wpis odpowiedzi w budowie: tekst złożony z fragmentów, numer fragmentu oczekiwanego i fragmenty, które przyszły przed swoją koleją. */
interface StrumienOdpowiedzi {
  wpis: HTMLElement;
  tekst: string;
  nastepny: number;
  odlozone: Map<number, FragmentStrumienia>;
}

/** Wiązanie jednej karty Studia: korzeń wnętrza i odłączenia jego uchwytów zdarzeń oraz nasłuchów. */
interface WiazanieKarty {
  korzen: Element;
  odlaczenia: Odsubskrybuj[];
}

/* Wiązania po identyfikatorze karty. Karty tego samego modułu stoją w płótnie
   obok siebie, więc stan wiązania nie może być jeden na moduł: każda karta ma
   własne okno, historię i uchwyty, a zamknięcie karty zwalnia tylko jej. */
const WIAZANIA = new Map<string, WiazanieKarty>();

/**
 * Wiąże wnętrze karty Studia z rdzeniem. Korzeń to element
 * `.cd-tresc--modul[data-karta]` albo identyfikator karty; zapytania o węzły
 * idą od niego. Okno stojące podaje wołający wznawiający sesję, pustka znaczy
 * okno zakładane pierwszą wiadomością. Prawda znaczy wiązanie założone; karta
 * związana już na tym korzeniu wraca fałszem.
 */
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
  WIAZANIA.set(idKarty, { korzen, odlaczenia });
  const wezly: WezlyStudia = znalezione;

  const wzory = zdejmijWzoryWpisow(wezly.historia);
  const wzorPozycji = zdejmijWzorPozycji(wezly.szyna);
  zdejmijTrescPrzykladowa(wezly, korzen);
  zdejmijCzynnosciWstazkiBezPokrycia(korzen);

  const wpisy = new Map<string, HTMLElement>();
  const strumienie = new Map<string, StrumienOdpowiedzi>();
  let idOkna = '';
  let idSesji = '';
  let dokument: StudioDocument | null = null;

  /** Opisuje wstążkę i pas tej karty dokumentem; węzły idą od korzenia karty, nie od dokumentu. */
  const opiszDokument = (opisywany: StudioDocument | null): void => {
    opiszDokumentKarty(korzen, opisywany);
  };

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

  /** Nanosi zmianę wiadomości na wpis już stojący; wpis nieznany dopisuje, usunięty zdejmuje wraz ze strumieniem, który do niego dopisywał. */
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

  /**
   * Stawia wpis odpowiedzi, do którego dopisują się fragmenty strumienia.
   * Wiadomość modelu rdzeń rozgłasza dopiero po turze, więc godziny nadania
   * nie ma jeszcze skąd wziąć; węzeł godziny zostaje pusty do tej chwili.
   */
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

  /** Przyjmuje fragment strumienia: nanosi ten, na który przyszła kolej, a przyszły odkłada, aż przyjdzie jego poprzednik. */
  function przyjmijFragment(tresc: StreamChunkEvent, numer: number, ostatni: boolean): void {
    let strumien = strumienie.get(tresc.messageId);
    if (strumien === undefined) {
      const wpis = wpisy.get(tresc.messageId) ?? wstawWpisStrumienia(tresc.messageId);
      if (wpis === null) return;
      /* Porządek liczy się od fragmentu pierwszego, jaki to okno zobaczyło:
         okno otwarte w trakcie tury zastaje strumień w połowie, a liczenie od
         jedynki zatrzymałoby każdy kolejny fragment w oczekiwaniu na
         poprzednika, który już przeszedł. */
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
    // Numer niższy od oczekiwanego znaczy fragment już naniesiony.
    if (numer !== 0 && numer < strumien.nastepny) return;
    naniesFragment(tresc.messageId, strumien, { tresc, ostatni });
  }

  /** Dokleja fragment do wpisu i wypuszcza za nim te odłożone, które właśnie doczekały swojej kolei; fragment domykający zamyka wpis. */
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

  /* Enter wysyła, Shift+Enter przechodzi do nowego wiersza. Pole jest obszarem
     tekstowym, a ten sam z siebie formularza nie zamyka — bez tego jedyną drogą
     wysłania byłby przycisk. */
  wezly.pole.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.key !== 'Enter' || zdarzenie.shiftKey) return;
    zdarzenie.preventDefault();
    wezly.formularz.requestSubmit();
  });

  /** Przycisk nowego dokumentu jest czynny wyłącznie przy stanowisku: dokument stoi w oknie komunikacji, a okno powstaje pierwszą wiadomością. */
  const ustawNowyDokument = (): void => {
    if (wezly.nowyDokument instanceof HTMLButtonElement) wezly.nowyDokument.disabled = idOkna === '';
  };

  /** Zdejmuje okno komunikacji z karty: rdzeń je zamknął albo usunął, a kolejna wiadomość zakłada stanowisko od nowa. */
  const zdejmijStanowisko = (): void => {
    idOkna = '';
    ustawNowyDokument();
  };

  /* Sesja powstaje dopiero pierwszą wiadomością wysłaną do modelu — nie
     otwarciem okna ani wejściem w moduł. Do tej chwili Operator chodzi po
     produkcie swobodnie i żadna sesja po nim nie zostaje. Wyjątkiem jest okno
     wskazane przez wołającego: ono stoi już w rejestrze i wchodzi przy montażu. */
  const zapewnijStanowisko = async (): Promise<string> => {
    if (idOkna !== '') return idOkna;
    const stanowisko = await otworzStanowisko(
      kanal, wezly, wzorPozycji, wstawWpis, idOknaStojacego, idKarty,
    );
    if (stanowisko === null) return '';
    idOkna = stanowisko.idOkna;
    idSesji = stanowisko.idSesji;
    ustawNowyDokument();
    /* Karta dostaje okno komunikacji: po nim zamknięcie karty zamyka okno
       w rdzeniu, a układ kart pasma ma do czego przylgnąć. Karta zdjęta, zanim
       rdzeń odpowiedział, zostaje bez okna — i mówi to wprost. */
    if (!przypiszOknoKomunikacji(idKarty, idOkna)) {
      oglos('Studio', 'Okno rozmowy stanęło, gdy jego karta zeszła z pasma — '
        + 'zamknięcie karty nie zamknie tego okna.', 'ostrzezenie');
    }
    zwiazUkladKart(kanal, idKarty, korzen, odlaczenia);
    zwiazPanele(kanal, idOkna, idKarty, korzen, odlaczenia);
    /* Okno stojące prowadzi już swój dokument. Zakładanie nowego przy powrocie
       do sesji mnożyłoby dokumenty przy każdym wejściu w moduł. */
    dokument = idOknaStojacego === ''
      ? await zalozDokument(kanal, idOkna, wezly)
      : await otworzDokumentOkna(kanal, idOkna, wezly);
    opiszDokument(dokument);
    return idOkna;
  };

  /**
   * Wysyła wiadomość i czyści pole dopiero po jej przyjęciu przez rdzeń.
   * Rdzeń odmawia wiadomości skierowanej do okna prowadzącego turę — pole
   * wyczyszczone przed wysyłką zabrałoby Operatorowi tekst bez słowa.
   */
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
    /* Pole czyści się tylko wtedy, gdy nadal niesie wysłany tekst: Operator,
       który zdążył dopisać kolejne zdanie, nie traci go. */
    if (wezly.pole.value.trim() === tresc) wezly.pole.value = '';
  };

  wezly.formularz.addEventListener('submit', (zdarzenie) => {
    zdarzenie.preventDefault();
    const tresc = wezly.pole.value.trim();
    if (tresc === '') return;
    void wyslijWiadomosc(tresc);
  });

  /* Zatrzymanie jest czynne niezależnie od stanu tury: rdzeń odpowiada na
     `message.stop` zawsze, a brak tury w biegu oddaje `stopped` równe fałszowi,
     nie odmowę. Okno bez stanowiska nie ma czego zatrzymać i mówi to wprost. */
  wezly.zatrzymaj?.addEventListener('click', () => {
    if (idOkna === '') {
      oglos('Studio', 'Okno nie prowadzi jeszcze rozmowy — nie ma czego zatrzymać.');
      return;
    }
    void wywolaj(kanal, Command.MessageStop, { windowId: idOkna });
  });

  /** Odkłada treść kanwy jako nową wersję w repozytorium sesji; okno bez dokumentu nie ma czego odłożyć i mówi to wprost. */
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

  // Ctrl+S zapisuje treść kanwy: pas stanu prototypu mówi o zapisie
  // automatycznym, ale odłożenie wersji w repozytorium sesji ma mieć wyzwalacz.
  wezly.kanwa.addEventListener('keydown', (zdarzenie) => {
    if (!zdarzenie.ctrlKey || zdarzenie.key.toLowerCase() !== 's') return;
    zdarzenie.preventDefault();
    odlozWersje();
  });

  wezly.zapiszWersje?.addEventListener('click', () => {
    odlozWersje();
  });

  wezly.porownajWersje?.addEventListener('click', () => {
    void porownajWersjeDokumentu(idKarty).then((porownane) => {
      if (porownane) return;
      oglos('Studio', 'Okno nie prowadzi dokumentu — nie ma czego porównać.');
    });
  });

  wezly.kanwa.addEventListener('input', () => {
    odswiezPasStanu(wezly, dokument);
  });

  /* Nowy dokument nie zakłada stanowiska: sesja powstaje wyłącznie pierwszą
     wiadomością, a dokument stoi w oknie komunikacji. Przycisk jest nieczynny
     do tej chwili; kliknięcie w kartę bez okna mówi to wprost. */
  wezly.nowyDokument?.addEventListener('click', () => {
    if (idOkna === '') {
      oglos('Studio', 'Dokument powstaje w oknie rozmowy — wyślij pierwszą wiadomość.');
      return;
    }
    void zalozDokument(kanal, idOkna, wezly).then((zalozony) => {
      dokument = zalozony;
      opiszDokument(zalozony);
    });
  });

  /* Zdarzenia idą rozdzielaczem wspólnym: uchwyt zgłoszony tam liczy się jako
     odbiorca, a zdarzenie cudzego okna albo cudzej sesji jest pomijane. */
  odlaczenia.push(
    zglosUchwyt(EventType.MessageChanged, (tresc) => {
      if (tresc.message.windowId !== idOkna) return;
      zmienWpis(tresc.change, tresc.message);
    }),
  );

  /* Odpowiedź modelu przyrasta w oknie fragmentami: pełną wiadomość rdzeń
     rozgłasza dopiero po turze. Numer fragmentu i znacznik domknięcia stoją
     w kopercie, nie w treści zdarzenia. */
  odlaczenia.push(
    zglosUchwyt(EventType.StreamChunk, (tresc, koperta) => {
      if (tresc.windowId !== idOkna) return;
      przyjmijFragment(tresc, koperta.seq ?? 0, koperta.done === true);
    }),
  );

  /* Sesja własna zmieniona z innego okna albo urządzenia: nazwa wraca na
     pozycję bieżącą szyny, a sesja usunięta zostawia okno bez stanowiska. */
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

  /* Okno własne zamknięte poza tym wiązaniem — usunięte albo w stanie
     zamkniętym po `session.close`: następna wiadomość zakłada stanowisko od
     nowa, zamiast wracać odmową do okna, które już nie przyjmuje. */
  odlaczenia.push(
    zglosUchwyt(EventType.WindowChanged, (tresc) => {
      if (idOkna === '' || tresc.window.id !== idOkna) return;
      if (tresc.change !== ChangeKind.Deleted && tresc.window.status !== WindowStatus.Closed) return;
      zdejmijStanowisko();
      oglos('Studio', 'Okno rozmowy zostało zamknięte w rdzeniu.', 'ostrzezenie');
    }),
  );

  /* Dokument prowadzony w oknie zmieniony przez model albo inne okno: nazwa
     i wersja na wstążce idą za rdzeniem, treść kanwy zostaje pracą Operatora. */
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

  /* Stan pusty dokumentu wchodzi od razu: nazwa pracy z prototypu jest treścią
     przykładową, a okno staje, zanim powstanie sesja i dokument. */
  opiszDokument(null);
  ustawNowyDokument();
  void opiszKanal(kanal, nazwaSrodowiska, korzen);
  void opiszWyborModelu(kanal, korzen);
  opiszWyborNakladu(korzen);
  /* Okno wskazane przez wołającego stoi już w rdzeniu, więc historia rozmowy
     wraca przy montażu, bez czekania na pierwszą wiadomość. */
  if (idOknaStojacego !== '') void zapewnijStanowisko();
  return true;
}

/** Zwalnia wiązanie karty: zdejmuje jej uchwyty zdarzeń i nasłuchy. Karta bez wiązania nie robi nic. */
export function zwolnijStudio(idKarty: string): void {
  const wiazanie = WIAZANIA.get(idKarty);
  if (wiazanie === undefined) return;
  for (const odlacz of wiazanie.odlaczenia) odlacz();
  WIAZANIA.delete(idKarty);
}

/** Korzeń karty ze wskazania: element wprost albo element odszukany w płótnie po identyfikatorze karty. */
function korzenKarty(wskazanie: Element | string): Element | null {
  if (typeof wskazanie !== 'string') return wskazanie;
  return document.querySelector(`.cd-tresc--modul[data-karta="${wskazanie}"]`);
}

/**
 * Treść wpisu po fragmencie strumienia. Tekstowy dokleja się do dotychczasowej,
 * a wersja ostateczna zastępuje ją w całości — tak stanowi kontrakt rodzaju
 * `final`. Fragment błędu domyka turę i jego zdanie jest jedyną treścią wpisu.
 * Pozostałe rodzaje nie zmieniają wpisu: nie ma go w nim na czym pokazać.
 */
function zlozTresc(dotychczasowa: string, fragment: StreamChunkEvent): string {
  if (fragment.kind === ChunkKind.Text) return dotychczasowa + (fragment.text ?? '');
  if (fragment.kind === ChunkKind.Final || fragment.kind === ChunkKind.Error) {
    return fragment.text ?? dotychczasowa;
  }
  return dotychczasowa;
}

/**
 * Zdejmuje treść przykładową okna wraz z treścią sześciu jego paneli. Woła się
 * przy montażu okna, przed powstaniem stanowiska: odczyt z rdzenia idzie osobno,
 * po założeniu okna komunikacji, a do tej chwili Operator nie ma prawa zobaczyć
 * dokumentu z prototypu i wziąć go za własną pracę.
 */
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
}

/**
 * Zdejmuje ze wstążki czynności, których wywołać nie ma czym. Podgląd wydruku
 * wymaga formatu wydania i miejsca na strony, wywołanie operacji — nazwy
 * operacji i jej zakresu, przekazanie do Biblioteki — profilu wydania; wstążka
 * nie ma węzła, w którym Operator poda którąkolwiek z tych wartości.
 */
function zdejmijCzynnosciWstazkiBezPokrycia(korzen: Element): void {
  korzen.querySelector('.st-wstazka [data-etykietka="Podgląd wydruku"]')?.remove();
  korzen.querySelector('.st-wstazka .st-wstazka-grupa--drugorzedna')?.remove();
}

/**
 * Zdejmuje znaczniki bez pokrycia w kontrakcie: nazwę gałęzi, ścieżkę
 * repozytorium i miarę różnicy w pasie czynności, wskazania źródła nad polem
 * wpisu oraz pas czytelności, dla którego rdzeń nie ma ani jednej miary.
 */
function zdejmijZnacznikiBezZrodla(korzen: Element): void {
  for (const znacznik of korzen.querySelectorAll('.sta-kontekst-akcji .sta-chip')) {
    if (znacznik.classList.contains('sta-chip--srodowisko')) continue;
    znacznik.remove();
  }
  for (const zrodlo of korzen.querySelectorAll('.sta-zrodlo')) zrodlo.remove();
  korzen.querySelector('.sta-kom-monitor')?.remove();
}

/** Wpisuje w wybór modelu kanały rejestru rdzenia; wiersz wzorcowy powiela się na każdy kanał, a kanał czynny nazywa sam znacznik wyboru. */
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

/**
 * Wpisuje w wybór nakładu wartości kontraktu. Prototyp niesie pięciostopniowy
 * suwak i nazwę spoza kontraktu; rdzeń zna trzy nakłady, więc suwak dostaje
 * granice trzech stopni, a nazwa bierze się z wybranego stopnia.
 */
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

/** Wpisuje w nagłówek okna komunikacji środowisko wejścia i model kanału; pole wysiłku znika, bo kontrakt nie niesie jego wartości. */
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

/** Nadaje wstążce i karcie edytora tytuł dokumentu; dokument bez nadanej nazwy dostaje nazwany stan pusty, nie własny identyfikator. */
function opiszDokumentKarty(korzen: Element, dokument: StudioDocument | null): void {
  const nazwa = dokument?.title ?? 'Dokument bez nazwy';
  const miejsca = '.st-wstazka-sesja span, .sta-okno-znacznik, .dn-karta--robocza .dn-karta-widoku-nazwa';
  for (const wezel of korzen.querySelectorAll(miejsca)) wezel.textContent = nazwa;
  wpiszWersje(korzen, dokument);
}

/** Nanosi wersję dokumentu na miarę wstążki; dokument bez wersji w repozytorium sesji nie ma czego pokazać, więc miara znika. */
function wpiszWersje(korzen: Element, dokument: StudioDocument | null): void {
  for (const wezel of korzen.querySelectorAll('.st-wstazka-stan .st-miara')) {
    if (dokument?.versionId === undefined) wezel.remove();
    else wezel.textContent = `wersja ${dokument.versionId}`;
  }
}

/** Okno komunikacji wraz z sesją, w której stoi; szyna dokumentów wskazuje po niej sesję bieżącą. */
interface Stanowisko {
  idOkna: string;
  idSesji: string;
}

/**
 * Wskazuje stanowisko karty, wczytuje jego historię i szynę sesji; pustka
 * znaczy odmowę rdzenia. Okno podane przez wołającego stoi już w rejestrze,
 * więc `window.create` nie pada — pada wyłącznie dla okna zakładanego
 * pierwszą wiadomością.
 */
async function otworzStanowisko(
  kanal: Kanal,
  wezly: WezlyStudia,
  wzorPozycji: HTMLElement | null,
  wstawWpis: (wiadomosc: Message) => void,
  idOknaStojacego: string,
  idKarty: string,
): Promise<Stanowisko | null> {
  const stanowisko = idOknaStojacego === ''
    ? await zalozStanowisko(kanal, idKarty)
    : await wskazStanowisko(kanal, idOknaStojacego);
  if (stanowisko === null) return null;

  const historia = await wywolaj(kanal, Command.MessageList, { windowId: stanowisko.idOkna });
  if (historia.udany && historia.wynik !== undefined) {
    for (const wiadomosc of historia.wynik.messages) wstawWpis(wiadomosc);
  }

  await wypelnijSzyne(kanal, wezly.szyna, wzorPozycji, stanowisko.idSesji);
  return stanowisko;
}

/**
 * Zakłada dla karty okno komunikacji Studia w sesji jej okna roboczego —
 * stojącej albo założonej tą wiadomością. Okno jest zawsze nowe: karta bez
 * okna to nowa praca, a okno stojące w sesji należy do innej karty i wchodzi
 * wyłącznie przy wznowieniu, jako okno wskazane wołającemu. Pustka znaczy
 * odmowę rdzenia na którymkolwiek kroku.
 */
async function zalozStanowisko(kanal: Kanal, idKarty: string): Promise<Stanowisko | null> {
  const idSesji = await zapewnijSesje(kanal, idKarty, 'Studio');
  if (idSesji === '') return null;

  const modul = await wskazModulStudia(kanal);
  if (modul === '') return null;

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

/** Odczytuje sesję okna stojącego z rejestru okien; pustka znaczy okno, którego rdzeń nie zna. */
async function wskazStanowisko(kanal: Kanal, idOkna: string): Promise<Stanowisko | null> {
  const wynik = await wywolaj(kanal, Command.WindowStateGet, { windowId: idOkna });
  if (!wynik.udany || wynik.wynik === undefined) return null;
  return { idOkna, idSesji: wynik.wynik.window.sessionId };
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

/**
 * Układ kart pasma okna roboczego idzie do rdzenia przy każdym przełączeniu
 * karty i wraca przy założeniu stanowiska. Sekcją jest karta pasma: kolejność
 * według pasma, ukryta znaczy nieotwartą. Przełącza biblioteka
 * `zakladki-paneli.js` na kliknięcie, więc odczyt stanu idzie po jej nasłuchu.
 */
function zwiazUkladKart(
  kanal: Kanal,
  idKarty: string,
  korzen: Element,
  odlaczenia: Odsubskrybuj[],
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
  korzen.addEventListener('click', naPrzelaczenie);
  odlaczenia.push(() => {
    korzen.removeEventListener('click', naPrzelaczenie);
  });
}

/** Układ kart pasma odczytany ze znacznika karty: kolejność według pasma, ukryta jest karta nieotwarta. */
function zbierzUkladKart(korzen: Element): PanelSection[] {
  const karty = korzen.querySelectorAll<HTMLElement>('.st-karty [role="tab"][data-karta]');
  return [...karty].map((karta, numer) => ({
    id: karta.dataset.karta ?? '',
    order: numer + 1,
    collapsed: false,
    hidden: karta.getAttribute('aria-selected') !== 'true',
  }));
}

/** Otwiera kartę pasma wskazaną układem z rdzenia; kliknięcie idzie przez bibliotekę, która przełącza panele. */
function naniesUkladKart(korzen: Element, sekcje: PanelSection[]): void {
  const otwarta = [...sekcje]
    .sort((pierwsza, druga) => pierwsza.order - druga.order)
    .find((sekcja) => sekcja.hidden !== true);
  if (otwarta === undefined) return;
  const karta = korzen.querySelector<HTMLElement>(
    `.st-karty [role="tab"][data-karta="${otwarta.id}"]`,
  );
  if (karta === null || karta.getAttribute('aria-selected') === 'true') return;
  karta.click();
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

  return wczytajDokumentDoKanwy(kanal, dokument, wezly);
}

/**
 * Otwiera dokument prowadzony w oknie stojącym. Pustka znaczy okno, w którym
 * dokumentu jeszcze nie ma — kanwa zostaje wtedy w stanie pustym, a dokument
 * powstaje przyciskiem nowego dokumentu, tak samo jak w oknie świeżym.
 */
async function otworzDokumentOkna(
  kanal: Kanal,
  idOkna: string,
  wezly: WezlyStudia,
): Promise<StudioDocument | null> {
  const otwarty = await wywolaj(kanal, Command.StudioDocumentOpen, { windowId: idOkna });
  if (!otwarty.udany || otwarty.wynik === undefined) return null;
  return wczytajDokumentDoKanwy(kanal, otwarty.wynik.document, wezly);
}

/** Wstawia treść dokumentu w kanwę wraz z pasem stanu i otwiera ją na pracę Operatora. */
async function wczytajDokumentDoKanwy(
  kanal: Kanal,
  dokument: StudioDocument,
  wezly: WezlyStudia,
): Promise<StudioDocument> {
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

/** Wskazuje węzły karty Studia od jej korzenia; pustka znaczy, że znacznik Studia w karcie nie stoi. */
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
  };
}

/** Węzeł znacznika jako element, gdy taki stoi; pustka znaczy węzeł, którego znacznik nie niesie. */
function wskazElement(wezel: Element | null): HTMLElement | null {
  return wezel instanceof HTMLElement ? wezel : null;
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
  wpiszTrescWpisu(wpis, wiadomosc.content);
}

/** Wpisuje treść w klon wpisu. Pasek działań należy do biblioteki i przeżywa podmianę treści; znika samo zdanie przykładowe, nie przyciski, które znacznik niesie. */
function wpiszTrescWpisu(wpis: HTMLElement, tekst: string): void {
  const tresc = wpis.querySelector('.sta-wpis-tresc');
  if (tresc === null) return;
  const akcje = tresc.querySelector('.sta-wpis-akcje');
  tresc.textContent = tekst;
  if (akcje !== null) tresc.appendChild(akcje);
}

/** Treść wpisu bez paska działań — od niej zaczyna dopisywanie strumień, który zastał wpis już stojący. */
function trescWpisu(wpis: HTMLElement): string {
  const tresc = wpis.querySelector('.sta-wpis-tresc');
  if (tresc === null) return '';
  const akcje = tresc.querySelector('.sta-wpis-akcje');
  return [...tresc.childNodes]
    .filter((wezel) => wezel !== akcje)
    .map((wezel) => wezel.textContent ?? '')
    .join('');
}

/** Godzina wpisu w zapisie, którego używa znacznik historii — godziny i minuty. */
function godzinaWpisu(znacznik: number): string {
  return new Date(znacznik).toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' });
}

/* Panele karty wiąże się dopiero po założeniu okna komunikacji: każdy z nich
   pyta rdzeń o treść tego okna, a przed jego powstaniem nie ma o co pytać.
   Węzły idą od korzenia karty, a odłączenia nasłuchów wchodzą do wiązania
   karty i schodzą razem z nim. Panel, którego znacznik nie stoi, nie robi nic. */
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
    zwiazNarzedzia(kanal, idOkna, korzen),
  ];
  for (const odlacz of zwiazane) {
    if (odlacz !== null) odlaczenia.push(odlacz);
  }
  zwiazPodglad(kanal, idOkna, korzen);
}
