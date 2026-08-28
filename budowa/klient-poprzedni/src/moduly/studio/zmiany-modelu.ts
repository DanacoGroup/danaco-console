import {
  StudioAuthor,
  StudioChangeDecision,
  type StudioTrackedChange,
} from '../../../../shared/contract';

/**
 * Zmiany śledzone naniesione na treść — wynik modelu widoczny w miejscu; odcinek treści niesie
 * tekst wraz z autorem zmiany, o ile jakiejś dotyczy.
 */
export interface OdcinekTresci {
  tekst: string;
  /** Zmiana, do której odcinek należy; `null` dla treści bez zmiany. */
  zmiana: StudioTrackedChange | null;
}

/** Zmiany oczekujące decyzji Operatora, uporządkowane w kolejności ich położenia w treści bieżącego dokumentu. */
export function zmianyOczekujace(
  zmiany: readonly StudioTrackedChange[],
): readonly StudioTrackedChange[] {
  return [...zmiany]
    .filter((zmiana) => zmiana.decision === StudioChangeDecision.Oczekuje)
    .sort((pierwsza, druga) => pierwsza.rangeStart - druga.rangeStart);
}

/**
 * Rozdziela treść na odcinki wedle zakresów zmian oczekujących: zakres spoza treści jest
 * pomijany, a przy nakładających się zakresach wygrywa pierwszy.
 */
export function rozdzielNaOdcinki(
  tresc: string,
  zmiany: readonly StudioTrackedChange[],
): OdcinekTresci[] {
  const znaki = [...tresc];
  const odcinki: OdcinekTresci[] = [];
  let pozycja = 0;

  for (const zmiana of zmianyOczekujace(zmiany)) {
    const od = zmiana.rangeStart;
    const doZnaku = zmiana.rangeEnd;
    if (od < pozycja || doZnaku > znaki.length || od > doZnaku) continue;
    if (od > pozycja) odcinki.push({ tekst: znaki.slice(pozycja, od).join(''), zmiana: null });
    odcinki.push({ tekst: znaki.slice(od, doZnaku).join(''), zmiana });
    pozycja = doZnaku;
  }
  if (pozycja < znaki.length) {
    odcinki.push({ tekst: znaki.slice(pozycja).join(''), zmiana: null });
  }
  if (odcinki.length === 0) odcinki.push({ tekst: tresc, zmiana: null });
  return odcinki;
}

/** Nazwa autora zmiany czytelna dla Operatora: rozróżnia zmianę wprowadzoną przez model od zmiany Operatora. */
export function nazwaAutora(autor: StudioAuthor): string {
  return autor === StudioAuthor.Model ? 'model' : 'Operator';
}

/**
 * Zdanie o zmianie — co model zrobił i co przyjęcie albo odrzucenie zrobi; liczby są policzone
 * z treści zmiany, nie oszacowane.
 */
export function opiszZmiane(zmiana: StudioTrackedChange): string {
  const przed = zmiana.before ?? '';
  const po = zmiana.after ?? '';
  const roznica = po.length - przed.length;
  const kierunek = roznica === 0 ? 'bez zmiany długości' : `${roznica > 0 ? '+' : ''}${roznica} znaków`;
  return (
    `${nazwaAutora(zmiana.author)} · ${zmiana.kind} · znaki ${zmiana.rangeStart}–${zmiana.rangeEnd} · ` +
    `przed ${przed.length}, po ${po.length} (${kierunek})`
  );
}

/** Zdanie zbiorcze dla paska stanu podające liczbę zmian oczekujących decyzji; brak zmian ma własne zdanie, nie zero. */
export function opiszZmiany(zmiany: readonly StudioTrackedChange[]): string {
  const oczekujace = zmianyOczekujace(zmiany);
  if (oczekujace.length === 0) return 'zmian modelu do decyzji: brak';
  const odModelu = oczekujace.filter((zmiana) => zmiana.author === StudioAuthor.Model).length;
  return `zmiany do decyzji: ${oczekujace.length} (modelu: ${odModelu})`;
}
