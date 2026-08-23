import type { KrawedzGrafu, OpisGrafu } from './graf-krokow';

/**
 * Eksport mapy zależności do dokumentacji procesu.
 *
 * Opracowanie wymienia cztery postacie zapisu; trzy z nich powstają po stronie
 * okna, bo treść jest już tutaj: zapis wektorowy bierze się wprost z kanwy
 * (`graf-krokow.ts`), a zapisy DOT i Mermaid składa ten plik. Czwarta postać —
 * mapa rastrowa — wymagałaby przerysowania kanwy na płótno i zostaje poza
 * oknem; postać wektorowa niesie to samo, a otwiera się w każdej przeglądarce.
 *
 * Żadna z tych postaci nie potrzebuje komendy kontraktu i mieć jej nie musi.
 *
 * Plik jest czysty: nie dotyka dokumentu i nie woła rdzenia.
 */

/** Nazwa układu w zapisach eksportu; oba formaty wymagają nazwy grafu. */
const NAZWA_UKLADU = 'zaleznosci_automatyki';

/**
 * Mapa w zapisie DOT.
 *
 * Kroki ścieżki krytycznej dostają grupę własną, a nie barwę wpisaną wprost:
 * barwa należy do narzędzia, które ten zapis wyrysuje, i do jego motywu.
 */
export function zapisDot(opis: OpisGrafu): string {
  const wiersze = [`digraph ${NAZWA_UKLADU} {`, '  rankdir=LR;', '  node [shape=box];'];
  for (const wezel of opis.wezly) {
    const nazwa = wezel.nazwa === '' ? wezel.kod : wezel.nazwa;
    const naSciezce = opis.sciezkaKrytyczna.has(wezel.kod);
    wiersze.push(
      `  ${cudzyslowDot(wezel.kod)} [label=${cudzyslowDot(nazwa)}` +
        (naSciezce ? ', group=sciezka_krytyczna, penwidth=2' : '') +
        '];',
    );
  }
  for (const krawedz of opis.krawedzie) {
    wiersze.push(
      `  ${cudzyslowDot(krawedz.odKroku)} -> ${cudzyslowDot(krawedz.doKroku)}` +
        (krawedz.podpis === '' ? ';' : ` [label=${cudzyslowDot(krawedz.podpis)}];`),
    );
  }
  wiersze.push('}', '');
  return wiersze.join('\n');
}

/**
 * Mapa w zapisie Mermaid.
 *
 * Identyfikatory kroków bywają w kontrakcie dowolnym napisem, a Mermaid czyta
 * identyfikator węzła jako nazwę bez cudzysłowu — dlatego identyfikator idzie
 * przez zastępnik złożony z liter i cyfr, a treść pierwotna zostaje w podpisie.
 */
export function zapisMermaid(opis: OpisGrafu): string {
  const zastepniki = new Map<string, string>();
  opis.wezly.forEach((wezel, miejsce) => zastepniki.set(wezel.kod, `krok${miejsce}`));

  const wiersze = ['flowchart LR'];
  for (const wezel of opis.wezly) {
    const zastepnik = zastepniki.get(wezel.kod) ?? wezel.kod;
    const nazwa = wezel.nazwa === '' ? wezel.kod : wezel.nazwa;
    wiersze.push(`  ${zastepnik}["${cudzyslowMermaid(nazwa)}"]`);
  }
  for (const krawedz of opis.krawedzie) {
    const od = zastepniki.get(krawedz.odKroku);
    const doo = zastepniki.get(krawedz.doKroku);
    if (od === undefined || doo === undefined) continue;
    wiersze.push(
      krawedz.podpis === ''
        ? `  ${od} --> ${doo}`
        : `  ${od} -- "${cudzyslowMermaid(krawedz.podpis)}" --> ${doo}`,
    );
  }
  const krytyczne = opis.wezly
    .filter((wezel) => opis.sciezkaKrytyczna.has(wezel.kod))
    .map((wezel) => zastepniki.get(wezel.kod) ?? wezel.kod);
  if (krytyczne.length > 0) {
    wiersze.push(`  classDef sciezkaKrytyczna stroke-width:2px;`);
    wiersze.push(`  class ${krytyczne.join(',')} sciezkaKrytyczna;`);
  }
  wiersze.push('');
  return wiersze.join('\n');
}

/** Podpis krawędzi: rodzaj zależności, a przy warunkowej także jej warunek. */
export function podpisKrawedzi(rodzaj: string, warunek?: string): KrawedzGrafu['podpis'] {
  return warunek === undefined || warunek.trim() === '' ? rodzaj : `${rodzaj}: ${warunek.trim()}`;
}

/** Napis w cudzysłowie zapisu DOT; cudzysłów i ukośnik w treści są chronione. */
function cudzyslowDot(tresc: string): string {
  return `"${tresc.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"`;
}

/** Treść podpisu Mermaid; cudzysłów zamykałby podpis w połowie. */
function cudzyslowMermaid(tresc: string): string {
  return tresc.replace(/"/g, "'");
}
