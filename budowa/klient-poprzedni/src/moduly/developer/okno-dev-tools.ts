import {
  ContainerActionKind,
  DataEngine,
  ScanKind,
} from '../../../../shared/contract';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { pole, przyciskAkcji, poleTresci, wybor, wiersz } from '../../modele/kontrolki-formularza-braki';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci } from './stany-okna';
import { rysujZaleznosci, zaleznosci } from './zaleznosci-zewnetrzne';
import {
  cialoZakladki,
  objasnienieZakladki,
  pasekZakladki,
  utworzZakladkiOkna,
} from './zakladki-okna';
import type { ZrodloWarsztatu } from './zrodlo-warsztatu';

/**
 * Dev Tools — kolumna czterech integracji deweloperskich modułu: API Client,
 * Data Console, Containers oraz Dependencies & Security.
 *
 * ── Co się tu zmieniło ─────────────────────────────────────────────────────
 * Okno stało wcześniej jako wykaz braków: kontrakt niósł nazwy komend, lecz
 * rdzeń nie miał dla nich uchwytów, więc każda zakładka nazywała, czego brak.
 * Rdzeń ma dziś wszystkie czternaście komend tych czterech rodzin, więc zakładki
 * naprawdę je wołają. Przycisk, który tylko tłumaczy swój brak, jest właściwy
 * dokładnie do chwili, w której brak zniknie — potem staje się nieprawdą.
 *
 * ── Czego okno nie robi po cichu ───────────────────────────────────────────
 * Nie zamienia odmowy w pustkę. Wykaz kontenerów pusty i wykaz nieodczytany to
 * dwa różne zdania i okno mówi każde z nich osobno; przy braku silnika
 * kontenerów rdzeń oddaje `engineAvailable: false`, a zakładka pisze wprost, że
 * pomiaru nie było, zamiast pokazywać zero kontenerów.
 *
 * ── Poświadczenia ──────────────────────────────────────────────────────────
 * Zakładka bazodanowa nie ma pola hasła i mieć go nie będzie: hasło leży
 * w sejfie, a wiersz połączenia niesie do niego ODWOŁANIE. Pole hasła
 * w przeglądarce byłoby kopią sekretu w miejscu, którego nikt nie rotuje.
 */
export interface OknoDevTools {
  element: HTMLElement;
  odswiez(): void;
  /** Okno nie zakłada nasłuchu; metoda stoi dla jednolitości złożenia modułu. */
  zamknij(): void;
}

export function utworzOknoDevTools(
  zrodlo: ZrodloWarsztatu,
  stan: StanDevelopera,
): OknoDevTools {
  const rama = utworzRameOkna({
    tytul: 'Dev Tools',
    rola: 'pomocnicze',
    kod: 'dev-tools',
    przeznaczenie:
      'Kolumna czterech integracji deweloperskich: klient API, konsola bazodanowa, kontenery ' +
      'oraz zależności i bezpieczeństwo.',
    modul: 'Developer',
    przedrostek: 'mdev',
  });

  const api = zakladkaApiClient(zrodlo, stan);
  const dane = zakladkaDataConsole(zrodlo, stan);
  const kontenery = zakladkaContainers(zrodlo, stan);
  const bezpieczenstwo = zakladkaBezpieczenstwa(zrodlo, stan);

  const zakladki = utworzZakladkiOkna('Integracje deweloperskie', [
    { kod: 'api-client', nazwa: 'API Client', element: api.element },
    { kod: 'data-console', nazwa: 'Data Console', element: dane.element },
    { kod: 'containers', nazwa: 'Containers', element: kontenery.element },
    { kod: 'security', nazwa: 'Dependencies & Security', element: bezpieczenstwo.element },
  ]);

  rama.narzedzia.append(
    objasnienieZakladki(
      'Widocznością zakładek steruje wedle opracowania klucz developer.narzedzia.* okna ' +
        'konfiguracji. Rdzeń nie zna dziś ani jednego ustawienia obszaru developer, więc widoczne ' +
        'są wszystkie cztery — ukrycie zakładki opierałoby się wtedy na wartości zmyślonej.',
    ),
  );
  rama.pasek.append(zakladki.pasek);
  rama.cialo.append(zakladki.obszary);

  return {
    element: rama.element,
    odswiez() {
      // Odświeżenie czyta to, co ma stan trwały po stronie rdzenia: kolekcje
      // zapytań i opisy połączeń. Zapytania, skanów ani czynności na kontenerach
      // nie uruchamia się samoczynnie — to są czynności Operatora, a nie odczyt.
      api.odswiez();
      dane.odswiez();
    },
    zamknij: () => undefined,
  };
}

/** Jedna zakładka okna wraz z jej odczytem. */
interface ZakladkaWarsztatu {
  element: HTMLElement;
  odswiez(): void;
}

/**
 * API Client — zapytania HTTP, kolekcje i kontrakty OpenAPI.
 *
 * Zapytanie jedzie przez rdzeń, nie z przeglądarki. Klient siedzi po drugiej
 * stronie gniazda niż katalog roboczy sesji, a dostęp sieciowy podlega zakresowi
 * izolacji okna — zapytanie wysłane z przeglądarki omijałoby jedno i drugie
 * i pokazywałoby wynik z cudzej sieci pod nazwą wyniku sesji.
 */
function zakladkaApiClient(zrodlo: ZrodloWarsztatu, stan: StanDevelopera): ZakladkaWarsztatu {
  const tresc = utworzStanTresci();
  const metoda = wybor('Metoda zapytania', [
    ['GET', 'GET'],
    ['POST', 'POST'],
    ['PUT', 'PUT'],
    ['PATCH', 'PATCH'],
    ['DELETE', 'DELETE'],
  ]);
  const adres = pole('Adres zapytania', 'https://usluga.example/zasob albo {{baseUrl}}/zasob');
  const naglowki = poleTresci('Nagłówki zapytania (JSON)', 3, '{"Authorization":"Bearer …"}');
  const cialo = poleTresci('Treść zapytania', 4, '{"nazwa":"wartość"}');
  const srodowisko = pole('Środowisko kolekcji', 'nazwa środowiska; puste znaczy bez podstawień');
  const nazwaKolekcji = pole('Nazwa kolekcji', 'Usługa rozliczeń');
  const kontrakt = pole('Plik kontraktu OpenAPI', 'otwarte.yaml');

  const wyslij = przyciskAkcji('Wyślij zapytanie');
  const zapisz = przyciskAkcji('Zapisz do kolekcji');
  const importuj = przyciskAkcji('Importuj kontrakt OpenAPI');
  const odczytaj = przyciskAkcji('Odczytaj kolekcje');

  wyslij.addEventListener('click', () => {
    void wykonajZapytanie();
  });
  zapisz.addEventListener('click', () => {
    void zapiszDoKolekcji();
  });
  importuj.addEventListener('click', () => {
    void importujKontrakt();
  });
  odczytaj.addEventListener('click', () => {
    void odczytajKolekcje();
  });

  async function wykonajZapytanie(): Promise<void> {
    const cel = adres.value.trim();
    if (cel === '') {
      tresc.blad('Zapytanie wymaga adresu — pole adresu jest puste.');
      return;
    }
    tresc.ladowanie('Zapytanie w toku…');
    const wynik = await zrodlo.zapytanieApi(zadanieZapytania(cel));
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Wykonanie zapytania', wynik.blad), wynik.blad);
      return;
    }
    const odpowiedz = wynik.wynik;
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapit(
        `Stan ${String(odpowiedz.status)}${odpowiedz.statusText === undefined ? '' : ` (${odpowiedz.statusText})`}` +
          ` · czas ${String(odpowiedz.durationMs)} ms` +
          (odpowiedz.sizeBytes === undefined ? '' : ` · ${String(odpowiedz.sizeBytes)} B`),
      ),
      blokTekstu(odpowiedz.body ?? '(odpowiedź bez treści)'),
    );
  }

  function zadanieZapytania(cel: string): Parameters<ZrodloWarsztatu['zapytanieApi']>[0] {
    const zadanie: Parameters<ZrodloWarsztatu['zapytanieApi']>[0] = {
      windowId: stan.okno(),
      method: metoda.value,
      url: cel,
    };
    const wpisy = naglowkiZadania(naglowki.value);
    if (wpisy !== undefined) zadanie.headers = wpisy;
    if (cialo.value.trim() !== '') {
      zadanie.body = cialo.value;
      zadanie.bodyKind = 'json';
    }
    if (srodowisko.value.trim() !== '') zadanie.environmentId = srodowisko.value.trim();
    return zadanie;
  }

  async function zapiszDoKolekcji(): Promise<void> {
    const nazwa = nazwaKolekcji.value.trim();
    if (nazwa === '') {
      tresc.blad('Kolekcja wymaga nazwy — pole nazwy kolekcji jest puste.');
      return;
    }
    tresc.ladowanie('Zapis kolekcji…');
    const wynik = await zrodlo.zapiszKolekcje({
      windowId: stan.okno(),
      name: nazwa,
      requests: [
        {
          name: `${metoda.value} ${adres.value.trim()}`,
          method: metoda.value,
          url: adres.value.trim(),
          body: cialo.value,
        },
      ],
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Zapis kolekcji', wynik.blad), wynik.blad);
      return;
    }
    tresc.potwierdzenie(`Kolekcja ${wynik.wynik.name} zapisana.`, true);
  }

  async function importujKontrakt(): Promise<void> {
    const sciezka = kontrakt.value.trim();
    if (sciezka === '') {
      tresc.blad('Import wymaga wskazania pliku kontraktu w repozytorium.');
      return;
    }
    tresc.ladowanie('Odczyt kontraktu…');
    const wynik = await zrodlo.importujOpenapi({ windowId: stan.okno(), path: sciezka });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Import kontraktu OpenAPI', wynik.blad), wynik.blad);
      return;
    }
    tresc.potwierdzenie(
      `Kolekcja ${wynik.wynik.collection.name} — ${String(wynik.wynik.requestCount)} zapytań z kontraktu.`,
      true,
    );
  }

  async function odczytajKolekcje(): Promise<void> {
    tresc.ladowanie('Odczyt kolekcji…');
    const wynik = await zrodlo.kolekcje({ windowId: stan.okno() });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt kolekcji zapytań', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.collections.length === 0) {
      tresc.pusto('Okno nie ma jeszcze ani jednej kolekcji zapytań.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.collections.map((kolekcja) =>
        akapit(`${kolekcja.name} — ${kolekcja.id}`),
      ),
    );
  }

  const element = cialoZakladki(
    objasnienieZakladki(
      'Zapytanie wykonuje rdzeń w sieci okna, nie przeglądarka: klient stoi po drugiej stronie ' +
        'gniazda niż katalog roboczy sesji, a dostęp sieciowy podlega zakresowi izolacji okna. ' +
        'Zmienne środowiska podstawiają się zapisem {{nazwa}} w adresie, nagłówkach i treści.',
    ),
    wiersz('Metoda', metoda, { klasa: 'dn-pole' }),
    wiersz('Adres', adres, { klasa: 'dn-pole' }),
    wiersz('Nagłówki (JSON)', naglowki, { klasa: 'dn-pole' }),
    wiersz('Treść', cialo, { klasa: 'dn-pole' }),
    wiersz('Środowisko', srodowisko, {
      klasa: 'dn-pole',
      objasnienie: 'Nazwa środowiska zapisanego w którejś kolekcji okna.',
    }),
    wiersz('Nazwa kolekcji', nazwaKolekcji, { klasa: 'dn-pole' }),
    wiersz('Plik kontraktu OpenAPI', kontrakt, {
      klasa: 'dn-pole',
      objasnienie: 'Ścieżka w katalogu roboczym okna; kontrakt wersjonowany z kodem.',
    }),
    pasekZakladki(wyslij, zapisz, importuj, odczytaj),
    tresc.element,
  );

  return { element, odswiez: () => void odczytajKolekcje() };
}

/** naglowkiZadania czyta nagłówki z pola; treść nieczytelna zostaje pominięta. */
function naglowkiZadania(zapis: string): Record<string, string> | undefined {
  const tresc = zapis.trim();
  if (tresc === '') return undefined;
  try {
    const wpisy: unknown = JSON.parse(tresc);
    if (typeof wpisy !== 'object' || wpisy === null) return undefined;
    return wpisy as Record<string, string>;
  } catch {
    return undefined;
  }
}

/**
 * Data Console — połączenia, schemat, zapytania i migracje.
 *
 * Zakładka nie ma pola hasła. Poświadczenie wchodzi ODWOŁANIEM do sejfu, tak
 * samo jak klucze dostawców modeli; hasło wpisane w przeglądarce byłoby kopią
 * sekretu w miejscu, którego nikt nie rotuje.
 */
function zakladkaDataConsole(zrodlo: ZrodloWarsztatu, stan: StanDevelopera): ZakladkaWarsztatu {
  const tresc = utworzStanTresci();
  const nazwa = pole('Nazwa połączenia', 'magazyn');
  const silnik = wybor('Silnik bazy', [
    [DataEngine.Postgres, 'PostgreSQL'],
    [DataEngine.Mysql, 'MySQL'],
    [DataEngine.Sqlite, 'SQLite'],
  ]);
  const host = pole('Host', 'localhost');
  const baza = pole('Baza albo plik', 'magazyn albo dane/magazyn.sqlite');
  const uzytkownik = pole('Użytkownik', 'danaco');
  const poswiadczenie = pole('Odwołanie do poświadczenia', 'nazwa wpisu sejfu');
  const wybranePolaczenie = wybor('Połączenie', []);
  const zapytanie = poleTresci('Polecenie SQL', 4, 'SELECT * FROM towar');
  const katalogMigracji = pole('Katalog migracji', 'migracje');

  const zapisz = przyciskAkcji('Zapisz połączenie');
  const schemat = przyciskAkcji('Odczytaj schemat');
  const wykonaj = przyciskAkcji('Wykonaj zapytanie');
  const plan = przyciskAkcji('Pokaż plan wykonania');
  const migruj = przyciskAkcji('Uruchom migracje');

  zapisz.addEventListener('click', () => void zapiszPolaczenie());
  schemat.addEventListener('click', () => void odczytajSchemat());
  wykonaj.addEventListener('click', () => void wykonajZapytanie(false));
  plan.addEventListener('click', () => void wykonajZapytanie(true));
  migruj.addEventListener('click', () => void uruchomMigracje());

  async function zapiszPolaczenie(): Promise<void> {
    if (nazwa.value.trim() === '' || baza.value.trim() === '') {
      tresc.blad('Połączenie wymaga nazwy oraz wskazania bazy albo pliku.');
      return;
    }
    tresc.ladowanie('Zapis połączenia…');
    const zadanie: Parameters<ZrodloWarsztatu['ustawPolaczenie']>[0] = {
      windowId: stan.okno(),
      name: nazwa.value.trim(),
      engine: silnik.value as DataEngine,
      database: baza.value.trim(),
    };
    if (host.value.trim() !== '') zadanie.host = host.value.trim();
    if (uzytkownik.value.trim() !== '') zadanie.user = uzytkownik.value.trim();
    if (poswiadczenie.value.trim() !== '') zadanie.credentialRef = poswiadczenie.value.trim();
    const wynik = await zrodlo.ustawPolaczenie(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Zapis połączenia', wynik.blad), wynik.blad);
      return;
    }
    tresc.potwierdzenie(`Połączenie ${wynik.wynik.name} zapisane.`, true);
    await odczytajPolaczenia();
  }

  async function odczytajPolaczenia(): Promise<void> {
    const wynik = await zrodlo.polaczenia({ windowId: stan.okno() });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt połączeń', wynik.blad), wynik.blad);
      return;
    }
    wybranePolaczenie.replaceChildren();
    for (const polaczenie of wynik.wynik.connections) {
      const pozycja = document.createElement('option');
      pozycja.value = polaczenie.id;
      pozycja.textContent = `${polaczenie.name} (${polaczenie.engine}${polaczenie.readOnly ? ', tylko odczyt' : ''})`;
      wybranePolaczenie.append(pozycja);
    }
  }

  async function odczytajSchemat(): Promise<void> {
    const kod = wybranePolaczenie.value;
    if (kod === '') {
      tresc.blad('Odczyt schematu wymaga wskazania połączenia.');
      return;
    }
    tresc.ladowanie('Odczyt schematu…');
    const wynik = await zrodlo.schemat({ connectionId: kod });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt schematu', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.nodes.length === 0) {
      tresc.pusto('Połączenie nie wykazało ani jednej tabeli.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.nodes.map((wezel) =>
        akapit(
          `${wezel.path} — ${wezel.kind}${wezel.dataType === undefined ? '' : ` (${wezel.dataType})`}`,
        ),
      ),
    );
  }

  async function wykonajZapytanie(zPlanem: boolean): Promise<void> {
    const kod = wybranePolaczenie.value;
    if (kod === '' || zapytanie.value.trim() === '') {
      tresc.blad('Wykonanie wymaga wskazania połączenia i treści polecenia SQL.');
      return;
    }
    tresc.ladowanie('Wykonanie polecenia…');
    const wynik = await zrodlo.zapytanieDanych({
      connectionId: kod,
      sql: zapytanie.value,
      explain: zPlanem,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Wykonanie polecenia SQL', wynik.blad), wynik.blad);
      return;
    }
    const siatka = wynik.wynik;
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapit(
        `Wierszy: ${String(siatka.rowCount)}` +
          (siatka.affectedRows === undefined ? '' : ` · zmieniono ${String(siatka.affectedRows)}`) +
          ` · czas ${String(siatka.durationMs)} ms` +
          (siatka.truncated ? ' · wynik przycięty' : ''),
      ),
      blokTekstu(
        siatka.columns.join(' | ') + '\n' + JSON.stringify(siatka.rows, null, 2),
      ),
    );
  }

  async function uruchomMigracje(): Promise<void> {
    const kod = wybranePolaczenie.value;
    if (kod === '') {
      tresc.blad('Migracje wymagają wskazania połączenia.');
      return;
    }
    tresc.ladowanie('Migracje w toku…');
    const wynik = await zrodlo.migracje({
      connectionId: kod,
      windowId: stan.okno(),
      direction: 'up',
      target: katalogMigracji.value.trim() === '' ? undefined : katalogMigracji.value.trim(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Uruchomienie migracji', wynik.blad), wynik.blad);
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapit(
        `Zastosowano: ${String(wynik.wynik.applied.length)} · oczekuje: ${String(wynik.wynik.pending.length)}`,
      ),
      blokTekstu(wynik.wynik.output ?? '(rdzeń nie podał opisu przebiegu)'),
    );
  }

  const element = cialoZakladki(
    objasnienieZakladki(
      'Hasła nie ma w tym oknie i nie będzie: wiersz połączenia niesie ODWOŁANIE do wpisu sejfu, ' +
        'a nie sam sekret. Połączenie zakładane jest jako tylko do odczytu — zdjęcie tej nastawy ' +
        'jest osobnym rozstrzygnięciem, bo konsola SQL nad produkcją bez niej jest jedną ' +
        'nieuwagą od szkody.',
    ),
    rysujZaleznosci(zaleznosci(['serwer-bazy'])),
    wiersz('Nazwa', nazwa, { klasa: 'dn-pole' }),
    wiersz('Silnik', silnik, { klasa: 'dn-pole' }),
    wiersz('Host', host, { klasa: 'dn-pole' }),
    wiersz('Baza albo plik', baza, { klasa: 'dn-pole' }),
    wiersz('Użytkownik', uzytkownik, { klasa: 'dn-pole' }),
    wiersz('Odwołanie do poświadczenia', poswiadczenie, {
      klasa: 'dn-pole',
      objasnienie: 'Nazwa wpisu sejfu; hasło nigdy nie przechodzi przez to okno.',
    }),
    wiersz('Połączenie', wybranePolaczenie, { klasa: 'dn-pole' }),
    wiersz('Polecenie SQL', zapytanie, { klasa: 'dn-pole' }),
    wiersz('Katalog migracji', katalogMigracji, { klasa: 'dn-pole' }),
    pasekZakladki(zapisz, schemat, wykonaj, plan, migruj),
    tresc.element,
  );

  return { element, odswiez: () => void odczytajPolaczenia() };
}

/** Containers — kontenery, obrazy i stosy usług silnika kontenerów. */
function zakladkaContainers(zrodlo: ZrodloWarsztatu, stan: StanDevelopera): ZakladkaWarsztatu {
  const tresc = utworzStanTresci();
  const kontener = pole('Kontener', 'identyfikator albo nazwa kontenera');
  const dockerfile = pole('Plik Dockerfile', 'Dockerfile');
  const znacznik = pole('Znacznik obrazu', 'danaco/usluga:1.0');
  const plikStosu = pole('Plik stosu', 'docker-compose.yml');

  const odczytaj = przyciskAkcji('Odczytaj kontenery i obrazy');
  const uruchom = przyciskAkcji('Uruchom kontener');
  const zatrzymaj = przyciskAkcji('Zatrzymaj kontener');
  const log = przyciskAkcji('Pokaż log kontenera');
  const zbuduj = przyciskAkcji('Zbuduj obraz');
  const podnies = przyciskAkcji('Podnieś stos usług');
  const opusc = przyciskAkcji('Zatrzymaj stos usług');

  odczytaj.addEventListener('click', () => void odczytajKontenery());
  uruchom.addEventListener('click', () => void czynnosc(ContainerActionKind.Start));
  zatrzymaj.addEventListener('click', () => void czynnosc(ContainerActionKind.Stop));
  log.addEventListener('click', () => void czynnosc(ContainerActionKind.Logs));
  zbuduj.addEventListener('click', () => void budujObraz());
  podnies.addEventListener('click', () => void stos(false));
  opusc.addEventListener('click', () => void stos(true));

  async function odczytajKontenery(): Promise<void> {
    tresc.ladowanie('Odczyt kontenerów…');
    const wynik = await zrodlo.kontenery({ windowId: stan.okno(), all: true, includeImages: true });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt kontenerów', wynik.blad), wynik.blad);
      return;
    }
    // Brak silnika NIE jest pustką: „nic nie biegnie” i „nie było czym zapytać”
    // to dwa różne zdania o tej samej maszynie.
    if (!wynik.wynik.engineAvailable) {
      tresc.pusto(
        'Na serwerze nie odpowiada żaden silnik kontenerów (Docker ani Podman), więc pomiaru nie ' +
          'było. To nie znaczy, że kontenerów nie ma — znaczy, że nie było kogo o nie zapytać. ' +
          'Instalacja silnika po stronie serwera jest rozstrzygnięciem Właściciela.',
      );
      return;
    }
    if (wynik.wynik.containers.length === 0 && (wynik.wynik.images ?? []).length === 0) {
      tresc.pusto('Silnik odpowiedział: nie ma ani jednego kontenera ani obrazu.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.containers.map((pozycja) =>
        akapit(
          `${pozycja.name} — ${pozycja.status}${pozycja.image === undefined ? '' : ` · ${pozycja.image}`}` +
            `${pozycja.ports === undefined ? '' : ` · porty ${pozycja.ports.join(', ')}`}`,
        ),
      ),
      ...(wynik.wynik.images ?? []).map((obraz) =>
        akapit(`obraz ${(obraz.tags ?? [obraz.id]).join(', ')}`),
      ),
    );
  }

  async function czynnosc(rodzaj: ContainerActionKind): Promise<void> {
    if (kontener.value.trim() === '') {
      tresc.blad('Czynność wymaga wskazania kontenera.');
      return;
    }
    tresc.ladowanie('Czynność na kontenerze…');
    const wynik = await zrodlo.czynnoscKontenera({
      containerId: kontener.value.trim(),
      action: rodzaj,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Czynność na kontenerze', wynik.blad), wynik.blad);
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapit(`${wynik.wynik.container.name} — ${wynik.wynik.container.status}`),
      blokTekstu(wynik.wynik.output ?? '(czynność bez wyjścia)'),
    );
  }

  async function budujObraz(): Promise<void> {
    if (dockerfile.value.trim() === '' || znacznik.value.trim() === '') {
      tresc.blad('Budowanie obrazu wymaga pliku Dockerfile i znacznika.');
      return;
    }
    tresc.ladowanie('Budowanie obrazu…');
    const wynik = await zrodlo.budujObraz({
      windowId: stan.okno(),
      dockerfile: dockerfile.value.trim(),
      tag: znacznik.value.trim(),
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Budowanie obrazu', wynik.blad), wynik.blad);
      return;
    }
    tresc.potwierdzenie(`Obraz zbudowany: ${wynik.wynik.imageId}`, true);
  }

  async function stos(opuszczenie: boolean): Promise<void> {
    if (plikStosu.value.trim() === '') {
      tresc.blad('Stos usług wymaga wskazania pliku opisu.');
      return;
    }
    tresc.ladowanie(opuszczenie ? 'Zatrzymanie stosu…' : 'Podnoszenie stosu…');
    const wynik = await zrodlo.stosUslug({
      windowId: stan.okno(),
      file: plikStosu.value.trim(),
      down: opuszczenie,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Czynność na stosie usług', wynik.blad), wynik.blad);
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.services.map((usluga) => akapit(`${usluga.name} — ${usluga.status}`)),
      blokTekstu(wynik.wynik.output ?? '(stos bez opisu przebiegu)'),
    );
  }

  const element = cialoZakladki(
    objasnienieZakladki(
      'Rozmowa idzie zestawem narzędziowym silnika wkompilowanym w rdzeń, przez gniazdo Dockera ' +
        'albo Podmana — nie przez program wiersza poleceń. Silnika nie da się wkompilować, więc ' +
        'jego brak wraca odpowiedzią mówiącą wprost, że pomiaru nie było.',
    ),
    objasnienieZakladki(
      'Powłoka w kontenerze nie należy do tej zakładki: dostarcza ją warstwa wykonawcza modułu ' +
        'Terminal, którego karta stoi w pasie okien pomocniczych tego modułu. Druga droga do ' +
        'powłoki byłaby drugą prawdą o tym samym procesie.',
    ),
    rysujZaleznosci(zaleznosci(['silnik-kontenerow'])),
    wiersz('Kontener', kontener, { klasa: 'dn-pole' }),
    wiersz('Plik Dockerfile', dockerfile, { klasa: 'dn-pole' }),
    wiersz('Znacznik obrazu', znacznik, { klasa: 'dn-pole' }),
    wiersz('Plik stosu', plikStosu, { klasa: 'dn-pole' }),
    pasekZakladki(odczytaj, uruchom, zatrzymaj, log, zbuduj, podnies, opusc),
    tresc.element,
  );

  return { element, odswiez: () => undefined };
}

/** Dependencies & Security — zależności, podatności, sekrety i licencje. */
function zakladkaBezpieczenstwa(
  zrodlo: ZrodloWarsztatu,
  stan: StanDevelopera,
): ZakladkaWarsztatu {
  const tresc = utworzStanTresci();
  const manifest = pole('Manifest zależności', 'go.mod, package.json, requirements.txt, Cargo.toml');

  const drzewo = przyciskAkcji('Odczytaj drzewo zależności');
  const skanZaleznosci = przyciskAkcji('Skan zależności');
  const skanSekretow = przyciskAkcji('Skan sekretów');
  const skanKodu = przyciskAkcji('Analiza bezpieczeństwa kodu');
  const skanLicencji = przyciskAkcji('Analiza licencji');
  const wykazZnalezisk = przyciskAkcji('Pokaż znaleziska okna');

  drzewo.addEventListener('click', () => void odczytajZaleznosci());
  skanZaleznosci.addEventListener('click', () => void uruchomSkan([ScanKind.Dependencies]));
  skanSekretow.addEventListener('click', () => void uruchomSkan([ScanKind.Secrets]));
  skanKodu.addEventListener('click', () => void uruchomSkan([ScanKind.Code]));
  skanLicencji.addEventListener('click', () => void uruchomSkan([ScanKind.Licenses]));
  wykazZnalezisk.addEventListener('click', () => void pokazZnaleziska());

  async function odczytajZaleznosci(): Promise<void> {
    tresc.ladowanie('Odczyt manifestu…');
    const zadanie: Parameters<ZrodloWarsztatu['zaleznosci']>[0] = { windowId: stan.okno() };
    if (manifest.value.trim() !== '') zadanie.manifest = manifest.value.trim();
    const wynik = await zrodlo.zaleznosci(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt zależności', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.dependencies.length === 0) {
      tresc.pusto(`Manifest ${wynik.wynik.manifest} nie wykazuje ani jednej zależności.`);
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapit(`Manifest: ${wynik.wynik.manifest}`),
      ...wynik.wynik.dependencies.map((pozycja) =>
        akapit(
          `${pozycja.name} ${pozycja.version}${pozycja.direct ? '' : ' (pośrednia)'}` +
            `${pozycja.license === undefined ? '' : ` · ${pozycja.license}`}`,
        ),
      ),
    );
  }

  async function uruchomSkan(rodzaje: readonly ScanKind[]): Promise<void> {
    tresc.ladowanie('Skanowanie repozytorium…');
    const wynik = await zrodlo.skan({ windowId: stan.okno(), kinds: [...rodzaje] });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Uruchomienie skanu', wynik.blad), wynik.blad);
      return;
    }
    await pokazZnaleziska(wynik.wynik.id);
  }

  async function pokazZnaleziska(skan?: string): Promise<void> {
    const zadanie: Parameters<ZrodloWarsztatu['znaleziska']>[0] =
      skan === undefined ? { windowId: stan.okno() } : { scanId: skan };
    const wynik = await zrodlo.znaleziska(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Odczyt znalezisk', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.findings.length === 0) {
      // Zdanie mówi o SKANIE, nie o repozytorium: pusty wykaz po przebiegu
      // znaczy „sprawdzono i nie znaleziono”, a to jest inne zdanie niż
      // „repozytorium jest czyste”, którego okno orzec nie może.
      tresc.pusto('Przebieg skanowania zakończył się bez znalezisk w sprawdzonym zakresie.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.findings.map((znalezisko) =>
        akapit(
          `[${znalezisko.severity}] ${znalezisko.title}` +
            `${znalezisko.path === undefined ? '' : ` — ${znalezisko.path}`}` +
            `${znalezisko.line === undefined ? '' : `:${String(znalezisko.line)}`}` +
            `${znalezisko.description === undefined ? '' : ` — ${znalezisko.description}`}`,
        ),
      ),
    );
  }

  const element = cialoZakladki(
    objasnienieZakladki(
      'Wykaz zależności powstaje z odczytu manifestu repozytorium, a nie z wywołania narzędzia ' +
        'języka: przeglądarka zależności ma się otwierać także w repozytorium, którego nikt ' +
        'jeszcze nie zbudował.',
    ),
    objasnienieZakladki(
      'Skan sekretów i kodu rozpoznaje kształty regułami wyrażeń regularnych, a nazwa reguły ' +
        'stoi w znalezisku — widać więc, co sprawdzono. Treści znalezionego sekretu rdzeń NIE ' +
        'zapisuje: baza produktu nie ma być drugim miejscem, w którym ten klucz leży.',
    ),
    rysujZaleznosci(zaleznosci(['skanery-bezpieczenstwa'])),
    wiersz('Manifest', manifest, {
      klasa: 'dn-pole',
      objasnienie: 'Puste znaczy manifest rozpoznany w katalogu roboczym okna.',
    }),
    pasekZakladki(drzewo, skanZaleznosci, skanSekretow, skanKodu, skanLicencji, wykazZnalezisk),
    tresc.element,
  );

  return { element, odswiez: () => undefined };
}

/** akapit składa jeden wiersz treści okna. */
function akapit(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mdev-wiersz';
  element.textContent = zdanie;
  return element;
}

/** blokTekstu składa miejsce na treść wielowierszową — log, wynik, plan. */
function blokTekstu(zawartosc: string): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'mdev-blok';
  element.textContent = zawartosc;
  return element;
}
