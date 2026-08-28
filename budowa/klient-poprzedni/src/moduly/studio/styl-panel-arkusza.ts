import {
  StudioCaseTransform,
  StudioLineSpacingRule,
  StudioStyleKind,
  StudioTextAlign,
  StudioTextEffect,
  StudioUnderlineStyle,
  type StudioCharacterFormat,
  type StudioFormatCaseSetRequest,
  type StudioFormatCharacterSetRequest,
  type StudioFormatClearRequest,
  type StudioFormatPainterApplyRequest,
  type StudioFormatPainterCopyRequest,
  type StudioFormatParagraphSetRequest,
  type StudioFormatReplaceRequest,
  type StudioFormatSimilarSelectRequest,
  type StudioNamedStyle,
  type StudioParagraphFormat,
  type StudioRulerTabstopSetRequest,
  type StudioStyleApplyRequest,
  type StudioStyleDeleteRequest,
  type StudioStyleSaveRequest,
  StudioTabKind,
  StudioTabLeader,
} from '../../../../shared/contract';
import {
  poleLogiczne,
  poleTekstowe,
  poleWyboru,
  przycisk,
  utworzWierszOdpowiedzi,
} from '../../modele/kontrolki-formularza';
import { BEZ_ZMIANY, liczbaPola, poleLiczbowe, tekstPola, wyborPola } from './strona-pola-postaci';

/**
 * Panel arkusza stylów i formatowania fragmentu: style, znak, akapit, malarz formatów i zamiana.
 */

/**
 * Treść żądania bez dokumentu i bez zakresu; oba pola dokłada warstwa wyżej, znająca zaznaczenie
 * i dokument bieżący.
 */
type BezZakresu<T> = Omit<T, 'documentId' | 'rangeStart' | 'rangeEnd'>;

/**
 * Czynności panelu arkusza stylów: zapis, stosowanie i usunięcie stylu, formatowanie znaku
 * i akapitu, malarz formatów, zamiana i tabulator.
 */
export interface CzynnosciStyluPanelu {
  naZapisStylu(zadanie: Omit<StudioStyleSaveRequest, 'documentId'>): void;
  naStosowanieStylu(zadanie: BezZakresu<StudioStyleApplyRequest>): void;
  naUsuniecieStylu(zadanie: Omit<StudioStyleDeleteRequest, 'documentId'>): void;
  naStylZnaku(zadanie: BezZakresu<StudioFormatCharacterSetRequest>): void;
  naStylAkapitu(zadanie: BezZakresu<StudioFormatParagraphSetRequest>): void;
  naCzyszczenie(zadanie: BezZakresu<StudioFormatClearRequest>): void;
  naWielkoscLiter(zadanie: BezZakresu<StudioFormatCaseSetRequest>): void;
  naPobraniePostaci(zadanie: BezZakresu<StudioFormatPainterCopyRequest>): void;
  naNalozeniePostaci(zadanie: BezZakresu<Omit<StudioFormatPainterApplyRequest, 'clipId'>>): void;
  naPodobne(zadanie: BezZakresu<StudioFormatSimilarSelectRequest>): void;
  naZamiane(zadanie: BezZakresu<StudioFormatReplaceRequest>): void;
  naTabulator(zadanie: BezZakresu<StudioRulerTabstopSetRequest>): void;
  /** Ponowny odczyt arkusza stylów oraz postaci znaku i akapitu zaznaczenia. */
  naOdczyt(): void;
}

/**
 * Panel arkusza stylów wraz z jego sterowaniem: pokazywanie stylów, postaci znaku i akapitu,
 * stanu malarza formatów oraz odpowiedzi rdzenia.
 */
export interface StylPanelArkusza {
  element: HTMLElement;
  przestawWidocznosc(): void;
  widoczny(): boolean;
  /** Arkusz stylów oddany przez rdzeń wraz z liczbą miejsc użycia. */
  pokazStyle(style: readonly StudioNamedStyle[]): void;
  /** Postać znaku zaznaczenia wraz z tym, co we fragmencie niejednolite. */
  pokazPostacZnaku(postac: StudioCharacterFormat, niejednolite: readonly string[]): void;
  /** Postać akapitu zaznaczenia. */
  pokazPostacAkapitu(postac: StudioParagraphFormat, niejednolite: readonly string[]): void;
  /** Czy malarz formatów ma coś pobrane; zdanie widoczne przy jego przyciskach. */
  pokazMalarza(pobrane: boolean, opis: string): void;
  pokazOdpowiedz(tresc: string, powodzenie: boolean): void;
  // Czy zamiana ma objąć wyłącznie zaznaczenie; wybór należy do Operatora, panel go nie zgaduje.
  zamianaTylkoWZaznaczeniu(): boolean;
}

export function utworzStylPanelArkusza(czynnosci: CzynnosciStyluPanelu): StylPanelArkusza {
  const odpowiedz = utworzWierszOdpowiedzi();

  /* ── Arkusz stylów ───────────────────────────────────────────────────────── */

  const wykazStylow = document.createElement('ul');
  wykazStylow.className = 'ms-postac__wykaz';

  const nazwaStylu = poleTekstowe({
    etykieta: 'Nazwa stylu',
    opis: 'Nazwa jest jedynym identyfikatorem stylu w dokumencie — nazwa zastana znaczy zmianę.',
  });
  const nazwaWidoczna = poleTekstowe({ etykieta: 'Nazwa widoczna dla Operatora' });
  const rodzajStylu = poleWyboru({ etykieta: 'Rodzaj stylu' }, [
    { wartosc: StudioStyleKind.Paragraph, etykieta: 'Styl akapitu' },
    { wartosc: StudioStyleKind.Character, etykieta: 'Styl znaku' },
    { wartosc: StudioStyleKind.Table, etykieta: 'Styl tabeli' },
    { wartosc: StudioStyleKind.List, etykieta: 'Styl listy' },
  ]);
  const stylNadrzedny = poleTekstowe({
    etykieta: 'Styl nadrzędny, po którym ten dziedziczy',
  });
  const stylNastepny = poleTekstowe({ etykieta: 'Styl akapitu następnego' });
  const stylZamienny = poleTekstowe({
    etykieta: 'Styl przejmujący miejsca użycia przy usunięciu',
  });

  const zapiszStyl = przycisk('Zapisz styl nazwany', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszStyl.dataset['czynnosc'] = 'zapisz-styl';
  zapiszStyl.title =
    'Idzie komendą studio.style.save. Zmiana stylu przestawia WSZYSTKIE miejsca dokumentu, które ' +
    'go używają — to jest sens stylu nazwanego. Postać znaku i akapitu bierze się z pól poniżej.';
  zapiszStyl.addEventListener('click', () => {
    const nazwa = tekstPola(nazwaStylu.kontrolka);
    if (nazwa === undefined) {
      odpowiedz.pokaz('Styl bez nazwy nie da się później wskazać ani zastosować — nazwij go.', false);
      return;
    }
    const zadanie: Omit<StudioStyleSaveRequest, 'documentId'> = {
      name: nazwa,
      kind: rodzajStylu.kontrolka.value as StudioStyleKind,
    };
    const widoczna = tekstPola(nazwaWidoczna.kontrolka);
    if (widoczna !== undefined) zadanie.displayName = widoczna;
    const nadrzedny = tekstPola(stylNadrzedny.kontrolka);
    if (nadrzedny !== undefined) zadanie.basedOn = nadrzedny;
    const nastepny = tekstPola(stylNastepny.kontrolka);
    if (nastepny !== undefined) zadanie.nextStyle = nastepny;
    const znak = zlozPostacZnaku();
    if (Object.keys(znak).length > 0) zadanie.character = znak;
    const akapit = zlozPostacAkapitu();
    if (Object.keys(akapit).length > 0) zadanie.paragraph = akapit;
    czynnosci.naZapisStylu(zadanie);
  });

  const zastosujStyl = przycisk('Zastosuj styl do zaznaczenia', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zastosujStyl.dataset['czynnosc'] = 'zastosuj-styl';
  zastosujStyl.title = 'Idzie komendą studio.style.apply na zaznaczonym fragmencie.';
  zastosujStyl.addEventListener('click', () => {
    const nazwa = tekstPola(nazwaStylu.kontrolka);
    if (nazwa === undefined) {
      odpowiedz.pokaz('Wskaż styl nazwą — bez niej nie ma czego zastosować.', false);
      return;
    }
    czynnosci.naStosowanieStylu({ name: nazwa });
  });

  const usunStyl = przycisk('Usuń styl własny', 'dn-btn dn-btn--sm dn-btn--ostrzezenie');
  usunStyl.dataset['czynnosc'] = 'usun-styl';
  usunStyl.title =
    'Idzie komendą studio.style.delete. Stylu FABRYCZNEGO rdzeń nie usuwa — odpowiada odmową ' +
    'nazywającą powód, i tak ma zostać.';
  usunStyl.addEventListener('click', () => {
    const nazwa = tekstPola(nazwaStylu.kontrolka);
    if (nazwa === undefined) {
      odpowiedz.pokaz('Wskaż styl nazwą — bez niej nie ma czego usunąć.', false);
      return;
    }
    const zadanie: Omit<StudioStyleDeleteRequest, 'documentId'> = { name: nazwa };
    const zamienny = tekstPola(stylZamienny.kontrolka);
    if (zamienny !== undefined) zadanie.replaceWith = zamienny;
    czynnosci.naUsuniecieStylu(zadanie);
  });

  /* ── Styl znaku ──────────────────────────────────────────────────────────── */

  const krój = poleTekstowe({ etykieta: 'Krój pisma', podpowiedz: 'np. Source Serif' });
  const stopien = poleLiczbowe({ etykieta: 'Stopień pisma w punktach' }, { dolna: 1, gorna: 400, krok: 0.5 });
  const krokStopnia = poleLiczbowe(
    {
      etykieta: 'Powiększ albo pomniejsz stopień o',
      opis: 'Wartość ujemna pomniejsza. Działa bez znajomości stopnia zastanego.',
    },
    { krok: 0.5 },
  );
  const pogrubienie = poleWyboru({ etykieta: 'Pogrubienie' }, TRZY_STANY);
  const kursywa = poleWyboru({ etykieta: 'Kursywa' }, TRZY_STANY);
  const podkreslenie = poleWyboru({ etykieta: 'Podkreślenie' }, [
    BEZ_ZMIANY,
    { wartosc: StudioUnderlineStyle.None, etykieta: 'Bez podkreślenia' },
    { wartosc: StudioUnderlineStyle.Single, etykieta: 'Pojedyncze' },
    { wartosc: StudioUnderlineStyle.Double, etykieta: 'Podwójne' },
    { wartosc: StudioUnderlineStyle.Thick, etykieta: 'Grube' },
    { wartosc: StudioUnderlineStyle.Dotted, etykieta: 'Kropkowane' },
    { wartosc: StudioUnderlineStyle.Dashed, etykieta: 'Kreskowane' },
    { wartosc: StudioUnderlineStyle.Wavy, etykieta: 'Faliste' },
    { wartosc: StudioUnderlineStyle.Words, etykieta: 'Tylko wyrazy, bez spacji' },
  ]);
  const przekreslenie = poleWyboru({ etykieta: 'Przekreślenie' }, TRZY_STANY);
  const indeksGorny = poleWyboru({ etykieta: 'Indeks górny' }, TRZY_STANY);
  const indeksDolny = poleWyboru({ etykieta: 'Indeks dolny' }, TRZY_STANY);
  const barwa = poleTekstowe({ etykieta: 'Barwa tekstu', podpowiedz: '#1a1a1a' });
  const wyroznienie = poleTekstowe({
    etykieta: 'Barwa wyróżnienia tła',
    opis: 'Wyróżnienie jest cechą POSTACI dokumentu, więc przeżywa zapis tak samo jak krój.',
    podpowiedz: '#fff3a3',
  });
  const odstepLiter = poleLiczbowe(
    { etykieta: 'Odstęp międzyliterowy w punktach' },
    { krok: 0.1 },
  );
  const kapitaliki = poleWyboru({ etykieta: 'Kapitaliki' }, TRZY_STANY);
  const wielkieLitery = poleWyboru({ etykieta: 'Wszystkie wielkie litery' }, TRZY_STANY);
  const efekt = poleWyboru({ etykieta: 'Efekt znaku' }, [
    BEZ_ZMIANY,
    { wartosc: StudioTextEffect.None, etykieta: 'Bez efektu' },
    { wartosc: StudioTextEffect.Shadow, etykieta: 'Cień' },
    { wartosc: StudioTextEffect.Outline, etykieta: 'Kontur' },
    { wartosc: StudioTextEffect.Emboss, etykieta: 'Wypukły' },
    { wartosc: StudioTextEffect.Engrave, etykieta: 'Wgłębiony' },
  ]);
  const jezyk = poleTekstowe({
    etykieta: 'Język fragmentu do sprawdzania pisowni',
    podpowiedz: 'pl-PL',
  });
  const niejednoliteZnaku = document.createElement('p');
  niejednoliteZnaku.className = 'dn-pole-opis';

  /** Postać znaku złożona z pól; pola puste zostają nietknięte. */
  function zlozPostacZnaku(): Record<string, unknown> {
    const postac: Record<string, unknown> = {};
    const nazwaKroju = tekstPola(krój.kontrolka);
    if (nazwaKroju !== undefined) postac['fontFamily'] = nazwaKroju;
    const stopienPt = liczbaPola(stopien.kontrolka);
    if (stopienPt !== undefined) postac['fontSizePt'] = stopienPt;
    dopiszLogiczna(postac, 'bold', pogrubienie.kontrolka);
    dopiszLogiczna(postac, 'italic', kursywa.kontrolka);
    const odmiana = wyborPola(podkreslenie.kontrolka);
    if (odmiana !== undefined) postac['underline'] = odmiana;
    dopiszLogiczna(postac, 'strikethrough', przekreslenie.kontrolka);
    dopiszLogiczna(postac, 'superscript', indeksGorny.kontrolka);
    dopiszLogiczna(postac, 'subscript', indeksDolny.kontrolka);
    const barwaTekstu = tekstPola(barwa.kontrolka);
    if (barwaTekstu !== undefined) postac['color'] = barwaTekstu;
    const barwaTla = tekstPola(wyroznienie.kontrolka);
    if (barwaTla !== undefined) postac['highlightColor'] = barwaTla;
    const odstep = liczbaPola(odstepLiter.kontrolka);
    if (odstep !== undefined) postac['letterSpacingPt'] = odstep;
    dopiszLogiczna(postac, 'smallCaps', kapitaliki.kontrolka);
    dopiszLogiczna(postac, 'allCaps', wielkieLitery.kontrolka);
    const wybranyEfekt = wyborPola(efekt.kontrolka);
    if (wybranyEfekt !== undefined) postac['effect'] = wybranyEfekt;
    const kodJezyka = tekstPola(jezyk.kontrolka);
    if (kodJezyka !== undefined) postac['language'] = kodJezyka;
    return postac;
  }

  const zapiszZnak = przycisk('Nanieś styl znaku na zaznaczenie', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszZnak.dataset['czynnosc'] = 'nanies-styl-znaku';
  zapiszZnak.title = 'Idzie komendą studio.format.character.set na zaznaczonym fragmencie.';
  zapiszZnak.addEventListener('click', () => {
    const zadanie = zlozPostacZnaku() as BezZakresu<StudioFormatCharacterSetRequest>;
    const krok = liczbaPola(krokStopnia.kontrolka);
    if (krok !== undefined) zadanie.fontSizeStepPt = krok;
    if (Object.keys(zadanie).length === 0) {
      odpowiedz.pokaz(
        'Żadne pole stylu znaku nie jest wypełnione — nie ma czego nanieść. Puste pole znaczy ' +
          '„nie ruszaj tej cechy", więc żądanie bez ani jednego pola nie zmieniłoby niczego.',
        false,
      );
      return;
    }
    czynnosci.naStylZnaku(zadanie);
  });

  /* ── Styl akapitu ────────────────────────────────────────────────────────── */

  const wyrownanie = poleWyboru({ etykieta: 'Wyrównanie' }, [
    BEZ_ZMIANY,
    { wartosc: StudioTextAlign.Left, etykieta: 'Do lewej' },
    { wartosc: StudioTextAlign.Center, etykieta: 'Wyśrodkowanie' },
    { wartosc: StudioTextAlign.Right, etykieta: 'Do prawej' },
    { wartosc: StudioTextAlign.Justify, etykieta: 'Wyjustowanie' },
  ]);
  const wciecieWiersza = poleLiczbowe(
    {
      etykieta: 'Wcięcie pierwszego wiersza w milimetrach',
      opis: 'Wartość ujemna daje wysunięcie.',
    },
    { krok: 0.5 },
  );
  const wciecieLewe = poleLiczbowe({ etykieta: 'Wcięcie lewe w milimetrach' }, { krok: 0.5 });
  const wciecePrawe = poleLiczbowe({ etykieta: 'Wcięcie prawe w milimetrach' }, { krok: 0.5 });
  const krokWciecia = poleLiczbowe(
    {
      etykieta: 'Zwiększ albo zmniejsz wcięcie lewe o',
      opis: 'Wartość ujemna zmniejsza — to jest droga przycisków wcięcia ze wstążki.',
    },
    { krok: 0.5 },
  );
  const odstepPrzed = poleLiczbowe({ etykieta: 'Odstęp przed akapitem w punktach' }, { dolna: 0, krok: 0.5 });
  const odstepPo = poleLiczbowe({ etykieta: 'Odstęp po akapicie w punktach' }, { dolna: 0, krok: 0.5 });
  const zasadaInterlinii = poleWyboru({ etykieta: 'Zasada interlinii' }, [
    BEZ_ZMIANY,
    { wartosc: StudioLineSpacingRule.Single, etykieta: 'Pojedyncza' },
    { wartosc: StudioLineSpacingRule.OneAndHalf, etykieta: 'Półtora' },
    { wartosc: StudioLineSpacingRule.Double, etykieta: 'Podwójna' },
    { wartosc: StudioLineSpacingRule.AtLeast, etykieta: 'Co najmniej podana wartość' },
    { wartosc: StudioLineSpacingRule.Exactly, etykieta: 'Dokładnie podana wartość' },
    { wartosc: StudioLineSpacingRule.Multiple, etykieta: 'Wielokrotność pojedynczej' },
  ]);
  const wartoscInterlinii = poleLiczbowe({ etykieta: 'Wartość interlinii' }, { dolna: 0, krok: 0.05 });
  const cieniowanie = poleTekstowe({ etykieta: 'Cieniowanie tła akapitu', podpowiedz: '#f4f4f4' });
  const wdowy = poleWyboru({ etykieta: 'Kontrola wdów i sierot' }, TRZY_STANY);
  const zNastepnym = poleWyboru({ etykieta: 'Razem z następnym akapitem' }, TRZY_STANY);
  const bezDzielenia = poleWyboru({ etykieta: 'Bez dzielenia akapitu między strony' }, TRZY_STANY);
  const poziomKonspektu = poleLiczbowe(
    { etykieta: 'Poziom konspektu', opis: 'Zero znaczy tekst zasadniczy.' },
    { dolna: 0, gorna: 9 },
  );
  const odPrawej = poleWyboru({ etykieta: 'Kierunek pisma od prawej do lewej' }, TRZY_STANY);
  const niejednoliteAkapitu = document.createElement('p');
  niejednoliteAkapitu.className = 'dn-pole-opis';

  function zlozPostacAkapitu(): Record<string, unknown> {
    const postac: Record<string, unknown> = {};
    const wybrane = wyborPola(wyrownanie.kontrolka);
    if (wybrane !== undefined) postac['align'] = wybrane;
    const pierwszy = liczbaPola(wciecieWiersza.kontrolka);
    if (pierwszy !== undefined) postac['firstLineIndentMm'] = pierwszy;
    const lewe = liczbaPola(wciecieLewe.kontrolka);
    if (lewe !== undefined) postac['indentLeftMm'] = lewe;
    const prawe = liczbaPola(wciecePrawe.kontrolka);
    if (prawe !== undefined) postac['indentRightMm'] = prawe;
    const przed = liczbaPola(odstepPrzed.kontrolka);
    if (przed !== undefined) postac['spaceBeforePt'] = przed;
    const po = liczbaPola(odstepPo.kontrolka);
    if (po !== undefined) postac['spaceAfterPt'] = po;
    const zasada = wyborPola(zasadaInterlinii.kontrolka);
    if (zasada !== undefined) postac['lineSpacingRule'] = zasada;
    const wartosc = liczbaPola(wartoscInterlinii.kontrolka);
    if (wartosc !== undefined) postac['lineSpacingValue'] = wartosc;
    const tlo = tekstPola(cieniowanie.kontrolka);
    if (tlo !== undefined) postac['shadingColor'] = tlo;
    dopiszLogiczna(postac, 'widowControl', wdowy.kontrolka);
    dopiszLogiczna(postac, 'keepWithNext', zNastepnym.kontrolka);
    dopiszLogiczna(postac, 'keepLines', bezDzielenia.kontrolka);
    const poziom = liczbaPola(poziomKonspektu.kontrolka);
    if (poziom !== undefined) postac['outlineLevel'] = poziom;
    dopiszLogiczna(postac, 'rightToLeft', odPrawej.kontrolka);
    return postac;
  }

  const zapiszAkapit = przycisk(
    'Nanieś styl akapitu na zaznaczenie',
    'dn-btn dn-btn--sm dn-btn--sygnal',
  );
  zapiszAkapit.dataset['czynnosc'] = 'nanies-styl-akapitu';
  zapiszAkapit.title = 'Idzie komendą studio.format.paragraph.set na zaznaczonym fragmencie.';
  zapiszAkapit.addEventListener('click', () => {
    const zadanie = zlozPostacAkapitu() as BezZakresu<StudioFormatParagraphSetRequest>;
    const krok = liczbaPola(krokWciecia.kontrolka);
    if (krok !== undefined) zadanie.indentStepMm = krok;
    if (Object.keys(zadanie).length === 0) {
      odpowiedz.pokaz(
        'Żadne pole stylu akapitu nie jest wypełnione — nie ma czego nanieść.',
        false,
      );
      return;
    }
    czynnosci.naStylAkapitu(zadanie);
  });

  /* ── Tabulator liczbą ────────────────────────────────────────────────────── */

  const polozenieTabulatora = poleLiczbowe(
    {
      etykieta: 'Tabulator — położenie od lewego marginesu w milimetrach',
      opis:
        'Ta sama komenda, którą jedzie chwyt tabulatora z linijki poziomej ' +
        '(studio.ruler.tabstop.set). Tu wpisuje się go liczbą, bo praca wyłącznie myszą odcięłaby ' +
        'część Operatorów.',
    },
    { dolna: 0, gorna: 1000, krok: 0.5 },
  );
  const rodzajTabulatora = poleWyboru({ etykieta: 'Rodzaj tabulatora' }, [
    BEZ_ZMIANY,
    { wartosc: StudioTabKind.Left, etykieta: 'Lewy' },
    { wartosc: StudioTabKind.Right, etykieta: 'Prawy' },
    { wartosc: StudioTabKind.Center, etykieta: 'Środkowy' },
    { wartosc: StudioTabKind.Decimal, etykieta: 'Dziesiętny' },
    { wartosc: StudioTabKind.Bar, etykieta: 'Kreska' },
  ]);
  const znakWiodacy = poleWyboru({ etykieta: 'Znak wiodący' }, [
    BEZ_ZMIANY,
    { wartosc: StudioTabLeader.None, etykieta: 'Bez znaku wiodącego' },
    { wartosc: StudioTabLeader.Dot, etykieta: 'Kropki' },
    { wartosc: StudioTabLeader.Dash, etykieta: 'Kreski' },
    { wartosc: StudioTabLeader.Underline, etykieta: 'Linia ciągła' },
  ]);
  const zdejmijTabulator = poleLogiczne({ etykieta: 'Zdejmij tabulator z tego położenia' });

  const zapiszTabulator = przycisk('Ustaw tabulator', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zapiszTabulator.dataset['czynnosc'] = 'ustaw-tabulator';
  zapiszTabulator.addEventListener('click', () => {
    const polozenie = liczbaPola(polozenieTabulatora.kontrolka);
    if (polozenie === undefined) {
      odpowiedz.pokaz(
        'Tabulator bez położenia nie ma gdzie stanąć — podaj odległość od lewego marginesu.',
        false,
      );
      return;
    }
    const zadanie: BezZakresu<StudioRulerTabstopSetRequest> = {
      positionMm: polozenie,
      remove: zdejmijTabulator.kontrolka.checked,
    };
    const rodzaj = wyborPola(rodzajTabulatora.kontrolka);
    if (rodzaj !== undefined) zadanie.kind = rodzaj as StudioTabKind;
    const znak = wyborPola(znakWiodacy.kontrolka);
    if (znak !== undefined) zadanie.leader = znak as StudioTabLeader;
    czynnosci.naTabulator(zadanie);
  });

  /* ── Czyszczenie, wielkość liter, malarz, podobne ────────────────────────── */

  const czyscZnak = poleLogiczne({ etykieta: 'Czyść styl znaku' });
  czyscZnak.kontrolka.checked = true;
  const czyscAkapit = poleLogiczne({ etykieta: 'Czyść styl akapitu' });
  czyscAkapit.kontrolka.checked = true;

  const wyczysc = przycisk('Wyczyść formatowanie zaznaczenia', 'dn-btn dn-btn--sm dn-btn--duch');
  wyczysc.dataset['czynnosc'] = 'wyczysc-format';
  wyczysc.title =
    'Idzie komendą studio.format.clear. Postać wraca do stylu nazwanego albo do postaci domyślnej — ' +
    'litery zostają.';
  wyczysc.addEventListener('click', () => {
    const zadanie: BezZakresu<StudioFormatClearRequest> = {
      character: czyscZnak.kontrolka.checked,
      paragraph: czyscAkapit.kontrolka.checked,
    };
    if (!zadanie.character && !zadanie.paragraph) {
      odpowiedz.pokaz(
        'Wskaż, co czyścić — ani styl znaku, ani styl akapitu nie jest zaznaczony, więc czynność ' +
          'nie zmieniłaby niczego.',
        false,
      );
      return;
    }
    czynnosci.naCzyszczenie(zadanie);
  });

  const wielkoscLiter = poleWyboru({ etykieta: 'Wielkość liter zaznaczenia' }, [
    { wartosc: StudioCaseTransform.Upper, etykieta: 'WSZYSTKIE WIELKIE' },
    { wartosc: StudioCaseTransform.Lower, etykieta: 'wszystkie małe' },
    { wartosc: StudioCaseTransform.Sentence, etykieta: 'Jak w zdaniu' },
    { wartosc: StudioCaseTransform.Capitalize, etykieta: 'Każde Słowo Wielką Literą' },
    { wartosc: StudioCaseTransform.Toggle, etykieta: 'Odwrócenie wielkości' },
  ]);

  const przestawLitery = przycisk('Przestaw wielkość liter', 'dn-btn dn-btn--sm dn-btn--duch');
  przestawLitery.dataset['czynnosc'] = 'przestaw-litery';
  przestawLitery.addEventListener('click', () => {
    const zadanie: BezZakresu<StudioFormatCaseSetRequest> = {
      transform: wielkoscLiter.kontrolka.value as StudioCaseTransform,
    };
    czynnosci.naWielkoscLiter(zadanie);
  });

  const malarzZAkapitem = poleLogiczne({ etykieta: 'Malarz obejmuje także styl akapitu' });
  malarzZAkapitem.kontrolka.checked = true;
  const stanMalarza = document.createElement('p');
  stanMalarza.className = 'dn-pole-opis';
  stanMalarza.textContent = 'Malarz formatów nie ma nic pobranego.';

  const pobierz = przycisk('Malarz: pobierz postać zaznaczenia', 'dn-btn dn-btn--sm dn-btn--zarys');
  pobierz.dataset['czynnosc'] = 'malarz-pobierz';
  pobierz.title = 'Idzie komendą studio.format.painter.copy — kopiuje POSTAĆ, nie treść.';
  pobierz.addEventListener('click', () => {
    czynnosci.naPobraniePostaci({ includeParagraph: malarzZAkapitem.kontrolka.checked });
  });

  const naloz = przycisk('Malarz: nanieś na zaznaczenie', 'dn-btn dn-btn--sm dn-btn--sygnal');
  naloz.dataset['czynnosc'] = 'malarz-naloz';
  naloz.title = 'Idzie komendą studio.format.painter.apply z uchwytem pobranej postaci.';
  naloz.addEventListener('click', () => {
    czynnosci.naNalozeniePostaci({ includeParagraph: malarzZAkapitem.kontrolka.checked });
  });

  const podobneStyl = poleTekstowe({
    etykieta: 'Wzór podobieństwa — styl nazwany',
    opis: 'Puste znaczy „wzorem jest postać zaznaczenia".',
  });
  const podobneZnak = poleLogiczne({ etykieta: 'Dopasuj styl znaku' });
  podobneZnak.kontrolka.checked = true;
  const podobneAkapit = poleLogiczne({ etykieta: 'Dopasuj styl akapitu' });

  const zaznaczPodobne = przycisk(
    'Zaznacz wedle podobnego formatowania',
    'dn-btn dn-btn--sm dn-btn--zarys',
  );
  zaznaczPodobne.dataset['czynnosc'] = 'podobne-formatowanie';
  zaznaczPodobne.title = 'Idzie komendą studio.format.similar.select.';
  zaznaczPodobne.addEventListener('click', () => {
    const zadanie: BezZakresu<StudioFormatSimilarSelectRequest> = {
      matchCharacter: podobneZnak.kontrolka.checked,
      matchParagraph: podobneAkapit.kontrolka.checked,
    };
    const styl = tekstPola(podobneStyl.kontrolka);
    if (styl !== undefined) zadanie.styleName = styl;
    czynnosci.naPodobne(zadanie);
  });

  /* ── Zamiana wraz z postacią ─────────────────────────────────────────────── */

  const szukanaTresc = poleTekstowe({ etykieta: 'Szukana treść' });
  const wstawianaTresc = poleTekstowe({ etykieta: 'Treść wstawiana' });
  const szukanyStyl = poleTekstowe({ etykieta: 'Szukany styl nazwany' });
  const nadawanyStyl = poleTekstowe({ etykieta: 'Styl nazwany nadawany trafieniom' });
  const wielkoscLiterZamiany = poleLogiczne({ etykieta: 'Odróżniaj wielkość liter' });
  const caleWyrazy = poleLogiczne({ etykieta: 'Dopasowuj całe wyrazy' });
  const wyrazenie = poleLogiczne({ etykieta: 'Szukana treść jest wyrażeniem' });
  const wszystkieTrafienia = poleLogiczne({ etykieta: 'Zamień wszystkie trafienia' });
  wszystkieTrafienia.kontrolka.checked = true;
  const tylkoZaznaczenie = poleLogiczne({
    etykieta: 'Tylko w zaznaczeniu',
    opis: 'Wyłączone znaczy „w całym dokumencie", niezależnie od tego, co jest zaznaczone.',
  });

  const zamien = przycisk('Znajdź i zamień wraz z postacią', 'dn-btn dn-btn--sm dn-btn--sygnal');
  zamien.dataset['czynnosc'] = 'zamien-z-postacia';
  zamien.title =
    'Idzie komendą studio.format.replace. Fragmenty pod blokadą NIE zatrzymują całości — zamiana ' +
    'wykonuje się poza blokadą i wraca bilansem nazywającym, co pominięto i przez którą blokadę.';
  zamien.addEventListener('click', () => {
    const zadanie: BezZakresu<StudioFormatReplaceRequest> = {
      matchCase: wielkoscLiterZamiany.kontrolka.checked,
      wholeWord: caleWyrazy.kontrolka.checked,
      regex: wyrazenie.kontrolka.checked,
      replaceAll: wszystkieTrafienia.kontrolka.checked,
    };
    const szukana = tekstPola(szukanaTresc.kontrolka);
    if (szukana !== undefined) zadanie.findText = szukana;
    const wstawiana = tekstPola(wstawianaTresc.kontrolka);
    if (wstawiana !== undefined) zadanie.replaceText = wstawiana;
    const szukanyNazwany = tekstPola(szukanyStyl.kontrolka);
    if (szukanyNazwany !== undefined) zadanie.findStyleName = szukanyNazwany;
    const nadawanyNazwany = tekstPola(nadawanyStyl.kontrolka);
    if (nadawanyNazwany !== undefined) zadanie.replaceStyleName = nadawanyNazwany;
    const szukanaPostac = zlozPostacZnaku();
    // Postać z pól służy naniesieniu i zamianie: przy zamianie jest postacią nadawaną trafieniom.
    if (Object.keys(szukanaPostac).length > 0) zadanie.replaceFormat = szukanaPostac;
    if (
      zadanie.findText === undefined &&
      zadanie.findStyleName === undefined &&
      zadanie.findFormat === undefined
    ) {
      odpowiedz.pokaz(
        'Nie ma czego szukać: podaj treść, styl nazwany albo postać. Zamiana bez wzoru objęłaby ' +
          'cały dokument i nie jest tym, o co Operator prosi.',
        false,
      );
      return;
    }
    czynnosci.naZamiane(zadanie);
  });

  const odczytaj = przycisk('Odczytaj arkusz stylów i postać zaznaczenia', 'dn-btn dn-btn--sm dn-btn--zarys');
  odczytaj.dataset['czynnosc'] = 'odczytaj-style';
  odczytaj.title =
    'Idzie komendami studio.style.list, studio.format.character.get i studio.format.paragraph.get.';
  odczytaj.addEventListener('click', () => czynnosci.naOdczyt());

  /* ── Układ panelu ────────────────────────────────────────────────────────── */

  const element = document.createElement('section');
  element.className = 'ms-postac ms-postac--styl';
  element.dataset['panel'] = 'arkusz-stylow';
  element.hidden = true;
  element.setAttribute('aria-label', 'Arkusz stylów i formatowanie fragmentu');

  element.append(
    grupa('Arkusz stylów nazwanych', [
      odczytaj,
      wykazStylow,
      nazwaStylu.element,
      nazwaWidoczna.element,
      rodzajStylu.element,
      stylNadrzedny.element,
      stylNastepny.element,
      stylZamienny.element,
      zapiszStyl,
      zastosujStyl,
      usunStyl,
    ]),
    grupa('Styl znaku', [
      niejednoliteZnaku,
      krój.element,
      stopien.element,
      krokStopnia.element,
      pogrubienie.element,
      kursywa.element,
      podkreslenie.element,
      przekreslenie.element,
      indeksGorny.element,
      indeksDolny.element,
      barwa.element,
      wyroznienie.element,
      odstepLiter.element,
      kapitaliki.element,
      wielkieLitery.element,
      efekt.element,
      jezyk.element,
      zapiszZnak,
    ]),
    grupa('Styl akapitu', [
      niejednoliteAkapitu,
      wyrownanie.element,
      wciecieWiersza.element,
      wciecieLewe.element,
      wciecePrawe.element,
      krokWciecia.element,
      odstepPrzed.element,
      odstepPo.element,
      zasadaInterlinii.element,
      wartoscInterlinii.element,
      cieniowanie.element,
      wdowy.element,
      zNastepnym.element,
      bezDzielenia.element,
      poziomKonspektu.element,
      odPrawej.element,
      zapiszAkapit,
    ]),
    grupa('Tabulatory', [
      polozenieTabulatora.element,
      rodzajTabulatora.element,
      znakWiodacy.element,
      zdejmijTabulator.element,
      zapiszTabulator,
    ]),
    grupa('Czyszczenie, wielkość liter, malarz formatów', [
      czyscZnak.element,
      czyscAkapit.element,
      wyczysc,
      wielkoscLiter.element,
      przestawLitery,
      malarzZAkapitem.element,
      stanMalarza,
      pobierz,
      naloz,
    ]),
    grupa('Zaznacz wedle podobnego formatowania', [
      podobneStyl.element,
      podobneZnak.element,
      podobneAkapit.element,
      zaznaczPodobne,
    ]),
    grupa('Znajdź i zamień wraz z postacią', [
      szukanaTresc.element,
      wstawianaTresc.element,
      szukanyStyl.element,
      nadawanyStyl.element,
      wielkoscLiterZamiany.element,
      caleWyrazy.element,
      wyrazenie.element,
      wszystkieTrafienia.element,
      tylkoZaznaczenie.element,
      zamien,
    ]),
    odpowiedz.element,
  );

  return {
    element,

    przestawWidocznosc() {
      element.hidden = !element.hidden;
    },

    widoczny: () => !element.hidden,

    pokazStyle(style) {
      wykazStylow.replaceChildren(
        ...style.map((styl) => {
          const pozycja = document.createElement('li');
          pozycja.className = 'ms-postac__styl';
          pozycja.dataset['styl'] = styl.name;
          pozycja.textContent =
            `${styl.displayName ?? styl.name} · ${styl.kind}` +
            `${styl.builtin === true ? ' · fabryczny' : ' · własny'}` +
            `${styl.basedOn === undefined ? '' : ` · po ${styl.basedOn}`}` +
            ` · miejsc użycia ${styl.usageCount ?? 0}`;
          const wybierz = przycisk('Wskaż', 'dn-btn dn-btn--sm dn-btn--zarys');
          wybierz.addEventListener('click', () => {
            nazwaStylu.kontrolka.value = styl.name;
            nazwaWidoczna.kontrolka.value = styl.displayName ?? '';
            rodzajStylu.kontrolka.value = styl.kind;
            stylNadrzedny.kontrolka.value = styl.basedOn ?? '';
            stylNastepny.kontrolka.value = styl.nextStyle ?? '';
          });
          pozycja.append(wybierz);
          return pozycja;
        }),
      );
      if (style.length === 0) {
        const pusty = document.createElement('li');
        pusty.className = 'dn-pole-opis';
        pusty.textContent =
          'Rdzeń nie oddał ani jednego stylu. Dokument bez arkusza stylów jest możliwy — ' +
          'załóż pierwszy nazwą i postacią poniżej.';
        wykazStylow.append(pusty);
      }
    },

    pokazPostacZnaku(postac, niejednolite) {
      krój.kontrolka.value = postac.fontFamily ?? '';
      stopien.kontrolka.value = postac.fontSizePt === undefined ? '' : String(postac.fontSizePt);
      ustawTrzyStany(pogrubienie.kontrolka, postac.bold);
      ustawTrzyStany(kursywa.kontrolka, postac.italic);
      podkreslenie.kontrolka.value = postac.underline ?? '';
      ustawTrzyStany(przekreslenie.kontrolka, postac.strikethrough);
      ustawTrzyStany(indeksGorny.kontrolka, postac.superscript);
      ustawTrzyStany(indeksDolny.kontrolka, postac.subscript);
      barwa.kontrolka.value = postac.color ?? '';
      wyroznienie.kontrolka.value = postac.highlightColor ?? '';
      odstepLiter.kontrolka.value =
        postac.letterSpacingPt === undefined ? '' : String(postac.letterSpacingPt);
      ustawTrzyStany(kapitaliki.kontrolka, postac.smallCaps);
      ustawTrzyStany(wielkieLitery.kontrolka, postac.allCaps);
      efekt.kontrolka.value = postac.effect ?? '';
      jezyk.kontrolka.value = postac.language ?? '';
      niejednoliteZnaku.textContent =
        niejednolite.length === 0
          ? 'Postać znaku jest we fragmencie jednolita.'
          : `We fragmencie NIEJEDNOLITE: ${niejednolite.join(', ')}. Pola tych cech pokazują ` +
            'wartość wiodącą, a nie jedyną — naniesienie ujednolici je w całym zaznaczeniu.';
    },

    pokazPostacAkapitu(postac, niejednolite) {
      wyrownanie.kontrolka.value = postac.align ?? '';
      wciecieWiersza.kontrolka.value =
        postac.firstLineIndentMm === undefined ? '' : String(postac.firstLineIndentMm);
      wciecieLewe.kontrolka.value =
        postac.indentLeftMm === undefined ? '' : String(postac.indentLeftMm);
      wciecePrawe.kontrolka.value =
        postac.indentRightMm === undefined ? '' : String(postac.indentRightMm);
      odstepPrzed.kontrolka.value =
        postac.spaceBeforePt === undefined ? '' : String(postac.spaceBeforePt);
      odstepPo.kontrolka.value = postac.spaceAfterPt === undefined ? '' : String(postac.spaceAfterPt);
      zasadaInterlinii.kontrolka.value = postac.lineSpacingRule ?? '';
      wartoscInterlinii.kontrolka.value =
        postac.lineSpacingValue === undefined ? '' : String(postac.lineSpacingValue);
      cieniowanie.kontrolka.value = postac.shadingColor ?? '';
      ustawTrzyStany(wdowy.kontrolka, postac.widowControl);
      ustawTrzyStany(zNastepnym.kontrolka, postac.keepWithNext);
      ustawTrzyStany(bezDzielenia.kontrolka, postac.keepLines);
      poziomKonspektu.kontrolka.value =
        postac.outlineLevel === undefined ? '' : String(postac.outlineLevel);
      ustawTrzyStany(odPrawej.kontrolka, postac.rightToLeft);
      niejednoliteAkapitu.textContent =
        niejednolite.length === 0
          ? 'Postać akapitu jest we fragmencie jednolita.'
          : `We fragmencie NIEJEDNOLITE: ${niejednolite.join(', ')}.`;
    },

    pokazMalarza(pobrane, opis) {
      stanMalarza.textContent = pobrane
        ? `Malarz formatów ma pobraną postać: ${opis}`
        : 'Malarz formatów nie ma nic pobranego.';
      naloz.disabled = !pobrane;
    },

    pokazOdpowiedz: (tresc, powodzenie) => odpowiedz.pokaz(tresc, powodzenie),
    zamianaTylkoWZaznaczeniu: () => tylkoZaznaczenie.kontrolka.checked,
  };
}

/**
 * Trzy stany cechy logicznej: bez zmiany, włącz, wyłącz; dwustanowy przełącznik nie umiałby
 * nazwać stanu bez zmiany.
 */
const TRZY_STANY = [
  BEZ_ZMIANY,
  { wartosc: 'tak', etykieta: 'Włącz' },
  { wartosc: 'nie', etykieta: 'Wyłącz' },
];

/** Dopisuje cechę logiczną do postaci znaku albo akapitu, o ile Operator wskazał jej wartość w formularzu. */
function dopiszLogiczna(
  postac: Record<string, unknown>,
  pole: string,
  kontrolka: HTMLSelectElement,
): void {
  const wartosc = wyborPola(kontrolka);
  if (wartosc === undefined) return;
  postac[pole] = wartosc === 'tak';
}

/** Ustawia listę trzech stanów pola formularza wedle wartości logicznej otrzymanej z rdzenia dokumentu. */
function ustawTrzyStany(kontrolka: HTMLSelectElement, wartosc: boolean | undefined): void {
  kontrolka.value = wartosc === undefined ? '' : wartosc ? 'tak' : 'nie';
}

/** Grupa pól formularza wyświetlana pod wspólnym tytułem sekcji wewnątrz panelu arkusza stylów dokumentu. */
function grupa(tytul: string, zawartosc: readonly HTMLElement[]): HTMLElement {
  const naglowek = document.createElement('h4');
  naglowek.className = 'ms-postac__tytul';
  naglowek.textContent = tytul;

  const element = document.createElement('div');
  element.className = 'ms-postac__grupa';
  element.append(naglowek, ...zawartosc);
  return element;
}
