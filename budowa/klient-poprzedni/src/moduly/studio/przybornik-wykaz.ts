import {
  StudioAuthor,
  StudioChangeDecision,
  type StudioAnnotation,
  type StudioComment,
  type StudioTrackedChange,
} from '../../../../shared/contract';
import type { PropozycjaZmiany } from './pola-stanu';
import { przybornikOpiszZnacznik, type ZnacznikWlasny } from './przybornik-znaczniki';

/** Stała RodzajZnakowania nazywa pięć rozłącznych rodzajów znakowania dokumentu: komentarz, propozycję, zmianę, adnotację i znacznik. */
export const RodzajZnakowania = {
  Komentarz: 'komentarz',
  Propozycja: 'propozycja',
  Zmiana: 'zmiana',
  Adnotacja: 'adnotacja',
  Znacznik: 'znacznik',
} as const;
export type RodzajZnakowania = (typeof RodzajZnakowania)[keyof typeof RodzajZnakowania];

/** Stała NAZWY_RODZAJOW niesie nazwę każdego rodzaju znakowania wraz ze zdaniem opisującym, czym różni się od pozostałych rodzajów. */
export const NAZWY_RODZAJOW: Readonly<Record<RodzajZnakowania, { nazwa: string; czym: string }>> = {
  komentarz: {
    nazwa: 'Komentarz',
    czym: 'Mówi o fragmencie i NIE niesie brzmienia — treści dokumentu nie zmienia ani teraz, ani po przyjęciu.',
  },
  propozycja: {
    nazwa: 'Propozycja zmiany',
    czym: 'Niesie brzmienie fragmentu, ale NIE weszła w treść. Operator ją przyjmuje, odrzuca albo poprawia.',
  },
  zmiana: {
    nazwa: 'Zmiana śledzona',
    czym: 'Jest JUŻ w treści dokumentu i czeka na decyzję: przyjęcie zostawia ją, odrzucenie wycofuje.',
  },
  adnotacja: {
    nazwa: 'Adnotacja różnicy',
    czym: 'Uwaga przy numerze fragmentu porównania wersji — nie przy miejscu w treści bieżącej.',
  },
  znacznik: {
    nazwa: 'Znacznik własny',
    czym: 'Nazwa i barwa nadana fragmentowi przez Operatora albo model; żyje przez sesję okna, bo kontrakt nie ma na nią pola.',
  },
};

/** Interfejs PozycjaZnakowania niesie jedną pozycję wykazu znakowań: kod, rodzaj, autora, zakres, tytuł, treść, podstawę i czas założenia. */
export interface PozycjaZnakowania {
  kod: string;
  rodzaj: RodzajZnakowania;
  autor: StudioAuthor;
  /** Czy pozycja jest otwarta — nierozwiązana, nierozstrzygnięta, nieodhaczona. */
  otwarta: boolean;
  /** Zakres w znakach treści; `null` znaczy „bez zakotwiczenia w treści". */
  zakres: { poczatek: number; koniec: number } | null;
  /** Jednowierszowe nazwanie pozycji. */
  tytul: string;
  /** Treść pozycji: komentarz, brzmienie propozycji, treść adnotacji. */
  tresc: string;
  /** Zdanie o pochodzeniu i stanie pozycji. */
  podstawa: string;
  /** Czas założenia w milisekundach epoki. */
  czas: number;
}

/** Interfejs FiltrZnakowan niesie zawężenie wykazu znakowań wedle trzech osi: rodzaju, autora i stanu otwarcia pozycji. */
export interface FiltrZnakowan {
  rodzaj: RodzajZnakowania | 'wszystkie';
  autor: StudioAuthor | 'wszyscy';
  stan: 'wszystkie' | 'otwarte' | 'zamkniete';
}

/** Funkcja przybornikFiltrPelny zwraca filtr niczego nie zawężający, będący stanem początkowym przybornika znakowań. */
export function przybornikFiltrPelny(): FiltrZnakowan {
  return { rodzaj: 'wszystkie', autor: 'wszyscy', stan: 'wszystkie' };
}

/** Interfejs MaterialZnakowan niesie materiał, z którego przybornik składa wykaz znakowań: komentarze, zmiany, adnotacje, znaczniki i propozycję. */
export interface MaterialZnakowan {
  komentarze: readonly StudioComment[];
  zmiany: readonly StudioTrackedChange[];
  adnotacje: readonly StudioAnnotation[];
  znaczniki: readonly ZnacznikWlasny[];
  /** Propozycja czekająca na decyzję; `null`, gdy żadnej nie ma. */
  propozycja: PropozycjaZmiany | null;
  /** Zaznaczenie, którego propozycja dotyczy — propozycja nie nosi zakresu sama. */
  zakresPropozycji: { poczatek: number; koniec: number } | null;
}

/**
 * Funkcja przybornikZlozZnakowania składa wykaz znakowań ze wszystkich pięciu źródeł, w kolejności położenia w treści; pozycje bez zakotwiczenia idą na koniec wykazu.
 */
export function przybornikZlozZnakowania(material: MaterialZnakowan): PozycjaZnakowania[] {
  const pozycje: PozycjaZnakowania[] = [];

  for (const komentarz of material.komentarze) {
    if (komentarz.parentCommentId !== undefined) continue;
    const odpowiedzi = material.komentarze.filter(
      (pozycja) => pozycja.parentCommentId === komentarz.id,
    ).length;
    pozycje.push({
      kod: komentarz.id,
      rodzaj: RodzajZnakowania.Komentarz,
      autor: komentarz.author,
      otwarta: komentarz.resolved !== true,
      zakres:
        komentarz.selectionStart === undefined
          ? null
          : {
              poczatek: komentarz.selectionStart,
              koniec: komentarz.selectionEnd ?? komentarz.selectionStart,
            },
      tytul: NAZWY_RODZAJOW.komentarz.nazwa,
      tresc: komentarz.body,
      podstawa:
        `studio.comment.add · odpowiedzi w wątku: ${odpowiedzi} · ` +
        (komentarz.resolved === true ? 'wątek rozwiązany' : 'wątek otwarty'),
      czas: komentarz.createdAt,
    });
  }

  if (material.propozycja !== null && material.propozycja.tresc !== '') {
    pozycje.push({
      kod:
        material.propozycja.idPropozycji === ''
          ? `propozycja-${material.propozycja.idAkcji}`
          : material.propozycja.idPropozycji,
      rodzaj: RodzajZnakowania.Propozycja,
      // Propozycja jest wynikiem operacji kontekstowej, a tę wykonuje model.
      autor: StudioAuthor.Model,
      otwarta: true,
      zakres: material.zakresPropozycji,
      tytul: NAZWY_RODZAJOW.propozycja.nazwa,
      tresc: material.propozycja.tresc,
      podstawa:
        `operacja ${material.propozycja.idAkcji} · ${material.propozycja.tresc.length} znaków · ` +
        (material.propozycja.idPropozycji === ''
          ? 'bez odwołania w rdzeniu — decyzja idzie treścią, nie identyfikatorem'
          : `studio.proposal.decide dla ${material.propozycja.idPropozycji}`) +
        ' · brzmienie stoi NA MARGINESIE, w treści go jeszcze nie ma',
      czas: Date.now(),
    });
  }

  for (const zmiana of material.zmiany) {
    pozycje.push({
      kod: zmiana.id,
      rodzaj: RodzajZnakowania.Zmiana,
      autor: zmiana.author,
      otwarta: zmiana.decision === StudioChangeDecision.Oczekuje,
      zakres: { poczatek: zmiana.rangeStart, koniec: zmiana.rangeEnd },
      tytul: NAZWY_RODZAJOW.zmiana.nazwa,
      tresc: zmiana.after ?? zmiana.before ?? '',
      podstawa:
        `studio.tracking.* · ${zmiana.kind} · przed ${(zmiana.before ?? '').length}, po ` +
        `${(zmiana.after ?? '').length} znaków · stan decyzji: ${zmiana.decision} · treść JUŻ ` +
        'zawiera tę zmianę',
      czas: zmiana.createdAt,
    });
  }

  for (const adnotacja of material.adnotacje) {
    pozycje.push({
      kod: adnotacja.id,
      rodzaj: RodzajZnakowania.Adnotacja,
      autor: adnotacja.author,
      otwarta: true,
      // Adnotacja wisi przy numerze fragmentu różnicy, nie przy znaku treści, więc zakres zostaje pusty.
      zakres: null,
      tytul: `${NAZWY_RODZAJOW.adnotacja.nazwa} — fragment ${adnotacja.hunkIndex}`,
      tresc: adnotacja.body,
      podstawa:
        `studio.annotation.add · fragment różnicy ${adnotacja.hunkIndex} · wersje ` +
        `${adnotacja.baseVersionId ?? 'bez odniesienia'} → ` +
        `${adnotacja.targetVersionId ?? adnotacja.proposalId ?? 'bieżąca'}`,
      czas: adnotacja.createdAt,
    });
  }

  for (const znacznik of material.znaczniki) {
    pozycje.push({
      kod: znacznik.kod,
      rodzaj: RodzajZnakowania.Znacznik,
      autor: znacznik.autor,
      otwarta: znacznik.otwarty,
      zakres: znacznik.zakres,
      tytul: `${NAZWY_RODZAJOW.znacznik.nazwa} — ${znacznik.nazwa}`,
      tresc: znacznik.nazwa,
      podstawa: `${przybornikOpiszZnacznik(znacznik)} · barwa ${znacznik.barwa}`,
      czas: znacznik.czas,
    });
  }

  return pozycje.sort(przybornikPorownajPolozeniem);
}

/** Funkcja przybornikPorownajPolozeniem porządkuje wykaz wedle położenia w treści, stawiając pozycje bez zakotwiczenia na końcu. */
function przybornikPorownajPolozeniem(
  pierwsza: PozycjaZnakowania,
  druga: PozycjaZnakowania,
): number {
  if (pierwsza.zakres === null && druga.zakres === null) return pierwsza.czas - druga.czas;
  if (pierwsza.zakres === null) return 1;
  if (druga.zakres === null) return -1;
  if (pierwsza.zakres.poczatek !== druga.zakres.poczatek) {
    return pierwsza.zakres.poczatek - druga.zakres.poczatek;
  }
  return pierwsza.czas - druga.czas;
}

/** Funkcja przybornikPrzefiltruj zawęża wykaz znakowań dokumentu wedle rodzaju, autora i stanu otwarcia pozycji. */
export function przybornikPrzefiltruj(
  pozycje: readonly PozycjaZnakowania[],
  filtr: FiltrZnakowan,
): PozycjaZnakowania[] {
  return pozycje.filter((pozycja) => {
    if (filtr.rodzaj !== 'wszystkie' && pozycja.rodzaj !== filtr.rodzaj) return false;
    if (filtr.autor !== 'wszyscy' && pozycja.autor !== filtr.autor) return false;
    if (filtr.stan === 'otwarte' && !pozycja.otwarta) return false;
    if (filtr.stan === 'zamkniete' && pozycja.otwarta) return false;
    return true;
  });
}

/**
 * Zdanie o wykazie — liczby policzone, nie zaokrąglone.
 *
 * Podaje osobno liczbę pozycji wszystkich, widocznych po zawężeniu i tych
 * autora `model`, bo to trzy różne pytania Operatora: ile jest, ile widzę
 * i ile z tego zrobił model.
 */
export function przybornikOpiszWykaz(
  wszystkie: readonly PozycjaZnakowania[],
  widoczne: readonly PozycjaZnakowania[],
): string {
  if (wszystkie.length === 0) {
    return 'Dokument nie ma ani jednego znakowania: ani komentarza, ani propozycji, ani zmiany ' +
      'śledzonej, ani adnotacji, ani znacznika.';
  }
  const odModelu = wszystkie.filter((pozycja) => pozycja.autor === StudioAuthor.Model).length;
  const otwarte = wszystkie.filter((pozycja) => pozycja.otwarta).length;
  return (
    `Znakowań ${wszystkie.length}, z tego otwartych ${otwarte} i autora „model" ${odModelu}. ` +
    `Po zawężeniu widocznych ${widoczne.length}.`
  );
}

/** Funkcja przybornikPolicz zwraca liczbę pozycji danego rodzaju znakowania, do plakietek przy zawężeniu wykazu. */
export function przybornikPolicz(
  pozycje: readonly PozycjaZnakowania[],
  rodzaj: RodzajZnakowania,
): number {
  return pozycje.filter((pozycja) => pozycja.rodzaj === rodzaj).length;
}
