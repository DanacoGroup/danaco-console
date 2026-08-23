/**
 * Katalog funkcji modułu Roundtable — pełny wykaz narzędzi opracowania wraz ze
 * stanem, w jakim moduł je dziś wykonuje.
 *
 * Katalog istnieje po to, żeby okno nie wyglądało na kompletne tam, gdzie nie
 * jest. Opracowanie modułu wymienia siedemdziesiąt cztery narzędzia — bez wykazu
 * Operator poznawałby rozjazd między nimi a stanem modułu dopiero po naciśnięciu
 * każdego przycisku z osobna.
 *
 * Po scaleniu kontraktu rozjazd zmienił naturę. Obszar `roundtable.*` niósł
 * cztery komendy i niesie ich dziś czterdzieści sześć, więc zdanie „tego nie
 * wykona żadna komenda obszaru” przestało być prawdziwe co do czterdziestu
 * dwóch narzędzi naraz. Brakuje ich obsługi, nie kontraktu — i to rozróżnienie
 * niesie stan pozycji.
 *
 * Wykaz jest danymi, nie wyglądem: rysuje go `warstwy-modulu.ts`, a każde okno
 * bierze z niego wyłącznie swoje pozycje. Nazwy narzędzi stoją w brzmieniu
 * opracowania modułu, bez numeracji rozdziałów — numer rozdziału nie jest
 * identyfikatorem niczego w produkcie.
 *
 * Warstwa pozycji jest warstwą widoczności z opracowania (rozdz. 3.1), przypisaną
 * po elemencie interfejsu, którym narzędzie się otwiera:
 *   1 — widoczne bez interakcji,
 *   2 — znacznik kontekstowy, przycisk, przełącznik,
 *   3 — menu zestawu akcji okna,
 *   4 — nastawa dostępna poleceniem, wyszukiwarką funkcji albo trybem
 *       administracyjnym; w oknie nie ma jej w stanie spoczynku.
 */

/** Okno, w którym narzędzie jest osiągalne. */
export type OknoModulu =
  | 'model-panels'
  | 'debate-panel'
  | 'argument-map-analysis'
  | 'voting-evaluation-center'
  | 'moderator-panel'
  | 'consensus-panel'
  | 'chat-window'
  | 'execution-loop-window';

/** Warstwa widoczności z opracowania modułu. */
export type WarstwaWidocznosci = 1 | 2 | 3 | 4;

/**
 * Stan wykonania narzędzia — cztery stany, z których każdy znaczy co innego.
 *
 * `wykonana` — moduł wykonuje narzędzie w zakresie opracowania.
 *
 * `czesciowa` — nie jest stanem pośrednim między porażką a sukcesem: znaczy, że
 * narzędzie działa w części zakresu, a nie działa w reszcie. Zdanie pozycji
 * mówi, gdzie przebiega granica.
 *
 * `bez-obslugi` — komenda JEST w kontrakcie, a moduł jej nie wywołuje.
 * To najliczniejszy stan modułu i najważniejszy do odróżnienia: nie brakuje
 * uzgodnienia, brakuje pracy. Zdanie pozycji nazywa komendę, którą narzędzie
 * wejdzie.
 *
 * `bez-pokrycia` — kontrakt nie ma czym narzędzia wykonać: nie ma komendy albo
 * ma komendę, lecz brakuje pola w jej żądaniu. Po scaleniu kontraktu stan ten
 * zszedł do pojedynczych pozycji i każda nazywa brakujące pole.
 *
 * `poza-modulem` — narzędzie należy do okna wspólnego platformy.
 */
export type StanFunkcji =
  | 'wykonana'
  | 'czesciowa'
  | 'bez-obslugi'
  | 'bez-pokrycia'
  | 'poza-modulem';

export interface FunkcjaKatalogu {
  /** Nazwa narzędzia w brzmieniu opracowania modułu. */
  nazwa: string;
  okno: OknoModulu;
  warstwa: WarstwaWidocznosci;
  stan: StanFunkcji;
  /** Co moduł wykonuje albo czego kontrakt nie niesie — zdanie pełne. */
  zdanie: string;
}

const SKLAD_I_PERSONY: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Participant Composer',
    okno: 'model-panels',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Uczestnika zakłada roundtable.model.add z kanału rejestru, nazwy tożsamości i promptu systemowego. Awatar, opis roli, rola, waga i agent wchodzą drugim krokiem — polami żądania roundtable.model.update, które okno wywołuje z paska akcji.',
  },
  {
    nazwa: 'Multiple Identities per Channel',
    okno: 'model-panels',
    warstwa: 2,
    stan: 'wykonana',
    zdanie:
      'Ten sam kanał wchodzi do debaty wielokrotnie pod odrębnymi tożsamościami; panel powtórzonego kanału jest oznaczony, żeby dwa podobne panele nie czytały się jako usterka.',
  },
  {
    nazwa: 'Role Library',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Gotowe role wraz z promptami systemowymi oddaje roundtable.role.list; okno wywołuje ją z zestawu akcji uczestnika, a prompt roli przenosi się do pola tożsamości.',
  },
  {
    nazwa: 'Persona Card & Profile',
    okno: 'model-panels',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Panel uczestnika pokazuje tożsamość, kanał, prompt systemowy, miejsce w kolejności głosu, wyciszenie, oznaczenie kluczowe i wypowiedzi widziane przez okno. Historię sprzed otwarcia okna wczytuje roundtable.debate.get z zestawu akcji moderacji; statystyk udziału między sesjami kontrakt nie niesie.',
  },
  {
    nazwa: 'Agent-as-Participant',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Zapisanego agenta wpina pole agentId żądania roundtable.model.update; uczestnik ma w kontrakcie odpowiadające mu pole.',
  },
  {
    nazwa: 'Panel Presets / Teams',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Cały skład jako nazwany zespół zapisuje roundtable.team.save, wykaz zespołów oddaje roundtable.team.list, a wnosi go do debaty roundtable.team.apply.',
  },
  {
    nazwa: 'Self-Consistency Ensemble',
    okno: 'model-panels',
    warstwa: 4,
    stan: 'bez-pokrycia',
    zdanie:
      'RoundtableParticipant ma pole sampleCount, ale nie ma go w żądaniu żadnej komendy obszaru — liczby wariantów jednej persony Operator nie ma dziś czym zapisać.',
  },
  {
    nazwa: 'Weighting by Expertise',
    okno: 'model-panels',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Wagę kompetencji ustawia pole weight żądania roundtable.model.update; uczestnik niesie ją w kontrakcie i wchodzi do głosowania oraz syntezy.',
  },
];

const ODPOWIEDZ_ROWNOLEGLA: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Parallel Broadcast',
    okno: 'model-panels',
    warstwa: 1,
    stan: 'wykonana',
    zdanie:
      'Jedno pytanie roundtable.debate.start idzie do całego składu, a odpowiedzi rosną równolegle we wszystkich panelach z fragmentów strumienia.',
  },
  {
    nazwa: 'Side-by-Side Compare',
    okno: 'model-panels',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Odpowiedzi stoją obok siebie w siatce paneli, a zbieżność zdań liczy Argument Map & Analysis — leksykalnie, nie semantycznie. Scalenie powtórzeń semantycznych zleca roundtable.analysis.run w odmianie dedup.',
  },
  {
    nazwa: 'Blind Mode',
    okno: 'model-panels',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Okno anonimizuje nagłówki paneli jako nastawę widoku po stronie klienta. RoundtableTurn ma pole anonymous, ale nie ma go w żądaniu żadnej komendy obszaru, więc anonimizacji po stronie rdzenia nie ma dziś czym zamówić.',
  },
  {
    nazwa: 'Response Metrics',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Okno mierzy długość wypowiedzi i czas od otwarcia tury do jej zapisania. Pewność uczestnika niesie pole confidence wypowiedzi — rdzeń wypełnia je klasyfikacją analizy, gdy uczestnik pewność zadeklarował; kosztu tokenów nie niesie żadne pole kontraktu.',
  },
  {
    nazwa: 'Follow-up to One Panel',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Pytanie doprecyzowujące do jednego uczestnika poza turą ogólną kieruje roundtable.debate.followup.',
  },
  {
    nazwa: 'Regenerate & Branch',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Powtórzenie wywołania kanału wykonuje roundtable.statement.regenerate, a rozgałęzienie tury jako wariantu — roundtable.debate.branch.',
  },
  {
    nazwa: 'Cross-Model Fact-Check',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Krzyżową weryfikację faktyczności twierdzeń zleca roundtable.analysis.run w odmianie factCheck.',
  },
];

const DEBATA_I_TURY: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Debate Formats',
    okno: 'debate-panel',
    warstwa: 2,
    stan: 'wykonana',
    zdanie:
      'Pole formatu niesie komplet sześciu wartości wyliczenia RoundtableFormat: debatę swobodną, strukturalną, oksfordzką, turę okrężną, rundy Delphi i panel ekspercki.',
  },
  {
    nazwa: 'Delphi Rounds',
    okno: 'debate-panel',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Format delphi stoi w polu wyboru i jedzie do rdzenia. Anonimowości rundy nie ma czym zamówić: pole anonymous jest w turze, ale nie ma go w żądaniu żadnej komendy obszaru.',
  },
  {
    nazwa: 'Cross-Examination',
    okno: 'debate-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Pytanie adresowane do jednego uczestnika kieruje roundtable.debate.followup, a powiązanie wypowiedzi niesie pole replyToId.',
  },
  {
    nazwa: 'Turn Orchestrator',
    okno: 'debate-panel',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'Liczbę tur ustawia turnLimit żądania tury, a kolejność głosu speakingOrder czynności moderatora. Pola timeLimitMs i maxStatementChars są w turze, ale nie ma ich w żądaniu żadnej komendy obszaru.',
  },
  {
    nazwa: 'Devil\'s Advocate Injection',
    okno: 'moderator-panel',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Kontrgłos wprowadza wartość injectCounterVoice wyliczenia ModeratorAction; progu zbieżności, po którym miałby wchodzić, nie niesie dziś żadne pole kontraktu.',
  },
  {
    nazwa: 'Side Threads',
    okno: 'moderator-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Wątek poboczny otwierają i scalają wartości openSideThread i mergeSideThread wyliczenia ModeratorAction, a turę nadrzędną wskazuje pole parentTurnId tury.',
  },
  {
    nazwa: 'Jury Mode',
    okno: 'voting-evaluation-center',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Podział na debatujących i jury zapisuje pole role żądania roundtable.model.update; wyliczenie ról niesie RoundtableParticipantRole.',
  },
  {
    nazwa: 'Position Drift Tracking',
    okno: 'debate-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Zmiany stanowisk uczestników między turami wraz z powodem oddaje roundtable.drift.get.',
  },
];

const ANALIZA_ARGUMENTOW: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Argument Map (graf)',
    okno: 'argument-map-analysis',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Węzły i krawędzie grafu oddaje roundtable.argument.list, a relacje niosą pola replyToId i speechAct wypowiedzi. Okno pokazuje dziś chronologię i mierzy, w ilu wypowiedziach rdzeń wypełnił już pola relacji.',
  },
  {
    nazwa: 'Argument Mining',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Wydobycie przesłanek i wniosków z wypowiedzi zleca roundtable.analysis.run w odmianie argumentMining.',
  },
  {
    nazwa: 'Speech Act Classification',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Klasyfikację aktów mowy zleca roundtable.analysis.run w odmianie speechAct, a rodzaj wypowiedzi niesie jej pole speechAct.',
  },
  {
    nazwa: 'Fallacy Detection',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Wykrywanie chwytów erystycznych zleca roundtable.analysis.run w odmianie fallacy; katalog wykrywanych błędów prowadzą roundtable.fallacy.catalog.get i .set.',
  },
  {
    nazwa: 'Steelman / Strawman Check',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Ocenę rzetelności kontry zleca roundtable.analysis.run w odmianie steelman; powiązanie dwóch wypowiedzi niesie pole replyToId.',
  },
  {
    nazwa: 'Evidence & Citation Ledger',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Dowody i cytowania wraz z oznaczeniem twierdzeń niepopartych oddaje roundtable.evidence.list.',
  },
  {
    nazwa: 'AIF / Argdown Export',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Graf w formatach DOT, GraphML, Argdown i AIF oraz jako obraz wydaje roundtable.argument.export. Okno eksportuje dziś chronologię w Markdown.',
  },
  {
    nazwa: 'Key-Argument Pinning',
    okno: 'model-panels',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Węzeł grafu oznacza jako kluczowy roundtable.argument.pin, a pole key uczestnika ustawia roundtable.model.update.',
  },
];

const OCENA_I_GLOSOWANIE: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Vote Engine (metody)',
    okno: 'voting-evaluation-center',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'Głosowanie otwiera roundtable.vote.start wraz z metodą agregacji, głos przyjmuje roundtable.vote.cast, a wynik oddaje roundtable.vote.get.',
  },
  {
    nazwa: 'Star & Preference Rating',
    okno: 'voting-evaluation-center',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Ocenę gwiazdkami i porównanie parami zapisuje roundtable.rating.set.',
  },
  {
    nazwa: 'LLM-as-Judge',
    okno: 'voting-evaluation-center',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Ocenę wypowiedzi modelami-sędziami według rubryki wraz z uzasadnieniem zleca roundtable.judge.run.',
  },
  {
    nazwa: 'Rubric Builder',
    okno: 'voting-evaluation-center',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Rubrykę z kryteriami i wagami zakłada roundtable.rubric.set, a rubryki dostępne oknu oddaje roundtable.rubric.list.',
  },
  {
    nazwa: 'Model Elo / TrueSkill Leaderboard',
    okno: 'voting-evaluation-center',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Ranking akumulowany między sesjami oddaje roundtable.leaderboard.get wraz z algorytmem i zakresem akumulacji.',
  },
  {
    nazwa: 'Multi-Criteria Decision Matrix',
    okno: 'voting-evaluation-center',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Warianty i kryteria z wagami zapisuje roundtable.decision.matrix.set, a macierz wraz z wynikiem oddaje roundtable.decision.matrix.get.',
  },
  {
    nazwa: 'Confidence & Calibration',
    okno: 'voting-evaluation-center',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Zestawienie deklarowanej pewności z trafnością oddaje roundtable.calibration.get; pewność uczestnika niesie pole confidence wypowiedzi.',
  },
  {
    nazwa: 'Quorum & Threshold',
    okno: 'voting-evaluation-center',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'Próg zgody niesie pole quorum żądania roundtable.vote.start, a jego osiągnięcie orzeka wynik z roundtable.vote.get. Okno pokazuje dziś udział zmierzony i nazywa go udziałem, nie kworum.',
  },
];

const ZGODA_I_SPOR: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Agreement / Dispute Detection',
    okno: 'argument-map-analysis',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'Okno wskazuje zdania powtarzające się u kilku mówców i zdania wyłącznie własne — jest to zbieżność leksykalna, nie zgodność stanowisk. Punkty zgody i punkty sporne oddaje roundtable.agreement.get.',
  },
  {
    nazwa: 'Agreement Matrix (heatmapa)',
    okno: 'argument-map-analysis',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Macierz par mówców liczona jest zbieżnością leksykalną ich wypowiedzi w turze. Macierz zgodności stanowisk oddaje roundtable.agreement.get w polu matrix.',
  },
  {
    nazwa: 'Opinion Clustering',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Obozy zbliżonych stanowisk oddaje roundtable.cluster.get.',
  },
  {
    nazwa: 'Semantic Dedup of Arguments',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Scalenie argumentów powtórzonych semantycznie zleca roundtable.analysis.run w odmianie dedup; okno umie dziś wskazać wyłącznie powtórzenie dosłowne.',
  },
  {
    nazwa: 'Consensus Convergence Meter',
    okno: 'argument-map-analysis',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Pomiar zbliżania stanowisk tura po turze oddaje roundtable.convergence.get; klient sam go nie policzy, bo widzi wyłącznie turę bieżącą.',
  },
  {
    nazwa: 'Crux Finder',
    okno: 'argument-map-analysis',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Kluczowy punkt sporny wskazuje roundtable.crux.get z zależności między węzłami grafu.',
  },
];

const SYNTEZA_I_KONSENSUS: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Consensus Generator',
    okno: 'consensus-panel',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'roundtable.consensus.get oddaje stanowisko złożone przez rdzeń — tury bieżącej albo całej debaty — i okno je czyta. Zlecenia wygenerowania stanowiska kontrakt nadal nie ma.',
  },
  {
    nazwa: 'Editable Position',
    okno: 'consensus-panel',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'Treść stanowiska po redakcji Operatora zapisuje roundtable.consensus.set; odczyt przestał być jednokierunkowy.',
  },
  {
    nazwa: 'Minority Report',
    okno: 'consensus-panel',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'Zdanie odrębne zapisuje i podpisuje roundtable.consensus.minority.set, a stanowisko niesie je polem minority.',
  },
  {
    nazwa: 'Weighted Support Tally',
    okno: 'consensus-panel',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Poparcie ważone niesie pole support stanowiska; wagi uczestników ustawia roundtable.model.update, a rozkład głosów oddaje roundtable.vote.get.',
  },
  {
    nazwa: 'Points of Agreement / Dispute Summary',
    okno: 'consensus-panel',
    warstwa: 2,
    stan: 'czesciowa',
    zdanie:
      'Punkty zgody i sporne niosą pola agreementPoints i disputePoints stanowiska, a osobno oddaje je roundtable.agreement.get.',
  },
  {
    nazwa: 'Position Versioning',
    okno: 'consensus-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Kolejne redakcje jako odrębne wersje do porównania oddaje roundtable.consensus.version.list; numer wersji nadaje rdzeń i okno pokazuje go w stopce.',
  },
  {
    nazwa: 'Decision Record (ADR)',
    okno: 'consensus-panel',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Okno składa zapis decyzji z pól stanowiska i wypisuje wprost części nieuzupełnione. Kontekst, warianty i konsekwencje zapisuje roundtable.consensus.set, a ryzyko mniejszości — roundtable.consensus.minority.set.',
  },
];

const MODERACJA: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Moderator Interventions',
    okno: 'moderator-panel',
    warstwa: 1,
    stan: 'wykonana',
    zdanie:
      'Interwencja jedzie czynnością moderatora do wszystkich uczestników tury; rdzeń wstrzykuje ją do kontekstu zadań.',
  },
  {
    nazwa: 'Turn Timer & Limits',
    okno: 'moderator-panel',
    warstwa: 1,
    stan: 'czesciowa',
    zdanie:
      'Okno liczy na żywo czas od chwili startedAt tury, a dla tury zamkniętej podaje czas jej trwania. Pola timeLimitMs i maxStatementChars są w turze, ale nie ma ich w żądaniu żadnej komendy obszaru, więc zegar mierzy zamiast odliczać.',
  },
  {
    nazwa: 'Speaking Order Control',
    okno: 'moderator-panel',
    warstwa: 1,
    stan: 'wykonana',
    zdanie:
      'Kolejność głosu przestawia się w wykazie składu i jedzie do rdzenia polem speakingOrder czynności moderatora. Wyliczenie ma dziś osobną wartość setSpeakingOrder, której okno jeszcze nie wysyła.',
  },
  {
    nazwa: 'Manual Turn Close',
    okno: 'moderator-panel',
    warstwa: 1,
    stan: 'wykonana',
    zdanie:
      'Turę zamyka czynność moderatora closeTurn, niezależnie od stanu wypowiedzi uczestników.',
  },
  {
    nazwa: 'Moderation Templates',
    okno: 'moderator-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Format, liczbę tur i kolejność głosu zapisuje jako nazwany szablon roundtable.moderation.template.save, a wykaz oddaje roundtable.moderation.template.list. Wniesienia szablonu do tury nie niesie dziś żadne pole żądania roundtable.debate.start.',
  },
  {
    nazwa: 'Bias & Tone Guardrails',
    okno: 'moderator-panel',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Sygnalizację tonu, offtopu i dominacji uczestnika zleca roundtable.analysis.run w odmianie tone.',
  },
];

const PETLA_WYKONAWCZA: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Debate Brief Decomposition',
    okno: 'execution-loop-window',
    warstwa: 1,
    stan: 'poza-modulem',
    zdanie:
      'Execution Loop Window jest oknem wspólnym platformy — nie należy do katalogu okien operacyjnych tego modułu i nie jest przez niego budowane.',
  },
  {
    nazwa: 'Multi-Executor Dispatch',
    okno: 'execution-loop-window',
    warstwa: 1,
    stan: 'poza-modulem',
    zdanie:
      'Przydział zadań tury wykonawcom prowadzi pętla wykonawcza platformy, poza tym modułem.',
  },
  {
    nazwa: 'Executor Health & Retry',
    okno: 'execution-loop-window',
    warstwa: 4,
    stan: 'poza-modulem',
    zdanie:
      'Reguły ponowienia zadań pętli należą do okna wspólnego platformy; obszar roundtable nie ma pola ani komendy ponowienia zadania.',
  },
  {
    nazwa: 'Loop Quality Gates',
    okno: 'execution-loop-window',
    warstwa: 4,
    stan: 'poza-modulem',
    zdanie:
      'Kontrola jakości wyniku zadania leży w pętli wykonawczej platformy, poza tym modułem.',
  },
  {
    nazwa: 'Loop Control',
    okno: 'execution-loop-window',
    warstwa: 1,
    stan: 'poza-modulem',
    zdanie:
      'Wstrzymanie, wznowienie i przerwanie pętli prowadzi okno wspólne platformy.',
  },
  {
    nazwa: 'Coordinator Message Log',
    okno: 'execution-loop-window',
    warstwa: 3,
    stan: 'poza-modulem',
    zdanie:
      'Rejestr komunikatów Koordynator ↔ Wykonawca prowadzi okno wspólne platformy.',
  },
];

const WYDANIE_I_INTEGRACJE: readonly FunkcjaKatalogu[] = [
  {
    nazwa: 'Transcript Export',
    okno: 'debate-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Okno pobiera transkrypt w Markdown i w JSON, z treścią złożoną z zapisu tury i z fragmentów strumienia. Wydanie w PDF i DOCX wykonuje roundtable.transcript.export, a zakres poszerza roundtable.debate.get.',
  },
  {
    nazwa: 'Argument Map Export',
    okno: 'argument-map-analysis',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Okno pobiera strukturę zapisu w Markdown. Graf w SVG, PNG, DOT, GraphML, Argdown i AIF wydaje roundtable.argument.export.',
  },
  {
    nazwa: 'Voting/Scorecard Report',
    okno: 'voting-evaluation-center',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Okno pobiera zestawienie udziału uczestników w Markdown. Wyniki głosowań i oceny rubrykowe wchodzą do raportu z roundtable.vote.get i roundtable.judge.run.',
  },
  {
    nazwa: 'TTS Debate Readback',
    okno: 'consensus-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Odczytanie przebiegu debaty głosem przypisanym uczestnikowi zleca roundtable.speech.synthesize.',
  },
  {
    nazwa: 'Send to Studio',
    okno: 'consensus-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Przekazanie stanowiska do modułu Studio wykonuje roundtable.consensus.handoff z wartością studio pola target.',
  },
  {
    nazwa: 'Send to Research',
    okno: 'consensus-panel',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Przekazanie stanowiska do modułu Research wykonuje roundtable.consensus.handoff z wartością research pola target.',
  },
  {
    nazwa: 'Research Intake',
    okno: 'chat-window',
    warstwa: 1,
    stan: 'poza-modulem',
    zdanie:
      'Wstępna synteza z modułu Research wchodzi pytaniem w oknie rozmowy, które jest bytem sesji i składa je warstwa rozmowy, nie ten moduł.',
  },
  {
    nazwa: 'Automations Step',
    okno: 'execution-loop-window',
    warstwa: 3,
    stan: 'czesciowa',
    zdanie:
      'Udostępnienie debaty jako kroku procesu wsadowego wykonuje roundtable.consensus.handoff z wartością automations pola target.',
  },
  {
    nazwa: 'Roundtable API / MCP',
    okno: 'execution-loop-window',
    warstwa: 4,
    stan: 'bez-pokrycia',
    zdanie:
      'Wystawienia debaty i głosowania na zewnątrz nie niesie kontrakt obszaru; komendy roundtable są wywoływane wyłącznie z klienta.',
  },
  {
    nazwa: 'External Model Connectors',
    okno: 'model-panels',
    warstwa: 4,
    stan: 'czesciowa',
    zdanie:
      'Uczestnikiem zostaje dowolny kanał rejestru rdzenia, także kanał modelu zewnętrznego. Kluczy dostępowych i konfiguracji kanału ten moduł nie prowadzi — należą do okna konfiguracji.',
  },
];

/** Pełny katalog narzędzi opracowania modułu. */
export const KATALOG_FUNKCJI: readonly FunkcjaKatalogu[] = [
  ...SKLAD_I_PERSONY,
  ...ODPOWIEDZ_ROWNOLEGLA,
  ...DEBATA_I_TURY,
  ...ANALIZA_ARGUMENTOW,
  ...OCENA_I_GLOSOWANIE,
  ...ZGODA_I_SPOR,
  ...SYNTEZA_I_KONSENSUS,
  ...MODERACJA,
  ...PETLA_WYKONAWCZA,
  ...WYDANIE_I_INTEGRACJE,
];

/** Narzędzia przypisane jednemu oknu, w kolejności katalogu. */
export function funkcjeOkna(okno: OknoModulu): FunkcjaKatalogu[] {
  return KATALOG_FUNKCJI.filter((funkcja) => funkcja.okno === okno);
}

/** Narzędzia okna należące do jednej warstwy widoczności. */
export function funkcjeWarstwy(okno: OknoModulu, warstwa: WarstwaWidocznosci): FunkcjaKatalogu[] {
  return funkcjeOkna(okno).filter((funkcja) => funkcja.warstwa === warstwa);
}

/** Bilans stanów wykonania — dla zdania nad wykazem. */
export interface BilansFunkcji {
  wszystkich: number;
  wykonanych: number;
  czesciowych: number;
  bezObslugi: number;
  bezPokrycia: number;
  pozaModulem: number;
}

export function bilansFunkcji(funkcje: readonly FunkcjaKatalogu[]): BilansFunkcji {
  const ile = (stan: StanFunkcji): number =>
    funkcje.filter((funkcja) => funkcja.stan === stan).length;
  return {
    wszystkich: funkcje.length,
    wykonanych: ile('wykonana'),
    czesciowych: ile('czesciowa'),
    bezObslugi: ile('bez-obslugi'),
    bezPokrycia: ile('bez-pokrycia'),
    pozaModulem: ile('poza-modulem'),
  };
}

/** Nazwa stanu w brzmieniu, w jakim czyta ją Operator. */
export function opisStanuFunkcji(stan: StanFunkcji): string {
  switch (stan) {
    case 'wykonana':
      return 'wykonane';
    case 'czesciowa':
      return 'wykonane częściowo';
    case 'bez-obslugi':
      return 'komenda w kontrakcie, obsługi jeszcze nie zbudowano';
    case 'bez-pokrycia':
      return 'bez pokrycia w kontrakcie';
    case 'poza-modulem':
      return 'poza tym modułem';
  }
}
