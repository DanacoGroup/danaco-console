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
 * Moduł Apps zestawia pięć okien operacyjnych w jednym układzie: Product
 * Builder jako wejście integrujące pozostałe, Architecture Designer przed
 * warsztatami, dwa warsztaty równolegle oraz Deployment Panel na końcu.
 */
export interface ModulApps {
  /** Element osadzany w obszarze roboczym powłoki. */
  element: HTMLElement;
  /** Zleca odczyt okna modułu i katalogu akcji rdzenia. */
  wczytaj(idSesji: string): Promise<void>;
  /** Odłącza subskrypcje zdarzeń rdzenia. */
  rozlacz(): void;
}

/** Nazwa modułu używana w etykietach dostępności pasa okien pomocniczych osadzonego wewnątrz tego modułu. */
const NAZWA_MODULU = 'Apps';

/** Etapy procesu w kolejności wyświetlania na osi etapów: architektura, warsztaty, a następnie wdrożenie. */
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
  // Dwa etapy strony dystrybucji: zbudowane rozszerzenie wraca do katalogu, z którego korzystają moduły.
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
  // Drugi stan, bo druga strona modułu stoi na innym obszarze kontraktu i nie zna okna.
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

  // Strona dystrybucji stoi w osobnym pasie: pierwsza buduje produkt, druga rozdaje gotowe rozszerzenia.
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
  // Druga subskrypcja, bo drugi stan: zmiana katalogu rozszerzeń nie dotyczy okien strony budowy.
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
    // Okno modułu przychodzi z rdzenia dopiero po odczycie wykazu okien, więc pas dostaje je tutaj.
    pomocnicze.ustawOkno(stan.idOkna());
    pomocnicze.odswiez();
  }

  return {
    element,

    async wczytaj(idSesji) {
      // Katalog akcji i wykaz okien dotyczą różnych obszarów; odmowa jednego zostaje w jego miejscu.
      akcje.wczytaj();
      // Wykaz komend rdzenia idzie raz na moduł i nie czeka na okno — powitanie dotyczy połączenia.
      void pokrycie.odczytaj();
      // Katalog rozszerzeń i punkty dostępu nie wymagają okna modułu, ruszają więc od razu.
      const rejestrGotowy = Promise.all([rejestr.odczytajKatalog(), rejestr.odczytajPunkty()]);
      await stan.odswiez(idSesji);
      // Odczyt okna modułu nie musi zakończyć się powiadomieniem obserwatorów; pas ma powiedzieć o braku.
      pomocnicze.ustawOkno(stan.idOkna());
      // Odczyty obszaru idą zaraz po ustaleniu okna, równolegle, bo są od siebie niezależne.
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
      // Wykaz pokrycia jest wspólny wszystkim modułom; odpinamy z niego kontrolki tego modułu.
      pokrycie.zamknij();
      rejestr.rozlacz();
      // Pas zamyka się pierwszy: jego panele trzymają subskrypcje żyjące niezależnie od stanu modułu.
      pomocnicze.zamknij();
      stan.rozlacz();
    },
  };
}
