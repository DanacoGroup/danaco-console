import { describe, expect, it } from 'vitest';

import { Command, DesignTokenKind, DesignTokenTarget } from '../../../../shared/contract';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { CZYNNOSCI_WARSZTATOW_DESIGNU } from './czynnosci-warsztatow-designu';
import { utworzZrodloDesignu } from './zrodlo-designu';
import { utworzZrodloWarsztatowDesignu } from './zrodlo-warsztatow-designu';

/**
 * Sprawdziany warstwy klienckiej obszaru `design.*`.
 *
 * Główny z nich pilnuje jednej rzeczy: czy KAŻDA komenda rodziny `design.*`
 * z kontraktu ma drogę z okna do rdzenia. Wykaz nie jest tu przepisany —
 * powstaje z wywołań źródła, a porównywany jest ze stałymi kontraktu, więc
 * komenda dołożona do kontraktu i pominięta w oknie zostanie tu nazwana.
 *
 * Dróg z okna do rdzenia są DWIE i sprawdzian przechodzi obie. Pierwsza to
 * źródło obszaru (`zrodlo-designu.ts`) — metoda na komendę, typowana kontraktem;
 * tą drogą jadą okna warstwy pierwszej. Druga to warsztaty
 * (`zrodlo-warsztatow-designu.ts`) — jedna droga na komendę wskazaną KATALOGIEM
 * czynności; tą jadą sześć okien warsztatowych, bo siedemdziesiąt pięć metod
 * w jednym źródle byłoby siedemdziesięcioma pięcioma miejscami na tę samą
 * pomyłkę. Wymaganie sprawdzianu jest w obu przypadkach to samo: komenda bez
 * drogi z okna jest funkcją, której Operator nie ma.
 *
 * Pozostałe sprawdziany dotyczą rozstrzygnięć, które warstwa kliencka
 * podejmuje sama i które łatwo cofnąć nieuważną poprawką: pola opcjonalne idą
 * do rdzenia WYŁĄCZNIE wskazane. Skala zero nie jest krotnością, granica zero
 * nie jest granicą, a `resolved: false` wysłane przy zapisie treści otwierałoby
 * wątek, którego nikt nie kazał otwierać. Każde z tych pól wysłane „na wszelki
 * wypadek" jest zdaniem o woli Operatora, którego Operator nie wypowiedział.
 */

/** Kanał próbny: zapamiętuje komendy wraz z żądaniami i oddaje odpowiedź pustą. */
function kanalProbny(odpowiedzi: Record<string, unknown> = {}): {
  kanal: Kanal;
  wyslane: string[];
  zadania: Record<string, unknown>;
} {
  const wyslane: string[] = [];
  const zadania: Record<string, unknown> = {};
  const kanal = {
    wyslij(komenda: string, zadanie: unknown, przyWyniku?: (wynik: Wynik<unknown>) => void) {
      wyslane.push(komenda);
      zadania[komenda] = zadanie;
      przyWyniku?.({ udany: true, wynik: odpowiedzi[komenda] ?? {} });
      return `zadanie-${wyslane.length}`;
    },
    naZdarzenie: () => () => undefined,
    naDowolny: () => () => undefined,
    sesja: () => ({}) as ReturnType<Kanal['sesja']>,
    dziennikNieznanych: () => ({}) as ReturnType<Kanal['dziennikNieznanych']>,
  } as unknown as Kanal;
  return { kanal, wyslane, zadania };
}

/** Wywołuje każdą czynność źródła raz — pełny przelot rodziny. */
async function przelotZrodla(kanal: Kanal): Promise<void> {
  const zrodlo = utworzZrodloDesignu(kanal);
  const prompt = { subject: 'ilustracja bohatera strony' };

  await zrodlo.zasoby({
    idOkna: 'okno-1',
    rodzaj: '',
    etykiety: [],
    tylkoUlubione: false,
    granica: 0,
  });
  await zrodlo.generuj({ idOkna: 'okno-1', prompt, idReferencji: '', idKanalu: '', rodzaj: '' });
  await zrodlo.zapiszKompozycje({ idOkna: 'okno-1', idKompozycji: '', nazwa: '', warstwy: [] });
  await zrodlo.ustawEtykiety({ idZasobu: 'zasob-1', etykiety: [] });
  await zrodlo.wgrajZasob({
    idOkna: 'okno-1',
    nazwa: '',
    rodzaj: 'image',
    trescBase64: 'AA==',
    format: '',
    szerokosc: 0,
    wysokosc: 0,
    etykiety: [],
  });
  await zrodlo.usunZasob('zasob-1');
  await zrodlo.ustawUlubiony('zasob-1', true);
  await zrodlo.kompozycje('okno-1');

  await zrodlo.trescZasobu('zasob-1', 0);
  await zrodlo.wydajZasob({ idZasobu: 'zasob-1', format: 'png', skala: 0, jakosc: 0 });
  await zrodlo.wydajPartie({ idZasobow: ['zasob-1'], format: 'png', skale: [], jakosc: 0 });

  await zrodlo.zalozKolekcje('okno-1', 'Kampania', '');
  await zrodlo.przypiszDoKolekcji({
    idKolekcji: 'kolekcja-1',
    idZasobow: ['zasob-1'],
    zdejmij: false,
  });
  await zrodlo.kolekcje('okno-1', '');

  await zrodlo.zapiszSzablonPromptu({ idOkna: 'okno-1', nazwa: 'Szablon', prompt, idSzablonu: '' });
  await zrodlo.szablonyPromptu('okno-1');
  await zrodlo.historiaPromptow('okno-1', 0);

  await zrodlo.zapiszWersje({ idKompozycji: 'plansza-1', nazwa: '', uzasadnienie: '' });
  await zrodlo.wersje('plansza-1', 0);
  await zrodlo.przywrocWersje('wersja-1');
  await zrodlo.wyrysujKompozycje({
    idKompozycji: 'plansza-1',
    format: 'png',
    skala: 0,
    obszar: null,
  });

  await zrodlo.ustawAdnotacje({
    idKompozycji: 'plansza-1',
    tresc: 'margines za wąski',
    idWarstwy: '',
    idAdnotacji: '',
    idNadrzednej: '',
    zamknieta: null,
  });
  await zrodlo.adnotacje('plansza-1', false);
  await zrodlo.zglosObecnosc({
    idKompozycji: 'plansza-1',
    x: 0,
    y: 0,
    zaznaczone: [],
    odchodzi: false,
  });

  await zrodlo.zapiszZestawZetonow({
    idOkna: 'okno-1',
    nazwa: 'System marki',
    zetony: [{ name: 'sygnal', kind: DesignTokenKind.Color, value: '#c8a24a' }],
    idZestawu: '',
    motyw: '',
  });
  await zrodlo.zestawyZetonow('okno-1', '');
  await zrodlo.wydajZetony('zestaw-1', DesignTokenTarget.CssVariables, '');
  await zrodlo.wczytajZetony({
    idOkna: 'okno-1',
    nazwa: 'Wczytany',
    trescBase64: 'e30=',
    postac: '',
  });
  await zrodlo.wydajPrzewodnik('zestaw-1', 'library', '');
}

/**
 * Przelot drugą drogą: każda czynność katalogu warsztatów woła rdzeń swoją
 * komendą.
 *
 * Żądanie jest tu puste, bo sprawdzian pyta o DROGĘ, nie o kształt żądania.
 * Kształt składa `zloz` czynności z formularza i sprawdza go rdzeń — a odmowa
 * walidacji z nazwą pola jest odpowiedzią, którą kanał próbny i tak zwraca jako
 * powodzenie. Ten sprawdzian ma wyłapać komendę, której NIKT nie woła.
 */
async function przelotWarsztatow(kanal: Kanal): Promise<void> {
  const zrodlo = utworzZrodloWarsztatowDesignu(kanal);
  for (const czynnosc of CZYNNOSCI_WARSZTATOW_DESIGNU) {
    await zrodlo.wykonaj(czynnosc.komenda, {});
  }
}

describe('źródło modułu Design', () => {
  it('ma drogę z okna do każdej komendy rodziny design.*', async () => {
    const { kanal, wyslane } = kanalProbny();
    await przelotZrodla(kanal);
    await przelotWarsztatow(kanal);

    // Wykaz oczekiwany bierze się z KONTRAKTU, nie z tego pliku: komenda
    // dołożona do rodziny i pominięta w źródle wypadnie tu jako brak, bez
    // dopisywania czegokolwiek w sprawdzianie.
    const zRodziny = Object.values(Command).filter((nazwa) => nazwa.startsWith('design.'));
    const bezDrogi = zRodziny.filter((nazwa) => !wyslane.includes(nazwa));

    expect(bezDrogi, `komendy design.* bez drogi z okna: ${bezDrogi.join(', ')}`).toEqual([]);
  });

  it('nie wysyła pól opcjonalnych, których Operator nie wskazał', async () => {
    const { kanal, zadania } = kanalProbny();
    await przelotZrodla(kanal);

    // Granica zero nie jest granicą — pole ma zostać pominięte.
    expect(zadania[Command.DesignAssetContentGet]).not.toHaveProperty('maxBytes');
    // Skala zero nie jest krotnością, jakość zero nie jest stopniem kompresji.
    expect(zadania[Command.DesignAssetExport]).not.toHaveProperty('scale');
    expect(zadania[Command.DesignAssetExport]).not.toHaveProperty('quality');
    // Pusty wykaz skal zostawia rdzeniowi skalę naturalną.
    expect(zadania[Command.DesignAssetExportBatch]).not.toHaveProperty('scales');
    // Dokładka jest brakiem pola `remove`, tak stanowi kontrakt.
    expect(zadania[Command.DesignCollectionAssign]).not.toHaveProperty('remove');
    // Stan zamknięcia idzie wyłącznie wskazany.
    expect(zadania[Command.DesignAnnotationSet]).not.toHaveProperty('resolved');
    // Zawężenie do wątków niezamkniętych idzie wyłącznie włączone.
    expect(zadania[Command.DesignAnnotationList]).not.toHaveProperty('openOnly');
    // Postać zapisu rozpoznaje rdzeń, gdy okno jej nie wskazało.
    expect(zadania[Command.DesignTokensetImport]).not.toHaveProperty('format');
    // Moduł docelowy pusty zostawia wydanie wołającemu.
    expect(zadania[Command.DesignTokensetExport]).not.toHaveProperty('targetModuleId');
  });

  it('wysyła położenie kursora także w rogu kanwy', async () => {
    const { kanal, zadania } = kanalProbny();
    await przelotZrodla(kanal);

    // Zero jest tu prawdziwym położeniem (lewy górny róg), więc rozstrzyga
    // skończoność liczby, nie jej wartość — inaczej kursor w rogu znikałby
    // z widoku pozostałych.
    expect(zadania[Command.DesignPresenceReport]).toMatchObject({ x: 0, y: 0 });
  });

  it('wysyła zestaw etykiet także pusty — inaczej nie da się zdjąć ostatniej', async () => {
    const { kanal, zadania } = kanalProbny();
    await przelotZrodla(kanal);

    expect(zadania[Command.DesignAssetTagSet]).toMatchObject({ tags: [] });
  });

  it('uznaje wydanie bez treści za odpowiedź niepełną', async () => {
    // Koperta udana z pustym `contentBase64` jest kopertą kłamiącą: plik
    // zerowej długości zapisze się u Operatora tak samo jak plik prawdziwy.
    const { kanal } = kanalProbny({
      [Command.DesignAssetExport]: { contentBase64: '', fileName: 'a.png', mediaType: 'image/png', sizeBytes: 0 },
    });
    const zrodlo = utworzZrodloDesignu(kanal);

    const wynik = await zrodlo.wydajZasob({
      idZasobu: 'zasob-1',
      format: 'png',
      skala: 0,
      jakosc: 0,
    });

    expect(wynik.udany).toBe(false);
  });

  it('uznaje wydanie przewodnika bez zasobu za odpowiedź niepełną', async () => {
    // `published: true` bez identyfikatora zasobu jest obietnicą wydania,
    // którego nie da się odnaleźć.
    const { kanal } = kanalProbny({
      [Command.DesignStyleguidePublish]: { assetId: '', published: true },
    });
    const zrodlo = utworzZrodloDesignu(kanal);

    const wynik = await zrodlo.wydajPrzewodnik('zestaw-1', 'library', '');

    expect(wynik.udany).toBe(false);
  });

  it('przyjmuje treść zasobu w postaci odsyłania, bez bajtów', async () => {
    // Postać odsyłania oddaje miarę i sumę bez bajtów — to jest odpowiedź
    // PEŁNA, a sprawdzenie kształtu nie ma prawa uznać jej za ubytek.
    const { kanal } = kanalProbny({
      [Command.DesignAssetContentGet]: {
        assetId: 'zasob-1',
        mediaType: 'image/png',
        sizeBytes: 4096,
        checksum: 'ab'.repeat(32),
      },
    });
    const zrodlo = utworzZrodloDesignu(kanal);

    const wynik = await zrodlo.trescZasobu('zasob-1', 0);

    expect(wynik.udany).toBe(true);
  });
});
