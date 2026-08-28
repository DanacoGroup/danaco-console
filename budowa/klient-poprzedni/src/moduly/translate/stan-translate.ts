import { ChangeKind, type TranslationPanel, type Window } from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { utworzMagazynTranslate, type FazaKontekstu, type MagazynTranslate } from './magazyn-translate';
import {
  utworzZrodloDokumentuTranslate,
  type ZrodloDokumentuTranslate,
} from './zrodlo-dokumentu-translate';
import { utworzZrodloGlosariusza, type ZrodloGlosariusza } from './zrodlo-glosariusza';
import {
  utworzZrodloKanalowTranslate,
  type ZrodloKanalowTranslate,
} from './zrodlo-kanalow-translate';
import { utworzZrodloOknaTranslate, type ZrodloOknaTranslate } from './zrodlo-okna-translate';
import { utworzZrodloPaneli, type ZrodloPaneli } from './zrodlo-paneli';
import {
  utworzZrodloWarsztatuTranslate,
  type ZrodloWarsztatuTranslate,
} from './zrodlo-warsztatu-translate';
import { utworzZrodloZrodla, type ZrodloZrodla } from './zrodlo-zrodla';

/**
 * Stan modułu Translate jest jednym źródłem kontraktu nad jednym zbiorem danych; tekst źródłowy
 * i panele są jednym bytem, więc zmiana źródła aktualizuje wszystkie panele naraz.
 */
export type { FazaKontekstu } from './magazyn-translate';

export interface StanTranslate {
  zrodlo: ZrodloZrodla;
  panele: ZrodloPaneli;
  glosariusz: ZrodloGlosariusza;
  okna: ZrodloOknaTranslate;
  /** Dokument wejściowy i jego zamiana formatu to obszar kontraktu, z którego korzysta Format Studio. */
  dokument: ZrodloDokumentuTranslate;
  /** Rejestr kanałów jest wspólny dla wszystkich sterów, żeby moduł nie pytał rdzenia wielokrotnie. */
  kanaly: ZrodloKanalowTranslate;
  /** Warsztat obejmuje rodziny komend spoza czterech okien: pamięć, segmentację, terminologię, korektę. */
  warsztat: ZrodloWarsztatuTranslate;
  /** Okno modułu wskazane przez rdzeń; pusty napis znaczy brak. */
  idOkna(): string;
  /** Okno wraz z parametrami wykonania albo `null`, gdy rdzeń go nie oddał. */
  oknoModulu(): Window | null;
  /** Faza odczytu kontekstu okna — trzy pustki są rozróżnialne. */
  fazaKontekstu(): FazaKontekstu;
  /** Powód nieudanego odczytu kontekstu; pusty, gdy odczyt się udał. */
  powodKontekstu(): string;
  tekstZrodlowy(): string;
  jezykZrodlowy(): string;
  segmenty(): readonly string[];
  /** Zapamiętuje tekst przyjęty przez rdzeń wraz z rozpoznanym językiem. */
  wchlonZrodlo(tekst: string, jezyk: string, segmenty: readonly string[]): void;
  /** Zapamiętuje segmenty oddane przez `translate.source.segment`. */
  wchlonSegmenty(segmenty: readonly string[]): void;
  /** Panele języków docelowych w kolejności nadanej przez rdzeń. */
  panelJezykow(): readonly TranslationPanel[];
  wchlonPanel(panel: TranslationPanel): void;
  wchlonPanele(panele: readonly TranslationPanel[]): void;
  /** Odczytuje okno modułu z rdzenia; wolno wołać wielokrotnie. */
  wczytajKontekst(idSesji: string): Promise<void>;
  obserwuj(sluchacz: () => void): () => void;
  rozlacz(): void;
}

export function utworzStanTranslate(kanal: Kanal): StanTranslate {
  const magazyn = utworzMagazynTranslate();
  const zrodlo = utworzZrodloZrodla(kanal);
  const panele = utworzZrodloPaneli(kanal);
  const glosariusz = utworzZrodloGlosariusza(kanal);
  const okna = utworzZrodloOknaTranslate(kanal);
  const kanaly = utworzZrodloKanalowTranslate(kanal);
  const dokument = utworzZrodloDokumentuTranslate(kanal);
  const warsztat = utworzZrodloWarsztatuTranslate(kanal);

  // Powrót rejestru kanałów przerysowuje moduł tą samą drogą, co każda inna zmiana stanu modułu.
  const odsubskrybujKanaly = kanaly.obserwuj(() => magazyn.ogloszZmiane());

  // Zdarzenie jest jedynym źródłem odświeżenia poza własnym działaniem modułu na panelach języków.
  const odsubskrybuj = panele.naZmiane((tresc) => {
    if (tresc.change === ChangeKind.Deleted) magazyn.usunPanel(tresc.panel.id);
    else magazyn.wchlonPanel(tresc.panel);
  });

  return {
    zrodlo,
    panele,
    glosariusz,
    okna,
    dokument,
    kanaly,
    warsztat,
    idOkna: () => magazyn.oknoModulu()?.id ?? '',
    oknoModulu: magazyn.oknoModulu,
    fazaKontekstu: magazyn.faza,
    powodKontekstu: magazyn.powod,
    tekstZrodlowy: magazyn.tekstZrodlowy,
    jezykZrodlowy: magazyn.jezykZrodlowy,
    segmenty: magazyn.segmenty,
    wchlonZrodlo: magazyn.ustawZrodlo,
    wchlonSegmenty: magazyn.ustawSegmenty,
    panelJezykow: magazyn.panele,
    wchlonPanel: magazyn.wchlonPanel,
    wchlonPanele: magazyn.ustawPanele,
    obserwuj: magazyn.obserwuj,
    wczytajKontekst: (idSesji) => wczytajKontekst(magazyn, okna, kanaly, idSesji),

    rozlacz() {
      odsubskrybuj();
      odsubskrybujKanaly();
      kanaly.zapomnijSluchaczy();
      magazyn.zapomnijSluchaczy();
    },
  };
}

/**
 * Odczyt okna modułu z rdzenia buduje dwie komendy jedna po drugiej; niepowodzenie odczytu stanu
 * nie unieważnia wybranego okna, bo identyfikator wystarcza modułowi do pracy.
 */
async function wczytajKontekst(
  magazyn: MagazynTranslate,
  okna: ZrodloOknaTranslate,
  kanaly: ZrodloKanalowTranslate,
  idSesji: string,
): Promise<void> {
  void kanaly.wczytaj();
  magazyn.ustawFaze('odczyt');
  const wykaz = await okna.okna(idSesji);
  if (!wykaz.udany || wykaz.wynik === undefined) {
    magazyn.ustawFaze('blad', wykaz.blad?.message ?? 'rdzeń nie podał powodu');
    return;
  }
  const wybrane = wybierzOkno(wykaz.wynik.windows);
  magazyn.ustawOkno(wybrane);
  magazyn.ustawFaze('gotowe');
  if (wybrane === null) return;

  const stan = await okna.stan(wybrane.id);
  if (!stan.udany || stan.wynik === undefined) {
    magazyn.ustawFaze('gotowe', stan.blad?.message ?? '');
    return;
  }
  magazyn.ustawOkno(stan.wynik.window);
}

/**
 * Okno modułu Translate spośród okien sesji ma pierwszeństwo przypisane do modułu translate przez
 * rdzeń, a bez takiego okna bierze się pierwsze okno sesji.
 */
function wybierzOkno(okna: readonly Window[]): Window | null {
  return okna.find((okno) => okno.moduleId === 'translate') ?? okna[0] ?? null;
}
