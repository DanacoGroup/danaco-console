import {
  StudioAuthor,
  StudioChangeDecision,
  type StudioTrackedChange,
} from '../../../../shared/contract';

/**
 * Zmiany śledzone naniesione na treść — wynik modelu widoczny w miejscu.
 *
 * Rdzeń rejestruje wynik operacji kontekstowej jako zmianę śledzoną autora
 * `model` o zakresie liczonym w znakach treści BIEŻĄCEJ. Treść ma więc już
 * w sobie tekst modelu, a zmiana mówi, który to fragment i co stało tam
 * wcześniej. Zadaniem tego pliku jest przełożyć to na odcinki treści, żeby
 * powierzchnia mogła oznaczyć fragment modelu w miejscu, a nie obok.
 *
 * Plik nie zna DOM ani rdzenia: wejściem jest napis i wykaz zmian, wyjściem
 * odcinki. Dzięki temu nakładanie sprawdza się bez stawiania okna.
 */

/** Jeden odcinek treści wraz z autorem zmiany, o ile jakiejś dotyczy. */
export interface OdcinekTresci {
  tekst: string;
  /** Zmiana, do której odcinek należy; `null` dla treści bez zmiany. */
  zmiana: StudioTrackedChange | null;
}

/** Zmiany oczekujące decyzji, w kolejności położenia w treści. */
export function zmianyOczekujace(
  zmiany: readonly StudioTrackedChange[],
): readonly StudioTrackedChange[] {
  return [...zmiany]
    .filter((zmiana) => zmiana.decision === StudioChangeDecision.Oczekuje)
    .sort((pierwsza, druga) => pierwsza.rangeStart - druga.rangeStart);
}

/**
 * Rozdziela treść na odcinki wedle zakresów zmian oczekujących.
 *
 * Zakresy sięgające poza treść są pomijane, a nie przycinane: zakres spoza
 * treści znaczy, że dokument zmienił się od czasu zarejestrowania zmiany, więc
 * oznaczenie fragmentu wskazywałoby niewłaściwe miejsce. Tak samo rozstrzyga
 * rdzeń przy decyzji o zmianie — treść zostawia nietkniętą.
 *
 * Zakresy nachodzące na siebie liczą się pierwszy wygrywa: dwie zmiany na tym
 * samym fragmencie to stan, którego rdzeń nie zakłada, a oznaczenie fragmentu
 * dwoma autorami naraz nie miałoby czego znaczyć.
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

/** Nazwa autora dla Operatora. */
export function nazwaAutora(autor: StudioAuthor): string {
  return autor === StudioAuthor.Model ? 'model' : 'Operator';
}

/**
 * Zdanie o zmianie — co model zrobił i co przyjęcie albo odrzucenie zrobi.
 *
 * Liczby są policzone z treści zmiany, nie oszacowane: tyle znaków stało przed,
 * tyle stoi po. Bez nich „zmiana modelu" nie mówi, czy chodzi o przecinek, czy
 * o przepisanie akapitu.
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

/** Zdanie zbiorcze dla paska stanu; brak zmian ma własne zdanie, nie zero. */
export function opiszZmiany(zmiany: readonly StudioTrackedChange[]): string {
  const oczekujace = zmianyOczekujace(zmiany);
  if (oczekujace.length === 0) return 'zmian modelu do decyzji: brak';
  const odModelu = oczekujace.filter((zmiana) => zmiana.author === StudioAuthor.Model).length;
  return `zmiany do decyzji: ${oczekujace.length} (modelu: ${odModelu})`;
}
