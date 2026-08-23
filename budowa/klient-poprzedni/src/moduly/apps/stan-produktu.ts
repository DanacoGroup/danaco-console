import {
  EventType,
  type AppArchitecture,
  type AppComponent,
  type AppDeployEnvironment,
  type AppDeployment,
  type AppStage,
  type AppWorkspaceLayer,
  type DeveloperFile,
  type ProgressChangedEvent,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';
import { utworzWarsztatPlikow, type WarsztatPlikow } from './warsztat-plikow';
import { utworzZbiorBudowy, type RachunekRamek, type ZbiorBudowy } from './zbior-budowy';
import { utworzZrodloApps, type ZrodloApps } from './zrodlo-apps';
import { utworzZrodloOknaModulu, wybierzOknoModulu } from './zrodlo-okna-modulu';

/** Faza odczytu okna modułu — trzy pustki rozróżnialne (`dostepy/stany-odczytu.ts`). */
export type FazaOdczytu = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

/**
 * Który z trzech odczytów obszaru — klucz powodu odmowy.
 *
 * Nazwy są własne modułu, nie nazwami komend: okno pyta o powód odczytu
 * wdrożeń, a nie o powód komendy `apps.deployment.list`, bo warstwa widoku
 * nazw komend nie zna.
 */
export type RodzajOdczytu = 'wdrozenia' | 'architektura';

/** Zdanie o odczycie, który nie mógł ruszyć, bo sesja nie ma jeszcze okna modułu. */
const BEZ_OKNA_DO_ODCZYTU =
  'Rdzeń nie wskazał jeszcze okna modułu Apps, a odczyt wymaga jego identyfikatora — ' +
  'żądanie nie zostało wysłane.';

/**
 * Jedno źródło prawdy modułu Apps: okno rdzenia, architektura, etapy budowy
 * i wdrożenia.
 *
 * Dwa równoległe stany dałyby dwie prawdy o tym samym produkcie. Architecture
 * Designer definiuje komponenty, oba warsztaty przypisują do nich pliki,
 * Product Builder rysuje z nich oś etapów, a Deployment Panel wdraża —
 * wszystkie patrzą tutaj.
 *
 * Trzy komendy odczytu obszaru (`apps.deployment.list`, `apps.architecture.get`,
 * `apps.workspace.list`) wypełniają stan tym, co rdzeń trzyma w bazie, zamiast
 * tym, co przeleciało gniazdem w bieżącej sesji. Bez nich odświeżenie okna
 * przeglądarki zerowałoby moduł, choć historia wdrożeń, architektura i pliki
 * warsztatu leżą w rdzeniu.
 *
 * Etapy budowy zostają wyłącznie przy zdarzeniu: komendy ich odczytu kontrakt
 * nie niesie, więc pusty wykaz etapów na starcie jest stanem prawdziwym, a nie
 * brakiem odczytu — i okno mówi o nim inaczej niż o wdrożeniach.
 *
 * Odczyt nie gasi stanu, który już jest: każdy z trzech odczytów jest
 * niezależny, odmowa jednego zostawia dwa pozostałe nietknięte i nie kasuje
 * tego, co moduł już wie. Powód odmowy trafia do osobnego pola dla każdego
 * odczytu, bo okna czytają je w różnych miejscach ekranu.
 */
export interface StanProduktu {
  zrodlo: ZrodloApps;
  /** Okno modułu w rdzeniu; pusty łańcuch, gdy rdzeń go jeszcze nie wskazał. */
  idOkna(): string;
  faza(): FazaOdczytu;
  powodNiepowodzenia(): string;
  architektura(): AppArchitecture | null;
  komponenty(): readonly AppComponent[];
  etapy(): readonly AppStage[];
  wdrozenia(): readonly AppDeployment[];
  /**
   * Rachunek ramek `apps.build.changed`, które przyszły do tego modułu.
   *
   * Okno pyta o niego, gdy wykaz etapów jest pusty: bez tych liczb nie da się
   * odróżnić „jeszcze nic nie przyszło" od „przyszło, ale etapów w tym nie
   * było", a różnica rozstrzyga, co wolno powiedzieć Operatorowi.
   */
  ramki(): RachunekRamek;
  /** Czy wykaz wdrożeń wrócił już z rdzenia — „pusto" to nie to samo, co „nie pytano". */
  czyWdrozeniaCzytane(): boolean;
  /** Pliki warsztatu jednej warstwy; puste, dopóki odczyt nie wrócił. */
  plikiWarsztatu(warstwa: AppWorkspaceLayer): readonly DeveloperFile[];
  czyWarsztatCzytany(warstwa: AppWorkspaceLayer): boolean;
  /**
   * Powód odmowy ostatniego odczytu danego rodzaju; pusty, gdy odczyt się udał
   * albo nie był jeszcze zlecony. Kluczem jest nazwa własna modułu, nie nazwa
   * komendy — zdanie o odmowie składa źródło, okno bierze je gotowe.
   */
  powodOdczytu(co: RodzajOdczytu): string;
  /**
   * Powód odmowy odczytu warsztatu jednej warstwy.
   *
   * Osobno dla każdej warstwy, bo oba warsztaty czytają równolegle i mają dwa
   * osobne okna. Jedno pole na oba znaczyłoby, że udany odczyt backendu kasuje
   * odmowę frontendu, a Frontend Workspace milczałby o niepowodzeniu, które
   * właśnie go spotkało.
   */
  powodOdczytuWarsztatu(warstwa: AppWorkspaceLayer): string;
  postep(): ProgressChangedEvent | null;
  /** Zleca `apps.deployment.list` dla okna modułu i wchłania wykaz. */
  odczytajWdrozenia(srodowisko: AppDeployEnvironment | '', granica: number): Promise<void>;
  /**
   * Zleca `apps.architecture.get` dla okna modułu i wchłania architekturę.
   *
   * Oddaje `true`, gdy rdzeń przysłał architekturę, a `false`, gdy odpowiedział
   * udanie i bez niej albo gdy odczyt się nie odbył. Okno mówi o tych dwóch
   * przypadkach różnymi zdaniami, więc rozróżnienie musi wyjść ze stanu, a nie
   * być zgadywane z tego, co zostało na kanwie.
   */
  odczytajArchitekture(): Promise<boolean>;
  /** Zleca `apps.workspace.list` dla jednej warstwy i wchłania jej pliki. */
  odczytajWarsztat(warstwa: AppWorkspaceLayer): Promise<void>;
  /** Zapamiętuje plik potwierdzony przez rdzeń po zapisie warsztatu. */
  wchlonZapisWarsztatu(warstwa: AppWorkspaceLayer, plik: DeveloperFile): void;
  /** Zapamiętuje komponent zbudowany w Architecture Designerze. */
  dodajKomponent(komponent: AppComponent): void;
  usunKomponent(idKomponentu: string): void;
  /** Wchłania architekturę potwierdzoną przez rdzeń. */
  wchlonArchitekture(architektura: AppArchitecture): void;
  /**
   * Zasiewa wdrożenie migawką z odpowiedzi na `apps.deployment.run`.
   *
   * Nazwa mówi o źródle, bo źródło rozstrzyga: odpowiedź komendy niesie stan
   * z chwili ruszenia przebiegu (`pending`), a postęp przychodzi wyłącznie
   * zdarzeniem `apps.build.changed`. Zbiór budowy nie pozwala migawce cofnąć
   * stanu, który zdarzenie już przyniosło (`zbior-budowy.ts`).
   */
  wchlonOdpowiedzWdrozenia(wdrozenie: AppDeployment): void;
  /** Odczytuje okno modułu z rdzenia; wolno wołać wielokrotnie. */
  odswiez(idSesji: string): Promise<void>;
  obserwuj(sluchacz: () => void): Odsubskrybuj;
  rozlacz(): void;
}

export function utworzStanProduktu(kanal: Kanal): StanProduktu {
  const zrodlo = utworzZrodloApps(kanal);
  const oknaModulu = utworzZrodloOknaModulu(kanal);
  const budowa: ZbiorBudowy = utworzZbiorBudowy();
  const warsztat: WarsztatPlikow = utworzWarsztatPlikow();
  const sluchacze = new Set<() => void>();
  /** Powody odmów odczytu — jeden na rodzaj, żeby jedna odmowa nie zamazała drugiej. */
  const powodyOdczytu: Record<RodzajOdczytu, string> = { wdrozenia: '', architektura: '' };
  /** Powody odmów odczytu warsztatu — jeden na warstwę; klucz jest wartością kontraktu. */
  const powodyWarsztatu = new Map<string, string>();

  let idOkna = '';
  let faza: FazaOdczytu = 'spoczynek';
  let powod = '';
  let architektura: AppArchitecture | null = null;
  let komponenty: AppComponent[] = [];
  let postep: ProgressChangedEvent | null = null;

  function oglos(): void {
    for (const sluchacz of sluchacze) sluchacz();
  }

  /**
   * Wspólne wejście trzech odczytów: czyści powód poprzedni i sprawdza okno.
   *
   * Okno modułu jest warunkiem wszystkich trzech żądań, bo `windowId` jest
   * w nich polem wymaganym. Bez niego żądanie poszłoby po dane niczyje, a odmowa
   * rdzenia mówiłaby o brakującym polu zamiast o brakującym oknie sesji.
   * Zwraca `false`, gdy odczytu nie wolno zlecić.
   */
  function zacznijOdczyt(co: RodzajOdczytu): boolean {
    powodyOdczytu[co] = '';
    if (idOkna !== '') return true;
    powodyOdczytu[co] = BEZ_OKNA_DO_ODCZYTU;
    oglos();
    return false;
  }

  const odsubskrybowania: Odsubskrybuj[] = [
    zrodlo.naZmianeBudowy((tresc) => {
      budowa.wchlonZdarzenie(tresc);
      oglos();
    }),
    // Drugie zdarzenie obszaru. Zmiana pliku warsztatu wchodzi tu niezależnie
    // od tego, gdzie zaszła — w drugim oknie tej sesji, w obcym połączeniu czy
    // po stronie rdzenia. Zawężamy ją do okna modułu, bo zbiór plików jest
    // zbiorem tego okna, a ramka o cudzym oknie opisuje inny warsztat.
    zrodlo.naZmianeWarsztatu((tresc) => {
      if (idOkna !== '' && tresc.windowId !== idOkna) return;
      warsztat.wchlonZdarzenie(tresc);
      oglos();
    }),
    // Telemetria postępu jest wspólna całej platformie: Product
    // Builder czyta z niej etap bieżący, bo obszar `apps` nie ma własnego
    // zdarzenia postępu. Bierzemy wyłącznie proces okna modułu.
    kanal.naZdarzenie(EventType.ProgressChanged, (tresc) => {
      if (idOkna !== '' && tresc.windowId !== undefined && tresc.windowId !== idOkna) return;
      postep = tresc;
      oglos();
    }),
  ];

  return {
    zrodlo,
    idOkna: () => idOkna,
    faza: () => faza,
    powodNiepowodzenia: () => powod,
    architektura: () => architektura,
    komponenty: () => komponenty,
    etapy: () => budowa.etapy(),
    wdrozenia: () => budowa.wdrozenia(),
    ramki: () => budowa.ramki(),
    czyWdrozeniaCzytane: () => budowa.czyWdrozeniaCzytane(),
    plikiWarsztatu: (warstwa) => warsztat.pliki(warstwa),
    czyWarsztatCzytany: (warstwa) => warsztat.czyCzytana(warstwa),
    powodOdczytu: (co) => powodyOdczytu[co],
    powodOdczytuWarsztatu: (warstwa) => powodyWarsztatu.get(warstwa) ?? '',
    postep: () => postep,

    async odczytajWdrozenia(srodowisko, granica) {
      if (!zacznijOdczyt('wdrozenia')) return;
      const wynik = await zrodlo.wdrozenia({ idOkna, srodowisko, granica });
      if (!wynik.udany || wynik.wynik === undefined) {
        powodyOdczytu.wdrozenia = opisOdmowyBledu('Odczyt wdrożeń', wynik.blad);
        oglos();
        return;
      }
      budowa.wchlonWykazWdrozen(wynik.wynik.deployments);
      oglos();
    },

    async odczytajArchitekture() {
      if (!zacznijOdczyt('architektura')) return false;
      const wynik = await zrodlo.architektura(idOkna);
      if (!wynik.udany || wynik.wynik === undefined) {
        powodyOdczytu.architektura = opisOdmowyBledu('Odczyt architektury', wynik.blad);
        oglos();
        return false;
      }
      // Brak pola `architecture` jest odpowiedzią udaną i znaczy „okno nie ma
      // jeszcze żadnej architektury". Nie wpisujemy wtedy niczego i nie kasujemy
      // komponentów zestawionych na kanwie — Operator mógł je właśnie ułożyć,
      // a odczyt nie jest poleceniem sprzątania.
      const oddana = wynik.wynik.architecture;
      if (oddana !== undefined) {
        architektura = oddana;
        if (oddana.components !== undefined) komponenty = [...oddana.components];
      }
      oglos();
      return oddana !== undefined;
    },

    async odczytajWarsztat(warstwa) {
      powodyWarsztatu.set(warstwa, '');
      if (idOkna === '') {
        powodyWarsztatu.set(warstwa, BEZ_OKNA_DO_ODCZYTU);
        oglos();
        return;
      }
      const wynik = await zrodlo.plikiWarsztatu(idOkna, warstwa);
      if (!wynik.udany || wynik.wynik === undefined) {
        powodyWarsztatu.set(warstwa, opisOdmowyBledu('Odczyt plików warsztatu', wynik.blad));
        oglos();
        return;
      }
      warsztat.wchlonWykaz(warstwa, wynik.wynik.files);
      oglos();
    },

    wchlonZapisWarsztatu(warstwa, plik) {
      warsztat.wchlonZapis(warstwa, plik);
      oglos();
    },

    dodajKomponent(komponent) {
      komponenty = [...komponenty.filter((inny) => inny.id !== komponent.id), komponent];
      oglos();
    },

    usunKomponent(idKomponentu) {
      komponenty = komponenty.filter((inny) => inny.id !== idKomponentu);
      oglos();
    },

    wchlonArchitekture(nowa) {
      architektura = nowa;
      if (nowa.components !== undefined) komponenty = [...nowa.components];
      oglos();
    },

    wchlonOdpowiedzWdrozenia(wdrozenie) {
      budowa.wchlonOdpowiedzWdrozenia(wdrozenie);
      oglos();
    },

    async odswiez(idSesji) {
      faza = 'odczyt';
      powod = '';
      oglos();
      const wynik = await oknaModulu.okna(idSesji);
      if (!wynik.udany) {
        faza = 'blad';
        powod = wynik.blad?.message ?? 'rdzeń nie podał przyczyny';
        oglos();
        return;
      }
      idOkna = wybierzOknoModulu(wynik.wynik?.windows ?? []);
      faza = 'gotowe';
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    rozlacz() {
      for (const odsubskrybuj of odsubskrybowania.splice(0)) odsubskrybuj();
      sluchacze.clear();
    },
  };
}
