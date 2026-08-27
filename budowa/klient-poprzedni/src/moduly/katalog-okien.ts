import { Command, type ModuleListResponse } from '../../../shared/contract';
import { EnvelopeStatus } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';

/** Katalog okien — jedno źródło prawdy o oknach operacyjnych, których moduł nie zbudował. */

/** Zdanie wypowiadane w pasku uczciwości, dopóki rdzeń nie odpowiedział na odczyt katalogu okien modułu. */
export const KATALOG_W_ODCZYCIE = 'Katalog okien rdzenia dla tego modułu — odczyt w toku…';

/** Rozstrzygnięcie o katalogu modułu wynikające z porównania z rdzeniem; ta sama wartość trafia do atrybutu przechowującego stan. */
export type StanKataloguOkien =
  | 'nieustalone'
  | 'modul-nieznany'
  | 'katalog-pelny'
  | 'katalog-szerszy';

/** Rozjazd katalogu okien rdzenia z oknami, które moduł naprawdę buduje, wraz z ewentualną odmową odczytu. */
export interface RozjazdOkien {
  /** Rozstrzygnięcie zbiorcze — bez pytania o szczegóły. */
  stan: StanKataloguOkien;
  /** Kody okien operacyjnych katalogu rdzenia bez okna rozmowy, w kolejności podanej przez rdzeń. */
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

/** Katalog okien modułu wraz z paskami uczciwości, które biorą z niego swoje zdanie o rozjeździe z rdzeniem. */
export interface KatalogOkien {
  /** Rozjazd bez budowania czegokolwiek — dla okien liczących coś własnego. */
  rozjazd(): RozjazdOkien;
  /** Samo zdanie paska uczciwości — dla modułów budujących nośnik własny. */
  zdanie(): string;
  /** Akapit paska uczciwości przerysowujący się po każdej zmianie katalogu; klasa należy do modułu. */
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
 * Czy kod wskazuje okno rozmowy, w obu postaciach niesionych przez rdzeń: bez przedrostka
 * modułu albo z przedrostkiem kodu modułu — obie postacie sprawdzane są łącznie.
 */
export function czyOknoRozmowy(kod: string): boolean {
  return kod === 'chat-window' || kod.endsWith('.chat-window');
}

/** Odczyt katalogu okien wspólny całemu kanałowi połączenia — jeden odczyt na połączenie, nie na moduł. */
interface Odczyt {
  /** Katalog po kodzie modułu; `null` = rdzeń jeszcze nie orzekł. */
  katalog: ReadonlyMap<string, readonly string[]> | null;
  /** Zdanie odmowy odczytu; puste, dopóki nic nie odmówiło. */
  odmowa: string;
}

/** Wpis pamięci podręcznej jednego kanału, niosący ostatni odczyt oraz zależne od niego kontrolki modułów. */
interface WpisPamieci {
  odczyt: Odczyt;
  /** Odczyt w drodze — wszystkie moduły kanału czekają na jedną odpowiedź. */
  wToku: Promise<void> | null;
  /** Kontrolki wszystkich modułów tego kanału, do przerysowania po odczycie. */
  zalezni: Set<() => void>;
}

const pamiec = new WeakMap<Kanal, WpisPamieci>();

/**
 * Katalog okien jednego modułu, budowany z kanału modułu, kodu modułu w katalogu rdzenia
 * i kodów okien operacyjnych, które moduł naprawdę buduje z pominięciem okna rozmowy.
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

/** Wpis pamięci podręcznej kanału wraz z nasłuchem cudzych odczytów katalogu — zakładany raz na cały kanał. */
function wpisPamieci(kanal: Kanal): WpisPamieci {
  const znany = pamiec.get(kanal);
  if (znany !== undefined) return znany;
  const swiezy: WpisPamieci = {
    odczyt: { katalog: null, odmowa: '' },
    wToku: null,
    zalezni: new Set(),
  };
  pamiec.set(kanal, swiezy);
  // Odczyt cudzy jest tą samą odpowiedzią co własny: katalog odświeża się sam, bez zapytania stąd.
  kanal.naDowolny((koperta) => {
    if (koperta.type !== Command.ModuleList || koperta.status !== EnvelopeStatus.Ok) return;
    const moduly = (koperta.payload as ModuleListResponse | undefined)?.modules;
    if (!czyTablica(moduly)) return;
    zapisz(swiezy, { katalog: zlozKatalog(moduly as ModuleListResponse['modules']), odmowa: '' });
  });
  return swiezy;
}

/** Pyta rdzeń o katalog okien wyłącznie wtedy, gdy nikt jeszcze nie zapytał ani odpowiedzi jeszcze nie zna. */
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

/** Jeden odczyt katalogu modułów rdzenia i zapis jego wyniku we wspólnej pamięci — udanego albo odmownego. */
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

/** Katalog okien operacyjnych po kodzie modułu rdzenia, z odjętym oknem rozmowy niemontowanym przez moduł. */
function zlozKatalog(moduly: ModuleListResponse['modules']): ReadonlyMap<string, readonly string[]> {
  const katalog = new Map<string, readonly string[]>();
  for (const modul of moduly) {
    katalog.set(modul.code, (modul.operationalWindowCodes ?? []).filter((kod) => !czyOknoRozmowy(kod)));
  }
  return katalog;
}

/** Zapis wyniku odczytu we wspólnej pamięci i przerysowanie wszystkich pasków wszystkich zależnych modułów. */
function zapisz(wspolny: WpisPamieci, odczyt: Odczyt): void {
  wspolny.odczyt = odczyt;
  for (const przerysuj of wspolny.zalezni) przerysuj();
}

/** Rozjazd katalogu okien rdzenia z oknami, które moduł buduje — same kody bez składania zdania czytelnego dla Operatora. */
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
    // Moduł nieznany rdzeniowi: kodów budowanych nie nazywamy „poza katalogiem” — nie ma z czym porównać.
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

/** Kody okien operacyjnych złożone jako część zdania paska, z rozróżnieniem liczby pojedynczej i mnogiej. */
function nazwijOkna(kody: readonly string[]): string {
  return kody.length === 1 ? `okno ${kody[0]}` : `okna ${kody.join(', ')}`;
}

/**
 * Liczba wraz z odmienionym rzeczownikiem „okno operacyjne”: jeden przypadek dla liczby jeden,
 * drugi dla dwóch do czterech z wyjątkiem nastek, trzeci dla pozostałych liczb.
 */
function odmien(ile: number): string {
  const nastka = ile % 100;
  if (ile === 1) return '1 okno operacyjne';
  if (ile % 10 >= 2 && ile % 10 <= 4 && (nastka < 12 || nastka > 14)) {
    return `${ile} okna operacyjne`;
  }
  return `${ile} okien operacyjnych`;
}
