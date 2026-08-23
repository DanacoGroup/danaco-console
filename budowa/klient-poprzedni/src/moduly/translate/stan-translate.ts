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
 * Stan modułu Translate — źródła kontraktu nad jednym zbiorem danych.
 *
 * Tekst źródłowy i panele są jednym bytem: zmiana źródła aktualizuje wszystkie
 * panele naraz, więc odpowiedź `translate.source.set` niesie od razu komplet
 * paneli. Gdyby Source Panel prowadził własną kopię tekstu, a Translation
 * Panels własną kopię paneli, ta jednoczesność musiałaby być odtwarzana ręcznie
 * w dwóch miejscach i rozjeżdżałaby się przy pierwszej odmowie.
 *
 * Identyfikator okna jest warunkiem wstępnym: bez niego moduł nie wysyła
 * `translate.source.set` ani `translate.target.add` — mówi o tym wprost zamiast
 * wysyłać żądanie z pustym polem i pokazywać odmowę walidacji jako własną
 * usterkę.
 */
export type { FazaKontekstu } from './magazyn-translate';

export interface StanTranslate {
  zrodlo: ZrodloZrodla;
  panele: ZrodloPaneli;
  glosariusz: ZrodloGlosariusza;
  okna: ZrodloOknaTranslate;
  /**
   * Dokument wejściowy i jego zamiana formatu — cudzy obszar kontraktu,
   * z którego korzysta Format Studio.
   *
   * Obszar `translate` nie ma ani jednej komendy przyjmującej plik, więc bez
   * tego źródła moduł nie miałby wejścia od strony dokumentu wcale.
   */
  dokument: ZrodloDokumentuTranslate;
  /**
   * Rejestr kanałów modelu — cudzy obszar, z którego biorą wartość stery
   * `channelId` przy dodaniu języka i przy tłumaczeniu zwrotnym.
   *
   * Jeden rejestr na moduł, nie jeden na ster. Sterów jest 1 + N (formularz
   * „+ Dodaj język" i po jednym w każdej instancji panelu); rejestr zakładany
   * osobno przez każdy z nich pytałby rdzeń N razy o ten sam wykaz.
   */
  kanaly: ZrodloKanalowTranslate;
  /**
   * Warsztat — wszystkie rodziny komend spoza czterech okien pierwotnych:
   * pamięć jako byt Operatora, segmentacja, terminologia, korekta, obieg,
   * dokumenty, lokalizacja, napisy, silniki i wymiana zewnętrzna.
   */
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

  // Powrót rejestru kanałów przerysowuje moduł tą samą drogą, co każda inna
  // zmiana — stery obsadzają się same, bez własnego nasłuchu w N instancjach.
  const odsubskrybujKanaly = kanaly.obserwuj(() => magazyn.ogloszZmiane());

  // Zdarzenie jest jedynym źródłem odświeżenia poza własnym działaniem: panel
  // przeliczony przez rdzeń dociera tą samą drogą co panel zmieniony ręcznie.
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
 * Odczyt okna modułu z rdzenia — dwie komendy zbudowane, jedna po drugiej.
 *
 * Niepowodzenie `window.state.get` nie unieważnia wybranego okna: identyfikator
 * wystarcza modułowi do pracy, a parametry wykonania są dodatkiem
 * informacyjnym (fail-open).
 *
 * Rejestr kanałów idzie obok kolejki komend. Wykaz kanałów modelu nie zależy
 * od okna i nie jest warunkiem żadnej czynności — obsadza wyłącznie ster
 * `channelId`, którego pusty wybór jest poprawną wartością. Czekanie na niego
 * opóźniałoby pas kontekstu, a jego odmowa nie zatrzymuje modułu, więc obietnica
 * jest tu świadomie nieoczekiwana. Powód niepowodzenia zapisuje sam rejestr
 * i mówi go ster.
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
 * Okno modułu Translate spośród okien sesji.
 *
 * Pierwszeństwo ma okno, które rdzeń przypisał do modułu `translate`. Gdy
 * takiego nie ma, bierzemy pierwsze okno sesji: moduł okna zmienia się komendą
 * `workspace.enter`, więc okno sesji bez modułu Translate i tak jest tym oknem,
 * w którym Operator właśnie pracuje.
 */
function wybierzOkno(okna: readonly Window[]): Window | null {
  return okna.find((okno) => okno.moduleId === 'translate') ?? okna[0] ?? null;
}
