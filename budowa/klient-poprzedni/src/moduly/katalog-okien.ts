import { Command, type ModuleListResponse } from '../../../shared/contract';
import { EnvelopeStatus } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/**
 * Katalog okien — jedno źródło prawdy o tym, których okien operacyjnych
 * swojego katalogu moduł nie zbudował.
 *
 * Liczba wpisana w pasek uczciwości na stałe rozjeżdża się z rdzeniem po cichu:
 * katalog okien zmienia migracja rdzenia, a napis nie jest z nim niczym
 * połączony. Ten byt liczy rozjazd z odczytu, więc pasek nie orzeka o czymś,
 * czego rdzeń nie powiedział.
 *
 * Byt stoi w korzeniu `moduly/`, a nie w module, bo tę samą potrzebę ma każdy
 * moduł kontraktu; kopia na moduł dałaby tyle samo rozjeżdżających się zdań
 * o jednym stanie produktu. Bliźniakiem jest `moduly/pokrycie-komend.ts`,
 * rozstrzygający to samo o komendach.
 *
 * Prawda pochodzi z `module.list`: rdzeń oddaje w
 * `Module.operationalWindowCodes` kody okien operacyjnych otwieranych wraz
 * z modułem, odpowiednik tabeli `okno_operacyjne`. Czego moduł buduje, rdzeń
 * nie wie i wiedzieć nie może — tę połowę podaje moduł.
 *
 * Okno rozmowy nie jest oknem operacyjnym modułu: katalog rdzenia niesie
 * `chat-window` w każdym module, a montuje je scena sesji, nie rejestr
 * modułów. Wykaz, który tego nie odejmuje, przypisuje każdemu modułowi brak
 * okna, które działa.
 *
 * Rozstrzygnięcia są cztery, bo każde znaczy co innego:
 *   `nieustalone`      rdzeń jeszcze nie odpowiedział albo odmówił; cisza nie
 *                      jest orzeczeniem o braku,
 *   `modul-nieznany`   rdzeń nie wymienia tego modułu wcale — katalogu jego
 *                      okien nie ma czym sprawdzić,
 *   `katalog-pelny`    moduł buduje każde okno operacyjne swojego katalogu,
 *   `katalog-szerszy`  katalog niesie kody, których moduł nie buduje — wykaz
 *                      okien niezbudowanych, liczony, nie napisany.
 * Niezależnie od nich wychodzi `pozaKatalogiem`: kod, który moduł buduje,
 * a rdzeń mu go nie przypisuje. To nie brak modułu, tylko rozjazd z rdzeniem.
 *
 * Odczyt jest jeden na połączenie, nie jeden na moduł — moduły pytające o ten
 * sam katalog dałyby tyle samo zbędnych zapytań. Katalog leży w pamięci
 * podręcznej przypisanej do kanału (`WeakMap`): pierwsze `odczytaj()` pyta,
 * pozostałe czekają na tę samą odpowiedź. Odmowa nie zostaje w pamięci, więc
 * kolejne `odczytaj()` ponawia pytanie.
 *
 * Unieważnienie: transport ponawia połączenie pod tym samym kanałem, więc po
 * zerwaniu można trafić na rdzeń o innym katalogu. Dlatego byt nasłuchuje
 * odpowiedzi `module.list` na całym kanale — każdy odczyt, także cudzy,
 * odświeża katalog bez ani jednego zapytania stąd. Powłoka może wymusić
 * zapomnienie wprost: `zapomnijKatalogOkien(kanal)`.
 *
 * Użycie:
 *
 *   const katalog = utworzKatalogOkien(kanal, 'design',
 *     ['design-board', 'preview-window', 'assets-panel', 'prompt-builder']);
 *   pasek.append(katalog.zdanieElement('dn-pole-opis md-uczciwosc__opis'));
 *   void katalog.odczytaj();   // raz po montażu modułu
 *   // przy zamknięciu modułu: katalog.zamknij();
 */

/** Zdanie wypowiadane, dopóki rdzeń nie odpowiedział. */
export const KATALOG_W_ODCZYCIE = 'Katalog okien rdzenia dla tego modułu — odczyt w toku…';

/** Rozstrzygnięcie o katalogu modułu; ta sama wartość idzie w `data-katalog`. */
export type StanKataloguOkien =
  | 'nieustalone'
  | 'modul-nieznany'
  | 'katalog-pelny'
  | 'katalog-szerszy';

/** Rozjazd katalogu rdzenia z oknami, które moduł naprawdę buduje. */
export interface RozjazdOkien {
  /** Rozstrzygnięcie zbiorcze — bez pytania o szczegóły. */
  stan: StanKataloguOkien;
  /**
   * Kody okien operacyjnych katalogu rdzenia bez okna rozmowy, w kolejności
   * podanej przez rdzeń. Puste, dopóki rdzeń nie orzekł.
   */
  wKatalogu: readonly string[];
  /** Kody katalogu, które moduł buduje. */
  zbudowane: readonly string[];
  /** Kody katalogu, których moduł nie buduje — wykaz okien niezbudowanych. */
  niezbudowane: readonly string[];
  /** Kody budowane przez moduł, których rdzeń mu nie przypisuje. */
  pozaKatalogiem: readonly string[];
  /** Zdanie odmowy odczytu; puste, dopóki nic nie odmówiło. */
  odmowa: string;
}

/** Katalog okien modułu wraz z paskami, które z niego biorą swoje zdanie. */
export interface KatalogOkien {
  /** Rozjazd bez budowania czegokolwiek — dla okien liczących coś własnego. */
  rozjazd(): RozjazdOkien;
  /** Samo zdanie paska uczciwości — dla modułów budujących nośnik własny. */
  zdanie(): string;
  /**
   * Akapit paska uczciwości przerysowujący się po każdej zmianie katalogu.
   *
   * @param klasa klasa rodziny modułu; nazwy klas należą do modułu,
   *   więc ten plik żadnej nie narzuca.
   */
  zdanieElement(klasa?: string): HTMLElement;
  /** Przerysowanie kontrolki własnej po każdej zmianie katalogu. Woła się od razu. */
  naOdczyt(przerysuj: () => void): void;
  /** Pyta rdzeń o katalog okien. Woła się raz, po montażu modułu. */
  odczytaj(): Promise<void>;
  /** Odczyt wymuszony — po ponowieniu połączenia albo migracji rdzenia. */
  odswiez(): Promise<void>;
  /** Odpina kontrolki tego modułu od wspólnego katalogu. Wołane z `zamknij()`. */
  zamknij(): void;
}

/**
 * Czy kod wskazuje okno rozmowy — w obu postaciach, jakie niesie rdzeń.
 *
 * Rdzeń podaje `chat-window` bez przedrostka, a dopuszcza także postać
 * `<kod modułu>.chat-window`. Obie są prawdziwe, więc sprawdzane są obie —
 * inaczej pasek wymieniałby okno rozmowy jako niezbudowane w każdym module,
 * choć okno rozmowy działa.
 */
export function czyOknoRozmowy(kod: string): boolean {
  return kod === 'chat-window' || kod.endsWith('.chat-window');
}

/** Odczyt katalogu okien — jeden na połączenie. */
interface Odczyt {
  /** Katalog po kodzie modułu; `null` = rdzeń jeszcze nie orzekł. */
  katalog: ReadonlyMap<string, readonly string[]> | null;
  /** Zdanie odmowy odczytu; puste, dopóki nic nie odmówiło. */
  odmowa: string;
}

/** Wpis pamięci podręcznej jednego kanału. */
interface WpisPamieci {
  odczyt: Odczyt;
  /** Odczyt w drodze — wszystkie moduły kanału czekają na jedną odpowiedź. */
  wToku: Promise<void> | null;
  /** Kontrolki wszystkich modułów tego kanału, do przerysowania po odczycie. */
  zalezni: Set<() => void>;
}

const pamiec = new WeakMap<Kanal, WpisPamieci>();

/**
 * Katalog okien jednego modułu.
 *
 * @param kanal kanał modułu; moduł nie sięga po niego sam.
 * @param kodModulu kod modułu w katalogu rdzenia — kolumna `modul.kod`, nie
 *   identyfikator wiersza i nie nazwa pozycji nawigacji.
 * @param zbudowane kody okien operacyjnych, które moduł naprawdę buduje. Tego
 *   rdzeń nie wie i wiedzieć nie może; okno rozmowy pomija się, bo montuje je
 *   scena sesji.
 */
export function utworzKatalogOkien(
  kanal: Kanal,
  kodModulu: string,
  zbudowane: readonly string[],
): KatalogOkien {
  const wspolny = wpisPamieci(kanal);
  /** Kontrolki tego modułu; wspólny wpis zna je jako jedno przerysowanie. */
  const moje: Array<() => void> = [];
  const przerysujMoje = (): void => {
    for (const przerysuj of moje) przerysuj();
  };
  wspolny.zalezni.add(przerysujMoje);

  /** Podpięcie kontrolki: przerysowanie teraz i po każdej zmianie katalogu. */
  function podepnij(przerysuj: () => void): void {
    moje.push(przerysuj);
    przerysuj();
  }

  const policz = (): RozjazdOkien => zlozRozjazd(wspolny.odczyt, kodModulu, zbudowane);

  return {
    rozjazd: policz,
    zdanie: () => zdanieRozjazdu(policz(), kodModulu),

    zdanieElement(klasa = '') {
      const element = document.createElement('p');
      if (klasa !== '') element.className = klasa;
      podepnij(() => {
        const rozjazd = policz();
        element.textContent = zdanieRozjazdu(rozjazd, kodModulu);
        element.dataset['katalog'] = rozjazd.stan;
      });
      return element;
    },

    naOdczyt: podepnij,
    odczytaj: () => zapewnijOdczyt(kanal),
    odswiez: () => {
      zapomnijKatalogOkien(kanal);
      return zapewnijOdczyt(kanal);
    },
    zamknij: () => {
      wspolny.zalezni.delete(przerysujMoje);
      moje.length = 0;
    },
  };
}

/**
 * Zapomnienie katalogu — po zerwaniu połączenia albo migracji rdzenia.
 *
 * Paski wracają do zdania „odczyt w toku”, a nie do ostatniego orzeczenia:
 * katalog rdzenia, którego już nie ma, nie jest wiedzą o rdzeniu, który
 * przyjdzie.
 */
export function zapomnijKatalogOkien(kanal: Kanal): void {
  zapisz(wpisPamieci(kanal), { katalog: null, odmowa: '' });
}

/** Wpis kanału wraz z nasłuchem odczytów — zakładany raz na kanał. */
function wpisPamieci(kanal: Kanal): WpisPamieci {
  const znany = pamiec.get(kanal);
  if (znany !== undefined) return znany;
  const swiezy: WpisPamieci = {
    odczyt: { katalog: null, odmowa: '' },
    wToku: null,
    zalezni: new Set(),
  };
  pamiec.set(kanal, swiezy);
  // Odczyt cudzy jest tą samą odpowiedzią co własny: powłoka pyta o moduły przy
  // budowie nawigacji, więc katalog odświeża się sam, bez zapytania stąd.
  kanal.naDowolny((koperta) => {
    if (koperta.type !== Command.ModuleList || koperta.status !== EnvelopeStatus.Ok) return;
    const moduly = (koperta.payload as ModuleListResponse | undefined)?.modules;
    if (!czyTablica(moduly)) return;
    zapisz(swiezy, { katalog: zlozKatalog(moduly as ModuleListResponse['modules']), odmowa: '' });
  });
  return swiezy;
}

/** Pyta rdzeń o katalog okien, jeśli nikt jeszcze nie zapytał ani nie wie. */
async function zapewnijOdczyt(kanal: Kanal): Promise<void> {
  const wspolny = wpisPamieci(kanal);
  if (wspolny.odczyt.katalog !== null) return;
  if (wspolny.wToku !== null) return wspolny.wToku;
  const bieg = odczytajModuly(kanal, wspolny);
  wspolny.wToku = bieg;
  await bieg;
  // Odmowa nie zostaje w pamięci: `katalog` dalej `null`, więc następne
  // `odczytaj()` ponowi pytanie.
  wspolny.wToku = null;
}

/** Jeden odczyt `module.list` i zapis jego wyniku — udanego albo odmownego. */
async function odczytajModuly(kanal: Kanal, wspolny: WpisPamieci): Promise<void> {
  const wynik = sprawdzKsztalt(
    await wywolaj(kanal, Command.ModuleList, {}),
    Command.ModuleList,
    (tresc) => czyTablica(tresc.modules),
  );
  if (!wynik.udany || wynik.wynik === undefined) {
    zapisz(wspolny, {
      katalog: null,
      odmowa: opisOdmowyBledu('Odczyt katalogu okien modułów', wynik.blad),
    });
    return;
  }
  zapisz(wspolny, { katalog: zlozKatalog(wynik.wynik.modules), odmowa: '' });
}

/** Katalog po kodzie modułu, z odjętym oknem rozmowy. */
function zlozKatalog(moduly: ModuleListResponse['modules']): ReadonlyMap<string, readonly string[]> {
  const katalog = new Map<string, readonly string[]>();
  for (const modul of moduly) {
    katalog.set(modul.code, (modul.operationalWindowCodes ?? []).filter((kod) => !czyOknoRozmowy(kod)));
  }
  return katalog;
}

/** Zapis odczytu i przerysowanie wszystkich pasków wszystkich modułów. */
function zapisz(wspolny: WpisPamieci, odczyt: Odczyt): void {
  wspolny.odczyt = odczyt;
  for (const przerysuj of wspolny.zalezni) przerysuj();
}

/** Rozjazd katalogu rdzenia z oknami modułu — same kody, bez składania zdania. */
function zlozRozjazd(
  odczyt: Odczyt,
  kodModulu: string,
  zbudowane: readonly string[],
): RozjazdOkien {
  const puste = {
    wKatalogu: [] as readonly string[],
    zbudowane: [] as readonly string[],
    niezbudowane: [] as readonly string[],
    pozaKatalogiem: [] as readonly string[],
  };
  if (odczyt.katalog === null) {
    return { stan: 'nieustalone', ...puste, odmowa: odczyt.odmowa };
  }
  const wKatalogu = odczyt.katalog.get(kodModulu);
  if (wKatalogu === undefined) {
    // Moduł nieznany rdzeniowi: kodów budowanych nie nazywamy wtedy
    // „poza katalogiem" — katalogu nie ma z czym porównać.
    return { stan: 'modul-nieznany', ...puste, odmowa: odczyt.odmowa };
  }
  const stoiWKatalogu = new Set(wKatalogu);
  const buduje = new Set(zbudowane.filter((kod) => !czyOknoRozmowy(kod)));
  const zbudowaneZKatalogu = wKatalogu.filter((kod) => buduje.has(kod));
  const niezbudowane = wKatalogu.filter((kod) => !buduje.has(kod));
  return {
    stan: niezbudowane.length === 0 ? 'katalog-pelny' : 'katalog-szerszy',
    wKatalogu,
    zbudowane: zbudowaneZKatalogu,
    niezbudowane,
    pozaKatalogiem: [...buduje].filter((kod) => !stoiWKatalogu.has(kod)),
    odmowa: odczyt.odmowa,
  };
}

/**
 * Zdanie paska uczciwości.
 *
 * Nie orzeka o niczym, czego nie powiedział rdzeń, i nie liczy braków przed
 * odpowiedzią.
 */
function zdanieRozjazdu(rozjazd: RozjazdOkien, kodModulu: string): string {
  const poza = rozjazd.pozaKatalogiem;
  const ogon =
    poza.length === 0
      ? ''
      : ` Moduł buduje przy tym ${nazwijOkna(poza)}, ` +
        `${poza.length === 1 ? 'którego' : 'których'} rdzeń mu nie przypisuje — ` +
        'rozjazd z katalogiem, nie brak modułu.';

  switch (rozjazd.stan) {
    case 'nieustalone':
      return rozjazd.odmowa === '' ? KATALOG_W_ODCZYCIE : `${rozjazd.odmowa}.`;
    case 'modul-nieznany':
      return (
        `Rdzeń nie wymienia modułu ${kodModulu} w wykazie modułów — katalogu jego okien ` +
        'nie ma czym sprawdzić.'
      );
    case 'katalog-pelny':
      return rozjazd.wKatalogu.length === 0
        ? `Rdzeń nie przypisuje modułowi ${kodModulu} ani jednego okna operacyjnego poza ` +
            `oknem rozmowy.${ogon}`
        : `Moduł buduje każde okno operacyjne swojego katalogu w rdzeniu ` +
            `(${rozjazd.wKatalogu.join(', ')}) — zmierzone odczytem module.list, nie wpisane ` +
            `na stałe.${ogon}`;
    default:
      return (
        `Katalog rdzenia przypisuje temu modułowi ${odmien(rozjazd.wKatalogu.length)}; ` +
        `moduł buduje ${rozjazd.zbudowane.length}. Niezbudowane: ` +
        `${rozjazd.niezbudowane.join(', ')} — zmierzone odczytem module.list.${ogon}`
      );
  }
}

/** Kody okien jako część zdania, z liczbą pojedynczą i mnogą. */
function nazwijOkna(kody: readonly string[]): string {
  return kody.length === 1 ? `okno ${kody[0]}` : `okna ${kody.join(', ')}`;
}

/**
 * Liczba wraz z odmienionym „okno operacyjne” — rzeczownik i przymiotnik.
 *
 * Zdanie składane z samej liczby brzmi „4 okien operacyjnych", czyli
 * niepoprawnie. Polszczyzna ma tu trzy przypadki: 1 → „okno operacyjne",
 * 2–4 → „okna operacyjne", reszta → „okien operacyjnych", z wyjątkiem nastek
 * (12, 13, 14), które idą jak reszta. Odmiana stoi w jednym miejscu, bo liczba
 * i końcówka rozjeżdżają się przy przepisywaniu.
 */
function odmien(ile: number): string {
  const nastka = ile % 100;
  if (ile === 1) return '1 okno operacyjne';
  if (ile % 10 >= 2 && ile % 10 <= 4 && (nastka < 12 || nastka > 14)) {
    return `${ile} okna operacyjne`;
  }
  return `${ile} okien operacyjnych`;
}
