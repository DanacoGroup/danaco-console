import {
  StudioBreakKind,
  StudioHeaderScope,
  StudioPageNumberFormat,
  StudioPageOrientation,
  StudioSectionStart,
  StudioWatermarkKind,
  type StudioHeaderFooter,
  type StudioPageEnvelopeSetRequest,
  type StudioPageBreakInsertRequest,
  type StudioPageHeaderfooterSetRequest,
  type StudioPageNumberingSetRequest,
  type StudioPageSetup,
  type StudioPageSetupSetRequest,
  type StudioPageWatermarkSetRequest,
  type StudioPaperFormat,
  type StudioSection,
  type StudioSectionSaveRequest,
} from '../../../../shared/contract';
import {
  poleLogiczne,
  poleTekstowe,
  poleWyboru,
  przycisk,
  ustawPozycje,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { NASTAWY_MARGINESOW } from './nastawy-strony';
import { BEZ_ZMIANY, liczbaPola, poleLiczbowe, tekstPola, wyborPola } from './strona-pola-postaci';

/**
 * Treść żądania bez dokumentu — dokument dokłada warstwa wołająca rdzeń, ponieważ
 * panel składa wyłącznie pola formularza i nie zna identyfikatora dokumentu, na
 * którym pracuje Operator.
 */
type BezDokumentu<T> = Omit<T, 'documentId'>;

/**
 * Czynności panelu nastaw strony, które wywołujący przekłada na komendy rdzenia;
 * panel nie woła rdzenia sam, tylko składa treść żądania z pól wypełnionych przez
 * Operatora.
 */
export interface CzynnosciStronyPanelu {
  naNastawyStrony(zadanie: BezDokumentu<StudioPageSetupSetRequest>): void;
  naNumeracje(zadanie: BezDokumentu<StudioPageNumberingSetRequest>): void;
  naNaglowek(zadanie: BezDokumentu<StudioPageHeaderfooterSetRequest>): void;
  naZnakWodny(zadanie: BezDokumentu<StudioPageWatermarkSetRequest>): void;
  naKoperte(zadanie: BezDokumentu<StudioPageEnvelopeSetRequest>): void;
  naPodzial(zadanie: BezDokumentu<StudioPageBreakInsertRequest>): void;
  naZapisSekcji(zadanie: BezDokumentu<StudioSectionSaveRequest>): void;
  naUsuniecieSekcji(idSekcji: string): void;
  /** Ponowny odczyt nastaw, sekcji, nagłówków i wykazu nośników z rdzenia. */
  naOdczyt(): void;
}

/**
 * Panel nastaw strony wraz z jego sterowaniem: przełączaniem widoczności,
 * wpisywaniem nastaw odczytanych z rdzenia oraz wykazem sekcji, nagłówków
 * i nośników do wyboru.
 */
export interface StronaPanelNastaw {
  element: HTMLElement;
  przestawWidocznosc(): void;
  widoczny(): boolean;
  /** Wpisuje w pola nastawy oddane przez rdzeń. */
  pokazNastawy(nastawy: StudioPageSetup): void;
  /** Wykaz sekcji do wybrania — pozycja pusta znaczy „cały dokument". */
  pokazSekcje(sekcje: readonly StudioSection[]): void;
  /** Nagłówki i stopki wedle zasięgu; panel pokazuje ten o zasięgu wybranym. */
  pokazNaglowki(naglowki: readonly StudioHeaderFooter[]): void;
  /** Wykaz nośników z rdzenia wraz z kopertami. */
  pokazNosniki(nosniki: readonly StudioPaperFormat[]): void;
  pokazOdpowiedz(tresc: string, powodzenie: boolean): void;
  /** Sekcja wskazana przez Operatora; `undefined` znaczy cały dokument. */
  sekcjaWybrana(): string | undefined;
}

/**
 * Panel nastaw strony udostępnia numerację, marginesy, nagłówki, znak wodny,
 * kopertę i podziały jako osobne nastawy zapisywane komendami rdzenia, nie jako
 * pojedynczy przełącznik profilu wydania.
 */
export function utworzStronePanelNastaw(
  czynnosci: CzynnosciStronyPanelu,
  /** Miejsce kursora w treści — podział wstawia się tam, gdzie stoi kursor. */
  miejsceKursora: () => number,
): StronaPanelNastaw {
  const odpowiedz = utworzWierszOdpowiedzi();
  let naglowkiBiezace: readonly StudioHeaderFooter[] = [];

  /* ── Sekcja, której nastawy dotyczą ──────────────────────────────────────── */

  const sekcja = poleWyboru(
    {
      etykieta: 'Sekcja, której nastawy dotyczą',
      opis:
        'Wybór dotyczy wszystkich nastaw tego panelu. „Cały dokument" ustawia je dla dokumentu; ' +
        'sekcja ustawia je tylko jej — tak pismo z załącznikiem poziomym zostaje jednym dokumentem.',
    },
    [{ wartosc: '', etykieta: 'Cały dokument' }],
  );

  function sekcjaWybrana(): string | undefined {
    return wyborPola(sekcja.kontrolka);
  }

  /** Pola sekcji dokładane do każdego żądania; brak pola znaczy cały dokument. */
  function zSekcja(): { sectionId?: string } {
    const wybrana = sekcjaWybrana();
    return wybrana === undefined ? {} : { sectionId: wybrana };
  }

  /* ── Nośnik i orientacja ─────────────────────────────────────────────────── */

  const nosnik = poleWyboru(
    {
      etykieta: 'Format nośnika',
      opis:
        'Wykaz pochodzi z rdzenia (studio.page.paper.list) i niesie szereg A i B, Letter, Legal, ' +
        'Tabloid oraz koperty DL, C4, C5 i C6. Formatu własnego nie wybiera się z wykazu — podaje ' +
        'się go wymiarami w milimetrach poniżej.',
    },
    [BEZ_ZMIANY],
  );

  const szerokosc = poleLiczbowe(
    { etykieta: 'Format własny — szerokość w milimetrach', podpowiedz: 'np. 250' },
    { dolna: 1, gorna: 2000 },
  );
  const wysokosc = poleLiczbowe(
    { etykieta: 'Format własny — wysokość w milimetrach', podpowiedz: 'np. 350' },
    { dolna: 1, gorna: 2000 },
  );

  const orientacja = poleWyboru({ etykieta: 'Orientacja' }, [
    BEZ_ZMIANY,
    { wartosc: StudioPageOrientation.Pionowa, etykieta: 'Pionowa' },
    { wartosc: StudioPageOrientation.Pozioma, etykieta: 'Pozioma' },
  ]);

  /* ── Marginesy ───────────────────────────────────────────────────────────── */

  const nastawaMarginesow = poleWyboru(
    {
      etykieta: 'Nastawa gotowa marginesów',
      opis:
        'Wąskie, normalne, szerokie i do oprawy. Wybór wpisuje liczby w cztery pola poniżej — ' +
        'dalej wolno je poprawić, bo nastawa gotowa jest punktem wyjścia, nie klatką.',
    },
    [
      BEZ_ZMIANY,
      ...NASTAWY_MARGINESOW.map((nastawa) => ({ wartosc: nastawa.kod, etykieta: nastawa.nazwa })),
    ],
  );

  const marginesGora = poleLiczbowe(
    { etykieta: 'Margines górny w milimetrach' },
    { dolna: 0, gorna: 200, krok: 0.5 },
  );
  const marginesDol = poleLiczbowe(
    { etykieta: 'Margines dolny w milimetrach' },
    { dolna: 0, gorna: 200, krok: 0.5 },
  );
  const marginesLewy = poleLiczbowe(
    { etykieta: 'Margines lewy w milimetrach' },
    { dolna: 0, gorna: 200, krok: 0.5 },
  );
  const marginesPrawy = poleLiczbowe(
    { etykieta: 'Margines prawy w milimetrach' },
    { dolna: 0, gorna: 200, krok: 0.5 },
  );
  const oprawa = poleLiczbowe(
    {
      etykieta: 'Margines na oprawę w milimetrach',
      opis: 'Pas doliczany do marginesu wewnętrznego — pod zszycie albo klejenie.',
    },
    { dolna: 0, gorna: 100, krok: 0.5 },
  );
  const odbicia = poleLogiczne({
    etykieta: 'Marginesy odbicia dla druku dwustronnego',
    opis:
      'Margines lewy staje się WEWNĘTRZNYM: na stronie nieparzystej stoi po lewej, na parzystej ' +
      'po prawej. Bez tego oprawa wypadałaby raz w rowku, raz na krawędzi.',
  });

  nastawaMarginesow.kontrolka.addEventListener('change', () => {
    const kod = wyborPola(nastawaMarginesow.kontrolka);
    if (kod === undefined) return;
    const nastawa = NASTAWY_MARGINESOW.find((pozycja) => pozycja.kod === kod);
    if (nastawa === undefined) return;
    marginesGora.kontrolka.value = String(nastawa.goraMm);
    marginesDol.kontrolka.value = String(nastawa.dolMm);
    marginesLewy.kontrolka.value = String(nastawa.lewyMm);
    marginesPrawy.kontrolka.value = String(nastawa.prawyMm);
    oprawa.kontrolka.value = String(nastawa.oprawaMm);
    if (nastawa.oprawaMm > 0) odbicia.kontrolka.checked = true;
    odpowiedz.pokaz(
      `Nastawa „${nastawa.nazwa}" wpisana w pola. ${nastawa.opis} Naciśnij „Zapisz nastawy strony", ` +
        'żeby dojechała do rdzenia — sama zmiana pól niczego jeszcze nie zapisuje.',
      true,
    );
  });

  /* ── Kolumny ─────────────────────────────────────────────────────────────── */

  const kolumny = poleLiczbowe({ etykieta: 'Liczba kolumn' }, { dolna: 1, gorna: 12 });
  const odstepKolumn = poleLiczbowe(
    { etykieta: 'Odstęp między kolumnami w milimetrach' },
    { dolna: 0, gorna: 100, krok: 0.5 },
  );
  const liniaKolumn = poleLogiczne({ etykieta: 'Linia między kolumnami' });

  const zapiszStrone = przycisk('Zapisz nastawy strony', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszStrone.dataset['czynnosc'] = 'zapisz-nastawy-strony';
  zapiszStrone.title =
    'Idzie komendą studio.page.setup.set. Pole puste znaczy „nie ruszaj tej cechy", nie zero: ' +
    'margines zerowy jest nastawą, którą Operator wybiera świadomie.';
  zapiszStrone.addEventListener('click', () => {
    const zadanie: BezDokumentu<StudioPageSetupSetRequest> = { ...zSekcja() };
    const wybranyNosnik = wyborPola(nosnik.kontrolka);
    if (wybranyNosnik !== undefined) zadanie.paperName = wybranyNosnik;
    const szerokoscMm = liczbaPola(szerokosc.kontrolka);
    const wysokoscMm = liczbaPola(wysokosc.kontrolka);
    // Format własny podaje się dwiema liczbami naraz — sama szerokość albo
    // wysokość nie jest nośnikiem.
    if ((szerokoscMm === undefined) !== (wysokoscMm === undefined)) {
      odpowiedz.pokaz(
        'Format własny podaje się DWIEMA liczbami — szerokością i wysokością. Jedna z nich bez ' +
          'drugiej nie jest nośnikiem, więc żądanie nie pojechało.',
        false,
      );
      return;
    }
    if (szerokoscMm !== undefined) zadanie.widthMm = szerokoscMm;
    if (wysokoscMm !== undefined) zadanie.heightMm = wysokoscMm;
    const wybranaOrientacja = wyborPola(orientacja.kontrolka);
    if (wybranaOrientacja !== undefined) {
      zadanie.orientation = wybranaOrientacja as StudioPageOrientation;
    }
    const gora = liczbaPola(marginesGora.kontrolka);
    if (gora !== undefined) zadanie.marginTopMm = gora;
    const dol = liczbaPola(marginesDol.kontrolka);
    if (dol !== undefined) zadanie.marginBottomMm = dol;
    const lewy = liczbaPola(marginesLewy.kontrolka);
    if (lewy !== undefined) zadanie.marginLeftMm = lewy;
    const prawy = liczbaPola(marginesPrawy.kontrolka);
    if (prawy !== undefined) zadanie.marginRightMm = prawy;
    const oprawaMm = liczbaPola(oprawa.kontrolka);
    if (oprawaMm !== undefined) zadanie.gutterMm = oprawaMm;
    zadanie.mirrorMargins = odbicia.kontrolka.checked;
    const kodNastawy = wyborPola(nastawaMarginesow.kontrolka);
    if (kodNastawy !== undefined) zadanie.marginPreset = kodNastawy;
    const ileKolumn = liczbaPola(kolumny.kontrolka);
    if (ileKolumn !== undefined) zadanie.columns = ileKolumn;
    const odstep = liczbaPola(odstepKolumn.kontrolka);
    if (odstep !== undefined) zadanie.columnGapMm = odstep;
    zadanie.columnRule = liniaKolumn.kontrolka.checked;
    czynnosci.naNastawyStrony(zadanie);
  });

  /* ── Numeracja stron ─────────────────────────────────────────────────────── */

  const numeracjaCzynna = poleLogiczne({
    etykieta: 'Numeracja stron włączona',
    opis: 'Wyłączenie zdejmuje numery ze stron; pozostałe nastawy numeracji zostają zapisane.',
  });

  const formatNumeru = poleWyboru(
    {
      etykieta: 'Styl numeru',
      opis: 'Cyfry arabskie, rzymskie wielkie i małe, litery wielkie i małe.',
    },
    [
      BEZ_ZMIANY,
      { wartosc: StudioPageNumberFormat.Arabic, etykieta: 'Cyfry arabskie — 1, 2, 3' },
      { wartosc: StudioPageNumberFormat.RomanUpper, etykieta: 'Cyfry rzymskie wielkie — I, II, III' },
      { wartosc: StudioPageNumberFormat.RomanLower, etykieta: 'Cyfry rzymskie małe — i, ii, iii' },
      { wartosc: StudioPageNumberFormat.LetterUpper, etykieta: 'Litery wielkie — A, B, C' },
      { wartosc: StudioPageNumberFormat.LetterLower, etykieta: 'Litery małe — a, b, c' },
    ],
  );

  const numerStartowy = poleLiczbowe(
    {
      etykieta: 'Numeracja od numeru',
      opis:
        'Punkt startu numeracji. Pismo z kartą tytułową liczoną osobno zaczyna się często od 2 ' +
        'albo od numeru dalszego, gdy stanowi część większej całości.',
    },
    { dolna: 0, gorna: 100000 },
  );

  const wznowienie = poleLogiczne({
    etykieta: 'Numeracja wznawia się w tej sekcji',
    opis:
      'Załącznik numerowany od nowa dostaje wznowienie; ciągłość z pismem głównym zostawia je ' +
      'wyłączone. Nastawa dotyczy sekcji wybranej u góry panelu.',
  });

  const liczbaStron = poleLogiczne({
    etykieta: 'Dopisz liczbę stron — „strona N z M"',
  });

  const polozenieNumeru = poleWyboru(
    {
      etykieta: 'Umiejscowienie numeru',
      opis: 'Nagłówek albo stopka, wraz z wyrównaniem do lewej, środka albo prawej.',
    },
    [
      BEZ_ZMIANY,
      { wartosc: 'naglowek-lewo', etykieta: 'Nagłówek, do lewej' },
      { wartosc: 'naglowek-srodek', etykieta: 'Nagłówek, na środku' },
      { wartosc: 'naglowek-prawo', etykieta: 'Nagłówek, do prawej' },
      { wartosc: 'stopka-lewo', etykieta: 'Stopka, do lewej' },
      { wartosc: 'stopka-srodek', etykieta: 'Stopka, na środku' },
      { wartosc: 'stopka-prawo', etykieta: 'Stopka, do prawej' },
    ],
  );

  const zapiszNumeracje = przycisk('Zapisz numerację stron', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszNumeracje.dataset['czynnosc'] = 'zapisz-numeracje';
  zapiszNumeracje.title = 'Idzie komendą studio.page.numbering.set.';
  zapiszNumeracje.addEventListener('click', () => {
    const zadanie: BezDokumentu<StudioPageNumberingSetRequest> = {
      ...zSekcja(),
      enabled: numeracjaCzynna.kontrolka.checked,
      restartInSection: wznowienie.kontrolka.checked,
      showTotal: liczbaStron.kontrolka.checked,
    };
    const format = wyborPola(formatNumeru.kontrolka);
    if (format !== undefined) zadanie.format = format as StudioPageNumberFormat;
    const start = liczbaPola(numerStartowy.kontrolka);
    if (start !== undefined) zadanie.startAt = start;
    const polozenie = wyborPola(polozenieNumeru.kontrolka);
    if (polozenie !== undefined) zadanie.position = polozenie;
    czynnosci.naNumeracje(zadanie);
  });

  /* ── Nagłówek i stopka ───────────────────────────────────────────────────── */

  const zasieg = poleWyboru(
    {
      etykieta: 'Zasięg nagłówka i stopki',
      opis:
        'Strony zwykłe, pierwsza strona osobno i strony parzyste osobno — trzy zasięgi, każdy ' +
        'z własną treścią. Zmiana zasięgu wpisuje w pola treść już zapisaną dla niego.',
    },
    [
      { wartosc: StudioHeaderScope.Default, etykieta: 'Strony zwykłe sekcji' },
      { wartosc: StudioHeaderScope.FirstPage, etykieta: 'Pierwsza strona sekcji' },
      { wartosc: StudioHeaderScope.EvenPages, etykieta: 'Strony parzyste' },
    ],
  );

  const trescNaglowka = poleTekstowe({ etykieta: 'Treść nagłówka' });
  const trescStopki = poleTekstowe({ etykieta: 'Treść stopki' });
  const odlegloscNaglowka = poleLiczbowe(
    { etykieta: 'Odległość nagłówka od krawędzi w milimetrach' },
    { dolna: 0, gorna: 100, krok: 0.5 },
  );
  const odlegloscStopki = poleLiczbowe(
    { etykieta: 'Odległość stopki od krawędzi w milimetrach' },
    { dolna: 0, gorna: 100, krok: 0.5 },
  );
  const przejecieNaglowka = poleLogiczne({
    etykieta: 'Sekcja przejmuje nagłówek sekcji poprzedniej',
  });

  function wpiszNaglowek(): void {
    const wybrany = naglowkiBiezace.find((pozycja) => pozycja.scope === zasieg.kontrolka.value);
    trescNaglowka.kontrolka.value = wybrany?.headerText ?? '';
    trescStopki.kontrolka.value = wybrany?.footerText ?? '';
    odlegloscNaglowka.kontrolka.value =
      wybrany?.headerDistanceMm === undefined ? '' : String(wybrany.headerDistanceMm);
    odlegloscStopki.kontrolka.value =
      wybrany?.footerDistanceMm === undefined ? '' : String(wybrany.footerDistanceMm);
    przejecieNaglowka.kontrolka.checked = wybrany?.linkedToPrevious ?? false;
  }
  zasieg.kontrolka.addEventListener('change', () => wpiszNaglowek());

  const zapiszNaglowek = przycisk('Zapisz nagłówek i stopkę', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszNaglowek.dataset['czynnosc'] = 'zapisz-naglowek';
  zapiszNaglowek.title = 'Idzie komendą studio.page.headerfooter.set.';
  zapiszNaglowek.addEventListener('click', () => {
    const zadanie: BezDokumentu<StudioPageHeaderfooterSetRequest> = {
      ...zSekcja(),
      scope: zasieg.kontrolka.value as StudioHeaderScope,
      linkedToPrevious: przejecieNaglowka.kontrolka.checked,
    };
    // Treść pusta jedzie jawnie: opróżnienie nagłówka jest świadomą czynnością,
    // nie brakiem wskazania.
    zadanie.headerText = trescNaglowka.kontrolka.value;
    zadanie.footerText = trescStopki.kontrolka.value;
    const odlegloscG = liczbaPola(odlegloscNaglowka.kontrolka);
    if (odlegloscG !== undefined) zadanie.headerDistanceMm = odlegloscG;
    const odlegloscD = liczbaPola(odlegloscStopki.kontrolka);
    if (odlegloscD !== undefined) zadanie.footerDistanceMm = odlegloscD;
    czynnosci.naNaglowek(zadanie);
  });

  /* ── Znak wodny ──────────────────────────────────────────────────────────── */

  const rodzajZnaku = poleWyboru({ etykieta: 'Znak wodny' }, [
    { wartosc: StudioWatermarkKind.None, etykieta: 'Bez znaku wodnego' },
    { wartosc: StudioWatermarkKind.Text, etykieta: 'Napis' },
    { wartosc: StudioWatermarkKind.Image, etykieta: 'Obraz z magazynu zasobów' },
  ]);
  const napisZnaku = poleTekstowe({
    etykieta: 'Napis znaku wodnego',
    podpowiedz: 'np. PROJEKT',
  });
  const zasobZnaku = poleTekstowe({
    etykieta: 'Zasób obrazu znaku wodnego',
    opis: 'Identyfikator zasobu magazynu rdzenia.',
  });
  const krycieZnaku = poleLiczbowe(
    { etykieta: 'Krycie od 0 do 1' },
    { dolna: 0, gorna: 1, krok: 0.05 },
  );
  const obrotZnaku = poleLiczbowe(
    { etykieta: 'Obrót w stopniach' },
    { dolna: -180, gorna: 180 },
  );
  const barwaZnaku = poleTekstowe({ etykieta: 'Barwa napisu', podpowiedz: '#c0c0c0' });
  const stopienZnaku = poleLiczbowe(
    { etykieta: 'Stopień pisma napisu w punktach' },
    { dolna: 1, gorna: 400 },
  );

  const zapiszZnak = przycisk('Zapisz znak wodny', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszZnak.dataset['czynnosc'] = 'zapisz-znak-wodny';
  zapiszZnak.title = 'Idzie komendą studio.page.watermark.set.';
  zapiszZnak.addEventListener('click', () => {
    const zadanie: BezDokumentu<StudioPageWatermarkSetRequest> = {
      ...zSekcja(),
      kind: rodzajZnaku.kontrolka.value as StudioWatermarkKind,
    };
    const napis = tekstPola(napisZnaku.kontrolka);
    if (napis !== undefined) zadanie.text = napis;
    const zasob = tekstPola(zasobZnaku.kontrolka);
    if (zasob !== undefined) zadanie.assetId = zasob;
    const krycie = liczbaPola(krycieZnaku.kontrolka);
    if (krycie !== undefined) zadanie.opacity = krycie;
    const obrot = liczbaPola(obrotZnaku.kontrolka);
    if (obrot !== undefined) zadanie.angleDeg = obrot;
    const barwa = tekstPola(barwaZnaku.kontrolka);
    if (barwa !== undefined) zadanie.color = barwa;
    const stopien = liczbaPola(stopienZnaku.kontrolka);
    if (stopien !== undefined) zadanie.fontSizePt = stopien;
    czynnosci.naZnakWodny(zadanie);
  });

  /* ── Nadruk koperty ──────────────────────────────────────────────────────── */

  const adresat = poleTekstowe({ etykieta: 'Adres adresata' });
  const nadawca = poleTekstowe({ etykieta: 'Adres nadawcy' });
  const zNadawca = poleLogiczne({ etykieta: 'Nadrukuj nadawcę' });
  const adresatX = poleLiczbowe(
    { etykieta: 'Adresat — od lewej krawędzi w milimetrach' },
    { dolna: 0, gorna: 500, krok: 0.5 },
  );
  const adresatY = poleLiczbowe(
    { etykieta: 'Adresat — od górnej krawędzi w milimetrach' },
    { dolna: 0, gorna: 500, krok: 0.5 },
  );
  const nadawcaX = poleLiczbowe(
    { etykieta: 'Nadawca — od lewej krawędzi w milimetrach' },
    { dolna: 0, gorna: 500, krok: 0.5 },
  );
  const nadawcaY = poleLiczbowe(
    { etykieta: 'Nadawca — od górnej krawędzi w milimetrach' },
    { dolna: 0, gorna: 500, krok: 0.5 },
  );

  const zapiszKoperte = przycisk('Zapisz nadruk koperty', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszKoperte.dataset['czynnosc'] = 'zapisz-koperte';
  zapiszKoperte.title =
    'Idzie komendą studio.page.envelope.set. Koperta bez adresów i ich położenia byłaby samym ' +
    'rozmiarem nośnika, nie funkcją nadruku.';
  zapiszKoperte.addEventListener('click', () => {
    const zadanie: BezDokumentu<StudioPageEnvelopeSetRequest> = {
      ...zSekcja(),
      includeSender: zNadawca.kontrolka.checked,
    };
    const doKogo = tekstPola(adresat.kontrolka);
    if (doKogo !== undefined) zadanie.recipient = doKogo;
    const odKogo = tekstPola(nadawca.kontrolka);
    if (odKogo !== undefined) zadanie.sender = odKogo;
    const ax = liczbaPola(adresatX.kontrolka);
    if (ax !== undefined) zadanie.recipientXMm = ax;
    const ay = liczbaPola(adresatY.kontrolka);
    if (ay !== undefined) zadanie.recipientYMm = ay;
    const nx = liczbaPola(nadawcaX.kontrolka);
    if (nx !== undefined) zadanie.senderXMm = nx;
    const ny = liczbaPola(nadawcaY.kontrolka);
    if (ny !== undefined) zadanie.senderYMm = ny;
    czynnosci.naKoperte(zadanie);
  });

  /* ── Podział ─────────────────────────────────────────────────────────────── */

  const rodzajPodzialu = poleWyboru({ etykieta: 'Rodzaj podziału' }, [
    { wartosc: StudioBreakKind.Page, etykieta: 'Podział strony' },
    { wartosc: StudioBreakKind.Column, etykieta: 'Podział kolumny' },
    { wartosc: StudioBreakKind.Section, etykieta: 'Podział sekcji' },
    { wartosc: StudioBreakKind.Line, etykieta: 'Podział wiersza bez nowego akapitu' },
  ]);

  const poczatekSekcji = poleWyboru(
    {
      etykieta: 'Sposób rozpoczęcia sekcji',
      opis: 'Dotyczy wyłącznie podziału sekcji.',
    },
    [
      BEZ_ZMIANY,
      { wartosc: StudioSectionStart.Continuous, etykieta: 'Dalej na tej stronie' },
      { wartosc: StudioSectionStart.NewPage, etykieta: 'Od nowej strony' },
      { wartosc: StudioSectionStart.EvenPage, etykieta: 'Od strony parzystej' },
      { wartosc: StudioSectionStart.OddPage, etykieta: 'Od strony nieparzystej' },
      { wartosc: StudioSectionStart.NewColumn, etykieta: 'Od nowej kolumny' },
    ],
  );

  const wstawPodzial = przycisk('Wstaw podział w miejscu kursora', 'dn-btn dn-btn--sm dn-btn--duch');
  wstawPodzial.dataset['czynnosc'] = 'wstaw-podzial';
  wstawPodzial.title = 'Idzie komendą studio.page.break.insert.';
  wstawPodzial.addEventListener('click', () => {
    const zadanie: BezDokumentu<StudioPageBreakInsertRequest> = {
      offset: miejsceKursora(),
      kind: rodzajPodzialu.kontrolka.value as StudioBreakKind,
    };
    const start = wyborPola(poczatekSekcji.kontrolka);
    if (start !== undefined) zadanie.sectionStart = start as StudioSectionStart;
    czynnosci.naPodzial(zadanie);
  });

  /* ── Sekcje: założenie, zmiana, usunięcie ────────────────────────────────── */

  const nazwaSekcji = poleTekstowe({
    etykieta: 'Nazwa sekcji',
    opis: 'Sekcja bez nazwy nie da się później wskazać w wykazie.',
  });
  const sekcjaOd = poleLiczbowe({ etykieta: 'Początek sekcji w znakach' }, { dolna: 0 });
  const sekcjaDo = poleLiczbowe({ etykieta: 'Koniec sekcji w znakach' }, { dolna: 0 });

  const zalozSekcje = przycisk('Zapisz sekcję', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zalozSekcje.dataset['czynnosc'] = 'zapisz-sekcje';
  zalozSekcje.title =
    'Idzie komendą studio.section.save. Sekcja wybrana u góry panelu jest sekcją ZMIENIANĄ; ' +
    'przy „Cały dokument" zakłada się sekcję nową.';
  zalozSekcje.addEventListener('click', () => {
    const zadanie: BezDokumentu<StudioSectionSaveRequest> = { ...zSekcja() };
    const nazwa = tekstPola(nazwaSekcji.kontrolka);
    if (nazwa !== undefined) zadanie.title = nazwa;
    const od = liczbaPola(sekcjaOd.kontrolka);
    if (od !== undefined) zadanie.rangeStart = od;
    const doZnaku = liczbaPola(sekcjaDo.kontrolka);
    if (doZnaku !== undefined) zadanie.rangeEnd = doZnaku;
    const start = wyborPola(poczatekSekcji.kontrolka);
    if (start !== undefined) zadanie.start = start as StudioSectionStart;
    czynnosci.naZapisSekcji(zadanie);
  });

  const usunSekcje = przycisk('Usuń wybraną sekcję', 'dn-btn dn-btn--sm dn-btn--ostrzezenie');
  usunSekcje.dataset['czynnosc'] = 'usun-sekcje';
  usunSekcje.title =
    'Idzie komendą studio.section.delete. Treść sekcji NIE ginie — przechodzi do sekcji ' +
    'poprzedniej wraz z jej nastawami.';
  usunSekcje.addEventListener('click', () => {
    const wybrana = sekcjaWybrana();
    if (wybrana === undefined) {
      odpowiedz.pokaz(
        'Usunięcie dotyczy sekcji, nie dokumentu — wskaż sekcję u góry panelu. „Cały dokument" ' +
          'nie jest sekcją i nie ma czego usunąć.',
        false,
      );
      return;
    }
    czynnosci.naUsuniecieSekcji(wybrana);
  });

  const odczytaj = przycisk('Odczytaj nastawy z rdzenia', 'dn-btn dn-btn--sm dn-btn--zarys');
  odczytaj.dataset['czynnosc'] = 'odczytaj-nastawy';
  odczytaj.title =
    'Idzie komendami studio.page.setup.get, studio.section.list, studio.page.headerfooter.get ' +
    'i studio.page.paper.list — pola pokazują to, co naprawdę stoi w dokumencie.';
  odczytaj.addEventListener('click', () => czynnosci.naOdczyt());

  /* ── Układ panelu ────────────────────────────────────────────────────────── */

  const element = document.createElement('section');
  element.className = 'ms-postac ms-postac--strona';
  element.dataset['panel'] = 'nastawy-strony';
  element.hidden = true;
  element.setAttribute('aria-label', 'Nastawy strony — nośnik, marginesy, numeracja, sekcje');

  element.append(
    grupa('Sekcja i odczyt', [sekcja.element, odczytaj]),
    grupa('Nośnik i orientacja', [
      nosnik.element,
      szerokosc.element,
      wysokosc.element,
      orientacja.element,
    ]),
    grupa('Marginesy', [
      nastawaMarginesow.element,
      marginesGora.element,
      marginesDol.element,
      marginesLewy.element,
      marginesPrawy.element,
      oprawa.element,
      odbicia.element,
    ]),
    grupa('Kolumny', [kolumny.element, odstepKolumn.element, liniaKolumn.element, zapiszStrone]),
    grupa('Numeracja stron', [
      numeracjaCzynna.element,
      formatNumeru.element,
      numerStartowy.element,
      wznowienie.element,
      liczbaStron.element,
      polozenieNumeru.element,
      zapiszNumeracje,
    ]),
    grupa('Nagłówek i stopka', [
      zasieg.element,
      trescNaglowka.element,
      trescStopki.element,
      odlegloscNaglowka.element,
      odlegloscStopki.element,
      przejecieNaglowka.element,
      zapiszNaglowek,
    ]),
    grupa('Znak wodny', [
      rodzajZnaku.element,
      napisZnaku.element,
      zasobZnaku.element,
      krycieZnaku.element,
      obrotZnaku.element,
      barwaZnaku.element,
      stopienZnaku.element,
      zapiszZnak,
    ]),
    grupa('Nadruk koperty', [
      adresat.element,
      nadawca.element,
      zNadawca.element,
      adresatX.element,
      adresatY.element,
      nadawcaX.element,
      nadawcaY.element,
      zapiszKoperte,
    ]),
    grupa('Podziały i sekcje', [
      rodzajPodzialu.element,
      poczatekSekcji.element,
      wstawPodzial,
      nazwaSekcji.element,
      sekcjaOd.element,
      sekcjaDo.element,
      zalozSekcje,
      usunSekcje,
    ]),
    odpowiedz.element,
  );

  return {
    element,

    przestawWidocznosc() {
      element.hidden = !element.hidden;
    },

    widoczny: () => !element.hidden,

    pokazNastawy(nastawy) {
      if (nastawy.pageSize !== undefined) nosnik.kontrolka.value = nastawy.pageSize;
      if (nastawy.orientation !== undefined) orientacja.kontrolka.value = nastawy.orientation;
      wpiszLiczbe(marginesGora.kontrolka, nastawy.marginTop);
      wpiszLiczbe(marginesDol.kontrolka, nastawy.marginBottom);
      wpiszLiczbe(marginesLewy.kontrolka, nastawy.marginLeft);
      wpiszLiczbe(marginesPrawy.kontrolka, nastawy.marginRight);
      wpiszLiczbe(szerokosc.kontrolka, nastawy.widthMm);
      wpiszLiczbe(wysokosc.kontrolka, nastawy.heightMm);
      // Nagłówek i stopka wchodzą wyłącznie bez wpisu zasięgu — źródło bogatsze
      // nie ustępuje uboższemu.
      if (naglowkiBiezace.length === 0) {
        trescNaglowka.kontrolka.value = nastawy.header ?? '';
        trescStopki.kontrolka.value = nastawy.footer ?? '';
      }
      if (nastawy.pageNumbers !== undefined) {
        numeracjaCzynna.kontrolka.checked = nastawy.pageNumbers;
      }
    },

    pokazSekcje(sekcje) {
      ustawPozycje(sekcja.kontrolka, [
        { wartosc: '', etykieta: 'Cały dokument' },
        ...sekcje.map((pozycja) => ({
          wartosc: pozycja.id,
          etykieta:
            `${pozycja.index + 1}. ${pozycja.title ?? 'sekcja bez nazwy'} ` +
            `(znaki ${pozycja.rangeStart}–${pozycja.rangeEnd})`,
        })),
      ]);
    },

    pokazNaglowki(naglowki) {
      naglowkiBiezace = naglowki;
      wpiszNaglowek();
    },

    pokazNosniki(nosniki) {
      ustawPozycje(nosnik.kontrolka, [
        BEZ_ZMIANY,
        ...nosniki.map((pozycja) => ({
          wartosc: pozycja.name,
          etykieta: `${pozycja.name} — ${pozycja.widthMm}×${pozycja.heightMm} mm`,
        })),
      ]);
    },

    pokazOdpowiedz: (tresc, powodzenie) => odpowiedz.pokaz(tresc, powodzenie),
    sekcjaWybrana,
  };
}

/**
 * Wpisuje liczbę w pole tekstowe; wartość nieokreślona zostawia pole puste
 * zamiast zera, ponieważ pole puste i margines zerowy niosą odrębne znaczenie
 * w żądaniu do rdzenia.
 */
function wpiszLiczbe(kontrolka: HTMLInputElement, wartosc: number | undefined): void {
  kontrolka.value = wartosc === undefined ? '' : String(wartosc);
}

/**
 * Buduje grupę pól panelu złożoną z nagłówka i przekazanej zawartości,
 * zachowując ten sam układ wizualny, jaki stosują grupy przycisków wstążki
 * narzędziowej.
 */
function grupa(tytul: string, zawartosc: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'ms-postac__tytul';
  naglowek.textContent = tytul;

  const element = document.createElement('div');
  element.className = 'ms-postac__grupa';
  element.append(naglowek, ...zawartosc);
  return element;
}
