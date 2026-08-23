import type { BrakFunkcji } from './braki-kontraktu';

/**
 * Teksty widoczne dla Operatora w oknach modułu Apps, trzymane osobno od
 * plików budujących elementy (wzór: `sterowanie/etykiety-sterowania.ts`).
 *
 * Pozycja wykazu braków niesie tylko to, co okno wie samo: nazwę czynności,
 * czego by do niej trzeba (`czego`), nazwę komendy, której wejście brak znosi
 * (`komendaZnoszaca` — pozycja znika wtedy z wykazu sama) i cudze drogi do
 * sprawdzenia w wykazie (`komendyCudze`). Zdanie o stanie kontraktu dokłada
 * `braki-kontraktu.ts`, czytając `KOMENDY` przy składaniu okna.
 */

/**
 * Kody okien z katalogu rdzenia (`okno_operacyjne.kod`). Kod jest bezmodułowy —
 * definicje okien mają kody bez przedrostka modułu, więc kody
 * z `Module.operationalWindowCodes` zestawiają się wprost z tymi wartościami.
 */
export const KODY_OKIEN = {
  ProductBuilder: 'product-builder',
  ArchitectureDesigner: 'architecture-designer',
  FrontendWorkspace: 'frontend-workspace',
  BackendWorkspace: 'backend-workspace',
  DeploymentPanel: 'deployment-panel',
  AppCatalog: 'app-catalog',
  InstalledAppsManager: 'installed-apps-manager',
  PermissionsTrustCenter: 'permissions-trust-center',
  IntegrationsHub: 'integrations-hub',
  McpConnectorConsole: 'mcp-connector-console',
  PublisherPanel: 'publisher-panel',
} as const;

/**
 * Kody okien strony dystrybucji i konsumpcji, których katalog rdzenia jeszcze
 * nie zna.
 *
 * Rdzeń wnosi definicje okien migracją słownikową i przypina je do modułów
 * macierzą; sześć okien tej strony nie ma tam ani wiersza definicji, ani
 * wiersza przypięcia. Okno zbudowane w kliencie działa mimo to, bo komendy
 * `extension.*` i `access.*` nie niosą pola `windowId` — katalog rozszerzeń
 * stoi poziom wyżej niż okno modułu. Wykaz stoi tutaj, żeby okno umiało
 * powiedzieć o tej granicy zamiast ją przemilczeć.
 */
export const OKNA_SPOZA_KATALOGU_RDZENIA: readonly string[] = [
  KODY_OKIEN.AppCatalog,
  KODY_OKIEN.InstalledAppsManager,
  KODY_OKIEN.PermissionsTrustCenter,
  KODY_OKIEN.IntegrationsHub,
  KODY_OKIEN.McpConnectorConsole,
  KODY_OKIEN.PublisherPanel,
];

/** Nazwy okien operacyjnych modułu, wyświetlane w nagłówkach. */
export const NAZWY_OKIEN: Readonly<Record<string, string>> = {
  [KODY_OKIEN.ProductBuilder]: 'Product Builder',
  [KODY_OKIEN.ArchitectureDesigner]: 'Architecture Designer',
  [KODY_OKIEN.FrontendWorkspace]: 'Frontend Workspace',
  [KODY_OKIEN.BackendWorkspace]: 'Backend Workspace',
  [KODY_OKIEN.DeploymentPanel]: 'Deployment Panel',
  [KODY_OKIEN.AppCatalog]: 'App Catalog',
  [KODY_OKIEN.InstalledAppsManager]: 'Installed Apps Manager',
  [KODY_OKIEN.PermissionsTrustCenter]: 'Permissions & Trust Center',
  [KODY_OKIEN.IntegrationsHub]: 'Integrations Hub',
  [KODY_OKIEN.McpConnectorConsole]: 'MCP & Connector Console',
  [KODY_OKIEN.PublisherPanel]: 'Publisher Panel',
};

/**
 * Nazwy komend zaproponowanych dla czynności opracowania, których dzisiejszy
 * kontrakt nie prowadzi.
 *
 * Nazwa proponowana nie jest nazwą kontraktu i nie idzie na drut: pozycja
 * z takim wskazaniem stoi w wykazie braków, a wpisanie komendy do kontraktu
 * zdejmuje pozycję z wykazu samo (`braki-kontraktu.ts` czyta `KOMENDY`
 * przy składaniu okna). Nazwy stoją w jednym miejscu, żeby definicja oddana
 * właścicielowi i zdanie widoczne Operatorowi mówiły o tej samej komendzie.
 *
 * Obszar nazwy idzie za bytem, którego dotyczy: czynności na katalogu
 * rozszerzeń należą do obszaru `extension`, bo katalog stoi poziom wyżej niż
 * moduł, a czynności na produkcie budowanym w module — do obszaru `apps`.
 */
export const KOMENDY_PROPONOWANE = {
  SzukanieKatalogu: 'extension.search',
  SzczegolPozycji: 'extension.detail.get',
  KolekcjeKuratorskie: 'extension.collection.list',
  RejestrOrganizacji: 'extension.registry.list',
  SprawdzenieAktualizacji: 'extension.update.check',
  PrzeslaniePaczki: 'extension.package.upload',
  PrzypiecieWersji: 'extension.version.pin',
  CofniecieWersji: 'extension.version.rollback',
  InstalacjaZestawu: 'extension.bundle.install',
  DziennikCyklu: 'extension.history.list',
  OdkrycieNarzedzi: 'extension.tool.list',
  WywolanieProbne: 'extension.tool.call',
  DziennikProtokolu: 'extension.protocol.log.list',
  PiaskownicaTestu: 'extension.sandbox.run',
  ImportDefinicji: 'extension.definition.import',
  WyborTransportu: 'extension.transport.set',
  PowiazanieDanychDostepowych: 'extension.credential.bind',
  Webhooki: 'extension.webhook.list',
  MapowanieDanych: 'extension.mapping.save',
  MetrykiUzycia: 'extension.usage.get',
  ZdrowieIntegracji: 'extension.health.check',
  ReferencjeSekretow: 'extension.secret.list',
  WspoldzielenieSekretu: 'extension.secret.share',
  WykazUprawnien: 'extension.permission.list',
  NadanieUprawnien: 'extension.permission.grant',
  WeryfikacjaPodpisu: 'extension.signature.verify',
  SkanerManifestu: 'extension.manifest.scan',
  AudytUprawnien: 'extension.audit.list',
  EtapyBudowy: 'apps.stage.list',
  ZapisEtapu: 'apps.stage.save',
  MetadaneProduktu: 'apps.product.get',
  KamienieMilowe: 'apps.milestone.list',
  EksportUkladu: 'apps.architecture.export',
  PodgladNaZywo: 'apps.preview.start',
  PunktyKoncowe: 'apps.endpoint.list',
  ZapytanieProbne: 'apps.endpoint.probe',
  ZmienneSrodowiska: 'apps.environment.variable.list',
  SchematBazy: 'apps.schema.get',
  ArtefaktyBudowania: 'apps.artifact.list',
  LogWdrozenia: 'apps.deployment.log.read',
  KondycjaProduktu: 'apps.deployment.health.get',
  DomenaWdrozenia: 'apps.deployment.domain.set',
  SkalowanieWdrozenia: 'apps.deployment.scale.set',
  BudowaPakietu: 'apps.package.build',
  ManifestPakietu: 'apps.package.manifest.save',
  WalidacjaPakietu: 'apps.package.validate',
  PodpisPakietu: 'apps.package.sign',
  PublikacjaPakietu: 'apps.package.publish',
} as const;

/** Zdanie stanu pustego, gdy rdzeń nie wskazał jeszcze okna modułu. */
export const BEZ_OKNA_MODULU =
  'Rdzeń nie wskazał jeszcze okna modułu Apps. Każda komenda obszaru wymaga pola ' +
  'windowId, więc okno nie wysyła nic, dopóki go nie ma.';

/**
 * Zdanie stanu pustego dla okien, których treść wypełnia sam strumień zdarzeń.
 *
 * Dotyczy osi etapów: etapy budowy przychodzą jedynie ramką `apps.build.changed`,
 * bo komendy ich odczytu kontrakt nie niesie. Wdrożenia, architektura i pliki
 * warsztatu mają własne komendy odczytu i mówią o swojej pustce inaczej —
 * zdaniem, które odróżnia „rdzeń nic nie ma" od „nikt nie pytał".
 */
export const BEZ_KOMENDY_ODCZYTU =
  'Etapy budowy nie mają w kontrakcie komendy odczytu — pojawią się, gdy rdzeń ' +
  'przyśle apps.build.changed. Pustka znaczy tu „nic jeszcze nie przyszło", ' +
  'a nie nieudany odczyt.';

export const BRAKI_PRODUCT_BUILDER: readonly BrakFunkcji[] = [
  {
    etykieta: 'Kamienie milowe',
    czego: 'Kamień milowy musiałby być osobnym bytem umowy z własną komendą zapisu.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.KamienieMilowe,
  },
  {
    etykieta: 'Odczyt etapów budowy',
    czego:
      'Etapy wypełnia wyłącznie zdarzenie apps.build.changed — wejście w moduł ' +
      'zaczyna od pustego wykazu, choć rdzeń trzyma etapy w bazie.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.EtapyBudowy,
  },
  {
    etykieta: 'Przypisanie wykonawców',
    czego:
      'Etap niesie pole ownerAgentId wyłącznie w zdarzeniu apps.build.changed, ' +
      'czyli w kierunku od rdzenia — ustawienie go wymagałoby komendy w drugą stronę.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.ZapisEtapu,
  },
  {
    etykieta: 'Widżety integracji',
    czego: 'Katalog integracji produktu musiałby być bytem umowy z komendą odczytu.',
  },
  {
    etykieta: 'Centrum dokumentacji',
    czego: 'Dokumenty produktu leżą w bibliotece, nie w architekturze ani we wdrożeniu.',
    komendyCudze: ['library.file.list', 'library.collection.create'],
  },
  {
    etykieta: 'Metadane produktu',
    czego:
      'Produkt jako całość nie jest bytem umowy — apps.architecture.define ' +
      'opisuje architekturę JEDNEGO okna, a nie produkt, do którego należy.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.MetadaneProduktu,
  },
];

export const BRAKI_ARCHITECTURE_DESIGNER: readonly BrakFunkcji[] = [
  {
    etykieta: 'Wizualizacja przepływu danych',
    czego:
      'AppComponent niesie zależności (dependsOn), a nie kierunek ani rodzaj ' +
      'przepływu — nie ma z czego narysować strzałki danych.',
  },
  {
    etykieta: 'Adnotacje',
    czego: 'Adnotacja musiałaby być polem komponentu albo osobnym bytem umowy.',
  },
  {
    etykieta: 'Eksport diagramu',
    czego: 'Eksport układu jako pliku musiałby mieć komendę wytwarzającą treść po stronie rdzenia.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.EksportUkladu,
    komendyCudze: ['research.report.export'],
  },
  {
    etykieta: 'Odczyt zapisanej architektury',
    czego: 'Powrót do układu zapisanego wcześniej wymaga komendy odczytu architektury okna.',
    komendaZnoszaca: 'apps.architecture.get',
  },
];

export const BRAKI_FRONTEND: readonly BrakFunkcji[] = [
  {
    etykieta: 'Podgląd na żywo (hot reload)',
    czego:
      'Podgląd na żywo wymaga komendy uruchamiającej serwer podglądu i zdarzenia ' +
      'jego odświeżenia — rdzeń nie stawia dziś serwera produktu.',
  },
  {
    etykieta: 'Podgląd responsywny',
    czego: 'Podgląd responsywny stoi na podglądzie na żywo, którego nie ma.',
  },
  {
    etykieta: 'Selektor zasobów z Design',
    czego: 'Przeniesienie zasobu między modułami jest czynnością nadawcy, czyli modułu Design.',
    komendyCudze: ['context.transfer'],
  },
  {
    etykieta: 'Edytor motywu',
    czego: 'Motyw produktu musiałby być bytem umowy z komendą zapisu.',
  },
  {
    etykieta: 'Mapa routingu',
    czego: 'Trasy produktu musiałyby być bytem umowy z komendą odczytu.',
  },
  {
    etykieta: 'Inspektor stanu',
    czego: 'Odczyt stanu URUCHOMIONEGO produktu wymaga uruchomionego produktu, którego rdzeń nie hostuje.',
  },
  {
    etykieta: 'Wskaźnik budowania i lintowania',
    czego: 'Budowanie kodu jest czynnością modułu Developer, a nie warsztatu Apps.',
    komendyCudze: ['developer.build.run'],
  },
  {
    etykieta: 'Integracja z Git Panel i Terminal',
    czego: 'Oba okna należą do modułów Developer i Terminal i mają własne komendy.',
    komendyCudze: ['developer.git.action', 'terminal.command.exec'],
  },
];

export const BRAKI_BACKEND: readonly BrakFunkcji[] = [
  {
    etykieta: 'Eksplorator punktów końcowych API',
    czego:
      'Komponent niesie apiContract jako jeden tekst — punkt końcowy nie jest ' +
      'bytem umowy, więc nie ma czego wyliczyć w wykazie.',
  },
  {
    etykieta: 'Konstruktor zapytań testowych',
    czego: 'Zapytanie próbne wobec usługi musiałoby mieć komendę wykonania po stronie rdzenia.',
  },
  {
    etykieta: 'Podgląd schematu bazy',
    czego: 'Schemat bazy produktu musiałby być bytem umowy z komendą odczytu.',
  },
  {
    etykieta: 'Zmienne środowiskowe',
    czego:
      'apps.deployment.run przyjmuje ŚRODOWISKO wdrożenia (dev, staging, prod), ' +
      'a nie zmienne w nim — zapis zmiennych byłby osobną komendą.',
  },
  {
    etykieta: 'Logi na żywo usług',
    czego: 'Strumień logów usługi wymagałby osobnego zdarzenia w obszarze apps.',
  },
  {
    etykieta: 'Powiązanie z kolejkami',
    czego:
      'Kolejka zakłada okna sesji, a usługa warsztatu jest plikiem ' +
      'warstwy — po żadnej ze stron nie ma pola, którym jedno wskazałoby drugie.',
    komendyCudze: ['queue.create', 'queue.link'],
  },
  {
    etykieta: 'Integracja z Terminal',
    czego: 'Okno terminala należy do modułu Terminal i ma własne komendy.',
    komendyCudze: ['terminal.session.open', 'terminal.command.exec'],
  },
];

export const BRAKI_DEPLOYMENT: readonly BrakFunkcji[] = [
  {
    etykieta: 'Ustawienia domeny i DNS',
    czego: 'Domena wdrożenia musiałaby być polem wdrożenia albo osobnym bytem z komendą zapisu.',
  },
  {
    etykieta: 'Lista artefaktów budowania',
    czego: 'Wdrożenie niesie wersję i odnośnik logu — artefakt nie jest bytem umowy.',
  },
  {
    etykieta: 'Podgląd logów wdrożenia na żywo',
    czego:
      'Wdrożenie niesie logRef — sam odnośnik, bez treści. Odczyt logu wymagałby ' +
      'komendy, a podgląd na żywo dodatkowo zdarzenia strumienia.',
  },
  {
    etykieta: 'Panel stanu produktu (health check)',
    czego: 'Sprawdzenie kondycji wymaga odpytania uruchomionego produktu, którego rdzeń nie hostuje.',
  },
  {
    etykieta: 'Zmienne i dane dostępowe środowiska',
    czego: 'Dane dostępowe są bytem obszaru access; wdrożenie nie ma pola, którym by je wskazało.',
    komendyCudze: ['access.point.list', 'access.grant.list'],
  },
  {
    etykieta: 'Ustawienia skalowania',
    czego: 'Parametry skalowania musiałyby być polem wdrożenia albo osobnym bytem z komendą zapisu.',
  },
  {
    etykieta: 'Powiązanie z Automations',
    czego: 'Wyzwalacz wdrożenia jest czynnością obszaru automation, nie apps.',
    komendyCudze: ['automation.workflow.save', 'automation.schedule.set'],
  },
  {
    etykieta: 'Przegląd statusu publikacji z rdzenia',
    czego: 'Historia wdrożeń niezależna od bieżącej sesji wymaga komendy odczytu wykazu.',
    komendaZnoszaca: 'apps.deployment.list',
  },
];

/**
 * Zdanie stanu pustego okien strony dystrybucji, zanim ktokolwiek zapytał
 * rdzeń o katalog.
 *
 * Katalog rozszerzeń stoi poziom wyżej niż okno modułu: żadna komenda obszaru
 * `extension` nie niesie pola windowId, więc pustka nie bierze się z braku okna
 * — bierze się stąd, że odczyt jeszcze nie wrócił.
 */
export const BEZ_ODCZYTU_KATALOGU =
  'Katalog rozszerzeń nie był jeszcze czytany z rdzenia. Naciśnij „Odczytaj katalog", ' +
  'żeby zobaczyć, co rejestr platformy zawiera.';

/** Zdanie o oknie, którego katalog okien operacyjnych rdzenia jeszcze nie zna. */
export const OKNO_SPOZA_KATALOGU =
  'Rdzeń nie ma tego okna w katalogu okien operacyjnych, więc nie postawi go w bocznej ' +
  'nawigacji ani nie policzy w wykazie okien modułu. Okno działa mimo to, bo komendy, ' +
  'na których stoi, nie wymagają identyfikatora okna.';

export const BRAKI_APP_CATALOG: readonly BrakFunkcji[] = [
  {
    etykieta: 'Szukanie po narzędziach i znacznikach',
    czego:
      'Pozycja katalogu niesie nazwę, opis, rodzaj i wersję — udostępnianych narzędzi ' +
      'ani znaczników nie niesie, więc okno filtruje wyłącznie po tym, co dostało.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.SzukanieKatalogu,
  },
  {
    etykieta: 'Karta szczegółów rozszerzenia',
    czego:
      'Dziennik zmian, wykaz udostępnianych narzędzi i zasobów, wymagane uprawnienia ' +
      'oraz zależności nie są polami pozycji katalogu.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.SzczegolPozycji,
  },
  {
    etykieta: 'Kolekcje kuratorskie',
    czego: 'Nazwany zestaw rozszerzeń musiałby być osobnym bytem umowy z własną komendą zapisu.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.KolekcjeKuratorskie,
  },
  {
    etykieta: 'Rekomendacje Wykonawcy',
    czego: 'Dobór rozszerzeń pod opisany cel jest pracą modelu, nie odczytem katalogu.',
    komendyCudze: ['message.send'],
  },
  {
    etykieta: 'Prywatny rejestr organizacji',
    czego:
      'Pochodzenie pozycji ma dwie wartości — danaco i personal. Rejestr organizacji ' +
      'byłby trzecim źródłem, którego wyliczenie kontraktu nie zna.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.RejestrOrganizacji,
  },
];

export const BRAKI_INSTALLED_APPS: readonly BrakFunkcji[] = [
  {
    etykieta: 'Wskaźnik dostępnej aktualizacji',
    czego:
      'Pozycja niesie swoją wersję, ale nie wersję dostępną w rejestrze — nie ma czego ' +
      'z czym zestawić.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.SprawdzenieAktualizacji,
  },
  {
    etykieta: 'Przesłanie paczki z urządzenia',
    czego:
      'Instalacja przyjmuje kod, rodzaj i deklarowane źródło — przesłania treści pliku ' +
      'kanałem kontrakt nie prowadzi.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.PrzeslaniePaczki,
  },
  {
    etykieta: 'Przypięcie wersji',
    czego: 'Wersja jest polem opisowym pozycji, nie bytem z własnym stanem przypięcia.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.PrzypiecieWersji,
  },
  {
    etykieta: 'Cofnięcie do wcześniejszej wersji',
    czego:
      'Powrót do poprzedniej wersji wymaga przechowywania wcześniejszych artefaktów, ' +
      'których rejestr nie prowadzi.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.CofniecieWersji,
  },
  {
    etykieta: 'Instalacja z manifestu zestawu',
    czego: 'Zestaw rozszerzeń opisany jednym plikiem musiałby być bytem umowy.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.InstalacjaZestawu,
  },
  {
    etykieta: 'Dziennik cyklu życia',
    czego:
      'Rejestr niesie stan bieżący pozycji, nie chronologię instalacji, włączeń ' +
      'i wyłączeń, które do niego doprowadziły.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.DziennikCyklu,
  },
];

export const BRAKI_PERMISSIONS: readonly BrakFunkcji[] = [
  {
    etykieta: 'Wykaz wymaganych uprawnień',
    czego:
      'Deklarowany zakres dostępu leży w manifeście rozszerzenia, a pozycja katalogu ' +
      'niesie wyłącznie nieprzezroczystą konfigurację.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.WykazUprawnien,
  },
  {
    etykieta: 'Nadanie zakresu uprawnień',
    czego: 'Bez wykazu uprawnień nie ma czego nadawać ani czego zapisać.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.NadanieUprawnien,
  },
  {
    etykieta: 'Weryfikacja podpisu i sumy kontrolnej',
    czego:
      'Pochodzenie pozycji jest znane, podpis cyfrowy i suma kontrolna pakietu — nie; ' +
      'nie są polami pozycji katalogu.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.WeryfikacjaPodpisu,
  },
  {
    etykieta: 'Skaner manifestu',
    czego: 'Ostrzeżenie o szerokich uprawnieniach zakłada odczyt manifestu, którego nie ma.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.SkanerManifestu,
  },
  {
    etykieta: 'Audyt uprawnień i użycia',
    czego:
      'Zestawienie, kiedy i przez którego eksperta rozszerzenie było użyte, wymaga ' +
      'dziennika użycia po stronie rdzenia.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.AudytUprawnien,
  },
  {
    etykieta: 'Zakres współdzielenia referencji sekretu',
    czego: 'Referencja sekretu nie jest bytem umowy, więc nie ma czego komu udostępnić.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.WspoldzielenieSekretu,
  },
];

export const BRAKI_INTEGRATIONS_HUB: readonly BrakFunkcji[] = [
  {
    etykieta: 'Wybór i test transportu MCP',
    czego:
      'Transport da się dziś wpisać wyłącznie do nieprzezroczystej konfiguracji pozycji — ' +
      'kontrakt nie nazywa tego pola, więc rdzeń i okno mogą rozumieć je inaczej.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.WyborTransportu,
  },
  {
    etykieta: 'Kreator uwierzytelniania',
    czego:
      'Sposób uwierzytelnienia i referencja poświadczenia nie są polami pozycji katalogu; ' +
      'punkt dostępu niesie własną referencję, ale dotyczy mostu, nie integracji.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.PowiazanieDanychDostepowych,
  },
  {
    etykieta: 'Odkrywanie narzędzi, zasobów i promptów',
    czego:
      'Odpowiedź serwera na tools/list, resources/list i prompts/list nie ma w kontrakcie ' +
      'ani komendy, ani struktury, w którą miałaby wejść.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.OdkrycieNarzedzi,
  },
  {
    etykieta: 'Import definicji z OpenAPI i GraphQL',
    czego: 'Zbudowanie integracji z opisu API wymaga komendy przetwarzającej ten opis w rdzeniu.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.ImportDefinicji,
  },
  {
    etykieta: 'Webhooki przychodzące i wychodzące',
    czego: 'Adres webhooka i weryfikacja podpisu zdarzeń musiałyby być bytem umowy.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.Webhooki,
  },
  {
    etykieta: 'Mapowanie i transformacja danych',
    czego: 'Odwzorowanie pól między systemem zewnętrznym a encjami platformy nie jest bytem umowy.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.MapowanieDanych,
  },
  {
    etykieta: 'Limity szybkości, metryki użycia i koszt',
    czego: 'Liczniki wywołań i cennik integracji leżą w rdzeniu, bez komendy odczytu.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.MetrykiUzycia,
  },
  {
    etykieta: 'Zdrowie integracji',
    czego:
      'Sprawdzenie punktu dostępu bada MOST, nie usługę za nim — czas odpowiedzi ' +
      'integracji, liczba jej narzędzi i błędy powitania wymagają własnej komendy.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.ZdrowieIntegracji,
  },
  {
    etykieta: 'Rejestr referencji sekretów i rotacja',
    czego: 'Ważność tokenu i planowa rotacja klucza musiałyby być bytem umowy.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.ReferencjeSekretow,
  },
  {
    etykieta: 'Alerty o awarii integracji',
    czego: 'Reguła powiadomienia jest czynnością obszaru automation, nie katalogu rozszerzeń.',
    komendyCudze: ['automation.workflow.save', 'automation.schedule.set'],
  },
];

export const BRAKI_MCP_CONSOLE: readonly BrakFunkcji[] = [
  {
    etykieta: 'Odkryte narzędzia, zasoby i prompty',
    czego:
      'Konsola nie ma czego wyliczyć: definicje udostępniane przez serwer nie mają ' +
      'w kontrakcie ani komendy odczytu, ani struktury.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.OdkrycieNarzedzi,
  },
  {
    etykieta: 'Formularz argumentów ze schematu',
    czego:
      'Pola formularza powstają ze schematu JSON Schema narzędzia — bez odkrycia narzędzi ' +
      'nie ma z czego ich zbudować.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.OdkrycieNarzedzi,
  },
  {
    etykieta: 'Próbne wywołanie narzędzia',
    czego: 'Wykonanie tools/call wymaga komendy prowadzącej wywołanie po stronie rdzenia.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.WywolanieProbne,
  },
  {
    etykieta: 'Piaskownica testu umiejętności i wtyczki',
    czego:
      'Uruchomienie pozycji na przykładowym wejściu bez podłączania jej do eksperta ' +
      'wymaga komendy uruchomienia w izolacji.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.PiaskownicaTestu,
  },
  {
    etykieta: 'Log JSON-RPC',
    czego:
      'Strumień żądań, odpowiedzi i notyfikacji między platformą a serwerem wymaga ' +
      'komendy odczytu i zdarzenia strumienia.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.DziennikProtokolu,
  },
  {
    etykieta: 'Przekazanie podzbioru narzędzi ekspertowi',
    czego: 'Podłączenie narzędzi do eksperta jest czynnością obszaru agent, nie konsoli.',
    komendyCudze: ['agent.connector.add', 'agent.plugin.add'],
  },
  {
    etykieta: 'Log diagnostyczny rozszerzenia',
    czego: 'Dziennik zdarzeń i błędów prowadzi moduł Diagnostics i ma na to własne komendy.',
    komendyCudze: ['diagnostics.log.query', 'diagnostics.error.list'],
  },
];

export const BRAKI_PUBLISHER: readonly BrakFunkcji[] = [
  {
    etykieta: 'Kreator pakietu rozszerzenia',
    czego:
      'Zapakowanie produktu w dystrybuowalne rozszerzenie wymaga komendy wytwarzającej ' +
      'archiwum po stronie rdzenia; klient nie ma dostępu do artefaktów budowania.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.BudowaPakietu,
  },
  {
    etykieta: 'Zapis manifestu pakietu',
    czego:
      'Formularz manifestu stoi i trzyma wpisane wartości, ale nie ma dokąd ich wysłać — ' +
      'manifest nie jest bytem umowy.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.ManifestPakietu,
  },
  {
    etykieta: 'Walidacja zgodności z kontraktem rozszerzenia',
    czego: 'Reguły jednolitego kontraktu rozszerzenia zna rdzeń, nie okno.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.WalidacjaPakietu,
  },
  {
    etykieta: 'Podpisanie pakietu kluczem wydawcy',
    czego: 'Klucz wydawcy leży w warstwie sekretów i nigdy nie opuszcza rdzenia.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.PodpisPakietu,
  },
  {
    etykieta: 'Publikacja do prywatnego rejestru',
    czego: 'Rejestr organizacji nie jest dziś źródłem pozycji katalogu.',
    komendaZnoszaca: KOMENDY_PROPONOWANE.PublikacjaPakietu,
  },
];
