import { Command, StudioPdfPageOperationKind } from '../../../../shared/contract';

/** Katalog czynności warsztatu dokumentu: piętnaście komend rdzenia opisanych danymi, nie formularzami. */

/** Rodzaj kontrolki, którą pole czynności zbiera od Operatora: tekst, hasło, liczba, przełącznik, wielowiersz, zasób, zasoby albo wybór z wykazu. */
export type RodzajPola =
  | 'tekst'
  | 'tajne'
  | 'liczba'
  | 'logiczne'
  | 'wielowiersz'
  | 'zasob'
  | 'zasoby'
  | 'wybor';

/** Jedno pole czynności warsztatu: kod pola kontraktu, etykieta widoczna Operatorowi, rodzaj kontrolki oraz opis i podpowiedź pomocnicze. */
export interface PoleCzynnosci {
  kod: string;
  etykieta: string;
  rodzaj: RodzajPola;
  opis?: string;
  podpowiedz?: string;
  pozycje?: readonly { wartosc: string; etykieta: string }[];
}

/** Wartości zebrane z formularza czynności, kod pola na wartość tekstową; puste pole nie trafia do żądania złożonego do rdzenia. */
export type WartosciCzynnosci = Readonly<Record<string, string>>;

/** Otoczenie żądania czynności: identyfikator okna modułu oraz identyfikator dokumentu wczytanego w edytorze, gdy czynność go wymaga. */
export interface OtoczenieCzynnosci {
  idOkna: string;
  idDokumentu: string | null;
}

/** Wynik składania żądania czynności: gotowe żądanie do wysłania rdzeniowi albo powód, dla którego nie dało się go złożyć. */
export type ZlozenieZadania =
  | { zadanie: Record<string, unknown> }
  | { odmowa: string };

/** Czynność warsztatu dokumentu: komenda rdzenia, wykaz pól formularza oraz funkcja składająca żądanie z wartości i otoczenia. */
export interface CzynnoscWarsztatu {
  komenda: Command;
  nazwa: string;
  opis: string;
  pola: readonly PoleCzynnosci[];
  zloz(wartosci: WartosciCzynnosci, otoczenie: OtoczenieCzynnosci): ZlozenieZadania;
}

/** Odczyt pola tekstowego formularza czynności; pole puste albo złożone z samych odstępów zwraca wartość pustą, nie pusty napis. */
function tekst(wartosci: WartosciCzynnosci, kod: string): string | undefined {
  const wartosc = (wartosci[kod] ?? '').trim();
  return wartosc === '' ? undefined : wartosc;
}

/** Odczyt pola liczbowego formularza czynności; wartość spoza liczb albo niedokończona zwraca wartość pustą zamiast liczby zgadniętej. */
function liczba(wartosci: WartosciCzynnosci, kod: string): number | undefined {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return undefined;
  const odczytana = Number(wartosc);
  return Number.isFinite(odczytana) ? odczytana : undefined;
}

/** Odczyt pola przełącznika formularza czynności: wartość zgodna ze znacznikiem zaznaczenia oznacza przełącznik włączony. */
function logiczne(wartosci: WartosciCzynnosci, kod: string): boolean {
  return wartosci[kod] === 'tak';
}

/** Rozbiór listy oddzielonej przecinkami na wykaz pozycji przyciętych z odstępów, z pominięciem pozycji pustych. */
function wykaz(wartosci: WartosciCzynnosci, kod: string): string[] {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return [];
  return wartosc
    .split(',')
    .map((pozycja) => pozycja.trim())
    .filter((pozycja) => pozycja !== '');
}

/** Wiersze pola wielowierszowego formularza czynności, przycięte z odstępów, z pominięciem wierszy pustych. */
function wiersze(wartosci: WartosciCzynnosci, kod: string): string[] {
  const wartosc = tekst(wartosci, kod);
  if (wartosc === undefined) return [];
  return wartosc
    .split('\n')
    .map((wiersz) => wiersz.trim())
    .filter((wiersz) => wiersz !== '');
}

/** Dokłada do żądania wyłącznie pola o wartości podanej — pole nieustawione albo wykaz pusty nie trafia do żądania rdzenia. */
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

/** Pole materiału powtarzające się w czternastu czynnościach z piętnastu — wspólny opis zasobu magazynu wejściowego. */
const POLE_MATERIALU: PoleCzynnosci = {
  kod: 'assetId',
  etykieta: 'Materiał',
  rodzaj: 'zasob',
  opis: 'Dokument PDF z magazynu okna. Materiał zostaje nietknięty — wynik jest nowym zasobem.',
};

export const CZYNNOSCI_WARSZTATU: readonly CzynnoscWarsztatu[] = [
  {
    komenda: Command.StudioPdfMerge,
    nazwa: 'Scal dokumenty',
    opis: 'Łączy wskazane dokumenty w jeden, w kolejności zaznaczenia.',
    pola: [
      {
        kod: 'assetIds',
        etykieta: 'Materiały',
        rodzaj: 'zasoby',
        opis: 'Co najmniej dwa dokumenty. Kolejność scalania jest kolejnością na liście.',
      },
      { kod: 'name', etykieta: 'Nazwa wyniku', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const materialy = wykaz(wartosci, 'assetIds');
      if (materialy.length < 2) {
        return { odmowa: 'Scalanie wymaga co najmniej dwóch dokumentów.' };
      }
      return {
        zadanie: zPolami(
          { assetIds: materialy, windowId: otoczenie.idOkna },
          { name: tekst(wartosci, 'name') },
        ),
      };
    },
  },
  {
    komenda: Command.StudioPdfSplit,
    nazwa: 'Podziel dokument',
    opis: 'Dzieli dokument na części. Brak zakresów dzieli każdą stronę osobno.',
    pola: [
      POLE_MATERIALU,
      {
        kod: 'ranges',
        etykieta: 'Zakresy stron',
        rodzaj: 'tekst',
        podpowiedz: '1-2, 3-4',
        opis: 'Zakresy oddzielone przecinkami. Puste pole dzieli dokument strona po stronie.',
      },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          { ranges: wykaz(wartosci, 'ranges') },
        ),
      };
    },
  },
  {
    komenda: Command.StudioPdfOptimize,
    nazwa: 'Odchudź dokument',
    opis:
      'Porządkuje strukturę i usuwa powielone zasoby. Obrazów nie przelicza — ' +
      'to wymagałoby silnika rasteryzacji spoza instalki. Odpowiedź niesie zysk zmierzony.',
    pola: [
      POLE_MATERIALU,
      { kod: 'imageQuality', etykieta: 'Jakość obrazów (1–100)', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          { imageQuality: liczba(wartosci, 'imageQuality') },
        ),
      };
    },
  },
  {
    komenda: Command.StudioPdfPagesReorder,
    nazwa: 'Ułóż strony',
    opis:
      'Jedna czynność na stronach dokumentu. Wstawienie stron z innego dokumentu ' +
      'idzie przez scalanie — biblioteka wstawia strony puste.',
    pola: [
      POLE_MATERIALU,
      {
        kod: 'kind',
        etykieta: 'Czynność',
        rodzaj: 'wybor',
        pozycje: [
          { wartosc: StudioPdfPageOperationKind.Wyodrebnienie, etykieta: 'Wyodrębnij strony' },
          { wartosc: StudioPdfPageOperationKind.Usuniecie, etykieta: 'Usuń strony' },
          { wartosc: StudioPdfPageOperationKind.Przeniesienie, etykieta: 'Zbierz w kolejności' },
          { wartosc: StudioPdfPageOperationKind.Obrot, etykieta: 'Obróć' },
        ],
      },
      { kod: 'pages', etykieta: 'Zakres stron', rodzaj: 'tekst', podpowiedz: '2-3' },
      {
        kod: 'degrees',
        etykieta: 'Kąt obrotu',
        rodzaj: 'liczba',
        opis: 'Wielokrotność 90 — dokument obraca się ćwiartkami.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const strony = tekst(wartosci, 'pages');
      if (strony === undefined) {
        return { odmowa: 'Czynność na stronach potrzebuje zakresu stron.' };
      }
      const czynnosc = zPolami(
        { kind: wartosci['kind'] ?? '', pages: strony },
        { degrees: liczba(wartosci, 'degrees') },
      );
      return {
        zadanie: {
          assetId: wartosci['assetId'] ?? '',
          windowId: otoczenie.idOkna,
          operations: [czynnosc],
        },
      };
    },
  },
  {
    komenda: Command.StudioPdfStamp,
    nazwa: 'Ostempluj',
    opis: 'Kładzie pieczęć nad treścią stron — tekstową albo graficzną, nie obie naraz.',
    pola: [
      POLE_MATERIALU,
      { kod: 'text', etykieta: 'Treść pieczęci', rodzaj: 'tekst' },
      { kod: 'stampAssetId', etykieta: 'Pieczęć graficzna', rodzaj: 'zasob' },
      { kod: 'pages', etykieta: 'Zakres stron', rodzaj: 'tekst', podpowiedz: 'puste — wszystkie' },
      { kod: 'opacity', etykieta: 'Krycie (%)', rodzaj: 'liczba' },
      { kod: 'rotation', etykieta: 'Obrót pieczęci', rodzaj: 'liczba' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            text: tekst(wartosci, 'text'),
            stampAssetId: tekst(wartosci, 'stampAssetId'),
            pages: tekst(wartosci, 'pages'),
            opacity: liczba(wartosci, 'opacity'),
            rotation: liczba(wartosci, 'rotation'),
          },
        ),
      };
    },
  },
  {
    komenda: Command.StudioPdfBates,
    nazwa: 'Numeracja Bates',
    opis: 'Nadaje numerację prawną wraz z nagłówkiem i stopką, strona po stronie.',
    pola: [
      POLE_MATERIALU,
      { kod: 'prefix', etykieta: 'Przedrostek', rodzaj: 'tekst' },
      { kod: 'startNumber', etykieta: 'Numer początkowy', rodzaj: 'liczba' },
      { kod: 'digits', etykieta: 'Liczba cyfr', rodzaj: 'liczba' },
      { kod: 'header', etykieta: 'Nagłówek', rodzaj: 'tekst' },
      { kod: 'footer', etykieta: 'Stopka', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            prefix: tekst(wartosci, 'prefix'),
            startNumber: liczba(wartosci, 'startNumber'),
            digits: liczba(wartosci, 'digits'),
            header: tekst(wartosci, 'header'),
            footer: tekst(wartosci, 'footer'),
          },
        ),
      };
    },
  },
  {
    komenda: Command.StudioPdfBookmarksSet,
    nazwa: 'Ustaw zakładki',
    opis: 'Zakłada drzewo nawigacji. Podane zakładki zastępują zastane w całości.',
    pola: [
      POLE_MATERIALU,
      {
        kod: 'bookmarks',
        etykieta: 'Zakładki',
        rodzaj: 'wielowiersz',
        podpowiedz: 'Wstęp | 1\nRozdział pierwszy | 4',
        opis: 'Wiersz na zakładkę: tytuł, pionowa kreska, numer strony.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const zakladki = [];
      for (const wiersz of wiersze(wartosci, 'bookmarks')) {
        const kreska = wiersz.lastIndexOf('|');
        if (kreska < 0) {
          return { odmowa: `Wiersz „${wiersz}" nie ma kreski oddzielającej numer strony.` };
        }
        const tytul = wiersz.slice(0, kreska).trim();
        const strona = Number(wiersz.slice(kreska + 1).trim());
        if (tytul === '' || !Number.isInteger(strona) || strona < 1) {
          return { odmowa: `Wiersz „${wiersz}" nie niesie tytułu i numeru strony.` };
        }
        zakladki.push({ title: tytul, page: strona });
      }
      if (zakladki.length === 0) {
        return { odmowa: 'Nie podano ani jednej zakładki.' };
      }
      return {
        zadanie: {
          assetId: wartosci['assetId'] ?? '',
          windowId: otoczenie.idOkna,
          bookmarks: zakladki,
        },
      };
    },
  },
  {
    komenda: Command.StudioPdfFormFill,
    nazwa: 'Formularz',
    opis: 'Bez wartości sam odczytuje pola. Z wartościami wypełnia i oddaje nowy dokument.',
    pola: [
      POLE_MATERIALU,
      {
        kod: 'values',
        etykieta: 'Wartości pól',
        rodzaj: 'wielowiersz',
        podpowiedz: 'nazwisko | Kowalski',
        opis: 'Wiersz na pole: nazwa, pionowa kreska, wartość. Puste pole znaczy sam odczyt.',
      },
      { kod: 'flatten', etykieta: 'Utrwal wypełnienie', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      const pola: Record<string, string> = {};
      for (const wiersz of wiersze(wartosci, 'values')) {
        const kreska = wiersz.indexOf('|');
        if (kreska < 0) {
          return { odmowa: `Wiersz „${wiersz}" nie ma kreski oddzielającej wartość.` };
        }
        pola[wiersz.slice(0, kreska).trim()] = wiersz.slice(kreska + 1).trim();
      }
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            values: Object.keys(pola).length === 0 ? undefined : pola,
            flatten: logiczne(wartosci, 'flatten') ? true : undefined,
          },
        ),
      };
    },
  },
  {
    komenda: Command.StudioPdfExtract,
    nazwa: 'Wyciągnij osadzone',
    opis: 'Wyjmuje obrazy i załączniki, zakładając każdemu osobny zasób magazynu.',
    pola: [
      POLE_MATERIALU,
      { kod: 'images', etykieta: 'Obrazy', rodzaj: 'logiczne' },
      { kod: 'attachments', etykieta: 'Załączniki', rodzaj: 'logiczne' },
    ],
    zloz(wartosci, otoczenie) {
      const obrazy = logiczne(wartosci, 'images');
      const zalaczniki = logiczne(wartosci, 'attachments');
      if (!obrazy && !zalaczniki) {
        return { odmowa: 'Zaznacz przynajmniej jedno: obrazy albo załączniki.' };
      }
      return {
        zadanie: {
          assetId: wartosci['assetId'] ?? '',
          windowId: otoczenie.idOkna,
          images: obrazy,
          attachments: zalaczniki,
        },
      };
    },
  },
  {
    komenda: Command.StudioSecurityEncrypt,
    nazwa: 'Szyfrowanie',
    opis:
      'Nakłada szyfrowanie AES. Puste hasło otwarcia je zdejmuje i wymaga wtedy ' +
      'hasła właściciela.',
    pola: [
      POLE_MATERIALU,
      { kod: 'userPassword', etykieta: 'Hasło otwarcia', rodzaj: 'tajne' },
      { kod: 'ownerPassword', etykieta: 'Hasło właściciela', rodzaj: 'tajne' },
      {
        kod: 'permissions',
        etykieta: 'Uprawnienia odbiorcy',
        rodzaj: 'tekst',
        podpowiedz: 'drukowanie, kopiowanie',
        opis: 'Przyjmowane: drukowanie, zmiana, kopiowanie, komentowanie, wypelnianie, skladanie.',
      },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          {
            userPassword: tekst(wartosci, 'userPassword'),
            ownerPassword: tekst(wartosci, 'ownerPassword'),
            permissions: wykaz(wartosci, 'permissions'),
          },
        ),
      };
    },
  },
  {
    komenda: Command.StudioSecurityMetadataStrip,
    nazwa: 'Wyczyść metadane',
    opis: 'Zdejmuje opis, właściwości własne, strumień XMP i ślad narzędzi.',
    pola: [
      POLE_MATERIALU,
      {
        kod: 'keepFields',
        etykieta: 'Pola zostawiane',
        rodzaj: 'tekst',
        podpowiedz: 'Title, CreationDate',
      },
    ],
    zloz(wartosci, otoczenie) {
      return {
        zadanie: zPolami(
          { assetId: wartosci['assetId'] ?? '', windowId: otoczenie.idOkna },
          { keepFields: wykaz(wartosci, 'keepFields') },
        ),
      };
    },
  },
  {
    komenda: Command.StudioSecurityRedact,
    nazwa: 'Redakcja',
    opis:
      'Wycina tekst z obszaru, dopiero potem kładzie czarne pole. Współrzędne ' +
      'liczą się od lewego DOLNEGO rogu strony, w punktach.',
    pola: [
      POLE_MATERIALU,
      {
        kod: 'regions',
        etykieta: 'Obszary',
        rodzaj: 'wielowiersz',
        podpowiedz: '1 | 72 | 640 | 200 | 24',
        opis: 'Wiersz na obszar: strona, x, y, szerokość, wysokość — oddzielone kreskami.',
      },
    ],
    zloz(wartosci, otoczenie) {
      const obszary = [];
      for (const wiersz of wiersze(wartosci, 'regions')) {
        const czesci = wiersz.split('|').map((czesc) => Number(czesc.trim()));
        if (czesci.length !== 5 || czesci.some((czesc) => !Number.isFinite(czesc))) {
          return {
            odmowa: `Wiersz „${wiersz}" nie niesie pięciu liczb: strona, x, y, szerokość, wysokość.`,
          };
        }
        obszary.push({
          page: czesci[0],
          x: czesci[1],
          y: czesci[2],
          width: czesci[3],
          height: czesci[4],
        });
      }
      if (obszary.length === 0) {
        return { odmowa: 'Nie wskazano ani jednego obszaru.' };
      }
      return {
        zadanie: {
          assetId: wartosci['assetId'] ?? '',
          windowId: otoczenie.idOkna,
          regions: obszary,
        },
      };
    },
  },
  {
    komenda: Command.StudioSecuritySign,
    nazwa: 'Podpisz',
    opis:
      'Podpis tego produktu, liczony po treści stron. Cudze czytniki go nie ' +
      'pokażą — podpisem kwalifikowanym nie jest.',
    pola: [
      POLE_MATERIALU,
      {
        kod: 'certificateId',
        etykieta: 'Certyfikat',
        rodzaj: 'zasob',
        opis: 'Zasób niosący certyfikat i klucz prywatny w postaci PEM.',
      },
      { kod: 'reason', etykieta: 'Powód podpisu', rodzaj: 'tekst' },
      { kod: 'location', etykieta: 'Miejsce', rodzaj: 'tekst' },
    ],
    zloz(wartosci, otoczenie) {
      const certyfikat = tekst(wartosci, 'certificateId');
      if (certyfikat === undefined) {
        return { odmowa: 'Podpis potrzebuje certyfikatu.' };
      }
      return {
        zadanie: zPolami(
          {
            assetId: wartosci['assetId'] ?? '',
            certificateId: certyfikat,
            windowId: otoczenie.idOkna,
          },
          { reason: tekst(wartosci, 'reason'), location: tekst(wartosci, 'location') },
        ),
      };
    },
  },
  {
    komenda: Command.StudioSecuritySignVerify,
    nazwa: 'Sprawdź podpisy',
    opis: 'Przelicza skrót treści i porównuje go ze skrótem podpisanym.',
    pola: [POLE_MATERIALU],
    zloz(wartosci) {
      return { zadanie: { assetId: wartosci['assetId'] ?? '' } };
    },
  },
  {
    komenda: Command.StudioSecuritySensitiveDetect,
    nazwa: 'Rozpoznaj dane wrażliwe',
    opis:
      'Czyta treść dokumentu wczytanego w edytorze i niczego nie zmienia. ' +
      'Położenie oddaje w znakach, więc wchodzi wprost jako zakres redakcji.',
    pola: [
      {
        kod: 'categories',
        etykieta: 'Kategorie',
        rodzaj: 'tekst',
        podpowiedz: 'poczta, karta-platnicza',
        opis:
          'Rozpoznawane: poczta, karta-platnicza, numer-ewidencyjny, numer-podatkowy, ' +
          'rachunek-bankowy, telefon. Puste bierze wszystkie.',
      },
    ],
    zloz(wartosci, otoczenie) {
      if (otoczenie.idDokumentu === null) {
        return {
          odmowa:
            'Rozpoznanie czyta dokument Studia — wczytaj dokument w Studio Editorze.',
        };
      }
      return {
        zadanie: zPolami(
          { documentId: otoczenie.idDokumentu },
          { categories: wykaz(wartosci, 'categories') },
        ),
      };
    },
  },
];

/** Zdanie o skutku czynności złożone z odpowiedzi rdzenia — każde pole odczytywane osobno, bo odpowiedź może go nie nieść. */
export function opiszSkutek(odpowiedz: unknown): string {
  if (typeof odpowiedz !== 'object' || odpowiedz === null) return 'Rdzeń przyjął czynność.';
  const tresc = odpowiedz as Record<string, unknown>;
  const czesci: string[] = [];

  const zasob = tresc['asset'];
  if (typeof zasob === 'object' && zasob !== null) {
    const identyfikator = (zasob as Record<string, unknown>)['id'];
    if (typeof identyfikator === 'string') czesci.push(`nowy zasób ${identyfikator}`);
  }
  if (typeof tresc['pages'] === 'number') czesci.push(`stron: ${tresc['pages']}`);
  if (typeof tresc['parts'] === 'number') czesci.push(`części: ${tresc['parts']}`);
  if (typeof tresc['extracted'] === 'number') czesci.push(`wyciągnięto: ${tresc['extracted']}`);
  if (typeof tresc['redacted'] === 'number') czesci.push(`obszarów zamazanych: ${tresc['redacted']}`);
  if (typeof tresc['lastNumber'] === 'number') czesci.push(`ostatni numer: ${tresc['lastNumber']}`);
  if (typeof tresc['savedBytes'] === 'number') czesci.push(`ubyło bajtów: ${tresc['savedBytes']}`);
  if (typeof tresc['encrypted'] === 'boolean') {
    czesci.push(tresc['encrypted'] === true ? 'dokument zaszyfrowany' : 'szyfrowanie zdjęte');
  }
  if (typeof tresc['allValid'] === 'boolean') {
    czesci.push(tresc['allValid'] === true ? 'wszystkie podpisy poprawne' : 'podpisy niepoprawne');
  }
  if (Array.isArray(tresc['assetIds'])) czesci.push(`zasobów: ${tresc['assetIds'].length}`);
  if (Array.isArray(tresc['fields'])) czesci.push(`pól formularza: ${tresc['fields'].length}`);
  if (Array.isArray(tresc['findings'])) czesci.push(`znalezisk: ${tresc['findings'].length}`);
  if (Array.isArray(tresc['signatures'])) czesci.push(`podpisów: ${tresc['signatures'].length}`);
  if (Array.isArray(tresc['removed'])) czesci.push(`pól usuniętych: ${tresc['removed'].length}`);

  return czesci.length === 0 ? 'Rdzeń przyjął czynność.' : czesci.join(', ') + '.';
}
