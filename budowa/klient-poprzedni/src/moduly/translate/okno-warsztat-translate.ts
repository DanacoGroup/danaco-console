import { Command } from '../../../../shared/contract';
import { opisOdmowy } from '../../komponenty/odmowa';
import {
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
  type PoleFormularza,
} from '../../modele/kontrolki-formularza';
import { naglowekOkna } from './kontrolki-translate';
import { utworzStanOkna, type StanOkna } from './stan-okna-translate';
import type { StanTranslate } from './stan-translate';
import { utworzRozwiniecie } from './warstwy-translate';
import type { WynikTranslate } from './wywolanie-translate';

/** Warsztat modułu Translate — okno, przez które przechodzą rodziny komend spoza czterech okien pierwotnych, ze wspólnymi wskazaniami na górze. */
export interface OknoWarsztatTranslate {
  element: HTMLElement;
  odswiez(): void;
}

/** Opis jednej czynności warsztatu: nazwa przycisku, pola żądania i funkcja jej wykonania w warsztacie. */
interface Czynnosc {
  nazwa: string;
  wyjasnienie: string;
  pola: HTMLElement[];
  wykonaj(): Promise<{ udane: boolean; zdanie: string }>;
}

export function utworzOknoWarsztatTranslate(stan: StanTranslate): OknoWarsztatTranslate {
  const okno: StanOkna = utworzStanOkna({
    tytul: 'Warsztat czeka na czynność',
    opis:
      'Wybierz czynność niżej. Każda z nich kończy się zdaniem o skutku: ' +
      'liczbą wierszy, które doszły, albo ścieżką pliku, który powstał.',
  });
  const odpowiedz = utworzWierszOdpowiedzi();
  const warsztat = stan.warsztat;

  const panel = poleWyboru({ etykieta: 'Panel języka docelowego' }, []);
  const sciezka = poleTekstowe({
    etykieta: 'Ścieżka pliku',
    podpowiedz: '/dane/materialy/umowa.docx',
    opis: 'Wskazanie wspólne dla wczytania materiału i dla zapisu wyniku.',
  });

  /** idOkna oddaje okno modułu albo pusty napis — czynność sama nazwie brak. */
  function idOkna(): string {
    return stan.idOkna();
  }

  /** wymagajOkna zwraca zdanie o braku okna albo pusty napis, gdy okno jest. */
  function wymagajOkna(): string {
    return idOkna() === ''
      ? 'Rdzeń nie wskazał okna tłumaczenia — czynność nie ma na czym pracować.'
      : '';
  }

  /** wymagajPanelu zwraca zdanie o braku panelu albo pusty napis. */
  function wymagajPanelu(): string {
    return panel.kontrolka.value === ''
      ? 'Wskaż panel języka docelowego — czynność pracuje na przekładzie, nie na źródle.'
      : '';
  }

  /** wymagajSciezki zwraca zdanie o braku ścieżki albo pusty napis. */
  function wymagajSciezki(): string {
    return sciezka.kontrolka.value.trim() === ''
      ? 'Podaj ścieżkę pliku — rdzeń nie zgaduje, skąd wziąć materiał ani gdzie zapisać wynik.'
      : '';
  }

  // ── Pamięć tłumaczeń ─────────────────────────────────────────────────────
  const jezykPary = poleTekstowe({ etykieta: 'Język pary', podpowiedz: 'angielski' });
  const segmentZrodlowy = poleTekstowe({ etykieta: 'Segment źródłowy' });
  const segmentDocelowy = poleTekstowe({ etykieta: 'Segment docelowy' });
  const projektPary = poleTekstowe({ etykieta: 'Projekt (nieobowiązkowo)' });
  const progDopasowania = poleTekstowe({
    etykieta: 'Próg dopasowania [%]',
    podpowiedz: '75',
    opis: 'Poniżej progu para nie wchodzi do tłumaczenia wstępnego.',
  });
  const rodzajUtrzymania = poleWyboru({ etykieta: 'Rodzaj utrzymania pamięci' }, [
    { wartosc: 'deduplicate', etykieta: 'Czyszczenie duplikatów' },
    { wartosc: 'merge', etykieta: 'Scalanie par o wspólnym źródle' },
    { wartosc: 'replace', etykieta: 'Masowa podmiana fragmentu' },
    { wartosc: 'filter', etykieta: 'Usunięcie par spełniających zawężenie' },
  ]);
  const szukanyFragment = poleTekstowe({ etykieta: 'Szukany fragment' });
  const zamiennik = poleTekstowe({ etykieta: 'Zamiennik' });
  const tekstZrodlowyWyrownania = poleTekstowe({ etykieta: 'Tekst źródłowy do wyrównania' });
  const tekstDocelowyWyrownania = poleTekstowe({ etykieta: 'Tekst docelowy do wyrównania' });

  // ── Segmentacja ──────────────────────────────────────────────────────────
  const numerySegmentow = poleTekstowe({
    etykieta: 'Numery segmentów do scalenia',
    podpowiedz: '0, 1',
  });
  const numerSegmentu = poleTekstowe({ etykieta: 'Numer segmentu do podziału', podpowiedz: '0' });
  const miejscePodzialu = poleTekstowe({ etykieta: 'Miejsce podziału [znak]', podpowiedz: '16' });
  const nazwaZestawu = poleTekstowe({ etykieta: 'Nazwa zestawu reguł' });

  // ── Kontrola i obieg ─────────────────────────────────────────────────────
  const etapObiegu = poleWyboru({ etykieta: 'Etap obiegu' }, [
    { wartosc: 'translation', etykieta: 'Tłumaczenie' },
    { wartosc: 'proofreading', etykieta: 'Korekta' },
    { wartosc: 'approved', etykieta: 'Zatwierdzone' },
    { wartosc: 'rejected', etykieta: 'Odrzucone do poprawy' },
  ]);
  const uwagaObiegu = poleTekstowe({ etykieta: 'Uwaga obiegu (nieobowiązkowo)' });
  const nazwaProfiluQa = poleTekstowe({ etykieta: 'Nazwa profilu kontroli jakości' });

  // ── Materiał ─────────────────────────────────────────────────────────────
  const formatNapisow = poleWyboru({ etykieta: 'Format napisów' }, [
    { wartosc: 'srt', etykieta: 'SubRip (SRT)' },
    { wartosc: 'webvtt', etykieta: 'WebVTT' },
    { wartosc: 'ttml', etykieta: 'TTML' },
    { wartosc: 'stl', etykieta: 'EBU STL' },
  ]);
  const kluczZasobu = poleTekstowe({ etykieta: 'Klucz zasobu' });
  const kontekstKlucza = poleTekstowe({ etykieta: 'Kontekst klucza' });
  const idDokumentu = poleTekstowe({ etykieta: 'Identyfikator dokumentu' });
  const idZasobu = poleTekstowe({ etykieta: 'Identyfikator zasobu lokalizacyjnego' });

  // ── Silniki i wymiana ────────────────────────────────────────────────────
  const nazwaProfiluSilnika = poleTekstowe({ etykieta: 'Nazwa profilu silnika' });
  const kanalyProfilu = poleTekstowe({
    etykieta: 'Kanały modelu',
    podpowiedz: 'kanal-a, kanal-b',
  });
  const segmentPorownania = poleTekstowe({ etykieta: 'Segment do porównania silników' });
  const jezykPivota = poleTekstowe({ etykieta: 'Domyślny język pośredni' });
  const rodzajWytworu = poleWyboru({ etykieta: 'Rodzaj wytworu' }, [
    { wartosc: 'targetFile', etykieta: 'Plik w języku docelowym' },
    { wartosc: 'bilingualFile', etykieta: 'Plik dwujęzyczny' },
    { wartosc: 'memory', etykieta: 'Pamięć tłumaczeń (TMX)' },
    { wartosc: 'termbase', etykieta: 'Baza terminologiczna (TBX)' },
  ]);
  const instrukcjeWykonawcy = poleTekstowe({ etykieta: 'Instrukcje dla wykonawcy' });
  const trybOdeslania = poleWyboru({ etykieta: 'Tryb odesłania wyniku' }, [
    { wartosc: 'bilingual', etykieta: 'Dwujęzycznie' },
    { wartosc: 'targetOnly', etykieta: 'Sam przekład' },
    { wartosc: 'appendix', etykieta: 'Jako załącznik' },
  ]);
  const idDokumentuMostu = poleTekstowe({ etykieta: 'Dokument mostu' });
  const trescMostu = poleTekstowe({ etykieta: 'Treść wnoszona mostem' });

  const wykaz = document.createElement('ol');
  wykaz.className = 'mt-warsztat__wykaz';

  /** pokazWykaz wypisuje wynik czynności wierszami — skutek widoczny, nie sam licznik. */
  function pokazWykaz(pozycje: readonly string[]): void {
    wykaz.replaceChildren(
      ...pozycje.map((tresc) => {
        const wiersz = document.createElement('li');
        wiersz.textContent = tresc;
        return wiersz;
      }),
    );
  }

  /** liczbaZPola czyta liczbę z pola albo oddaje `undefined`, gdy pole puste. */
  function liczbaZPola(pole: PoleFormularza<HTMLInputElement>): number | undefined {
    const tresc = pole.kontrolka.value.trim();
    if (tresc === '') return undefined;
    const wartosc = Number.parseInt(tresc, 10);
    return Number.isNaN(wartosc) ? undefined : wartosc;
  }

  /** tekstZPola czyta treść pola albo oddaje `undefined`, gdy pole puste. */
  function tekstZPola(pole: PoleFormularza<HTMLInputElement>): string | undefined {
    const tresc = pole.kontrolka.value.trim();
    return tresc === '' ? undefined : tresc;
  }

  /** wykazZPola rozbija treść pola po przecinkach na wykaz wskazań. */
  function wykazZPola(pole: PoleFormularza<HTMLInputElement>): string[] {
    return pole.kontrolka.value
      .split(',')
      .map((czesc) => czesc.trim())
      .filter((czesc) => czesc !== '');
  }

  /** zdanieOdmowy składa jedno brzmienie odmowy dla wszystkich czynności. */
  function zdanieOdmowy<T>(nazwa: string, wynik: WynikTranslate<T>): string {
    if (wynik.odmowa !== undefined) {
      return `${nazwa}: rdzeń nie obsługuje komendy ${wynik.odmowa.zadanyTyp}.`;
    }
    return opisOdmowy(nazwa, wynik.blad?.code, wynik.blad?.message);
  }

  const czynnosci: Czynnosc[] = [
    {
      nazwa: 'Wykaz pamięci tłumaczeń',
      wyjasnienie: `${Command.TranslateMemoryList} — pary zawężone językiem i projektem.`,
      pola: [jezykPary.element, projektPary.element],
      async wykonaj() {
        const wynik = await warsztat.wykazPamieci({
          language: tekstZPola(jezykPary),
          project: tekstZPola(projektPary),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wykaz pamięci', wynik) };
        }
        pokazWykaz(
          wynik.wynik.entries.map(
            (wpis) => `${wpis.language}: „${wpis.sourceSegment}” → „${wpis.targetSegment}”`,
          ),
        );
        return {
          udane: true,
          zdanie: `Pamięć niesie ${String(wynik.wynik.total)} par; pokazano ${String(
            wynik.wynik.entries.length,
          )}.`,
        };
      },
    },
    {
      nazwa: 'Dopisz parę do pamięci',
      wyjasnienie: `${Command.TranslateMemorySet} — para wniesiona wprost przez Operatora.`,
      pola: [jezykPary.element, segmentZrodlowy.element, segmentDocelowy.element, projektPary.element],
      async wykonaj() {
        const wynik = await warsztat.zapiszPare({
          language: jezykPary.kontrolka.value.trim(),
          sourceSegment: segmentZrodlowy.kontrolka.value,
          targetSegment: segmentDocelowy.kontrolka.value,
          project: tekstZPola(projektPary),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Zapis pary pamięci', wynik) };
        }
        return { udane: true, zdanie: `Para zapisana pod identyfikatorem ${wynik.wynik.entry.id}.` };
      },
    },
    {
      nazwa: 'Import pamięci z pliku',
      wyjasnienie: `${Command.TranslateMemoryImport} — plik TMX albo CSV spod wskazanej ścieżki.`,
      pola: [sciezka.element, projektPary.element],
      async wykonaj() {
        const brak = wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.importujPamiec({
          path: sciezka.kontrolka.value.trim(),
          project: tekstZPola(projektPary),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Import pamięci', wynik) };
        }
        return {
          udane: true,
          zdanie: `Wniesiono ${String(wynik.wynik.importedCount)} par, pominięto ${String(
            wynik.wynik.skippedCount,
          )}.`,
        };
      },
    },
    {
      nazwa: 'Eksport pamięci do pliku',
      wyjasnienie: `${Command.TranslateMemoryExport} — plik TMX (albo CSV wedle końcówki nazwy).`,
      pola: [sciezka.element, jezykPary.element],
      async wykonaj() {
        const brak = wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.eksportujPamiec({
          path: sciezka.kontrolka.value.trim(),
          language: tekstZPola(jezykPary),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Eksport pamięci', wynik) };
        }
        return {
          udane: true,
          zdanie: `Zapisano ${String(wynik.wynik.exportedCount)} par pod ${wynik.wynik.path}.`,
        };
      },
    },
    {
      nazwa: 'Utrzymanie pamięci',
      wyjasnienie:
        `${Command.TranslateMemoryMaintain} — czyszczenie, scalanie, podmiana i usuwanie par. ` +
        'Przebieg próbny liczy to samo, czego nie zapisuje.',
      pola: [rodzajUtrzymania.element, szukanyFragment.element, zamiennik.element],
      async wykonaj() {
        const wynik = await warsztat.utrzymajPamiec({
          kind: rodzajUtrzymania.kontrolka.value as never,
          search: tekstZPola(szukanyFragment),
          replacement: tekstZPola(zamiennik),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Utrzymanie pamięci', wynik) };
        }
        return {
          udane: true,
          zdanie: `Czynność objęła ${String(wynik.wynik.affectedCount)} par.`,
        };
      },
    },
    {
      nazwa: 'Tłumaczenie wstępne z pamięci',
      wyjasnienie: `${Command.TranslateMemoryPretranslate} — wypełnia panele parami powyżej progu.`,
      pola: [progDopasowania.element],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.tlumaczWstepnie({
          windowId: idOkna(),
          threshold: liczbaZPola(progDopasowania),
          onlyEmpty: false,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Tłumaczenie wstępne', wynik) };
        }
        return {
          udane: true,
          zdanie: `Wypełniono ${String(wynik.wynik.filledCount)} paneli.`,
        };
      },
    },
    {
      nazwa: 'Wyrównanie tekstów dwujęzycznych',
      wyjasnienie: `${Command.TranslateMemoryAlign} — pary segmentów zapisane w pamięci.`,
      pola: [jezykPary.element, tekstZrodlowyWyrownania.element, tekstDocelowyWyrownania.element],
      async wykonaj() {
        const wynik = await warsztat.wyrownajTeksty({
          sourceText: tekstZrodlowyWyrownania.kontrolka.value,
          targetText: tekstDocelowyWyrownania.kontrolka.value,
          language: jezykPary.kontrolka.value.trim(),
          commit: true,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wyrównanie tekstów', wynik) };
        }
        pokazWykaz(
          wynik.wynik.pairs.map(
            (para) => `${String(para.score)}%: „${para.sourceSegment}” → „${para.targetSegment}”`,
          ),
        );
        return {
          udane: true,
          zdanie: `Sparowano ${String(wynik.wynik.pairs.length)} segmentów, zapisano ${String(
            wynik.wynik.committedCount,
          )}.`,
        };
      },
    },
    {
      nazwa: 'Polityka pamięci okna',
      wyjasnienie: `${Command.TranslateMemoryPolicySet} — próg dopasowania obowiązujący w tym oknie.`,
      pola: [progDopasowania.element],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.ustawPolitykePamieci({
          windowId: idOkna(),
          threshold: liczbaZPola(progDopasowania),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Polityka pamięci', wynik) };
        }
        return {
          udane: true,
          zdanie: `Próg okna: ${String(wynik.wynik.policy.threshold)}%, zasięg: ${wynik.wynik.policy.scope}.`,
        };
      },
    },
    {
      nazwa: 'Scal segmenty',
      wyjasnienie: `${Command.TranslateSegmentMerge} — trwała zmiana podziału okna.`,
      pola: [numerySegmentow.element],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.scalSegmenty({
          windowId: idOkna(),
          segmentIndexes: wykazZPola(numerySegmentow),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Scalenie segmentów', wynik) };
        }
        pokazWykaz(wynik.wynik.segments);
        return { udane: true, zdanie: `Okno ma ${String(wynik.wynik.segments.length)} segmentów.` };
      },
    },
    {
      nazwa: 'Podziel segment',
      wyjasnienie: `${Command.TranslateSegmentSplit} — trwała zmiana podziału okna.`,
      pola: [numerSegmentu.element, miejscePodzialu.element],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.podzielSegment({
          windowId: idOkna(),
          segmentIndex: liczbaZPola(numerSegmentu) ?? 0,
          offset: liczbaZPola(miejscePodzialu) ?? 0,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Podział segmentu', wynik) };
        }
        pokazWykaz(wynik.wynik.segments);
        return { udane: true, zdanie: `Okno ma ${String(wynik.wynik.segments.length)} segmentów.` };
      },
    },
    {
      nazwa: 'Zestaw reguł segmentacji',
      wyjasnienie: `${Command.TranslateSegmentationRulesSet} — nastawa wielokrotnego użytku.`,
      pola: [nazwaZestawu.element, jezykPary.element],
      async wykonaj() {
        const wynik = await warsztat.ustawRegulySegmentacji({
          name: nazwaZestawu.kontrolka.value.trim(),
          language: tekstZPola(jezykPary),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Zestaw reguł segmentacji', wynik) };
        }
        return { udane: true, zdanie: `Zestaw zapisany pod ${wynik.wynik.ruleset.id}.` };
      },
    },
    {
      nazwa: 'Kandydaci na termin',
      wyjasnienie: `${Command.TranslateTermExtract} — częstość słów i zbitek w tekście źródłowym.`,
      pola: [],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wyjmijTerminy({ windowId: idOkna() });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wyjęcie terminów', wynik) };
        }
        pokazWykaz(
          wynik.wynik.candidates.map(
            (kandydat) =>
              `${kandydat.source} — ${String(kandydat.frequency)} wystąpień` +
              (kandydat.known ? ' (już w słowniku)' : ''),
          ),
        );
        return {
          udane: true,
          zdanie: `Kandydatów na termin: ${String(wynik.wynik.candidates.length)}.`,
        };
      },
    },
    {
      nazwa: 'Korekta językowa',
      wyjasnienie: `${Command.TranslateProofreadRun} — ustalenia wraz z propozycjami poprawek.`,
      pola: [],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.sprawdzKorekte({ panelId: panel.kontrolka.value });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Korekta językowa', wynik) };
        }
        pokazWykaz([
          ...wynik.wynik.findings.map(
            (ustalenie) => `${ustalenie.severity} · ${ustalenie.kind}: ${ustalenie.detail}`,
          ),
          ...wynik.wynik.readability.map(
            (miara) => `czytelność ${miara.metric}: ${String(miara.value)}`,
          ),
        ]);
        return {
          udane: true,
          zdanie: `Ustaleń korekty: ${String(wynik.wynik.findings.length)}.`,
        };
      },
    },
    {
      nazwa: 'Zastosuj pierwszą poprawkę',
      wyjasnienie:
        `${Command.TranslateProofreadApply} — wstawia propozycję ustalenia w treść panelu. ` +
        'Bierze pierwsze ustalenie z ostatniej korekty, bo tylko ono ma pewny identyfikator.',
      pola: [],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const korekta = await warsztat.sprawdzKorekte({ panelId: panel.kontrolka.value });
        if (!korekta.udany || korekta.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Korekta językowa', korekta) };
        }
        const zPropozycja = korekta.wynik.findings.find(
          (ustalenie) => ustalenie.suggestion !== undefined && ustalenie.suggestion !== '',
        );
        if (zPropozycja === undefined) {
          return { udane: false, zdanie: 'Żadne ustalenie korekty nie niesie propozycji poprawki.' };
        }
        const wynik = await warsztat.zastosujKorekte({
          panelId: panel.kontrolka.value,
          findingIds: [zPropozycja.id],
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Zastosowanie poprawki', wynik) };
        }
        return {
          udane: true,
          zdanie: `Zastosowano ${String(wynik.wynik.appliedCount)} poprawek; panel zaktualizowany.`,
        };
      },
    },
    {
      nazwa: 'Kontrola spójności',
      wyjasnienie: `${Command.TranslateConsistencyCheck} — zdania i terminy przełożone różnie.`,
      pola: [],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.sprawdzSpojnosc({ windowId: idOkna() });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Kontrola spójności', wynik) };
        }
        pokazWykaz(
          wynik.wynik.findings.map(
            (ustalenie) => `${ustalenie.kind}: „${ustalenie.sourceText}” → ${ustalenie.variants.join(' | ')}`,
          ),
        );
        return {
          udane: true,
          zdanie: `Niespójności: ${String(wynik.wynik.findings.length)}.`,
        };
      },
    },
    {
      nazwa: 'Profil kontroli jakości',
      wyjasnienie: `${Command.TranslateQaProfileSet} — profil z kompletem sześciu kontroli.`,
      pola: [nazwaProfiluQa.element],
      async wykonaj() {
        const wynik = await warsztat.ustawProfilQa({
          name: nazwaProfiluQa.kontrolka.value.trim(),
          checks: (['number', 'date', 'currency', 'placeholder', 'length', 'omission'] as const).map(
            (rodzaj) => ({ kind: rodzaj, severity: 'warning', enabled: true }),
          ),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Profil kontroli jakości', wynik) };
        }
        return { udane: true, zdanie: `Profil zapisany pod ${wynik.wynik.profile.id}.` };
      },
    },
    {
      nazwa: 'Wykaz profili kontroli jakości',
      wyjasnienie: `${Command.TranslateQaProfileList}`,
      pola: [],
      async wykonaj() {
        const wynik = await warsztat.wykazProfiliQa({});
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wykaz profili', wynik) };
        }
        pokazWykaz(
          wynik.wynik.profiles.map(
            (profil) => `${profil.name} — kontroli: ${String(profil.checks.length)} (${profil.id})`,
          ),
        );
        return { udane: true, zdanie: `Profili: ${String(wynik.wynik.profiles.length)}.` };
      },
    },
    {
      nazwa: 'Etap zatwierdzenia panelu',
      wyjasnienie: `${Command.TranslateApprovalSet} — krok obiegu wraz z migawką na panelu.`,
      pola: [etapObiegu.element, uwagaObiegu.element],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.ustawZatwierdzenie({
          panelId: panel.kontrolka.value,
          stage: etapObiegu.kontrolka.value as never,
          note: tekstZPola(uwagaObiegu),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Etap zatwierdzenia', wynik) };
        }
        return {
          udane: true,
          zdanie: `Panel na etapie ${wynik.wynik.record.stage}; autor: ${wynik.wynik.record.author}.`,
        };
      },
    },
    {
      nazwa: 'Obieg zatwierdzeń okna',
      wyjasnienie: `${Command.TranslateApprovalList}`,
      pola: [],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wykazZatwierdzen({ windowId: idOkna() });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Obieg zatwierdzeń', wynik) };
        }
        pokazWykaz(
          wynik.wynik.records.map(
            (zapis) => `${zapis.stage} · ${zapis.author} · panel ${zapis.panelId}`,
          ),
        );
        return { udane: true, zdanie: `Kroków obiegu: ${String(wynik.wynik.records.length)}.` };
      },
    },
    {
      nazwa: 'Wczytaj dokument',
      wyjasnienie:
        `${Command.TranslateDocumentLoad} — DOCX, PDF, PPTX, XLSX, ODT, Markdown i HTML; ` +
        'treść dokumentu staje się tekstem źródłowym okna.',
      pola: [sciezka.element],
      async wykonaj() {
        const brak = wymagajOkna() || wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wczytajDokument({
          windowId: idOkna(),
          path: sciezka.kontrolka.value.trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wczytanie dokumentu', wynik) };
        }
        idDokumentu.kontrolka.value = wynik.wynik.document.id;
        pokazWykaz(wynik.wynik.segments.map((segment) => segment.text));
        return {
          udane: true,
          zdanie:
            `Dokument ${wynik.wynik.document.id}: segmentów ${String(wynik.wynik.document.segmentCount)}` +
            (wynik.wynik.usedOcr ? ' (treść z rozpoznania pisma)' : ''),
        };
      },
    },
    {
      nazwa: 'Złóż dokument wyniku',
      wyjasnienie: `${Command.TranslateDocumentRender} — plik z treści panelu.`,
      pola: [idDokumentu.element, sciezka.element],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.zlozDokument({
          documentId: idDokumentu.kontrolka.value.trim(),
          panelId: panel.kontrolka.value,
          path: tekstZPola(sciezka),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Złożenie dokumentu', wynik) };
        }
        return {
          udane: true,
          zdanie: `Zapisano ${String(wynik.wynik.renderedCount)} akapitów pod ${wynik.wynik.path}.`,
        };
      },
    },
    {
      nazwa: 'Porównaj układ',
      wyjasnienie: `${Command.TranslateDocumentLayoutCompare} — przepełnienia, przesunięcia i braki.`,
      pola: [idDokumentu.element],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.porownajUklad({
          documentId: idDokumentu.kontrolka.value.trim(),
          panelId: panel.kontrolka.value,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Porównanie układu', wynik) };
        }
        pokazWykaz(
          wynik.wynik.differences.map(
            (roznica) => `strona ${String(roznica.page)} · ${roznica.kind}: ${roznica.detail}`,
          ),
        );
        return {
          udane: true,
          zdanie: `Różnic układu: ${String(wynik.wynik.differences.length)} na ${String(
            wynik.wynik.comparedPages,
          )} stronach.`,
        };
      },
    },
    {
      nazwa: 'Wczytaj zasób lokalizacyjny',
      wyjasnienie: `${Command.TranslateResourceImport} — JSON, YAML, properties, XML, strings, RESX, PO.`,
      pola: [sciezka.element],
      async wykonaj() {
        const brak = wymagajOkna() || wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wczytajZasob({
          windowId: idOkna(),
          path: sciezka.kontrolka.value.trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wczytanie zasobu', wynik) };
        }
        idZasobu.kontrolka.value = wynik.wynik.resource.id;
        pokazWykaz(wynik.wynik.keys.map((klucz) => `${klucz.key} = ${klucz.text}`));
        return {
          udane: true,
          zdanie: `Zasób ${wynik.wynik.resource.id}: kluczy ${String(wynik.wynik.resource.keyCount)}.`,
        };
      },
    },
    {
      nazwa: 'Wydaj zasób lokalizacyjny',
      wyjasnienie: `${Command.TranslateResourceExport} — plik kluczy w języku panelu.`,
      pola: [idZasobu.element, sciezka.element],
      async wykonaj() {
        const brak = wymagajPanelu() || wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wydajZasob({
          resourceId: idZasobu.kontrolka.value.trim(),
          panelId: panel.kontrolka.value,
          path: sciezka.kontrolka.value.trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wydanie zasobu', wynik) };
        }
        return {
          udane: true,
          zdanie: `Zapisano ${String(wynik.wynik.exportedCount)} kluczy pod ${wynik.wynik.path}.`,
        };
      },
    },
    {
      nazwa: 'Formy mnogie języka panelu',
      wyjasnienie: `${Command.TranslateResourcePluralApply} — warianty CLDR wymagane przez język.`,
      pola: [idZasobu.element],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.zastosujFormyMnogie({
          resourceId: idZasobu.kontrolka.value.trim(),
          panelId: panel.kontrolka.value,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Formy mnogie', wynik) };
        }
        return {
          udane: true,
          zdanie: `Zmieniono ${String(wynik.wynik.changedCount)} kluczy.`,
        };
      },
    },
    {
      nazwa: 'Kontekst klucza',
      wyjasnienie: `${Command.TranslateResourceKeyContextSet} — opis miejsca użycia klucza.`,
      pola: [idZasobu.element, kluczZasobu.element, kontekstKlucza.element],
      async wykonaj() {
        const wynik = await warsztat.ustawKontekstKlucza({
          resourceId: idZasobu.kontrolka.value.trim(),
          key: kluczZasobu.kontrolka.value.trim(),
          context: tekstZPola(kontekstKlucza),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Kontekst klucza', wynik) };
        }
        return { udane: true, zdanie: `Klucz ${wynik.wynik.key.key} ma zapisany kontekst.` };
      },
    },
    {
      nazwa: 'Import XLIFF',
      wyjasnienie: `${Command.TranslateXliffImport} — jednostki wchodzą do panelu języka docelowego.`,
      pola: [sciezka.element],
      async wykonaj() {
        const brak = wymagajOkna() || wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wczytajXliff({
          windowId: idOkna(),
          path: sciezka.kontrolka.value.trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Import XLIFF', wynik) };
        }
        pokazWykaz(wynik.wynik.notes);
        return {
          udane: true,
          zdanie: `Wniesiono ${String(wynik.wynik.importedCount)} jednostek do ${String(
            wynik.wynik.panels.length,
          )} paneli.`,
        };
      },
    },
    {
      nazwa: 'Import napisów',
      wyjasnienie: `${Command.TranslateSubtitleImport} — SRT, WebVTT, TTML albo EBU STL.`,
      pola: [sciezka.element],
      async wykonaj() {
        const brak = wymagajOkna() || wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wczytajNapisy({
          windowId: idOkna(),
          path: sciezka.kontrolka.value.trim(),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Import napisów', wynik) };
        }
        pokazWykaz(
          wynik.wynik.cues.map(
            (kwestia) =>
              `${String(kwestia.startMs)}–${String(kwestia.endMs)} ms: ${kwestia.text}`,
          ),
        );
        return { udane: true, zdanie: `Wniesiono ${String(wynik.wynik.importedCount)} kwestii.` };
      },
    },
    {
      nazwa: 'Wydaj napisy',
      wyjasnienie: `${Command.TranslateSubtitleExport} — taktowanie źródła, treść przekładu.`,
      pola: [formatNapisow.element, sciezka.element],
      async wykonaj() {
        const brak = wymagajPanelu() || wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wydajNapisy({
          panelId: panel.kontrolka.value,
          path: sciezka.kontrolka.value.trim(),
          format: formatNapisow.kontrolka.value as never,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wydanie napisów', wynik) };
        }
        return {
          udane: true,
          zdanie: `Zapisano ${String(wynik.wynik.exportedCount)} kwestii pod ${wynik.wynik.path}.`,
        };
      },
    },
    {
      nazwa: 'Kontrola taktowania napisów',
      wyjasnienie: `${Command.TranslateSubtitleTimingCheck} — tempo, długość linii, nakładanie.`,
      pola: [],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.sprawdzTaktowanie({ panelId: panel.kontrolka.value });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Kontrola taktowania', wynik) };
        }
        pokazWykaz(
          wynik.wynik.issues.map(
            (zastrzezenie) =>
              `kwestia ${String(zastrzezenie.index)} · ${zastrzezenie.kind}: ` +
              `${String(zastrzezenie.value)} przy granicy ${String(zastrzezenie.limit)}`,
          ),
        );
        return { udane: true, zdanie: `Zastrzeżeń taktowania: ${String(wynik.wynik.issues.length)}.` };
      },
    },
    {
      nazwa: 'Scenariusz dubbingu',
      wyjasnienie: `${Command.TranslateDubbingScriptBuild} — kwestie wraz z długością wypowiedzenia.`,
      pola: [],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.zlozScenariuszDubbingu({ panelId: panel.kontrolka.value });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Scenariusz dubbingu', wynik) };
        }
        pokazWykaz(
          wynik.wynik.lines.map(
            (linia) => `${String(linia.index)} · ${String(linia.durationMs)} ms: ${linia.text}`,
          ),
        );
        return {
          udane: true,
          zdanie: `Kwestii: ${String(wynik.wynik.lines.length)}, ponad limit: ${String(
            wynik.wynik.overLimitCount,
          )}.`,
        };
      },
    },
    {
      nazwa: 'Profil silnika',
      wyjasnienie: `${Command.TranslateEngineProfileSet} — kanały modelu wraz z kolejnością wyboru.`,
      pola: [nazwaProfiluSilnika.element, kanalyProfilu.element],
      async wykonaj() {
        const wynik = await warsztat.ustawProfilSilnika({
          name: nazwaProfiluSilnika.kontrolka.value.trim(),
          channelIds: wykazZPola(kanalyProfilu),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Profil silnika', wynik) };
        }
        return {
          udane: true,
          zdanie: `Profil ${wynik.wynik.profile.id} z ${String(
            wynik.wynik.profile.channelIds.length,
          )} kanałami.`,
        };
      },
    },
    {
      nazwa: 'Wykaz profili silników',
      wyjasnienie: `${Command.TranslateEngineProfileList}`,
      pola: [],
      async wykonaj() {
        const wynik = await warsztat.wykazProfiliSilnikow({});
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wykaz profili silników', wynik) };
        }
        pokazWykaz(
          wynik.wynik.profiles.map(
            (profil) => `${profil.name} — kanały: ${profil.channelIds.join(', ')}`,
          ),
        );
        return { udane: true, zdanie: `Profili silników: ${String(wynik.wynik.profiles.length)}.` };
      },
    },
    {
      nazwa: 'Porównaj silniki',
      wyjasnienie: `${Command.TranslateEngineCompare} — ten sam segment kanałem po kanale.`,
      pola: [segmentPorownania.element, kanalyProfilu.element],
      async wykonaj() {
        const brak = wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.porownajSilniki({
          panelId: panel.kontrolka.value,
          segment: segmentPorownania.kontrolka.value,
          channelIds: wykazZPola(kanalyProfilu),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Porównanie silników', wynik) };
        }
        pokazWykaz(wynik.wynik.variants.map((wariant) => `${wariant.channelId}: ${wariant.text}`));
        return { udane: true, zdanie: `Wariantów: ${String(wynik.wynik.variants.length)}.` };
      },
    },
    {
      nazwa: 'Polityka pivota',
      wyjasnienie: `${Command.TranslatePivotPolicySet} — język pośredni dla par bez pary wprost.`,
      pola: [jezykPivota.element],
      async wykonaj() {
        const wynik = await warsztat.ustawPolitykePivota({
          defaultPivot: tekstZPola(jezykPivota),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Polityka pivota', wynik) };
        }
        return {
          udane: true,
          zdanie: `Polityka zasięgu ${wynik.wynik.policy.scope}; par: ${String(
            wynik.wynik.policy.pairs.length,
          )}.`,
        };
      },
    },
    {
      nazwa: 'Przebieg pakietowy',
      wyjasnienie: `${Command.TranslateBatchRun} — kontrola jakości wszystkich paneli okna.`,
      pola: [],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.uruchomPakiet({
          windowId: idOkna(),
          operations: ['qualityCheck'],
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Przebieg pakietowy', wynik) };
        }
        return {
          udane: true,
          zdanie: `Przebieg ${wynik.wynik.queueId}: pozycji ${String(wynik.wynik.itemCount)}.`,
        };
      },
    },
    {
      nazwa: 'Złóż pakiet przekazania',
      wyjasnienie: `${Command.TranslateHandoffBuild} — archiwum z XLIFF, TMX, TBX i instrukcjami.`,
      pola: [instrukcjeWykonawcy.element, sciezka.element],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.zlozPakietPrzekazania({
          windowId: idOkna(),
          contents: ['xliff', 'memory', 'termbase', 'instructions'],
          instructions: instrukcjeWykonawcy.kontrolka.value,
          path: tekstZPola(sciezka),
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Pakiet przekazania', wynik) };
        }
        return {
          udane: true,
          zdanie: `Pakiet ${wynik.wynik.package.id} leży pod ${wynik.wynik.path}.`,
        };
      },
    },
    {
      nazwa: 'Przyjmij zwrot wykonawcy',
      wyjasnienie: `${Command.TranslateHandoffReceive} — archiwum albo pojedynczy XLIFF.`,
      pola: [sciezka.element],
      async wykonaj() {
        const brak = wymagajOkna() || wymagajSciezki();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.przyjmijZwrot({
          windowId: idOkna(),
          path: sciezka.kontrolka.value.trim(),
          runQualityCheck: true,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Zwrot wykonawcy', wynik) };
        }
        return {
          udane: true,
          zdanie: `Przyjęto ${String(wynik.wynik.importedCount)} jednostek do ${String(
            wynik.wynik.panels.length,
          )} paneli.`,
        };
      },
    },
    {
      nazwa: 'Most: przyjmij źródło',
      wyjasnienie: `${Command.TranslateBridgeSourceReceive} — materiał z dokumentu innego modułu.`,
      pola: [idDokumentuMostu.element, trescMostu.element],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.przyjmijZrodloMostu({
          windowId: idOkna(),
          documentId: idDokumentuMostu.kontrolka.value.trim(),
          text: trescMostu.kontrolka.value,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Most — źródło', wynik) };
        }
        return {
          udane: true,
          zdanie: `Okno ${wynik.wynik.windowId} ma ${String(wynik.wynik.segmentCount)} segmentów.`,
        };
      },
    },
    {
      nazwa: 'Most: odeślij wynik',
      wyjasnienie: `${Command.TranslateBridgeResultSend} — przekład wraca do dokumentu źródłowego.`,
      pola: [trybOdeslania.element],
      async wykonaj() {
        const brak = wymagajOkna() || wymagajPanelu();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.odesljWynikMostu({
          windowId: idOkna(),
          panelIds: [panel.kontrolka.value],
          mode: trybOdeslania.kontrolka.value as never,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Most — wynik', wynik) };
        }
        return {
          udane: true,
          zdanie: `Do dokumentu ${wynik.wynik.documentId} poszło ${String(wynik.wynik.sentCount)} paneli.`,
        };
      },
    },
    {
      nazwa: 'Wydaj wytwór do biblioteki',
      wyjasnienie: `${Command.TranslateArtifactPublish} — plik w magazynie rdzenia wraz z wierszem pliku.`,
      pola: [rodzajWytworu.element],
      async wykonaj() {
        const brak = wymagajOkna();
        if (brak !== '') return { udane: false, zdanie: brak };
        const wynik = await warsztat.wydajWytwor({
          windowId: idOkna(),
          kind: rodzajWytworu.kontrolka.value as never,
        });
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wydanie wytworu', wynik) };
        }
        return {
          udane: true,
          zdanie: `Plik biblioteki ${wynik.wynik.fileId} leży pod ${wynik.wynik.path}.`,
        };
      },
    },
    {
      nazwa: 'Kroki pracy modułu',
      wyjasnienie: `${Command.TranslateStepList} — czynności rdzenia wraz z komendami i polami.`,
      pola: [],
      async wykonaj() {
        const wynik = await warsztat.wykazKrokow({});
        if (!wynik.udany || wynik.wynik === undefined) {
          return { udane: false, zdanie: zdanieOdmowy('Wykaz kroków', wynik) };
        }
        pokazWykaz(
          wynik.wynik.steps.map((krok) => `${krok.name} — ${krok.command} (${krok.parameters.join(', ')})`),
        );
        return { udane: true, zdanie: `Kroków pracy: ${String(wynik.wynik.steps.length)}.` };
      },
    },
  ];

  const wskazania = document.createElement('div');
  wskazania.className = 'mt-warsztat__wskazania';
  wskazania.append(panel.element, sciezka.element);

  okno.tresc.append(wskazania, odpowiedz.element, wykaz);
  for (const czynnosc of czynnosci) {
    okno.tresc.append(sekcjaCzynnosci(czynnosc, okno, odpowiedz));
  }

  const element = document.createElement('section');
  element.className = 'mt-okno mt-okno--zarzadca';
  element.dataset['okno'] = 'translation-workshop';
  element.append(naglowekOkna('Warsztat tłumaczenia', 'zarządca'), okno.element);

  function odswiez(): void {
    ustawPozycje(
      panel.kontrolka,
      stan
        .panelJezykow()
        .map((wpis) => ({ wartosc: wpis.id, etykieta: `${wpis.language} — panel ${wpis.id}` })),
    );
  }

  odswiez();
  return { element, odswiez };
}

/** Sekcja jednej czynności: rozwinięcie z polami i przyciskiem, zawsze klikalnym — brak wskazania kończy się zdaniem, czego brakuje. */
function sekcjaCzynnosci(
  czynnosc: Czynnosc,
  okno: StanOkna,
  odpowiedz: ReturnType<typeof utworzWierszOdpowiedzi>,
): HTMLElement {
  const rozwiniecie = utworzRozwiniecie({
    warstwa: 2,
    nazwa: czynnosc.nazwa,
    wyjasnienie: czynnosc.wyjasnienie,
    znacznik: '▼',
  });

  const guzik = przycisk(czynnosc.nazwa, 'dn-btn dn-btn--sm dn-btn--atrament');
  guzik.addEventListener('click', () => {
    void (async () => {
      okno.ladowanie(`${czynnosc.nazwa}: rdzeń pracuje.`);
      odpowiedz.pokaz(`${czynnosc.nazwa}…`, true);
      const wynik = await czynnosc.wykonaj();
      odpowiedz.pokaz(wynik.zdanie, wynik.udane);
      if (wynik.udane) {
        okno.gotowe();
        return;
      }
      okno.blad(wynik.zdanie);
    })();
  });

  const pasek = document.createElement('div');
  pasek.className = 'mt-pasek';
  pasek.append(guzik);

  rozwiniecie.tresc.append(...czynnosc.pola, pasek);
  return rozwiniecie.element;
}
