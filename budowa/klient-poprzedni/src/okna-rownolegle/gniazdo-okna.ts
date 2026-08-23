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

/** Jedno miejsce na scenie: nagłówek roli i okno komunikacji pod nim. */
export interface GniazdoOkna {
  id: IdGniazda;
  /** Element gniazda montowany w torze układu. */
  element: HTMLElement;
  /** Nagłówek gniazda — rola, kierunek zlecenia, stan pętli. */
  naglowek: NaglowekGniazda;
  /**
   * Stała fasada okna komunikacji.
   *
   * Przepływ komunikatów podpina się raz i przeżywa przebudowę widoku przy
   * zmianie roli — nie musi znać życia okna pod spodem.
   */
  fasada: OknoKomunikacji;
  /** Bieżąca rola okna. */
  rola(): WindowRole;
  /** Nadaje rolę: przestawia nagłówek i przebudowuje widok okna. */
  ustawRole(rola: WindowRole): void;
  /**
   * Zgłasza zmianę roli — także tę przychodzącą z rdzenia.
   *
   * Bez tego zgłoszenia układ nie wiedziałby o zmianie nadanej gniazdu z boku
   * (zdarzenie `window.changed`), a więź koordynator–wykonawca i pas relacji
   * zostawałyby przy poprzedniej figurze.
   */
  naZmianeRoli(sluchacz: (rola: WindowRole) => void): Odsubskrybuj;
  /**
   * Zgłasza moduł nadany gniazdu — także ten przychodzący z rdzenia.
   *
   * Moduł rozstrzyga o figurze rozmowy sceny (ile okien i w jakich rolach —
   * `figura-modulu.ts`), a przychodzi zdarzeniem `window.changed` na fasadę
   * gniazda. Bez tego zgłoszenia układ pytałby o figurę raz, przy montażu,
   * kiedy moduł okna jeszcze nie jest znany.
   */
  naZmianeModulu(sluchacz: (kod: string) => void): Odsubskrybuj;
  /** Włącza albo wyłącza gniazdo ze sceny. */
  pokaz(widoczne: boolean): void;
  /** Zgłasza, czy gniazdo jest obecnie na scenie. */
  widoczne(): boolean;
  /** Ustawia stan pętli pokazywany w nagłówku. */
  ustawStanPary(stan: StanPary | null): void;
  /**
   * Podaje nagłówkowi odczyt łączności z rdzeniem.
   *
   * Idzie do nagłówka, nie przez fasadę okna: fasadowe `pokazStan` bez
   * wbudowanej rozmowy jest cichym no-opem (`okno === null`), więc łączność
   * poprowadzona tamtędy nic by nie pokazała.
   */
  ustawLacznosc(odczyt: OdczytLacznosci): void;
  /**
   * Podaje oknu wykonania gniazda kod nadany przez rdzeń.
   *
   * Panele stoją na komendach rdzenia, a te żądają `windowId` okna otwartego.
   * Kod przychodzi dopiero po uzgodnieniu, więc gniazdo powstaje bez niego
   * i dostaje go później — pusty napis do tej chwili nie jest awarią, tylko
   * stanem, który panel nazywa wprost.
   */
  ustawOknoWykonania(kod: string): void;
  /**
   * Podaje sekcji czynności port trybów widoku transkryptu — albo go zabiera.
   *
   * Wejście dla powłoki, nie dla układu. Tryby zapisu należą do rozmowy
   * (`rozmowa/widok-zapisu.ts`), a rozmowę osadza w gnieździe dopiero
   * `aplikacja/wiazanie-gniazda.ts` — po uzgodnieniu z rdzeniem, czyli po
   * powstaniu gniazda. Gniazdo nie zna ani jednego trybu: przenosi port dalej,
   * do sekcji czynności.
   */
  ustawWidokZapisu(port: PortWidokuZapisu | null): void;
  /**
   * Przerysowuje sekcję czynności sesji w menu `⋮`.
   *
   * Woła to układ po każdej zmianie sceny: pozycja `Otwórz w nowym oknie` żyje
   * dopóty, dopóki scena ma wolne gniazdo, a sufit przestawia się razem
   * z figurą modułu przychodzącą z rdzenia.
   */
  odswiezCzynnosci(): void;
  /**
   * Zamyka panele gniazda wraz z ich subskrypcjami rdzenia.
   *
   * Obowiązkowe przy zdejmowaniu sceny. Panel trzyma subskrypcję strumienia
   * (`stream.chunk`), która przeżywa usunięcie węzła z drzewa dokumentu —
   * gniazdo zdjęte bez tego wywołania zostawia ją żywą.
   */
  zamknij(): void;
}

/** Ustawienia gniazda; jedyne ma wartość domyślną. */
export interface OpcjeGniazda {
  /**
   * Czy gniazdo buduje własny widok rozmowy (`okno-komunikacji/utworzOkno`)
   * pod swoim nagłówkiem.
   *
   * Domyślnie wyłączone. Scena sesji (`aplikacja/scena-sesji.ts`) osadza
   * właściwy widok rozmowy z zewnątrz, przez `aplikacja/wiazanie-gniazda.ts`;
   * drugi widok tego samego miejsca byłby podwójną implementacją. Stanowisko
   * podglądu układu (`podglad-ukladu.html`) nie ma ani rdzenia, ani takiego
   * wiązania, więc włącza to ustawienie jawnie: bez wbudowanego widoku jego
   * gniazda stałyby puste.
   */
  wbudowanaRozmowa?: boolean;
  /**
   * Kanał do rdzenia, którym panele pomocnicze wołają swoje komendy.
   *
   * Pominięty znaczy, że gniazdo nie ma czym otworzyć ani jednego panelu —
   * menu `⋮` jest wtedy puste i mówi to wprost. Stanowisko podglądu układu
   * pracuje bez rdzenia i tak właśnie ma wyglądać.
   */
  kanal?: Kanal;
  /**
   * Kod okna wykonania nadany przez rdzeń, jeśli już jest znany.
   *
   * Zwykle nie jest — okno powstaje po uzgodnieniu, więc gniazdo dostaje kod
   * później, przez `ustawOknoWykonania`.
   */
  okno?: string;
  /**
   * Dojście do liczby gniazd sceny — pod pozycję `Otwórz w nowym oknie`
   * w menu `⋮`.
   *
   * Podaje je układ, bo tylko on wie, ile gniazd scena ma i ile jeszcze
   * zniesie. Gniazdo tego nie liczy: przenosi port dalej, do sekcji czynności.
   * Pominięty znaczy „nie ma sceny, na którą dałoby się dostawić okno" —
   * gniazdo stojące samotnie na stanowisku podglądu pozycji nie pokazuje.
   */
  noweOkno?: PortNowegoOkna;
  /**
   * Dojście do przebiegu ponowienia dla plakietki łączności w nagłówku.
   *
   * Pominięte znaczy „transport nie wystawia numeru próby ani czasu do
   * następnej" — plakietka nazywa wtedy ten brak zamiast liczyć próby po raz
   * drugi u siebie. Patrz `zrodlo-lacznosci.ts → portPonawiania`.
   */
  ponowienie?: PortPonawiania;
}

/**
 * Gniazdo układu okien równoległych.
 *
 * Jedna odpowiedzialność: utrzymanie jednego miejsca na scenie — jego opisu
 * i nagłówka roli — wraz z opcjonalnym własnym oknem komunikacji
 * (patrz `OpcjeGniazda.wbudowanaRozmowa`). Gniazdo nie zna ani pary, ani
 * liczby okien na scenie — te należą do układu.
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
  // Kod katalogu okien aplikacji pada tutaj, bo tutaj powstaje okno rozmowy.
  // Ramy okna operacyjnego (`komponenty/rama-okna.ts`) to okno nie ma — niesie
  // ją gniazdo układu równoległego, a druga rama byłaby drugim nagłówkiem nad
  // jednym oknem.
  //
  // Nazwa idzie z `opis.tytul`, nie ze stałej: gniazd jest cztery („Okno
  // komunikacji 1"…„4"), a jedna nazwa na wszystkie odebrałaby czytnikowi
  // ekranu rozróżnienie widoczne wzrokiem.
  oznaczOknoAplikacji({ element, kod: 'chat-window', nazwa: opis.tytul });
  element.append(naglowek.element);

  naglowek.ustawRole(rolaPoczatkowa);
  // Moduł znany już w chwili montażu (z opisu okna) trafia do nagłówka od razu;
  // późniejsze przestawienia idą przez `fasada.ustawModul`, czyli tą samą
  // drogą, którą przychodzi `window.changed` z rdzenia.
  naglowek.ustawModul(opis.modul);

  /**
   * Sekcja czynności sesji doklejana w menu `⋮` za kreską.
   *
   * Identyfikator sesji przychodzi kanałem, nie osobnym parametrem: warstwa
   * rozmowy pracuje na `idOkna`, a czynności sesji żądają `idSesji`, który
   * gniazdo ma przez `kanal.sesja()` — to samo źródło, z którego
   * `aplikacja/scena-sesji.ts` liczy `sesjaGotowa()`, a
   * `aplikacja/otwarcie-okna.ts` bierze `sessionId` dla `window.create`. Sekcja
   * czyta identyfikator przy każdym pytaniu, bo gniazdo powstaje przed
   * uzgodnieniem i sesji wtedy jeszcze nie ma.
   *
   * Bez kanału sekcji nie ma: stanowisko podglądu układu pracuje bez rdzenia,
   * a `menu-paneli.ts` stawia kreskę wyłącznie przy niepustym `sekcjeDalsze`.
   */
  const czynnosciSesji: SekcjaCzynnosciSesji | null =
    opcje.kanal === undefined
      ? null
      : utworzSekcjeCzynnosciSesji({
          kanal: opcje.kanal,
          ...(opcje.noweOkno === undefined ? {} : { noweOkno: opcje.noweOkno }),
        });

  /**
   * Panele tej rozmowy — sterowanie w nagłówku, stos i uchwyt obok rozmowy.
   *
   * Przedrostek `dnp` należy do obszaru okien pomocniczych, nie do modułu:
   * panel stoi obok rozmowy, a nie w złożeniu Developera czy Diagnostics, więc
   * klasy modułu nie miałyby arkusza, który je pokrywa.
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
   * Podział szerokości gniazda na kolumnę rozmowy i kolumnę paneli.
   *
   * Liczony w TypeScripcie, nie w arkuszu, bo kontrakt szerokości (600 px) musi
   * być tą samą liczbą, o którą opiera się uchwyt. Dwa progi — jeden w CSS,
   * drugi w kodzie — rozjechałyby się przy pierwszej poprawce.
   */
  function przeliczKolumny(): void {
    const podzial = ustalPodzial(panele.szerokosc(), element.clientWidth);
    element.style.gridTemplateColumns = kolumnyGniazda(podzial);
  }

  panele.naZmiane(przeliczKolumny);
  przeliczKolumny();

  // Dziennik, okno i jego subskrypcja istnieją wyłącznie z wbudowaną rozmową
  // — bez niej nie ma czyjej historii pilnować ani czego przebudowywać po
  // zmianie roli (patrz `odbuduj` niżej).
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
   * Przebudowa widoku okna po zmianie roli.
   *
   * Rola widnieje w nagłówku okna komunikacji, więc po jej zmianie widok musi
   * powstać na nowo — inaczej gniazdo pokazywałoby jedną rolę, a okno pod nim
   * drugą. Historia i stan łączności wracają z dziennika, a subskrypcja
   * treści wpisanej przechodzi na nowy widok.
   *
   * Bez wbudowanej rozmowy nie ma czego przebudować: rolę widzi operator już
   * w nagłówku gniazda, który `ustawRole` przestawia niezależnie od tej
   * funkcji — funkcja wtedy jawnie nic nie robi.
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

  // Bez wbudowanej rozmowy fasada jest jawnie nieczynna dla treści i stanu:
  // operator widzi rozmowę w widoku osadzonym z zewnątrz, a stan łączności
  // płynie do niego inną drogą (transport → widok rozmowy), więc drugi zapis
  // tej samej historii tutaj nie miałby odbiorcy.
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

    // Moduł trafia także do opisu gniazda, bo przebudowa widoku po zmianie roli
    // stawia okno na nowo z opisu. Bez tego okno wracałoby do modułu sprzed
    // przestawienia, mimo że rdzeń przestawił je na inny.
    //
    // Odczyt idzie z opisu, nie z widoku: na scenie sesji gniazdo nie buduje
    // wbudowanego okna (rozmowę osadza wiązanie), więc widok bywa nieobecny,
    // a moduł gniazda musi być znany zawsze.
    ustawModul(kod) {
      if (kod === opis.modul) return;
      opis = { ...opis, modul: kod };
      okno?.ustawModul(kod);
      // Nagłówek gniazda niesie nazwę modułu jako wartość bieżącą okna, więc
      // przestawia się razem z nią — inaczej pokazywałby moduł sprzed zmiany.
      naglowek.ustawModul(kod);
      // Spis okien pomocniczych idzie po module, więc menu `⋮` musi przestawić
      // się razem z nim — inaczej pokazywałoby panele modułu poprzedniego.
      panele.ustawModul(kod);
      zmianyModulu.oglos(kod);
    },
    modul: () => opis.modul,

    // Dyktowanie przechodzi wprost do okna wbudowanego i nie jest tu dublowane.
    // Gniazdo bez wbudowanego okna nie ma warstwy mowy z tego samego powodu, dla
    // którego nie ma treści ani stanu: rozmowę osadza wtedy wiązanie z zewnątrz
    // i to ona niesie mikrofon.
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
      // Sekcja czynności trzyma subskrypcję `session.changed` — przeżywa
      // usunięcie węzła z dokumentu tak samo jak subskrypcje paneli.
      czynnosciSesji?.zamknij();
      // Nagłówek trzyma zegar błysku i odliczanie plakietki łączności; oba
      // przeżywają usunięcie węzła z dokumentu.
      naglowek.zamknij();
    },
  };
}
