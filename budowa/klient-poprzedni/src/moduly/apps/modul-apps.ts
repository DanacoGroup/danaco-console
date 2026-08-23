import './apps.css';

import { AppWorkspaceLayer } from '../../../../shared/contract';

import { utworzPasPomocniczych, type PasPomocniczych } from '../../okna-pomocnicze/indeks';
import type { Kanal } from '../../protokol/kanal';
import { utworzPokrycieKomend, type PokrycieKomend } from '../pokrycie-komend';
import { KODY_OKIEN, NAZWY_OKIEN } from './etykiety-apps';
import { utworzOknoAppCatalog } from './okno-app-catalog';
import { utworzOknoArchitectureDesigner } from './okno-architecture-designer';
import { utworzOknoDeploymentPanel } from './okno-deployment-panel';
import { utworzOknoInstalledAppsManager } from './okno-installed-apps-manager';
import { utworzOknoIntegrationsHub } from './okno-integrations-hub';
import { utworzOknoMcpConnectorConsole } from './okno-mcp-connector-console';
import { utworzOknoPermissionsTrustCenter } from './okno-permissions-trust-center';
import { utworzOknoProductBuilder } from './okno-product-builder';
import { utworzOknoPublisherPanel } from './okno-publisher-panel';
import { utworzOknoBackendWorkspace } from './okno-backend-workspace';
import { utworzOknoFrontendWorkspace } from './okno-frontend-workspace';
import type { PozycjaOsi } from './os-etapow';
import { utworzPanelAkcji } from './panel-akcji';
import { utworzStanProduktu, type StanProduktu } from './stan-produktu';
import { utworzStanRozszerzen, type StanRozszerzen } from './stan-rozszerzen';
import { utworzZrodloIzolacjiApps } from './zrodlo-izolacji-apps';
import { KOD_MODULU } from './zrodlo-okna-modulu';

/**
 * Moduł Apps — pięć okien operacyjnych osadzonych w jednym układzie.
 *
 * Układ wynika z roli okna. Product Builder jest punktem wejścia integrującym
 * pozostałe okna, więc stoi w pasie pierwszym na całą szerokość. Architecture
 * Designer poprzedza pracę w warsztatach, więc stoi nad nimi. Oba warsztaty
 * pracują równolegle, więc stoją obok siebie. Deployment Panel zamyka proces,
 * więc stoi na końcu.
 *
 * Stan produktu jest jeden na cały moduł: komponent zestawiony w Architecture
 * Designerze pojawia się natychmiast w wykazach obu warsztatów, a wdrożenie
 * potwierdzone przez rdzeń — w dzienniku wydań Product Buildera.
 *
 * Komendy obszaru mają uchwyt w rdzeniu: zapis architektury, zapis pliku
 * warsztatu i uruchomienie wdrożenia wpina `zarejestrujAplikacje`
 * (`server/internal/core/adapter_modul_aplikacje_uchwyty.go`), wołane
 * z `kompozycja.go`; port `Aplikacje` jest wypełniony w `montaz_porty.go`.
 * Wdrożenie zmienia stan (`pending` → `running` → `succeeded`/`failed`),
 * a przejścia przychodzą zdarzeniem `apps.build.changed`. Okna pokazują
 * odpowiedź rdzenia, a odmowę merytoryczną — jako odmowę, nie jako pustkę
 * i nie jako sukces.
 *
 * Na końcu układu stoi pas okien pomocniczych
 * (`okna-pomocnicze/pas-pomocniczych.ts`) — ten sam, którym stoją Developer
 * i Diagnostics. Pas składa panele z jednej wytwórni (`wytwornia-paneli.ts`)
 * po spisie modułu (`rejestr-pomocniczych.ts`), a pozycje niezbudowane wypisuje
 * wraz z powodem, więc brak zostaje widoczny. Terminal wchodzi tą samą drogą co
 * każdy inny panel — wpisem w wytwórni i pozycją w spisie modułu `apps` — a nie
 * osobnym wywołaniem w tym pliku.
 *
 * Pas dostaje okno później, niż powstaje. Developer i Diagnostics montują się
 * przez `widokZOknaSesji`, więc znają okno już przy montażu. Apps montuje się
 * z samym kanałem (`indeks.ts` → `utworzModulApps(kanal)`), a okno modułu
 * poznaje dopiero z `window.list` w `wczytaj`. Pas stoi więc od początku
 * z oknem pustym — panele mówią wtedy wprost, czego brakuje — i przyjmuje
 * właściwe okno wywołaniem `ustawOkno`, które samo zamyka panele stojące
 * i stawia je na nowym oknie.
 */
export interface ModulApps {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt okna modułu i katalogu akcji rdzenia. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje zdarzeń rdzenia. */
  rozlacz(): void;
}

/** Nazwa modułu w etykietach dostępności pasa okien pomocniczych. */
const NAZWA_MODULU = 'Apps';

/** Etapy procesu w kolejności: architektura → warsztaty → wdrożenie. */
const POZYCJE_OSI: readonly PozycjaOsi[] = [
  {
    kod: KODY_OKIEN.ArchitectureDesigner,
    tytul: NAZWY_OKIEN[KODY_OKIEN.ArchitectureDesigner] ?? '',
    opis: 'Definiowanie komponentów rozwiązania i ustalenie zależności',
  },
  {
    kod: KODY_OKIEN.FrontendWorkspace,
    tytul: NAZWY_OKIEN[KODY_OKIEN.FrontendWorkspace] ?? '',
    opis: 'Edycja warstwy interfejsu i podgląd wyniku',
  },
  {
    kod: KODY_OKIEN.BackendWorkspace,
    tytul: NAZWY_OKIEN[KODY_OKIEN.BackendWorkspace] ?? '',
    opis: 'Edycja logiki serwerowej i konfiguracja usług',
  },
  {
    kod: KODY_OKIEN.DeploymentPanel,
    tytul: NAZWY_OKIEN[KODY_OKIEN.DeploymentPanel] ?? '',
    opis: 'Uruchomienie wdrożenia i przegląd statusu publikacji',
  },
  // Dwa etapy strony dystrybucji. Oś prowadzi do nich, bo opracowanie wymienia
  // App Catalog i Integrations Hub wśród kafli nawigacyjnych Product Buildera,
  // a proces produktu nie kończy się na wdrożeniu: zbudowane rozszerzenie wraca
  // do katalogu, z którego korzystają moduły i eksperci.
  {
    kod: KODY_OKIEN.AppCatalog,
    tytul: NAZWY_OKIEN[KODY_OKIEN.AppCatalog] ?? '',
    opis: 'Przegląd katalogu rozszerzeń i instalacja pozycji',
  },
  {
    kod: KODY_OKIEN.IntegrationsHub,
    tytul: NAZWY_OKIEN[KODY_OKIEN.IntegrationsHub] ?? '',
    opis: 'Konektory i serwery MCP wraz ze sprawdzeniem mostów',
  },
];

export function utworzModulApps(kanal: Kanal): ModulApps {
  const stan: StanProduktu = utworzStanProduktu(kanal);
  // Drugi stan, bo druga strona modułu stoi na innym obszarze kontraktu:
  // produkt budowany w module jest bytem okna (`apps.*` niesie `windowId`),
  // a katalog rozszerzeń stoi poziom wyżej i okna nie zna. Jeden stan na oba
  // znaczyłby, że odczyt katalogu czeka na okno, którego nie potrzebuje.
  const rejestr: StanRozszerzen = utworzStanRozszerzen(kanal);
  const izolacja = utworzZrodloIzolacjiApps(kanal, KOD_MODULU);
  const pokrycie: PokrycieKomend = utworzPokrycieKomend(kanal);
  const akcje = utworzPanelAkcji(kanal);

  const architektura = utworzOknoArchitectureDesigner(stan);
  // Tożsamość każdego okna rozstrzyga jego własny plik; ten plik tylko je
  // zestawia w układ.
  const frontend = utworzOknoFrontendWorkspace(stan);
  const backend = utworzOknoBackendWorkspace(stan);
  const wdrozenia = utworzOknoDeploymentPanel(stan);
  const builder = utworzOknoProductBuilder(stan, POZYCJE_OSI, (kod) => przeniesOgnisko(kod));

  const katalog = utworzOknoAppCatalog(rejestr);
  const zainstalowane = utworzOknoInstalledAppsManager(rejestr);
  const integracje = utworzOknoIntegrationsHub(rejestr);
  const uprawnienia = utworzOknoPermissionsTrustCenter(rejestr, izolacja, () => stan.idOkna());
  const konsola = utworzOknoMcpConnectorConsole(rejestr, pokrycie);
  const wydawca = utworzOknoPublisherPanel(stan, pokrycie);

  const wszystkie = [
    builder.element,
    architektura.element,
    frontend.element,
    backend.element,
    wdrozenia.element,
    katalog.element,
    zainstalowane.element,
    integracje.element,
    uprawnienia.element,
    konsola.element,
    wydawca.element,
  ];

  function przeniesOgnisko(kodOkna: string): void {
    for (const okno of wszystkie) {
      if (okno.dataset['okno'] === kodOkna) okno.dataset['ognisko'] = 'tak';
      else delete okno.dataset['ognisko'];
    }
    const cel = wszystkie.find((okno) => okno.dataset['okno'] === kodOkna);
    cel?.scrollIntoView({ block: 'nearest' });
  }

  // Okno puste na starcie: panele mówią wtedy wprost, że rdzeń nie dał jeszcze
  // modułowi okna.
  const pomocnicze: PasPomocniczych = utworzPasPomocniczych({
    kanal,
    modul: KOD_MODULU,
    nazwaModulu: NAZWA_MODULU,
    okno: '',
    przedrostek: 'mp',
  });

  const pasWarsztatow = document.createElement('div');
  pasWarsztatow.className = 'mp-modul__pas mp-modul__pas--warsztaty';
  pasWarsztatow.append(frontend.element, backend.element);

  // Strona dystrybucji stoi w osobnym pasie, bo to inna strona cyklu życia
  // oprogramowania: pierwsza buduje produkt, druga rozdaje i konsumuje gotowe
  // rozszerzenia. Panele boczne (uprawnienia, konsola, wydawca) stoją obok
  // siebie, bo każdy otwiera się na pozycji wskazanej w oknie wiodącym.
  const pasRejestru = document.createElement('div');
  pasRejestru.className = 'mp-modul__pas mp-modul__pas--rejestr';
  pasRejestru.append(katalog.element, zainstalowane.element);

  const pasPaneli = document.createElement('div');
  pasPaneli.className = 'mp-modul__pas mp-modul__pas--panele';
  pasPaneli.append(uprawnienia.element, konsola.element, wydawca.element);

  const naglowekDystrybucji = document.createElement('p');
  naglowekDystrybucji.className = 'mp-modul__rozdzial';
  naglowekDystrybucji.textContent = 'Dystrybucja i konsumpcja';

  const element = document.createElement('div');
  element.className = 'mp-modul';
  element.dataset['modul'] = KOD_MODULU;
  element.setAttribute('aria-label', 'Moduł Apps — okna operacyjne');
  element.append(
    akcje.element,
    builder.element,
    architektura.element,
    pasWarsztatow,
    wdrozenia.element,
    naglowekDystrybucji,
    pasRejestru,
    integracje.element,
    pasPaneli,
    pomocnicze.element,
  );

  const odsubskrybuj = stan.obserwuj(() => odswiezOkna());
  // Druga subskrypcja, bo drugi stan: zmiana katalogu rozszerzeń nie dotyczy
  // okien strony budowy i odwrotnie. Wspólne przerysowanie kosztowałoby
  // przebieg wszystkich jedenastu okien przy każdej ramce jednego obszaru.
  const odsubskrybujRejestr = rejestr.obserwuj(() => odswiezOknaRejestru());

  function odswiezOknaRejestru(): void {
    katalog.odswiez();
    zainstalowane.odswiez();
    integracje.odswiez();
    uprawnienia.odswiez();
    konsola.odswiez();
  }

  function odswiezOkna(): void {
    builder.odswiez();
    architektura.odswiez();
    frontend.odswiez();
    backend.odswiez();
    wdrozenia.odswiez();
    // Dziennik wydań Publisher Panelu czyta wdrożenia, więc należy do tej strony.
    wydawca.odswiez();
    // Okno modułu przychodzi z rdzenia dopiero po `window.list`, więc pas
    // dostaje je tutaj — przy każdej zmianie stanu, bo to samo okno niczego nie
    // przebudowuje, a na pustym panele mówią wprost, czego brakuje.
    pomocnicze.ustawOkno(stan.idOkna());
    pomocnicze.odswiez();
  }

  return {
    element,

    async wczytaj(idSesji) {
      // Katalog akcji i wykaz okien dotyczą różnych obszarów kontraktu i żaden
      // nie warunkuje drugiego; odmowa jednego zostaje w jego miejscu.
      akcje.wczytaj();
      // Wykaz komend rdzenia wypełnia powody kontrolek strony dystrybucji.
      // Idzie raz na moduł i nie czeka na okno — powitanie dotyczy połączenia,
      // nie sesji.
      void pokrycie.odczytaj();
      // Katalog rozszerzeń i punkty dostępu nie wymagają okna modułu, więc
      // ruszają od razu, równolegle z odczytem okna. Odmowa któregokolwiek
      // zostaje w jego oknie i nie zabiera strony budowy.
      const rejestrGotowy = Promise.all([rejestr.odczytajKatalog(), rejestr.odczytajPunkty()]);
      await stan.odswiez(idSesji);
      // Odczyt okna modułu nie musi zakończyć się powiadomieniem obserwatorów
      // (odmowa `window.list` zostawia stan bez zmiany), a pas ma wtedy
      // powiedzieć o braku okna, nie milczeć.
      pomocnicze.ustawOkno(stan.idOkna());
      // Odczyty obszaru idą zaraz po ustaleniu okna, bo każdy wymaga
      // `windowId` — stąd miejsce po `stan.odswiez`. Idą równolegle, bo są od
      // siebie niezależne: odmowa jednego nie zabiera pozostałych. Wynik
      // każdego ląduje w stanie modułu, a okna dowiadują się o nim obserwacją.
      await Promise.all([
        stan.odczytajWdrozenia('', 0),
        stan.odczytajArchitekture(),
        stan.odczytajWarsztat(AppWorkspaceLayer.Frontend),
        stan.odczytajWarsztat(AppWorkspaceLayer.Backend),
        rejestrGotowy,
      ]);
    },

    rozlacz() {
      odsubskrybuj();
      odsubskrybujRejestr();
      // Wykaz pokrycia jest wspólny wszystkim modułom; odpinamy z niego
      // kontrolki tego modułu, żeby nie przerysowywał martwych elementów.
      pokrycie.zamknij();
      rejestr.rozlacz();
      // Pas zamyka się pierwszy: jego panele trzymają subskrypcje żyjące
      // niezależnie od stanu modułu (podgląd w tle — `stream.chunk`, Terminal —
      // zdarzenia karty powłoki) i po zejściu ze sceny nikt by ich nie zdjął.
      pomocnicze.zamknij();
      stan.rozlacz();
    },
  };
}
