import { utworzKontekstCzytelnosci, znacznikMowcy } from './czytelnosc-glosow';
import { ostatniaWypowiedz, wyciszony, zlozGlosBiezacy } from './glos-biezacy';
import { nazwaUczestnika, type StanDebaty } from './stan-debaty';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';

/** Zestawia udział uczestników debaty w turze bieżącej na podstawie liczby i objętości wypowiedzi widocznych przez okno, jedną pozycją na uczestnika. */
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
    // Wypowiedź otwarta liczy się z widocznej treści, nie z content, bo rdzeń dopisuje ją strumieniem.
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
 * Zdanie opisowe o zmierzonym udziale w turze, jawnie odróżnione od progu quorum przyjmowanego
 * przy głosowaniu.
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

/** Zestawienie udziału w formacie Markdown, jedyny raport tego modułu mający pokrycie w dostępnych danych źródłowych. */
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
