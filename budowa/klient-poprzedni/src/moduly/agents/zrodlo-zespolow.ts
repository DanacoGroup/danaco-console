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
 * Cztery komendy obszaru `team.*` wraz ze zdarzeniem `team.changed` widziane
 * przez panel zespołów ekspertów. Nazwy `Command.*` kończą się w tym pliku,
 * a żadna czynność nie rzuca: odmowa wraca polem `blad` wyniku.
 */
export interface ZrodloZespolow {
  /** Zespoły zapisane przez Operatora; pusta fraza znaczy wykaz pełny. */
  wykaz(fraza: string): Promise<Wynik<{ teams: Team[]; total: number }>>;
  /** Zapisuje zespół; pusty identyfikator zakłada nowy, podany zmienia stary. */
  zapisz(zapis: ZapisZespolu): Promise<Wynik<{ team: Team }>>;
  /** Czyta jeden zespół wraz ze składem — wykaz nie musi go nieść w całości. */
  wczytaj(idZespolu: string): Promise<Wynik<{ team: Team }>>;
  /** Zakłada kopię zespołu; pusta nazwa zostawia nazwanie rdzeniowi. */
  powiel(idZespolu: string, nazwaKopii: string): Promise<Wynik<{ team: Team }>>;
  /** Subskrypcja zdarzenia `team.changed` — jedynego zdarzenia obszaru. */
  naZmiane(sluchacz: (tresc: TeamChangedEvent) => void): Odsubskrybuj;
}

/**
 * Treść zapisu zespołu w postaci, w jakiej trzyma ją formularz panelu: zespół
 * zmieniany, nazwa, opis oraz skład w kolejności nadanej przez Operatora.
 * Źródło przekłada tę postać na żądanie kontraktu przy każdym zapisie.
 */
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
      // Opis jest w kontrakcie opcjonalny, więc pusty nie idzie do żądania.
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
 * Zdanie o zespole na wykazie podaje nazwę, liczebność składu oraz opis, gdy
 * zespół go ma. Stoi w źródle, a nie w widoku, ponieważ mówi o kształcie bytu
 * `Team`, a nie o układzie kontrolek panelu.
 */
export function zdanieOZespole(zespol: Team): string {
  const liczba = zespol.agentIds.length;
  const sklad =
    liczba === 0 ? 'skład pusty' : liczba === 1 ? '1 ekspert w składzie' : `${liczba} ekspertów w składzie`;
  const opis = (zespol.description ?? '').trim();
  return opis === '' ? sklad : `${sklad} — ${opis}`;
}
