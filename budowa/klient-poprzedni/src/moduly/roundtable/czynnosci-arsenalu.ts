import {
  RoundtableAnalysisKind,
  RoundtableArgumentFormat,
  RoundtableHandoffTarget,
  RoundtableLeaderboardScope,
  RoundtableRatingKind,
  RoundtableTranscriptFormat,
  RoundtableVoteMethod,
} from '../../../../shared/contract';
import { przyciskAkcji as przycisk } from '../../modele/kontrolki-formularza-braki';
import type { StanDebaty } from './stan-debaty';
import type { ZrodloArsenaluRoundtable } from './zrodlo-arsenalu';

/**
 * Czynności arsenału modułu wpięte w przyciski okien.
 *
 * Plik istnieje po to, żeby żadna komenda obszaru nie została bez drogi z okna.
 * Wcześniej pozycje bez obsługi stały jako przyciski nazywające brak; teraz
 * każda z nich naprawdę woła swoją komendę, a wynik melduje w oknie, z którego
 * padła.
 *
 * ── Skąd biorą się dane żądania ──────────────────────────────────────────────
 * Ze stanu debaty wspólnego oknom modułu: okno, tura bieżąca, skład, wypowiedzi.
 * Czynność, dla której stan nie ma jeszcze wskazania — nie ma tury, nie ma
 * uczestników, nie ma wypowiedzi — nie idzie do rdzenia po to, żeby dostać
 * odmowę: melduje brak wskazania od razu i nazywa, czego brakuje. Rdzeń odmówiłby
 * tak samo, tylko po podróży tam i z powrotem.
 *
 * ── Czego tu nie ma ──────────────────────────────────────────────────────────
 * Wartości domyślnych podstawianych za Operatora. Tam, gdzie komenda potrzebuje
 * treści (stanowisko, zdanie odrębne, warianty głosowania), treść przychodzi
 * z pola okna. Przycisk bez wypełnionego pola melduje, czego brakuje.
 */

/** Meldunek okna: zdanie i to, czy czynność się powiodła. */
export type Meldunek = (zdanie: string, powodzenie: boolean) => void;

/** Wynik czynności rozpoznany po `Wynik` warstwy protokołu. */
type WynikCzynnosci = { udany: boolean; blad?: { message?: string } };

/**
 * Melduje wynik wywołania jednym zdaniem. Odmowa idzie z treścią odmowy
 * rdzenia — to ona mówi Operatorowi, co zrobić dalej; zdanie zastępcze
 * („nie udało się") zabrałoby mu tę wiedzę.
 */
function zamelduj(meldunek: Meldunek, wynik: WynikCzynnosci, powodzenie: string): void {
  if (wynik.udany) {
    meldunek(powodzenie, true);
    return;
  }
  meldunek(wynik.blad?.message ?? 'Rdzeń odmówił bez podania powodu.', false);
}

/** Buduje przycisk wołający czynność; blokada podwójnego naciśnięcia w środku. */
function czynnosc(etykieta: string, praca: () => Promise<void>): HTMLButtonElement {
  const kontrolka = przycisk(etykieta);
  kontrolka.addEventListener('click', () => {
    if (kontrolka.dataset['wtoku'] === 'tak') return;
    kontrolka.dataset['wtoku'] = 'tak';
    void praca().finally(() => {
      delete kontrolka.dataset['wtoku'];
    });
  });
  return kontrolka;
}

/** Pierwsza wypowiedź tury bieżącej; pusta znaczy „okno nic nie widziało". */
function pierwszaWypowiedz(stan: StanDebaty): string {
  const wykaz = stan.wypowiedzi();
  return wykaz.length > 0 ? wykaz[0]!.id : '';
}

/** Pierwszy uczestnik składu; pusty znaczy „debata bez składu". */
function pierwszyUczestnik(stan: StanDebaty): string {
  const wykaz = stan.uczestnicy();
  return wykaz.length > 0 ? wykaz[0]!.id : '';
}

/**
 * Czynności Model Panels: skład, zespoły i biblioteka ról.
 */
export function czynnosciSkladu(
  zrodlo: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  meldunek: Meldunek,
  wybranyUczestnik: () => string,
  nazwaZespolu: () => string,
): HTMLButtonElement[] {
  return [
    czynnosc('Usuń uczestnika', async () => {
      const kod = wybranyUczestnik();
      if (kod === '') {
        meldunek('Wskaż uczestnika w panelu, zanim go usuniesz.', false);
        return;
      }
      const wynik = await zrodlo.usunModel({ windowId: stan.okno(), participantId: kod });
      if (wynik.udany && wynik.wynik) stan.ustawUczestnikow(wynik.wynik.participants);
      zamelduj(meldunek, wynik, 'Uczestnik zdjęty ze składu; jego wypowiedzi zostają w zapisie tury.');
    }),
    czynnosc('Oznacz jako kluczowego', async () => {
      const kod = wybranyUczestnik();
      if (kod === '') {
        meldunek('Wskaż uczestnika, którego chcesz oznaczyć jako kluczowego.', false);
        return;
      }
      const wynik = await zrodlo.zmienModel({
        windowId: stan.okno(),
        participantId: kod,
        key: true,
      });
      if (wynik.udany && wynik.wynik) stan.dodajUczestnika(wynik.wynik.participant);
      zamelduj(meldunek, wynik, 'Uczestnik oznaczony jako kluczowy.');
    }),
    czynnosc('Zapisz skład jako zespół', async () => {
      const nazwa = nazwaZespolu().trim();
      if (nazwa === '') {
        meldunek('Nazwij zespół w polu obok, zanim go zapiszesz.', false);
        return;
      }
      const wynik = await zrodlo.zapiszZespol({
        windowId: stan.okno(),
        name: nazwa,
        includeFormat: true,
      });
      zamelduj(
        meldunek,
        wynik,
        `Zespół „${nazwa}" zapisany wraz z formatem debaty — wnieś go do innego okna przyciskiem obok.`,
      );
    }),
    czynnosc('Wnieś zapisany zespół', async () => {
      const zespoly = await zrodlo.zespoly({ limit: 1 });
      if (!zespoly.udany || !zespoly.wynik) {
        zamelduj(meldunek, zespoly, '');
        return;
      }
      const pierwszy = zespoly.wynik.teams[0];
      if (pierwszy === undefined) {
        meldunek('Nie ma ani jednego zapisanego zespołu — zapisz skład, zanim go wniesiesz.', false);
        return;
      }
      const wynik = await zrodlo.wniesZespol({ windowId: stan.okno(), teamId: pierwszy.id });
      if (wynik.udany && wynik.wynik) stan.ustawUczestnikow(wynik.wynik.participants);
      zamelduj(meldunek, wynik, `Zespół „${pierwszy.name}" wniesiony do składu debaty.`);
    }),
    czynnosc('Biblioteka ról debaty', async () => {
      const wynik = await zrodlo.role({});
      if (wynik.udany && wynik.wynik) {
        const nazwy = wynik.wynik.roles.map((rola) => rola.name).join(', ');
        meldunek(`Role dostępne w bibliotece: ${nazwy}.`, true);
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Pytanie doprecyzowujące', async () => {
      const kod = wybranyUczestnik() === '' ? pierwszyUczestnik(stan) : wybranyUczestnik();
      if (kod === '') {
        meldunek('Debata nie ma uczestnika, do którego można skierować pytanie.', false);
        return;
      }
      const wynik = await zrodlo.doprecyzuj({
        windowId: stan.okno(),
        participantId: kod,
        question: 'Rozwiń swoje ostatnie stanowisko i podaj, na czym je opierasz.',
        ...(stan.tura() === '' ? {} : { turnId: stan.tura() }),
      });
      zamelduj(meldunek, wynik, 'Wątek boczny zamknięty — odpowiedź uczestnika jest w zapisie debaty.');
    }),
    czynnosc('Regeneracja odpowiedzi', async () => {
      const kod = pierwszaWypowiedz(stan);
      if (kod === '') {
        meldunek('Okno nie widziało jeszcze wypowiedzi, którą można powtórzyć.', false);
        return;
      }
      const wynik = await zrodlo.powtorz({ windowId: stan.okno(), statementId: kod });
      zamelduj(meldunek, wynik, 'Wypowiedź zastąpiona nową redakcją.');
    }),
    czynnosc('Rozgałęź turę', async () => {
      if (stan.tura() === '') {
        meldunek('Żadna tura nie biegnie — nie ma czego rozgałęziać.', false);
        return;
      }
      const wynik = await zrodlo.rozgalez({
        windowId: stan.okno(),
        turnId: stan.tura(),
        label: 'wariant do porównania',
      });
      zamelduj(meldunek, wynik, 'Wariant tury założony; obie gałęzie zostają w zapisie.');
    }),
  ];
}

/**
 * Czynności Argument Map & Analysis: graf, analiza, dowody i katalog błędów.
 */
export function czynnosciAnalizy(
  zrodlo: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  meldunek: Meldunek,
  rodzajAnalizy: () => string,
  formatGrafu: () => string,
): HTMLButtonElement[] {
  return [
    czynnosc('Uruchom analizę', async () => {
      const wynik = await zrodlo.analizuj({
        windowId: stan.okno(),
        kind: (rodzajAnalizy() as RoundtableAnalysisKind) || RoundtableAnalysisKind.ArgumentMining,
        ...(stan.tura() === '' ? {} : { turnId: stan.tura() }),
      });
      if (wynik.udany && wynik.wynik) {
        meldunek(`Analiza zapisana: ${wynik.wynik.findings.length} ustaleń.`, true);
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Odczytaj graf argumentów', async () => {
      const wynik = await zrodlo.graf({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        const graf = wynik.wynik.graph;
        meldunek(
          `Graf: ${graf.nodes.length} węzłów, ${graf.edges.length} krawędzi` +
            (graf.crux === undefined ? '.' : `; punkt sporny: ${graf.crux.text}`),
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Oznacz pierwszy argument jako kluczowy', async () => {
      const graf = await zrodlo.graf({ windowId: stan.okno() });
      if (!graf.udany || !graf.wynik) {
        zamelduj(meldunek, graf, '');
        return;
      }
      const wezel = graf.wynik.graph.nodes[0];
      if (wezel === undefined) {
        meldunek('Graf jest pusty — uruchom najpierw wydobycie argumentów.', false);
        return;
      }
      const wynik = await zrodlo.przypnij({
        windowId: stan.okno(),
        nodeId: wezel.id,
        pinned: true,
      });
      zamelduj(meldunek, wynik, 'Argument oznaczony jako kluczowy; wejdzie do syntezy stanowiska.');
    }),
    czynnosc('Wydaj graf argumentów', async () => {
      const wynik = await zrodlo.wydajGraf({
        windowId: stan.okno(),
        format: (formatGrafu() as RoundtableArgumentFormat) || RoundtableArgumentFormat.Svg,
      });
      if (wynik.udany && wynik.wynik) {
        meldunek(
          `Graf wydany jako artefakt ${wynik.wynik.artifactId}` +
            (wynik.wynik.uri === undefined ? '.' : ` (${wynik.wynik.uri}).`),
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Rejestr dowodów', async () => {
      const wynik = await zrodlo.dowody({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        const niepoparte = wynik.wynik.evidence.filter((pozycja) => !pozycja.supported).length;
        meldunek(
          `Rejestr dowodów: ${wynik.wynik.evidence.length} twierdzeń, w tym ${niepoparte} niepopartych.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Katalog błędów logicznych', async () => {
      const wynik = await zrodlo.katalogBledow({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        const wykrywane = wynik.wynik.definitions.filter((pozycja) => pozycja.enabled).length;
        meldunek(
          `Katalog: ${wynik.wynik.definitions.length} pozycji, wykrywanych w tym oknie ${wykrywane}.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Punkty zgody i sporu', async () => {
      const wynik = await zrodlo.zgodnosc({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        const zgodne = wynik.wynik.points.filter((punkt) => punkt.agreed).length;
        meldunek(
          `Punkty zgody: ${zgodne}, punkty sporne: ${wynik.wynik.points.length - zgodne}.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Obozy opinii', async () => {
      const wynik = await zrodlo.klastry({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        meldunek(`Rozpoznane obozy opinii: ${wynik.wynik.clusters.length}.`, true);
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Zbieżność i dryf stanowisk', async () => {
      const zbieznosc = await zrodlo.zbieznosc({ windowId: stan.okno() });
      if (!zbieznosc.udany || !zbieznosc.wynik) {
        zamelduj(meldunek, zbieznosc, '');
        return;
      }
      const dryf = await zrodlo.dryf({ windowId: stan.okno() });
      if (!dryf.udany || !dryf.wynik) {
        zamelduj(meldunek, dryf, '');
        return;
      }
      const ostatni = zbieznosc.wynik.points[zbieznosc.wynik.points.length - 1];
      meldunek(
        `Zbieżność po ostatniej turze: ${ostatni === undefined ? 'brak pomiaru' : ostatni.convergence.toFixed(2)}; ` +
          `zmian stanowiska: ${dryf.wynik.entries.length}.`,
        true,
      );
    }),
    czynnosc('Kluczowy punkt sporny', async () => {
      const wynik = await zrodlo.punktSporny({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        meldunek(
          wynik.wynik.crux === undefined
            ? 'Żadnego punktu spornego nie rozpoznano — w grafie nie ma podważeń.'
            : `Punkt sporny: ${wynik.wynik.crux.text}`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
  ];
}

/**
 * Czynności Voting & Evaluation Center: głosowanie, oceny, rubryki i ranking.
 */
export function czynnosciOceny(
  zrodlo: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  meldunek: Meldunek,
  warianty: () => string,
  metoda: () => string,
): HTMLButtonElement[] {
  let otwarte = '';

  return [
    czynnosc('Otwórz głosowanie', async () => {
      const wykaz = warianty()
        .split('\n')
        .map((pozycja) => pozycja.trim())
        .filter((pozycja) => pozycja !== '');
      if (wykaz.length < 2) {
        meldunek('Podaj co najmniej dwa warianty — po jednym w wierszu.', false);
        return;
      }
      const wynik = await zrodlo.otworzGlosowanie({
        windowId: stan.okno(),
        method: (metoda() as RoundtableVoteMethod) || RoundtableVoteMethod.Approval,
        options: wykaz,
        ...(stan.tura() === '' ? {} : { turnId: stan.tura() }),
      });
      if (wynik.udany && wynik.wynik) otwarte = wynik.wynik.vote.id;
      zamelduj(meldunek, wynik, `Głosowanie otwarte nad ${wykaz.length} wariantami.`);
    }),
    czynnosc('Oddaj głos na pierwszy wariant', async () => {
      const glosowanie = await zrodlo.glosowanie({
        windowId: stan.okno(),
        ...(otwarte === '' ? {} : { voteId: otwarte }),
      });
      if (!glosowanie.udany || !glosowanie.wynik) {
        zamelduj(meldunek, glosowanie, '');
        return;
      }
      const wariant = glosowanie.wynik.vote.options[0];
      if (wariant === undefined) {
        meldunek('Głosowanie nie ma wariantów, na które można oddać głos.', false);
        return;
      }
      const wynik = await zrodlo.oddajGlos({
        windowId: stan.okno(),
        voteId: glosowanie.wynik.vote.id,
        voterId: 'operator',
        approvals: [wariant.id],
      });
      zamelduj(meldunek, wynik, `Głos oddany na wariant „${wariant.label}".`);
    }),
    czynnosc('Wynik głosowania', async () => {
      const wynik = await zrodlo.glosowanie({
        windowId: stan.okno(),
        ...(otwarte === '' ? {} : { voteId: otwarte }),
      });
      if (wynik.udany && wynik.wynik) {
        const rezultat = wynik.wynik.result;
        meldunek(
          rezultat === undefined
            ? 'Głosowanie otwarte — nie padł jeszcze ani jeden głos.'
            : `Oddanych głosów: ${rezultat.ballots}; zwycięzca: ${rezultat.winnerOptionId ?? 'remis'}.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Wskaż wypowiedź bardziej przekonującą', async () => {
      const wykaz = stan.wypowiedzi();
      if (wykaz.length < 2) {
        meldunek('Porównanie parami potrzebuje dwóch wypowiedzi w turze.', false);
        return;
      }
      const wynik = await zrodlo.ocen({
        windowId: stan.okno(),
        kind: RoundtableRatingKind.Pairwise,
        targetStatementId: wykaz[1]!.id,
        preferredId: wykaz[0]!.id,
      });
      zamelduj(meldunek, wynik, 'Ocena zapisana; ranking tożsamości przesunięty o ten pojedynek.');
    }),
    czynnosc('Rubryki oceny', async () => {
      const wynik = await zrodlo.rubryki({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        meldunek(
          wynik.wynik.rubrics.length === 0
            ? 'Nie ma ani jednej rubryki — załóż ją, zanim uruchomisz ocenę sędziowską.'
            : `Rubryki dostępne: ${wynik.wynik.rubrics.map((pozycja) => pozycja.name).join(', ')}.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Ocena modelami-sędziami', async () => {
      const rubryki = await zrodlo.rubryki({ windowId: stan.okno() });
      if (!rubryki.udany || !rubryki.wynik) {
        zamelduj(meldunek, rubryki, '');
        return;
      }
      const rubryka = rubryki.wynik.rubrics[0];
      if (rubryka === undefined) {
        meldunek('Nie ma rubryki, według której sędziowie mieliby oceniać.', false);
        return;
      }
      const sedzia = pierwszyUczestnik(stan);
      if (sedzia === '') {
        meldunek('Debata nie ma uczestnika, który mógłby zasiąść jako sędzia.', false);
        return;
      }
      const wynik = await zrodlo.osadz({
        windowId: stan.okno(),
        rubricId: rubryka.id,
        judgeParticipantIds: [sedzia],
        ...(stan.tura() === '' ? {} : { turnId: stan.tura() }),
      });
      if (wynik.udany && wynik.wynik) {
        meldunek(`Sędzia wystawił ${wynik.wynik.judgements.length} ocen.`, true);
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Ranking tożsamości', async () => {
      const wynik = await zrodlo.ranking({ scope: RoundtableLeaderboardScope.Environment });
      if (wynik.udany && wynik.wynik) {
        const czolo = wynik.wynik.entries[0];
        meldunek(
          czolo === undefined
            ? 'Ranking jest pusty — nie rozstrzygnięto jeszcze ani jednego pojedynku.'
            : `Na czele ${czolo.displayName} z punktacją ${czolo.rating.toFixed(0)} po ${czolo.matches} pojedynkach.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Macierz decyzyjna', async () => {
      const wynik = await zrodlo.macierz({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        const najlepszy = [...wynik.wynik.matrix.options].sort(
          (pierwszy, drugi) => (drugi.total ?? 0) - (pierwszy.total ?? 0),
        )[0];
        meldunek(
          najlepszy === undefined
            ? 'Macierz nie ma wariantów.'
            : `Najwyższy wynik ważony: ${najlepszy.label} (${(najlepszy.total ?? 0).toFixed(2)}).`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Kalibracja pewności', async () => {
      const wynik = await zrodlo.kalibracja({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        meldunek(
          wynik.wynik.calibrations.length === 0
            ? 'Nie ma ani jednej ocenionej wypowiedzi — trafności nie ma z czego zmierzyć.'
            : `Kalibrację zmierzono dla ${wynik.wynik.calibrations.length} uczestników.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
  ];
}

/**
 * Czynności Moderator Panelu: szablony moderacji i wydanie zapisu debaty.
 */
export function czynnosciModeracji(
  zrodlo: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  meldunek: Meldunek,
  nazwaSzablonu: () => string,
  formatTranskryptu: () => string,
): HTMLButtonElement[] {
  return [
    czynnosc('Zapisz szablon moderacji', async () => {
      const nazwa = nazwaSzablonu().trim();
      if (nazwa === '') {
        meldunek('Nazwij szablon, zanim go zapiszesz.', false);
        return;
      }
      const wynik = await zrodlo.zapiszSzablon({ windowId: stan.okno(), name: nazwa });
      zamelduj(meldunek, wynik, `Szablon „${nazwa}" zapisany wraz z formatem i kolejnością głosu.`);
    }),
    czynnosc('Szablony moderacji', async () => {
      const wynik = await zrodlo.szablony({});
      if (wynik.udany && wynik.wynik) {
        meldunek(
          wynik.wynik.templates.length === 0
            ? 'Nie ma ani jednego zapisanego szablonu moderacji.'
            : `Szablony: ${wynik.wynik.templates.map((pozycja) => pozycja.name).join(', ')}.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Wydaj transkrypt', async () => {
      const wynik = await zrodlo.wydajTranskrypt({
        windowId: stan.okno(),
        format:
          (formatTranskryptu() as RoundtableTranscriptFormat) || RoundtableTranscriptFormat.Markdown,
        includeAnalysis: true,
      });
      if (wynik.udany && wynik.wynik) {
        meldunek(
          `Transkrypt wydany jako artefakt ${wynik.wynik.artifactId}` +
            (wynik.wynik.uri === undefined ? '.' : ` (${wynik.wynik.uri}).`),
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Odsłuch debaty', async () => {
      const wynik = await zrodlo.odsluch({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        meldunek(
          `Nagranie gotowe (${wynik.wynik.durationMs ?? 0} ms): ${wynik.wynik.uri ?? wynik.wynik.artifactId}.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Odczytaj stan debaty z rdzenia', async () => {
      const wynik = await zrodlo.stanDebaty({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        const migawka = wynik.wynik.snapshot;
        stan.ustawUczestnikow(migawka.participants);
        meldunek(
          `Odczytano ${migawka.turns.length} tur i ${migawka.statements.length} wypowiedzi.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
  ];
}

/**
 * Czynności Consensus Panelu: redakcja stanowiska, zdanie odrębne, wersje
 * i przekazanie do modułu docelowego.
 */
export function czynnosciStanowiska(
  zrodlo: ZrodloArsenaluRoundtable,
  stan: StanDebaty,
  meldunek: Meldunek,
  tresc: () => string,
  modulDocelowy: () => string,
): HTMLButtonElement[] {
  let ostatnie = '';

  return [
    czynnosc('Zapisz stanowisko', async () => {
      const zapis = tresc().trim();
      if (zapis === '') {
        meldunek('Stanowisko bez treści skasowałoby zapis debaty — wpisz treść.', false);
        return;
      }
      const wynik = await zrodlo.zapiszStanowisko({ windowId: stan.okno(), content: zapis });
      if (wynik.udany && wynik.wynik) ostatnie = wynik.wynik.consensus.id;
      zamelduj(
        meldunek,
        wynik,
        'Stanowisko zapisane; od tej chwili należy do Ciebie i nie wraca do złożenia z tur.',
      );
    }),
    czynnosc('Zapisz zdanie odrębne', async () => {
      if (ostatnie === '') {
        meldunek('Zapisz najpierw stanowisko, wobec którego zdanie jest odrębne.', false);
        return;
      }
      const uczestnik = pierwszyUczestnik(stan);
      if (uczestnik === '') {
        meldunek('Debata nie ma uczestnika, który mógłby podpisać zdanie odrębne.', false);
        return;
      }
      const zapis = tresc().trim();
      if (zapis === '') {
        meldunek('Wpisz treść zdania odrębnego.', false);
        return;
      }
      const wynik = await zrodlo.zapiszZdanieOdrebne({
        windowId: stan.okno(),
        consensusId: ostatnie,
        participantId: uczestnik,
        content: zapis,
      });
      zamelduj(meldunek, wynik, 'Zdanie odrębne podpisane i dołączone do stanowiska.');
    }),
    czynnosc('Wersje stanowiska', async () => {
      const wynik = await zrodlo.wersjeStanowiska({ windowId: stan.okno() });
      if (wynik.udany && wynik.wynik) {
        meldunek(`Redakcji stanowiska: ${wynik.wynik.versions.length}.`, true);
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
    czynnosc('Przekaż stanowisko', async () => {
      if (ostatnie === '') {
        meldunek('Zapisz najpierw stanowisko, które ma zostać przekazane.', false);
        return;
      }
      const wynik = await zrodlo.przekazStanowisko({
        windowId: stan.okno(),
        consensusId: ostatnie,
        target: (modulDocelowy() as RoundtableHandoffTarget) || RoundtableHandoffTarget.Studio,
        includeTranscript: true,
      });
      if (wynik.udany && wynik.wynik) {
        meldunek(
          `Stanowisko przekazane; artefakt ${wynik.wynik.handoff.artifactId ?? 'bez identyfikatora'}.`,
          true,
        );
        return;
      }
      zamelduj(meldunek, wynik, '');
    }),
  ];
}
