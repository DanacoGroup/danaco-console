import { Command } from '../../../../shared/contract';
import type { AkcjaOkna } from './panel-akcji';

/**
 * Katalog akcji siedmiu okien Research.
 *
 * Jedna odpowiedzialność: nazwanie akcji i wskazanie drogi ich wykonania. Droga
 * stoi tu, a nie w oknie, bo rozstrzyga o odbiorze: akcja bez wykonawcy idzie
 * generycznym `window.action`, zamiast znikać z paska albo udawać wykonanie.
 *
 * Nazwa komendy NIGDY nie jest tu napisem. Kod akcji, która ma w kontrakcie
 * własną komendę, bierze się ze stałej `Command.*` — inaczej zmiana nazwy
 * w `contract.json` zostawiłaby w panelu martwe wywołanie, którego kompilator
 * nie wychwyci. Zaporę na to niesie `src/kontrakt.test.ts`.
 *
 * Wykaz nie jest kopią rejestru akcji rdzenia. Rejestr (`action.list`) opisuje
 * akcje modułu w bazie, a zaczyn katalogu akcji
 * (`budowa/server/internal/store/migracja_009_zaczyn_akcji.sql`) zakłada dla
 * modułu Research wyłącznie trzy pozycje paska promptu — pozycje wiadomości —
 * bo powstają złączeniem z tabelą `modul`. Pozostałe kody z tego pliku wiersza
 * w katalogu nie mają, więc panel zbudowany wyłącznie z rejestru byłby pusty.
 * Wykaz znika, gdy rejestr odda te same pozycje.
 */

/** Czym akcja się dziś kończy. */
export type DrogaAkcji =
  /** Okno wykonuje ją u siebie albo komendą, której rdzeń już słucha. */
  | 'okno'
  /**
   * Zdolność ma w kontrakcie własną komendę, którą rdzeń obsługuje. Kod akcji
   * jest nazwą tej komendy, wziętą ze stałej `Command.*`, a wywołanie składa
   * `wywolania-komend.ts`.
   */
  | 'komenda'
  /**
   * Zdolności nie odpowiada osobna komenda kontraktu — jest nastawą pola
   * komendy istniejącej albo należy do okna konfiguracji. Wychodzi generycznym
   * `window.action`; rdzeń tę komendę zna, więc odmowa jest merytoryczna:
   * `not_found` z katalogu akcji albo `conflict` braku wykonawcy
   * (`budowa/server/internal/core/adapter_okno_akcja.go`).
   */
  | 'akcja';

export interface AkcjaBadania extends AkcjaOkna {
  droga: DrogaAkcji;
}

function okno(kod: string, nazwa: string, objasnienie: string): AkcjaBadania {
  return { kod, nazwa, droga: 'okno', objasnienie };
}

function akcja(kod: string, nazwa: string, objasnienie: string): AkcjaBadania {
  return { kod, nazwa, droga: 'akcja', objasnienie };
}

/**
 * Akcja, której zdolność ma już komendę w kontrakcie.
 *
 * Kod bierze się z `Command.*`, więc jest dokładnie nazwą komendy docelowej —
 * a katalog akcji rdzenia jest właśnie odwzorowaniem kodu akcji na komendę.
 */
function doKomendy(komenda: Command, nazwa: string, czynnosc: string): AkcjaBadania {
  return { kod: komenda, nazwa, droga: 'komenda', objasnienie: `${czynnosc} ${BEZ_UCHWYTU}` };
}

/**
 * Zdanie o akcji, którą wykonuje komenda kontraktu.
 *
 * Rozróżnienie wobec `BEZ_KOMENDY` nie jest odcieniem: tam osobnej komendy nie
 * ma i mieć nie musi, więc pozycja wychodzi generycznym `window.action`. Tu
 * komenda istnieje, rdzeń ma jej uchwyt, a okno wywołuje ją wprost — i odmowa,
 * którą Operator zobaczy, mówi o jego badaniu, nie o brakach rdzenia.
 */
const BEZ_UCHWYTU =
  'Komenda jest w kontrakcie i rdzeń ma dla niej uchwyt — kod tej akcji jest jej nazwą, ' +
  'a wywołanie składa wywolania-komend.ts z okna badania, zaznaczenia i wskazań okien. ' +
  'Odmowa, która tu przyjdzie, jest odmową merytoryczną rdzenia (brak wskazania, brak ' +
  'materiału, dostawca bez odpowiedzi), nie brakiem uchwytu.';

/**
 * Zdanie o akcji, której osobnej komendy w kontrakcie nie ma.
 *
 * Zostaje ono przy pozycjach, które komendą nie są i nie mają nią być: nastawy
 * będące polem komendy istniejącej oraz ustawienia należące do okna
 * konfiguracji. Treść ma się zgadzać ze słowami odmowy rdzenia — zapowiedź
 * braku wykonawcy przy odmowie płynącej z katalogu akcji prowadziłaby w złą
 * stronę przy szukaniu przyczyny.
 */
const BEZ_KOMENDY =
  'Osobnej komendy kontrakt dla tej pozycji nie ma i mieć nie musi, więc wychodzi ona ' +
  'generycznym window.action. Rdzeń odmówi na pierwszym kroku — „akcja … nie istnieje ' +
  'w katalogu akcji" — bo zaczyn katalogu nie ma wiersza dla tego kodu.';

/** Nastawa będąca polem komendy istniejącej — wskazanie, którym polem. */
function poleKomendy(komenda: Command, pole: string): string {
  return (
    `Ta pozycja jest nastawą, nie osobną czynnością: niesie ją pole ${pole} komendy ` +
    `${komenda}. Dobudowa polega na wystawieniu pola w oknie, nie na nowej komendzie. ` +
    `${BEZ_KOMENDY}`
  );
}

/**
 * Zdanie o oknie, którego katalog okien rdzenia jeszcze nie zna.
 *
 * Discovery Panel i Reading View mają wiersz w opracowaniu modułu, a nie mają go
 * w migracjach 030 i 031. Odmowa przyjdzie więc o krok wcześniej niż przy akcji
 * bez wiersza w katalogu akcji — na samym oknie — i akcja ma to zapowiadać,
 * zamiast obiecywać drogę, której dziś nie ma nawet do połowy.
 */
const OKNO_SPOZA_KATALOGU =
  'Uwaga: katalog okien rdzenia nie zna jeszcze tego okna (migracje 030 i 031 zakładają dla ' +
  'modułu Research pięć okien, a opracowanie wylicza siedem), więc odmowa przyjdzie na oknie, ' +
  'zanim rdzeń dojdzie do katalogu akcji.';

/** Research Workspace — okno wiodące, punkt wejścia modułu. */
export const AKCJE_WORKSPACE: readonly AkcjaBadania[] = [
  okno('research.scope.edit', 'Edytuj zakres', 'Przenosi ognisko do pola zakresu badania. Zapis idzie komendą research.workspace.set, której rdzeń słucha.'),
  okno('research.scope.investigate', 'Zbadaj →', 'Zapisuje zakres i przenosi ognisko do Discovery Panel, żeby wyszukać materiał.'),
  doKomendy(Command.ResearchWorkspaceGet, 'Odczytaj przestrzeń badania z rdzenia', 'Zakres, etapy, pytania badawcze, odbiorca raportu, protokół i notatka robocza — w postaci, w której rdzeń je trzyma. Zakres i etapy wchodzą do pól tego okna.'),
  doKomendy(Command.ResearchWorkspaceQuestionSet, 'Pytania badawcze', 'Wykaz pytań badawczych badania.'),
  doKomendy(Command.ResearchWorkspaceCoverage, 'Pokrycie pytań', 'Które pytania mają poparcie w źródłach, a które nie.'),
  doKomendy(Command.ResearchWorkspaceNoteSet, 'Notatka robocza', 'Hipotezy i pytania otwarte niezwiązane z konkretnym źródłem.'),
  doKomendy(Command.ResearchGapFind, 'Znajdź luki badawcze', 'Obszary tematu niepokryte dotychczasowymi źródłami.'),
  doKomendy(Command.ResearchWorkspaceFreshness, 'Wskaźnik świeżości', 'Data ostatniego pozyskania materiału wraz z sugestią odświeżenia.'),
  doKomendy(Command.ResearchEvidenceGraph, 'Evidence Map', 'Graf źródeł, ustaleń i twierdzeń wraz z relacjami poparcia i sprzeczności.'),
  doKomendy(Command.ResearchPrismaGet, 'Diagram PRISMA', 'Liczniki przesiewu: zidentyfikowane, przesiane, włączone.'),
  akcja('research.scope.timeline', 'Oś czasu badania', `Chronologia dodawania źródeł i ustaleń. Jest reprezentacją materiału, który okno już ma — składa się z pól createdAt i acquiredAt, bez pytania rdzenia. ${BEZ_KOMENDY}`),
  akcja('research.scope.kanban', 'Kanban etapów', `Tablica etapów dla pracy zespołowej. Etapy niesie już odpowiedź zapisu zakresu; brakuje stanu etapu, a ten jest polem, nie komendą. ${BEZ_KOMENDY}`),
  akcja('research.scope.evidenceTable', 'Tabela dowodów', `Zestawienie ustaleń ze źródłami i oceną. Jako blok raportu niesie ją pole kind komendy research.report.insert; jako widok przestrzeni składa się z materiału, który okno ma. ${BEZ_KOMENDY}`),
];

/** Discovery Panel — wyszukiwanie i odkrywanie źródeł. */
export const AKCJE_ODKRYWANIA: readonly AkcjaBadania[] = [
  okno('research.discovery.run', 'Szukaj', 'Uruchamia zapytanie w trybie bieżącym. Tryb semantyczny idzie komendą knowledge.search, pełnotekstowy — library.file.search; obu rdzeń słucha.'),
  okno('research.discovery.toSources', '→ Sources Manager', 'Przenosi ognisko do katalogu źródeł; pozycje wstawia przycisk „→ Dodaj do źródeł" przy wyniku, komendą research.source.add.'),
  doKomendy(Command.ResearchDiscoverySearch, 'Wyszukiwanie webowe i naukowe', `Zapytanie do wyszukiwarki sieciowej albo baz publikacji; tryb rozstrzyga pole mode. ${OKNO_SPOZA_KATALOGU}`),
  doKomendy(Command.ResearchDiscoveryAssist, 'Query Assistant', 'Przekształcenie pytania badawczego w zapytania z operatorami.'),
  doKomendy(Command.ResearchSourceResolve, 'Rozstrzygnij identyfikator', 'Pełne metadane pozycji po numerze DOI, ISBN, PMID albo arXiv.'),
  doKomendy(Command.ResearchDiscoverySnowball, 'Cytowania', 'Prace cytowane i cytujące wybraną pozycję.'),
  doKomendy(Command.ResearchDiscoveryReject, 'Odrzuć z uzasadnieniem', 'Oznaczenie pozycji jako odrzuconej; zasila liczniki diagramu PRISMA.'),
  doKomendy(Command.ResearchMonitorList, 'Monitory tematów i RSS', 'Zapisane tematy i kanały wraz ze skrzynką „nowe źródła".'),
  doKomendy(Command.ResearchBatchImport, 'Import wsadowy adresów', 'Zbiorcze pozyskanie wielu źródeł jednym zleceniem, w tle.'),
  akcja('research.discovery.filters', 'Filtry wyników', poleKomendy(Command.ResearchDiscoverySearch, 'yearFrom, yearTo, field, openAccessOnly, itemType i providers')),
  akcja('research.discovery.providers', 'Dostawcy wyszukiwania', `Zestaw dostawców webowych i naukowych oraz ich klucze dostępu. Należy do okna konfiguracji, nie do panelu badania. ${BEZ_KOMENDY}`),
];

/** Sources Manager — zarządca źródeł badania. */
export const AKCJE_ZRODEL: readonly AkcjaBadania[] = [
  okno('research.source.new', '+ Dodaj źródło', 'Przenosi ognisko do formularza katalogowania źródła.'),
  okno('research.source.link', 'Powiąż z ustaleniem', 'Przenosi zaznaczone źródła do formularza ustalenia w Findings Panel; powiązanie zapisuje research.finding.add.'),
  okno('research.source.read', 'Czytaj', 'Otwiera zaznaczone źródło w Reading View. Treść dochodzi komendą library.file.preview, której rdzeń słucha.'),
  doKomendy(Command.ResearchSourceList, 'Odczytaj wykaz z rdzenia', 'Źródła skatalogowane w oknie badania. Do czasu dobudowy wykaz narasta wyłącznie z odpowiedzi tej sesji.'),
  doKomendy(Command.ResearchSourceRemove, 'Usuń z listy', 'Usunięcie źródła wraz z jego powiązaniami.'),
  doKomendy(Command.ResearchSourceDuplicates, 'Duplikaty', 'Wskazanie duplikatów i niemal-duplikatów w katalogu.'),
  doKomendy(Command.ResearchSourceMerge, 'Scal duplikat', 'Złączenie źródeł powtórzonych w jedno, z przeniesieniem powiązań.'),
  doKomendy(Command.ResearchSourceTag, 'Etykiety i kolekcje', 'Katalogowanie wielowymiarowe: tagi tematyczne i kolekcje.'),
  doKomendy(Command.ResearchSourceUpdate, 'Stan lektury i metadane', 'Zmiana metadanych, oceny wiarygodności i stanu lektury źródła.'),
  doKomendy(Command.ResearchSourceAttachmentList, 'Załączniki', 'Pełny tekst, migawka i notatki źródła wraz z kontrolą kompletności.'),
  doKomendy(Command.ResearchSourceAttachmentAdd, 'Dołącz pełny tekst', 'Załącznik pełnego tekstu przy zaznaczonym źródle: plik wskazany w polu tytułu formularza albo dokument repozytorium, który źródło już wskazuje.'),
  doKomendy(Command.ResearchSourceImport, 'Import zbiorczy referencji', 'Wczytanie bibliografii z pliku BibTeX, RIS, CSL-JSON, EndNote XML albo CSV.'),
  doKomendy(Command.ResearchSourceCapture, 'Zapisz stronę i migawkę', 'Web-clipping wraz z niezmienną migawką strony.'),
  doKomendy(Command.ResearchSourceTranscribe, 'Transkrybuj nagranie', 'Zamiana nagrania w cytowalny transkrypt ze znacznikami czasu.'),
  doKomendy(Command.ResearchCitationRender, 'Cytuj i eksportuj bibliografię', 'Cytat w tekście i pozycja bibliograficzna w wybranym stylu CSL.'),
  doKomendy(Command.ResearchCitationStyles, 'Styl cytatu', 'Style dostępne w repozytorium CSL wraz z wariantami własnymi.'),
  doKomendy(Command.ResearchCitationCheck, 'Kontrola cytowań', 'Kompletność metadanych, źródła niecytowane i cytowania bez wpisu.'),
  doKomendy(Command.ResearchRetractionCheck, 'Sprawdź wycofania', 'Czy cytowana praca została wycofana albo skorygowana.'),
  doKomendy(Command.ResearchReadingOpen, 'Podgląd', 'Treść źródła bez opuszczania katalogu.'),
];

/** Reading View — lektura materiału i wypisy z niego. */
export const AKCJE_LEKTURY: readonly AkcjaBadania[] = [
  okno('research.reading.toFinding', '→ ustalenie', 'Zapisuje zaznaczony fragment jako ustalenie z cytatem i powiązaniem do czytanego źródła — komenda research.finding.add, której rdzeń słucha.'),
  okno('research.reading.reload', 'Wczytaj ponownie', 'Powtarza odczyt bieżącej strony materiału komendą library.file.preview.'),
  okno('research.reading.toFindings', '→ Findings Panel', 'Przenosi ognisko do wykazu ustaleń; wypisy z tego okna stoją tam razem z ustaleniami wpisanymi ręcznie.'),
  doKomendy(Command.ResearchAnnotationAdd, 'Podświetl i notuj', `Trwałe podświetlenie albo notatka na marginesie, z kotwicą pozycji. ${OKNO_SPOZA_KATALOGU}`),
  doKomendy(Command.ResearchAnnotationList, 'Adnotacje źródła', 'Podświetlenia, notatki i zakładki zapisane przy materiale.'),
  doKomendy(Command.ResearchAnnotationRemove, 'Zdejmij adnotację', 'Zdjęcie najnowszej adnotacji czytanego materiału. CZYNNOŚĆ NIEODWRACALNA: pierwsze naciśnięcie nazywa adnotację i ostrzega, dopiero drugie ją zdejmuje — rdzeń nie ma komendy przywrócenia.'),
  doKomendy(Command.ResearchExcerptList, 'Wypisy', 'Zbiorczy przegląd podświetleń i notatek z jednego lub wielu źródeł.'),
  doKomendy(Command.ResearchSourceSummarize, 'Streść źródło', 'Abstrakt roboczy, tezy, metodologia i wnioski pojedynczego źródła.'),
  doKomendy(Command.ResearchSourceExtractTable, 'Wyodrębnij tabelę', 'Tabele i dane liczbowe do postaci ustrukturyzowanej.'),
  doKomendy(Command.ResearchSourceExtractClaims, 'Wydobądź twierdzenia', 'Kluczowe twierdzenia, dane liczbowe i podmioty ze wskazaniem fragmentu.'),
  doKomendy(Command.ResearchSourceOcr, 'OCR skanu', 'Rozpoznanie tekstu w skanie bez warstwy tekstowej.'),
  doKomendy(Command.ResearchCorpusAsk, 'Zapytaj o korpus', 'Rozmowa oparta na treści wybranych źródeł, zakotwiczona w cytatach.'),
  doKomendy(Command.ResearchSourceUpdate, 'Stan lektury', 'Oznaczenie materiału jako przeczytanego, przejrzanego albo pominiętego.'),
  akcja('research.reading.find', 'Szukaj w treści', `Wyszukiwanie w treści materiału z podświetleniem trafień. Treść okno już ma — szukanie odbywa się w niej, bez pytania rdzenia. ${BEZ_KOMENDY}`),
];

/** Findings Panel — ustalenia narastające w toku badania. */
export const AKCJE_USTALEN: readonly AkcjaBadania[] = [
  okno('research.finding.new', '+ Nowe ustalenie', 'Czyści formularz i przenosi do niego ognisko; zapis idzie research.finding.add bez pola findingId.'),
  okno('research.finding.edit', 'Edytuj', 'Wciąga zaznaczone ustalenie do formularza; zapis idzie research.finding.add z polem findingId.'),
  okno('research.finding.link', 'Powiąż źródło', 'Przenosi ognisko do wyboru źródeł; powiązanie jedzie w polu sourceIds tej samej komendy.'),
  okno('research.finding.toReport', '→ Report Builder', 'Przekazuje zaznaczone ustalenia do kreatora raportu (pole findingIds komendy research.report.build).'),
  okno('research.finding.resolve', 'Oznacz jako rozstrzygnięte', 'Zapisuje ustalenie ze stanem resolved — pole status komendy research.finding.add.'),
  doKomendy(Command.ResearchFindingList, 'Odczytaj wykaz z rdzenia', 'Ustalenia zapisane w oknie badania wraz z wyszukiwaniem pełnotekstowym w ich treści.'),
  doKomendy(Command.ResearchFindingUpdate, 'Klasyfikuj i waż', 'Rodzaj ustalenia (fakt, hipoteza, opinia, dana, cytat) oraz waga kluczowa albo poboczna.'),
  doKomendy(Command.ResearchFindingCode, 'Koduj jakościowo', 'Kody tematyczne przypisywane ustaleniu.'),
  doKomendy(Command.ResearchCodebookGet, 'Książka kodów', 'Kody badania wraz z ich definicjami i licznością wystąpień.'),
  doKomendy(Command.ResearchCodebookSet, 'Zapisz książkę kodów', 'Kody z pola treści — jeden w wierszu, po znaku | definicja kodu. Zapis idzie po odczycie i jest złączeniem: kod zastany, którego teraz nie wpisano, zostaje w rdzeniu.'),
  doKomendy(Command.ResearchFindingMatrix, 'Macierz kod × źródło', 'Rozkład kodów w źródłach.'),
  doKomendy(Command.ResearchFindingContradictions, 'Wykryj sprzeczności', 'Rozbieżności między ustaleniami z różnych źródeł.'),
  doKomendy(Command.ResearchContradictionResolve, 'Rozstrzygnij sprzeczność', 'Oznaczenie sprzeczności jako rozstrzygniętej wraz z uzasadnieniem.'),
  doKomendy(Command.ResearchFindingFactCheck, 'Fact-Check', 'Weryfikacja twierdzenia: poparcie, brak poparcia albo sprzeczność.'),
  doKomendy(Command.ResearchFindingCluster, 'Grupuj w wątki', 'Klastrowanie ustaleń w wątki tematyczne.'),
  doKomendy(Command.ResearchFindingMerge, 'Scal powtórzone', 'Złączenie powtórzonych ustaleń w jedno z wieloma odnośnikami.'),
  doKomendy(Command.ResearchFindingProvenance, 'Prowenancja', 'Ślad pochodzenia: źródło, fragment, sprawca i czas.'),
  doKomendy(Command.ResearchFindingRemove, 'Usuń ustalenie', 'Usunięcie ustalenia wraz z jego powiązaniami.'),
  akcja('research.finding.compare', 'Porównaj źródła', `Zestawienie źródeł ustalenia obok siebie. Powstaje z materiału, który okno ma — powiązania niesie już samo ustalenie. ${BEZ_KOMENDY}`),
];

/** Report Builder — kreator dokumentu końcowego. */
export const AKCJE_RAPORTU: readonly AkcjaBadania[] = [
  okno('research.report.section.add', '+ Dodaj sekcję', 'Dokłada pustą sekcję do redakcji; sekcje jadą w polu sections komendy research.report.build.'),
  okno('research.report.toExport', '→ Export Panel', 'Przenosi ognisko do Export Panel. Panel jest czynny zawsze; brak raportu nazywa komunikatem, nie blokadą.'),
  doKomendy(Command.ResearchReportGet, 'Odczytaj raport z rdzenia', 'Raport okna wraz z sekcjami. Do czasu dobudowy raport dochodzi odpowiedzią budowy albo zdarzeniem zmiany.'),
  doKomendy(Command.ResearchReportTemplateList, 'Szablon raportu', 'Wzorce struktury: analiza konkurencyjna, SWOT, benchmark, nota badawcza, przegląd literatury.'),
  doKomendy(Command.ResearchReportSummarize, 'Streszczenie zarządcze', 'Synteza z ustaleń kluczowych, wstawiana na początek dokumentu.'),
  doKomendy(Command.ResearchReportBibliography, 'Bibliografia końcowa', 'Lista literatury w wybranym stylu CSL, wszystkie źródła albo tylko cytowane.'),
  doKomendy(Command.ResearchReportFootnoteSet, 'Menedżer przypisów', 'Przypisy dolne i końcowe, numeracja i przypisy skrócone.'),
  doKomendy(Command.ResearchReportInsert, 'Wstaw zestawienie', 'Macierz porównawcza, oś czasu, wykres albo tabela dowodów w sekcji.'),
  doKomendy(Command.ResearchReportVersionList, 'Wersje raportu', 'Kolejne kompletacje dokumentu.'),
  doKomendy(Command.ResearchReportDiff, 'Porównaj wersje', 'Zestawienie różnicowe treści dwóch wersji raportu.'),
  doKomendy(Command.ResearchReportCommentList, 'Tryb recenzji', 'Komentarze do fragmentów raportu wraz z wątkami dyskusji.'),
  doKomendy(Command.ResearchReportContextualOp, 'Operacje na zaznaczeniu', 'Korekta, streszczenie, zmiana stylu i rozwinięcie zaznaczonego fragmentu.'),
];

/** Export Panel — wydanie raportu w formacie dokumentowym. */
export const AKCJE_EKSPORTU: readonly AkcjaBadania[] = [
  okno('research.export.run', 'Eksportuj teraz', 'Wydaje raport komendą research.report.export w wybranym formacie i miejscu docelowym; rdzeń jej słucha.'),
  okno('research.export.again', 'Pobierz ponownie', 'Powtarza ostatni eksport tą samą komendą; powtórzenie jest bezpieczne po stronie logiki, nie przez blokadę kontrolki.'),
  doKomendy(Command.ResearchExportPreview, 'Podgląd przed eksportem', 'Dokument w formacie docelowym przed jego wygenerowaniem.'),
  doKomendy(Command.ResearchExportList, 'Historia eksportów rdzenia', 'Ślady wykonanych wydań. Historia w oknie obejmuje dziś wyłącznie bieżącą sesję.'),
  doKomendy(Command.ResearchExportTemplateList, 'Szablony eksportu', 'Zapisane kombinacje formatu, miejsca docelowego i składu dokumentu.'),
  doKomendy(Command.ResearchExportTemplateSet, 'Zapisz szablon eksportu', 'Utrwalenie bieżących nastaw do ponownego użycia.'),
  doKomendy(Command.ResearchExportShare, 'Udostępnij odnośnik', 'Odnośnik do raportu przekazywany odbiorcy bez pobierania pliku.'),
  akcja('research.export.content', 'Zawartość dokumentu', poleKomendy(Command.ResearchReportExport, 'content')),
  akcja('research.export.target', 'Cel: Studio i Roundtable', poleKomendy(Command.ResearchReportExport, 'target')),
  akcja('research.export.formats', 'Formaty rozszerzone', `PPTX, XLSX oraz LaTeX z plikiem BibTeX. ${poleKomendy(Command.ResearchReportExport, 'format')}`),
];
