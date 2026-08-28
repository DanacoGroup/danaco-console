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

/** Układ okien komunikacji pokazuje od jednego do trzech okien obok siebie na wspólnej scenie równoległej. */
export interface UkladOkien {
  /** Element montowany na scenie powłoki. */
  element: HTMLElement;
  /** Ustawia liczbę okien na scenie. */
  ustawLiczbe(liczba: number): void;
  /** Bieżąca liczba okien na scenie. */
  liczba(): number;
  // Przestawia scenę na moduł o podanym kodzie, licząc jego liczbę okien i zdanie o niej dla operatora.
  ustawModulSceny(kod: string): void;
  /** Figura rozmowy modułu, w którym scena pracuje w tej chwili. */
  figura(): FiguraModulu;
  /** Nadaje oknu rolę w pętli koordynator–wykonawca. */
  nadajRole(idOkna: string, rola: WindowRole): void;
  /** Pokazuje przekazanie zlecenia między dwoma oknami — w obu naraz. */
  pokazPrzekazanie(od: string, do_: string): void;
  /** Ustawia stan pętli pokazywany w nagłówkach pary i na pasie relacji. */
  ustawStanPary(stan: StanPary): void;
  // Rozsyła odczyt łączności do wszystkich gniazd sceny, bo łącze jest jedno na całego klienta.
  ustawLacznosc(odczyt: OdczytLacznosci): void;
  /** Gniazdo po identyfikatorze; `null` dla identyfikatora spoza układu. */
  gniazdo(idOkna: string): GniazdoOkna | null;
  /** Fasada okna gniazda — punkt podpięcia przepływu komunikatów. */
  fasadaOkna(idOkna: string): OknoKomunikacji | null;
  // Podaje gniazdu kod okna wykonania nadany przez rdzeń dopiero po uzgodnieniu, po zmontowaniu układu.
  ustawOknoWykonania(idGniazda: string, kod: string): void;
  // Zamyka panele wszystkich gniazd wraz z ich subskrypcjami rdzenia, żeby nie zostały żywe nasłuchy.
  zamknij(): void;
}

/** Ustawienia początkowe układu okien mają wartość domyślną, więc żadne pole nie jest wymagane do podania. */
export interface OpcjeUkladu {
  /** Opis wspólny, z którego powstają opisy okien gniazd. */
  podstawa?: OpisOkna;
  // Liczba okien na starcie zostaje sprowadzona do figury modułu sceny, nie przyjmowana wprost.
  liczbaPoczatkowa?: number;
  /** Role nadane z góry; pozostałe gniazda biorą rolę domyślną. */
  rolePoczatkowe?: Partial<Record<IdGniazda, WindowRole>>;
  // Czy każde gniazdo buduje własny widok rozmowy; domyślnie wyłączone, bo scenę osadza się z zewnątrz.
  wbudowanaRozmowa?: boolean;
  // Kanał do rdzenia przekazywany gniazdom dla paneli pomocniczych; pominięty czyni menu paneli pustym.
  kanal?: Kanal;
  // Transport, z którego scena bierze stan łącza i długość kolejki; pominięty wycisza plakietki łącza.
  transport?: Transport;
}

/** Układ okien równoległych ma jedną odpowiedzialność: kompozycję sceny wraz z liczbą, rolami i więzią koordynator–wykonawca. */
export function utworzUkladOkien(opcje: OpcjeUkladu = {}): UkladOkien {
  const podstawa = opcje.podstawa ?? opisPoczatkowy();

  // Figura rozmowy modułu sceny mówi, ile okien moduł prowadzi i w jakich rolach; nadaje ją gniazdo.
  let figuraSceny = figuraModulu(podstawa.modul);
  let liczbaOkien = ograniczLiczbe(
    Math.min(opcje.liczbaPoczatkowa ?? LICZBA_MIN, pojemnoscModulu()),
  );
  let stanBiezacy: StanPary = 'brak-pary';

  // Ile okien scena otworzy w module bieżącym; moduł bez rozmowy nie zostawia sceny pustej.
  function pojemnoscModulu(): number {
    return Math.max(figuraSceny.liczbaOkien, LICZBA_MIN);
  }

  // Czy dostawienie gniazda cokolwiek zmieni, liczone tą samą rachubą, którą wykona ustawienie liczby.
  function wolneGniazdo(): boolean {
    return ograniczLiczbe(Math.min(liczbaOkien + 1, pojemnoscModulu())) > liczbaOkien;
  }

  /** Gniazda o roli nadanej z zewnątrz — zmiana liczby okien jej nie odbiera. */
  const reczne = new Set<IdGniazda>();

  // Znacznik przestawiania ról przez sam układ odróżnia rolę domyślną od roli nadanej z zewnątrz.
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

  // Przebieg ponowienia bierze się z transportu; brak transportu daje port pusty, nazwany wprost.
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
      // Pozycja otwarcia nowego okna idzie tą samą drogą co przełącznik liczby, nie zakłada okna z boku.
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

  // Moduł sceny idzie z gniazda pierwszego, bo to ono niesie okno uzgodnione z rdzeniem.
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
      // Menu gniazda niesie pozycję zależną od wolnego miejsca na scenie, więc przelicza się razem z nią.
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

  // Co scena ma dziś do powiedzenia o liczbie okien modułu, w trzech zdaniach, tylko gdy prawdziwe.
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

  // Zmiana liczby okien zachowuje role nadane z zewnątrz, a pozostałym gniazdom nadaje rolę domyślną.
  function ustawLiczbe(liczba: number): void {
    // Żądanie ponad figurę modułu nie jest odmową: scena bierze tyle okien, ile moduł prowadzi.
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

  // Przestawienie sceny na moduł nie zdejmuje okien samo, nawet gdy jest ich więcej, niż moduł prowadzi.
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

  // Rozesłanie odczytu łącza idzie do wszystkich gniazd, także tych chwilowo zdjętych ze sceny.
  function ustawLacznosc(odczyt: OdczytLacznosci): void {
    for (const [, miejsce] of gniazda) miejsce.ustawLacznosc(odczyt);
  }

  // Wiązanie transportu żyje tak długo jak scena, więc odłączenie idzie razem z zamknięciem paneli.
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

    // Znacznik trwałości i przeliczenie sceny biorą na siebie zgłoszenie gniazda o zmianie roli.
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
