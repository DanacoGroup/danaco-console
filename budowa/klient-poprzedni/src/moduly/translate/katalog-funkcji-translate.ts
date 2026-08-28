import { Command } from '../../../../shared/contract';
import type { WarstwaWidocznosci } from './warstwy-translate';

/**
 * Katalog funkcji modułu Translate: wykaz z opracowania wraz z komendą, która
 * pozycję wykonuje.
 */

/**
 * Nazwy okien operacyjnych modułu Translate w jednym brzmieniu na cały moduł,
 * aby pozycja katalogu i okno, które ją pokazuje, nazywały to samo miejsce
 * jednym słowem.
 */
export const OKNA = {
  chat: 'Chat Window',
  petla: 'Execution Loop Window',
  zrodlo: 'Source Panel',
  panele: 'Translation Panels',
  glosariusz: 'Glossary Manager',
  pamiec: 'Translation Memory Panel',
  formaty: 'Format Studio',
  jakosc: 'QA & Review Center',
} as const;

/**
 * Jedna pozycja katalogu funkcji: nazwa, grupa i opis przejęte z opracowania
 * modułu, okno i warstwa widoczności w tej budowie oraz komenda kontraktu albo
 * nazwany brak pokrycia.
 */
export interface PozycjaKatalogu {
  /** Nazwa własna pozycji, dokładnie jak w opracowaniu modułu. */
  readonly nazwa: string;
  /** Grupa tematyczna opracowania. */
  readonly grupa: string;
  /** Co pozycja robi. */
  readonly opis: string;
  /** Okno modułu, w którym pozycja jest osiągalna. */
  readonly okno: string;
  /** Warstwa widoczności pozycji. */
  readonly warstwa: WarstwaWidocznosci;
  /** Komendy kontraktu wykonujące pozycję; wykaz pusty znaczy brak komendy. */
  readonly komendy: readonly string[];
  /** Czego brakuje — obowiązkowe przy braku komendy, dopuszczalne przy pokryciu częściowym. */
  readonly brak?: string;
}

const SILNIKI = 'Silniki tłumaczenia i tryby pracy';
const SEGMENTACJA = 'Segmentacja, wyrównanie i struktura tekstu';
const PAMIEC = 'Pamięć tłumaczeń';
const TERMINOLOGIA = 'Terminologia i glosariusze';
const KOREKTA = 'Korekta i redakcja językowa';
const JAKOSC = 'Kontrola jakości tłumaczenia';
const FORMATY = 'Tłumaczenie dokumentów z zachowaniem formatu';
const LOKALIZACJA = 'Lokalizacja oprogramowania i treści cyfrowych';
const MOWA = 'Napisy, treść mówiona i głos';
const WYDANIE = 'Praca zespołowa, wydanie i integracje';

export const KATALOG_FUNKCJI: readonly PozycjaKatalogu[] = [
  {
    nazwa: 'Multi-Engine Translation',
    grupa: SILNIKI,
    opis: 'Tłumaczenie jednego materiału przez wiele silników z porównaniem wariantów.',
    okno: OKNA.panele,
    warstwa: 2,
    komendy: [Command.TranslateTargetAdd],
    brak:
      'Żądanie niesie jeden kanał modelu na panel, więc silnik wybiera się per kolumna. ' +
      'Porównania kilku silników dla jednego segmentu nie ma czym zlecić.',
  },
  {
    nazwa: 'Parallel Panels (N języków)',
    grupa: SILNIKI,
    opis: 'Ten sam tekst źródłowy tłumaczony jednocześnie na wiele języków w siatce paneli.',
    okno: OKNA.panele,
    warstwa: 1,
    komendy: [Command.TranslateTargetAdd, Command.TranslateSourceSet],
  },
  {
    nazwa: 'Post-Editing Mode (MTPE)',
    grupa: SILNIKI,
    opis: 'Korekta tłumaczenia maszynowego z rejestrem różnic wobec wersji finalnej.',
    okno: OKNA.panele,
    warstwa: 1,
    komendy: [Command.TranslateTranslationSet],
    brak:
      'Korekta zapisuje się w rdzeniu, ale rejestru różnic wynik maszynowy → wersja finalna ' +
      'kontrakt nie prowadzi: odpowiedź niesie panel po zmianie, nie jego wersję poprzednią.',
  },
  {
    nazwa: 'Pivot Translation',
    grupa: SILNIKI,
    opis: 'Tłumaczenie przez język pośredni, gdy brak pary bezpośredniej, z oznaczeniem trasy.',
    okno: OKNA.panele,
    warstwa: 4,
    komendy: [],
    brak:
      'Żądanie dodania panelu niesie okno, język, ton i kanał — żadnego pola na język ' +
      'pośredni ani na trasę przekładu.',
  },
  {
    nazwa: 'Back-Translation',
    grupa: SILNIKI,
    opis: 'Tłumaczenie zwrotne wyniku na język źródłowy jako kontrola sensu.',
    okno: OKNA.panele,
    warstwa: 3,
    komendy: [Command.TranslateBacktranslationRun],
  },
  {
    nazwa: 'Register/Tone per Panel',
    grupa: SILNIKI,
    opis: 'Niezależny rejestr każdego panelu: formalny, techniczny, marketingowy, prawniczy.',
    okno: OKNA.panele,
    warstwa: 2,
    komendy: [Command.TranslatePanelToneSet, Command.TranslateTargetAdd],
  },
  {
    nazwa: 'Bulk/Batch Translate',
    grupa: SILNIKI,
    opis: 'Tłumaczenie zbioru plików albo segmentów jednym zleceniem, z kolejką postępu.',
    okno: OKNA.petla,
    warstwa: 3,
    komendy: [],
    brak:
      'Kontrakt nie ma komendy zlecenia wsadowego w obszarze translate. Powtórzenie czynności ' +
      'panel po panelu wykonuje QA & Review Center, ale kolejki zadań to nie zakłada.',
  },
  {
    nazwa: 'Adaptive/Domain MT',
    grupa: SILNIKI,
    opis: 'Dostrojenie przekładu do dziedziny i do zatwierdzonej pamięci tłumaczeń.',
    okno: OKNA.panele,
    warstwa: 4,
    komendy: [],
    brak: 'Żądanie przekładu nie niesie ani dziedziny, ani wskazania pamięci jako kontekstu.',
  },

  {
    nazwa: 'Sentence/Segment Engine (SRX)',
    grupa: SEGMENTACJA,
    opis: 'Podział tekstu na jednostki tłumaczeniowe z numeracją wspólną wszystkim panelom.',
    okno: OKNA.zrodlo,
    warstwa: 1,
    komendy: [Command.TranslateSourceSegment, Command.TranslateSourceSet],
    brak:
      'Podział wykonuje rdzeń. Żądanie niesie sam tekst — pola na reguły segmentacji SRX ' +
      'nie ma ani w nim, ani w żądaniu zapisu źródła, więc granic nie da się dostroić z okna.',
  },
  {
    nazwa: 'Bitext Aligner',
    grupa: SEGMENTACJA,
    opis: 'Wyrównanie gotowego tekstu źródłowego i jego przekładu w pary segmentów.',
    okno: OKNA.pamiec,
    warstwa: 3,
    komendy: [],
    brak: 'Kontrakt nie ma komendy przyjmującej parę tekstów do wyrównania.',
  },
  {
    nazwa: 'Concordance Search',
    grupa: SEGMENTACJA,
    opis: 'Wyszukanie fragmentu w całej pamięci tłumaczeń wraz z kontekstem.',
    okno: OKNA.pamiec,
    warstwa: 1,
    komendy: [Command.TranslateMemorySuggest],
    brak:
      'Podpowiedź wraca jako wykaz napisów bez kontekstu, bez metadanych i bez procentu ' +
      'dopasowania; przeglądu całej pamięci kontrakt nie ma.',
  },
  {
    nazwa: 'Re-segmentation & Merge/Split',
    grupa: SEGMENTACJA,
    opis: 'Ręczne łączenie i dzielenie segmentów bez utraty powiązań z panelami.',
    okno: OKNA.zrodlo,
    warstwa: 3,
    komendy: [Command.TranslateSourceSegment],
    brak:
      'Ponowny podział całości jest wykonalny; łączenia i dzielenia pojedynczego segmentu ' +
      'kontrakt nie niesie — segment nie jest w nim bytem o własnym identyfikatorze.',
  },
  {
    nazwa: 'Tag/Placeholder Protection',
    grupa: SEGMENTACJA,
    opis: 'Wykrycie znaczników formatu, zmiennych i symboli zastępczych chronionych przed przekładem.',
    okno: OKNA.zrodlo,
    warstwa: 1,
    komendy: [Command.TranslateQualityCheck],
    brak:
      'Wykrycie wzorców w tekście źródłowym liczy okno; rdzeń zgłasza niezgodność symbolu ' +
      'zastępczego dopiero w kontroli jakości panelu.',
  },

  {
    nazwa: 'TM Store & Fuzzy Match',
    grupa: PAMIEC,
    opis: 'Zapamiętane pary segmentów podpowiadane z dopasowaniem przybliżonym.',
    okno: OKNA.pamiec,
    warstwa: 1,
    komendy: [Command.TranslateMemorySuggest],
    brak:
      'Progu dopasowania żądanie nie niesie, a odpowiedź nie oddaje procentu podobieństwa ' +
      'ani rekordu pamięci — same napisy podpowiedzi.',
  },
  {
    nazwa: 'TMX Import/Export',
    grupa: PAMIEC,
    opis: 'Wymiana pamięci tłumaczeń z narzędziami zewnętrznymi w standardzie TMX.',
    okno: OKNA.pamiec,
    warstwa: 3,
    komendy: [],
    brak:
      'Kontrakt ma import i eksport glosariusza, nie ma ich dla pamięci tłumaczeń — ' +
      'ani wczytania pliku TMX, ani zapisu do niego.',
  },
  {
    nazwa: 'Context (101%) Match',
    grupa: PAMIEC,
    opis: 'Odróżnienie dopasowania w identycznym kontekście od zwykłego dopasowania pełnego.',
    okno: OKNA.pamiec,
    warstwa: 4,
    komendy: [],
    brak: 'Odpowiedź podpowiedzi nie niesie metadanych kontekstu ani stopnia dopasowania.',
  },
  {
    nazwa: 'TM Maintenance',
    grupa: PAMIEC,
    opis: 'Czyszczenie duplikatów, scalanie pamięci, masowa podmiana, filtrowanie wg pola.',
    okno: OKNA.pamiec,
    warstwa: 4,
    komendy: [],
    brak:
      'Żadna z czterech operacji nie ma komendy: kontrakt nie wystawia pamięci do odczytu ' +
      'ani do zapisu, jedynym wejściem jest podpowiedź dla segmentu.',
  },
  {
    nazwa: 'Leverage/Pre-translate',
    grupa: PAMIEC,
    opis: 'Wstępne wypełnienie materiału trafieniami z pamięci przed przekładem maszynowym.',
    okno: OKNA.pamiec,
    warstwa: 3,
    komendy: [],
    brak:
      'Komendy wstępnego wypełnienia nie ma. Podpowiedź da się wstawić do panelu ręcznie ' +
      'i zapisać korektą, ale to czynność na jeden panel, nie wypełnienie materiału.',
  },

  {
    nazwa: 'Termbase Manager',
    grupa: TERMINOLOGIA,
    opis: 'Baza terminów z odpowiednikiem per język, uwagą i statusem terminu.',
    okno: OKNA.glosariusz,
    warstwa: 1,
    komendy: [Command.TranslateGlossarySet],
    brak:
      'Zapis i edycja działają; odczytu bazy kontrakt nie ma, więc okno pokazuje wyłącznie ' +
      'terminy zapisane w tej sesji.',
  },
  {
    nazwa: 'TBX/CSV Import/Export',
    grupa: TERMINOLOGIA,
    opis: 'Wymiana bazy terminologicznej w standardzie TBX oraz w formacie CSV.',
    okno: OKNA.glosariusz,
    warstwa: 3,
    komendy: [Command.TranslateGlossaryImport, Command.TranslateGlossaryExport],
  },
  {
    nazwa: 'Term Enforcement',
    grupa: TERMINOLOGIA,
    opis: 'Ujednolicenie terminologii paneli według zatwierdzonych odpowiedników.',
    okno: OKNA.glosariusz,
    warstwa: 3,
    komendy: [Command.TranslateGlossaryApply],
  },
  {
    nazwa: 'Inconsistency Detection',
    grupa: TERMINOLOGIA,
    opis: 'Wskazanie miejsc, w których ten sam termin przetłumaczono w jednym języku różnie.',
    okno: OKNA.glosariusz,
    warstwa: 3,
    komendy: [Command.TranslateGlossaryOccurrences],
    brak:
      'Rdzeń oddaje wystąpienia terminu wraz z otoczeniem; wykazu niespójności nie liczy — ' +
      'ocena rozbieżności zostaje po stronie Operatora.',
  },
  {
    nazwa: 'Do-Not-Translate (DNT)',
    grupa: TERMINOLOGIA,
    opis: 'Oznaczenie nazw własnych, jednostek i kodów jako niepodlegających tłumaczeniu.',
    okno: OKNA.glosariusz,
    warstwa: 1,
    komendy: [Command.TranslateGlossarySet],
  },
  {
    nazwa: 'Term Extraction',
    grupa: TERMINOLOGIA,
    opis: 'Wskazanie kandydatów na terminy z tekstu źródłowego do zasilenia bazy.',
    okno: OKNA.glosariusz,
    warstwa: 3,
    komendy: [],
    brak: 'Kontrakt nie ma komendy wyprowadzającej kandydatów na terminy z tekstu.',
  },
  {
    nazwa: 'Forbidden/Preferred Terms',
    grupa: TERMINOLOGIA,
    opis: 'Listy terminów zabronionych i preferowanych zgodne z przewodnikiem stylu.',
    okno: OKNA.glosariusz,
    warstwa: 3,
    komendy: [],
    brak:
      'Termin ma w kontrakcie jeden znacznik stanu — „nie tłumacz". Statusu zatwierdzony, ' +
      'kandydat ani zabroniony nie ma czym zapisać.',
  },

  {
    nazwa: 'Grammar & Spelling Check',
    grupa: KOREKTA,
    opis: 'Sprawdzanie gramatyki, ortografii i interpunkcji osobno dla każdego języka.',
    okno: OKNA.jakosc,
    warstwa: 3,
    komendy: [],
    brak: 'Kontrola jakości kontraktu zna sześć rodzajów niezgodności; korekty językowej wśród nich nie ma.',
  },
  {
    nazwa: 'Style & Register Check',
    grupa: KOREKTA,
    opis: 'Ocena zgodności rejestru z tonem zamierzonym dla panelu.',
    okno: OKNA.jakosc,
    warstwa: 3,
    komendy: [],
    brak: 'Ton panelu da się ustawić, ale zgodności przekładu z nim kontrakt nie ocenia.',
  },
  {
    nazwa: 'Readability Scoring',
    grupa: KOREKTA,
    opis: 'Wskaźniki czytelności liczone osobno dla każdego języka docelowego.',
    okno: OKNA.jakosc,
    warstwa: 4,
    komendy: [],
    brak: 'Kontrakt nie ma komendy liczącej wskaźniki czytelności.',
  },
  {
    nazwa: 'False-Friends & Interference',
    grupa: KOREKTA,
    opis: 'Wykrywanie kalek, fałszywych przyjaciół i interferencji z języka źródłowego.',
    okno: OKNA.jakosc,
    warstwa: 4,
    komendy: [],
    brak: 'Kontrakt nie ma komendy wykrywającej interferencję językową.',
  },
  {
    nazwa: 'Punctuation & Typography Locale',
    grupa: KOREKTA,
    opis: 'Poprawa cudzysłowów, spacji i myślników zgodnie z konwencją języka docelowego.',
    okno: OKNA.jakosc,
    warstwa: 4,
    komendy: [],
    brak: 'Kontrakt nie ma komendy stosującej reguły typograficzne języka docelowego.',
  },
  {
    nazwa: 'Consistency Checker',
    grupa: KOREKTA,
    opis: 'Wykrywanie niespójnego przekładu identycznych zdań i wariantów jednego terminu.',
    okno: OKNA.jakosc,
    warstwa: 3,
    komendy: [Command.TranslateGlossaryOccurrences],
    brak:
      'Wystąpienia terminu rdzeń pokazuje; niespójności zdań identycznych nie liczy ' +
      'żadna komenda kontraktu.',
  },

  {
    nazwa: 'Numeric/Date/Currency Check',
    grupa: JAKOSC,
    opis: 'Zgodność liczb, dat, walut i jednostek między źródłem a przekładem.',
    okno: OKNA.jakosc,
    warstwa: 1,
    komendy: [Command.TranslateQualityCheck],
  },
  {
    nazwa: 'Placeholder/Tag Integrity',
    grupa: JAKOSC,
    opis: 'Kompletność i kolejność symboli zastępczych oraz znaczników formatu.',
    okno: OKNA.jakosc,
    warstwa: 1,
    komendy: [Command.TranslateQualityCheck],
  },
  {
    nazwa: 'Length/Fit Check',
    grupa: JAKOSC,
    opis: 'Kontrola limitu znaków i sygnalizacja przekroczeń długości.',
    okno: OKNA.jakosc,
    warstwa: 1,
    komendy: [Command.TranslateQualityCheck],
    brak: 'Profilu limitów żądanie nie niesie — próg długości ustala rdzeń, nie okno.',
  },
  {
    nazwa: 'Omission/Untranslated Detect',
    grupa: JAKOSC,
    opis: 'Wykrycie segmentów pominiętych, nieprzetłumaczonych albo tożsamych ze źródłem.',
    okno: OKNA.jakosc,
    warstwa: 1,
    komendy: [Command.TranslateQualityCheck],
  },
  {
    nazwa: 'LQA Scorecard (MQM/DQF)',
    grupa: JAKOSC,
    opis: 'Ocena jakości wedle modelu błędów MQM/DQF: kategoria, waga, wynik końcowy.',
    okno: OKNA.jakosc,
    warstwa: 4,
    komendy: [],
    brak:
      'Kontrakt nie ma komendy zapisu oceny. Kartę wypełnia się i przelicza w oknie, ' +
      'a wynik zostaje w nim — rdzeń go nie przechowuje i nie wyda go w raporcie.',
  },
  {
    nazwa: 'QA Profiles',
    grupa: JAKOSC,
    opis: 'Zestawy reguł kontroli włączane per projekt albo klient.',
    okno: OKNA.jakosc,
    warstwa: 2,
    komendy: [],
    brak: 'Żądanie kontroli niesie sam panel; pola na profil reguł w nim nie ma.',
  },

  {
    nazwa: 'DOCX Round-Trip',
    grupa: FORMATY,
    opis: 'Wydobycie tekstu z dokumentu, przekład i wydanie dokumentu wynikowego.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [Command.DocumentTextExtract, Command.TranslatePanelExport],
    brak:
      'Wydobycie tekstu i wydanie pliku działają osobno; odtworzenia stylów, tabel ' +
      'i osadzeń dokumentu wejściowego kontrakt nie obiecuje.',
  },
  {
    nazwa: 'PDF Extract & Reflow',
    grupa: FORMATY,
    opis: 'Wydobycie tekstu i układu z dokumentu PDF, przekład i złożenie wyniku.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [Command.DocumentTextExtract, Command.TranslatePanelExport],
    brak: 'Odpowiedź niesie sam tekst i liczbę stron — układu strony nie oddaje.',
  },
  {
    nazwa: 'OCR for Scanned Docs',
    grupa: FORMATY,
    opis: 'Rozpoznanie pisma ze skanów i obrazów przed przekładem.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [Command.DocumentTextExtract],
  },
  {
    nazwa: 'PPTX/XLSX Translate',
    grupa: FORMATY,
    opis: 'Przekład prezentacji i arkuszy z zachowaniem slajdów, notatek i komórek.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [],
    brak:
      'Zamiana formatów kontraktu wymienia markdown, html, docx, odt, pdf, epub, rtf i csv — ' +
      'prezentacji ani arkusza wśród nich nie ma.',
  },
  {
    nazwa: 'Markdown/HTML Translate',
    grupa: FORMATY,
    opis: 'Przekład treści z ochroną składni, odnośników i bloków kodu.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [Command.DocumentConvert, Command.TranslatePanelExport],
  },
  {
    nazwa: 'ODT/ODF Translate',
    grupa: FORMATY,
    opis: 'Przekład dokumentów OpenDocument z zachowaniem stylów.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [Command.DocumentConvert],
    brak: 'Zamiana formatu obejmuje odt; zachowania stylów przy przekładzie kontrakt nie obiecuje.',
  },
  {
    nazwa: 'Layout Diff / Visual Compare',
    grupa: FORMATY,
    opis: 'Porównanie wyglądu dokumentu przed i po przekładzie: przepełnienia, złamany układ.',
    okno: OKNA.formaty,
    warstwa: 3,
    komendy: [],
    brak: 'Kontrakt nie ma komendy renderującej dokument ani porównującej jego układ.',
  },
  {
    nazwa: 'Bilingual Export',
    grupa: FORMATY,
    opis: 'Wydanie pliku dwujęzycznego: źródło obok przekładu.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [],
    brak:
      'Eksport panelu przyjmuje pięć formatów kontraktu i wydaje sam przekład. Formatu XLIFF ' +
      'ani wariantu dwujęzycznego w wykazie nie ma.',
  },

  {
    nazwa: 'Localization File Formats',
    grupa: LOKALIZACJA,
    opis: 'Wczytanie i wydanie zasobów lokalizacyjnych z ochroną kluczy i zmiennych.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [],
    brak: 'Kontrakt nie zna formatów zasobów lokalizacyjnych ani pojęcia klucza zasobu.',
  },
  {
    nazwa: 'XLIFF 1.2 / 2.1 Round-Trip',
    grupa: LOKALIZACJA,
    opis: 'Pełen obieg standardu wymiany lokalizacyjnej wraz z metadanymi stanu i notatkami.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [],
    brak: 'XLIFF nie jest wartością formatu eksportu i nie ma komendy przyjmującej plik XLIFF.',
  },
  {
    nazwa: 'Pseudolocalization',
    grupa: LOKALIZACJA,
    opis: 'Pseudotłumaczenie do testu interfejsu: wydłużenie, znaki diakrytyczne, ramki.',
    okno: OKNA.formaty,
    warstwa: 4,
    komendy: [],
    brak:
      'Kontrakt nie ma komendy pseudolokalizacji. Przekształcenie liczy okno i pokazuje wynik ' +
      'do przeniesienia ręcznego — w rdzeniu nic się przy tym nie zapisuje.',
  },
  {
    nazwa: 'Plural & Gender Rules',
    grupa: LOKALIZACJA,
    opis: 'Obsługa form liczby mnogiej i rodzaju wedle reguł standardu CLDR.',
    okno: OKNA.formaty,
    warstwa: 4,
    komendy: [],
    brak: 'Kontrakt nie ma komendy stosującej reguły liczby mnogiej ani rodzaju.',
  },
  {
    nazwa: 'Placeholder Style Mapping',
    grupa: LOKALIZACJA,
    opis: 'Rozpoznanie i przekład stylu zmiennych między platformami.',
    okno: OKNA.panele,
    warstwa: 4,
    komendy: [],
    brak:
      'Kontrakt nie ma komendy zamiany stylu zmiennych. Mapowanie liczy okno na tekście ' +
      'źródłowym i oddaje wynik do przeniesienia ręcznego.',
  },
  {
    nazwa: 'Key Context & Screenshots',
    grupa: LOKALIZACJA,
    opis: 'Powiązanie klucza lokalizacyjnego z kontekstem i zrzutem ekranu miejsca użycia.',
    okno: OKNA.formaty,
    warstwa: 4,
    komendy: [],
    brak: 'Kontrakt nie zna klucza lokalizacyjnego, więc nie ma czego wiązać z kontekstem.',
  },

  {
    nazwa: 'Subtitle Formats',
    grupa: MOWA,
    opis: 'Wczytanie, przekład i wydanie napisów z zachowaniem taktowania.',
    okno: OKNA.formaty,
    warstwa: 2,
    komendy: [],
    brak: 'Formatów napisowych nie ma ani w wykazie eksportu, ani w zamianie formatów dokumentu.',
  },
  {
    nazwa: 'Timing & CPS Control',
    grupa: MOWA,
    opis: 'Kontrola liczby znaków na sekundę, długości linii i czasu wyświetlania.',
    okno: OKNA.jakosc,
    warstwa: 4,
    komendy: [],
    brak: 'Kontrakt nie zna taktowania napisów, więc nie ma czego kontrolować.',
  },
  {
    nazwa: 'Text-to-Speech Readback',
    grupa: MOWA,
    opis: 'Odsłuch przekładu syntezą mowy jako kontrola brzmieniowa.',
    okno: OKNA.panele,
    warstwa: 3,
    komendy: [Command.TranslateSpeechSynthesize],
  },
  {
    nazwa: 'Speech-to-Text Dictation',
    grupa: MOWA,
    opis: 'Dyktowanie przekładu głosem oraz transkrypcja nagrań źródłowych.',
    okno: OKNA.chat,
    warstwa: 3,
    komendy: [Command.SpeechTranscribe],
    brak:
      'Komenda transkrypcji jest w kontrakcie, ale nagranie dźwięku prowadzi okno komunikacji ' +
      'platformy — okna modułu Translate mikrofonu nie obsługują.',
  },
  {
    nazwa: 'Dubbing/Voiceover Script',
    grupa: MOWA,
    opis: 'Skrypt dubbingu z segmentacją kwestii, oznaczeniem mówców i długością pod lip-sync.',
    okno: OKNA.formaty,
    warstwa: 4,
    komendy: [],
    brak: 'Kontrakt nie zna kwestii dialogowej ani mówcy, więc skryptu nie ma z czego złożyć.',
  },

  {
    nazwa: 'Vendor Handoff Package',
    grupa: WYDANIE,
    opis: 'Pakiet dla wykonawcy zewnętrznego: materiał, pamięć, baza terminów i instrukcje.',
    okno: OKNA.jakosc,
    warstwa: 3,
    komendy: [Command.TranslateGlossaryExport, Command.TranslatePanelExport],
    brak:
      'Baza terminów i panele wychodzą osobnymi komendami. Pamięci tłumaczeń nie ma czym ' +
      'wyeksportować, a spakowania całości w jeden pakiet kontrakt nie przewiduje.',
  },
  {
    nazwa: 'Review & Approval Workflow',
    grupa: WYDANIE,
    opis: 'Przebieg tłumaczenie → korekta → zatwierdzenie ze śladem autora zmiany.',
    okno: OKNA.jakosc,
    warstwa: 1,
    komendy: [],
    brak:
      'Panel ma w kontrakcie cztery stany wykonania: oczekuje, tłumaczenie w toku, gotowe, błąd. ' +
      'Etapu akceptacji ani autora zmiany nie ma w nim czym zapisać.',
  },
  {
    nazwa: 'Word/Char Count & Cost Estimate',
    grupa: WYDANIE,
    opis: 'Statystyka objętości z analizą względem pamięci i wyceną wedle stawek.',
    okno: OKNA.zrodlo,
    warstwa: 2,
    komendy: [],
    brak:
      'Słowa, znaki i segmenty liczy okno z tekstu, który ma przed sobą. Analizy względem ' +
      'pamięci ani siatki stawek kontrakt nie niesie, więc wyceny nie ma z czego złożyć.',
  },
  {
    nazwa: 'Studio Bridge',
    grupa: WYDANIE,
    opis: 'Odbiór zaznaczenia z modułu Studio jako tekstu źródłowego i zwrot wersji wielojęzycznej.',
    okno: OKNA.zrodlo,
    warstwa: 4,
    komendy: [],
    brak:
      'Powiązanie pary modułów ustanawia okno konfiguracji, nie moduł. Komendy przekazującej ' +
      'zaznaczenie między modułami w obszarze translate nie ma.',
  },
  {
    nazwa: 'Library Sync',
    grupa: WYDANIE,
    opis: 'Zapis wydanych plików dwujęzycznych, pamięci i bazy terminów jako artefaktów.',
    okno: OKNA.jakosc,
    warstwa: 4,
    komendy: [],
    brak:
      'Eksport panelu oddaje ścieżkę pliku, nie zasób magazynu — powiązania z modułem Library ' +
      'nie ma czym ustanowić z tego okna.',
  },
  {
    nazwa: 'Automations Steps',
    grupa: WYDANIE,
    opis: 'Udostępnienie operacji modułu jako kroków procesu wsadowego.',
    okno: OKNA.petla,
    warstwa: 4,
    komendy: [],
    brak:
      'Kroki procesu definiuje moduł Automations własnymi komendami; obszar translate nie ma ' +
      'komendy zgłaszającej swoje operacje jako kroki.',
  },
  {
    nazwa: 'Translation API/Webhook',
    grupa: WYDANIE,
    opis: 'Wystawienie przekładu i kontroli jakości jako operacji wywoływanych z zewnątrz.',
    okno: OKNA.jakosc,
    warstwa: 4,
    komendy: [Command.SessionToolAttach, Command.ToolsCatalogList],
    brak:
      'Operacje modułu są zadeklarowane w kontrakcie jako narzędzia modelu i dokładają się ' +
      'do sesji. Webhooka wywoływanego z zewnątrz kontrakt nie wystawia.',
  },
];

/**
 * Grupy katalogu w kolejności opracowania, bez powtórzeń — kolejność, w jakiej
 * wyszukiwarka funkcji i okna modułu układają pozycje na widoku.
 */
export function grupyKatalogu(): readonly string[] {
  const widziane = new Set<string>();
  for (const pozycja of KATALOG_FUNKCJI) widziane.add(pozycja.grupa);
  return [...widziane];
}

/**
 * Ile pozycji katalogu ma komendę kontraktu — liczba pozycji rzeczywiście
 * wykonywanych, a nie zadeklarowana wielkość modułu.
 */
export function liczbaZKomenda(): number {
  return KATALOG_FUNKCJI.filter((pozycja) => pozycja.komendy.length > 0).length;
}
