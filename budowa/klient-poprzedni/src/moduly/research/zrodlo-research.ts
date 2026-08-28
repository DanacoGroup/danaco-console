import {
  Command,
  EventType,
  type RequestOf,
  type ResponseOf,
  type ResearchFindingAddRequest,
  type ResearchFindingAddResponse,
  type ResearchReportBuildRequest,
  type ResearchReportBuildResponse,
  type ResearchReportChangedEvent,
  type ResearchReportExportRequest,
  type ResearchReportExportResponse,
  type ResearchSourceAddRequest,
  type ResearchSourceAddResponse,
  type ResearchWorkspaceSetRequest,
  type ResearchWorkspaceSetResponse,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../../protokol/kanal';
import { czyObiekt, czyTekst, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { utworzNasluchOdmow, type NasluchOdmow, type OdpowiedzBadania } from './nasluch-odmow';
import type {
  ZlecenieEksportu,
  ZlecenieRaportu,
  ZlecenieUstalenia,
  ZlecenieZakresu,
  ZlecenieZrodla,
} from './zlecenia-badania';

/**
 * Pięć komend obszaru research widzianych przez okna modułu: źródło nie ma stanu, jest
 * warstwą wywołań i sprawdzianu kształtu odpowiedzi rdzenia.
 */
export interface ZrodloResearch {
  ustawZakres(zlecenie: ZlecenieZakresu): Promise<OdpowiedzBadania<ResearchWorkspaceSetResponse>>;
  dodajZrodlo(zlecenie: ZlecenieZrodla): Promise<OdpowiedzBadania<ResearchSourceAddResponse>>;
  zapiszUstalenie(
    zlecenie: ZlecenieUstalenia,
  ): Promise<OdpowiedzBadania<ResearchFindingAddResponse>>;
  zlozRaport(zlecenie: ZlecenieRaportu): Promise<OdpowiedzBadania<ResearchReportBuildResponse>>;
  eksportuj(zlecenie: ZlecenieEksportu): Promise<OdpowiedzBadania<ResearchReportExportResponse>>;
  /** Wywołanie dowolnej komendy obszaru research wraz z rozpoznaniem odmowy nieznanej komendy. */
  wywolaj<K extends Command>(
    komenda: K,
    zadanie: RequestOf<K>,
  ): Promise<OdpowiedzBadania<ResponseOf<K>>>;
  /** Subskrypcja `research.report.changed` — pierwszego ze zdarzeń obszaru. */
  naZmianeRaportu(sluchacz: (tresc: ResearchReportChangedEvent) => void): Odsubskrybuj;
  rozlacz(): void;
}

export function utworzZrodloResearch(kanal: Kanal): ZrodloResearch {
  const nasluch: NasluchOdmow = utworzNasluchOdmow(kanal);

  return {
    async ustawZakres(zlecenie) {
      const zadanie: ResearchWorkspaceSetRequest = { scope: zlecenie.zakres.trim() };
      if (zlecenie.etapy.length > 0) zadanie.stages = [...zlecenie.etapy];
      return sprawdz(
        await nasluch.wyslij(Command.ResearchWorkspaceSet, zadanie),
        Command.ResearchWorkspaceSet,
        (tresc) => czyTekst(tresc.scope),
      );
    },

    async dodajZrodlo(zlecenie) {
      return sprawdz(
        await nasluch.wyslij(Command.ResearchSourceAdd, zadanieZrodla(zlecenie)),
        Command.ResearchSourceAdd,
        (tresc) => czyObiekt(tresc.source),
      );
    },

    async zapiszUstalenie(zlecenie) {
      return sprawdz(
        await nasluch.wyslij(Command.ResearchFindingAdd, zadanieUstalenia(zlecenie)),
        Command.ResearchFindingAdd,
        (tresc) => czyObiekt(tresc.finding),
      );
    },

    async zlozRaport(zlecenie) {
      return sprawdz(
        await nasluch.wyslij(Command.ResearchReportBuild, zadanieRaportu(zlecenie)),
        Command.ResearchReportBuild,
        (tresc) => czyObiekt(tresc.report),
      );
    },

    async eksportuj(zlecenie) {
      const zadanie: ResearchReportExportRequest = {
        reportId: zlecenie.idRaportu,
        format: zlecenie.format,
        toLibrary: zlecenie.doRepozytorium,
      };
      if (zlecenie.sciezka.trim() !== '') zadanie.targetPath = zlecenie.sciezka.trim();
      return sprawdz(
        await nasluch.wyslij(Command.ResearchReportExport, zadanie),
        Command.ResearchReportExport,
        (tresc) => czyTekst(tresc.format),
      );
    },

    wywolaj(komenda, zadanie) {
      return nasluch.wyslij(komenda, zadanie);
    },

    naZmianeRaportu(sluchacz) {
      return kanal.naZdarzenie(EventType.ResearchReportChanged, (tresc) => sluchacz(tresc));
    },

    rozlacz: () => nasluch.rozlacz(),
  };
}

/**
 * Sprawdzian kształtu zachowujący nazwę nieznanego typu, którego sprawdzian kształtu warstwy
 * protokołu nie zna, by okno mogło ją wypisać.
 */
function sprawdz<T>(
  odpowiedz: OdpowiedzBadania<T>,
  komenda: string,
  sprawdzian: (tresc: T) => boolean,
): OdpowiedzBadania<T> {
  const wynik = sprawdzKsztalt(odpowiedz, komenda, sprawdzian);
  if (wynik.udany || odpowiedz.nieznanyTyp === undefined) return wynik;
  return { ...wynik, nieznanyTyp: odpowiedz.nieznanyTyp };
}

/** Treść żądania dodania źródła złożona ze zlecenia okna; pola puste nie idą do rdzenia jako treść żądania. */
function zadanieZrodla(zlecenie: ZlecenieZrodla): ResearchSourceAddRequest {
  const zadanie: ResearchSourceAddRequest = {
    windowId: zlecenie.idOkna,
    title: zlecenie.tytul.trim(),
    kind: zlecenie.rodzaj,
    credibility: zlecenie.wiarygodnosc,
  };
  if (zlecenie.adres.trim() !== '') zadanie.url = zlecenie.adres.trim();
  if (zlecenie.pochodzenie.trim() !== '') zadanie.origin = zlecenie.pochodzenie.trim();
  if (zlecenie.idPlikuRepozytorium.trim() !== '') {
    zadanie.libraryFileId = zlecenie.idPlikuRepozytorium.trim();
  }
  return zadanie;
}

/** Treść żądania zapisu ustalenia złożona ze zlecenia okna; brak identyfikatora ustalenia zakłada wpis nowy. */
function zadanieUstalenia(zlecenie: ZlecenieUstalenia): ResearchFindingAddRequest {
  const zadanie: ResearchFindingAddRequest = {
    windowId: zlecenie.idOkna,
    content: zlecenie.tresc.trim(),
    status: zlecenie.stan,
  };
  if (zlecenie.idUstalenia !== '') zadanie.findingId = zlecenie.idUstalenia;
  if (zlecenie.idZrodel.length > 0) zadanie.sourceIds = [...zlecenie.idZrodel];
  return zadanie;
}

/** Treść żądania złożenia raportu złożona ze zlecenia okna; brak identyfikatora raportu zakłada raport nowy. */
function zadanieRaportu(zlecenie: ZlecenieRaportu): ResearchReportBuildRequest {
  const zadanie: ResearchReportBuildRequest = { windowId: zlecenie.idOkna };
  if (zlecenie.idRaportu !== '') zadanie.reportId = zlecenie.idRaportu;
  if (zlecenie.tytul.trim() !== '') zadanie.title = zlecenie.tytul.trim();
  if (zlecenie.idUstalen.length > 0) zadanie.findingIds = [...zlecenie.idUstalen];
  if (zlecenie.sekcje.length > 0) zadanie.sections = [...zlecenie.sekcje];
  return zadanie;
}
