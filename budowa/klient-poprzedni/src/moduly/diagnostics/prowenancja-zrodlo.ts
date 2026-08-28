import {
  Command,
  type ModelCallTrace,
  type ProvenanceCallGetRequest,
  type ProvenanceCallGetResponse,
  type ProvenanceCallListRequest,
  type ProvenanceCallRateRequest,
  type ProvenanceTraceExportRequest,
  type ProvenanceTraceExportResponse,
} from '../../../../shared/contract';
import type { Kanal } from '../../protokol/kanal';
import { czyObiekt, czyTablica, czyTekst } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';
import { rozstrzygnij, type WynikDiagnostyki } from './zrodlo-diagnostics';

/**
 * Cztery komendy obszaru `provenance.*` widziane przez zakładkę Provenance
 * Explorer: wykaz wywołań kanału modelu, odczyt jednego wywołania, ocena
 * odpowiedzi nadana przez Operatora oraz wydanie śladu w formacie maszynowym.
 */
export interface ZrodloProwenancji {
  /** `provenance.call.list` — rejestr wywołań zawężony filtrami okna. */
  wykazWywolan(
    zadanie: ProvenanceCallListRequest,
  ): Promise<WynikDiagnostyki<{ calls: ModelCallTrace[]; total?: number; truncated?: boolean }>>;
  /** `provenance.call.get` — jedno wywołanie wraz z odcinkami i treścią. */
  odczytajWywolanie(
    zadanie: ProvenanceCallGetRequest,
  ): Promise<WynikDiagnostyki<ProvenanceCallGetResponse>>;
  /** `provenance.call.rate` — zapis oceny trafności odpowiedzi. */
  ocenWywolanie(zadanie: ProvenanceCallRateRequest): Promise<WynikDiagnostyki<ModelCallTrace>>;
  /** `provenance.trace.export` — wydanie śladów w formacie maszynowym. */
  wydajSlad(
    zadanie: ProvenanceTraceExportRequest,
  ): Promise<WynikDiagnostyki<ProvenanceTraceExportResponse>>;
}

export function utworzZrodloProwenancji(kanal: Kanal): ZrodloProwenancji {
  return {
    async wykazWywolan(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.ProvenanceCallList, zadanie),
        Command.ProvenanceCallList,
        (tresc) => czyTablica(tresc.calls),
        (tresc) => ({
          calls: tresc.calls,
          ...(tresc.total === undefined ? {} : { total: tresc.total }),
          ...(tresc.truncated === undefined ? {} : { truncated: tresc.truncated }),
        }),
      );
    },

    async odczytajWywolanie(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.ProvenanceCallGet, zadanie),
        Command.ProvenanceCallGet,
        // Odcinki są w kontrakcie obowiązkowe, więc ich brak jest odpowiedzią nieczytelną.
        (tresc) => czyObiekt(tresc.call) && czyTablica(tresc.spans),
        (tresc) => tresc,
      );
    },

    async ocenWywolanie(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.ProvenanceCallRate, zadanie),
        Command.ProvenanceCallRate,
        (tresc) => czyObiekt(tresc.call),
        (tresc) => tresc.call,
      );
    },

    async wydajSlad(zadanie) {
      return rozstrzygnij(
        await wywolaj(kanal, Command.ProvenanceTraceExport, zadanie),
        Command.ProvenanceTraceExport,
        // Format i licznik muszą przyjść; sama treść wydania bywa pusta legalnie.
        (tresc) => czyTekst(tresc.content) && czyTekst(tresc.format),
        (tresc) => tresc,
      );
    },
  };
}
