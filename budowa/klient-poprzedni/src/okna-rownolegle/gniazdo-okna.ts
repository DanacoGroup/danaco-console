import './gniazdo.css';

import type { WindowRole } from '../../../shared/contract';
import { oznaczOknoAplikacji } from '../komponenty/okno-aplikacji';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { StanPolaczenia } from '../polaczenie/stan-polaczenia';
import { utworzOkno, type OknoKomunikacji } from '../okno-komunikacji/okno';
import type { OpisOkna } from '../okno-komunikacji/opis-okna';
import type { Kanal } from '../protokol/kanal';
import {
  utworzSekcjeCzynnosciSesji,
  type SekcjaCzynnosciSesji,
} from './czynnosci-sesji-menu';
import { utworzDziennikWpisow, type DziennikWpisow } from './dziennik-wpisow';
import type { IdGniazda } from './identyfikatory';
import type { OdczytLacznosci, PortPonawiania } from './lacznosc-okna';
import { utworzNaglowekGniazda, type NaglowekGniazda } from './naglowek-gniazda';
import { opisGniazda } from './opis-gniazda';
import { utworzPaneleGniazda, type PaneleGniazda } from './panele-gniazda';
import type { PortWidokuZapisu } from './podmenu-widoku-zapisu';
import type { PortNowegoOkna } from './pozycja-nowego-okna';
import { kolumnyGniazda, ustalPodzial } from './szerokosci-gniazda';
import type { StanPary } from './stan-pary';

/**
 * Jedno miejsce na scenie okien równoległych: nagłówek roli wraz z oknem komunikacji tuż pod tym nagłówkiem.
 */
export interface GniazdoOkna {
  id: IdGniazda;
  /** Element gniazda montowany w torze układu. */
  element: HTMLElement;
  /** Nagłówek gniazda — rola, kierunek zlecenia, stan pętli. */
  naglowek: NaglowekGniazda;
  /**
   * Stała fasada okna komunikacji; przepływ komunikatów podpina się raz i przeżywa przebudowę.
   */
  fasada: OknoKomunikacji;
  /** Bieżąca rola okna. */
  rola(): WindowRole;
  /** Nadaje rolę: przestawia nagłówek i przebudowuje widok okna. */
  ustawRole(rola: WindowRole): void;
  /**
   * Zgłasza zmianę roli — także tę przychodzącą z rdzenia zdarzeniem window.changed.
   */
  naZmianeRoli(sluchacz: (rola: WindowRole) => void): Odsubskrybuj;
  /**
   * Zgłasza moduł nadany gniazdu, także ten przychodzący z rdzenia zdarzeniem window.changed.
   */
  naZmianeModulu(sluchacz: (kod: string) => void): Odsubskrybuj;
  /** Włącza albo wyłącza gniazdo ze sceny. */
  pokaz(widoczne: boolean): void;
  /** Zgłasza, czy gniazdo jest obecnie na scenie. */
  widoczne(): boolean;
  /** Ustawia stan pętli pokazywany w nagłówku. */
  ustawStanPary(stan: StanPary | null): void;
  /**
   * Podaje nagłówkowi odczyt łączności z rdzeniem; fasada bez rozmowy jest cichym no-opem.
   */
  ustawLacznosc(odczyt: OdczytLacznosci): void;
  /**
   * Podaje oknu wykonania gniazda kod nadany przez rdzeń, dostępny dopiero po uzgodnieniu.
   */
  ustawOknoWykonania(kod: string): void;
  /**
   * Podaje sekcji czynności port trybów widoku transkryptu; gniazdo samo nie zna żadnego trybu.
   */
  ustawWidokZapisu(port: PortWidokuZapisu | null): void;
  /**
   * Przerysowuje sekcję czynności sesji w menu bocznym po każdej zmianie sceny.
   */
  odswiezCzynnosci(): void;
  /**
   * Zamyka panele gniazda wraz z ich subskrypcjami rdzenia; obowiązkowe przy zdejmowaniu sceny.
   */
  zamknij(): void;
}

/**
 * Ustawienia gniazda układu okien równoległych; jedyne pole tych ustawień ma wartość domyślną gniazda.
 */
export interface OpcjeGniazda {
  /**
   * Czy gniazdo buduje własny widok rozmowy pod swoim nagłówkiem; domyślnie wyłączone.
   */
  wbudowanaRozmowa?: boolean;
  /**
   * Kanał do rdzenia, którym panele pomocnicze wołają swoje komendy; pominięty ogranicza gniazdo.
   */
  kanal?: Kanal;
  /**
   * Kod okna wykonania nadany przez rdzeń, jeśli już jest znany w chwili powstania gniazda.
   */
  okno?: string;
  /**
   * Dojście do liczby gniazd sceny, potrzebne do pozycji otwarcia nowego okna w menu bocznym.
   */
  noweOkno?: PortNowegoOkna;
  /**
   * Dojście do przebiegu ponowienia dla plakietki łączności widocznej w nagłówku gniazda.
   */
  ponowienie?: PortPonawiania;
}

/**
 * Gniazdo układu okien równoległych utrzymuje jedno miejsce na scenie z opcjonalnym oknem komunikacji.
 */
export function utworzGniazdo(
  id: IdGniazda,
  podstawa: OpisOkna,
  rolaPoczatkowa: WindowRole,
  opcje: OpcjeGniazda = {},
): GniazdoOkna {
  const wbudowanaRozmowa = opcje.wbudowanaRozmowa ?? false;

  let opis = opisGniazda(podstawa, id, rolaPoczatkowa);

  const wpisane = utworzMagistrale<string>();
  const zmianyRoli = utworzMagistrale<WindowRole>();
  const zmianyModulu = utworzMagistrale<string>();
  const naglowek = utworzNaglowekGniazda(
    id,
    rolaPoczatkowa,
    opcje.ponowienie === undefined ? {} : { ponowienie: opcje.ponowienie },
  );

  const gospodarz = document.createElement('div');
  gospodarz.className = 'dn-okna__okno';

  const element = document.createElement('section');
  element.className = 'dn-okna__gniazdo';
  element.dataset.gniazdo = id;
  // Kod katalogu okien aplikacji, bo tutaj powstaje okno rozmowy tego gniazda.
  oznaczOknoAplikacji({ element, kod: 'chat-window', nazwa: opis.tytul });
  element.append(naglowek.element);

  naglowek.ustawRole(rolaPoczatkowa);
  // Moduł znany już w chwili montażu trafia do nagłówka od razu, dalsze zmiany idą inną drogą.
  naglowek.ustawModul(opis.modul);

  /**
   * Sekcja czynności sesji doklejana w menu bocznym; identyfikator sesji przychodzi kanałem gniazda.
   */
  const czynnosciSesji: SekcjaCzynnosciSesji | null =
    opcje.kanal === undefined
      ? null
      : utworzSekcjeCzynnosciSesji({
          kanal: opcje.kanal,
          ...(opcje.noweOkno === undefined ? {} : { noweOkno: opcje.noweOkno }),
        });

  /**
   * Panele tej rozmowy: sterowanie w nagłówku, stos i uchwyt obok rozmowy tego gniazda.
   */
  const panele: PaneleGniazda = utworzPaneleGniazda({
    kanal: opcje.kanal ?? null,
    modul: opis.modul,
    okno: opcje.okno ?? '',
    przedrostek: 'dnp',
    dostepnaSzerokosc: () => element.clientWidth,
    ...(czynnosciSesji === null ? {} : { sekcjeDalsze: [czynnosciSesji.element] }),
  });

  naglowek.sterowanie.append(panele.sterowanie);
  element.append(panele.uchwyt, panele.kolumna, panele.pelnyEkran);

  /**
   * Podział szerokości gniazda na kolumnę rozmowy i kolumnę paneli, liczony w kodzie, nie w arkuszu.
   */
  function przeliczKolumny(): void {
    const podzial = ustalPodzial(panele.szerokosc(), element.clientWidth);
    element.style.gridTemplateColumns = kolumnyGniazda(podzial);
  }

  panele.naZmiane(przeliczKolumny);
  przeliczKolumny();

  // Dziennik, okno i jego subskrypcja istnieją wyłącznie z wbudowaną rozmową tego gniazda.
  const dziennik: DziennikWpisow | null = wbudowanaRozmowa ? utworzDziennikWpisow() : null;
  let okno: OknoKomunikacji | null = null;
  let odsubskrybuj: Odsubskrybuj | null = null;
  let ostatniStan: { stan: StanPolaczenia; oczekujace: number } | null = null;

  if (wbudowanaRozmowa) {
    element.append(gospodarz);
    okno = utworzOkno(opis);
    odsubskrybuj = okno.naWpisanie((tresc) => wpisane.oglos(tresc));
    gospodarz.append(okno.element);
  }

  /**
   * Przebudowa widoku okna po zmianie roli, bo rola widnieje w nagłówku okna komunikacji.
   */
  function odbuduj(): void {
    if (okno === null || odsubskrybuj === null || dziennik === null) return;
    odsubskrybuj();
    const nowe = utworzOkno(opis);
    dziennik.odtworz(nowe);
    if (ostatniStan !== null) nowe.pokazStan(ostatniStan.stan, ostatniStan.oczekujace);
    gospodarz.replaceChildren(nowe.element);
    okno = nowe;
    odsubskrybuj = okno.naWpisanie((tresc) => wpisane.oglos(tresc));
  }

  // Bez wbudowanej rozmowy fasada jest jawnie nieczynna dla treści i stanu łączności.
  const fasada: OknoKomunikacji = {
    element: gospodarz,
    dopisz(wpis) {
      if (okno === null || dziennik === null) return;
      dziennik.zapisz(wpis);
      okno.dopisz(wpis);
    },
    dopiszFragment(persona, fragment) {
      if (okno === null || dziennik === null) return;
      dziennik.zapiszFragment(persona, fragment);
      okno.dopiszFragment(persona, fragment);
    },
    pokazStan(stan, oczekujace) {
      if (okno === null) return;
      ostatniStan = { stan, oczekujace };
      okno.pokazStan(stan, oczekujace);
    },
    naWpisanie: (sluchacz) => wpisane.subskrybuj(sluchacz),
    ustawOgnisko: () => okno?.ustawOgnisko(),

    // Moduł trafia także do opisu gniazda, bo przebudowa widoku stawia okno na nowo z opisu.
    ustawModul(kod) {
      if (kod === opis.modul) return;
      opis = { ...opis, modul: kod };
      okno?.ustawModul(kod);
      // Nagłówek gniazda niesie nazwę modułu jako wartość bieżącą okna, przestawianą razem z nim.
      naglowek.ustawModul(kod);
      // Spis okien pomocniczych idzie po module, więc menu boczne przestawia się razem z nim.
      panele.ustawModul(kod);
      zmianyModulu.oglos(kod);
    },
    modul: () => opis.modul,

    // Dyktowanie przechodzi wprost do okna wbudowanego i nie jest tu dublowane.
    dyktowanie: () => okno?.dyktowanie() ?? null,
    podlaczDyktowanie: (warstwa) => okno?.podlaczDyktowanie(warstwa),
  };

  return {
    id,
    element,
    naglowek,
    fasada,
    rola: () => opis.rola,

    ustawRole(rola) {
      if (rola === opis.rola) return;
      opis = { ...opis, rola };
      naglowek.ustawRole(rola);
      odbuduj();
      zmianyRoli.oglos(rola);
    },

    naZmianeRoli: (sluchacz) => zmianyRoli.subskrybuj(sluchacz),
    naZmianeModulu: (sluchacz) => zmianyModulu.subskrybuj(sluchacz),

    pokaz(widoczne) {
      element.hidden = !widoczne;
    },

    widoczne: () => !element.hidden,
    ustawStanPary: (stan) => naglowek.ustawStan(stan),
    ustawLacznosc: (odczyt) => naglowek.ustawLacznosc(odczyt),
    ustawOknoWykonania: (kod) => panele.ustawOkno(kod),
    ustawWidokZapisu: (port) => czynnosciSesji?.ustawWidokZapisu(port),
    odswiezCzynnosci: () => czynnosciSesji?.odswiez(),

    zamknij() {
      panele.zamknij();
      // Sekcja czynności trzyma subskrypcję zmiany sesji, przeżywającą usunięcie węzła z dokumentu.
      czynnosciSesji?.zamknij();
      // Nagłówek trzyma zegar błysku i licznik plakietki, oba przeżywają usunięcie węzła.
      naglowek.zamknij();
    },
  };
}
