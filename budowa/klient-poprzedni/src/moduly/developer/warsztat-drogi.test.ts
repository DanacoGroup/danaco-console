import { beforeEach, describe, expect, it } from 'vitest';

import { Command, KOMENDY } from '../../../../shared/contract';
import { utworzOknoCodeEditora } from './okno-code-editor';
import { utworzOknoDevTools } from './okno-dev-tools';
import { utworzOknoWarsztatuKodu } from './okno-warsztatu-kodu';
import { utworzPanelHistoriiBudowan } from './panel-historii-budowan';
import { utworzPanelRunDebug } from './panel-run-debug';
import { utworzStanDevelopera, type StanDevelopera } from './stan-developer';
import { utworzZrodloProbne } from './zrodlo-probne';
import {
  utworzZrodloWarsztatuProbne,
  type ZrodloWarsztatuProbne,
} from './zrodlo-warsztatu-probne';

/**
 * Sprawdziany DRÓG z okna do rdzenia dla warsztatu modułu Developer.
 *
 * ── Czego pilnują ──────────────────────────────────────────────────────────
 * Przycisk, który niczego nie woła, wygląda dokładnie tak samo jak przycisk
 * działający — i to jest usterka, którą widać dopiero u Operatora. Sprawdziany
 * naciskają kontrolki okien i pytają, czy z każdej rodziny komend naprawdę
 * wyszło wywołanie: nie o to, czy okno się narysowało.
 *
 * ── Czego NIE sprawdzają ───────────────────────────────────────────────────
 * Nie sprawdzają treści rdzenia — od tego są sprawdziany skutku po stronie
 * serwera, które schodzą do bazy i na dysk. Tutaj mierzy się jedno: czy droga
 * istnieje i czy okno mówi prawdę, gdy rdzeń odpowiada brakiem narzędzia.
 */

/** Naciska przycisk o podanej etykiecie i oddaje obietnicę obiegu zdarzeń. */
async function nacisnij(gospodarz: HTMLElement, etykieta: string): Promise<void> {
  const przyciski = [...gospodarz.querySelectorAll('button')];
  const znaleziony = przyciski.find((przycisk) => przycisk.textContent === etykieta);
  expect(znaleziony, `w oknie nie ma przycisku „${etykieta}”`).toBeDefined();
  znaleziony?.click();
  // Czynności okien są asynchroniczne; dwa obiegi wystarczają, bo atrapa
  // odpowiada natychmiast, a okna rysują wynik w kolejnym kroku.
  await Promise.resolve();
  await Promise.resolve();
  await Promise.resolve();
}

/** Wpisuje wartość do pola o podanej etykiecie dostępności. */
function wpisz(gospodarz: HTMLElement, etykieta: string, wartosc: string): void {
  const kontrolka = gospodarz.querySelector<HTMLInputElement | HTMLTextAreaElement>(
    `[aria-label="${etykieta}"]`,
  );
  expect(kontrolka, `w oknie nie ma pola „${etykieta}”`).not.toBeNull();
  if (kontrolka !== null) kontrolka.value = wartosc;
}

/** Zdanie stanu treści okna — po nim poznajemy, co okno naprawdę powiedziało. */
function zdanieOkna(gospodarz: HTMLElement): string {
  return gospodarz.textContent ?? '';
}

describe('warsztat modułu Developer — drogi z okien do rdzenia', () => {
  let warsztat: ZrodloWarsztatuProbne;
  let stan: StanDevelopera;

  beforeEach(() => {
    warsztat = utworzZrodloWarsztatuProbne();
    stan = utworzStanDevelopera(utworzZrodloProbne(), { okno: 'okno-1', sciezka: 'suma.go' });
  });

  it('warsztat kodu prowadzi do serwera języka, formatera, wersji i operacji modelu', async () => {
    const okno = utworzOknoWarsztatuKodu(warsztat, stan);
    document.body.replaceChildren(okno.element);

    await nacisnij(okno.element, 'Nawiguj po symbolu');
    await nacisnij(okno.element, 'Formatuj plik');
    await nacisnij(okno.element, 'Uruchom analizę statyczną');
    await nacisnij(okno.element, 'Podgląd refaktoryzacji');
    await nacisnij(okno.element, 'Historia wersji pliku');
    await nacisnij(okno.element, 'Przywróć ostatnią wersję');
    await nacisnij(okno.element, 'Wykonaj operację kontekstową');
    await nacisnij(okno.element, 'Sprawdź programy warsztatu');

    expect(warsztat.wywolania).toEqual([
      'nawigujDoSymbolu',
      'formatuj',
      'analiza',
      'refaktoryzuj',
      'wersjePliku',
      'przywrocWersje',
      'operacjaKontekstowa',
      'warsztat',
    ]);
    // Plik bierze się ze stanu modułu, a nie z pola okna — tak jak wskazał go
    // Project Tree.
    expect(warsztat.ladunki['nawigujDoSymbolu']?.[0]).toMatchObject({ path: 'suma.go' });
  });

  it('brak serwera języka wraca zdaniem o braku pomiaru, a nie o braku symbolu', async () => {
    warsztat.programyWarsztatu = false;
    const okno = utworzOknoWarsztatuKodu(warsztat, stan);
    document.body.replaceChildren(okno.element);

    await nacisnij(okno.element, 'Nawiguj po symbolu');
    // To jest sedno: „nie było komu szukać” to inne zdanie niż „nie ma symbolu”.
    expect(zdanieOkna(okno.element)).toContain('nie było komu go szukać');

    await nacisnij(okno.element, 'Uruchom analizę statyczną');
    expect(zdanieOkna(okno.element)).toContain('niesprawdzone to inne zdanie');
  });

  it('okno warsztatu kodu pokazuje odmowę rdzenia, a nie pustkę', async () => {
    warsztat.odmawiaj = true;
    const okno = utworzOknoWarsztatuKodu(warsztat, stan);
    document.body.replaceChildren(okno.element);

    await nacisnij(okno.element, 'Formatuj plik');
    expect(zdanieOkna(okno.element)).toContain('atrapa odmawia czynności formatuj');
  });

  it('Run & Debug prowadzi do punktów przerwania, sesji, kroku, stosu i wyrażenia', async () => {
    const panel = utworzPanelRunDebug(warsztat, stan);
    document.body.replaceChildren(panel.element);
    wpisz(panel.element, 'Wiersz punktu przerwania', '10');
    wpisz(panel.element, 'Wyrażenie do obliczenia', 'suma');

    await nacisnij(panel.element, 'Ustaw punkt przerwania');
    await nacisnij(panel.element, 'Rozpocznij sesję debugowania');
    await nacisnij(panel.element, 'Przejdź (krok)');
    await nacisnij(panel.element, 'Wykonaj wyrażenie');
    await nacisnij(panel.element, 'Zakończ sesję');

    // Krok pociąga za sobą odczyt stosu: stan po kroku bierze się z rdzenia,
    // a nie z naciśniętego przycisku.
    expect(warsztat.wywolania).toEqual([
      'punktPrzerwania',
      'rozpocznijDebugowanie',
      'sterujDebugowaniem',
      'zakresDebugowania',
      'obliczWyrazenie',
      'sterujDebugowaniem',
    ]);
    expect(zdanieOkna(panel.element)).toContain('Sesja debugowania zakończona');
    panel.zamknij();
  });

  it('Run & Debug odmawia kroku bez sesji zamiast udawać, że go wykonał', async () => {
    const panel = utworzPanelRunDebug(warsztat, stan);
    document.body.replaceChildren(panel.element);

    await nacisnij(panel.element, 'Kontynuuj');
    expect(warsztat.wywolania).toEqual([]);
    expect(zdanieOkna(panel.element)).toContain('sesja debugowania nie została rozpoczęta');
  });

  it('historia przebiegów prowadzi do wykazu, logu, testów i pokrycia', async () => {
    const panel = utworzPanelHistoriiBudowan(warsztat, stan);
    document.body.replaceChildren(panel.element);

    await nacisnij(panel.element, 'Odczytaj historię przebiegów');
    await nacisnij(panel.element, 'Pokaż log przebiegu');
    await nacisnij(panel.element, 'Pokaż wynik testów');
    await nacisnij(panel.element, 'Pokaż pokrycie kodu');

    expect(warsztat.wywolania).toEqual([
      'wykazBudowan',
      'logBudowania',
      'wynikTestow',
      'pokrycie',
    ]);
    // Identyfikator przebiegu wypełnił się sam z wykazu — kolejne pytania
    // dotyczą tego samego biegu.
    expect(warsztat.ladunki['logBudowania']?.[0]).toMatchObject({ buildId: 'build-1' });
  });

  it('historia mówi o przycięciu logu zamiast podawać ogon jako całość', async () => {
    const panel = utworzPanelHistoriiBudowan(warsztat, stan);
    document.body.replaceChildren(panel.element);

    await nacisnij(panel.element, 'Odczytaj historię przebiegów');
    await nacisnij(panel.element, 'Pokaż log przebiegu');
    expect(zdanieOkna(panel.element)).toContain('Log jest przycięty');
  });

  it('Dev Tools prowadzi do wszystkich czternastu komend swoich czterech zakładek', async () => {
    const okno = utworzOknoDevTools(warsztat, stan);
    document.body.replaceChildren(okno.element);
    wpisz(okno.element, 'Adres zapytania', 'https://usluga.example/towary');
    wpisz(okno.element, 'Nazwa kolekcji', 'Usługa magazynu');
    wpisz(okno.element, 'Plik kontraktu OpenAPI', 'otwarte.yaml');
    wpisz(okno.element, 'Nazwa połączenia', 'magazyn');
    wpisz(okno.element, 'Baza albo plik', 'magazyn.sqlite');
    wpisz(okno.element, 'Polecenie SQL', 'SELECT 1');
    wpisz(okno.element, 'Kontener', 'kon-1');
    wpisz(okno.element, 'Plik Dockerfile', 'Dockerfile');
    wpisz(okno.element, 'Znacznik obrazu', 'danaco/usluga:1.0');
    wpisz(okno.element, 'Plik stosu', 'docker-compose.yml');

    await nacisnij(okno.element, 'Wyślij zapytanie');
    await nacisnij(okno.element, 'Zapisz do kolekcji');
    await nacisnij(okno.element, 'Importuj kontrakt OpenAPI');
    await nacisnij(okno.element, 'Odczytaj kolekcje');
    await nacisnij(okno.element, 'Zapisz połączenie');
    await nacisnij(okno.element, 'Odczytaj schemat');
    await nacisnij(okno.element, 'Wykonaj zapytanie');
    await nacisnij(okno.element, 'Uruchom migracje');
    await nacisnij(okno.element, 'Odczytaj kontenery i obrazy');
    await nacisnij(okno.element, 'Uruchom kontener');
    await nacisnij(okno.element, 'Zbuduj obraz');
    await nacisnij(okno.element, 'Podnieś stos usług');
    await nacisnij(okno.element, 'Odczytaj drzewo zależności');
    await nacisnij(okno.element, 'Skan sekretów');

    const oczekiwane = [
      'zapytanieApi',
      'zapiszKolekcje',
      'importujOpenapi',
      'kolekcje',
      'ustawPolaczenie',
      'polaczenia',
      'schemat',
      'zapytanieDanych',
      'migracje',
      'kontenery',
      'czynnoscKontenera',
      'budujObraz',
      'stosUslug',
      'zaleznosci',
      'skan',
      'znaleziska',
    ];
    for (const czynnosc of oczekiwane) {
      expect(warsztat.wywolania, `okno nie wywołało czynności ${czynnosc}`).toContain(czynnosc);
    }
  });

  it('brak silnika kontenerów wraca zdaniem o braku pomiaru, a nie pustym wykazem', async () => {
    warsztat.silnikKontenerow = false;
    const okno = utworzOknoDevTools(warsztat, stan);
    document.body.replaceChildren(okno.element);

    await nacisnij(okno.element, 'Odczytaj kontenery i obrazy');
    const zdanie = zdanieOkna(okno.element);
    expect(zdanie).toContain('nie odpowiada żaden silnik kontenerów');
    expect(zdanie).toContain('nie było kogo o nie zapytać');
  });

  it('panel akcji Code Editora prowadzi do operacji kontekstowej rdzenia', async () => {
    const zrodlo = utworzZrodloProbne();
    const stanEdytora = utworzStanDevelopera(zrodlo, { okno: 'okno-1', sciezka: 'suma.go' });
    const okno = utworzOknoCodeEditora(zrodlo, warsztat, stanEdytora);
    document.body.replaceChildren(okno.element);

    await nacisnij(okno.element, 'Wyjaśnij');
    expect(warsztat.wywolania).toEqual(['operacjaKontekstowa']);
    expect(warsztat.ladunki['operacjaKontekstowa']?.[0]).toMatchObject({
      windowId: 'okno-1',
      path: 'suma.go',
      operation: 'explain',
    });
    // Wynik jest propozycją — okno mówi to wprost, zamiast wstawiać go do pola.
    expect(zdanieOkna(okno.element)).toContain('okno niczego nie zapisało');
    okno.zamknij();
  });

  it('panel akcji odmawia operacji bez wskazanego pliku', async () => {
    const zrodlo = utworzZrodloProbne();
    const bezPliku = utworzStanDevelopera(zrodlo, { okno: 'okno-1' });
    const okno = utworzOknoCodeEditora(zrodlo, warsztat, bezPliku);
    document.body.replaceChildren(okno.element);

    await nacisnij(okno.element, 'Udokumentuj');
    expect(warsztat.wywolania).toEqual([]);
    expect(zdanieOkna(okno.element)).toContain('Operacja dotyczy pliku');
    okno.zamknij();
  });

  it('każda komenda warsztatu z kontraktu ma nazwę używaną przez klienta', () => {
    // Sprawdzian pilnuje, że klient woła komendy PO NAZWACH KONTRAKTU, a nie po
    // napisach wpisanych w kod: literówka w nazwie komendy wraca odmową
    // „nieznana komenda” dopiero u Operatora.
    const warsztatoweKomendy = [
      Command.DeveloperSymbolNavigate,
      Command.DeveloperFormatRun,
      Command.DeveloperLintGet,
      Command.DeveloperRefactorApply,
      Command.DeveloperFileVersionList,
      Command.DeveloperFileVersionRestore,
      Command.DeveloperBuildList,
      Command.DeveloperBuildLogGet,
      Command.DeveloperTestResultGet,
      Command.DeveloperCoverageGet,
      Command.DeveloperDebugSessionStart,
      Command.DeveloperDebugSessionControl,
      Command.DeveloperBreakpointSet,
      Command.DeveloperDebugScopeGet,
      Command.DeveloperDebugEvaluate,
      Command.DeveloperApiRequest,
      Command.DeveloperApiCollectionSave,
      Command.DeveloperApiCollectionList,
      Command.DeveloperApiOpenapiImport,
      Command.DeveloperDataConnectionSet,
      Command.DeveloperDataConnectionList,
      Command.DeveloperDataSchemaGet,
      Command.DeveloperDataQueryRun,
      Command.DeveloperDataMigrationRun,
      Command.DeveloperContainerList,
      Command.DeveloperContainerAction,
      Command.DeveloperImageBuild,
      Command.DeveloperComposeUp,
      Command.DeveloperDependencyList,
      Command.DeveloperScanRun,
      Command.DeveloperScanResultList,
      Command.DeveloperContextualOp,
      Command.DeveloperToolchainCheck,
    ];
    expect(warsztatoweKomendy).toHaveLength(33);
    for (const komenda of warsztatoweKomendy) {
      expect(KOMENDY as readonly string[]).toContain(komenda);
    }
  });
});
