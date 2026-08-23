import {
  AccessMode,
  AccessPointKind,
  ChangeKind,
  type AccessGrant,
  type AccessGrantAddResponse,
  type AccessGrantUpdateRequest,
  type AccessGrantUpdateResponse,
  type AccessPoint,
  type AccessPointAddResponse,
  type AccessPointCheckResponse,
  type AccessPointRemoveResponse,
} from '../../../shared/contract';
import type { Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../protokol/kanal';
import { nazwaZeSciezki } from './nazwy-dostepow';
import { utworzZapisyNadan } from './zapisy-nadan';
import { uporzadkujNadania, utworzZrodloNadan } from './zrodlo-nadan';
import { uporzadkujPunkty, utworzZrodloPunktow } from './zrodlo-punktow';

/**
 * Stan sekcji dostępów — jedno źródło prawdy dla wykazu punktów i zbioru
 * nadań jednego okna rozmowy.
 *
 * Punkty i nadania trzymamy razem, ponieważ wiersz nadania nie da się
 * narysować bez punktu, na który się powołuje: z punktu pochodzą korzenie,
 * rodzaj i nazwa maszyny. Dwa równoległe stany dałyby dwie prawdy o tym samym
 * nadaniu.
 *
 * Nadanie żyje per okno. Zmiana okna nie przebudowuje sekcji — zmienia zbiór
 * nadań i ogłasza przeliczenie; wykaz punktów jest wspólny dla platformy
 * i zostaje.
 *
 * Żadna ścieżka nie zatrzymuje sekcji. Rdzeń, który nie odda wykazu, zostawia
 * go pustym; sekcja pozostaje czynna i pozwala spytać ponownie.
 */
export type FazaOdczytu = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface StanDostepow {
  /**
   * Faza odczytu wykazów z rdzenia.
   *
   * Bez niej pusty wykaz znaczy trzy rzeczy naraz: „jeszcze nie pytałem",
   * „pytam" i „rdzeń nie zna ani jednego punktu". Widok musi je rozróżnić,
   * bo każdej należy się inny stan: nic, wskaźnik odczytu, stan pusty.
   */
  faza(): FazaOdczytu;
  /** Powód ostatniego niepowodzenia odczytu; pusty, gdy odczyt się powiódł. */
  powodNiepowodzenia(): string;
  /** Punkty dostępu w porządku: rodzaj, potem nazwa. */
  punkty(): readonly AccessPoint[];
  /** Punkt o wskazanym identyfikatorze; `null`, gdy wykaz go nie zna. */
  punkt(punktID: string): AccessPoint | null;
  /** Nadania okna czynnego, w kolejności. */
  nadania(): readonly AccessGrant[];
  /** Okno rozmowy, którego nadania pokazuje sekcja. */
  oknoID(): string;
  /** Wiąże sekcję z innym oknem i wczytuje jego nadania. */
  ustawOkno(oknoID: string): void;
  /** Wczytuje wykaz punktów oraz nadania okna czynnego. */
  odswiez(): Promise<void>;
  /** `access.point.check` wraz z naniesieniem wyniku na wykaz. */
  sprawdz(punktID: string): Promise<Wynik<AccessPointCheckResponse>>;
  /**
   * Zakłada punkt rodzaju `localDirectory` na wskazanej ścieżce.
   *
   * `urzadzenieID` wskazuje maszynę, na której katalog istnieje — schemat bazy
   * wymaga go dla punktu rodzaju `localDirectory`. Katalog dodany z „Mój
   * komputer" należy do maszyny bieżącej; jej identyfikator poda komenda
   * `device.list`, gdy trafi do kontraktu. Do tego czasu wywołanie bez
   * urządzenia wraca z odmową merytoryczną rdzenia, nie z fałszywym sukcesem.
   */
  zalozKatalogLokalny(
    sciezka: string,
    urzadzenieID?: string,
  ): Promise<Wynik<AccessPointAddResponse>>;
  /** Usuwa punkt wraz z nadaniami, które się na niego powoływały. */
  usunPunkt(punktID: string): Promise<Wynik<AccessPointRemoveResponse>>;
  /** Nadaje oknu czynnemu dostęp do punktu. */
  nadaj(
    punktID: string,
    tryb: AccessMode,
    korzenie: readonly string[],
  ): Promise<Wynik<AccessGrantAddResponse>>;
  /** Zmienia nadanie: tryb, korzenie, kolejność albo oznaczenie głównego. */
  zmienNadanie(zadanie: AccessGrantUpdateRequest): Promise<Wynik<AccessGrantUpdateResponse>>;
  /** Odbiera oknu nadanie. */
  odbierz(nadanieID: string): Promise<Wynik<{ removed: boolean; grants: AccessGrant[] }>>;
  /** Subskrypcja przeliczenia stanu. */
  naZmiane(sluchacz: () => void): void;
  /** Odłącza subskrypcje kanału. */
  rozlacz(): void;
}

export function utworzStanDostepow(kanal: Kanal, oknoPoczatkowe = ''): StanDostepow {
  // Powód niepowodzenia zbiera się z obu odczytów jednego odświeżenia: wykaz
  // punktów i nadania okna jadą równolegle, a Operatorowi należy się zdanie
  // o każdym, który nie dojechał.
  let powody: string[] = [];
  const zapiszPowod = (powod: string): void => void powody.push(powod);

  const zrodloPunktow = utworzZrodloPunktow(kanal, zapiszPowod);
  const zrodloNadan = utworzZrodloNadan(kanal, zapiszPowod);

  let punkty: AccessPoint[] = [];
  let nadania: AccessGrant[] = [];
  let okno = oknoPoczatkowe;
  let faza: FazaOdczytu = 'spoczynek';

  const sluchacze: Array<() => void> = [];
  const oglos = (): void => {
    for (const sluchacz of [...sluchacze]) sluchacz();
  };

  /** Zmiana punktu z rdzenia — także dokonana w innym oknie albo urządzeniu. */
  const odsubskrybujPunkty: Odsubskrybuj = zrodloPunktow.naZmiane((zdarzenie) => {
    punkty =
      zdarzenie.change === ChangeKind.Deleted
        ? punkty.filter((punkt) => punkt.id !== zdarzenie.point.id)
        : uporzadkujPunkty([
            ...punkty.filter((punkt) => punkt.id !== zdarzenie.point.id),
            zdarzenie.point,
          ]);
    oglos();
  });

  /** Zmiana nadania; zdarzenie z innego okna nie rusza zbioru okna czynnego. */
  const odsubskrybujNadania: Odsubskrybuj = zrodloNadan.naZmiane((zdarzenie) => {
    if (zdarzenie.windowId !== okno) return;
    przyjmijNadanie(zdarzenie.grant, zdarzenie.change === ChangeKind.Deleted);
    oglos();
  });

  function przyjmijNadanie(nadanie: AccessGrant, usuniete: boolean): void {
    const pozostale = nadania.filter((istniejace) => istniejace.id !== nadanie.id);
    nadania = usuniete ? pozostale : uporzadkujNadania([...pozostale, nadanie]);
  }

  /** Komplet nadań z odpowiedzi zapisu zastępuje zbiór okna w całości. */
  function przyjmijKomplet(komplet: readonly AccessGrant[] | undefined): void {
    if (komplet !== undefined) nadania = uporzadkujNadania(komplet);
    oglos();
  }

  const zapisy = utworzZapisyNadan({
    zrodlo: zrodloNadan,
    oknoID: () => okno,
    przyjmijKomplet,
  });

  async function wczytajNadania(): Promise<void> {
    nadania = okno === '' ? [] : await zrodloNadan.lista({ windowId: okno });
  }

  return {
    faza: () => faza,

    powodNiepowodzenia: () => powody.join(' '),

    punkty: () => punkty,

    punkt: (punktID) => punkty.find((punkt) => punkt.id === punktID) ?? null,

    nadania: () => nadania,

    oknoID: () => okno,

    ustawOkno(nowe) {
      if (nowe === okno) return;
      okno = nowe;
      nadania = [];
      oglos();
      void wczytajNadania().then(oglos);
    },

    async odswiez() {
      powody = [];
      faza = 'odczyt';
      oglos();
      const [pobranePunkty] = await Promise.all([zrodloPunktow.lista({}), wczytajNadania()]);
      punkty = pobranePunkty;
      faza = powody.length > 0 ? 'blad' : 'gotowe';
      oglos();
    },

    async sprawdz(punktID) {
      const wynik = await zrodloPunktow.sprawdz(punktID);
      const tresc = wynik.wynik;
      if (wynik.udany && tresc !== undefined) {
        punkty = punkty.map((punkt) =>
          punkt.id === punktID
            ? { ...punkt, status: tresc.status, checkedAt: tresc.checkedAt }
            : punkt,
        );
      }
      oglos();
      return wynik;
    },

    async zalozKatalogLokalny(sciezka, urzadzenieID) {
      const wynik = await zrodloPunktow.dodaj({
        name: nazwaZeSciezki(sciezka),
        kind: AccessPointKind.LocalDirectory,
        roots: [sciezka],
        defaultMode: AccessMode.Read,
        description: `Katalog wskazany oknem powłoki: ${sciezka}`,
        // Pole opcjonalne kontraktu: wysyłamy je tylko, gdy znamy urządzenie.
        // Puste `deviceId` nie przechodzi więzu schematu, więc pominięcie jest
        // uczciwsze niż napis pusty — rdzeń odmówi z powodem, a nie z błędu bazy.
        ...(urzadzenieID !== undefined && urzadzenieID !== '' ? { deviceId: urzadzenieID } : {}),
      });
      const punkt = wynik.wynik?.point;
      if (wynik.udany && punkt !== undefined) {
        punkty = uporzadkujPunkty([...punkty.filter((istniejacy) => istniejacy.id !== punkt.id), punkt]);
      }
      oglos();
      return wynik;
    },

    async usunPunkt(punktID) {
      const wynik = await zrodloPunktow.usun(punktID);
      if (wynik.udany && wynik.wynik?.removed === true) {
        punkty = punkty.filter((punkt) => punkt.id !== punktID);
        nadania = nadania.filter((nadanie) => nadanie.accessPointId !== punktID);
      }
      oglos();
      return wynik;
    },

    nadaj: zapisy.nadaj,

    zmienNadanie: zapisy.zmien,

    odbierz: zapisy.odbierz,

    naZmiane: (sluchacz) => void sluchacze.push(sluchacz),

    rozlacz() {
      odsubskrybujPunkty();
      odsubskrybujNadania();
    },
  };
}
