import { ETYKIETA_BRAKU_ZRODLA } from './etykiety-pulpitu';
import { utworzKafelLiczby } from './kafel-liczby';
import type { AktywnoscAI } from './model-danych';
import { utworzSekcje } from './naglowek-sekcji';

/**
 * Sekcja aktywności AI — cztery kafle otwierające pulpit, odpowiadające na
 * pytanie, co biegnie w tle.
 *
 * Dwie miary pochodzą z telemetrii rdzenia; dwie nie mają źródła w kontrakcie
 * i kafel mówi to wprost zamiast pokazać wartość zastępczą. Nazwy miar są
 * dosłowne, żeby nie mylić wysycenia kanałów z obciążeniem maszyny.
 *
 * Pas decyzji stoi pod kaflami i dokłada go `mission-control.ts`.
 */
export interface SekcjaAktywnosci {
  element: HTMLElement;
  /** Miejsce montażu pasa decyzji, pod rzędem kafli. */
  podKaflami: HTMLElement;
  /** Podmienia liczby bez przebudowy sekcji. */
  odswiez(aktywnosc: AktywnoscAI): void;
}

/** Buduje sekcję aktywności wraz z miejscem na pas decyzji. */
export function utworzSekcjeAktywnosci(aktywnosc: AktywnoscAI): SekcjaAktywnosci {
  const { element, cialo } = utworzSekcje(
    'mc-sekcja--aktywnosc',
    {
      tytul: 'Aktywność AI',
      dopisek: 'Co biegnie w tej chwili — także wtedy, gdy patrzysz gdzie indziej.',
      ikona: 'uruchom',
    },
    'mc-tytul-aktywnosc',
  );

  const kafle = document.createElement('div');
  kafle.className = 'mc-kafle';

  const podKaflami = document.createElement('div');
  podKaflami.className = 'mc-sekcja--aktywnosc__pod';

  cialo.append(kafle, podKaflami);

  const odswiez = (dane: AktywnoscAI): void => {
    kafle.replaceChildren(...wyrysujKafle(dane));
  };
  odswiez(aktywnosc);

  return { element, podKaflami, odswiez };
}

/**
 * Cztery kafle sekcji, w kolejności wymaganej przez makietę. Miara bez źródła
 * w kontrakcie dostaje kreskę i etykietę braku zamiast liczby.
 */
function wyrysujKafle(aktywnosc: AktywnoscAI): HTMLElement[] {
  return [
    utworzKafelLiczby({
      wartosc: String(aktywnosc.aktywneProcesy),
      etykieta: 'aktywnych procesów',
      objasnienie: 'Procesy telemetrii progress.changed w stanie „w biegu"',
      ikona: 'uruchom',
    }),
    utworzKafelLiczby({
      wartosc: String(aktywnosc.agenciPracuja),
      etykieta: 'agentów pracuje',
      objasnienie: 'Okna komunikacji, w których biegnie proces telemetrii',
      ikona: 'uzytkownik',
    }),
    utworzKafelLiczby({
      wartosc: aktywnosc.modeleAnalizuja === null ? '—' : String(aktywnosc.modeleAnalizuja),
      etykieta: 'modele analizują',
      objasnienie:
        aktywnosc.modeleAnalizuja === null
          ? ETYKIETA_BRAKU_ZRODLA
          : 'Kanały modelu przetwarzające zlecenie',
      ikona: 'oko',
    }),
    utworzKafelLiczby({
      wartosc: aktywnosc.walidatoryOczekuja === null ? '—' : String(aktywnosc.walidatoryOczekuja),
      etykieta: 'walidatory oczekują',
      objasnienie:
        aktywnosc.walidatoryOczekuja === null
          ? ETYKIETA_BRAKU_ZRODLA
          : 'Wyniki czekające na werdykt walidatora',
      ikona: 'ptaszek-kolo',
    }),
  ];
}
