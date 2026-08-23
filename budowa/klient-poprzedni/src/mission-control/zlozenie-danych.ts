import {
  ProgressStatus,
  SessionStatus,
  WindowStatus,
  type ProgressChangedEvent,
  type Session,
  type SessionPresence,
  type Window,
} from '../../../shared/contract';
import {
  ZrodloDanych,
  type AgentZespolu,
  type DanePulpitu,
  type KolumnaSrodowiska,
  type SesjaMatrycy,
} from './model-danych';
import { pustyStanZrodla, type StanZrodla } from './stan-zrodla';

/**
 * Złożenie kompletu `DanePulpitu` z surowych odczytów i zdarzeń rdzenia.
 *
 * Czyste przełożenie stanu zebranego przez `zrodlo-pulpitu.ts` na kształt,
 * który widok umie wyrysować. Każda wartość pochodzi z odczytu kontraktu,
 * a miara bez źródła zostaje `null` i widok pokazuje ją jako stan pusty.
 *
 * Kolumny matrycy pochodzą z rdzenia w całości — kod, nazwa, motto i kolejność
 * z odczytu `environment.list`. Przed pierwszą odpowiedzią kolumn nie ma:
 * nagłówek pulpitu mówi „oczekiwanie na rdzeń", a matryca pokazuje ten sam
 * stan zamiast zastępczego kompletu kolumn zbudowanego z kodów modułów.
 */

/** Komplet sprzed pierwszego odczytu — stan oczekiwania, bez wartości zastępczych. */
export function pustyKomplet(): DanePulpitu {
  return zlozDanePulpitu(pustyStanZrodla());
}

/**
 * Motto kolumny, gdy rdzeń go nie podał — jawna, pusta wartość.
 *
 * Pole `Environment.motto` jest w kontrakcie, więc wartość pusta znaczy „ten
 * wiersz nie ma motta", a nie brak metadanej. Pustego motta nie zastępuje się
 * tekstem ułożonym po stronie widoku; treść karty środowiska strony głównej
 * jest osobna i nie jest tu powielana.
 */
const MOTTO_BEZ_ZRODLA = '';

/** Komplet pulpitu z surowego stanu — funkcja czysta. */
export function zlozDanePulpitu(stan: StanZrodla): DanePulpitu {
  const oknaSesji = pogrupujOknaPoSesji(stan.okna);
  const { matryca, pozaSrodowiskami } = zlozMatryce(stan, oknaSesji);

  const procesyBiegnace = [...stan.procesy.values()].filter(
    (proces) => proces.status === ProgressStatus.Running,
  );
  const oknaZProcesem = new Set(
    procesyBiegnace.flatMap((proces) => (proces.windowId === undefined ? [] : [proces.windowId])),
  );

  return {
    zrodloDanych: stan.odczytano ? ZrodloDanych.Rdzen : ZrodloDanych.Oczekiwanie,
    aktywnosc: {
      aktywneProcesy: procesyBiegnace.length,
      agenciPracuja: oknaZProcesem.size,
      modeleAnalizuja: null,
      walidatoryOczekuja: null,
    },
    decyzje: null,
    kanaly: stan.kanaly.map((kanal) => ({
      id: kanal.id,
      nazwa: kanal.name,
      rodzaj: kanal.kind,
      model: kanal.model ?? null,
      czynny: kanal.enabled,
    })),
    matryca,
    pozaSrodowiskami,
    relacje: null,
    procesy: [...stan.procesy.values()].map((proces) => zlozProces(proces, stan.okna)),
    kolejki: [...stan.kolejki.values()].map((kolejka) => ({
      id: kolejka.id,
      nazwa: kolejka.name ?? kolejka.id,
      status: kolejka.status,
      okna: kolejka.windowIds?.length ?? null,
      obiegi: kolejka.cycle ?? null,
    })),
    zespol: zlozZespol(stan.okna, oknaZProcesem, stan.procesy),
  };
}

/** Okna pogrupowane po sesji nadrzędnej. */
function pogrupujOknaPoSesji(okna: Map<string, Window>): Map<string, Window[]> {
  const grupy = new Map<string, Window[]>();
  for (const okno of okna.values()) {
    const grupa = grupy.get(okno.sessionId) ?? [];
    grupa.push(okno);
    grupy.set(okno.sessionId, grupa);
  }
  return grupy;
}

/** Kolumny środowisk i wykaz sesji, których środowiska odczyt nie wskazał. */
function zlozMatryce(
  stan: StanZrodla,
  oknaSesji: Map<string, Window[]>,
): { matryca: KolumnaSrodowiska[]; pozaSrodowiskami: SesjaMatrycy[] } {
  // Kolejność ustawia rdzeń: `Environment.order` odpowiada kolumnie
  // `srodowisko.kolejnosc`. Klient nie sortuje po nazwie ani po kodzie, bo
  // porządek kart jest zapisany w bazie.
  const matryca: KolumnaSrodowiska[] = [...stan.srodowiska.values()]
    .sort((pierwsze, drugie) => pierwsze.order - drugie.order)
    .map((srodowisko) => ({
      id: srodowisko.code,
      // Nazwa z rdzenia; pusta wraca do kodu, tak jak sesja bez tytułu wraca
      // do identyfikatora — nigdy do napisu ułożonego tutaj.
      nazwa: srodowisko.name !== '' ? srodowisko.name : srodowisko.code,
      // Motto z rdzenia (`environment.list`, pole `motto`). Brak pola i pole
      // puste znaczą to samo: kolumna bez motta, bez zastępnika.
      motto: srodowisko.motto ?? MOTTO_BEZ_ZRODLA,
      sesje: [],
    }));
  const kolumny = new Map(matryca.map((kolumna) => [kolumna.id, kolumna]));
  const pozaSrodowiskami: SesjaMatrycy[] = [];

  for (const sesja of stan.sesje.values()) {
    if (sesja.status === SessionStatus.Archived) continue;
    const obecnosc = stan.obecnosc.get(sesja.id);
    const pozycja = zlozSesje(sesja, obecnosc, oknaSesji.get(sesja.id) ?? []);
    const kolumna =
      obecnosc?.environmentCode === undefined ? undefined : kolumny.get(obecnosc.environmentCode);
    if (kolumna === undefined) {
      pozaSrodowiskami.push(pozycja);
    } else {
      kolumna.sesje.push(pozycja);
    }
  }

  return { matryca, pozaSrodowiskami };
}

/** Jedna sesja matrycy z odczytu sesji, jej obecności i okien. */
function zlozSesje(
  sesja: Session,
  obecnosc: SessionPresence | undefined,
  okna: Window[],
): SesjaMatrycy {
  return {
    id: sesja.id,
    tytul: sesja.title ?? sesja.id,
    status: sesja.status,
    rolaOkna: rolaOknaWiodacego(obecnosc, okna),
    okna: okna.length > 0 ? okna.length : (sesja.windowIds?.length ?? 0),
    pracaWTle: (obecnosc?.live ?? false) && (obecnosc?.streamingWindowCount ?? 0) > 0,
  };
}

/** Rola okna wiodącego: okno ogniskowane, w braku — pierwsze otwarte. */
function rolaOknaWiodacego(
  obecnosc: SessionPresence | undefined,
  okna: Window[],
): SesjaMatrycy['rolaOkna'] {
  const ogniskowane = okna.find((okno) => okno.id === obecnosc?.focusedWindowId);
  if (ogniskowane !== undefined) return ogniskowane.windowRole;
  const otwarte = okna.find((okno) => okno.status === WindowStatus.Open);
  return otwarte?.windowRole ?? null;
}

/** Jeden proces w tle ze zdarzenia telemetrii i znanego okna. */
function zlozProces(
  proces: ProgressChangedEvent,
  okna: Map<string, Window>,
): DanePulpitu['procesy'][number] {
  const okno = proces.windowId === undefined ? undefined : okna.get(proces.windowId);
  return {
    id: proces.processId,
    nazwa: okno?.title ?? okno?.moduleId ?? proces.processId,
    sesjaId: okno?.sessionId ?? null,
    etapBiezacy: proces.currentStep,
    etapowRazem: proces.totalSteps,
    etykietaEtapu: proces.stepLabel ?? null,
    status: proces.status,
  };
}

/** Zespół agentów: otwarte okna komunikacji wraz z ich telemetrią. */
function zlozZespol(
  okna: Map<string, Window>,
  oknaZProcesem: Set<string>,
  procesy: Map<string, ProgressChangedEvent>,
): AgentZespolu[] {
  const etapyOkien = new Map<string, string>();
  for (const proces of procesy.values()) {
    if (proces.windowId !== undefined && proces.stepLabel !== undefined) {
      etapyOkien.set(proces.windowId, proces.stepLabel);
    }
  }

  return [...okna.values()]
    .filter((okno) => okno.status === WindowStatus.Open)
    .map((okno) => ({
      id: okno.id,
      nazwa: okno.title ?? okno.moduleId,
      zajecie: etapyOkien.get(okno.id) ?? null,
      rola: okno.windowRole,
      podagenci: null,
      czynny: oknaZProcesem.has(okno.id),
    }));
}
