import type { RoundtableStatement, RoundtableTurn } from '../../../../shared/contract';
import {
  chwila,
  opisMowcy,
  utworzKontekstCzytelnosci,
  utworzWypowiedz,
} from './czytelnosc-glosow';
import {
  ostatniaWypowiedz,
  ostatnieGlosy,
  wyciszony,
  zlozGlosBiezacy,
  type GlosBiezacy,
} from './glos-biezacy';
import type { StanDebaty } from './stan-debaty';
import type { StanTresci } from './stany-okna';
import type { GlosNieprzypisany, StrumienWypowiedzi } from './strumien-wypowiedzi';

/**
 * Rysuje przebieg tury albo stan pustki — czysta funkcja danych, biorąca gromadzenie
 * strumienia tego okna debaty jako parametr.
 */
export function rysujDebate(
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  tresc: StanTresci,
  filtr: string,
): void {
  const definicja = stan.definicjaTury();
  const nieprzypisane = strumien.nieprzypisane();
  if (definicja === null && nieprzypisane.length === 0) {
    tresc.pusto('Żadnej tury nie uruchomiono jeszcze w tym oknie. Wypełnij pytanie i uruchom turę.');
    return;
  }
  const wszystkie = stan.wypowiedzi();
  const wypowiedzi =
    filtr === '' ? wszystkie : wszystkie.filter((w) => pasujeDoFiltra(w, stan, strumien, filtr));
  const widok = tresc.tresc();
  if (definicja !== null) widok.append(zbudujTure(definicja));
  widok.append(
    zbudujGlosyTury(wypowiedzi, stan, strumien),
    zbudujListeWypowiedzi(wypowiedzi, stan, strumien, wszystkie.length, filtr),
  );
  if (nieprzypisane.length > 0) widok.append(zbudujNieprzypisane(nieprzypisane));
}

/**
 * Transkrypt Markdown z tego, co okno ma w pamięci — nic więcej.
 *
 * Treść bierze się z głosu bieżącego, nie wprost z `content`: eksport w trakcie
 * mówienia modelu zapisywałby inaczej pustą wypowiedź jako pustą, a Operator
 * ma w oknie przed sobą jej słowa.
 */
export function transkryptMarkdown(stan: StanDebaty, strumien: StrumienWypowiedzi): string {
  const definicja = stan.definicjaTury();
  const wiersze = [
    '# Transkrypt debaty',
    '',
    '_Obejmuje wyłącznie wypowiedzi, które to okno widziało od otwarcia. Historię sprzed wejścia oddaje komenda roundtable.debate.get, ale okno jej jeszcze nie wywołuje._',
    '',
  ];
  if (definicja !== null) {
    wiersze.push(`Tura ${definicja.index} · format ${definicja.format} · stan ${definicja.status}`, '');
  }
  const wypowiedzi = stan.wypowiedzi();
  const rosnace = ostatnieGlosy(wypowiedzi);
  for (const wypowiedz of wypowiedzi) {
    const mowca = opisMowcy(wypowiedz.participantId, stan.uczestnik(wypowiedz.participantId), stan);
    const biezacy = rosnace.has(wypowiedz.id)
      ? glosDlaWypowiedzi(wypowiedz, stan, strumien, wypowiedzi)
      : null;
    const slowa = biezacy === null ? wypowiedz.content : biezacy.tekst;
    const ogon = biezacy !== null && biezacy.rosnie ? ' _(strumień jeszcze otwarty)_' : '';
    wiersze.push(`## ${mowca} — ${chwila(wypowiedz.createdAt)}${ogon}`, slowa, '');
  }
  for (const wpis of strumien.nieprzypisane()) {
    wiersze.push(`## Mówca nieustalony (strumień ${wpis.identyfikator})`, wpis.glos.tekst, '');
  }
  return wiersze.join('\n');
}

/**
 * Transkrypt w JSON — ta sama treść co Markdown, jawnie oznaczona jako zapis okna, nie zapis
 * rdzenia, do dalszego przetwarzania.
 */
export function transkryptJson(stan: StanDebaty, strumien: StrumienWypowiedzi): string {
  const definicja = stan.definicjaTury();
  const wypowiedzi = stan.wypowiedzi();
  const rosnace = ostatnieGlosy(wypowiedzi);
  const zapis = {
    zakres:
      'Wypowiedzi widziane przez okno debaty od jego otwarcia. Historię sprzed otwarcia oddaje komenda roundtable.debate.get; okno jej jeszcze nie wywołuje, bo obsługi odczytu nie zbudowano.',
    tura:
      definicja === null
        ? null
        : {
            identyfikator: definicja.id,
            numer: definicja.index,
            format: definicja.format,
            stan: definicja.status,
            zagadnienie: definicja.topic ?? null,
          },
    wypowiedzi: wypowiedzi.map((wypowiedz) => {
      const biezacy = rosnace.has(wypowiedz.id)
        ? glosDlaWypowiedzi(wypowiedz, stan, strumien, wypowiedzi)
        : null;
      return {
        identyfikator: wypowiedz.id,
        mowca: opisMowcy(wypowiedz.participantId, stan.uczestnik(wypowiedz.participantId), stan),
        identyfikatorMowcy: wypowiedz.participantId,
        chwila: chwila(wypowiedz.createdAt),
        tresc: biezacy === null ? wypowiedz.content : biezacy.tekst,
        strumienOtwarty: biezacy !== null && biezacy.rosnie,
      };
    }),
    glosyNieprzypisane: strumien.nieprzypisane().map((wpis) => ({
      identyfikatorStrumienia: wpis.identyfikator,
      tresc: wpis.glos.tekst,
    })),
  };
  return JSON.stringify(zapis, null, 2);
}

/** Głos bieżący mówcy tej wypowiedzi — zapis i strumień złożone w jedno, gotowe do wyświetlenia w bloku wypowiedzi. */
function glosDlaWypowiedzi(
  wypowiedz: RoundtableStatement,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  wszystkie: readonly RoundtableStatement[],
): GlosBiezacy {
  const uczestnik = stan.uczestnik(wypowiedz.participantId);
  return zlozGlosBiezacy(
    strumien.glos(wypowiedz.participantId),
    ostatniaWypowiedz(wszystkie, wypowiedz.participantId),
    wyciszony(uczestnik),
  );
}

/**
 * Filtr czyta także treść rosnącą. Bez tego wypowiedź, której słowa właśnie
 * płyną strumieniem, znikałaby z wykazu przy każdym filtrze — jej `content` jest
 * przez cały czas mówienia modelu pustym napisem, więc nie pasowałaby do niczego.
 */
function pasujeDoFiltra(
  wypowiedz: RoundtableStatement,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  filtr: string,
): boolean {
  if (wypowiedz.content.toLowerCase().includes(filtr)) return true;
  const glos = strumien.glos(wypowiedz.participantId);
  if (glos !== null && glos.tekst.toLowerCase().includes(filtr)) return true;
  const mowca = opisMowcy(wypowiedz.participantId, stan.uczestnik(wypowiedz.participantId), stan);
  return mowca.toLowerCase().includes(filtr);
}

function zbudujTure(tura: RoundtableTurn): HTMLElement {
  const element = document.createElement('div');
  element.className = 'dr-tura';
  element.dataset['stanTury'] = tura.status;
  const temat = tura.topic === undefined || tura.topic === '' ? '' : ` · ${tura.topic}`;
  element.textContent = `Tura ${tura.index} · ${tura.format} · ${tura.status}${temat}`;
  return element;
}

/**
 * Spis głosów tej tury — ilu mówców i ile od każdego, liczone wyłącznie z wypowiedzi, które
 * okno widziało.
 */
function zbudujGlosyTury(
  wypowiedzi: readonly RoundtableStatement[],
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
): HTMLElement {
  const spis = document.createElement('p');
  spis.className = 'dr-glosy-tury';

  const ile = new Map<string, number>();
  for (const wypowiedz of wypowiedzi) {
    ile.set(wypowiedz.participantId, (ile.get(wypowiedz.participantId) ?? 0) + 1);
  }
  const milczacy = stan.uczestnicy().filter((uczestnik) => !ile.has(uczestnik.id));
  const mowiacy = liczMowiacych(stan, strumien);
  // Liczba mówiących stoi pierwsza, bo to jedyna wartość wiersza, która zmienia się w trakcie tury.
  spis.dataset['mowiacych'] = String(mowiacy);
  const czolo = mowiacy === 0 ? 'Nikt nie mówi w tej chwili.' : `Mówi teraz: ${mowiacy}.`;

  const czesci = [...ile.entries()].map(([idMowcy, liczba]) => {
    const mowca = opisMowcy(idMowcy, stan.uczestnik(idMowcy), stan);
    return `${mowca} — ${liczba}`;
  });

  if (czesci.length === 0) {
    spis.textContent = `${czolo} Żaden głos tej tury nie dotarł jeszcze do tego okna.`;
    return spis;
  }
  const ogon =
    milczacy.length === 0
      ? ''
      : ` Bez wypowiedzi widzianej przez to okno: ${milczacy.length} ze składu.`;
  spis.textContent = `${czolo} Głosów w tej turze — ${czesci.join(' · ')}.${ogon}`;
  return spis;
}

/** Ilu uczestników ma w tej chwili strumień otwarty i niezerwany, licząc tylko głosy znane temu oknu debaty. */
function liczMowiacych(stan: StanDebaty, strumien: StrumienWypowiedzi): number {
  let ilu = 0;
  for (const uczestnik of stan.uczestnicy()) {
    const glos = strumien.glos(uczestnik.id);
    if (glos !== null && !glos.domkniety && glos.przyczyna === '') ilu += 1;
  }
  return ilu;
}

/**
 * Głosy, których identyfikatora nie zna ani skład, ani wykaz wypowiedzi tury — pokazane, nie
 * porzucone.
 */
function zbudujNieprzypisane(wpisy: readonly GlosNieprzypisany[]): HTMLElement {
  const lista = document.createElement('ul');
  lista.className = 'dr-wypowiedzi';
  lista.setAttribute('aria-label', 'Głosy o nieustalonym mówcy');
  for (const wpis of wpisy) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-wypowiedz';
    pozycja.dataset['rodzajGlosu'] = 'nieznany';
    pozycja.dataset['rolaMowcy'] = 'nieznany';

    const mowca = document.createElement('span');
    mowca.className = 'dr-wypowiedz__mowca';
    mowca.textContent = `Mówca nieustalony (identyfikator strumienia: ${wpis.identyfikator})`;

    const tresc = document.createElement('p');
    tresc.className = 'dr-wypowiedz__tresc';
    tresc.textContent = wpis.glos.tekst;

    const stanGlosu = document.createElement('p');
    stanGlosu.className = 'dr-wypowiedz__stan';
    stanGlosu.dataset['stanGlosu'] = wpis.glos.domkniety ? 'strumień domknięty' : 'mówi teraz';
    stanGlosu.textContent =
      'Ten fragment przyszedł z okna tej debaty, ale jego identyfikatora nie zna ani skład, ' +
      'ani wykaz wypowiedzi tury — przypisanie nastąpi samo, gdy rdzeń rozgłosi tożsamość.';

    pozycja.append(mowca, tresc, stanGlosu);
    lista.append(pozycja);
  }
  return lista;
}

/**
 * Wykaz wypowiedzi. Kontekst czytelności powstaje na jedno rysowanie: znaczniki
 * mówców i role wstęgi mają odpowiadać temu, co widać teraz, a nie narastać
 * między przerysowaniami.
 */
function zbudujListeWypowiedzi(
  wypowiedzi: readonly RoundtableStatement[],
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  ilosciowoZnanych: number,
  filtr: string,
): HTMLElement {
  const lista = document.createElement('ul');
  lista.className = 'dr-wypowiedzi';
  lista.setAttribute('aria-label', 'Wypowiedzi tury bieżącej');
  if (wypowiedzi.length === 0) {
    const pozycja = document.createElement('li');
    pozycja.className = 'dr-wypowiedz';
    pozycja.dataset['rodzajGlosu'] = 'nieznany';
    pozycja.textContent =
      filtr === ''
        ? 'Ta tura nie ma jeszcze żadnej wypowiedzi widzianej przez to okno.'
        : `Filtr „${filtr}" nie pasuje do żadnej z ${ilosciowoZnanych} widzianych wypowiedzi.`;
    lista.append(pozycja);
    return lista;
  }
  const kontekst = utworzKontekstCzytelnosci();
  // Treść rosnąca należy do wypowiedzi otwartej przez mówcę ostatnio, nie do każdej jego wypowiedzi.
  const rosnace = ostatnieGlosy(wypowiedzi);
  for (const wypowiedz of wypowiedzi) {
    const biezacy = rosnace.has(wypowiedz.id)
      ? glosDlaWypowiedzi(wypowiedz, stan, strumien, wypowiedzi)
      : undefined;
    lista.append(utworzWypowiedz(wypowiedz, stan, kontekst, biezacy));
  }
  return lista;
}
