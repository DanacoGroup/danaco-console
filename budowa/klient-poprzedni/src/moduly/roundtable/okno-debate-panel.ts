import {
  RoundtableFormat,
  type RoundtableDebateStartRequest,
  type RoundtableParticipant,
  type RoundtableTurn,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pobierzPlik,
  pole,
  poleLiczbowe,
  poleTresci,
  przyciskAkcji as przycisk,
  wiersz,
  wybor,
} from '../../modele/kontrolki-formularza';
import { rysujDebate, transkryptJson, transkryptMarkdown } from './przebieg-debaty';
import type { StanDebaty } from './stan-debaty';
import { utworzStanTresci } from './stany-okna';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';
import { czynnosciModeracji } from './czynnosci-arsenalu';
import { utworzWykazFunkcji, utworzZestawAkcji } from './warstwy-modulu';
import type { ZrodloArsenaluRoundtable } from './zrodlo-arsenalu';
import type { ZrodloRoundtable } from './zrodlo-roundtable';

/**
 * Debate Panel — okno monitora modułu Roundtable: uruchamia kolejną turę debaty i pokazuje
 * narastające argumenty jednocześnie od wszystkich uczestników.
 */
export interface OknoDebatePanel {
  element: HTMLElement;
  odswiez(): void;
  /** Przerysowanie po fragmencie strumienia — wołane przez złożenie modułu, kilkanaście razy na sekundę. */
  odswiezGlosy(): void;
  /** Zamyka nasłuch `stan.naZmiane(...)` założony przez to okno. */
  zamknij(): void;
}

/**
 * @param przyciskiRozszerzen przyciski otwierające okna warstwy drugiej, wytworzone przez pas
 *   rozszerzeń złożenia modułu; Debate Panel dostaje je gotowe, mechanizmu nie zna.
 */
export function utworzOknoDebatePanel(
  zrodlo: ZrodloRoundtable,
  arsenal: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  przyciskiRozszerzen: readonly HTMLElement[],
): OknoDebatePanel {
  const rama = utworzRameOkna({
    tytul: 'Debate Panel',
    rola: 'monitor',
    kod: 'debate-panel',
    przeznaczenie:
      'Uruchomienie kolejnej tury debaty i przegląd narastających wypowiedzi uczestników.',
    przedrostek: 'dr',
  });
  const tresc = utworzStanTresci();
  // Czynności nad zapisem debaty wołają komendy obszaru wprost; format bierze się z pól okna.
  const powierzchnia = zlozPowierzchnieDebaty(
    rama,
    tresc.element,
    przyciskiRozszerzen,
    czynnosciModeracji(
      arsenal,
      stan,
      (zdanie, powodzenie) => tresc.potwierdzenie(zdanie, powodzenie),
      () => 'Szablon z Debate Panelu',
      () => 'markdown',
    ),
  );
  let filtrZapisu = '';

  function rysuj(): void {
    rysujDebate(stan, strumien, tresc, filtrZapisu);
  }

  /** Stany, których przerysowanie po fragmencie nie ma prawa zdjąć z ekranu. */
  function stanNietykalny(): boolean {
    const rodzaj = tresc.rodzaj();
    return rodzaj === 'blad' || rodzaj === 'ladowanie';
  }

  function odswiezGlosy(): void {
    if (stanNietykalny()) return;
    rysuj();
  }

  function uruchomTure(): void {
    const idOkna = stan.okno();
    if (idOkna === '') {
      tresc.blad('Nie wiadomo, które okno debaty uruchamiać — brak identyfikatora okna.');
      return;
    }
    const pytanie = powierzchnia.pytanie.value.trim();
    if (pytanie === '') {
      // Puste pytanie nie jedzie do rdzenia — to walidacja kliencka, nie odmowa serwera.
      tresc.blad('Pytanie jest wymagane. Puste pytanie nie jedzie do rdzenia.');
      return;
    }
    const zadanie = zlozZadanieDebaty(idOkna, pytanie, powierzchnia);
    tresc.ladowanie('Pytanie wysyłane jednocześnie do wszystkich uczestników…');
    void zrodlo.uruchomDebate(zadanie).then((wynik) => {
      if (!wynik.udany || wynik.wynik === undefined) {
        // Nie rdzeń odmówił: niepowodzenie bywa odpowiedzią bez pola obowiązkowego, którego nie orzekł.
        tresc.blad('Tura debaty nie została uruchomiona.', wynik.blad);
        return;
      }
      stan.ustawTure(wynik.wynik.turn);
      const potwierdzenie = opisUczestnictwa(
        wynik.wynik.turn,
        wynik.wynik.participantIds,
        stan.uczestnicy(),
      );
      tresc.potwierdzenie(potwierdzenie.zdanie, potwierdzenie.udane);
    });
  }

  /** Czy okno ma cokolwiek do wydania — wspólne dla obu postaci transkryptu. */
  function jestCoEksportowac(): boolean {
    return stan.wypowiedzi().length > 0 || strumien.nieprzypisane().length > 0;
  }

  function eksportujTranskrypt(): void {
    if (!jestCoEksportowac()) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — okno nie widziało jeszcze wypowiedzi.', false);
      return;
    }
    pobierzPlik(
      `debata-${stan.tura() === '' ? 'bez-tury' : stan.tura()}-transkrypt.md`,
      transkryptMarkdown(stan, strumien),
      'text/markdown',
    );
    tresc.potwierdzenie(
      'Transkrypt pobrany. Obejmuje wyłącznie wypowiedzi widziane przez to okno od jego otwarcia — historii sprzed wejścia okno jeszcze nie wczytuje.',
      true,
    );
  }

  function eksportujTranskryptJson(): void {
    if (!jestCoEksportowac()) {
      tresc.potwierdzenie('Nie ma czego wyeksportować — okno nie widziało jeszcze wypowiedzi.', false);
      return;
    }
    pobierzPlik(
      `debata-${stan.tura() === '' ? 'bez-tury' : stan.tura()}-transkrypt.json`,
      transkryptJson(stan, strumien),
      'application/json',
    );
    tresc.potwierdzenie(
      'Transkrypt pobrany w postaci JSON. Zakres jest ten sam co w Markdown — wypowiedzi widziane przez to okno od jego otwarcia.',
      true,
    );
  }

  podepnijAkcjeDebaty(powierzchnia, {
    uruchomTure,
    eksportujTranskrypt,
    eksportujTranskryptJson,
    filtruj: (wartosc) => {
      filtrZapisu = wartosc;
      rysuj();
    },
  });

  const odsubskrybuj = stan.naZmiane(() => rysuj());
  rysuj();

  return { element: rama.element, odswiez: rysuj, odswiezGlosy, zamknij: odsubskrybuj };
}

/** Zadanie uruchomienia debaty złożone z kontrolek okna — pola puste nie trafiają do żądania wysyłanego rdzeniowi. */
function zlozZadanieDebaty(
  idOkna: string,
  pytanie: string,
  powierzchnia: PowierzchniaDebaty,
): RoundtableDebateStartRequest {
  const zadanie: RoundtableDebateStartRequest = { windowId: idOkna, question: pytanie };
  const temat = powierzchnia.temat.value.trim();
  if (temat !== '') zadanie.topic = temat;
  if (powierzchnia.format.value !== '') zadanie.format = powierzchnia.format.value as RoundtableFormat;
  const limit = Number.parseInt(powierzchnia.limitTur.value, 10);
  if (Number.isInteger(limit) && limit > 0) zadanie.turnLimit = limit;
  return zadanie;
}

/**
 * Zdanie o tym, którą turę rdzeń otworzył i kto w niej faktycznie odpowiada, licząc różnicę
 * wobec pełnego składu debaty.
 */
function opisUczestnictwa(
  tura: RoundtableTurn,
  participantIds: readonly string[],
  sklad: readonly RoundtableParticipant[],
): { zdanie: string; udane: boolean } {
  const naglowek = `Rdzeń otworzył turę #${tura.index} (${tura.status}).`;
  if (participantIds.length === 0) {
    return {
      zdanie: `${naglowek} Rdzeń nie wskazał ANI JEDNEGO uczestnika odpowiadającego w tej turze.`,
      udane: false,
    };
  }
  if (sklad.length === 0) {
    return {
      zdanie: `${naglowek} Odpowiada ${participantIds.length} uczestników (składu debaty to okno jeszcze nie zna).`,
      udane: true,
    };
  }
  if (participantIds.length >= sklad.length) {
    return {
      zdanie: `${naglowek} Pytanie poszło do całego składu (${participantIds.length} z ${sklad.length}).`,
      udane: true,
    };
  }
  return {
    zdanie:
      `${naglowek} Pytanie poszło TYLKO do ${participantIds.length} z ${sklad.length} uczestników — ` +
      'reszta składu nie odpowiada w tej turze.',
    udane: true,
  };
}

interface AkcjeDebaty {
  uruchom: HTMLButtonElement;
  eksport: HTMLButtonElement;
  eksportJson: HTMLButtonElement;
}

/**
 * Pasek akcji Debate Panel: uruchomienie tury, dwie postaci transkryptu i przyciski
 * otwierające okna warstwy drugiej.
 */
function zlozAkcjeDebaty(
  gospodarz: HTMLElement,
  przyciskiRozszerzen: readonly HTMLElement[],
): AkcjeDebaty {
  const uruchom = przycisk('Uruchom kolejną turę', 'dn-btn dn-btn--atrament');
  const eksport = przycisk('Eksportuj transkrypt (Markdown)');
  const eksportJson = przycisk('Eksportuj transkrypt (JSON)');
  gospodarz.append(uruchom, eksport, eksportJson, ...przyciskiRozszerzen);
  return { uruchom, eksport, eksportJson };
}

/**
 * Zestaw akcji warstwy trzeciej — czynności nad zapisem debaty, meldujące wynik w stanie
 * treści okna, tą samą drogą co reszta odpowiedzi.
 */
function zlozZestawDebaty(czynnosci: HTMLButtonElement[]): HTMLElement {
  return utworzZestawAkcji('Zestaw akcji i filtrów zapisu debaty', czynnosci);
}

interface PowierzchniaDebaty extends AkcjeDebaty {
  pytanie: HTMLTextAreaElement;
  temat: HTMLInputElement;
  format: HTMLSelectElement;
  limitTur: HTMLInputElement;
  filtr: HTMLInputElement;
}

/**
 * Formaty debaty — komplet wyliczenia kontraktu, nie jego podzbiór, żeby pole wyboru nie
 * odbierało Operatorowi formatów, których rdzeń nie odmawia.
 */
const OPISY_FORMATU: ReadonlyArray<readonly [string, string]> = [
  ['', 'Format domyślny rdzenia'],
  [RoundtableFormat.Free, 'Debata swobodna'],
  [RoundtableFormat.Structured, 'Debata strukturalna'],
  [RoundtableFormat.Oxford, 'Debata oksfordzka'],
  [RoundtableFormat.RoundRobin, 'Tura okrężna'],
  [RoundtableFormat.Delphi, 'Rundy Delphi'],
  [RoundtableFormat.ExpertPanel, 'Panel ekspercki'],
];

/** Składa kontrolki okna, pasek akcji i ciało ramy — czysta konstrukcja bez podpięcia zdarzeń, które robi funkcja osobna. */
function zlozPowierzchnieDebaty(
  rama: { akcje: HTMLElement; cialo: HTMLElement },
  stanTresci: HTMLElement,
  przyciskiRozszerzen: readonly HTMLElement[],
  czynnosci: HTMLButtonElement[],
): PowierzchniaDebaty {
  const pytanie = poleTresci('Pytanie do wszystkich uczestników', 3, 'Pytanie kierowane jednocześnie do całego składu…');
  const temat = pole('Zagadnienie tury', 'opcjonalne');
  const format = wybor('Format debaty', OPISY_FORMATU);
  const limitTur = poleLiczbowe('Liczba tur', 'domyślnie bez limitu');
  const filtr = pole('Filtruj zapis', 'mówca albo fragment treści');
  const akcje = zlozAkcjeDebaty(rama.akcje, przyciskiRozszerzen);

  rama.cialo.append(
    wiersz('Pytanie', pytanie, {
      klasa: 'dr-wiersz',
      objasnienie: 'Wymagane — puste pytanie nie jedzie do rdzenia.',
    }),
    wiersz('Zagadnienie', temat, { klasa: 'dr-wiersz' }),
    wiersz('Format', format, { klasa: 'dr-wiersz' }),
    wiersz('Limit tur', limitTur, { klasa: 'dr-wiersz' }),
    wiersz('Filtr', filtr, {
      klasa: 'dr-wiersz',
      objasnienie: 'Filtrowanie działa wyłącznie na wypowiedziach, które to okno już widziało.',
    }),
    zlozZestawDebaty(czynnosci),
    stanTresci,
    utworzWykazFunkcji('debate-panel'),
  );
  return { pytanie, temat, format, limitTur, filtr, ...akcje };
}

function podepnijAkcjeDebaty(
  powierzchnia: PowierzchniaDebaty,
  obsluga: {
    uruchomTure: () => void;
    eksportujTranskrypt: () => void;
    eksportujTranskryptJson: () => void;
    filtruj: (wartosc: string) => void;
  },
): void {
  powierzchnia.uruchom.addEventListener('click', obsluga.uruchomTure);
  powierzchnia.eksport.addEventListener('click', obsluga.eksportujTranskrypt);
  powierzchnia.eksportJson.addEventListener('click', obsluga.eksportujTranskryptJson);
  powierzchnia.filtr.addEventListener('input', () => {
    obsluga.filtruj(powierzchnia.filtr.value.trim().toLowerCase());
  });
}

