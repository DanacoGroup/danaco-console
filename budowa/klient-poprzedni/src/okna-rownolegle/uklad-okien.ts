import './uklad.css';

import type { WindowRole } from '../../../shared/contract';
import type { OknoKomunikacji } from '../okno-komunikacji/okno';
import { opisPoczatkowy, type OpisOkna } from '../okno-komunikacji/opis-okna';
import type { Transport } from '../polaczenie/gniazdo';
import type { Kanal } from '../protokol/kanal';
import { zdanieNadmiaruSceny, zdanieNadwyzkiPozycji } from './etykiety-ukladu';
import { figuraModulu, type FiguraModulu } from './figura-modulu';
import { utworzGniazdo, type GniazdoOkna } from './gniazdo-okna';
import {
  czyIdGniazda,
  ID_GNIAZD,
  LICZBA_MAX,
  LICZBA_MIN,
  ograniczLiczbe,
  type IdGniazda,
} from './identyfikatory';
import type { OdczytLacznosci } from './lacznosc-okna';
import { utworzPasRelacji } from './pas-relacji';
import { pokazPrzekazanieZlecenia } from './przekazanie-zlecenia';
import { utworzPrzelacznikLiczby } from './przelacznik-liczby';
import { rolaDomyslna } from './role-domyslne';
import type { StanPary } from './stan-pary';
import { opisKierunku, ustalWiez, type Wiez } from './wiez-koordynacji';
import { portPonawiania, zwiazLacznoscUkladu } from './zrodlo-lacznosci';

/** Układ od jednego do trzech okien komunikacji obok siebie. */
export interface UkladOkien {
  /** Element montowany na scenie powłoki. */
  element: HTMLElement;
  /** Ustawia liczbę okien na scenie. */
  ustawLiczbe(liczba: number): void;
  /** Bieżąca liczba okien na scenie. */
  liczba(): number;
  /**
   * Przestawia scenę na moduł o podanym kodzie: liczbę okien, jaką ten moduł
   * prowadzi, i zdanie o niej dla Operatora.
   *
   * Sama okien nie zdejmuje — patrz `zdanieNadmiaruSceny`.
   */
  ustawModulSceny(kod: string): void;
  /** Figura rozmowy modułu, w którym scena pracuje w tej chwili. */
  figura(): FiguraModulu;
  /** Nadaje oknu rolę w pętli koordynator–wykonawca. */
  nadajRole(idOkna: string, rola: WindowRole): void;
  /** Pokazuje przekazanie zlecenia między dwoma oknami — w obu naraz. */
  pokazPrzekazanie(od: string, do_: string): void;
  /** Ustawia stan pętli pokazywany w nagłówkach pary i na pasie relacji. */
  ustawStanPary(stan: StanPary): void;
  /**
   * Rozsyła odczyt łączności do wszystkich gniazd sceny.
   *
   * Jedna prawda o nastawie: łącze jest jedno na całego klienta, więc gniazdo
   * pierwsze i czwarte muszą po jednej zmianie transportu pokazać dokładnie to
   * samo. Rozesłanie idzie stąd, a nie z każdego gniazda osobno, właśnie po to.
   *
   * Woła to samo wiązanie transportu (`OpcjeUkladu.transport`); wejście zostaje
   * publiczne, bo scena podglądu układu pracuje bez rdzenia i podaje odczyty
   * ręcznie.
   */
  ustawLacznosc(odczyt: OdczytLacznosci): void;
  /** Gniazdo po identyfikatorze; `null` dla identyfikatora spoza układu. */
  gniazdo(idOkna: string): GniazdoOkna | null;
  /** Fasada okna gniazda — punkt podpięcia przepływu komunikatów. */
  fasadaOkna(idOkna: string): OknoKomunikacji | null;
  /**
   * Podaje gniazdu kod okna wykonania nadany przez rdzeń.
   *
   * Panele pomocnicze wołają komendy żądające `windowId` okna otwartego, a ten
   * kod powstaje dopiero po uzgodnieniu — układ montuje się wcześniej. Dopóki
   * kodu nie ma, panel mówi wprost, czego mu brakuje, zamiast pokazywać pustkę.
   */
  ustawOknoWykonania(idGniazda: string, kod: string): void;
  /**
   * Zamyka panele wszystkich gniazd wraz z ich subskrypcjami rdzenia.
   *
   * Subskrypcja strumienia przeżywa usunięcie węzła z drzewa dokumentu —
   * scena zdjęta bez tego wywołania zostawia po sobie żywe nasłuchy.
   */
  zamknij(): void;
}

/** Ustawienia początkowe układu; każde ma wartość domyślną. */
export interface OpcjeUkladu {
  /** Opis wspólny, z którego powstają opisy okien gniazd. */
  podstawa?: OpisOkna;
  /**
   * Liczba okien na starcie.
   *
   * Zostaje sprowadzona do figury modułu sceny: moduł prowadzący dwa okna nie
   * wstanie z czterema tylko dlatego, że ktoś je tu wpisał.
   */
  liczbaPoczatkowa?: number;
  /** Role nadane z góry; pozostałe gniazda biorą rolę domyślną. */
  rolePoczatkowe?: Partial<Record<IdGniazda, WindowRole>>;
  /**
   * Czy każde gniazdo buduje własny widok rozmowy — patrz
   * `GniazdoOkna` / `OpcjeGniazda.wbudowanaRozmowa`.
   *
   * Domyślnie wyłączone (scena sesji osadza właściwy widok z zewnątrz);
   * stanowisko podglądu układu włącza je jawnie, bo nie ma z zewnątrz niczego
   * do osadzenia.
   */
  wbudowanaRozmowa?: boolean;
  /**
   * Kanał do rdzenia przekazywany gniazdom dla ich paneli pomocniczych.
   *
   * Pominięty znaczy „nie ma czym otworzyć ani jednego panelu" — menu paneli
   * jest wtedy puste i mówi to wprost. Stanowisko podglądu układu pracuje bez
   * rdzenia i właśnie tak ma wyglądać (brak nazwany, nie ukryty).
   */
  kanal?: Kanal;
  /**
   * Transport, z którego scena bierze stan łącza i długość kolejki.
   *
   * Pominięty znaczy „nie ma czego pokazywać" — plakietki łączności milczą,
   * bo scena nie zna łącza. Tak pracuje stanowisko podglądu układu.
   */
  transport?: Transport;
}

/**
 * Układ okien równoległych.
 *
 * Jedna odpowiedzialność: kompozycja sceny — przełącznik liczby, tor gniazd
 * i pas relacji — oraz utrzymanie tego, co między oknami wspólne: liczby,
 * ról i więzi koordynator–wykonawca. Samo okno komunikacji nie powstaje tutaj;
 * układ składa gotowe okna z `okno-komunikacji/`.
 */
export function utworzUkladOkien(opcje: OpcjeUkladu = {}): UkladOkien {
  const podstawa = opcje.podstawa ?? opisPoczatkowy();

  /**
   * Figura rozmowy modułu sceny — ile okien moduł prowadzi i w jakich rolach.
   *
   * Moduł bierze się z opisu okna, a przy każdej odpowiedzi rdzenia
   * (`window.changed`) przestawia go gniazdo pierwsze — patrz subskrypcja niżej.
   */
  let figuraSceny = figuraModulu(podstawa.modul);
  let liczbaOkien = ograniczLiczbe(
    Math.min(opcje.liczbaPoczatkowa ?? LICZBA_MIN, pojemnoscModulu()),
  );
  let stanBiezacy: StanPary = 'brak-pary';

  /**
   * Ile okien scena otworzy w module bieżącym.
   *
   * Moduł bez rozmowy (`liczbaOkien === 0`) nie zostawia sceny pustej: scena
   * trzyma `LICZBA_MIN`, a zdanie przełącznika mówi wprost, że rozmowy w tym
   * module nie ma. Pusta scena czytałaby się jak awaria, a to nie jest awaria.
   */
  function pojemnoscModulu(): number {
    return Math.max(figuraSceny.liczbaOkien, LICZBA_MIN);
  }

  /**
   * Czy dostawienie gniazda cokolwiek zmieni.
   *
   * Pytanie idzie tą samą rachubą, która potem wykona: odpowiedź to wprost
   * wynik `ustawLiczbe` policzony na sucho, bo sufit składa się z figury modułu
   * i z `LICZBA_MAX`, a figura przestawia się z rdzenia (`window.changed`).
   * Osobne „liczba < 4" rozjechałoby się z tamtym przycięciem przy pierwszym
   * module węższym niż scena.
   *
   * Odpowiedź „nie" znaczy krótszą listę w menu `⋮`, a nie wiersz wygaszony —
   * patrz `pozycja-nowego-okna.ts`.
   */
  function wolneGniazdo(): boolean {
    return ograniczLiczbe(Math.min(liczbaOkien + 1, pojemnoscModulu())) > liczbaOkien;
  }

  /** Gniazda o roli nadanej z zewnątrz — zmiana liczby okien jej nie odbiera. */
  const reczne = new Set<IdGniazda>();

  /**
   * Znacznik przestawiania ról przez sam układ. Gniazdo zgłasza tylko to, że
   * rola się zmieniła, a źródło rozstrzyga o trwałości: rola z zewnątrz (operator
   * albo rdzeń) jest decyzją i zostaje, rola z `ustawLiczbe` jest domyślną
   * i ustępuje następnej. Nadawanie ról jest synchroniczne, więc znacznik nie
   * przecieka poza swoją pętlę.
   */
  let ukladPrzestawia = false;

  const element = document.createElement('section');
  element.className = 'dn-okna';
  element.setAttribute('aria-label', 'Okna komunikacji sesji');

  const listwa = document.createElement('div');
  listwa.className = 'dn-listwa dn-okna__listwa';

  const tor = document.createElement('div');
  tor.className = 'dn-okna__tor';

  const pas = utworzPasRelacji();

  const przelacznik = utworzPrzelacznikLiczby(liczbaOkien, (wybor) => ustawLiczbe(wybor));
  listwa.append(przelacznik.element);

  // Przebieg ponowienia bierze się z transportu, jeśli transport go wystawia.
  // Gdy nie wystawia, port jest `null`, a plakietka nazywa ten brak zamiast
  // go zasłaniać.
  const ponowienie =
    opcje.transport === undefined ? null : portPonawiania(opcje.transport);

  const gniazda = new Map<IdGniazda, GniazdoOkna>();
  for (const id of ID_GNIAZD) {
    const zadana = opcje.rolePoczatkowe?.[id];
    if (zadana !== undefined) reczne.add(id);
    const gniazdo = utworzGniazdo(id, podstawa, zadana ?? rolaDomyslna(id, liczbaOkien), {
      wbudowanaRozmowa: opcje.wbudowanaRozmowa,
      ...(opcje.kanal === undefined ? {} : { kanal: opcje.kanal }),
      ...(ponowienie === null ? {} : { ponowienie }),
      // Pozycja `Otwórz w nowym oknie` idzie tą samą drogą co przełącznik
      // liczby. Podniesienie liczby wprowadza gniazdo na scenę, a wejście
      // gniazda na scenę zamawia dla niego okno rdzenia
      // (`aplikacja/scena-sesji.ts`, obserwator `hidden`). Okno założone
      // `window.create` z boku byłoby niewidoczne. Tu stoją same wywołania:
      // układ nie wie, że po drugiej stronie jest wiersz menu.
      noweOkno: { wolneGniazdo, naNoweOkno: () => ustawLiczbe(liczbaOkien + 1) },
    });
    gniazda.set(id, gniazdo);
    tor.append(gniazdo.element);
    // Każda zmiana roli — czyjakolwiek — przelicza figurę pętli na całej scenie.
    gniazdo.naZmianeRoli(() => {
      if (!ukladPrzestawia) reczne.add(id);
      odswiez();
    });
  }

  // Moduł sceny idzie z gniazda pierwszego. To ono niesie okno uzgodnione
  // z rdzeniem (`montaz-ukladu.ts`), więc to na nie przychodzi `window.changed`
  // razem z `moduleId` (`okno-komunikacji/przeplyw-komunikatow.ts`) — także po
  // `workspace.enter`, które przestawia moduł okna zamiast zakładać drugie.
  // Scena z oknami w różnych modułach bierze figurę z gniazda pierwszego.
  gniazda.get(ID_GNIAZD[0])?.naZmianeModulu((kod) => ustawModulSceny(kod));

  element.append(listwa, tor, pas.element);

  /** Gniazda obecne na scenie, w kolejności od lewej. */
  function widoczne(): IdGniazda[] {
    return ID_GNIAZD.slice(0, liczbaOkien);
  }

  /** Mapa ról wszystkich gniazd — podstawa ustalenia więzi. */
  function role(): Map<IdGniazda, WindowRole> {
    return new Map([...gniazda].map(([id, gniazdo]) => [id, gniazdo.rola()]));
  }

  /** Rozesłanie stanu pętli do nagłówków pary i na pas relacji. */
  function rozeslijStan(wiez: Wiez | null): void {
    for (const [id, gniazdo] of gniazda) {
      const wPetli = wiez !== null && (id === wiez.koordynator || id === wiez.wykonawca);
      gniazdo.ustawStanPary(wPetli ? stanBiezacy : null);
    }
    pas.ustawStan(stanBiezacy);
  }

  /** Przeliczenie całej sceny po każdej zmianie liczby albo roli. */
  function odswiez(): void {
    element.dataset.liczba = String(liczbaOkien);
    tor.dataset.liczba = String(liczbaOkien);

    const naScenie = widoczne();
    for (const [id, gniazdo] of gniazda) {
      gniazdo.pokaz(naScenie.includes(id));
      // Menu `⋮` każdego gniazda niesie pozycję żyjącą z wolnego gniazda na
      // scenie. Zmiana liczby albo figury modułu przesuwa sufit, więc sekcja
      // czynności musi przeliczyć się razem ze sceną — inaczej wiersz przeżyłby
      // własne pokrycie.
      gniazdo.odswiezCzynnosci();
    }

    const wiez = ustalWiez(role(), naScenie);
    if (wiez === null) stanBiezacy = 'brak-pary';
    else if (stanBiezacy === 'brak-pary') stanBiezacy = 'gotowa';

    for (const [id, gniazdo] of gniazda) {
      gniazdo.naglowek.ustawKierunek(opisKierunku(wiez, id));
    }

    pas.ustawWiez(wiez, liczbaOkien);
    rozeslijStan(wiez);
    przelacznik.ustawFigure(pojemnoscModulu(), zdanieSceny());
  }

  /**
   * Co scena ma dziś do powiedzenia o liczbie okien swojego modułu.
   *
   * Trzy zdania, każde tylko wtedy, gdy jest prawdziwe: czym jest figura
   * modułu; że pozycja wyższa okna nie dokłada; że scena trzyma okien więcej,
   * niż moduł prowadzi. Ostatnie bierze się stąd, że układ okien nie zdejmuje
   * sam — zdjęcie okna zamyka je również w rdzeniu i jest decyzją Operatora.
   */
  function zdanieSceny(): string {
    const zdania = [figuraSceny.zdanie];
    const pojemnosc = pojemnoscModulu();

    if (figuraSceny.kod.length > 0 && pojemnosc < LICZBA_MAX) {
      zdania.push(zdanieNadwyzkiPozycji(pojemnosc));
    }
    if (liczbaOkien > pojemnosc) {
      zdania.push(zdanieNadmiaruSceny(liczbaOkien, pojemnosc));
    }
    return zdania.join(' ');
  }

  /**
   * Zmiana liczby okien. Gniazda o roli nadanej z zewnątrz zachowują ją;
   * pozostałe biorą rolę domyślną właściwą dla nowej liczby okien — przejście
   * z jednego okna na dwa od razu pokazuje figurę koordynator–wykonawca,
   * a decyzja operatora ani rdzenia nigdy nie zostaje cofnięta.
   */
  function ustawLiczbe(liczba: number): void {
    // Żądanie ponad figurę modułu nie jest odmową: pozycja przełącznika zostaje
    // czynna, scena bierze tyle okien, ile moduł prowadzi, a zdanie pod grupą
    // mówi, dlaczego wyżej się nie da. Cicha zgoda na czwarte okno w module
    // dwuokiennym byłaby sukcesem udawanym.
    liczbaOkien = ograniczLiczbe(Math.min(liczba, pojemnoscModulu()));
    ukladPrzestawia = true;
    for (const [id, gniazdo] of gniazda) {
      if (reczne.has(id)) continue;
      gniazdo.ustawRole(rolaDomyslna(id, liczbaOkien));
    }
    ukladPrzestawia = false;
    przelacznik.pokaz(liczbaOkien);
    odswiez();
  }

  /**
   * Przestawienie sceny na moduł.
   *
   * Okien nie zdejmuje, nawet gdy jest ich więcej, niż moduł prowadzi. Zdjęcie
   * okna ze sceny zamyka je również w rdzeniu (`aplikacja/scena-sesji.ts` →
   * `zamknijOstatnie`); zrobione samoczynnie przy przestawieniu modułu byłoby
   * zamknięciem okna, o które nikt nie prosił. Scena mówi o nadmiarze wprost
   * i zostawia zdjęcie Operatorowi.
   */
  function ustawModulSceny(kod: string): void {
    if (kod === figuraSceny.kod) return;
    figuraSceny = figuraModulu(kod);
    odswiez();
  }

  /** Gniazdo po identyfikatorze; identyfikator spoza układu daje `null`. */
  function gniazdo(idOkna: string): GniazdoOkna | null {
    return czyIdGniazda(idOkna) ? (gniazda.get(idOkna) ?? null) : null;
  }

  /** Przestawienie stanu pętli i rozesłanie go po scenie — jedna droga. */
  function przestawStan(nowy: StanPary): void {
    stanBiezacy = nowy;
    rozeslijStan(ustalWiez(role(), widoczne()));
  }

  /**
   * Rozesłanie odczytu łącza do wszystkich gniazd — także tych poza sceną.
   *
   * Gniazdo zdjęte ze sceny przełącznikiem liczby wraca na nią bez ponownego
   * montażu, więc pominięcie go tutaj dałoby okno pokazujące stan łącza sprzed
   * schowania. Rozesłanie do czterech elementów jest tańsze niż druga ścieżka
   * doganiania stanu przy pokazaniu gniazda.
   */
  function ustawLacznosc(odczyt: OdczytLacznosci): void {
    for (const [, miejsce] of gniazda) miejsce.ustawLacznosc(odczyt);
  }

  // Wiązanie transportu żyje tak długo jak scena. Subskrypcja przeżywa
  // usunięcie węzła z dokumentu, więc odłączenie idzie w `zamknij()` razem
  // z panelami gniazd.
  const odlaczLacznosc =
    opcje.transport === undefined
      ? null
      : zwiazLacznoscUkladu(opcje.transport, ustawLacznosc);

  odswiez();

  return {
    element,
    ustawLiczbe,
    liczba: () => liczbaOkien,
    ustawModulSceny,
    figura: () => figuraSceny,

    // Znacznik trwałości i przeliczenie sceny bierze na siebie zgłoszenie
    // gniazda — tą samą drogą wchodzi rola nadana przez rdzeń.
    nadajRole: (idOkna, rola) => gniazdo(idOkna)?.ustawRole(rola),

    pokazPrzekazanie(od, do_) {
      const zrodlo = gniazdo(od);
      const cel = gniazdo(do_);
      if (zrodlo === null || cel === null || zrodlo === cel) return;
      pokazPrzekazanieZlecenia({ od: zrodlo, do_: cel, pas, ustawStan: przestawStan });
    },

    ustawStanPary: przestawStan,
    ustawLacznosc,

    gniazdo,
    fasadaOkna: (idOkna) => gniazdo(idOkna)?.fasada ?? null,

    ustawOknoWykonania: (idGniazda, kod) => gniazdo(idGniazda)?.ustawOknoWykonania(kod),

    zamknij() {
      odlaczLacznosc?.();
      for (const [, miejsce] of gniazda) miejsce.zamknij();
    },
  };
}
