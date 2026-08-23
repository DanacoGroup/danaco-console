import {
  Command,
  EventType,
  type Team,
  type TeamChangedEvent,
  type TeamDuplicateRequest,
  type TeamListRequest,
  type TeamSaveRequest,
} from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { Kanal, Wynik } from '../../protokol/kanal';
import { czyObiekt, czyTablica, sprawdzKsztalt } from '../../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../../protokol/wywolanie';

/**
 * Cztery komendy obszaru `team.*` widziane przez panel „Zespoły ekspertów”.
 *
 * Źródło jest osobne od `zrodlo-agentow.ts`, bo zespół nie jest polem eksperta
 * ani jego wersją — jest własnym bytem kontraktu (`Team`) z własnym zdarzeniem
 * `team.changed` i własnym cyklem życia. Wsadzenie tych czterech komend do
 * źródła obszaru `agent.*` dałoby jeden plik o dwóch odpowiedzialnościach.
 *
 * Nazwy `Command.*` kończą się tutaj. Panel i widok składu wołają czynności
 * tego interfejsu — po polsku, o bytach dziedziny — i nigdy nie dotykają
 * kontraktu. Dzięki temu przemianowanie komendy w `contract.json` przerywa
 * kompilację w jednym pliku, nie w trzech widokach.
 *
 * Żadna czynność nie rzuca. Odmowa wraca polem `blad` wyniku i pokazuje się
 * w wierszu odpowiedzi panelu, w którym ją wywołano.
 */
export interface ZrodloZespolow {
  /** Zespoły zapisane przez Operatora; pusta fraza znaczy wykaz pełny. */
  wykaz(fraza: string): Promise<Wynik<{ teams: Team[]; total: number }>>;
  /**
   * Zapisuje zespół; pusty identyfikator zakłada nowy, podany zmienia istniejący.
   *
   * Rozróżnienie „nowy czy zmiana” robi wyłącznie obecność identyfikatora, więc
   * pole idzie do żądania tylko wtedy, gdy naprawdę jest. Wysłanie pustego
   * napisu byłoby dla rdzenia wskazaniem zespołu o nazwie pustej, nie brakiem
   * wskazania.
   */
  zapisz(zapis: ZapisZespolu): Promise<Wynik<{ team: Team }>>;
  /** Czyta jeden zespół wraz ze składem — wykaz nie musi go nieść w całości. */
  wczytaj(idZespolu: string): Promise<Wynik<{ team: Team }>>;
  /** Zakłada kopię zespołu; pusta nazwa zostawia nazwanie rdzeniowi. */
  powiel(idZespolu: string, nazwaKopii: string): Promise<Wynik<{ team: Team }>>;
  /** Subskrypcja zdarzenia `team.changed` — jedynego zdarzenia obszaru. */
  naZmiane(sluchacz: (tresc: TeamChangedEvent) => void): Odsubskrybuj;
}

/** Treść zapisu zespołu w postaci, w jakiej trzyma ją formularz panelu. */
export interface ZapisZespolu {
  /** Zespół zmieniany; pusty zakłada nowy. */
  idZespolu: string;
  nazwa: string;
  opis: string;
  /** Skład w kolejności nadanej przez Operatora. */
  idEkspertow: readonly string[];
}

export function utworzZrodloZespolow(kanal: Kanal): ZrodloZespolow {
  return {
    async wykaz(fraza) {
      const zadanie: TeamListRequest = {};
      if (fraza.trim() !== '') zadanie.query = fraza.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamList, zadanie),
        Command.TeamList,
        (tresc) => czyTablica(tresc.teams),
      );
    },

    async zapisz(zapis) {
      const zadanie: TeamSaveRequest = {
        name: zapis.nazwa.trim(),
        agentIds: [...zapis.idEkspertow],
      };
      if (zapis.idZespolu !== '') zadanie.teamId = zapis.idZespolu;
      // Opis jest w kontrakcie polem opcjonalnym, więc pusty nie jedzie: przy
      // zmianie zespołu wysłanie pustego napisu skasowałoby opis już zapisany,
      // choć formularz o skasowanie nie prosił.
      if (zapis.opis.trim() !== '') zadanie.description = zapis.opis.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamSave, zadanie),
        Command.TeamSave,
        (tresc) => czyObiekt(tresc.team),
      );
    },

    async wczytaj(idZespolu) {
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamLoad, { teamId: idZespolu }),
        Command.TeamLoad,
        (tresc) => czyObiekt(tresc.team),
      );
    },

    async powiel(idZespolu, nazwaKopii) {
      const zadanie: TeamDuplicateRequest = { teamId: idZespolu };
      if (nazwaKopii.trim() !== '') zadanie.name = nazwaKopii.trim();
      return sprawdzKsztalt(
        await wywolaj(kanal, Command.TeamDuplicate, zadanie),
        Command.TeamDuplicate,
        (tresc) => czyObiekt(tresc.team),
      );
    },

    naZmiane(sluchacz) {
      return kanal.naZdarzenie(EventType.TeamChanged, (tresc) => sluchacz(tresc));
    },
  };
}

/**
 * Zdanie o zespole na wykazie — nazwa, liczebność składu i opis, gdy jest.
 *
 * Stoi tutaj, a nie w widoku, bo mówi o kształcie bytu `Team`, nie o układzie
 * kontrolek: liczebność składu jest jedyną rzeczą, po której Operator poznaje
 * na wykazie różnicę między zespołem a jego kopią świeżo powieloną.
 */
export function zdanieOZespole(zespol: Team): string {
  const liczba = zespol.agentIds.length;
  const sklad =
    liczba === 0 ? 'skład pusty' : liczba === 1 ? '1 ekspert w składzie' : `${liczba} ekspertów w składzie`;
  const opis = (zespol.description ?? '').trim();
  return opis === '' ? sklad : `${sklad} — ${opis}`;
}
