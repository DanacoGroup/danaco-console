/**
 * Wiązanie okna Studia z rdzeniem. Znacznik niesie biblioteka Właściciela —
 * ten plik nic nie buduje: słucha zdarzeń, woła komendy kontraktu i wypełnia
 * stojące węzły, a wpisy i pozycje powiela z wzorów zdjętych ze znacznika.
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
  type Message,
  type Session,
  type StreamChunkEvent,
  type StudioDocument,
} from '../../../shared/contract.ts';
import type { Kanal } from '../protokol/kanal.ts';
import { wywolaj } from '../protokol/wywolanie.ts';
import { oglos } from './ogloszenie.ts';
import { wpiszPole } from './okno-modulu.ts';
import { zapewnijSesje } from './sesja-biezaca.ts';
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

/**
 * Węzeł historii, na którym stoi wiązanie bieżące. Powłoka wstawia wnętrze okna
 * na nowo przy każdym wejściu w moduł, a wiązanie ma objąć znacznik świeży:
 * znacznik ten sam znaczy wiązanie już założone, znacznik inny — okno wstawione
 * ponownie i wiązanie do założenia od nowa.
 */
let zwiazanaHistoria: Element | null = null;

/**
 * Wiąże okno Studia z rdzeniem. Kanał domyślnie z obiektu globalnego, bo
 * skrypty biblioteki są funkcjami domkniętymi i importu z nich nie ma.
 * Identyfikator okna podaje wołający wznawiający sesję: takie okno stoi już
 * w rejestrze rdzenia, a pustka znaczy okno zakładane pierwszą wiadomością.
 * Prawda znaczy, że znacznik okna stał i wiązanie zostało założone.
 */
export function zwiazStudio(
  kanal: Kanal | undefined = globalThis.DanacoKanal,
  nazwaSrodowiska = '',
  idOknaStojacego = '',
): boolean {
  if (kanal === undefined) return false;
  const znalezione = zbierzWezly();
  if (znalezione === null) return false;
  if (znalezione.historia === zwiazanaHistoria) return false;
  zwiazanaHistoria = znalezione.historia;
  const wezly: WezlyStudia = znalezione;

  const wzory = zdejmijWzoryWpisow(wezly.historia);
  const wzorPozycji = zdejmijWzorPozycji(wezly.szyna);
  zdejmijTrescPrzykladowa(wezly);
  zdejmijCzynnosciWstazkiBezPokrycia();

  const wpisy = new Map<string, HTMLElement>();
  const strumienie = new Map<string, StrumienOdpowiedzi>();
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

  /* Sesja powstaje dopiero pierwszą wiadomością wysłaną do modelu — nie
     otwarciem okna ani wejściem w moduł. Do tej chwili Operator chodzi po
     produkcie swobodnie i żadna sesja po nim nie zostaje. Wyjątkiem jest okno
     wskazane przez wołającego: ono stoi już w rejestrze i wchodzi przy montażu. */
  const zapewnijStanowisko = async (): Promise<string> => {
    if (idOkna !== '') return idOkna;
    idOkna = await otworzStanowisko(kanal, wezly, wzorPozycji, wstawWpis, idOknaStojacego);
    if (idOkna === '') return '';
    zwiazPanele(kanal, idOkna);
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
    void porownajWersjeDokumentu().then((porownane) => {
      if (porownane) return;
      oglos('Studio', 'Okno nie prowadzi dokumentu — nie ma czego porównać.');
    });
  });

  wezly.kanwa.addEventListener('input', () => {
    odswiezPasStanu(wezly, dokument);
  });

  /* Nowy dokument jest pracą z modelem w edytorze, więc zakłada stanowisko
     tak samo jak pierwsza wiadomość. */
  wezly.nowyDokument?.addEventListener('click', () => {
    void zapewnijStanowisko().then((okno) => {
      if (okno === '') return;
      void zalozDokument(kanal, okno, wezly).then((zalozony) => {
        dokument = zalozony;
        opiszDokument(zalozony);
      });
    });
  });

  kanal.naZdarzenie(EventType.MessageChanged, (tresc) => {
    if (tresc.message.windowId !== idOkna) return;
    zmienWpis(tresc.change, tresc.message);
  });

  /* Odpowiedź modelu przyrasta w oknie fragmentami: pełną wiadomość rdzeń
     rozgłasza dopiero po turze. Numer fragmentu i znacznik domknięcia stoją
     w kopercie, nie w treści zdarzenia. */
  kanal.naZdarzenie(EventType.StreamChunk, (tresc, koperta) => {
    if (tresc.windowId !== idOkna) return;
    przyjmijFragment(tresc, koperta.seq ?? 0, koperta.done === true);
  });

  /* Stan pusty dokumentu wchodzi od razu: nazwa pracy z prototypu jest treścią
     przykładową, a okno staje, zanim powstanie sesja i dokument. */
  opiszDokument(null);
  void opiszKanal(kanal, nazwaSrodowiska);
  void opiszWyborModelu(kanal);
  opiszWyborNakladu();
  /* Okno wskazane przez wołającego stoi już w rdzeniu, więc historia rozmowy
     wraca przy montażu, bez czekania na pierwszą wiadomość. */
  if (idOknaStojacego !== '') void zapewnijStanowisko();
  return true;
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
function zdejmijTrescPrzykladowa(wezly: WezlyStudia): void {
  wezly.historia.replaceChildren();
  wezly.szyna.replaceChildren();
  wezly.kanwa.replaceChildren();
  zdejmijZnacznikiBezZrodla();
  zdejmijTrescPrzykladowaPlanu();
  zdejmijTrescPrzykladowaRoznic();
  zdejmijTrescPrzykladowaRepozytorium();
  zdejmijTrescPrzykladowaPlikow();
  zdejmijTrescPrzykladowaNarzedzi();
  zdejmijTrescPrzykladowaPodgladu();
}

/**
 * Zdejmuje ze wstążki czynności, których wywołać nie ma czym. Podgląd wydruku
 * wymaga formatu wydania i miejsca na strony, wywołanie operacji — nazwy
 * operacji i jej zakresu, przekazanie do Biblioteki — profilu wydania; wstążka
 * nie ma węzła, w którym Operator poda którąkolwiek z tych wartości.
 */
function zdejmijCzynnosciWstazkiBezPokrycia(): void {
  document.querySelector('.st-wstazka [data-etykietka="Podgląd wydruku"]')?.remove();
  document.querySelector('.st-wstazka .st-wstazka-grupa--drugorzedna')?.remove();
}

/**
 * Zdejmuje znaczniki bez pokrycia w kontrakcie: nazwę gałęzi, ścieżkę
 * repozytorium i miarę różnicy w pasie czynności, wskazania źródła nad polem
 * wpisu oraz pas czytelności, dla którego rdzeń nie ma ani jednej miary.
 */
function zdejmijZnacznikiBezZrodla(): void {
  for (const znacznik of document.querySelectorAll('.sta-kontekst-akcji .sta-chip')) {
    if (znacznik.classList.contains('sta-chip--srodowisko')) continue;
    znacznik.remove();
  }
  for (const zrodlo of document.querySelectorAll('.sta-zrodlo')) zrodlo.remove();
  document.querySelector('.sta-kom-monitor')?.remove();
}

/** Wpisuje w wybór modelu kanały rejestru rdzenia; wiersz wzorcowy powiela się na każdy kanał, a kanał czynny nazywa sam znacznik wyboru. */
async function opiszWyborModelu(kanal: Kanal): Promise<void> {
  const znak = document.querySelector('.sta-chip--model');
  const spis = document.getElementById('pop-model');
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
function opiszWyborNakladu(): void {
  const spis = document.getElementById('pop-wysilek');
  const znak = document.querySelector('[data-popover="pop-wysilek"]');
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

/** Okno komunikacji wraz z sesją, w której stoi; szyna dokumentów wskazuje po niej sesję bieżącą. */
interface Stanowisko {
  idOkna: string;
  idSesji: string;
}

/**
 * Wskazuje stanowisko okna, wczytuje jego historię i szynę sesji; zwraca
 * identyfikator okna albo pustkę, gdy rdzeń odmówił. Okno podane przez
 * wołającego stoi już w rejestrze, więc `window.create` nie pada — pada
 * wyłącznie dla okna zakładanego pierwszą wiadomością.
 */
async function otworzStanowisko(
  kanal: Kanal,
  wezly: WezlyStudia,
  wzorPozycji: HTMLElement | null,
  wstawWpis: (wiadomosc: Message) => void,
  idOknaStojacego: string,
): Promise<string> {
  const stanowisko = idOknaStojacego === ''
    ? await zalozStanowisko(kanal)
    : await wskazStanowisko(kanal, idOknaStojacego);
  if (stanowisko === null) return '';

  const historia = await wywolaj(kanal, Command.MessageList, { windowId: stanowisko.idOkna });
  if (historia.udany && historia.wynik !== undefined) {
    for (const wiadomosc of historia.wynik.messages) wstawWpis(wiadomosc);
  }

  await wypelnijSzyne(kanal, wezly.szyna, wzorPozycji, stanowisko.idSesji);
  return stanowisko.idOkna;
}

/** Zakłada sesję i okno komunikacji; pustka znaczy odmowę rdzenia na którymkolwiek z trzech kroków. */
async function zalozStanowisko(kanal: Kanal): Promise<Stanowisko | null> {
  // Stanowisko staje w karcie sesji bieżącej, tak samo jak każde inne okno modułu.
  const idSesji = await zapewnijSesje(kanal, 'Studio');
  if (idSesji === '') return null;

  const modul = await wskazModulStudia(kanal);
  const kanalModelu = await wskazKanalModelu(kanal);
  if (modul === '' || kanalModelu === '') return null;

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
    zatrzymaj: wskazElement(kom?.querySelector('.sta-prompt-zatrzymaj') ?? null),
    szyna,
    nowyDokument: nowyDokument instanceof HTMLElement ? nowyDokument : null,
    kanwa,
    status,
    zapiszWersje: wskazElement(document.querySelector('.st-wstazka [data-etykietka="Zapisz wersję"]')),
    porownajWersje: wskazElement(
      document.querySelector('.st-wstazka [data-etykietka="Porównaj wersje"]'),
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
