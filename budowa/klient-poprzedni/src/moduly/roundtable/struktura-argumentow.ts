import type { RoundtableStatement } from '../../../../shared/contract';
import {
  chwila,
  opisMowcy,
  utworzKontekstCzytelnosci,
  znacznikMowcy,
} from './czytelnosc-glosow';
import { ostatniaWypowiedz, ostatnieGlosy, wyciszony, zlozGlosBiezacy } from './glos-biezacy';
import type { StanDebaty } from './stan-debaty';
import type { StrumienWypowiedzi } from './strumien-wypowiedzi';

/**
 * Struktura zapisu debaty — to, co da się policzyć z wypowiedzi tury bieżącej.
 *
 * Opracowanie modułu opisuje w oknie Argument Map & Analysis graf
 * teza→argument→kontrargument→riposta. Kontrakt ma dziś wszystko, czego graf
 * wymaga: `RoundtableStatement` niesie `replyToId` i `speechAct`, a wezły
 * i krawędzie oddaje `roundtable.argument.list`. Brakuje obsługi — ani rdzeń
 * tych pól jeszcze nie wypełnia, ani to okno komendy grafu nie wywołuje.
 *
 * Dlatego plik liczy to, co da się policzyć z samego zapisu tury, i osobno
 * mierzy, ile wypowiedzi niesie już pola relacji:
 *   • strukturę chronologiczną — kto, kiedy i ile powiedział w tej turze,
 *   • zbieżność LEKSYKALNĄ wypowiedzi — zdania powtarzające się dosłownie
 *     u kilku mówców i zdania wyłącznie własne,
 *   • pokrycie pól relacji — ile wypowiedzi ma `replyToId`, a ile `speechAct`.
 *
 * Pomiar trzeci jest odpowiedzią na pytanie, którego okno nie ma jak zadać
 * kontraktowi: czy rdzeń pola relacji już wypełnia. Zero znaczy „jeszcze nie”,
 * a wartość niezerowa jest sygnałem, że graf da się zbudować na prawdziwych
 * krawędziach zamiast na chronologii.
 *
 * Zbieżność leksykalna nie jest zgodnością stanowisk i nie udaje jej być.
 * Dwa zdania o tym samym znaczeniu i innych słowach są dla tej miary różne,
 * a dwa zdania identyczne w brzmieniu bywają przytoczeniem cudzej tezy po to,
 * żeby ją obalić. Miara mówi, gdzie szukać, nie co znaleziono — i tak jest
 * nazwana w oknie.
 */

/**
 * Najkrótsze zdanie wchodzące do porównania, liczone w wyrazach.
 *
 * Zdania krótsze są w debacie potwierdzeniami i zwrotami wiążącymi („Zgadzam
 * się.”, „Z drugiej strony.”), więc zbiegałyby się u wszystkich mówców naraz
 * i zalałyby wykaz zbieżności treścią bez wartości rozpoznawczej.
 */
const MIN_WYRAZOW_ZDANIA = 4;

/** Węzeł struktury — jedna wypowiedź tury. */
export interface WezelStruktury {
  idWypowiedzi: string;
  idMowcy: string;
  /** Znacznik mówcy — ten sam, którym podpisuje go skład i przebieg debaty. */
  znacznik: string;
  mowca: string;
  chwila: string;
  znakow: number;
  /** Miejsce w chronologii tury, liczone od 1. */
  kolejnosc: number;
}

/** Zdanie powtórzone dosłownie u więcej niż jednego mówcy. */
export interface ZdanieWspolne {
  /** Postać, w jakiej padło po raz pierwszy — porównanie idzie po postaci znormalizowanej. */
  tresc: string;
  /** Mówcy, u których zdanie padło, w kolejności pierwszego wystąpienia. */
  mowcy: string[];
}

/** Para mówców wraz z liczbą zdań wspólnych — wiersz macierzy zbieżności. */
export interface ParaMowcow {
  pierwszy: string;
  drugi: string;
  wspolnych: number;
}

export interface StrukturaDebaty {
  wezly: WezelStruktury[];
  /** Ilu mówców odezwało się w turze wśród wypowiedzi widzianych przez okno. */
  mowcow: number;
  znakow: number;
  /** Zdania powtórzone u kilku mówców, od najczęstszych. */
  wspolne: ZdanieWspolne[];
  /** Liczba zdań wyłącznie własnych, po mówcy. */
  wylacznieWlasne: ReadonlyMap<string, number>;
  /** Pary mówców o niezerowej liczbie zdań wspólnych, od najliczniejszych. */
  pary: ParaMowcow[];
  /** Zdania odrzucone jako zbyt krótkie do porównania. */
  pominietychZdan: number;
  /** Wypowiedzi niosące `replyToId` — krawędź grafu wskazaną przez rdzeń. */
  zOdniesieniem: number;
  /** Wypowiedzi niosące `speechAct` — rodzaj jednostki argumentacyjnej. */
  zAktemMowy: number;
}

/**
 * Składa strukturę tury bieżącej.
 *
 * Treść bierze się z głosu bieżącego, nie wprost z `content`: rdzeń rozgłasza
 * wypowiedź najpierw pustą, a słowa nadaje strumieniem, więc analiza czytająca
 * sam zapis liczyłaby zera przez cały czas mówienia modeli.
 */
export function zlozStrukture(stan: StanDebaty, strumien: StrumienWypowiedzi): StrukturaDebaty {
  const wypowiedzi = stan.wypowiedzi();
  const kontekst = utworzKontekstCzytelnosci();
  const rosnace = ostatnieGlosy(wypowiedzi);

  const wezly: WezelStruktury[] = [];
  const trescPoMowcy = new Map<string, string[]>();
  let znakow = 0;

  wypowiedzi.forEach((wypowiedz, indeks) => {
    const tekst = trescWypowiedzi(wypowiedz, stan, strumien, wypowiedzi, rosnace.has(wypowiedz.id));
    znakow += tekst.length;
    wezly.push({
      idWypowiedzi: wypowiedz.id,
      idMowcy: wypowiedz.participantId,
      znacznik: znacznikMowcy(wypowiedz.participantId, stan, kontekst),
      mowca: opisMowcy(wypowiedz.participantId, stan.uczestnik(wypowiedz.participantId), stan),
      chwila: chwila(wypowiedz.createdAt),
      znakow: tekst.length,
      kolejnosc: indeks + 1,
    });
    const zebrane = trescPoMowcy.get(wypowiedz.participantId) ?? [];
    zebrane.push(tekst);
    trescPoMowcy.set(wypowiedz.participantId, zebrane);
  });

  const zbieznosc = policzZbieznosc(trescPoMowcy, (idMowcy) =>
    opisMowcy(idMowcy, stan.uczestnik(idMowcy), stan),
  );

  return {
    wezly,
    mowcow: trescPoMowcy.size,
    znakow,
    wspolne: zbieznosc.wspolne,
    wylacznieWlasne: zbieznosc.wylacznieWlasne,
    pary: zbieznosc.pary,
    pominietychZdan: zbieznosc.pominietychZdan,
    zOdniesieniem: wypowiedzi.filter((wypowiedz) => wypowiedz.replyToId !== undefined).length,
    zAktemMowy: wypowiedzi.filter((wypowiedz) => wypowiedz.speechAct !== undefined).length,
  };
}

/**
 * Zdanie o pokryciu pól relacji — mierzone, nie zakładane.
 *
 * Okno nie ma jak zapytać kontraktu, czy rdzeń wypełnia `replyToId`
 * i `speechAct`; policzenie ich w wypowiedziach, które już przyszły, jest
 * jedyną odpowiedzią opartą na faktach. Zdanie mówi wprost, co z tego wynika
 * dla grafu.
 */
export function zdaniePolRelacji(struktura: StrukturaDebaty): string {
  const ile = struktura.wezly.length;
  if (struktura.zOdniesieniem === 0 && struktura.zAktemMowy === 0) {
    return (
      `Pola relacji są w kontrakcie (RoundtableStatement.replyToId i .speechAct), ale rdzeń nie wypełnił ich ` +
      `w żadnej z ${ile} wypowiedzi widzianych przez to okno. Dopóki tak jest, krawędzi grafu nie ma z czego ` +
      'poprowadzić i okno pokazuje chronologię.'
    );
  }
  return (
    `Rdzeń wypełnia już pola relacji: odniesienie do innej wypowiedzi niesie ${struktura.zOdniesieniem} ` +
    `z ${ile} wypowiedzi, rodzaj aktu mowy — ${struktura.zAktemMowy}. Graf argumentów da się na nich ` +
    'zbudować; to okno jeszcze go nie rysuje i nie wywołuje komendy roundtable.argument.list.'
  );
}

/** Treść jednej wypowiedzi — rosnąca, gdy to ostatni głos mówcy, inaczej zapis. */
function trescWypowiedzi(
  wypowiedz: RoundtableStatement,
  stan: StanDebaty,
  strumien: StrumienWypowiedzi,
  wszystkie: readonly RoundtableStatement[],
  rosnaca: boolean,
): string {
  if (!rosnaca) return wypowiedz.content;
  const uczestnik = stan.uczestnik(wypowiedz.participantId);
  return zlozGlosBiezacy(
    strumien.glos(wypowiedz.participantId),
    ostatniaWypowiedz(wszystkie, wypowiedz.participantId),
    wyciszony(uczestnik),
  ).tekst;
}

/** Wynik pomiaru zbieżności leksykalnej. */
interface WynikZbieznosci {
  wspolne: ZdanieWspolne[];
  wylacznieWlasne: Map<string, number>;
  pary: ParaMowcow[];
  pominietychZdan: number;
}

/**
 * Zbieżność leksykalna wypowiedzi.
 *
 * Porównanie idzie po zdaniach sprowadzonych do jednej postaci: małe litery,
 * bez znaków przestankowych, bez wielokrotnych odstępów. Dzięki temu ta sama
 * teza zapisana z innym znakiem końca zdania liczy się raz, a nie dwa razy.
 */
function policzZbieznosc(
  trescPoMowcy: ReadonlyMap<string, readonly string[]>,
  opis: (idMowcy: string) => string,
): WynikZbieznosci {
  /** Zdanie znormalizowane → mówcy, u których padło, i postać pierwotna. */
  const wystapienia = new Map<string, { mowcy: Set<string>; tresc: string }>();
  let pominietychZdan = 0;

  for (const [idMowcy, fragmenty] of trescPoMowcy) {
    for (const zdanie of naZdania(fragmenty.join(' '))) {
      const klucz = znormalizuj(zdanie);
      if (klucz === '') continue;
      if (klucz.split(' ').length < MIN_WYRAZOW_ZDANIA) {
        pominietychZdan += 1;
        continue;
      }
      const wpis = wystapienia.get(klucz);
      if (wpis === undefined) {
        wystapienia.set(klucz, { mowcy: new Set([idMowcy]), tresc: zdanie });
        continue;
      }
      wpis.mowcy.add(idMowcy);
    }
  }

  const wspolne: ZdanieWspolne[] = [];
  const wylacznieWlasne = new Map<string, number>();
  const pary = new Map<string, ParaMowcow>();

  for (const wpis of wystapienia.values()) {
    const mowcy = [...wpis.mowcy];
    if (mowcy.length === 1) {
      const jedyny = mowcy[0]!;
      wylacznieWlasne.set(opis(jedyny), (wylacznieWlasne.get(opis(jedyny)) ?? 0) + 1);
      continue;
    }
    wspolne.push({ tresc: wpis.tresc, mowcy: mowcy.map(opis) });
    for (let i = 0; i < mowcy.length; i += 1) {
      for (let j = i + 1; j < mowcy.length; j += 1) {
        const pierwszy = opis(mowcy[i]!);
        const drugi = opis(mowcy[j]!);
        const klucz = `${pierwszy} ${drugi}`;
        const znana = pary.get(klucz);
        if (znana === undefined) pary.set(klucz, { pierwszy, drugi, wspolnych: 1 });
        else znana.wspolnych += 1;
      }
    }
  }

  wspolne.sort((pierwsze, drugie) => drugie.mowcy.length - pierwsze.mowcy.length);
  return {
    wspolne,
    wylacznieWlasne,
    pary: [...pary.values()].sort((pierwsza, druga) => druga.wspolnych - pierwsza.wspolnych),
    pominietychZdan,
  };
}

/** Podział treści na zdania po kropce, pytajniku, wykrzykniku i końcu wiersza. */
function naZdania(tresc: string): string[] {
  return tresc
    .split(/(?<=[.!?])\s+|\n+/u)
    .map((zdanie) => zdanie.trim())
    .filter((zdanie) => zdanie !== '');
}

/** Postać porównawcza zdania: małe litery, bez znaków przestankowych i nadmiaru odstępów. */
function znormalizuj(zdanie: string): string {
  return zdanie
    .toLocaleLowerCase('pl-PL')
    .replace(/[^\p{L}\p{N}\s]/gu, ' ')
    .replace(/\s+/gu, ' ')
    .trim();
}

/** Transkrypt struktury w Markdown — chronologia i zbieżność, bez udawanego grafu. */
export function strukturaMarkdown(struktura: StrukturaDebaty, opisTury: string): string {
  const wiersze = [
    '# Struktura zapisu debaty',
    '',
    `Dotyczy: ${opisTury}.`,
    '',
    '_Zapis obejmuje wyłącznie wypowiedzi widziane przez okno od jego otwarcia. Struktura jest chronologią, nie grafem argumentów: pola relacji są w kontrakcie, ale okno grafu jeszcze nie buduje._',
    '',
    `Wypowiedzi: ${struktura.wezly.length} · mówców: ${struktura.mowcow} · znaków treści: ${struktura.znakow}.`,
    '',
    zdaniePolRelacji(struktura),
    '',
    '## Chronologia tury',
    '',
  ];
  for (const wezel of struktura.wezly) {
    wiersze.push(
      `${wezel.kolejnosc}. [${wezel.znacznik}] ${wezel.mowca} — ${wezel.chwila} — ${wezel.znakow} znaków`,
    );
  }
  wiersze.push('', '## Zbieżność leksykalna wypowiedzi', '');
  if (struktura.wspolne.length === 0) {
    wiersze.push('Żadne zdanie nie powtórzyło się dosłownie u więcej niż jednego mówcy.');
  } else {
    for (const zdanie of struktura.wspolne) {
      wiersze.push(`- „${zdanie.tresc}” — mówcy: ${zdanie.mowcy.join(', ')}`);
    }
  }
  wiersze.push('', '## Zdania wyłącznie własne', '');
  if (struktura.wylacznieWlasne.size === 0) {
    wiersze.push('Brak zdań dłuższych od progu porównania.');
  } else {
    for (const [mowca, ile] of struktura.wylacznieWlasne) wiersze.push(`- ${mowca}: ${ile}`);
  }
  wiersze.push('');
  return wiersze.join('\n');
}
