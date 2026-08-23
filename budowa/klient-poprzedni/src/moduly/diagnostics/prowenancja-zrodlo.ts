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
 * odpowiedzi nadana przez Operatora i wydanie śladu w formacie maszynowym.
 *
 * Osobne źródło, a nie dołożenie metod do `zrodlo-diagnostics.ts`: tamten plik
 * należy do rodziny `diagnostics.*`, a prowenancja jest własną rodziną komend
 * kontraktu. Rozstrzygnięcie odpowiedzi jest jednak to samo — `rozstrzygnij`
 * z tamtego pliku — bo rozróżnienie odmowy rdzenia od odpowiedzi nieczytelnej
 * i od odpowiedzi bez treści ma tu tę samą wagę: zakładka poświęcona jawności
 * pracy modeli nie może zamilczeć własnego potknięcia.
 *
 * Piąta komenda rodziny, `provenance.call.replay`, do tego źródła nie należy:
 * powtórzenie wywołania wydaje pieniądze Operatora i jest czynnością sprawczą
 * osobnego odcinka, nie odczytem. Zakładka mówi o niej wprost, zamiast wołać ją
 * po cichu.
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
        // Odcinki są w kontrakcie obowiązkowe i puste dla wywołania bez drzewa,
        // więc ich brak jest odpowiedzią nieczytelną, nie wywołaniem prostym.
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
        // Treść wydania bywa pusta legalnie (zakres bez wywołań), ale format
        // i licznik muszą przyjść — bez nich okno nie ma czego powiedzieć
        // o tym, co plik niesie.
        (tresc) => czyTekst(tresc.content) && czyTekst(tresc.format),
        (tresc) => tresc,
      );
    },
  };
}
