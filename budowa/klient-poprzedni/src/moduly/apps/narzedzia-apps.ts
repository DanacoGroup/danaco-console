import {
  AppDeployEnvironment,
  Command,
  AppEndpointMethod,
  AppExportFormat,
  AppMilestoneStatus,
  AppPackageFormat,
  AppPackageVisibility,
  AppProductPlatform,
  AppStageStatus,
  AppWorkspaceLayer,
  ExtensionKind,
  type AppPackageManifest,
} from '../../../../shared/contract';
import {
  bladOdmowyNarzedzia,
  jsonZPola,
  liczbaZPola,
  wymagajOknaModulu,
  wymagajPola,
  type NarzedzieApps,
} from './przybornik-apps';
import type { StanProduktu } from './stan-produktu';

// Narzędzia przyborników modułu Apps: po jednym na każdą komendę obszaru.

/** Skrót do okna modułu — funkcja pobiera identyfikator z bieżącego stanu produktu; każda komenda obszaru go wymaga w żądaniu. */
function okno(stan: StanProduktu): string {
  return wymagajOknaModulu(stan.idOkna());
}

/** Narzędzia Product Buildera modułu Apps: metadane produktu, powiązania między bytami, etapy realizacji, kamienie milowe i oś czasu. */
export function narzedziaProductBuilder(stan: StanProduktu): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Metadane produktu — odczyt',
      komenda: Command.AppsProductGet,
      async wykonaj() {
        const wynik = await stan.zrodlo.produkt(okno(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt produktu', wynik.blad);
        }
        const produkt = wynik.wynik.product;
        if (produkt === undefined) {
          return 'Okno nie ma jeszcze zapisanego produktu — to stan prawdziwy, nie odmowa.';
        }
        return `Produkt „${produkt.name}" (${produkt.id}), platform: ${produkt.platforms?.length ?? 0}.`;
      },
    },
    {
      etykieta: 'Metadane produktu — zapis',
      komenda: Command.AppsProductSave,
      pola: [
        { klucz: 'nazwa', etykieta: 'Nazwa produktu', wartosc: 'Portal klienta' },
        { klucz: 'opis', etykieta: 'Opis' },
        { klucz: 'platformy', etykieta: 'Platformy (web,mobile,desktop)', wartosc: 'web' },
        { klucz: 'repozytorium', etykieta: 'Repozytorium źródłowe' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.zapiszProdukt({
          idOkna: okno(stan),
          nazwa: wymagajPola(wartosci['nazwa'] ?? '', 'Nazwa produktu'),
          opis: (wartosci['opis'] ?? '').trim(),
          platformy: platformyZPola(wartosci['platformy'] ?? ''),
          repozytorium: (wartosci['repozytorium'] ?? '').trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis produktu', wynik.blad);
        }
        return `Zapisano produkt ${wynik.wynik.product.id} („${wynik.wynik.product.name}").`;
      },
    },
    {
      etykieta: 'Powiązania z modułami',
      komenda: Command.AppsProductLinkList,
      async wykonaj() {
        const wynik = await stan.zrodlo.powiazaniaProduktu(okno(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt powiązań produktu', wynik.blad);
        }
        const czynne = wynik.wynik.links.filter((powiazanie) => powiazanie.enabled);
        if (czynne.length === 0) {
          return `Żadne z ${wynik.wynik.links.length} powiązań nie jest jeszcze czynne — powiązanie staje się czynne, gdy w oknie pojawi się to, czym żyje.`;
        }
        return `Czynne powiązania (${czynne.length} z ${wynik.wynik.links.length}): ${czynne
          .map((powiazanie) => `${powiazanie.moduleCode} — ${powiazanie.detail ?? 'bez szczegółu'}`)
          .join('; ')}.`;
      },
    },
    {
      etykieta: 'Etapy budowy — odczyt',
      komenda: Command.AppsStageList,
      pola: [{ klucz: 'stan', etykieta: 'Stan (pending/active/done/blocked)' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.etapy(
          okno(stan),
          stanEtapuZPola(wartosci['stan'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt etapów', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'Okno nie ma jeszcze ani jednego etapu budowy.';
        return `Etapów: ${wynik.wynik.total} — ${wynik.wynik.stages
          .map((etap) => `${etap.name} (${etap.status})`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Etap budowy — zapis',
      komenda: Command.AppsStageSave,
      pola: [
        { klucz: 'idEtapu', etykieta: 'Etap zmieniany (pusty zakłada nowy)' },
        { klucz: 'nazwa', etykieta: 'Nazwa etapu', wartosc: 'Architektura' },
        { klucz: 'kolejnosc', etykieta: 'Kolejność' },
        { klucz: 'stan', etykieta: 'Stan etapu' },
        { klucz: 'wykonawca', etykieta: 'Wykonawca' },
      ],
      async wykonaj(wartosci) {
        // Pole wykonawcy: puste przy etapie nowym nie zmienia przypisania; pusty ciąg przy zmianie usuwa.
        const idEtapu = (wartosci['idEtapu'] ?? '').trim();
        const wykonawca = wartosci['wykonawca'] ?? '';
        const wynik = await stan.zrodlo.zapiszEtap({
          idOkna: okno(stan),
          idEtapu,
          nazwa: (wartosci['nazwa'] ?? '').trim(),
          kolejnosc: liczbaZPola(wartosci['kolejnosc'] ?? ''),
          stan: stanEtapuZPola(wartosci['stan'] ?? ''),
          wykonawca: idEtapu === '' && wykonawca.trim() === '' ? null : wykonawca.trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis etapu', wynik.blad);
        }
        const etap = wynik.wynik.stage;
        return `Etap ${etap.id} („${etap.name}") w stanie ${etap.status}, wykonawca: ${
          etap.ownerAgentId ?? 'bez przypisania'
        }.`;
      },
    },
    {
      etykieta: 'Kamienie milowe — odczyt',
      komenda: Command.AppsMilestoneList,
      pola: [{ klucz: 'stan', etykieta: 'Stan (planned/active/reached/missed)' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.kamienieMilowe(
          okno(stan),
          stanKamieniaZPola(wartosci['stan'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt kamieni milowych', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'Okno nie ma jeszcze ani jednego kamienia milowego.';
        return `Kamieni milowych: ${wynik.wynik.total} — ${wynik.wynik.milestones
          .map((kamien) => `${kamien.name} (${kamien.status})`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Kamień milowy — zapis',
      komenda: Command.AppsMilestoneSave,
      pola: [
        { klucz: 'idKamienia', etykieta: 'Kamień zmieniany (pusty zakłada nowy)' },
        { klucz: 'nazwa', etykieta: 'Nazwa kamienia', wartosc: 'Wydanie 1.0' },
        { klucz: 'termin', etykieta: 'Termin (ms epoki)' },
        { klucz: 'stan', etykieta: 'Stan realizacji' },
        { klucz: 'etapy', etykieta: 'Etapy (identyfikatory po przecinku)' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.zapiszKamien({
          idOkna: okno(stan),
          idKamienia: (wartosci['idKamienia'] ?? '').trim(),
          nazwa: wymagajPola(wartosci['nazwa'] ?? '', 'Nazwa kamienia'),
          termin: liczbaZPola(wartosci['termin'] ?? ''),
          stan: stanKamieniaZPola(wartosci['stan'] ?? ''),
          idEtapow: rozdzielPrzecinkami(wartosci['etapy'] ?? ''),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis kamienia milowego', wynik.blad);
        }
        const kamien = wynik.wynik.milestone;
        return `Kamień ${kamien.id} („${kamien.name}"), etapów: ${kamien.stageIds?.length ?? 0}.`;
      },
    },
    {
      etykieta: 'Kamień milowy — usunięcie',
      komenda: Command.AppsMilestoneDelete,
      pola: [{ klucz: 'idKamienia', etykieta: 'Kamień usuwany' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.usunKamien(
          okno(stan),
          wymagajPola(wartosci['idKamienia'] ?? '', 'Kamień usuwany'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Usunięcie kamienia milowego', wynik.blad);
        }
        return wynik.wynik.deleted
          ? 'Kamień milowy usunięty wraz z jego związkiem z etapami.'
          : 'Rdzeń zameldował, że nie było czego usunąć.';
      },
    },
    {
      etykieta: 'Oś czasu projektu',
      komenda: Command.AppsTimelineList,
      pola: [{ klucz: 'granica', etykieta: 'Górna granica pozycji' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.osCzasu(
          okno(stan),
          liczbaZPola(wartosci['granica'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt osi czasu', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'W tym oknie nic jeszcze nie zaszło.';
        const najnowsze = wynik.wynik.entries[0];
        return `Zdarzeń: ${wynik.wynik.total}; najnowsze — ${najnowsze?.summary ?? 'bez opisu'}.`;
      },
    },
  ];
}

/** Narzędzia Architecture Designera modułu Apps: walidacja układu, zapis i odczyt wersji, adnotacje oraz eksport do pliku. */
export function narzedziaArchitectureDesigner(stan: StanProduktu): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Waliduj układ',
      komenda: Command.AppsArchitectureValidate,
      async wykonaj() {
        const wynik = await stan.zrodlo.sprawdzArchitekture(okno(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Walidacja architektury', wynik.blad);
        }
        const zastrzezenia = wynik.wynik.issues;
        if (zastrzezenia.length === 0) return 'Układ bez zastrzeżeń.';
        // Zastrzeżenie nie blokuje niczego; zdanie podaje, co znaleziono, nie że nie można iść dalej.
        return `Zastrzeżeń: ${zastrzezenia.length} (ostrzeżenia, nie brama) — ${zastrzezenia
          .map((zastrzezenie) => `${zastrzezenie.severity}: ${zastrzezenie.message}`)
          .join('; ')}.`;
      },
    },
    {
      etykieta: 'Historia wersji układu',
      komenda: Command.AppsArchitectureVersionList,
      pola: [{ klucz: 'granica', etykieta: 'Górna granica wersji' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.wersjeArchitektury(
          okno(stan),
          liczbaZPola(wartosci['granica'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt historii wersji', wynik.blad);
        }
        return `Wersji: ${wynik.wynik.total} — ${wynik.wynik.versions
          .map((wersja) => `v${wersja.version} (${wersja.componentCount} komponentów)`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Adnotacja projektowa',
      komenda: Command.AppsArchitectureAnnotationSave,
      pola: [
        { klucz: 'idAdnotacji', etykieta: 'Adnotacja zmieniana (pusta zakłada nową)' },
        { klucz: 'komponent', etykieta: 'Komponent' },
        { klucz: 'zaleznoscZ', etykieta: 'Zależność — komponent źródłowy' },
        { klucz: 'zaleznoscDo', etykieta: 'Zależność — komponent docelowy' },
        { klucz: 'tresc', etykieta: 'Treść notatki', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.zapiszAdnotacje({
          idOkna: okno(stan),
          idAdnotacji: (wartosci['idAdnotacji'] ?? '').trim(),
          idKomponentu: (wartosci['komponent'] ?? '').trim(),
          zaleznoscZ: (wartosci['zaleznoscZ'] ?? '').trim(),
          zaleznoscDo: (wartosci['zaleznoscDo'] ?? '').trim(),
          tresc: wymagajPola(wartosci['tresc'] ?? '', 'Treść notatki'),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis adnotacji', wynik.blad);
        }
        return `Adnotacja ${wynik.wynik.annotation.id} zapisana.`;
      },
    },
    {
      etykieta: 'Eksport diagramu',
      komenda: Command.AppsArchitectureExport,
      pola: [{ klucz: 'format', etykieta: 'Format (svg/png/mermaid/markdown)', wartosc: 'mermaid' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.wyeksportujArchitekture(
          okno(stan),
          formatEksportuZPola(wartosci['format'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Eksport diagramu', wynik.blad);
        }
        return `Plik w magazynie rdzenia: ${wynik.wynik.artifactRef} (${wynik.wynik.sizeBytes} bajtów).`;
      },
    },
  ];
}

/**
 * Narzędzia warsztatu. Warstwa rozstrzyga, co okno pokazuje: podgląd i mapa
 * routingu należą do frontendu, eksplorator punktów końcowych i schemat bazy —
 * do backendu. Motyw stoi przy frontendzie, bo tam go Operator edytuje.
 */
export function narzedziaWarsztatu(
  stan: StanProduktu,
  warstwa: AppWorkspaceLayer,
): readonly NarzedzieApps[] {
  const wspolne: NarzedzieApps[] = [
    {
      etykieta: 'Podgląd na żywo — uruchom',
      komenda: Command.AppsPreviewStart,
      async wykonaj() {
        const wynik = await stan.zrodlo.uruchomPodglad(okno(stan), warstwa);
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Uruchomienie podglądu', wynik.blad);
        }
        return `Podgląd (${wynik.wynik.status}) stoi pod ${wynik.wynik.previewUrl}.`;
      },
    },
    {
      etykieta: 'Podgląd na żywo — zatrzymaj',
      komenda: Command.AppsPreviewStop,
      async wykonaj() {
        const wynik = await stan.zrodlo.zatrzymajPodglad(okno(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zatrzymanie podglądu', wynik.blad);
        }
        return wynik.wynik.stopped
          ? 'Serwer podglądu zatrzymany, nasłuch zamknięty.'
          : 'Rdzeń zameldował, że nie było czego zatrzymywać.';
      },
    },
    {
      etykieta: 'Dziennik usług',
      komenda: Command.AppsServiceLogRead,
      pola: [
        { klucz: 'komponent', etykieta: 'Komponent' },
        { klucz: 'granica', etykieta: 'Górna granica wierszy' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.dziennikUslugi(
          okno(stan),
          (wartosci['komponent'] ?? '').trim(),
          liczbaZPola(wartosci['granica'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt dziennika usług', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'W tym oknie nic jeszcze nie zapisało się w dzienniku.';
        return `Wierszy: ${wynik.wynik.total}; ostatni — ${
          wynik.wynik.lines[wynik.wynik.lines.length - 1] ?? ''
        }.`;
      },
    },
  ];

  if (warstwa === AppWorkspaceLayer.Frontend) {
    return [
      ...wspolne,
      {
        etykieta: 'Mapa routingu',
        komenda: Command.AppsRouteList,
        async wykonaj() {
          const wynik = await stan.zrodlo.trasy(okno(stan));
          if (!wynik.udany || wynik.wynik === undefined) {
            throw bladOdmowyNarzedzia('Odczyt mapy routingu', wynik.blad);
          }
          if (wynik.wynik.total === 0) {
            return 'W plikach warstwy interfejsu nie ma jeszcze ani jednej deklaracji trasy.';
          }
          return `Tras: ${wynik.wynik.total} — ${wynik.wynik.routes
            .map((trasa) => `${trasa.path}${trasa.viewName === undefined ? '' : ` → ${trasa.viewName}`}`)
            .join(', ')}.`;
        },
      },
      {
        etykieta: 'Motyw produktu — odczyt',
        komenda: Command.AppsThemeGet,
        async wykonaj() {
          const wynik = await stan.zrodlo.motyw(okno(stan));
          if (!wynik.udany || wynik.wynik === undefined) {
            throw bladOdmowyNarzedzia('Odczyt motywu', wynik.blad);
          }
          if (wynik.wynik.theme === undefined) return 'Produkt nie ma jeszcze zapisanego motywu.';
          return `Motyw: ${JSON.stringify(wynik.wynik.theme)}`;
        },
      },
      {
        etykieta: 'Motyw produktu — zapis',
        komenda: Command.AppsThemeSet,
        pola: [
          {
            klucz: 'motyw',
            etykieta: 'Motyw (JSON)',
            wartosc: '{"kolorGlowny":"#c8a24a"}',
            obszerne: true,
          },
        ],
        async wykonaj(wartosci) {
          const motyw = jsonZPola(wartosci['motyw'] ?? '', 'Motyw (JSON)');
          if (motyw === undefined) {
            throw new Error('Pole „Motyw (JSON)" jest wymagane przez kontrakt — żądanie nie poszło.');
          }
          const wynik = await stan.zrodlo.ustawMotyw(okno(stan), motyw);
          if (!wynik.udany || wynik.wynik === undefined) {
            throw bladOdmowyNarzedzia('Zapis motywu', wynik.blad);
          }
          return 'Motyw produktu zapisany.';
        },
      },
    ];
  }

  return [
    ...wspolne,
    {
      etykieta: 'Eksplorator punktów końcowych',
      komenda: Command.AppsEndpointList,
      pola: [{ klucz: 'komponent', etykieta: 'Komponent' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.punktyKoncowe(
          okno(stan),
          (wartosci['komponent'] ?? '').trim(),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt punktów końcowych', wynik.blad);
        }
        if (wynik.wynik.total === 0) {
          return 'Żaden komponent architektury nie ma jeszcze wypełnionego kontraktu API.';
        }
        return `Punktów: ${wynik.wynik.total} — ${wynik.wynik.endpoints
          .map((punkt) => `${punkt.method.toUpperCase()} ${punkt.path} (${punkt.status ?? 'bez stanu'})`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Zapytanie testowe',
      komenda: Command.AppsEndpointProbe,
      pola: [
        { klucz: 'metoda', etykieta: 'Metoda', wartosc: 'get' },
        { klucz: 'sciezka', etykieta: 'Ścieżka', wartosc: '/zamowienia' },
        { klucz: 'srodowisko', etykieta: 'Środowisko' },
        { klucz: 'naglowki', etykieta: 'Nagłówki (JSON)' },
        { klucz: 'tresc', etykieta: 'Treść zapytania', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.zapytajPunkt({
          idOkna: okno(stan),
          metoda: metodaZPola(wartosci['metoda'] ?? ''),
          sciezka: wymagajPola(wartosci['sciezka'] ?? '', 'Ścieżka'),
          naglowki: jsonZPola(wartosci['naglowki'] ?? '', 'Nagłówki (JSON)'),
          tresc: (wartosci['tresc'] ?? '').trim(),
          srodowisko: srodowiskoZPola(wartosci['srodowisko'] ?? ''),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapytanie testowe', wynik.blad);
        }
        const rezultat = wynik.wynik.result;
        if (rezultat.errorDetail !== undefined) {
          return `Usługa nie odpowiedziała po ${rezultat.durationMs} ms: ${rezultat.errorDetail}`;
        }
        return `${rezultat.statusCode} w ${rezultat.durationMs} ms.`;
      },
    },
    {
      etykieta: 'Schemat bazy produktu',
      komenda: Command.AppsSchemaGet,
      pola: [{ klucz: 'komponent', etykieta: 'Komponent bazy' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.schemat(
          okno(stan),
          (wartosci['komponent'] ?? '').trim(),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt schematu bazy', wynik.blad);
        }
        const schemat = wynik.wynik.schema;
        if (schemat === undefined) {
          return 'W plikach warstwy backendu nie ma jeszcze ani jednego polecenia CREATE TABLE.';
        }
        return `Tabel: ${schemat.tables.length} — ${schemat.tables
          .map((tabela) => `${tabela.name} (${tabela.columns.length} kolumn)`)
          .join(', ')}.`;
      },
    },
  ];
}

/** Narzędzia Deployment Panelu modułu Apps: środowiska wdrożeniowe, zmienne konfiguracyjne, domena, skalowanie oraz kondycja usługi. */
export function narzedziaDeploymentPanel(stan: StanProduktu): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Środowiska wdrożeniowe',
      komenda: Command.AppsEnvironmentList,
      async wykonaj() {
        const wynik = await stan.zrodlo.srodowiska(okno(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt środowisk', wynik.blad);
        }
        return `Środowisk: ${wynik.wynik.total} — ${wynik.wynik.environments
          .map((srodowisko) => `${srodowisko.code}${srodowisko.domain === undefined ? '' : ` (${srodowisko.domain})`}`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Zmienne środowiska — odczyt',
      komenda: Command.AppsEnvironmentVariableList,
      pola: [{ klucz: 'srodowisko', etykieta: 'Środowisko', wartosc: 'dev' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.zmienneSrodowiska(
          okno(stan),
          wymaganeSrodowiskoZPola(wartosci['srodowisko'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt zmiennych środowiska', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'To środowisko nie ma jeszcze ani jednej zmiennej.';
        return `Zmiennych: ${wynik.wynik.total} — ${wynik.wynik.variables
          .map((zmienna) => `${zmienna.name}${zmienna.secretRef === undefined ? '' : ' (sekret)'}`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Zmienna środowiska — zapis',
      komenda: Command.AppsEnvironmentVariableSet,
      pola: [
        { klucz: 'srodowisko', etykieta: 'Środowisko', wartosc: 'dev' },
        { klucz: 'nazwa', etykieta: 'Nazwa zmiennej' },
        { klucz: 'wartosc', etykieta: 'Wartość jawna' },
        { klucz: 'sekret', etykieta: 'Odwołanie do sekretu' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.ustawZmienna({
          idOkna: okno(stan),
          srodowisko: wymaganeSrodowiskoZPola(wartosci['srodowisko'] ?? ''),
          nazwa: wymagajPola(wartosci['nazwa'] ?? '', 'Nazwa zmiennej'),
          wartosc: (wartosci['wartosc'] ?? '').trim(),
          odwolanieSekretu: (wartosci['sekret'] ?? '').trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis zmiennej środowiska', wynik.blad);
        }
        const zmienna = wynik.wynik.variable;
        return `Zmienna ${zmienna.name} zapisana${
          zmienna.secretRef === undefined ? '' : ' jako odwołanie do sekretu'
        }.`;
      },
    },
    {
      etykieta: 'Domena i DNS',
      komenda: Command.AppsDeploymentDomainSet,
      pola: [
        { klucz: 'srodowisko', etykieta: 'Środowisko', wartosc: 'production' },
        { klucz: 'domena', etykieta: 'Domena produktu' },
        { klucz: 'dns', etykieta: 'Wpisy DNS (JSON)', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.ustawDomene({
          idOkna: okno(stan),
          srodowisko: wymaganeSrodowiskoZPola(wartosci['srodowisko'] ?? ''),
          domena: wymagajPola(wartosci['domena'] ?? '', 'Domena produktu'),
          wpisyDns: jsonZPola(wartosci['dns'] ?? '', 'Wpisy DNS (JSON)'),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Nadanie domeny', wynik.blad);
        }
        return `Domena ${wynik.wynik.domain} zapisana; nazwa ${
          wynik.wynik.verified ? 'rozwiązuje się' : 'jeszcze się nie rozwiązuje'
        }.`;
      },
    },
    {
      etykieta: 'Skalowanie usługi',
      komenda: Command.AppsDeploymentScaleSet,
      pola: [
        { klucz: 'srodowisko', etykieta: 'Środowisko', wartosc: 'production' },
        { klucz: 'instancje', etykieta: 'Instancje stałe' },
        { klucz: 'minimum', etykieta: 'Dolna granica' },
        { klucz: 'maksimum', etykieta: 'Górna granica' },
        { klucz: 'reguly', etykieta: 'Reguły (JSON)', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.ustawSkalowanie({
          idOkna: okno(stan),
          srodowisko: wymaganeSrodowiskoZPola(wartosci['srodowisko'] ?? ''),
          instancje: liczbaZPola(wartosci['instancje'] ?? ''),
          minimum: liczbaZPola(wartosci['minimum'] ?? ''),
          maksimum: liczbaZPola(wartosci['maksimum'] ?? ''),
          reguly: jsonZPola(wartosci['reguly'] ?? '', 'Reguły (JSON)'),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Nastawa skalowania', wynik.blad);
        }
        return `Nastawa zapisana; instancji obowiązujących: ${
          wynik.wynik.effectiveInstances ?? 'bez wskazania'
        }.`;
      },
    },
    {
      etykieta: 'Kondycja produktu',
      komenda: Command.AppsDeploymentHealthGet,
      pola: [{ klucz: 'srodowisko', etykieta: 'Środowisko' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.kondycja(
          okno(stan),
          srodowiskoZPola(wartosci['srodowisko'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Sprawdzenie kondycji', wynik.blad);
        }
        const kondycja = wynik.wynik.health;
        return `${kondycja.environment}: ${
          kondycja.available ? 'usługa odpowiada' : 'usługa nie odpowiada'
        }, dostępność ${kondycja.availabilityPercent ?? 'bez pomiaru'}% — ${
          kondycja.detail ?? 'bez szczegółu'
        }.`;
      },
    },
    {
      etykieta: 'Dziennik wdrożenia',
      komenda: Command.AppsDeploymentLogRead,
      pola: [
        { klucz: 'wdrozenie', etykieta: 'Wdrożenie' },
        { klucz: 'granica', etykieta: 'Górna granica wierszy' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.dziennikWdrozenia(
          okno(stan),
          wymagajPola(wartosci['wdrozenie'] ?? '', 'Wdrożenie'),
          liczbaZPola(wartosci['granica'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt dziennika wdrożenia', wynik.blad);
        }
        return `Wierszy: ${wynik.wynik.total}; ostatni — ${
          wynik.wynik.lines[wynik.wynik.lines.length - 1] ?? 'brak'
        }.`;
      },
    },
    {
      etykieta: 'Artefakty budowania',
      komenda: Command.AppsArtifactList,
      pola: [{ klucz: 'wdrozenie', etykieta: 'Wdrożenie' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.artefakty(
          okno(stan),
          (wartosci['wdrozenie'] ?? '').trim(),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt artefaktów', wynik.blad);
        }
        if (wynik.wynik.total === 0) {
          return 'Okno nie ma jeszcze artefaktów — powstają po udanym wdrożeniu.';
        }
        return `Artefaktów: ${wynik.wynik.total} — ${wynik.wynik.artifacts
          .map((artefakt) => `${artefakt.kind} ${artefakt.sizeBytes ?? 0} B`)
          .join(', ')}.`;
      },
    },
  ];
}

/** Narzędzia Publisher Panelu modułu Apps: pakietowanie, manifest, walidacja pakietu, podpis cyfrowy oraz publikacja do kanału. */
export function narzedziaPublisherPanel(stan: StanProduktu): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Zbuduj pakiet',
      komenda: Command.AppsPackageBuild,
      pola: [
        { klucz: 'artefakt', etykieta: 'Artefakt (pusty bierze ostatni udany)' },
        { klucz: 'format', etykieta: 'Format (zip/targz)', wartosc: 'zip' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.zbudujPakiet(
          okno(stan),
          (wartosci['artefakt'] ?? '').trim(),
          formatPakietuZPola(wartosci['format'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Budowa pakietu', wynik.blad);
        }
        const pakiet = wynik.wynik.package;
        return `Pakiet ${pakiet.id} (${pakiet.format}, ${pakiet.sizeBytes ?? 0} bajtów).`;
      },
    },
    {
      etykieta: 'Zapisz manifest',
      komenda: Command.AppsPackageManifestSave,
      pola: [
        { klucz: 'pakiet', etykieta: 'Pakiet (pusty zakłada nowy)' },
        { klucz: 'identyfikator', etykieta: 'Identyfikator', wartosc: 'portal-klienta' },
        { klucz: 'nazwa', etykieta: 'Nazwa', wartosc: 'Portal klienta' },
        { klucz: 'wersja', etykieta: 'Wersja semantyczna', wartosc: '1.0.0' },
        { klucz: 'rodzaj', etykieta: 'Rodzaj (mcp/plugin/api/skill)', wartosc: 'plugin' },
        { klucz: 'opis', etykieta: 'Opis' },
      ],
      async wykonaj(wartosci) {
        const manifest: AppPackageManifest = {
          identifier: wymagajPola(wartosci['identyfikator'] ?? '', 'Identyfikator'),
          name: wymagajPola(wartosci['nazwa'] ?? '', 'Nazwa'),
          version: wymagajPola(wartosci['wersja'] ?? '', 'Wersja semantyczna'),
          kind: rodzajRozszerzeniaZPola(wartosci['rodzaj'] ?? ''),
        };
        const opis = (wartosci['opis'] ?? '').trim();
        if (opis !== '') manifest.description = opis;
        const wynik = await stan.zrodlo.zapiszManifest(
          okno(stan),
          (wartosci['pakiet'] ?? '').trim(),
          manifest,
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis manifestu', wynik.blad);
        }
        return `Manifest pakietu ${wynik.wynik.package.id} zapisany.`;
      },
    },
    {
      etykieta: 'Waliduj pakiet',
      komenda: Command.AppsPackageValidate,
      pola: [{ klucz: 'pakiet', etykieta: 'Pakiet sprawdzany' }],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.sprawdzPakiet(
          okno(stan),
          wymagajPola(wartosci['pakiet'] ?? '', 'Pakiet sprawdzany'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Walidacja pakietu', wynik.blad);
        }
        if (wynik.wynik.issues.length === 0) return 'Pakiet bez zastrzeżeń.';
        return `Zastrzeżeń: ${wynik.wynik.issues.length} (nie wstrzymują publikacji) — ${wynik.wynik.issues
          .map((zastrzezenie) => `${zastrzezenie.severity}: ${zastrzezenie.message}`)
          .join('; ')}.`;
      },
    },
    {
      etykieta: 'Podpisz pakiet',
      komenda: Command.AppsPackageSign,
      pola: [
        { klucz: 'pakiet', etykieta: 'Pakiet podpisywany' },
        { klucz: 'klucz', etykieta: 'Odwołanie do klucza wydawcy' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.podpiszPakiet(
          okno(stan),
          wymagajPola(wartosci['pakiet'] ?? '', 'Pakiet podpisywany'),
          // Do rdzenia idzie odwołanie, nie treść klucza; klucz leży w warstwie sekretów.
          wymagajPola(wartosci['klucz'] ?? '', 'Odwołanie do klucza wydawcy'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Podpisanie pakietu', wynik.blad);
        }
        const podpis = wynik.wynik.signature;
        return `Podpis ${podpis.algorithm ?? 'bez nazwy algorytmu'}, weryfikacja: ${
          podpis.verified ? 'przeszła' : 'nie przeszła'
        }, suma ${podpis.checksumSha256 ?? 'bez sumy'}.`;
      },
    },
    {
      etykieta: 'Opublikuj do rejestru',
      komenda: Command.AppsPackagePublish,
      pola: [
        { klucz: 'pakiet', etykieta: 'Pakiet publikowany' },
        { klucz: 'notatki', etykieta: 'Notatki wydania', obszerne: true },
        { klucz: 'widocznosc', etykieta: 'Widoczność (organization/restricted)' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.zrodlo.opublikujPakiet({
          idOkna: okno(stan),
          idPakietu: wymagajPola(wartosci['pakiet'] ?? '', 'Pakiet publikowany'),
          notatki: (wartosci['notatki'] ?? '').trim(),
          widocznosc: widocznoscZPola(wartosci['widocznosc'] ?? ''),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Publikacja pakietu', wynik.blad);
        }
        const pozycja = wynik.wynik.extension;
        return `Pozycja katalogu ${pozycja.code} (${pozycja.kind}) w wersji ${
          pozycja.version ?? 'bez wersji'
        }; stan wyjściowy: ${pozycja.enabled ? 'włączona' : 'wyłączona'}.`;
      },
    },
  ];
}

// ── Przekłady wartości pól na wartości kontraktu ─────────────────────────────

// Każdy przekład odmawia przy wartości spoza kontraktu; odmowa pada przy polu, które wypełniano.

/** Rozdziela wartości pola wpisane po przecinku na wykaz pozycji, pomijając pozycje puste powstałe z odstępów wokół przecinków. */
function rozdzielPrzecinkami(wartosc: string): readonly string[] {
  return wartosc
    .split(',')
    .map((czesc) => czesc.trim())
    .filter((czesc) => czesc !== '');
}

function platformyZPola(wartosc: string): readonly AppProductPlatform[] {
  const znane = Object.values(AppProductPlatform) as readonly string[];
  return rozdzielPrzecinkami(wartosc).map((czesc) => {
    if (!znane.includes(czesc)) {
      throw new Error(`Platforma „${czesc}" nie jest wartością kontraktu (${znane.join(', ')}).`);
    }
    return czesc as AppProductPlatform;
  });
}

function stanEtapuZPola(wartosc: string): AppStageStatus | '' {
  return wartoscWyliczenia(wartosc, AppStageStatus, 'stan etapu') as AppStageStatus | '';
}

function stanKamieniaZPola(wartosc: string): AppMilestoneStatus | '' {
  return wartoscWyliczenia(wartosc, AppMilestoneStatus, 'stan kamienia milowego') as
    | AppMilestoneStatus
    | '';
}

function formatEksportuZPola(wartosc: string): AppExportFormat {
  const wybrana = wartoscWyliczenia(wartosc, AppExportFormat, 'format eksportu');
  if (wybrana === '') {
    throw new Error('Pole „Format" jest wymagane przez kontrakt — żądanie nie poszło.');
  }
  return wybrana as AppExportFormat;
}

function metodaZPola(wartosc: string): AppEndpointMethod {
  const wybrana = wartoscWyliczenia(wartosc.toLowerCase(), AppEndpointMethod, 'metoda zapytania');
  if (wybrana === '') {
    throw new Error('Pole „Metoda" jest wymagane przez kontrakt — żądanie nie poszło.');
  }
  return wybrana as AppEndpointMethod;
}

function srodowiskoZPola(wartosc: string): AppDeployEnvironment | '' {
  return wartoscWyliczenia(wartosc, AppDeployEnvironment, 'środowisko wdrożenia') as
    | AppDeployEnvironment
    | '';
}

function wymaganeSrodowiskoZPola(wartosc: string): AppDeployEnvironment {
  const wybrane = srodowiskoZPola(wartosc);
  if (wybrane === '') {
    throw new Error('Pole „Środowisko" jest wymagane przez kontrakt — żądanie nie poszło.');
  }
  return wybrane;
}

function formatPakietuZPola(wartosc: string): AppPackageFormat | '' {
  return wartoscWyliczenia(wartosc, AppPackageFormat, 'format pakietu') as AppPackageFormat | '';
}

function widocznoscZPola(wartosc: string): AppPackageVisibility | '' {
  return wartoscWyliczenia(wartosc, AppPackageVisibility, 'widoczność pakietu') as
    | AppPackageVisibility
    | '';
}

function rodzajRozszerzeniaZPola(wartosc: string): ExtensionKind {
  const wybrany = wartoscWyliczenia(wartosc, ExtensionKind, 'rodzaj rozszerzenia');
  if (wybrany === '') {
    throw new Error('Pole „Rodzaj" jest wymagane przez kontrakt — żądanie nie poszło.');
  }
  return wybrany as ExtensionKind;
}

/**
 * Wspólny przekład pola na wartość wyliczenia kontraktu. Wykaz dopuszczalnych
 * wartości bierze się z wyliczenia, nie z listy wpisanej tutaj — dopisanie
 * wartości do kontraktu przepisuje odmowę samo.
 */
function wartoscWyliczenia(
  wartosc: string,
  wyliczenie: Record<string, string>,
  nazwa: string,
): string {
  const przyciete = wartosc.trim();
  if (przyciete === '') return '';
  const znane = Object.values(wyliczenie);
  if (!znane.includes(przyciete)) {
    throw new Error(
      `Wartość „${przyciete}" nie jest wartością kontraktu dla pola ${nazwa} (${znane.join(', ')}).`,
    );
  }
  return przyciete;
}
