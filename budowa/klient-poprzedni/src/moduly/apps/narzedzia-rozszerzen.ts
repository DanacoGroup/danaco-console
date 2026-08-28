import {
  Command,
  ExtensionAuthKind,
  ExtensionBulkAction,
  ExtensionDefinitionFormat,
  ExtensionKind,
  ExtensionOrigin,
  ExtensionPermissionScope,
  ExtensionToolKind,
  ExtensionWebhookDirection,
  McpTransport,
  type ExtensionPermission,
} from '../../../../shared/contract';
import {
  bladOdmowyNarzedzia,
  jsonZPola,
  liczbaZPola,
  wymagajPola,
  type NarzedzieApps,
} from './przybornik-apps';
import type { StanRozszerzen } from './stan-rozszerzen';

// Narzędzia przyborników strony rozszerzeń: po jednym na każdą komendę extension.*.

/** Pozycja wskazana w oknie stanu rozszerzeń, pobrana funkcją stan.wybrane(); brak wskazania kończy się odmową z podanym powodem. */
function wskazana(stan: StanRozszerzen): string {
  const pozycja = stan.wybrane();
  if (pozycja !== null) return pozycja.id;
  throw new Error(
    'Żadna pozycja katalogu nie jest wskazana — wybierz ją najpierw z listy okna, ' +
      'a potem powtórz czynność.',
  );
}

/** Narzędzia App Catalog modułu Apps: wyszukiwarka, karta szczegółów, kolekcje oraz rejestr organizacji. */
export function narzedziaAppCatalog(stan: StanRozszerzen): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Wyszukiwarka katalogu',
      komenda: Command.ExtensionSearch,
      pola: [
        { klucz: 'fraza', etykieta: 'Fraza' },
        { klucz: 'rodzaj', etykieta: 'Rodzaj (mcp/plugin/api/skill)' },
        { klucz: 'pochodzenie', etykieta: 'Źródło (danaco/personal)' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.szukaj({
          fraza: (wartosci['fraza'] ?? '').trim(),
          rodzaj: wartoscWyliczeniaRozszerzen(wartosci['rodzaj'] ?? '', ExtensionKind, 'rodzaj'),
          pochodzenie: wartoscWyliczeniaRozszerzen(
            wartosci['pochodzenie'] ?? '',
            ExtensionOrigin,
            'źródło',
          ),
          tylkoZainstalowane: false,
          granica: 0,
          odsuniecie: 0,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Wyszukiwanie w katalogu', wynik.blad);
        }
        if (wynik.wynik.total === 0) {
          return 'Żadna pozycja katalogu nie pasuje do zawężenia — to odpowiedź, nie brak.';
        }
        return `Trafień: ${wynik.wynik.total} — ${wynik.wynik.extensions
          .map((pozycja) => `${pozycja.name} (${pozycja.kind})`)
          .join(', ')}. Podpowiedzi: ${(wynik.wynik.suggestions ?? []).join(', ')}.`;
      },
    },
    {
      etykieta: 'Karta szczegółów pozycji',
      komenda: Command.ExtensionDetailGet,
      async wykonaj() {
        const wynik = await stan.dobudowa.szczegol(wskazana(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt szczegółów pozycji', wynik.blad);
        }
        const szczegol = wynik.wynik.detail;
        return `Narzędzi: ${szczegol.tools?.length ?? 0}, uprawnień deklarowanych: ${
          szczegol.permissions?.length ?? 0
        }, zależności: ${szczegol.dependencies?.length ?? 0}, podpis: ${
          szczegol.signature === undefined ? 'brak' : szczegol.signature.trustLevel
        }.`;
      },
    },
    {
      etykieta: 'Kolekcje kuratorskie',
      komenda: Command.ExtensionCollectionList,
      pola: [{ klucz: 'kolekcja', etykieta: 'Kolekcja (pusta: wszystkie)' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.kolekcje((wartosci['kolekcja'] ?? '').trim());
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt kolekcji', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'Nie ma jeszcze ani jednej kolekcji kuratorskiej.';
        return `Kolekcji: ${wynik.wynik.total} — ${wynik.wynik.collections
          .map((kolekcja) => `${kolekcja.name} (${kolekcja.extensionIds.length} pozycji)`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Kolekcja — zapis',
      komenda: Command.ExtensionCollectionSave,
      pola: [
        { klucz: 'kolekcja', etykieta: 'Kolekcja zmieniana (pusta zakłada nową)' },
        { klucz: 'nazwa', etykieta: 'Nazwa kolekcji', wartosc: 'Zestaw wdrożeniowy' },
        { klucz: 'opis', etykieta: 'Opis' },
        { klucz: 'barwa', etykieta: 'Oznaczenie barwne' },
        { klucz: 'pozycje', etykieta: 'Pozycje (identyfikatory po przecinku)' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.zapiszKolekcje({
          idKolekcji: (wartosci['kolekcja'] ?? '').trim(),
          nazwa: wymagajPola(wartosci['nazwa'] ?? '', 'Nazwa kolekcji'),
          opis: (wartosci['opis'] ?? '').trim(),
          oznaczenieBarwne: (wartosci['barwa'] ?? '').trim(),
          idPozycji: rozdzielPrzecinkamiRozszerzen(wartosci['pozycje'] ?? ''),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis kolekcji', wynik.blad);
        }
        return `Kolekcja ${wynik.wynik.collection.id} („${wynik.wynik.collection.name}"), pozycji: ${wynik.wynik.collection.extensionIds.length}.`;
      },
    },
    {
      etykieta: 'Kolekcja — zastosuj grupowo',
      komenda: Command.ExtensionCollectionApply,
      pola: [
        { klucz: 'kolekcja', etykieta: 'Kolekcja' },
        { klucz: 'wlacz', etykieta: 'Włączyć? (tak/nie)', wartosc: 'tak' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.zastosujKolekcje(
          wymagajPola(wartosci['kolekcja'] ?? '', 'Kolekcja'),
          czyTak(wartosci['wlacz'] ?? 'tak'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zastosowanie kolekcji', wynik.blad);
        }
        await stan.odczytajKatalog();
        return `Przestawiono ${wynik.wynik.applied.length} pozycji; odrzucono ${
          wynik.wynik.rejected.length
        }${
          wynik.wynik.rejected.length === 0
            ? ''
            : ` (${wynik.wynik.rejected.map((odrzucona) => odrzucona.reason).join('; ')})`
        }.`;
      },
    },
    {
      etykieta: 'Prywatny rejestr organizacji',
      komenda: Command.ExtensionRegistryList,
      pola: [{ klucz: 'fraza', etykieta: 'Fraza' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.rejestrOrganizacji((wartosci['fraza'] ?? '').trim(), '');
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt rejestru organizacji', wynik.blad);
        }
        if (wynik.wynik.total === 0) {
          return 'Rejestr organizacji jest pusty — trafiają do niego pozycje opublikowane z Publisher Panelu.';
        }
        return `Pozycji w rejestrze: ${wynik.wynik.total} — ${wynik.wynik.extensions
          .map((pozycja) => pozycja.code)
          .join(', ')}.`;
      },
    },
  ];
}

/** Narzędzia Installed Apps Managera: aktualizacje, przesyłka paczki, wersjonowanie i dziennik cyklu życia. */
export function narzedziaInstalledApps(stan: StanRozszerzen): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Sprawdź aktualizacje',
      komenda: Command.ExtensionUpdateCheck,
      async wykonaj() {
        const wynik = await stan.dobudowa.aktualizacje('');
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Sprawdzenie aktualizacji', wynik.blad);
        }
        if (wynik.wynik.updates.length === 0) {
          return 'Żadna zainstalowana pozycja nie ma nowszego wydania w rejestrze.';
        }
        return `Aktualizacji: ${wynik.wynik.updates.length} — ${wynik.wynik.updates
          .map(
            (aktualizacja) =>
              `${aktualizacja.extensionId}: ${aktualizacja.currentVersion} → ${aktualizacja.availableVersion}${
                aktualizacja.breaking === true ? ' (łamie zgodność)' : ''
              }`,
          )
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Prześlij paczkę (instalacja Personal)',
      komenda: Command.ExtensionPackageUpload,
      pola: [
        { klucz: 'nazwa', etykieta: 'Nazwa pliku', wartosc: 'wtyczka.zip' },
        { klucz: 'tresc', etykieta: 'Treść w base64', obszerne: true },
        { klucz: 'suma', etykieta: 'Suma kontrolna SHA-256' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.przeslijPaczke(
          wymagajPola(wartosci['nazwa'] ?? '', 'Nazwa pliku'),
          wymagajPola(wartosci['tresc'] ?? '', 'Treść w base64'),
          (wartosci['suma'] ?? '').trim(),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Przesłanie paczki', wynik.blad);
        }
        return `Paczka w magazynie rdzenia: ${wynik.wynik.uploadRef} (${wynik.wynik.sizeBytes} bajtów).`;
      },
    },
    {
      etykieta: 'Przypnij wersję',
      komenda: Command.ExtensionVersionPin,
      pola: [{ klucz: 'wersja', etykieta: 'Wersja (pusta zdejmuje przypięcie)' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.przypnijWersje(
          wskazana(stan),
          (wartosci['wersja'] ?? '').trim(),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Przypięcie wersji', wynik.blad);
        }
        stan.wchlonPozycje(wynik.wynik.extension);
        const wersja = wynik.wynik.pinnedVersion ?? '';
        return wersja === ''
          ? 'Przypięcie wersji zdjęte — pozycja idzie za rejestrem.'
          : `Wersja przypięta: ${wersja}.`;
      },
    },
    {
      etykieta: 'Cofnij do wcześniejszej wersji',
      komenda: Command.ExtensionVersionRollback,
      pola: [{ klucz: 'wersja', etykieta: 'Wersja docelowa' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.cofnijWersje(
          wskazana(stan),
          wymagajPola(wartosci['wersja'] ?? '', 'Wersja docelowa'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Cofnięcie wersji', wynik.blad);
        }
        stan.wchlonPozycje(wynik.wynik.extension);
        return `Cofnięto z wersji ${wynik.wynik.rolledBackFromVersion} do ${
          wynik.wynik.extension.version ?? 'bez wersji'
        }.`;
      },
    },
    {
      etykieta: 'Instalacja z manifestu zestawu',
      komenda: Command.ExtensionBundleInstall,
      pola: [
        {
          klucz: 'manifest',
          etykieta: 'Manifest zestawu (JSON)',
          wartosc: '{"extensions":[{"code":"mcp-repozytoria","kind":"mcp"}]}',
          obszerne: true,
        },
        { klucz: 'wlacz', etykieta: 'Włączyć po instalacji? (tak/nie)', wartosc: 'nie' },
      ],
      async wykonaj(wartosci) {
        const manifest = jsonZPola(wartosci['manifest'] ?? '', 'Manifest zestawu (JSON)');
        if (manifest === undefined) {
          throw new Error('Pole „Manifest zestawu (JSON)" jest wymagane — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.zainstalujZestaw(
          manifest,
          czyTak(wartosci['wlacz'] ?? 'nie'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Instalacja zestawu', wynik.blad);
        }
        await stan.odczytajKatalog();
        return `Zainstalowano ${wynik.wynik.installed.length} pozycji; odrzucono ${
          wynik.wynik.rejected.length
        }${
          wynik.wynik.rejected.length === 0
            ? ''
            : ` (${wynik.wynik.rejected
                .map((odrzucona) => `${odrzucona.code}: ${odrzucona.reason}`)
                .join('; ')})`
        }.`;
      },
    },
    {
      etykieta: 'Dziennik cyklu życia',
      komenda: Command.ExtensionHistoryList,
      pola: [{ klucz: 'granica', etykieta: 'Górna granica wpisów' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.historia('', liczbaZPola(wartosci['granica'] ?? ''));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt dziennika cyklu życia', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'W dzienniku cyklu życia nie ma jeszcze ani jednego wpisu.';
        const najnowszy = wynik.wynik.entries[0];
        return `Wpisów: ${wynik.wynik.total}; najnowszy — ${najnowszy?.action ?? 'bez czynności'} pozycji ${
          najnowszy?.extensionId ?? ''
        }.`;
      },
    },
    {
      etykieta: 'Operacja zbiorcza rejestru',
      komenda: Command.ExtensionAdminBulk,
      pola: [
        { klucz: 'pozycje', etykieta: 'Pozycje (identyfikatory po przecinku)' },
        { klucz: 'czynnosc', etykieta: 'Czynność (enable/disable/uninstall)', wartosc: 'enable' },
      ],
      async wykonaj(wartosci) {
        const czynnosc = wartoscWyliczeniaRozszerzen(
          wartosci['czynnosc'] ?? '',
          ExtensionBulkAction,
          'czynność zbiorcza',
        );
        if (czynnosc === '') {
          throw new Error('Pole „Czynność" jest wymagane przez kontrakt — żądanie nie poszło.');
        }
        const pozycje = rozdzielPrzecinkamiRozszerzen(wartosci['pozycje'] ?? '');
        if (pozycje.length === 0) {
          throw new Error('Operacja zbiorcza wymaga co najmniej jednej pozycji.');
        }
        const wynik = await stan.dobudowa.operacjaZbiorcza(pozycje, czynnosc);
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Operacja zbiorcza', wynik.blad);
        }
        await stan.odczytajKatalog();
        return `Zmieniono ${wynik.wynik.affected.length} pozycji; odrzucono ${wynik.wynik.rejected.length}.`;
      },
    },
  ];
}

/** Narzędzia Integrations Hub oraz MCP i Connector Console modułu Apps: transport, poświadczenia i protokół. */
export function narzedziaIntegracji(stan: StanRozszerzen): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Transport integracji',
      komenda: Command.ExtensionTransportSet,
      pola: [
        { klucz: 'transport', etykieta: 'Transport (stdio/sse/http/streamableHttp)', wartosc: 'http' },
        { klucz: 'adres', etykieta: 'Adres serwera' },
        { klucz: 'polecenie', etykieta: 'Polecenie procesu (stdio)' },
        { klucz: 'sprawdz', etykieta: 'Sprawdzić powitaniem? (tak/nie)', wartosc: 'tak' },
      ],
      async wykonaj(wartosci) {
        const transport = wartoscWyliczeniaRozszerzen(
          wartosci['transport'] ?? '',
          McpTransport,
          'transport',
        );
        if (transport === '') {
          throw new Error('Pole „Transport" jest wymagane przez kontrakt — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.ustawTransport({
          idRozszerzenia: wskazana(stan),
          transport,
          adres: (wartosci['adres'] ?? '').trim(),
          polecenie: (wartosci['polecenie'] ?? '').trim(),
          sprawdz: czyTak(wartosci['sprawdz'] ?? 'tak'),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Ustawienie transportu', wynik.blad);
        }
        stan.wchlonPozycje(wynik.wynik.extension);
        return `Transport zapisany; próbne powitanie: ${wynik.wynik.probeStatus ?? 'niesprawdzane'}.`;
      },
    },
    {
      etykieta: 'Poświadczenie integracji',
      komenda: Command.ExtensionCredentialBind,
      pola: [
        { klucz: 'sposob', etykieta: 'Sposób (oauth2/apiKey/token/basic/none)', wartosc: 'apiKey' },
        { klucz: 'odwolanie', etykieta: 'Odwołanie do sekretu (klucz jawny)' },
        { klucz: 'zakresy', etykieta: 'Zakresy OAuth2 (po przecinku)' },
      ],
      async wykonaj(wartosci) {
        const sposob = wartoscWyliczeniaRozszerzen(
          wartosci['sposob'] ?? '',
          ExtensionAuthKind,
          'sposób uwierzytelnienia',
        );
        if (sposob === '') {
          throw new Error('Pole „Sposób" jest wymagane przez kontrakt — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.powiazPoswiadczenie({
          idRozszerzenia: wskazana(stan),
          sposob,
          // To pole przyjmuje wyłącznie klucz jawny, nigdy treść poświadczenia.
          odwolanie: wymagajPola(wartosci['odwolanie'] ?? '', 'Odwołanie do sekretu'),
          zakresy: rozdzielPrzecinkamiRozszerzen(wartosci['zakresy'] ?? ''),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Powiązanie poświadczenia', wynik.blad);
        }
        stan.wchlonPozycje(wynik.wynik.extension);
        return wynik.wynik.authorizationUrl === undefined
          ? 'Odwołanie do poświadczenia powiązane z integracją.'
          : `Odwołanie powiązane; adres zgody: ${wynik.wynik.authorizationUrl}`;
      },
    },
    {
      etykieta: 'Odkryj narzędzia i zasoby',
      komenda: Command.ExtensionToolList,
      pola: [
        { klucz: 'rodzaj', etykieta: 'Rodzaj (tool/resource/prompt)' },
        { klucz: 'odswiez', etykieta: 'Zapytać serwer? (tak/nie)', wartosc: 'tak' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.narzedzia(
          wskazana(stan),
          wartoscWyliczeniaRozszerzen(wartosci['rodzaj'] ?? '', ExtensionToolKind, 'rodzaj wpisu'),
          czyTak(wartosci['odswiez'] ?? 'tak'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odkrywanie narzędzi', wynik.blad);
        }
        if (wynik.wynik.entries.length === 0) {
          return 'Serwer nie udostępnia jeszcze ani jednego wpisu.';
        }
        return `Wpisów: ${wynik.wynik.entries.length} (protokół ${
          wynik.wynik.protocolVersion ?? 'bez wersji'
        }) — ${wynik.wynik.entries.map((wpis) => `${wpis.kind}:${wpis.name}`).join(', ')}.`;
      },
    },
    {
      etykieta: 'Próbne wywołanie narzędzia',
      komenda: Command.ExtensionToolCall,
      pola: [
        { klucz: 'narzedzie', etykieta: 'Narzędzie' },
        { klucz: 'argumenty', etykieta: 'Argumenty (JSON)', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.wywolajNarzedzie(
          wskazana(stan),
          wymagajPola(wartosci['narzedzie'] ?? '', 'Narzędzie'),
          jsonZPola(wartosci['argumenty'] ?? '', 'Argumenty (JSON)'),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Próbne wywołanie', wynik.blad);
        }
        // Niepowodzenie serwera jest wynikiem, nie odmową — inspektor pokazuje
        // je tak samo jak powodzenie.
        if (!wynik.wynik.ok) {
          return `Serwer odmówił po ${wynik.wynik.durationMs} ms: ${
            wynik.wynik.errorDetail ?? 'bez powodu'
          }`;
        }
        return `Odpowiedź w ${wynik.wynik.durationMs} ms: ${wynik.wynik.text ?? '(bez treści tekstowej)'}`;
      },
    },
    {
      etykieta: 'Log protokołu JSON-RPC',
      komenda: Command.ExtensionProtocolLogList,
      pola: [{ klucz: 'granica', etykieta: 'Górna granica ramek' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.logProtokolu(
          wskazana(stan),
          liczbaZPola(wartosci['granica'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt logu protokołu', wynik.blad);
        }
        if (wynik.wynik.total === 0) {
          return 'Z tą integracją nie odbyła się jeszcze ani jedna rozmowa.';
        }
        const najnowsza = wynik.wynik.frames[0];
        return `Ramek: ${wynik.wynik.total}; najnowsza — ${najnowsza?.direction ?? ''} ${
          najnowsza?.method ?? '(odpowiedź)'
        }.`;
      },
    },
    {
      etykieta: 'Piaskownica przebiegu próbnego',
      komenda: Command.ExtensionSandboxRun,
      pola: [
        { klucz: 'wejscie', etykieta: 'Wejście (JSON)', wartosc: '{}', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const wejscie = jsonZPola(wartosci['wejscie'] ?? '', 'Wejście (JSON)');
        if (wejscie === undefined) {
          throw new Error('Pole „Wejście (JSON)" jest wymagane przez kontrakt — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.piaskownica(wskazana(stan), wejscie);
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Przebieg piaskownicy', wynik.blad);
        }
        return `Przebieg ${wynik.wynik.ok ? 'udany' : 'nieudany'} w ${
          wynik.wynik.durationMs
        } ms; dziennik: ${wynik.wynik.logRef ?? 'brak'}.`;
      },
    },
    {
      etykieta: 'Import definicji API',
      komenda: Command.ExtensionDefinitionImport,
      pola: [
        { klucz: 'format', etykieta: 'Format (openapi3/graphql)', wartosc: 'openapi3' },
        { klucz: 'kod', etykieta: 'Kod pozycji', wartosc: 'konektor-api' },
        { klucz: 'zrodlo', etykieta: 'Adres opisu' },
        { klucz: 'tresc', etykieta: 'Treść opisu', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const format = wartoscWyliczeniaRozszerzen(
          wartosci['format'] ?? '',
          ExtensionDefinitionFormat,
          'format definicji',
        );
        if (format === '') {
          throw new Error('Pole „Format" jest wymagane przez kontrakt — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.zaimportujDefinicje({
          format,
          kod: wymagajPola(wartosci['kod'] ?? '', 'Kod pozycji'),
          zrodlo: (wartosci['zrodlo'] ?? '').trim(),
          tresc: (wartosci['tresc'] ?? '').trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Import definicji', wynik.blad);
        }
        await stan.odczytajKatalog();
        return `Pozycja ${wynik.wynik.extension.code}: operacji ${wynik.wynik.operations.length}.`;
      },
    },
    {
      etykieta: 'Webhooki integracji',
      komenda: Command.ExtensionWebhookList,
      pola: [{ klucz: 'kierunek', etykieta: 'Kierunek (inbound/outbound)' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.webhooki(
          wskazana(stan),
          wartoscWyliczeniaRozszerzen(
            wartosci['kierunek'] ?? '',
            ExtensionWebhookDirection,
            'kierunek webhooka',
          ),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt webhooków', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'Ta integracja nie ma jeszcze ani jednego webhooka.';
        return `Webhooków: ${wynik.wynik.total} — ${wynik.wynik.webhooks
          .map((webhook) => `${webhook.direction}${webhook.enabled ? ' (czynny)' : ''}`)
          .join(', ')}.`;
      },
    },
    {
      etykieta: 'Webhook — zapis',
      komenda: Command.ExtensionWebhookSave,
      pola: [
        { klucz: 'webhook', etykieta: 'Webhook zmieniany (pusty zakłada nowy)' },
        { klucz: 'kierunek', etykieta: 'Kierunek (inbound/outbound)', wartosc: 'outbound' },
        { klucz: 'adres', etykieta: 'Adres docelowy' },
        { klucz: 'zdarzenia', etykieta: 'Zdarzenia (po przecinku)' },
        { klucz: 'sekret', etykieta: 'Odwołanie do sekretu HMAC' },
        { klucz: 'czynny', etykieta: 'Czynny? (tak/nie)', wartosc: 'tak' },
      ],
      async wykonaj(wartosci) {
        const kierunek = wartoscWyliczeniaRozszerzen(
          wartosci['kierunek'] ?? '',
          ExtensionWebhookDirection,
          'kierunek webhooka',
        );
        if (kierunek === '') {
          throw new Error('Pole „Kierunek" jest wymagane przez kontrakt — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.zapiszWebhook({
          idWebhooka: (wartosci['webhook'] ?? '').trim(),
          idRozszerzenia: wskazana(stan),
          kierunek,
          adres: (wartosci['adres'] ?? '').trim(),
          zdarzenia: rozdzielPrzecinkamiRozszerzen(wartosci['zdarzenia'] ?? ''),
          odwolanieSekretu: (wartosci['sekret'] ?? '').trim(),
          czynny: czyTak(wartosci['czynny'] ?? 'tak'),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis webhooka', wynik.blad);
        }
        const webhook = wynik.wynik.webhook;
        return `Webhook ${webhook.id} (${webhook.direction})${
          webhook.receiveUrl === undefined ? '' : `, adres nasłuchu: ${webhook.receiveUrl}`
        }.`;
      },
    },
    {
      etykieta: 'Odwzorowanie danych',
      komenda: Command.ExtensionMappingSave,
      pola: [
        { klucz: 'mapowanie', etykieta: 'Odwzorowanie zmieniane (puste zakłada nowe)' },
        { klucz: 'nazwa', etykieta: 'Nazwa odwzorowania', wartosc: 'klient → kontrahent' },
        { klucz: 'reguly', etykieta: 'Reguły (JSON)', wartosc: '{"name":"nazwa"}', obszerne: true },
      ],
      async wykonaj(wartosci) {
        const reguly = jsonZPola(wartosci['reguly'] ?? '', 'Reguły (JSON)');
        if (reguly === undefined) {
          throw new Error('Pole „Reguły (JSON)" jest wymagane przez kontrakt — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.zapiszMapowanie(
          wskazana(stan),
          (wartosci['mapowanie'] ?? '').trim(),
          wymagajPola(wartosci['nazwa'] ?? '', 'Nazwa odwzorowania'),
          reguly,
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Zapis odwzorowania', wynik.blad);
        }
        const zastrzezenia = wynik.wynik.validationIssues ?? [];
        return zastrzezenia.length === 0
          ? `Odwzorowanie ${wynik.wynik.mapping.id} zapisane bez zastrzeżeń.`
          : `Odwzorowanie ${wynik.wynik.mapping.id} zapisane; zastrzeżenia (nie wstrzymują): ${zastrzezenia.join('; ')}.`;
      },
    },
    {
      etykieta: 'Metryki użycia i kosztu',
      komenda: Command.ExtensionUsageGet,
      async wykonaj() {
        const wynik = await stan.dobudowa.uzycie('');
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt metryk użycia', wynik.blad);
        }
        if (wynik.wynik.usage.length === 0) {
          return 'W oknie pomiaru nie odbyło się ani jedno wywołanie integracji.';
        }
        return `Integracji z użyciem: ${wynik.wynik.usage.length} — ${wynik.wynik.usage
          .map(
            (wpis) =>
              `${wpis.extensionId}: ${wpis.calls} wywołań, ${wpis.failures} niepowodzeń, ${
                wpis.avgLatencyMs ?? '—'
              } ms`,
          )
          .join('; ')}.`;
      },
    },
    {
      etykieta: 'Panel zdrowia integracji',
      komenda: Command.ExtensionHealthCheck,
      async wykonaj() {
        const wynik = await stan.dobudowa.kondycja('');
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Sprawdzenie kondycji', wynik.blad);
        }
        if (wynik.wynik.results.length === 0) {
          return 'Żadna integracja nie jest włączona — nie ma czego sprawdzać.';
        }
        return wynik.wynik.results
          .map(
            (rezultat) =>
              `${rezultat.extensionId}: ${rezultat.status}${
                rezultat.latencyMs === undefined ? '' : ` (${rezultat.latencyMs} ms)`
              }`,
          )
          .join('; ');
      },
    },
    {
      etykieta: 'Audyt użycia',
      komenda: Command.ExtensionAuditList,
      pola: [
        { klucz: 'ekspert', etykieta: 'Ekspert' },
        { klucz: 'granica', etykieta: 'Górna granica wpisów' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.audyt(
          '',
          (wartosci['ekspert'] ?? '').trim(),
          liczbaZPola(wartosci['granica'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt audytu', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'Audyt nie ma jeszcze ani jednego wpisu użycia.';
        return `Wpisów audytu: ${wynik.wynik.total}; najnowszy — pozycja ${
          wynik.wynik.entries[0]?.extensionId ?? ''
        }, narzędzie ${wynik.wynik.entries[0]?.toolName ?? '—'}.`;
      },
    },
  ];
}

/** Narzędzia Permissions and Trust Center modułu Apps: uprawnienia, podpis, skaner bezpieczeństwa i sekrety. */
export function narzedziaZaufania(stan: StanRozszerzen): readonly NarzedzieApps[] {
  return [
    {
      etykieta: 'Uprawnienia deklarowane i nadane',
      komenda: Command.ExtensionPermissionList,
      async wykonaj() {
        const wynik = await stan.dobudowa.uprawnienia(wskazana(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt uprawnień', wynik.blad);
        }
        const nadmiar = wynik.wynik.excessive ?? [];
        return `Deklarowanych: ${(wynik.wynik.declared ?? []).length}, nadanych: ${
          (wynik.wynik.granted ?? []).length
        }${nadmiar.length === 0 ? '' : `, nadmiarowych: ${nadmiar.join(', ')}`}.`;
      },
    },
    {
      etykieta: 'Nadaj zakres uprawnień',
      komenda: Command.ExtensionPermissionGrant,
      pola: [
        {
          klucz: 'zakresy',
          etykieta: 'Zakresy (network:domena, fileRead:/dane …)',
          wartosc: 'network:example.com',
        },
        { klucz: 'ekspert', etykieta: 'Ekspert (pusty: cała platforma)' },
      ],
      async wykonaj(wartosci) {
        const uprawnienia = uprawnieniaZPola(wartosci['zakresy'] ?? '');
        if (uprawnienia.length === 0) {
          throw new Error('Nadanie wymaga co najmniej jednego zakresu — żądanie nie poszło.');
        }
        const wynik = await stan.dobudowa.nadajUprawnienia(
          wskazana(stan),
          uprawnienia,
          (wartosci['ekspert'] ?? '').trim(),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Nadanie uprawnień', wynik.blad);
        }
        return `Nadano ${wynik.wynik.granted.length} uprawnień; nadanie WYMIENIA komplet, więc zakresy pominięte zostały cofnięte.`;
      },
    },
    {
      etykieta: 'Weryfikacja podpisu i pochodzenia',
      komenda: Command.ExtensionSignatureVerify,
      async wykonaj() {
        const wynik = await stan.dobudowa.zweryfikujPodpis(wskazana(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Weryfikacja podpisu', wynik.blad);
        }
        const podpis = wynik.wynik.signature;
        return `Podpisana: ${podpis.signed ? 'tak' : 'nie'}, weryfikacja: ${
          podpis.verified ? 'przeszła' : 'nie przeszła'
        }, zaufanie: ${podpis.trustLevel}${
          podpis.detail === undefined ? '' : ` — ${podpis.detail}`
        }.`;
      },
    },
    {
      etykieta: 'Skaner manifestu',
      komenda: Command.ExtensionManifestScan,
      async wykonaj() {
        const wynik = await stan.dobudowa.skanujManifest(wskazana(stan));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Skan manifestu', wynik.blad);
        }
        if (wynik.wynik.findings.length === 0) {
          return 'Skaner nie ma zastrzeżeń do tej pozycji.';
        }
        // Skaner jest sygnałem, nie bramą — zdanie mówi to wprost, żeby nikt
        // nie wziął spostrzeżenia za zakaz.
        return `Spostrzeżeń: ${wynik.wynik.findings.length} (sygnał, nie brama) — ${wynik.wynik.findings
          .map((spostrzezenie) => `${spostrzezenie.severity}: ${spostrzezenie.message}`)
          .join('; ')}.`;
      },
    },
    {
      etykieta: 'Rejestr referencji sekretów',
      komenda: Command.ExtensionSecretList,
      pola: [{ klucz: 'dni', etykieta: 'Wygasające w ciągu dni' }],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.sekrety('', liczbaZPola(wartosci['dni'] ?? ''));
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Odczyt rejestru sekretów', wynik.blad);
        }
        if (wynik.wynik.total === 0) return 'Rejestr referencji sekretów jest pusty.';
        return `Referencji: ${wynik.wynik.total} — ${wynik.wynik.secrets
          .map((sekret) => sekret.ref)
          .join(', ')}. Rejestr niesie klucze jawne, nigdy treść poświadczeń.`;
      },
    },
    {
      etykieta: 'Zakres współdzielenia sekretu',
      komenda: Command.ExtensionSecretShare,
      pola: [
        { klucz: 'odwolanie', etykieta: 'Odwołanie do sekretu' },
        { klucz: 'pozycje', etykieta: 'Pozycje (identyfikatory po przecinku)' },
        { klucz: 'role', etykieta: 'Role (identyfikatory po przecinku)' },
      ],
      async wykonaj(wartosci) {
        const wynik = await stan.dobudowa.udostepnijSekret(
          wymagajPola(wartosci['odwolanie'] ?? '', 'Odwołanie do sekretu'),
          rozdzielPrzecinkamiRozszerzen(wartosci['pozycje'] ?? ''),
          rozdzielPrzecinkamiRozszerzen(wartosci['role'] ?? ''),
        );
        if (!wynik.udany || wynik.wynik === undefined) {
          throw bladOdmowyNarzedzia('Udostępnienie sekretu', wynik.blad);
        }
        const sekret = wynik.wynik.secret;
        return `Referencja ${sekret.ref}: pozycji ${
          sekret.sharedWithExtensionIds?.length ?? 0
        }, ról ${sekret.sharedWithRoleIds?.length ?? 0}.`;
      },
    },
  ];
}

// ── Przekłady wartości pól ───────────────────────────────────────────────────

/** Rozdziela wartości pola wpisane po przecinku na wykaz pozycji, pomijając pozycje puste powstałe z odstępów. */
function rozdzielPrzecinkamiRozszerzen(wartosc: string): readonly string[] {
  return wartosc
    .split(',')
    .map((czesc) => czesc.trim())
    .filter((czesc) => czesc !== '');
}

/**
 * Rozstrzyga pole logiczne wpisane słowem. „tak" i „nie" są tym, co Operator
 * naprawdę wpisuje; wartość inna jest odmową, a nie cichym „nie".
 */
function czyTak(wartosc: string): boolean {
  const przyciete = wartosc.trim().toLowerCase();
  if (przyciete === 'tak' || przyciete === 'true' || przyciete === '1') return true;
  if (przyciete === 'nie' || przyciete === 'false' || przyciete === '0' || przyciete === '') {
    return false;
  }
  throw new Error(`Wartość „${wartosc}" nie jest odpowiedzią tak/nie — żądanie nie poszło.`);
}

/**
 * Rozkłada zakresy uprawnień wpisane w postaci `zakres:byt`. Byt pominięty
 * znaczy „bez zawężenia" — i to właśnie wychwytuje potem skaner manifestu.
 */
function uprawnieniaZPola(wartosc: string): readonly ExtensionPermission[] {
  const znane = Object.values(ExtensionPermissionScope) as readonly string[];
  return rozdzielPrzecinkamiRozszerzen(wartosc).map((wpis) => {
    const [zakres, byt] = wpis.split(':');
    const nazwaZakresu = (zakres ?? '').trim();
    if (!znane.includes(nazwaZakresu)) {
      throw new Error(
        `Zakres „${nazwaZakresu}" nie jest wartością kontraktu (${znane.join(', ')}).`,
      );
    }
    const uprawnienie: ExtensionPermission = {
      scope: nazwaZakresu as ExtensionPermissionScope,
    };
    const wskazanyByt = (byt ?? '').trim();
    if (wskazanyByt !== '') uprawnienie.target = wskazanyByt;
    return uprawnienie;
  });
}

/**
 * Wspólny przekład pola na wartość wyliczenia kontraktu. Wykaz dopuszczalnych
 * wartości bierze się z wyliczenia, nie z listy wpisanej tutaj.
 */
function wartoscWyliczeniaRozszerzen<T extends string>(
  wartosc: string,
  wyliczenie: Record<string, T>,
  nazwa: string,
): T | '' {
  const przyciete = wartosc.trim();
  if (przyciete === '') return '';
  const znane = Object.values(wyliczenie) as readonly string[];
  if (!znane.includes(przyciete)) {
    throw new Error(
      `Wartość „${przyciete}" nie jest wartością kontraktu dla pola ${nazwa} (${znane.join(', ')}).`,
    );
  }
  return przyciete as T;
}
