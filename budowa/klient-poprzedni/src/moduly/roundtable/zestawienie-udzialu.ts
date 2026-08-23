import { utworzKontekstCzytelnosci, znacznikMowcy } from './czytelnosc-glosow';
import { ostatniaWypowiedz, wyciszony, zlozGlosBiezacy } from './glos-biezacy';
import { nazwaUczestnika, type StanDebaty } from './stan-debaty';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';

/**
 * Zestawienie udziału uczestników w turze — jedyna ocena, którą da się w tym
 * module zmierzyć.
 *
 * Opracowanie modułu opisuje w oknie Voting & Evaluation Center głosowania,
 * rubryki, ocenę modelami-sędziami i ranking akumulowany między sesjami.
 * Wszystkie cztery są dziś w kontrakcie — `roundtable.vote.start`, `.cast`,
 * `.get`, `roundtable.rubric.set`, `roundtable.judge.run`,
 * `roundtable.leaderboard.get` — i żadnej z tych komend moduł jeszcze nie
 * wywołuje. Głosowanie zbudowane bez obsługi byłoby czynnością, której wynik
 * nie dociera nigdzie — czyli pozorem oceny.
 *
 * Zmierzyć da się udział: ilu uczestników składu zabrało głos w turze, ile razy
 * i jak obszernie. To nie jest ocena jakości wypowiedzi i tak jest nazwane
 * w oknie — udział, nie ranking, i nie quorum.
 */

/** Udział jednego uczestnika w turze bieżącej. */
export interface UdzialUczestnika {
  idUczestnika: string;
  znacznik: string;
  nazwa: string;
  /** Wypowiedzi tego uczestnika widziane przez okno w turze bieżącej. */
  wypowiedzi: number;
  znakow: number;
  wyciszony: boolean;
  /** Miejsce w kolejności głosu wskazane przez rdzeń; `null`, gdy nie wskazał. */
  kolejnosc: number | null;
  /** Czy rdzeń oznaczył uczestnika jako kluczowego. */
  kluczowy: boolean;
  /** Nazwa stanu głosu — wprost z `glos-biezacy.ts`. */
  stanGlosu: string;
}

export interface ZestawienieUdzialu {
  wiersze: UdzialUczestnika[];
  /** Uczestnicy składu znanego oknu. */
  skladu: number;
  /** Ilu z nich zabrało głos w turze bieżącej. */
  mowiacych: number;
  wypowiedzi: number;
  znakow: number;
}

export function zlozZestawienieUdzialu(
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
): ZestawienieUdzialu {
  const wypowiedziTury = stan.wypowiedzi();
  const kontekst = utworzKontekstCzytelnosci();

  const wiersze = stan.uczestnicy().map((uczestnik) => {
    const swoje = wypowiedziTury.filter((wypowiedz) => wypowiedz.participantId === uczestnik.id);
    const cisza = wyciszony(uczestnik);
    const ostatnia = ostatniaWypowiedz(wypowiedziTury, uczestnik.id);
    const biezacy = zlozGlosBiezacy(strumien.glos(uczestnik.id), ostatnia, cisza);
    // Wypowiedź otwarta ostatnio liczy się z treści widocznej w oknie, a nie
    // z pola `content`: rdzeń rozgłasza ją najpierw pustą i dopisuje słowa
    // strumieniem, więc sam zapis dawałby zero przez cały czas mówienia modelu.
    // Wypowiedzi wcześniejsze są już utrwalone, więc liczą się z zapisu.
    const znakow = swoje.reduce(
      (suma, wypowiedz) =>
        suma +
        (ostatnia !== null && wypowiedz.id === ostatnia.id
          ? biezacy.tekst.length
          : wypowiedz.content.length),
      0,
    );
    return {
      idUczestnika: uczestnik.id,
      znacznik: znacznikMowcy(uczestnik.id, stan, kontekst),
      nazwa: nazwaUczestnika(uczestnik, stan.opisKanalu(uczestnik.channelId)),
      wypowiedzi: swoje.length,
      znakow,
      wyciszony: cisza,
      kolejnosc: uczestnik.order ?? null,
      kluczowy: uczestnik.key === true,
      stanGlosu: biezacy.stan,
    };
  });

  return {
    wiersze,
    skladu: wiersze.length,
    mowiacych: wiersze.filter((udzial) => udzial.wypowiedzi > 0).length,
    wypowiedzi: wypowiedziTury.length,
    znakow: wiersze.reduce((suma, udzial) => suma + udzial.znakow, 0),
  };
}

/**
 * Zdanie o udziale zmierzonym — jawnie odróżnione od progu quorum.
 *
 * Progu zgody nie niesie żadne pole kontraktu, więc zdanie nie orzeka
 * o osiągnięciu czegokolwiek; podaje stosunek i mówi, czego w nim nie ma.
 */
export function zdanieUdzialu(zestawienie: ZestawienieUdzialu): string {
  if (zestawienie.skladu === 0) {
    return 'Debata nie ma jeszcze uczestników znanych temu oknu, więc udziału nie ma z czego policzyć.';
  }
  return (
    `Głos w tej turze zabrało ${zestawienie.mowiacych} z ${zestawienie.skladu} uczestników składu znanego oknu ` +
    `(${zestawienie.wypowiedzi} wypowiedzi, ${zestawienie.znakow} znaków). ` +
    'Jest to udział zmierzony, nie kworum: próg zgody niesie pole quorum żądania roundtable.vote.start, którego okno jeszcze nie wywołuje.'
  );
}

/** Zestawienie udziału w Markdown — jedyny raport, który ma pokrycie w danych. */
export function zestawienieMarkdown(
  zestawienie: ZestawienieUdzialu,
  opisZakresu: string,
): string {
  const wiersze = [
    '# Zestawienie udziału w debacie',
    '',
    `Dotyczy: ${opisZakresu}.`,
    '',
    '_Raport nie zawiera wyników głosowań ani ocen rubrykowych: ich komendy są w kontrakcie, ale okno jeszcze ich nie wywołuje. Zawiera udział zmierzony z wypowiedzi widzianych przez okno od jego otwarcia._',
    '',
    zdanieUdzialu(zestawienie),
    '',
    '| Znacznik | Uczestnik | Wypowiedzi | Znaków | Kolejność głosu | Stan w turze | Stan głosu |',
    '| --- | --- | --- | --- | --- | --- | --- |',
  ];
  for (const udzial of zestawienie.wiersze) {
    wiersze.push(
      `| ${udzial.znacznik} | ${udzial.nazwa} | ${udzial.wypowiedzi} | ${udzial.znakow} | ` +
        `${udzial.kolejnosc === null ? 'nie wskazana przez rdzeń' : `#${udzial.kolejnosc}`} | ` +
        `${udzial.wyciszony ? 'wyciszony' : 'słyszany'} | ${udzial.stanGlosu} |`,
    );
  }
  wiersze.push('');
  return wiersze.join('\n');
}
