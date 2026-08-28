import { Command, ComponentKind, type Component } from '../../../shared/contract';
import { opisOdmowyBledu } from '../komponenty/odmowa';
import type { Kanal } from '../protokol/kanal';
import { czyObiekt, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';
import type { KodKomponentu } from './pozycje-komponentow';

/** Zakładanie komponentu własnego ze strony głównej obejmuje trzy z czterech rodzajów; profil asystenta zostaje wyłączony, dopóki platforma nie ma dla niego magazynu. */
export interface ZalozenieKomponentu {
  /** `component.create` — zakłada komponent i oddaje go wraz z metadanymi. */
  zaloz(rodzaj: KodKomponentu, nazwa: string, opis: string): Promise<WynikZalozenia>;
}

/** Wynik czynności założenia komponentu w postaci, którą widok pokazuje bez dopowiadania niczego od siebie. */
export interface WynikZalozenia {
  udane: boolean;
  /** Zdanie dla Operatora — powodzenie albo odmowa rdzenia. */
  zdanie: string;
  /** Komponent założony; wyłącznie przy powodzeniu. */
  komponent?: Component;
}

/** Rodzaje, które rdzeń zakłada; wykaz stoi tutaj, a nie w widoku, bo mówi o zdolnościach rdzenia, nie o wyglądzie. */
export const RODZAJE_DO_ZALOZENIA: readonly { kod: KodKomponentu; nazwa: string }[] = [
  { kod: ComponentKind.Automations, nazwa: 'Automatyka' },
  { kod: ComponentKind.Agents, nazwa: 'Ekspert' },
  { kod: ComponentKind.Workspace, nazwa: 'Projekt' },
];

/** Rodzaj, którego rdzeń dzisiaj nie zakłada, wraz z powodem przeznaczonym do pokazania w zdaniu widoku. */
export const RODZAJ_BEZ_MAGAZYNU = {
  kod: ComponentKind.Assistant,
  powod:
    'Profilu asystenta nie da się dziś założyć — rdzeń odmawia, bo platforma nie ' +
    'ma magazynu profili. Pozostałe trzy rodzaje zakładają się normalnie.',
} as const;

export function utworzZalozenieKomponentu(kanal: Kanal): ZalozenieKomponentu {
  return {
    async zaloz(rodzaj, nazwa, opis) {
      const wpisana = nazwa.trim();
      if (wpisana === '') {
        // Jedyny sprawdzian przed wysyłką: puste imię nie jest odmową, lecz formularzem niedokończonym.
        return { udane: false, zdanie: 'Wpisz nazwę — komponent bez nazwy nie powstanie.' };
      }

      const wynik = sprawdzKsztalt(
        await wywolaj(kanal, Command.ComponentCreate, {
          kind: rodzaj,
          name: wpisana,
          ...(opis.trim() === '' ? {} : { description: opis.trim() }),
        }),
        Command.ComponentCreate,
        (tresc) => czyObiekt(tresc.component),
      );

      if (!wynik.udany || wynik.wynik === undefined) {
        return { udane: false, zdanie: opisOdmowyBledu('Założenie komponentu', wynik.blad) };
      }

      const komponent = wynik.wynik.component;
      return {
        udane: true,
        // Potwierdzenie mówi o tym, co oddał rdzeń: nazwa bywa inna niż wpisana, gdy rdzeń ją przytnie.
        zdanie:
          komponent.targetId === undefined
            ? `Komponent „${komponent.name}" założony.`
            : `Komponent „${komponent.name}" założony (${komponent.targetId}).`,
        komponent,
      };
    },
  };
}
