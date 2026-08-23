import { Command, ConfigScope, StudioDocumentFormat } from '../../../../shared/contract';
import type {
  CzynnoscWarsztatu,
  OtoczenieCzynnosci,
  PoleCzynnosci,
  WartosciCzynnosci,
  ZlozenieZadania,
} from './czynnosci-warsztatu';

/**
 * Katalog czynności redakcyjnych — piętnaście komend, które prowadzą dokument
 * Studia poza sam edytor: gałęzie, wydanie, podgląd układu, wsad, wczytywanie.
 *
 * Katalog jest osobny od `czynnosci-warsztatu`, choć typy dzieli. Powód jest
 * w materiale, na którym te dwa zbiory pracują: warsztat bierze DOKUMENT PDF
 * z magazynu okna, redakcja bierze DOKUMENT STUDIA wczytany w edytorze. Jeden
 * wspólny wykaz kazałby Operatorowi wybierać czynność, która nie ma na czym
 * pracować, i dowiadywać się o tym dopiero z odmowy.
 *
 * Nazwy pól są nazwami kontraktu — one jadą do rdzenia. Etykiety są zdaniem
 * Operatora i z nazwami się nie pokrywają.
 */

/** Odczyt pola tekstowego; puste zwraca `undefined`, nie pusty napis. */
function tekst(wartosci: WartosciCzynnosci, kod: string): string | undefined {
  const wartosc = (wartosci[kod] ?? '').trim();
  return wartosc === '' ? undefined : wartosc;
}

/** Odczyt liczby; wartość spoza liczb zwraca `undefined`. */
function liczba(wartosci: WartosciCzynnosci, kod: string): number | undefined {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return undefined;
  const odczytana = Number(wartosc);
  return Number.isFinite(odczytana) ? odczytana : undefined;
}

/**
 * Odczyt pola trójstanowego: „domyślnie", „tak", „nie".
 *
 * Przełącznik dwustanowy byłby tu pomyłką. Kontrakt paczki redakcyjnej mówi
 * o polach, których BRAK znaczy „tak" — przełącznik niezaznaczony wysyłałby
 * `false` i wyłączał człon, o którego wyłączenie Operator nie prosił.
 */
function trojstanowe(wartosci: WartosciCzynnosci, kod: string): boolean | undefined {
  const wartosc = wartosci[kod] ?? '';
  if (wartosc === 'tak') return true;
  if (wartosc === 'nie') return false;
  return undefined;
}

/** Rozbiór listy oddzielonej przecinkami na wykaz bez pozycji pustych. */
function wykaz(wartosci: WartosciCzynnosci, kod: string): string[] {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return [];
  return wartosc
    .split(',')
    .map((pozycja) => pozycja.trim())
    .filter((pozycja) => pozycja !== '');
}

/** Dokłada do żądania wyłącznie pola o wartości podanej. */
function zPolami(
  podstawa: Record<string, unknown>,
  dodatki: Record<string, unknown>,
): Record<string, unknown> {
  const zadanie = { ...podstawa };
  for (const [nazwa, wartosc] of Object.entries(dodatki)) {
    if (wartosc === undefined) continue;
    if (Array.isArray(wartosc) && wartosc.length === 0) continue;
    zadanie[nazwa] = wartosc;
  }
  return zadanie;
}

/**
 * Wymaga dokumentu wczytanego w edytorze.
 *
 * Odmowa pada tutaj, a nie w rdzeniu, bo tutaj da się powiedzieć Operatorowi,
 * CO ma zrobić: wczytać dokument. Rdzeń odmówiłby brakiem pola `documentId`,
 * z czego Operator nie wyczyta ani przyczyny, ani drogi wyjścia.
 */
function zDokumentem(
  otoczenie: OtoczenieCzynnosci,
  zloz: (idDokumentu: string) => Record<string, unknown>,
): ZlozenieZadania {
  if (otoczenie.idDokumentu === null) {
    return { odmowa: 'Czynność pracuje na dokumencie Studia — wczytaj dokument w Studio Editorze.' };
  }
  return { zadanie: zloz(otoczenie.idDokumentu) };
}

/** Pozycje pola trójstanowego — jedne dla wszystkich członów paczki. */
const POZYCJE_TROJSTANOWE = [
  { wartosc: '', etykieta: 'domyślnie (dołącz)' },
  { wartosc: 'tak', etykieta: 'dołącz' },
  { wartosc: 'nie', etykieta: 'pomiń' },
] as const;

/** Pozycje formatu dokumentu — słownik zamiany formatu z kontraktu. */
const POZYCJE_FORMATU = [
  { wartosc: StudioDocumentFormat.Pdf, etykieta: 'PDF' },
  { wartosc: StudioDocumentFormat.Docx, etykieta: 'DOCX' },
  { wartosc: StudioDocumentFormat.Markdown, etykieta: 'Markdown' },
  { wartosc: StudioDocumentFormat.Txt, etykieta: 'Tekst' },
] as const;

/** Pole wersji odniesienia — powtarza się w trzech czynnościach różnicy. */
const POLE_WERSJI_ODNIESIENIA: PoleCzynnosci = {
  kod: 'baseVersionId',
  etykieta: 'Wersja odniesienia',
  rodzaj: 'tekst',
  opis: 'Identyfikator wersji z Session Repository.',
};

export const CZYNNOSCI_REDAKCJI: readonly CzynnoscWarsztatu[] = [
  {
    komenda: Command.StudioBranchCreate,
    nazwa: 'Rozgałęź dokument',
    opis:
      'Zakłada wariant dokumentu od wskazanej wersji i otwiera go jako treść bieżącą. ' +
      'Historia pnia zostaje nietknięta.',
    pola: [
      { kod: 'fromVersionId', etykieta: 'Wersja startowa', rodzaj: 'tekst' },
      { kod: 'name', etykieta: 'Nazwa gałęzi', rodzaj: 'tekst', podpowiedz: 'redakcja klienta' },
    ],
    zloz(wartosci, otoczenie) {
      const wersja = tekst(wartosci, 'fromVersionId');
      const nazwa = tekst(wartosci, 'name');
      if (wersja === undefined || nazwa === undefined) {
        return { odmowa: 'Gałąź potrzebuje wersji startowej i nazwy.' };
      }
      return zDokumentem(otoczenie, (idDokumentu) => ({
        documentId: idDokumentu,
        fromVersionId: wersja,
        name: nazwa,
      }));
    },
  },
  {
    komenda: Command.StudioBranchList,
    nazwa: 'Pokaż gałęzie',
    opis: 'Wykaz gałęzi dokumentu wraz z ich punktami startowymi.',
    pola: [],
    zloz(_wartosci, otoczenie) {
      return zDokumentem(otoczenie, (idDokumentu) => ({ documentId: idDokumentu }));
    },
  },
  {
    komenda: Command.StudioBranchMerge,
    nazwa: 'Scal gałęzie',
    opis:
      'Scala gałąź z gałęzią docelową. Fragment zmieniony po obu stronach wraca ' +
      'konfliktem — rdzeń nie wybiera za Operatora.',
    pola: [
      { kod: 'sourceBranchId', etykieta: 'Gałąź scalana', rodzaj: 'tekst' },
      { kod: 'targetBranchId', etykieta: 'Gałąź docelowa', rodzaj: 'tekst' },
      {
        kod: 'resolutionIndex',
        etykieta: 'Numer konfliktu do rozstrzygnięcia',
        rodzaj: 'liczba',
        opis: 'Numer z wykazu konfliktów zwróconego przez poprzednie scalenie.',
      },
      {
        kod: 'resolutionSide',
        etykieta: 'Strona, która zostaje',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'bez rozstrzygnięcia' },
          { wartosc: 'scalana', etykieta: 'gałąź scalana' },
          { wartosc: 'docelowa', etykieta: 'gałąź docelowa' },
          { wartosc: 'obie', etykieta: 'obie, jedna pod drugą' },
        ],
      },
    ],
    zloz(wartosci) {
      const scalana = tekst(wartosci, 'sourceBranchId');
      const docelowa = tekst(wartosci, 'targetBranchId');
      if (scalana === undefined || docelowa === undefined) {
        return { odmowa: 'Scalenie potrzebuje gałęzi scalanej i docelowej.' };
      }
      const numer = liczba(wartosci, 'resolutionIndex');
      const strona = tekst(wartosci, 'resolutionSide');
      const rozstrzygniecia =
        numer === undefined || strona === undefined
          ? []
          : [{ index: numer, side: strona }];
      return {
        zadanie: zPolami(
          { sourceBranchId: scalana, targetBranchId: docelowa },
          { resolutions: rozstrzygniecia },
        ),
      };
    },
  },
  {
    komenda: Command.StudioVersionReferenceCreate,
    nazwa: 'Odwołanie do wersji',
    opis: 'Tworzy adres wersji w obrębie platformy o wskazanym zasięgu widoczności.',
    pola: [
      { kod: 'versionId', etykieta: 'Wersja', rodzaj: 'tekst' },
      {
        kod: 'scope',
        etykieta: 'Zasięg widoczności',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'domyślnie (sesja)' },
          { wartosc: ConfigScope.Session, etykieta: 'sesja' },
          { wartosc: ConfigScope.Project, etykieta: 'projekt' },
          { wartosc: ConfigScope.Global, etykieta: 'globalny' },
        ],
      },
      {
        kod: 'expiresAt',
        etykieta: 'Wygasa (ms epoki)',
        rodzaj: 'liczba',
        opis: 'Puste znaczy bez wygasania.',
      },
    ],
    zloz(wartosci) {
      const wersja = tekst(wartosci, 'versionId');
      if (wersja === undefined) return { odmowa: 'Odwołanie potrzebuje wskazania wersji.' };
      return {
        zadanie: zPolami(
          { versionId: wersja },
          { scope: tekst(wartosci, 'scope'), expiresAt: liczba(wartosci, 'expiresAt') },
        ),
      };
    },
  },
  {
    komenda: Command.StudioRepositoryExport,
    nazwa: 'Wydaj historię',
    opis: 'Pakuje wskazane wersje albo całą historię w archiwum z manifestem.',
    pola: [
      {
        kod: 'versionIds',
        etykieta: 'Wersje',
        rodzaj: 'tekst',
        opis: 'Identyfikatory oddzielone przecinkami. Puste bierze całą historię.',
      },
      {
        kod: 'format',
        etykieta: 'Postać archiwum',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: '', etykieta: 'domyślnie (zip)' },
          { wartosc: 'zip', etykieta: 'ZIP' },
          { wartosc: 'tar', etykieta: 'TAR' },
        ],
      },
    ],
    zloz(wartosci, otoczenie) {
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          { documentId: idDokumentu, windowId: otoczenie.idOkna },
          { versionIds: wykaz(wartosci, 'versionIds'), format: tekst(wartosci, 'format') },
        ),
      );
    },
  },
  {
    komenda: Command.StudioPackageExport,
    nazwa: 'Paczka redakcyjna',
    opis: 'Pakuje dokument finalny, historię, raport zmian i adnotacje w jedno archiwum przekazania.',
    pola: [
      { kod: 'versionId', etykieta: 'Wersja finalna', rodzaj: 'tekst' },
      {
        kod: 'documentFormat',
        etykieta: 'Format dokumentu finalnego',
        rodzaj: 'wybor',
        pozycje: [{ wartosc: '', etykieta: 'format dokumentu' }, ...POZYCJE_FORMATU],
      },
      { kod: 'includeHistory', etykieta: 'Historia wersji', rodzaj: 'wybor', pozycje: [...POZYCJE_TROJSTANOWE] },
      { kod: 'includeDiffReport', etykieta: 'Raport zmian', rodzaj: 'wybor', pozycje: [...POZYCJE_TROJSTANOWE] },
      { kod: 'includeAnnotations', etykieta: 'Adnotacje', rodzaj: 'wybor', pozycje: [...POZYCJE_TROJSTANOWE] },
    ],
    zloz(wartosci, otoczenie) {
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          { documentId: idDokumentu, windowId: otoczenie.idOkna },
          {
            versionId: tekst(wartosci, 'versionId'),
            documentFormat: tekst(wartosci, 'documentFormat'),
            includeHistory: trojstanowe(wartosci, 'includeHistory'),
            includeDiffReport: trojstanowe(wartosci, 'includeDiffReport'),
            includeAnnotations: trojstanowe(wartosci, 'includeAnnotations'),
          },
        ),
      );
    },
  },
  {
    komenda: Command.StudioDiffReportExport,
    nazwa: 'Wydaj raport zmian',
    opis: 'Składa widok różnicowy jako osobny dokument redakcji wraz ze statystyką i adnotacjami.',
    pola: [
      POLE_WERSJI_ODNIESIENIA,
      {
        kod: 'targetVersionId',
        etykieta: 'Wersja porównywana',
        rodzaj: 'tekst',
        opis: 'Puste bierze treść bieżącą.',
      },
      { kod: 'format', etykieta: 'Format raportu', rodzaj: 'wybor', pozycje: [...POZYCJE_FORMATU] },
      { kod: 'includeAnnotations', etykieta: 'Adnotacje', rodzaj: 'wybor', pozycje: [...POZYCJE_TROJSTANOWE] },
    ],
    zloz(wartosci, otoczenie) {
      const baza = tekst(wartosci, 'baseVersionId');
      const format = tekst(wartosci, 'format');
      if (baza === undefined || format === undefined) {
        return { odmowa: 'Raport potrzebuje wersji odniesienia i formatu.' };
      }
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          { documentId: idDokumentu, baseVersionId: baza, format, windowId: otoczenie.idOkna },
          {
            targetVersionId: tekst(wartosci, 'targetVersionId'),
            includeAnnotations: trojstanowe(wartosci, 'includeAnnotations'),
          },
        ),
      );
    },
  },
  {
    komenda: Command.StudioDiffVisual,
    nazwa: 'Różnica wizualna',
    opis:
      'Porównuje WYGLĄD dwóch wersji: renderuje obie i wskazuje obszary różniące się ' +
      'graficznie. Widzi zmianę układu, której różnica tekstowa nie widzi.',
    pola: [
      POLE_WERSJI_ODNIESIENIA,
      { kod: 'targetVersionId', etykieta: 'Wersja porównywana', rodzaj: 'tekst' },
      { kod: 'pageFrom', etykieta: 'Pierwsza strona', rodzaj: 'liczba' },
      { kod: 'pageTo', etykieta: 'Ostatnia strona', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const baza = tekst(wartosci, 'baseVersionId');
      const cel = tekst(wartosci, 'targetVersionId');
      if (baza === undefined || cel === undefined) {
        return { odmowa: 'Porównanie wizualne potrzebuje obu wersji.' };
      }
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          {
            documentId: idDokumentu,
            baseVersionId: baza,
            targetVersionId: cel,
            windowId: otoczenie.idOkna,
          },
          { pageFrom: liczba(wartosci, 'pageFrom'), pageTo: liczba(wartosci, 'pageTo') },
        ),
      );
    },
  },
  {
    komenda: Command.StudioDiffSource,
    nazwa: 'Porównaj ze źródłem',
    opis:
      'Zestawia dokument roboczy z materiałem wejściowym z Library albo z magazynu ' +
      'i wskazuje rozbieżności. Brak wskazania bierze plik, z którego dokument otwarto.',
    pola: [
      { kod: 'sourceLibraryFileId', etykieta: 'Plik repozytorium', rodzaj: 'tekst' },
      { kod: 'sourceAssetId', etykieta: 'Zasób magazynu', rodzaj: 'zasob' },
      { kod: 'versionId', etykieta: 'Wersja dokumentu', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          { documentId: idDokumentu },
          {
            sourceLibraryFileId: tekst(wartosci, 'sourceLibraryFileId'),
            sourceAssetId: tekst(wartosci, 'sourceAssetId'),
            versionId: tekst(wartosci, 'versionId'),
          },
        ),
      );
    },
  },
  {
    komenda: Command.StudioPreviewRender,
    nazwa: 'Wyrenderuj podgląd',
    opis:
      'Składa dokument w formacie docelowym i oddaje strony jako zasoby — podgląd ' +
      'pokazuje UKŁAD, a nie sam tekst.',
    pola: [
      { kod: 'versionId', etykieta: 'Wersja', rodzaj: 'tekst', opis: 'Puste bierze treść bieżącą.' },
      { kod: 'format', etykieta: 'Format docelowy', rodzaj: 'wybor', pozycje: [...POZYCJE_FORMATU] },
      { kod: 'pageFrom', etykieta: 'Pierwsza strona', rodzaj: 'liczba' },
      { kod: 'pageTo', etykieta: 'Ostatnia strona', rodzaj: 'liczba' },
      { kod: 'profileId', etykieta: 'Profil wydania', rodzaj: 'tekst' },
      {
        kod: 'watermark',
        etykieta: 'Znak wodny',
        rodzaj: 'tekst',
        opis: 'Dotyczy wyłącznie podglądu — treści dokumentu nie zmienia.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const format = tekst(wartosci, 'format');
      if (format === undefined) return { odmowa: 'Render potrzebuje formatu docelowego.' };
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          { documentId: idDokumentu, format, windowId: otoczenie.idOkna },
          {
            versionId: tekst(wartosci, 'versionId'),
            pageFrom: liczba(wartosci, 'pageFrom'),
            pageTo: liczba(wartosci, 'pageTo'),
            profileId: tekst(wartosci, 'profileId'),
            watermark: tekst(wartosci, 'watermark'),
          },
        ),
      );
    },
  },
  {
    komenda: Command.StudioSearchSemantic,
    nazwa: 'Szukaj znaczeniowo',
    opis:
      'Znajduje fragmenty bliskie zapytaniu, nie tylko dopasowania dosłowne. ' +
      'Odpowiedź mówi, którą drogą poszła: kanałem modelu czy miarą rdzenia.',
    pola: [
      { kod: 'query', etykieta: 'Zapytanie', rodzaj: 'wielowiersz' },
      { kod: 'versionId', etykieta: 'Wersja', rodzaj: 'tekst' },
      { kod: 'limit', etykieta: 'Ile fragmentów', rodzaj: 'liczba' },
      { kod: 'minScore', etykieta: 'Próg bliskości (0–1)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      const zapytanie = tekst(wartosci, 'query');
      if (zapytanie === undefined) return { odmowa: 'Wyszukiwanie potrzebuje zapytania.' };
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          { documentId: idDokumentu, query: zapytanie },
          {
            versionId: tekst(wartosci, 'versionId'),
            limit: liczba(wartosci, 'limit'),
            minScore: liczba(wartosci, 'minScore'),
          },
        ),
      );
    },
  },
  {
    komenda: Command.StudioAssetEmbed,
    nazwa: 'Osadź zasób',
    opis: 'Wstawia grafikę z magazynu w treść dokumentu — odwołaniem, nie kopią bajtów.',
    pola: [
      { kod: 'assetId', etykieta: 'Zasób', rodzaj: 'zasob' },
      {
        kod: 'position',
        etykieta: 'Położenie (w znakach)',
        rodzaj: 'liczba',
        opis: 'Puste wstawia na końcu treści.',
      },
      { kod: 'caption', etykieta: 'Podpis', rodzaj: 'tekst' },
      { kod: 'altText', etykieta: 'Tekst alternatywny', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const zasob = tekst(wartosci, 'assetId');
      if (zasob === undefined) return { odmowa: 'Osadzenie potrzebuje wskazania zasobu.' };
      return zDokumentem(otoczenie, (idDokumentu) =>
        zPolami(
          { documentId: idDokumentu, assetId: zasob },
          {
            position: liczba(wartosci, 'position'),
            caption: tekst(wartosci, 'caption'),
            altText: tekst(wartosci, 'altText'),
          },
        ),
      );
    },
  },
  {
    komenda: Command.StudioBatchRun,
    nazwa: 'Wsad na dokumentach',
    opis:
      'Wykonuje tę samą operację kontekstową na wielu dokumentach. Odmowa jednego ' +
      'nie wstrzymuje pozostałych — wynik niesie wykaz odrzuconych wraz z powodem.',
    pola: [
      {
        kod: 'documentIds',
        etykieta: 'Dokumenty',
        rodzaj: 'tekst',
        opis: 'Identyfikatory dokumentów Studia oddzielone przecinkami.',
      },
      { kod: 'actionId', etykieta: 'Operacja', rodzaj: 'tekst', podpowiedz: 'skroc' },
    ],
    zloz(wartosci, otoczenie) {
      const dokumenty = wykaz(wartosci, 'documentIds');
      const operacja = tekst(wartosci, 'actionId');
      if (dokumenty.length === 0 || operacja === undefined) {
        return { odmowa: 'Wsad potrzebuje wykazu dokumentów i operacji.' };
      }
      return {
        zadanie: { windowId: otoczenie.idOkna, documentIds: dokumenty, actionId: operacja },
      };
    },
  },
  {
    komenda: Command.StudioIngestUrl,
    nazwa: 'Wczytaj stronę',
    opis:
      'Pobiera stronę sieciową, oczyszcza ją i dokłada jako pozycję kolejki wczytywania. ' +
      'Strona zbudowana wyłącznie skryptem odda mało treści — migawkę wykonanej strony ' +
      'oddaje moduł Browser.',
    pola: [
      { kod: 'url', etykieta: 'Adres strony', rodzaj: 'tekst', podpowiedz: 'https://…' },
      { kod: 'readability', etykieta: 'Oczyść do postaci czytelnej', rodzaj: 'wybor', pozycje: [...POZYCJE_TROJSTANOWE] },
      { kod: 'includeImages', etykieta: 'Wciągnij obrazy do magazynu', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      const adres = tekst(wartosci, 'url');
      if (adres === undefined) return { odmowa: 'Wczytanie potrzebuje adresu strony.' };
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna, url: adres },
          {
            readability: trojstanowe(wartosci, 'readability'),
            includeImages: wartosci['includeImages'] === 'tak' ? true : undefined,
          },
        ),
      };
    },
  },
  {
    komenda: Command.StudioIngestDeviceScan,
    nazwa: 'Skanuj z urządzenia',
    opis:
      'Skaner stoi po stronie Operatora, a czynność wykonuje się na serwerze rdzenia — ' +
      'to dwie różne maszyny. Rdzeń odmówi, nazywając brak i wskazując drogę, ' +
      'która działa: dołożenie zeskanowanego pliku do kolejki wczytywania.',
    pola: [
      { kod: 'deviceId', etykieta: 'Urządzenie', rodzaj: 'tekst' },
      { kod: 'resolutionDpi', etykieta: 'Rozdzielczość (dpi)', rodzaj: 'liczba' },
      { kod: 'colorMode', etykieta: 'Tryb barwny', rodzaj: 'tekst' },
      { kod: 'pages', etykieta: 'Liczba stron z podajnika', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { windowId: otoczenie.idOkna },
          {
            deviceId: tekst(wartosci, 'deviceId'),
            resolutionDpi: liczba(wartosci, 'resolutionDpi'),
            colorMode: tekst(wartosci, 'colorMode'),
            pages: liczba(wartosci, 'pages'),
          },
        ),
      };
    },
  },
];

/**
 * Zdanie o skutku czynności redakcyjnej złożone z odpowiedzi rdzenia.
 *
 * Każde pole jest sprawdzane osobno, bo każda z piętnastu odpowiedzi ma inny
 * kształt. Meldunek „gotowe" bez liczby nie odróżniałby scalenia, które doszło
 * do skutku, od scalenia zatrzymanego konfliktem.
 */
export function opiszSkutekRedakcji(odpowiedz: unknown): string {
  if (typeof odpowiedz !== 'object' || odpowiedz === null) return 'Rdzeń przyjął czynność.';
  const tresc = odpowiedz as Record<string, unknown>;
  const czesci: string[] = [];

  const zasob = tresc['asset'];
  if (typeof zasob === 'object' && zasob !== null) {
    const identyfikator = (zasob as Record<string, unknown>)['id'];
    if (typeof identyfikator === 'string') czesci.push(`nowy zasób ${identyfikator}`);
  }
  const galaz = tresc['branch'];
  if (typeof galaz === 'object' && galaz !== null) {
    const nazwa = (galaz as Record<string, unknown>)['name'];
    if (typeof nazwa === 'string') czesci.push(`gałąź „${nazwa}"`);
  }
  if (Array.isArray(tresc['branches'])) czesci.push(`gałęzi: ${tresc['branches'].length}`);
  if (typeof tresc['merged'] === 'boolean') {
    czesci.push(tresc['merged'] === true ? 'scalone' : 'scalenie wstrzymane konfliktem');
  }
  if (Array.isArray(tresc['conflicts']) && tresc['conflicts'].length > 0) {
    czesci.push(`konfliktów do rozstrzygnięcia: ${tresc['conflicts'].length}`);
  }
  if (typeof tresc['reference'] === 'string') czesci.push(`odwołanie ${tresc['reference']}`);
  if (typeof tresc['entries'] === 'number') czesci.push(`pozycji: ${tresc['entries']}`);
  if (typeof tresc['sizeBytes'] === 'number') czesci.push(`bajtów: ${tresc['sizeBytes']}`);
  if (typeof tresc['pages'] === 'number') czesci.push(`stron: ${tresc['pages']}`);
  if (Array.isArray(tresc['pageAssetIds'])) czesci.push(`zasobów stron: ${tresc['pageAssetIds'].length}`);
  if (Array.isArray(tresc['regions'])) czesci.push(`obszarów różnicy: ${tresc['regions'].length}`);
  if (Array.isArray(tresc['overlayAssetIds'])) czesci.push(`nakładek: ${tresc['overlayAssetIds'].length}`);
  if (Array.isArray(tresc['matches'])) czesci.push(`fragmentów: ${tresc['matches'].length}`);
  if (typeof tresc['mode'] === 'string') czesci.push(`droga: ${tresc['mode']}`);
  if (Array.isArray(tresc['hunks'])) czesci.push(`rozbieżności: ${tresc['hunks'].length}`);
  if (typeof tresc['sourceResolved'] === 'boolean' && tresc['sourceResolved'] === false) {
    czesci.push('materiału wejściowego nie udało się odczytać');
  }
  if (typeof tresc['runId'] === 'string') czesci.push(`przebieg ${tresc['runId']}`);
  if (typeof tresc['accepted'] === 'number') czesci.push(`przyjętych: ${tresc['accepted']}`);
  if (Array.isArray(tresc['rejected']) && tresc['rejected'].length > 0) {
    czesci.push(`odrzuconych: ${tresc['rejected'].length}`);
  }
  const pozycja = tresc['item'];
  if (typeof pozycja === 'object' && pozycja !== null) {
    const identyfikator = (pozycja as Record<string, unknown>)['id'];
    if (typeof identyfikator === 'string') czesci.push(`pozycja kolejki ${identyfikator}`);
  }

  return czesci.length === 0 ? 'Rdzeń przyjął czynność.' : czesci.join(', ') + '.';
}
