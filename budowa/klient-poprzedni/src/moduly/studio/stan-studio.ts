import {
  ChangeKind,
  type StudioDocument,
  type StudioOperationScope,
  type StudioVersion,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import {
  migawkaPolStanu,
  przyjmijPropozycje,
  przywrocPolaStanu,
  pustePolaStanu,
  ustawZaznaczenieStanu,
  wchlonDokument,
  zakresSkuteczny,
  zakresZadany,
  zapomnijDokument,
  type ParaPorownania,
  type MigawkaDokumentu,
  type PolaStanu,
  type PropozycjaZmiany,
  type Zaznaczenie,
} from './pola-stanu';
import { utworzZrodloStudio, type ZrodloStudio } from './zrodlo-studio';

export type { MigawkaDokumentu, ParaPorownania, PropozycjaZmiany, Zaznaczenie } from './pola-stanu';

/**
 * Stan modułu Studio prowadzi jeden dokument czynny współdzielony przez pięć
 * okien modułu, z buforem treści roboczej oddzielonym od treści rdzenia
 * i odświeżeniem przez zdarzenie zmiany dokumentu zamiast odpytywania.
 */
export interface StanStudio {
  /** Źródło komend obszaru `studio.*` — okna wołają je wprost. */
  zrodlo: ZrodloStudio;
  /** Okno komunikacji sesji, w którym pracuje moduł; puste przed osadzeniem. */
  idOkna(): string;
  ustawOkno(idOkna: string): void;
  /** Dokument czynny albo `null`, gdy żadnego nie wczytano. */
  dokument(): StudioDocument | null;
  /** Treść robocza edytora — bufor lokalny przed zapisem. */
  trescRobocza(): string;
  ustawTresc(tresc: string): void;
  /** Zaznaczenie w edytorze; `null` znaczy zakres „cały dokument". */
  zaznaczenie(): Zaznaczenie | null;
  ustawZaznaczenie(zakres: Zaznaczenie | null): void;
  /** Zakres pokazywany przez ster panelu narzędzi — wybór ręczny albo ślad zaznaczenia. */
  zakresZadany(): StudioOperationScope;
  /** Zakres wysyłany do rdzenia operacją kontekstową; jedyne źródło dla panelu i edytora. */
  zakresSkuteczny(): StudioOperationScope;
  /** Zapisuje wybór ręczny zakresu; `null` wraca do podążania za zaznaczeniem. */
  ustawZakresReczny(zakres: StudioOperationScope | null): void;
  /** Wersje repozytorium sesji w kolejności oddanej przez rdzeń. */
  wersje(): readonly StudioVersion[];
  ustawWersje(wersje: readonly StudioVersion[]): void;
  /** Propozycja zmiany czekająca na decyzję; `null`, gdy żadnej nie ma. */
  propozycja(): PropozycjaZmiany | null;
  ustawPropozycje(propozycja: PropozycjaZmiany | null): void;
  /** Powód odmowy ostatniej operacji kontekstowej; `null`, gdy żadna nie odmówiła. */
  odmowaOperacji(): string | null;
  ustawOdmoweOperacji(powod: string | null): void;
  /** Para wersji wskazana do porównania; `null`, gdy żadnej nie wskazano. */
  paraPorownania(): ParaPorownania | null;
  ustawParePorownania(para: ParaPorownania | null): void;
  /** Treść ostatniej zmiany zaakceptowanej — źródło podglądu. */
  trescZaakceptowana(): string;
  /** Wciąga dokument oddany przez rdzeń i przyjmuje jego treść za bieżącą. */
  wchlon(dokument: StudioDocument): void;
  /** Przyjmuje propozycję: jej treść staje się treścią roboczą i podglądem. */
  przyjmijPropozycje(): void;
  /** Zdejmuje migawkę pól dokumentu czynnego do niezależnego przechowania przez okno. */
  migawka(): MigawkaDokumentu;
  /** Wstawia migawkę jako dokument czynny. */
  przywrocMigawke(migawka: MigawkaDokumentu): void;
  /** Powiadamia widoki o każdej zmianie stanu. */
  obserwuj(sluchacz: () => void): () => void;
  /** Odłącza subskrypcję zdarzeń rdzenia. */
  rozlacz(): void;
}

export function utworzStanStudio(kanal: Kanal): StanStudio {
  const zrodlo = utworzZrodloStudio(kanal);
  const sluchacze = new Set<() => void>();
  const pola: PolaStanu = pustePolaStanu();

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  const odsubskrybuj = zrodlo.naZmianeDokumentu((tresc) => {
    if (pola.dokument !== null && tresc.document.id !== pola.dokument.id) return;
    if (tresc.change === ChangeKind.Deleted) zapomnijDokument(pola);
    else wchlonDokument(pola, tresc.document);
    oglos();
  });

  return {
    zrodlo,
    idOkna: () => pola.okno,
    dokument: () => pola.dokument,
    trescRobocza: () => pola.robocza,
    zaznaczenie: () => pola.zakres,
    zakresZadany: () => zakresZadany(pola),
    zakresSkuteczny: () => zakresSkuteczny(pola),
    wersje: () => pola.historia,
    propozycja: () => pola.propozycja,
    trescZaakceptowana: () => pola.zaakceptowana,
    odmowaOperacji: () => pola.odmowaOperacji,
    paraPorownania: () => pola.paraPorownania,

    ustawOkno(idOkna) {
      if (pola.okno === idOkna) return;
      pola.okno = idOkna;
      oglos();
    },

    ustawTresc(tresc) {
      if (pola.robocza === tresc) return;
      pola.robocza = tresc;
      oglos();
    },

    ustawZaznaczenie(zakres) {
      ustawZaznaczenieStanu(pola, zakres);
      oglos();
    },

    ustawZakresReczny(zakres) {
      if (pola.zakresReczny === zakres) return;
      pola.zakresReczny = zakres;
      oglos();
    },

    ustawWersje(nowe) {
      pola.historia = nowe;
      oglos();
    },

    ustawPropozycje(nowa) {
      pola.propozycja = nowa;
      // Nadejście wyniku unieważnia odmowę poprzednią.
      if (nowa !== null) pola.odmowaOperacji = null;
      oglos();
    },

    ustawOdmoweOperacji(powod) {
      if (pola.odmowaOperacji === powod) return;
      pola.odmowaOperacji = powod;
      oglos();
    },

    ustawParePorownania(para) {
      pola.paraPorownania = para;
      oglos();
    },

    wchlon(dokument) {
      wchlonDokument(pola, dokument);
      oglos();
    },

    przyjmijPropozycje() {
      przyjmijPropozycje(pola);
      oglos();
    },

    migawka: () => migawkaPolStanu(pola),

    przywrocMigawke(migawka) {
      przywrocPolaStanu(pola, migawka);
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    rozlacz() {
      odsubskrybuj();
      sluchacze.clear();
    },
  };
}
