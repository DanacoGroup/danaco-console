import {
  ContextualOpKind,
  RefactorKind,
  SymbolNavigationKind,
} from '../../../../shared/contract';
import { opisOdmowyBledu } from '../../komponenty/odmowa';
import { utworzRameOkna } from '../../komponenty/rama-okna';
import {
  pole,
  poleLiczbowe,
  poleTresci,
  przyciskAkcji,
  wybor,
  wiersz,
} from '../../modele/kontrolki-formularza-braki';
import type { StanDevelopera } from './stan-developer';
import { utworzStanTresci } from './stany-okna';
import { objasnienieZakladki, pasekZakladki } from './zakladki-okna';
import { rysujZaleznosci, zaleznosci } from './zaleznosci-zewnetrzne';
import type { ZrodloWarsztatu } from './zrodlo-warsztatu';

/**
 * Warsztat kodu — pasek operacji Code Editora wyniesiony do własnego okna:
 * nawigacja po symbolach, formatowanie, analiza statyczna, refaktoryzacje,
 * historia wersji pliku, operacje kontekstowe modelu oraz sonda programów
 * warsztatu na serwerze.
 *
 * ── Dlaczego osobne okno, a nie pasek w edytorze ───────────────────────────
 * Opracowanie opisuje te czynności jako pasek pływający nad zaznaczeniem.
 * Pasek pływający wymaga zaznaczenia w polu edycji, a pole edycji Code Editora
 * jest zwykłym obszarem tekstu bez modelu dokumentu — zaznaczenie znika przy
 * pierwszym kliknięciu poza nim. Okno jest tym samym zestawem czynności podanym
 * w postaci, która nie potrzebuje żywego zaznaczenia: plik bierze się ze stanu
 * modułu, a zaznaczenie wpisuje się jawnie.
 *
 * ── Wynik jest propozycją, nie zapisem ─────────────────────────────────────
 * Refaktoryzacja i operacje modelu wracają jako podgląd zmiany. Zapis na dysk
 * jest osobnym rozstrzygnięciem Operatora, bo to ta jedna chwila, w której da
 * się pracę modelu odrzucić.
 *
 * ── Brak programu nie jest odmową ──────────────────────────────────────────
 * Nawigacja po symbolach i analiza statyczna niosą w odpowiedzi jawne pole
 * mówiące, czy serwer języka i linter były osiągalne. Okno rozróżnia „nie ma
 * wystąpień” od „nie było czym sprawdzić” — pierwsze naprawia się w kodzie,
 * drugie instalacją po stronie serwera.
 */
export interface OknoWarsztatuKodu {
  element: HTMLElement;
  odswiez(): void;
  zamknij(): void;
}

export function utworzOknoWarsztatuKodu(
  zrodlo: ZrodloWarsztatu,
  stan: StanDevelopera,
): OknoWarsztatuKodu {
  const rama = utworzRameOkna({
    tytul: 'Warsztat kodu',
    rola: 'pomocnicze',
    kod: 'warsztat-kodu',
    przeznaczenie:
      'Nawigacja po symbolach, formatowanie, analiza statyczna, refaktoryzacje, historia wersji ' +
      'pliku oraz operacje kontekstowe modelu nad plikiem otwartym w Code Editorze.',
    modul: 'Developer',
    przedrostek: 'mdev',
  });
  const tresc = utworzStanTresci();

  const wierszSymbolu = poleLiczbowe('Wiersz', '1');
  const kolumnaSymbolu = poleLiczbowe('Kolumna', '1');
  const rodzajNawigacji = wybor('Rodzaj nawigacji', [
    [SymbolNavigationKind.Definition, 'Przejdź do definicji'],
    [SymbolNavigationKind.Implementation, 'Przejdź do implementacji'],
    [SymbolNavigationKind.References, 'Znajdź wystąpienia'],
    [SymbolNavigationKind.DocumentSymbol, 'Symbole dokumentu'],
  ]);
  const rodzajRefaktoryzacji = wybor('Rodzaj refaktoryzacji', [
    [RefactorKind.Rename, 'Zmień nazwę symbolu'],
    [RefactorKind.OrganizeImports, 'Uporządkuj importy'],
    [RefactorKind.ExtractFunction, 'Wyodrębnij funkcję'],
    [RefactorKind.ExtractVariable, 'Wyodrębnij zmienną'],
    [RefactorKind.Inline, 'Wstaw w miejscu'],
    [RefactorKind.Move, 'Przenieś symbol'],
  ]);
  const nowaNazwa = pole('Nowa nazwa symbolu', 'NowaNazwa');
  const rodzajOperacji = wybor('Operacja kontekstowa', [
    [ContextualOpKind.Explain, 'Wyjaśnij'],
    [ContextualOpKind.Refactor, 'Refaktoryzuj'],
    [ContextualOpKind.Document, 'Udokumentuj'],
    [ContextualOpKind.Fix, 'Napraw'],
    [ContextualOpKind.WriteTest, 'Napisz sprawdzian'],
    [ContextualOpKind.Optimize, 'Zoptymalizuj'],
    [ContextualOpKind.Generate, 'Wygeneruj'],
    [ContextualOpKind.Convert, 'Przepisz na inny język'],
  ]);
  const zaznaczenie = poleTresci('Zaznaczenie', 4, 'fragment kodu, którego dotyczy operacja');
  const wskazowka = pole('Wskazanie Operatora', 'np. nowa nazwa albo język docelowy');

  const nawiguj = przyciskAkcji('Nawiguj po symbolu');
  const formatuj = przyciskAkcji('Formatuj plik');
  const analizuj = przyciskAkcji('Uruchom analizę statyczną');
  const refaktoryzuj = przyciskAkcji('Podgląd refaktoryzacji');
  const wersje = przyciskAkcji('Historia wersji pliku');
  const przywroc = przyciskAkcji('Przywróć ostatnią wersję');
  const operacja = przyciskAkcji('Wykonaj operację kontekstową');
  const warsztat = przyciskAkcji('Sprawdź programy warsztatu');

  let ostatniaWersja = '';

  nawiguj.addEventListener('click', () => void nawigujPoSymbolu());
  formatuj.addEventListener('click', () => void formatujPlik());
  analizuj.addEventListener('click', () => void uruchomAnalize());
  refaktoryzuj.addEventListener('click', () => void podgladRefaktoryzacji());
  wersje.addEventListener('click', () => void odczytajWersje());
  przywroc.addEventListener('click', () => void przywrocWersje());
  operacja.addEventListener('click', () => void wykonajOperacje());
  warsztat.addEventListener('click', () => void sprawdzWarsztat());

  /** sciezkaPliku bierze plik ze stanu modułu — tego, który wskazał Project Tree. */
  function sciezkaPliku(): string {
    return stan.sciezka();
  }

  function brakPliku(): boolean {
    if (sciezkaPliku() !== '') return false;
    tresc.blad('Czynność dotyczy pliku — wskaż go w Project Tree albo otwórz w Code Editorze.');
    return true;
  }

  async function nawigujPoSymbolu(): Promise<void> {
    if (brakPliku()) return;
    tresc.ladowanie('Pytanie do serwera języka…');
    const wynik = await zrodlo.nawigujDoSymbolu({
      windowId: stan.okno(),
      path: sciezkaPliku(),
      line: liczbaPola(wierszSymbolu, 1),
      column: liczbaPola(kolumnaSymbolu, 1),
      kind: rodzajNawigacji.value as SymbolNavigationKind,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Nawigacja po symbolu', wynik.blad), wynik.blad);
      return;
    }
    // Dwa różne zdania o pustym wykazie — i to jest cała wartość tego pola.
    if (!wynik.wynik.serverAvailable) {
      tresc.pusto(
        'Serwera języka nie ma na serwerze aplikacji, więc pytania nie zadano. To nie znaczy, że ' +
          'symbolu nie ma — znaczy, że nie było komu go szukać.',
      );
      return;
    }
    if (wynik.wynik.symbols.length === 0) {
      tresc.pusto('Serwer języka odpowiedział: pod tym miejscem nie ma symbolu do pokazania.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.symbols.map((symbol) => {
        const element = akapitWarsztatu(
          `${symbol.name} — ${symbol.path}:${String(symbol.line)}` +
            `${symbol.column === undefined ? '' : `:${String(symbol.column)}`}`,
        );
        // Wskazanie wiedzie do Code Editora tą samą drogą, co kliknięcie
        // w Project Tree: stan modułu jest jedną prawdą o otwartym pliku.
        element.addEventListener('click', () => stan.wskazPlik(symbol.path));
        return element;
      }),
    );
  }

  async function formatujPlik(): Promise<void> {
    if (brakPliku()) return;
    tresc.ladowanie('Formatowanie…');
    const wynik = await zrodlo.formatuj({
      windowId: stan.okno(),
      path: sciezkaPliku(),
      writeToDisk: true,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Formatowanie pliku', wynik.blad), wynik.blad);
      return;
    }
    tresc.potwierdzenie(
      wynik.wynik.changed
        ? `Plik sformatowany programem ${wynik.wynik.formatter ?? 'serwera'}.`
        : 'Plik był już sformatowany — nic się nie zmieniło.',
      true,
    );
  }

  async function uruchomAnalize(): Promise<void> {
    tresc.ladowanie('Analiza statyczna repozytorium…');
    const wynik = await zrodlo.analiza({ windowId: stan.okno() });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Analiza statyczna', wynik.blad), wynik.blad);
      return;
    }
    if (!wynik.wynik.linterAvailable) {
      tresc.pusto(
        'Na serwerze nie ma programu analizy statycznej, więc przebiegu nie było. Repozytorium ' +
          'niesprawdzone to inne zdanie niż repozytorium czyste.',
      );
      return;
    }
    if (wynik.wynik.diagnostics.length === 0) {
      tresc.pusto('Analiza przeszła bez zgłoszeń.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.diagnostics.map((zgloszenie) => {
        const element = akapitWarsztatu(
          `[${zgloszenie.severity}] ${zgloszenie.path}:${String(zgloszenie.line)} — ${zgloszenie.message}` +
            `${zgloszenie.source === undefined ? '' : ` (${zgloszenie.source})`}`,
        );
        element.addEventListener('click', () => stan.wskazPlik(zgloszenie.path));
        return element;
      }),
    );
  }

  async function podgladRefaktoryzacji(): Promise<void> {
    if (brakPliku()) return;
    tresc.ladowanie('Refaktoryzacja — podgląd…');
    const zadanie: Parameters<ZrodloWarsztatu['refaktoryzuj']>[0] = {
      windowId: stan.okno(),
      path: sciezkaPliku(),
      line: liczbaPola(wierszSymbolu, 1),
      column: liczbaPola(kolumnaSymbolu, 1),
      kind: rodzajRefaktoryzacji.value as RefactorKind,
      preview: true,
    };
    if (nowaNazwa.value.trim() !== '') zadanie.newName = nowaNazwa.value.trim();
    const wynik = await zrodlo.refaktoryzuj(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Refaktoryzacja', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.edits.length === 0) {
      tresc.pusto('Serwer języka nie zaproponował żadnej zmiany dla tego miejsca.');
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapitWarsztatu(
        `Podgląd — zmiany w ${String((wynik.wynik.changedPaths ?? []).length)} plikach; ` +
          'zapisu nie było.',
      ),
      ...wynik.wynik.edits.map((zmiana) =>
        akapitWarsztatu(
          `${zmiana.path}: wiersze ${String(zmiana.startLine)}–${String(zmiana.endLine)}`,
        ),
      ),
    );
  }

  async function odczytajWersje(): Promise<void> {
    if (brakPliku()) return;
    tresc.ladowanie('Odczyt historii pliku…');
    const wynik = await zrodlo.wersjePliku({ windowId: stan.okno(), path: sciezkaPliku() });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Historia wersji pliku', wynik.blad), wynik.blad);
      return;
    }
    if (wynik.wynik.versions.length === 0) {
      tresc.pusto(
        'Plik nie ma migawek. Wersja powstaje przy zapisie z zaznaczonym „załóż wersję” — to zapis ' +
          'roboczy edytora, a nie historia repozytorium.',
      );
      return;
    }
    ostatniaWersja = wynik.wynik.versions[0]?.id ?? '';
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.versions.map((wersja) =>
        akapitWarsztatu(
          `${wersja.id} — ${new Date(wersja.createdAt).toLocaleString('pl-PL')}` +
            `${wersja.sizeBytes === undefined ? '' : ` · ${String(wersja.sizeBytes)} B`}`,
        ),
      ),
    );
  }

  async function przywrocWersje(): Promise<void> {
    if (ostatniaWersja === '') {
      tresc.blad('Najpierw odczytaj historię wersji — przywrócenie wskazuje konkretną migawkę.');
      return;
    }
    tresc.ladowanie('Przywracanie wersji…');
    const wynik = await zrodlo.przywrocWersje({
      windowId: stan.okno(),
      versionId: ostatniaWersja,
      writeToDisk: true,
    });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Przywrócenie wersji pliku', wynik.blad), wynik.blad);
      return;
    }
    // Stan modułu dostaje przywrócony plik, żeby Code Editor pokazał to, co
    // naprawdę leży teraz na dysku, a nie treść sprzed przywrócenia.
    stan.ustawPlik(wynik.wynik);
    tresc.potwierdzenie(
      `Plik ${wynik.wynik.path} przywrócony; stan sprzed przywrócenia rdzeń odłożył jako nową migawkę.`,
      true,
    );
  }

  async function wykonajOperacje(): Promise<void> {
    if (brakPliku()) return;
    tresc.ladowanie('Operacja kontekstowa…');
    const zadanie: Parameters<ZrodloWarsztatu['operacjaKontekstowa']>[0] = {
      windowId: stan.okno(),
      operation: rodzajOperacji.value as ContextualOpKind,
      path: sciezkaPliku(),
    };
    if (zaznaczenie.value.trim() !== '') zadanie.selection = zaznaczenie.value;
    if (wskazowka.value.trim() !== '') zadanie.instruction = wskazowka.value.trim();
    const wynik = await zrodlo.operacjaKontekstowa(zadanie);
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Operacja kontekstowa', wynik.blad), wynik.blad);
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      akapitWarsztatu(
        'Wynik jest propozycją do przyjęcia — okno niczego nie zapisało na dysku.',
      ),
      blokWarsztatu(wynik.wynik.result),
    );
  }

  async function sprawdzWarsztat(): Promise<void> {
    tresc.ladowanie('Sonda programów warsztatu…');
    const wynik = await zrodlo.warsztat({ programs: [] });
    if (!wynik.udany || wynik.wynik === undefined) {
      tresc.blad(opisOdmowyBledu('Sonda programów warsztatu', wynik.blad), wynik.blad);
      return;
    }
    const miejsce = tresc.tresc();
    miejsce.replaceChildren(
      ...wynik.wynik.programs.map((pozycja) =>
        akapitWarsztatu(
          `${pozycja.program} — ${pozycja.present ? 'stoi' : 'brak na serwerze'}` +
            `${pozycja.version === undefined ? '' : ` · ${pozycja.version}`}`,
        ),
      ),
    );
  }

  rama.narzedzia.append(
    rysujZaleznosci(zaleznosci(['serwer-jezyka', 'formater-linter'])),
    objasnienieZakladki(
      'Programy warsztatu stoją NA SERWERZE, nie u Operatora: instalka jest cienka, a arsenał ' +
        'ma serwer. Sonda programów mówi, czym serwer dziś dysponuje — zanim czynność go zażąda.',
    ),
  );
  rama.cialo.append(
    wiersz('Wiersz', wierszSymbolu, { klasa: 'dn-pole' }),
    wiersz('Kolumna', kolumnaSymbolu, { klasa: 'dn-pole' }),
    wiersz('Rodzaj nawigacji', rodzajNawigacji, { klasa: 'dn-pole' }),
    wiersz('Rodzaj refaktoryzacji', rodzajRefaktoryzacji, { klasa: 'dn-pole' }),
    wiersz('Nowa nazwa', nowaNazwa, { klasa: 'dn-pole' }),
    wiersz('Operacja kontekstowa', rodzajOperacji, { klasa: 'dn-pole' }),
    wiersz('Zaznaczenie', zaznaczenie, {
      klasa: 'dn-pole',
      objasnienie: 'Puste znaczy operację nad całym plikiem wskazanym w Project Tree.',
    }),
    wiersz('Wskazanie Operatora', wskazowka, { klasa: 'dn-pole' }),
    pasekZakladki(nawiguj, formatuj, analizuj, refaktoryzuj),
    pasekZakladki(wersje, przywroc, operacja, warsztat),
    tresc.element,
  );

  return {
    element: rama.element,
    // Okno nie odpytuje rdzenia samo: każda z tych czynności jest poleceniem
    // Operatora nad wskazanym plikiem, a nie stanem do odtworzenia.
    odswiez: () => undefined,
    zamknij: () => undefined,
  };
}

/** liczbaPola czyta liczbę z pola; wartość nieczytelna spada na wartość zastępczą. */
function liczbaPola(kontrolka: HTMLInputElement, zastepcza: number): number {
  const liczba = Number.parseInt(kontrolka.value, 10);
  return Number.isNaN(liczba) || liczba < 1 ? zastepcza : liczba;
}

/** akapitWarsztatu składa jeden wiersz treści okna. */
function akapitWarsztatu(zdanie: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'mdev-wiersz';
  element.textContent = zdanie;
  return element;
}

/** blokWarsztatu składa miejsce na treść wielowierszową — wynik operacji. */
function blokWarsztatu(zawartosc: string): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'mdev-blok';
  element.textContent = zawartosc;
  return element;
}
